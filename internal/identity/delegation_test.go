package identity_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/identity"
)

func TestDelegationRequestCanonicalDigestAndDefensiveProjections(t *testing.T) {
	input := validDelegationInput()
	first, err := identity.NewDelegationRequest(input)
	if err != nil {
		t.Fatal(err)
	}
	input.IdempotencyKey = "replay-metadata-changed"
	input.CorrelationID = "correlation-metadata-changed"
	second, err := identity.NewDelegationRequest(input)
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest() != second.Digest() {
		t.Fatal("non-authority replay metadata changed the authority digest")
	}
	actions := first.Actions()
	actions[0] = "forge.release"
	if first.Actions()[0] != identity.DelegatedActionAssembleContext {
		t.Fatal("caller mutated immutable request actions")
	}

	mutations := []struct {
		name   string
		mutate func(*identity.DelegationRequestInput)
	}{
		{"human", func(value *identity.DelegationRequestInput) { value.HumanSubjectID = "other-human" }},
		{"agent", func(value *identity.DelegationRequestInput) { value.AgentSubjectID = "other-agent" }},
		{"tenant", func(value *identity.DelegationRequestInput) { value.TenantID = "other-tenant" }},
		{"workspace", func(value *identity.DelegationRequestInput) { value.WorkspaceID = "other-workspace" }},
		{"project", func(value *identity.DelegationRequestInput) { value.ProjectID = "other-project" }},
		{"repository", func(value *identity.DelegationRequestInput) {
			value.RepositoryID = "other-repository"
			value.Resources = []string{"other-repository"}
		}},
		{"run", func(value *identity.DelegationRequestInput) { value.DeliveryRunID = "other-run" }},
		{"classification", func(value *identity.DelegationRequestInput) { value.Classifications = []string{"Restricted"} }},
		{"expiry", func(value *identity.DelegationRequestInput) {
			value.RequestedExpiresAt = value.RequestedExpiresAt.Add(time.Minute)
		}},
	}
	for _, test := range mutations {
		t.Run(test.name, func(t *testing.T) {
			changed := validDelegationInput()
			test.mutate(&changed)
			request, requestErr := identity.NewDelegationRequest(changed)
			if requestErr != nil {
				t.Fatal(requestErr)
			}
			if request.Digest() == first.Digest() {
				t.Fatalf("%s mutation did not change digest", test.name)
			}
		})
	}
}

func TestDelegationRequestRejectsBroaderM1Authority(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*identity.DelegationRequestInput)
	}{
		{"purpose", func(value *identity.DelegationRequestInput) { value.Purpose = "forge.release" }},
		{"action", func(value *identity.DelegationRequestInput) { value.Actions = []string{"forge.release"} }},
		{"destination", func(value *identity.DelegationRequestInput) { value.Destinations = []string{"external"} }},
		{"environment", func(value *identity.DelegationRequestInput) { value.Environment = "production" }},
		{"autonomy", func(value *identity.DelegationRequestInput) { value.MaximumAutonomy = "unbounded" }},
		{"child", func(value *identity.DelegationRequestInput) { value.ParentDelegationID = "parent" }},
		{"resource mismatch", func(value *identity.DelegationRequestInput) { value.Resources = []string{"other-repository"} }},
		{"duplicate action", func(value *identity.DelegationRequestInput) {
			value.Actions = []string{identity.DelegatedActionAssembleContext, identity.DelegatedActionAssembleContext}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := validDelegationInput()
			test.mutate(&input)
			if _, err := identity.NewDelegationRequest(input); err == nil {
				t.Fatal("broadened delegation request was accepted")
			}
		})
	}
}

func TestDelegationLifecycleIsExactIdempotentAndRevocable(t *testing.T) {
	now := testNow
	human := delegationPrincipal(t, "human-1", identity.HumanPrincipal, testNow.Add(30*time.Minute), true)
	agent := delegationPrincipal(t, "agent-1", identity.AgentPrincipal, testNow.Add(20*time.Minute), true)
	request := delegationRequest(t, validDelegationInput())
	decision := permitDelegation(request, human, testNow.Add(15*time.Minute))
	service, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	command := identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: request,
		Human: human, Agent: agent, ParentDecision: decision,
		EnforcedObligationTypes: []string{"evidence-required"},
	}
	first, err := service.CreateDelegation(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	if first.Replayed || first.Delegation.ParentAuthorityID() != decision.ID || first.Delegation.RequestDigest() != request.Digest() {
		t.Fatalf("unexpected delegation: %#v", first)
	}
	if !first.Delegation.ExpiresAt().Equal(testNow.Add(10 * time.Minute)) {
		t.Fatalf("expiry = %v", first.Delegation.ExpiresAt())
	}
	replayed, err := service.CreateDelegation(context.Background(), command)
	if err != nil || !replayed.Replayed || replayed.Delegation.ID() != first.Delegation.ID() {
		t.Fatalf("idempotent replay = %#v, %v", replayed, err)
	}
	chain, err := service.ResolveDelegationChain(context.Background(), identity.ResolveDelegationRequest{
		ContractVersion: identity.ResolveDelegationContractVersion, CorrelationID: "correlation-delegation-1",
		Human: human, Agent: agent, DelegationID: first.Delegation.ID(),
	})
	if err != nil || len(chain) != 1 || chain[0].ParentAuthorityID() != decision.ID {
		t.Fatalf("resolved chain = %#v, %v", chain, err)
	}
	projected := chain[0].Actions()
	projected[0] = "forge.release"
	if chain[0].Actions()[0] != identity.DelegatedActionAssembleContext {
		t.Fatal("caller mutated resolved delegation")
	}
	if err := service.RevokeDelegation(context.Background(), identity.RevokeDelegationRequest{
		ContractVersion: identity.RevokeDelegationContractVersion, CorrelationID: "correlation-delegation-1",
		Human: human, DelegationID: first.Delegation.ID(), Reason: "authority withdrawn",
	}); err != nil {
		t.Fatal(err)
	}
	_, err = service.ResolveDelegationChain(context.Background(), identity.ResolveDelegationRequest{
		ContractVersion: identity.ResolveDelegationContractVersion, CorrelationID: "correlation-delegation-1",
		Human: human, Agent: agent, DelegationID: first.Delegation.ID(),
	})
	assertFailureCode(t, err, failure.Forbidden)
	_, err = service.CreateDelegation(context.Background(), command)
	assertFailureCode(t, err, failure.Forbidden)

	now = testNow.Add(11 * time.Minute)
	assertFailureCode(t, first.Delegation.ValidateFor(human, agent, now, "correlation-delegation-1"), failure.Expired)
}

func TestDelegationCreationRejectsInexactOrInactiveAuthority(t *testing.T) {
	human := delegationPrincipal(t, "human-1", identity.HumanPrincipal, testNow.Add(30*time.Minute), true)
	agent := delegationPrincipal(t, "agent-1", identity.AgentPrincipal, testNow.Add(30*time.Minute), true)
	request := delegationRequest(t, validDelegationInput())
	base := permitDelegation(request, human, testNow.Add(10*time.Minute))
	tests := []struct {
		name     string
		mutate   func(*identity.CreateDelegationRequest)
		wantCode failure.Code
	}{
		{"deny", func(value *identity.CreateDelegationRequest) { value.ParentDecision.Outcome = "deny" }, failure.Forbidden},
		{"indeterminate", func(value *identity.CreateDelegationRequest) { value.ParentDecision.Outcome = "indeterminate" }, failure.Indeterminate},
		{"expired", func(value *identity.CreateDelegationRequest) {
			value.ParentDecision.EvaluatedAt = testNow.Add(-2 * time.Minute)
			expiresAt := testNow
			value.ParentDecision.ExpiresAt = &expiresAt
		}, failure.Expired},
		{"wrong point", func(value *identity.CreateDelegationRequest) { value.ParentDecision.DecisionPoint = "retrieval" }, failure.Forbidden},
		{"wrong subject", func(value *identity.CreateDelegationRequest) { value.ParentDecision.SubjectID = "other" }, failure.Forbidden},
		{"wrong action", func(value *identity.CreateDelegationRequest) { value.ParentDecision.Action = "forge.release" }, failure.Forbidden},
		{"wrong digest", func(value *identity.CreateDelegationRequest) { value.ParentDecision.Resource = "other" }, failure.Forbidden},
		{"wrong evidence correlation", func(value *identity.CreateDelegationRequest) { value.ParentDecision.EvidenceCorrelationID = "other" }, failure.Forbidden},
		{"obligation not enforced", func(value *identity.CreateDelegationRequest) { value.EnforcedObligationTypes = nil }, failure.Forbidden},
		{"wrong agent", func(value *identity.CreateDelegationRequest) {
			value.Agent = delegationPrincipal(t, "other-agent", identity.AgentPrincipal, testNow.Add(30*time.Minute), true)
		}, failure.Forbidden},
		{"non synthetic agent", func(value *identity.CreateDelegationRequest) {
			value.Agent = delegationPrincipal(t, "agent-1", identity.AgentPrincipal, testNow.Add(30*time.Minute), false)
		}, failure.Forbidden},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			command := identity.CreateDelegationRequest{
				ContractVersion: identity.CreateDelegationContractVersion, Request: request,
				Human: human, Agent: agent, ParentDecision: base, EnforcedObligationTypes: []string{"evidence-required"},
			}
			test.mutate(&command)
			_, err = service.CreateDelegation(context.Background(), command)
			assertFailureCode(t, err, test.wantCode)
		})
	}
}

func TestDelegationIdempotencyIsAtomicAndConflictingAuthorityFails(t *testing.T) {
	human := delegationPrincipal(t, "human-1", identity.HumanPrincipal, testNow.Add(time.Hour), true)
	agent := delegationPrincipal(t, "agent-1", identity.AgentPrincipal, testNow.Add(time.Hour), true)
	request := delegationRequest(t, validDelegationInput())
	decision := permitDelegation(request, human, testNow.Add(time.Hour))
	service, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	command := identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: request,
		Human: human, Agent: agent, ParentDecision: decision,
		EnforcedObligationTypes: []string{"evidence-required"},
	}
	const callers = 16
	ids := make(chan string, callers)
	errs := make(chan error, callers)
	var group sync.WaitGroup
	for range callers {
		group.Add(1)
		go func() {
			defer group.Done()
			result, createErr := service.CreateDelegation(context.Background(), command)
			if createErr != nil {
				errs <- createErr
				return
			}
			ids <- result.Delegation.ID()
		}()
	}
	group.Wait()
	close(ids)
	close(errs)
	for createErr := range errs {
		t.Fatal(createErr)
	}
	var first string
	for id := range ids {
		if first == "" {
			first = id
		}
		if id != first {
			t.Fatalf("atomic replay produced different IDs: %q != %q", id, first)
		}
	}

	changed := validDelegationInput()
	changed.RequestedExpiresAt = changed.RequestedExpiresAt.Add(time.Minute)
	changedRequest := delegationRequest(t, changed)
	changedCommand := command
	changedCommand.Request = changedRequest
	changedCommand.ParentDecision = permitDelegation(changedRequest, human, testNow.Add(time.Hour))
	_, err = service.CreateDelegation(context.Background(), changedCommand)
	assertFailureCode(t, err, failure.Conflict)
}

func TestDelegationStoreFailuresAreDisclosureSafe(t *testing.T) {
	human := delegationPrincipal(t, "human-1", identity.HumanPrincipal, testNow.Add(time.Hour), true)
	agent := delegationPrincipal(t, "agent-1", identity.AgentPrincipal, testNow.Add(time.Hour), true)
	request := delegationRequest(t, validDelegationInput())
	for _, test := range []struct {
		name string
		err  error
		want failure.Code
	}{
		{"unavailable", identity.ErrDelegationStoreUnavailable, failure.DependencyUnavailable},
		{"unknown sensitive adapter error", errors.New("secret database diagnostic"), failure.Indeterminate},
	} {
		t.Run(test.name, func(t *testing.T) {
			service, err := identity.NewDelegationService(failingDelegationStore{err: test.err}, time.Minute, func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			_, err = service.CreateDelegation(context.Background(), identity.CreateDelegationRequest{
				ContractVersion: identity.CreateDelegationContractVersion, Request: request,
				Human: human, Agent: agent, ParentDecision: permitDelegation(request, human, testNow.Add(time.Hour)),
				EnforcedObligationTypes: []string{"evidence-required"},
			})
			assertFailureCode(t, err, test.want)
			if err != nil && errors.Is(err, test.err) {
				t.Fatalf("adapter error escaped capability boundary: %v", err)
			}
		})
	}
}

func BenchmarkDelegationRequestCanonicalization(b *testing.B) {
	input := validDelegationInput()
	for index := 0; index < b.N; index++ {
		if _, err := identity.NewDelegationRequest(input); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDelegationLifecycle(b *testing.B) {
	human := delegationPrincipal(b, "human-1", identity.HumanPrincipal, testNow.Add(time.Hour), true)
	agent := delegationPrincipal(b, "agent-1", identity.AgentPrincipal, testNow.Add(time.Hour), true)
	baseInput := validDelegationInput()
	request := delegationRequest(b, baseInput)
	decision := permitDelegation(request, human, testNow.Add(time.Hour))
	newService := func() *identity.DelegationService {
		service, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return testNow })
		if err != nil {
			b.Fatal(err)
		}
		return service
	}
	command := identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: request,
		Human: human, Agent: agent, ParentDecision: decision,
		EnforcedObligationTypes: []string{"evidence-required"},
	}

	b.Run("Create", func(b *testing.B) {
		service := newService()
		b.ResetTimer()
		for index := 0; index < b.N; index++ {
			input := baseInput
			input.IdempotencyKey = "benchmark-create-" + strconv.Itoa(index)
			current, err := identity.NewDelegationRequest(input)
			if err != nil {
				b.Fatal(err)
			}
			currentCommand := command
			currentCommand.Request = current
			if _, err := service.CreateDelegation(context.Background(), currentCommand); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Replay", func(b *testing.B) {
		service := newService()
		if _, err := service.CreateDelegation(context.Background(), command); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for index := 0; index < b.N; index++ {
			if _, err := service.CreateDelegation(context.Background(), command); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ResolveActive", func(b *testing.B) {
		service := newService()
		created, err := service.CreateDelegation(context.Background(), command)
		if err != nil {
			b.Fatal(err)
		}
		resolve := identity.ResolveDelegationRequest{
			ContractVersion: identity.ResolveDelegationContractVersion, CorrelationID: baseInput.CorrelationID,
			Human: human, Agent: agent, DelegationID: created.Delegation.ID(),
		}
		b.ResetTimer()
		for index := 0; index < b.N; index++ {
			if _, err := service.ResolveDelegationChain(context.Background(), resolve); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("RevokeIdempotent", func(b *testing.B) {
		service := newService()
		created, err := service.CreateDelegation(context.Background(), command)
		if err != nil {
			b.Fatal(err)
		}
		revoke := identity.RevokeDelegationRequest{
			ContractVersion: identity.RevokeDelegationContractVersion, CorrelationID: baseInput.CorrelationID,
			Human: human, DelegationID: created.Delegation.ID(), Reason: "benchmark",
		}
		if err := service.RevokeDelegation(context.Background(), revoke); err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		for index := 0; index < b.N; index++ {
			if err := service.RevokeDelegation(context.Background(), revoke); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ResolveRevoked", func(b *testing.B) {
		service := newService()
		created, err := service.CreateDelegation(context.Background(), command)
		if err != nil {
			b.Fatal(err)
		}
		if err := service.RevokeDelegation(context.Background(), identity.RevokeDelegationRequest{
			ContractVersion: identity.RevokeDelegationContractVersion, CorrelationID: baseInput.CorrelationID,
			Human: human, DelegationID: created.Delegation.ID(), Reason: "benchmark",
		}); err != nil {
			b.Fatal(err)
		}
		resolve := identity.ResolveDelegationRequest{
			ContractVersion: identity.ResolveDelegationContractVersion, CorrelationID: baseInput.CorrelationID,
			Human: human, Agent: agent, DelegationID: created.Delegation.ID(),
		}
		b.ResetTimer()
		for index := 0; index < b.N; index++ {
			if _, err := service.ResolveDelegationChain(context.Background(), resolve); err == nil {
				b.Fatal("revoked delegation unexpectedly resolved")
			}
		}
	})
}

type failingDelegationStore struct{ err error }

func (store failingDelegationStore) Create(context.Context, identity.Delegation, string, string) (identity.Delegation, bool, error) {
	return identity.Delegation{}, false, store.err
}
func (store failingDelegationStore) Get(context.Context, string) (identity.StoredDelegation, bool, error) {
	return identity.StoredDelegation{}, false, store.err
}
func (store failingDelegationStore) Revoke(context.Context, string, string, time.Time) (identity.StoredDelegation, bool, error) {
	return identity.StoredDelegation{}, false, store.err
}

func validDelegationInput() identity.DelegationRequestInput {
	return identity.DelegationRequestInput{
		ContractVersion: identity.DelegationRequestContractVersion,
		IdempotencyKey:  "delegation-key-1", CorrelationID: "correlation-delegation-1",
		HumanSubjectID: "human-1", AgentSubjectID: "agent-1",
		TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1",
		RepositoryID: "repository-forge", DeliveryRunID: "run-1",
		Purpose: identity.M1ChangePurpose, Actions: []string{identity.DelegatedActionAssembleContext},
		Resources: []string{"repository-forge"}, Classifications: []string{"Internal"}, Destinations: []string{},
		Environment: identity.DeveloperTrialEnvironment, MaximumAutonomy: identity.ContextPreparationOnly,
		RequestedExpiresAt: testNow.Add(25 * time.Minute),
	}
}

func delegationRequest(t testing.TB, input identity.DelegationRequestInput) identity.DelegationRequest {
	t.Helper()
	request, err := identity.NewDelegationRequest(input)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func delegationPrincipal(t testing.TB, subject string, kind identity.PrincipalKind, expiresAt time.Time, synthetic bool) identity.Principal {
	t.Helper()
	provider := &providerStub{assertion: identity.Assertion{
		SubjectID: subject, Kind: kind, Issuer: "forge.synthetic", Audience: "forge-developer-trial",
		Assurance: "synthetic-local", ProfileID: "developer-trial", IssuedAt: testNow.Add(-time.Minute),
		ExpiresAt: expiresAt, Synthetic: synthetic,
	}}
	verifier, err := identity.NewVerifier(provider, "forge-developer-trial", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque", CorrelationID: "correlation-delegation-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

func permitDelegation(request identity.DelegationRequest, human identity.Principal, expiresAt time.Time) trust.PolicyDecision {
	return trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, ID: "policy-delegation-1", DecisionPoint: "execution",
		SubjectID: human.SubjectID(), Action: identity.DelegationActionCreate, Resource: request.Digest(),
		EvaluatedAt: testNow, ExpiresAt: &expiresAt, PolicyVersions: []string{"policy.m1.v1"}, Outcome: "permit",
		Obligations: []map[string]any{{"type": "evidence-required"}}, ReasonCodes: []string{"policy-permit"}, EvidenceCorrelationID: "correlation-delegation-1",
	}
}
