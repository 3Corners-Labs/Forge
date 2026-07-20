# ADR-0001: Begin as a modular platform

Status: Proposed
Owner: Forge Architecture
Date: 2026-07-20

## Context

Forge has an ambitious capability model but no production implementation or measured scaling boundaries. Starting with many network services would introduce distributed failure modes, contract overhead, and operational cost before those boundaries are validated.

## Decision

Build the initial control plane as a modular monolith with enforced internal capability boundaries. Run untrusted or resource-intensive work in isolated worker processes. Persist owned data behind module interfaces. Use versioned contracts at external boundaries.

Extract services only under the criteria in the deployment architecture and through a separate accepted ADR.

## Consequences

This reduces initial operational complexity and enables refactoring while the domain is learned. It requires automated dependency rules, explicit module ownership, and avoidance of shared-domain shortcuts.

## Validation

Architecture tests SHALL detect forbidden module dependencies. Operational tests SHALL demonstrate isolation between the control plane and workers.
