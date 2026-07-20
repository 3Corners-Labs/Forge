# ADR-0003: Separate authority objects from context objects

Status: Proposed
Owner: Forge Security and Knowledge Architecture
Date: 2026-07-20

## Context

The earlier architecture allowed identity, policy, model, tool and classification fields to appear in both the Execution Envelope and Context Package. Divergent copies could create stale authorization or confused-deputy behaviour.

## Decision

The Execution Envelope is the immutable authorization artefact for one execution. Policy Decisions and Approval Records grant or constrain specific actions. The Context Package is a derived information manifest that references one envelope and cannot grant authority.

Use the Proposed schemas under `schemas/forge/trust/v1/` as the contract baseline. Secrets appear only as authorized references. Every consequential publication binds policy, verification and approval to an immutable candidate digest.

## Alternatives

- One combined request/context object: rejected because untrusted content and authorization would share a mutable boundary.
- Let workers evaluate their own final authority: rejected because workers process untrusted inputs and cannot widen delegation.
- Provider-native authorization only: rejected because Forge requires portable policy semantics and evidence.

## Consequences

Policy is evaluated at multiple gates. Context assembly, workers, tools and adapters must receive a validated envelope. Child delegation produces a narrower envelope with parent lineage.

## Validation

- Context mutation cannot add an action, destination, model, tool or secret reference.
- Expired, mismatched-tenant and mismatched-digest objects are denied.
- Approval replay against a different action, environment or digest fails.
- Schemas are validated in architecture CI.

## Supersession

None.
