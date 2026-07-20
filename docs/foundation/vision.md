# Forge product vision

Status: Proposed
Owner: Forge Product and Architecture

## Purpose

Forge is an AI-native application development team implemented as an application and platform. It turns human intent into secure, operable, understandable software and application development environments.

Forge is the team. It coordinates a virtual workforce across the complete software lifecycle: discovery, product management, architecture, user research, experience design, engineering, verification, security, release, operations, support, and continuous improvement.

The goal is not to remove human judgment. The goal is to automate repeatable work, expose consequential decisions, and give people effective control over outcomes.

## Product promise

A user can describe an outcome in natural language and collaborate with Forge to:

- clarify the problem and intended users;
- explore alternatives and trade-offs;
- produce requirements, architecture, journeys, prototypes, code, tests, and operational assets;
- inspect how the system works before approving it;
- deploy to Forge SaaS or a self-hosted environment;
- operate and evolve the result with traceable evidence.

## Intent-driven development

Forge supports the immediacy commonly associated with vibe coding, but makes it dependable. An informal prompt is the start of a structured conversation, not an authorization to ship opaque code.

Forge progressively converts intent into linked, versioned artefacts:

~~~text
intent -> outcomes -> requirements -> decisions -> design -> code -> evidence -> release -> learning
~~~

Each transition remains understandable, reversible where practical, and governed according to risk.

## Experience principles

Forge SHALL:

- meet users at their level of expertise;
- show what was understood, assumed, changed, and verified;
- provide visual explanations of architecture, flows, data, and runtime behaviour;
- let experts inspect and override details within policy;
- offer safe defaults without concealing important trade-offs;
- support collaboration between humans and specialist agents;
- keep deployment ownership portable between SaaS and self-hosted operation.

## Boundaries

Forge is not an autonomous authority for production risk. Humans or delegated organizational controls remain accountable for policy, budget, material architectural decisions, and production release.

Forge does not promise one universal development method. It implements composable delivery controls that can express methods such as continuous delivery, trunk-based development, Scrum, Kanban, DevSecOps, platform engineering, SRE, and regulated stage gates.
