// Package intent owns promotion of an exact Experience proposal into canonical
// intent. It does not own interaction drafts, delivery state, or approval.
package intent

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
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const (
	PromotionRequestVersion = "forge.intent-promote-request.v1"
	CanonicalIntentVersion  = "forge.canonical-intent.v1"
)

type Question struct {
	ID          string
	Text        string
	Consequence string
	Material    bool
	Resolved    bool
}

type PromotionSubject struct {
	ProposalID           string
	ProposalVersion      uint64
	ProposalDigest       string
	Current              bool
	ActorID              string
	TenantID             string
	WorkspaceID          string
	ProjectID            string
	RepositoryID         string
	Outcome              string
	AffectedUsers        []string
	Problem              string
	FunctionalBehavior   []string
	UXExpectations       []string
	Accessibility        []string
	Privacy              []string
	Security             []string
	Operations           []string
	Constraints          []string
	Exclusions           []string
	AcceptanceCriteria   []string
	DirectUserStatements []string
	Assumptions          []string
	Hypotheses           []string
	EvidenceReferences   []string
	Questions            []Question
	Experiments          []string
	Provenance           []string
	ChangeExplanation    string
	CreatedAt            time.Time
	CorrelationID        string
}

func BindProposalDigest(value PromotionSubject) PromotionSubject {
	value = canonicalSubject(value)
	value.ProposalDigest = proposalDigest(value)
	return value
}

func VerifyProposalDigest(value PromotionSubject) bool {
	return value.ProposalDigest != "" && value.ProposalDigest == proposalDigest(canonicalSubject(value))
}

type ProposalResolver interface {
	ResolvePromotionSubject(context.Context, string, string) (PromotionSubject, bool, error)
}

type DeliveryCapturer interface {
	CaptureIntent(context.Context, delivery.CaptureIntentRequest) (delivery.CaptureIntentResult, error)
}

type Record struct {
	ContractVersion string
	ID              string
	ProposalID      string
	ProposalDigest  string
	DeliveryRunID   string
	ActorID         string
	TenantID        string
	WorkspaceID     string
	ProjectID       string
	RepositoryID    string
	Outcome         string
	CreatedAt       time.Time
}

type Store interface {
	Create(context.Context, string, string, string, Record) (Record, bool, error)
}

type PromoteRequest struct {
	ContractVersion   string
	CorrelationID     string
	IdempotencyKey    string
	ProposalID        string
	ProposalDigest    string
	Principal         identity.Principal
	Scope             tenant.Scope
	RepositoryBinding tenant.RepositoryBinding
}

type PromoteResult struct {
	Intent   Record
	Run      delivery.DeliveryRun
	Replayed bool
}

type Service struct {
	proposals ProposalResolver
	delivery  DeliveryCapturer
	store     Store
	now       func() time.Time
	newID     func() string
}

func NewService(proposals ProposalResolver, deliveryService DeliveryCapturer, store Store, now func() time.Time, newID func() string) (*Service, error) {
	if proposals == nil || deliveryService == nil || store == nil || now == nil || newID == nil {
		return nil, errors.New("intent promotion configuration is invalid")
	}
	return &Service{proposals: proposals, delivery: deliveryService, store: store, now: now, newID: newID}, nil
}

func (service *Service) Promote(ctx context.Context, request PromoteRequest) (PromoteResult, error) {
	if service == nil || service.proposals == nil || service.delivery == nil || service.store == nil || service.now == nil || service.newID == nil {
		return PromoteResult{}, unavailable(request.CorrelationID)
	}
	if ctx == nil || request.ContractVersion != PromotionRequestVersion || request.CorrelationID == "" ||
		request.IdempotencyKey == "" || request.ProposalID == "" || request.ProposalDigest == "" {
		return PromoteResult{}, invalid(request.CorrelationID)
	}
	at := service.now()
	if err := request.Principal.ValidateAt(at, request.CorrelationID); err != nil {
		return PromoteResult{}, err
	}
	if request.Principal.Kind() != identity.HumanPrincipal {
		return PromoteResult{}, forbidden(request.CorrelationID)
	}
	if err := request.Scope.ValidateFor(request.Principal, at, request.CorrelationID); err != nil {
		return PromoteResult{}, err
	}
	if err := request.RepositoryBinding.ValidateFor(request.Scope, request.CorrelationID); err != nil {
		return PromoteResult{}, err
	}
	subject, found, err := service.proposals.ResolvePromotionSubject(ctx, request.Scope.TenantID(), request.ProposalID)
	if err != nil {
		return PromoteResult{}, unavailable(request.CorrelationID)
	}
	if !found || !subject.Current || subject.ProposalID != request.ProposalID ||
		subject.ProposalDigest != request.ProposalDigest || !VerifyProposalDigest(subject) ||
		subject.ActorID != request.Principal.SubjectID() ||
		subject.TenantID != request.Scope.TenantID() ||
		subject.WorkspaceID != request.Scope.WorkspaceID() ||
		subject.ProjectID != request.Scope.ProjectID() ||
		subject.RepositoryID != request.RepositoryBinding.RepositoryID() ||
		subject.CorrelationID != request.CorrelationID {
		return PromoteResult{}, forbidden(request.CorrelationID)
	}
	for _, question := range subject.Questions {
		if question.Material && !question.Resolved && question.Consequence != "" {
			return PromoteResult{}, failure.New(failure.Conflict, "material requirements question remains unresolved", request.CorrelationID, failure.RetryReconcile)
		}
	}
	captured, err := service.delivery.CaptureIntent(ctx, delivery.CaptureIntentRequest{
		ContractVersion: delivery.CaptureIntentRequestContractVersion,
		CorrelationID:   request.CorrelationID, IdempotencyKey: "proposal:" + request.ProposalDigest,
		Principal: request.Principal, DelegationChain: []string{}, Scope: request.Scope,
		RepositoryBinding: request.RepositoryBinding, Objective: subject.Outcome,
		Constraints: slices.Clone(subject.Constraints), AcceptanceCriteria: slices.Clone(subject.AcceptanceCriteria),
	})
	if err != nil {
		return PromoteResult{}, err
	}
	record := Record{
		ContractVersion: CanonicalIntentVersion, ID: service.newID(),
		ProposalID: subject.ProposalID, ProposalDigest: subject.ProposalDigest,
		DeliveryRunID: captured.Run.ID(), ActorID: subject.ActorID,
		TenantID: subject.TenantID, WorkspaceID: subject.WorkspaceID,
		ProjectID: subject.ProjectID, RepositoryID: subject.RepositoryID,
		Outcome: subject.Outcome, CreatedAt: at,
	}
	stored, replayed, err := service.store.Create(ctx, subject.TenantID, request.IdempotencyKey, promotionFingerprint(record), record)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			return PromoteResult{}, failure.New(failure.Conflict, "intent promotion identity belongs to different input", request.CorrelationID, failure.RetryReconcile)
		}
		return PromoteResult{}, unavailable(request.CorrelationID)
	}
	if !validRecord(stored, subject, captured.Run) {
		return PromoteResult{}, failure.New(failure.Indeterminate, "canonical intent could not be verified", request.CorrelationID, failure.RetryNone)
	}
	return PromoteResult{Intent: stored, Run: captured.Run, Replayed: replayed || captured.Replayed}, nil
}

var ErrConflict = errors.New("intent promotion conflict")

type memoryEntry struct {
	fingerprint string
	record      Record
}

type MemoryStore struct {
	mu      sync.Mutex
	entries map[string]memoryEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{entries: make(map[string]memoryEntry)}
}

func (store *MemoryStore) Create(_ context.Context, tenantID, idempotencyKey, fingerprint string, record Record) (Record, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := tenantID + "\x00" + idempotencyKey
	if existing, found := store.entries[key]; found {
		if existing.fingerprint != fingerprint {
			return Record{}, false, ErrConflict
		}
		return existing.record, true, nil
	}
	store.entries[key] = memoryEntry{fingerprint: fingerprint, record: record}
	return record, false, nil
}

func canonicalSubject(value PromotionSubject) PromotionSubject {
	value.ProposalDigest = ""
	value.AffectedUsers = canonicalStrings(value.AffectedUsers)
	value.FunctionalBehavior = canonicalStrings(value.FunctionalBehavior)
	value.UXExpectations = canonicalStrings(value.UXExpectations)
	value.Accessibility = canonicalStrings(value.Accessibility)
	value.Privacy = canonicalStrings(value.Privacy)
	value.Security = canonicalStrings(value.Security)
	value.Operations = canonicalStrings(value.Operations)
	value.Constraints = canonicalStrings(value.Constraints)
	value.Exclusions = canonicalStrings(value.Exclusions)
	value.AcceptanceCriteria = canonicalStrings(value.AcceptanceCriteria)
	value.DirectUserStatements = canonicalStrings(value.DirectUserStatements)
	value.Assumptions = canonicalStrings(value.Assumptions)
	value.Hypotheses = canonicalStrings(value.Hypotheses)
	value.EvidenceReferences = canonicalStrings(value.EvidenceReferences)
	value.Experiments = canonicalStrings(value.Experiments)
	value.Provenance = canonicalStrings(value.Provenance)
	sort.Slice(value.Questions, func(left, right int) bool { return value.Questions[left].ID < value.Questions[right].ID })
	return value
}

func proposalDigest(value PromotionSubject) string {
	value.ProposalDigest = ""
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func promotionFingerprint(value Record) string {
	value.ID = ""
	value.CreatedAt = time.Time{}
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func canonicalStrings(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return result
}

func validRecord(value Record, subject PromotionSubject, run delivery.DeliveryRun) bool {
	return value.ContractVersion == CanonicalIntentVersion && value.ID != "" &&
		value.ProposalID == subject.ProposalID && value.ProposalDigest == subject.ProposalDigest &&
		value.DeliveryRunID == run.ID() && value.ActorID == subject.ActorID &&
		value.TenantID == subject.TenantID && value.WorkspaceID == subject.WorkspaceID &&
		value.ProjectID == subject.ProjectID && value.RepositoryID == subject.RepositoryID &&
		value.Outcome == subject.Outcome && !value.CreatedAt.IsZero()
}

func invalid(correlationID string) error {
	return failure.New(failure.InvalidRequest, "intent promotion request is invalid", correlationID, failure.RetryNone)
}

func forbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "requirements proposal is not available for promotion", correlationID, failure.RetryNone)
}

func unavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "intent promotion is unavailable", correlationID, failure.RetryBounded)
}
