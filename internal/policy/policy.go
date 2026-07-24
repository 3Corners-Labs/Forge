// Package policy evaluates trusted, versioned policy against immutable facts.
// It owns decisions and obligations, but not identity, scope, delivery state,
// execution authority, enforcement, approval, or business state.
package policy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const (
	EvaluateRequestContractVersion           = "forge.policy-evaluate-request.v1"
	EvaluateDelegationRequestContractVersion = "forge.policy-evaluate-delegation-request.v1"
	EvaluateApprovalRequestContractVersion   = "forge.policy-evaluate-approval-request.v1"
	DeveloperTrialEnvironment                = Environment("developer-trial")
	ActionAssembleContext                    = Action("forge.context.assemble")
	ActionCreateDelegation                   = Action("forge.delegation.create")
	ActionDecideApproval                     = Action("forge.approval.decide")
	ActionProposeChange                      = Action("forge.intelligence.propose-change")
	ActionPrepareChange                      = Action("forge.worker.prepare-change")
	ActionPublishPullRequest                 = Action("forge.source-control.publish-pr")
)

type Action string
type Environment string
type Risk string
type Effect string

const (
	RiskLow      Risk = "low"
	RiskModerate Risk = "moderate"
	RiskHigh     Risk = "high"
	RiskCritical Risk = "critical"

	EffectPermit Effect = "permit"
	EffectDeny   Effect = "deny"
)

// Obligation is a mandatory runtime requirement. Optional policy advice is
// intentionally outside the Policy Decision obligation list.
type Obligation struct {
	Type       string
	Parameters map[string]any
}

// Rule is an immutable, trusted policy definition supplied to an Evaluator.
// Empty match lists are wildcards. More-specific permit rules can add
// obligations, but no permit can override a matching deny.
type Rule struct {
	ID              string
	Effect          Effect
	Mandatory       bool
	Actions         []Action
	Environments    []Environment
	Risks           []Risk
	SubjectIDs      []string
	ResourceIDs     []string
	DeliveryRunIDs  []string
	Purposes        []string
	Classifications []string
	TenantIDs       []string
	WorkspaceIDs    []string
	ProjectIDs      []string
	Obligations     []Obligation
}

type PolicySet struct {
	Version string
	Rules   []Rule
}

// EvaluateRequest carries untrusted requested action/environment/risk facts
// alongside capability-minted identity, scope, and delivery facts. Policy
// validates and binds the exact values; it does not trust them as authority.
type EvaluateRequest struct {
	ContractVersion string
	CorrelationID   string
	Principal       identity.Principal
	Scope           tenant.Scope
	Run             delivery.DeliveryRun
	Action          Action
	Environment     Environment
	Risk            Risk
}

// EvaluateDelegationRequest binds a separately verified task agent and the
// immutable Identity-owned request to human, Tenant, and Delivery authority.
// The delegation request remains non-authoritative until this exact digest is
// permitted and Identity mints the resulting delegation.
type EvaluateDelegationRequest struct {
	ContractVersion string
	CorrelationID   string
	Human           identity.Principal
	Agent           identity.Principal
	Scope           tenant.Scope
	Run             delivery.DeliveryRun
	Delegation      identity.DelegationRequest
	Environment     Environment
	Risk            Risk
}

type EvaluateApprovalRequest struct {
	ContractVersion string
	CorrelationID   string
	Principal       identity.Principal
	Scope           tenant.Scope
	Run             delivery.DeliveryRun
	Request         trust.ApprovalRequest
	Environment     Environment
	Risk            Risk
}

type Evaluator struct {
	set                  PolicySet
	supportedObligations map[string]struct{}
	decisionTTL          time.Duration
	now                  func() time.Time
	decisionMu           sync.RWMutex
	decisions            map[string]trust.PolicyDecision
}

func NewEvaluator(set PolicySet, supportedObligations []string, decisionTTL time.Duration, now func() time.Time) (*Evaluator, error) {
	if set.Version == "" {
		return nil, errors.New("policy set version is required")
	}
	if len(set.Rules) == 0 {
		return nil, errors.New("policy set rules are required")
	}
	if decisionTTL <= 0 {
		return nil, errors.New("policy decision TTL must be positive")
	}
	if now == nil {
		return nil, errors.New("policy clock is required")
	}
	rules := make([]Rule, len(set.Rules))
	seenRuleIDs := make(map[string]struct{}, len(set.Rules))
	for index, rule := range set.Rules {
		if err := validateRule(rule, seenRuleIDs); err != nil {
			return nil, err
		}
		rules[index] = cloneRule(rule)
	}
	sort.Slice(rules, func(left, right int) bool { return rules[left].ID < rules[right].ID })
	supported := make(map[string]struct{}, len(supportedObligations))
	for _, obligationType := range supportedObligations {
		if obligationType == "" {
			return nil, errors.New("supported obligation type cannot be empty")
		}
		if _, duplicate := supported[obligationType]; duplicate {
			return nil, errors.New("supported obligation types must be unique")
		}
		supported[obligationType] = struct{}{}
	}
	return &Evaluator{
		set:                  PolicySet{Version: set.Version, Rules: rules},
		supportedObligations: supported,
		decisionTTL:          decisionTTL,
		now:                  now,
		decisions:            make(map[string]trust.PolicyDecision),
	}, nil
}

func (evaluator *Evaluator) Evaluate(ctx context.Context, request EvaluateRequest) (trust.PolicyDecision, error) {
	if evaluator == nil || evaluator.now == nil || evaluator.set.Version == "" || len(evaluator.set.Rules) == 0 {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	if ctx == nil || request.ContractVersion != EvaluateRequestContractVersion || request.CorrelationID == "" {
		return trust.PolicyDecision{}, invalidRequest(request.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	at := evaluator.now()
	if at.IsZero() {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	if err := request.Principal.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := request.Scope.ValidateFor(request.Principal, at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if !validRequestFacts(request) {
		return trust.PolicyDecision{}, invalidRequest(request.CorrelationID)
	}
	if !validRunBinding(request) {
		return trust.PolicyDecision{}, failure.New(
			failure.Forbidden,
			"policy facts are outside the authenticated delivery scope",
			request.CorrelationID,
			failure.RetryNone,
		)
	}

	matching := make([]Rule, 0, len(evaluator.set.Rules))
	for _, rule := range evaluator.set.Rules {
		if matches(rule, request) {
			matching = append(matching, rule)
		}
	}
	outcome, obligations, reasonCodes := evaluator.resolve(matching)
	expiresAt := at.Add(evaluator.decisionTTL)
	if principalExpiry := request.Principal.ExpiresAt(); principalExpiry.Before(expiresAt) {
		expiresAt = principalExpiry
	}
	decision := trust.PolicyDecision{
		SchemaVersion:         trust.PolicyDecisionSchemaVersion,
		DecisionPoint:         "retrieval",
		SubjectID:             request.Principal.SubjectID(),
		Action:                string(request.Action),
		Resource:              request.Run.RepositoryID(),
		EvaluatedAt:           at,
		ExpiresAt:             &expiresAt,
		PolicyVersions:        []string{evaluator.set.Version},
		Outcome:               outcome,
		Obligations:           obligations,
		ReasonCodes:           reasonCodes,
		EvidenceCorrelationID: request.CorrelationID,
	}
	decision.ID = decisionID(request, decision, matching)
	if err := decision.Validate(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, failure.New(
			failure.Indeterminate,
			"policy decision could not be verified",
			request.CorrelationID,
			failure.RetryNone,
		)
	}
	evaluator.recordDecision(decision)
	return decision, nil
}

// EvaluateDelegation evaluates only the bounded M1 root-delegation action.
// It deliberately does not mint or persist delegation authority.
func (evaluator *Evaluator) EvaluateDelegation(ctx context.Context, request EvaluateDelegationRequest) (trust.PolicyDecision, error) {
	if evaluator == nil || evaluator.now == nil || evaluator.set.Version == "" || len(evaluator.set.Rules) == 0 {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	if ctx == nil || request.ContractVersion != EvaluateDelegationRequestContractVersion || request.CorrelationID == "" {
		return trust.PolicyDecision{}, invalidRequest(request.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	at := evaluator.now()
	if at.IsZero() {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	if err := request.Human.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := request.Agent.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := request.Scope.ValidateFor(request.Human, at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := request.Delegation.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if !validDelegationFacts(request) {
		return trust.PolicyDecision{}, invalidRequest(request.CorrelationID)
	}
	if !validDelegationBinding(request) {
		return trust.PolicyDecision{}, failure.New(
			failure.Forbidden,
			"delegation facts are outside the authenticated delivery scope",
			request.CorrelationID,
			failure.RetryNone,
		)
	}

	matching := make([]Rule, 0, len(evaluator.set.Rules))
	for _, rule := range evaluator.set.Rules {
		if matchesDelegation(rule, request) {
			matching = append(matching, rule)
		}
	}
	outcome, obligations, reasonCodes := evaluator.resolve(matching)
	expiresAt := at.Add(evaluator.decisionTTL)
	for _, candidate := range []time.Time{
		request.Human.ExpiresAt(), request.Agent.ExpiresAt(), request.Delegation.RequestedExpiresAt(),
	} {
		if candidate.Before(expiresAt) {
			expiresAt = candidate
		}
	}
	decision := trust.PolicyDecision{
		SchemaVersion:         trust.PolicyDecisionSchemaVersion,
		DecisionPoint:         "execution",
		SubjectID:             request.Human.SubjectID(),
		Action:                string(ActionCreateDelegation),
		Resource:              request.Delegation.Digest(),
		EvaluatedAt:           at,
		ExpiresAt:             &expiresAt,
		PolicyVersions:        []string{evaluator.set.Version},
		Outcome:               outcome,
		Obligations:           obligations,
		ReasonCodes:           reasonCodes,
		EvidenceCorrelationID: request.CorrelationID,
	}
	decision.ID = delegationDecisionID(request, decision, matching)
	if err := decision.Validate(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, failure.New(
			failure.Indeterminate,
			"policy decision could not be verified",
			request.CorrelationID,
			failure.RetryNone,
		)
	}
	evaluator.recordDecision(decision)
	return decision, nil
}

// EvaluateApproval evaluates only whether the exact active human may decide
// one digest-bound plan-review request. It does not decide for the human and
// cannot mint an Approval Record.
func (evaluator *Evaluator) EvaluateApproval(ctx context.Context, request EvaluateApprovalRequest) (trust.PolicyDecision, error) {
	if evaluator == nil || evaluator.now == nil || evaluator.set.Version == "" || len(evaluator.set.Rules) == 0 {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	if ctx == nil || request.ContractVersion != EvaluateApprovalRequestContractVersion || request.CorrelationID == "" {
		return trust.PolicyDecision{}, invalidRequest(request.CorrelationID)
	}
	if err := ctx.Err(); err != nil {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	at := evaluator.now()
	if at.IsZero() {
		return trust.PolicyDecision{}, policyUnavailable(request.CorrelationID)
	}
	if err := request.Principal.ValidateAt(at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := request.Scope.ValidateFor(request.Principal, at, request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if err := request.Request.VerifyDigest(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, err
	}
	if !validApprovalFacts(request, at) {
		return trust.PolicyDecision{}, failure.New(failure.Forbidden, "approval facts are outside the authenticated plan review", request.CorrelationID, failure.RetryNone)
	}
	matching := make([]Rule, 0, len(evaluator.set.Rules))
	for _, rule := range evaluator.set.Rules {
		if matchesApproval(rule, request) {
			matching = append(matching, rule)
		}
	}
	outcome, obligations, reasonCodes := evaluator.resolve(matching)
	expiresAt := at.Add(evaluator.decisionTTL)
	for _, candidate := range []time.Time{request.Principal.ExpiresAt(), request.Request.ExpiresAt} {
		if candidate.Before(expiresAt) {
			expiresAt = candidate
		}
	}
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion, DecisionPoint: "approval",
		SubjectID: request.Principal.SubjectID(), Action: string(ActionDecideApproval),
		Resource: request.Request.RequestDigest, EvaluatedAt: at, ExpiresAt: &expiresAt,
		PolicyVersions: []string{evaluator.set.Version}, Outcome: outcome, Obligations: obligations,
		ReasonCodes: reasonCodes, EvidenceCorrelationID: request.CorrelationID,
	}
	decision.ID = approvalDecisionID(request, decision, matching)
	if err := decision.Validate(request.CorrelationID); err != nil {
		return trust.PolicyDecision{}, failure.New(failure.Indeterminate, "approval eligibility decision could not be verified", request.CorrelationID, failure.RetryNone)
	}
	evaluator.recordDecision(decision)
	return decision, nil
}

// ResolvePolicyDecision returns only decisions minted by this evaluator. It
// lets enforcement capabilities consume a Policy-owned fact without importing
// Policy implementation or trusting a caller-constructed decision object.
func (evaluator *Evaluator) ResolvePolicyDecision(ctx context.Context, id string) (trust.PolicyDecision, bool, error) {
	if evaluator == nil || evaluator.decisions == nil || evaluator.now == nil {
		return trust.PolicyDecision{}, false, errors.New("policy decision source unavailable")
	}
	if ctx == nil || id == "" {
		return trust.PolicyDecision{}, false, errors.New("policy decision lookup invalid")
	}
	if err := ctx.Err(); err != nil {
		return trust.PolicyDecision{}, false, err
	}
	evaluator.decisionMu.RLock()
	defer evaluator.decisionMu.RUnlock()
	decision, found := evaluator.decisions[id]
	return clonePolicyDecision(decision), found, nil
}

func (evaluator *Evaluator) recordDecision(decision trust.PolicyDecision) {
	evaluator.decisionMu.Lock()
	defer evaluator.decisionMu.Unlock()
	evaluator.decisions[decision.ID] = clonePolicyDecision(decision)
}

func clonePolicyDecision(value trust.PolicyDecision) trust.PolicyDecision {
	value.PolicyVersions = slices.Clone(value.PolicyVersions)
	value.ReasonCodes = slices.Clone(value.ReasonCodes)
	obligations := value.Obligations
	value.Obligations = make([]map[string]any, len(obligations))
	for index, obligation := range obligations {
		value.Obligations[index] = cloneMap(obligation)
	}
	if value.ExpiresAt != nil {
		expiresAt := *value.ExpiresAt
		value.ExpiresAt = &expiresAt
	}
	return value
}

func (evaluator *Evaluator) resolve(rules []Rule) (string, []map[string]any, []string) {
	if len(rules) == 0 {
		return "indeterminate", []map[string]any{}, []string{"no-applicable-policy"}
	}
	mandatoryDeny := false
	deny := false
	permit := false
	for _, rule := range rules {
		switch rule.Effect {
		case EffectDeny:
			deny = true
			mandatoryDeny = mandatoryDeny || rule.Mandatory
		case EffectPermit:
			permit = true
		}
	}
	if mandatoryDeny {
		return "deny", []map[string]any{}, []string{"mandatory-policy-deny"}
	}
	if deny {
		return "deny", []map[string]any{}, []string{"policy-deny"}
	}
	if !permit {
		return "indeterminate", []map[string]any{}, []string{"policy-conflict"}
	}

	byType := make(map[string][]byte)
	obligations := make(map[string]map[string]any)
	for _, rule := range rules {
		if rule.Effect != EffectPermit {
			continue
		}
		for _, obligation := range rule.Obligations {
			if _, supported := evaluator.supportedObligations[obligation.Type]; !supported {
				return "deny", []map[string]any{}, []string{"unknown-mandatory-obligation"}
			}
			value := obligationMap(obligation)
			encoded, _ := json.Marshal(value)
			if existing, found := byType[obligation.Type]; found && !slices.Equal(existing, encoded) {
				return "indeterminate", []map[string]any{}, []string{"policy-obligation-conflict"}
			}
			byType[obligation.Type] = encoded
			obligations[obligation.Type] = value
		}
	}
	types := make([]string, 0, len(obligations))
	for obligationType := range obligations {
		types = append(types, obligationType)
	}
	sort.Strings(types)
	result := make([]map[string]any, 0, len(types))
	for _, obligationType := range types {
		result = append(result, cloneMap(obligations[obligationType]))
	}
	return "permit", result, []string{"policy-permit"}
}

func validateRule(rule Rule, seen map[string]struct{}) error {
	if rule.ID == "" {
		return errors.New("policy rule ID is required")
	}
	if _, duplicate := seen[rule.ID]; duplicate {
		return errors.New("policy rule IDs must be unique")
	}
	seen[rule.ID] = struct{}{}
	if rule.Effect != EffectPermit && rule.Effect != EffectDeny {
		return errors.New("policy rule effect is invalid")
	}
	for _, action := range rule.Actions {
		if action != ActionAssembleContext && action != ActionCreateDelegation && action != ActionDecideApproval &&
			action != ActionProposeChange && action != ActionPrepareChange && action != ActionPublishPullRequest {
			return errors.New("policy rule action is invalid")
		}
	}
	for _, environment := range rule.Environments {
		if environment != DeveloperTrialEnvironment {
			return errors.New("policy rule environment is invalid")
		}
	}
	for _, risk := range rule.Risks {
		if !validRisk(risk) {
			return errors.New("policy rule risk is invalid")
		}
	}
	if !allUniqueNonZero(rule.Actions) || !allUniqueNonZero(rule.Environments) || !allUniqueNonZero(rule.Risks) {
		return errors.New("policy rule typed match values must be unique")
	}
	if !allUniqueNonZero(rule.SubjectIDs) || !allUniqueNonZero(rule.ResourceIDs) ||
		!allUniqueNonZero(rule.DeliveryRunIDs) || !allUniqueNonZero(rule.Purposes) ||
		!allUniqueNonZero(rule.Classifications) || !allUniqueNonZero(rule.TenantIDs) ||
		!allUniqueNonZero(rule.WorkspaceIDs) || !allUniqueNonZero(rule.ProjectIDs) {
		return errors.New("policy rule match values must be unique and non-empty")
	}
	for _, classification := range rule.Classifications {
		if !validClassification(classification) {
			return errors.New("policy rule classification is invalid")
		}
	}
	if rule.Effect == EffectDeny && len(rule.Obligations) != 0 {
		return errors.New("deny rules cannot declare obligations")
	}
	seenObligations := make(map[string]struct{}, len(rule.Obligations))
	for _, obligation := range rule.Obligations {
		if obligation.Type == "" {
			return errors.New("policy obligation type is required")
		}
		if _, duplicate := seenObligations[obligation.Type]; duplicate {
			return errors.New("policy rule obligation types must be unique")
		}
		seenObligations[obligation.Type] = struct{}{}
		if _, reserved := obligation.Parameters["type"]; reserved {
			return errors.New("policy obligation parameters cannot replace type")
		}
		if _, err := json.Marshal(obligation.Parameters); err != nil {
			return errors.New("policy obligation parameters must be JSON serializable")
		}
	}
	return nil
}

func validRequestFacts(request EvaluateRequest) bool {
	return request.Action == ActionAssembleContext &&
		request.Environment == DeveloperTrialEnvironment &&
		validRisk(request.Risk) &&
		request.Run.ContractVersion() == delivery.DeliveryRunContractVersion &&
		request.Run.State() == delivery.CapturedState &&
		request.Run.ID() != "" && request.Run.RepositoryID() != "" &&
		request.Run.Purpose() != "" && request.Run.Classification() != ""
}

func validRunBinding(request EvaluateRequest) bool {
	return request.Run.TenantID() == request.Scope.TenantID() &&
		request.Run.WorkspaceID() == request.Scope.WorkspaceID() &&
		request.Run.ProjectID() == request.Scope.ProjectID() &&
		request.Run.ActorID() == request.Principal.SubjectID()
}

func validDelegationFacts(request EvaluateDelegationRequest) bool {
	return request.Environment == DeveloperTrialEnvironment &&
		validRisk(request.Risk) &&
		request.Delegation.CorrelationID() == request.CorrelationID &&
		request.Human.Kind() == identity.HumanPrincipal && request.Human.Synthetic() &&
		request.Agent.Kind() == identity.AgentPrincipal && request.Agent.Synthetic() &&
		request.Human.SubjectID() != request.Agent.SubjectID() &&
		request.Run.ContractVersion() == delivery.DeliveryRunContractVersion &&
		validDelegationRunState(request) && request.Run.ID() != "" &&
		request.Run.RepositoryID() != "" && request.Run.Purpose() != "" &&
		request.Run.Classification() != ""
}

func validDelegationRunState(request EvaluateDelegationRequest) bool {
	if slices.Equal(request.Delegation.Actions(), []string{identity.DelegatedActionAssembleContext}) {
		return request.Run.State() == delivery.CapturedState
	}
	return slices.Equal(request.Delegation.Actions(), []string{identity.DelegatedActionProposeChange}) &&
		request.Run.State() == delivery.PlanApprovedState &&
		request.Run.Plan() != nil && request.Run.Context() != nil && request.Run.Approval() != nil
}

func validDelegationBinding(request EvaluateDelegationRequest) bool {
	delegation := request.Delegation
	return request.Run.TenantID() == request.Scope.TenantID() &&
		request.Run.WorkspaceID() == request.Scope.WorkspaceID() &&
		request.Run.ProjectID() == request.Scope.ProjectID() &&
		request.Run.ActorID() == request.Human.SubjectID() &&
		delegation.HumanSubjectID() == request.Human.SubjectID() &&
		delegation.AgentSubjectID() == request.Agent.SubjectID() &&
		delegation.TenantID() == request.Scope.TenantID() &&
		delegation.WorkspaceID() == request.Scope.WorkspaceID() &&
		delegation.ProjectID() == request.Scope.ProjectID() &&
		delegation.RepositoryID() == request.Run.RepositoryID() &&
		delegation.DeliveryRunID() == request.Run.ID() &&
		validDelegationCapabilityScope(delegation, request.Run.RepositoryID()) &&
		slices.Equal(delegation.Classifications(), []string{request.Run.Classification()})
}

func validDelegationCapabilityScope(delegation identity.DelegationRequest, repositoryID string) bool {
	if slices.Equal(delegation.Actions(), []string{identity.DelegatedActionAssembleContext}) {
		return slices.Equal(delegation.Resources(), []string{repositoryID}) &&
			len(delegation.Destinations()) == 0 &&
			delegation.MaximumAutonomy() == identity.ContextPreparationOnly
	}
	resources, destinations := delegation.Resources(), delegation.Destinations()
	return slices.Equal(delegation.Actions(), []string{identity.DelegatedActionProposeChange}) &&
		len(resources) == 1 && strings.HasPrefix(resources[0], "sha256:") &&
		len(destinations) == 1 && strings.HasPrefix(destinations[0], "model:") &&
		delegation.MaximumAutonomy() == identity.ChangeProposalOnly
}

func matches(rule Rule, request EvaluateRequest) bool {
	return containsAction(rule.Actions, request.Action) &&
		containsEnvironment(rule.Environments, request.Environment) &&
		containsRisk(rule.Risks, request.Risk) &&
		containsString(rule.SubjectIDs, request.Principal.SubjectID()) &&
		containsString(rule.ResourceIDs, request.Run.RepositoryID()) &&
		containsString(rule.DeliveryRunIDs, request.Run.ID()) &&
		containsString(rule.Purposes, request.Run.Purpose()) &&
		containsString(rule.Classifications, request.Run.Classification()) &&
		containsString(rule.TenantIDs, request.Scope.TenantID()) &&
		containsString(rule.WorkspaceIDs, request.Scope.WorkspaceID()) &&
		containsString(rule.ProjectIDs, request.Scope.ProjectID())
}

func matchesDelegation(rule Rule, request EvaluateDelegationRequest) bool {
	return containsAction(rule.Actions, ActionCreateDelegation) &&
		containsEnvironment(rule.Environments, request.Environment) &&
		containsRisk(rule.Risks, request.Risk) &&
		containsString(rule.SubjectIDs, request.Human.SubjectID()) &&
		containsString(rule.ResourceIDs, request.Delegation.Digest()) &&
		containsString(rule.DeliveryRunIDs, request.Run.ID()) &&
		containsString(rule.Purposes, request.Run.Purpose()) &&
		containsString(rule.Classifications, request.Run.Classification()) &&
		containsString(rule.TenantIDs, request.Scope.TenantID()) &&
		containsString(rule.WorkspaceIDs, request.Scope.WorkspaceID()) &&
		containsString(rule.ProjectIDs, request.Scope.ProjectID())
}

func validApprovalFacts(request EvaluateApprovalRequest, at time.Time) bool {
	if request.Environment != DeveloperTrialEnvironment || !validRisk(request.Risk) ||
		request.Principal.Kind() != identity.HumanPrincipal ||
		request.Run.ContractVersion() != delivery.DeliveryRunContractVersion ||
		request.Run.ActorID() != request.Principal.SubjectID() ||
		request.Run.TenantID() != request.Scope.TenantID() ||
		request.Run.WorkspaceID() != request.Scope.WorkspaceID() ||
		request.Run.ProjectID() != request.Scope.ProjectID() ||
		request.Request.IntendedApproverID != request.Principal.SubjectID() ||
		request.Request.Scope != (trust.ApprovalScope{
			TenantID: request.Run.TenantID(), WorkspaceID: request.Run.WorkspaceID(),
			ProjectID: request.Run.ProjectID(), RepositoryID: request.Run.RepositoryID(),
			DeliveryRunID: request.Run.ID(),
		}) || request.Request.Environment != string(request.Environment) ||
		request.Request.Purpose != request.Run.Purpose() ||
		request.Request.Classification != request.Run.Classification() ||
		request.Request.RiskClass != string(request.Risk) || !at.Before(request.Request.ExpiresAt) {
		return false
	}
	switch request.Request.Action {
	case "forge.delivery.execute-plan":
		plan := request.Run.Plan()
		return (request.Run.State() == delivery.PlannedState || request.Run.State() == delivery.AwaitingPlanApprovalState) &&
			plan != nil && request.Request.SubjectType == "delivery-plan" &&
			request.Request.SubjectID == plan.ID && request.Request.SubjectDigest == plan.Digest
	case "forge.source-control.publish-pr":
		candidate, verification := request.Run.Candidate(), request.Run.Verification()
		return request.Run.State() == delivery.AwaitingPublicationApprovalState &&
			candidate != nil && verification != nil && verification.Outcome == "passed" &&
			request.Request.SubjectType == "verified-candidate" &&
			request.Request.SubjectID == request.Run.ID()+".candidate" &&
			request.Request.SubjectDigest == candidate.Digest
	default:
		return false
	}
}

func matchesApproval(rule Rule, request EvaluateApprovalRequest) bool {
	return containsAction(rule.Actions, ActionDecideApproval) &&
		containsEnvironment(rule.Environments, request.Environment) &&
		containsRisk(rule.Risks, request.Risk) &&
		containsString(rule.SubjectIDs, request.Principal.SubjectID()) &&
		containsString(rule.ResourceIDs, request.Request.RequestDigest) &&
		containsString(rule.DeliveryRunIDs, request.Run.ID()) &&
		containsString(rule.Purposes, request.Run.Purpose()) &&
		containsString(rule.Classifications, request.Run.Classification()) &&
		containsString(rule.TenantIDs, request.Scope.TenantID()) &&
		containsString(rule.WorkspaceIDs, request.Scope.WorkspaceID()) &&
		containsString(rule.ProjectIDs, request.Scope.ProjectID())
}

func decisionID(request EvaluateRequest, decision trust.PolicyDecision, rules []Rule) string {
	ruleIDs := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleIDs = append(ruleIDs, rule.ID)
	}
	encoded, _ := json.Marshal(struct {
		Decision        trust.PolicyDecision
		DeliveryRunID   string
		TenantID        string
		WorkspaceID     string
		ProjectID       string
		Purpose         string
		Classification  string
		Risk            Risk
		Environment     Environment
		MatchingRuleIDs []string
	}{
		Decision: decision, DeliveryRunID: request.Run.ID(),
		TenantID: request.Scope.TenantID(), WorkspaceID: request.Scope.WorkspaceID(),
		ProjectID: request.Scope.ProjectID(), Purpose: request.Run.Purpose(),
		Classification: request.Run.Classification(), Risk: request.Risk,
		Environment: request.Environment, MatchingRuleIDs: ruleIDs,
	})
	digest := sha256.Sum256(encoded)
	return "forge.policy-decision." + hex.EncodeToString(digest[:])
}

func delegationDecisionID(request EvaluateDelegationRequest, decision trust.PolicyDecision, rules []Rule) string {
	ruleIDs := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleIDs = append(ruleIDs, rule.ID)
	}
	encoded, _ := json.Marshal(struct {
		Decision        trust.PolicyDecision
		AgentSubjectID  string
		DeliveryRunID   string
		TenantID        string
		WorkspaceID     string
		ProjectID       string
		Purpose         string
		Classification  string
		Risk            Risk
		Environment     Environment
		MatchingRuleIDs []string
	}{
		Decision: decision, AgentSubjectID: request.Agent.SubjectID(),
		DeliveryRunID: request.Run.ID(), TenantID: request.Scope.TenantID(),
		WorkspaceID: request.Scope.WorkspaceID(), ProjectID: request.Scope.ProjectID(),
		Purpose: request.Run.Purpose(), Classification: request.Run.Classification(),
		Risk: request.Risk, Environment: request.Environment, MatchingRuleIDs: ruleIDs,
	})
	digest := sha256.Sum256(encoded)
	return "forge.policy-decision." + hex.EncodeToString(digest[:])
}

func approvalDecisionID(request EvaluateApprovalRequest, decision trust.PolicyDecision, rules []Rule) string {
	ruleIDs := make([]string, 0, len(rules))
	for _, rule := range rules {
		ruleIDs = append(ruleIDs, rule.ID)
	}
	encoded, _ := json.Marshal(struct {
		Decision        trust.PolicyDecision
		RequestDigest   string
		DeliveryRunID   string
		TenantID        string
		WorkspaceID     string
		ProjectID       string
		Purpose         string
		Classification  string
		Risk            Risk
		Environment     Environment
		MatchingRuleIDs []string
	}{
		Decision: decision, RequestDigest: request.Request.RequestDigest,
		DeliveryRunID: request.Run.ID(), TenantID: request.Scope.TenantID(),
		WorkspaceID: request.Scope.WorkspaceID(), ProjectID: request.Scope.ProjectID(),
		Purpose: request.Run.Purpose(), Classification: request.Run.Classification(),
		Risk: request.Risk, Environment: request.Environment, MatchingRuleIDs: ruleIDs,
	})
	digest := sha256.Sum256(encoded)
	return "forge.policy-decision." + hex.EncodeToString(digest[:])
}

func cloneRule(rule Rule) Rule {
	obligations := rule.Obligations
	rule.Actions = slices.Clone(rule.Actions)
	rule.Environments = slices.Clone(rule.Environments)
	rule.Risks = slices.Clone(rule.Risks)
	rule.SubjectIDs = slices.Clone(rule.SubjectIDs)
	rule.ResourceIDs = slices.Clone(rule.ResourceIDs)
	rule.DeliveryRunIDs = slices.Clone(rule.DeliveryRunIDs)
	rule.Purposes = slices.Clone(rule.Purposes)
	rule.Classifications = slices.Clone(rule.Classifications)
	rule.TenantIDs = slices.Clone(rule.TenantIDs)
	rule.WorkspaceIDs = slices.Clone(rule.WorkspaceIDs)
	rule.ProjectIDs = slices.Clone(rule.ProjectIDs)
	rule.Obligations = make([]Obligation, len(obligations))
	for index, obligation := range obligations {
		rule.Obligations[index] = Obligation{Type: obligation.Type, Parameters: cloneMap(obligation.Parameters)}
	}
	return rule
}

func obligationMap(obligation Obligation) map[string]any {
	value := cloneMap(obligation.Parameters)
	if value == nil {
		value = make(map[string]any)
	}
	value["type"] = obligation.Type
	return value
}

func cloneMap(value map[string]any) map[string]any {
	if value == nil {
		return nil
	}
	encoded, _ := json.Marshal(value)
	var result map[string]any
	_ = json.Unmarshal(encoded, &result)
	return result
}

func allUniqueNonZero[T comparable](values []T) bool {
	seen := make(map[T]struct{}, len(values))
	var zero T
	for _, value := range values {
		if value == zero {
			return false
		}
		if _, duplicate := seen[value]; duplicate {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func containsAction(values []Action, value Action) bool {
	return len(values) == 0 || slices.Contains(values, value)
}

func containsEnvironment(values []Environment, value Environment) bool {
	return len(values) == 0 || slices.Contains(values, value)
}

func containsRisk(values []Risk, value Risk) bool {
	return len(values) == 0 || slices.Contains(values, value)
}

func containsString(values []string, value string) bool {
	return len(values) == 0 || slices.Contains(values, value)
}

func validRisk(value Risk) bool {
	return value == RiskLow || value == RiskModerate || value == RiskHigh || value == RiskCritical
}

func validClassification(value string) bool {
	return value == "Public" || value == "Internal" || value == "Confidential" || value == "Restricted"
}

func invalidRequest(correlationID string) error {
	return failure.New(failure.InvalidRequest, "policy evaluation request is invalid", correlationID, failure.RetryNone)
}

func policyUnavailable(correlationID string) error {
	return failure.New(failure.DependencyUnavailable, "policy evaluation is unavailable", correlationID, failure.RetryBounded)
}
