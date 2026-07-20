# Intelligence, model, and resource governance

Status: Proposed
Owner: Forge Intelligence Architecture

This design consolidates the existing intelligence model, model selection policy, LLM bootstrap, token stewardship, and compute stewardship documents.

## Principle

Intelligence is a governed capability, not a hard dependency on one model or provider. Deterministic software owns identity, policy, state transitions, budgets, and verification. Models propose transformations within those controls.

## Registry and routing

The model registry records provider, model version, modalities, context limits, deployment location, data handling, approved classifications, quality evaluations, latency, availability, energy indicators, and cost.

Routing evaluates, in order:

1. policy and data eligibility;
2. required capability and evaluated quality;
3. deployment and residency constraints;
4. latency and availability objective;
5. token, compute, financial, and sustainability budget.

A provider gateway MAY implement the routing contract, but no specific gateway is part of the domain contract without an accepted ADR.

## Resource stewardship

- Use the smallest sufficient model and context.
- Reuse verified artefacts and deterministic computation before generating again.
- Apply request, tenant, project, and organization budgets.
- Expose estimated and actual token, compute, latency, energy, and financial cost.
- Stop or request approval when a budget or risk threshold is reached.
- Never trade required safety or correctness for lower cost silently.

## Bootstrap and degradation

Forge SHALL be able to start without model connectivity for administration, inspection, policy, and recovery. Provider failures degrade affected intelligent functions rather than corrupting workflow state. Retry, fallback, and provider switching remain bounded by data policy and evidence requirements.
