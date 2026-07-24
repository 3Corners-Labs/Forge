// Package discovery defines the provider-neutral product-discovery contract.
// Its values are information only and never carry authority.
package discovery

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	RequestVersion  = "forge.discovery-turn-request.v1"
	ProposalVersion = "forge.discovery-turn-proposal.v1"
	RoleProfileV1   = "forge.role.product-manager.discovery.v1"
	PurposeV1       = "forge.product.discovery"
)

const (
	maxStatements = 32
	maxItems      = 64
	maxTextBytes  = 16 << 10
)

type Limits struct {
	MaxInputBytes  int64         `json:"maxInputBytes"`
	MaxOutputBytes int64         `json:"maxOutputBytes"`
	MaxTokens      int64         `json:"maxTokens"`
	MaxLatency     time.Duration `json:"maxLatency"`
	MaxCostMicros  int64         `json:"maxCostMicros"`
}

type Draft struct {
	Outcome            string   `json:"outcome"`
	AffectedUsers      []string `json:"affectedUsers"`
	Problem            string   `json:"problem"`
	FunctionalBehavior []string `json:"functionalBehavior"`
	UXExpectations     []string `json:"uxExpectations"`
	Accessibility      []string `json:"accessibility"`
	Privacy            []string `json:"privacy"`
	Security           []string `json:"security"`
	Operations         []string `json:"operations"`
	Constraints        []string `json:"constraints"`
	Exclusions         []string `json:"exclusions"`
	AcceptanceCriteria []string `json:"acceptanceCriteria"`
	Inferences         []string `json:"inferences"`
	Assumptions        []string `json:"assumptions"`
	Hypotheses         []string `json:"hypotheses"`
	Experiments        []string `json:"experiments"`
	ProductAreaRefs    []string `json:"productAreaRefs"`
}

type Question struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Consequence string `json:"consequence"`
	Material    bool   `json:"material"`
}

type TurnRequest struct {
	SchemaVersion     string     `json:"schemaVersion"`
	SessionID         string     `json:"sessionId"`
	CorrelationID     string     `json:"correlationId"`
	TenantID          string     `json:"tenantId"`
	WorkspaceID       string     `json:"workspaceId"`
	ProjectID         string     `json:"projectId"`
	HumanSubjectID    string     `json:"humanSubjectId"`
	TaskAgentID       string     `json:"taskAgentId"`
	RoleProfileID     string     `json:"roleProfileId"`
	Purpose           string     `json:"purpose"`
	Classification    string     `json:"classification"`
	UserStatements    []string   `json:"userStatements"`
	CurrentDraft      Draft      `json:"currentDraft"`
	PreviousQuestions []Question `json:"previousQuestions"`
	ProvenanceRefs    []string   `json:"provenanceRefs"`
	Limits            Limits     `json:"limits"`
	ExpiresAt         time.Time  `json:"expiresAt"`
}

type Usage struct {
	InputTokens  int64 `json:"inputTokens"`
	OutputTokens int64 `json:"outputTokens"`
	LatencyMs    int64 `json:"latencyMs"`
	CostMicros   int64 `json:"costMicros"`
}

type Receipt struct {
	ActivationID     string `json:"activationId"`
	ActivationDigest string `json:"activationDigest"`
	Provider         string `json:"provider"`
	ModelID          string `json:"modelId"`
	AdapterVersion   string `json:"adapterVersion"`
	RoleProfileID    string `json:"roleProfileId"`
	PromptVersion    string `json:"promptVersion"`
	ResponseSchema   string `json:"responseSchema"`
	ProviderRequest  string `json:"providerRequest"`
}

type TurnProposal struct {
	SchemaVersion  string     `json:"schemaVersion"`
	SessionID      string     `json:"sessionId"`
	CorrelationID  string     `json:"correlationId"`
	Response       string     `json:"response"`
	Draft          Draft      `json:"draft"`
	Questions      []Question `json:"questions"`
	Uncertainty    []string   `json:"uncertainty"`
	Limitations    []string   `json:"limitations"`
	ProvenanceRefs []string   `json:"provenanceRefs"`
	Refusal        bool       `json:"refusal"`
	Incomplete     bool       `json:"incomplete"`
	Usage          Usage      `json:"usage"`
	Receipt        Receipt    `json:"receipt"`
}

func (value TurnRequest) Validate(now time.Time) error {
	if value.SchemaVersion != RequestVersion {
		return fmt.Errorf("unsupported discovery request version %q", value.SchemaVersion)
	}
	if blank(value.SessionID, value.CorrelationID, value.TenantID, value.WorkspaceID, value.ProjectID,
		value.HumanSubjectID, value.TaskAgentID, value.RoleProfileID, value.Purpose, value.Classification) {
		return errors.New("discovery request identity, scope, role, purpose, and classification are required")
	}
	if value.RoleProfileID != RoleProfileV1 || value.Purpose != PurposeV1 {
		return errors.New("unsupported discovery role or purpose")
	}
	if len(value.UserStatements) == 0 || len(value.UserStatements) > maxStatements ||
		!validTexts(value.UserStatements, maxTextBytes) {
		return errors.New("discovery user statements are invalid")
	}
	if len(value.ProvenanceRefs) == 0 || len(value.ProvenanceRefs) > maxItems ||
		!validRefs(value.ProvenanceRefs) {
		return errors.New("discovery provenance references are invalid")
	}
	if !validQuestions(value.PreviousQuestions) || !value.CurrentDraft.valid(false) {
		return errors.New("discovery request draft or prior questions are invalid")
	}
	if value.Limits.MaxInputBytes <= 0 || value.Limits.MaxOutputBytes <= 0 ||
		value.Limits.MaxTokens <= 0 || value.Limits.MaxLatency <= 0 ||
		value.Limits.MaxCostMicros < 0 {
		return errors.New("discovery request limits are invalid")
	}
	if value.ExpiresAt.IsZero() || !value.ExpiresAt.After(now) {
		return errors.New("discovery request is expired")
	}
	return nil
}

func (value TurnProposal) Validate(request TurnRequest) error {
	if value.SchemaVersion != ProposalVersion {
		return fmt.Errorf("unsupported discovery proposal version %q", value.SchemaVersion)
	}
	if value.SessionID != request.SessionID || value.CorrelationID != request.CorrelationID {
		return errors.New("discovery proposal scope does not match the request")
	}
	if strings.TrimSpace(value.Response) == "" || len(value.Response) > maxTextBytes {
		return errors.New("discovery proposal response is invalid")
	}
	if !value.Draft.valid(!value.Refusal && !value.Incomplete && noMaterialQuestions(value.Questions)) ||
		!validQuestions(value.Questions) || !validTexts(value.Uncertainty, maxTextBytes) ||
		!validTexts(value.Limitations, maxTextBytes) || len(value.ProvenanceRefs) == 0 ||
		!validRefs(value.ProvenanceRefs) {
		return errors.New("discovery proposal content is invalid")
	}
	if !subset(value.ProvenanceRefs, request.ProvenanceRefs) {
		return errors.New("discovery proposal invented provenance")
	}
	if blank(value.Receipt.ActivationID, value.Receipt.ActivationDigest,
		value.Receipt.Provider, value.Receipt.ModelID, value.Receipt.AdapterVersion,
		value.Receipt.RoleProfileID, value.Receipt.PromptVersion, value.Receipt.ResponseSchema,
		value.Receipt.ProviderRequest) {
		return errors.New("discovery proposal receipt is incomplete")
	}
	if !validDigest(value.Receipt.ActivationDigest) {
		return errors.New("discovery proposal activation digest is invalid")
	}
	if value.Receipt.RoleProfileID != request.RoleProfileID || value.Receipt.ResponseSchema != ProposalVersion {
		return errors.New("discovery proposal receipt does not match the request")
	}
	if value.Usage.InputTokens < 0 || value.Usage.OutputTokens < 0 || value.Usage.LatencyMs < 0 ||
		value.Usage.CostMicros < 0 || value.Usage.OutputTokens > request.Limits.MaxTokens ||
		value.Usage.InputTokens > request.Limits.MaxTokens-value.Usage.OutputTokens ||
		value.Usage.CostMicros > request.Limits.MaxCostMicros ||
		time.Duration(value.Usage.LatencyMs)*time.Millisecond > request.Limits.MaxLatency {
		return errors.New("discovery proposal exceeded its resource limits")
	}
	return nil
}

func (value Draft) Complete() bool {
	return value.valid(true)
}

func (value Draft) Clone() Draft {
	value.AffectedUsers = slices.Clone(value.AffectedUsers)
	value.FunctionalBehavior = slices.Clone(value.FunctionalBehavior)
	value.UXExpectations = slices.Clone(value.UXExpectations)
	value.Accessibility = slices.Clone(value.Accessibility)
	value.Privacy = slices.Clone(value.Privacy)
	value.Security = slices.Clone(value.Security)
	value.Operations = slices.Clone(value.Operations)
	value.Constraints = slices.Clone(value.Constraints)
	value.Exclusions = slices.Clone(value.Exclusions)
	value.AcceptanceCriteria = slices.Clone(value.AcceptanceCriteria)
	value.Inferences = slices.Clone(value.Inferences)
	value.Assumptions = slices.Clone(value.Assumptions)
	value.Hypotheses = slices.Clone(value.Hypotheses)
	value.Experiments = slices.Clone(value.Experiments)
	value.ProductAreaRefs = slices.Clone(value.ProductAreaRefs)
	return value
}

func (value Draft) valid(requireComplete bool) bool {
	collections := [][]string{
		value.AffectedUsers, value.FunctionalBehavior, value.UXExpectations, value.Accessibility,
		value.Privacy, value.Security, value.Operations, value.Constraints, value.Exclusions,
		value.AcceptanceCriteria, value.Inferences, value.Assumptions, value.Hypotheses, value.Experiments,
		value.ProductAreaRefs,
	}
	for _, items := range collections {
		if len(items) > maxItems || !validTexts(items, maxTextBytes) {
			return false
		}
	}
	if len(value.Outcome) > maxTextBytes || len(value.Problem) > maxTextBytes {
		return false
	}
	return !requireComplete || (strings.TrimSpace(value.Outcome) != "" &&
		strings.TrimSpace(value.Problem) != "" && len(value.AffectedUsers) > 0 &&
		len(value.AcceptanceCriteria) > 0 && value.Constraints != nil)
}

func validQuestions(values []Question) bool {
	if len(values) > maxItems {
		return false
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if blank(value.ID, value.Text) || len(value.ID)+len(value.Text)+len(value.Consequence) > maxTextBytes {
			return false
		}
		if value.Material && strings.TrimSpace(value.Consequence) == "" {
			return false
		}
		if _, exists := seen[value.ID]; exists {
			return false
		}
		seen[value.ID] = struct{}{}
	}
	return true
}

func noMaterialQuestions(values []Question) bool {
	for _, value := range values {
		if value.Material {
			return false
		}
	}
	return true
}

func validTexts(values []string, totalLimit int) bool {
	total := 0
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return false
		}
		total += len(value)
	}
	return total <= totalLimit
}

func validRefs(values []string) bool {
	return validTexts(values, maxTextBytes) && len(values) == len(unique(values))
}

func unique(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func subset(values, allowed []string) bool {
	index := make(map[string]struct{}, len(allowed))
	for _, value := range allowed {
		index[value] = struct{}{}
	}
	for _, value := range values {
		if _, exists := index[value]; !exists {
			return false
		}
	}
	return true
}

func blank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != len("sha256:")+64 {
		return false
	}
	for _, character := range strings.TrimPrefix(value, "sha256:") {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}
