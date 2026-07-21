# Evidence and observability architecture

Status: Accepted
Owner: Forge Assurance and Reliability Architecture
Version: 0.1

## Purpose

Define separate but correlated lifecycles for audit evidence and operational telemetry. Evidence supports trust and accountability; telemetry supports operation and diagnosis.

## Audit Evidence capability

`forge.capability.evidence` owns append-oriented Evidence Records and integrity relationships. An Evidence Record contains subject references and digests rather than unrestricted copies of source, prompts, secrets or personal data.

Evidence properties:

- attributable actor and scope;
- immutable subject identity and digest;
- source and decision references;
- classification and retention class;
- time and sequence metadata;
- integrity digest and optional signature or previous-record link;
- controlled disclosure and export.

“Tamper-evident” requires an implemented integrity mechanism, protected signing or anchoring keys, independent verification and tests. An append-only database claim alone is insufficient.

## Observability capability

`forge.capability.telemetry` owns metrics, logs, traces and profiles for health, performance, capacity and cost. Telemetry may be sampled, aggregated or expired according to operational policy.

Telemetry uses tenant-safe correlation identifiers. Sensitive source, prompts, personal data, credentials, secret references, ciphertext, key identifiers and unrestricted tool results are excluded or redacted before export.

## Promotion from telemetry to evidence

Telemetry contributes to evidence only through a promotion record containing:

- source system and immutable query or selection description;
- time range and correlation scope;
- transformation and redaction version;
- result digest;
- actor and policy decision;
- evidence classification and retention.

The promoted artefact is independently verifiable. Later telemetry expiry does not silently invalidate a conformance claim; the evidence record declares what supporting material remains available.

## Availability and failure semantics

- Routine telemetry failure must not synchronously stop every read path.
- Evidence required for a consequential side effect is committed before or atomically coordinated through a durable outbox.
- If the evidence sink is unavailable, publication or release stops unless an accepted policy permits deferred evidence with protected local persistence.
- Evidence export is a separately authorized disclosure action.
- Backpressure protects control and execution paths from telemetry overload.

## Retention and privacy

Audit, security, delivery, support and product-learning evidence use different retention classes. Legal hold and deletion restrictions are explicit. Evidence minimization is preferred over indefinite storage of raw prompts or source.

## Initial conformance scenarios

- An Evidence Record fails verification after subject or reference tampering.
- Telemetry from one tenant cannot be queried or exported by another.
- A trace containing a seeded secret is blocked or redacted before export.
- A pull-request publication cannot claim complete evidence when its required record was not durably committed.
- A promoted telemetry result is reproducible from preserved inputs or clearly marked non-reproducible with reason and approval.

The initial Evidence Record shape is `schemas/forge/trust/v1/evidence-record.schema.json`.
