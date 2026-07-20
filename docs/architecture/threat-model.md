# Forge platform threat model

Status: Proposed
Owner: Forge Security Architecture
Version: 0.1

## Scope

This threat model covers the first intent-to-PR vertical slice and its reusable platform boundaries. It is a living architecture artefact, not proof that controls are implemented.

## Protected assets

- human, workload, agent and provider identities;
- tenant, workspace, project and repository boundaries;
- source, requirements, architecture and context content;
- policies, delegations, approvals and execution envelopes;
- secrets, credentials, keys and provider tokens;
- candidate artefacts, tests, provenance and signatures;
- evidence, audit integrity and operational telemetry;
- model, tool, extension and dependency supply chains;
- financial, token, compute and human-attention budgets.

## Trust boundaries

~~~mermaid
flowchart LR
    Human --> Experience
    Experience --> API[Public API boundary]
    API --> Control[Control plane]
    Control --> Trust[Trust plane]
    Control --> Knowledge[Knowledge plane]
    Control --> Worker[Isolated worker boundary]
    Worker --> Repository[Repository boundary]
    Worker --> Tools[Tool and extension boundary]
    Worker --> Models[Model provider boundary]
    Control --> Evidence[Evidence sink boundary]
    Control --> Telemetry[Telemetry pipeline]
~~~

Every arrow crossing a boundary requires authenticated identity where applicable, authorization, schema validation, classification handling, bounded resources and evidence proportional to side effects.

## Threat actors and failure sources

- external attacker;
- malicious or compromised tenant user;
- over-privileged or compromised operator;
- compromised agent, worker, tool, extension or dependency;
- malicious repository or document author;
- compromised model or infrastructure provider;
- accidental human or configuration error;
- faulty deterministic software or probabilistic output;
- resource exhaustion, outage or inconsistent external state.

## Priority threat scenarios

| Threat | Primary boundary | Required controls | Required evidence/test |
|---|---|---|---|
| Direct or indirect prompt injection | Context, repository, tool output | Instruction/data separation, envelope enforcement, output schema, side-effect gates | Adversarial corpus cannot widen actions or disclose seeded data |
| Confused deputy | API, worker, adapter | Immutable scope, audience-bound delegation, action/resource policy | Cross-purpose and cross-tenant negative tests |
| Cross-tenant leakage | Storage, cache, queue, index, telemetry | Server-derived scope, partition keys, encryption domains, query enforcement | Automated negative isolation suite |
| Malicious repository execution | Worker | Ephemeral isolation, no default secrets, egress allow-list, resource/time limit | Escape, exfiltration and fork-bomb tests |
| Tool abuse or unexpected side effect | Tool broker | Risk metadata, validated arguments, dry-run, idempotency, approval | Unauthorized destination and replay tests |
| Approval replay | Publication/release | Digest, action, actor, environment and expiry binding | Changed-digest and expired-approval denial |
| Poisoned knowledge or memory | Knowledge | Provenance, lifecycle, confidence, independent approval, revocation | Superseded and malicious sources excluded or labeled |
| Model/provider data leakage | Provider adapter | Classification eligibility, minimization, provider policy, regional routing | Seeded restricted content blocked from ineligible provider |
| Evidence tampering | Evidence sink | Digests, signatures or chaining, restricted writers, independent verification | Mutation detection and key-compromise procedure |
| Denial of wallet or capacity | Control/provider | Admission control, quotas, budgets, deadlines, bounded fallback | Budget exhaustion stops cleanly without cascade |
| Supply-chain compromise | Build/extension | Provenance, SBOM, digest pinning, signature policy, revocation | Tampered artefact and untrusted publisher rejected |
| Candidate self-attestation | Verification | Independent runner/oracle and separation of duties | Candidate cannot alter required external check result |
| Operator misuse | SaaS operations | Just-in-time access, approval, session evidence, customer visibility | Unauthorized support access denied and alerted |
| Scenario oracle capture | Journey assurance | Intent-owned expectations, independent oracle, reviewed changes, bounded semantic evaluation | Candidate cannot rewrite expected outcome or self-approve scenario |
| Synthetic-user data leakage | Browser, evidence, test data | Synthetic identities/data, redaction, classification, retention, restricted capture | Seeded secrets and personal data absent from evidence |
| Offensive testing escapes scope | Security runner, DNS, redirect, target | External allow-list, target revalidation, egress control, rate/impact bounds, kill switch | Out-of-scope and rebound targets rejected; hard stop verified |
| Untrusted scanner or template | Tool and supply chain | Version pinning, provenance, permission review, isolation, output validation | Unreviewed template and tampered scanner rejected |

## Security assumptions requiring validation

- The identity provider and source-control provider can issue sufficiently scoped, short-lived credentials.
- Worker isolation can enforce filesystem, process, network and resource limits on supported hosts.
- Policy and evidence stores can meet required availability without unsafe bypass.
- Provider contracts expose enough data-handling information for eligibility decisions.
- Self-hosted operators can protect signing, encryption and recovery keys.

An invalidated assumption triggers an ADR or profile restriction.

## Abuse and misuse review

Each agent, model, tool and extension definition declares intended use, prohibited use, maximum autonomy, accessible classifications, external effects, budgets, monitoring and emergency disablement. Safety evaluation covers harmful generated content where the application domain or user population requires it; generic moderation does not replace authorization or secure tool design. Active and autonomous offensive testing additionally requires versioned rules of engagement and external scope enforcement as defined in [Continuous Security Testing](continuous-security-testing.md).

## Review triggers

Review this threat model when a new trust boundary, data class, deployment profile, extension type, external provider, production action, learning use or autonomous capability is proposed, and after material incidents or failed adversarial tests.
