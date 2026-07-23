package sourcecontrol

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

var publicationNow = time.Date(2026, time.July, 23, 15, 0, 0, 0, time.UTC)
var candidateDigest = "sha256:" + strings.Repeat("a", 64)
var treeDigest = "sha256:" + strings.Repeat("b", 64)
var verificationDigest = "sha256:" + strings.Repeat("c", 64)

type subjectStub struct {
	subject trust.PublicationApprovalSubject
}

func (stub *subjectStub) ResolvePublicationSubject(context.Context, string, string) (trust.PublicationApprovalSubject, bool, error) {
	return stub.subject, true, nil
}

type approvalStub struct{ record trust.ApprovalRecordV2 }

func (stub approvalStub) ResolveApproval(context.Context, string, string) (trust.ApprovalRecordV2, bool, error) {
	return stub.record, true, nil
}

type policyStub struct{ decision trust.PolicyDecision }

func (stub policyStub) ResolvePolicyDecision(context.Context, string) (trust.PolicyDecision, bool, error) {
	return stub.decision, true, nil
}

type evidenceStub struct{ err error }

func (stub evidenceStub) VerifySubject(context.Context, trust.EvidenceScope, trust.EvidenceSubject, []string, string) error {
	return stub.err
}

type identityProvider struct{}

func (identityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: "human-1", Kind: identity.HumanPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: publicationNow.Add(-time.Minute), ExpiresAt: publicationNow.Add(time.Hour), Synthetic: true,
	}, nil
}

type bindingSource struct{}

func (bindingSource) ResolveRepositoryBinding(context.Context, string, string, string, string) (tenant.RepositoryBindingDefinition, bool, error) {
	return tenant.RepositoryBindingDefinition{
		TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1",
		RepositoryID: "repository-1", Selector: "forge", CanonicalLocator: "contract-test://forge",
		Classification: "Internal",
	}, true, nil
}

func TestPublicationCreatesOnePullRequestAndReplaysExactOutcome(t *testing.T) {
	service, request, provider, _ := publicationFixture(t, FailNone)
	first, replayed, err := service.PublishPullRequest(context.Background(), request)
	if err != nil || replayed {
		t.Fatalf("first publication replay=%t receipt=%#v err=%v", replayed, first, err)
	}
	if first.PullRequestID == "" || first.PullRequestURL == "" || first.Branch == "" {
		t.Fatalf("receipt is incomplete: %#v", first)
	}
	second, replayed, err := service.PublishPullRequest(context.Background(), request)
	if err != nil || !replayed || second.PullRequestID != first.PullRequestID || second.OperationID != first.OperationID {
		t.Fatalf("replay=%t receipt=%#v err=%v", replayed, second, err)
	}
	if provider.APIRequests() > 8 {
		t.Fatalf("unbounded provider calls: %d", provider.APIRequests())
	}
}

func TestPublicationReconcilesTimeoutAfterRemoteMutation(t *testing.T) {
	for _, point := range []FailurePoint{FailBeforeBranch, FailAfterBranch, FailBeforePullRequest, FailAfterPullRequest} {
		t.Run(string(point), func(t *testing.T) {
			service, request, _, _ := publicationFixture(t, point)
			receipt, _, err := service.PublishPullRequest(context.Background(), request)
			if err != nil {
				t.Fatal(err)
			}
			if receipt.PullRequestID == "" || !receipt.Reconciled {
				t.Fatalf("ambiguous outcome was not reconciled: %#v", receipt)
			}
		})
	}
}

func TestPublicationFailsClosedForWrongApprovalEvidenceAndScope(t *testing.T) {
	service, request, _, subjects := publicationFixture(t, FailNone)
	planApproval := request
	planApproval.ApprovalRecordID = "plan-approval"
	approval, _ := validApproval(subjects.subject)
	approval.Action = "forge.delivery.execute-plan"
	approval.SubjectType = "delivery-plan"
	approval.ID = ""
	approval.Integrity = trust.Integrity{}
	approval, _ = trust.BindApprovalRecordV2Integrity(approval, "correlation-1")
	service.approvals = approvalStub{record: approval}
	_, _, err := service.PublishPullRequest(context.Background(), planApproval)
	assertFailure(t, err, failure.Forbidden)

	service, request, _, _ = publicationFixture(t, FailNone)
	service.evidence = evidenceStub{err: errors.New("not ready")}
	_, _, err = service.PublishPullRequest(context.Background(), request)
	assertFailure(t, err, failure.EvidenceUnavailable)

	service, request, _, subjects = publicationFixture(t, FailNone)
	subjects.subject.Scope.RepositoryID = "other-repository"
	_, _, err = service.PublishPullRequest(context.Background(), request)
	assertFailure(t, err, failure.Forbidden)
}

func TestPublicationFailsClosedForChangedCandidateMaterialAndMetadata(t *testing.T) {
	tests := map[string]func(*PublishRequest){
		"tree": func(request *PublishRequest) {
			request.TreeDigest = "sha256:" + strings.Repeat("d", 64)
		},
		"paths": func(request *PublishRequest) {
			request.ChangedPaths = []string{"internal/other/other.go"}
		},
		"title": func(request *PublishRequest) {
			request.Title = "Unreviewed title"
		},
		"body": func(request *PublishRequest) {
			request.Body = "Unreviewed body"
		},
		"base": func(request *PublishRequest) {
			request.BaseRevision = "different-base"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			service, request, _, _ := publicationFixture(t, FailNone)
			mutate(&request)
			_, _, err := service.PublishPullRequest(context.Background(), request)
			assertFailure(t, err, failure.Forbidden)
		})
	}
}

func TestPublicationConflictingRemoteStateRequiresReconciliation(t *testing.T) {
	service, request, provider, subjects := publicationFixture(t, FailNone)
	operation := buildOperation(request, subjects.subject)
	if err := provider.ReplaceState(operation.ID, RemoteState{
		BranchExists: true, Branch: operation.Branch, CommitID: "foreign",
		CandidateDigest: "sha256:" + strings.Repeat("d", 64), TreeDigest: operation.TreeDigest,
		BaseBranch: operation.BaseBranch, Marker: operation.Marker, RepositoryID: operation.RepositoryID,
	}); err != nil {
		t.Fatal(err)
	}
	_, _, err := service.PublishPullRequest(context.Background(), request)
	assertFailure(t, err, failure.Conflict)
}

func publicationFixture(t *testing.T, failurePoint FailurePoint) (*Service, PublishRequest, *ContractDouble, *subjectStub) {
	t.Helper()
	principal := verifiedPrincipal(t)
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-1", ID: "workspace-1", Selector: "workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-1", ID: "project-1", Selector: "forge"}},
		[]tenant.MembershipDefinition{{PrincipalID: "human-1", WorkspaceID: "workspace-1"}},
		bindingSource{}, func() time.Time { return publicationNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	scope, err := directory.ResolveWorkspaceScope(principal, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: "correlation-1",
		WorkspaceSelector: "workspace", ProjectSelector: "forge", RepositorySelector: "forge",
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
	subjects := &subjectStub{subject: validSubject()}
	approval, err := validApproval(subjects.subject)
	if err != nil {
		t.Fatal(err)
	}
	decision := validDecision()
	envelope, err := validEnvelope(approval, decision)
	if err != nil {
		t.Fatal(err)
	}
	provider := NewContractDouble(failurePoint)
	profile := PublicationProfile{
		Provider: "contract-double", BaseBranch: "main", BaseRevision: "base-revision-1",
		Title: "Forge M1 candidate",
		Body:  "Verified M1 candidate. Human review and merge remain required.",
	}
	service, err := NewService(subjects, approvalStub{record: approval}, policyStub{decision: decision}, evidenceStub{}, provider, NewMemoryStore(), profile, func() time.Time { return publicationNow })
	if err != nil {
		t.Fatal(err)
	}
	return service, PublishRequest{
		ContractVersion: PublishRequestContractVersion, CorrelationID: "correlation-1",
		IdempotencyKey: "publication-key", TenantID: "tenant-1", RunID: "run-1",
		Envelope: envelope, PolicyDecisionID: decision.ID, ApprovalRecordID: approval.ID,
		RepositoryBinding: binding, Provider: "contract-double", BaseBranch: "main",
		BaseRevision: "base-revision-1", TreeDigest: treeDigest,
		ChangedPaths: []string{"internal/example/example.go"}, Title: "Forge M1 candidate",
		Body: "Verified M1 candidate. Human review and merge remain required.",
	}, provider, subjects
}

func validSubject() trust.PublicationApprovalSubject {
	return trust.PublicationApprovalSubject{
		RunID: "run-1", CandidateID: "run-1.candidate", CandidateDigest: candidateDigest,
		TreeDigest: treeDigest, ChangedPaths: []string{"internal/example/example.go"},
		Scope: trust.ApprovalScope{
			TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1",
			RepositoryID: "repository-1", DeliveryRunID: "run-1",
		},
		RequestingActorID: "human-1", Purpose: "forge.m1.prepare-change",
		Classification: "Internal", RiskClass: "low", VerificationResultID: "verification-1",
		VerificationResultDigest: verificationDigest,
		EvidenceReferences:       []string{"candidate:" + candidateDigest, "verification:" + verificationDigest},
		PolicyDecisionIDs:        []string{"candidate-policy"}, ExpiresAt: publicationNow.Add(30 * time.Minute),
		CorrelationID: "correlation-1",
	}
}

func validApproval(subject trust.PublicationApprovalSubject) (trust.ApprovalRecordV2, error) {
	request := trust.ApprovalRequest{
		SchemaVersion: trust.ApprovalRequestSchemaVersion, Issuer: "forge.approval",
		RequestingActorID: "human-1", IntendedApproverID: "human-1", Scope: subject.Scope,
		Environment: DeveloperTrialEnvironment, Purpose: subject.Purpose, Classification: subject.Classification,
		RiskClass: subject.RiskClass, Action: PublishAction, SubjectType: "verified-candidate",
		SubjectID: subject.CandidateID, SubjectDigest: subject.CandidateDigest,
		EvidenceReferences: subject.EvidenceReferences, PermittedConditions: []string{},
		PolicyDecisionIDs: subject.PolicyDecisionIDs, CreatedAt: publicationNow.Add(-time.Minute),
		ExpiresAt: publicationNow.Add(29 * time.Minute), MaxDecisionLifetime: 1800,
		Lifecycle: "open", EvidenceCorrelationID: "correlation-1",
	}
	request, err := trust.BindApprovalRequestDigest(request, "correlation-1")
	if err != nil {
		return trust.ApprovalRecordV2{}, err
	}
	return trust.BindApprovalRecordV2Integrity(trust.ApprovalRecordV2{
		SchemaVersion: trust.ApprovalRecordV2SchemaVersion, Issuer: "forge.approval",
		RequestID: request.ID, RequestDigest: request.RequestDigest, ApproverID: "human-1",
		Scope: subject.Scope, Environment: DeveloperTrialEnvironment, Purpose: subject.Purpose,
		Classification: subject.Classification, RiskClass: subject.RiskClass,
		Action: PublishAction, SubjectType: "verified-candidate", SubjectID: subject.CandidateID,
		SubjectDigest: subject.CandidateDigest, Decision: "approved", DecidedAt: publicationNow.Add(-time.Minute),
		ExpiresAt: publicationNow.Add(20 * time.Minute), Conditions: []string{},
		PolicyDecisionIDs: []string{"approval-policy"}, EvidenceReferences: subject.EvidenceReferences,
		EvidenceCorrelationID: "correlation-1",
	}, "correlation-1")
}

func validDecision() trust.PolicyDecision {
	expires := publicationNow.Add(15 * time.Minute)
	return trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, ID: "publication-policy",
		DecisionPoint: "publication", SubjectID: "workload-1", Action: PublishAction,
		Resource: candidateDigest, EvaluatedAt: publicationNow.Add(-time.Minute), ExpiresAt: &expires,
		PolicyVersions: []string{"policy-v1"}, Outcome: "permit",
		Obligations:           []map[string]any{{"type": "evidence-ready"}, {"type": "human-publication-approval"}},
		EvidenceCorrelationID: "correlation-1",
	}
}

func validEnvelope(approval trust.ApprovalRecordV2, decision trust.PolicyDecision) (trust.ExecutionEnvelope, error) {
	return trust.BindExecutionEnvelopeIntegrity(trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion,
		Issuer:        "forge.policy", IssuedAt: publicationNow.Add(-time.Minute), ExpiresAt: publicationNow.Add(10 * time.Minute),
		Scope:     trust.ExecutionScope{TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", DeliveryRunID: "run-1"},
		Principal: trust.EnvelopePrincipal{SubjectID: "workload-1", AgentID: "source-control-contract-double", DelegationChain: []string{}},
		Purpose:   "forge.m1.prepare-change", RiskClass: "low", Actions: []string{PublishAction},
		Classifications: []string{"Internal"}, Destinations: []string{"contract-double:repository-1"},
		PolicyDecisionIDs: []string{decision.ID}, ApprovalRecordIDs: []string{approval.ID},
		Limits:         trust.ExecutionLimits{MaxDurationSeconds: 60, MaxCostMinorUnits: 0, MaxConcurrency: 1},
		IdempotencyKey: "publication-key", EvidenceCorrelationID: "correlation-1",
	}, "correlation-1")
}

func verifiedPrincipal(t *testing.T) identity.Principal {
	t.Helper()
	verifier, err := identity.NewVerifier(identityProvider{}, "forge-cli", func() time.Time { return publicationNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque", CorrelationID: "correlation-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	return principal
}

func assertFailure(t *testing.T, err error, code failure.Code) {
	t.Helper()
	var normalized *failure.Error
	if err == nil || !errors.As(err, &normalized) || normalized.Code != code {
		t.Fatalf("error=%v want=%s", err, code)
	}
}
