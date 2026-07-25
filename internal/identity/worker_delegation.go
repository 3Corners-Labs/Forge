package identity

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
)

const (
	WorkerDelegationRequestVersion = "forge.worker-delegation-request.v1"
	CreateWorkerDelegationVersion  = "forge.create-worker-delegation.v1"
	ResolveWorkerDelegationVersion = "forge.resolve-worker-delegation.v1"
	WorkerDelegationVersion        = "forge.worker-delegation.v1"
	WorkerPrepareChangeAction      = "forge.worker.prepare-change"
	WorkerMaximumAutonomy          = "approved-plan-change-preparation-only"
	WorkerM1Purpose                = "forge.m1.prepare-change"
)

type WorkerDelegationRequest struct {
	SchemaVersion         string
	ID                    string
	Digest                string
	IdempotencyKey        string
	CorrelationID         string
	HumanSubjectID        string
	TaskAgentSubjectID    string
	WorkerSubjectID       string
	TenantID              string
	WorkspaceID           string
	ProjectID             string
	RepositoryID          string
	DeliveryRunID         string
	PlanID                string
	PlanDigest            string
	ChangeProposalDigest  string
	ContextPackageID      string
	ContextManifestDigest string
	ApprovalRequestID     string
	ApprovalRequestDigest string
	ApprovalRecordID      string
	ApprovalRecordDigest  string
	SourceRevision        string
	ChangePaths           []string
	Classification        string
	RiskClass             string
	Environment           string
	Action                string
	Purpose               string
	MaximumAutonomy       string
	MaxDurationSeconds    int64
	MaxProcesses          int64
	MaxOutputBytes        int64
	MaxChangedFiles       int64
	MaxAttempts           int64
	RequestedExpiresAt    time.Time
}

func BindWorkerDelegationRequest(value WorkerDelegationRequest) (WorkerDelegationRequest, error) {
	value.ID, value.Digest = "", ""
	sort.Strings(value.ChangePaths)
	value.ChangePaths = slices.Compact(value.ChangePaths)
	if value.SchemaVersion != WorkerDelegationRequestVersion || value.IdempotencyKey == "" ||
		value.CorrelationID == "" || value.HumanSubjectID == "" || value.TaskAgentSubjectID == "" ||
		value.WorkerSubjectID == "" || value.HumanSubjectID == value.TaskAgentSubjectID ||
		value.HumanSubjectID == value.WorkerSubjectID || value.TaskAgentSubjectID == value.WorkerSubjectID ||
		value.TenantID == "" || value.WorkspaceID == "" || value.ProjectID == "" ||
		value.RepositoryID == "" || value.DeliveryRunID == "" || value.PlanID == "" ||
		!workerDigest(value.PlanDigest) || !workerDigest(value.ChangeProposalDigest) || value.ContextPackageID == "" ||
		!workerDigest(value.ContextManifestDigest) || value.ApprovalRequestID == "" ||
		!workerDigest(value.ApprovalRequestDigest) || value.ApprovalRecordID == "" ||
		!workerDigest(value.ApprovalRecordDigest) || value.SourceRevision == "" ||
		len(value.ChangePaths) == 0 || value.Classification == "" || value.RiskClass == "" ||
		value.Environment != DeveloperTrialEnvironment || value.Action != WorkerPrepareChangeAction ||
		value.Purpose != WorkerM1Purpose || value.MaximumAutonomy != WorkerMaximumAutonomy ||
		value.MaxDurationSeconds < 1 || value.MaxProcesses < 1 || value.MaxOutputBytes < 1 ||
		value.MaxChangedFiles < 1 || value.MaxAttempts < 1 || value.RequestedExpiresAt.IsZero() {
		return WorkerDelegationRequest{}, errors.New("worker delegation request is outside M1")
	}
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	value.Digest = "sha256:" + hex.EncodeToString(sum[:])
	value.ID = "forge.worker-delegation-request." + hex.EncodeToString(sum[:])
	return value, nil
}

func (value WorkerDelegationRequest) Validate(at time.Time) error {
	bound, err := BindWorkerDelegationRequest(value)
	if err != nil || bound.ID != value.ID || bound.Digest != value.Digest ||
		at.IsZero() || !at.Before(value.RequestedExpiresAt) {
		return errors.New("worker delegation request is invalid")
	}
	return nil
}

type WorkerDelegation struct {
	SchemaVersion          string
	ID                     string
	RequestID              string
	RequestDigest          string
	ParentPolicyDecisionID string
	HumanSubjectID         string
	TaskAgentSubjectID     string
	WorkerSubjectID        string
	TenantID               string
	WorkspaceID            string
	ProjectID              string
	RepositoryID           string
	DeliveryRunID          string
	PlanDigest             string
	ChangeProposalDigest   string
	ContextManifestDigest  string
	ApprovalRecordID       string
	ApprovalRecordDigest   string
	SourceRevision         string
	ChangePaths            []string
	Classification         string
	RiskClass              string
	Environment            string
	Action                 string
	Purpose                string
	MaximumAutonomy        string
	MaxDurationSeconds     int64
	MaxProcesses           int64
	MaxOutputBytes         int64
	MaxChangedFiles        int64
	MaxAttempts            int64
	IssuedAt               time.Time
	ExpiresAt              time.Time
	EvidenceCorrelationID  string
}

func (value WorkerDelegation) ValidateFor(human, taskAgent, worker Principal, at time.Time, correlationID string) error {
	for _, principal := range []Principal{human, taskAgent, worker} {
		if err := principal.ValidateAt(at, correlationID); err != nil {
			return err
		}
	}
	if human.Kind() != HumanPrincipal || taskAgent.Kind() != AgentPrincipal ||
		worker.Kind() != WorkloadPrincipal || !human.Synthetic() || !taskAgent.Synthetic() ||
		!worker.Synthetic() || value.SchemaVersion != WorkerDelegationVersion || value.ID == "" ||
		value.RequestID == "" || !workerDigest(value.RequestDigest) || value.ParentPolicyDecisionID == "" ||
		value.HumanSubjectID != human.SubjectID() || value.TaskAgentSubjectID != taskAgent.SubjectID() ||
		value.WorkerSubjectID != worker.SubjectID() || value.TenantID == "" || value.WorkspaceID == "" ||
		value.ProjectID == "" || value.RepositoryID == "" || value.DeliveryRunID == "" ||
		!workerDigest(value.PlanDigest) || !workerDigest(value.ChangeProposalDigest) || !workerDigest(value.ContextManifestDigest) ||
		value.ApprovalRecordID == "" || !workerDigest(value.ApprovalRecordDigest) ||
		value.SourceRevision == "" || len(value.ChangePaths) == 0 ||
		value.Environment != DeveloperTrialEnvironment || value.Action != WorkerPrepareChangeAction ||
		value.Purpose != WorkerM1Purpose || value.MaximumAutonomy != WorkerMaximumAutonomy ||
		value.EvidenceCorrelationID == "" || value.IssuedAt.IsZero() ||
		!value.ExpiresAt.After(value.IssuedAt) || at.Before(value.IssuedAt) || !at.Before(value.ExpiresAt) {
		return failure.New(failure.Forbidden, "worker delegation is invalid", correlationID, failure.RetryNewAuthority)
	}
	return nil
}

type StoredWorkerDelegation struct {
	Delegation WorkerDelegation
	Revoked    bool
	RevokedAt  time.Time
}

type WorkerDelegationStore interface {
	CreateWorker(context.Context, WorkerDelegation, string) (WorkerDelegation, bool, error)
	GetWorker(context.Context, string) (StoredWorkerDelegation, bool, error)
	RevokeWorker(context.Context, string, time.Time) (StoredWorkerDelegation, bool, error)
}

type WorkerDelegationService struct {
	store       WorkerDelegationStore
	maxLifetime time.Duration
	now         func() time.Time
}

func NewWorkerDelegationService(store WorkerDelegationStore, maxLifetime time.Duration, now func() time.Time) (*WorkerDelegationService, error) {
	if store == nil || maxLifetime <= 0 || now == nil {
		return nil, errors.New("worker delegation service configuration is invalid")
	}
	return &WorkerDelegationService{store: store, maxLifetime: maxLifetime, now: now}, nil
}

type CreateWorkerDelegationRequest struct {
	SchemaVersion           string
	Request                 WorkerDelegationRequest
	Human                   Principal
	TaskAgent               Principal
	Worker                  Principal
	ParentDecision          trust.PolicyDecision
	EnforcedObligationTypes []string
}

func (service *WorkerDelegationService) Create(ctx context.Context, command CreateWorkerDelegationRequest) (WorkerDelegation, bool, error) {
	at := service.now()
	correlationID := command.Request.CorrelationID
	if ctx == nil || command.SchemaVersion != CreateWorkerDelegationVersion || at.IsZero() {
		return WorkerDelegation{}, false, failure.New(failure.InvalidRequest, "worker delegation command is invalid", correlationID, failure.RetryNone)
	}
	if err := command.Request.Validate(at); err != nil {
		return WorkerDelegation{}, false, failure.New(failure.InvalidRequest, "worker delegation request is invalid", correlationID, failure.RetryNone)
	}
	for _, principal := range []Principal{command.Human, command.TaskAgent, command.Worker} {
		if err := principal.ValidateAt(at, correlationID); err != nil {
			return WorkerDelegation{}, false, err
		}
	}
	if command.Human.Kind() != HumanPrincipal || command.TaskAgent.Kind() != AgentPrincipal ||
		command.Worker.Kind() != WorkloadPrincipal || command.Request.HumanSubjectID != command.Human.SubjectID() ||
		command.Request.TaskAgentSubjectID != command.TaskAgent.SubjectID() ||
		command.Request.WorkerSubjectID != command.Worker.SubjectID() ||
		command.ParentDecision.ExpiresAt == nil {
		return WorkerDelegation{}, false, failure.New(failure.Forbidden, "worker delegation principals are invalid", correlationID, failure.RetryNewAuthority)
	}
	if err := command.ParentDecision.Authorize("execution", command.Human.SubjectID(), DelegationActionCreate, command.Request.Digest, at, command.EnforcedObligationTypes, correlationID); err != nil {
		return WorkerDelegation{}, false, err
	}
	expiresAt := earliest(command.Request.RequestedExpiresAt, command.Human.ExpiresAt(), command.TaskAgent.ExpiresAt(), command.Worker.ExpiresAt(), *command.ParentDecision.ExpiresAt, at.Add(service.maxLifetime))
	value := mintWorkerDelegation(command.Request, command.ParentDecision.ID, at, expiresAt)
	stored, replayed, err := service.store.CreateWorker(ctx, value, command.Request.IdempotencyKey)
	if err != nil {
		return WorkerDelegation{}, false, failure.New(failure.DependencyUnavailable, "worker delegation persistence is unavailable", correlationID, failure.RetryBounded)
	}
	if err := stored.ValidateFor(command.Human, command.TaskAgent, command.Worker, at, correlationID); err != nil {
		return WorkerDelegation{}, false, err
	}
	return cloneWorkerDelegation(stored), replayed, nil
}

type ResolveWorkerDelegationRequest struct {
	SchemaVersion string
	CorrelationID string
	Human         Principal
	TaskAgent     Principal
	Worker        Principal
	DelegationID  string
}

func (service *WorkerDelegationService) Resolve(ctx context.Context, command ResolveWorkerDelegationRequest) (WorkerDelegation, error) {
	at := service.now()
	if ctx == nil || command.SchemaVersion != ResolveWorkerDelegationVersion || command.DelegationID == "" {
		return WorkerDelegation{}, failure.New(failure.InvalidRequest, "worker delegation resolution is invalid", command.CorrelationID, failure.RetryNone)
	}
	stored, found, err := service.store.GetWorker(ctx, command.DelegationID)
	if err != nil {
		return WorkerDelegation{}, failure.New(failure.DependencyUnavailable, "worker delegation resolution is unavailable", command.CorrelationID, failure.RetryBounded)
	}
	if !found || stored.Revoked {
		return WorkerDelegation{}, failure.New(failure.Forbidden, "worker delegation is unavailable", command.CorrelationID, failure.RetryNewAuthority)
	}
	if err := stored.Delegation.ValidateFor(command.Human, command.TaskAgent, command.Worker, at, command.CorrelationID); err != nil {
		return WorkerDelegation{}, err
	}
	return cloneWorkerDelegation(stored.Delegation), nil
}

func (service *WorkerDelegationService) ResolveWorkerDelegation(ctx context.Context, command ResolveWorkerDelegationRequest) (WorkerDelegation, error) {
	return service.Resolve(ctx, command)
}

type memoryWorkerEntry struct {
	value          WorkerDelegation
	idempotencyKey string
	revoked        bool
	revokedAt      time.Time
}

type MemoryWorkerDelegationStore struct {
	mu            sync.RWMutex
	byID          map[string]memoryWorkerEntry
	byIdempotency map[string]string
}

func NewMemoryWorkerDelegationStore() *MemoryWorkerDelegationStore {
	return &MemoryWorkerDelegationStore{byID: make(map[string]memoryWorkerEntry), byIdempotency: make(map[string]string)}
}

func (store *MemoryWorkerDelegationStore) CreateWorker(_ context.Context, value WorkerDelegation, idempotencyKey string) (WorkerDelegation, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := value.TenantID + "\x00" + idempotencyKey
	if id, found := store.byIdempotency[key]; found {
		existing := store.byID[id]
		if existing.value.RequestDigest != value.RequestDigest || existing.revoked {
			return WorkerDelegation{}, false, ErrDelegationConflict
		}
		return cloneWorkerDelegation(existing.value), true, nil
	}
	store.byID[value.ID] = memoryWorkerEntry{value: cloneWorkerDelegation(value), idempotencyKey: idempotencyKey}
	store.byIdempotency[key] = value.ID
	return cloneWorkerDelegation(value), false, nil
}

func (store *MemoryWorkerDelegationStore) GetWorker(_ context.Context, id string) (StoredWorkerDelegation, bool, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	entry, found := store.byID[id]
	return StoredWorkerDelegation{Delegation: cloneWorkerDelegation(entry.value), Revoked: entry.revoked, RevokedAt: entry.revokedAt}, found, nil
}

func (store *MemoryWorkerDelegationStore) RevokeWorker(_ context.Context, id string, at time.Time) (StoredWorkerDelegation, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	entry, found := store.byID[id]
	if !found {
		return StoredWorkerDelegation{}, false, nil
	}
	entry.revoked, entry.revokedAt = true, at
	store.byID[id] = entry
	return StoredWorkerDelegation{Delegation: cloneWorkerDelegation(entry.value), Revoked: true, RevokedAt: at}, true, nil
}

func mintWorkerDelegation(request WorkerDelegationRequest, parent string, issuedAt, expiresAt time.Time) WorkerDelegation {
	value := WorkerDelegation{
		SchemaVersion: WorkerDelegationVersion, RequestID: request.ID, RequestDigest: request.Digest,
		ParentPolicyDecisionID: parent, HumanSubjectID: request.HumanSubjectID,
		TaskAgentSubjectID: request.TaskAgentSubjectID, WorkerSubjectID: request.WorkerSubjectID,
		TenantID: request.TenantID, WorkspaceID: request.WorkspaceID, ProjectID: request.ProjectID,
		RepositoryID: request.RepositoryID, DeliveryRunID: request.DeliveryRunID,
		PlanDigest: request.PlanDigest, ChangeProposalDigest: request.ChangeProposalDigest,
		ContextManifestDigest: request.ContextManifestDigest,
		ApprovalRecordID:      request.ApprovalRecordID, ApprovalRecordDigest: request.ApprovalRecordDigest,
		SourceRevision: request.SourceRevision, ChangePaths: slices.Clone(request.ChangePaths),
		Classification: request.Classification, RiskClass: request.RiskClass,
		Environment: request.Environment, Action: request.Action, Purpose: request.Purpose,
		MaximumAutonomy: request.MaximumAutonomy, MaxDurationSeconds: request.MaxDurationSeconds,
		MaxProcesses: request.MaxProcesses, MaxOutputBytes: request.MaxOutputBytes,
		MaxChangedFiles: request.MaxChangedFiles, MaxAttempts: request.MaxAttempts,
		IssuedAt: issuedAt, ExpiresAt: expiresAt, EvidenceCorrelationID: request.CorrelationID,
	}
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	value.ID = "forge.worker-delegation." + hex.EncodeToString(sum[:])
	return value
}

func cloneWorkerDelegation(value WorkerDelegation) WorkerDelegation {
	value.ChangePaths = slices.Clone(value.ChangePaths)
	return value
}

func workerDigest(value string) bool {
	if len(value) != 71 || value[:7] != "sha256:" {
		return false
	}
	_, err := hex.DecodeString(value[7:])
	return err == nil
}
