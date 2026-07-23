package experience_test

import (
	"context"
	"testing"
	"time"

	"github.com/3corners/forge/internal/experience"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/intent"
	"github.com/3corners/forge/internal/tenant"
)

var proposalTime = time.Date(2026, time.July, 23, 19, 0, 0, 0, time.UTC)

type identityProvider struct{ subject string }

func (provider identityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: provider.subject, Kind: identity.HumanPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: proposalTime.Add(-time.Minute), ExpiresAt: proposalTime.Add(time.Hour), Synthetic: true,
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

func TestRequirementsProposalSeparatesInputFromAuthorityAndBindsDigest(t *testing.T) {
	session := proposalSession(t, "human-1", "correlation-proposal")
	service, err := experience.NewService(experience.NewMemoryStore(), func() time.Time { return proposalTime }, func() string { return "proposal-1" })
	if err != nil {
		t.Fatal(err)
	}
	request := validProposal()
	proposal, err := service.Create(context.Background(), session, request)
	if err != nil {
		t.Fatal(err)
	}
	if proposal.ContractVersion != experience.RequirementsProposalVersion || proposal.Lifecycle != "proposal" ||
		!proposal.Current || proposal.Digest == "" ||
		proposal.Statements[0].Kind != "direct-user-statement" ||
		proposal.Statements[1].Kind != "model-inference" {
		t.Fatalf("unexpected proposal: %#v", proposal)
	}
	subject, found, err := service.ResolvePromotionSubject(context.Background(), proposal.TenantID, proposal.ID)
	if err != nil || !found || !intent.VerifyProposalDigest(subject) {
		t.Fatalf("promotion subject found=%t err=%v subject=%#v", found, err, subject)
	}
	if subject.ProposalDigest != proposal.Digest || subject.DirectUserStatements[0] != request.RawInput ||
		subject.Inferences[0] != request.Inferences[0] ||
		subject.ProductAreaRefs[0] != request.ProductAreaRefs[0] {
		t.Fatal("proposal digest, direct user statement, inference, or product area was not preserved")
	}
}

func TestProposalDigestChangesWithDirectUserStatementAndIsStableForSetOrdering(t *testing.T) {
	session := proposalSession(t, "human-1", "correlation-proposal")
	store := experience.NewMemoryStore()
	ids := []string{"proposal-1", "proposal-2", "proposal-3"}
	index := 0
	service, err := experience.NewService(store, func() time.Time { return proposalTime }, func() string {
		value := ids[index]
		index++
		return value
	})
	if err != nil {
		t.Fatal(err)
	}
	firstRequest := validProposal()
	first, err := service.Create(context.Background(), session, firstRequest)
	if err != nil {
		t.Fatal(err)
	}
	reordered := validProposal()
	reordered.Constraints = []string{"Synthetic data only", "No production access"}
	second, err := service.Create(context.Background(), session, reordered)
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest == second.Digest {
		t.Fatal("proposal identity and version must keep distinct proposals digest-bound")
	}
	firstSubject, _, _ := service.ResolvePromotionSubject(context.Background(), first.TenantID, first.ID)
	secondSubject, _, _ := service.ResolvePromotionSubject(context.Background(), second.TenantID, second.ID)
	firstSubject.ProposalID, secondSubject.ProposalID = "same", "same"
	if intent.BindProposalDigest(firstSubject).ProposalDigest != intent.BindProposalDigest(secondSubject).ProposalDigest {
		t.Fatal("canonical set ordering changed the proposal digest")
	}
	changed := validProposal()
	changed.RawInput = "Build a different application"
	third, err := service.Create(context.Background(), session, changed)
	if err != nil {
		t.Fatal(err)
	}
	if third.Digest == first.Digest {
		t.Fatal("changed direct user statement did not change the digest")
	}
}

func validProposal() experience.CreateProposal {
	return experience.CreateProposal{
		RawInput: "Build a guided Forge journey", Outcome: "A user can guide one bounded Forge change",
		AffectedUsers: []string{"developer"}, Problem: "The user cannot yet try Forge",
		FunctionalBehavior: []string{"Create a requirements proposal", "Require exact promotion"},
		UXExpectations:     []string{"Immediate acknowledgement", "Visible authority gates"},
		Accessibility:      []string{"Keyboard-operable CLI"}, Privacy: []string{"No prompt content in telemetry"},
		Security: []string{"Raw text is not authority"}, Operations: []string{"Deterministic retry"},
		Constraints:        []string{"No production access", "Synthetic data only"},
		Exclusions:         []string{"Live GitHub credentials"},
		AcceptanceCriteria: []string{"Exact proposal digest is promoted"},
		Inferences:         []string{"The user likely needs a resumable session"},
		Assumptions:        []string{"The Developer Trial profile is active"},
		Hypotheses:         []string{"A guided CLI reduces requirement ambiguity"},
		EvidenceReferences: []string{"adr:0020"},
		Questions:          []intent.Question{{ID: "q-1", Text: "Which provider?", Consequence: "", Material: false, Resolved: false}},
		Experiments:        []string{"Disposable contract-double preview"},
		ProductAreaRefs:    []string{"forge.discovery"}, Provenance: []string{"channel:cli"},
		ChangeExplanation: "Initial proposal",
	}
}

func proposalSession(t *testing.T, subject, correlation string) experience.Session {
	t.Helper()
	verifier, err := identity.NewVerifier(identityProvider{subject: subject}, "forge-cli", func() time.Time { return proposalTime })
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
		bindingSource{}, func() time.Time { return proposalTime },
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
	session, err := experience.NewSession(principal, scope, binding, proposalTime, correlation)
	if err != nil {
		t.Fatal(err)
	}
	return session
}
