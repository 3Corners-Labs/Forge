package tenant_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

var testNow = time.Date(2026, 7, 22, 8, 0, 0, 0, time.UTC)

type identityProvider struct {
	subjectID string
	expiresAt time.Time
}

func (provider identityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: provider.subjectID, Kind: identity.HumanPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "forge.identity.synthetic.test",
		IssuedAt: testNow.Add(-time.Minute), ExpiresAt: provider.expiresAt, Synthetic: true,
	}, nil
}

type bindingSource struct {
	definition tenant.RepositoryBindingDefinition
	found      bool
	err        error
	calls      int
	seen       [4]string
}

type selectorBindingSource struct{}

func (selectorBindingSource) ResolveRepositoryBinding(_ context.Context, tenantID, workspaceID, projectID, selector string) (tenant.RepositoryBindingDefinition, bool, error) {
	return tenant.RepositoryBindingDefinition{
		TenantID: tenantID, WorkspaceID: workspaceID, ProjectID: projectID,
		RepositoryID: "repository-" + selector, Selector: selector,
		CanonicalLocator: "contract-test://" + selector, Classification: "Internal",
	}, true, nil
}

func (source *bindingSource) ResolveRepositoryBinding(_ context.Context, tenantID, workspaceID, projectID, selector string) (tenant.RepositoryBindingDefinition, bool, error) {
	source.calls++
	source.seen = [4]string{tenantID, workspaceID, projectID, selector}
	return source.definition, source.found, source.err
}

func TestMemberResolvesImmutableScopeThenRepositoryBinding(t *testing.T) {
	t.Parallel()
	source := validBindingSource()
	directory := validDirectory(t, "synthetic-developer-allowed", source, func() time.Time { return testNow })
	principal := verifiedPrincipal(t, "synthetic-developer-allowed", testNow.Add(30*time.Minute))
	scope, err := directory.ResolveWorkspaceScope(principal, validRequest())
	if err != nil {
		t.Fatalf("ResolveWorkspaceScope() error = %v", err)
	}
	if scope.ContractVersion() != tenant.ScopeContractVersion || scope.TenantID() != "tenant-fixture" || scope.WorkspaceID() != "workspace-fixture" || scope.ProjectID() != "project-forge" || scope.PrincipalID() != principal.SubjectID() {
		t.Fatalf("unexpected scope: tenant=%q workspace=%q project=%q principal=%q", scope.TenantID(), scope.WorkspaceID(), scope.ProjectID(), scope.PrincipalID())
	}
	if source.calls != 0 {
		t.Fatalf("repository discovery occurred during scope resolution: calls=%d", source.calls)
	}
	binding, err := directory.GetRepositoryBinding(context.Background(), bindingRequest(scope))
	if err != nil {
		t.Fatalf("GetRepositoryBinding() error = %v", err)
	}
	if source.calls != 1 || source.seen != [4]string{"tenant-fixture", "workspace-fixture", "project-forge", "forge"} {
		t.Fatalf("unexpected binding lookup: calls=%d seen=%v", source.calls, source.seen)
	}
	if binding.ContractVersion() != tenant.RepositoryBindingContractVersion || binding.RepositoryID() != "repository-forge" || binding.CanonicalLocator() != "contract-test://forge" || binding.Classification() != "Internal" {
		t.Fatalf("unexpected binding: id=%q locator=%q class=%q", binding.RepositoryID(), binding.CanonicalLocator(), binding.Classification())
	}
}

func TestNonMemberIsDeniedBeforeRepositoryDiscovery(t *testing.T) {
	t.Parallel()
	source := validBindingSource()
	directory := validDirectory(t, "another-principal", source, func() time.Time { return testNow })
	principal := verifiedPrincipal(t, "synthetic-observer-denied", testNow.Add(30*time.Minute))
	_, err := directory.ResolveWorkspaceScope(principal, validRequest())
	assertFailureCodeAndMessage(t, err, failure.Forbidden, "requested scope is not available to the authenticated principal")
	if source.calls != 0 {
		t.Fatalf("repository discovery occurred for non-member: calls=%d", source.calls)
	}
}

func TestCallerRequestHasNoTenantOverride(t *testing.T) {
	t.Parallel()
	typeOfRequest := reflect.TypeOf(tenant.ScopeRequest{})
	for index := 0; index < typeOfRequest.NumField(); index++ {
		if typeOfRequest.Field(index).Name == "TenantID" || typeOfRequest.Field(index).Name == "TenantSelector" {
			t.Fatalf("caller scope request exposes authoritative tenant field %q", typeOfRequest.Field(index).Name)
		}
	}
}

func TestUnmintedScopeCannotTriggerRepositoryDiscovery(t *testing.T) {
	t.Parallel()
	source := validBindingSource()
	directory := validDirectory(t, "synthetic-developer-allowed", source, func() time.Time { return testNow })
	_, err := directory.GetRepositoryBinding(context.Background(), bindingRequest(tenant.Scope{}))
	assertFailureCodeAndMessage(t, err, failure.Forbidden, "requested scope is not available to the authenticated principal")
	if source.calls != 0 {
		t.Fatalf("repository discovery occurred for unminted scope: calls=%d", source.calls)
	}
}

func TestRepositoryBindingCannotCrossTenantAuthorities(t *testing.T) {
	t.Parallel()
	principal := verifiedPrincipal(t, "synthetic-developer-allowed", testNow.Add(30*time.Minute))
	newDirectory := func(authorityID string) *tenant.Directory {
		directory, err := tenant.NewDirectory(
			authorityID,
			[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
			[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
			[]tenant.MembershipDefinition{{PrincipalID: "synthetic-developer-allowed", WorkspaceID: "workspace-fixture"}},
			validBindingSource(),
			func() time.Time { return testNow },
		)
		if err != nil {
			t.Fatal(err)
		}
		return directory
	}
	directoryA := newDirectory("tenant-authority-a")
	directoryB := newDirectory("tenant-authority-b")
	scopeA, err := directoryA.ResolveWorkspaceScope(principal, validRequest())
	if err != nil {
		t.Fatal(err)
	}
	scopeB, err := directoryB.ResolveWorkspaceScope(principal, validRequest())
	if err != nil {
		t.Fatal(err)
	}
	bindingB, err := directoryB.GetRepositoryBinding(context.Background(), bindingRequest(scopeB))
	if err != nil {
		t.Fatal(err)
	}
	assertFailureCodeAndMessage(t, bindingB.ValidateFor(scopeA, "correlation-1"), failure.Forbidden, "requested scope is not available to the authenticated principal")
}

func TestRepositoryBindingCannotCrossRepositorySelectors(t *testing.T) {
	t.Parallel()
	principal := verifiedPrincipal(t, "synthetic-developer-allowed", testNow.Add(30*time.Minute))
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: "synthetic-developer-allowed", WorkspaceID: "workspace-fixture"}},
		selectorBindingSource{},
		func() time.Time { return testNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	requestA := validRequest()
	requestA.RepositorySelector = "forge-a"
	scopeA, err := directory.ResolveWorkspaceScope(principal, requestA)
	if err != nil {
		t.Fatal(err)
	}
	requestB := validRequest()
	requestB.RepositorySelector = "forge-b"
	scopeB, err := directory.ResolveWorkspaceScope(principal, requestB)
	if err != nil {
		t.Fatal(err)
	}
	bindingB, err := directory.GetRepositoryBinding(context.Background(), bindingRequest(scopeB))
	if err != nil {
		t.Fatal(err)
	}
	assertFailureCodeAndMessage(t, bindingB.ValidateFor(scopeA, "correlation-1"), failure.Forbidden, "requested scope is not available to the authenticated principal")
}

func TestUnsupportedRequestVersionsFailBeforeRepositoryDiscovery(t *testing.T) {
	t.Parallel()
	source := validBindingSource()
	directory := validDirectory(t, "synthetic-developer-allowed", source, func() time.Time { return testNow })
	principal := verifiedPrincipal(t, "synthetic-developer-allowed", testNow.Add(30*time.Minute))
	scopeRequest := validRequest()
	scopeRequest.ContractVersion = "other"
	_, err := directory.ResolveWorkspaceScope(principal, scopeRequest)
	assertFailureCode(t, err, failure.InvalidRequest)
	_, err = directory.GetRepositoryBinding(context.Background(), tenant.BindingRequest{ContractVersion: "other", CorrelationID: "correlation-1"})
	assertFailureCode(t, err, failure.InvalidRequest)
	if source.calls != 0 {
		t.Fatalf("repository discovery occurred for unsupported contract versions: calls=%d", source.calls)
	}
}

func TestUnknownAndUnauthorizedScopesHaveSameDisclosureSafeFailure(t *testing.T) {
	t.Parallel()
	principal := verifiedPrincipal(t, "synthetic-observer-denied", testNow.Add(30*time.Minute))
	tests := []tenant.ScopeRequest{
		validRequest(),
		mutateRequest(func(request *tenant.ScopeRequest) { request.WorkspaceSelector = "unknown" }),
		mutateRequest(func(request *tenant.ScopeRequest) { request.ProjectSelector = "unknown" }),
	}
	for _, request := range tests {
		source := validBindingSource()
		directory := validDirectory(t, "another-principal", source, func() time.Time { return testNow })
		_, err := directory.ResolveWorkspaceScope(principal, request)
		assertFailureCodeAndMessage(t, err, failure.Forbidden, "requested scope is not available to the authenticated principal")
		if source.calls != 0 {
			t.Fatalf("repository discovery occurred for denied request %+v", request)
		}
	}
}

func TestRepositoryBindingResolutionFailsClosed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*bindingSource)
		want   failure.Code
	}{
		{"missing", func(source *bindingSource) { source.found = false }, failure.Forbidden},
		{"unavailable", func(source *bindingSource) { source.err = errors.New("store diagnostic") }, failure.DependencyUnavailable},
		{"scope mismatch", func(source *bindingSource) { source.definition.TenantID = "other-tenant" }, failure.Indeterminate},
		{"unknown classification", func(source *bindingSource) { source.definition.Classification = "Unknown" }, failure.Indeterminate},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			source := validBindingSource()
			test.mutate(source)
			directory := validDirectory(t, "synthetic-developer-allowed", source, func() time.Time { return testNow })
			principal := verifiedPrincipal(t, "synthetic-developer-allowed", testNow.Add(30*time.Minute))
			scope, err := directory.ResolveWorkspaceScope(principal, validRequest())
			if err != nil {
				t.Fatal(err)
			}
			_, err = directory.GetRepositoryBinding(context.Background(), bindingRequest(scope))
			assertFailureCode(t, err, test.want)
		})
	}
}

func TestExpiredPrincipalBlocksScopeAndBindingResolution(t *testing.T) {
	t.Parallel()
	source := validBindingSource()
	directory := validDirectory(t, "synthetic-developer-allowed", source, func() time.Time { return testNow })
	expired := verifiedPrincipalAt(t, "synthetic-developer-allowed", testNow.Add(-time.Minute), testNow.Add(-time.Second))
	_, err := directory.ResolveWorkspaceScope(expired, validRequest())
	assertFailureCode(t, err, failure.Expired)
	if source.calls != 0 {
		t.Fatalf("repository discovery occurred for expired principal: calls=%d", source.calls)
	}

	active := verifiedPrincipal(t, "synthetic-developer-allowed", testNow.Add(time.Minute))
	scope, err := directory.ResolveWorkspaceScope(active, validRequest())
	if err != nil {
		t.Fatal(err)
	}
	directoryAtExpiry := validDirectory(t, "synthetic-developer-allowed", source, func() time.Time { return testNow.Add(time.Minute) })
	_, err = directoryAtExpiry.GetRepositoryBinding(context.Background(), bindingRequest(scope))
	assertFailureCode(t, err, failure.Expired)
}

func validDirectory(t *testing.T, memberPrincipal string, source *bindingSource, now func() time.Time) *tenant.Directory {
	t.Helper()
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: memberPrincipal, WorkspaceID: "workspace-fixture"}},
		source,
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	return directory
}

func validBindingSource() *bindingSource {
	return &bindingSource{found: true, definition: tenant.RepositoryBindingDefinition{
		TenantID: "tenant-fixture", WorkspaceID: "workspace-fixture", ProjectID: "project-forge",
		RepositoryID: "repository-forge", Selector: "forge", CanonicalLocator: "contract-test://forge", Classification: "Internal",
	}}
}

func validRequest() tenant.ScopeRequest {
	return tenant.ScopeRequest{ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: "correlation-1", WorkspaceSelector: "fixture-workspace", ProjectSelector: "forge", RepositorySelector: "forge"}
}

func mutateRequest(mutate func(*tenant.ScopeRequest)) tenant.ScopeRequest {
	request := validRequest()
	mutate(&request)
	return request
}

func bindingRequest(scope tenant.Scope) tenant.BindingRequest {
	return tenant.BindingRequest{ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: "correlation-1", Scope: scope}
}

func verifiedPrincipal(t *testing.T, subjectID string, expiresAt time.Time) identity.Principal {
	t.Helper()
	return verifiedPrincipalAt(t, subjectID, testNow.Add(-time.Minute), expiresAt)
}

func verifiedPrincipalAt(t *testing.T, subjectID string, issuedAt, expiresAt time.Time) identity.Principal {
	t.Helper()
	provider := identityProvider{subjectID: subjectID, expiresAt: expiresAt}
	verifier, err := identity.NewVerifier(provider, "forge-cli", func() time.Time { return issuedAt.Add(time.Second) })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{ContractVersion: identity.VerifyRequestContractVersion, Credential: "credential", CorrelationID: "correlation-1"})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

func assertFailureCodeAndMessage(t *testing.T, err error, code failure.Code, message string) {
	t.Helper()
	assertFailureCode(t, err, code)
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.SafeMessage != message {
		t.Fatalf("safe message = %q, want %q", normalized.SafeMessage, message)
	}
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
