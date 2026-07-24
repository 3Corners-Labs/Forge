package policy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"time"

	changecontract "github.com/3corners/forge/internal/contracts/change"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const IssueChangeAuthorityVersion = "forge.issue-change-proposal-authority.v1"

type ChangeAuthorityProfile struct {
	Issuer                  string
	Provider                string
	ModelID                 string
	CredentialRef           string
	MaximumLifetime         time.Duration
	EnforcedObligationTypes []string
}

type ChangeAuthorityIssuer struct {
	evaluator *Evaluator
	resolver  DelegationResolver
	store     EnvelopeStore
	profile   ChangeAuthorityProfile
	now       func() time.Time
}

func NewChangeAuthorityIssuer(evaluator *Evaluator, resolver DelegationResolver, store EnvelopeStore, profile ChangeAuthorityProfile, now func() time.Time) (*ChangeAuthorityIssuer, error) {
	if evaluator == nil || resolver == nil || store == nil || now == nil || profile.Issuer == "" ||
		profile.Provider == "" || profile.ModelID == "" || profile.CredentialRef == "" ||
		profile.MaximumLifetime <= 0 ||
		!slices.Equal(profile.EnforcedObligationTypes, []string{"model-output-non-authoritative", "purpose-minimized-context"}) {
		return nil, errors.New("change-proposal authority configuration is invalid")
	}
	return &ChangeAuthorityIssuer{evaluator: evaluator, resolver: resolver, store: store, profile: profile, now: now}, nil
}

type IssueChangeAuthorityRequest struct {
	SchemaVersion  string
	CorrelationID  string
	IdempotencyKey string
	Human          identity.Principal
	TaskAgent      identity.Principal
	Scope          tenant.Scope
	Run            delivery.DeliveryRun
	Request        changecontract.Request
	DelegationID   string
	Environment    Environment
	Risk           Risk
}

type ChangeAuthorityResult struct {
	Decision trust.PolicyDecision
	Envelope trust.ExecutionEnvelope
	Replayed bool
}

func (issuer *ChangeAuthorityIssuer) Issue(ctx context.Context, command IssueChangeAuthorityRequest) (ChangeAuthorityResult, error) {
	at := issuer.now()
	if ctx == nil || command.SchemaVersion != IssueChangeAuthorityVersion ||
		command.CorrelationID == "" || command.IdempotencyKey == "" || command.DelegationID == "" ||
		command.Environment != DeveloperTrialEnvironment || !validRisk(command.Risk) ||
		at.IsZero() || command.Run.State() != delivery.PlanApprovedState ||
		command.Run.Plan() == nil || command.Run.Context() == nil || command.Run.Approval() == nil {
		return ChangeAuthorityResult{}, failure.New(failure.InvalidRequest, "change-proposal authority request is invalid", command.CorrelationID, failure.RetryNone)
	}
	if err := command.Human.ValidateAt(at, command.CorrelationID); err != nil {
		return ChangeAuthorityResult{}, err
	}
	if err := command.TaskAgent.ValidateAt(at, command.CorrelationID); err != nil {
		return ChangeAuthorityResult{}, err
	}
	if err := command.Scope.ValidateFor(command.Human, at, command.CorrelationID); err != nil {
		return ChangeAuthorityResult{}, err
	}
	chain, err := issuer.resolver.ResolveDelegationChain(ctx, identity.ResolveDelegationRequest{
		ContractVersion: identity.ResolveDelegationContractVersion, CorrelationID: command.CorrelationID,
		Human: command.Human, Agent: command.TaskAgent, DelegationID: command.DelegationID,
	})
	if err != nil || len(chain) != 1 || !validChangeDelegation(chain[0], command, issuer.profile.Provider) {
		return ChangeAuthorityResult{}, failure.New(failure.Forbidden, "change-proposal delegation is invalid", command.CorrelationID, failure.RetryNewAuthority)
	}
	if err := command.Request.Validate(at); err != nil || !validChangeBindings(command) {
		return ChangeAuthorityResult{}, failure.New(failure.Forbidden, "change-proposal facts exceed the approved plan", command.CorrelationID, failure.RetryNewAuthority)
	}
	matching := make([]Rule, 0, len(issuer.evaluator.set.Rules))
	for _, rule := range issuer.evaluator.set.Rules {
		if matchesChangeRule(rule, command) {
			matching = append(matching, rule)
		}
	}
	outcome, obligations, reasons := issuer.evaluator.resolve(matching)
	expiresAt := earliestPolicyExpiry(
		at.Add(issuer.evaluator.decisionTTL), command.Human.ExpiresAt(),
		command.TaskAgent.ExpiresAt(), command.Request.ExpiresAt, chain[0].ExpiresAt(),
	)
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, DecisionPoint: "intelligence",
		SubjectID: command.TaskAgent.SubjectID(), Action: string(ActionProposeChange),
		Resource: command.Request.Digest, EvaluatedAt: at, ExpiresAt: &expiresAt,
		PolicyVersions: []string{issuer.evaluator.set.Version}, Outcome: outcome,
		Obligations: obligations, ReasonCodes: reasons, EvidenceCorrelationID: command.CorrelationID,
	}
	decision.ID = changeDecisionID(decision, command.Run.ID(), matching)
	if err := decision.Validate(command.CorrelationID); err != nil {
		return ChangeAuthorityResult{}, failure.New(failure.Indeterminate, "change-proposal policy decision is invalid", command.CorrelationID, failure.RetryNone)
	}
	if err := decision.Authorize(
		"intelligence", command.TaskAgent.SubjectID(), string(ActionProposeChange),
		command.Request.Digest, at, issuer.profile.EnforcedObligationTypes, command.CorrelationID,
	); err != nil {
		return ChangeAuthorityResult{}, err
	}
	issuer.evaluator.recordDecision(decision)
	envelopeExpires := earliestPolicyExpiry(expiresAt, at.Add(issuer.profile.MaximumLifetime))
	envelope, err := trust.BindExecutionEnvelopeIntegrity(trust.ExecutionEnvelope{
		SchemaVersion: trust.ExecutionEnvelopeSchemaVersion, Issuer: issuer.profile.Issuer,
		IssuedAt: at, ExpiresAt: envelopeExpires,
		Scope: trust.ExecutionScope{
			TenantID: command.Run.TenantID(), WorkspaceID: command.Run.WorkspaceID(),
			ProjectID: command.Run.ProjectID(), DeliveryRunID: command.Run.ID(),
		},
		Principal: trust.EnvelopePrincipal{
			SubjectID: command.Human.SubjectID(), AgentID: command.TaskAgent.SubjectID(),
			DelegationChain: []string{command.DelegationID},
		},
		Purpose: changecontract.PurposeV1, RiskClass: string(command.Risk),
		Actions: []string{string(ActionProposeChange)}, Classifications: []string{command.Run.Classification()},
		Destinations: []string{"model:" + issuer.profile.Provider}, ModelIDs: []string{issuer.profile.ModelID},
		ToolIDs: []string{}, SecretReferenceIDs: []string{issuer.profile.CredentialRef},
		PolicyDecisionIDs: []string{decision.ID}, ApprovalRecordIDs: []string{command.Run.Approval().RecordID},
		Limits: trust.ExecutionLimits{
			MaxDurationSeconds: int64(envelopeExpires.Sub(at).Seconds()),
			MaxCostMinorUnits:  command.Request.Limits.MaxCostMicros, MaxInputTokens: command.Request.Limits.MaxTokens,
			MaxOutputTokens: command.Request.Limits.MaxTokens, MaxConcurrency: 1,
		},
		IdempotencyKey: command.IdempotencyKey, EvidenceCorrelationID: command.CorrelationID,
	}, command.CorrelationID)
	if err != nil {
		return ChangeAuthorityResult{}, err
	}
	fingerprint := changeAuthorityFingerprint(envelope, command.Request)
	stored, replayed, err := issuer.store.Create(ctx, envelope, command.Run.TenantID(), command.IdempotencyKey, fingerprint)
	if err != nil {
		return ChangeAuthorityResult{}, failure.New(failure.DependencyUnavailable, "change-proposal authority persistence is unavailable", command.CorrelationID, failure.RetryBounded)
	}
	if stored.VerifyIntegrity(command.CorrelationID) != nil || stored.ID != envelope.ID {
		return ChangeAuthorityResult{}, failure.New(failure.Indeterminate, "change-proposal authority could not be verified", command.CorrelationID, failure.RetryNone)
	}
	return ChangeAuthorityResult{Decision: decision, Envelope: stored, Replayed: replayed}, nil
}

func validChangeDelegation(delegation identity.Delegation, command IssueChangeAuthorityRequest, provider string) bool {
	return delegation.ID() == command.DelegationID &&
		delegation.HumanSubjectID() == command.Human.SubjectID() &&
		delegation.AgentSubjectID() == command.TaskAgent.SubjectID() &&
		delegation.DeliveryRunID() == command.Run.ID() &&
		delegation.RepositoryID() == command.Run.RepositoryID() &&
		slices.Equal(delegation.Actions(), []string{identity.DelegatedActionProposeChange}) &&
		slices.Equal(delegation.Resources(), []string{command.Request.Digest}) &&
		slices.Equal(delegation.Destinations(), []string{"model:" + provider}) &&
		delegation.MaximumAutonomy() == identity.ChangeProposalOnly
}

func validChangeBindings(command IssueChangeAuthorityRequest) bool {
	plan, contextRecord := command.Run.Plan(), command.Run.Context()
	request := command.Request
	return command.Human.Kind() == identity.HumanPrincipal &&
		command.TaskAgent.Kind() == identity.AgentPrincipal &&
		command.Run.ActorID() == command.Human.SubjectID() &&
		command.Scope.TenantID() == command.Run.TenantID() &&
		command.Scope.WorkspaceID() == command.Run.WorkspaceID() &&
		command.Scope.ProjectID() == command.Run.ProjectID() &&
		request.CorrelationID == command.CorrelationID &&
		request.TenantID == command.Run.TenantID() && request.WorkspaceID == command.Run.WorkspaceID() &&
		request.ProjectID == command.Run.ProjectID() && request.RepositoryID == command.Run.RepositoryID() &&
		request.DeliveryRunID == command.Run.ID() && request.TaskAgentID == command.TaskAgent.SubjectID() &&
		request.PlanID == plan.ID && request.PlanDigest == plan.Digest &&
		request.ContextPackageID == contextRecord.PackageID &&
		request.ContextManifestDigest == contextRecord.ManifestDigest &&
		request.SourceRevision == contextRecord.SourceRevisions[0] &&
		request.Classification == command.Run.Classification() &&
		slices.Equal(request.AllowedWritePaths, plan.ChangePaths)
}

func matchesChangeRule(rule Rule, command IssueChangeAuthorityRequest) bool {
	return containsAction(rule.Actions, ActionProposeChange) &&
		containsEnvironment(rule.Environments, command.Environment) &&
		containsRisk(rule.Risks, command.Risk) &&
		containsString(rule.SubjectIDs, command.TaskAgent.SubjectID()) &&
		containsString(rule.ResourceIDs, command.Request.Digest) &&
		containsString(rule.DeliveryRunIDs, command.Run.ID()) &&
		containsString(rule.Purposes, command.Run.Purpose()) &&
		containsString(rule.Classifications, command.Run.Classification()) &&
		containsString(rule.TenantIDs, command.Scope.TenantID()) &&
		containsString(rule.WorkspaceIDs, command.Scope.WorkspaceID()) &&
		containsString(rule.ProjectIDs, command.Scope.ProjectID())
}

func changeDecisionID(decision trust.PolicyDecision, runID string, rules []Rule) string {
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	encoded, _ := json.Marshal(struct {
		Decision trust.PolicyDecision
		RunID    string
		RuleIDs  []string
	}{decision, runID, ids})
	sum := sha256.Sum256(encoded)
	return "forge.policy-decision." + hex.EncodeToString(sum[:])
}

func changeAuthorityFingerprint(envelope trust.ExecutionEnvelope, request changecontract.Request) string {
	encoded, _ := json.Marshal(struct {
		EnvelopeID    string
		RequestDigest string
	}{envelope.ID, request.Digest})
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}
