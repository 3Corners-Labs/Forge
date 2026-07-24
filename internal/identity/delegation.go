package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
)

const (
	DelegationRequestContractVersion = "forge.delegation-request.v1"
	CreateDelegationContractVersion  = "forge.create-delegation-request.v1"
	RevokeDelegationContractVersion  = "forge.revoke-delegation-request.v1"
	ResolveDelegationContractVersion = "forge.resolve-delegation-request.v1"
	DelegationContractVersion        = "forge.delegation.v1"

	DelegationActionCreate         = "forge.delegation.create"
	DelegatedActionAssembleContext = "forge.context.assemble"
	DelegatedActionProposeChange   = "forge.intelligence.propose-change"
	DeveloperTrialEnvironment      = "developer-trial"
	ContextPreparationOnly         = "context-preparation-only"
	ChangeProposalOnly             = "change-proposal-only"
	M1ChangePurpose                = "forge.delivery.m1-change"
	delegationRevocationSource     = "forge.identity.delegation-store"
)

var (
	ErrDelegationStoreUnavailable = errors.New("delegation store unavailable")
	ErrDelegationConflict         = errors.New("delegation conflict")
	ErrDelegationRevoked          = errors.New("delegation revoked")
)

// DelegationRequestInput contains requested authority, not granted authority.
// Identity canonicalizes it before Policy evaluates its digest.
type DelegationRequestInput struct {
	ContractVersion    string
	IdempotencyKey     string
	CorrelationID      string
	HumanSubjectID     string
	AgentSubjectID     string
	TenantID           string
	WorkspaceID        string
	ProjectID          string
	RepositoryID       string
	DeliveryRunID      string
	Purpose            string
	Actions            []string
	Resources          []string
	Classifications    []string
	Destinations       []string
	Environment        string
	MaximumAutonomy    string
	RequestedExpiresAt time.Time
	ParentDelegationID string
}

// DelegationRequest is immutable requested authority. It grants nothing until
// an exact permit Policy Decision is verified and Identity mints a Delegation.
type DelegationRequest struct {
	contractVersion    string
	idempotencyKey     string
	correlationID      string
	humanSubjectID     string
	agentSubjectID     string
	tenantID           string
	workspaceID        string
	projectID          string
	repositoryID       string
	deliveryRunID      string
	purpose            string
	actions            []string
	resources          []string
	classifications    []string
	destinations       []string
	environment        string
	maximumAutonomy    string
	requestedExpiresAt time.Time
	parentDelegationID string
	digest             string
}

func NewDelegationRequest(input DelegationRequestInput) (DelegationRequest, error) {
	actions, err := canonicalRequired(input.Actions)
	if err != nil {
		return DelegationRequest{}, errors.New("delegation actions are invalid")
	}
	resources, err := canonicalRequired(input.Resources)
	if err != nil {
		return DelegationRequest{}, errors.New("delegation resources are invalid")
	}
	classifications, err := canonicalRequired(input.Classifications)
	if err != nil {
		return DelegationRequest{}, errors.New("delegation classifications are invalid")
	}
	destinations, err := canonicalOptional(input.Destinations)
	if err != nil {
		return DelegationRequest{}, errors.New("delegation destinations are invalid")
	}
	if input.ContractVersion != DelegationRequestContractVersion || input.IdempotencyKey == "" || input.CorrelationID == "" ||
		input.HumanSubjectID == "" || input.AgentSubjectID == "" || input.HumanSubjectID == input.AgentSubjectID ||
		input.TenantID == "" || input.WorkspaceID == "" || input.ProjectID == "" || input.RepositoryID == "" ||
		input.DeliveryRunID == "" || input.Purpose != M1ChangePurpose || input.Environment != DeveloperTrialEnvironment ||
		input.RequestedExpiresAt.IsZero() || input.ParentDelegationID != "" ||
		!validRootDelegationShape(actions, resources, destinations, input.RepositoryID, input.MaximumAutonomy) ||
		len(classifications) != 1 || !validDelegationClassification(classifications[0]) {
		return DelegationRequest{}, errors.New("delegation request is outside the M1 root contract")
	}
	value := DelegationRequest{
		contractVersion: input.ContractVersion, idempotencyKey: input.IdempotencyKey,
		correlationID: input.CorrelationID, humanSubjectID: input.HumanSubjectID,
		agentSubjectID: input.AgentSubjectID, tenantID: input.TenantID,
		workspaceID: input.WorkspaceID, projectID: input.ProjectID,
		repositoryID: input.RepositoryID, deliveryRunID: input.DeliveryRunID,
		purpose: input.Purpose, actions: actions, resources: resources,
		classifications: classifications, destinations: destinations,
		environment: input.Environment, maximumAutonomy: input.MaximumAutonomy,
		requestedExpiresAt: input.RequestedExpiresAt,
	}
	value.digest = delegationRequestDigest(value)
	return value, nil
}

func (value DelegationRequest) ContractVersion() string       { return value.contractVersion }
func (value DelegationRequest) IdempotencyKey() string        { return value.idempotencyKey }
func (value DelegationRequest) CorrelationID() string         { return value.correlationID }
func (value DelegationRequest) HumanSubjectID() string        { return value.humanSubjectID }
func (value DelegationRequest) AgentSubjectID() string        { return value.agentSubjectID }
func (value DelegationRequest) TenantID() string              { return value.tenantID }
func (value DelegationRequest) WorkspaceID() string           { return value.workspaceID }
func (value DelegationRequest) ProjectID() string             { return value.projectID }
func (value DelegationRequest) RepositoryID() string          { return value.repositoryID }
func (value DelegationRequest) DeliveryRunID() string         { return value.deliveryRunID }
func (value DelegationRequest) Purpose() string               { return value.purpose }
func (value DelegationRequest) Actions() []string             { return slices.Clone(value.actions) }
func (value DelegationRequest) Resources() []string           { return slices.Clone(value.resources) }
func (value DelegationRequest) Classifications() []string     { return slices.Clone(value.classifications) }
func (value DelegationRequest) Destinations() []string        { return slices.Clone(value.destinations) }
func (value DelegationRequest) Environment() string           { return value.environment }
func (value DelegationRequest) MaximumAutonomy() string       { return value.maximumAutonomy }
func (value DelegationRequest) RequestedExpiresAt() time.Time { return value.requestedExpiresAt }
func (value DelegationRequest) ParentDelegationID() string    { return value.parentDelegationID }
func (value DelegationRequest) Digest() string                { return value.digest }

func (value DelegationRequest) ValidateAt(at time.Time, correlationID string) error {
	if at.IsZero() {
		return delegationInvalid(correlationID)
	}
	rebuilt, err := NewDelegationRequest(DelegationRequestInput{
		ContractVersion: value.contractVersion, IdempotencyKey: value.idempotencyKey,
		CorrelationID: value.correlationID, HumanSubjectID: value.humanSubjectID,
		AgentSubjectID: value.agentSubjectID, TenantID: value.tenantID,
		WorkspaceID: value.workspaceID, ProjectID: value.projectID,
		RepositoryID: value.repositoryID, DeliveryRunID: value.deliveryRunID,
		Purpose: value.purpose, Actions: value.actions, Resources: value.resources,
		Classifications: value.classifications, Destinations: value.destinations,
		Environment: value.environment, MaximumAutonomy: value.maximumAutonomy,
		RequestedExpiresAt: value.requestedExpiresAt, ParentDelegationID: value.parentDelegationID,
	})
	if err != nil || rebuilt.digest != value.digest {
		return delegationInvalid(correlationID)
	}
	if !value.requestedExpiresAt.After(at) {
		return failure.New(failure.Expired, "delegation request is not active", correlationID, failure.RetryNewAuthority)
	}
	return nil
}

// Delegation is immutable Identity-minted authority derived from one exact
// permit decision. It is not an Execution Envelope.
type Delegation struct {
	contractVersion       string
	id                    string
	parentAuthorityID     string
	issuerSubjectID       string
	humanSubjectID        string
	agentSubjectID        string
	tenantID              string
	workspaceID           string
	projectID             string
	repositoryID          string
	deliveryRunID         string
	purpose               string
	actions               []string
	resources             []string
	classifications       []string
	destinations          []string
	environment           string
	maximumAutonomy       string
	issuedAt              time.Time
	expiresAt             time.Time
	revocationSource      string
	evidenceCorrelationID string
	requestDigest         string
}

func (value Delegation) ContractVersion() string       { return value.contractVersion }
func (value Delegation) ID() string                    { return value.id }
func (value Delegation) ParentAuthorityID() string     { return value.parentAuthorityID }
func (value Delegation) IssuerSubjectID() string       { return value.issuerSubjectID }
func (value Delegation) HumanSubjectID() string        { return value.humanSubjectID }
func (value Delegation) AgentSubjectID() string        { return value.agentSubjectID }
func (value Delegation) TenantID() string              { return value.tenantID }
func (value Delegation) WorkspaceID() string           { return value.workspaceID }
func (value Delegation) ProjectID() string             { return value.projectID }
func (value Delegation) RepositoryID() string          { return value.repositoryID }
func (value Delegation) DeliveryRunID() string         { return value.deliveryRunID }
func (value Delegation) Purpose() string               { return value.purpose }
func (value Delegation) Actions() []string             { return slices.Clone(value.actions) }
func (value Delegation) Resources() []string           { return slices.Clone(value.resources) }
func (value Delegation) Classifications() []string     { return slices.Clone(value.classifications) }
func (value Delegation) Destinations() []string        { return slices.Clone(value.destinations) }
func (value Delegation) Environment() string           { return value.environment }
func (value Delegation) MaximumAutonomy() string       { return value.maximumAutonomy }
func (value Delegation) IssuedAt() time.Time           { return value.issuedAt }
func (value Delegation) ExpiresAt() time.Time          { return value.expiresAt }
func (value Delegation) RevocationSource() string      { return value.revocationSource }
func (value Delegation) EvidenceCorrelationID() string { return value.evidenceCorrelationID }
func (value Delegation) RequestDigest() string         { return value.requestDigest }

func (value Delegation) ValidateFor(human, agent Principal, at time.Time, correlationID string) error {
	if err := human.ValidateAt(at, correlationID); err != nil {
		return err
	}
	if err := agent.ValidateAt(at, correlationID); err != nil {
		return err
	}
	if value.contractVersion != DelegationContractVersion || value.id == "" || value.parentAuthorityID == "" ||
		value.issuerSubjectID == "" || value.humanSubjectID != human.SubjectID() || value.agentSubjectID != agent.SubjectID() ||
		human.Kind() != HumanPrincipal || agent.Kind() != AgentPrincipal || !human.Synthetic() || !agent.Synthetic() ||
		value.issuerSubjectID != human.SubjectID() ||
		value.tenantID == "" || value.workspaceID == "" || value.projectID == "" || value.repositoryID == "" ||
		value.deliveryRunID == "" || value.purpose != M1ChangePurpose || value.environment != DeveloperTrialEnvironment ||
		value.revocationSource != delegationRevocationSource ||
		value.evidenceCorrelationID == "" || value.requestDigest == "" || value.issuedAt.IsZero() ||
		!value.expiresAt.After(value.issuedAt) ||
		!validRootDelegationShape(value.actions, value.resources, value.destinations, value.repositoryID, value.maximumAutonomy) ||
		len(value.classifications) != 1 || !validDelegationClassification(value.classifications[0]) {
		return failure.New(failure.Forbidden, "task delegation is invalid", correlationID, failure.RetryNewAuthority)
	}
	if at.Before(value.issuedAt) || !at.Before(value.expiresAt) {
		return failure.New(failure.Expired, "task delegation is not active", correlationID, failure.RetryNewAuthority)
	}
	return nil
}

type StoredDelegation struct {
	Delegation Delegation
	Revoked    bool
	RevokedAt  time.Time
}

type DelegationStore interface {
	Create(context.Context, Delegation, string, string) (Delegation, bool, error)
	Get(context.Context, string) (StoredDelegation, bool, error)
	Revoke(context.Context, string, string, time.Time) (StoredDelegation, bool, error)
}

type DelegationService struct {
	store       DelegationStore
	maxLifetime time.Duration
	now         func() time.Time
}

func NewDelegationService(store DelegationStore, maxLifetime time.Duration, now func() time.Time) (*DelegationService, error) {
	if store == nil {
		return nil, errors.New("delegation store is required")
	}
	if maxLifetime <= 0 {
		return nil, errors.New("delegation maximum lifetime must be positive")
	}
	if now == nil {
		return nil, errors.New("delegation clock is required")
	}
	return &DelegationService{store: store, maxLifetime: maxLifetime, now: now}, nil
}

type CreateDelegationRequest struct {
	ContractVersion         string
	Request                 DelegationRequest
	Human                   Principal
	Agent                   Principal
	ParentDecision          trust.PolicyDecision
	EnforcedObligationTypes []string
}

type CreateDelegationResult struct {
	Delegation Delegation
	Replayed   bool
}

func (service *DelegationService) CreateDelegation(ctx context.Context, command CreateDelegationRequest) (CreateDelegationResult, error) {
	correlationID := command.Request.CorrelationID()
	if service == nil || service.store == nil || service.now == nil {
		return CreateDelegationResult{}, delegationUnavailable(correlationID)
	}
	if ctx == nil || command.ContractVersion != CreateDelegationContractVersion {
		return CreateDelegationResult{}, delegationInvalid(correlationID)
	}
	if err := ctx.Err(); err != nil {
		return CreateDelegationResult{}, delegationUnavailable(correlationID)
	}
	at := service.now()
	if at.IsZero() {
		return CreateDelegationResult{}, delegationInvalid(correlationID)
	}
	if err := command.Request.ValidateAt(at, correlationID); err != nil {
		return CreateDelegationResult{}, err
	}
	if err := validateDelegationPrincipals(command.Request, command.Human, command.Agent, at, correlationID); err != nil {
		return CreateDelegationResult{}, err
	}
	if command.ParentDecision.ExpiresAt == nil {
		return CreateDelegationResult{}, failure.New(failure.Forbidden, "delegation parent authority is invalid", correlationID, failure.RetryNewAuthority)
	}
	if command.ParentDecision.EvidenceCorrelationID != correlationID {
		return CreateDelegationResult{}, failure.New(failure.Forbidden, "delegation parent authority is invalid", correlationID, failure.RetryNewAuthority)
	}
	if err := command.ParentDecision.Authorize(
		"execution", command.Human.SubjectID(), DelegationActionCreate, command.Request.Digest(), at,
		command.EnforcedObligationTypes, correlationID,
	); err != nil {
		return CreateDelegationResult{}, err
	}
	expiresAt := earliest(
		command.Request.RequestedExpiresAt(), command.Human.ExpiresAt(), command.Agent.ExpiresAt(),
		*command.ParentDecision.ExpiresAt, at.Add(service.maxLifetime),
	)
	delegation := mintDelegation(command.Request, command.Human, command.ParentDecision.ID, at, expiresAt)
	stored, replayed, err := service.store.Create(ctx, delegation, command.Request.IdempotencyKey(), command.Request.Digest())
	if err != nil {
		return CreateDelegationResult{}, normalizeDelegationStoreError(err, correlationID)
	}
	if stored.RequestDigest() != command.Request.Digest() {
		return CreateDelegationResult{}, failure.New(failure.Indeterminate, "delegation result could not be verified", correlationID, failure.RetryNone)
	}
	if err := stored.ValidateFor(command.Human, command.Agent, at, correlationID); err != nil {
		return CreateDelegationResult{}, err
	}
	return CreateDelegationResult{Delegation: cloneDelegation(stored), Replayed: replayed}, nil
}

type RevokeDelegationRequest struct {
	ContractVersion string
	CorrelationID   string
	Human           Principal
	DelegationID    string
	Reason          string
}

func (service *DelegationService) RevokeDelegation(ctx context.Context, command RevokeDelegationRequest) error {
	if service == nil || service.store == nil || service.now == nil {
		return delegationUnavailable(command.CorrelationID)
	}
	if ctx == nil || command.ContractVersion != RevokeDelegationContractVersion || command.CorrelationID == "" || command.DelegationID == "" || command.Reason == "" {
		return delegationInvalid(command.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return delegationUnavailable(command.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return delegationUnavailable(command.CorrelationID)
	}
	if err := command.Human.ValidateAt(at, command.CorrelationID); err != nil {
		return err
	}
	current, found, err := service.store.Get(ctx, command.DelegationID)
	if err != nil {
		return normalizeDelegationStoreError(err, command.CorrelationID)
	}
	if !found || current.Delegation.HumanSubjectID() != command.Human.SubjectID() || command.Human.Kind() != HumanPrincipal || !command.Human.Synthetic() {
		return delegationForbidden(command.CorrelationID)
	}
	_, found, err = service.store.Revoke(ctx, command.DelegationID, command.Reason, at)
	if err != nil {
		return normalizeDelegationStoreError(err, command.CorrelationID)
	}
	if !found {
		return delegationForbidden(command.CorrelationID)
	}
	return nil
}

type ResolveDelegationRequest struct {
	ContractVersion string
	CorrelationID   string
	Human           Principal
	Agent           Principal
	DelegationID    string
}

func (service *DelegationService) ResolveDelegationChain(ctx context.Context, command ResolveDelegationRequest) ([]Delegation, error) {
	if service == nil || service.store == nil || service.now == nil {
		return nil, delegationUnavailable(command.CorrelationID)
	}
	if ctx == nil || command.ContractVersion != ResolveDelegationContractVersion || command.CorrelationID == "" || command.DelegationID == "" {
		return nil, delegationInvalid(command.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return nil, delegationUnavailable(command.CorrelationID)
	}
	at := service.now()
	if at.IsZero() {
		return nil, delegationUnavailable(command.CorrelationID)
	}
	stored, found, err := service.store.Get(ctx, command.DelegationID)
	if err != nil {
		return nil, normalizeDelegationStoreError(err, command.CorrelationID)
	}
	if !found || stored.Revoked {
		return nil, delegationForbidden(command.CorrelationID)
	}
	if err := stored.Delegation.ValidateFor(command.Human, command.Agent, at, command.CorrelationID); err != nil {
		return nil, err
	}
	return []Delegation{cloneDelegation(stored.Delegation)}, nil
}

type memoryDelegationEntry struct {
	delegation    Delegation
	requestDigest string
	revoked       bool
	revokedAt     time.Time
	reason        string
}

type MemoryDelegationStore struct {
	mu            sync.RWMutex
	byID          map[string]memoryDelegationEntry
	byIdempotency map[string]string
}

func NewMemoryDelegationStore() *MemoryDelegationStore {
	return &MemoryDelegationStore{byID: make(map[string]memoryDelegationEntry), byIdempotency: make(map[string]string)}
}

func (store *MemoryDelegationStore) Create(_ context.Context, value Delegation, idempotencyKey, requestDigest string) (Delegation, bool, error) {
	if store == nil {
		return Delegation{}, false, ErrDelegationStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	key := value.TenantID() + "\x00" + idempotencyKey
	if existingID, found := store.byIdempotency[key]; found {
		entry := store.byID[existingID]
		if entry.requestDigest != requestDigest {
			return Delegation{}, false, ErrDelegationConflict
		}
		if entry.revoked {
			return Delegation{}, false, ErrDelegationRevoked
		}
		return cloneDelegation(entry.delegation), true, nil
	}
	if existing, found := store.byID[value.ID()]; found && existing.requestDigest != requestDigest {
		return Delegation{}, false, ErrDelegationConflict
	}
	store.byID[value.ID()] = memoryDelegationEntry{delegation: cloneDelegation(value), requestDigest: requestDigest}
	store.byIdempotency[key] = value.ID()
	return cloneDelegation(value), false, nil
}

func (store *MemoryDelegationStore) Get(_ context.Context, id string) (StoredDelegation, bool, error) {
	if store == nil {
		return StoredDelegation{}, false, ErrDelegationStoreUnavailable
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	entry, found := store.byID[id]
	if !found {
		return StoredDelegation{}, false, nil
	}
	return StoredDelegation{Delegation: cloneDelegation(entry.delegation), Revoked: entry.revoked, RevokedAt: entry.revokedAt}, true, nil
}

func (store *MemoryDelegationStore) Revoke(_ context.Context, id, reason string, at time.Time) (StoredDelegation, bool, error) {
	if store == nil {
		return StoredDelegation{}, false, ErrDelegationStoreUnavailable
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	entry, found := store.byID[id]
	if !found {
		return StoredDelegation{}, false, nil
	}
	if !entry.revoked {
		entry.revoked = true
		entry.revokedAt = at
		entry.reason = reason
		store.byID[id] = entry
	}
	return StoredDelegation{Delegation: cloneDelegation(entry.delegation), Revoked: entry.revoked, RevokedAt: entry.revokedAt}, true, nil
}

func validateDelegationPrincipals(request DelegationRequest, human, agent Principal, at time.Time, correlationID string) error {
	if err := human.ValidateAt(at, correlationID); err != nil {
		return err
	}
	if err := agent.ValidateAt(at, correlationID); err != nil {
		return err
	}
	if human.Kind() != HumanPrincipal || agent.Kind() != AgentPrincipal || !human.Synthetic() || !agent.Synthetic() ||
		request.HumanSubjectID() != human.SubjectID() || request.AgentSubjectID() != agent.SubjectID() {
		return delegationForbidden(correlationID)
	}
	return nil
}

func mintDelegation(request DelegationRequest, human Principal, parentDecisionID string, issuedAt, expiresAt time.Time) Delegation {
	value := Delegation{
		contractVersion: DelegationContractVersion, parentAuthorityID: parentDecisionID,
		issuerSubjectID: human.SubjectID(), humanSubjectID: request.HumanSubjectID(),
		agentSubjectID: request.AgentSubjectID(), tenantID: request.TenantID(),
		workspaceID: request.WorkspaceID(), projectID: request.ProjectID(),
		repositoryID: request.RepositoryID(), deliveryRunID: request.DeliveryRunID(),
		purpose: request.Purpose(), actions: request.Actions(), resources: request.Resources(),
		classifications: request.Classifications(), destinations: request.Destinations(),
		environment: request.Environment(), maximumAutonomy: request.MaximumAutonomy(),
		issuedAt: issuedAt, expiresAt: expiresAt, revocationSource: delegationRevocationSource,
		evidenceCorrelationID: request.CorrelationID(), requestDigest: request.Digest(),
	}
	encoded, _ := json.Marshal(value.digestInput())
	digest := sha256.Sum256(encoded)
	value.id = "forge.delegation." + hex.EncodeToString(digest[:])
	return value
}

func (value Delegation) digestInput() any {
	return struct {
		ContractVersion, ParentAuthorityID, IssuerSubjectID, HumanSubjectID, AgentSubjectID string
		TenantID, WorkspaceID, ProjectID, RepositoryID, DeliveryRunID                       string
		Purpose, Environment, MaximumAutonomy, RevocationSource, EvidenceCorrelationID      string
		RequestDigest                                                                       string
		Actions, Resources, Classifications, Destinations                                   []string
		IssuedAt, ExpiresAt                                                                 time.Time
	}{
		value.contractVersion, value.parentAuthorityID, value.issuerSubjectID, value.humanSubjectID, value.agentSubjectID,
		value.tenantID, value.workspaceID, value.projectID, value.repositoryID, value.deliveryRunID,
		value.purpose, value.environment, value.maximumAutonomy, value.revocationSource, value.evidenceCorrelationID,
		value.requestDigest, value.actions, value.resources, value.classifications, value.destinations,
		value.issuedAt, value.expiresAt,
	}
}

func delegationRequestDigest(value DelegationRequest) string {
	encoded, _ := json.Marshal(struct {
		ContractVersion, HumanSubjectID, AgentSubjectID               string
		TenantID, WorkspaceID, ProjectID, RepositoryID, DeliveryRunID string
		Purpose, Environment, MaximumAutonomy, ParentDelegationID     string
		Actions, Resources, Classifications, Destinations             []string
		RequestedExpiresAt                                            time.Time
	}{
		value.contractVersion, value.humanSubjectID, value.agentSubjectID,
		value.tenantID, value.workspaceID, value.projectID, value.repositoryID, value.deliveryRunID,
		value.purpose, value.environment, value.maximumAutonomy, value.parentDelegationID,
		value.actions, value.resources, value.classifications, value.destinations,
		value.requestedExpiresAt,
	})
	digest := sha256.Sum256(encoded)
	return "forge.delegation-request." + hex.EncodeToString(digest[:])
}

func cloneDelegation(value Delegation) Delegation {
	value.actions = slices.Clone(value.actions)
	value.resources = slices.Clone(value.resources)
	value.classifications = slices.Clone(value.classifications)
	value.destinations = slices.Clone(value.destinations)
	return value
}

func canonicalRequired(values []string) ([]string, error) {
	if len(values) == 0 {
		return nil, errors.New("values are required")
	}
	return canonicalOptional(values)
}

func canonicalOptional(values []string) ([]string, error) {
	result := slices.Clone(values)
	sort.Strings(result)
	for index, value := range result {
		if value == "" || (index > 0 && value == result[index-1]) {
			return nil, errors.New("values must be unique and non-empty")
		}
	}
	if result == nil {
		result = []string{}
	}
	return result, nil
}

func earliest(values ...time.Time) time.Time {
	result := values[0]
	for _, value := range values[1:] {
		if value.Before(result) {
			result = value
		}
	}
	return result
}

func validDelegationClassification(value string) bool {
	return value == "Public" || value == "Internal" || value == "Confidential" || value == "Restricted"
}

func validRootDelegationShape(actions, resources, destinations []string, repositoryID, autonomy string) bool {
	if len(resources) != 1 {
		return false
	}
	if slices.Equal(actions, []string{DelegatedActionAssembleContext}) {
		return autonomy == ContextPreparationOnly && resources[0] == repositoryID && len(destinations) == 0
	}
	return slices.Equal(actions, []string{DelegatedActionProposeChange}) &&
		autonomy == ChangeProposalOnly && validSHA256Digest(resources[0]) &&
		len(destinations) == 1 && strings.HasPrefix(destinations[0], "model:")
}

func validSHA256Digest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != len("sha256:")+sha256.Size*2 {
		return false
	}
	decoded, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
	return err == nil && len(decoded) == sha256.Size
}

func normalizeDelegationStoreError(err error, correlationID string) error {
	switch {
	case errors.Is(err, ErrDelegationRevoked):
		return delegationForbidden(correlationID)
	case errors.Is(err, ErrDelegationConflict):
		return failure.New(failure.Conflict, "delegation request conflicts with existing authority", correlationID, failure.RetryNone)
	case errors.Is(err, ErrDelegationStoreUnavailable):
		return delegationUnavailable(correlationID)
	default:
		return failure.New(failure.Indeterminate, "delegation persistence could not be verified", correlationID, failure.RetryNone)
	}
}

func delegationInvalid(correlationID string) error {
	return failure.New(failure.InvalidRequest, "delegation request is invalid", correlationID, failure.RetryNone)
}

func delegationForbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "task delegation is not authorized", correlationID, failure.RetryNone)
}

func delegationUnavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "delegation service is unavailable", correlationID, failure.RetryBounded)
}
