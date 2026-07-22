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
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/tenant"
)

const (
	EvaluateRequestContractVersion = "forge.policy-evaluate-request.v1"
	DeveloperTrialEnvironment      = Environment("developer-trial")
	ActionAssembleContext          = Action("forge.context.assemble")
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

type Evaluator struct {
	set                  PolicySet
	supportedObligations map[string]struct{}
	decisionTTL          time.Duration
	now                  func() time.Time
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
	return decision, nil
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
		if action != ActionAssembleContext {
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
