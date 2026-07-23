// Package approval owns immutable approval requests, human decisions,
// revocation chains, and current-state resolution. It never constructs plans,
// evaluates policy, or grants execution authority.
package approval

import (
	"context"
	"errors"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/identity"
)

const (
	RequestPlanReviewContractVersion = "forge.request-plan-review.v1"
	DecideApprovalContractVersion    = "forge.decide-approval.v1"
	RevokeApprovalContractVersion    = "forge.revoke-approval.v1"
)

var (
	ErrStoreUnavailable    = errors.New("approval store unavailable")
	ErrIdempotencyConflict = errors.New("approval idempotency conflict")
	ErrDecisionConflict    = errors.New("approval decision conflict")
)

type PlanSubject = trust.PlanApprovalSubject

type SubjectResolver interface {
	ResolvePlanSubject(context.Context, string, string) (PlanSubject, bool, error)
}

type PolicyResolver interface {
	ResolvePolicyDecision(context.Context, string) (trust.PolicyDecision, bool, error)
}

type Store interface {
	CreateRequest(context.Context, string, string, trust.ApprovalRequest) (trust.ApprovalRequest, bool, error)
	ResolveRequest(context.Context, string, string) (trust.ApprovalRequest, bool, error)
	CreateDecision(context.Context, string, string, string, trust.ApprovalRecordV2) (trust.ApprovalRecordV2, bool, error)
	ResolveCurrent(context.Context, string, string) (trust.ApprovalRecordV2, bool, error)
}

type Profile struct {
	Issuer                  string
	MaximumDecisionLifetime time.Duration
	SupportedConditions     []string
	EnforcedObligationTypes []string
}

type Service struct {
	subjects SubjectResolver
	policies PolicyResolver
	store    Store
	profile  Profile
	now      func() time.Time
}

func NewService(subjects SubjectResolver, policies PolicyResolver, store Store, profile Profile, now func() time.Time) (*Service, error) {
	if subjects == nil || policies == nil || store == nil || now == nil {
		return nil, errors.New("approval dependencies are required")
	}
	if profile.Issuer == "" || profile.MaximumDecisionLifetime <= 0 || profile.MaximumDecisionLifetime > 2*time.Hour {
		return nil, errors.New("approval profile is outside the M1 validation envelope")
	}
	conditions := slices.Clone(profile.SupportedConditions)
	obligations := slices.Clone(profile.EnforcedObligationTypes)
	sort.Strings(conditions)
	sort.Strings(obligations)
	if !uniqueNonEmpty(conditions) || !uniqueNonEmpty(obligations) {
		return nil, errors.New("approval profile collections must be unique and non-empty")
	}
	profile.SupportedConditions = conditions
	profile.EnforcedObligationTypes = obligations
	return &Service{subjects: subjects, policies: policies, store: store, profile: profile, now: now}, nil
}

type RequestPlanReview struct {
	ContractVersion    string
	CorrelationID      string
	IdempotencyKey     string
	TenantID           string
	RunID              string
	RequestingActorID  string
	IntendedApproverID string
	RequestedExpiry    time.Time
}

func (service *Service) RequestPlanReview(ctx context.Context, command RequestPlanReview) (trust.ApprovalRequest, bool, error) {
	if service == nil || service.subjects == nil || service.store == nil || service.now == nil ||
		ctx == nil || command.ContractVersion != RequestPlanReviewContractVersion ||
		command.CorrelationID == "" || command.IdempotencyKey == "" || command.TenantID == "" ||
		command.RunID == "" || command.RequestingActorID == "" || command.IntendedApproverID == "" {
		return trust.ApprovalRequest{}, false, invalid(command.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return trust.ApprovalRequest{}, false, unavailable(command.CorrelationID)
	}
	subject, found, err := service.subjects.ResolvePlanSubject(ctx, command.TenantID, command.RunID)
	if err != nil {
		return trust.ApprovalRequest{}, false, unavailable(command.CorrelationID)
	}
	if !found || !validSubject(subject, command) {
		return trust.ApprovalRequest{}, false, forbidden(command.CorrelationID)
	}
	expiresAt := command.RequestedExpiry
	for _, candidate := range []time.Time{at.Add(service.profile.MaximumDecisionLifetime), subject.ExpiresAt} {
		if expiresAt.IsZero() || candidate.Before(expiresAt) {
			expiresAt = candidate
		}
	}
	evidence := canonical(subject.EvidenceReferences)
	evidence = canonical(append(evidence,
		"context:"+subject.ContextPackageID,
		"context-manifest:"+subject.ContextManifest,
		"plan:"+subject.PlanDigest,
	))
	value := trust.ApprovalRequest{
		SchemaVersion: trust.ApprovalRequestSchemaVersion, Issuer: service.profile.Issuer,
		RequestingActorID: command.RequestingActorID, IntendedApproverID: command.IntendedApproverID,
		Scope: subject.Scope, Environment: "developer-trial", Purpose: subject.Purpose,
		Classification: subject.Classification, RiskClass: subject.RiskClass,
		Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan",
		SubjectID: subject.PlanID, SubjectDigest: subject.PlanDigest,
		EvidenceReferences: evidence, PermittedConditions: slices.Clone(service.profile.SupportedConditions),
		PolicyDecisionIDs: canonical(subject.PolicyDecisionIDs), CreatedAt: at, ExpiresAt: expiresAt,
		MaxDecisionLifetime: int64(service.profile.MaximumDecisionLifetime / time.Second),
		Lifecycle:           "open", EvidenceCorrelationID: command.CorrelationID,
	}
	value, err = trust.BindApprovalRequestDigest(value, command.CorrelationID)
	if err != nil {
		return trust.ApprovalRequest{}, false, err
	}
	stored, replayed, err := service.store.CreateRequest(ctx, command.TenantID, command.IdempotencyKey, value)
	if err != nil {
		return trust.ApprovalRequest{}, false, normalizeStore(err, command.CorrelationID)
	}
	if err := stored.VerifyDigest(command.CorrelationID); err != nil {
		return trust.ApprovalRequest{}, false, err
	}
	if stored.ID != value.ID || stored.RequestDigest != value.RequestDigest {
		return trust.ApprovalRequest{}, false, failure.New(failure.Indeterminate, "stored approval request differs from reviewed facts", command.CorrelationID, failure.RetryNone)
	}
	return cloneRequest(stored), replayed, nil
}

type DecideApproval struct {
	ContractVersion    string
	CorrelationID      string
	IdempotencyKey     string
	TenantID           string
	RequestID          string
	Approver           identity.Principal
	PolicyDecisionID   string
	Decision           string
	Conditions         []string
	EvidenceReferences []string
	RequestedExpiry    time.Time
}

func (service *Service) DecideApproval(ctx context.Context, command DecideApproval) (trust.ApprovalRecordV2, bool, error) {
	if service == nil || service.policies == nil || service.store == nil || service.now == nil ||
		ctx == nil || command.ContractVersion != DecideApprovalContractVersion ||
		command.CorrelationID == "" || command.IdempotencyKey == "" || command.TenantID == "" ||
		command.RequestID == "" || command.PolicyDecisionID == "" ||
		(command.Decision != "approved" && command.Decision != "rejected") {
		return trust.ApprovalRecordV2{}, false, invalid(command.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return trust.ApprovalRecordV2{}, false, unavailable(command.CorrelationID)
	}
	if err := command.Approver.ValidateAt(at, command.CorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	request, found, err := service.store.ResolveRequest(ctx, command.TenantID, command.RequestID)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, unavailable(command.CorrelationID)
	}
	if !found || request.Scope.TenantID != command.TenantID || request.IntendedApproverID != command.Approver.SubjectID() ||
		request.EvidenceCorrelationID != command.CorrelationID {
		return trust.ApprovalRecordV2{}, false, forbidden(command.CorrelationID)
	}
	if err := request.VerifyDigest(command.CorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	if !at.Before(request.ExpiresAt) {
		return trust.ApprovalRecordV2{}, false, failure.New(failure.Expired, "approval request is no longer active", command.CorrelationID, failure.RetryNewAuthority)
	}
	subject, subjectFound, err := service.subjects.ResolvePlanSubject(ctx, command.TenantID, request.Scope.DeliveryRunID)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, unavailable(command.CorrelationID)
	}
	if !subjectFound || !matchesRequestSubject(subject, request) {
		return trust.ApprovalRecordV2{}, false, forbidden(command.CorrelationID)
	}
	if !subset(command.Conditions, request.PermittedConditions) ||
		!sameSet(command.EvidenceReferences, request.EvidenceReferences) {
		return trust.ApprovalRecordV2{}, false, forbidden(command.CorrelationID)
	}
	decision, found, err := service.policies.ResolvePolicyDecision(ctx, command.PolicyDecisionID)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, unavailable(command.CorrelationID)
	}
	if !found || decision.ID != command.PolicyDecisionID ||
		decision.EvidenceCorrelationID != command.CorrelationID {
		return trust.ApprovalRecordV2{}, false, forbidden(command.CorrelationID)
	}
	if err := decision.Authorize("approval", command.Approver.SubjectID(), "forge.approval.decide",
		request.RequestDigest, at, service.profile.EnforcedObligationTypes, command.CorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	expiresAt := command.RequestedExpiry
	for _, candidate := range []time.Time{request.ExpiresAt, command.Approver.ExpiresAt()} {
		if decision.ExpiresAt != nil && decision.ExpiresAt.Before(candidate) {
			candidate = *decision.ExpiresAt
		}
		if expiresAt.IsZero() || candidate.Before(expiresAt) {
			expiresAt = candidate
		}
	}
	if !expiresAt.After(at) {
		return trust.ApprovalRecordV2{}, false, failure.New(failure.Expired, "approval decision has no active lifetime", command.CorrelationID, failure.RetryNewAuthority)
	}
	record := trust.ApprovalRecordV2{
		SchemaVersion: trust.ApprovalRecordV2SchemaVersion, Issuer: service.profile.Issuer,
		RequestID: request.ID, RequestDigest: request.RequestDigest, ApproverID: command.Approver.SubjectID(),
		Scope: request.Scope, Environment: request.Environment, Purpose: request.Purpose,
		Classification: request.Classification, RiskClass: request.RiskClass,
		Action: request.Action, SubjectType: request.SubjectType, SubjectID: request.SubjectID,
		SubjectDigest: request.SubjectDigest, Decision: command.Decision, DecidedAt: at, ExpiresAt: expiresAt,
		Conditions: canonical(command.Conditions), PolicyDecisionIDs: []string{decision.ID},
		EvidenceReferences: canonical(command.EvidenceReferences), EvidenceCorrelationID: command.CorrelationID,
	}
	record, err = trust.BindApprovalRecordV2Integrity(record, command.CorrelationID)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	stored, replayed, err := service.store.CreateDecision(ctx, command.TenantID, command.IdempotencyKey, request.ID, record)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, normalizeStore(err, command.CorrelationID)
	}
	if err := stored.VerifyIntegrity(command.CorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	return cloneRecord(stored), replayed, nil
}

type RevokeApproval struct {
	ContractVersion  string
	CorrelationID    string
	IdempotencyKey   string
	TenantID         string
	RecordID         string
	Approver         identity.Principal
	PolicyDecisionID string
}

func (service *Service) RevokeApproval(ctx context.Context, command RevokeApproval) (trust.ApprovalRecordV2, bool, error) {
	if service == nil || command.ContractVersion != RevokeApprovalContractVersion ||
		command.RecordID == "" || command.IdempotencyKey == "" || command.TenantID == "" {
		return trust.ApprovalRecordV2{}, false, invalid(command.CorrelationID)
	}
	current, found, err := service.ResolveApproval(ctx, command.TenantID, command.RecordID)
	if err != nil || !found {
		return trust.ApprovalRecordV2{}, false, forbidden(command.CorrelationID)
	}
	at := service.now()
	if err := command.Approver.ValidateAt(at, command.CorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	if current.Decision != "approved" || current.ApproverID != command.Approver.SubjectID() || !at.Before(current.ExpiresAt) {
		return trust.ApprovalRecordV2{}, false, forbidden(command.CorrelationID)
	}
	decision, found, err := service.policies.ResolvePolicyDecision(ctx, command.PolicyDecisionID)
	if err != nil || !found {
		return trust.ApprovalRecordV2{}, false, unavailable(command.CorrelationID)
	}
	if err := decision.Authorize("approval", command.Approver.SubjectID(), "forge.approval.decide",
		current.RequestDigest, at, service.profile.EnforcedObligationTypes, command.CorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	revocation := current
	revocation.ID = ""
	revocation.Integrity = trust.Integrity{}
	revocation.Decision = "revoked"
	revocation.DecidedAt = at
	revocation.PreviousRecordID = current.ID
	revocation.PolicyDecisionIDs = []string{decision.ID}
	revocation, err = trust.BindApprovalRecordV2Integrity(revocation, command.CorrelationID)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	stored, replayed, err := service.store.CreateDecision(ctx, command.TenantID, command.IdempotencyKey, current.RequestID, revocation)
	if err != nil {
		return trust.ApprovalRecordV2{}, false, normalizeStore(err, command.CorrelationID)
	}
	return cloneRecord(stored), replayed, nil
}

func (service *Service) ResolveApproval(ctx context.Context, tenantID, recordID string) (trust.ApprovalRecordV2, bool, error) {
	if service == nil || service.store == nil || ctx == nil || tenantID == "" || recordID == "" {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	value, found, err := service.store.ResolveCurrent(ctx, tenantID, recordID)
	if err != nil || !found {
		return trust.ApprovalRecordV2{}, found, err
	}
	if err := value.VerifyIntegrity(value.EvidenceCorrelationID); err != nil {
		return trust.ApprovalRecordV2{}, false, err
	}
	return cloneRecord(value), true, nil
}

func (service *Service) ResolveApprovalRequest(ctx context.Context, tenantID, requestID string) (trust.ApprovalRequest, bool, error) {
	if service == nil || service.store == nil || ctx == nil || tenantID == "" || requestID == "" {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	value, found, err := service.store.ResolveRequest(ctx, tenantID, requestID)
	if err != nil || !found {
		return trust.ApprovalRequest{}, found, err
	}
	if err := value.VerifyDigest(value.EvidenceCorrelationID); err != nil {
		return trust.ApprovalRequest{}, false, err
	}
	return cloneRequest(value), true, nil
}

type requestEntry struct {
	fingerprint string
	value       trust.ApprovalRequest
}

type decisionEntry struct {
	fingerprint string
	value       trust.ApprovalRecordV2
}

type MemoryStore struct {
	mu               sync.RWMutex
	requestKeys      map[string]requestEntry
	requests         map[string]trust.ApprovalRequest
	decisionKeys     map[string]decisionEntry
	records          map[string]trust.ApprovalRecordV2
	currentByRequest map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		requestKeys: make(map[string]requestEntry), requests: make(map[string]trust.ApprovalRequest),
		decisionKeys: make(map[string]decisionEntry), records: make(map[string]trust.ApprovalRecordV2),
		currentByRequest: make(map[string]string),
	}
}

func (store *MemoryStore) CreateRequest(ctx context.Context, tenantID, key string, value trust.ApprovalRequest) (trust.ApprovalRequest, bool, error) {
	if store == nil || ctx == nil {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	idempotency := tenantID + "\x00" + key
	if existing, found := store.requestKeys[idempotency]; found {
		if existing.fingerprint != value.RequestDigest {
			return trust.ApprovalRequest{}, false, ErrIdempotencyConflict
		}
		return cloneRequest(existing.value), true, nil
	}
	requestKey := tenantID + "\x00" + value.ID
	if _, exists := store.requests[requestKey]; exists {
		return trust.ApprovalRequest{}, false, ErrDecisionConflict
	}
	store.requests[requestKey] = cloneRequest(value)
	store.requestKeys[idempotency] = requestEntry{fingerprint: value.RequestDigest, value: cloneRequest(value)}
	return cloneRequest(value), false, nil
}

func (store *MemoryStore) ResolveRequest(ctx context.Context, tenantID, requestID string) (trust.ApprovalRequest, bool, error) {
	if store == nil || ctx == nil {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	value, found := store.requests[tenantID+"\x00"+requestID]
	return cloneRequest(value), found, nil
}

func (store *MemoryStore) CreateDecision(ctx context.Context, tenantID, key, requestID string, value trust.ApprovalRecordV2) (trust.ApprovalRecordV2, bool, error) {
	if store == nil || ctx == nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	idempotency := tenantID + "\x00" + key
	if existing, found := store.decisionKeys[idempotency]; found {
		if existing.fingerprint != value.Integrity.Digest {
			return trust.ApprovalRecordV2{}, false, ErrIdempotencyConflict
		}
		return cloneRecord(existing.value), true, nil
	}
	currentKey := tenantID + "\x00" + requestID
	if currentID, exists := store.currentByRequest[currentKey]; exists {
		current := store.records[tenantID+"\x00"+currentID]
		if value.Decision != "revoked" || value.PreviousRecordID != current.ID || current.Decision != "approved" {
			return trust.ApprovalRecordV2{}, false, ErrDecisionConflict
		}
	}
	store.records[tenantID+"\x00"+value.ID] = cloneRecord(value)
	store.currentByRequest[currentKey] = value.ID
	store.decisionKeys[idempotency] = decisionEntry{fingerprint: value.Integrity.Digest, value: cloneRecord(value)}
	return cloneRecord(value), false, nil
}

func (store *MemoryStore) ResolveCurrent(ctx context.Context, tenantID, recordID string) (trust.ApprovalRecordV2, bool, error) {
	if store == nil || ctx == nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	value, found := store.records[tenantID+"\x00"+recordID]
	if !found || store.currentByRequest[tenantID+"\x00"+value.RequestID] != recordID {
		return trust.ApprovalRecordV2{}, false, nil
	}
	return cloneRecord(value), true, nil
}

func validSubject(subject PlanSubject, command RequestPlanReview) bool {
	return subject.RunID == command.RunID && subject.Scope.TenantID == command.TenantID &&
		subject.Scope.DeliveryRunID == command.RunID && subject.PlanID != "" && subject.PlanDigest != "" &&
		subject.RequestingActorID == command.RequestingActorID &&
		subject.RequestingActorID == command.IntendedApproverID &&
		subject.CorrelationID == command.CorrelationID && subject.Purpose != "" &&
		subject.Classification != "" && subject.RiskClass != "" && subject.ContextPackageID != "" &&
		subject.ContextManifest != "" && len(subject.EvidenceReferences) > 0 &&
		len(subject.PolicyDecisionIDs) > 0 && !subject.ExpiresAt.IsZero()
}

func matchesRequestSubject(subject PlanSubject, request trust.ApprovalRequest) bool {
	return subject.RunID == request.Scope.DeliveryRunID && subject.PlanID == request.SubjectID &&
		subject.PlanDigest == request.SubjectDigest && subject.Scope == request.Scope &&
		subject.RequestingActorID == request.RequestingActorID &&
		subject.Purpose == request.Purpose && subject.Classification == request.Classification &&
		subject.RiskClass == request.RiskClass && subject.CorrelationID == request.EvidenceCorrelationID &&
		subject.ContextPackageID != "" && subject.ContextManifest != "" &&
		slices.Contains(request.EvidenceReferences, "context:"+subject.ContextPackageID) &&
		slices.Contains(request.EvidenceReferences, "context-manifest:"+subject.ContextManifest) &&
		slices.Contains(request.EvidenceReferences, "plan:"+subject.PlanDigest) &&
		!subject.ExpiresAt.Before(request.ExpiresAt) &&
		subset(request.PolicyDecisionIDs, subject.PolicyDecisionIDs) &&
		subset(subject.EvidenceReferences, request.EvidenceReferences)
}

func cloneRequest(value trust.ApprovalRequest) trust.ApprovalRequest {
	value.EvidenceReferences = slices.Clone(value.EvidenceReferences)
	value.PermittedConditions = slices.Clone(value.PermittedConditions)
	value.PolicyDecisionIDs = slices.Clone(value.PolicyDecisionIDs)
	return value
}

func cloneRecord(value trust.ApprovalRecordV2) trust.ApprovalRecordV2 {
	value.Conditions = slices.Clone(value.Conditions)
	value.PolicyDecisionIDs = slices.Clone(value.PolicyDecisionIDs)
	value.EvidenceReferences = slices.Clone(value.EvidenceReferences)
	return value
}

func canonical(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return slices.Compact(result)
}

func uniqueNonEmpty(values []string) bool {
	for index, value := range values {
		if value == "" || (index > 0 && value == values[index-1]) {
			return false
		}
	}
	return true
}

func subset(values, allowed []string) bool {
	for _, value := range values {
		if !slices.Contains(allowed, value) {
			return false
		}
	}
	return true
}

func sameSet(left, right []string) bool {
	return slices.Equal(canonical(left), canonical(right))
}

func normalizeStore(err error, correlationID string) error {
	switch {
	case errors.Is(err, ErrIdempotencyConflict), errors.Is(err, ErrDecisionConflict):
		return failure.New(failure.Conflict, "approval command conflicts with current state", correlationID, failure.RetryReconcile)
	default:
		return unavailable(correlationID)
	}
}

func invalid(correlationID string) error {
	return failure.New(failure.InvalidRequest, "approval command is invalid", correlationID, failure.RetryNone)
}

func forbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "approval command is not authorized", correlationID, failure.RetryNone)
}

func unavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "approval capability is unavailable", correlationID, failure.RetryBounded)
}
