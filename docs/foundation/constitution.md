# Forge engineering constitution

Status: Accepted
Owner: Forge Architecture
Version: 1.1

This constitution translates the [Forge enduring principles](principles.md) into non-negotiable engineering guardrails. It becomes normative only when its lifecycle status is Accepted.

## 1. Human authority

Humans define outcomes, risk appetite, policy, and approval boundaries. Production-impacting actions require an authorized human or an explicitly delegated organizational policy.

## 2. Repository truth

Version-controlled artefacts and accepted decisions are the source of truth. Generated views SHALL be reproducible from those sources.

## 3. Intent and evidence

Every material change SHALL be traceable from intent through decisions and implementation to verification evidence.

## 4. Secure by design

Identity, least privilege, data minimization, isolation, policy enforcement, and auditability SHALL be architectural properties rather than optional integrations.

## 5. Untrusted intelligence

Model output, prompts, retrieved content, external events, and tool results SHALL be treated as untrusted. Models propose; deterministic controls authorize and verify.

## 6. Explainability

Forge SHALL explain consequential behaviour, assumptions, uncertainty, data use, policy decisions, and verification results to the appropriate audience.

## 7. Portability

Core capability contracts SHALL remain deployment- and vendor-neutral. SaaS and self-hosted editions share the same logical architecture and policy semantics.

## 8. Modular evolution

Forge begins as a modular platform. A capability becomes a separately deployed service only when evidence justifies the operational cost.

## 9. Open contracts

Boundaries use versioned, testable contracts and open standards where suitable. Provider-specific behaviour stays behind adapters.

## 10. Reversible delivery

Changes SHALL be small, reviewable, observable, and reversible where practical. Irreversible operations require additional approval and recovery planning.

## 11. Responsible resource use

Compute, tokens, storage, and human attention are governed resources. Quality, privacy, latency, carbon impact, and cost are explicit selection factors.

## 12. Continuous assurance

Architecture, security, accessibility, quality, and operational fitness SHALL be evaluated throughout delivery rather than at the end.

## 13. Identity, context, and policy before execution

Every execution SHALL have an attributable identity, one declared objective, an authorized and reproducible context package, and an evaluated policy decision. Runtime components enforce policy obligations and SHALL NOT reinterpret or expand them.

## 14. Explicit capability contracts

Every platform capability SHALL have a stable identifier, owner, responsibility boundary, versioned inputs and outputs, lifecycle, evidence obligations, security considerations, observability, and conformance criteria before it is represented as complete.

## 15. Knowledge and memory integrity

Repositories store information; accepted governance establishes knowledge. Published knowledge SHALL preserve provenance and authoritative source. Durable memory SHALL require explicit scope, ownership, validation, retention, and deletion policy.

## 16. Bounded execution

Tools, models, workers, and agents SHALL operate inside task-scoped limits for identity, data, filesystem, network, secrets, time, compute, cost, and side effects. They SHALL NOT expand their own authority.

## 17. Independent trust

No consequential output SHALL approve, certify, or provide the sole evidence for itself. Trust decisions require independent deterministic validation, separation of duties, or accountable human review proportional to risk.

## 18. Observable and recoverable operation

Every material state transition and side effect SHALL be observable, attributable, and correlated with its intent and evidence. Failure, cancellation, retry, compensation, rollback, backup, restore, and manual recovery SHALL be designed according to risk.

## 19. Governed evolution

Architecture, policies, prompts, models, tools, knowledge, schemas, and standards SHALL use explicit lifecycle states. Breaking changes require impact analysis, migration, compatibility policy, and review; no implementation may claim conformance without objective evidence.

## 20. Safe self-development

Forge MAY develop Forge only through an isolated candidate workspace, branch protection, independent checks, and human-controlled promotion. A candidate version SHALL NOT weaken its own controls, approve its own consequential change, or attest its own trustworthiness.

## 21. Preserve accepted user outcomes

Material user journeys and expected outcomes SHALL be versioned, role-aware, traceable to product intent, and tested through the same public interfaces available to real users. An implementation, test generator, or assurance agent SHALL NOT silently redefine accepted behaviour to make a change pass.

## 22. Bound offensive assurance

Active or autonomous security testing SHALL operate only within externally enforced scope, explicit rules of engagement, isolated identity and environment limits, impact controls, evidence policy, and emergency termination. Automated testing SHALL NOT be represented as proof that no vulnerabilities exist or as a universal replacement for qualified human assurance.
