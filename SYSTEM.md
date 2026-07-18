---
title: AI Engineering Instructions
version: 1.0
status: Proposed
---

# AI Engineering Instructions

These instructions define how AI-assisted engineering work is performed in Forge. They do not define the platform architecture.

## Required workflow

1. Confirm the task and its authority.
2. Read the canonical constitution, relevant ADRs, and affected contracts.
3. Distinguish accepted facts from proposals and assumptions.
4. Challenge contradictions, boundary violations, and security risks.
5. Obtain architecture review before changing architecture, security, deployment boundaries, or public contracts.
6. Implement the smallest approved change.
7. Validate behavior with deterministic evidence.
8. Update affected documentation and decisions.
9. Report limitations and unresolved risks.

## Invariants

- Product applications do not own reusable platform governance capabilities.
- Authentication and authorization remain separate responsibilities.
- Policy decisions are independent from execution technology.
- External side effects require explicit authorization and audit evidence.
- Provider and protocol products remain replaceable behind Forge-owned contracts.
- Runtime state has an explicit owner, lifecycle, and retention policy.
- Documentation never claims implementation or validation without evidence.
- Humans approve consequential architectural, security, and release decisions.

## Stop conditions

Pause implementation when:

- normative documents conflict;
- the security or data boundary is unknown;
- a consequential action lacks explicit authority;
- the required contract or acceptance criterion is absent;
- validation would require weakening an existing control.
