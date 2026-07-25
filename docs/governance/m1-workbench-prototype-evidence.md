# M1 Workbench prototype evidence

Status: Proposed
Owner: Forge Product, Experience, Delivery, Quality, and Engineering
Verified: 2026-07-25

## Iteration outcome

This iteration turns the accepted M1 kernel into a product that can be tried by
a human. `forge workbench serve` starts one loopback-only API and an embedded,
responsive web cockpit. `forge workbench chat` provides a conversational
terminal client over that same API. Both surfaces project the existing guided
M1 journey rather than implementing a parallel workflow.

A user can now:

1. describe a desired outcome in free text;
2. inspect the Requirements Proposal that Forge derived;
3. promote the exact proposal to planning authority;
4. inspect and approve the bounded plan;
5. watch worker preparation and independent verification;
6. inspect changed paths, test status, digests, evidence, agent state, and
   dependencies;
7. separately approve publication; and
8. reach one reviewable contract-double pull request with merge still
   forbidden.

The cockpit deliberately labels itself `Trial` and `Exploratory`. It has no
remote listener, login, durable Workbench store, live GitHub effect, merge,
release, deployment, or production authority.

## Why

The kernel already proved its contracts and authority chain, but the product
team could not yet experience Forge as an application builder. The Workbench
creates the first feedback surface without weakening M1: human judgement is
requested at the three consequential transitions, while routine projections,
verification, and status updates remain automatic.

The API-first shape is intentional. Web, terminal, and a future Webex
Messaging adapter can share one interaction model instead of diverging into
channel-specific delivery engines.

## Applicable architecture

- The Engineering Constitution requires human intent, bounded delegated
  authority, independent verification, evidence, accessibility, and privacy.
- ADR-0001 and the accepted module graph keep Experience as a surface over
  capability APIs.
- ADR-0009 and ADR-0021 preserve the guided public journey and its separate
  gates.
- ADR-0017 preserves bounded worker execution and independent verification.
- ADR-0019 forbids merge and keeps publication provider-neutral.
- ADR-0020 requires free text to become a proposal rather than authority.
- `docs/architecture/experience.md` requires consistent state and decisions
  across interfaces.
- `docs/architecture/agentic-delivery.md` and
  `docs/architecture/evidence-and-observability.md` govern agent, dependency,
  progress, and evidence projections.

This is a new Trial experience projection over accepted contracts. It does not
change a canonical trust boundary, capability owner, provider commitment, or
production topology. No ADR is required.

## Capability effects

| Capability | Workbench projection |
|---|---|
| Experience | Owns the HTTP surface, web assets, terminal conversation, status wording, and decision presentation |
| Intent | Supplies the proposal and exact digest shown before promotion |
| Delivery | Supplies run state, plan, candidate, and lifecycle progress |
| Identity and Policy | Remain responsible for admitted principals, scope, delegation, and decisions inside the composed journey |
| Approval | Retains separate plan and publication approval records |
| Execution | Supplies bounded candidate preparation state; the UI grants no worker authority |
| Verification | Supplies independent result and verification digest |
| Evidence | Supplies the traceability chain and activity projection |
| Source Control | Supplies the reviewable contract-double receipt with `mergePermitted: false` |

No Workbench package reads another capability's store. `internal/cli` composes
the existing public journey, and `internal/experience/workbench` contains only
the interface assets and API description.

## Implementation

- `GET /api/v1/overview` lists local runs.
- `POST /api/v1/runs` converts bounded free text into a Requirements Proposal.
- `GET /api/v1/runs/{id}` returns the current projection.
- `POST /api/v1/runs/{id}/actions` admits only `promote`,
  `approve-plan`, and `approve-publication` in the correct state.
- `GET /api/v1/runs/{id}/events` streams evidence-backed activity with
  server-sent events.
- The server rejects non-loopback listeners, disables caching for API data,
  limits JSON input, rejects unknown fields, and sends browser security
  headers.
- The web cockpit uses React and React Flow for the agent/dependency view and
  is compiled into the Go binary. Its exact Trial API is described by
  `internal/experience/workbench/openapi.yaml`.
- The terminal client uses the same HTTP resources, starts an ephemeral local
  API when none is supplied, and requires the word `approve` at every
  consequential gate.

## Generated verification

Automated tests cover:

- the complete API journey and exact capability call order;
- fail-closed invalid, unsupported, and out-of-order actions;
- preservation of `mergePermitted: false`;
- browser security headers and embedded asset serving;
- the terminal journey and three explicit decisions;
- UI presentation helpers;
- a real Chrome journey through every decision, build/test/evidence
  inspection, and final reviewable state; and
- an axe-core WCAG scan of the completed cockpit.

The browser test exposed and the candidate fixed two issues before evidence
was recorded: a final empty decision list was serialized as `null`, and two
small sidebar labels did not meet WCAG AA contrast.

Observed local commands:

~~~text
go test ./internal/cli ./internal/experience/workbench
pnpm test
pnpm run build
pnpm run test:e2e
pnpm audit --audit-level high
python3 scripts/check_architecture.py
~~~

The dependency audit reported no known vulnerabilities.
The exact committed revision must still pass the full repository and protected
CI suites.

## How to try it

Web:

~~~text
go run ./cmd/forge workbench serve
open http://127.0.0.1:7331
~~~

Terminal:

~~~text
go run ./cmd/forge workbench chat
~~~

The Trial uses the deterministic, credential-free contract composition. It is
for product feedback, not a claim of durable, multi-user, live-provider, or
production readiness.

## Next product tranche

The next tranche SHOULD use this Workbench to prepare one real bounded Forge
self-change, then add durable run resume and the model-assisted virtual Product
Manager conversation. Webex Messaging can subsequently become another client
of the same API. Live source control, merge, release, deployment, and
production remain separate human-controlled scopes.
