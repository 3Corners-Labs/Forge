package policy_test

import (
	"context"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/policy"
	"github.com/3corners/forge/internal/tenant"
)

func TestEvaluateApprovalBindsExactPlanReviewFacts(t *testing.T) {
	base := validRequest(t)
	run := plannedPolicyRun(t, base)
	request := approvalPolicyRequest(t, run)
	evaluator := newEvaluator(t, policy.PolicySet{Version: "policy.approval.v1", Rules: []policy.Rule{{
		ID: "permit-plan-review", Effect: policy.EffectPermit,
		Actions:      []policy.Action{policy.ActionDecideApproval},
		Environments: []policy.Environment{policy.DeveloperTrialEnvironment}, Risks: []policy.Risk{policy.RiskLow},
		SubjectIDs: []string{base.Principal.SubjectID()}, ResourceIDs: []string{request.RequestDigest},
		DeliveryRunIDs: []string{run.ID()}, Purposes: []string{run.Purpose()},
		Classifications: []string{run.Classification()}, TenantIDs: []string{run.TenantID()},
		WorkspaceIDs: []string{run.WorkspaceID()}, ProjectIDs: []string{run.ProjectID()},
		Obligations: []policy.Obligation{{Type: "human-interaction"}},
	}}}, []string{"human-interaction"})
	command := policy.EvaluateApprovalRequest{
		ContractVersion: policy.EvaluateApprovalRequestContractVersion, CorrelationID: run.CorrelationID(),
		Principal: base.Principal, Scope: base.Scope, Run: run, Request: request,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	}
	decision, err := evaluator.EvaluateApproval(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "permit" || decision.DecisionPoint != "approval" ||
		decision.Action != "forge.approval.decide" || decision.Resource != request.RequestDigest {
		t.Fatalf("unexpected eligibility decision: %#v", decision)
	}
	if err := decision.Authorize("approval", base.Principal.SubjectID(), "forge.approval.decide",
		request.RequestDigest, testNow, []string{"human-interaction"}, run.CorrelationID()); err != nil {
		t.Fatal(err)
	}

	crossScope := request
	crossScope.ID, crossScope.RequestDigest = "", ""
	crossScope.Scope.ProjectID = "other-project"
	crossScope, err = trust.BindApprovalRequestDigest(crossScope, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	command.Request = crossScope
	_, err = evaluator.EvaluateApproval(context.Background(), command)
	assertFailureCode(t, err, failure.Forbidden)
}

func plannedPolicyRun(t *testing.T, base policy.EvaluateRequest) delivery.DeliveryRun {
	t.Helper()
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: base.Principal.SubjectID(), WorkspaceID: "workspace-fixture"}},
		bindingSource{}, func() time.Time { return testNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	binding, err := directory.GetRepositoryBinding(context.Background(), tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: base.CorrelationID, Scope: base.Scope,
	})
	if err != nil {
		t.Fatal(err)
	}
	store, err := delivery.NewMemoryStore(func() string { return "approval-run-1" })
	if err != nil {
		t.Fatal(err)
	}
	service, err := delivery.NewService(store, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	captured, err := service.CaptureIntent(context.Background(), delivery.CaptureIntentRequest{
		ContractVersion: delivery.CaptureIntentRequestContractVersion, CorrelationID: base.CorrelationID,
		IdempotencyKey: "approval-intent", Principal: base.Principal, DelegationChain: []string{},
		Scope: base.Scope, RepositoryBinding: binding, Objective: "Review one exact plan",
		Constraints: []string{"Developer Trial only"}, AcceptanceCriteria: []string{"Plan approval is digest bound"},
	})
	if err != nil {
		t.Fatal(err)
	}
	authorizing, _, err := service.BeginContextAuthorization(context.Background(), delivery.TransitionRequest{
		CorrelationID: base.CorrelationID, IdempotencyKey: "approval-context-authorizing",
		Run: captured.Run, ExpectedVersion: captured.Run.Version(),
	})
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := trust.BindExecutionEnvelopeIntegrity(trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: "forge.policy",
		IssuedAt: testNow.Add(-time.Minute), ExpiresAt: testNow.Add(time.Hour),
		Scope:     trust.ExecutionScope{TenantID: authorizing.TenantID(), WorkspaceID: authorizing.WorkspaceID(), ProjectID: authorizing.ProjectID(), DeliveryRunID: authorizing.ID()},
		Principal: trust.EnvelopePrincipal{SubjectID: authorizing.ActorID(), AgentID: "agent-1", DelegationChain: []string{"delegation-1"}},
		Purpose:   authorizing.Purpose(), RiskClass: "low", Actions: []string{"forge.context.assemble"},
		Classifications: []string{authorizing.Classification()}, Destinations: []string{}, ModelIDs: []string{},
		ToolIDs: []string{}, SecretReferenceIDs: []string{}, PolicyDecisionIDs: []string{"context-decision"},
		ApprovalRecordIDs: []string{}, Limits: trust.ExecutionLimits{MaxDurationSeconds: 600, MaxConcurrency: 1},
		EvidenceCorrelationID: base.CorrelationID,
	}, base.CorrelationID)
	if err != nil {
		t.Fatal(err)
	}
	expires := testNow.Add(30 * time.Minute)
	contextPackage, err := trust.BindContextPackageManifest(trust.ContextPackage{
		SchemaVersion: trust.ContextPackageSchemaVersion, ExecutionEnvelopeID: envelope.ID,
		Objective: authorizing.Objective(), AssembledAt: testNow, ExpiresAt: &expires,
		PolicyDecisionIDs: []string{"context-decision"}, Items: []trust.ContextItem{{
			SourceID: "architecture", SourceRevision: "revision-1", Locator: "docs/architecture",
			Classification: authorizing.Classification(), Authority: "accepted-knowledge",
			Purpose: authorizing.Purpose(), ContentDigest: "sha256:content",
		}}, Omissions: []string{}, Conflicts: []string{},
	}, base.CorrelationID)
	if err != nil {
		t.Fatal(err)
	}
	ready, _, err := service.RecordContext(context.Background(), delivery.RecordContextRequest{
		TransitionRequest: delivery.TransitionRequest{
			CorrelationID: base.CorrelationID, IdempotencyKey: "approval-context-ready",
			Run: authorizing, ExpectedVersion: authorizing.Version(),
		}, Package: contextPackage, Envelope: envelope,
	})
	if err != nil {
		t.Fatal(err)
	}
	planProfile := delivery.PlanProfile{
		Version:             "forge.m1-plan-profile.v1",
		Steps:               []delivery.PlanStep{{ID: "verify", Dependencies: []string{}, OutputKinds: []string{"verification-result"}}},
		AllowedCapabilities: []string{"forge.capability.delivery"},
		AllowedPathPrefixes: []string{"internal"}, ChangePaths: []string{"internal/delivery"},
		Assumptions: []string{"accepted architecture remains unchanged"},
		Risks:       []string{"in-process trial state"}, Limitations: []string{"no worker execution"},
		StopConditions:     []string{"architecture decision required"},
		Budgets:            delivery.PlanBudgets{MaxDurationSeconds: 600, MaxConcurrency: 1, MaxAttempts: 1},
		VerificationIntent: []string{"go test ./..."}, GuardrailReferences: []string{"ADR-0015"},
		Rollback: "Discard branch.", ManualRecovery: "Return to ContextReady.", RiskClass: "low",
		MaximumLifetime: 20 * time.Minute,
	}
	profiles, err := delivery.NewStaticPlanProfiles([]delivery.PlanProfile{planProfile})
	if err != nil {
		t.Fatal(err)
	}
	planned, _, err := service.RecordPlan(context.Background(), profiles, delivery.RecordPlanRequest{
		TransitionRequest: delivery.TransitionRequest{
			CorrelationID: base.CorrelationID, IdempotencyKey: "approval-planned",
			Run: ready, ExpectedVersion: ready.Version(),
		},
		ProfileVersion: planProfile.Version,
	})
	if err != nil {
		t.Fatal(err)
	}
	return planned
}

func approvalPolicyRequest(t *testing.T, run delivery.DeliveryRun) trust.ApprovalRequest {
	t.Helper()
	plan := run.Plan()
	value, err := trust.BindApprovalRequestDigest(trust.ApprovalRequest{
		SchemaVersion: trust.ApprovalRequestSchemaVersion, Issuer: "forge.approval",
		RequestingActorID: run.ActorID(), IntendedApproverID: run.ActorID(),
		Scope:       trust.ApprovalScope{TenantID: run.TenantID(), WorkspaceID: run.WorkspaceID(), ProjectID: run.ProjectID(), RepositoryID: run.RepositoryID(), DeliveryRunID: run.ID()},
		Environment: "developer-trial", Purpose: run.Purpose(), Classification: run.Classification(),
		RiskClass: plan.RiskClass, Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan",
		SubjectID: plan.ID, SubjectDigest: plan.Digest,
		EvidenceReferences:  []string{"context:" + plan.ContextPackageID, "plan:" + plan.Digest},
		PermittedConditions: []string{"require-independent-verification"},
		PolicyDecisionIDs:   []string{"context-decision"}, CreatedAt: testNow,
		ExpiresAt: testNow.Add(20 * time.Minute), MaxDecisionLifetime: 1200,
		Lifecycle: "open", EvidenceCorrelationID: run.CorrelationID(),
	}, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	return value
}
