package policy_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/policy"
	"github.com/3corners/forge/internal/tenant"
)

var testNow = time.Date(2026, time.July, 22, 14, 30, 0, 0, time.UTC)

type identityProvider struct {
	assertion identity.Assertion
}

func (provider identityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return provider.assertion, nil
}

type bindingSource struct{}

func (bindingSource) ResolveRepositoryBinding(context.Context, string, string, string, string) (tenant.RepositoryBindingDefinition, bool, error) {
	return tenant.RepositoryBindingDefinition{
		TenantID: "tenant-fixture", WorkspaceID: "workspace-fixture", ProjectID: "project-forge",
		RepositoryID: "repository-forge", Selector: "forge", CanonicalLocator: "https://example.invalid/forge.git",
		Classification: "Internal",
	}, true, nil
}

func TestEvaluatePermitIsDeterministicAndContractValid(t *testing.T) {
	request := validRequest(t)
	evaluator := newEvaluator(t, policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{{
		ID: "permit-context", Effect: policy.EffectPermit,
		Actions:      []policy.Action{policy.ActionAssembleContext},
		Environments: []policy.Environment{policy.DeveloperTrialEnvironment},
		Risks:        []policy.Risk{policy.RiskLow},
		SubjectIDs:   []string{"synthetic-developer-allowed"}, ResourceIDs: []string{"repository-forge"},
		DeliveryRunIDs: []string{"delivery-run-policy-1"}, Purposes: []string{delivery.M1ChangePurpose},
		Classifications: []string{"Internal"},
		TenantIDs:       []string{"tenant-fixture"}, WorkspaceIDs: []string{"workspace-fixture"},
		ProjectIDs: []string{"project-forge"},
		Obligations: []policy.Obligation{
			{Type: "evidence-required", Parameters: map[string]any{"classification": "Internal"}},
			{Type: "independent-verification", Parameters: map[string]any{"required": true}},
		},
	}}}, []string{"independent-verification", "evidence-required"})

	first, err := evaluator.Evaluate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	second, err := evaluator.Evaluate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID {
		t.Fatalf("decision IDs differ: %q != %q", first.ID, second.ID)
	}
	if !strings.HasPrefix(first.ID, "forge.policy-decision.") || first.Outcome != "permit" || first.DecisionPoint != "retrieval" {
		t.Fatalf("unexpected decision: %#v", first)
	}
	if first.SubjectID != request.Principal.SubjectID() || first.Resource != request.Run.RepositoryID() || first.Action != string(request.Action) {
		t.Fatalf("decision is not bound to exact immutable facts: %#v", first)
	}
	if first.ExpiresAt == nil || !first.ExpiresAt.Equal(testNow.Add(5*time.Minute)) {
		t.Fatalf("expiry = %v, want %v", first.ExpiresAt, testNow.Add(5*time.Minute))
	}
	if got := first.Obligations[0]["type"]; got != "evidence-required" {
		t.Fatalf("first obligation type = %v", got)
	}
	if err := first.Authorize(
		"retrieval", request.Principal.SubjectID(), string(request.Action), request.Run.RepositoryID(),
		testNow, []string{"evidence-required", "independent-verification"}, request.CorrelationID,
	); err != nil {
		t.Fatalf("authorize exact decision: %v", err)
	}
}

func TestPrincipalExpiryNarrowsDecision(t *testing.T) {
	request := validRequestWithExpiry(t, testNow.Add(2*time.Minute))
	evaluator := newEvaluator(t, permitSet(), []string{"independent-verification"})
	decision, err := evaluator.Evaluate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.ExpiresAt == nil || !decision.ExpiresAt.Equal(request.Principal.ExpiresAt()) {
		t.Fatalf("decision expiry = %v, principal expiry = %v", decision.ExpiresAt, request.Principal.ExpiresAt())
	}
}

func TestEvaluatorDefensivelyCopiesPolicyDefinitions(t *testing.T) {
	request := validRequest(t)
	set := permitSet()
	evaluator := newEvaluator(t, set, []string{"independent-verification"})
	set.Rules[0].Effect = policy.EffectDeny
	set.Rules[0].ProjectIDs[0] = "mutated-project"
	set.Rules[0].Obligations[0].Type = "mutated-obligation"
	set.Rules[0].Obligations[0].Parameters["required"] = false

	decision, err := evaluator.Evaluate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "permit" || decision.Obligations[0]["type"] != "independent-verification" || decision.Obligations[0]["required"] != true {
		t.Fatalf("policy constructor input mutated active evaluation: %#v", decision)
	}
}

func TestDenyPrecedenceAndUnknownObligationFailClosed(t *testing.T) {
	request := validRequest(t)
	tests := []struct {
		name   string
		set    policy.PolicySet
		reason string
	}{
		{
			name: "mandatory deny overrides specific permit",
			set: policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{
				{ID: "global-deny", Effect: policy.EffectDeny, Mandatory: true, Actions: []policy.Action{policy.ActionAssembleContext}},
				{ID: "project-permit", Effect: policy.EffectPermit, ProjectIDs: []string{"project-forge"}},
			}},
			reason: "mandatory-policy-deny",
		},
		{
			name: "ordinary deny overrides permit",
			set: policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{
				{ID: "deny", Effect: policy.EffectDeny},
				{ID: "permit", Effect: policy.EffectPermit},
			}},
			reason: "policy-deny",
		},
		{
			name: "unknown obligation denies",
			set: policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{{
				ID: "permit", Effect: policy.EffectPermit,
				Obligations: []policy.Obligation{{Type: "unimplemented-control"}},
			}}},
			reason: "unknown-mandatory-obligation",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := newEvaluator(t, test.set, nil)
			decision, err := evaluator.Evaluate(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			if decision.Outcome != "deny" || len(decision.ReasonCodes) != 1 || decision.ReasonCodes[0] != test.reason {
				t.Fatalf("unexpected decision: %#v", decision)
			}
			assertFailureCode(t, decision.Authorize(
				"retrieval", decision.SubjectID, decision.Action, decision.Resource, testNow, nil, request.CorrelationID,
			), failure.Forbidden)
		})
	}
}

func TestNoMatchAndObligationConflictAreIndeterminate(t *testing.T) {
	request := validRequest(t)
	tests := []struct {
		name   string
		set    policy.PolicySet
		reason string
	}{
		{
			name: "no applicable rule",
			set: policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{{
				ID: "other-project", Effect: policy.EffectPermit, ProjectIDs: []string{"another-project"},
			}}},
			reason: "no-applicable-policy",
		},
		{
			name: "conflicting obligation parameters",
			set: policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{
				{ID: "permit-a", Effect: policy.EffectPermit, Obligations: []policy.Obligation{{Type: "independent-verification", Parameters: map[string]any{"runner": "a"}}}},
				{ID: "permit-b", Effect: policy.EffectPermit, Obligations: []policy.Obligation{{Type: "independent-verification", Parameters: map[string]any{"runner": "b"}}}},
			}},
			reason: "policy-obligation-conflict",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluator := newEvaluator(t, test.set, []string{"independent-verification"})
			decision, err := evaluator.Evaluate(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			if decision.Outcome != "indeterminate" || decision.ReasonCodes[0] != test.reason {
				t.Fatalf("unexpected decision: %#v", decision)
			}
			assertFailureCode(t, decision.Authorize(
				"retrieval", decision.SubjectID, decision.Action, decision.Resource, testNow, nil, request.CorrelationID,
			), failure.Indeterminate)
		})
	}
}

func TestEveryVariablePolicyMatchDimensionIsFailClosed(t *testing.T) {
	request := validRequest(t)
	tests := []struct {
		name   string
		mutate func(*policy.Rule)
	}{
		{"risk", func(rule *policy.Rule) { rule.Risks = []policy.Risk{policy.RiskHigh} }},
		{"subject", func(rule *policy.Rule) { rule.SubjectIDs = []string{"other-subject"} }},
		{"resource", func(rule *policy.Rule) { rule.ResourceIDs = []string{"other-repository"} }},
		{"delivery run", func(rule *policy.Rule) { rule.DeliveryRunIDs = []string{"other-run"} }},
		{"purpose", func(rule *policy.Rule) { rule.Purposes = []string{"other-purpose"} }},
		{"classification", func(rule *policy.Rule) { rule.Classifications = []string{"Restricted"} }},
		{"tenant", func(rule *policy.Rule) { rule.TenantIDs = []string{"other-tenant"} }},
		{"workspace", func(rule *policy.Rule) { rule.WorkspaceIDs = []string{"other-workspace"} }},
		{"project", func(rule *policy.Rule) { rule.ProjectIDs = []string{"other-project"} }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := permitSet()
			test.mutate(&set.Rules[0])
			decision, err := newEvaluator(t, set, []string{"independent-verification"}).Evaluate(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			if decision.Outcome != "indeterminate" || decision.ReasonCodes[0] != "no-applicable-policy" {
				t.Fatalf("unexpected decision: %#v", decision)
			}
		})
	}
}

func TestEvaluateRejectsInvalidOrUnboundFactsBeforeDecision(t *testing.T) {
	base := validRequest(t)
	tests := []struct {
		name   string
		ctx    context.Context
		mutate func(*policy.EvaluateRequest)
		want   failure.Code
	}{
		{"nil context", nil, func(*policy.EvaluateRequest) {}, failure.InvalidRequest},
		{"contract", context.Background(), func(value *policy.EvaluateRequest) { value.ContractVersion = "other" }, failure.InvalidRequest},
		{"correlation", context.Background(), func(value *policy.EvaluateRequest) { value.CorrelationID = "" }, failure.InvalidRequest},
		{"principal", context.Background(), func(value *policy.EvaluateRequest) { value.Principal = identity.Principal{} }, failure.Unauthenticated},
		{"scope", context.Background(), func(value *policy.EvaluateRequest) { value.Scope = tenant.Scope{} }, failure.Forbidden},
		{"run", context.Background(), func(value *policy.EvaluateRequest) { value.Run = delivery.DeliveryRun{} }, failure.InvalidRequest},
		{"action", context.Background(), func(value *policy.EvaluateRequest) { value.Action = "forge.release" }, failure.InvalidRequest},
		{"environment", context.Background(), func(value *policy.EvaluateRequest) { value.Environment = "production" }, failure.InvalidRequest},
		{"risk", context.Background(), func(value *policy.EvaluateRequest) { value.Risk = "unknown" }, failure.InvalidRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := base
			test.mutate(&request)
			evaluator := newEvaluator(t, permitSet(), []string{"independent-verification"})
			_, err := evaluator.Evaluate(test.ctx, request)
			assertFailureCode(t, err, test.want)
		})
	}
}

func TestUnavailableEvaluatorAndClockNormalizeFailure(t *testing.T) {
	request := validRequest(t)
	var unavailable *policy.Evaluator
	_, err := unavailable.Evaluate(context.Background(), request)
	assertFailureCode(t, err, failure.DependencyUnavailable)

	evaluator, err := policy.NewEvaluator(permitSet(), []string{"independent-verification"}, 5*time.Minute, func() time.Time { return time.Time{} })
	if err != nil {
		t.Fatal(err)
	}
	_, err = evaluator.Evaluate(context.Background(), request)
	assertFailureCode(t, err, failure.DependencyUnavailable)

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = newEvaluator(t, permitSet(), []string{"independent-verification"}).Evaluate(canceled, request)
	assertFailureCode(t, err, failure.DependencyUnavailable)
}

func TestNewEvaluatorRejectsMalformedPolicyDefinitions(t *testing.T) {
	tests := []struct {
		name string
		set  policy.PolicySet
	}{
		{name: "version", set: policy.PolicySet{Rules: []policy.Rule{{ID: "rule", Effect: policy.EffectPermit}}}},
		{name: "rules", set: policy.PolicySet{Version: "v1"}},
		{name: "duplicate IDs", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{
			{ID: "same", Effect: policy.EffectPermit}, {ID: "same", Effect: policy.EffectDeny},
		}}},
		{name: "effect", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{ID: "rule", Effect: "allow"}}}},
		{name: "action", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit, Actions: []policy.Action{"unknown"},
		}}}},
		{name: "environment", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit, Environments: []policy.Environment{"production"},
		}}}},
		{name: "risk", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit, Risks: []policy.Risk{"unknown"},
		}}}},
		{name: "duplicate typed match", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit,
			Actions: []policy.Action{policy.ActionAssembleContext, policy.ActionAssembleContext},
		}}}},
		{name: "classification", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit, Classifications: []string{"Unknown"},
		}}}},
		{name: "empty match", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit, TenantIDs: []string{""},
		}}}},
		{name: "duplicate obligations", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit,
			Obligations: []policy.Obligation{{Type: "one"}, {Type: "one"}},
		}}}},
		{name: "reserved obligation field", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectPermit,
			Obligations: []policy.Obligation{{Type: "one", Parameters: map[string]any{"type": "two"}}},
		}}}},
		{name: "deny obligation", set: policy.PolicySet{Version: "v1", Rules: []policy.Rule{{
			ID: "rule", Effect: policy.EffectDeny, Obligations: []policy.Obligation{{Type: "one"}},
		}}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := policy.NewEvaluator(test.set, nil, 5*time.Minute, func() time.Time { return testNow }); err == nil {
				t.Fatal("expected constructor error")
			}
		})
	}
	if _, err := policy.NewEvaluator(permitSet(), []string{"one", "one"}, 5*time.Minute, func() time.Time { return testNow }); err == nil {
		t.Fatal("expected duplicate supported obligation error")
	}
}

func permitSet() policy.PolicySet {
	return policy.PolicySet{Version: "policy.m1.v1", Rules: []policy.Rule{{
		ID: "permit-context", Effect: policy.EffectPermit,
		Actions:      []policy.Action{policy.ActionAssembleContext},
		Environments: []policy.Environment{policy.DeveloperTrialEnvironment},
		Risks:        []policy.Risk{policy.RiskLow},
		SubjectIDs:   []string{"synthetic-developer-allowed"}, ResourceIDs: []string{"repository-forge"},
		DeliveryRunIDs: []string{"delivery-run-policy-1"}, Purposes: []string{delivery.M1ChangePurpose},
		Classifications: []string{"Internal"},
		TenantIDs:       []string{"tenant-fixture"}, WorkspaceIDs: []string{"workspace-fixture"},
		ProjectIDs:  []string{"project-forge"},
		Obligations: []policy.Obligation{{Type: "independent-verification", Parameters: map[string]any{"required": true}}},
	}}}
}

func newEvaluator(t *testing.T, set policy.PolicySet, supported []string) *policy.Evaluator {
	t.Helper()
	evaluator, err := policy.NewEvaluator(set, supported, 5*time.Minute, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	return evaluator
}

func validRequest(t *testing.T) policy.EvaluateRequest {
	return validRequestWithExpiry(t, testNow.Add(30*time.Minute))
}

func validRequestWithExpiry(t *testing.T, expiresAt time.Time) policy.EvaluateRequest {
	t.Helper()
	provider := identityProvider{assertion: identity.Assertion{
		SubjectID: "synthetic-developer-allowed", Kind: identity.HumanPrincipal,
		Issuer: "fixture-issuer", Audience: "forge-developer-trial", Assurance: "mfa",
		ProfileID: "developer-trial", IssuedAt: testNow.Add(-time.Minute), ExpiresAt: expiresAt, Synthetic: true,
	}}
	verifier, err := identity.NewVerifier(provider, "forge-developer-trial", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque-fixture", CorrelationID: "correlation-policy-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: principal.SubjectID(), WorkspaceID: "workspace-fixture"}},
		bindingSource{}, func() time.Time { return testNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	scope, err := directory.ResolveWorkspaceScope(principal, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: "correlation-policy-1",
		WorkspaceSelector: "fixture-workspace", ProjectSelector: "forge", RepositorySelector: "forge",
	})
	if err != nil {
		t.Fatal(err)
	}
	binding, err := directory.GetRepositoryBinding(context.Background(), tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: "correlation-policy-1", Scope: scope,
	})
	if err != nil {
		t.Fatal(err)
	}
	store, err := delivery.NewMemoryStore(func() string { return "delivery-run-policy-1" })
	if err != nil {
		t.Fatal(err)
	}
	service, err := delivery.NewService(store, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.CaptureIntent(context.Background(), delivery.CaptureIntentRequest{
		ContractVersion: delivery.CaptureIntentRequestContractVersion,
		CorrelationID:   "correlation-policy-1", IdempotencyKey: "intent-policy-1",
		Principal: principal, DelegationChain: []string{}, Scope: scope, RepositoryBinding: binding,
		Objective: "Prepare one bounded Forge change", Constraints: []string{"Developer Trial only"},
		AcceptanceCriteria: []string{"Policy decision is deterministic"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return policy.EvaluateRequest{
		ContractVersion: policy.EvaluateRequestContractVersion, CorrelationID: "correlation-policy-1",
		Principal: principal, Scope: scope, Run: result.Run,
		Action: policy.ActionAssembleContext, Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	}
}

func assertFailureCode(t *testing.T, err error, want failure.Code) {
	t.Helper()
	var normalized *failure.Error
	if !errors.As(err, &normalized) {
		t.Fatalf("error %T %v is not normalized", err, err)
	}
	if normalized.Code != want {
		t.Fatalf("failure code = %q, want %q", normalized.Code, want)
	}
	if validationErr := normalized.Validate(); validationErr != nil {
		t.Fatalf("invalid normalized error: %v", validationErr)
	}
}
