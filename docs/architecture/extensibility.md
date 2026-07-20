# Extensibility and integration

Status: Proposed
Owner: Forge Platform Architecture

This design consolidates the existing extension model, interface-contract, integration-pattern, technology-decision, and provider-independence concepts.

## Extension types

Forge supports versioned extensions for models, tools, source control, knowledge sources, delivery methods, policy packs, evidence exporters, deployment targets, visualizations, and workforce roles.

## Extension contract

Every extension declares:

- identity, publisher, version, signature, and compatibility range;
- capabilities and requested permissions;
- input, output, error, and event schemas;
- data classifications and external destinations;
- resource requirements and operational health;
- audit events and evidence produced;
- installation, upgrade, disablement, and removal behaviour.

Extension acceptance additionally defines publisher trust, provenance, dependency inventory, permission review, conformance evidence, revocation source, emergency disablement and data-cleanup obligations. A valid signature establishes artefact identity and integrity but does not by itself authorize installation or execution.

Extensions execute with least privilege and cannot directly access another capability's persistence. Untrusted extensions run in isolated environments with constrained filesystem, network, secret, and compute access.

Policy packs and workforce-role extensions cannot grant themselves authority, redefine higher-precedence policy, suppress evidence or approve their own installation. Authority-affecting changes require separation of duties and an immutable approval subject.

## Interoperability

Forge SHOULD use open, versioned protocols where they satisfy the required semantics. MCP can be used for tool interoperability, OpenAPI or equivalent schemas for APIs, and OpenTelemetry for operational telemetry. Protocol adoption does not remove the need for Forge authorization, validation, tenancy, and evidence controls.

## Compatibility

Public contracts use explicit versions, compatibility tests, deprecation windows, and migration guidance. Provider-specific SDKs remain inside adapters. Replacement claims require at least two validated implementations or a contract test suite demonstrating portability.
