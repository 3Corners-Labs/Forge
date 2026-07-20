# Policy and authorization architecture

Status: Proposed
Owner: Forge Security and Governance Architecture
Version: 0.1

## Purpose

Policy transforms organizational intent into deterministic decisions. Authorization enforcement ensures that each receiving capability applies those decisions without expanding them.

## Logical responsibilities

The Policy Decision capability:

- resolves applicable versioned policies;
- evaluates verified subject, action, resource, scope, classification, risk and environment facts;
- produces permit, deny or indeterminate plus obligations and reason codes;
- detects and resolves conflicts deterministically;
- records evaluator and policy versions.

Policy Enforcement Points exist at API, knowledge retrieval, secret resolution, context assembly, model and tool selection, worker execution, evidence disclosure, publication and release boundaries.

The Runtime enforces obligations. It does not reinterpret policy or treat provider success as authorization.

## Evaluation sequence

1. Verify principal and immutable tenant/workspace/project scope.
2. Identify the action gate and requested resource.
3. Resolve global mandatory policy followed by organization, tenant, workspace, project and task policy.
4. Evaluate classifications, purpose, environment, risk and current resource state.
5. Return a versioned Policy Decision.
6. Enforce obligations before the protected action.
7. Emit an Evidence Record reference without leaking sensitive input.

More-specific policy can narrow broader policy but cannot override a mandatory deny or weaken a higher-precedence obligation.

## Obligations

Obligations may require:

- human approval bound to action and digest;
- model, tool, destination or provider restriction;
- data redaction or local processing;
- network, filesystem, time, token, compute or cost limits;
- independent verification;
- additional evidence, notification or retention;
- cancellation or manual recovery behaviour.

Unknown mandatory obligations cause denial. Optional advice is represented separately and cannot be confused with enforcement.

## Failure behaviour

- Policy unavailability fails closed for retrieval, secret access, mutation, external disclosure, publication and release.
- Cached decisions are usable only within explicit subject, action, resource, version and expiry bounds.
- Indeterminate is not permit.
- Retries preserve the original decision reference or trigger re-evaluation; they never silently switch semantics.
- Break-glass is a separately governed action, not an error fallback.

## Change governance

Policy changes require authenticated authorship, peer or accountable-owner review, version increment, regression and conflict tests, impact preview, rollout and rollback. A policy pack extension cannot grant itself precedence or suppress audit evidence.

## Initial conformance scenarios

- Retrieval is denied before content return when project membership is absent.
- A tool invocation with an unlisted destination is denied at enforcement.
- A changed candidate digest invalidates publication approval.
- Policy version changes during a run trigger the declared continue, re-evaluate or cancel behaviour.
- Deny and obligation decisions are reproducible from recorded non-sensitive inputs and versions.
- Enforcement cannot convert an unknown obligation into a no-op.

The canonical decision shape is `schemas/forge/trust/v1/policy-decision.schema.json`.
