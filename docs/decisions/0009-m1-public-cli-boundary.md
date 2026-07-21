# ADR-0009: Establish the M1 public CLI boundary

Status: Accepted
Owner: Forge Product and Architecture
Date: 2026-07-21

## Context

G9 requires the allowed and denied M1 journeys to use the same public interface available to a real user. The accepted first-slice contract names the Experience command `SubmitIntent`, but it does not fix CLI syntax. A stable boundary is required before independent journey oracles can be reviewed without binding them to internal services.

Sijbren Beukenkamp approved the recommended `forge delivery start` boundary on 2026-07-21. This decision defines only the public contract needed for F0 assurance design; it does not authorize or claim an implementation.

## Decision

The M1 public CLI entry point for `SubmitIntent` SHALL be:

~~~text
forge delivery start --workspace <workspace> --project <project> --repository <repository> --intent-file <path> --output json
~~~

The workspace, project, and repository arguments are selectors. They are untrusted input and SHALL be resolved against the authenticated principal; they never establish authoritative tenant or resource scope. No tenant override is exposed.

The JSON response SHALL use the normalized success or error contract. A permitted request may capture intent and return a delivery-run reference, but cannot itself approve a plan, execute a worker, publish a branch, or open a pull request. A denied or indeterminate request SHALL fail closed before repository discovery and return a safe stable error code.

The command SHALL use the same Experience and capability APIs as any later Workbench. Tests may replace external providers through contract doubles, but cannot call private services, databases, or hidden test endpoints to claim journey success.

## Alternatives

- Leave syntax to M1 implementation: rejected because the G9 oracle would not bind a real public interface.
- Use `forge intent submit`: rejected because the user initiates a delivery journey, while Intent remains an internal capability boundary.
- Permit a caller-supplied tenant flag: rejected because tenant scope is server-derived.
- Test internal capability commands directly: rejected by ADR-0005 and the accepted journey-assurance architecture.

## Consequences

- M1 must add the `delivery start` command without changing the existing developer-only `up` and `down` safety boundary.
- CLI contract tests and both G9 journeys can share one public entry point.
- Later additive options may follow normal contract evolution. Renaming the command, weakening scope resolution, or changing its consequential side effects requires architecture review and a superseding ADR.
- This decision does not select identity, repository, worker, model, or source-control providers.

## Validation

- Architecture checks SHALL reject G9 scenarios that use another entry point or omit the public CLI command.
- One permitted and one denied role SHALL exercise the boundary with deterministic expected and forbidden outcomes.
- M1-exit tests SHALL prove the public command cannot bypass identity, workspace, policy, approval, evidence, or Developer-profile controls.

## Supersession

None.
