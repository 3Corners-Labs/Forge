# Software supply-chain architecture

Status: Accepted
Owner: Forge Engineering and Security Architecture
Version: 0.1

## Purpose

Ensure that Forge and Forge-generated applications can identify, reproduce, verify, promote, revoke and recover every executable artefact and dependency.

## Build and promotion model

~~~text
reviewed source revision
  -> isolated reproducible build
  -> tests and policy gates
  -> SBOM and provenance
  -> immutable artefact digest
  -> signature or attestation
  -> environment promotion without rebuild
~~~

Promotion references immutable digests. Mutable tags are discovery aliases only and cannot be the sole release identity.

## Required artefacts

Each release candidate includes:

- source revision and repository identity;
- build definition and builder identity;
- dependency lock or resolved dependency evidence;
- test and evaluation references;
- software bill of materials;
- provenance attestation;
- vulnerability and license assessment;
- artefact digest and signatures;
- deployment manifest and compatibility range;
- rollback and revocation metadata.

## Dependency governance

Dependencies, build actions, base images, tools, model artefacts and extensions have owner, source, version, integrity, license, update policy and support status. CI actions and production artefacts SHOULD be pinned to immutable revisions or digests with an automated update process.

Version tags are temporarily permitted in the Developer Trial profile only when no production trust claim is made and automated monitoring records the gap. Production-shaped evidence requires digest resolution and verification.

## Build isolation

Builds run without production credentials. Network access is deny-by-default or dependency-proxy restricted. Untrusted build scripts cannot access signing keys directly. Signing occurs after policy and evidence validation through a narrowly scoped signing service or trusted workflow identity.

## Provenance and verification

Consumers verify artefact digest, expected builder, source revision, policy identity, signature trust root, revocation state and profile compatibility before execution. A valid signature from an untrusted publisher is not sufficient.

Offline bundles include the artefacts, verification material, trust-root update procedure, revocation snapshot, SBOM, provenance, installation plan and rollback assets needed without Forge SaaS.

## Vulnerability and update handling

Findings are evaluated against reachability, deployment profile, exploitability and compensating controls. Exceptions have owner, expiry and evidence. Emergency updates preserve provenance and review; urgency does not authorize an untracked rebuild.

Filesystem, dependency, SBOM, image and IaC scanning is part of the [continuous security testing](continuous-security-testing.md) portfolio. Trivy is an initial Trial candidate, not a hard platform dependency. Scanner databases, policies and versions require provenance, freshness, offline-update and rollback handling just like other executable supply-chain inputs.

## Initial conformance scenarios

- Changing one artefact byte causes verification failure.
- A release cannot be promoted by mutable tag alone.
- The same digest is promoted between compatible test and reference environments.
- A revoked extension or signing identity is rejected.
- An offline bundle verifies and installs without network access.
- Rollback uses a previously verified artefact and compatible data migration state.
- CI cannot expose signing material to untrusted pull-request code.

## Current bootstrap status

The existing Compose file and GitHub workflows are Developer Trial assets. Several images and actions still use version tags. They SHALL NOT be represented as production provenance until digest pinning, signature policy, SBOM generation and verification are implemented.
