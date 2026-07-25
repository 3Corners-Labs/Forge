// Package intelligence owns bounded model eligibility and adapter invocation.
// Models propose information; this package never grants authority.
package intelligence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	changecontract "github.com/3corners/forge/internal/contracts/change"
	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
)

const ProviderActivationVersion = "forge.model-provider-activation.v1"
const maxTrackedDiscoverySessions = 256
const maxTrackedDiscoveryProjects = 256

type ProviderActivation struct {
	SchemaVersion          string
	ID                     string
	Digest                 string
	Provider               string
	Endpoint               string
	AdapterVersion         string
	AllowedModelIDs        []string
	AllowedPurposes        []string
	AllowedClassifications []string
	CredentialRef          string
	HumanApprovalRef       string
	PromptVersion          string
	Retention              string
	TrainingUse            string
	Region                 string
	StoreProviderState     bool
	ReasoningEffort        string
	Pricing                Pricing
	Budgets                discovery.Limits
	MaxTurns               int
	MaxSessionTokens       int64
	MaxSessionCostMicros   int64
	MaxProjectTokensPerDay int64
	MaxProjectCostPerDay   int64
	MaxConcurrentRequests  int
	ActivatedAt            time.Time
	ExpiresAt              time.Time
}

type Pricing struct {
	Currency                     string
	InputMicrosPerMillionTokens  int64
	OutputMicrosPerMillionTokens int64
}

type Adapter interface {
	Provider() string
	Propose(context.Context, ProviderActivation, discovery.TurnRequest) (discovery.TurnProposal, error)
}

type Service struct {
	activation ProviderActivation
	adapter    Adapter
	now        func() time.Time
	turnMu     sync.Mutex
	turns      map[string]trackedTurns
	projects   map[string]trackedUsage
	inFlight   int
}

type trackedTurns struct {
	count int
	trackedUsage
}

type trackedUsage struct {
	usedTokens     int64
	usedCost       int64
	reservedTokens int64
	reservedCost   int64
	expiresAt      time.Time
}

type reservation struct {
	sessionID  string
	projectKey string
	tokens     int64
	cost       int64
}

func NewService(activation ProviderActivation, adapter Adapter, now func() time.Time) (*Service, error) {
	if adapter == nil || now == nil {
		return nil, errors.New("intelligence discovery configuration is invalid")
	}
	if err := activation.validate(now()); err != nil {
		return nil, err
	}
	if adapter.Provider() != activation.Provider {
		return nil, errors.New("intelligence adapter does not match the provider activation")
	}
	return &Service{
		activation: cloneActivation(activation), adapter: adapter, now: now,
		turns: make(map[string]trackedTurns), projects: make(map[string]trackedUsage),
	}, nil
}

func (service *Service) Propose(ctx context.Context, request discovery.TurnRequest) (discovery.TurnProposal, error) {
	if service == nil || service.adapter == nil || service.now == nil || ctx == nil {
		return discovery.TurnProposal{}, failure.New(
			failure.DependencyUnavailable, "model-assisted discovery is unavailable",
			request.CorrelationID, failure.RetryBounded,
		)
	}
	now := service.now()
	if err := request.Validate(now); err != nil {
		return discovery.TurnProposal{}, failure.New(
			failure.InvalidRequest, "model-assisted discovery request is invalid",
			request.CorrelationID, failure.RetryNone,
		)
	}
	if err := service.activation.validate(now); err != nil {
		return discovery.TurnProposal{}, failure.New(
			failure.Forbidden, "model provider activation is not eligible",
			request.CorrelationID, failure.RetryNewAuthority,
		)
	}
	if !contains(service.activation.AllowedPurposes, request.Purpose) ||
		!contains(service.activation.AllowedClassifications, request.Classification) {
		return discovery.TurnProposal{}, failure.New(
			failure.Forbidden, "model provider is not eligible for this purpose or classification",
			request.CorrelationID, failure.RetryNewAuthority,
		)
	}
	if request.Limits.MaxInputBytes > service.activation.Budgets.MaxInputBytes ||
		request.Limits.MaxOutputBytes > service.activation.Budgets.MaxOutputBytes ||
		request.Limits.MaxTokens > service.activation.Budgets.MaxTokens ||
		request.Limits.MaxLatency > service.activation.Budgets.MaxLatency ||
		request.Limits.MaxCostMicros > service.activation.Budgets.MaxCostMicros {
		return discovery.TurnProposal{}, failure.New(
			failure.BudgetExhausted, "model-assisted discovery request exceeds provider activation budgets",
			request.CorrelationID, failure.RetryNewAuthority,
		)
	}
	reserved, ok := service.reserve(request, now)
	if !ok {
		return discovery.TurnProposal{}, failure.New(
			failure.BudgetExhausted, "model-assisted discovery session or project budget is exhausted",
			request.CorrelationID, failure.RetryNone,
		)
	}
	release := true
	defer func() {
		if release {
			service.release(reserved)
		}
	}()
	encoded, err := json.Marshal(request)
	if err != nil || int64(len(encoded)) > request.Limits.MaxInputBytes {
		return discovery.TurnProposal{}, failure.New(
			failure.BudgetExhausted, "model-assisted discovery input exceeds its budget",
			request.CorrelationID, failure.RetryNone,
		)
	}
	deadline := now.Add(request.Limits.MaxLatency)
	if request.ExpiresAt.Before(deadline) {
		deadline = request.ExpiresAt
	}
	bounded, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	proposal, err := service.adapter.Propose(bounded, cloneActivation(service.activation), request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(bounded.Err(), context.DeadlineExceeded) {
			return discovery.TurnProposal{}, failure.New(
				failure.DependencyUnavailable, "model-assisted discovery timed out",
				request.CorrelationID, failure.RetryBounded,
			)
		}
		return discovery.TurnProposal{}, failure.New(
			failure.DependencyUnavailable, "model-assisted discovery provider is unavailable",
			request.CorrelationID, failure.RetryBounded,
		)
	}
	if err := proposal.Validate(request); err != nil {
		return discovery.TurnProposal{}, failure.New(
			failure.Indeterminate, "model-assisted discovery returned an invalid proposal",
			request.CorrelationID, failure.RetryBounded,
		)
	}
	if proposal.Receipt.Provider != service.activation.Provider ||
		proposal.Receipt.ActivationID != service.activation.ID ||
		proposal.Receipt.ActivationDigest != service.activation.Digest ||
		!contains(service.activation.AllowedModelIDs, proposal.Receipt.ModelID) ||
		proposal.Receipt.AdapterVersion != service.activation.AdapterVersion ||
		proposal.Receipt.PromptVersion != service.activation.PromptVersion {
		return discovery.TurnProposal{}, failure.New(
			failure.Indeterminate, "model-assisted discovery receipt is not activation-bound",
			request.CorrelationID, failure.RetryNone,
		)
	}
	encoded, err = json.Marshal(proposal)
	if err != nil || int64(len(encoded)) > request.Limits.MaxOutputBytes {
		return discovery.TurnProposal{}, failure.New(
			failure.BudgetExhausted, "model-assisted discovery output exceeds its budget",
			request.CorrelationID, failure.RetryNone,
		)
	}
	service.commit(reserved, proposal.Usage)
	release = false
	return cloneProposal(proposal), nil
}

func (service *Service) reserve(request discovery.TurnRequest, now time.Time) (reservation, bool) {
	service.turnMu.Lock()
	defer service.turnMu.Unlock()
	for id, tracked := range service.turns {
		if !tracked.expiresAt.After(now) {
			delete(service.turns, id)
		}
	}
	today := now.UTC().Format("2006-01-02")
	for key := range service.projects {
		if !strings.HasSuffix(key, ":"+today) {
			delete(service.projects, key)
		}
	}
	tracked, found := service.turns[request.SessionID]
	if !found && len(service.turns) >= maxTrackedDiscoverySessions {
		return reservation{}, false
	}
	projectKey := request.ProjectID + ":" + today
	project, projectFound := service.projects[projectKey]
	if !projectFound && len(service.projects) >= maxTrackedDiscoveryProjects {
		return reservation{}, false
	}
	if tracked.count >= service.activation.MaxTurns ||
		exceedsLimit(tracked.usedTokens+tracked.reservedTokens, request.Limits.MaxTokens, service.activation.MaxSessionTokens) ||
		exceedsLimit(tracked.usedCost+tracked.reservedCost, request.Limits.MaxCostMicros, service.activation.MaxSessionCostMicros) ||
		exceedsLimit(project.usedTokens+project.reservedTokens, request.Limits.MaxTokens, service.activation.MaxProjectTokensPerDay) ||
		exceedsLimit(project.usedCost+project.reservedCost, request.Limits.MaxCostMicros, service.activation.MaxProjectCostPerDay) ||
		service.inFlight >= service.activation.MaxConcurrentRequests {
		return reservation{}, false
	}
	tracked.count++
	tracked.reservedTokens += request.Limits.MaxTokens
	tracked.reservedCost += request.Limits.MaxCostMicros
	tracked.expiresAt = service.activation.ExpiresAt
	project.reservedTokens += request.Limits.MaxTokens
	project.reservedCost += request.Limits.MaxCostMicros
	project.expiresAt = endOfUTCDate(now)
	service.turns[request.SessionID] = tracked
	service.projects[projectKey] = project
	service.inFlight++
	return reservation{
		sessionID: request.SessionID, projectKey: projectKey,
		tokens: request.Limits.MaxTokens, cost: request.Limits.MaxCostMicros,
	}, true
}

func (service *Service) release(value reservation) {
	service.turnMu.Lock()
	defer service.turnMu.Unlock()
	service.finishReservation(value, discovery.Usage{})
}

func (service *Service) commit(value reservation, usage discovery.Usage) {
	service.turnMu.Lock()
	defer service.turnMu.Unlock()
	service.finishReservation(value, usage)
}

func (service *Service) finishReservation(value reservation, usage discovery.Usage) {
	tracked := service.turns[value.sessionID]
	tracked.reservedTokens -= value.tokens
	tracked.reservedCost -= value.cost
	tracked.usedTokens += usage.InputTokens + usage.OutputTokens
	tracked.usedCost += usage.CostMicros
	service.turns[value.sessionID] = tracked
	project := service.projects[value.projectKey]
	project.reservedTokens -= value.tokens
	project.reservedCost -= value.cost
	project.usedTokens += usage.InputTokens + usage.OutputTokens
	project.usedCost += usage.CostMicros
	service.projects[value.projectKey] = project
	service.inFlight--
}

func exceedsLimit(current, requested, maximum int64) bool {
	return current < 0 || requested < 0 || maximum < 0 ||
		current > maximum || requested > maximum-current
}

func endOfUTCDate(now time.Time) time.Time {
	year, month, day := now.UTC().Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC)
}

func (value ProviderActivation) validate(now time.Time) error {
	if value.SchemaVersion != ProviderActivationVersion {
		return errors.New("provider activation version is unsupported")
	}
	if blank(value.ID, value.Digest, value.Provider, value.Endpoint, value.AdapterVersion, value.CredentialRef,
		value.HumanApprovalRef, value.PromptVersion, value.Retention, value.TrainingUse, value.Region) {
		return errors.New("provider activation is incomplete")
	}
	if !validSHA256(value.Digest) {
		return errors.New("provider activation digest is invalid")
	}
	if value.Provider != "openai" && value.Provider != "contract-double" {
		return fmt.Errorf("unsupported discovery provider %q", value.Provider)
	}
	if len(value.AllowedModelIDs) == 0 || !allNonEmpty(value.AllowedModelIDs) ||
		len(value.AllowedPurposes) == 0 || !allNonEmpty(value.AllowedPurposes) ||
		len(value.AllowedClassifications) == 0 || !allNonEmpty(value.AllowedClassifications) {
		return errors.New("provider activation allow-lists are invalid")
	}
	if !allUnique(value.AllowedModelIDs) || !allUnique(value.AllowedPurposes) ||
		!allUnique(value.AllowedClassifications) {
		return errors.New("provider activation allow-lists contain duplicates")
	}
	for _, purpose := range value.AllowedPurposes {
		if purpose != discovery.PurposeV1 && purpose != changecontract.PurposeV1 {
			return errors.New("provider activation contains an unsupported purpose")
		}
	}
	for _, classification := range value.AllowedClassifications {
		if classification != "Public" && classification != "Internal" {
			return errors.New("provider activation contains an unsupported classification")
		}
	}
	if value.StoreProviderState || value.TrainingUse != "disabled" {
		return errors.New("provider storage and training use must be disabled")
	}
	if !contains([]string{"none", "low", "medium", "high", "xhigh", "max"}, value.ReasoningEffort) {
		return errors.New("provider activation reasoning effort is invalid")
	}
	if value.Pricing.Currency != "USD" || value.Pricing.InputMicrosPerMillionTokens < 0 ||
		value.Pricing.OutputMicrosPerMillionTokens < 0 {
		return errors.New("provider activation pricing is invalid")
	}
	if value.Provider == "openai" &&
		(value.Pricing.InputMicrosPerMillionTokens == 0 ||
			value.Pricing.OutputMicrosPerMillionTokens == 0) {
		return errors.New("live provider activation pricing is required")
	}
	if value.Provider == "openai" &&
		(!validOpenAIAdapterPurpose(value) ||
			len(value.AllowedModelIDs) != 1 ||
			!validEnvironmentCredentialReference(value.CredentialRef) ||
			!allowedOpenAIEndpoint(value.Endpoint, value.Region)) {
		return errors.New("OpenAI provider activation is invalid")
	}
	if value.Budgets.MaxInputBytes <= 0 || value.Budgets.MaxOutputBytes <= 0 ||
		value.Budgets.MaxTokens <= 0 || value.Budgets.MaxLatency <= 0 ||
		value.Budgets.MaxCostMicros < 0 || value.MaxTurns <= 0 || value.MaxTurns > 20 ||
		value.MaxSessionTokens < value.Budgets.MaxTokens ||
		value.MaxSessionCostMicros < value.Budgets.MaxCostMicros ||
		value.MaxProjectTokensPerDay < value.MaxSessionTokens ||
		value.MaxProjectCostPerDay < value.MaxSessionCostMicros ||
		value.MaxConcurrentRequests <= 0 || value.MaxConcurrentRequests > 16 {
		return errors.New("provider activation budgets are invalid")
	}
	if value.ActivatedAt.IsZero() || value.ExpiresAt.IsZero() ||
		value.ExpiresAt.Before(value.ActivatedAt) || !value.ExpiresAt.After(now) {
		return errors.New("provider activation lifetime is invalid")
	}
	return nil
}

func validOpenAIAdapterPurpose(value ProviderActivation) bool {
	if slices.Equal(value.AllowedPurposes, []string{discovery.PurposeV1}) {
		return value.AdapterVersion == openAIAdapterVersion && value.PromptVersion == discoveryPromptVersion
	}
	if slices.Equal(value.AllowedPurposes, []string{changecontract.PurposeV1}) {
		return value.AdapterVersion == openAIChangeAdapterVersion && value.PromptVersion == changePromptVersion
	}
	return false
}

func validSHA256(value string) bool {
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

func validEnvironmentCredentialReference(value string) bool {
	const prefix = "env:"
	if !strings.HasPrefix(value, prefix) {
		return false
	}
	name := strings.TrimPrefix(value, prefix)
	if name == "" || !((name[0] >= 'A' && name[0] <= 'Z') ||
		(name[0] >= 'a' && name[0] <= 'z') || name[0] == '_') {
		return false
	}
	for _, character := range name[1:] {
		if !((character >= 'A' && character <= 'Z') ||
			(character >= 'a' && character <= 'z') ||
			(character >= '0' && character <= '9') || character == '_') {
			return false
		}
	}
	return true
}

func cloneActivation(value ProviderActivation) ProviderActivation {
	value.AllowedModelIDs = slices.Clone(value.AllowedModelIDs)
	value.AllowedPurposes = slices.Clone(value.AllowedPurposes)
	value.AllowedClassifications = slices.Clone(value.AllowedClassifications)
	return value
}

func cloneProposal(value discovery.TurnProposal) discovery.TurnProposal {
	value.Draft = value.Draft.Clone()
	value.Questions = slices.Clone(value.Questions)
	value.Uncertainty = slices.Clone(value.Uncertainty)
	value.Limitations = slices.Clone(value.Limitations)
	value.ProvenanceRefs = slices.Clone(value.ProvenanceRefs)
	return value
}

func contains(values []string, expected string) bool {
	return slices.Contains(values, expected)
}

func blank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func allNonEmpty(values []string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return false
		}
	}
	return true
}

func allUnique(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, found := seen[value]; found {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}
