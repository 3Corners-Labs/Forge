package context_test

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	forgectx "github.com/3corners/forge/internal/context"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/policy"
	"github.com/3corners/forge/internal/tenant"
)

var contextNow = time.Date(2026, time.July, 22, 21, 0, 0, 0, time.UTC)

func TestAuthorityToContextSliceProducesMinimalReproducibleManifest(t *testing.T) {
	fixture := newContextFixture(t)
	first, err := fixture.service.AssembleContext(context.Background(), fixture.request)
	if err != nil {
		t.Fatal(err)
	}
	if first.Replayed || fixture.catalog.callCount() != 1 {
		t.Fatalf("first assembly replay=%t catalog calls=%d", first.Replayed, fixture.catalog.callCount())
	}
	value := first.Package
	if err := value.VerifyManifest(fixture.request.CorrelationID); err != nil {
		t.Fatal(err)
	}
	if err := value.ValidateFor(fixture.request.Envelope, contextNow, fixture.request.CorrelationID); err != nil {
		t.Fatal(err)
	}
	if value.Objective != fixture.request.Run.Objective() || value.ExecutionEnvelopeID != fixture.request.Envelope.ID ||
		len(value.Items) != 2 || value.Items[0].SourceID != "architecture" || value.Items[1].SourceID != "intent" ||
		!slices.Equal(value.Omissions, []string{"requested-source-unavailable:missing"}) || len(value.Conflicts) != 0 {
		t.Fatalf("unexpected context manifest: %#v", value)
	}
	for _, item := range value.Items {
		if item.Classification != fixture.request.Run.Classification() || item.ContentDigest == "" || item.SourceRevision == "" {
			t.Fatalf("unverified context item: %#v", item)
		}
	}

	replayed, err := fixture.service.AssembleContext(context.Background(), fixture.request)
	if err != nil || !replayed.Replayed || replayed.Package.ID != value.ID || fixture.catalog.callCount() != 1 {
		t.Fatalf("replay = %#v err=%v catalog calls=%d", replayed, err, fixture.catalog.callCount())
	}
	projected := replayed.Package.Items
	projected[0].SourceRevision = "mutated"
	again, err := fixture.service.AssembleContext(context.Background(), fixture.request)
	if err != nil || again.Package.Items[0].SourceRevision == "mutated" {
		t.Fatal("caller mutated the stored manifest")
	}

	mutated := value
	mutated.Items[0].ContentDigest = "sha256:changed"
	assertContextFailure(t, mutated.VerifyManifest(fixture.request.CorrelationID), failure.VerificationFailed)
}

func TestContextDenialsOccurBeforeSourceDiscovery(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*contextFixture)
		want   failure.Code
	}{
		{"mutated envelope", func(fixture *contextFixture) { fixture.request.Envelope.Scope.ProjectID = "other" }, failure.VerificationFailed},
		{"forged decision reference", func(fixture *contextFixture) { fixture.request.PolicyDecisionID = "forged" }, failure.Forbidden},
		{"unverified human", func(fixture *contextFixture) { fixture.request.Human = identity.Principal{} }, failure.Unauthenticated},
		{"revoked delegation", func(fixture *contextFixture) {
			if err := fixture.delegations.RevokeDelegation(context.Background(), identity.RevokeDelegationRequest{
				ContractVersion: identity.RevokeDelegationContractVersion, CorrelationID: fixture.request.CorrelationID,
				Human: fixture.request.Human, DelegationID: fixture.request.DelegationID, Reason: "test",
			}); err != nil {
				t.Fatal(err)
			}
		}, failure.Forbidden},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newContextFixture(t)
			test.mutate(&fixture)
			_, err := fixture.service.AssembleContext(context.Background(), fixture.request)
			assertContextFailure(t, err, test.want)
			if fixture.catalog.callCount() != 0 {
				t.Fatalf("source discovery occurred before authority: calls=%d", fixture.catalog.callCount())
			}
		})
	}
}

func TestContextRequiresTrustedObligationEnforcementBeforeDiscovery(t *testing.T) {
	fixture := newContextFixture(t)
	profile := validAssemblyProfile()
	profile.EnforcedObligationTypes = nil
	service, err := forgectx.NewService(
		fixture.delegations, fixture.evaluator, fixture.catalog, forgectx.NewMemoryStore(), profile,
		func() time.Time { return contextNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.AssembleContext(context.Background(), fixture.request)
	assertContextFailure(t, err, failure.Forbidden)
	if fixture.catalog.callCount() != 0 {
		t.Fatal("source discovery occurred before mandatory obligation enforcement")
	}
}

func TestContextRejectsUntrustedOrCrossScopeSourceMetadata(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*forgectx.SourceRecord)
	}{
		{"tenant", func(value *forgectx.SourceRecord) { value.TenantID = "other" }},
		{"workspace", func(value *forgectx.SourceRecord) { value.WorkspaceID = "other" }},
		{"project", func(value *forgectx.SourceRecord) { value.ProjectID = "other" }},
		{"repository", func(value *forgectx.SourceRecord) { value.RepositoryID = "other" }},
		{"classification", func(value *forgectx.SourceRecord) { value.Classification = "Public" }},
		{"authority", func(value *forgectx.SourceRecord) { value.Authority = "executable-instruction" }},
		{"digest", func(value *forgectx.SourceRecord) { value.ContentDigest = "" }},
		{"unrequested source", func(value *forgectx.SourceRecord) { value.SourceID = "secret-source" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newContextFixture(t)
			test.mutate(&fixture.catalog.records[0])
			_, err := fixture.service.AssembleContext(context.Background(), fixture.request)
			assertContextFailure(t, err, failure.Indeterminate)
		})
	}
}

func TestContextIdempotencyConflictsBeforeRediscovery(t *testing.T) {
	fixture := newContextFixture(t)
	if _, err := fixture.service.AssembleContext(context.Background(), fixture.request); err != nil {
		t.Fatal(err)
	}
	changed := fixture.request
	changed.RequestedSourceIDs = []string{"intent"}
	_, err := fixture.service.AssembleContext(context.Background(), changed)
	assertContextFailure(t, err, failure.Conflict)
	if fixture.catalog.callCount() != 1 {
		t.Fatalf("conflict triggered rediscovery: calls=%d", fixture.catalog.callCount())
	}
}

func TestContextNormalizesSourceAndStoreFailures(t *testing.T) {
	fixture := newContextFixture(t)
	fixture.catalog.err = errors.New("sensitive repository diagnostic")
	_, err := fixture.service.AssembleContext(context.Background(), fixture.request)
	assertContextFailure(t, err, failure.DependencyUnavailable)
	if errors.Is(err, fixture.catalog.err) {
		t.Fatal("source diagnostic escaped capability boundary")
	}

	fixture = newContextFixture(t)
	service, err := forgectx.NewService(
		fixture.delegations, fixture.evaluator, fixture.catalog, failingContextStore{err: errors.New("sensitive store diagnostic")},
		validAssemblyProfile(), func() time.Time { return contextNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.AssembleContext(context.Background(), fixture.request)
	assertContextFailure(t, err, failure.Indeterminate)
}

func TestContextManifestBudgetFailsClosed(t *testing.T) {
	fixture := newContextFixture(t)
	profile := validAssemblyProfile()
	profile.MaximumSerializedBytes = 10
	service, err := forgectx.NewService(fixture.delegations, fixture.evaluator, fixture.catalog, forgectx.NewMemoryStore(), profile, func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.AssembleContext(context.Background(), fixture.request)
	assertContextFailure(t, err, failure.BudgetExhausted)
}

func TestContextReproductionReportsExactChangedAndCrossRunFacts(t *testing.T) {
	fixture := newContextFixture(t)
	assembled, err := fixture.service.AssembleContext(context.Background(), fixture.request)
	if err != nil {
		t.Fatal(err)
	}
	authorizing, _, err := fixture.deliveryService.BeginContextAuthorization(context.Background(), delivery.TransitionRequest{
		CorrelationID: fixture.request.CorrelationID, IdempotencyKey: "context-authorizing",
		Run: fixture.request.Run, ExpectedVersion: fixture.request.Run.Version(),
	})
	if err != nil {
		t.Fatal(err)
	}
	ready, _, err := fixture.deliveryService.RecordContext(context.Background(), delivery.RecordContextRequest{
		TransitionRequest: delivery.TransitionRequest{
			CorrelationID: fixture.request.CorrelationID, IdempotencyKey: "context-ready",
			Run: authorizing, ExpectedVersion: authorizing.Version(),
		},
		Package: assembled.Package, Envelope: fixture.request.Envelope,
	})
	if err != nil {
		t.Fatal(err)
	}
	request := forgectx.ReproduceRequest{
		ContractVersion: forgectx.ReproduceRequestContractVersion, CorrelationID: fixture.request.CorrelationID,
		Human: fixture.request.Human, Scope: fixture.request.Scope, Run: ready,
		PackageID: assembled.Package.ID, ExpectedManifestDigest: assembled.Package.ManifestDigest,
	}
	result, err := fixture.service.ReproduceContext(context.Background(), request)
	if err != nil || !result.Exact || len(result.ChangedSourceIDs) != 0 || len(result.StaleSourceIDs) != 0 {
		t.Fatalf("exact reproduction=%#v err=%v", result, err)
	}

	fixture.catalog.mu.Lock()
	fixture.catalog.records[0].SourceRevision = "changed-revision"
	fixture.catalog.mu.Unlock()
	changed, err := fixture.service.ReproduceContext(context.Background(), request)
	if err != nil || changed.Exact || !slices.Equal(changed.ChangedSourceIDs, []string{"intent"}) {
		t.Fatalf("changed reproduction=%#v err=%v", changed, err)
	}

	crossRun := request
	crossRun.ExpectedManifestDigest = "sha256:other"
	_, err = fixture.service.ReproduceContext(context.Background(), crossRun)
	assertContextFailure(t, err, failure.Forbidden)
}

func BenchmarkAuthorityToContext(b *testing.B) {
	b.Run("Assemble", func(b *testing.B) {
		fixture := newContextFixture(b)
		for index := 0; index < b.N; index++ {
			request := fixture.request
			request.IdempotencyKey = "context-benchmark-" + strconv.Itoa(index)
			if _, err := fixture.service.AssembleContext(context.Background(), request); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Replay", func(b *testing.B) {
		fixture := newContextFixture(b)
		if _, err := fixture.service.AssembleContext(context.Background(), fixture.request); err != nil {
			b.Fatal(err)
		}
		for index := 0; index < b.N; index++ {
			if _, err := fixture.service.AssembleContext(context.Background(), fixture.request); err != nil {
				b.Fatal(err)
			}
		}
	})
}

type contextFixture struct {
	evaluator       *policy.Evaluator
	delegations     *identity.DelegationService
	catalog         *catalogStub
	service         *forgectx.Service
	deliveryService *delivery.Service
	request         forgectx.AssembleRequest
}

func newContextFixture(t testing.TB) contextFixture {
	t.Helper()
	human := contextPrincipal(t, "synthetic-human", identity.HumanPrincipal)
	agent := contextPrincipal(t, "synthetic-agent", identity.AgentPrincipal)
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-1", ID: "workspace-1", Selector: "workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-1", ID: "project-1", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: human.SubjectID(), WorkspaceID: "workspace-1"}},
		contextBindingSource{}, func() time.Time { return contextNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	scope, err := directory.ResolveWorkspaceScope(human, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: "context-correlation",
		WorkspaceSelector: "workspace", ProjectSelector: "forge", RepositorySelector: "forge",
	})
	if err != nil {
		t.Fatal(err)
	}
	binding, err := directory.GetRepositoryBinding(context.Background(), tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: "context-correlation", Scope: scope,
	})
	if err != nil {
		t.Fatal(err)
	}
	deliveryStore, err := delivery.NewMemoryStore(func() string { return "run-context-1" })
	if err != nil {
		t.Fatal(err)
	}
	deliveryService, err := delivery.NewService(deliveryStore, func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	captured, err := deliveryService.CaptureIntent(context.Background(), delivery.CaptureIntentRequest{
		ContractVersion: delivery.CaptureIntentRequestContractVersion, CorrelationID: "context-correlation", IdempotencyKey: "intent-context-1",
		Principal: human, DelegationChain: []string{}, Scope: scope, RepositoryBinding: binding,
		Objective: "Assemble the minimum Forge context", Constraints: []string{"Developer Trial only"},
		AcceptanceCriteria: []string{"Authorization precedes discovery"},
	})
	if err != nil {
		t.Fatal(err)
	}
	delegationRequest, err := identity.NewDelegationRequest(identity.DelegationRequestInput{
		ContractVersion: identity.DelegationRequestContractVersion, IdempotencyKey: "delegation-context-1", CorrelationID: "context-correlation",
		HumanSubjectID: human.SubjectID(), AgentSubjectID: agent.SubjectID(), TenantID: scope.TenantID(), WorkspaceID: scope.WorkspaceID(),
		ProjectID: scope.ProjectID(), RepositoryID: captured.Run.RepositoryID(), DeliveryRunID: captured.Run.ID(),
		Purpose: identity.M1ChangePurpose, Actions: []string{identity.DelegatedActionAssembleContext},
		Resources: []string{captured.Run.RepositoryID()}, Classifications: []string{captured.Run.Classification()}, Destinations: []string{},
		Environment: identity.DeveloperTrialEnvironment, MaximumAutonomy: identity.ContextPreparationOnly,
		RequestedExpiresAt: contextNow.Add(20 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	set := policy.PolicySet{Version: "policy.context.v1", Rules: []policy.Rule{
		{
			ID: "permit-delegation", Effect: policy.EffectPermit, Actions: []policy.Action{policy.ActionCreateDelegation},
			Environments: []policy.Environment{policy.DeveloperTrialEnvironment}, Risks: []policy.Risk{policy.RiskLow},
			SubjectIDs: []string{human.SubjectID()}, ResourceIDs: []string{delegationRequest.Digest()}, DeliveryRunIDs: []string{captured.Run.ID()},
			Purposes: []string{captured.Run.Purpose()}, Classifications: []string{captured.Run.Classification()},
			TenantIDs: []string{scope.TenantID()}, WorkspaceIDs: []string{scope.WorkspaceID()}, ProjectIDs: []string{scope.ProjectID()},
		},
		{
			ID: "permit-context", Effect: policy.EffectPermit, Actions: []policy.Action{policy.ActionAssembleContext},
			Environments: []policy.Environment{policy.DeveloperTrialEnvironment}, Risks: []policy.Risk{policy.RiskLow},
			SubjectIDs: []string{human.SubjectID()}, ResourceIDs: []string{captured.Run.RepositoryID()}, DeliveryRunIDs: []string{captured.Run.ID()},
			Purposes: []string{captured.Run.Purpose()}, Classifications: []string{captured.Run.Classification()},
			TenantIDs: []string{scope.TenantID()}, WorkspaceIDs: []string{scope.WorkspaceID()}, ProjectIDs: []string{scope.ProjectID()},
			Obligations: []policy.Obligation{{Type: "context-evidence", Parameters: map[string]any{"required": true}}},
		},
	}}
	evaluator, err := policy.NewEvaluator(set, []string{"context-evidence"}, 5*time.Minute, func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	delegationDecision, err := evaluator.EvaluateDelegation(context.Background(), policy.EvaluateDelegationRequest{
		ContractVersion: policy.EvaluateDelegationRequestContractVersion, CorrelationID: "context-correlation",
		Human: human, Agent: agent, Scope: scope, Run: captured.Run, Delegation: delegationRequest,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		t.Fatal(err)
	}
	delegations, err := identity.NewDelegationService(identity.NewMemoryDelegationStore(), 10*time.Minute, func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	delegation, err := delegations.CreateDelegation(context.Background(), identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: delegationRequest,
		Human: human, Agent: agent, ParentDecision: delegationDecision,
	})
	if err != nil {
		t.Fatal(err)
	}
	contextDecision, err := evaluator.Evaluate(context.Background(), policy.EvaluateRequest{
		ContractVersion: policy.EvaluateRequestContractVersion, CorrelationID: "context-correlation",
		Principal: human, Scope: scope, Run: captured.Run, Action: policy.ActionAssembleContext,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		t.Fatal(err)
	}
	envelopeIssuer, err := policy.NewEnvelopeIssuer(evaluator, delegations, policy.NewMemoryEnvelopeStore(), policy.EnvelopeProfile{
		Issuer: "forge.policy.developer-trial", Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
		AllowedPolicyVersions: []string{"policy.context.v1"}, EnforcedObligationTypes: []string{"context-evidence"},
		MaximumLifetime: 10 * time.Minute,
		Limits:          trust.ExecutionLimits{MaxDurationSeconds: 600, MaxConcurrency: 1},
	}, func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := envelopeIssuer.IssueExecutionEnvelope(context.Background(), policy.IssueExecutionEnvelopeRequest{
		ContractVersion: policy.IssueExecutionEnvelopeContractVersion, CorrelationID: "context-correlation", IdempotencyKey: "envelope-context-1",
		Human: human, Agent: agent, Scope: scope, Run: captured.Run, DelegationID: delegation.Delegation.ID(), PolicyDecisionID: contextDecision.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	catalog := &catalogStub{records: []forgectx.SourceRecord{
		contextSource("intent", "docs/intent.md"), contextSource("architecture", "docs/architecture/overview.md"),
	}}
	contextStore := forgectx.NewMemoryStore()
	service, err := forgectx.NewService(delegations, evaluator, catalog, contextStore, validAssemblyProfile(), func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	return contextFixture{
		evaluator: evaluator, delegations: delegations, catalog: catalog, service: service, deliveryService: deliveryService,
		request: forgectx.AssembleRequest{
			ContractVersion: forgectx.AssembleRequestContractVersion, CorrelationID: "context-correlation", IdempotencyKey: "context-1",
			Human: human, Agent: agent, Scope: scope, Run: captured.Run, DelegationID: delegation.Delegation.ID(),
			Envelope: envelope.Envelope, PolicyDecisionID: contextDecision.ID,
			RequestedSourceIDs: []string{"intent", "missing", "architecture"},
		},
	}
}

type contextProvider struct{ assertion identity.Assertion }

func (provider contextProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return provider.assertion, nil
}

func contextPrincipal(t testing.TB, subject string, kind identity.PrincipalKind) identity.Principal {
	t.Helper()
	verifier, err := identity.NewVerifier(contextProvider{identity.Assertion{
		SubjectID: subject, Kind: kind, Issuer: "forge.synthetic", Audience: "forge-context",
		Assurance: "synthetic-local", ProfileID: "developer-trial", IssuedAt: contextNow.Add(-time.Minute),
		ExpiresAt: contextNow.Add(30 * time.Minute), Synthetic: true,
	}}, "forge-context", func() time.Time { return contextNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque", CorrelationID: "context-correlation",
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

type contextBindingSource struct{}

func (contextBindingSource) ResolveRepositoryBinding(context.Context, string, string, string, string) (tenant.RepositoryBindingDefinition, bool, error) {
	return tenant.RepositoryBindingDefinition{
		TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", RepositoryID: "repository-forge",
		Selector: "forge", CanonicalLocator: "https://example.invalid/forge.git", Classification: "Internal",
	}, true, nil
}

type catalogStub struct {
	mu      sync.Mutex
	records []forgectx.SourceRecord
	err     error
	calls   int
}

func (catalog *catalogStub) Discover(_ context.Context, request forgectx.SourceDiscoveryRequest) ([]forgectx.SourceRecord, error) {
	catalog.mu.Lock()
	defer catalog.mu.Unlock()
	catalog.calls++
	if request.TenantID != "tenant-1" || request.WorkspaceID != "workspace-1" || request.ProjectID != "project-1" || request.RepositoryID != "repository-forge" {
		return nil, errors.New("scope was not authority-derived")
	}
	return slices.Clone(catalog.records), catalog.err
}

func (catalog *catalogStub) callCount() int {
	catalog.mu.Lock()
	defer catalog.mu.Unlock()
	return catalog.calls
}

func contextSource(id, locator string) forgectx.SourceRecord {
	confidence := 1.0
	freshness := contextNow.Add(-time.Minute)
	return forgectx.SourceRecord{
		TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", RepositoryID: "repository-forge",
		SourceID: id, SourceRevision: "revision-" + id, Locator: locator, Classification: "Internal",
		Authority: "accepted-knowledge", Purpose: "bounded task context", FreshnessAt: &freshness,
		Confidence: &confidence, ContentDigest: "sha256:" + id,
	}
}

func validAssemblyProfile() forgectx.AssemblyProfile {
	return forgectx.AssemblyProfile{
		MaximumItems: 100, MaximumSerializedBytes: 1024 * 1024, MaximumLifetime: 5 * time.Minute,
		EnforcedObligationTypes: []string{"context-evidence"},
	}
}

type failingContextStore struct{ err error }

func (store failingContextStore) Get(context.Context, string, string, string) (trust.ContextPackage, bool, error) {
	return trust.ContextPackage{}, false, store.err
}
func (store failingContextStore) Create(context.Context, trust.ContextPackage, string, string, string) (trust.ContextPackage, bool, error) {
	return trust.ContextPackage{}, false, store.err
}

func assertContextFailure(t testing.TB, err error, want failure.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("error=nil, want %q", want)
	}
	var normalized *failure.Error
	if !errors.As(err, &normalized) {
		t.Fatalf("error %T %v is not normalized", err, err)
	}
	if normalized.Code != want {
		t.Fatalf("code=%q, want %q", normalized.Code, want)
	}
	if validationErr := normalized.Validate(); validationErr != nil {
		t.Fatalf("invalid normalized error: %v", validationErr)
	}
}
