package intent_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/experience"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/intent"
	"github.com/3corners/forge/internal/tenant"
)

var promotionTime = time.Date(2026, time.July, 23, 19, 30, 0, 0, time.UTC)

type identityProvider struct{ subject string }

func (provider identityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: provider.subject, Kind: identity.HumanPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: promotionTime.Add(-time.Minute), ExpiresAt: promotionTime.Add(time.Hour), Synthetic: true,
	}, nil
}

type bindingSource struct{}

func (bindingSource) ResolveRepositoryBinding(context.Context, string, string, string, string) (tenant.RepositoryBindingDefinition, bool, error) {
	return tenant.RepositoryBindingDefinition{
		TenantID: "tenant-fixture", WorkspaceID: "workspace-fixture", ProjectID: "project-forge",
		RepositoryID: "repository-forge", Selector: "forge", CanonicalLocator: "contract-test://forge",
		Classification: "Internal",
	}, true, nil
}

func TestExactCurrentProposalPromotionCreatesCanonicalIntentAndDeliveryRun(t *testing.T) {
	fixture := newFixture(t, nil)
	result, err := fixture.intents.Promote(context.Background(), fixture.request)
	if err != nil {
		t.Fatal(err)
	}
	if result.Intent.ContractVersion != intent.CanonicalIntentVersion ||
		result.Intent.ProposalDigest != fixture.request.ProposalDigest ||
		result.Run.ID() != result.Intent.DeliveryRunID || result.Run.State() != delivery.CapturedState {
		t.Fatalf("unexpected promotion result: %#v", result)
	}
	replayed, err := fixture.intents.Promote(context.Background(), fixture.request)
	if err != nil || !replayed.Replayed || replayed.Intent.ID != result.Intent.ID || replayed.Run.ID() != result.Run.ID() {
		t.Fatalf("promotion replay=%#v err=%v", replayed, err)
	}
}

func TestPromotionFailsClosedForDigestScopeAndMaterialQuestion(t *testing.T) {
	tests := []struct {
		name     string
		question []intent.Question
		mutate   func(*intent.PromoteRequest)
		want     failure.Code
	}{
		{
			name: "changed digest", mutate: func(value *intent.PromoteRequest) {
				value.ProposalDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			}, want: failure.Forbidden,
		},
		{
			name: "material unresolved question",
			question: []intent.Question{{
				ID: "q-risk", Text: "May customer data be used?", Consequence: "Changes data classification",
				Material: true, Resolved: false,
			}},
			mutate: func(*intent.PromoteRequest) {}, want: failure.Conflict,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t, test.question)
			request := fixture.request
			test.mutate(&request)
			_, err := fixture.intents.Promote(context.Background(), request)
			assertFailure(t, err, test.want)
		})
	}
}

func TestPromotionRejectsSupersededProposalAndAcceptsExactRevision(t *testing.T) {
	fixture := newFixture(t, nil)
	originalDigest := fixture.request.ProposalDigest
	revised, err := fixture.proposals.Revise(
		context.Background(), fixture.session, fixture.request.ProposalID, originalDigest,
		experience.ProposalEdit{
			Field: "security", Item: 1, Replacement: "Treat raw input and model output as untrusted",
			ChangeExplanation: "Clarify the trust requirement.",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fixture.intents.Promote(context.Background(), fixture.request)
	assertFailure(t, err, failure.Forbidden)
	fixture.request.ProposalDigest = revised.Digest
	result, err := fixture.intents.Promote(context.Background(), fixture.request)
	if err != nil || result.Intent.ProposalDigest != revised.Digest {
		t.Fatalf("revised promotion=%#v err=%v", result, err)
	}
}

type fixture struct {
	intents   *intent.Service
	proposals *experience.Service
	session   experience.Session
	request   intent.PromoteRequest
}

func newFixture(t *testing.T, questions []intent.Question) fixture {
	t.Helper()
	principal, scope, binding := authority(t, "human-1", "correlation-promotion")
	session, err := experience.NewSession(principal, scope, binding, promotionTime, "correlation-promotion")
	if err != nil {
		t.Fatal(err)
	}
	proposals, err := experience.NewService(experience.NewMemoryStore(), func() time.Time { return promotionTime }, func() string { return "proposal-1" })
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := proposals.Create(context.Background(), session, experience.CreateProposal{
		RawInput: "Build the guided Forge journey", Outcome: "A user can guide one bounded Forge change",
		AffectedUsers: []string{"developer"}, Problem: "Forge has no public delivery journey",
		FunctionalBehavior: []string{"Create a proposal"}, UXExpectations: []string{"Show the digest"},
		Accessibility: []string{}, Privacy: []string{}, Security: []string{"Raw input is not authority"},
		Operations: []string{}, Constraints: []string{"No production access"},
		Exclusions: []string{"Live GitHub credentials"}, AcceptanceCriteria: []string{"Create one canonical intent"},
		Questions: questions, Provenance: []string{"channel:cli"}, ChangeExplanation: "Initial proposal",
	})
	if err != nil {
		t.Fatal(err)
	}
	deliveryStore, err := delivery.NewMemoryStore(func() string { return "delivery-run-1" })
	if err != nil {
		t.Fatal(err)
	}
	deliveryService, err := delivery.NewService(deliveryStore, func() time.Time { return promotionTime })
	if err != nil {
		t.Fatal(err)
	}
	intents, err := intent.NewService(
		proposals, deliveryService, intent.NewMemoryStore(), func() time.Time { return promotionTime },
		func() string { return "intent-1" },
	)
	if err != nil {
		t.Fatal(err)
	}
	return fixture{
		intents: intents, proposals: proposals, session: session,
		request: intent.PromoteRequest{
			ContractVersion: intent.PromotionRequestVersion, CorrelationID: "correlation-promotion",
			IdempotencyKey: "promotion-1", ProposalID: proposal.ID, ProposalDigest: proposal.Digest,
			Principal: principal, Scope: scope, RepositoryBinding: binding,
		},
	}
}

func authority(t *testing.T, subject, correlation string) (identity.Principal, tenant.Scope, tenant.RepositoryBinding) {
	t.Helper()
	verifier, err := identity.NewVerifier(identityProvider{subject: subject}, "forge-cli", func() time.Time { return promotionTime })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "synthetic", CorrelationID: correlation,
	})
	if err != nil {
		t.Fatal(err)
	}
	directory, err := tenant.NewDirectory(
		"tenant-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-fixture", ID: "workspace-fixture", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-fixture", ID: "project-forge", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: subject, WorkspaceID: "workspace-fixture"}},
		bindingSource{}, func() time.Time { return promotionTime },
	)
	if err != nil {
		t.Fatal(err)
	}
	scope, err := directory.ResolveWorkspaceScope(principal, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: correlation,
		WorkspaceSelector: "fixture-workspace", ProjectSelector: "forge", RepositorySelector: "forge",
	})
	if err != nil {
		t.Fatal(err)
	}
	binding, err := directory.GetRepositoryBinding(context.Background(), tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: correlation, Scope: scope,
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal, scope, binding
}

func assertFailure(t *testing.T, err error, want failure.Code) {
	t.Helper()
	var normalized *failure.Error
	if err == nil || !errors.As(err, &normalized) || normalized.Code != want {
		t.Fatalf("error=%v want code=%s", err, want)
	}
}
