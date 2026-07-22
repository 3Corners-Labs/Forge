# ADR-0011: Establish bounded autonomous contributor delivery

Status: Accepted
Owner: Forge Architecture Board
Date: 2026-07-22

## Context

Forge contributors currently obtain frequent conversational confirmation for routine, reversible delivery steps such as creating a feature branch, committing an already reviewed iteration, pushing that branch, opening a draft pull request, and monitoring CI. Those pauses preserve human control but consume attention and interrupt delivery even when the task scope, architecture, checks, and landing boundary are already explicit.

The desired outcome is a temporary mode in which an AI contributor can carry one approved task through routine implementation and pull-request preparation without waiting after every reversible step. Human control must remain mandatory for consequential authority changes, architecture acceptance, destructive actions, production effects, and landing on the protected default branch.

A repository file or shared code cannot safely act as a password or sole authorization. Candidate-controlled content is untrusted, is visible to the contributor, and can be changed on a feature branch. Platform sandbox prompts, repository protections, identity controls, and higher-authority instructions also cannot be waived by repository content.

This decision governs **contributors working on the Forge repository**. It does not implement or grant authority to Forge runtime agents and does not replace the Identity, Delegation, Policy Decision, Execution Envelope, Approval Record, or Evidence Record architecture.

## Decision

Adopt **Bounded Autonomous Delivery** as an explicit, temporary contributor operating mode.

The mode SHALL require both:

1. a valid human-approved autonomy grant loaded from a protected default-branch revision; and
2. a task-scoped human activation that references the grant identifier and its immutable base revision or digest.

Neither condition alone grants authority. Silence, conversational inference, a feature-branch file, generated content, a copied activation phrase, or a grant whose protected revision cannot be verified SHALL NOT activate the mode.

### Authority and precedence

An autonomy grant is a narrow delivery delegation, not a replacement operating contract. It SHALL remain subordinate to:

1. system and platform safety controls;
2. the Accepted Forge constitution;
3. Accepted ADRs and architecture;
4. the active repository `AGENTS.md` contract;
5. protected-branch and required-review policy; and
6. the exact task scope activated by the human.

A grant cannot waive or reinterpret a higher-authority requirement. Tool or sandbox approval prompts may still occur and do not indicate that Forge policy has been relaxed.

### Grant contract

The implementation SHALL define a closed, versioned grant schema. A grant SHALL include at least:

- grant ID, schema version, lifecycle status, issuer, accountable owner, issue time, expiry, and revocation state;
- repository and protected base revision or digest;
- task, milestone, capability, path, and branch scope;
- allowed routine actions;
- actions that always require human confirmation;
- external systems and destinations in scope;
- time, cost, attempt, concurrency, and change-size limits where applicable;
- required checks, evidence, update cadence, and terminal hand-off state;
- stop conditions and safe recovery instructions; and
- integrity metadata sufficient to detect candidate or base-revision substitution.

The canonical grant location SHALL be declared by the implementation and covered by required architecture checks. Grants SHALL be version-controlled, reviewable, English-language artefacts. A grant identifier, nonce, or digest is not a secret and SHALL NOT be treated as authentication material.

### Activation

The accountable human SHALL activate one grant for one task by referencing the grant ID and protected base revision or digest through the current authenticated collaboration boundary. Activation SHALL also identify the requested terminal state, normally `pull-request-merge-ready`.

Activation expires at the earliest of:

- the grant expiry;
- task completion or cancellation;
- the requested terminal state;
- human revocation;
- a stop condition;
- protected-base or head substitution that invalidates the grant binding; or
- loss of evidence needed to verify the active authority.

The contributor SHALL report activation, material progress, stop conditions, and terminal hand-off. Updates are informational and do not require a response while work remains within the active grant.

### Initially permitted routine actions

Subject to the exact grant, the initial contributor mode MAY authorize:

- load the task-relevant accepted context;
- create or switch to a non-default feature branch;
- edit task-scoped source, tests, schemas, and documentation;
- run formatting, tests, architecture checks, builds, static analysis, and bounded diagnostics;
- repair in-scope deterministic failures and CI failures;
- create task-scoped commits without rewriting published history;
- push ordinary feature-branch commits without force;
- open and update a draft pull request;
- address actionable review feedback without changing accepted product intent;
- mark a pull request ready when the declared checks pass;
- monitor checks, reviews, conflicts, and merge readiness; and
- stop at a documented merge-ready hand-off.

Permission is additive only within the grant. An action omitted from the allowed set remains unauthorized.

### Mandatory human gates

The initial mode SHALL always stop for explicit human direction before:

- accepting, rejecting, superseding, or materially changing an ADR;
- changing an Accepted architecture, capability, trust, deployment, evidence, policy, or public-interface boundary without prior authorization that exactly covers the change;
- expanding task, repository, capability, data, destination, credential, or environmental scope;
- merging or otherwise landing a pull request on a protected default branch;
- releasing, deploying, promoting, or changing production or customer state;
- acquiring or using production credentials, production data, customer data, or an undeclared external destination;
- destructive or difficult-to-recover operations outside an exact prior instruction;
- disabling, weakening, bypassing, or administratively overriding policy, CI, review, branch protection, evidence, or security controls;
- changing accepted user outcomes or independent verification oracles;
- resolving materially ambiguous or contradictory product requirements;
- incurring cost, time, change size, attempts, or side effects beyond the grant; or
- proceeding after authority, evidence, identity, repository state, or safe recovery becomes indeterminate.

Landing SHALL continue to require a fresh, exact pull-request and head-SHA confirmation after all current checks and required reviews pass. A broad grant, earlier approval, or confirmation for another head is never landing approval.

### Candidate non-interference

A candidate branch SHALL NOT activate, widen, renew, replace, or revoke the grant under which it is being evaluated. A proposed grant or operating-contract change becomes eligible only after independent review, human approval, protected landing, and a later activation bound to the new protected revision.

The contributor SHALL resolve active grant authority from the protected revision, not from the candidate worktree. Architecture fitness SHALL test this invariant and reject an active grant that depends on candidate-controlled state.

### Evidence and pull requests

Work performed under this mode SHALL identify the grant ID and protected revision or digest in the pull-request evidence. The pull request SHALL still describe intent, affected capabilities, architecture and security impact, verification, exclusions, rollback, and human gates.

The grant is evidence of delegated delivery scope, not evidence that the change is correct. Required CI, independent review, and exact landing approval remain separate.

### Revocation and recovery

Human revocation or grant expiry SHALL prevent new delegated actions at the next enforceable boundary. In-flight reversible work SHALL stop safely, preserve a reviewable worktree and evidence, and report the last completed action. The contributor SHALL NOT push, publish, or land additional changes after revocation unless separately authorized for recovery.

If automatic continuation cannot prove scope, idempotency, repository state, or authority, the result is a human hand-off rather than an inferred permission.

## Implementation sequence

This Accepted ADR does not by itself activate Bounded Autonomous Delivery. It
becomes implementation authority only after protected landing.

After accountable acceptance and protected landing, implementation proceeds in a separate pull request:

1. add the closed grant schema and canonical protected location;
2. add deterministic validation and negative candidate-substitution fixtures;
3. update `AGENTS.md` to describe resolution, activation, reporting, expiry, revocation, and stop behaviour;
4. add architecture fitness checks for lifecycle, scope, protected revision, mandatory gates, expiry, and non-interference;
5. add a reviewed example grant that is inactive by default;
6. verify documentation, checks, and branch protection; and
7. require a separate human activation before first use.

No grant becomes active merely because the implementation pull request merges.

## Alternatives

- **Unrestricted “god mode”:** rejected because it has no least-privilege boundary, expiry, or reliable stop condition.
- **A code stored in the repository:** rejected because it is visible, copyable, candidate-controlled, and not authenticated human approval.
- **A file that activates itself when present:** rejected because a candidate could expand its own authority.
- **Conversation-only standing permission:** permitted as a temporary task-scoped instruction under the existing operating contract, but rejected as the durable mechanism because scope, expiry, evidence, and protected-revision binding are not machine-verifiable.
- **Reuse Forge runtime delegation objects immediately:** rejected because contributor governance and Forge product runtime authority are separate trust boundaries and runtime delegation is not yet implemented.
- **Include merge, release, or production actions:** rejected because independent human promotion and production controls are constitutional boundaries.
- **Require confirmation after every reversible step:** retained as the default when no grant is active, but rejected as the only mode because it consumes human attention without proportionate risk reduction.

## Consequences

- Humans can delegate routine delivery progression without delegating architecture acceptance, promotion, production authority, or control bypass.
- Contributors can continue through ordinary branch, implementation, verification, commit, push, draft-PR, CI-repair, and review cycles until merge-ready.
- The repository gains another governance contract that requires schema, lifecycle, validation, evidence, and maintenance.
- Candidate non-interference and protected-revision resolution become mandatory fitness properties.
- Platform approval prompts and external repository protections remain visible even during autonomous delivery.
- The mechanism does not prove correctness and cannot be presented as Forge runtime-agent conformance.

## Security and misuse scenarios

Implementation SHALL demonstrate at least:

- a grant changed only on the candidate branch does not activate or widen authority;
- an expired, revoked, unknown-version, malformed, wrong-repository, wrong-base, wrong-branch, wrong-task, or over-budget grant fails closed;
- a copied grant ID or activation phrase without authenticated human activation grants nothing;
- an action outside the explicit allow-list stops before mutation or external side effect;
- a request to merge, release, deploy, use production credentials, or weaken a control always requires a separate human gate;
- a changed pull-request head invalidates any previous landing confirmation;
- revocation stops new delegated actions and produces a recoverable hand-off; and
- the audit trail distinguishes grant issuance, task activation, contributor actions, review, and landing approval.

## Validation

- This ADR became `Accepted` after the accountable human explicitly accepted the contributor-authority boundary recorded below.
- The decision-only pull request SHALL change no active operating instruction or autonomy implementation.
- Documentation checks SHALL verify ADR metadata, register consistency, navigation, and local links.
- The implementation pull request SHALL include the schema, validator, negative fixtures, architecture fitness, operating-contract update, documentation, and rollback evidence defined above.
- First activation SHALL reference the accepted grant ID and protected base revision or digest and stop at merge-ready.

## Accountable review

| Role | Reviewer | Outcome | Date |
|---|---|---|---|
| Product and delivery owner | Sijbren Beukenkamp, temporary | Accept | 2026-07-22 |
| Architecture | Sijbren Beukenkamp, temporary | Accept | 2026-07-22 |
| Security and privacy | Sijbren Beukenkamp, temporary | Accept | 2026-07-22 |
| Engineering | Sijbren Beukenkamp, temporary | Accept | 2026-07-22 |

Sijbren Beukenkamp accepted this decision as drafted at Git object
`1283ef805bb6d3c4ce6e6c4f74458eb23cb0e399` on 2026-07-22. The AI
contributor prepared the proposal and did not provide accountable acceptance.

## Supersession

None. This decision supplements ADR-0010 for contributor delivery only and does not change the M1 product capability boundary.
