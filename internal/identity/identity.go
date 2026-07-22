// Package identity verifies provider assertions and returns normalized,
// immutable principals. It does not own membership or final authorization.
package identity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
)

const (
	VerifyRequestContractVersion = "forge.identity-verify-request.v1"
	PrincipalContractVersion     = "forge.normalized-principal.v1"
)

type PrincipalKind string

const (
	HumanPrincipal    PrincipalKind = "human"
	WorkloadPrincipal PrincipalKind = "workload"
	AgentPrincipal    PrincipalKind = "agent"
)

// ErrProviderUnavailable lets an identity adapter distinguish an unavailable
// verifier from an invalid assertion without exposing provider diagnostics.
var ErrProviderUnavailable = errors.New("identity provider unavailable")

// Assertion contains verified facts returned by an identity adapter. It must
// never contain the original credential or other reusable secret material.
type Assertion struct {
	SubjectID string
	Kind      PrincipalKind
	Issuer    string
	Audience  string
	Assurance string
	ProfileID string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Synthetic bool
}

// Provider verifies an opaque credential and returns normalized assertion
// facts. Provider selection and credential acquisition remain outside this
// contract.
type Provider interface {
	Verify(context.Context, string) (Assertion, error)
}

type VerifyRequest struct {
	ContractVersion string
	Credential      string
	CorrelationID   string
}

// Principal can only be created by a successful Verifier call. Its fields are
// intentionally private so caller input cannot construct authoritative facts.
type Principal struct {
	contractVersion string
	subjectID       string
	kind            PrincipalKind
	issuer          string
	audience        string
	assurance       string
	profileID       string
	issuedAt        time.Time
	expiresAt       time.Time
	synthetic       bool
}

func (value Principal) ContractVersion() string { return value.contractVersion }
func (value Principal) SubjectID() string       { return value.subjectID }
func (value Principal) Kind() PrincipalKind     { return value.kind }
func (value Principal) Issuer() string          { return value.issuer }
func (value Principal) Audience() string        { return value.audience }
func (value Principal) Assurance() string       { return value.assurance }
func (value Principal) ProfileID() string       { return value.profileID }
func (value Principal) IssuedAt() time.Time     { return value.issuedAt }
func (value Principal) ExpiresAt() time.Time    { return value.expiresAt }
func (value Principal) Synthetic() bool         { return value.synthetic }

// ValidateAt lets receiving capabilities recheck principal lifecycle without
// gaining access to assertion verification or credential material.
func (value Principal) ValidateAt(at time.Time, correlationID string) error {
	if value.contractVersion != PrincipalContractVersion || value.subjectID == "" || value.issuer == "" || value.audience == "" || value.assurance == "" || value.profileID == "" || !validKind(value.kind) {
		return failure.New(failure.Unauthenticated, "authenticated principal is invalid", correlationID, failure.RetryReauthenticate)
	}
	if at.IsZero() {
		return failure.New(failure.InvalidRequest, "principal validation time is required", correlationID, failure.RetryNone)
	}
	if at.Before(value.issuedAt) || !at.Before(value.expiresAt) {
		return failure.New(failure.Expired, "authenticated principal is not active", correlationID, failure.RetryReauthenticate)
	}
	return nil
}

type Verifier struct {
	provider Provider
	audience string
	now      func() time.Time
}

func NewVerifier(provider Provider, audience string, now func() time.Time) (*Verifier, error) {
	if provider == nil {
		return nil, errors.New("identity provider is required")
	}
	if audience == "" {
		return nil, errors.New("identity audience is required")
	}
	if now == nil {
		return nil, errors.New("identity clock is required")
	}
	return &Verifier{provider: provider, audience: audience, now: now}, nil
}

func (verifier *Verifier) Verify(ctx context.Context, request VerifyRequest) (Principal, error) {
	if verifier == nil || verifier.provider == nil || verifier.now == nil || verifier.audience == "" {
		return Principal{}, failure.New(failure.DependencyUnavailable, "identity verification is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if request.ContractVersion != VerifyRequestContractVersion {
		return Principal{}, failure.New(failure.InvalidRequest, "identity verification request version is invalid", request.CorrelationID, failure.RetryNone)
	}
	if ctx == nil || request.Credential == "" {
		return Principal{}, failure.New(failure.Unauthenticated, "identity assertion is required", request.CorrelationID, failure.RetryReauthenticate)
	}
	assertion, err := verifier.provider.Verify(ctx, request.Credential)
	if err != nil {
		if errors.Is(err, ErrProviderUnavailable) {
			return Principal{}, failure.New(failure.DependencyUnavailable, "identity verification is unavailable", request.CorrelationID, failure.RetryBounded)
		}
		return Principal{}, failure.New(failure.Unauthenticated, "identity assertion is invalid", request.CorrelationID, failure.RetryReauthenticate)
	}
	if err := validateAssertion(assertion, verifier.audience, verifier.now()); err != nil {
		return Principal{}, failure.New(failure.Unauthenticated, "identity assertion is invalid", request.CorrelationID, failure.RetryReauthenticate)
	}
	return Principal{
		contractVersion: PrincipalContractVersion,
		subjectID:       assertion.SubjectID,
		kind:            assertion.Kind,
		issuer:          assertion.Issuer,
		audience:        assertion.Audience,
		assurance:       assertion.Assurance,
		profileID:       assertion.ProfileID,
		issuedAt:        assertion.IssuedAt,
		expiresAt:       assertion.ExpiresAt,
		synthetic:       assertion.Synthetic,
	}, nil
}

func validateAssertion(value Assertion, audience string, at time.Time) error {
	if value.SubjectID == "" || value.Issuer == "" || value.Audience == "" || value.Assurance == "" || value.ProfileID == "" {
		return errors.New("assertion has an empty required field")
	}
	if !validKind(value.Kind) {
		return fmt.Errorf("unsupported principal kind %q", value.Kind)
	}
	if value.Audience != audience {
		return errors.New("assertion audience mismatch")
	}
	if value.IssuedAt.IsZero() || value.ExpiresAt.IsZero() || !value.ExpiresAt.After(value.IssuedAt) {
		return errors.New("assertion validity window is invalid")
	}
	if at.IsZero() || at.Before(value.IssuedAt) || !at.Before(value.ExpiresAt) {
		return errors.New("assertion is not active")
	}
	return nil
}

func validKind(value PrincipalKind) bool {
	switch value {
	case HumanPrincipal, WorkloadPrincipal, AgentPrincipal:
		return true
	default:
		return false
	}
}
