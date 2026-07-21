# First-slice data flow and privacy analysis

Status: Proposed
Owner: Forge Security and Privacy Architecture
Version: 0.1

## Scope

This document describes data movement, ownership, trust boundaries, minimization, and privacy obligations for the trustworthy intent-to-pull-request slice. It is a design and review artefact, not evidence that controls are implemented.

The slice ends after branch and pull-request publication. Merge, production deployment, customer application data, product learning, billing, and Workplace-specific data are excluded.

## Logical data flow

~~~mermaid
flowchart LR
    Human[Authenticated human] --> Experience[Public CLI or Workbench]
    Experience --> Control[Delivery control]
    Control --> Identity[Identity and workspace]
    Control --> Policy[Policy decision]
    Policy --> Envelope[Execution Envelope]
    Envelope --> Context[Context assembly]
    Context --> Worker[Isolated worker]
    Repository[Source repository] --> Context
    Repository --> Worker
    Worker --> Candidate[Candidate digest and artefacts]
    Candidate --> Verification[Independent verification]
    Verification --> Approval[Human approval]
    Approval --> Publication[Source-control adapter]
    Publication --> Provider[Source-control provider]
    Control --> Evidence[Evidence sink]
    Worker --> Evidence
    Verification --> Evidence
    Publication --> Evidence
    Control --> Telemetry[Operational telemetry]
    Worker --> Telemetry
~~~

Each arrow crossing a capability or external boundary requires schema validation, immutable scope, policy enforcement appropriate to the action, classification handling, and safe correlation.

## Data inventory

| Data | Canonical owner | Default class | Included content | Explicit exclusions |
|---|---|---|---|---|
| Identity assertion | External identity provider; normalized facts owned by Identity | Confidential | subject, issuer, audience, assurance, expiry | reusable credentials in prompts or logs |
| Workspace scope | Tenant and workspace | Internal | tenant, workspace, project, repository binding, membership | caller-authoritative tenant values |
| Intent and acceptance criteria | Intent and delivery | Internal; may become Confidential | objective, constraints, expected outcomes, author | unrelated conversation history |
| Policy Decision | Policy | Internal | action, resource, outcome, obligations, versions, reason codes | unrestricted sensitive policy inputs |
| Execution Envelope | Delivery and trust issuance boundary | Confidential | identifiers, granted actions, limits, references, integrity | secret values and retrieved content |
| Context Package | Context assembly | Inherits highest included item | source references, versions, purpose, omissions, conflicts, manifest digest | authority grants and unrestricted repository dump |
| Repository checkout | Source repository; temporary worker copy | Usually Confidential | authorized revision and required files | unrelated repositories and hidden credentials |
| Worker credentials | Credential or secret provider | Restricted | short-lived, task-scoped references resolved at execution | credential material in durable workflow state |
| Candidate artefacts | Intent and delivery references; artefact store owns bytes | Same as source or higher | changed files, build outputs, candidate digest | production secrets or unrelated workspace data |
| Verification result | Verification | Internal; findings may be Confidential | criteria, outcomes, findings, versions, evidence refs | uncontrolled full logs or seeded secrets |
| Approval Record | Approval | Internal | actor, action, digest, decision, conditions, expiry | authentication secrets or copied candidate content |
| Publication receipt | Source-control adapter | Internal | repository, branch, pull request, provider receipt, observed digest | provider token |
| Evidence Record | Evidence | Classification declared per record | actor, scope, subject digest, references, outcome, integrity | raw prompts, unrestricted source, secrets, bulk telemetry |
| Operational telemetry | Telemetry | Internal after minimization | correlation, operation class, latency, resource and error category | source, prompt, credential, personal content, key identifiers |

Classification is evaluated from actual tenant policy and source metadata. Table defaults never downgrade a stronger classification.

## Processing sequence and gates

| Step | Data entering boundary | Required gate | Data leaving boundary |
|---|---|---|---|
| Authenticate | identity assertion | issuer, signature, audience, expiry, assurance | normalized principal |
| Resolve workspace | principal and requested resource | membership and tenant isolation | immutable scope and repository binding |
| Capture intent | human request and scope | input schema, classification, purpose | intent record and delivery-run ID |
| Issue authority | principal, scope, action, risk, classification | deterministic policy evaluation | Policy Decision and Execution Envelope |
| Assemble context | envelope, objective, authorized source refs | authorize before discovery and retrieval | minimal Context Package manifest |
| Execute change | envelope, plan, context, source revision | worker isolation, secret, filesystem, network, resource, and expiry enforcement | candidate digest, artefact refs, structured result |
| Verify | candidate digest and declared criteria | independent runner and evidence policy | findings and Verification Result |
| Approve | exact digest, action, evidence, actor | approver authority, separation of duties, expiry | Approval Record |
| Publish | envelope, decision, approval, candidate, binding | exact-digest, repository, branch, idempotency, and evidence readiness | remote branch and pull-request receipt |
| Record evidence | subject and referenced outcomes | classification, minimization, retention, integrity | append-oriented Evidence Record |

Policy is re-evaluated when subject, action, resource, classification, destination, candidate digest, environment, policy version, or approval validity changes in a way defined by policy.

## Trust and processing boundaries

- The human client is untrusted until authentication and scope resolution complete.
- The public interface cannot set authoritative tenant or project identifiers.
- Repository files and history are untrusted content even when access is authorized.
- Context assembly is an authorization enforcement point, not merely a ranking stage.
- The worker is an untrusted execution boundary and receives no control-plane database access.
- Model and tool providers, if introduced, are external processing destinations requiring separate eligibility. The first deterministic slice does not require a model.
- The source-control adapter is an external side-effect boundary and cannot decide Forge policy.
- Evidence and telemetry are separate stores with different writers, retention, and disclosure paths.
- A self-hosted profile cannot make an undeclared call to Forge SaaS.

## Purpose limitation and minimization

The permitted first-slice purposes are capturing one requested change, preparing one candidate, verifying it, obtaining publication approval, publishing one pull request, and preserving proportionate evidence.

The following are separate, excluded purposes and require later review:

- model training or provider improvement;
- employee analytics or productivity scoring;
- organization-wide memory creation;
- cross-tenant analytics;
- support access to source or execution content;
- product learning from prompts or repository content.

Context uses exact files and sections where practical. Telemetry records measurements and safe categories rather than payload content. Evidence uses digests and references rather than copying canonical data.

## Retention and deletion obligations

Before Trial, every data category SHALL receive an owner-approved retention class with an expiry or review rule. The implementation SHALL support:

- immediate expiry of task credentials at completion, cancellation, or revocation;
- deterministic cleanup of ephemeral worker storage;
- deletion or invalidation propagation to indexes, cached context, and derived projections;
- preservation of required change and evidence records according to repository and organizational policy;
- legal-hold handling without silently converting all data to indefinite retention;
- backup reconciliation so logically deleted data does not silently become active after restore;
- provider-held data deletion or explicit documented limitations.

No fixed retention duration is invented by this architecture document. Product, legal, privacy, security, and operations owners must approve the profile-specific schedule.

## Privacy and security scenarios

- A caller-supplied tenant identifier conflicting with verified membership is denied before repository discovery.
- A Context Package for one project cannot reference or retrieve another project's content.
- A seeded credential or personal identifier is absent from telemetry, routine evidence, screenshots, and provider requests.
- Restricted content cannot reach a model, tool, or source-control description without explicit destination eligibility.
- Cancellation or delegation revocation prevents subsequent credential resolution and publication.
- Evidence export requires a separate disclosure decision and reveals no content unavailable to the requesting actor.
- Worker cleanup removes the checkout and resolved credentials while retaining only authorized artefact and evidence references.
- Restore processing preserves logical deletion and tenant isolation.

## Required evidence before Trial

- reviewed data-flow diagram with capability and external boundaries;
- profile-specific data inventory and retention schedule;
- negative cross-tenant tests for API, store, cache, queue, context, telemetry, and evidence paths;
- secret and personal-data seeding tests across logs, traces, errors, evidence, screenshots, and provider calls;
- worker cleanup, cancellation, revocation, and restore reconciliation tests;
- provider destination and residency policy tests if any external intelligence or tool is enabled;
- privacy-owner approval of purposes, retention, deletion, support access, and evidence disclosure.

## Open decisions

- identity provider and assurance requirements;
- developer and production-shaped tenancy topology;
- repository and artefact storage locations;
- retention schedule and legal-hold ownership;
- evidence integrity and signing-key design;
- permitted external providers, regions, retention, and training-use terms;
- privacy request, export, deletion, and support-access procedures;
- backup scope, restore reconciliation, and verified recovery objectives.

These decisions must be resolved for the applicable profile before a production-readiness claim.
