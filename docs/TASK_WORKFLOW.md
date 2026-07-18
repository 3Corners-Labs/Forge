---
title: Engineering Workflow
version: 1.0
status: Proposed
owner: Engineering
---

# Engineering Workflow

```text
Explicit task or RFC
        |
        v
Discover constraints and evidence
        |
        v
Challenge architecture and security
        |
        v
Architecture review when required
        |
        v
Approve scope and acceptance criteria
        |
        v
Implement the smallest approved change
        |
        v
Validate behavior and controls
        |
        v
Update contracts, decisions, and guidance
        |
        v
Human review and merge decision
```

## Architecture review triggers

Review is required before changing:

- trust or data boundaries;
- platform capability ownership;
- deployment units or state ownership;
- identity, policy, audit, or security controls;
- public interfaces and compatibility;
- provider or protocol requirements;
- service objectives or recovery behavior.

## Definition of done

- acceptance criteria are met;
- deterministic validation passes;
- security and privacy impacts are addressed;
- documentation and decisions agree;
- evidence and limitations are reported;
- accountable human approval is recorded.
