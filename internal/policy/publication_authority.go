package policy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"sort"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const (
	EvaluatePublicationVersion      = "forge.policy-evaluate-publication.v1"
	IssuePublicationEnvelopeVersion = "forge.issue-publication-envelope.v1"
)

type PublicationSubjectResolver interface {
	ResolvePublicationSubject(context.Context, string, string) (trust.PublicationApprovalSubject, bool, error)
}

type EvaluatePublicationRequest struct {
	SchemaVersion     string
	CorrelationID     string
	Human             identity.Principal
	Workload          identity.Principal
	Scope             tenant.Scope
	RepositoryBinding tenant.RepositoryBinding
	Run               delivery.DeliveryRun
	ApprovalRecordID  string
	Provider          string
	BaseBranch        string
	BaseRevision      string
	Environment       Environment
	Risk              Risk
}

func (evaluator *Evaluator) EvaluatePublication(ctx context.Context, subjects PublicationSubjectResolver, approvals ApprovalResolver, request EvaluatePublicationRequest) (trust.PolicyDecision, trust.PublicationApprovalSubject, trust.ApprovalRecordV2, error) {
	at, err := evaluator.workerEvaluationTime(ctx, request.SchemaVersion, EvaluatePublicationVersion, request.CorrelationID)
	if err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, err
	}
	if subjects == nil || approvals == nil || request.ApprovalRecordID == "" ||
		request.Provider != "contract-double" || request.BaseBranch == "" || request.BaseRevision == "" ||
		request.Environment != DeveloperTrialEnvironment || !validRisk(request.Risk) {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, publicationForbidden(request.CorrelationID)
	}
	if err := request.Human.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, err
	}
	if err := request.Workload.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, err
	}
	if request.Human.Kind() != identity.HumanPrincipal || request.Workload.Kind() != identity.WorkloadPrincipal ||
		request.Run.State() != delivery.AwaitingPublicationApprovalState ||
		request.Run.ActorID() != request.Human.SubjectID() ||
		request.Run.TenantID() != request.Scope.TenantID() ||
		request.Run.WorkspaceID() != request.Scope.WorkspaceID() ||
		request.Run.ProjectID() != request.Scope.ProjectID() {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, publicationForbidden(request.CorrelationID)
	}
	if err := request.Scope.ValidateFor(request.Human, at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, err
	}
	if err := request.RepositoryBinding.ValidateFor(request.Scope, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, err
	}
	subject, found, err := subjects.ResolvePublicationSubject(ctx, request.Run.TenantID(), request.Run.ID())
	if err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, failure.New(failure.DependencyUnavailable, "publication subject resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || subject.RunID != request.Run.ID() || subject.Scope.RepositoryID != request.RepositoryBinding.RepositoryID() ||
		subject.CandidateDigest == "" || subject.VerificationResultDigest == "" {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, publicationForbidden(request.CorrelationID)
	}
	record, found, err := approvals.ResolveApproval(ctx, request.Run.TenantID(), request.ApprovalRecordID)
	if err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, failure.New(failure.DependencyUnavailable, "publication approval resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found || record.ID != request.ApprovalRecordID {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, publicationForbidden(request.CorrelationID)
	}
	expected := trust.ApprovalExpectedFacts{
		RequestID: record.RequestID, RequestDigest: record.RequestDigest, ApproverID: request.Human.SubjectID(),
		Scope: subject.Scope, Environment: string(DeveloperTrialEnvironment), Purpose: subject.Purpose,
		Classification: subject.Classification, RiskClass: subject.RiskClass,
		Action: string(ActionPublishPullRequest), SubjectType: "verified-candidate",
		SubjectID: subject.CandidateID, SubjectDigest: subject.CandidateDigest,
		EvidenceReferences:  []string{"candidate:" + subject.CandidateDigest, "verification:" + subject.VerificationResultDigest},
		PermittedConditions: []string{},
	}
	if err := record.Authorize(expected, at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, err
	}
	matching := evaluator.workerRules(request.Workload.SubjectID(), ActionPublishPullRequest, subject.CandidateDigest, request.Run, request.Scope, request.Environment, request.Risk)
	outcome, obligations, reasons := evaluator.resolve(matching)
	expiresAt := earliestPolicyExpiry(at.Add(evaluator.decisionTTL), request.Human.ExpiresAt(), request.Workload.ExpiresAt(), subject.ExpiresAt, record.ExpiresAt)
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, DecisionPoint: "publication",
		SubjectID: request.Workload.SubjectID(), Action: string(ActionPublishPullRequest),
		Resource: subject.CandidateDigest, EvaluatedAt: at, ExpiresAt: &expiresAt,
		PolicyVersions: []string{evaluator.set.Version}, Outcome: outcome, Obligations: obligations,
		ReasonCodes: reasons, EvidenceCorrelationID: request.CorrelationID,
	}
	decision.ID = workerDecisionID(decision, request.Run.ID(), record.Integrity.Digest+"\x00"+request.BaseRevision, matching)
	if err := decision.Validate(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, trust.PublicationApprovalSubject{}, trust.ApprovalRecordV2{}, failure.New(failure.Indeterminate, "publication policy decision could not be verified", request.CorrelationID, failure.RetryNone)
	}
	evaluator.recordDecision(decision)
	return decision, subject, record, nil
}

type PublicationEnvelopeProfile struct {
	Issuer                  string
	MaximumLifetime         time.Duration
	EnforcedObligationTypes []string
}

type PublicationEnvelopeIssuer struct {
	evaluator *Evaluator
	subjects  PublicationSubjectResolver
	approvals ApprovalResolver
	store     EnvelopeStore
	profile   PublicationEnvelopeProfile
	now       func() time.Time
}

func NewPublicationEnvelopeIssuer(evaluator *Evaluator, subjects PublicationSubjectResolver, approvals ApprovalResolver, store EnvelopeStore, profile PublicationEnvelopeProfile, now func() time.Time) (*PublicationEnvelopeIssuer, error) {
	if evaluator == nil || subjects == nil || approvals == nil || store == nil || now == nil ||
		profile.Issuer == "" || profile.MaximumLifetime <= 0 ||
		!slices.Equal(canonicalPublicationStrings(profile.EnforcedObligationTypes), []string{"evidence-ready", "human-publication-approval"}) {
		return nil, errors.New("publication envelope issuer configuration is invalid")
	}
	profile.EnforcedObligationTypes = canonicalPublicationStrings(profile.EnforcedObligationTypes)
	return &PublicationEnvelopeIssuer{evaluator: evaluator, subjects: subjects, approvals: approvals, store: store, profile: profile, now: now}, nil
}

type IssuePublicationEnvelopeRequest struct {
	SchemaVersion     string
	CorrelationID     string
	IdempotencyKey    string
	Human             identity.Principal
	Workload          identity.Principal
	Scope             tenant.Scope
	RepositoryBinding tenant.RepositoryBinding
	Run               delivery.DeliveryRun
	ApprovalRecordID  string
	Provider          string
	BaseBranch        string
	BaseRevision      string
	Environment       Environment
	Risk              Risk
}

func (issuer *PublicationEnvelopeIssuer) Issue(ctx context.Context, request IssuePublicationEnvelopeRequest) (trust.ExecutionEnvelope, bool, error) {
	if request.SchemaVersion != IssuePublicationEnvelopeVersion || request.IdempotencyKey == "" {
		return trust.ExecutionEnvelope{}, false, failure.New(failure.InvalidRequest, "publication envelope request is invalid", request.CorrelationID, failure.RetryNone)
	}
	decision, subject, approval, err := issuer.evaluator.EvaluatePublication(ctx, issuer.subjects, issuer.approvals, EvaluatePublicationRequest{
		SchemaVersion: EvaluatePublicationVersion, CorrelationID: request.CorrelationID,
		Human: request.Human, Workload: request.Workload, Scope: request.Scope,
		RepositoryBinding: request.RepositoryBinding, Run: request.Run,
		ApprovalRecordID: request.ApprovalRecordID, Provider: request.Provider,
		BaseBranch: request.BaseBranch, BaseRevision: request.BaseRevision,
		Environment: request.Environment, Risk: request.Risk,
	})
	if err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	at := issuer.now()
	if decision.ExpiresAt == nil {
		return trust.ExecutionEnvelope{}, false, publicationForbidden(request.CorrelationID)
	}
	if err := decision.Authorize("publication", request.Workload.SubjectID(), string(ActionPublishPullRequest), subject.CandidateDigest, at, issuer.profile.EnforcedObligationTypes, request.CorrelationID); err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	expiresAt := earliestPolicyExpiry(*decision.ExpiresAt, approval.ExpiresAt, subject.ExpiresAt, at.Add(issuer.profile.MaximumLifetime))
	envelope := trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: issuer.profile.Issuer,
		IssuedAt: at, ExpiresAt: expiresAt,
		Scope:     trust.ExecutionScope{TenantID: subject.Scope.TenantID, WorkspaceID: subject.Scope.WorkspaceID, ProjectID: subject.Scope.ProjectID, DeliveryRunID: subject.RunID},
		Principal: trust.EnvelopePrincipal{SubjectID: request.Workload.SubjectID(), AgentID: request.Workload.SubjectID(), DelegationChain: []string{}},
		Purpose:   subject.Purpose, RiskClass: subject.RiskClass, Actions: []string{string(ActionPublishPullRequest)},
		Classifications: []string{subject.Classification},
		Destinations:    []string{"contract-double:" + subject.Scope.RepositoryID},
		ModelIDs:        []string{}, ToolIDs: []string{}, SecretReferenceIDs: []string{},
		PolicyDecisionIDs: []string{decision.ID}, ApprovalRecordIDs: []string{approval.ID},
		Limits: trust.ExecutionLimits{
			MaxDurationSeconds: minInt64(300, int64(expiresAt.Sub(at).Seconds())),
			MaxCostMinorUnits:  0, MaxConcurrency: 1,
		},
		IdempotencyKey: request.IdempotencyKey, EvidenceCorrelationID: request.CorrelationID,
	}
	envelope, err = trust.BindExecutionEnvelopeIntegrity(envelope, request.CorrelationID)
	if err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	fingerprint := publicationEnvelopeFingerprint(envelope, subject.CandidateDigest, request.BaseBranch, request.BaseRevision)
	stored, replayed, err := issuer.store.Create(ctx, envelope, subject.Scope.TenantID, request.IdempotencyKey, fingerprint)
	if err != nil {
		return trust.ExecutionEnvelope{}, false, envelopeUnavailable(request.CorrelationID)
	}
	if err := stored.VerifyIntegrity(request.CorrelationID); err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	return stored, replayed, nil
}

func publicationEnvelopeFingerprint(envelope trust.ExecutionEnvelope, candidateDigest, baseBranch, baseRevision string) string {
	envelope.ID, envelope.IdempotencyKey, envelope.EvidenceCorrelationID = "", "", ""
	envelope.Integrity = nil
	encoded, _ := json.Marshal(struct {
		Envelope        trust.ExecutionEnvelope
		CandidateDigest string
		BaseBranch      string
		BaseRevision    string
	}{envelope, candidateDigest, baseBranch, baseRevision})
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func canonicalPublicationStrings(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return slices.Compact(result)
}

func publicationForbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "publication authority is not permitted", correlationID, failure.RetryNewAuthority)
}
