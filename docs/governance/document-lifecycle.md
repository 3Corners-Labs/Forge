# Documentation lifecycle

Status: Proposed
Owner: Forge Architecture

## Status values

- Draft: incomplete and non-normative
- Proposed: ready for review and non-normative
- Accepted: approved and normative within its scope
- Superseded: retained for history and linked to its replacement
- Deprecated: valid temporarily with a stated removal path
- Archive: preserved source or historical evidence with no normative authority

Every architecture document declares status and owner. Accepted documents must not contradict higher-precedence material.

## Change rules

- Durable architectural choices are recorded as ADRs.
- Target-state documents state what exists and what does not.
- A document has one canonical location; indexes link rather than copy.
- Plain-text references to repository files must resolve.
- Diagrams are maintained alongside the text they explain.
- Documentation CI checks links, duplicate titles, required metadata, and navigation coverage.

## Separate maturity vocabulary

Document authority SHALL NOT be overloaded with implementation or technology maturity. Capabilities, technologies, models, prompts, agents, tools, and profiles use:

- Research: problem or option is being investigated;
- Trial: bounded implementation or proof is authorized;
- Supported: implementation has a declared support and conformance envelope;
- Reference: implementation has validated portability or interoperability value;
- Deprecated: supported temporarily with a migration path;
- Retired: no longer supported or selectable.

`Accepted` is a document decision status. `Supported` and `Reference` are implementation maturity states. Neither implies the other.

## Evidence labels

Evidence strength is recorded separately:

- Hypothesis: plausible and untested;
- Standard: supported by an applicable published standard;
- Proven: demonstrated in comparable production use;
- Validated: tested in a Forge-controlled evaluation;
- Vision: intended future direction.

Evidence labels do not grant document authority or implementation support.
