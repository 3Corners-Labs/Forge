package experience_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/experience"
	"github.com/3corners/forge/internal/identity"
)

type discoveryAdvisorStub struct {
	calls    int
	request  discovery.TurnRequest
	proposal discovery.TurnProposal
}

func (advisor *discoveryAdvisorStub) Propose(_ context.Context, request discovery.TurnRequest) (discovery.TurnProposal, error) {
	advisor.calls++
	advisor.request = request
	advisor.proposal.SchemaVersion = discovery.ProposalVersion
	advisor.proposal.SessionID = request.SessionID
	advisor.proposal.CorrelationID = request.CorrelationID
	advisor.proposal.ProvenanceRefs = []string{"interaction:user:1"}
	advisor.proposal.Receipt.RoleProfileID = request.RoleProfileID
	advisor.proposal.Receipt.ResponseSchema = discovery.ProposalVersion
	return advisor.proposal, nil
}

func TestDiscoveryTurnUsesResolvedScopeAndReturnsProposalOnly(t *testing.T) {
	session := proposalSession(t, "human-1", "correlation-discovery")
	taskAgent := agentPrincipal(t)
	advisor := &discoveryAdvisorStub{proposal: discovery.TurnProposal{
		Response: "Who is the primary user?",
		Questions: []discovery.Question{{
			ID: "q-1", Text: "Who is the primary user?", Consequence: "Actor remains unclear", Material: true,
		}},
		Usage: discovery.Usage{InputTokens: 10, OutputTokens: 10, LatencyMs: 10},
		Receipt: discovery.Receipt{
			Provider: "contract-double", ModelID: "discovery-double", AdapterVersion: "1",
			PromptVersion: "prompt.v1", ProviderRequest: "provider-request-1",
		},
	}}
	service, err := experience.NewDiscoveryService(
		advisor, func() time.Time { return proposalTime }, func() string { return "discovery-session-1" },
		discovery.Limits{
			MaxInputBytes: 16 << 10, MaxOutputBytes: 16 << 10, MaxTokens: 4096,
			MaxLatency: time.Second, MaxCostMicros: 1000,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := service.Turn(context.Background(), session, taskAgent, experience.DiscoveryTurn{
		UserStatements: []string{"Build something useful"},
		ProvenanceRefs: []string{"interaction:user:1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if advisor.calls != 1 || advisor.request.TenantID != "tenant-fixture" ||
		advisor.request.WorkspaceID != "workspace-fixture" ||
		advisor.request.HumanSubjectID != "human-1" ||
		proposal.SessionID != "discovery-session-1" {
		t.Fatalf("unexpected request=%#v proposal=%#v", advisor.request, proposal)
	}
}

func TestDiscoveryRejectsWorkloadAgentBeforeAdvisor(t *testing.T) {
	session := proposalSession(t, "human-1", "correlation-discovery")
	advisor := &discoveryAdvisorStub{}
	service, err := experience.NewDiscoveryService(
		advisor, func() time.Time { return proposalTime }, func() string { return "session-1" },
		discovery.Limits{
			MaxInputBytes: 1, MaxOutputBytes: 1, MaxTokens: 1,
			MaxLatency: time.Second, MaxCostMicros: 0,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Turn(context.Background(), session, identity.Principal{}, experience.DiscoveryTurn{
		UserStatements: []string{"input"}, ProvenanceRefs: []string{"interaction:user:1"},
	})
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Code != failure.InvalidRequest || advisor.calls != 0 {
		t.Fatalf("invalid agent reached advisor: calls=%d err=%v", advisor.calls, err)
	}
}

func agentPrincipal(t *testing.T) identity.Principal {
	t.Helper()
	agentVerifier, err := identity.NewVerifier(agentIdentityProvider{}, "forge-cli", func() time.Time { return proposalTime })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := agentVerifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "synthetic",
		CorrelationID: "correlation-discovery",
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

type agentIdentityProvider struct{}

func (agentIdentityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: "agent-1", Kind: identity.AgentPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: proposalTime.Add(-time.Minute), ExpiresAt: proposalTime.Add(time.Hour), Synthetic: true,
	}, nil
}
