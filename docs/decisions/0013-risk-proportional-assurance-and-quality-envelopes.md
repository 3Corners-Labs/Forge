# ADR-0013: Establish risk-proportional assurance and application quality envelopes

Status: Accepted
Owner: Forge Quality, Architecture, and Reliability
Date: 2026-07-22

## Context

Forge must remain responsive and operationally comprehensible while preserving security, privacy, correctness, evidence, and human accountability. It must also help design applications whose performance and scalability fit their actual users and workloads rather than reproducing Forge's own validation assumptions.

The Accepted architecture already requires responsible resource use, measurable capability envelopes, risk-based verification, backpressure, profile-specific capacity, telemetry separation, and evidence before trust. The M1 vertical slice includes bounded latency, concurrency, repository-size, duration, and cost hypotheses. These requirements do not yet define one decision mechanism that connects:

- immediate user-interface feedback and potentially long-running operations;
- application-specific workload and quality requirements;
- module, tool, library, framework, and provider selection;
- applicable security, privacy, quality, and operational guardrails;
- assurance scheduling, parallelism, reuse, and publication eligibility; and
- performance, capacity, compute, cost, sustainability, and human-attention evidence.

Without that mechanism, Forge could run every possible check on every interaction, silently skip slow controls, select technology from vendor benchmarks, over-engineer unknown scale, or postpone performance until replacement is expensive. Counting checks is not a meaningful quality measure: one irrelevant provider call can dominate thousands of cheap deterministic checks, while one omitted authority check can invalidate the complete outcome.

The accountable product direction is:

- every user interaction provides immediate visible feedback, although a query or background operation may take longer to complete;
- publication is blocked whenever an applicable mandatory guardrail is not satisfied;
- background duration and compute are application- and change-specific rather than one global Forge constant; and
- the first Forge Trial addresses one project context at a time, without imposing that scale assumption on applications designed through Forge.

## Decision

Adopt two linked, versioned control artefacts: an **Application Quality Envelope** and a per-change **Assurance Plan**.

### Application Quality Envelope

Every application and every deployment profile seeking Trial or Supported maturity SHALL define an Application Quality Envelope before architecture, technology, or support claims are accepted. The envelope records:

- application purpose, critical journeys, actor classes, and interaction types;
- immediate-feedback, operation-completion, API, background-job, and batch latency objectives;
- expected users, tenants, projects, concurrency, data volume, repository or artefact size, query patterns, peak shape, and growth assumptions;
- availability, consistency, continuity, recovery, cancellation, and degraded-mode behaviour;
- security, privacy, classification, residency, destination, retention, and audit requirements;
- compute, memory, storage, network, token, financial, sustainability, and human-attention budgets;
- applicable mandatory guardrails, advisory checks, publication requirements, and evidence retention; and
- validation workloads, environments, owners, evidence strength, uncertainty, review date, and invalidation triggers.

Unknown requirements remain explicit hypotheses with owners and review dates. Forge SHALL NOT silently convert an unknown scale requirement into either an enterprise-scale topology or an unbounded implementation.

The initial cross-application experience hypotheses are:

- local interaction feedback targets p95 below 100 ms;
- a server-backed interaction targets visible acknowledgement or changed status at p95 below 500 ms;
- an operation expected to exceed one second exposes progress or an explicit waiting state; and
- a long-running operation declares timeout, cancellation, recovery, and completion behaviour.

These are design hypotheses, not universal completion SLOs or production commitments. Each application validates or replaces them in its own envelope. Immediate feedback never means that a query, build, analysis, or verification must complete synchronously.

The first Forge Trial remains bounded to one active project context at a time. Existing M1 validation hypotheses for delivery runs, workers, repository size, and duration remain applicable within that project context. They do not become defaults for applications built with Forge.

### Proven performance as an eligibility gate

A new material module, tool, library, framework, provider, or technology selection or maturity promotion beyond Research MAY occur only when reproducible evidence demonstrates that it meets the applicable latency, throughput, capacity, resource, reliability, and cost envelope without weakening required security, privacy, correctness, isolation, portability, recovery, or operability controls. Material means that the choice can affect an envelope objective, trust boundary, deployment unit, external dependency, or significant operating cost.

Required controls are eligibility gates, not weighted criteria that a faster option can compensate for. Among eligible options, Forge compares measured performance, resource efficiency, cost, operational simplicity, portability, maintainability, and user experience.

For Trial, performance evidence SHALL be Validated in a Forge-controlled representative evaluation. Comparable production evidence labelled Proven MAY strengthen confidence but does not replace application-specific validation. Supported maturity requires validation under the named deployment profile and published capacity envelope.

Existing bounded M1 Trial components are not retroactively represented as performance-validated. They SHALL produce evidence for the applicable M1 validation envelope before M1 exit, renewed selection, wider Trial use, or any Supported claim. Until then the missing evidence remains an explicit limitation.

Vendor claims and generic benchmarks are insufficient. Selection evidence records:

- workload, dataset, scale, hardware, topology, versions, configuration, warm-up, and test method;
- p50, p95, and p99 latency, throughput, concurrency, saturation point, and error rate;
- CPU, memory, storage, network, token, energy indicator, and financial cost where applicable;
- cold-start, peak-load, noisy-neighbour, backpressure, timeout, failure, and recovery behaviour;
- raw or reproducibly derived results, uncertainty, owner, evidence strength, and date; and
- invalidation triggers such as version, configuration, workload, profile, topology, or scale changes.

Evidence depth is proportional to impact. A small deterministic in-process module can demonstrate eligibility through complexity analysis, focused benchmarks, resource bounds, and regression tests. A shared datastore, runtime, external provider, framework, worker, or deployment-boundary choice requires representative load, failure, recovery, and operating-cost evidence. This proportionality does not waive an applicable envelope objective.

A module is extracted into a separate deployment unit only after measured evidence meets the extraction criteria in the Deployment architecture and shows that the benefit justifies added operational cost. Hypothetical future scale is not sufficient.

### Risk-proportional Assurance Plan

Every material change receives a versioned Assurance Plan derived from the Application Quality Envelope, affected capabilities and dependencies, data classification, destinations, authority, side effects, change type, and risk. The candidate, generated code, model, tool, or test generator SHALL NOT be the sole selector or oracle for its own assurance.

Every registered control declares:

- owner, purpose, protected outcome, and covered threat, quality attribute, or journey;
- deterministic applicability conditions and trigger;
- mandatory, conditional, advisory, exploratory, or scheduled role;
- input digests, policy and tool versions, execution environment, and evidence output;
- latency, compute, concurrency, financial, and human-review budget;
- timeout, cancellation, retry, backpressure, and unavailable-result semantics; and
- validity window, reuse conditions, invalidation conditions, and retention class.

Assurance is staged:

| Stage | Purpose | Critical-path rule |
|---|---|---|
| Interaction | Identity, scope, schema, policy, admission, and immediate input safety | Cheap deterministic controls run synchronously and return visible feedback |
| Pre-operation | Authority, classification, limits, isolation, source binding, and prerequisites | The protected operation cannot start until applicable controls pass |
| Candidate verification | Relevant contract, unit, integration, journey, security, privacy, performance, and recovery evidence | Independent controls run as early and parallel as budgets permit; required results block publication |
| Scheduled assurance | Full regression, capacity, endurance, active security, portability, and recovery exercises | Runs outside ordinary interaction paths in an authorized environment; policy defines promotion impact |
| Publication | Candidate digest, applicable evidence, approval, destination, and evidence readiness | Reconciles existing results; it does not postpone all expensive checks until the publication request |

An applicable mandatory guardrail that fails, is missing, expired, stale, bound to different inputs, or returns indeterminate SHALL block the protected operation. A non-applicable guardrail requires a deterministic explanation. Advisory findings remain visible but do not become hidden blockers.

Evidence MAY be reused only when its candidate or input digest, policy version, tool version, configuration, environment, validity window, and applicability facts remain exact. Any authority- or outcome-relevant change invalidates affected evidence. Authorization, tenant isolation, secret access, integrity, and publication-approval checks are never probabilistically sampled away.

Independent controls SHOULD execute in parallel within declared concurrency and resource budgets. Queues, admission control, backpressure, bulkheads, circuit breakers, bounded retries, and cancellation prevent assurance from exhausting the application it protects. Telemetry MAY be sampled and aggregated, but telemetry failure SHALL follow the separate failure semantics in Evidence and Observability architecture.

### Capability ownership

- `forge.capability.intent` owns the application outcomes, critical user journeys, and accountable quality expectations supplied by product stakeholders.
- `forge.capability.architecture` owns the integrated, versioned Application Quality Envelope and traceability from those expectations to architecture constraints and selection decisions.
- `forge.capability.verification` owns the derived Assurance Plan, independent control selection, execution-result evaluation, and verification findings.
- Journey Assurance and Continuous Security Testing own their governed scenario and safety contracts; Verification references rather than redefines them.
- `forge.capability.policy` evaluates whether the protected action is authorized from trusted facts, required obligations, and verification state. It does not select its own tests or own the assurance catalogue.
- `forge.capability.telemetry` measures health, performance, capacity, and cost without becoming audit evidence by default.
- `forge.capability.reliability` owns profile SLO, capacity, continuity, and recovery operation once the applicable envelope is adopted.

### Generated application outcomes

When Forge proposes or prepares an application architecture, module, tool, or technology choice, the reviewable result SHALL include:

- the applicable Application Quality Envelope and unresolved hypotheses;
- alternatives considered and their eligibility evidence;
- benchmark and load-test definitions suitable for independent reproduction;
- the resulting Assurance Plan and publication-blocking guardrails;
- observed performance, capacity, resource, cost, security, privacy, and recovery evidence; and
- limitations, residual risks, scale boundaries, monitoring, rollback, and revalidation triggers.

Performance regressions are evaluated against the declared envelope and risk, not an undocumented global threshold. Forge SHALL explain whether a change remains inside the envelope and SHALL NOT describe an application or component as scalable without a measured range and reproducible evidence.

This decision defines architecture and evidence requirements only. It does not authorize a provider, production profile, assurance engine, automatic publication, service extraction, or relaxation of an existing guardrail.

## Alternatives considered

- **Run every available check on every interaction:** rejected because relevance, latency, cost, operational stability, and human attention would degrade without increasing proportional confidence.
- **Optimize for speed and treat safety checks as weighted trade-offs:** rejected because required security, privacy, correctness, authority, and evidence controls are eligibility gates.
- **Use one universal scale and performance profile:** rejected because Forge and the applications it builds have different users, workloads, data, risks, and deployment profiles.
- **Accept vendor benchmarks as proof:** rejected because they rarely reproduce the applicable workload, configuration, tail latency, failure mode, or total operating cost.
- **Defer performance and scalability until production:** rejected because architecture, data model, interaction design, and technology choices may already be costly to reverse.
- **Extract services pre-emptively for future scale:** rejected because the modular-monolith strategy requires evidence before accepting distributed-system cost.
- **Let each candidate choose its own minimal checks:** rejected because candidate-controlled assurance creates circular trust.

## Consequences

- User interactions remain responsive through immediate feedback, staged work, progress, and cancellation rather than synchronous completion of every operation.
- Publication remains fail-closed for every applicable mandatory guardrail without forcing every possible control onto every change.
- Application design begins with explicit quality and workload hypotheses instead of generic enterprise-scale assumptions.
- Technology selection requires reproducible, workload-relevant evidence and can be revisited when that evidence expires.
- Parallelism and evidence reuse can reduce elapsed time, but require deterministic applicability, cache invalidation, lineage, and resource controls.
- Teams incur earlier benchmarking, load-test, telemetry, and envelope-maintenance effort; this is accepted to reduce late architectural replacement and unsupported scale claims.
- Safety can still increase latency where risk requires it. Forge exposes that cost and seeks architectural improvement rather than silently skipping the control.
- Different applications can make different performance and topology choices while sharing Forge capability, trust, and evidence semantics.

## Validation

Implementation and conformance work SHALL demonstrate:

- visible feedback within each accepted interaction hypothesis and explicit progress for longer operations;
- independent, deterministic Assurance Plan selection from changed capabilities and risk;
- publication denial for failed, missing, stale, indeterminate, or mismatched applicable guardrails;
- explicit non-applicability and visible advisory outcomes without hidden bypasses;
- digest- and version-bound evidence reuse plus invalidation after a relevant change;
- representative module and tool comparisons with reproducible p50, p95, p99, saturation, resource, cost, and recovery results;
- bounded parallelism, admission control, cancellation, and backpressure under assurance load;
- tenant and workload isolation under declared concurrency and noisy-neighbour scenarios;
- no sampling of authority, isolation, secret, integrity, or exact-approval controls; and
- revalidation after declared version, configuration, workload, profile, topology, or scale changes.

Architecture fitness SHALL require lifecycle metadata, owners, applicable quality envelopes, evidence references, and explicit maturity claims when these artefacts become executable. Each supported profile publishes its validation range and SLOs. Review compares measured outcomes with the applicable envelope rather than counting controls.

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-22, acting as the temporary accountable Product, Architecture, Security and Privacy, Quality, Operations, and Engineering reviewer established for bounded M1. The decision becomes operationally normative when it reaches protected `main` through the required review path.

The interaction targets remain validation hypotheses rather than support commitments. Implementing an Assurance Plan engine, Application Quality Envelope contract, provider selection, production profile, or relaxed guardrail requires the applicable separately authorized iteration and evidence.

## Supersession

This decision specializes the responsible-resource, continuous-assurance, capability-envelope, deployment, and observability requirements already Accepted. It does not supersede an existing ADR or weaken any accepted guardrail.
