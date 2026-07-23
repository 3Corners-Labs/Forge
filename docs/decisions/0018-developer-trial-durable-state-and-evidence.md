# ADR-0018: Select Developer Trial durable state and evidence adapters

Status: Accepted
Owner: Forge Data, Evidence, Reliability, and Architecture
Date: 2026-07-23

## Context

M1 requires durable delivery state, optimistic concurrency, cancellation, checkpoints, bounded retry, transactional outbox, manual recovery, append-oriented evidence, candidate artefact references, and publication blocking when required evidence is unavailable. Current Identity, Delivery, Policy, Context, and planned Approval modules use in-process stores that prove ownership and idempotency but do not survive process loss.

The Developer Compose profile already includes PostgreSQL and MinIO. Their presence is development bootstrap, not authority to make them canonical state or evidence providers. F0 deliberately left persistence providers unselected and prohibits production durability, recovery, or support claims.

Autonomous delivery requires reliable checkpoints and reconciliation. Without durable state, every interruption becomes a human restart. Excessive synchronous evidence writes, however, would place every read and low-risk transformation behind an audit bottleneck. The selected design must support autonomy and recovery while keeping evidence proportional and separate from telemetry.

## Decision

Use PostgreSQL as the Developer Trial adapter for capability-owned control state, transactional outboxes, idempotency, and append-oriented evidence metadata. Use MinIO as the Developer Trial adapter for immutable candidate and verification artefact bytes identified by SHA-256 digest.

Both remain replaceable adapters behind capability contracts. This is a Developer Trial profile choice, not a production topology, cloud, portability, capacity, backup, RTO, or RPO commitment.

### Ownership inside one PostgreSQL instance

The modular control plane uses one PostgreSQL instance with capability-owned schemas and database roles:

- Tenant and Identity own their admitted scope and delegation state;
- Delivery owns intents, plans, run versions, transitions, checkpoints, cancellation, recovery state, candidate references, and its outbox;
- Policy owns decision and envelope records;
- Context owns Context Package records and reproduction metadata;
- Approval owns Approval Requests, Approval Record v2 chains, idempotency, and revocation;
- Verification owns Assurance Plans, findings, and Verification Results;
- Evidence owns Evidence Records, chain heads, verification state, and disclosure metadata;
- source-control publication later owns provider operation and reconciliation records.

Modules do not read another capability's tables. Cross-capability access uses explicit ports, immutable identifiers and digests, or owned outbox events. Database ownership does not collapse capability authority.

Migrations are versioned per capability, run in a declared order, and refuse unknown or incompatible schema versions. Candidate-controlled code, workers, verifiers, repository content, and tools receive no database credentials or network route to PostgreSQL.

### Transaction and outbox model

Each capability command commits only its owned state and outbox entries in one local transaction. No transaction spans capability schemas, MinIO, Docker, GitHub, or another external provider.

Consumers process outbox events at least once and record event IDs or use naturally idempotent commands. An event is never visible before its source state commits. Delivery coordinates cross-capability progress through explicit states and revalidates resolved facts before each transition.

The outbox enables autonomous continuation after transient process loss. Human intervention is required only when:

- an external outcome is ambiguous and cannot be reconciled;
- immutable authority or evidence no longer matches;
- a migration or integrity check cannot prove safety; or
- retry, duration, cost, or recovery budgets are exhausted.

Ordinary progress, retry, checkpoint, and cleanup events do not require human approval.

### Optimistic concurrency and idempotency

Tenant and delivery-run identifiers partition state. Delivery updates require the exact current version and increment it atomically. Approval, Policy, Context, Evidence, and later publication commands use tenant-scoped action-and-target idempotency keys.

Unique constraints prevent duplicate:

- run transitions for one version and command identity;
- Approval decisions or decision-chain branches;
- Policy Decisions, delegations, envelopes, Context Packages, Assurance Plans, Verification Results, and Evidence Records under their canonical command identity;
- candidate artefacts under one digest; and
- provider operations under one publication identity.

An idempotency-key replay with identical immutable input returns the existing result. Changed immutable input returns an explicit conflict and never overwrites canonical state.

### Artefact storage

MinIO stores only bounded immutable artefact bytes needed for M1:

- candidate patch or bundle;
- changed-path and tree manifests;
- structured worker result attachments;
- bounded verification output and reproducibility fixtures; and
- provider receipts or evidence attachments only when references are insufficient.

Every object has a server-derived key, content digest, size, media type, classification, tenant/project/run scope, retention class, creator capability, and creation time in PostgreSQL metadata. Object bytes are uploaded to a temporary key, verified, and promoted to an immutable digest key. A metadata transaction references only a successfully verified object.

Source checkouts, unrestricted repository content, prompts, credentials, environment dumps, raw terminal streams, and personal data are not durable artefacts. Workers and verifiers use task-scoped pre-authorized object operations and cannot list unrelated prefixes.

### Evidence integrity

Evidence Records remain distinct from telemetry and domain state. For M1, Evidence uses:

- canonical SHA-256 subject and record digests;
- exact references to Policy Decisions, envelopes, context, plans, approvals, worker results, candidates, verification, publication, and transitions;
- a previous-record digest within a tenant/project/run evidence stream;
- serializable chain-head update or equivalent compare-and-swap to prevent undetected forks;
- immutable record rows and separate append-only correction or revocation records;
- independent `VerifyEvidence` recomputation from stored canonical fields and referenced immutable digests; and
- evidence-readiness projections used by protected publication gates.

SHA-256 chain integrity and capability-owned resolution are sufficient for the bounded single-host Developer Trial. No signature, timestamping authority, transparency service, hardware key, or production non-repudiation claim is made. Those require the accepted cryptography decision and a later profile-specific ADR.

If required evidence cannot be committed or verified, consequential publication blocks. Non-consequential telemetry or advisory evidence failure follows its declared semantics and does not stall unrelated read paths.

### Retention, deletion, backup, and recovery

Every row and object declares a candidate retention class from the reviewed Developer Trial decisions. No fixed production duration is invented.

M1 implements:

- cleanup of ephemeral checkout and temporary object keys;
- expiry or revocation of task authority and credentials;
- invalidation of derived context, plans, approvals, verification, and publication readiness when an authority-relevant source changes;
- deletion markers and reference checks that prevent silently active dangling objects;
- deterministic local reset for synthetic fixtures; and
- documented manual recovery from an interrupted control-plane process using committed state and outbox entries.

M1 does not claim validated backup, restore, legal hold, disaster recovery, multi-host availability, or production deletion propagation. Before Trial beyond the bounded development profile, the owner-approved retention schedule and backup/restore reconciliation tests remain mandatory.

### Performance and autonomy

Persistence is staged so that:

- cheap admission and state updates target the ADR-0013 control-operation envelope;
- independent writes with no transition dependency may batch or run asynchronously;
- mandatory evidence needed by a later gate is produced early and reconciled at that gate;
- indexes follow immutable tenant/run, state/version, idempotency, subject digest, outbox readiness, evidence chain, and artefact digest access paths;
- bounded workers use backpressure and per-run concurrency rather than unbounded database connections; and
- telemetry measures latency, contention, queue depth, storage, error rate, and saturation without copying protected payloads.

No cache may become an authority source. Evidence reuse requires exact digest, version, environment, and validity matching under ADR-0013.

### Failure and recovery semantics

- PostgreSQL unavailable: block state mutation and protected actions; bounded retry may resume from committed state.
- MinIO unavailable before required artefact commit: do not record candidate completion.
- outbox delivery interrupted: retain committed event and retry idempotently.
- evidence chain conflict: fail closed and enter investigation or manual recovery.
- metadata exists without verified object: treat as incomplete and never publish.
- object exists without committed metadata: quarantine or delete through bounded reconciliation.
- process crashes after an external request: later provider adapters reconcile remote state before retry.
- migration unknown or partially applied: stop startup and require recovery; never guess schema state.

### Iteration boundary

The bounded M1 implementation MAY use PostgreSQL and MinIO for the path through candidate verification and publication evidence. It SHALL preserve in-memory adapters for fast contract tests. It SHALL NOT represent the Developer Compose volumes as production durability, use customer data, expose storage publicly, give workers database access, or claim Supported maturity.

## Alternatives considered

- **Keep all M1 state in memory:** rejected because process loss would erase committed delivery and force repeated human intervention.
- **Use local files as canonical state:** rejected because optimistic concurrency, transactional outbox, queryable recovery, and integrity constraints would be reimplemented poorly.
- **Use SQLite for every capability:** rejected for this profile because PostgreSQL already exists and the M1 validation envelope includes concurrent delivery and worker operations; SQLite remains possible for another profile.
- **Store artefact bytes in PostgreSQL:** rejected because bounded object storage separates immutable bytes from queryable control and evidence metadata.
- **Store every event synchronously in Evidence:** rejected because evidence is not a synchronous dependency for every read path and unnecessary coupling would harm responsiveness.
- **Use telemetry as evidence:** rejected because telemetry may be sampled, aggregated, dropped, or retained differently.
- **Require cryptographic signatures now:** rejected because the bounded Trial can detect mutation and forks without prematurely selecting production key infrastructure.
- **Share tables across capabilities for convenience:** rejected because one database does not authorize cross-capability data ownership.

## Consequences

- Autonomous runs can resume from committed checkpoints instead of restarting from human intent.
- PostgreSQL and MinIO become explicit Developer Trial dependencies with migration, health, capacity, and cleanup obligations.
- Capability isolation is logical and role-enforced inside one instance; it is not a claim of physical tenant isolation.
- Evidence mutation and chain forks are detectable without production signatures.
- More state improves recovery but increases retention, deletion, migration, and test responsibilities.

## Validation

Evidence SHALL prove:

- capability-owned schema and role boundaries with negative cross-schema access tests;
- optimistic-concurrency conflicts, atomic state/outbox commit, at-least-once consumption, duplicate delivery, and idempotent recovery;
- process interruption before and after each state, outbox, object, evidence, and provider boundary;
- no worker or verifier database access;
- tenant/project/run/object-prefix isolation and disclosure-safe errors;
- candidate upload, digest verification, immutable promotion, orphan quarantine, and missing-object denial;
- evidence canonicalization, subject/reference mutation detection, previous-digest chain verification, fork prevention, correction lineage, and unavailable-sink publication blocking;
- retention-class presence, ephemeral cleanup, fixture reset, invalidation, and no unrestricted source, prompt, secret, personal data, or environment dump in durable records;
- migration forward, incompatible-version refusal, rollback/manual-recovery documentation, and interrupted-migration handling;
- p50, p95, p99, throughput, connection use, lock contention, queue depth, storage growth, saturation, and failure recovery under the M1 validation envelope; and
- no production durability, availability, backup, recovery, portability, or support claim.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-23, acting as the temporary accountable Architecture, Security and Privacy, Data, Evidence, Reliability, Operations, Delivery, and Engineering reviewer established for bounded M1.

Functional implementation remains blocked until the decision reaches protected `main` and a separate bounded grant is activated.

## Supersession

If accepted, this decision specializes ADR-0001, ADR-0004, ADR-0010, and ADR-0013 for Developer Trial persistence and evidence. It does not select a production datastore, object store, signing system, backup platform, or deployment topology.
