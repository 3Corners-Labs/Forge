# Trust and security architecture

Status: Accepted
Owner: Forge Security Architecture

## Trust model

The following are untrusted by default: users, agents, models, prompts, retrieved content, repositories, dependencies, tool output, web content, and inbound events. Internal network location does not establish trust.

## Security invariants

- Every action is attributable to a human, workload, or delegated agent identity.
- Authorization is evaluated before context retrieval, secret access, tool use, data mutation, and release.
- Tool arguments and model outputs are schema-validated and policy-evaluated.
- Secrets are short-lived, scoped, injected at execution time, and excluded from prompts and evidence.
- Execution is isolated by tenant, project, risk, and data classification.
- Egress is denied by default and explicitly allow-listed.
- Sensitive data is minimized, classified, encrypted, retained, and deleted according to policy.
- Audit evidence is integrity-protected and access-controlled.
- Approval is bound to an artefact digest, action, actor, environment, and expiry time.

The trust-object schemas under `schemas/forge/trust/v1/` form the Accepted design baseline under ADR-0003 and separate Execution Envelope, Context Package, Policy Decision, Approval Record and Evidence Record. Their presence is not an implemented-conformance claim.

## Mandatory threat scenarios

Threat modelling SHALL cover prompt injection, confused deputy behaviour, cross-tenant leakage, malicious repositories, dependency compromise, tool abuse, data exfiltration, approval replay, poisoned memory, unsafe generated code, model-provider compromise, denial of wallet, and evidence tampering.

## Deployment posture

Development conveniences such as default credentials, broad port binding, debug exporters, and dev-mode identity or secret services SHALL be technically prevented from running under a production profile.

SaaS requires tenant isolation, regional data controls, managed key separation, abuse protection, and operator access governance. Self-hosted deployment requires offline-capable installation, customer-controlled identity and keys, signed artefacts, documented backup and recovery, and explicit support boundaries.

Detailed key, envelope-encryption, provider, lifecycle, and recovery rules are defined in the [cryptographic architecture](cryptography.md).
