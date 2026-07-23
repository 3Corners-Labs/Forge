package delivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"sort"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
)

var (
	ErrTransitionConflict  = errors.New("delivery transition conflict")
	ErrDeliveryUnavailable = errors.New("delivery lifecycle store unavailable")
)

type ContextRecord struct {
	PackageID         string
	ManifestDigest    string
	EnvelopeID        string
	PolicyDecisionIDs []string
	SourceRevisions   []string
	Omissions         []string
	Conflicts         []string
	Classification    string
	ExpiresAt         time.Time
}

type PlanApprovalRecord struct {
	RequestID     string
	RequestDigest string
	RecordID      string
}

type OutboxEvent struct {
	ID            string
	RunID         string
	FromState     string
	ToState       string
	RunVersion    uint64
	CorrelationID string
	OccurredAt    time.Time
}

type TransitionCommand struct {
	RunID           string
	ScopeKey        string
	ExpectedVersion uint64
	FromState       string
	ToState         string
	IdempotencyKey  string
	Fingerprint     string
	CorrelationID   string
	OccurredAt      time.Time
	Context         *ContextRecord
	Plan            *DeliveryPlan
	Approval        *PlanApprovalRecord
}

type transitionEntry struct {
	fingerprint string
	record      StoredRecord
}

type LifecycleStore interface {
	Load(context.Context, string) (StoredRecord, bool, error)
	Transition(context.Context, TransitionCommand) (StoredRecord, bool, error)
}

type TransitionRequest struct {
	CorrelationID   string
	IdempotencyKey  string
	Run             DeliveryRun
	ExpectedVersion uint64
}

type RecordContextRequest struct {
	TransitionRequest
	Package  trust.ContextPackage
	Envelope trust.ExecutionEnvelope
}

type RecordPlanRequest struct {
	TransitionRequest
	ProfileVersion string
}

type AwaitPlanApprovalRequest struct {
	TransitionRequest
	ApprovalRequestID string
}

type RecordPlanApprovalRequest struct {
	TransitionRequest
	ApprovalRecordID string
}

func (service *Service) BeginContextAuthorization(ctx context.Context, request TransitionRequest) (DeliveryRun, bool, error) {
	return service.transition(ctx, request, CapturedState, ContextAuthorizingState, nil, nil, nil)
}

func (service *Service) RecordContext(ctx context.Context, request RecordContextRequest) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request.TransitionRequest)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	if err := request.Package.VerifyManifest(request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	if err := request.Envelope.VerifyIntegrity(request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	if err := request.Package.ValidateFor(request.Envelope, at, request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	if request.Run.State() != ContextAuthorizingState || request.Package.Objective != request.Run.Objective() ||
		request.Envelope.Scope != (trust.ExecutionScope{
			TenantID: request.Run.TenantID(), WorkspaceID: request.Run.WorkspaceID(),
			ProjectID: request.Run.ProjectID(), DeliveryRunID: request.Run.ID(),
		}) || request.Package.ExecutionEnvelopeID != request.Envelope.ID ||
		!slices.Equal(request.Package.PolicyDecisionIDs, request.Envelope.PolicyDecisionIDs) {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "context package is outside the delivery run", request.CorrelationID, failure.RetryNone)
	}
	revisions := make([]string, 0, len(request.Package.Items))
	classification := request.Run.Classification()
	for _, item := range request.Package.Items {
		if item.Classification != classification {
			return DeliveryRun{}, false, failure.New(failure.Forbidden, "context classification differs from the delivery run", request.CorrelationID, failure.RetryNone)
		}
		revisions = append(revisions, item.SourceID+"@"+item.SourceRevision)
	}
	sort.Strings(revisions)
	if len(revisions) == 0 {
		return DeliveryRun{}, false, failure.New(failure.InvalidRequest, "context package contains no reproducible source revision", request.CorrelationID, failure.RetryNone)
	}
	expiresAt := request.Envelope.ExpiresAt
	if request.Package.ExpiresAt != nil && request.Package.ExpiresAt.Before(expiresAt) {
		expiresAt = *request.Package.ExpiresAt
	}
	record := &ContextRecord{
		PackageID: request.Package.ID, ManifestDigest: request.Package.ManifestDigest,
		EnvelopeID: request.Envelope.ID, PolicyDecisionIDs: slices.Clone(request.Package.PolicyDecisionIDs),
		SourceRevisions: revisions, Omissions: canonicalStrings(request.Package.Omissions),
		Conflicts: canonicalStrings(request.Package.Conflicts), Classification: classification, ExpiresAt: expiresAt,
	}
	return service.transition(ctx, request.TransitionRequest, ContextAuthorizingState, ContextReadyState, record, nil, nil)
}

func (service *Service) RecordPlan(ctx context.Context, profiles PlanProfileResolver, request RecordPlanRequest) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request.TransitionRequest)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	if profiles == nil || request.ProfileVersion == "" {
		return DeliveryRun{}, false, failure.New(failure.InvalidRequest, "trusted delivery plan profile is required", request.CorrelationID, failure.RetryNone)
	}
	profile, found, err := profiles.ResolvePlanProfile(ctx, request.ProfileVersion)
	if err != nil {
		return DeliveryRun{}, false, failure.New(failure.DependencyUnavailable, "delivery plan profile resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || profile.Version != request.ProfileVersion {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "delivery plan profile is not registered", request.CorrelationID, failure.RetryNone)
	}
	profile, err = NewPlanProfile(profile)
	if err != nil {
		return DeliveryRun{}, false, failure.New(failure.Indeterminate, "registered delivery plan profile is invalid", request.CorrelationID, failure.RetryNone)
	}
	plan, err := buildPlan(request.Run, profile, at, request.CorrelationID)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	return service.transition(ctx, request.TransitionRequest, ContextReadyState, PlannedState, nil, &plan, nil)
}

func (service *Service) AwaitPlanApproval(ctx context.Context, resolver ApprovalRequestResolver, request AwaitPlanApprovalRequest) (DeliveryRun, bool, error) {
	if resolver == nil || request.ApprovalRequestID == "" {
		return DeliveryRun{}, false, failure.New(failure.InvalidRequest, "approval request resolver and identifier are required", request.CorrelationID, failure.RetryNone)
	}
	approvalRequest, found, err := resolver.ResolveApprovalRequest(ctx, request.Run.TenantID(), request.ApprovalRequestID)
	if err != nil {
		return DeliveryRun{}, false, failure.New(failure.DependencyUnavailable, "approval request resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || approvalRequest.ID != request.ApprovalRequestID {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "current approval request was not found", request.CorrelationID, failure.RetryNone)
	}
	if err := approvalRequest.VerifyDigest(request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	plan := request.Run.Plan()
	if request.Run.State() != PlannedState || plan == nil ||
		approvalRequest.Scope != approvalScope(request.Run) ||
		approvalRequest.RequestingActorID != request.Run.ActorID() ||
		approvalRequest.IntendedApproverID != request.Run.ActorID() ||
		approvalRequest.SubjectID != plan.ID ||
		approvalRequest.SubjectDigest != plan.Digest ||
		approvalRequest.EvidenceCorrelationID != request.CorrelationID {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "approval request is outside the current delivery plan", request.CorrelationID, failure.RetryNone)
	}
	approval := &PlanApprovalRecord{RequestID: approvalRequest.ID, RequestDigest: approvalRequest.RequestDigest}
	return service.transition(ctx, request.TransitionRequest, PlannedState, AwaitingPlanApprovalState, nil, nil, approval)
}

func (service *Service) RecordPlanApproval(ctx context.Context, resolver ApprovalResolver, request RecordPlanApprovalRequest) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request.TransitionRequest)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	if resolver == nil || request.ApprovalRecordID == "" {
		return DeliveryRun{}, false, failure.New(failure.InvalidRequest, "approval resolver and record are required", request.CorrelationID, failure.RetryNone)
	}
	plan := request.Run.Plan()
	pending := request.Run.Approval()
	if request.Run.State() != AwaitingPlanApprovalState || plan == nil || pending == nil {
		return DeliveryRun{}, false, failure.New(failure.Conflict, "delivery run is not awaiting plan approval", request.CorrelationID, failure.RetryReconcile)
	}
	record, found, err := resolver.ResolveApproval(ctx, request.Run.TenantID(), request.ApprovalRecordID)
	if err != nil {
		return DeliveryRun{}, false, failure.New(failure.DependencyUnavailable, "approval resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || record.ID != request.ApprovalRecordID {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "current approval record was not found", request.CorrelationID, failure.RetryNone)
	}
	expected := trust.ApprovalExpectedFacts{
		RequestID: pending.RequestID, RequestDigest: pending.RequestDigest, ApproverID: request.Run.ActorID(),
		Scope: approvalScope(request.Run), Environment: "developer-trial", Purpose: request.Run.Purpose(),
		Classification: request.Run.Classification(), RiskClass: plan.RiskClass,
		Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan",
		SubjectID: plan.ID, SubjectDigest: plan.Digest,
		EvidenceReferences: planEvidence(plan), PermittedConditions: []string{"require-independent-verification"},
	}
	if err := record.Authorize(expected, at, request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	approved := &PlanApprovalRecord{RequestID: pending.RequestID, RequestDigest: pending.RequestDigest, RecordID: record.ID}
	return service.transition(ctx, request.TransitionRequest, AwaitingPlanApprovalState, PlanApprovedState, nil, nil, approved)
}

type PlanSubject = trust.PlanApprovalSubject

func (service *Service) ResolvePlanSubject(ctx context.Context, tenantID, runID string) (PlanSubject, bool, error) {
	store, ok := service.store.(LifecycleStore)
	if !ok || ctx == nil || tenantID == "" || runID == "" {
		return PlanSubject{}, false, ErrDeliveryUnavailable
	}
	stored, found, err := store.Load(ctx, runID)
	if err != nil || !found {
		return PlanSubject{}, found, err
	}
	if stored.TenantID != tenantID ||
		(stored.State != PlannedState && stored.State != AwaitingPlanApprovalState) ||
		stored.Plan == nil || stored.Context == nil {
		return PlanSubject{}, false, nil
	}
	plan := clonePlan(stored.Plan)
	return PlanSubject{
		RunID: stored.ID, PlanID: plan.ID, PlanDigest: plan.Digest,
		Scope:             trust.ApprovalScope{TenantID: stored.TenantID, WorkspaceID: stored.WorkspaceID, ProjectID: stored.ProjectID, RepositoryID: stored.RepositoryID, DeliveryRunID: stored.ID},
		RequestingActorID: stored.ActorID, Purpose: stored.Purpose, Classification: stored.Classification,
		RiskClass: plan.RiskClass, ContextPackageID: plan.ContextPackageID, ContextManifest: plan.ContextManifest,
		EvidenceReferences: planEvidence(plan), ExpiresAt: plan.ExpiresAt, CorrelationID: stored.CorrelationID,
		PolicyDecisionIDs: slices.Clone(stored.Context.PolicyDecisionIDs),
	}, true, nil
}

func (service *Service) transition(ctx context.Context, request TransitionRequest, from, to string, contextRecord *ContextRecord, plan *DeliveryPlan, approval *PlanApprovalRecord) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	store, ok := service.store.(LifecycleStore)
	if !ok {
		return DeliveryRun{}, false, failure.New(failure.DependencyUnavailable, "delivery lifecycle persistence is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if request.Run.State() != from || request.Run.Version() != request.ExpectedVersion {
		return DeliveryRun{}, false, failure.New(failure.Conflict, "delivery run version or state changed", request.CorrelationID, failure.RetryReconcile)
	}
	command := TransitionCommand{
		RunID: request.Run.ID(), ScopeKey: storedScopeKey(request.Run), ExpectedVersion: request.ExpectedVersion,
		FromState: from, ToState: to, IdempotencyKey: request.IdempotencyKey,
		CorrelationID: request.CorrelationID, OccurredAt: at,
		Context: cloneContextRecord(contextRecord), Plan: clonePlan(plan), Approval: clonePlanApprovalRecord(approval),
	}
	command.Fingerprint = transitionFingerprint(command)
	stored, replayed, err := store.Transition(ctx, command)
	if err != nil {
		if errors.Is(err, ErrTransitionConflict) {
			return DeliveryRun{}, false, failure.New(failure.Conflict, "delivery transition conflicts with current state", request.CorrelationID, failure.RetryReconcile)
		}
		return DeliveryRun{}, false, failure.New(failure.DependencyUnavailable, "delivery transition is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if stored.ID != request.Run.ID() || stored.State != to || stored.Version != request.ExpectedVersion+1 {
		return DeliveryRun{}, false, failure.New(failure.Indeterminate, "delivery transition could not be verified", request.CorrelationID, failure.RetryNone)
	}
	return runFromStored(stored), replayed, nil
}

func (service *Service) transitionTime(request TransitionRequest) (time.Time, error) {
	if service == nil || service.store == nil || service.now == nil || request.CorrelationID == "" ||
		request.IdempotencyKey == "" || request.ExpectedVersion == 0 || request.Run.ID() == "" {
		return time.Time{}, failure.New(failure.InvalidRequest, "delivery transition request is invalid", request.CorrelationID, failure.RetryNone)
	}
	at := service.now()
	if at.IsZero() {
		return time.Time{}, failure.New(failure.DependencyUnavailable, "delivery transition clock is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	return at, nil
}

func (store *MemoryStore) Load(ctx context.Context, runID string) (StoredRecord, bool, error) {
	if store == nil || store.records == nil || store.byID == nil || ctx == nil || runID == "" {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	if err := ctx.Err(); err != nil {
		return StoredRecord{}, false, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	key, found := store.byID[runID]
	if !found {
		return StoredRecord{}, false, nil
	}
	return cloneStoredRecord(store.records[key]), true, nil
}

func (store *MemoryStore) Transition(ctx context.Context, command TransitionCommand) (StoredRecord, bool, error) {
	if store == nil || store.records == nil || store.byID == nil || store.changes == nil || ctx == nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	if err := ctx.Err(); err != nil {
		return StoredRecord{}, false, err
	}
	changeKey := command.ScopeKey + "\x00" + command.IdempotencyKey
	store.mu.Lock()
	defer store.mu.Unlock()
	if existing, found := store.changes[changeKey]; found {
		if existing.fingerprint != command.Fingerprint {
			return StoredRecord{}, false, ErrTransitionConflict
		}
		return cloneStoredRecord(existing.record), true, nil
	}
	recordKey, found := store.byID[command.RunID]
	if !found {
		return StoredRecord{}, false, ErrTransitionConflict
	}
	current := store.records[recordKey]
	if current.ScopeKey != command.ScopeKey || current.Version != command.ExpectedVersion || current.State != command.FromState {
		return StoredRecord{}, false, ErrTransitionConflict
	}
	next := cloneStoredRecord(current)
	next.State = command.ToState
	next.Version++
	if command.Context != nil {
		next.Context = cloneContextRecord(command.Context)
	}
	if command.Plan != nil {
		next.Plan = clonePlan(command.Plan)
	}
	if command.Approval != nil {
		next.Approval = clonePlanApprovalRecord(command.Approval)
	}
	next.Outbox = append(next.Outbox, OutboxEvent{
		ID: command.Fingerprint, RunID: command.RunID, FromState: command.FromState, ToState: command.ToState,
		RunVersion: next.Version, CorrelationID: command.CorrelationID, OccurredAt: command.OccurredAt,
	})
	store.records[recordKey] = cloneStoredRecord(next)
	store.changes[changeKey] = transitionEntry{fingerprint: command.Fingerprint, record: cloneStoredRecord(next)}
	return cloneStoredRecord(next), false, nil
}

func transitionFingerprint(command TransitionCommand) string {
	value := command
	value.Fingerprint = ""
	value.CorrelationID = ""
	value.OccurredAt = time.Time{}
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func storedScopeKey(run DeliveryRun) string {
	return run.TenantID() + "\x00" + run.WorkspaceID() + "\x00" + run.ProjectID() + "\x00" + run.RepositoryID() + "\x00" + run.ActorID()
}

func approvalScope(run DeliveryRun) trust.ApprovalScope {
	return trust.ApprovalScope{
		TenantID: run.TenantID(), WorkspaceID: run.WorkspaceID(), ProjectID: run.ProjectID(),
		RepositoryID: run.RepositoryID(), DeliveryRunID: run.ID(),
	}
}

func planEvidence(plan *DeliveryPlan) []string {
	return []string{
		"context:" + plan.ContextPackageID,
		"context-manifest:" + plan.ContextManifest,
		"intent:version:1",
		"plan:" + plan.Digest,
	}
}

func cloneContextRecord(value *ContextRecord) *ContextRecord {
	if value == nil {
		return nil
	}
	result := *value
	result.PolicyDecisionIDs = slices.Clone(value.PolicyDecisionIDs)
	result.SourceRevisions = slices.Clone(value.SourceRevisions)
	result.Omissions = slices.Clone(value.Omissions)
	result.Conflicts = slices.Clone(value.Conflicts)
	return &result
}

func clonePlanApprovalRecord(value *PlanApprovalRecord) *PlanApprovalRecord {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
