package intelligence

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"sync"
	"time"

	changecontract "github.com/3corners/forge/internal/contracts/change"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
)

const ProposeChangeAction = "forge.intelligence.propose-change"

type ChangeAdapter interface {
	Provider() string
	Prepare(context.Context, ProviderActivation, changecontract.Request) (changecontract.Proposal, error)
}

// ChangeService admits structured model output as information only. It neither
// issues worker authority nor materializes candidate bytes.
type ChangeService struct {
	activation ProviderActivation
	adapter    ChangeAdapter
	now        func() time.Time
	mu         sync.Mutex
	turns      int
	usedTokens int64
	usedCost   int64
	inFlight   int
}

func NewChangeService(activation ProviderActivation, adapter ChangeAdapter, now func() time.Time) (*ChangeService, error) {
	if adapter == nil || now == nil {
		return nil, errors.New("intelligence change-preparation configuration is invalid")
	}
	if err := activation.validate(now()); err != nil {
		return nil, err
	}
	if adapter.Provider() != activation.Provider ||
		!slices.Equal(activation.AllowedPurposes, []string{changecontract.PurposeV1}) {
		return nil, errors.New("change-preparation adapter does not match its activation")
	}
	return &ChangeService{activation: cloneActivation(activation), adapter: adapter, now: now}, nil
}

type ChangePreparationCommand struct {
	Request  changecontract.Request
	Decision trust.PolicyDecision
	Envelope trust.ExecutionEnvelope
}

func (service *ChangeService) Prepare(ctx context.Context, command ChangePreparationCommand) (changecontract.Proposal, error) {
	request := command.Request
	if service == nil || service.adapter == nil || service.now == nil || ctx == nil {
		return changecontract.Proposal{}, changeFailure(
			failure.DependencyUnavailable, "model-assisted change preparation is unavailable",
			request.CorrelationID, failure.RetryBounded,
		)
	}
	now := service.now()
	if err := request.Validate(now); err != nil {
		return changecontract.Proposal{}, changeFailure(
			failure.InvalidRequest, "change-preparation request is invalid",
			request.CorrelationID, failure.RetryNone,
		)
	}
	if err := service.activation.validate(now); err != nil ||
		!contains(service.activation.AllowedPurposes, request.Purpose) ||
		!contains(service.activation.AllowedClassifications, request.Classification) {
		return changecontract.Proposal{}, changeFailure(
			failure.Forbidden, "model provider activation is not eligible for change preparation",
			request.CorrelationID, failure.RetryNewAuthority,
		)
	}
	if err := command.Decision.Authorize(
		"intelligence", request.TaskAgentID, ProposeChangeAction, request.Digest, now,
		[]string{"purpose-minimized-context", "model-output-non-authoritative"}, request.CorrelationID,
	); err != nil || command.Envelope.VerifyIntegrity(request.CorrelationID) != nil ||
		command.Envelope.Scope != (trust.ExecutionScope{
			TenantID: request.TenantID, WorkspaceID: request.WorkspaceID,
			ProjectID: request.ProjectID, DeliveryRunID: request.DeliveryRunID,
		}) ||
		command.Envelope.Principal.AgentID != request.TaskAgentID ||
		command.Envelope.Purpose != changecontract.PurposeV1 ||
		!slices.Equal(command.Envelope.Actions, []string{ProposeChangeAction}) ||
		!slices.Equal(command.Envelope.PolicyDecisionIDs, []string{command.Decision.ID}) ||
		!slices.Equal(command.Envelope.Destinations, []string{"model:" + service.activation.Provider}) ||
		!slices.Equal(command.Envelope.ModelIDs, service.activation.AllowedModelIDs) ||
		!slices.Equal(command.Envelope.SecretReferenceIDs, []string{service.activation.CredentialRef}) ||
		len(command.Envelope.ToolIDs) != 0 ||
		command.Envelope.Limits.MaxConcurrency != 1 ||
		command.Envelope.Limits.MaxInputTokens < request.Limits.MaxTokens ||
		command.Envelope.Limits.MaxOutputTokens < request.Limits.MaxTokens ||
		command.Envelope.Limits.MaxCostMinorUnits < request.Limits.MaxCostMicros ||
		!now.Before(command.Envelope.ExpiresAt) {
		return changecontract.Proposal{}, changeFailure(
			failure.Forbidden, "model-proposal authority does not match the exact request",
			request.CorrelationID, failure.RetryNewAuthority,
		)
	}
	if request.Limits.MaxInputBytes > service.activation.Budgets.MaxInputBytes ||
		request.Limits.MaxOutputBytes > service.activation.Budgets.MaxOutputBytes ||
		request.Limits.MaxTokens > service.activation.Budgets.MaxTokens ||
		request.Limits.MaxLatency > service.activation.Budgets.MaxLatency ||
		request.Limits.MaxCostMicros > service.activation.Budgets.MaxCostMicros {
		return changecontract.Proposal{}, changeFailure(
			failure.BudgetExhausted, "change-preparation request exceeds provider activation budgets",
			request.CorrelationID, failure.RetryNewAuthority,
		)
	}
	encoded, err := json.Marshal(request)
	if err != nil || int64(len(encoded)) > request.Limits.MaxInputBytes {
		return changecontract.Proposal{}, changeFailure(
			failure.BudgetExhausted, "change-preparation input exceeds its budget",
			request.CorrelationID, failure.RetryNone,
		)
	}
	if !service.reserve(request) {
		return changecontract.Proposal{}, changeFailure(
			failure.BudgetExhausted, "change-preparation activation budget is exhausted",
			request.CorrelationID, failure.RetryNone,
		)
	}
	defer service.finish()

	deadline := now.Add(request.Limits.MaxLatency)
	if request.ExpiresAt.Before(deadline) {
		deadline = request.ExpiresAt
	}
	bounded, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	proposal, err := service.adapter.Prepare(bounded, cloneActivation(service.activation), request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(bounded.Err(), context.DeadlineExceeded) {
			return changecontract.Proposal{}, changeFailure(
				failure.DependencyUnavailable, "model-assisted change preparation timed out",
				request.CorrelationID, failure.RetryBounded,
			)
		}
		return changecontract.Proposal{}, changeFailure(
			failure.DependencyUnavailable, "model-assisted change provider is unavailable",
			request.CorrelationID, failure.RetryBounded,
		)
	}
	if err := proposal.Validate(request); err != nil {
		return changecontract.Proposal{}, changeFailure(
			failure.Indeterminate, "model-assisted change provider returned an invalid proposal",
			request.CorrelationID, failure.RetryBounded,
		)
	}
	if proposal.Receipt.Provider != service.activation.Provider ||
		proposal.Receipt.ActivationID != service.activation.ID ||
		proposal.Receipt.ActivationDigest != service.activation.Digest ||
		!contains(service.activation.AllowedModelIDs, proposal.Receipt.ModelID) ||
		proposal.Receipt.AdapterVersion != service.activation.AdapterVersion ||
		proposal.Receipt.PromptVersion != service.activation.PromptVersion {
		return changecontract.Proposal{}, changeFailure(
			failure.Indeterminate, "change proposal receipt is not activation-bound",
			request.CorrelationID, failure.RetryNone,
		)
	}
	encoded, err = json.Marshal(proposal)
	if err != nil || int64(len(encoded)) > request.Limits.MaxOutputBytes {
		return changecontract.Proposal{}, changeFailure(
			failure.BudgetExhausted, "change proposal exceeds its output budget",
			request.CorrelationID, failure.RetryNone,
		)
	}
	service.commit(request, proposal.Usage)
	return cloneChangeProposal(proposal), nil
}

func (service *ChangeService) reserve(request changecontract.Request) bool {
	service.mu.Lock()
	defer service.mu.Unlock()
	if service.turns >= service.activation.MaxTurns ||
		service.inFlight >= service.activation.MaxConcurrentRequests ||
		service.usedTokens+request.Limits.MaxTokens > service.activation.MaxSessionTokens ||
		service.usedCost+request.Limits.MaxCostMicros > service.activation.MaxSessionCostMicros {
		return false
	}
	service.turns++
	service.inFlight++
	service.usedTokens += request.Limits.MaxTokens
	service.usedCost += request.Limits.MaxCostMicros
	return true
}

func (service *ChangeService) finish() {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.inFlight--
}

func (service *ChangeService) commit(request changecontract.Request, usage changecontract.Usage) {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.usedTokens -= request.Limits.MaxTokens
	service.usedCost -= request.Limits.MaxCostMicros
	service.usedTokens += usage.InputTokens + usage.OutputTokens
	service.usedCost += usage.CostMicros
}

func cloneChangeProposal(value changecontract.Proposal) changecontract.Proposal {
	value.Replacements = append([]changecontract.FileReplacement(nil), value.Replacements...)
	value.ProposedChecks = append([]string(nil), value.ProposedChecks...)
	value.Assumptions = append([]string(nil), value.Assumptions...)
	value.Uncertainty = append([]string(nil), value.Uncertainty...)
	value.Omissions = append([]string(nil), value.Omissions...)
	return value
}

func changeFailure(code failure.Code, message, correlationID string, retry failure.Retry) error {
	return failure.New(code, message, correlationID, retry)
}
