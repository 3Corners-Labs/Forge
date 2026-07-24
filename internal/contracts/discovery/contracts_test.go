package discovery_test

import (
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
)

func TestProposalRequiresCompleteDraftWhenNoMaterialQuestionsRemain(t *testing.T) {
	now := time.Date(2026, time.July, 23, 20, 0, 0, 0, time.UTC)
	request := validRequest(now)
	proposal := validProposal(request)
	if err := request.Validate(now); err != nil {
		t.Fatal(err)
	}
	if err := proposal.Validate(request); err != nil {
		t.Fatal(err)
	}
	proposal.Draft.Outcome = ""
	if err := proposal.Validate(request); err == nil {
		t.Fatal("expected an incomplete final draft to fail")
	}
	proposal.Questions = []discovery.Question{{
		ID: "q-1", Text: "What outcome matters?", Consequence: "The outcome remains unclear", Material: true,
	}}
	if err := proposal.Validate(request); err != nil {
		t.Fatalf("a material question should allow an incomplete draft: %v", err)
	}
}

func TestProposalCannotInventProvenanceOrExceedBudget(t *testing.T) {
	now := time.Date(2026, time.July, 23, 20, 0, 0, 0, time.UTC)
	request := validRequest(now)
	proposal := validProposal(request)
	proposal.ProvenanceRefs = []string{"provider:invented"}
	if err := proposal.Validate(request); err == nil {
		t.Fatal("invented provenance was accepted")
	}
	proposal = validProposal(request)
	proposal.Usage.OutputTokens = request.Limits.MaxTokens
	if err := proposal.Validate(request); err == nil {
		t.Fatal("token budget overrun was accepted")
	}
	proposal = validProposal(request)
	proposal.ProvenanceRefs = nil
	if err := proposal.Validate(request); err == nil {
		t.Fatal("proposal without provenance was accepted")
	}
	proposal = validProposal(request)
	proposal.Questions = []discovery.Question{{
		ID: "q-1", Text: "Who chooses?", Material: true,
	}}
	if err := proposal.Validate(request); err == nil {
		t.Fatal("material question without consequence was accepted")
	}
}

func TestRequestRejectsExpiredOrUnsupportedRole(t *testing.T) {
	now := time.Date(2026, time.July, 23, 20, 0, 0, 0, time.UTC)
	request := validRequest(now)
	request.ExpiresAt = now
	if err := request.Validate(now); err == nil {
		t.Fatal("expired request was accepted")
	}
	request = validRequest(now)
	request.RoleProfileID = "forge.role.unbounded"
	if err := request.Validate(now); err == nil {
		t.Fatal("unsupported role was accepted")
	}
}

func validRequest(now time.Time) discovery.TurnRequest {
	return discovery.TurnRequest{
		SchemaVersion: discovery.RequestVersion, SessionID: "session-1",
		CorrelationID: "correlation-1", TenantID: "tenant-1", WorkspaceID: "workspace-1",
		ProjectID: "project-1", HumanSubjectID: "human-1", TaskAgentID: "agent-1",
		RoleProfileID: discovery.RoleProfileV1, Purpose: discovery.PurposeV1,
		Classification: "Internal", UserStatements: []string{"Help employees find a desk"},
		ProvenanceRefs: []string{"interaction:user:1"},
		Limits: discovery.Limits{
			MaxInputBytes: 16 << 10, MaxOutputBytes: 16 << 10, MaxTokens: 4096,
			MaxLatency: 10 * time.Second, MaxCostMicros: 100_000,
		},
		ExpiresAt: now.Add(time.Minute),
	}
}

func validProposal(request discovery.TurnRequest) discovery.TurnProposal {
	return discovery.TurnProposal{
		SchemaVersion: discovery.ProposalVersion, SessionID: request.SessionID,
		CorrelationID: request.CorrelationID, Response: "I have enough to form a proposal.",
		Draft: discovery.Draft{
			Outcome: "Find a desk quickly", AffectedUsers: []string{"employee"},
			Problem: "Desk choice is slow", Constraints: []string{"advice only"},
			AcceptanceCriteria: []string{"recommendation within three seconds"},
			Inferences:         []string{"The employee chooses their own desk"},
		},
		ProvenanceRefs: []string{"interaction:user:1"},
		Usage:          discovery.Usage{InputTokens: 120, OutputTokens: 80, LatencyMs: 500, CostMicros: 200},
		Receipt: discovery.Receipt{
			ActivationID:     "activation-1",
			ActivationDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Provider:         "contract-double", ModelID: "discovery-double", AdapterVersion: "1",
			RoleProfileID: discovery.RoleProfileV1, PromptVersion: "prompt.v1",
			ResponseSchema: discovery.ProposalVersion, ProviderRequest: "request-1",
		},
	}
}
