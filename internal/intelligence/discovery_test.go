package intelligence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/intelligence"
)

var intelligenceTime = time.Date(2026, time.July, 23, 21, 0, 0, 0, time.UTC)

type adapterStub struct {
	calls    int
	proposal discovery.TurnProposal
	err      error
}

func (*adapterStub) Provider() string { return "contract-double" }

func (adapter *adapterStub) Propose(context.Context, intelligence.ProviderActivation, discovery.TurnRequest) (discovery.TurnProposal, error) {
	adapter.calls++
	return adapter.proposal, adapter.err
}

func TestServiceReturnsValidatedNonAuthoritativeProposal(t *testing.T) {
	request := requestFixture()
	adapter := &adapterStub{proposal: proposalFixture(request)}
	service, err := intelligence.NewService(activationFixture(), adapter, func() time.Time { return intelligenceTime })
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := service.Propose(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if adapter.calls != 1 || proposal.Response == "" || proposal.Receipt.Provider != "contract-double" {
		t.Fatalf("unexpected proposal: %#v", proposal)
	}
}

func TestServiceBlocksIneligibleClassificationBeforeAdapter(t *testing.T) {
	request := requestFixture()
	request.Classification = "Restricted"
	adapter := &adapterStub{}
	service, err := intelligence.NewService(activationFixture(), adapter, func() time.Time { return intelligenceTime })
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Propose(context.Background(), request)
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Code != failure.Forbidden || adapter.calls != 0 {
		t.Fatalf("ineligible request reached adapter: calls=%d err=%v", adapter.calls, err)
	}
	request = requestFixture()
	request.Classification = "internal"
	_, err = service.Propose(context.Background(), request)
	if !errors.As(err, &normalized) || normalized.Code != failure.Forbidden || adapter.calls != 0 {
		t.Fatalf("non-canonical classification reached adapter: calls=%d err=%v", adapter.calls, err)
	}
}

func TestServiceRejectsInventedProvenanceAndProviderFailureSafely(t *testing.T) {
	request := requestFixture()
	invalid := proposalFixture(request)
	invalid.ProvenanceRefs = []string{"invented:source"}
	adapter := &adapterStub{proposal: invalid}
	service, err := intelligence.NewService(activationFixture(), adapter, func() time.Time { return intelligenceTime })
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Propose(context.Background(), request)
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Code != failure.Indeterminate {
		t.Fatalf("invalid provider output was not normalized: %v", err)
	}
	adapter.proposal = discovery.TurnProposal{}
	adapter.err = errors.New("provider leaked detail FORGE_TEST_SECRET_DO_NOT_RECORD")
	_, err = service.Propose(context.Background(), request)
	if !errors.As(err, &normalized) || normalized.Code != failure.DependencyUnavailable ||
		contains(err.Error(), "FORGE_TEST_SECRET_DO_NOT_RECORD") {
		t.Fatalf("provider failure was not disclosure-safe: %v", err)
	}
}

func TestServiceEnforcesActivationBudgetsAndReceiptBinding(t *testing.T) {
	request := requestFixture()
	activation := activationFixture()
	activation.Budgets.MaxTokens = request.Limits.MaxTokens - 1
	adapter := &adapterStub{proposal: proposalFixture(request)}
	service, err := intelligence.NewService(activation, adapter, func() time.Time { return intelligenceTime })
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Propose(context.Background(), request)
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Code != failure.BudgetExhausted || adapter.calls != 0 {
		t.Fatalf("activation budget was not enforced before adapter: calls=%d err=%v", adapter.calls, err)
	}

	activation = activationFixture()
	adapter = &adapterStub{proposal: proposalFixture(request)}
	adapter.proposal.Receipt.ModelID = "unapproved-model"
	service, err = intelligence.NewService(activation, adapter, func() time.Time { return intelligenceTime })
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Propose(context.Background(), request)
	if !errors.As(err, &normalized) || normalized.Code != failure.Indeterminate {
		t.Fatalf("unbound provider receipt was accepted: %v", err)
	}
}

func TestServiceEnforcesActivationTurnBudgetBeforeAdapter(t *testing.T) {
	request := requestFixture()
	activation := activationFixture()
	activation.MaxTurns = 1
	adapter := &adapterStub{proposal: proposalFixture(request)}
	service, err := intelligence.NewService(activation, adapter, func() time.Time { return intelligenceTime })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Propose(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	_, err = service.Propose(context.Background(), request)
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Code != failure.BudgetExhausted ||
		adapter.calls != 1 {
		t.Fatalf("turn budget did not stop the adapter: calls=%d err=%v", adapter.calls, err)
	}
}

func TestServiceRejectsActivationPurposeClassificationAndDuplicateModels(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*intelligence.ProviderActivation)
	}{
		{
			name: "unsupported purpose",
			mutate: func(value *intelligence.ProviderActivation) {
				value.AllowedPurposes = []string{"forge.unbounded"}
			},
		},
		{
			name: "unsupported classification",
			mutate: func(value *intelligence.ProviderActivation) {
				value.AllowedClassifications = []string{"Restricted"}
			},
		},
		{
			name: "duplicate models",
			mutate: func(value *intelligence.ProviderActivation) {
				value.AllowedModelIDs = []string{"discovery-double", "discovery-double"}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			activation := activationFixture()
			test.mutate(&activation)
			if _, err := intelligence.NewService(
				activation, &adapterStub{}, func() time.Time { return intelligenceTime },
			); err == nil {
				t.Fatal("ineligible provider activation was accepted")
			}
		})
	}
}

func BenchmarkServiceProposeContractDouble(b *testing.B) {
	request := requestFixture()
	adapter := &adapterStub{proposal: proposalFixture(request)}
	activation := activationFixture()
	activation.MaxTurns = 20
	var service *intelligence.Service
	b.ReportAllocs()
	b.ResetTimer()
	for index := range b.N {
		if index%activation.MaxTurns == 0 {
			var err error
			service, err = intelligence.NewService(activation, adapter, func() time.Time { return intelligenceTime })
			if err != nil {
				b.Fatal(err)
			}
		}
		if _, err := service.Propose(context.Background(), request); err != nil {
			b.Fatal(err)
		}
	}
}

func activationFixture() intelligence.ProviderActivation {
	return intelligence.ProviderActivation{
		SchemaVersion: intelligence.ProviderActivationVersion,
		ID:            "activation-1", Provider: "contract-double", Endpoint: "contract-double://discovery",
		AdapterVersion: "1", AllowedModelIDs: []string{"discovery-double"},
		AllowedPurposes: []string{discovery.PurposeV1}, AllowedClassifications: []string{"Internal"},
		CredentialRef: "none:contract-double", HumanApprovalRef: "approval:activation-1",
		PromptVersion: "prompt.v1", Retention: "none", TrainingUse: "disabled",
		Budgets: discovery.Limits{
			MaxInputBytes: 16 << 10, MaxOutputBytes: 16 << 10, MaxTokens: 4096,
			MaxLatency: 10 * time.Second, MaxCostMicros: 100_000,
		},
		MaxTurns: 8,
		Region:   "local", ActivatedAt: intelligenceTime.Add(-time.Minute),
		ExpiresAt: intelligenceTime.Add(time.Hour),
	}
}

func requestFixture() discovery.TurnRequest {
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
		ExpiresAt: intelligenceTime.Add(time.Minute),
	}
}

func proposalFixture(request discovery.TurnRequest) discovery.TurnProposal {
	return discovery.TurnProposal{
		SchemaVersion: discovery.ProposalVersion, SessionID: request.SessionID,
		CorrelationID: request.CorrelationID, Response: "Who should receive a recommendation?",
		Draft: discovery.Draft{}, Questions: []discovery.Question{{
			ID: "q-1", Text: "Who chooses?", Consequence: "The primary actor is unclear", Material: true,
		}},
		ProvenanceRefs: []string{"interaction:user:1"},
		Usage:          discovery.Usage{InputTokens: 100, OutputTokens: 50, LatencyMs: 200, CostMicros: 100},
		Receipt: discovery.Receipt{
			Provider: "contract-double", ModelID: "discovery-double", AdapterVersion: "1",
			RoleProfileID: discovery.RoleProfileV1, PromptVersion: "prompt.v1",
			ResponseSchema: discovery.ProposalVersion, ProviderRequest: "request-1",
		},
	}
}

func contains(value, substring string) bool {
	for index := 0; index+len(substring) <= len(value); index++ {
		if value[index:index+len(substring)] == substring {
			return true
		}
	}
	return false
}
