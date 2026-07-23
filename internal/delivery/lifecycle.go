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
	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
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

type ExecutionAuthorityRecord struct {
	EnvelopeID        string
	EnvelopeDigest    string
	DelegationID      string
	WorkerSubjectID   string
	PolicyDecisionIDs []string
	ExpiresAt         time.Time
}

type CandidateRecord struct {
	Digest              string
	TreeDigest          string
	BundleDigest        string
	ChangedPaths        []string
	WorkerRequestID     string
	WorkerRequestDigest string
	WorkerResultDigest  string
	ArtefactReferences  []string
}

type VerificationRecord struct {
	AssurancePlanID     string
	AssurancePlanDigest string
	ResultID            string
	ResultDigest        string
	Outcome             string
	EvidenceReferences  []string
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
	Execution       *ExecutionAuthorityRecord
	Candidate       *CandidateRecord
	Verification    *VerificationRecord
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

type RecordExecutionAuthorityRequest struct {
	TransitionRequest
	Envelope     trust.ExecutionEnvelope
	DelegationID string
	WorkerID     string
}

type RecordCandidateRequest struct {
	TransitionRequest
	WorkerRequest      workercontract.Request
	Bundle             workercontract.CandidateBundle
	Result             workercontract.Result
	ArtefactReferences []string
}

type RecordVerificationRequest struct {
	TransitionRequest
	AssurancePlan      verificationcontract.AssurancePlan
	Result             verificationcontract.Result
	EvidenceReferences []string
}

func (service *Service) BeginContextAuthorization(ctx context.Context, request TransitionRequest) (DeliveryRun, bool, error) {
	return service.transition(ctx, request, CapturedState, ContextAuthorizingState, nil, nil, nil, nil, nil, nil)
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
	return service.transition(ctx, request.TransitionRequest, ContextAuthorizingState, ContextReadyState, record, nil, nil, nil, nil, nil)
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
	return service.transition(ctx, request.TransitionRequest, ContextReadyState, PlannedState, nil, &plan, nil, nil, nil, nil)
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
	return service.transition(ctx, request.TransitionRequest, PlannedState, AwaitingPlanApprovalState, nil, nil, approval, nil, nil, nil)
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
	return service.transition(ctx, request.TransitionRequest, AwaitingPlanApprovalState, PlanApprovedState, nil, nil, approved, nil, nil, nil)
}

func (service *Service) BeginExecutionAuthorization(ctx context.Context, request TransitionRequest) (DeliveryRun, bool, error) {
	return service.transition(ctx, request, PlanApprovedState, ExecutionAuthorizingState, nil, nil, nil, nil, nil, nil)
}

func (service *Service) RecordExecutionAuthority(ctx context.Context, request RecordExecutionAuthorityRequest) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request.TransitionRequest)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	plan, approval := request.Run.Plan(), request.Run.Approval()
	if plan == nil || approval == nil || request.DelegationID == "" || request.WorkerID == "" ||
		request.Run.State() != ExecutionAuthorizingState {
		return DeliveryRun{}, false, failure.New(failure.InvalidRequest, "execution authority input is invalid", request.CorrelationID, failure.RetryNone)
	}
	if err := request.Envelope.VerifyIntegrity(request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	if err := request.Envelope.AuthorizeAction(workercontract.PrepareChangeAction, at, request.CorrelationID); err != nil {
		return DeliveryRun{}, false, err
	}
	if request.Envelope.Scope != (trust.ExecutionScope{
		TenantID: request.Run.TenantID(), WorkspaceID: request.Run.WorkspaceID(),
		ProjectID: request.Run.ProjectID(), DeliveryRunID: request.Run.ID(),
	}) || request.Envelope.Principal.SubjectID != request.Run.ActorID() ||
		request.Envelope.Principal.AgentID != request.WorkerID ||
		!slices.Equal(request.Envelope.Principal.DelegationChain, []string{request.DelegationID}) ||
		!slices.Equal(request.Envelope.ApprovalRecordIDs, []string{approval.RecordID}) ||
		request.Envelope.Purpose != M1ChangePurpose || request.Envelope.RiskClass != plan.RiskClass ||
		!slices.Equal(request.Envelope.Classifications, []string{request.Run.Classification()}) ||
		len(request.Envelope.ModelIDs) != 0 || len(request.Envelope.ToolIDs) != 0 ||
		len(request.Envelope.SecretReferenceIDs) != 0 || len(request.Envelope.Destinations) != 0 ||
		request.Envelope.Integrity == nil {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "execution authority is outside the approved plan", request.CorrelationID, failure.RetryNewAuthority)
	}
	record := &ExecutionAuthorityRecord{
		EnvelopeID: request.Envelope.ID, EnvelopeDigest: request.Envelope.Integrity.Digest,
		DelegationID: request.DelegationID, WorkerSubjectID: request.WorkerID,
		PolicyDecisionIDs: slices.Clone(request.Envelope.PolicyDecisionIDs), ExpiresAt: request.Envelope.ExpiresAt,
	}
	return service.transition(ctx, request.TransitionRequest, ExecutionAuthorizingState, ExecutingState, nil, nil, nil, record, nil, nil)
}

func (service *Service) RecordCandidate(ctx context.Context, request RecordCandidateRequest) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request.TransitionRequest)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	plan, contextRecord := request.Run.Plan(), request.Run.Context()
	approval, authority := request.Run.Approval(), request.Run.Execution()
	if request.Run.State() != ExecutingState || plan == nil || contextRecord == nil || approval == nil || authority == nil {
		return DeliveryRun{}, false, failure.New(failure.Conflict, "delivery run is not executing", request.CorrelationID, failure.RetryReconcile)
	}
	if err := request.WorkerRequest.Validate(at); err != nil {
		return DeliveryRun{}, false, failure.New(failure.VerificationFailed, "worker request could not be verified", request.CorrelationID, failure.RetryNone)
	}
	if err := request.Bundle.Validate(); err != nil {
		return DeliveryRun{}, false, failure.New(failure.VerificationFailed, "candidate bundle could not be verified", request.CorrelationID, failure.RetryNone)
	}
	if request.WorkerRequest.EnvelopeID != authority.EnvelopeID ||
		request.WorkerRequest.DelegationID != authority.DelegationID ||
		request.WorkerRequest.PlanID != plan.ID || request.WorkerRequest.PlanDigest != plan.Digest ||
		request.WorkerRequest.PlanApprovalRecordID != approval.RecordID ||
		request.WorkerRequest.ContextPackageID != contextRecord.PackageID ||
		request.WorkerRequest.ContextManifestDigest != contextRecord.ManifestDigest ||
		request.WorkerRequest.SourceRevision != contextRecord.SourceRevisions[0] ||
		request.WorkerRequest.Scope != (workercontract.Scope{
			TenantID: request.Run.TenantID(), WorkspaceID: request.Run.WorkspaceID(),
			ProjectID: request.Run.ProjectID(), RepositoryID: request.Run.RepositoryID(), DeliveryRunID: request.Run.ID(),
		}) || !slices.Equal(request.WorkerRequest.AllowedWritePaths, plan.ChangePaths) ||
		request.Result.SchemaVersion != workercontract.ResultVersion ||
		request.Result.RequestID != request.WorkerRequest.ID ||
		request.Result.RequestDigest != request.WorkerRequest.Digest ||
		request.Result.CandidateDigest != request.Bundle.Digest ||
		request.Result.BundleDigest != request.Bundle.Digest ||
		request.Result.TreeDigest != workercontract.CandidateTreeDigest(request.Bundle) ||
		!slices.Equal(request.Result.ChangedPaths, bundlePaths(request.Bundle)) ||
		len(request.ArtefactReferences) == 0 {
		return DeliveryRun{}, false, failure.New(failure.Forbidden, "candidate is outside worker authority", request.CorrelationID, failure.RetryNone)
	}
	record := &CandidateRecord{
		Digest: request.Bundle.Digest, TreeDigest: request.Result.TreeDigest, BundleDigest: request.Result.BundleDigest,
		ChangedPaths: slices.Clone(request.Result.ChangedPaths), WorkerRequestID: request.WorkerRequest.ID,
		WorkerRequestDigest: request.WorkerRequest.Digest, WorkerResultDigest: digestValue(request.Result),
		ArtefactReferences: canonicalStrings(request.ArtefactReferences),
	}
	return service.transition(ctx, request.TransitionRequest, ExecutingState, VerifyingState, nil, nil, nil, nil, record, nil)
}

func (service *Service) RecordVerification(ctx context.Context, request RecordVerificationRequest) (DeliveryRun, bool, error) {
	at, err := service.transitionTime(request.TransitionRequest)
	if err != nil {
		return DeliveryRun{}, false, err
	}
	plan, candidate := request.Run.Plan(), request.Run.Candidate()
	if request.Run.State() != VerifyingState || plan == nil || candidate == nil {
		return DeliveryRun{}, false, failure.New(failure.Conflict, "delivery run is not verifying", request.CorrelationID, failure.RetryReconcile)
	}
	if err := request.AssurancePlan.Validate(at); err != nil {
		return DeliveryRun{}, false, failure.New(failure.VerificationFailed, "assurance plan could not be verified", request.CorrelationID, failure.RetryNone)
	}
	if err := request.Result.Validate(request.AssurancePlan); err != nil ||
		request.AssurancePlan.CandidateDigest != candidate.Digest ||
		request.AssurancePlan.PlanDigest != plan.Digest ||
		!slices.Equal(request.AssurancePlan.ChangedPaths, candidate.ChangedPaths) ||
		request.Result.Outcome != "passed" || len(request.EvidenceReferences) == 0 {
		return DeliveryRun{}, false, failure.New(failure.VerificationFailed, "candidate verification did not pass independently", request.CorrelationID, failure.RetryNone)
	}
	record := &VerificationRecord{
		AssurancePlanID: request.AssurancePlan.ID, AssurancePlanDigest: request.AssurancePlan.Digest,
		ResultID: request.Result.ID, ResultDigest: request.Result.Digest, Outcome: request.Result.Outcome,
		EvidenceReferences: canonicalStrings(request.EvidenceReferences),
	}
	return service.transition(ctx, request.TransitionRequest, VerifyingState, AwaitingPublicationApprovalState, nil, nil, nil, nil, nil, record)
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

func (service *Service) transition(ctx context.Context, request TransitionRequest, from, to string, contextRecord *ContextRecord, plan *DeliveryPlan, approval *PlanApprovalRecord, execution *ExecutionAuthorityRecord, candidate *CandidateRecord, verification *VerificationRecord) (DeliveryRun, bool, error) {
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
		Execution: cloneExecutionRecord(execution), Candidate: cloneCandidateRecord(candidate), Verification: cloneVerificationRecord(verification),
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
	if command.Execution != nil {
		next.Execution = cloneExecutionRecord(command.Execution)
	}
	if command.Candidate != nil {
		next.Candidate = cloneCandidateRecord(command.Candidate)
	}
	if command.Verification != nil {
		next.Verification = cloneVerificationRecord(command.Verification)
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

func cloneExecutionRecord(value *ExecutionAuthorityRecord) *ExecutionAuthorityRecord {
	if value == nil {
		return nil
	}
	result := *value
	result.PolicyDecisionIDs = slices.Clone(value.PolicyDecisionIDs)
	return &result
}

func cloneCandidateRecord(value *CandidateRecord) *CandidateRecord {
	if value == nil {
		return nil
	}
	result := *value
	result.ChangedPaths = slices.Clone(value.ChangedPaths)
	result.ArtefactReferences = slices.Clone(value.ArtefactReferences)
	return &result
}

func cloneVerificationRecord(value *VerificationRecord) *VerificationRecord {
	if value == nil {
		return nil
	}
	result := *value
	result.EvidenceReferences = slices.Clone(value.EvidenceReferences)
	return &result
}

func bundlePaths(bundle workercontract.CandidateBundle) []string {
	result := make([]string, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		result = append(result, file.Path)
	}
	return result
}

func digestValue(value any) string {
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}
