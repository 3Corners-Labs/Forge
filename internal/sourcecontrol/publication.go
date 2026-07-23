// Package sourcecontrol owns provider-neutral, idempotent branch and
// pull-request publication. It does not decide policy, approve publication,
// verify candidates, own delivery state, or merge pull requests.
package sourcecontrol

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/tenant"
)

const (
	PublishRequestContractVersion = "forge.publish-pull-request.v1"
	ProviderContractVersion       = "forge.source-control-provider.v1"
	PublicationReceiptVersion     = "forge.publication-receipt.v1"
	PublishAction                 = "forge.source-control.publish-pr"
	DeveloperTrialEnvironment     = "developer-trial"
)

var (
	ErrStoreUnavailable    = errors.New("source-control operation store unavailable")
	ErrStoreConflict       = errors.New("source-control operation conflict")
	ErrProviderAmbiguous   = errors.New("source-control provider outcome is ambiguous")
	ErrProviderConflict    = errors.New("source-control provider state conflicts with the authorized operation")
	ErrProviderUnavailable = errors.New("source-control provider unavailable")
)

type SubjectResolver interface {
	ResolvePublicationSubject(context.Context, string, string) (trust.PublicationApprovalSubject, bool, error)
}

type ApprovalResolver interface {
	ResolveApproval(context.Context, string, string) (trust.ApprovalRecordV2, bool, error)
}

type PolicyResolver interface {
	ResolvePolicyDecision(context.Context, string) (trust.PolicyDecision, bool, error)
}

type EvidenceVerifier interface {
	VerifySubject(context.Context, trust.EvidenceScope, trust.EvidenceSubject, []string, string) error
}

type PublicationProfile struct {
	Provider     string
	BaseBranch   string
	BaseRevision string
	Title        string
	Body         string
}

type PublishRequest struct {
	ContractVersion   string
	CorrelationID     string
	IdempotencyKey    string
	TenantID          string
	RunID             string
	Envelope          trust.ExecutionEnvelope
	PolicyDecisionID  string
	ApprovalRecordID  string
	RepositoryBinding tenant.RepositoryBinding
	Provider          string
	BaseBranch        string
	BaseRevision      string
	TreeDigest        string
	ChangedPaths      []string
	Title             string
	Body              string
}

type Operation struct {
	ContractVersion  string
	ID               string
	TenantID         string
	RunID            string
	RepositoryID     string
	Provider         string
	BaseBranch       string
	BaseRevision     string
	Branch           string
	CandidateID      string
	CandidateDigest  string
	TreeDigest       string
	ChangedPaths     []string
	MetadataDigest   string
	ApprovalRecordID string
	EnvelopeID       string
	PolicyDecisionID string
	Marker           string
}

type RemoteState struct {
	BranchExists    bool
	Branch          string
	CommitID        string
	CandidateDigest string
	TreeDigest      string
	PullRequestID   string
	PullRequestURL  string
	PullRequestOpen bool
	BaseBranch      string
	Marker          string
	RepositoryID    string
}

type Provider interface {
	Reconcile(context.Context, Operation) (RemoteState, error)
	CreateBranch(context.Context, Operation) error
	OpenPullRequest(context.Context, Operation, string, string) error
}

type Receipt struct {
	ContractVersion string
	OperationID     string
	Provider        string
	RepositoryID    string
	BaseBranch      string
	BaseRevision    string
	Branch          string
	CommitID        string
	PullRequestID   string
	PullRequestURL  string
	CandidateDigest string
	TreeDigest      string
	PublishedAt     time.Time
	Reconciled      bool
}

type Store interface {
	Create(context.Context, string, string, Operation) (Operation, bool, error)
	SaveReceipt(context.Context, string, Receipt) (Receipt, bool, error)
}

type Service struct {
	subjects  SubjectResolver
	approvals ApprovalResolver
	policies  PolicyResolver
	evidence  EvidenceVerifier
	provider  Provider
	store     Store
	profile   PublicationProfile
	now       func() time.Time
}

func NewService(subjects SubjectResolver, approvals ApprovalResolver, policies PolicyResolver, evidence EvidenceVerifier, provider Provider, store Store, profile PublicationProfile, now func() time.Time) (*Service, error) {
	if subjects == nil || approvals == nil || policies == nil || evidence == nil || provider == nil || store == nil || now == nil ||
		profile.Provider != "contract-double" || profile.BaseBranch == "" || profile.BaseRevision == "" ||
		profile.Title == "" || profile.Body == "" || !safeMetadata(profile.Title) || !safeMetadata(profile.Body) {
		return nil, errors.New("source-control dependencies are required")
	}
	return &Service{
		subjects: subjects, approvals: approvals, policies: policies, evidence: evidence,
		provider: provider, store: store, profile: profile, now: now,
	}, nil
}

func (service *Service) PublishPullRequest(ctx context.Context, request PublishRequest) (Receipt, bool, error) {
	if !validRequest(ctx, request) {
		return Receipt{}, false, sourceFailure(failure.InvalidRequest, "publication request is invalid", request.CorrelationID, failure.RetryNone)
	}
	if request.Provider != service.profile.Provider || request.BaseBranch != service.profile.BaseBranch ||
		request.BaseRevision != service.profile.BaseRevision || request.Title != service.profile.Title ||
		request.Body != service.profile.Body {
		return Receipt{}, false, forbidden(request.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return Receipt{}, false, sourceFailure(failure.DependencyUnavailable, "publication clock is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	subject, found, err := service.subjects.ResolvePublicationSubject(ctx, request.TenantID, request.RunID)
	if err != nil {
		return Receipt{}, false, sourceFailure(failure.DependencyUnavailable, "publication subject is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || !requestMatchesSubject(request, subject) {
		return Receipt{}, false, forbidden(request.CorrelationID)
	}
	approval, found, err := service.approvals.ResolveApproval(ctx, request.TenantID, request.ApprovalRecordID)
	if err != nil {
		return Receipt{}, false, sourceFailure(failure.DependencyUnavailable, "publication approval is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || !approvalMatchesSubject(approval, subject, at) {
		return Receipt{}, false, forbidden(request.CorrelationID)
	}
	decision, found, err := service.policies.ResolvePolicyDecision(ctx, request.PolicyDecisionID)
	if err != nil {
		return Receipt{}, false, sourceFailure(failure.DependencyUnavailable, "publication policy decision is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || decision.ID != request.PolicyDecisionID {
		return Receipt{}, false, forbidden(request.CorrelationID)
	}
	if err := decision.Authorize("publication", request.Envelope.Principal.SubjectID, PublishAction, subject.CandidateDigest, at, []string{"evidence-ready", "human-publication-approval"}, request.CorrelationID); err != nil {
		return Receipt{}, false, err
	}
	if err := validateEnvelope(request.Envelope, request, subject, approval, decision, at); err != nil {
		return Receipt{}, false, err
	}
	evidenceScope := trust.EvidenceScope{
		TenantID: subject.Scope.TenantID, WorkspaceID: subject.Scope.WorkspaceID,
		ProjectID: subject.Scope.ProjectID, DeliveryRunID: subject.Scope.DeliveryRunID,
	}
	if err := service.evidence.VerifySubject(
		ctx, evidenceScope,
		trust.EvidenceSubject{Type: "verified-candidate", ID: subject.CandidateID, Digest: subject.CandidateDigest},
		subject.EvidenceReferences, request.CorrelationID,
	); err != nil {
		return Receipt{}, false, sourceFailure(failure.EvidenceUnavailable, "publication evidence is not ready", request.CorrelationID, failure.RetryReconcile)
	}

	operation := buildOperation(request, subject)
	stored, replayedOperation, err := service.store.Create(ctx, request.TenantID, request.IdempotencyKey, operation)
	if errors.Is(err, ErrStoreConflict) {
		return Receipt{}, false, sourceFailure(failure.Conflict, "publication idempotency conflicts with existing authority", request.CorrelationID, failure.RetryReconcile)
	}
	if err != nil || stored.ID != operation.ID {
		return Receipt{}, false, sourceFailure(failure.DependencyUnavailable, "publication operation persistence is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	state, reconciled, err := reconcileExact(ctx, service.provider, stored)
	if err != nil {
		return Receipt{}, false, normalizeProvider(err, request.CorrelationID)
	}
	if state.PullRequestID == "" {
		if !state.BranchExists {
			for attempt := 0; attempt < 2 && !state.BranchExists; attempt++ {
				mutationErr := service.provider.CreateBranch(ctx, stored)
				if mutationErr != nil && !errors.Is(mutationErr, ErrProviderAmbiguous) {
					return Receipt{}, false, normalizeProvider(mutationErr, request.CorrelationID)
				}
				reconciled = reconciled || mutationErr != nil
				state, _, err = reconcileExact(ctx, service.provider, stored)
				if err != nil {
					return Receipt{}, false, normalizeProvider(err, request.CorrelationID)
				}
			}
		}
		if !state.BranchExists {
			return Receipt{}, false, sourceFailure(failure.ManualRecoveryRequired, "publication branch outcome is indeterminate", request.CorrelationID, failure.RetryReconcile)
		}
		if state.PullRequestID == "" {
			for attempt := 0; attempt < 2 && state.PullRequestID == ""; attempt++ {
				mutationErr := service.provider.OpenPullRequest(ctx, stored, request.Title, request.Body)
				if mutationErr != nil && !errors.Is(mutationErr, ErrProviderAmbiguous) {
					return Receipt{}, false, normalizeProvider(mutationErr, request.CorrelationID)
				}
				reconciled = reconciled || mutationErr != nil
				state, _, err = reconcileExact(ctx, service.provider, stored)
				if err != nil {
					return Receipt{}, false, normalizeProvider(err, request.CorrelationID)
				}
			}
		}
	}
	if state.PullRequestID == "" {
		return Receipt{}, false, sourceFailure(failure.ManualRecoveryRequired, "pull-request outcome is indeterminate", request.CorrelationID, failure.RetryReconcile)
	}
	receipt := Receipt{
		ContractVersion: PublicationReceiptVersion, OperationID: stored.ID,
		Provider: stored.Provider, RepositoryID: stored.RepositoryID, BaseBranch: stored.BaseBranch,
		BaseRevision: stored.BaseRevision, Branch: stored.Branch, CommitID: state.CommitID,
		PullRequestID: state.PullRequestID, PullRequestURL: state.PullRequestURL,
		CandidateDigest: stored.CandidateDigest, TreeDigest: stored.TreeDigest,
		PublishedAt: at, Reconciled: reconciled || replayedOperation,
	}
	saved, replayedReceipt, err := service.store.SaveReceipt(ctx, stored.ID, receipt)
	if errors.Is(err, ErrStoreConflict) {
		return Receipt{}, false, sourceFailure(failure.Conflict, "publication receipt conflicts with remote state", request.CorrelationID, failure.RetryReconcile)
	}
	if err != nil {
		return Receipt{}, false, sourceFailure(failure.DependencyUnavailable, "publication receipt persistence is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	return saved, replayedOperation || replayedReceipt, nil
}

func validRequest(ctx context.Context, request PublishRequest) bool {
	return ctx != nil && request.ContractVersion == PublishRequestContractVersion &&
		request.CorrelationID != "" && request.IdempotencyKey != "" && request.TenantID != "" &&
		request.RunID != "" && request.PolicyDecisionID != "" && request.ApprovalRecordID != "" &&
		request.Provider == "contract-double" && request.BaseBranch != "" && request.BaseRevision != "" &&
		request.TreeDigest != "" && validDigest(request.TreeDigest) && len(request.ChangedPaths) > 0 &&
		request.Title != "" && request.Body != "" && safeMetadata(request.Title) && safeMetadata(request.Body)
}

func requestMatchesSubject(request PublishRequest, subject trust.PublicationApprovalSubject) bool {
	return subject.RunID == request.RunID && subject.Scope.TenantID == request.TenantID &&
		subject.Scope.DeliveryRunID == request.RunID &&
		request.RepositoryBinding.ContractVersion() == tenant.RepositoryBindingContractVersion &&
		request.RepositoryBinding.TenantID() == subject.Scope.TenantID &&
		request.RepositoryBinding.WorkspaceID() == subject.Scope.WorkspaceID &&
		request.RepositoryBinding.ProjectID() == subject.Scope.ProjectID &&
		request.RepositoryBinding.RepositoryID() == subject.Scope.RepositoryID &&
		request.RepositoryBinding.CanonicalLocator() != "" &&
		request.RepositoryBinding.Classification() == subject.Classification &&
		subject.CandidateDigest != "" && validDigest(subject.CandidateDigest) &&
		request.TreeDigest == subject.TreeDigest &&
		slices.Equal(canonical(request.ChangedPaths), canonical(subject.ChangedPaths)) &&
		subject.VerificationResultDigest != "" && len(subject.EvidenceReferences) > 0 &&
		!subject.ExpiresAt.IsZero()
}

func approvalMatchesSubject(record trust.ApprovalRecordV2, subject trust.PublicationApprovalSubject, at time.Time) bool {
	if record.VerifyIntegrity(subject.CorrelationID) != nil {
		return false
	}
	return record.Decision == "approved" && at.Before(record.ExpiresAt) &&
		record.Action == PublishAction && record.SubjectType == "verified-candidate" &&
		record.SubjectID == subject.CandidateID && record.SubjectDigest == subject.CandidateDigest &&
		record.Scope == subject.Scope && record.ApproverID == subject.RequestingActorID &&
		record.EvidenceCorrelationID == subject.CorrelationID &&
		slices.Contains(record.EvidenceReferences, "candidate:"+subject.CandidateDigest) &&
		slices.Contains(record.EvidenceReferences, "verification:"+subject.VerificationResultDigest)
}

func validateEnvelope(envelope trust.ExecutionEnvelope, request PublishRequest, subject trust.PublicationApprovalSubject, approval trust.ApprovalRecordV2, decision trust.PolicyDecision, at time.Time) error {
	if envelope.VerifyIntegrity(request.CorrelationID) != nil ||
		envelope.Scope != (trust.ExecutionScope{TenantID: subject.Scope.TenantID, WorkspaceID: subject.Scope.WorkspaceID, ProjectID: subject.Scope.ProjectID, DeliveryRunID: subject.RunID}) ||
		envelope.Principal.SubjectID == "" || envelope.Principal.AgentID == "" ||
		!at.Before(envelope.ExpiresAt) || envelope.IssuedAt.After(at) ||
		!slices.Equal(envelope.Actions, []string{PublishAction}) ||
		!slices.Equal(envelope.ApprovalRecordIDs, []string{approval.ID}) ||
		!slices.Contains(envelope.PolicyDecisionIDs, decision.ID) ||
		!slices.Equal(envelope.Destinations, []string{"contract-double:" + request.RepositoryBinding.RepositoryID()}) ||
		envelope.IdempotencyKey != request.IdempotencyKey ||
		envelope.EvidenceCorrelationID != request.CorrelationID ||
		envelope.Limits.MaxDurationSeconds < 1 || envelope.Limits.MaxDurationSeconds > 300 ||
		envelope.Limits.MaxConcurrency != 1 || envelope.Limits.MaxCostMinorUnits != 0 ||
		len(envelope.ToolIDs) != 0 || len(envelope.ModelIDs) != 0 || len(envelope.SecretReferenceIDs) != 0 {
		return forbidden(request.CorrelationID)
	}
	return nil
}

func buildOperation(request PublishRequest, subject trust.PublicationApprovalSubject) Operation {
	changedPaths := canonical(request.ChangedPaths)
	identityInput := strings.Join([]string{
		subject.Scope.TenantID, subject.RunID, subject.Scope.RepositoryID, request.BaseRevision,
		subject.CandidateDigest, PublishAction,
	}, "\x00")
	sum := sha256.Sum256([]byte(identityInput))
	id := "publication-" + hex.EncodeToString(sum[:16])
	safeRun := safeSegment(subject.RunID)
	digestPrefix := strings.TrimPrefix(subject.CandidateDigest, "sha256:")
	branch := "forge/" + safeRun + "-" + digestPrefix[:12]
	metadata, _ := json.Marshal(struct{ Title, Body string }{request.Title, request.Body})
	metadataSum := sha256.Sum256(metadata)
	return Operation{
		ContractVersion: ProviderContractVersion, ID: id, TenantID: subject.Scope.TenantID,
		RunID: subject.RunID, RepositoryID: subject.Scope.RepositoryID, Provider: request.Provider,
		BaseBranch: request.BaseBranch, BaseRevision: request.BaseRevision, Branch: branch,
		CandidateID: subject.CandidateID, CandidateDigest: subject.CandidateDigest,
		TreeDigest: request.TreeDigest, ChangedPaths: changedPaths,
		MetadataDigest:   "sha256:" + hex.EncodeToString(metadataSum[:]),
		ApprovalRecordID: request.ApprovalRecordID, EnvelopeID: request.Envelope.ID,
		PolicyDecisionID: request.PolicyDecisionID, Marker: "<!-- forge-publication:" + id + " -->",
	}
}

func reconcileExact(ctx context.Context, provider Provider, operation Operation) (RemoteState, bool, error) {
	state, err := provider.Reconcile(ctx, operation)
	if err != nil {
		return RemoteState{}, false, err
	}
	if state == (RemoteState{}) {
		return state, false, nil
	}
	if state.RepositoryID != operation.RepositoryID || state.Branch != operation.Branch ||
		state.BaseBranch != operation.BaseBranch || state.CandidateDigest != operation.CandidateDigest ||
		state.TreeDigest != operation.TreeDigest || state.Marker != operation.Marker ||
		(state.PullRequestID != "" && (!state.BranchExists || !state.PullRequestOpen)) {
		return RemoteState{}, false, ErrProviderConflict
	}
	return state, true, nil
}

func safeMetadata(value string) bool {
	lower := strings.ToLower(value)
	return len(value) <= 4096 && !strings.Contains(lower, "authorization:") &&
		!strings.Contains(lower, "bearer ") && !strings.Contains(lower, "private key") &&
		!strings.ContainsRune(value, '\x00')
}

var unsafeSegment = regexp.MustCompile(`[^a-z0-9-]+`)

func safeSegment(value string) string {
	value = strings.ToLower(value)
	value = strings.Trim(unsafeSegment.ReplaceAllString(value, "-"), "-")
	if value == "" {
		return "run"
	}
	if len(value) > 32 {
		return value[:32]
	}
	return value
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
		return false
	}
	for _, character := range value[7:] {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

func canonical(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return slices.Compact(result)
}

func forbidden(correlationID string) error {
	return sourceFailure(failure.Forbidden, "publication authority does not match the current verified candidate", correlationID, failure.RetryNone)
}

func normalizeProvider(err error, correlationID string) error {
	switch {
	case errors.Is(err, ErrProviderConflict):
		return sourceFailure(failure.Conflict, "provider state conflicts with the authorized publication", correlationID, failure.RetryReconcile)
	case errors.Is(err, ErrProviderAmbiguous):
		return sourceFailure(failure.ManualRecoveryRequired, "provider outcome remains ambiguous", correlationID, failure.RetryReconcile)
	default:
		return sourceFailure(failure.DependencyUnavailable, "source-control provider is unavailable", correlationID, failure.RetryBounded)
	}
}

func sourceFailure(code failure.Code, message, correlationID string, retry failure.Retry) error {
	return failure.New(code, message, correlationID, retry)
}

type memoryOperation struct {
	fingerprint string
	operation   Operation
}

type MemoryStore struct {
	mu         sync.Mutex
	operations map[string]memoryOperation
	receipts   map[string]Receipt
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{operations: make(map[string]memoryOperation), receipts: make(map[string]Receipt)}
}

func (store *MemoryStore) Create(ctx context.Context, tenantID, key string, operation Operation) (Operation, bool, error) {
	if store == nil || ctx == nil || tenantID == "" || key == "" {
		return Operation{}, false, ErrStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	mapKey := tenantID + "\x00" + key
	fingerprint := operationFingerprint(operation)
	if existing, found := store.operations[mapKey]; found {
		if existing.fingerprint != fingerprint {
			return Operation{}, false, ErrStoreConflict
		}
		return cloneOperation(existing.operation), true, nil
	}
	store.operations[mapKey] = memoryOperation{fingerprint: fingerprint, operation: cloneOperation(operation)}
	return cloneOperation(operation), false, nil
}

func (store *MemoryStore) SaveReceipt(ctx context.Context, operationID string, receipt Receipt) (Receipt, bool, error) {
	if store == nil || ctx == nil || operationID == "" {
		return Receipt{}, false, ErrStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if existing, found := store.receipts[operationID]; found {
		if existing.RepositoryID != receipt.RepositoryID || existing.PullRequestID != receipt.PullRequestID ||
			existing.CandidateDigest != receipt.CandidateDigest || existing.Branch != receipt.Branch {
			return Receipt{}, false, ErrStoreConflict
		}
		return existing, true, nil
	}
	store.receipts[operationID] = receipt
	return receipt, false, nil
}

func operationFingerprint(value Operation) string {
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func cloneOperation(value Operation) Operation {
	value.ChangedPaths = slices.Clone(value.ChangedPaths)
	return value
}
