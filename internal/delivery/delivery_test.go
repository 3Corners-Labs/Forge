package delivery_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

var testNow = time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)

type identityProvider struct {
	subjectID string
	kind      identity.PrincipalKind
	expiresAt time.Time
}

func (provider identityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: provider.subjectID, Kind: provider.kind, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "forge.identity.synthetic.developer-permitted",
		IssuedAt: testNow.Add(-time.Minute), ExpiresAt: provider.expiresAt, Synthetic: true,
	}, nil
}

type bindingSource struct {
	classification string
}

func (source bindingSource) ResolveRepositoryBinding(context.Context, string, string, string, string) (tenant.RepositoryBindingDefinition, bool, error) {
	return tenant.RepositoryBindingDefinition{
		TenantID: "tenant-fixture", WorkspaceID: "workspace-fixture", ProjectID: "project-forge",
		RepositoryID: "repository-forge", Selector: "forge", CanonicalLocator: "contract-test://forge", Classification: source.classification,
	}, true, nil
}

type storeSpy struct {
	mu       sync.Mutex
	calls    int
	err      error
	stored   delivery.StoredRecord
	replayed bool
}

func (store *storeSpy) CreateCaptured(_ context.Context, _ delivery.CaptureRecord) (delivery.StoredRecord, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.calls++
	return store.stored, store.replayed, store.err
}

func (store *storeSpy) callCount() int {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.calls
}

func TestCaptureIntentCreatesImmutableCapturedRun(t *testing.T) {
	request := validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	store, err := delivery.NewMemoryStore(func() string { return "delivery-run-1" })
	if err != nil {
		t.Fatal(err)
	}
	service, err := delivery.NewService(store, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.CaptureIntent(context.Background(), request)
	if err != nil {
		t.Fatalf("CaptureIntent() error = %v", err)
	}
	if result.ContractVersion != delivery.CaptureIntentResultContractVersion || result.Replayed {
		t.Fatalf("unexpected result: version=%q replayed=%t", result.ContractVersion, result.Replayed)
	}
	run := result.Run
	if run.ContractVersion() != delivery.DeliveryRunContractVersion || run.ID() != "delivery-run-1" || run.State() != delivery.CapturedState || run.Version() != 1 {
		t.Fatalf("unexpected run identity: version=%q id=%q state=%q revision=%d", run.ContractVersion(), run.ID(), run.State(), run.Version())
	}
	if run.TenantID() != "tenant-fixture" || run.WorkspaceID() != "workspace-fixture" || run.ProjectID() != "project-forge" || run.RepositoryID() != "repository-forge" {
		t.Fatalf("unexpected authoritative scope: tenant=%q workspace=%q project=%q repository=%q", run.TenantID(), run.WorkspaceID(), run.ProjectID(), run.RepositoryID())
	}
	if run.ActorID() != "synthetic-developer-allowed" || run.Purpose() != delivery.M1ChangePurpose || run.Classification() != "Internal" || run.CreatedAt() != testNow {
		t.Fatalf("unexpected captured metadata: actor=%q purpose=%q classification=%q created=%s", run.ActorID(), run.Purpose(), run.Classification(), run.CreatedAt())
	}
	request.Constraints[0] = "mutated"
	request.AcceptanceCriteria[0] = "mutated"
	constraints := run.Constraints()
	criteria := run.AcceptanceCriteria()
	constraints[0] = "caller mutation"
	criteria[0] = "caller mutation"
	if run.Constraints()[0] != "No production access" || run.AcceptanceCriteria()[0] != "Intent is captured for plan review" {
		t.Fatal("captured run slices are mutable through caller-owned data")
	}
}

func TestCaptureIntentIsScopeBoundAndIdempotent(t *testing.T) {
	var generated atomic.Int64
	store, err := delivery.NewMemoryStore(func() string {
		generated.Add(1)
		return "delivery-run-1"
	})
	if err != nil {
		t.Fatal(err)
	}
	service, err := delivery.NewService(store, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	request := validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	first, err := service.CaptureIntent(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	retry := validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	retry.CorrelationID = "correlation-retry"
	second, err := service.CaptureIntent(context.Background(), retry)
	if err != nil {
		t.Fatal(err)
	}
	if !second.Replayed || first.Run.ID() != second.Run.ID() || generated.Load() != 1 {
		t.Fatalf("idempotent retry created another run: first=%q second=%q replayed=%t generated=%d", first.Run.ID(), second.Run.ID(), second.Replayed, generated.Load())
	}
	changed := retry
	changed.Objective = "FORGE_TEST_SECRET_DO_NOT_RECORD"
	_, err = service.CaptureIntent(context.Background(), changed)
	assertFailureCode(t, err, failure.Conflict)
	if strings.Contains(err.Error(), changed.Objective) {
		t.Fatal("conflict disclosed intent content")
	}
}

func TestCaptureIntentRejectsNilContextBeforePersistence(t *testing.T) {
	store := &storeSpy{}
	service, err := delivery.NewService(store, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.CaptureIntent(nil, validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute)))
	assertFailureCode(t, err, failure.InvalidRequest)
	if store.callCount() != 0 {
		t.Fatalf("store calls = %d, want 0", store.callCount())
	}
}

func TestConcurrentCaptureCommitsOneRun(t *testing.T) {
	var generated atomic.Int64
	store, err := delivery.NewMemoryStore(func() string {
		generated.Add(1)
		return "delivery-run-1"
	})
	if err != nil {
		t.Fatal(err)
	}
	service, err := delivery.NewService(store, func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	request := validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	const callers = 16
	results := make(chan delivery.CaptureIntentResult, callers)
	errorsSeen := make(chan error, callers)
	var wait sync.WaitGroup
	for range callers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			result, captureErr := service.CaptureIntent(context.Background(), request)
			results <- result
			errorsSeen <- captureErr
		}()
	}
	wait.Wait()
	close(results)
	close(errorsSeen)
	for captureErr := range errorsSeen {
		if captureErr != nil {
			t.Fatalf("CaptureIntent() error = %v", captureErr)
		}
	}
	for result := range results {
		if result.Run.ID() != "delivery-run-1" {
			t.Fatalf("concurrent run ID = %q", result.Run.ID())
		}
	}
	if generated.Load() != 1 {
		t.Fatalf("generated run IDs = %d, want 1", generated.Load())
	}
}

func TestCaptureIntentFailsClosedBeforePersistence(t *testing.T) {
	base := validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	otherPrincipal := verifiedPrincipal(t, "other-developer", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	tests := []struct {
		name   string
		at     time.Time
		mutate func(*delivery.CaptureIntentRequest)
		want   failure.Code
	}{
		{"unsupported version", testNow, func(value *delivery.CaptureIntentRequest) { value.ContractVersion = "other" }, failure.InvalidRequest},
		{"missing correlation", testNow, func(value *delivery.CaptureIntentRequest) { value.CorrelationID = "" }, failure.InvalidRequest},
		{"missing idempotency", testNow, func(value *delivery.CaptureIntentRequest) { value.IdempotencyKey = "" }, failure.InvalidRequest},
		{"missing delegation chain", testNow, func(value *delivery.CaptureIntentRequest) { value.DelegationChain = nil }, failure.InvalidRequest},
		{"delegated actor", testNow, func(value *delivery.CaptureIntentRequest) { value.DelegationChain = []string{"delegation-1"} }, failure.Forbidden},
		{"workload actor", testNow, func(value *delivery.CaptureIntentRequest) {
			value.Principal = verifiedPrincipal(t, "workload", identity.WorkloadPrincipal, testNow.Add(30*time.Minute))
		}, failure.Forbidden},
		{"expired actor", testNow.Add(31 * time.Minute), func(*delivery.CaptureIntentRequest) {}, failure.Expired},
		{"scope actor mismatch", testNow, func(value *delivery.CaptureIntentRequest) { value.Principal = otherPrincipal }, failure.Forbidden},
		{"unminted scope", testNow, func(value *delivery.CaptureIntentRequest) { value.Scope = tenant.Scope{} }, failure.Forbidden},
		{"unminted binding", testNow, func(value *delivery.CaptureIntentRequest) { value.RepositoryBinding = tenant.RepositoryBinding{} }, failure.Forbidden},
		{"missing objective", testNow, func(value *delivery.CaptureIntentRequest) { value.Objective = "" }, failure.InvalidRequest},
		{"empty constraint", testNow, func(value *delivery.CaptureIntentRequest) { value.Constraints = []string{""} }, failure.InvalidRequest},
		{"missing acceptance criteria", testNow, func(value *delivery.CaptureIntentRequest) { value.AcceptanceCriteria = nil }, failure.InvalidRequest},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			request := cloneRequest(base)
			test.mutate(&request)
			store := &storeSpy{}
			service, err := delivery.NewService(store, func() time.Time { return test.at })
			if err != nil {
				t.Fatal(err)
			}
			_, err = service.CaptureIntent(context.Background(), request)
			assertFailureCode(t, err, test.want)
			if store.callCount() != 0 {
				t.Fatalf("store calls = %d, want 0", store.callCount())
			}
		})
	}
}

func TestCaptureIntentNormalizesStoreFailures(t *testing.T) {
	request := validRequest(t, "synthetic-developer-allowed", identity.HumanPrincipal, testNow.Add(30*time.Minute))
	tests := []struct {
		name  string
		store *storeSpy
		want  failure.Code
	}{
		{"unavailable", &storeSpy{err: errors.New("sensitive database diagnostic")}, failure.DependencyUnavailable},
		{"idempotency conflict", &storeSpy{err: delivery.ErrIdempotencyConflict}, failure.Conflict},
		{"malformed stored record", &storeSpy{}, failure.Indeterminate},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			service, err := delivery.NewService(test.store, func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			_, err = service.CaptureIntent(context.Background(), request)
			assertFailureCode(t, err, test.want)
			if strings.Contains(err.Error(), "sensitive database diagnostic") {
				t.Fatal("store diagnostic escaped capability boundary")
			}
		})
	}
}

func validRequest(t *testing.T, subjectID string, kind identity.PrincipalKind, expiresAt time.Time) delivery.CaptureIntentRequest {
	t.Helper()
	principal := verifiedPrincipal(t, subjectID, kind, expiresAt)
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: subjectID, WorkspaceID: "workspace-fixture"}},
		bindingSource{classification: "Internal"},
		func() time.Time { return testNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	scope, err := directory.ResolveWorkspaceScope(principal, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: "correlation-1",
		WorkspaceSelector: "fixture-workspace", ProjectSelector: "forge", RepositorySelector: "forge",
	})
	if err != nil {
		t.Fatal(err)
	}
	binding, err := directory.GetRepositoryBinding(context.Background(), tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: "correlation-1", Scope: scope,
	})
	if err != nil {
		t.Fatal(err)
	}
	return delivery.CaptureIntentRequest{
		ContractVersion: delivery.CaptureIntentRequestContractVersion,
		CorrelationID:   "correlation-1", IdempotencyKey: "intent-key-1", Principal: principal,
		DelegationChain: []string{}, Scope: scope, RepositoryBinding: binding,
		Objective:          "Add a documentation-only marker to the fixture repository",
		Constraints:        []string{"No production access", "No source publication"},
		AcceptanceCriteria: []string{"Intent is captured for plan review"},
	}
}

func verifiedPrincipal(t *testing.T, subjectID string, kind identity.PrincipalKind, expiresAt time.Time) identity.Principal {
	t.Helper()
	verifier, err := identity.NewVerifier(identityProvider{subjectID: subjectID, kind: kind, expiresAt: expiresAt}, "forge-cli", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "credential", CorrelationID: "correlation-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

func cloneRequest(value delivery.CaptureIntentRequest) delivery.CaptureIntentRequest {
	value.DelegationChain = append([]string{}, value.DelegationChain...)
	value.Constraints = append([]string{}, value.Constraints...)
	value.AcceptanceCriteria = append([]string{}, value.AcceptanceCriteria...)
	return value
}

func assertFailureCode(t *testing.T, err error, want failure.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want %q", want)
	}
	var normalized *failure.Error
	if !errors.As(err, &normalized) {
		t.Fatalf("error type = %T, want *failure.Error", err)
	}
	if normalized.Code != want {
		t.Fatalf("error code = %q, want %q (error: %v)", normalized.Code, want, err)
	}
	if validationErr := normalized.Validate(); validationErr != nil {
		t.Fatalf("normalized error invalid: %v", validationErr)
	}
}
