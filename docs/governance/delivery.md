# Delivery governance

Status: Proposed
Owner: Forge Engineering

Every material change follows this minimum lifecycle:

1. Capture intent, scope, risks, and acceptance evidence.
2. Challenge architectural and security impact.
3. Work on a reviewable branch in an isolated environment.
4. Generate or modify artefacts with traceability.
5. Run proportionate automated and human verification, including affected accepted user journeys through public interfaces.
6. Record unresolved risks and required approvals.
7. Release through a reversible strategy.
8. Observe outcomes and feed learning back into product intent.

Pull requests SHALL identify affected capabilities, decisions, security impact, user impact, verification evidence, rollout, and rollback. Repository protection SHALL prevent direct changes to the default branch.

Exploratory scenarios may inform product work but do not become release gates until reviewed. Active or autonomous security testing requires separately authorized rules of engagement and cannot inherit permission merely from a delivery run.
