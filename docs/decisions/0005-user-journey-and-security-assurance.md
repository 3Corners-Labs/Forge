# ADR-0005: Separate user journey and continuous security assurance

Status: Accepted
Owner: Forge Product, Quality, and Security Architecture
Date: 2026-07-20

## Context

Forge required journey, regression, accessibility, adversarial, and security tests, but did not define who owns accepted end-user scenarios, how exploratory scenarios gain authority, or how active and autonomous security testing remains safely bounded. Treating these concerns as generic test scripts would lose product intent and create an unsafe path for offensive automation.

## Decision

Introduce two explicit capabilities:

1. **User Journey Assurance** owns versioned role-based scenarios, expected and forbidden outcomes, public-interface execution semantics, and journey evidence. An End-User Assurance Agent is one governed executor of that capability, not its authority owner.
2. **Continuous Security Testing** owns risk-based assurance selection, rules of engagement, offensive-test scope, normalized findings, and coverage disclosure. Any autonomous security-testing agent is treated as an untrusted offensive component constrained outside the model runtime.

Accepted scenarios can contribute to release eligibility. Exploratory scenarios remain non-blocking candidates until accountable review. Automated security testing complements human-led penetration testing and cannot support an absence-of-vulnerability claim.

## Alternatives

- Keep all journeys under generic quality automation: rejected because expected product behaviour and persona authority would have no durable owner.
- Let an AI agent infer expected results on each run: rejected because implementation behaviour could redefine intent and create circular validation.
- Run the broadest active scan after every build: rejected because mutation, duration, cost, third-party effects, and attack risk differ by change and environment.
- Treat scanner output as release policy: rejected because findings require validation, context, ownership, confidence, and exception governance.

## Consequences

- Product, experience, quality, security, and release ownership must be explicit for each accepted journey.
- Preview environments require synthetic identities, fixtures, reset, evidence redaction, and public interface stability.
- CI needs impact-based selection and visible quarantine governance.
- Offensive tooling requires stronger identity, network, rate, target, data, audit, and kill-switch controls than ordinary test runners.
- Forge gains a reusable assurance model for Forge itself and every generated application.

## Validation

- Execute one accepted happy path and one denied path as different roles through a public Forge interface.
- Prove that the agent cannot activate or change the expected result without review.
- Bind each result to scenario, application artefact, deployment, identity profile, fixtures, runner, and evidence.
- Reject out-of-scope targets and prohibited techniques outside the agent runtime.
- Demonstrate emergency termination, environment cleanup, finding validation, and coverage disclosure.
- Record measured repeatability, flake rate, cost, duration, false positives, and evidence leakage.

## Supersession

None.
