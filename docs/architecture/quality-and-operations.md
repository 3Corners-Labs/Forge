# Quality, resilience, and operations

Status: Proposed
Owner: Forge Quality and Reliability Architecture

This design consolidates the existing quality engineering, continuous validation, test automation, observability, resilience, localization, accessibility, and operational-model concepts.

## Continuous assurance

Quality evidence is produced at every lifecycle transition. Verification is selected by affected capability and risk, and can include:

- contract, unit, integration, system, and journey tests;
- accepted multi-role user journeys executed through public interfaces, including denied and interrupted paths;
- architecture and dependency fitness functions;
- security analysis and adversarial evaluation;
- accessibility and localization checks;
- performance, capacity, reliability, and recovery tests;
- model and agent evaluations using versioned datasets;
- human review of product, UX, architecture, and residual risk.

Generated code and generated tests do not independently validate each other. High-risk behaviour requires an independent oracle, invariant, reference implementation, or qualified human review.

Journey scenarios are governed product artefacts rather than disposable test scripts. Accepted expected outcomes cannot be rewritten by the implementation or agent under test. Exploratory scenarios generate candidate knowledge and do not become release gates without review. See [User Journey Assurance](user-journey-assurance.md).

Security assurance is layered. Safe source, dependency, configuration and applicable passive checks run continuously; active and autonomous testing requires an authorized, resettable target and bounded rules of engagement. Automated scanners complement rather than replace qualified penetration testing. See [Continuous Security Testing](continuous-security-testing.md).

## Observability

Forge correlates intent, delivery run, agent, model, tool, artefact, deployment, and runtime outcome. Metrics, traces, logs, profiles, audit evidence, and cost records have separate access and retention controls. Sensitive prompts, secrets, personal data, and source content are redacted or excluded by default.

Operational telemetry is not audit evidence by default. A telemetry item becomes evidence only through a deliberate promotion that binds its source, time window, query or transformation, integrity digest, classification and retention class. Evidence storage SHALL NOT be used as an unrestricted high-volume telemetry backend.

## Resilience

Durable workflows support deadlines, cancellation, retry policy, idempotency, checkpoints, compensation, and manual recovery. Queues enforce backpressure and fairness. External providers use circuit breakers, bulkheads, quotas, and health-aware routing.

## Operability

Every supported deployment profile defines SLOs, alerts, runbooks, capacity assumptions, backup, restore, upgrade, rollback, incident response, and disaster recovery. A feature is not complete when it cannot be operated or explained.

## Inclusive delivery

Accessibility and localization are design inputs. User-facing artefacts SHALL support keyboard and assistive technology, comprehensible language, adaptable detail, locale-aware presentation, and testing against declared accessibility targets.
