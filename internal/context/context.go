// Package context assembles bounded information manifests only after exact
// execution authority is verified. It owns no source content or authority.
package context

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

const AssembleRequestContractVersion = "forge.assemble-context-request.v1"

var (
	ErrSourceUnavailable       = errors.New("context source unavailable")
	ErrContextStoreUnavailable = errors.New("context store unavailable")
	ErrContextConflict         = errors.New("context idempotency conflict")
)

type DelegationResolver interface {
	ResolveDelegationChain(context.Context, identity.ResolveDelegationRequest) ([]identity.Delegation, error)
}

type DecisionResolver interface {
	ResolvePolicyDecision(context.Context, string) (trust.PolicyDecision, bool, error)
}

// SourceDiscoveryRequest contains only authority-derived scope plus bounded
// source identifiers. It never contains credentials or unrestricted content.
type SourceDiscoveryRequest struct {
	TenantID           string
	WorkspaceID        string
	ProjectID          string
	RepositoryID       string
	RequestedSourceIDs []string
}

// SourceRecord is provider metadata used to build a manifest. Content bytes
// are deliberately absent from this contract.
type SourceRecord struct {
	TenantID       string
	WorkspaceID    string
	ProjectID      string
	RepositoryID   string
	SourceID       string
	SourceRevision string
	Locator        string
	Classification string
	Authority      string
	Purpose        string
	FreshnessAt    *time.Time
	Confidence     *float64
	ContentDigest  string
}

type SourceCatalog interface {
	Discover(context.Context, SourceDiscoveryRequest) ([]SourceRecord, error)
}

type Store interface {
	Get(context.Context, string, string, string) (trust.ContextPackage, bool, error)
	Create(context.Context, trust.ContextPackage, string, string, string) (trust.ContextPackage, bool, error)
}

type ReproductionStore interface {
	Resolve(context.Context, string, string) (trust.ContextPackage, bool, error)
}

type AssemblyProfile struct {
	MaximumItems            int
	MaximumSerializedBytes  int
	MaximumLifetime         time.Duration
	EnforcedObligationTypes []string
}

type Service struct {
	delegations DelegationResolver
	decisions   DecisionResolver
	catalog     SourceCatalog
	store       Store
	profile     AssemblyProfile
	now         func() time.Time
}

func NewService(delegations DelegationResolver, decisions DecisionResolver, catalog SourceCatalog, store Store, profile AssemblyProfile, now func() time.Time) (*Service, error) {
	if delegations == nil {
		return nil, errors.New("identity delegation resolver is required")
	}
	if decisions == nil {
		return nil, errors.New("policy decision resolver is required")
	}
	if catalog == nil {
		return nil, errors.New("context source catalog is required")
	}
	if store == nil {
		return nil, errors.New("context store is required")
	}
	if profile.MaximumItems < 1 || profile.MaximumItems > 2000 ||
		profile.MaximumSerializedBytes < 1 || profile.MaximumSerializedBytes > 10*1024*1024 || profile.MaximumLifetime <= 0 {
		return nil, errors.New("context assembly profile is outside the M1 validation envelope")
	}
	obligations := slices.Clone(profile.EnforcedObligationTypes)
	sort.Strings(obligations)
	if !uniqueNonEmpty(obligations) {
		return nil, errors.New("context assembly obligation profile is invalid")
	}
	profile.EnforcedObligationTypes = obligations
	if now == nil {
		return nil, errors.New("context assembly clock is required")
	}
	return &Service{delegations: delegations, decisions: decisions, catalog: catalog, store: store, profile: profile, now: now}, nil
}

type AssembleRequest struct {
	ContractVersion    string
	CorrelationID      string
	IdempotencyKey     string
	Human              identity.Principal
	Agent              identity.Principal
	Scope              tenant.Scope
	Run                delivery.DeliveryRun
	DelegationID       string
	Envelope           trust.ExecutionEnvelope
	PolicyDecisionID   string
	RequestedSourceIDs []string
}

type AssembleResult struct {
	Package  trust.ContextPackage
	Replayed bool
}

func (service *Service) AssembleContext(ctx context.Context, request AssembleRequest) (AssembleResult, error) {
	if service == nil || service.delegations == nil || service.decisions == nil || service.catalog == nil || service.store == nil || service.now == nil {
		return AssembleResult{}, unavailable(request.CorrelationID)
	}
	requested, err := canonicalSourceIDs(request.RequestedSourceIDs)
	if ctx == nil || request.ContractVersion != AssembleRequestContractVersion || request.CorrelationID == "" ||
		request.IdempotencyKey == "" || request.DelegationID == "" || request.PolicyDecisionID == "" ||
		err != nil || len(requested) > service.profile.MaximumItems {
		return AssembleResult{}, invalid(request.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return AssembleResult{}, unavailable(request.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return AssembleResult{}, unavailable(request.CorrelationID)
	}
	if err := request.Human.ValidateAt(at, request.CorrelationID); err != nil {
		return AssembleResult{}, err
	}
	if err := request.Agent.ValidateAt(at, request.CorrelationID); err != nil {
		return AssembleResult{}, err
	}
	if err := request.Scope.ValidateFor(request.Human, at, request.CorrelationID); err != nil {
		return AssembleResult{}, err
	}
	if !validRunBinding(request) {
		return AssembleResult{}, forbidden(request.CorrelationID)
	}
	if err := request.Envelope.VerifyIntegrity(request.CorrelationID); err != nil {
		return AssembleResult{}, err
	}
	if err := request.Envelope.AuthorizeAction(identity.DelegatedActionAssembleContext, at, request.CorrelationID); err != nil {
		return AssembleResult{}, err
	}
	chain, err := service.delegations.ResolveDelegationChain(ctx, identity.ResolveDelegationRequest{
		ContractVersion: identity.ResolveDelegationContractVersion,
		CorrelationID:   request.CorrelationID, Human: request.Human, Agent: request.Agent, DelegationID: request.DelegationID,
	})
	if err != nil {
		return AssembleResult{}, err
	}
	if len(chain) != 1 || chain[0].ID() != request.DelegationID || !validAuthorityBinding(request, chain[0]) {
		return AssembleResult{}, forbidden(request.CorrelationID)
	}
	decision, found, err := service.decisions.ResolvePolicyDecision(ctx, request.PolicyDecisionID)
	if err != nil {
		return AssembleResult{}, unavailable(request.CorrelationID)
	}
	if !found || decision.ID != request.PolicyDecisionID || decision.EvidenceCorrelationID != request.CorrelationID ||
		!slices.Contains(request.Envelope.PolicyDecisionIDs, decision.ID) ||
		!slices.Contains(request.Envelope.PolicyDecisionIDs, chain[0].ParentAuthorityID()) {
		return AssembleResult{}, forbidden(request.CorrelationID)
	}
	if err := decision.Authorize(
		"retrieval", request.Human.SubjectID(), identity.DelegatedActionAssembleContext,
		request.Run.RepositoryID(), at, service.profile.EnforcedObligationTypes, request.CorrelationID,
	); err != nil {
		return AssembleResult{}, err
	}
	fingerprint := assemblyFingerprint(request, requested, service.profile)
	existing, found, err := service.store.Get(ctx, request.Scope.TenantID(), request.IdempotencyKey, fingerprint)
	if err != nil {
		return AssembleResult{}, normalizeStoreError(err, request.CorrelationID)
	}
	if found {
		if err := validateStoredPackage(existing, request, requested, service.profile, at); err != nil {
			return AssembleResult{}, err
		}
		return AssembleResult{Package: clonePackage(existing), Replayed: true}, nil
	}

	// No source adapter call occurs before all authority checks above succeed.
	records, err := service.catalog.Discover(ctx, SourceDiscoveryRequest{
		TenantID: request.Scope.TenantID(), WorkspaceID: request.Scope.WorkspaceID(),
		ProjectID: request.Scope.ProjectID(), RepositoryID: request.Run.RepositoryID(), RequestedSourceIDs: requested,
	})
	if err != nil {
		return AssembleResult{}, sourceUnavailable(request.CorrelationID)
	}
	items, omissions, err := validateRecords(records, requested, request)
	if err != nil {
		return AssembleResult{}, err
	}
	metadata, err := json.Marshal(struct {
		Items     []trust.ContextItem
		Omissions []string
	}{items, omissions})
	if err != nil || len(metadata) > service.profile.MaximumSerializedBytes {
		return AssembleResult{}, failure.New(failure.BudgetExhausted, "context manifest exceeds its metadata budget", request.CorrelationID, failure.RetryNone)
	}
	expiresAt := request.Envelope.ExpiresAt
	if profileExpiry := at.Add(service.profile.MaximumLifetime); profileExpiry.Before(expiresAt) {
		expiresAt = profileExpiry
	}
	value := trust.ContextPackage{
		SchemaVersion: trust.ContextPackageSchemaVersion, ExecutionEnvelopeID: request.Envelope.ID,
		Objective: request.Run.Objective(), AssembledAt: at, ExpiresAt: &expiresAt,
		PolicyDecisionIDs: slices.Clone(request.Envelope.PolicyDecisionIDs), Items: items,
		Omissions: omissions, Conflicts: []string{},
	}
	value, err = trust.BindContextPackageManifest(value, request.CorrelationID)
	if err != nil {
		return AssembleResult{}, err
	}
	stored, replayed, err := service.store.Create(ctx, value, request.Scope.TenantID(), request.IdempotencyKey, fingerprint)
	if err != nil {
		return AssembleResult{}, normalizeStoreError(err, request.CorrelationID)
	}
	if !replayed && (stored.ID != value.ID || stored.ManifestDigest != value.ManifestDigest) {
		return AssembleResult{}, failure.New(failure.Indeterminate, "stored context manifest differs from the assembled result", request.CorrelationID, failure.RetryNone)
	}
	if err := validateStoredPackage(stored, request, requested, service.profile, at); err != nil {
		return AssembleResult{}, err
	}
	return AssembleResult{Package: clonePackage(stored), Replayed: replayed}, nil
}

type memoryEntry struct {
	value       trust.ContextPackage
	fingerprint string
}

type MemoryStore struct {
	mu            sync.RWMutex
	byIdempotency map[string]memoryEntry
	byID          map[string]trust.ContextPackage
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		byIdempotency: make(map[string]memoryEntry),
		byID:          make(map[string]trust.ContextPackage),
	}
}

func (store *MemoryStore) Get(_ context.Context, tenantID, idempotencyKey, fingerprint string) (trust.ContextPackage, bool, error) {
	if store == nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	existing, found := store.byIdempotency[tenantID+"\x00"+idempotencyKey]
	if !found {
		return trust.ContextPackage{}, false, nil
	}
	if existing.fingerprint != fingerprint {
		return trust.ContextPackage{}, false, ErrContextConflict
	}
	return clonePackage(existing.value), true, nil
}

func (store *MemoryStore) Create(_ context.Context, value trust.ContextPackage, tenantID, idempotencyKey, fingerprint string) (trust.ContextPackage, bool, error) {
	if store == nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	key := tenantID + "\x00" + idempotencyKey
	if existing, found := store.byIdempotency[key]; found {
		if existing.fingerprint != fingerprint {
			return trust.ContextPackage{}, false, ErrContextConflict
		}
		return clonePackage(existing.value), true, nil
	}
	store.byIdempotency[key] = memoryEntry{value: clonePackage(value), fingerprint: fingerprint}
	store.byID[tenantID+"\x00"+value.ID] = clonePackage(value)
	return clonePackage(value), false, nil
}

func (store *MemoryStore) Resolve(ctx context.Context, tenantID, packageID string) (trust.ContextPackage, bool, error) {
	if store == nil || store.byID == nil || ctx == nil || tenantID == "" || packageID == "" {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	if err := ctx.Err(); err != nil {
		return trust.ContextPackage{}, false, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	value, found := store.byID[tenantID+"\x00"+packageID]
	return clonePackage(value), found, nil
}

type ReproduceRequest struct {
	ContractVersion        string
	CorrelationID          string
	Human                  identity.Principal
	Scope                  tenant.Scope
	Run                    delivery.DeliveryRun
	PackageID              string
	ExpectedManifestDigest string
}

const ReproduceRequestContractVersion = "forge.reproduce-context-request.v1"

type ReproductionResult struct {
	PackageID        string
	ManifestDigest   string
	Exact            bool
	SourceRevisions  []string
	Omissions        []string
	ChangedSourceIDs []string
	StaleSourceIDs   []string
	CheckedAt        time.Time
}

// ReproduceContext compares current authorized source metadata with the
// capability-owned manifest. It does not retrieve or return source content.
func (service *Service) ReproduceContext(ctx context.Context, request ReproduceRequest) (ReproductionResult, error) {
	if service == nil || service.catalog == nil || service.store == nil || service.now == nil {
		return ReproductionResult{}, unavailable(request.CorrelationID)
	}
	store, ok := service.store.(ReproductionStore)
	if !ok || ctx == nil || request.ContractVersion != ReproduceRequestContractVersion ||
		request.CorrelationID == "" || request.PackageID == "" || request.ExpectedManifestDigest == "" {
		return ReproductionResult{}, invalid(request.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return ReproductionResult{}, unavailable(request.CorrelationID)
	}
	if err := request.Human.ValidateAt(at, request.CorrelationID); err != nil {
		return ReproductionResult{}, err
	}
	if err := request.Scope.ValidateFor(request.Human, at, request.CorrelationID); err != nil {
		return ReproductionResult{}, err
	}
	contextRecord := request.Run.Context()
	if request.Run.ActorID() != request.Human.SubjectID() ||
		request.Run.TenantID() != request.Scope.TenantID() ||
		request.Run.WorkspaceID() != request.Scope.WorkspaceID() ||
		request.Run.ProjectID() != request.Scope.ProjectID() ||
		contextRecord == nil || contextRecord.PackageID != request.PackageID ||
		contextRecord.ManifestDigest != request.ExpectedManifestDigest {
		return ReproductionResult{}, forbidden(request.CorrelationID)
	}
	value, found, err := store.Resolve(ctx, request.Scope.TenantID(), request.PackageID)
	if err != nil {
		return ReproductionResult{}, unavailable(request.CorrelationID)
	}
	if !found || value.ID != request.PackageID || value.ManifestDigest != request.ExpectedManifestDigest {
		return ReproductionResult{}, failure.New(failure.Indeterminate, "context manifest could not be resolved", request.CorrelationID, failure.RetryNone)
	}
	if err := value.VerifyManifest(request.CorrelationID); err != nil {
		return ReproductionResult{}, err
	}
	if value.ExpiresAt == nil || !at.Before(*value.ExpiresAt) {
		return ReproductionResult{}, failure.New(failure.Expired, "context manifest is no longer current", request.CorrelationID, failure.RetryNewAuthority)
	}
	sourceIDs := make([]string, 0, len(value.Items))
	expected := make(map[string]trust.ContextItem, len(value.Items))
	revisions := make([]string, 0, len(value.Items))
	for _, item := range value.Items {
		sourceIDs = append(sourceIDs, item.SourceID)
		expected[item.SourceID] = item
		revisions = append(revisions, item.SourceID+"@"+item.SourceRevision)
	}
	sort.Strings(sourceIDs)
	sort.Strings(revisions)
	records, err := service.catalog.Discover(ctx, SourceDiscoveryRequest{
		TenantID: request.Scope.TenantID(), WorkspaceID: request.Scope.WorkspaceID(),
		ProjectID: request.Scope.ProjectID(), RepositoryID: request.Run.RepositoryID(),
		RequestedSourceIDs: sourceIDs,
	})
	if err != nil {
		return ReproductionResult{}, sourceUnavailable(request.CorrelationID)
	}
	seen := make(map[string]struct{}, len(records))
	changed := []string{}
	stale := []string{}
	for _, record := range records {
		item, expectedSource := expected[record.SourceID]
		if !expectedSource || record.TenantID != request.Scope.TenantID() ||
			record.WorkspaceID != request.Scope.WorkspaceID() || record.ProjectID != request.Scope.ProjectID() ||
			record.RepositoryID != request.Run.RepositoryID() {
			return ReproductionResult{}, failure.New(failure.Indeterminate, "context reproduction returned unauthorized metadata", request.CorrelationID, failure.RetryNone)
		}
		if _, duplicate := seen[record.SourceID]; duplicate {
			return ReproductionResult{}, failure.New(failure.Indeterminate, "context reproduction returned duplicate metadata", request.CorrelationID, failure.RetryNone)
		}
		seen[record.SourceID] = struct{}{}
		if record.SourceRevision != item.SourceRevision || record.ContentDigest != item.ContentDigest ||
			record.Locator != item.Locator || record.Classification != item.Classification ||
			record.Authority != item.Authority || record.Purpose != item.Purpose ||
			record.SourceRevision == "" || record.ContentDigest == "" || record.Locator == "" {
			changed = append(changed, record.SourceID)
		}
		if item.FreshnessAt != nil && (record.FreshnessAt == nil || record.FreshnessAt.Before(*item.FreshnessAt)) {
			stale = append(stale, record.SourceID)
		}
	}
	omissions := slices.Clone(value.Omissions)
	for _, sourceID := range sourceIDs {
		if _, found := seen[sourceID]; !found {
			omissions = append(omissions, "reproduction-source-unavailable:"+sourceID)
		}
	}
	sort.Strings(omissions)
	sort.Strings(changed)
	sort.Strings(stale)
	return ReproductionResult{
		PackageID: value.ID, ManifestDigest: value.ManifestDigest,
		Exact:           len(omissions) == len(value.Omissions) && len(changed) == 0 && len(stale) == 0,
		SourceRevisions: revisions, Omissions: omissions, ChangedSourceIDs: changed,
		StaleSourceIDs: stale, CheckedAt: at,
	}, nil
}

func validRunBinding(request AssembleRequest) bool {
	return request.Human.Kind() == identity.HumanPrincipal && request.Human.Synthetic() &&
		request.Agent.Kind() == identity.AgentPrincipal && request.Agent.Synthetic() &&
		request.Human.SubjectID() != request.Agent.SubjectID() &&
		request.Run.ContractVersion() == delivery.DeliveryRunContractVersion && request.Run.State() == delivery.CapturedState &&
		request.Run.ID() != "" && request.Run.RepositoryID() != "" && request.Run.Objective() != "" &&
		request.Run.TenantID() == request.Scope.TenantID() && request.Run.WorkspaceID() == request.Scope.WorkspaceID() &&
		request.Run.ProjectID() == request.Scope.ProjectID() && request.Run.ActorID() == request.Human.SubjectID()
}

func validAuthorityBinding(request AssembleRequest, delegation identity.Delegation) bool {
	return request.Envelope.Scope == (trust.ExecutionScope{
		TenantID: request.Scope.TenantID(), WorkspaceID: request.Scope.WorkspaceID(),
		ProjectID: request.Scope.ProjectID(), DeliveryRunID: request.Run.ID(),
	}) && request.Envelope.Principal.SubjectID == request.Human.SubjectID() && request.Envelope.Principal.AgentID == request.Agent.SubjectID() &&
		slices.Equal(request.Envelope.Principal.DelegationChain, []string{delegation.ID()}) &&
		request.Envelope.EvidenceCorrelationID == request.CorrelationID && request.Envelope.Purpose == delegation.Purpose() &&
		request.Envelope.RiskClass == "low" && slices.Equal(request.Envelope.Actions, delegation.Actions()) &&
		slices.Equal(request.Envelope.Classifications, delegation.Classifications()) && len(request.Envelope.Destinations) == 0 &&
		len(request.Envelope.ModelIDs) == 0 && len(request.Envelope.ToolIDs) == 0 && len(request.Envelope.SecretReferenceIDs) == 0 &&
		delegation.HumanSubjectID() == request.Human.SubjectID() && delegation.AgentSubjectID() == request.Agent.SubjectID() &&
		delegation.TenantID() == request.Scope.TenantID() && delegation.WorkspaceID() == request.Scope.WorkspaceID() &&
		delegation.ProjectID() == request.Scope.ProjectID() && delegation.RepositoryID() == request.Run.RepositoryID() &&
		delegation.DeliveryRunID() == request.Run.ID() && slices.Equal(delegation.Resources(), []string{request.Run.RepositoryID()}) &&
		slices.Equal(delegation.Classifications(), []string{request.Run.Classification()})
}

func validateRecords(records []SourceRecord, requested []string, request AssembleRequest) ([]trust.ContextItem, []string, error) {
	requestedSet := make(map[string]struct{}, len(requested))
	for _, id := range requested {
		requestedSet[id] = struct{}{}
	}
	seen := make(map[string]struct{}, len(records))
	items := make([]trust.ContextItem, 0, len(records))
	for _, record := range records {
		if _, requestedSource := requestedSet[record.SourceID]; !requestedSource {
			return nil, nil, failure.New(failure.Indeterminate, "context source returned unauthorized metadata", request.CorrelationID, failure.RetryNone)
		}
		if _, duplicate := seen[record.SourceID]; duplicate {
			return nil, nil, failure.New(failure.Indeterminate, "context source returned duplicate metadata", request.CorrelationID, failure.RetryNone)
		}
		seen[record.SourceID] = struct{}{}
		if record.TenantID != request.Scope.TenantID() || record.WorkspaceID != request.Scope.WorkspaceID() ||
			record.ProjectID != request.Scope.ProjectID() || record.RepositoryID != request.Run.RepositoryID() ||
			record.Classification != request.Run.Classification() || record.SourceRevision == "" || record.Locator == "" ||
			record.Purpose == "" || record.ContentDigest == "" || !validAuthority(record.Authority) ||
			(record.FreshnessAt != nil && record.FreshnessAt.IsZero()) ||
			(record.Confidence != nil && (*record.Confidence < 0 || *record.Confidence > 1)) {
			return nil, nil, failure.New(failure.Indeterminate, "context source metadata could not be verified", request.CorrelationID, failure.RetryNone)
		}
		items = append(items, trust.ContextItem{
			SourceID: record.SourceID, SourceRevision: record.SourceRevision, Locator: record.Locator,
			Classification: record.Classification, Authority: record.Authority, Purpose: record.Purpose,
			FreshnessAt: cloneTime(record.FreshnessAt), Confidence: cloneFloat(record.Confidence), ContentDigest: record.ContentDigest,
		})
	}
	sort.Slice(items, func(left, right int) bool {
		if items[left].SourceID == items[right].SourceID {
			return items[left].Locator < items[right].Locator
		}
		return items[left].SourceID < items[right].SourceID
	})
	omissions := make([]string, 0, len(requested)-len(items))
	for _, id := range requested {
		if _, found := seen[id]; !found {
			omissions = append(omissions, "requested-source-unavailable:"+id)
		}
	}
	return items, omissions, nil
}

func validateStoredPackage(value trust.ContextPackage, request AssembleRequest, requested []string, profile AssemblyProfile, at time.Time) error {
	if err := value.VerifyManifest(request.CorrelationID); err != nil {
		return err
	}
	if err := value.ValidateFor(request.Envelope, at, request.CorrelationID); err != nil {
		return err
	}
	if value.ExecutionEnvelopeID != request.Envelope.ID || value.Objective != request.Run.Objective() ||
		!slices.Equal(value.PolicyDecisionIDs, request.Envelope.PolicyDecisionIDs) || len(value.Conflicts) != 0 ||
		len(value.Items) > profile.MaximumItems || value.ExpiresAt == nil ||
		value.ExpiresAt.After(value.AssembledAt.Add(profile.MaximumLifetime)) {
		return failure.New(failure.Indeterminate, "stored context manifest is outside the authorized request", request.CorrelationID, failure.RetryNone)
	}
	requestedSet := make(map[string]struct{}, len(requested))
	for _, id := range requested {
		requestedSet[id] = struct{}{}
	}
	accounted := make(map[string]struct{}, len(requested))
	for _, item := range value.Items {
		if _, allowed := requestedSet[item.SourceID]; !allowed || item.Classification != request.Run.Classification() {
			return failure.New(failure.Indeterminate, "stored context manifest contains unauthorized metadata", request.CorrelationID, failure.RetryNone)
		}
		if _, duplicate := accounted[item.SourceID]; duplicate {
			return failure.New(failure.Indeterminate, "stored context manifest contains duplicate metadata", request.CorrelationID, failure.RetryNone)
		}
		accounted[item.SourceID] = struct{}{}
	}
	const omissionPrefix = "requested-source-unavailable:"
	for _, omission := range value.Omissions {
		if len(omission) <= len(omissionPrefix) || omission[:len(omissionPrefix)] != omissionPrefix {
			return failure.New(failure.Indeterminate, "stored context manifest contains an invalid omission", request.CorrelationID, failure.RetryNone)
		}
		id := omission[len(omissionPrefix):]
		if _, allowed := requestedSet[id]; !allowed {
			return failure.New(failure.Indeterminate, "stored context manifest contains an unauthorized omission", request.CorrelationID, failure.RetryNone)
		}
		if _, duplicate := accounted[id]; duplicate {
			return failure.New(failure.Indeterminate, "stored context manifest contains duplicate source outcomes", request.CorrelationID, failure.RetryNone)
		}
		accounted[id] = struct{}{}
	}
	if len(accounted) != len(requestedSet) {
		return failure.New(failure.Indeterminate, "stored context manifest omits a requested source outcome", request.CorrelationID, failure.RetryNone)
	}
	metadata, err := json.Marshal(struct {
		Items     []trust.ContextItem
		Omissions []string
	}{value.Items, value.Omissions})
	if err != nil || len(metadata) > profile.MaximumSerializedBytes {
		return failure.New(failure.Indeterminate, "stored context manifest exceeds its authorized budget", request.CorrelationID, failure.RetryNone)
	}
	return nil
}

func canonicalSourceIDs(values []string) ([]string, error) {
	if len(values) == 0 {
		return nil, errors.New("requested sources are required")
	}
	result := slices.Clone(values)
	sort.Strings(result)
	for index, value := range result {
		if value == "" || (index > 0 && value == result[index-1]) {
			return nil, errors.New("requested sources must be unique and non-empty")
		}
	}
	return result, nil
}

func assemblyFingerprint(request AssembleRequest, requested []string, profile AssemblyProfile) string {
	encoded, _ := json.Marshal(struct {
		EnvelopeID, DecisionID, DelegationID, Objective string
		RequestedSourceIDs, EnforcedObligationTypes     []string
		MaximumItems, MaximumSerializedBytes            int
		MaximumLifetimeNanoseconds                      int64
	}{
		request.Envelope.ID, request.PolicyDecisionID, request.DelegationID, request.Run.Objective(),
		requested, profile.EnforcedObligationTypes, profile.MaximumItems, profile.MaximumSerializedBytes,
		int64(profile.MaximumLifetime),
	})
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func clonePackage(value trust.ContextPackage) trust.ContextPackage {
	value.PolicyDecisionIDs = slices.Clone(value.PolicyDecisionIDs)
	value.Omissions = slices.Clone(value.Omissions)
	value.Conflicts = slices.Clone(value.Conflicts)
	items := value.Items
	value.Items = make([]trust.ContextItem, len(items))
	for index, item := range items {
		item.FreshnessAt = cloneTime(item.FreshnessAt)
		item.Confidence = cloneFloat(item.Confidence)
		value.Items[index] = item
	}
	if value.ExpiresAt != nil {
		expiresAt := *value.ExpiresAt
		value.ExpiresAt = &expiresAt
	}
	return value
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func cloneFloat(value *float64) *float64 {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func validAuthority(value string) bool {
	return value == "none" || value == "reference" || value == "accepted-knowledge" || value == "task-instruction"
}

func uniqueNonEmpty(values []string) bool {
	for index, value := range values {
		if value == "" || (index > 0 && value == values[index-1]) {
			return false
		}
	}
	return true
}

func normalizeStoreError(err error, correlationID string) error {
	switch {
	case errors.Is(err, ErrContextConflict):
		return failure.New(failure.Conflict, "context request conflicts with an existing manifest", correlationID, failure.RetryNone)
	case errors.Is(err, ErrContextStoreUnavailable):
		return unavailable(correlationID)
	default:
		return failure.New(failure.Indeterminate, "context manifest persistence could not be verified", correlationID, failure.RetryNone)
	}
}

func invalid(correlationID string) error {
	return failure.New(failure.InvalidRequest, "context assembly request is invalid", correlationID, failure.RetryNone)
}

func forbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "context assembly is not authorized", correlationID, failure.RetryNone)
}

func unavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "context assembly is unavailable", correlationID, failure.RetryBounded)
}

func sourceUnavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "authorized context source is unavailable", correlationID, failure.RetryBounded)
}
