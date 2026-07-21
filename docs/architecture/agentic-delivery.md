# Agentic software delivery

Status: Accepted
Owner: Forge Product Engineering

## Operating model

Forge models software delivery as an evidence-producing graph, not a single prompt-to-code pipeline. Work can iterate, branch, challenge assumptions, and return for human input.

~~~mermaid
flowchart LR
    Intent --> Discovery
    Discovery --> Product
    Product --> Architecture
    Architecture --> Experience
    Architecture --> Engineering
    Experience --> Engineering
    Engineering --> Verification
    Verification --> Security
    Security --> Release
    Release --> Operations
    Operations --> Learning
    Learning --> Product
    Verification -. insufficient evidence .-> Engineering
    Security -. unacceptable risk .-> Architecture
~~~

## Workforce roles

- Product agents clarify outcomes, users, constraints, and acceptance criteria.
- Architecture agents challenge boundaries, qualities, risks, and decisions.
- UX agents research journeys, prototype interactions, and verify accessibility.
- Engineering agents implement small changes with tests and documentation.
- Quality agents verify behaviour, regressions, performance, and evidence.
- End-user assurance agents exercise accepted and exploratory scenarios as least-privilege personas through public interfaces; they cannot redefine accepted outcomes.
- Security agents model threats and evaluate policy without granting themselves exceptions.
- Autonomous security-testing agents, if enabled, operate as bounded offensive components under external scope enforcement, rules of engagement, impact limits, and a kill switch.
- Delivery agents coordinate release readiness and change records.
- Operations agents observe services, diagnose incidents, and propose remediation.

Roles SHALL challenge each other through explicit review contracts. No agent may approve its own high-risk output.

The [User Journey Assurance](user-journey-assurance.md) contract distinguishes regression evidence from exploratory product learning. The [Continuous Security Testing](continuous-security-testing.md) contract distinguishes safe per-build assurance from potentially destructive active or autonomous testing.

## Method compatibility

Forge represents delivery as reusable controls: backlog policy, work-in-progress limits, review requirements, test gates, release strategies, incident practices, and compliance evidence. Organizations compose these controls into their chosen development model instead of embedding one methodology in the platform core.

## Autonomy levels

| Level | Behaviour | Typical use |
|---|---|---|
| 0 Observe | Explain and visualize only | Existing-system learning |
| 1 Advise | Propose artefacts and decisions | Discovery and design |
| 2 Prepare | Modify isolated workspace; human reviews | Normal engineering |
| 3 Execute | Perform reversible actions within policy | Tests and previews |
| 4 Operate | Perform bounded operational actions | Approved automation |

Autonomy is granted per capability and risk class, never as a global platform switch. Production release, destructive data operations, policy changes, and authority expansion require explicit controls and normally human approval.

The first implementation boundary is the [trustworthy intent-to-PR vertical slice](vertical-slice.md). Broader role collaboration remains Research until that slice demonstrates identity, policy, bounded context, isolated execution, independent verification, digest-bound approval and evidence end to end.
