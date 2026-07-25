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
	EvaluateWorkerDelegationVersion = "forge.policy-evaluate-worker-delegation.v1"
	EvaluateWorkerExecutionVersion  = "forge.policy-evaluate-worker-execution.v1"
	IssueWorkerEnvelopeVersion      = "forge.issue-worker-envelope.v1"
)

type ApprovalResolver interface {
	ResolveApproval(context.Context, string, string) (trust.ApprovalRecordV2, bool, error)
}

type EvaluateWorkerDelegationRequest struct {
	SchemaVersion string
	CorrelationID string
	Human         identity.Principal
	TaskAgent     identity.Principal
	Worker        identity.Principal
	Scope         tenant.Scope
	Run           delivery.DeliveryRun
	Request       identity.WorkerDelegationRequest
	Environment   Environment
	Risk          Risk
}

func (evaluator *Evaluator) EvaluateWorkerDelegation(ctx context.Context, approvals ApprovalResolver, request EvaluateWorkerDelegationRequest) (trust.PolicyDecision, error) {
	at, err := evaluator.workerEvaluationTime(ctx, request.SchemaVersion, EvaluateWorkerDelegationVersion, request.CorrelationID)
	if err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := validateWorkerFacts(request, at); err != nil {
		return trust.PolicyDecision{}, err
	}
	record, err := resolvePlanApproval(ctx, approvals, request.Run, request.Request.ApprovalRecordID, at, request.CorrelationID)
	if err != nil {
		return trust.PolicyDecision{}, err
	}
	if request.Request.ApprovalRecordDigest != record.Integrity.Digest {
		return trust.PolicyDecision{}, workerForbidden(request.CorrelationID)
	}
	matching := evaluator.workerRules(request.Human.SubjectID(), ActionCreateDelegation, request.Request.Digest, request.Run, request.Scope, request.Environment, request.Risk)
	outcome, obligations, reasons := evaluator.resolve(matching)
	expiresAt := earliestPolicyExpiry(at.Add(evaluator.decisionTTL), request.Human.ExpiresAt(), request.TaskAgent.ExpiresAt(), request.Worker.ExpiresAt(), request.Request.RequestedExpiresAt, record.ExpiresAt)
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, DecisionPoint: "execution",
		SubjectID: request.Human.SubjectID(), Action: string(ActionCreateDelegation),
		Resource: request.Request.Digest, EvaluatedAt: at, ExpiresAt: &expiresAt,
		PolicyVersions: []string{evaluator.set.Version}, Outcome: outcome, Obligations: obligations,
		ReasonCodes: reasons, EvidenceCorrelationID: request.CorrelationID,
	}
	decision.ID = workerDecisionID(decision, request.Run.ID(), request.Request.Digest, matching)
	if err := decision.Validate(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, failure.New(failure.Indeterminate, "worker delegation decision could not be verified", request.CorrelationID, failure.RetryNone)
	}
	evaluator.recordDecision(decision)
	return decision, nil
}

type WorkerDelegationResolver interface {
	ResolveWorkerDelegation(context.Context, identity.ResolveWorkerDelegationRequest) (identity.WorkerDelegation, error)
}

type EvaluateWorkerExecutionRequest struct {
	SchemaVersion string
	CorrelationID string
	Human         identity.Principal
	TaskAgent     identity.Principal
	Worker        identity.Principal
	Scope         tenant.Scope
	Run           delivery.DeliveryRun
	DelegationID  string
	Environment   Environment
	Risk          Risk
}

func (evaluator *Evaluator) EvaluateWorkerExecution(ctx context.Context, approvals ApprovalResolver, resolver WorkerDelegationResolver, request EvaluateWorkerExecutionRequest) (trust.PolicyDecision, identity.WorkerDelegation, error) {
	at, err := evaluator.workerEvaluationTime(ctx, request.SchemaVersion, EvaluateWorkerExecutionVersion, request.CorrelationID)
	if err != nil {
		return trust.PolicyDecision{}, identity.WorkerDelegation{}, err
	}
	if resolver == nil || approvals == nil || request.DelegationID == "" {
		return trust.PolicyDecision{}, identity.WorkerDelegation{}, workerForbidden(request.CorrelationID)
	}
	if err := validateWorkerPrincipalsAndRun(request.Human, request.TaskAgent, request.Worker, request.Scope, request.Run, at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, identity.WorkerDelegation{}, err
	}
	delegation, err := resolver.ResolveWorkerDelegation(ctx, identity.ResolveWorkerDelegationRequest{
		SchemaVersion: identity.ResolveWorkerDelegationVersion, CorrelationID: request.CorrelationID,
		Human: request.Human, TaskAgent: request.TaskAgent, Worker: request.Worker, DelegationID: request.DelegationID,
	})
	if err != nil {
		return trust.PolicyDecision{}, identity.WorkerDelegation{}, err
	}
	record, err := resolvePlanApproval(ctx, approvals, request.Run, delegation.ApprovalRecordID, at, request.CorrelationID)
	if err != nil || record.Integrity.Digest != delegation.ApprovalRecordDigest ||
		delegation.PlanDigest != request.Run.Plan().Digest ||
		delegation.ContextManifestDigest != request.Run.Context().ManifestDigest ||
		delegation.SourceRevision != request.Run.Context().SourceRevisions[0] {
		return trust.PolicyDecision{}, identity.WorkerDelegation{}, workerForbidden(request.CorrelationID)
	}
	matching := evaluator.workerRules(request.Worker.SubjectID(), ActionPrepareChange, delegation.ChangeProposalDigest, request.Run, request.Scope, request.Environment, request.Risk)
	outcome, obligations, reasons := evaluator.resolve(matching)
	expiresAt := earliestPolicyExpiry(at.Add(evaluator.decisionTTL), request.Human.ExpiresAt(), request.TaskAgent.ExpiresAt(), request.Worker.ExpiresAt(), delegation.ExpiresAt, record.ExpiresAt)
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, DecisionPoint: "execution",
		SubjectID: request.Worker.SubjectID(), Action: string(ActionPrepareChange),
		Resource: delegation.ChangeProposalDigest, EvaluatedAt: at, ExpiresAt: &expiresAt,
		PolicyVersions: []string{evaluator.set.Version}, Outcome: outcome, Obligations: obligations,
		ReasonCodes: reasons, EvidenceCorrelationID: request.CorrelationID,
	}
	decision.ID = workerDecisionID(decision, request.Run.ID(), delegation.ID, matching)
	if err := decision.Validate(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, identity.WorkerDelegation{}, failure.New(failure.Indeterminate, "worker execution decision could not be verified", request.CorrelationID, failure.RetryNone)
	}
	evaluator.recordDecision(decision)
	return decision, delegation, nil
}

type WorkerEnvelopeProfile struct {
	Issuer                  string
	MaximumLifetime         time.Duration
	EnforcedObligationTypes []string
}

type WorkerEnvelopeIssuer struct {
	evaluator *Evaluator
	approvals ApprovalResolver
	resolver  WorkerDelegationResolver
	store     EnvelopeStore
	profile   WorkerEnvelopeProfile
	now       func() time.Time
}

func NewWorkerEnvelopeIssuer(evaluator *Evaluator, approvals ApprovalResolver, resolver WorkerDelegationResolver, store EnvelopeStore, profile WorkerEnvelopeProfile, now func() time.Time) (*WorkerEnvelopeIssuer, error) {
	if evaluator == nil || approvals == nil || resolver == nil || store == nil || now == nil ||
		profile.Issuer == "" || profile.MaximumLifetime <= 0 {
		return nil, errors.New("worker envelope issuer configuration is invalid")
	}
	return &WorkerEnvelopeIssuer{evaluator: evaluator, approvals: approvals, resolver: resolver, store: store, profile: profile, now: now}, nil
}

type IssueWorkerEnvelopeRequest struct {
	SchemaVersion  string
	CorrelationID  string
	IdempotencyKey string
	Human          identity.Principal
	TaskAgent      identity.Principal
	Worker         identity.Principal
	Scope          tenant.Scope
	Run            delivery.DeliveryRun
	DelegationID   string
	Environment    Environment
	Risk           Risk
}

func (issuer *WorkerEnvelopeIssuer) Issue(ctx context.Context, request IssueWorkerEnvelopeRequest) (trust.ExecutionEnvelope, bool, error) {
	if request.SchemaVersion != IssueWorkerEnvelopeVersion || request.IdempotencyKey == "" {
		return trust.ExecutionEnvelope{}, false, failure.New(failure.InvalidRequest, "worker envelope request is invalid", request.CorrelationID, failure.RetryNone)
	}
	decision, delegation, err := issuer.evaluator.EvaluateWorkerExecution(ctx, issuer.approvals, issuer.resolver, EvaluateWorkerExecutionRequest{
		SchemaVersion: EvaluateWorkerExecutionVersion, CorrelationID: request.CorrelationID,
		Human: request.Human, TaskAgent: request.TaskAgent, Worker: request.Worker,
		Scope: request.Scope, Run: request.Run, DelegationID: request.DelegationID,
		Environment: request.Environment, Risk: request.Risk,
	})
	if err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	at := issuer.now()
	if decision.ExpiresAt == nil {
		return trust.ExecutionEnvelope{}, false, workerForbidden(request.CorrelationID)
	}
	if err := decision.Authorize(
		"execution", request.Worker.SubjectID(), string(ActionPrepareChange),
		delegation.ChangeProposalDigest, at, issuer.profile.EnforcedObligationTypes,
		request.CorrelationID,
	); err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	expiresAt := earliestPolicyExpiry(*decision.ExpiresAt, delegation.ExpiresAt, at.Add(issuer.profile.MaximumLifetime))
	envelope := trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: issuer.profile.Issuer,
		IssuedAt: at, ExpiresAt: expiresAt,
		Scope:     trust.ExecutionScope{TenantID: request.Run.TenantID(), WorkspaceID: request.Run.WorkspaceID(), ProjectID: request.Run.ProjectID(), DeliveryRunID: request.Run.ID()},
		Principal: trust.EnvelopePrincipal{SubjectID: request.Human.SubjectID(), AgentID: request.Worker.SubjectID(), DelegationChain: []string{delegation.ID}},
		Purpose:   identity.WorkerM1Purpose, RiskClass: string(request.Risk),
		Actions: []string{identity.WorkerPrepareChangeAction}, Classifications: []string{request.Run.Classification()},
		Destinations: []string{}, ModelIDs: []string{}, ToolIDs: []string{}, SecretReferenceIDs: []string{},
		PolicyDecisionIDs: []string{delegation.ParentPolicyDecisionID, decision.ID},
		ApprovalRecordIDs: []string{delegation.ApprovalRecordID},
		Limits:            trust.ExecutionLimits{MaxDurationSeconds: minInt64(delegation.MaxDurationSeconds, int64(expiresAt.Sub(at).Seconds())), MaxCostMinorUnits: 0, MaxConcurrency: 1},
		IdempotencyKey:    request.IdempotencyKey, EvidenceCorrelationID: request.CorrelationID,
	}
	sort.Strings(envelope.PolicyDecisionIDs)
	envelope, err = trust.BindExecutionEnvelopeIntegrity(envelope, request.CorrelationID)
	if err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	fingerprint := workerEnvelopeFingerprint(envelope, delegation)
	stored, replayed, err := issuer.store.Create(ctx, envelope, request.Run.TenantID(), request.IdempotencyKey, fingerprint)
	if err != nil {
		return trust.ExecutionEnvelope{}, false, envelopeUnavailable(request.CorrelationID)
	}
	if err := stored.VerifyIntegrity(request.CorrelationID); err != nil {
		return trust.ExecutionEnvelope{}, false, err
	}
	if stored.Scope != envelope.Scope || stored.Principal.AgentID != request.Worker.SubjectID() ||
		!slices.Equal(stored.Actions, []string{identity.WorkerPrepareChangeAction}) ||
		!slices.Equal(stored.ApprovalRecordIDs, []string{delegation.ApprovalRecordID}) {
		return trust.ExecutionEnvelope{}, false, failure.New(failure.Indeterminate, "worker envelope could not be verified", request.CorrelationID, failure.RetryNone)
	}
	return stored, replayed, nil
}

func (evaluator *Evaluator) workerEvaluationTime(ctx context.Context, actual, expected, correlationID string) (time.Time, error) {
	if evaluator == nil || evaluator.now == nil || ctx == nil || actual != expected || correlationID == "" {
		return time.Time{}, policyUnavailable(correlationID)
	}
	if err := ctx.Err(); err != nil {
		return time.Time{}, policyUnavailable(correlationID)
	}
	at := evaluator.now()
	if at.IsZero() {
		return time.Time{}, policyUnavailable(correlationID)
	}
	return at, nil
}

func validateWorkerFacts(request EvaluateWorkerDelegationRequest, at time.Time) error {
	if err := validateWorkerPrincipalsAndRun(request.Human, request.TaskAgent, request.Worker, request.Scope, request.Run, at, request.CorrelationID); err != nil {
		return err
	}
	if err := request.Request.Validate(at); err != nil {
		return failure.New(failure.InvalidRequest, "worker delegation facts are invalid", request.CorrelationID, failure.RetryNone)
	}
	plan, contextRecord, approval := request.Run.Plan(), request.Run.Context(), request.Run.Approval()
	if request.Environment != DeveloperTrialEnvironment || !validRisk(request.Risk) ||
		request.Request.HumanSubjectID != request.Human.SubjectID() ||
		request.Request.TaskAgentSubjectID != request.TaskAgent.SubjectID() ||
		request.Request.WorkerSubjectID != request.Worker.SubjectID() ||
		request.Request.TenantID != request.Run.TenantID() || request.Request.WorkspaceID != request.Run.WorkspaceID() ||
		request.Request.ProjectID != request.Run.ProjectID() || request.Request.RepositoryID != request.Run.RepositoryID() ||
		request.Request.DeliveryRunID != request.Run.ID() || request.Request.PlanID != plan.ID ||
		request.Request.PlanDigest != plan.Digest || request.Request.ContextPackageID != contextRecord.PackageID ||
		request.Request.ContextManifestDigest != contextRecord.ManifestDigest ||
		request.Request.ApprovalRequestID != approval.RequestID ||
		request.Request.ApprovalRequestDigest != approval.RequestDigest ||
		request.Request.ApprovalRecordID != approval.RecordID ||
		request.Request.SourceRevision != contextRecord.SourceRevisions[0] ||
		!slices.Equal(request.Request.ChangePaths, plan.ChangePaths) ||
		request.Request.Classification != request.Run.Classification() ||
		request.Request.RiskClass != plan.RiskClass {
		return workerForbidden(request.CorrelationID)
	}
	return nil
}

func validateWorkerPrincipalsAndRun(human, taskAgent, worker identity.Principal, scope tenant.Scope, run delivery.DeliveryRun, at time.Time, correlationID string) error {
	for _, principal := range []identity.Principal{human, taskAgent, worker} {
		if err := principal.ValidateAt(at, correlationID); err != nil {
			return err
		}
	}
	if human.Kind() != identity.HumanPrincipal || taskAgent.Kind() != identity.AgentPrincipal ||
		worker.Kind() != identity.WorkloadPrincipal || run.State() != delivery.PlanApprovedState ||
		run.ActorID() != human.SubjectID() || run.Plan() == nil || run.Context() == nil || run.Approval() == nil ||
		run.TenantID() != scope.TenantID() || run.WorkspaceID() != scope.WorkspaceID() ||
		run.ProjectID() != scope.ProjectID() {
		return workerForbidden(correlationID)
	}
	return nil
}

func resolvePlanApproval(ctx context.Context, resolver ApprovalResolver, run delivery.DeliveryRun, recordID string, at time.Time, correlationID string) (trust.ApprovalRecordV2, error) {
	if resolver == nil || recordID == "" {
		return trust.ApprovalRecordV2{}, workerForbidden(correlationID)
	}
	record, found, err := resolver.ResolveApproval(ctx, run.TenantID(), recordID)
	if err != nil {
		return trust.ApprovalRecordV2{}, failure.New(failure.DependencyUnavailable, "plan approval resolution is unavailable", correlationID, failure.RetryBounded)
	}
	plan, approval := run.Plan(), run.Approval()
	if !found || plan == nil || approval == nil || record.ID != recordID {
		return trust.ApprovalRecordV2{}, workerForbidden(correlationID)
	}
	expected := trust.ApprovalExpectedFacts{
		RequestID: approval.RequestID, RequestDigest: approval.RequestDigest, ApproverID: run.ActorID(),
		Scope:       trust.ApprovalScope{TenantID: run.TenantID(), WorkspaceID: run.WorkspaceID(), ProjectID: run.ProjectID(), RepositoryID: run.RepositoryID(), DeliveryRunID: run.ID()},
		Environment: "developer-trial", Purpose: run.Purpose(), Classification: run.Classification(),
		RiskClass: plan.RiskClass, Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan",
		SubjectID: plan.ID, SubjectDigest: plan.Digest,
		EvidenceReferences:  []string{"context:" + plan.ContextPackageID, "context-manifest:" + plan.ContextManifest, "intent:version:1", "plan:" + plan.Digest},
		PermittedConditions: []string{"require-independent-verification"},
	}
	if err := record.Authorize(expected, at, correlationID); err != nil {
		return trust.ApprovalRecordV2{}, err
	}
	return record, nil
}

func (evaluator *Evaluator) workerRules(subject string, action Action, resource string, run delivery.DeliveryRun, scope tenant.Scope, environment Environment, risk Risk) []Rule {
	result := make([]Rule, 0, len(evaluator.set.Rules))
	for _, rule := range evaluator.set.Rules {
		if containsAction(rule.Actions, action) && containsEnvironment(rule.Environments, environment) &&
			containsRisk(rule.Risks, risk) && containsString(rule.SubjectIDs, subject) &&
			containsString(rule.ResourceIDs, resource) && containsString(rule.DeliveryRunIDs, run.ID()) &&
			containsString(rule.Purposes, run.Purpose()) && containsString(rule.Classifications, run.Classification()) &&
			containsString(rule.TenantIDs, scope.TenantID()) && containsString(rule.WorkspaceIDs, scope.WorkspaceID()) &&
			containsString(rule.ProjectIDs, scope.ProjectID()) {
			result = append(result, rule)
		}
	}
	return result
}

func workerDecisionID(decision trust.PolicyDecision, runID, binding string, rules []Rule) string {
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	encoded, _ := json.Marshal(struct {
		Decision trust.PolicyDecision
		RunID    string
		Binding  string
		RuleIDs  []string
	}{decision, runID, binding, ids})
	sum := sha256.Sum256(encoded)
	return "forge.policy-decision." + hex.EncodeToString(sum[:])
}

func workerEnvelopeFingerprint(envelope trust.ExecutionEnvelope, delegation identity.WorkerDelegation) string {
	envelope.ID, envelope.IdempotencyKey, envelope.EvidenceCorrelationID = "", "", ""
	envelope.Integrity = nil
	encoded, _ := json.Marshal(struct {
		Envelope   trust.ExecutionEnvelope
		Delegation string
	}{envelope, delegation.ID})
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func earliestPolicyExpiry(values ...time.Time) time.Time {
	result := values[0]
	for _, value := range values[1:] {
		if value.Before(result) {
			result = value
		}
	}
	return result
}

func minInt64(left, right int64) int64 {
	if left < right {
		return left
	}
	return right
}

func workerForbidden(correlationID string) error {
	return failure.New(failure.Forbidden, "worker authority is not permitted", correlationID, failure.RetryNewAuthority)
}
