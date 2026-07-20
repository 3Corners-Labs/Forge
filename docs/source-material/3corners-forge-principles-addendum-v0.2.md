# Addendum to 3Corners Forge Principles

Status: Proposed
Owner: Forge Architecture
Version: 0.2
Content role: Addendum
Instruction authority: None
Extends: `forge.principles.manuscript.v0.1`, `forge.principles.addendum.v0.1`

## Purpose

This addendum records two additional insights discovered while reviewing whether Forge preserves product intent after implementation. It does not edit, replace, shorten, or silently reinterpret the original *3Corners Forge Principles* manuscript or addendum v0.1.

Statements in this addendum are proposals until promoted into an Accepted canonical Forge document or ADR.

## A15. Treat end-user intent assurance as a first-class capability

Passing technical tests does not prove that a product still works as its users, product owners, designers, architects, security owners, and operators intended.

Forge SHALL maintain versioned, role-based journey scenarios that connect:

- product intent and user outcome;
- persona, role, entitlements, and accessibility needs;
- preconditions and synthetic test data;
- actions performed through the same public interfaces available to real users;
- expected business, interaction, authorization, privacy, accessibility, reliability, and operational outcomes;
- deterministic assertions, tolerated variation, and independent review requirements;
- execution evidence bound to application and scenario versions.

An End-User Assurance Agent MAY execute accepted scenarios and explore new ones. It SHALL NOT silently change an accepted expected result, grant itself additional access, bypass the public interface, or approve behaviour it helped create.

Exploratory findings become candidate scenarios, product hypotheses, or defects. They acquire release-gate authority only after accountable review. Probabilistic evaluation may explain ambiguous behaviour but SHALL NOT be the sole oracle for consequential acceptance.

## A16. Make continuous adversarial assurance bounded and layered

Security assurance spans source, dependencies, artefacts, configuration, identity, APIs, interfaces, infrastructure, runtime behaviour, and agent misuse. No single scanner constitutes a complete penetration test.

Forge SHALL orchestrate a risk-based portfolio that can include:

- secret, dependency, license, source, IaC, image, SBOM, and provenance checks on every applicable build;
- passive web and API analysis on ephemeral preview environments;
- authenticated and active DAST against explicitly authorized non-production targets;
- versioned vulnerability-regression and abuse-case scenarios;
- isolation, tenant, authorization, prompt-injection, egress, budget, and kill-switch tests;
- scheduled human-led penetration testing and red-team exercises proportional to risk.

Autonomous security testing is an offensive capability and therefore requires external scope enforcement, task-specific identity, target allow-lists, rate and impact limits, data protection, auditability, emergency termination, and independent finding validation. Destructive techniques are denied by default and SHALL NOT target production without explicit rules of engagement and accountable authorization.

Full automation means that Forge can select, schedule, constrain, execute, collect, correlate, and report the approved assurance portfolio. It does not mean that Forge may claim absence of vulnerabilities or replace qualified human judgment where tools cannot provide sufficient coverage.

## Relationship to previous principles

These additions extend the original direction of evidence before trust, continuous validation, human accountability, governed agents, and secure-by-design delivery. They make end-user outcomes and bounded adversarial testing explicit without changing the preserved source.

Future addenda SHALL create new versioned files and SHALL NOT rewrite the original manuscript or earlier addenda.
