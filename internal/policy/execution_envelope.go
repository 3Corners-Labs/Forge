package policy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const IssueExecutionEnvelopeContractVersion = "forge.issue-execution-envelope-request.v1"

var (
	ErrEnvelopeStoreUnavailable = errors.New("execution envelope store unavailable")
	ErrEnvelopeConflict         = errors.New("execution envelope conflict")
)

// DelegationResolver is the narrow Identity boundary needed at issuance.
// Policy cannot construct, persist, revoke, or reinterpret delegations.
type DelegationResolver interface {
	ResolveDelegationChain(context.Context, identity.ResolveDelegationRequest) ([]identity.Delegation, error)
}

// EnvelopeProfile contains trusted Developer Trial issuance constraints. It
// is configuration supplied when the issuer is constructed, never caller
// authority carried on an issuance request.
type EnvelopeProfile struct {
	Issuer                  string
	Environment             Environment
	Risk                    Risk
	AllowedPolicyVersions   []string
	EnforcedObligationTypes []string
	MaximumLifetime         time.Duration
	Limits                  trust.ExecutionLimits
}

type EnvelopeStore interface {
	Create(context.Context, trust.ExecutionEnvelope, string, string, string) (trust.ExecutionEnvelope, bool, error)
}

type EnvelopeIssuer struct {
	evaluator *Evaluator
	resolver  DelegationResolver
	store     EnvelopeStore
	profile   EnvelopeProfile
	now       func() time.Time
}

func NewEnvelopeIssuer(evaluator *Evaluator, resolver DelegationResolver, store EnvelopeStore, profile EnvelopeProfile, now func() time.Time) (*EnvelopeIssuer, error) {
	if evaluator == nil {
		return nil, errors.New("policy evaluator is required")
	}
	if resolver == nil {
		return nil, errors.New("identity delegation resolver is required")
	}
	if store == nil {
		return nil, errors.New("execution envelope store is required")
	}
	if now == nil {
		return nil, errors.New("execution envelope clock is required")
	}
	versions := slices.Clone(profile.AllowedPolicyVersions)
	sort.Strings(versions)
	obligations := slices.Clone(profile.EnforcedObligationTypes)
	sort.Strings(obligations)
	if profile.Issuer == "" || profile.Environment != DeveloperTrialEnvironment || profile.Risk != RiskLow ||
		profile.MaximumLifetime <= 0 || len(versions) == 0 || !allUniqueNonZero(versions) ||
		!allUniqueNonZero(obligations) ||
		profile.Limits.MaxDurationSeconds < 1 || profile.Limits.MaxCostMinorUnits != 0 ||
		profile.Limits.MaxInputTokens != 0 || profile.Limits.MaxOutputTokens != 0 || profile.Limits.MaxConcurrency != 1 {
		return nil, errors.New("execution envelope profile is outside the M1 Developer Trial contract")
	}
	profile.AllowedPolicyVersions = versions
	profile.EnforcedObligationTypes = obligations
	return &EnvelopeIssuer{evaluator: evaluator, resolver: resolver, store: store, profile: profile, now: now}, nil
}

type IssueExecutionEnvelopeRequest struct {
	ContractVersion  string
	CorrelationID    string
	IdempotencyKey   string
	Human            identity.Principal
	Agent            identity.Principal
	Scope            tenant.Scope
	Run              delivery.DeliveryRun
	DelegationID     string
	PolicyDecisionID string
}

type IssueExecutionEnvelopeResult struct {
	Envelope trust.ExecutionEnvelope
	Replayed bool
}

func (issuer *EnvelopeIssuer) IssueExecutionEnvelope(ctx context.Context, request IssueExecutionEnvelopeRequest) (IssueExecutionEnvelopeResult, error) {
	if issuer == nil || issuer.evaluator == nil || issuer.resolver == nil || issuer.store == nil || issuer.now == nil {
		return IssueExecutionEnvelopeResult{}, envelopeUnavailable(request.CorrelationID)
	}
	if ctx == nil || request.ContractVersion != IssueExecutionEnvelopeContractVersion || request.CorrelationID == "" ||
		request.IdempotencyKey == "" || request.DelegationID == "" || request.PolicyDecisionID == "" {
		return IssueExecutionEnvelopeResult{}, envelopeInvalid(request.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return IssueExecutionEnvelopeResult{}, envelopeUnavailable(request.CorrelationID)
	}
	at := issuer.now()
	if at.IsZero() {
		return IssueExecutionEnvelopeResult{}, envelopeUnavailable(request.CorrelationID)
	}
	if err := request.Human.ValidateAt(at, request.CorrelationID); err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	if err := request.Agent.ValidateAt(at, request.CorrelationID); err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	if err := request.Scope.ValidateFor(request.Human, at, request.CorrelationID); err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	if !validEnvelopeRunBinding(request) {
		return IssueExecutionEnvelopeResult{}, envelopeForbidden(request.CorrelationID)
	}
	chain, err := issuer.resolver.ResolveDelegationChain(ctx, identity.ResolveDelegationRequest{
		ContractVersion: identity.ResolveDelegationContractVersion,
		CorrelationID:   request.CorrelationID,
		Human:           request.Human,
		Agent:           request.Agent,
		DelegationID:    request.DelegationID,
	})
	if err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	if len(chain) != 1 || chain[0].ID() != request.DelegationID {
		return IssueExecutionEnvelopeResult{}, envelopeForbidden(request.CorrelationID)
	}
	delegation := chain[0]
	if !validEnvelopeDelegationBinding(request, delegation, issuer.profile) {
		return IssueExecutionEnvelopeResult{}, envelopeForbidden(request.CorrelationID)
	}
	decision, err := issuer.evaluator.Evaluate(ctx, EvaluateRequest{
		ContractVersion: EvaluateRequestContractVersion, CorrelationID: request.CorrelationID,
		Principal: request.Human, Scope: request.Scope, Run: request.Run,
		Action: ActionAssembleContext, Environment: issuer.profile.Environment, Risk: issuer.profile.Risk,
	})
	if err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	if decision.ID != request.PolicyDecisionID || decision.ExpiresAt == nil ||
		decision.EvidenceCorrelationID != request.CorrelationID ||
		!allowedPolicyVersions(decision.PolicyVersions, issuer.profile.AllowedPolicyVersions) {
		return IssueExecutionEnvelopeResult{}, envelopeForbidden(request.CorrelationID)
	}
	if err := decision.Authorize(
		"retrieval", request.Human.SubjectID(), string(ActionAssembleContext), request.Run.RepositoryID(), at,
		issuer.profile.EnforcedObligationTypes, request.CorrelationID,
	); err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}

	expiresAt := earliestEnvelopeExpiry(
		request.Human.ExpiresAt(), request.Agent.ExpiresAt(), delegation.ExpiresAt(),
		*decision.ExpiresAt, at.Add(issuer.profile.MaximumLifetime),
	)
	remainingSeconds := int64(expiresAt.Sub(at) / time.Second)
	if remainingSeconds < 1 {
		return IssueExecutionEnvelopeResult{}, failure.New(failure.Expired, "execution authority is not active", request.CorrelationID, failure.RetryNewAuthority)
	}
	limits := issuer.profile.Limits
	if limits.MaxDurationSeconds > remainingSeconds {
		limits.MaxDurationSeconds = remainingSeconds
	}
	policyDecisionIDs := []string{delegation.ParentAuthorityID(), decision.ID}
	sort.Strings(policyDecisionIDs)
	policyDecisionIDs = slices.Compact(policyDecisionIDs)
	envelope := trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion,
		Issuer:        issuer.profile.Issuer,
		IssuedAt:      at,
		ExpiresAt:     expiresAt,
		Scope: trust.ExecutionScope{
			TenantID: request.Scope.TenantID(), WorkspaceID: request.Scope.WorkspaceID(),
			ProjectID: request.Scope.ProjectID(), DeliveryRunID: request.Run.ID(),
		},
		Principal: trust.EnvelopePrincipal{
			SubjectID: request.Human.SubjectID(), AgentID: request.Agent.SubjectID(),
			DelegationChain: []string{delegation.ID()},
		},
		Purpose: delegation.Purpose(), RiskClass: string(issuer.profile.Risk),
		Actions: delegation.Actions(), Classifications: delegation.Classifications(),
		Destinations: []string{}, ModelIDs: []string{}, ToolIDs: []string{}, SecretReferenceIDs: []string{},
		PolicyDecisionIDs: policyDecisionIDs, ApprovalRecordIDs: []string{}, Limits: limits,
		IdempotencyKey: request.IdempotencyKey, EvidenceCorrelationID: request.CorrelationID,
	}
	envelope, err = trust.BindExecutionEnvelopeIntegrity(envelope, request.CorrelationID)
	if err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	fingerprint := envelopeIssuanceFingerprint(envelope, issuer.profile)
	stored, replayed, err := issuer.store.Create(ctx, envelope, request.Scope.TenantID(), request.IdempotencyKey, fingerprint)
	if err != nil {
		return IssueExecutionEnvelopeResult{}, normalizeEnvelopeStoreError(err, request.CorrelationID)
	}
	if err := validateIssuedEnvelope(stored, request, delegation, decision, issuer.profile, at); err != nil {
		return IssueExecutionEnvelopeResult{}, err
	}
	return IssueExecutionEnvelopeResult{Envelope: cloneEnvelope(stored), Replayed: replayed}, nil
}

type memoryEnvelopeEntry struct {
	envelope    trust.ExecutionEnvelope
	fingerprint string
}

type MemoryEnvelopeStore struct {
	mu            sync.RWMutex
	byIdempotency map[string]memoryEnvelopeEntry
}

func NewMemoryEnvelopeStore() *MemoryEnvelopeStore {
	return &MemoryEnvelopeStore{byIdempotency: make(map[string]memoryEnvelopeEntry)}
}

func (store *MemoryEnvelopeStore) Create(_ context.Context, value trust.ExecutionEnvelope, tenantID, idempotencyKey, fingerprint string) (trust.ExecutionEnvelope, bool, error) {
	if store == nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	key := tenantID + "\x00" + idempotencyKey
	if existing, found := store.byIdempotency[key]; found {
		if existing.fingerprint != fingerprint {
			return trust.ExecutionEnvelope{}, false, ErrEnvelopeConflict
		}
		return cloneEnvelope(existing.envelope), true, nil
	}
	store.byIdempotency[key] = memoryEnvelopeEntry{envelope: cloneEnvelope(value), fingerprint: fingerprint}
	return cloneEnvelope(value), false, nil
}

func validEnvelopeRunBinding(request IssueExecutionEnvelopeRequest) bool {
	return request.Human.Kind() == identity.HumanPrincipal && request.Human.Synthetic() &&
		request.Agent.Kind() == identity.AgentPrincipal && request.Agent.Synthetic() &&
		request.Human.SubjectID() != request.Agent.SubjectID() &&
		request.Run.ContractVersion() == delivery.DeliveryRunContractVersion && request.Run.State() == delivery.CapturedState &&
		request.Run.ID() != "" && request.Run.RepositoryID() != "" && request.Run.Purpose() != "" && request.Run.Classification() != "" &&
		request.Run.TenantID() == request.Scope.TenantID() && request.Run.WorkspaceID() == request.Scope.WorkspaceID() &&
		request.Run.ProjectID() == request.Scope.ProjectID() && request.Run.ActorID() == request.Human.SubjectID()
}

func validEnvelopeDelegationBinding(request IssueExecutionEnvelopeRequest, delegation identity.Delegation, profile EnvelopeProfile) bool {
	return delegation.HumanSubjectID() == request.Human.SubjectID() && delegation.AgentSubjectID() == request.Agent.SubjectID() &&
		delegation.TenantID() == request.Scope.TenantID() && delegation.WorkspaceID() == request.Scope.WorkspaceID() &&
		delegation.ProjectID() == request.Scope.ProjectID() && delegation.RepositoryID() == request.Run.RepositoryID() &&
		delegation.DeliveryRunID() == request.Run.ID() && delegation.Environment() == string(profile.Environment) &&
		delegation.MaximumAutonomy() == identity.ContextPreparationOnly && delegation.EvidenceCorrelationID() == request.CorrelationID &&
		slices.Equal(delegation.Actions(), []string{identity.DelegatedActionAssembleContext}) &&
		slices.Equal(delegation.Resources(), []string{request.Run.RepositoryID()}) &&
		slices.Equal(delegation.Classifications(), []string{request.Run.Classification()}) && len(delegation.Destinations()) == 0
}

func validateIssuedEnvelope(value trust.ExecutionEnvelope, request IssueExecutionEnvelopeRequest, delegation identity.Delegation, decision trust.PolicyDecision, profile EnvelopeProfile, at time.Time) error {
	if err := value.VerifyIntegrity(request.CorrelationID); err != nil {
		return err
	}
	if err := value.AuthorizeAction(identity.DelegatedActionAssembleContext, at, request.CorrelationID); err != nil {
		return err
	}
	expectedPolicyDecisionIDs := []string{delegation.ParentAuthorityID(), decision.ID}
	sort.Strings(expectedPolicyDecisionIDs)
	expectedPolicyDecisionIDs = slices.Compact(expectedPolicyDecisionIDs)
	if value.ParentID != "" || value.Issuer != profile.Issuer || value.Scope.TenantID != request.Scope.TenantID() ||
		value.Scope.WorkspaceID != request.Scope.WorkspaceID() || value.Scope.ProjectID != request.Scope.ProjectID() ||
		value.Scope.DeliveryRunID != request.Run.ID() || value.Principal.SubjectID != request.Human.SubjectID() ||
		value.Principal.AgentID != request.Agent.SubjectID() ||
		!slices.Equal(value.Principal.DelegationChain, []string{delegation.ID()}) || value.Purpose != delegation.Purpose() ||
		value.RiskClass != string(profile.Risk) || !slices.Equal(value.Actions, delegation.Actions()) ||
		!slices.Equal(value.Classifications, delegation.Classifications()) || len(value.Destinations) != 0 ||
		len(value.ModelIDs) != 0 || len(value.ToolIDs) != 0 || len(value.SecretReferenceIDs) != 0 || len(value.ApprovalRecordIDs) != 0 ||
		!slices.Equal(value.PolicyDecisionIDs, expectedPolicyDecisionIDs) || value.IdempotencyKey != request.IdempotencyKey ||
		value.EvidenceCorrelationID != request.CorrelationID ||
		value.Limits.MaxDurationSeconds > profile.Limits.MaxDurationSeconds || value.Limits.MaxCostMinorUnits != 0 ||
		value.Limits.MaxInputTokens != 0 || value.Limits.MaxOutputTokens != 0 || value.Limits.MaxConcurrency != 1 {
		return failure.New(failure.Indeterminate, "issued execution envelope could not be verified", request.CorrelationID, failure.RetryNone)
	}
	return nil
}

func allowedPolicyVersions(actual, allowed []string) bool {
	return len(actual) > 0 && allUniqueNonZero(actual) &&
		!slices.ContainsFunc(actual, func(value string) bool { return !slices.Contains(allowed, value) })
}

func earliestEnvelopeExpiry(values ...time.Time) time.Time {
	result := values[0]
	for _, value := range values[1:] {
		if value.Before(result) {
			result = value
		}
	}
	return result
}

func envelopeIssuanceFingerprint(value trust.ExecutionEnvelope, profile EnvelopeProfile) string {
	value.ID = ""
	value.IssuedAt = time.Time{}
	value.ExpiresAt = time.Time{}
	value.IdempotencyKey = ""
	value.EvidenceCorrelationID = ""
	value.Integrity = nil
	encoded, _ := json.Marshal(struct {
		Envelope                   trust.ExecutionEnvelope
		Environment                Environment
		AllowedPolicyVersions      []string
		EnforcedObligationTypes    []string
		MaximumLifetimeNanoseconds int64
	}{
		Envelope: value, Environment: profile.Environment,
		AllowedPolicyVersions:      profile.AllowedPolicyVersions,
		EnforcedObligationTypes:    profile.EnforcedObligationTypes,
		MaximumLifetimeNanoseconds: int64(profile.MaximumLifetime),
	})
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func cloneEnvelope(value trust.ExecutionEnvelope) trust.ExecutionEnvelope {
	value.Principal.DelegationChain = slices.Clone(value.Principal.DelegationChain)
	value.Actions = slices.Clone(value.Actions)
	value.Classifications = slices.Clone(value.Classifications)
	value.Destinations = slices.Clone(value.Destinations)
	value.ModelIDs = slices.Clone(value.ModelIDs)
	value.ToolIDs = slices.Clone(value.ToolIDs)
	value.SecretReferenceIDs = slices.Clone(value.SecretReferenceIDs)
	value.PolicyDecisionIDs = slices.Clone(value.PolicyDecisionIDs)
	value.ApprovalRecordIDs = slices.Clone(value.ApprovalRecordIDs)
	if value.Integrity != nil {
		integrity := *value.Integrity
		value.Integrity = &integrity
	}
	return value
}

func normalizeEnvelopeStoreError(err error, correlationID string) error {
	switch {
	case errors.Is(err, ErrEnvelopeConflict):
		return failure.New(failure.Conflict, "execution envelope request conflicts with existing authority", correlationID, failure.RetryNone)
	case errors.Is(err, ErrEnvelopeStoreUnavailable):
		return envelopeUnavailable(correlationID)
	default:
		return failure.New(failure.Indeterminate, "execution envelope persistence could not be verified", correlationID, failure.RetryNone)
	}
}

func envelopeInvalid(correlationID string) error {
	return failure.New(failure.InvalidRequest, "execution envelope request is invalid", correlationID, failure.RetryNone)
}

func envelopeForbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "execution envelope is not authorized", correlationID, failure.RetryNone)
}

func envelopeUnavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "execution envelope issuance is unavailable", correlationID, failure.RetryBounded)
}
