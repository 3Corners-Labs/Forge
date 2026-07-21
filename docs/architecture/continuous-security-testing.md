# Continuous security testing

Status: Accepted
Owner: Forge Security and Quality Architecture
Version: 0.1

## Decision

Forge orchestrates a layered security-assurance portfolio across every applicable build, preview, release candidate, and supported runtime profile. Security testing is selected by risk, target, change impact, environment, and approved rules of engagement.

Forge SHALL NOT equate automated scanning with a complete penetration test or claim that a passing pipeline proves absence of vulnerabilities.

## Assurance layers

| Layer | Typical cadence | Scope | Initial open-source candidate |
|---|---|---|---|
| Source and secrets | Every change | Source, generated code, configuration, credentials | Existing CI controls; candidate selection pending |
| Dependencies and artefacts | Every build | Filesystem, lockfiles, SBOM, images, licenses, known vulnerabilities | Trivy as Trial candidate |
| Infrastructure and configuration | Every build | IaC, container and deployment configuration | Trivy plus architecture fitness checks |
| Passive application analysis | Every deployed preview | Browser and API traffic without attack payloads | OWASP ZAP Automation Framework |
| Authenticated journey security | Affected change and release candidate | Role boundaries, sessions, workflows, APIs | Journey runner plus ZAP |
| Active DAST | Scheduled, release, or elevated risk | Explicitly authorized non-production targets | OWASP ZAP Automation Framework |
| Vulnerability regression | Every applicable preview or schedule | Curated, version-pinned templates and known past findings | Nuclei as Trial candidate |
| Adversarial agent and isolation tests | Scheduled and elevated risk | Prompt injection, authority, egress, tenant, budget, sandbox, kill switch | Forge conformance harness |
| Human-led assurance | Risk-based and before material production exposure | Business logic, chained attacks, design and residual risk | Qualified penetration test/red-team provider |

OWASP Web Security Testing Guide scenario identifiers SHOULD provide a versioned coverage taxonomy. Tool output is normalized into Forge findings and Evidence Records; the tool-specific report remains linked provenance rather than the canonical risk decision.

## Why not run every attack after every build

Active scans can mutate state, consume capacity, trigger third parties, create test data, or exploit a real weakness. Running the broadest attack suite after every build would increase blast radius, duration, cost, and false-positive fatigue while providing little additional evidence for unrelated changes.

Forge therefore requires a safe baseline on every applicable build and expands coverage based on impact and risk. Broad active suites run against isolated, resettable targets on schedule or for release candidates. This preserves continuous assurance without pretending all techniques have equal safety or value.

## Autonomous Security Testing Agent

An autonomous security-testing role is an untrusted offensive component. Its model, prompts, browser, scanner, templates, and discovered target content do not define scope. An external deterministic controller enforces:

- exact target hosts, ports, protocols, tenant, paths, accounts, and exclusions;
- permitted test classes and explicitly denied techniques;
- task identity, expiry, action allow-list, egress, rate, concurrency, time, cost, and payload limits;
- data classification, credential handling, evidence redaction, and retention;
- cumulative impact monitoring, hard stops, a human-accessible kill switch, and incident escalation;
- immutable tool and template versions, provenance, and reproducible configuration;
- independent validation, deduplication, confidence, severity, ownership, and exception expiry.

The Forge controller validates scope before every consequential action. The agent runtime receives no credential that can override controller, network, sandbox, or target boundaries. Production and third-party targets are denied by default.

The OWASP Autonomous Penetration Testing Standard is a candidate governance reference for any future autonomous offensive capability. It complements, but does not replace, WSTG testing methodology, DAST tools, qualified human testing, or Forge's own risk and policy contracts.

## Rules of engagement

Every active or autonomous campaign requires a versioned authorization record containing:

- accountable requester, approver, operator, and emergency contact;
- business purpose and linked threat or release risk;
- target ownership proof and environment classification;
- in-scope and out-of-scope assets and data;
- allowed techniques, prohibited actions, schedule, rate, and stop conditions;
- test-account, cleanup, recovery, notification, and incident procedure;
- evidence classification, recipients, retention, and disclosure policy.

Rules of engagement are bound to target and configuration digests. A DNS, ownership, redirect, environment, or scope change invalidates authorization until re-evaluated.

## Finding lifecycle

~~~text
observed -> normalized -> independently validated -> triaged
                                             |          |
                                             v          v
                                        false positive  accepted risk / remediate
                                                            |
                                                            v
                                                     verified closure
~~~

A scanner does not grant severity, release-blocking status, or an exception by itself. Findings link to exact request/response or other evidence, affected artefact and deployment, detection rule, confidence, data exposure, exploit preconditions, owner, remediation, and re-test.

Suppression requires rationale, owner, expiry, and scope. A previous finding becomes a versioned vulnerability-regression scenario where safe and valuable.

## Initial safety and conformance scenarios

- The runner rejects a target outside the signed allow-list, including redirect and DNS rebinding paths.
- The kill switch terminates execution within the declared bound and prevents queued work from starting.
- A scanner or language model cannot enable a prohibited technique through generated configuration.
- Credentials are scoped to synthetic identities and never appear in logs, prompts, screenshots, or reports.
- Active scans cannot address production in the baseline profile.
- Curated template versions and checksums are recorded; unreviewed template updates cannot run automatically.
- A finding cannot be closed by the same agent that produced it without independent evidence.
- The environment resets or reconciles all mutations after testing.
- The final report discloses coverage, exclusions, failures, confidence, and residual risk.

## Adoption sequence

1. Establish per-build secret, dependency, filesystem, image, IaC, SBOM, and provenance checks.
2. Add passive ZAP analysis to an ephemeral preview and normalize its evidence.
3. Add authenticated role-negative journey scenarios.
4. Trial bounded active ZAP and curated Nuclei regression checks on resettable targets.
5. Establish WSTG coverage mapping and human penetration-test intake.
6. Consider an autonomous security-testing agent only after scope enforcement, isolation, kill switch, audit integrity, and finding validation pass independently.

All candidates remain Research until an ADR records versions, licensing, update and template trust, operational ownership, validation evidence, and supported deployment profiles.

## External references

- [OWASP Web Security Testing Guide](https://owasp.org/www-project-web-security-testing-guide/) — versioned web-security coverage taxonomy and methodology.
- [OWASP ZAP Automation Framework](https://www.zaproxy.org/docs/automate/automation-framework/) — candidate passive, authenticated and active DAST orchestration.
- [Nuclei documentation](https://docs.projectdiscovery.io/opensource/nuclei/overview) — candidate curated vulnerability-regression runner.
- [Trivy documentation](https://trivy.dev/docs/latest/guide/) — candidate filesystem, image, dependency and configuration scanner.
- [OWASP Autonomous Penetration Testing Standard](https://owasp.org/APTS/standard/Introduction.html) — candidate governance reference for future autonomous offensive capability; it is not a testing methodology.
