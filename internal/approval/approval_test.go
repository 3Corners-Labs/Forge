package approval

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/identity"
)

var approvalNow = time.Date(2026, time.July, 23, 12, 0, 0, 0, time.UTC)
var approvalSubjectDigest = "sha256:" + strings.Repeat("a", 64)

type subjectResolverStub struct {
	subject PlanSubject
	found   bool
	err     error
}

func (resolver subjectResolverStub) ResolvePlanSubject(context.Context, string, string) (PlanSubject, bool, error) {
	return resolver.subject, resolver.found, resolver.err
}

type policyResolverStub struct {
	decisions map[string]trust.PolicyDecision
	err       error
}

func (resolver *policyResolverStub) ResolvePolicyDecision(_ context.Context, id string) (trust.PolicyDecision, bool, error) {
	if resolver.err != nil {
		return trust.PolicyDecision{}, false, resolver.err
	}
	value, found := resolver.decisions[id]
	return value, found, nil
}

type identityProviderStub struct{}

func (identityProviderStub) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: "human-1", Kind: identity.HumanPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: approvalNow.Add(-time.Minute), ExpiresAt: approvalNow.Add(time.Hour), Synthetic: true,
	}, nil
}

func TestApprovalLifecycleBindsHumanPolicyScopeAndCurrentChain(t *testing.T) {
	service, policies, store, principal := approvalFixture(t)
	ctx := context.Background()
	request, replayed, err := service.RequestPlanReview(ctx, planReviewCommand("request-key"))
	if err != nil || replayed {
		t.Fatalf("request replay=%t err=%v", replayed, err)
	}
	policies.decisions["eligibility-1"] = eligibilityDecision(request, "permit")
	command := DecideApproval{
		ContractVersion: DecideApprovalContractVersion, CorrelationID: "correlation-1",
		IdempotencyKey: "decision-key", TenantID: "tenant-1", RequestID: request.ID,
		Approver: principal, PolicyDecisionID: "eligibility-1", Decision: "approved",
		Conditions:         []string{"require-independent-verification"},
		EvidenceReferences: request.EvidenceReferences, RequestedExpiry: approvalNow.Add(20 * time.Minute),
	}
	record, replayed, err := service.DecideApproval(ctx, command)
	if err != nil || replayed || record.Decision != "approved" {
		t.Fatalf("decision=%#v replay=%t err=%v", record, replayed, err)
	}
	if err := record.VerifyIntegrity("correlation-1"); err != nil {
		t.Fatal(err)
	}
	replayedRecord, replayed, err := service.DecideApproval(ctx, command)
	if err != nil || !replayed || replayedRecord.ID != record.ID {
		t.Fatalf("decision replay=%t id=%q err=%v", replayed, replayedRecord.ID, err)
	}
	if resolved, found, err := service.ResolveApproval(ctx, "tenant-1", record.ID); err != nil || !found || resolved.ID != record.ID {
		t.Fatalf("resolve found=%t id=%q err=%v", found, resolved.ID, err)
	}
	if _, found, err := service.ResolveApproval(ctx, "other-tenant", record.ID); err != nil || found {
		t.Fatalf("cross-tenant resolution found=%t err=%v", found, err)
	}

	revoked, replayed, err := service.RevokeApproval(ctx, RevokeApproval{
		ContractVersion: RevokeApprovalContractVersion, CorrelationID: "correlation-1",
		IdempotencyKey: "revoke-key", TenantID: "tenant-1", RecordID: record.ID,
		Approver: principal, PolicyDecisionID: "eligibility-1",
	})
	if err != nil || replayed || revoked.Decision != "revoked" || revoked.PreviousRecordID != record.ID {
		t.Fatalf("revocation=%#v replay=%t err=%v", revoked, replayed, err)
	}
	if _, found, err := service.ResolveApproval(ctx, "tenant-1", record.ID); err != nil || found {
		t.Fatalf("stale branch remained current: found=%t err=%v", found, err)
	}
	if resolved, found, err := store.ResolveCurrent(ctx, "tenant-1", revoked.ID); err != nil || !found || resolved.Decision != "revoked" {
		t.Fatalf("revoked current state found=%t decision=%q err=%v", found, resolved.Decision, err)
	}
}

func TestApprovalFailsClosedForReplayEvidenceAndEligibilityMismatch(t *testing.T) {
	service, policies, _, principal := approvalFixture(t)
	ctx := context.Background()
	request, _, err := service.RequestPlanReview(ctx, planReviewCommand("request-key"))
	if err != nil {
		t.Fatal(err)
	}
	policies.decisions["eligibility-1"] = eligibilityDecision(request, "permit")
	base := DecideApproval{
		ContractVersion: DecideApprovalContractVersion, CorrelationID: "correlation-1",
		IdempotencyKey: "decision-key", TenantID: "tenant-1", RequestID: request.ID,
		Approver: principal, PolicyDecisionID: "eligibility-1", Decision: "approved",
		Conditions: []string{"require-independent-verification"}, EvidenceReferences: request.EvidenceReferences,
	}
	changedEvidence := base
	changedEvidence.EvidenceReferences = []string{"plan:" + approvalSubjectDigest}
	_, _, err = service.DecideApproval(ctx, changedEvidence)
	assertApprovalCode(t, err, failure.Forbidden)

	unsupportedCondition := base
	unsupportedCondition.Conditions = []string{"widen-scope"}
	_, _, err = service.DecideApproval(ctx, unsupportedCondition)
	assertApprovalCode(t, err, failure.Forbidden)

	indeterminate := base
	indeterminate.PolicyDecisionID = "indeterminate"
	policies.decisions["indeterminate"] = eligibilityDecision(request, "indeterminate")
	_, _, err = service.DecideApproval(ctx, indeterminate)
	assertApprovalCode(t, err, failure.Indeterminate)

	if _, _, err := service.DecideApproval(ctx, base); err != nil {
		t.Fatal(err)
	}
	conflicting := base
	conflicting.Decision = "rejected"
	_, _, err = service.DecideApproval(ctx, conflicting)
	assertApprovalCode(t, err, failure.Conflict)
}

func BenchmarkApprovalResolution(b *testing.B) {
	service, record := approvedFixture(b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, found, err := service.ResolveApproval(context.Background(), "tenant-1", record.ID); err != nil || !found {
			b.Fatal(err)
		}
	}
}

func BenchmarkApprovalResolutionParallel(b *testing.B) {
	service, record := approvedFixture(b)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, found, err := service.ResolveApproval(context.Background(), "tenant-1", record.ID); err != nil || !found {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRevokedApprovalResolution(b *testing.B) {
	service, record := approvedFixture(b)
	// Revocation performance is covered by current-record lookup after a
	// directly stored, integrity-bound immutable revocation.
	revocation := record
	revocation.ID = ""
	revocation.Integrity = trust.Integrity{}
	revocation.Decision = "revoked"
	revocation.PreviousRecordID = record.ID
	revocation.DecidedAt = approvalNow.Add(time.Second)
	revocation, err := trust.BindApprovalRecordV2Integrity(revocation, "correlation-1")
	if err != nil {
		b.Fatal(err)
	}
	if _, _, err := service.store.CreateDecision(context.Background(), "tenant-1", "benchmark-revoke", record.RequestID, revocation); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, found, err := service.ResolveApproval(context.Background(), "tenant-1", revocation.ID); err != nil || !found {
			b.Fatal(err)
		}
	}
}

func TestApprovalResolutionLatencyDistribution(t *testing.T) {
	service, record := approvedFixture(t)
	const samples = 2000
	durations := make([]time.Duration, samples)
	for index := range durations {
		started := time.Now()
		if _, found, err := service.ResolveApproval(context.Background(), "tenant-1", record.ID); err != nil || !found {
			t.Fatal(err)
		}
		durations[index] = time.Since(started)
	}
	sort.Slice(durations, func(left, right int) bool { return durations[left] < durations[right] })
	p50, p95, p99 := durations[samples*50/100], durations[samples*95/100], durations[samples*99/100]
	t.Logf("approval resolution latency p50=%s p95=%s p99=%s samples=%d", p50, p95, p99, samples)
	if p99 >= 500*time.Millisecond {
		t.Fatalf("approval resolution p99=%s exceeds the M1 control-operation hypothesis", p99)
	}
}

func approvedFixture(t testing.TB) (*Service, trust.ApprovalRecordV2) {
	t.Helper()
	service, policies, _, principal := approvalFixture(t)
	request, _, err := service.RequestPlanReview(context.Background(), planReviewCommand("request-key"))
	if err != nil {
		t.Fatal(err)
	}
	policies.decisions["eligibility-1"] = eligibilityDecision(request, "permit")
	record, _, err := service.DecideApproval(context.Background(), DecideApproval{
		ContractVersion: DecideApprovalContractVersion, CorrelationID: "correlation-1",
		IdempotencyKey: "decision-key", TenantID: "tenant-1", RequestID: request.ID,
		Approver: principal, PolicyDecisionID: "eligibility-1", Decision: "approved",
		Conditions: []string{"require-independent-verification"}, EvidenceReferences: request.EvidenceReferences,
	})
	if err != nil {
		t.Fatal(err)
	}
	return service, record
}

func approvalFixture(t testing.TB) (*Service, *policyResolverStub, *MemoryStore, identity.Principal) {
	t.Helper()
	policies := &policyResolverStub{decisions: make(map[string]trust.PolicyDecision)}
	store := NewMemoryStore()
	service, err := NewService(
		subjectResolverStub{subject: validSubjectFixture(), found: true}, policies, store,
		Profile{
			Issuer: "forge.approval.developer-trial", MaximumDecisionLifetime: 30 * time.Minute,
			SupportedConditions:     []string{"require-independent-verification"},
			EnforcedObligationTypes: []string{"human-interaction"},
		},
		func() time.Time { return approvalNow },
	)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := identity.NewVerifier(identityProviderStub{}, "forge-cli", func() time.Time { return approvalNow })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque", CorrelationID: "correlation-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	return service, policies, store, principal
}

func validSubjectFixture() PlanSubject {
	return PlanSubject{
		RunID: "run-1", PlanID: "plan-1", PlanDigest: approvalSubjectDigest,
		Scope:             trust.ApprovalScope{TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", RepositoryID: "repository-1", DeliveryRunID: "run-1"},
		RequestingActorID: "human-1", Purpose: "forge.m1.prepare-change", Classification: "Internal",
		RiskClass: "low", ContextPackageID: "context-1", ContextManifest: "sha256:manifest",
		EvidenceReferences: []string{"intent:version:1", "plan:" + approvalSubjectDigest},
		PolicyDecisionIDs:  []string{"context-decision"}, ExpiresAt: approvalNow.Add(30 * time.Minute),
		CorrelationID: "correlation-1",
	}
}

func planReviewCommand(key string) RequestPlanReview {
	return RequestPlanReview{
		ContractVersion: RequestPlanReviewContractVersion, CorrelationID: "correlation-1",
		IdempotencyKey: key, TenantID: "tenant-1", RunID: "run-1",
		RequestingActorID: "human-1", IntendedApproverID: "human-1",
		RequestedExpiry: approvalNow.Add(30 * time.Minute),
	}
}

func eligibilityDecision(request trust.ApprovalRequest, outcome string) trust.PolicyDecision {
	expires := approvalNow.Add(25 * time.Minute)
	return trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, ID: map[string]string{"permit": "eligibility-1", "indeterminate": "indeterminate"}[outcome],
		DecisionPoint: "approval", SubjectID: "human-1", Action: "forge.approval.decide",
		Resource: request.RequestDigest, EvaluatedAt: approvalNow, ExpiresAt: &expires,
		PolicyVersions: []string{"policy-v1"}, Outcome: outcome,
		Obligations: []map[string]any{{"type": "human-interaction"}},
		ReasonCodes: []string{"test"}, EvidenceCorrelationID: "correlation-1",
	}
}

func assertApprovalCode(t *testing.T, err error, code failure.Code) {
	t.Helper()
	var normalized *failure.Error
	if err == nil || !errors.As(err, &normalized) || normalized.Code != code {
		t.Fatalf("error = %v, want %s", err, code)
	}
}
