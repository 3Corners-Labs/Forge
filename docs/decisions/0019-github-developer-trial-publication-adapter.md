# ADR-0019: Select the GitHub Developer Trial publication adapter

Status: Accepted
Owner: Forge Source Control, Security, Delivery, and Architecture
Date: 2026-07-23

## Context

M1 is complete only when Forge prepares one real reviewable pull request against the Forge repository. The Accepted source-control contract requires exact repository binding, candidate digest, independent verification, human publication approval, Policy authority, evidence readiness, idempotency, and reconciliation of ambiguous provider outcomes.

F0 permits only a contract double before M1 and leaves the live provider and credential model unselected. The Forge repository is hosted on GitHub, but conversation access, a developer's `gh` session, or the contributor autonomy grant is not Forge product authority and cannot silently become the runtime publication credential.

Publication is the only external side effect in M1. Requiring human confirmation for every API call would defeat bounded autonomy, while letting Forge merge, change protection, broaden repository access, or retry blindly after timeout would exceed the milestone.

## Decision

Implement one replaceable GitHub publication adapter for the Developer Trial. Use a dedicated GitHub App installation identity that can write one approved branch and create one pull request in the exact Forge repository. It receives a short-lived installation access token only after publication authority and evidence readiness pass.

The adapter contract remains provider-neutral. GitHub identifiers and API semantics are confined to the adapter and its provider receipt.

### Publication authority

Publication requires a new phase-specific Policy Decision, delegation where required, and Execution Envelope for action `forge.source-control.publish-pr`. They bind:

- authenticated initiating human and publication workload identity;
- exact tenant, workspace, project, delivery run, and Forge repository binding;
- immutable base revision and allowed base branch;
- exact candidate identifier, tree or commit digest, and changed-path manifest;
- independent Verification Result and applicable Assurance Plan digests;
- publication Approval Request and Approval Record v2 for the exact candidate digest;
- required Evidence Record identifiers and readiness digest;
- exact provider `github`, installation, repository owner/name or immutable repository identifier, and destination classification;
- deterministic branch name and pull-request operation identity;
- title/body metadata digest after disclosure and minimization policy;
- environment `developer-trial`;
- no merge, release, deployment, policy, protection, administration, secret, workflow, issue, or unrelated repository action; and
- bounded expiry, attempts, API calls, duration, and cost.

Publication approval uses action `forge.source-control.publish-pr`, subject type `verified-candidate`, and the exact verified candidate digest. It is separate from plan approval and cannot authorize worker execution or merge.

The adapter resolves all referenced canonical records from their owners. Caller-provided approval, evidence, verification, policy, repository, branch, or candidate values are untrusted hints.

### GitHub App identity and credentials

The Developer Trial uses a dedicated GitHub App installed only on the Forge repository or an explicitly enumerated test repository. Minimum permissions are:

- repository metadata: read;
- contents: read and write, only as required to create or update the approved branch reference;
- pull requests: read and write, only as required to create or reconcile the approved pull request.

The App receives no administration, organization, members, actions, secrets, environments, deployments, releases, packages, checks-write, security-event-write, or branch-protection permission.

[GitHub documents](https://docs.github.com/en/rest/apps/apps) that its pull-request permission covers pull requests, related operations, and merges. The credential permission alone therefore cannot prove that an App with `pull_requests:write` is unable to call merge. Forge treats this as an explicit provider limitation and enforces the M1 no-merge boundary through all of:

- an adapter endpoint allowlist containing no merge, review-decision, close, protection, release, or deployment operation;
- Policy and Execution Envelopes containing no merge action;
- repository rules on the protected default branch that deny updates by the GitHub App and grant it no bypass;
- a separately owned human landing path outside Forge product authority; and
- negative provider tests that attempt merge and protected-branch update with the App identity and require denial.

If protected-branch rules cannot prove that denial for the selected repository, the live M1 publication journey remains blocked.

The App private key or token-minting credential is held by the Developer Trial credential provider as a secret reference. It is not stored in Delivery, Approval, Policy, Evidence, logs, prompts, context, artefacts, or the worker. The publication adapter obtains a short-lived installation token after admission, keeps it in memory only, uses it for the exact operation, and destroys the handle on completion, cancellation, expiry, or failure.

Installation-token lifetime and repository permissions bound credential use. M1 does not claim per-request cryptographic token scope when GitHub cannot express it; Forge additionally enforces exact operation scope in its envelope and adapter. That residual risk is documented and measured rather than hidden.

Developer `gh` credentials, Codex connector credentials, personal access tokens, SSH agent state, and environment tokens are rejected as automatic Forge runtime credentials. A contract test may use synthetic credentials only.

### Candidate materialization

The publication adapter does not run candidate code. It materializes the already verified candidate into a fresh authorized checkout or provider commit operation and proves that:

- base revision equals the approved immutable base;
- applied candidate bytes reproduce the verified candidate digest;
- resulting tree and changed paths equal the verified manifest;
- no file outside the approved path set changed;
- commit metadata is minimized and attributable to the GitHub App plus human approval references; and
- branch and pull-request metadata contain no secret, unrestricted source, prompt, personal data, or unsafe generated instruction.

The candidate creator, worker, and repository content cannot choose repository owner, base branch, destination, credentials, approval, required evidence, or idempotency identity.

### Idempotency and reconciliation

GitHub does not provide one portable idempotency key for the complete branch-and-pull-request operation. Forge therefore uses a deterministic publication identity derived from tenant, delivery run, repository, base revision, candidate digest, and protected action.

The adapter derives:

- a collision-resistant Forge-owned branch name containing a safe run identifier and candidate-digest prefix;
- an immutable provider-operation record in PostgreSQL;
- a machine-readable publication marker in pull-request metadata that exposes only safe operation identity; and
- expected base, head, candidate, repository, and approval digests for reconciliation.

Before mutation and after every timeout, transport error, process loss, or retry, the adapter queries GitHub for:

1. the exact branch reference;
2. its observed commit and candidate/tree digest;
3. open and closed pull requests from that branch to the approved base; and
4. the exact publication marker and repository identity.

Outcomes are:

- no remote state: perform the next authorized mutation;
- exact existing branch or pull request: return the existing receipt idempotently;
- partial exact state: continue only the missing authorized step;
- conflicting branch, digest, repository, base, marker, or multiple pull requests: stop with conflict or manual recovery;
- indeterminate visibility or unavailable evidence: bounded retry, then manual recovery; never create a second branch or pull request blindly.

The provider receipt records safe immutable repository, branch, commit, pull-request, observed digest, operation, timestamps, and response-correlation references. It contains no token or unrestricted response body.

### Bounded autonomous publication

One valid human publication approval authorizes Forge to perform the exact branch and pull-request API sequence plus bounded reconciliation without asking for approval between provider calls. This autonomy ends when:

- any authorized digest, scope, destination, approval, evidence, policy, credential, base revision, or branch constraint changes;
- GitHub exposes conflicting or irreconcilable state;
- attempt, API-call, time, or cost budgets are exhausted; or
- cancellation or revocation occurs.

Forge MAY automatically recover an exact already-completed outcome. It SHALL NOT automatically choose a different branch, overwrite conflicting remote work, force-push, merge, close, reopen, approve, label beyond an explicit safe allowlist, modify protection, rerun workflows, release, or deploy.

### Public CLI and explanation

The public `forge delivery start` journey presents:

- intent, scope, context omissions, plan, risk, budgets, checks, candidate digest, verification, evidence readiness, publication subject, limitations, and rollback;
- immediate acknowledgement and visible progress for long-running operations;
- explicit plan and publication approval waits;
- cancellation and manual-recovery state; and
- final branch, pull-request URL, provider receipt, remaining human review responsibility, and statement that Forge cannot merge.

The CLI invokes public capability contracts; it does not call GitHub directly or bypass Policy, Approval, Verification, Evidence, or Delivery state.

### Contract double and provider isolation

A deterministic contract double remains the primary unit and failure-injection adapter. The same conformance suite runs against the GitHub adapter for:

- creation, exact replay, partial completion, timeout before and after remote mutation, duplicate visibility, conflict, cancellation, revocation, expired approval, changed digest, and provider unavailability;
- credential denial, minimization, cleanup, rate-limit, pagination, and disclosure-safe error handling; and
- replacement without Delivery, Approval, Verification, Evidence, or Policy domain changes.

GitHub API versions, SDKs, media types, and webhook details remain adapter implementation choices. A material SDK or framework selection still requires ADR-0013 eligibility evidence.

### Iteration boundary

After this decision is Accepted and upstream M1 states are implemented, a bounded iteration MAY create one branch and one pull request in the exact approved Forge repository. It SHALL NOT merge, approve its own pull request, change branch protection, use production credentials or data, publish to another repository, release, deploy, or claim general Git provider support.

## Alternatives considered

- **Use the developer's existing `gh` login:** rejected because ambient personal authority is broader, longer-lived, and outside Forge product evidence.
- **Use a fine-grained personal access token:** rejected as the default because it remains a reusable human credential and weakens workload attribution and revocation.
- **Use SSH push plus separate PR API:** rejected because it introduces two credential and reconciliation surfaces for one M1 outcome.
- **Require a human click for each GitHub call:** rejected because one exact publication approval should authorize the bounded idempotent operation.
- **Give Forge merge permission for convenience:** rejected because merge is explicitly outside M1.
- **Use GitHub-specific objects in Delivery:** rejected because provider replacement must not change canonical delivery or trust contracts.
- **Retry blindly after timeout:** rejected because it can create duplicate or conflicting remote state.
- **Postpone the live provider until M2:** rejected because M1 exit explicitly requires one real reviewable Forge pull request.

## Consequences

- M1 gains one real external publication path while keeping provider semantics at the edge.
- A GitHub App must be created, installed, permission-reviewed, rotated, and tested with protected-branch rules before the live journey.
- Short-lived installation tokens reduce but do not eliminate risk from the App's installation-level repository permissions.
- Reconciliation adds provider reads and durable operation state but permits safe autonomous recovery.
- GitHub remains a Trial adapter, not a platform-wide vendor commitment.

## Validation

Evidence SHALL prove:

- exact Policy, delegation, envelope, Approval v2, Verification, Evidence, repository, base, candidate, path, branch, provider, environment, and expiry binding;
- denial for plan-approval replay, changed candidate, stale verification, missing evidence, cross-repository, cross-branch, cross-environment, unauthorized metadata, revoked authority, and expired publication approval;
- GitHub App installation restricted to enumerated repositories and minimum permissions, with no administration permission and explicit acknowledgement that `pull_requests:write` also covers merge endpoints;
- protected-default-branch rules denying App updates or bypass, adapter endpoint allowlisting, and negative merge/protected-update tests;
- token minting only after admission, in-memory use, no ambient credential fallback, expiry, cancellation, revocation, and cleanup;
- deterministic branch and publication identity, exact replay, partial completion, timeout-before-mutation, timeout-after-mutation, duplicate result, conflict, pagination, rate-limit, and manual recovery;
- no force-push, overwrite, merge, release, deployment, branch-protection, workflow, secret, or unrelated repository mutation;
- candidate and resulting tree digest reproduction plus changed-path equality;
- seeded secret and personal-data markers absent from branch, commit, pull-request metadata, receipts, telemetry, evidence, and errors;
- contract-double and GitHub adapter conformance through the same domain contract;
- allowed and denied journeys through `forge delivery start`; and
- p50, p95, p99, API-call count, throughput, rate-limit, timeout, reconciliation, credential, and cost evidence under ADR-0013.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-23, acting as the temporary accountable Product, Architecture, Security and Privacy, Identity, Policy, Approval, Verification, Evidence, Source Control, Operations, and Engineering reviewer established for bounded M1.

Live provider implementation remains blocked until the decision reaches protected `main`, the GitHub App permission review is complete, and a separate bounded grant is activated.

## Supersession

If accepted, this decision specializes ADR-0009, ADR-0010, ADR-0013, ADR-0015, ADR-0016, and the F0 Developer Trial security decisions for one GitHub publication adapter. It does not authorize merge, production use, general Git provider support, release, or deployment.
