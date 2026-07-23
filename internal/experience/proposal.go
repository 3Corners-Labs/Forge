// Package experience owns interaction drafts, requirements proposals, and
// role-appropriate projections. Interaction input never becomes authority.
package experience

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/intent"
	"github.com/3corners/forge/internal/tenant"
)

const RequirementsProposalVersion = "forge.requirements-proposal.v1"

type Session struct {
	actorID        string
	tenantID       string
	workspaceID    string
	projectID      string
	repositoryID   string
	classification string
	correlation    string
}

func NewSession(principal identity.Principal, scope tenant.Scope, binding tenant.RepositoryBinding, at time.Time, correlationID string) (Session, error) {
	if correlationID == "" {
		return Session{}, failure.New(failure.InvalidRequest, "interaction correlation is required", correlationID, failure.RetryNone)
	}
	if err := principal.ValidateAt(at, correlationID); err != nil {
		return Session{}, err
	}
	if principal.Kind() != identity.HumanPrincipal {
		return Session{}, failure.New(failure.Forbidden, "guided interaction requires an authenticated human", correlationID, failure.RetryNone)
	}
	if err := scope.ValidateFor(principal, at, correlationID); err != nil {
		return Session{}, err
	}
	if err := binding.ValidateFor(scope, correlationID); err != nil {
		return Session{}, err
	}
	return Session{
		actorID: principal.SubjectID(), tenantID: scope.TenantID(), workspaceID: scope.WorkspaceID(),
		projectID: scope.ProjectID(), repositoryID: binding.RepositoryID(),
		classification: binding.Classification(), correlation: correlationID,
	}, nil
}

type Statement struct {
	Kind string
	Text string
}

type RequirementsProposal struct {
	ContractVersion    string
	ID                 string
	Version            uint64
	Digest             string
	Lifecycle          string
	Current            bool
	ActorID            string
	TenantID           string
	WorkspaceID        string
	ProjectID          string
	RepositoryID       string
	Outcome            string
	AffectedUsers      []string
	Problem            string
	FunctionalBehavior []string
	UXExpectations     []string
	Accessibility      []string
	Privacy            []string
	Security           []string
	Operations         []string
	Constraints        []string
	Exclusions         []string
	AcceptanceCriteria []string
	Inferences         []string
	Statements         []Statement
	Assumptions        []string
	Hypotheses         []string
	EvidenceReferences []string
	Questions          []intent.Question
	Experiments        []string
	ProductAreaRefs    []string
	Provenance         []string
	PreviousDigest     string
	ChangeExplanation  string
	CreatedAt          time.Time
	CorrelationID      string
}

type CreateProposal struct {
	RawInput           string
	Outcome            string
	AffectedUsers      []string
	Problem            string
	FunctionalBehavior []string
	UXExpectations     []string
	Accessibility      []string
	Privacy            []string
	Security           []string
	Operations         []string
	Constraints        []string
	Exclusions         []string
	AcceptanceCriteria []string
	Inferences         []string
	Assumptions        []string
	Hypotheses         []string
	EvidenceReferences []string
	Questions          []intent.Question
	Experiments        []string
	ProductAreaRefs    []string
	Provenance         []string
	ChangeExplanation  string
}

type Store interface {
	Create(context.Context, RequirementsProposal) error
	Resolve(context.Context, string, string) (RequirementsProposal, bool, error)
}

type Service struct {
	store Store
	now   func() time.Time
	newID func() string
}

func NewService(store Store, now func() time.Time, newID func() string) (*Service, error) {
	if store == nil || now == nil || newID == nil {
		return nil, errors.New("experience proposal configuration is invalid")
	}
	return &Service{store: store, now: now, newID: newID}, nil
}

func (service *Service) Create(ctx context.Context, session Session, request CreateProposal) (RequirementsProposal, error) {
	if service == nil || service.store == nil || service.now == nil || service.newID == nil {
		return RequirementsProposal{}, unavailable(session.correlation)
	}
	if ctx == nil || session.actorID == "" || session.tenantID == "" || session.repositoryID == "" ||
		request.RawInput == "" || request.Outcome == "" || request.Problem == "" ||
		len(request.AffectedUsers) == 0 || len(request.AcceptanceCriteria) == 0 ||
		!allNonEmpty(request.AffectedUsers) || !allNonEmpty(request.AcceptanceCriteria) {
		return RequirementsProposal{}, failure.New(failure.InvalidRequest, "requirements proposal input is invalid", session.correlation, failure.RetryNone)
	}
	at, id := service.now(), service.newID()
	if at.IsZero() || id == "" {
		return RequirementsProposal{}, unavailable(session.correlation)
	}
	proposal := RequirementsProposal{
		ContractVersion: RequirementsProposalVersion, ID: id, Version: 1,
		Lifecycle: "proposal", Current: true, ActorID: session.actorID,
		TenantID: session.tenantID, WorkspaceID: session.workspaceID,
		ProjectID: session.projectID, RepositoryID: session.repositoryID,
		Outcome: request.Outcome, AffectedUsers: slices.Clone(request.AffectedUsers),
		Problem: request.Problem, FunctionalBehavior: slices.Clone(request.FunctionalBehavior),
		UXExpectations: slices.Clone(request.UXExpectations), Accessibility: slices.Clone(request.Accessibility),
		Privacy: slices.Clone(request.Privacy), Security: slices.Clone(request.Security),
		Operations: slices.Clone(request.Operations), Constraints: slices.Clone(request.Constraints),
		Exclusions: slices.Clone(request.Exclusions), AcceptanceCriteria: slices.Clone(request.AcceptanceCriteria),
		Inferences: slices.Clone(request.Inferences),
		Statements: append(
			[]Statement{{Kind: "direct-user-statement", Text: request.RawInput}},
			statements("model-inference", request.Inferences)...,
		),
		Assumptions: slices.Clone(request.Assumptions), Hypotheses: slices.Clone(request.Hypotheses),
		EvidenceReferences: slices.Clone(request.EvidenceReferences), Questions: slices.Clone(request.Questions),
		Experiments: slices.Clone(request.Experiments), ProductAreaRefs: slices.Clone(request.ProductAreaRefs),
		Provenance:        slices.Clone(request.Provenance),
		ChangeExplanation: request.ChangeExplanation, CreatedAt: at, CorrelationID: session.correlation,
	}
	subject := intent.BindProposalDigest(toPromotionSubject(proposal))
	proposal.Digest = subject.ProposalDigest
	if err := service.store.Create(ctx, proposal); err != nil {
		return RequirementsProposal{}, unavailable(session.correlation)
	}
	return cloneProposal(proposal), nil
}

func (service *Service) ResolvePromotionSubject(ctx context.Context, tenantID, proposalID string) (intent.PromotionSubject, bool, error) {
	if service == nil || service.store == nil || ctx == nil || tenantID == "" || proposalID == "" {
		return intent.PromotionSubject{}, false, errors.New("proposal resolution is unavailable")
	}
	proposal, found, err := service.store.Resolve(ctx, tenantID, proposalID)
	if err != nil || !found {
		return intent.PromotionSubject{}, found, err
	}
	subject := toPromotionSubject(proposal)
	if !intent.VerifyProposalDigest(subject) {
		return intent.PromotionSubject{}, false, errors.New("proposal integrity is invalid")
	}
	return subject, true, nil
}

func toPromotionSubject(value RequirementsProposal) intent.PromotionSubject {
	return intent.PromotionSubject{
		ProposalID: value.ID, ProposalVersion: value.Version, ProposalDigest: value.Digest,
		Current: value.Current, ActorID: value.ActorID, TenantID: value.TenantID,
		WorkspaceID: value.WorkspaceID, ProjectID: value.ProjectID, RepositoryID: value.RepositoryID,
		Outcome: value.Outcome, AffectedUsers: slices.Clone(value.AffectedUsers), Problem: value.Problem,
		FunctionalBehavior: slices.Clone(value.FunctionalBehavior), UXExpectations: slices.Clone(value.UXExpectations),
		Accessibility: slices.Clone(value.Accessibility), Privacy: slices.Clone(value.Privacy),
		Security: slices.Clone(value.Security), Operations: slices.Clone(value.Operations),
		Constraints: slices.Clone(value.Constraints), Exclusions: slices.Clone(value.Exclusions),
		AcceptanceCriteria:   slices.Clone(value.AcceptanceCriteria),
		DirectUserStatements: statementTexts(value.Statements, "direct-user-statement"),
		Inferences:           slices.Clone(value.Inferences),
		Assumptions:          slices.Clone(value.Assumptions),
		Hypotheses:           slices.Clone(value.Hypotheses), EvidenceReferences: slices.Clone(value.EvidenceReferences),
		Questions: slices.Clone(value.Questions), Experiments: slices.Clone(value.Experiments),
		ProductAreaRefs: slices.Clone(value.ProductAreaRefs), Provenance: slices.Clone(value.Provenance),
		ChangeExplanation: value.ChangeExplanation,
		CreatedAt:         value.CreatedAt, CorrelationID: value.CorrelationID,
	}
}

func statementTexts(values []Statement, kind string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value.Kind == kind {
			result = append(result, value.Text)
		}
	}
	return result
}

func statements(kind string, values []string) []Statement {
	result := make([]Statement, len(values))
	for index, value := range values {
		result[index] = Statement{Kind: kind, Text: value}
	}
	return result
}

type MemoryStore struct {
	mu      sync.RWMutex
	records map[string]RequirementsProposal
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[string]RequirementsProposal)}
}

func (store *MemoryStore) Create(_ context.Context, value RequirementsProposal) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := value.TenantID + "\x00" + value.ID
	if _, found := store.records[key]; found {
		return errors.New("proposal already exists")
	}
	store.records[key] = cloneProposal(value)
	return nil
}

func (store *MemoryStore) Resolve(_ context.Context, tenantID, proposalID string) (RequirementsProposal, bool, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	value, found := store.records[tenantID+"\x00"+proposalID]
	return cloneProposal(value), found, nil
}

func cloneProposal(value RequirementsProposal) RequirementsProposal {
	value.AffectedUsers = slices.Clone(value.AffectedUsers)
	value.FunctionalBehavior = slices.Clone(value.FunctionalBehavior)
	value.UXExpectations = slices.Clone(value.UXExpectations)
	value.Accessibility = slices.Clone(value.Accessibility)
	value.Privacy = slices.Clone(value.Privacy)
	value.Security = slices.Clone(value.Security)
	value.Operations = slices.Clone(value.Operations)
	value.Constraints = slices.Clone(value.Constraints)
	value.Exclusions = slices.Clone(value.Exclusions)
	value.AcceptanceCriteria = slices.Clone(value.AcceptanceCriteria)
	value.Inferences = slices.Clone(value.Inferences)
	value.Statements = slices.Clone(value.Statements)
	value.Assumptions = slices.Clone(value.Assumptions)
	value.Hypotheses = slices.Clone(value.Hypotheses)
	value.EvidenceReferences = slices.Clone(value.EvidenceReferences)
	value.Questions = slices.Clone(value.Questions)
	value.Experiments = slices.Clone(value.Experiments)
	value.ProductAreaRefs = slices.Clone(value.ProductAreaRefs)
	value.Provenance = slices.Clone(value.Provenance)
	return value
}

func allNonEmpty(values []string) bool {
	for _, value := range values {
		if value == "" {
			return false
		}
	}
	return true
}

func unavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "requirements proposal is unavailable", correlationID, failure.RetryBounded)
}
