# User journey assurance

Status: Proposed
Owner: Forge Product, Experience, and Quality Architecture
Version: 0.1

## Decision

Forge introduces **User Journey Assurance** as a first-class capability and an **End-User Assurance Agent** as one governed workforce role. Together they test whether an application still delivers accepted outcomes to different users through its real public interfaces.

This capability goes beyond conventional system regression testing, but it does not replace unit, contract, integration, performance, accessibility, security, or human usability testing.

## Why this capability is distinct

Technical checks can prove that an API returns a value or a component renders. They do not by themselves prove that:

- the right user can complete the intended outcome;
- the wrong user is prevented from doing so;
- terminology and interaction remain comprehensible;
- privacy, accessibility, localization, interruption, and recovery behaviour remain correct;
- the application still matches the accepted product intent after multiple changes.

User Journey Assurance owns the versioned relationship between intent, personas, accepted scenarios, executable journeys, expected outcomes, and observed evidence. It does not own product intent, identity, authorization policy, test environments, or release approval.

## End-User Assurance Agent

The agent performs journeys as explicit personas with least-privilege synthetic identities. It uses the same public browser, mobile, CLI, API, or portal boundary as a real user. Private database access, internal service calls, hidden test endpoints, and administrative bypasses cannot be used to claim journey success.

The role can:

- execute accepted regression and smoke scenarios;
- compare observed behaviour with versioned expected outcomes;
- collect screenshots, traces, accessibility information, console and network signals, and referenced backend evidence;
- exercise positive, negative, interrupted, concurrent, and recovery paths;
- propose exploratory scenarios and identify contradictions or missing acceptance criteria;
- explain results from the perspective of each persona.

The role cannot:

- create or expand its own identity, entitlements, environment scope, or tool authority;
- change an accepted expectation because implementation behaviour differs;
- approve its own generated scenario as a release gate;
- use a language-model score as the sole oracle for a consequential outcome;
- expose secrets or personal data in recordings, prompts, logs, or evidence.

## Scenario contract

The canonical machine-readable contract begins with `schemas/forge/assurance/v1/journey-scenario.schema.json`. Each scenario records:

- stable scenario ID, owner, lifecycle, intent references, and risk class;
- persona, role, entitlements, denied capabilities, locale, viewport, and accessibility profile;
- interface type and public entry point;
- environment constraints and synthetic-data fixtures;
- preconditions, steps, checkpoints, expected outcomes, and forbidden outcomes;
- deterministic assertions, optional secondary semantic evaluation, and tolerances;
- evidence, privacy, retention, retry, and quarantine policy;
- application, interface, scenario, dataset, runner, browser, and model versions required for replay.

Expected outcomes SHALL describe observable product behaviour rather than internal implementation structure unless internal evidence is itself the accepted requirement.

## Scenario lifecycle and authority

~~~text
discovered -> proposed -> reviewed -> accepted -> active
                                     |          |
                                     v          v
                                  rejected   deprecated
~~~

- **Accepted scenarios** protect established intent and may contribute to release eligibility.
- **Exploratory scenarios** probe unanticipated paths and generate candidates, findings, or hypotheses; they do not block release by themselves.
- **Quarantined scenarios** remain visible with owner and expiry. Flakiness cannot be hidden by deletion or unlimited retry.
- A changed expected result is a product decision with impact analysis, not routine test maintenance.

Product owns intended outcomes. Experience owns interaction and accessibility expectations. Security owns abuse and authorization expectations. Quality owns scenario construction and execution independence. Release governance decides how accepted results contribute to a gate.

## Assertion model

Use the strongest available oracle in this order:

1. deterministic state, contract, invariant, or property assertion;
2. independently generated reference result or approved baseline;
3. bounded statistical comparison with declared confidence;
4. qualified human evaluation;
5. probabilistic semantic evaluation as supporting evidence.

Visual snapshots and model judgments are useful diagnostics but become brittle or self-confirming when treated as universal truth. Material scenarios therefore require independent evidence and explicit tolerances.

## Execution model

Each build produces an impact set from changed capabilities, interfaces, intent links, dependencies, risk, and recent failures. Forge selects:

- a deterministic smoke pack for every deployable preview;
- all accepted scenarios directly affected by the change;
- mandatory cross-cutting scenarios for identity, tenant isolation, accessibility, privacy, and recovery;
- broader packs on schedule, release candidates, or elevated risk;
- exploratory missions within an explicit time, cost, data, and action budget.

Tests run in ephemeral or dedicated non-production environments with synthetic or explicitly approved test data. Results use `schemas/forge/assurance/v1/journey-result.schema.json` and bind the exact scenario, application artefact, deployment, runner, identity profile, data fixture, and evidence records.

## Evidence and visualization

The Workbench presents a matrix of intent, persona, scenario, application version, environment, result, confidence, freshness, and evidence. A user can replay a journey, inspect failure steps, compare releases, and see which accepted intent lacks coverage.

Evidence is minimized and classified. Screenshots, video, traces, DOM/accessibility trees, and network records follow explicit retention and redaction policy. Journey telemetry becomes audit evidence only through the normal evidence-promotion contract.

## Initial technology direction

Technology selection remains a Trial decision, not an architecture dependency. For browser-facing applications, Playwright is the leading open-source candidate because it supports cross-browser end-to-end execution, isolation, parallel CI, traces, screenshots, and accessibility integrations. Non-browser interfaces require adapters that implement the same scenario and evidence contracts.

## Initial conformance scenarios

- A persona completes one accepted happy path through the public interface.
- A user without an entitlement is denied and sees a safe, comprehensible response.
- A scenario cannot pass by calling an internal endpoint unavailable to the persona.
- Changing an expected result requires review and creates visible history.
- An exploratory agent can propose but cannot activate a release-gating scenario.
- A failure result binds the exact artefact, environment, identity profile, fixture, and trace.
- Evidence capture redacts seeded secrets and personal data.
- A quarantined flaky scenario has an accountable owner, reason, expiry, and visible risk.

## Delivery maturity

The capability and agent role are **Research**. The first Trial should execute one accepted Forge Workbench journey with two authorization roles in a preview environment. It becomes Supported only after repeatability, flake rate, evidence privacy, execution cost, negative authorization behaviour, and release-gate integration are measured.

## External reference

- [Playwright documentation](https://playwright.dev/docs/intro) — candidate browser end-to-end runner; selection remains subject to ADR, security, licensing, portability and operational evaluation.
