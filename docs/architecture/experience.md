# Experience and explainability

Status: Accepted
Owner: Forge Experience Architecture

## Forge Workbench

Forge SHOULD provide a dedicated workbench that lets people collaborate with and inspect the agentic workforce. The workbench is a projection of versioned platform state, not a separate source of truth.

## Audience modes

- End-user mode explains journeys, visible behaviour, privacy, and limitations.
- Product mode shows outcomes, scope, acceptance criteria, and delivery confidence.
- UX mode shows flows, prototypes, accessibility, research evidence, and design rationale.
- Architecture mode shows capabilities, dependencies, decisions, data, qualities, and threats.
- Engineering mode links intent to code, tests, runtime traces, and change impact.
- Operations mode shows releases, health, cost, incidents, and recovery actions.

## Core visualizations

The workbench SHOULD visualize:

- intent-to-evidence traceability;
- system context, containers, components, and deployment topology;
- user journeys and interactive previews;
- scenario coverage by persona, intent, application version, expected result, observed result, freshness, and confidence;
- data lineage, classification, residency, and retention;
- agent plan, delegation, authority, progress, and hand-offs;
- code impact and dependency changes;
- policy decisions and approval state;
- runtime requests, failures, cost, and service health.

## Explanation contract

For every consequential change Forge SHALL expose:

- what the user asked for;
- what Forge understood and assumed;
- alternatives considered and why one was selected;
- artefacts and behaviour changed;
- data, tools, models, and policies used;
- tests and evaluations performed;
- known limitations, uncertainty, cost, and residual risk;
- approvals required and available rollback.

Explanations must be generated from structured evidence. A model-written narrative alone is not proof.

## Journey assurance view

The Workbench SHALL let accountable owners distinguish accepted, exploratory, quarantined, and deprecated scenarios. It exposes exactly which public interface, synthetic identity, entitlements, fixtures, runner and application artefact produced a result. A changed expected outcome is shown as a product-intent change requiring review, never hidden as routine test maintenance.
