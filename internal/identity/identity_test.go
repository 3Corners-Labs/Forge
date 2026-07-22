package identity_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/identity"
)

var testNow = time.Date(2026, 7, 22, 8, 0, 0, 0, time.UTC)

type providerStub struct {
	assertion identity.Assertion
	err       error
	seen      string
}

func (provider *providerStub) Verify(_ context.Context, credential string) (identity.Assertion, error) {
	provider.seen = credential
	return provider.assertion, provider.err
}

func TestVerifierReturnsNormalizedImmutablePrincipal(t *testing.T) {
	t.Parallel()
	provider := &providerStub{assertion: validAssertion()}
	verifier, err := identity.NewVerifier(provider, "forge-cli", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), verifyRequest("opaque-credential"))
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if provider.seen != "opaque-credential" {
		t.Fatalf("provider credential = %q", provider.seen)
	}
	if principal.ContractVersion() != identity.PrincipalContractVersion || principal.SubjectID() != "synthetic-developer-allowed" || principal.Kind() != identity.HumanPrincipal || !principal.Synthetic() {
		t.Fatalf("unexpected principal: version=%q subject=%q kind=%q synthetic=%t", principal.ContractVersion(), principal.SubjectID(), principal.Kind(), principal.Synthetic())
	}
	if err := principal.ValidateAt(testNow.Add(time.Minute), "correlation-1"); err != nil {
		t.Fatalf("ValidateAt() error = %v", err)
	}
}

func TestVerifierFailsClosedWithoutDisclosingProviderErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		provider   *providerStub
		credential string
		want       failure.Code
	}{
		{"missing credential", &providerStub{assertion: validAssertion()}, "", failure.Unauthenticated},
		{"invalid provider assertion", &providerStub{err: errors.New("sensitive provider diagnostic")}, "credential", failure.Unauthenticated},
		{"unavailable provider", &providerStub{err: identity.ErrProviderUnavailable}, "credential", failure.DependencyUnavailable},
		{"wrong audience", &providerStub{assertion: mutateAssertion(func(value *identity.Assertion) { value.Audience = "other" })}, "credential", failure.Unauthenticated},
		{"expired assertion", &providerStub{assertion: mutateAssertion(func(value *identity.Assertion) { value.ExpiresAt = testNow })}, "credential", failure.Unauthenticated},
		{"future assertion", &providerStub{assertion: mutateAssertion(func(value *identity.Assertion) {
			value.IssuedAt = testNow.Add(time.Minute)
			value.ExpiresAt = testNow.Add(time.Hour)
		})}, "credential", failure.Unauthenticated},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			verifier, err := identity.NewVerifier(test.provider, "forge-cli", func() time.Time { return testNow })
			if err != nil {
				t.Fatal(err)
			}
			_, err = verifier.Verify(context.Background(), verifyRequest(test.credential))
			assertFailureCode(t, err, test.want)
			if err != nil && errors.Is(err, test.provider.err) && test.provider.err != nil {
				t.Fatalf("provider error escaped boundary: %v", err)
			}
		})
	}
}

func TestPrincipalExpiresAtReceivingCapability(t *testing.T) {
	t.Parallel()
	verifier, err := identity.NewVerifier(&providerStub{assertion: validAssertion()}, "forge-cli", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), verifyRequest("credential"))
	if err != nil {
		t.Fatal(err)
	}
	assertFailureCode(t, principal.ValidateAt(testNow.Add(31*time.Minute), "correlation-1"), failure.Expired)
}

func TestVerifierRejectsUnsupportedRequestVersionBeforeProviderCall(t *testing.T) {
	t.Parallel()
	provider := &providerStub{assertion: validAssertion()}
	verifier, err := identity.NewVerifier(provider, "forge-cli", func() time.Time { return testNow })
	if err != nil {
		t.Fatal(err)
	}
	request := verifyRequest("credential")
	request.ContractVersion = "other"
	_, err = verifier.Verify(context.Background(), request)
	assertFailureCode(t, err, failure.InvalidRequest)
	if provider.seen != "" {
		t.Fatal("identity provider was called for an unsupported request version")
	}
}

func validAssertion() identity.Assertion {
	return identity.Assertion{
		SubjectID: "synthetic-developer-allowed", Kind: identity.HumanPrincipal,
		Issuer: "forge.synthetic", Audience: "forge-cli", Assurance: "synthetic-local",
		ProfileID: "forge.identity.synthetic.developer-permitted", IssuedAt: testNow.Add(-time.Minute),
		ExpiresAt: testNow.Add(30 * time.Minute), Synthetic: true,
	}
}

func mutateAssertion(mutate func(*identity.Assertion)) identity.Assertion {
	value := validAssertion()
	mutate(&value)
	return value
}

func verifyRequest(credential string) identity.VerifyRequest {
	return identity.VerifyRequest{ContractVersion: identity.VerifyRequestContractVersion, Credential: credential, CorrelationID: "correlation-1"}
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
		t.Fatalf("error code = %q, want %q", normalized.Code, want)
	}
	if validationErr := normalized.Validate(); validationErr != nil {
		t.Fatalf("normalized error invalid: %v", validationErr)
	}
}
