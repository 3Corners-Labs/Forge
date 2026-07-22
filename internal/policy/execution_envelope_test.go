package policy_test

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/policy"
)

func TestIssueExecutionEnvelopeBindsExactAuthorityAndReplays(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	first, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	if err != nil {
		t.Fatal(err)
	}
	if first.Replayed {
		t.Fatal("first issuance was reported as replayed")
	}
	value := first.Envelope
	if err := value.VerifyIntegrity(fixture.request.CorrelationID); err != nil {
		t.Fatal(err)
	}
	if value.Scope.TenantID != fixture.request.Scope.TenantID() || value.Scope.WorkspaceID != fixture.request.Scope.WorkspaceID() ||
		value.Scope.ProjectID != fixture.request.Scope.ProjectID() || value.Scope.DeliveryRunID != fixture.request.Run.ID() ||
		value.Principal.SubjectID != fixture.request.Human.SubjectID() || value.Principal.AgentID != fixture.request.Agent.SubjectID() ||
		!slices.Equal(value.Principal.DelegationChain, []string{fixture.request.DelegationID}) ||
		!slices.Equal(value.Actions, []string{identity.DelegatedActionAssembleContext}) ||
		!slices.Equal(value.Classifications, []string{fixture.request.Run.Classification()}) ||
		len(value.Destinations) != 0 || len(value.ModelIDs) != 0 || len(value.ToolIDs) != 0 || len(value.SecretReferenceIDs) != 0 {
		t.Fatalf("envelope widened or changed server-derived authority: %#v", value)
	}
	if !slices.Contains(value.PolicyDecisionIDs, fixture.request.PolicyDecisionID) ||
		!value.ExpiresAt.Equal(testNow.Add(5*time.Minute)) || value.Limits.MaxDurationSeconds != 300 ||
		value.Limits.MaxCostMinorUnits != 0 || value.Limits.MaxConcurrency != 1 {
		t.Fatalf("unexpected decision, expiry, or limits: %#v", value)
	}

	replayed, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	if err != nil || !replayed.Replayed || replayed.Envelope.ID != value.ID {
		t.Fatalf("replay = %#v, %v", replayed, err)
	}
	projected := replayed.Envelope.Actions
	projected[0] = "forge.release"
	again, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	if err != nil || again.Envelope.Actions[0] != identity.DelegatedActionAssembleContext {
		t.Fatal("caller mutated stored execution authority")
	}

	mutated := value
	mutated.Scope.ProjectID = "other-project"
	assertFailureCode(t, mutated.VerifyIntegrity(fixture.request.CorrelationID), failure.VerificationFailed)
}

func TestIssueExecutionEnvelopeFailsClosedBeforeMinting(t *testing.T) {
	base := newEnvelopeFixture(t)
	tests := []struct {
		name   string
		ctx    context.Context
		mutate func(*policy.IssueExecutionEnvelopeRequest)
		want   failure.Code
	}{
		{"nil context", nil, func(*policy.IssueExecutionEnvelopeRequest) {}, failure.InvalidRequest},
		{"contract", context.Background(), func(value *policy.IssueExecutionEnvelopeRequest) { value.ContractVersion = "other" }, failure.InvalidRequest},
		{"missing idempotency", context.Background(), func(value *policy.IssueExecutionEnvelopeRequest) { value.IdempotencyKey = "" }, failure.InvalidRequest},
		{"missing decision", context.Background(), func(value *policy.IssueExecutionEnvelopeRequest) { value.PolicyDecisionID = "" }, failure.InvalidRequest},
		{"forged decision reference", context.Background(), func(value *policy.IssueExecutionEnvelopeRequest) { value.PolicyDecisionID = "forged" }, failure.Forbidden},
		{"human", context.Background(), func(value *policy.IssueExecutionEnvelopeRequest) { value.Human = identity.Principal{} }, failure.Unauthenticated},
		{"agent", context.Background(), func(value *policy.IssueExecutionEnvelopeRequest) { value.Agent = identity.Principal{} }, failure.Unauthenticated},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newEnvelopeFixture(t)
			request := base.request
			test.mutate(&request)
			_, err := fixture.issuer.IssueExecutionEnvelope(test.ctx, request)
			assertFailureCode(t, err, test.want)
		})
	}
}

func TestIssueExecutionEnvelopeRequiresLiveDelegationAndObligations(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	if err := fixture.delegations.RevokeDelegation(context.Background(), identity.RevokeDelegationRequest{
		ContractVersion: identity.RevokeDelegationContractVersion, CorrelationID: fixture.request.CorrelationID,
		Human: fixture.request.Human, DelegationID: fixture.request.DelegationID, Reason: "test revocation",
	}); err != nil {
		t.Fatal(err)
	}
	_, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	assertFailureCode(t, err, failure.Forbidden)

	fixture = newEnvelopeFixture(t)
	profile := validEnvelopeProfile()
	profile.EnforcedObligationTypes = nil
	issuer, err := policy.NewEnvelopeIssuer(fixture.evaluator, fixture.delegations, policy.NewMemoryEnvelopeStore(), profile, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	_, err = issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	assertFailureCode(t, err, failure.Forbidden)
}

func TestIssueExecutionEnvelopeRejectsDenyIndeterminateAndExpiredDecisions(t *testing.T) {
	for _, test := range []struct {
		name   string
		effect policy.Effect
		match  bool
		want   failure.Code
	}{
		{"deny", policy.EffectDeny, true, failure.Forbidden},
		{"indeterminate", policy.EffectPermit, false, failure.Indeterminate},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newEnvelopeFixture(t)
			project := "other-project"
			if test.match {
				project = fixture.request.Scope.ProjectID()
			}
			set := policy.PolicySet{Version: "policy.outcome.v1", Rules: []policy.Rule{{
				ID: "outcome", Effect: test.effect, Actions: []policy.Action{policy.ActionAssembleContext},
				ProjectIDs: []string{project},
			}}}
			evaluator, err := policy.NewEvaluator(set, nil, 5*time.Minute, func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			decision, err := evaluator.Evaluate(context.Background(), policy.EvaluateRequest{
				ContractVersion: policy.EvaluateRequestContractVersion, CorrelationID: fixture.request.CorrelationID,
				Principal: fixture.request.Human, Scope: fixture.request.Scope, Run: fixture.request.Run,
				Action: policy.ActionAssembleContext, Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
			})
			if err != nil {
				t.Fatal(err)
			}
			profile := validEnvelopeProfile()
			profile.AllowedPolicyVersions = []string{"policy.outcome.v1"}
			profile.EnforcedObligationTypes = nil
			issuer, err := policy.NewEnvelopeIssuer(evaluator, fixture.delegations, policy.NewMemoryEnvelopeStore(), profile, func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			request := fixture.request
			request.PolicyDecisionID = decision.ID
			_, err = issuer.IssueExecutionEnvelope(context.Background(), request)
			assertFailureCode(t, err, test.want)
		})
	}

	fixture := newEnvelopeFixture(t)
	issuer, err := policy.NewEnvelopeIssuer(
		fixture.evaluator, fixture.delegations, policy.NewMemoryEnvelopeStore(), validEnvelopeProfile(),
		func() time.Time { return testNow.Add(6 * time.Minute) },
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	assertFailureCode(t, err, failure.Expired)
}

func TestEnvelopeIdempotencyConflictsOnChangedTrustedLimits(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	if _, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request); err != nil {
		t.Fatal(err)
	}
	profile := validEnvelopeProfile()
	profile.Limits.MaxDurationSeconds = 120
	issuer, err := policy.NewEnvelopeIssuer(fixture.evaluator, fixture.delegations, fixture.store, profile, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	_, err = issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	assertFailureCode(t, err, failure.Conflict)
}

func TestEnvelopeIdentityChangesWithAuthorityRelevantLimits(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	first, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
	if err != nil {
		t.Fatal(err)
	}
	profile := validEnvelopeProfile()
	profile.Limits.MaxDurationSeconds = 120
	issuer, err := policy.NewEnvelopeIssuer(fixture.evaluator, fixture.delegations, policy.NewMemoryEnvelopeStore(), profile, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	request := fixture.request
	request.IdempotencyKey = "different-authority"
	second, err := issuer.IssueExecutionEnvelope(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if first.Envelope.ID == second.Envelope.ID || first.Envelope.Integrity.Digest == second.Envelope.Integrity.Digest {
		t.Fatal("authority-relevant limit change did not change envelope identity")
	}
}

func TestEnvelopeIssuerNormalizesStoreFailures(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	for _, test := range []struct {
		name string
		err  error
		want failure.Code
	}{
		{"unavailable", policy.ErrEnvelopeStoreUnavailable, failure.DependencyUnavailable},
		{"unknown", errors.New("sensitive persistence diagnostic"), failure.Indeterminate},
	} {
		t.Run(test.name, func(t *testing.T) {
			issuer, err := policy.NewEnvelopeIssuer(fixture.evaluator, fixture.delegations, failingEnvelopeStore{err: test.err}, validEnvelopeProfile(), func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			_, err = issuer.IssueExecutionEnvelope(context.Background(), fixture.request)
			assertFailureCode(t, err, test.want)
			if errors.Is(err, test.err) {
				t.Fatalf("adapter diagnostic escaped boundary: %v", err)
			}
		})
	}
}

func TestPolicyDecisionLedgerIsCapabilityOwnedAndDefensive(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	decision, found, err := fixture.evaluator.ResolvePolicyDecision(context.Background(), fixture.request.PolicyDecisionID)
	if err != nil || !found {
		t.Fatalf("resolve = %#v, %t, %v", decision, found, err)
	}
	decision.PolicyVersions[0] = "mutated"
	decision.Obligations[0]["required"] = false
	again, found, err := fixture.evaluator.ResolvePolicyDecision(context.Background(), fixture.request.PolicyDecisionID)
	if err != nil || !found || again.PolicyVersions[0] != "policy.m1.v1" || again.Obligations[0]["required"] != true {
		t.Fatalf("caller mutated Policy-owned decision: %#v, %v", again, err)
	}
	if _, found, err := fixture.evaluator.ResolvePolicyDecision(context.Background(), "unknown"); err != nil || found {
		t.Fatalf("unknown decision disclosed: found=%t err=%v", found, err)
	}
}

func TestNewEnvelopeIssuerRejectsUnboundedProfile(t *testing.T) {
	fixture := newEnvelopeFixture(t)
	tests := []struct {
		name   string
		mutate func(*policy.EnvelopeProfile)
	}{
		{"issuer", func(value *policy.EnvelopeProfile) { value.Issuer = "" }},
		{"environment", func(value *policy.EnvelopeProfile) { value.Environment = "production" }},
		{"risk", func(value *policy.EnvelopeProfile) { value.Risk = policy.RiskCritical }},
		{"policy versions", func(value *policy.EnvelopeProfile) { value.AllowedPolicyVersions = nil }},
		{"cost", func(value *policy.EnvelopeProfile) { value.Limits.MaxCostMinorUnits = 1 }},
		{"tokens", func(value *policy.EnvelopeProfile) { value.Limits.MaxInputTokens = 1 }},
		{"concurrency", func(value *policy.EnvelopeProfile) { value.Limits.MaxConcurrency = 2 }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := validEnvelopeProfile()
			test.mutate(&profile)
			if _, err := policy.NewEnvelopeIssuer(fixture.evaluator, fixture.delegations, policy.NewMemoryEnvelopeStore(), profile, func() time.Time { return testNow }); err == nil {
				t.Fatal("unbounded profile was accepted")
			}
		})
	}
}

func BenchmarkIssueExecutionEnvelope(b *testing.B) {
	b.Run("Create", func(b *testing.B) {
		fixture := newEnvelopeFixture(b)
		for index := 0; index < b.N; index++ {
			request := fixture.request
			request.IdempotencyKey = "envelope-benchmark-" + strconv.Itoa(index)
			if _, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), request); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Replay", func(b *testing.B) {
		fixture := newEnvelopeFixture(b)
		if _, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request); err != nil {
			b.Fatal(err)
		}
		for index := 0; index < b.N; index++ {
			if _, err := fixture.issuer.IssueExecutionEnvelope(context.Background(), fixture.request); err != nil {
				b.Fatal(err)
			}
		}
	})
}

type envelopeFixture struct {
	evaluator   *policy.Evaluator
	delegations *identity.DelegationService
	store       *policy.MemoryEnvelopeStore
	issuer      *policy.EnvelopeIssuer
	request     policy.IssueExecutionEnvelopeRequest
}

func newEnvelopeFixture(t testing.TB) envelopeFixture {
	t.Helper()
	base := validRequest(t)
	agent := verifiedPolicyPrincipal(t, "synthetic-task-agent", identity.AgentPrincipal)
	delegationEvaluation := policy.EvaluateDelegationRequest{
		ContractVersion: policy.EvaluateDelegationRequestContractVersion, CorrelationID: base.CorrelationID,
		Human: base.Principal, Agent: agent, Scope: base.Scope, Run: base.Run,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	}
	delegationEvaluation.Delegation = delegationRequestForEvaluation(
		t, &delegationEvaluation, base.Principal.SubjectID(), agent.SubjectID(), base.Run.RepositoryID(), base.Run.Classification(),
	)
	set := permitSet()
	set.Rules = append(set.Rules, delegationPermitSet(delegationEvaluation.Delegation.Digest()).Rules[0])
	evaluator := newEvaluator(t, set, []string{"evidence-required", "independent-verification"})
	delegationDecision, err := evaluator.EvaluateDelegation(context.Background(), delegationEvaluation)
	if err != nil {
		t.Fatal(err)
	}
	delegations, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	created, err := delegations.CreateDelegation(context.Background(), identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: delegationEvaluation.Delegation,
		Human: base.Principal, Agent: agent, ParentDecision: delegationDecision,
		EnforcedObligationTypes: []string{"evidence-required"},
	})
	if err != nil {
		t.Fatal(err)
	}
	contextDecision, err := evaluator.Evaluate(context.Background(), base)
	if err != nil {
		t.Fatal(err)
	}
	store := policy.NewMemoryEnvelopeStore()
	issuer, err := policy.NewEnvelopeIssuer(evaluator, delegations, store, validEnvelopeProfile(), func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	return envelopeFixture{
		evaluator: evaluator, delegations: delegations, store: store, issuer: issuer,
		request: policy.IssueExecutionEnvelopeRequest{
			ContractVersion: policy.IssueExecutionEnvelopeContractVersion,
			CorrelationID:   base.CorrelationID, IdempotencyKey: "envelope-1",
			Human: base.Principal, Agent: agent, Scope: base.Scope, Run: base.Run,
			DelegationID: created.Delegation.ID(), PolicyDecisionID: contextDecision.ID,
		},
	}
}

func validEnvelopeProfile() policy.EnvelopeProfile {
	return policy.EnvelopeProfile{
		Issuer: "forge.policy.developer-trial", Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
		AllowedPolicyVersions: []string{"policy.m1.v1"}, EnforcedObligationTypes: []string{"independent-verification"},
		MaximumLifetime: 10 * time.Minute,
		Limits:          trust.ExecutionLimits{MaxDurationSeconds: 600, MaxCostMinorUnits: 0, MaxConcurrency: 1},
	}
}

type failingEnvelopeStore struct{ err error }

func (store failingEnvelopeStore) Create(context.Context, trust.ExecutionEnvelope, string, string, string) (trust.ExecutionEnvelope, bool, error) {
	return trust.ExecutionEnvelope{}, false, store.err
}
