package delivery

import (
	"context"
	"errors"
	"slices"
	"sort"
	"testing"
	"time"

	approvalcap "github.com/3corners/forge/internal/approval"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
	"github.com/3corners/forge/internal/identity"
)

var lifecycleTime = time.Date(2026, time.July, 23, 11, 0, 0, 0, time.UTC)

var (
	_ approvalcap.SubjectResolver = (*Service)(nil)
	_ ApprovalRequestResolver     = (*approvalcap.Service)(nil)
	_ ApprovalResolver            = (*approvalcap.Service)(nil)
)

type approvalResolverStub struct {
	record trust.ApprovalRecordV2
	found  bool
	err    error
}

func (resolver approvalResolverStub) ResolveApproval(context.Context, string, string) (trust.ApprovalRecordV2, bool, error) {
	return resolver.record, resolver.found, resolver.err
}

type approvalRequestResolverStub struct {
	request trust.ApprovalRequest
	found   bool
	err     error
}

func (resolver approvalRequestResolverStub) ResolveApprovalRequest(context.Context, string, string) (trust.ApprovalRequest, bool, error) {
	return resolver.request, resolver.found, resolver.err
}

type approvalPolicyResolverStub struct {
	decision trust.PolicyDecision
}

func (resolver *approvalPolicyResolverStub) ResolvePolicyDecision(_ context.Context, id string) (trust.PolicyDecision, bool, error) {
	return resolver.decision, resolver.decision.ID == id, nil
}

type lifecycleIdentityProvider struct{}

func (lifecycleIdentityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return identity.Assertion{
		SubjectID: "human-1", Kind: identity.HumanPrincipal, Issuer: "forge.synthetic",
		Audience: "forge-cli", Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: lifecycleTime.Add(-time.Minute), ExpiresAt: lifecycleTime.Add(time.Hour), Synthetic: true,
	}, nil
}

func TestContextToApprovedPlanLifecycleIsAtomicAndDigestBound(t *testing.T) {
	service, store, run := lifecycleFixture(t)
	ctx := context.Background()

	contextAuthorizing, replayed, err := service.BeginContextAuthorization(ctx, transitionRequest(run, "transition-context-authorizing"))
	if err != nil || replayed || contextAuthorizing.State() != ContextAuthorizingState || contextAuthorizing.Version() != 2 {
		t.Fatalf("context authorization = state %q version %d replay=%t err=%v", contextAuthorizing.State(), contextAuthorizing.Version(), replayed, err)
	}
	replayedRun, replayed, err := service.BeginContextAuthorization(ctx, transitionRequest(run, "transition-context-authorizing"))
	if err != nil || !replayed || replayedRun.Version() != 2 {
		t.Fatalf("context authorization replay = %#v replay=%t err=%v", replayedRun, replayed, err)
	}

	envelope, contextPackage := lifecycleContext(t, contextAuthorizing)
	contextReady, _, err := service.RecordContext(ctx, RecordContextRequest{
		TransitionRequest: transitionRequest(contextAuthorizing, "transition-context-ready"),
		Package:           contextPackage, Envelope: envelope,
	})
	if err != nil || contextReady.State() != ContextReadyState || contextReady.Context() == nil {
		t.Fatalf("record context = state %q err=%v", contextReady.State(), err)
	}

	planned, _, err := service.RecordPlan(ctx, planProfiles(t), RecordPlanRequest{
		TransitionRequest: transitionRequest(contextReady, "transition-planned"),
		ProfileVersion:    validPlanProfile().Version,
	})
	if err != nil || planned.State() != PlannedState || planned.Plan() == nil {
		t.Fatalf("record plan = state %q err=%v", planned.State(), err)
	}
	if err := planned.Plan().Validate("correlation-lifecycle"); err != nil {
		t.Fatal(err)
	}

	approvalRequest := planApprovalRequest(t, planned)
	awaiting, _, err := service.AwaitPlanApproval(ctx, approvalRequestResolverStub{request: approvalRequest, found: true}, AwaitPlanApprovalRequest{
		TransitionRequest: transitionRequest(planned, "transition-awaiting-approval"),
		ApprovalRequestID: approvalRequest.ID,
	})
	if err != nil || awaiting.State() != AwaitingPlanApprovalState || awaiting.Approval() == nil {
		t.Fatalf("await approval = state %q err=%v", awaiting.State(), err)
	}
	record := planApprovalRecord(t, approvalRequest)
	approved, _, err := service.RecordPlanApproval(ctx, approvalResolverStub{record: record, found: true}, RecordPlanApprovalRequest{
		TransitionRequest: transitionRequest(awaiting, "transition-plan-approved"),
		ApprovalRecordID:  record.ID,
	})
	if err != nil || approved.State() != PlanApprovedState || approved.Approval().RecordID != record.ID {
		t.Fatalf("record approval = state %q approval=%#v err=%v", approved.State(), approved.Approval(), err)
	}
	awaitingPublication := completeM1BDelivery(t, service, approved)
	if awaitingPublication.State() != AwaitingPublicationApprovalState ||
		awaitingPublication.Verification() == nil || awaitingPublication.Verification().Outcome != "passed" {
		t.Fatalf("M1-B completion = state %q verification=%#v", awaitingPublication.State(), awaitingPublication.Verification())
	}

	stored, found, err := store.Load(ctx, run.ID())
	if err != nil || !found || len(stored.Outbox) != 9 {
		t.Fatalf("stored outbox count=%d found=%t err=%v", len(stored.Outbox), found, err)
	}
	for index, event := range stored.Outbox {
		if event.RunVersion != uint64(index+2) || event.RunID != run.ID() || event.CorrelationID != run.CorrelationID() {
			t.Fatalf("invalid outbox event %d: %#v", index, event)
		}
	}
}

func completeM1BDelivery(t testing.TB, service *Service, approved DeliveryRun) DeliveryRun {
	t.Helper()
	ctx := context.Background()
	authorizing, _, err := service.BeginExecutionAuthorization(ctx, transitionRequest(approved, "transition-execution-authorizing"))
	if err != nil {
		t.Fatal(err)
	}
	plan, contextRecord, approval := authorizing.Plan(), authorizing.Context(), authorizing.Approval()
	envelope, err := trust.BindExecutionEnvelopeIntegrity(trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: "forge.policy",
		IssuedAt: lifecycleTime.Add(-time.Minute), ExpiresAt: lifecycleTime.Add(time.Hour),
		Scope: trust.ExecutionScope{
			TenantID: authorizing.TenantID(), WorkspaceID: authorizing.WorkspaceID(),
			ProjectID: authorizing.ProjectID(), DeliveryRunID: authorizing.ID(),
		},
		Principal: trust.EnvelopePrincipal{
			SubjectID: authorizing.ActorID(), AgentID: "worker-1",
			DelegationChain: []string{"worker-delegation-1"},
		},
		Purpose: M1ChangePurpose, RiskClass: plan.RiskClass,
		Actions: []string{workercontract.PrepareChangeAction}, Classifications: []string{authorizing.Classification()},
		Destinations: []string{}, ModelIDs: []string{}, ToolIDs: []string{}, SecretReferenceIDs: []string{},
		PolicyDecisionIDs:     []string{"decision-delegation", "decision-execution"},
		ApprovalRecordIDs:     []string{approval.RecordID},
		Limits:                trust.ExecutionLimits{MaxDurationSeconds: 60, MaxConcurrency: 1},
		EvidenceCorrelationID: authorizing.CorrelationID(),
	}, authorizing.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	executing, _, err := service.RecordExecutionAuthority(ctx, RecordExecutionAuthorityRequest{
		TransitionRequest: transitionRequest(authorizing, "transition-executing"),
		Envelope:          envelope, DelegationID: "worker-delegation-1", WorkerID: "worker-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	replacements := make([]workercontract.FileReplacement, 0, len(plan.ChangePaths))
	for _, changedPath := range plan.ChangePaths {
		replacements = append(replacements, workercontract.FileReplacement{Path: changedPath, Content: []byte("bounded candidate\n")})
	}
	workerRequest, err := workercontract.BindRequest(workercontract.Request{
		SchemaVersion: workercontract.RequestVersion, EnvelopeID: envelope.ID,
		DelegationID: "worker-delegation-1", PlanID: plan.ID, PlanDigest: plan.Digest,
		PlanApprovalRecordID: approval.RecordID, ContextPackageID: contextRecord.PackageID,
		ContextManifestDigest: contextRecord.ManifestDigest, SourceRevision: contextRecord.SourceRevisions[0],
		Scope: workercontract.Scope{
			TenantID: executing.TenantID(), WorkspaceID: executing.WorkspaceID(), ProjectID: executing.ProjectID(),
			RepositoryID: executing.RepositoryID(), DeliveryRunID: executing.ID(),
		},
		Action: workercontract.PrepareChangeAction, TransformationID: workercontract.TransformationV1,
		TransformationVersion: "1.0.0", AllowedReadPaths: slices.Clone(plan.ChangePaths),
		AllowedWritePaths: slices.Clone(plan.ChangePaths), DependencyDigests: []string{plan.Digest, contextRecord.ManifestDigest},
		Replacements: replacements,
		Limits: workercontract.Limits{
			MaxDurationSeconds: 60, MaxProcesses: 4, MaxOutputBytes: 1 << 20,
			MaxChangedFiles: int64(len(replacements)), MaxAttempts: plan.Budgets.MaxAttempts,
		},
		Attempt: 1, ExpiresAt: envelope.ExpiresAt, EvidenceCorrelationID: executing.CorrelationID(),
	})
	if err != nil {
		t.Fatal(err)
	}
	files := make([]workercontract.CandidateFile, 0, len(replacements))
	for _, replacement := range replacements {
		files = append(files, workercontract.CandidateFile{Path: replacement.Path, Content: replacement.Content})
	}
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: workerRequest.SourceRevision,
		PlanDigest: plan.Digest, Files: files,
	})
	if err != nil {
		t.Fatal(err)
	}
	changedPaths := make([]string, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		changedPaths = append(changedPaths, file.Path)
	}
	workerResult := workercontract.Result{
		SchemaVersion: workercontract.ResultVersion, RequestID: workerRequest.ID, RequestDigest: workerRequest.Digest,
		CandidateDigest: bundle.Digest, TreeDigest: workercontract.CandidateTreeDigest(bundle),
		ChangedPaths: changedPaths, BundleDigest: bundle.Digest,
		ResourceUse:    workercontract.ResourceUse{DurationMillis: 1, OutputBytes: 128, ChangedFiles: int64(len(changedPaths))},
		CleanupOutcome: "completed", EvidenceReferences: []string{}, CompletedAt: lifecycleTime,
	}
	verifying, _, err := service.RecordCandidate(ctx, RecordCandidateRequest{
		TransitionRequest: transitionRequest(executing, "transition-verifying"),
		WorkerRequest:     workerRequest, Bundle: bundle, Result: workerResult,
		ArtefactReferences: []string{"artifact:" + bundle.Digest},
	})
	if err != nil {
		t.Fatal(err)
	}
	assurancePlan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion: verificationcontract.AssurancePlanVersion, CandidateDigest: bundle.Digest,
		PlanDigest: plan.Digest, RiskClass: plan.RiskClass, Classification: verifying.Classification(),
		AffectedCapabilities: slices.Clone(plan.AffectedCapabilities), ChangedPaths: changedPaths,
		Controls:    []verificationcontract.Control{{ID: "forge.candidate-integrity", Version: "1", Consequential: true}},
		Environment: "developer-trial", ExpiresAt: lifecycleTime.Add(time.Hour),
		EvidenceCorrelationID: verifying.CorrelationID(),
	})
	if err != nil {
		t.Fatal(err)
	}
	verificationResult, err := verificationcontract.BindResult(verificationcontract.Result{
		SchemaVersion: verificationcontract.ResultVersion, AssurancePlanID: assurancePlan.ID,
		AssurancePlanDigest: assurancePlan.Digest, CandidateDigest: bundle.Digest,
		Environment: "developer-trial", Outcome: "passed", Independent: true, CompletedAt: lifecycleTime,
		EvidenceCorrelationID: verifying.CorrelationID(),
		Findings: []verificationcontract.Finding{{
			ControlID: "forge.candidate-integrity", Outcome: "passed", Code: "passed",
			EvidenceDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Consequential:  true,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	completed, _, err := service.RecordVerification(ctx, RecordVerificationRequest{
		TransitionRequest: transitionRequest(verifying, "transition-awaiting-publication"),
		AssurancePlan:     assurancePlan, Result: verificationResult,
		EvidenceReferences: []string{"verification:" + verificationResult.Digest},
	})
	if err != nil {
		t.Fatal(err)
	}
	return completed
}

func TestDeliveryAndApprovalCapabilitiesComposeThroughOwnedResolvers(t *testing.T) {
	service, _, run := lifecycleFixture(t)
	ctx := context.Background()
	authorizing, _, err := service.BeginContextAuthorization(ctx, transitionRequest(run, "compose-authorizing"))
	if err != nil {
		t.Fatal(err)
	}
	envelope, contextPackage := lifecycleContext(t, authorizing)
	ready, _, err := service.RecordContext(ctx, RecordContextRequest{
		TransitionRequest: transitionRequest(authorizing, "compose-context"),
		Package:           contextPackage, Envelope: envelope,
	})
	if err != nil {
		t.Fatal(err)
	}
	planned, _, err := service.RecordPlan(ctx, planProfiles(t), RecordPlanRequest{
		TransitionRequest: transitionRequest(ready, "compose-plan"),
		ProfileVersion:    validPlanProfile().Version,
	})
	if err != nil {
		t.Fatal(err)
	}
	policies := &approvalPolicyResolverStub{}
	approvals, err := approvalcap.NewService(
		service, policies, approvalcap.NewMemoryStore(),
		approvalcap.Profile{
			Issuer: "forge.approval.developer-trial", MaximumDecisionLifetime: 20 * time.Minute,
			SupportedConditions:     []string{"require-independent-verification"},
			EnforcedObligationTypes: []string{"human-interaction"},
		},
		func() time.Time { return lifecycleTime },
	)
	if err != nil {
		t.Fatal(err)
	}
	approvalRequest, _, err := approvals.RequestPlanReview(ctx, approvalcap.RequestPlanReview{
		ContractVersion: approvalcap.RequestPlanReviewContractVersion, CorrelationID: run.CorrelationID(),
		IdempotencyKey: "compose-request", TenantID: run.TenantID(), RunID: run.ID(),
		RequestingActorID: run.ActorID(), IntendedApproverID: run.ActorID(),
		RequestedExpiry: lifecycleTime.Add(20 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	awaiting, _, err := service.AwaitPlanApproval(ctx, approvals, AwaitPlanApprovalRequest{
		TransitionRequest: transitionRequest(planned, "compose-awaiting"),
		ApprovalRequestID: approvalRequest.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	expires := lifecycleTime.Add(10 * time.Minute)
	policies.decision = trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, ID: "eligibility-compose",
		DecisionPoint: "approval", SubjectID: run.ActorID(), Action: "forge.approval.decide",
		Resource: approvalRequest.RequestDigest, EvaluatedAt: lifecycleTime, ExpiresAt: &expires,
		PolicyVersions: []string{"policy.compose.v1"}, Outcome: "permit",
		Obligations: []map[string]any{{"type": "human-interaction"}},
		ReasonCodes: []string{"policy-permit"}, EvidenceCorrelationID: run.CorrelationID(),
	}
	verifier, err := identity.NewVerifier(lifecycleIdentityProvider{}, "forge-cli", func() time.Time { return lifecycleTime })
	if err != nil {
		t.Fatal(err)
	}
	principal, err := verifier.Verify(ctx, identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque",
		CorrelationID: run.CorrelationID(),
	})
	if err != nil {
		t.Fatal(err)
	}
	record, _, err := approvals.DecideApproval(ctx, approvalcap.DecideApproval{
		ContractVersion: approvalcap.DecideApprovalContractVersion, CorrelationID: run.CorrelationID(),
		IdempotencyKey: "compose-decision", TenantID: run.TenantID(), RequestID: approvalRequest.ID,
		Approver: principal, PolicyDecisionID: policies.decision.ID, Decision: "approved",
		Conditions:         []string{"require-independent-verification"},
		EvidenceReferences: approvalRequest.EvidenceReferences,
	})
	if err != nil {
		t.Fatal(err)
	}
	approved, _, err := service.RecordPlanApproval(ctx, approvals, RecordPlanApprovalRequest{
		TransitionRequest: transitionRequest(awaiting, "compose-approved"), ApprovalRecordID: record.ID,
	})
	if err != nil || approved.State() != PlanApprovedState {
		t.Fatalf("composed state=%q err=%v", approved.State(), err)
	}
}

func TestPlanAndApprovalNegativeBoundariesFailClosed(t *testing.T) {
	service, _, run := lifecycleFixture(t)
	ctx := context.Background()
	authorizing, _, _ := service.BeginContextAuthorization(ctx, transitionRequest(run, "begin"))
	envelope, contextPackage := lifecycleContext(t, authorizing)
	changed := contextPackage
	changed.Items = slices.Clone(contextPackage.Items)
	changed.Items[0].SourceRevision = "changed"
	_, _, err := service.RecordContext(ctx, RecordContextRequest{
		TransitionRequest: transitionRequest(authorizing, "context"),
		Package:           changed, Envelope: envelope,
	})
	assertLifecycleCode(t, err, failure.VerificationFailed)

	ready, _, err := service.RecordContext(ctx, RecordContextRequest{
		TransitionRequest: transitionRequest(authorizing, "context"),
		Package:           contextPackage, Envelope: envelope,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = service.RecordPlan(ctx, planProfiles(t), RecordPlanRequest{
		TransitionRequest: transitionRequest(ready, "plan-invalid"), ProfileVersion: "caller-controlled-profile",
	})
	assertLifecycleCode(t, err, failure.Forbidden)

	planned, _, err := service.RecordPlan(ctx, planProfiles(t), RecordPlanRequest{
		TransitionRequest: transitionRequest(ready, "plan"), ProfileVersion: validPlanProfile().Version,
	})
	if err != nil {
		t.Fatal(err)
	}
	request := planApprovalRequest(t, planned)
	request.Scope.ProjectID = "other-project"
	_, _, err = service.AwaitPlanApproval(ctx, approvalRequestResolverStub{request: request, found: true}, AwaitPlanApprovalRequest{
		TransitionRequest: transitionRequest(planned, "await"), ApprovalRequestID: request.ID,
	})
	assertLifecycleCode(t, err, failure.VerificationFailed)
}

func TestDeliveryPlanDigestChangesWithNormativeInput(t *testing.T) {
	_, _, run := lifecycleFixture(t)
	contextRecord := &ContextRecord{
		PackageID: "package-1", ManifestDigest: "sha256:manifest", EnvelopeID: "envelope-1",
		PolicyDecisionIDs: []string{"decision-1"}, SourceRevisions: []string{"source@revision"},
		Omissions: []string{}, Conflicts: []string{}, Classification: "Internal",
		ExpiresAt: lifecycleTime.Add(time.Hour),
	}
	run.state = ContextReadyState
	run.context = contextRecord
	left, err := buildPlan(run, validPlanProfile(), lifecycleTime, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	profile := validPlanProfile()
	profile.Budgets.MaxAttempts++
	right, err := buildPlan(run, profile, lifecycleTime, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	if left.Digest == right.Digest {
		t.Fatal("normative budget mutation did not change the plan digest")
	}
}

func BenchmarkDeterministicPlan(b *testing.B) {
	_, _, run := lifecycleFixture(b)
	run.state = ContextReadyState
	run.context = &ContextRecord{
		PackageID: "package-1", ManifestDigest: "sha256:manifest", EnvelopeID: "envelope-1",
		PolicyDecisionIDs: []string{"decision-1"}, SourceRevisions: []string{"source@revision"},
		Omissions: []string{}, Conflicts: []string{}, Classification: "Internal",
		ExpiresAt: lifecycleTime.Add(time.Hour),
	}
	profile := validPlanProfile()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := buildPlan(run, profile, lifecycleTime, run.CorrelationID()); err != nil {
			b.Fatal(err)
		}
	}
}

func TestDeterministicPlanLatencyDistribution(t *testing.T) {
	_, _, run := lifecycleFixture(t)
	run.state = ContextReadyState
	run.context = &ContextRecord{
		PackageID: "package-1", ManifestDigest: "sha256:manifest", EnvelopeID: "envelope-1",
		PolicyDecisionIDs: []string{"decision-1"}, SourceRevisions: []string{"source@revision"},
		Omissions: []string{}, Conflicts: []string{}, Classification: "Internal",
		ExpiresAt: lifecycleTime.Add(time.Hour),
	}
	profile := validPlanProfile()
	const samples = 2000
	durations := make([]time.Duration, samples)
	for index := range durations {
		started := time.Now()
		if _, err := buildPlan(run, profile, lifecycleTime, run.CorrelationID()); err != nil {
			t.Fatal(err)
		}
		durations[index] = time.Since(started)
	}
	sort.Slice(durations, func(left, right int) bool { return durations[left] < durations[right] })
	p50, p95, p99 := durations[samples*50/100], durations[samples*95/100], durations[samples*99/100]
	t.Logf("deterministic plan latency p50=%s p95=%s p99=%s samples=%d", p50, p95, p99, samples)
	if p99 >= 500*time.Millisecond {
		t.Fatalf("deterministic plan p99=%s exceeds the M1 control-operation hypothesis", p99)
	}
}

func lifecycleFixture(t testing.TB) (*Service, *MemoryStore, DeliveryRun) {
	t.Helper()
	store, err := NewMemoryStore(func() string { return "unused" })
	if err != nil {
		t.Fatal(err)
	}
	record := StoredRecord{
		CaptureRecord: CaptureRecord{
			ScopeKey:           "tenant-1\x00workspace-1\x00project-1\x00repository-1\x00human-1",
			RequestFingerprint: "sha256:capture", IdempotencyKey: "capture", CorrelationID: "correlation-lifecycle",
			TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1",
			RepositoryID: "repository-1", ActorID: "human-1", DelegationChain: []string{},
			Purpose: M1ChangePurpose, Classification: "Internal", Objective: "Prepare one bounded change",
			Constraints: []string{"No production access"}, AcceptanceCriteria: []string{"Tests pass"},
			CreatedAt: lifecycleTime.Add(-time.Minute),
		},
		ContractVersion: DeliveryRunContractVersion, ID: "run-1", State: CapturedState, Version: 1,
	}
	store.records[record.ScopeKey] = cloneStoredRecord(record)
	store.byID[record.ID] = record.ScopeKey
	service, err := NewService(store, func() time.Time { return lifecycleTime })
	if err != nil {
		t.Fatal(err)
	}
	return service, store, runFromStored(record)
}

func lifecycleContext(t testing.TB, run DeliveryRun) (trust.ExecutionEnvelope, trust.ContextPackage) {
	t.Helper()
	envelope, err := trust.BindExecutionEnvelopeIntegrity(trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: "forge.policy",
		IssuedAt: lifecycleTime.Add(-time.Minute), ExpiresAt: lifecycleTime.Add(time.Hour),
		Scope:     trust.ExecutionScope{TenantID: run.TenantID(), WorkspaceID: run.WorkspaceID(), ProjectID: run.ProjectID(), DeliveryRunID: run.ID()},
		Principal: trust.EnvelopePrincipal{SubjectID: run.ActorID(), AgentID: "agent-1", DelegationChain: []string{"delegation-1"}},
		Purpose:   run.Purpose(), RiskClass: "low", Actions: []string{"forge.context.assemble"},
		Classifications: []string{run.Classification()}, Destinations: []string{}, ModelIDs: []string{},
		ToolIDs: []string{}, SecretReferenceIDs: []string{}, PolicyDecisionIDs: []string{"decision-context"},
		ApprovalRecordIDs: []string{}, Limits: trust.ExecutionLimits{MaxDurationSeconds: 600, MaxConcurrency: 1},
		EvidenceCorrelationID: run.CorrelationID(),
	}, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	expires := lifecycleTime.Add(30 * time.Minute)
	value, err := trust.BindContextPackageManifest(trust.ContextPackage{
		SchemaVersion: trust.ContextPackageSchemaVersion, ExecutionEnvelopeID: envelope.ID,
		Objective: run.Objective(), AssembledAt: lifecycleTime, ExpiresAt: &expires,
		PolicyDecisionIDs: []string{"decision-context"},
		Items: []trust.ContextItem{{
			SourceID: "architecture", SourceRevision: "revision-1", Locator: "docs/architecture",
			Classification: run.Classification(), Authority: "accepted-knowledge", Purpose: run.Purpose(),
			ContentDigest: "sha256:content",
		}}, Omissions: []string{}, Conflicts: []string{},
	}, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	return envelope, value
}

func validPlanProfile() PlanProfile {
	return PlanProfile{
		Version: "forge.m1-plan-profile.v1",
		Steps: []PlanStep{
			{ID: "implement", Dependencies: []string{}, OutputKinds: []string{"source-change"}},
			{ID: "verify", Dependencies: []string{"implement"}, OutputKinds: []string{"verification-result"}},
		},
		AllowedCapabilities: []string{"forge.capability.approval", "forge.capability.delivery"},
		AllowedPathPrefixes: []string{"internal", "schemas"},
		ChangePaths:         []string{"internal/delivery", "schemas/forge/trust"},
		Assumptions:         []string{"accepted architecture remains unchanged"},
		Risks:               []string{"bounded trial state is in-process only"},
		Limitations:         []string{"no worker execution"},
		StopConditions:      []string{"architecture decision required"},
		Budgets:             PlanBudgets{MaxDurationSeconds: 1800, MaxConcurrency: 1, MaxAttempts: 3},
		VerificationIntent:  []string{"go test ./..."},
		GuardrailReferences: []string{"ADR-0015", "ADR-0016"},
		Rollback:            "Discard the candidate branch.", ManualRecovery: "Return to ContextReady.",
		RiskClass: "low", MaximumLifetime: 30 * time.Minute,
	}
}

func planProfiles(t testing.TB) *StaticPlanProfiles {
	t.Helper()
	profiles, err := NewStaticPlanProfiles([]PlanProfile{validPlanProfile()})
	if err != nil {
		t.Fatal(err)
	}
	return profiles
}

func transitionRequest(run DeliveryRun, key string) TransitionRequest {
	return TransitionRequest{
		CorrelationID: run.CorrelationID(), IdempotencyKey: key,
		Run: run, ExpectedVersion: run.Version(),
	}
}

func planApprovalRequest(t testing.TB, run DeliveryRun) trust.ApprovalRequest {
	t.Helper()
	plan := run.Plan()
	value, err := trust.BindApprovalRequestDigest(trust.ApprovalRequest{
		SchemaVersion: trust.ApprovalRequestSchemaVersion, Issuer: "forge.approval",
		RequestingActorID: run.ActorID(), IntendedApproverID: run.ActorID(), Scope: approvalScope(run),
		Environment: "developer-trial", Purpose: run.Purpose(), Classification: run.Classification(),
		RiskClass: plan.RiskClass, Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan",
		SubjectID: plan.ID, SubjectDigest: plan.Digest, EvidenceReferences: planEvidence(plan),
		PermittedConditions: []string{"require-independent-verification"},
		PolicyDecisionIDs:   []string{"decision-context"}, CreatedAt: lifecycleTime,
		ExpiresAt: lifecycleTime.Add(20 * time.Minute), MaxDecisionLifetime: 1200,
		Lifecycle: "open", EvidenceCorrelationID: run.CorrelationID(),
	}, run.CorrelationID())
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func planApprovalRecord(t testing.TB, request trust.ApprovalRequest) trust.ApprovalRecordV2 {
	t.Helper()
	value, err := trust.BindApprovalRecordV2Integrity(trust.ApprovalRecordV2{
		SchemaVersion: trust.ApprovalRecordV2SchemaVersion, Issuer: request.Issuer,
		RequestID: request.ID, RequestDigest: request.RequestDigest, ApproverID: request.IntendedApproverID,
		Scope: request.Scope, Environment: request.Environment, Purpose: request.Purpose,
		Classification: request.Classification, RiskClass: request.RiskClass, Action: request.Action,
		SubjectType: request.SubjectType, SubjectID: request.SubjectID, SubjectDigest: request.SubjectDigest,
		Decision: "approved", DecidedAt: lifecycleTime, ExpiresAt: lifecycleTime.Add(10 * time.Minute),
		Conditions: []string{"require-independent-verification"}, PolicyDecisionIDs: []string{"decision-approval"},
		EvidenceReferences: request.EvidenceReferences, EvidenceCorrelationID: request.EvidenceCorrelationID,
	}, request.EvidenceCorrelationID)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func assertLifecycleCode(t *testing.T, err error, code failure.Code) {
	t.Helper()
	var normalized *failure.Error
	if err == nil || !errors.As(err, &normalized) || normalized.Code != code {
		t.Fatalf("error = %v, want %s", err, code)
	}
}
