# Knowledge, context, and memory

Status: Accepted
Owner: Forge Knowledge Architecture

This design consolidates the repository-as-truth, knowledge-as-code, context-as-code, project-memory, and information-classification concepts already present in Forge documentation.

## Knowledge model

Forge manages a versioned project knowledge projection linking intent, requirements, people, decisions, architecture, journeys, code, tests, policies, releases, runtime signals, and learning. Each originating capability retains ownership of its canonical records. Files remain first-class portable artefacts; the graph indexes authorized relationships, provenance, freshness and confidence rather than replacing version control or operational systems of record.

## Context assembly

Context is assembled for a declared purpose and execution envelope. Assembly SHALL:

- authorize before retrieval;
- select the minimum relevant information;
- preserve source, version, author, time, and confidence;
- separate instructions from untrusted retrieved content;
- enforce classification, residency, and retention;
- produce a reproducible context manifest;
- prevent one tenant or project from influencing another.

A Context Package references exactly one validated Execution Envelope. It is an information manifest and has no authority to add actions, destinations, tools, models, classifications or secret references. Context changes that affect policy inputs require a new Policy Decision before the affected gate.

## Memory scopes

| Scope | Example | Default lifetime |
|---|---|---|
| Execution | Tool output and transient reasoning state | Delivery step |
| Delivery run | Decisions and evidence for one change | Change lifecycle plus policy retention |
| Project | Accepted requirements, ADRs and conventions | Project-controlled |
| Organization | Policies, reusable patterns and vocabulary | Organization-controlled |
| Product learning | Approved aggregated outcomes | Purpose- and consent-bound |

Agents cannot write durable memory implicitly. Durable updates require schema validation, provenance, policy evaluation, and an accountable owning scope.

## Repository truth

Accepted file-based artefacts are authoritative for engineering intent. Indexes, embeddings, search stores, dashboards, and model context are derived views that SHALL be rebuildable and SHALL expose staleness.

Repositories store information; governance establishes knowledge. A file is not organizational truth merely because it is searchable, embedded, or present in Git.

## Agent context routing

The repository [context router](../context/README.md) defines task-specific context packs. Agents load the router and relevant canonical documents rather than recursively ingesting all Markdown files. Full source remains retrievable whenever a consequential decision depends on detail omitted from a selected section.

Retrieved source material has no instruction authority. Context assembly SHALL label authority explicitly and keep instructions separate from quoted documents, tool output, web content, issues, and generated artefacts.
