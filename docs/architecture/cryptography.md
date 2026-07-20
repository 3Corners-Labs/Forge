# Cryptographic architecture

Status: Proposed
Owner: Forge Security Architecture

## Purpose

Forge requires portable cryptographic controls for platform data, generated applications, software supply chains, and self-hosted or SaaS operation. Cryptography supplements identity, authorization, isolation, privacy, backup protection, and secure operations; it does not replace them.

## Assets and separation

The following are separate governed assets:

- plaintext and encrypted application data;
- data-encryption keys and wrapped data keys;
- key-encryption and signing keys;
- logical key references, versions, algorithms, and provider bindings;
- credentials and other secrets;
- data backups, key backups, and recovery material;
- audit and approval evidence.

A data-store backup alone SHALL NOT be sufficient to decrypt classified application data. Access to data and key recovery uses separated roles and controls.

## Protection layers

| Layer | Purpose |
|---|---|
| Infrastructure | Protect disks, volumes, snapshots, hosts, and backup media |
| Platform | Protect orchestration state, databases, object storage, and messaging storage |
| Application | Protect selected classified fields or documents from partial data-store or operator compromise |
| Transport | Authenticate and protect external and service-to-service communication |
| Supply chain | Sign and verify source attestations, artefacts, releases, and offline bundles |

Application-level encryption is selective and classification-driven. Encrypting every field without access, search, retention, recovery, and availability design creates risk rather than assurance.

## Envelope encryption

For classified application data, Forge SHOULD support envelope encryption:

1. An authorized workload requests a data key for a tenant, project, and purpose.
2. A cryptographically secure provider generates the data key.
3. The payload is encrypted using centrally governed authenticated encryption.
4. The plaintext data key is cleared from application memory as far as the runtime permits.
5. Ciphertext, wrapped data key, and non-secret metadata are stored together.
6. Decryption reauthorizes the workload, tenant, purpose, operation, and key state.

Authenticated associated data binds ciphertext to tenant, project, object type, immutable object identifier, field or purpose, schema version, and crypto version. Moving ciphertext to a different scope must cause authentication failure.

## Key domains

Key domains follow security purpose, retention, recovery, residency, and destruction boundaries rather than table layout. Potential Forge domains include:

- organization and workspace confidential data;
- project knowledge and restricted source material;
- model and tool provider credentials;
- personal and user-research data;
- evidence exports and support bundles;
- software and release signing;
- analytics pseudonymization.

The final hierarchy requires a threat model and data-classification matrix. One key per database and one key per field are both rejected as universal defaults.

## Crypto Service Port

Domain capabilities SHALL NOT call provider-specific KMS SDKs. A versioned Crypto Service Port defines normalized operations such as data-key generation, encrypt, decrypt, rotate, rewrap, sign, verify, metadata inspection, and provider health.

An encrypted envelope includes:

- crypto and schema version;
- logical provider alias and key reference;
- key version and purpose;
- policy-selected algorithm identifier;
- nonce or equivalent public parameters;
- wrapped data key and ciphertext;
- associated-data schema identifier.

Provider endpoints, accounts, and credentials do not appear in domain data.

## Provider interoperability

Provider adapters MAY use native cloud KMS APIs, KMIP profiles, or a KMS backed by an HSM through PKCS #11. OpenBao Transit is a Trial candidate for a development or open self-hosted reference profile. Candidate status is not acceptance.

Interoperability requires contract and conformance tests covering supported operations, authentication, key attributes, lifecycle, rotation, rewrap, error normalization, availability, audit, and migration. At least two independent providers must pass before provider portability is claimed.

## Key ownership profiles

- Platform-managed: Forge or the self-hosted operator manages key domains and lifecycle.
- Bring your own key: customer material is imported through a controlled ownership and recovery process.
- Customer-controlled KMS: keys remain in a customer system and Forge receives narrowly scoped operation rights.

Customer-controlled KMS introduces an external availability and latency dependency. Onboarding must define shared responsibility, revocation behaviour, caching restrictions, support, and accidental data-unavailability risk.

## Workload identity and authorization

Long-lived shared KMS administrator credentials are prohibited. Workloads use short-lived identity such as federated service accounts, mTLS identities, or an accepted workload-identity system.

Authorization binds workload identity, environment, tenant or partition, project where applicable, key purpose, operation, and approved break-glass context. Normal workloads never receive decrypt-all-tenants authority.

## Key lifecycle

~~~text
Proposed -> Generated or Imported -> Active -> Deactivated -> Revoked
         -> Destroy scheduled -> Destroyed
~~~

Compromised and Recovery-only are explicit exceptional states. Activation records owner, custodian, purpose, policy, rotation, retention, and recoverability. Rotation creates a new active version. Rewrap changes data-key protection without persisting application plaintext and is resumable, idempotent, observable, and auditable.

Revocation and suspected compromise trigger scope analysis, new-encryption stop, credential rotation, rewrap or re-encryption planning, and preserved incident evidence.

Destruction requires dual control, waiting period, dependency inventory, impact preview, legal and retention evaluation, and recoverability confirmation.

## Secrets versus keys

Signing keys and key-encryption keys belong in KMS/HSM boundaries. Provider refresh tokens and database credentials belong in a secrets capability or encrypted credential store. Wrapped data keys may live beside ciphertext. Non-secret endpoints and feature configuration belong in configuration management.

Secrets are short-lived where practical and never embedded in source, images, prompts, logs, traces, or evidence payloads.

## Availability and recovery

Each capability defines behaviour during provider latency or outage. Sensitive operations generally fail closed; unrelated non-sensitive views may degrade and continue. Timeouts, circuit breakers, bounded queues, and explicit retry budgets prevent resource exhaustion.

Plaintext data-key caching is disabled by default. Any exception is process-local, short-lived, non-persistent, revocation-aware, measured without sensitive labels, and explicitly risk-accepted.

Recovery testing covers application data, object storage, workflow state, encrypted credentials, provider bindings, policy, key metadata or supported backup, and evidence. Tests prove correct-key recovery, wrong-tenant denial, preserved key state, new-environment restore, and declared RTO/RPO.

## Audit and observability

High-value evidence includes key generation or import, activation, rotation, rewrap, revocation, destruction, policy and provider-binding changes, failed authorization, break-glass, abnormal operation volume, and backup/restore.

Telemetry can record operation class, latency, normalized error, circuit state, health, and rewrap backlog. It never records plaintext, ciphertext, wrapped keys, nonces, credentials, tokens, or identifying key material.

## Required quality scenarios

Before acceptance, Forge SHALL validate:

- a workload scoped to one tenant cannot decrypt another tenant's envelope;
- provider outage fails sensitive operations safely without resource exhaustion;
- rotation preserves reads and resumable rewrap reaches complete coverage;
- provider migration works through contracts without domain schema changes;
- destructive action cannot be completed by one actor immediately;
- SaaS and self-hosted profiles verify the same signed artefacts and crypto metadata;
- backup and restore recover both data and valid key state without weakening isolation;
- telemetry and crash evidence contain no key or plaintext leakage.
