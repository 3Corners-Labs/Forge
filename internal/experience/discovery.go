package experience

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/identity"
)

type DiscoveryAdvisor interface {
	Propose(context.Context, discovery.TurnRequest) (discovery.TurnProposal, error)
}

type DiscoveryTurn struct {
	SessionID         string
	UserStatements    []string
	CurrentDraft      discovery.Draft
	PreviousQuestions []discovery.Question
	ProvenanceRefs    []string
}

type DiscoveryService struct {
	advisor DiscoveryAdvisor
	now     func() time.Time
	newID   func() string
	limits  discovery.Limits
}

func NewDiscoveryService(advisor DiscoveryAdvisor, now func() time.Time, newID func() string, limits discovery.Limits) (*DiscoveryService, error) {
	if advisor == nil || now == nil || newID == nil || limits.MaxInputBytes <= 0 ||
		limits.MaxOutputBytes <= 0 || limits.MaxTokens <= 0 || limits.MaxLatency <= 0 ||
		limits.MaxCostMicros < 0 {
		return nil, errors.New("experience discovery configuration is invalid")
	}
	return &DiscoveryService{advisor: advisor, now: now, newID: newID, limits: limits}, nil
}

func (service *DiscoveryService) Turn(ctx context.Context, session Session, taskAgent identity.Principal, turn DiscoveryTurn) (discovery.TurnProposal, error) {
	if service == nil || service.advisor == nil || service.now == nil || service.newID == nil ||
		ctx == nil || session.actorID == "" || session.tenantID == "" || session.repositoryID == "" ||
		taskAgent.SubjectID() == "" || taskAgent.Kind() != identity.AgentPrincipal ||
		len(turn.UserStatements) == 0 || len(turn.ProvenanceRefs) == 0 {
		return discovery.TurnProposal{}, failure.New(
			failure.InvalidRequest, "guided discovery turn is invalid",
			session.correlation, failure.RetryNone,
		)
	}
	now := service.now()
	if err := taskAgent.ValidateAt(now, session.correlation); err != nil {
		return discovery.TurnProposal{}, err
	}
	sessionID := turn.SessionID
	if sessionID == "" {
		sessionID = service.newID()
	}
	if sessionID == "" {
		return discovery.TurnProposal{}, failure.New(
			failure.DependencyUnavailable, "guided discovery session is unavailable",
			session.correlation, failure.RetryBounded,
		)
	}
	request := discovery.TurnRequest{
		SchemaVersion: discovery.RequestVersion, SessionID: sessionID,
		CorrelationID: session.correlation, TenantID: session.tenantID,
		WorkspaceID: session.workspaceID, ProjectID: session.projectID,
		HumanSubjectID: session.actorID, TaskAgentID: taskAgent.SubjectID(),
		RoleProfileID: discovery.RoleProfileV1, Purpose: discovery.PurposeV1,
		Classification: session.classification,
		UserStatements: slices.Clone(turn.UserStatements), CurrentDraft: turn.CurrentDraft.Clone(),
		PreviousQuestions: slices.Clone(turn.PreviousQuestions),
		ProvenanceRefs:    slices.Clone(turn.ProvenanceRefs), Limits: service.limits,
		ExpiresAt: now.Add(service.limits.MaxLatency + time.Minute),
	}
	if err := request.Validate(now); err != nil {
		return discovery.TurnProposal{}, failure.New(
			failure.InvalidRequest, "guided discovery request is invalid",
			session.correlation, failure.RetryNone,
		)
	}
	proposal, err := service.advisor.Propose(ctx, request)
	if err != nil {
		return discovery.TurnProposal{}, err
	}
	if err := proposal.Validate(request); err != nil {
		return discovery.TurnProposal{}, failure.New(
			failure.Indeterminate, "guided discovery proposal is invalid",
			session.correlation, failure.RetryBounded,
		)
	}
	return proposal, nil
}
