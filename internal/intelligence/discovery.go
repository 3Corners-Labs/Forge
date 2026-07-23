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

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
)

const ProviderActivationVersion = "forge.model-provider-activation.v1"
const maxTrackedDiscoverySessions = 256

type ProviderActivation struct {
	SchemaVersion          string
	ID                     string
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
	Budgets                discovery.Limits
	MaxTurns               int
	ActivatedAt            time.Time
	ExpiresAt              time.Time
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
}

type trackedTurns struct {
	count     int
	expiresAt time.Time
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
		turns: make(map[string]trackedTurns),
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
	if !service.reserveTurn(request.SessionID, request.ExpiresAt, now) {
		return discovery.TurnProposal{}, failure.New(
			failure.BudgetExhausted, "model-assisted discovery turn budget is exhausted",
			request.CorrelationID, failure.RetryNone,
		)
	}
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
	return cloneProposal(proposal), nil
}

func (service *Service) reserveTurn(sessionID string, expiresAt, now time.Time) bool {
	service.turnMu.Lock()
	defer service.turnMu.Unlock()
	for id, tracked := range service.turns {
		if !tracked.expiresAt.After(now) {
			delete(service.turns, id)
		}
	}
	tracked, found := service.turns[sessionID]
	if !found && len(service.turns) >= maxTrackedDiscoverySessions {
		return false
	}
	if tracked.count >= service.activation.MaxTurns {
		return false
	}
	tracked.count++
	if expiresAt.After(tracked.expiresAt) {
		tracked.expiresAt = expiresAt
	}
	service.turns[sessionID] = tracked
	return true
}

func (value ProviderActivation) validate(now time.Time) error {
	if value.SchemaVersion != ProviderActivationVersion {
		return errors.New("provider activation version is unsupported")
	}
	if blank(value.ID, value.Provider, value.Endpoint, value.AdapterVersion, value.CredentialRef,
		value.HumanApprovalRef, value.PromptVersion, value.Retention, value.TrainingUse, value.Region) {
		return errors.New("provider activation is incomplete")
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
		if purpose != discovery.PurposeV1 {
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
	if value.Budgets.MaxInputBytes <= 0 || value.Budgets.MaxOutputBytes <= 0 ||
		value.Budgets.MaxTokens <= 0 || value.Budgets.MaxLatency <= 0 ||
		value.Budgets.MaxCostMicros < 0 || value.MaxTurns <= 0 || value.MaxTurns > 20 {
		return errors.New("provider activation budgets are invalid")
	}
	if value.ActivatedAt.IsZero() || value.ExpiresAt.IsZero() ||
		value.ExpiresAt.Before(value.ActivatedAt) || !value.ExpiresAt.After(now) {
		return errors.New("provider activation lifetime is invalid")
	}
	return nil
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
