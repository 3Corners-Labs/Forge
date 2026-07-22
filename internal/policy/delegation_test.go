package policy_test

import (
	"context"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/policy"
)

func TestEvaluateDelegationBindsExactRequestAndSupportsIdentityIssuance(t *testing.T) {
	request := validDelegationEvaluation(t)
	set := delegationPermitSet(request.Delegation.Digest())
	evaluator := newEvaluator(t, set, []string{"evidence-required"})

	first, err := evaluator.EvaluateDelegation(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	second, err := evaluator.EvaluateDelegation(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID || first.Outcome != "permit" || first.DecisionPoint != "execution" {
		t.Fatalf("unexpected delegation decision: %#v", first)
	}
	if first.SubjectID != request.Human.SubjectID() || first.Action != string(policy.ActionCreateDelegation) || first.Resource != request.Delegation.Digest() {
		t.Fatalf("decision is not exactly digest-bound: %#v", first)
	}
	if first.ExpiresAt == nil || !first.ExpiresAt.Equal(testNow.Add(5*time.Minute)) {
		t.Fatalf("decision expiry = %v", first.ExpiresAt)
	}
	service, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.CreateDelegation(context.Background(), identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: request.Delegation,
		Human: request.Human, Agent: request.Agent, ParentDecision: first,
		EnforcedObligationTypes: []string{"evidence-required"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Delegation.ParentAuthorityID() != first.ID || result.Delegation.RequestDigest() != request.Delegation.Digest() {
		t.Fatalf("Identity did not retain exact Policy authority: %#v", result)
	}
}

func TestEvaluateDelegationIsFailClosedAcrossBindingsAndPolicyDimensions(t *testing.T) {
	base := validDelegationEvaluation(t)
	tests := []struct {
		name   string
		mutate func(*policy.EvaluateDelegationRequest, *policy.PolicySet)
		want   failure.Code
		result string
	}{
		{"wrong human binding", func(request *policy.EvaluateDelegationRequest, _ *policy.PolicySet) {
			request.Delegation = delegationRequestForEvaluation(t, request, "other-human", request.Agent.SubjectID(), request.Run.RepositoryID(), request.Run.Classification())
		}, failure.Forbidden, ""},
		{"wrong agent binding", func(request *policy.EvaluateDelegationRequest, _ *policy.PolicySet) {
			request.Delegation = delegationRequestForEvaluation(t, request, request.Human.SubjectID(), "other-agent", request.Run.RepositoryID(), request.Run.Classification())
		}, failure.Forbidden, ""},
		{"cross repository", func(request *policy.EvaluateDelegationRequest, _ *policy.PolicySet) {
			request.Delegation = delegationRequestForEvaluation(t, request, request.Human.SubjectID(), request.Agent.SubjectID(), "other-repository", request.Run.Classification())
		}, failure.Forbidden, ""},
		{"wrong classification", func(request *policy.EvaluateDelegationRequest, _ *policy.PolicySet) {
			request.Delegation = delegationRequestForEvaluation(t, request, request.Human.SubjectID(), request.Agent.SubjectID(), request.Run.RepositoryID(), "Restricted")
		}, failure.Forbidden, ""},
		{"wrong environment", func(request *policy.EvaluateDelegationRequest, _ *policy.PolicySet) {
			request.Environment = "production"
		}, failure.InvalidRequest, ""},
		{"wrong evidence correlation", func(request *policy.EvaluateDelegationRequest, _ *policy.PolicySet) {
			request.CorrelationID = "other-correlation"
		}, failure.InvalidRequest, ""},
		{"wrong risk rule", func(_ *policy.EvaluateDelegationRequest, set *policy.PolicySet) {
			set.Rules[0].Risks = []policy.Risk{policy.RiskCritical}
		}, "", "indeterminate"},
		{"wrong digest rule", func(_ *policy.EvaluateDelegationRequest, set *policy.PolicySet) {
			set.Rules[0].ResourceIDs = []string{"other-digest"}
		}, "", "indeterminate"},
		{"deny precedence", func(_ *policy.EvaluateDelegationRequest, set *policy.PolicySet) {
			set.Rules = append(set.Rules, policy.Rule{ID: "mandatory-deny", Effect: policy.EffectDeny, Mandatory: true, Actions: []policy.Action{policy.ActionCreateDelegation}})
		}, "", "deny"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := base
			set := delegationPermitSet(base.Delegation.Digest())
			test.mutate(&request, &set)
			decision, err := newEvaluator(t, set, []string{"evidence-required"}).EvaluateDelegation(context.Background(), request)
			if test.want != "" {
				assertFailureCode(t, err, test.want)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if decision.Outcome != test.result {
				t.Fatalf("outcome = %q, want %q", decision.Outcome, test.result)
			}
		})
	}
}

func TestEvaluateDelegationRejectsUnverifiedOrBroadenedInputs(t *testing.T) {
	base := validDelegationEvaluation(t)
	tests := []struct {
		name   string
		ctx    context.Context
		mutate func(*policy.EvaluateDelegationRequest)
		want   failure.Code
	}{
		{"nil context", nil, func(*policy.EvaluateDelegationRequest) {}, failure.InvalidRequest},
		{"contract", context.Background(), func(value *policy.EvaluateDelegationRequest) { value.ContractVersion = "other" }, failure.InvalidRequest},
		{"human", context.Background(), func(value *policy.EvaluateDelegationRequest) { value.Human = identity.Principal{} }, failure.Unauthenticated},
		{"agent", context.Background(), func(value *policy.EvaluateDelegationRequest) { value.Agent = identity.Principal{} }, failure.Unauthenticated},
		{"request", context.Background(), func(value *policy.EvaluateDelegationRequest) { value.Delegation = identity.DelegationRequest{} }, failure.InvalidRequest},
		{"risk", context.Background(), func(value *policy.EvaluateDelegationRequest) { value.Risk = "unknown" }, failure.InvalidRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := base
			test.mutate(&request)
			_, err := newEvaluator(t, delegationPermitSet(base.Delegation.Digest()), []string{"evidence-required"}).EvaluateDelegation(test.ctx, request)
			assertFailureCode(t, err, test.want)
		})
	}
}

func BenchmarkEvaluateDelegation(b *testing.B) {
	request := validDelegationEvaluation(b)
	evaluator := newEvaluator(b, delegationPermitSet(request.Delegation.Digest()), []string{"evidence-required"})
	for index := 0; index < b.N; index++ {
		if _, err := evaluator.EvaluateDelegation(context.Background(), request); err != nil {
			b.Fatal(err)
		}
	}
}

func validDelegationEvaluation(t testing.TB) policy.EvaluateDelegationRequest {
	t.Helper()
	contextRequest := validRequest(t)
	agent := verifiedPolicyPrincipal(t, "synthetic-task-agent", identity.AgentPrincipal)
	request := policy.EvaluateDelegationRequest{
		ContractVersion: policy.EvaluateDelegationRequestContractVersion,
		CorrelationID:   contextRequest.CorrelationID,
		Human:           contextRequest.Principal,
		Agent:           agent,
		Scope:           contextRequest.Scope,
		Run:             contextRequest.Run,
		Environment:     policy.DeveloperTrialEnvironment,
		Risk:            policy.RiskLow,
	}
	request.Delegation = delegationRequestForEvaluation(
		t, &request, request.Human.SubjectID(), request.Agent.SubjectID(), request.Run.RepositoryID(), request.Run.Classification(),
	)
	return request
}

func delegationRequestForEvaluation(t testing.TB, request *policy.EvaluateDelegationRequest, human, agent, repository, classification string) identity.DelegationRequest {
	t.Helper()
	value, err := identity.NewDelegationRequest(identity.DelegationRequestInput{
		ContractVersion: identity.DelegationRequestContractVersion,
		IdempotencyKey:  "delegation-policy-1", CorrelationID: request.CorrelationID,
		HumanSubjectID: human, AgentSubjectID: agent,
		TenantID: request.Scope.TenantID(), WorkspaceID: request.Scope.WorkspaceID(), ProjectID: request.Scope.ProjectID(),
		RepositoryID: repository, DeliveryRunID: request.Run.ID(), Purpose: identity.M1ChangePurpose,
		Actions: []string{identity.DelegatedActionAssembleContext}, Resources: []string{repository},
		Classifications: []string{classification}, Destinations: []string{},
		Environment: identity.DeveloperTrialEnvironment, MaximumAutonomy: identity.ContextPreparationOnly,
		RequestedExpiresAt: testNow.Add(20 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func verifiedPolicyPrincipal(t testing.TB, subject string, kind identity.PrincipalKind) identity.Principal {
	t.Helper()
	provider := identityProvider{assertion: identity.Assertion{
		SubjectID: subject, Kind: kind, Issuer: "fixture-issuer", Audience: "forge-developer-trial",
		Assurance: "synthetic-local", ProfileID: "developer-trial", IssuedAt: testNow.Add(-time.Minute),
		ExpiresAt: testNow.Add(30 * time.Minute), Synthetic: true,
	}}
	verifier, err := identity.NewVerifier(provider, "forge-developer-trial", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque", CorrelationID: "correlation-policy-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

func delegationPermitSet(requestDigest string) policy.PolicySet {
	return policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{{
		ID: "permit-root-delegation", Effect: policy.EffectPermit,
		Actions: []policy.Action{policy.ActionCreateDelegation}, Environments: []policy.Environment{policy.DeveloperTrialEnvironment},
		Risks: []policy.Risk{policy.RiskLow}, SubjectIDs: []string{"synthetic-developer-allowed"},
		ResourceIDs: []string{requestDigest}, DeliveryRunIDs: []string{"delivery-run-policy-1"},
		Purposes: []string{delivery.M1ChangePurpose}, Classifications: []string{"Internal"},
		TenantIDs: []string{"tenant-fixture"}, WorkspaceIDs: []string{"workspace-fixture"}, ProjectIDs: []string{"project-forge"},
		Obligations: []policy.Obligation{{Type: "evidence-required", Parameters: map[string]any{"classification": "Internal"}}},
	}}}
}
