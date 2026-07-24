package intelligence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	changecontract "github.com/3corners/forge/internal/contracts/change"
	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
)

type changeAdapterDouble struct {
	crossScope bool
}

type failingChangeAdapter struct {
	calls int
}

func (changeAdapterDouble) Provider() string   { return "contract-double" }
func (*failingChangeAdapter) Provider() string { return "contract-double" }

func (adapter *failingChangeAdapter) Prepare(context.Context, ProviderActivation, changecontract.Request) (changecontract.Proposal, error) {
	adapter.calls++
	return changecontract.Proposal{}, errors.New("provider unavailable")
}

func (adapter changeAdapterDouble) Prepare(_ context.Context, activation ProviderActivation, request changecontract.Request) (changecontract.Proposal, error) {
	proposal, err := changecontract.BindProposal(changecontract.Proposal{
		SchemaVersion: changecontract.ProposalVersion,
		RequestID:     request.ID, RequestDigest: request.Digest, SourceRevision: request.SourceRevision,
		PlanDigest: request.PlanDigest, ContextDigest: request.ContextManifestDigest,
		Replacements: []changecontract.FileReplacement{{
			Path: request.AllowedWritePaths[0], Content: "package example\n", Rationale: "Satisfies AC-1.",
		}},
		ProposedChecks: []string{"go test ./internal/example"},
		Receipt: changecontract.Receipt{
			ActivationID: activation.ID, ActivationDigest: activation.Digest,
			Provider: activation.Provider, ModelID: activation.AllowedModelIDs[0],
			AdapterVersion: activation.AdapterVersion, RoleProfileID: request.RoleProfileID,
			PromptVersion: activation.PromptVersion, ResponseSchema: changecontract.ProposalVersion,
			ProviderRequest: "request-1",
		},
	}, request)
	if err == nil && adapter.crossScope {
		proposal.Replacements[0].Path = ".github/workflows/escape.yml"
	}
	return proposal, err
}

func TestChangeServiceRequiresExactPolicyAuthorityAndRejectsCrossScopeOutput(t *testing.T) {
	now := time.Date(2026, time.July, 25, 12, 0, 0, 0, time.UTC)
	request := intelligenceChangeRequest(t, now)
	activation := intelligenceChangeActivation(now)
	decision, envelope := intelligenceChangeAuthority(t, now, request, activation)
	service, err := NewChangeService(activation, changeAdapterDouble{}, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := service.Prepare(context.Background(), ChangePreparationCommand{
		Request: request, Decision: decision, Envelope: envelope,
	})
	if err != nil || proposal.Digest == "" {
		t.Fatalf("proposal=%#v err=%v", proposal, err)
	}

	wrong := decision
	wrong.Resource = strings.Repeat("x", 71)
	if _, err := service.Prepare(context.Background(), ChangePreparationCommand{
		Request: request, Decision: wrong, Envelope: envelope,
	}); failureCode(err) != failure.Forbidden {
		t.Fatalf("wrong authority err=%v", err)
	}
	broadened := envelope
	broadened.Destinations = append(broadened.Destinations, "model:other")
	if _, err := service.Prepare(context.Background(), ChangePreparationCommand{
		Request: request, Decision: decision, Envelope: broadened,
	}); failureCode(err) != failure.Forbidden {
		t.Fatalf("broadened authority err=%v", err)
	}

	crossScope, err := NewChangeService(activation, changeAdapterDouble{crossScope: true}, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := crossScope.Prepare(context.Background(), ChangePreparationCommand{
		Request: request, Decision: decision, Envelope: envelope,
	}); failureCode(err) != failure.Indeterminate {
		t.Fatalf("cross-scope output err=%v", err)
	}
}

func TestChangeServiceCountsFailedProviderInvocationAgainstActivationBudget(t *testing.T) {
	now := time.Date(2026, time.July, 25, 12, 0, 0, 0, time.UTC)
	request := intelligenceChangeRequest(t, now)
	activation := intelligenceChangeActivation(now)
	activation.MaxTurns = 1
	decision, envelope := intelligenceChangeAuthority(t, now, request, activation)
	adapter := &failingChangeAdapter{}
	service, err := NewChangeService(activation, adapter, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	command := ChangePreparationCommand{Request: request, Decision: decision, Envelope: envelope}
	if _, err := service.Prepare(context.Background(), command); failureCode(err) != failure.DependencyUnavailable {
		t.Fatalf("first provider failure err=%v", err)
	}
	if _, err := service.Prepare(context.Background(), command); failureCode(err) != failure.BudgetExhausted {
		t.Fatalf("failed invocation did not consume the bounded turn: %v", err)
	}
	if adapter.calls != 1 {
		t.Fatalf("adapter calls = %d, want 1", adapter.calls)
	}
}

func intelligenceChangeRequest(t *testing.T, now time.Time) changecontract.Request {
	t.Helper()
	request, err := changecontract.BindRequest(changecontract.Request{
		SchemaVersion: changecontract.RequestVersion, CorrelationID: "correlation-1",
		TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1",
		RepositoryID: "repository-1", DeliveryRunID: "run-1", TaskAgentID: "agent-1",
		PlanID: "plan-1", PlanDigest: testDigest("plan"), ContextPackageID: "context-1",
		ContextManifestDigest: testDigest("context"), RequirementsRecordID: "requirements-1",
		SourceRevision: strings.Repeat("a", 40), RoleProfileID: changecontract.RoleProfileV1,
		Purpose: changecontract.PurposeV1, Classification: "Internal",
		AllowedReadPaths: []string{"internal/example"}, AllowedWritePaths: []string{"internal/example/example.go"},
		AcceptanceCriteria: []string{"AC-1"}, TestIntent: []string{"unit test"},
		ControlReferences: []string{"ADR-0023"},
		Limits: changecontract.Limits{
			MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20, MaxTokens: 1000,
			MaxLatency: time.Second, MaxCostMicros: 1000, MaxChangedFiles: 1,
		},
		ExpiresAt: now.Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func intelligenceChangeActivation(now time.Time) ProviderActivation {
	return ProviderActivation{
		SchemaVersion: ProviderActivationVersion, ID: "activation-1", Digest: testDigest("activation"),
		Provider: "contract-double", Endpoint: "contract-double://change",
		AdapterVersion: "adapter-1", AllowedModelIDs: []string{"model-1"},
		AllowedPurposes:        []string{changecontract.PurposeV1},
		AllowedClassifications: []string{"Internal"}, CredentialRef: "none:test",
		HumanApprovalRef: "approval:test", PromptVersion: "prompt-1",
		Retention: "none", TrainingUse: "disabled", Region: "local",
		ReasoningEffort: "low", Pricing: Pricing{Currency: "USD"},
		Budgets: discovery.Limits{
			MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20, MaxTokens: 1000,
			MaxLatency: time.Second, MaxCostMicros: 1000,
		},
		MaxTurns: 3, MaxSessionTokens: 3000, MaxSessionCostMicros: 3000,
		MaxProjectTokensPerDay: 3000, MaxProjectCostPerDay: 3000,
		MaxConcurrentRequests: 1, ActivatedAt: now.Add(-time.Minute), ExpiresAt: now.Add(time.Hour),
	}
}

func intelligenceChangeAuthority(t *testing.T, now time.Time, request changecontract.Request, activation ProviderActivation) (trust.PolicyDecision, trust.ExecutionEnvelope) {
	t.Helper()
	expires := now.Add(30 * time.Minute)
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, ID: "decision-1",
		DecisionPoint: "intelligence", SubjectID: request.TaskAgentID,
		Action: ProposeChangeAction, Resource: request.Digest, EvaluatedAt: now.Add(-time.Second),
		ExpiresAt: &expires, PolicyVersions: []string{"policy-v1"}, Outcome: "permit",
		Obligations: []map[string]any{
			{"type": "model-output-non-authoritative"},
			{"type": "purpose-minimized-context"},
		},
		EvidenceCorrelationID: request.CorrelationID,
	}
	envelope, err := trust.BindExecutionEnvelopeIntegrity(trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: "forge.policy",
		IssuedAt: now.Add(-time.Second), ExpiresAt: expires,
		Scope: trust.ExecutionScope{
			TenantID: request.TenantID, WorkspaceID: request.WorkspaceID,
			ProjectID: request.ProjectID, DeliveryRunID: request.DeliveryRunID,
		},
		Principal: trust.EnvelopePrincipal{
			SubjectID: "human-1", AgentID: request.TaskAgentID,
			DelegationChain: []string{"delegation-1"},
		},
		Purpose: request.Purpose, RiskClass: "low", Actions: []string{ProposeChangeAction},
		Classifications: []string{request.Classification},
		Destinations:    []string{"model:" + activation.Provider}, ModelIDs: activation.AllowedModelIDs,
		SecretReferenceIDs: []string{activation.CredentialRef}, PolicyDecisionIDs: []string{decision.ID},
		Limits: trust.ExecutionLimits{
			MaxDurationSeconds: 60, MaxCostMinorUnits: request.Limits.MaxCostMicros,
			MaxInputTokens: request.Limits.MaxTokens, MaxOutputTokens: request.Limits.MaxTokens,
			MaxConcurrency: 1,
		},
		IdempotencyKey: "change-1", EvidenceCorrelationID: request.CorrelationID,
	}, request.CorrelationID)
	if err != nil {
		t.Fatal(err)
	}
	return decision, envelope
}

func testDigest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func failureCode(err error) failure.Code {
	var normalized *failure.Error
	if !errors.As(err, &normalized) {
		return ""
	}
	return normalized.Code
}
