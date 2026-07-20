# Architecture roadmap

Status: Proposed
Owner: Forge Product and Architecture

## Horizon 0: coherent foundation

- Accept the constitution, capability model, trust model, and first ADR.
- Close the mandatory gates in the architecture readiness register.
- Reconcile document status, capability maturity, and evidence labels.
- Accept or revise the metamodel, trust-object, and deployment-profile ADRs.
- Define non-functional requirements and deployment profiles.
- Establish documentation and architecture fitness checks.
- Define accepted journey, exploratory scenario, and offensive-test rules-of-engagement contracts.
- Replace development defaults that could be mistaken for production posture.

## Horizon 1: trustworthy vertical slice

- Create workspace and authenticated principal.
- Capture intent and produce a delivery plan.
- Execute one isolated repository change.
- Apply policy before execution and approval before publication.
- Link the change to tests and evidence.
- Execute one public-interface journey with an allowed and a denied role where the interface exists.
- Explain and visualize the result in a minimal workbench.
- Meet the exit criteria in `architecture/vertical-slice.md` before broadening capability scope.

## Horizon 2: collaborative workforce

- Introduce product, architecture, UX, engineering, quality, and security roles.
- Trial the End-User Assurance Agent against versioned scenarios and synthetic identities.
- Add structured hand-offs, challenge reviews, and autonomy levels.
- Add model and tool registries with routing, quotas, and evaluation.
- Support preview environments and reversible releases.
- Add passive application scanning and bounded authenticated security journeys to previews.

## Horizon 3: platform deployment

- Validate self-hosted installation, upgrades, backup, and recovery.
- Validate SaaS tenant isolation, regional controls, and managed operations.
- Publish portable project and evidence formats.
- Establish SLOs and production readiness evidence.
- Establish scheduled active DAST, curated vulnerability regressions, WSTG coverage mapping, and human penetration-test intake.

## Horizon 4: continuously improving software company

- Connect runtime and user feedback to product intent.
- Automate bounded maintenance, remediation, and optimization.
- Let exploratory end-user agents propose new product scenarios while accepted expectations remain human-governed.
- Consider autonomous offensive testing only after independent scope, isolation, kill-switch, audit, and finding-validation conformance.
- Expand method and compliance profiles without changing core semantics.
- Increase autonomy only where outcome and safety evidence supports it.
