// Package delivery owns the bounded M1 intent and delivery-run record. This
// iteration captures only the initial state; it does not evaluate policy,
// issue execution authority, schedule work, or publish source changes.
package delivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const (
	CaptureIntentRequestContractVersion = "forge.capture-intent-request.v1"
	CaptureIntentResultContractVersion  = "forge.capture-intent-result.v1"
	DeliveryRunContractVersion          = "forge.delivery-run.v1"
	CapturedState                       = "Captured"
	M1ChangePurpose                     = "forge.m1.prepare-change"
)

// ErrIdempotencyConflict is returned by Store implementations when a key was
// already committed for different immutable intent input in the same scope.
var ErrIdempotencyConflict = errors.New("delivery idempotency conflict")

type CaptureIntentRequest struct {
	ContractVersion    string
	CorrelationID      string
	IdempotencyKey     string
	Principal          identity.Principal
	DelegationChain    []string
	Scope              tenant.Scope
	RepositoryBinding  tenant.RepositoryBinding
	Objective          string
	Constraints        []string
	AcceptanceCriteria []string
}

// CaptureRecord is the capability-owned persistence command. Store adapters
// must commit it atomically by ScopeKey and IdempotencyKey.
type CaptureRecord struct {
	ScopeKey           string
	RequestFingerprint string
	IdempotencyKey     string
	CorrelationID      string
	TenantID           string
	WorkspaceID        string
	ProjectID          string
	RepositoryID       string
	ActorID            string
	DelegationChain    []string
	Purpose            string
	Classification     string
	Objective          string
	Constraints        []string
	AcceptanceCriteria []string
	CreatedAt          time.Time
}

// StoredRecord is returned across the replaceable persistence boundary and is
// revalidated before it becomes an immutable DeliveryRun.
type StoredRecord struct {
	CaptureRecord
	ContractVersion string
	ID              string
	State           string
	Version         uint64
}

type Store interface {
	CreateCaptured(context.Context, CaptureRecord) (StoredRecord, bool, error)
}

// DeliveryRun exposes canonical captured state through immutable fields and
// defensive slice copies.
type DeliveryRun struct {
	contractVersion    string
	id                 string
	state              string
	version            uint64
	tenantID           string
	workspaceID        string
	projectID          string
	repositoryID       string
	actorID            string
	delegationChain    []string
	purpose            string
	classification     string
	objective          string
	constraints        []string
	acceptanceCriteria []string
	createdAt          time.Time
	correlationID      string
}

func (value DeliveryRun) ContractVersion() string      { return value.contractVersion }
func (value DeliveryRun) ID() string                   { return value.id }
func (value DeliveryRun) State() string                { return value.state }
func (value DeliveryRun) Version() uint64              { return value.version }
func (value DeliveryRun) TenantID() string             { return value.tenantID }
func (value DeliveryRun) WorkspaceID() string          { return value.workspaceID }
func (value DeliveryRun) ProjectID() string            { return value.projectID }
func (value DeliveryRun) RepositoryID() string         { return value.repositoryID }
func (value DeliveryRun) ActorID() string              { return value.actorID }
func (value DeliveryRun) DelegationChain() []string    { return slices.Clone(value.delegationChain) }
func (value DeliveryRun) Purpose() string              { return value.purpose }
func (value DeliveryRun) Classification() string       { return value.classification }
func (value DeliveryRun) Objective() string            { return value.objective }
func (value DeliveryRun) Constraints() []string        { return slices.Clone(value.constraints) }
func (value DeliveryRun) AcceptanceCriteria() []string { return slices.Clone(value.acceptanceCriteria) }
func (value DeliveryRun) CreatedAt() time.Time         { return value.createdAt }
func (value DeliveryRun) CorrelationID() string        { return value.correlationID }

type CaptureIntentResult struct {
	ContractVersion string
	Run             DeliveryRun
	Replayed        bool
}

type Service struct {
	store Store
	now   func() time.Time
}

func NewService(store Store, now func() time.Time) (*Service, error) {
	if store == nil {
		return nil, errors.New("delivery store is required")
	}
	if now == nil {
		return nil, errors.New("delivery clock is required")
	}
	return &Service{store: store, now: now}, nil
}

func (service *Service) CaptureIntent(ctx context.Context, request CaptureIntentRequest) (CaptureIntentResult, error) {
	if service == nil || service.store == nil || service.now == nil {
		return CaptureIntentResult{}, failure.New(failure.DependencyUnavailable, "intent capture is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if request.ContractVersion != CaptureIntentRequestContractVersion {
		return CaptureIntentResult{}, failure.New(failure.InvalidRequest, "intent capture request version is invalid", request.CorrelationID, failure.RetryNone)
	}
	if ctx == nil || request.CorrelationID == "" || request.IdempotencyKey == "" {
		return CaptureIntentResult{}, failure.New(failure.InvalidRequest, "intent capture context, correlation, and idempotency are required", request.CorrelationID, failure.RetryNone)
	}
	if request.DelegationChain == nil {
		return CaptureIntentResult{}, failure.New(failure.InvalidRequest, "intent capture delegation chain is required", request.CorrelationID, failure.RetryNone)
	}
	at := service.now()
	if at.IsZero() {
		return CaptureIntentResult{}, failure.New(failure.DependencyUnavailable, "intent capture is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if err := request.Principal.ValidateAt(at, request.CorrelationID); err != nil {
		return CaptureIntentResult{}, err
	}
	if request.Principal.Kind() != identity.HumanPrincipal || len(request.DelegationChain) != 0 {
		return CaptureIntentResult{}, failure.New(failure.Forbidden, "intent capture requires a direct authenticated human", request.CorrelationID, failure.RetryNone)
	}
	if request.Objective == "" || request.Constraints == nil || len(request.AcceptanceCriteria) == 0 {
		return CaptureIntentResult{}, failure.New(failure.InvalidRequest, "objective, constraints, and acceptance criteria are required", request.CorrelationID, failure.RetryNone)
	}
	if !allNonEmpty(request.Constraints) || !allNonEmpty(request.AcceptanceCriteria) {
		return CaptureIntentResult{}, failure.New(failure.InvalidRequest, "constraints and acceptance criteria cannot contain empty values", request.CorrelationID, failure.RetryNone)
	}
	if err := request.Scope.ValidateFor(request.Principal, at, request.CorrelationID); err != nil {
		return CaptureIntentResult{}, err
	}
	if err := request.RepositoryBinding.ValidateFor(request.Scope, request.CorrelationID); err != nil {
		return CaptureIntentResult{}, err
	}

	record := CaptureRecord{
		ScopeKey:           scopeKey(request),
		IdempotencyKey:     request.IdempotencyKey,
		CorrelationID:      request.CorrelationID,
		TenantID:           request.Scope.TenantID(),
		WorkspaceID:        request.Scope.WorkspaceID(),
		ProjectID:          request.Scope.ProjectID(),
		RepositoryID:       request.RepositoryBinding.RepositoryID(),
		ActorID:            request.Principal.SubjectID(),
		DelegationChain:    slices.Clone(request.DelegationChain),
		Purpose:            M1ChangePurpose,
		Classification:     request.RepositoryBinding.Classification(),
		Objective:          request.Objective,
		Constraints:        slices.Clone(request.Constraints),
		AcceptanceCriteria: slices.Clone(request.AcceptanceCriteria),
		CreatedAt:          at,
	}
	record.RequestFingerprint = fingerprint(record)
	stored, replayed, err := service.store.CreateCaptured(ctx, cloneCaptureRecord(record))
	if err != nil {
		if errors.Is(err, ErrIdempotencyConflict) {
			return CaptureIntentResult{}, failure.New(failure.Conflict, "idempotency key belongs to different intent", request.CorrelationID, failure.RetryReconcile)
		}
		return CaptureIntentResult{}, failure.New(failure.DependencyUnavailable, "intent capture is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !validStoredRecord(stored, record, replayed) {
		return CaptureIntentResult{}, failure.New(failure.Indeterminate, "captured intent could not be verified", request.CorrelationID, failure.RetryNone)
	}
	return CaptureIntentResult{
		ContractVersion: CaptureIntentResultContractVersion,
		Run:             runFromStored(stored),
		Replayed:        replayed,
	}, nil
}

func validStoredRecord(stored StoredRecord, requested CaptureRecord, replayed bool) bool {
	if stored.ContractVersion != DeliveryRunContractVersion || stored.ID == "" || stored.State != CapturedState || stored.Version != 1 || stored.CreatedAt.IsZero() || stored.CorrelationID == "" {
		return false
	}
	if !replayed && (stored.CreatedAt != requested.CreatedAt || stored.CorrelationID != requested.CorrelationID) {
		return false
	}
	return stored.ScopeKey == requested.ScopeKey &&
		stored.RequestFingerprint == requested.RequestFingerprint &&
		stored.IdempotencyKey == requested.IdempotencyKey &&
		stored.TenantID == requested.TenantID &&
		stored.WorkspaceID == requested.WorkspaceID &&
		stored.ProjectID == requested.ProjectID &&
		stored.RepositoryID == requested.RepositoryID &&
		stored.ActorID == requested.ActorID &&
		slices.Equal(stored.DelegationChain, requested.DelegationChain) &&
		stored.Purpose == requested.Purpose &&
		stored.Classification == requested.Classification &&
		stored.Objective == requested.Objective &&
		slices.Equal(stored.Constraints, requested.Constraints) &&
		slices.Equal(stored.AcceptanceCriteria, requested.AcceptanceCriteria)
}

func runFromStored(value StoredRecord) DeliveryRun {
	return DeliveryRun{
		contractVersion:    value.ContractVersion,
		id:                 value.ID,
		state:              value.State,
		version:            value.Version,
		tenantID:           value.TenantID,
		workspaceID:        value.WorkspaceID,
		projectID:          value.ProjectID,
		repositoryID:       value.RepositoryID,
		actorID:            value.ActorID,
		delegationChain:    slices.Clone(value.DelegationChain),
		purpose:            value.Purpose,
		classification:     value.Classification,
		objective:          value.Objective,
		constraints:        slices.Clone(value.Constraints),
		acceptanceCriteria: slices.Clone(value.AcceptanceCriteria),
		createdAt:          value.CreatedAt,
		correlationID:      value.CorrelationID,
	}
}

func fingerprint(value CaptureRecord) string {
	encoded, _ := json.Marshal(struct {
		TenantID           string
		WorkspaceID        string
		ProjectID          string
		RepositoryID       string
		ActorID            string
		DelegationChain    []string
		Purpose            string
		Classification     string
		Objective          string
		Constraints        []string
		AcceptanceCriteria []string
	}{
		TenantID:           value.TenantID,
		WorkspaceID:        value.WorkspaceID,
		ProjectID:          value.ProjectID,
		RepositoryID:       value.RepositoryID,
		ActorID:            value.ActorID,
		DelegationChain:    value.DelegationChain,
		Purpose:            value.Purpose,
		Classification:     value.Classification,
		Objective:          value.Objective,
		Constraints:        value.Constraints,
		AcceptanceCriteria: value.AcceptanceCriteria,
	})
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:])
}

func scopeKey(request CaptureIntentRequest) string {
	return request.Scope.TenantID() + "\x00" + request.Scope.WorkspaceID() + "\x00" + request.Scope.ProjectID() + "\x00" + request.RepositoryBinding.RepositoryID() + "\x00" + request.Principal.SubjectID()
}

func allNonEmpty(values []string) bool {
	for _, value := range values {
		if value == "" {
			return false
		}
	}
	return true
}

func cloneCaptureRecord(value CaptureRecord) CaptureRecord {
	value.DelegationChain = slices.Clone(value.DelegationChain)
	value.Constraints = slices.Clone(value.Constraints)
	value.AcceptanceCriteria = slices.Clone(value.AcceptanceCriteria)
	return value
}
