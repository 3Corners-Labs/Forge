Ik zou niet beginnen met een voorwoord. Ook niet met “What is Forge?”.

Ik zou de lezer vanaf de eerste pagina laten voelen dat er een fundamenteel probleem is dat vrijwel iedereen in AI negeert.

Zo zou ik hoofdstuk 1 openen.

⸻

Part I — Why Forge Exists

Chapter 1

The AI Engineering Problem

Software engineering learned long ago that complexity cannot be eliminated. It can only be managed.

Artificial Intelligence has reintroduced complexity at a scale that traditional software engineering was never designed to handle.

Forge exists because AI systems deserve engineering discipline.

⸻

1.1 The Second Software Revolution

Software engineering has experienced several fundamental shifts.

The first revolution separated hardware from software.

The second introduced structured programming.

The third embraced object-oriented design.

The fourth established cloud-native architectures.

Artificial Intelligence represents the next revolution—not because computers have become intelligent, but because software has become probabilistic.

For the first time in modern software engineering, developers routinely build systems whose outputs cannot be predicted with complete certainty.

This changes everything.

Traditional software systems execute deterministic algorithms.

Given the same inputs, they produce the same outputs.

Large Language Models do not.

Instead, they generate statistically probable responses based on learned patterns.

This difference is not merely technical.

It fundamentally changes how software must be designed, tested, secured, governed and maintained.

⸻

1.2 The Industry Response

The software industry reacted quickly.

Thousands of projects appeared almost overnight.

Frameworks.

Prompt libraries.

Agent platforms.

Model routers.

Workflow engines.

Retrieval systems.

Memory layers.

Most solved an immediate technical problem.

Very few addressed the engineering problem.

As a result, many AI applications evolved into collections of prompts, configuration files and framework-specific logic that work remarkably well during demonstrations but become increasingly difficult to understand, operate and evolve.

The challenge is no longer building AI applications.

The challenge is engineering AI systems that remain understandable five years after they were created.

⸻

1.3 The Hidden Complexity

Artificial Intelligence introduces entirely new dimensions of complexity.

Traditional software engineers think in terms of APIs, databases, message queues and deployment pipelines.

AI engineers must additionally reason about:

* context construction;
* prompt composition;
* knowledge retrieval;
* model selection;
* token economics;
* probabilistic behaviour;
* policy enforcement;
* tool orchestration;
* conversational state;
* memory management;
* evidence generation;
* human oversight.

None of these concerns exist in isolation.

Every request traverses multiple decision layers before a model generates a single token.

The complexity is architectural rather than algorithmic.

⸻

Figure 1 — Complexity Explosion

flowchart LR
Request
Request --> Identity
Identity --> Policy
Policy --> Context
Context --> Memory
Memory --> Knowledge
Knowledge --> Runtime
Runtime --> Tools
Tools --> Providers
Providers --> Response
Audit -.-> Identity
Audit -.-> Policy
Audit -.-> Context
Audit -.-> Runtime
Audit -.-> Response
Observability -.-> Runtime
Observability -.-> Providers

The language model represents only one component in a much larger engineering system.

Treating the model as the system is the architectural mistake that Forge seeks to correct.

⸻

1.4 The Framework Trap

Most AI platforms describe themselves using implementation terminology.

They are:

* orchestration frameworks;
* agent frameworks;
* prompt frameworks;
* workflow engines;
* model gateways.

While technically accurate, these descriptions reveal a deeper assumption:

The framework is considered the centre of the architecture.

Forge rejects this assumption.

Applications should not be designed around frameworks.

Frameworks should be replaceable.

Engineering principles should remain.

Architecture should survive technology.

The platform should outlive every dependency it contains.

⸻

1.5 Engineering Before Intelligence

Forge adopts a deliberately different starting point.

The primary engineering problem is not intelligence.

It is complexity.

Intelligence without engineering becomes experimentation.

Engineering without intelligence becomes automation.

Sustainable AI systems require both.

This distinction leads to the central thesis of this document.

Artificial Intelligence is not the product.

Engineering Intelligence is the product.

⸻

1.6 The Core Thesis

Forge is founded on one observation:

Organizations do not struggle because they lack language models.

Organizations struggle because they lack a repeatable engineering model for building AI systems.

Models will continue to evolve.

Providers will change.

Frameworks will disappear.

New orchestration engines will emerge.

Yet the engineering disciplines required to govern, understand, secure and evolve AI systems will remain.

Forge therefore treats language models as replaceable implementation details.

The enduring value lies elsewhere:

* Architecture.
* Context.
* Knowledge.
* Memory.
* Governance.
* Evidence.
* Policy.
* Engineering.

Everything else is infrastructure.

⸻

Closing Thought

The history of software engineering demonstrates that every technological revolution is eventually followed by an engineering revolution.

Artificial Intelligence has already transformed what computers can do.

The next transformation will determine how responsibly, predictably and sustainably those capabilities can be engineered.

Forge is an attempt to define that engineering discipline.

⸻
Mooi. Dit is precies het hoofdstuk waar Forge echt een eigen identiteit krijgt.

Hoofdstuk 1 beschreef het probleem.

Hoofdstuk 2 mag niet direct de oplossing beschrijven. Het moet eerst een nieuw denkkader introduceren. Ik wil dat de lezer na dit hoofdstuk nooit meer over AI-platformen denkt zoals daarvoor.

⸻

Chapter 2

Engineering Operating Systems

Every major advance in computing has been accompanied by a new operating model.

Artificial Intelligence is no different.

Forge is not another AI framework. Forge is an Engineering Operating System.

⸻

2.1 Beyond Frameworks

Software engineers naturally reach for frameworks.

Web applications have web frameworks.

Mobile applications have mobile frameworks.

Data engineering has data platforms.

Machine learning has ML frameworks.

When Large Language Models became widely available, the industry responded in exactly the same way.

Agent frameworks.

Prompt frameworks.

Workflow frameworks.

Orchestration frameworks.

Model gateways.

This response was predictable.

It was also incomplete.

Frameworks solve implementation problems.

They do not solve engineering problems.

As organizations began deploying AI beyond prototypes, a different class of challenges emerged.

Not how to call a model.

But how to operate an ecosystem of models, policies, knowledge, tools and governance across hundreds of applications.

The challenge was no longer orchestration.

The challenge was engineering.

⸻

2.2 A Historical Perspective

Operating systems did not emerge because CPUs became faster.

They emerged because hardware became too complex for every application to manage independently.

Before operating systems, every application managed:

* memory;
* storage;
* devices;
* scheduling;
* interrupts.

Eventually engineers realized that these concerns should become shared platform capabilities.

Applications should focus on business logic.

The operating system should manage infrastructure.

Artificial Intelligence is approaching the same inflection point.

Today, many AI applications independently implement:

* prompt construction;
* context retrieval;
* conversation memory;
* provider selection;
* retries;
* tool invocation;
* security controls;
* audit logging;
* cost management;
* observability.

Each application rebuilds the same capabilities.

Each application solves them differently.

This is architectural duplication.

⸻

Figure 2.1 — From AI Frameworks to an Engineering Operating System

flowchart LR
subgraph Traditional
A1["Application A"]
A2["Application B"]
A3["Application C"]
A1 --> F1["Framework"]
A2 --> F2["Framework"]
A3 --> F3["Framework"]
end
subgraph Forge
B1["Application A"]
B2["Application B"]
B3["Application C"]
Platform["Forge Engineering Operating System"]
B1 --> Platform
B2 --> Platform
B3 --> Platform
end

The objective is not to centralize applications.

The objective is to centralize engineering capabilities.

⸻

2.3 Defining an Engineering Operating System

An Engineering Operating System is a platform that provides reusable engineering capabilities independently of any individual application.

Unlike traditional operating systems, Forge does not abstract hardware.

It abstracts engineering complexity.

Applications delegate responsibilities that should never have been implemented repeatedly.

Those responsibilities include:

* identity;
* authorization;
* context construction;
* knowledge retrieval;
* memory management;
* provider abstraction;
* policy evaluation;
* audit generation;
* observability;
* governance.

Applications retain ownership of business logic.

Forge owns the engineering infrastructure.

⸻

2.4 The Capability Model

Forge is intentionally organized around capabilities rather than technologies.

Technologies change.

Capabilities endure.

A capability represents a stable engineering responsibility.

For example:

Capability	Responsibility
Identity	Establish trust
Policy	Decide what is permitted
Context	Build execution context
Memory	Preserve state and knowledge
Runtime	Coordinate execution
Tooling	Access external systems
Audit	Produce evidence
Observability	Explain system behaviour

These capabilities define the platform.

LiteLLM, LangGraph and MCP merely implement portions of that platform today.

They are replaceable.

The capability model is not.

⸻

2.5 Architectural Inversion

Most AI systems are designed from the model outward.

The language model becomes the centre of the architecture.

Everything else is attached around it.

Forge deliberately reverses this relationship.

The language model occupies the outer edge of the system.

The platform itself becomes the architectural centre.

flowchart TB
subgraph Forge Platform
Identity
Policy
Context
Memory
Runtime
Audit
Observability
end
Runtime --> LiteLLM
LiteLLM --> Provider["Language Model"]
Runtime --> MCP
MCP --> Enterprise["Enterprise Systems"]

The model is a dependency.

Not the architecture.

⸻

2.6 Stable Capabilities, Replaceable Technology

History suggests that every technology stack eventually changes.

Programming languages evolve.

Databases change.

Cloud providers mature.

Frameworks rise and fall.

Large Language Models will follow the same pattern.

An architecture that depends on a particular provider, orchestration framework or prompting strategy will eventually require substantial redesign.

Forge therefore treats technology as an implementation concern.

Engineering capabilities remain stable even when implementations change.

This separation protects long-term investments.

⸻

2.7 The Platform Contract

Every capability within Forge SHALL satisfy the same architectural contract.

A capability SHALL be:

* Independently deployable.
* Independently observable.
* Independently testable.
* Independently replaceable.
* Independently evolvable.
* Accessible through documented interfaces.
* Governed through policy.
* Auditable by default.

These characteristics establish consistency across the entire platform.

⸻

2.8 The Central Principle

Forge can ultimately be summarized in a single architectural principle:

Applications should solve business problems.

Forge should solve AI engineering problems.

Everything else follows from that decision.

⸻

Closing Thought

The software industry did not advance because every application became better at managing memory, storage and scheduling.

It advanced because those concerns became shared platform capabilities.

Artificial Intelligence now faces the same transition.

The future does not belong to applications that embed increasingly complex AI frameworks.

It belongs to platforms that make sophisticated AI engineering appear deceptively simple.

Forge is an attempt to define that platform.

⸻

Waarom ik dit hoofdstuk belangrijk vind

Ik denk dat dit hoofdstuk uiteindelijk het “manifest” van Forge wordt. Als hoofdstuk 1 uitlegt waarom de huidige AI-aanpak tekortschiet, legt hoofdstuk 2 uit hoe we er fundamenteel anders naar kijken. Vanaf hoofdstuk 3 kunnen we dan de visie en de ontwerpfilosofie uitwerken, omdat de lezer eerst het onderliggende paradigma begrijpt. Dat geeft de rest van de Architecture Bible een veel steviger fundament.
Mooi. Nu de lezer begrijpt waarom AI een engineeringprobleem is en waarom Forge een Engineering Operating System is, kunnen we de visie neerzetten.

Dit hoofdstuk moet niet aanvoelen als marketing. Het moet een architectonisch kompas zijn waar iedere latere beslissing aan getoetst kan worden.

⸻

Chapter 3

Vision

Vision is not a prediction of the future.

It is a decision about the future we choose to build.

⸻

3.1 Why Forge Exists

Every generation of software engineering has solved the limitations of the generation before it.

Structured programming reduced uncontrolled complexity.

Object-oriented design improved modularity.

Cloud computing abstracted infrastructure.

DevOps integrated development and operations.

Artificial Intelligence introduces a new class of complexity that cannot be addressed by existing engineering practices alone.

Organizations increasingly rely on systems that:

* reason instead of calculate;
* generate instead of retrieve;
* infer instead of execute predefined logic;
* evolve continuously rather than through discrete releases.

Traditional software engineering provides few answers for governing systems whose behavior is partially probabilistic.

Forge exists to close that gap.

Its purpose is not to make AI more intelligent.

Its purpose is to make AI systems more engineerable.

⸻

3.2 The Future We Believe In

Forge is built upon several long-term assumptions.

Intelligence will become ubiquitous.

Language models will eventually become as common as relational databases or HTTP servers.

Selecting an AI provider will become an implementation decision rather than a strategic one.

The competitive advantage will no longer come from access to intelligence.

It will come from engineering that intelligence effectively.

⸻

Models will become commodities.

Today’s frontier models will eventually be replaced.

History has repeatedly shown that technologies mature, standardize and commoditize.

Organizations that tightly couple their architecture to specific models or vendors will continuously pay the cost of migration.

Organizations that engineer around stable capabilities will evolve with significantly less friction.

⸻

Engineering will become the differentiator.

As intelligence becomes inexpensive, engineering becomes valuable.

Organizations will compete through:

* trusted context;
* organizational knowledge;
* governance;
* security;
* operational excellence;
* engineering discipline.

The quality of an AI system will increasingly depend on everything surrounding the model.

⸻

3.3 A New Engineering Discipline

Software engineering emerged because programming alone proved insufficient.

Artificial Intelligence is reaching a similar inflection point.

Prompt engineering is valuable.

Model engineering is valuable.

But neither provides a complete discipline for building enterprise-grade AI systems.

Forge proposes a broader field.

AI Engineering is the discipline of designing, governing, operating and evolving intelligent systems throughout their entire lifecycle.

This discipline extends beyond models.

It encompasses architecture, policy, context, memory, economics, observability and human accountability.

⸻

Figure 3.1 — AI Engineering as a System

flowchart TB
subgraph AI Engineering
Architecture
Context
Knowledge
Memory
Policy
Identity
Security
Runtime
Observability
Economics
Governance
HumanOversight
end
AIEngineering --> LanguageModel["Language Models"]

The language model is only one component within a much larger engineering ecosystem.

⸻

3.4 Our Vision

We envision a future where AI systems are engineered with the same rigor as mission-critical software.

Every intelligent system should be:

* understandable;
* testable;
* observable;
* secure;
* explainable;
* governed;
* economically sustainable;
* continuously evolvable.

These characteristics should be inherent to the platform rather than optional additions implemented by individual applications.

⸻

3.5 Engineering Before Features

Many platforms compete by adding more features.

Forge competes by reducing complexity.

Every capability introduced into the platform must satisfy a simple question:

Does this reduce engineering complexity for applications?

If the answer is no, the capability does not belong in Forge.

The goal is not to create the largest platform.

The goal is to create the most coherent one.

⸻

3.6 Platform Thinking

Applications are temporary.

Platforms endure.

Individual business applications may exist for a few years.

The engineering platform supporting them should evolve for decades.

This distinction influences every architectural decision.

Forge optimizes for long-term platform evolution rather than short-term feature delivery.

⸻

3.7 Human Accountability

Artificial Intelligence must never obscure responsibility.

Models may assist.

Agents may automate.

Workflows may orchestrate.

Only humans remain accountable.

Every architectural decision within Forge reflects this principle.

The platform is designed to amplify human decision-making rather than replace it.

Automation without accountability is not engineering.

It is delegation without governance.

⸻

3.8 Success Criteria

Forge will be successful when organizations no longer think about AI infrastructure.

Instead, they focus entirely on solving business problems.

Applications should consume engineering capabilities in the same way they consume networking, storage or authentication today.

AI engineering should become invisible.

Its quality should be measured not by the sophistication of its implementation, but by the simplicity it enables.

⸻

3.9 The Long View

Forge is not being designed for the current generation of language models.

It is being designed for the next decade of intelligent systems.

The platform should remain relevant even if:

* today’s frontier models disappear;
* orchestration frameworks are replaced;
* new reasoning architectures emerge;
* AI hardware changes fundamentally;
* autonomous systems become commonplace.

Architecture should outlive technology.

Principles should outlive architecture.

⸻

Closing Thought

Every enduring engineering discipline is founded on a simple idea.

Civil engineering is not about concrete.

Electrical engineering is not about copper.

Software engineering is not about programming languages.

Likewise, AI Engineering is not about language models.

It is about designing systems in which intelligence can be trusted, governed, understood and continuously improved.

Forge exists to help define that discipline.

⸻

Mijn voorstel voor hoofdstuk 4

Nu we het probleem (hoofdstuk 1), het paradigma (hoofdstuk 2) en de visie (hoofdstuk 3) hebben neergezet, wil ik hoofdstuk 4 uitzonderlijk sterk maken. Dat wordt niet zomaar een lijstje ontwerpprincipes, maar een Design Philosophy waarin we uitspraken formuleren die richting geven aan álle toekomstige architectuur.

Voorbeelden daarvan zijn:

* Complexity belongs inside the platform.
* Context is engineered, not prompted.
* Knowledge is curated, not collected.
* Memory is intentional, not incidental.
* Capabilities outlive technologies.
* Evidence is a first-class artifact.
* Engineering is the product; AI is an implementation detail.

Ik verwacht dat dit hoofdstuk uiteindelijk één van de meest geciteerde delen van de Architecture Bible zal worden, omdat het de filosofische basis vormt voor alle technische keuzes die daarna volgen.

Perfect. Dit is het hoofdstuk waar Forge zijn DNA krijgt. Als iemand over tien jaar vraagt “waar staat Forge voor?”, dan moet dit hoofdstuk het antwoord zijn.

Ik wil dit niet schrijven als “10 design principles”. Ik wil het schrijven als een verzameling engineering-axioma’s. Dingen waarvan je later kunt zeggen: deze beslissing klopt, want hij volgt uit de filosofie van Forge.

⸻

Chapter 4

Design Philosophy

Architecture is the consequence of philosophy.

Every engineering decision embodies a belief about how systems should be built.

Forge is built upon explicit beliefs.

⸻

4.1 Philosophy Before Technology

Technologies evolve.

Frameworks appear and disappear.

Programming languages change.

Cloud providers compete.

Language models improve.

Engineering philosophy changes far more slowly.

Forge deliberately separates enduring engineering principles from temporary implementation choices.

Technology answers how.

Philosophy answers why.

Forge is designed around the conviction that sustainable AI systems are built by making good architectural decisions long before selecting technologies.

Technology follows architecture.

Architecture follows philosophy.

⸻

4.2 Complexity Belongs Inside the Platform

Complexity cannot be eliminated.

It can only be moved.

Poor platforms push complexity into applications.

Good platforms absorb complexity so applications remain simple.

Forge exists to absorb engineering complexity.

Applications should never need to understand:

* prompt construction;
* provider routing;
* policy evaluation;
* context assembly;
* memory retrieval;
* audit generation;
* model failover;
* retry strategies.

Those concerns belong to the platform.

Business applications should remain focused on business logic.

⸻

Figure 4.1 — Complexity Migration

flowchart LR
subgraph Traditional
Application1["Application"]
Prompting
Policies
Memory
Security
Providers
end
subgraph Forge
BusinessLogic
ForgePlatform["Forge Platform"]
end
BusinessLogic --> ForgePlatform
ForgePlatform --> Context
ForgePlatform --> Policy
ForgePlatform --> Memory
ForgePlatform --> Providers

Engineering quality increases when applications become simpler.

⸻

4.3 Context Is Engineered, Not Prompted

Prompts are implementation artifacts.

Context is an engineering artifact.

Forge rejects the idea that AI systems should be designed by writing increasingly sophisticated prompts.

Instead, prompts are considered the final serialization of a carefully engineered context.

Context consists of:

* identity;
* intent;
* policies;
* organizational knowledge;
* execution state;
* historical interactions;
* environmental constraints.

Prompt engineering optimizes text.

Context engineering optimizes understanding.

This distinction is fundamental.

⸻

4.4 Knowledge Is Curated, Not Collected

Organizations rarely suffer from a lack of information.

They suffer from an excess of unmanaged information.

Adding more documents to a vector database does not create organizational knowledge.

Knowledge requires:

* ownership;
* versioning;
* validation;
* provenance;
* lifecycle management.

Forge therefore treats knowledge as a governed engineering asset.

Knowledge without governance eventually becomes noise.

⸻

4.5 Memory Is Intentional

Memory should never emerge accidentally.

Every memory exists for a specific engineering purpose.

Forge distinguishes between multiple memory domains.

For example:

Memory Type	Purpose
Execution Memory	Runtime state
Conversation Memory	Interaction continuity
Engineering Memory	Design knowledge
Repository Memory	Source-controlled knowledge
Organizational Memory	Institutional knowledge

Each memory type has different:

* retention;
* ownership;
* governance;
* security;
* lifecycle;
* retrieval strategies.

Memory is architecture.

Not storage.

⸻

4.6 Capabilities Outlive Technologies

Capabilities represent engineering responsibilities.

Technologies merely implement them.

For example:

Capability	Possible Implementations
Model Execution	OpenAI, Anthropic, Ollama
Workflow	LangGraph, Temporal, Custom Runtime
Identity	Entra ID, Keycloak
Policy	OPA, Cedar
Messaging	Kafka, NATS, RabbitMQ

Applications depend on capabilities.

Platforms choose implementations.

This separation dramatically increases longevity.

⸻

4.7 Every Decision Produces Evidence

Engineering without evidence eventually becomes opinion.

Every significant platform action SHOULD leave behind evidence.

Evidence includes:

* audit events;
* policy decisions;
* architectural decisions;
* deployment records;
* test results;
* evaluation outcomes;
* model selections;
* context composition.

Evidence is not logging.

Evidence explains why something happened.

⸻

Figure 4.2 — Evidence Pipeline

flowchart LR
Decision
Policy
Context
Runtime
Execution
Audit
EvidenceStore
Decision --> Audit
Policy --> Audit
Context --> Audit
Runtime --> Audit
Audit --> EvidenceStore

The objective is reconstructability.

An engineer should always be able to answer:

“Why did the platform behave this way?”

⸻

4.8 Engineering Over Automation

Automation is not the goal.

Good engineering is.

Automation without architecture creates fragile systems.

Automation without governance creates unsafe systems.

Automation without evidence creates untrustworthy systems.

Forge automates only where engineering discipline already exists.

Engineering always precedes automation.

⸻

4.9 Economic Intelligence

Engineering decisions consume resources.

Every token generated.

Every model invocation.

Every workflow execution.

Every retrieval.

Every tool call.

All have measurable cost.

Forge therefore considers economics a first-class architectural concern.

Optimization is not premature.

Waste is.

The platform continuously evaluates:

* cost;
* latency;
* energy consumption;
* provider efficiency;
* token utilization.

The objective is not minimizing cost.

The objective is maximizing engineering value.

⸻

4.10 Human Accountability

Artificial Intelligence expands human capability.

It does not replace human responsibility.

Humans remain accountable for:

* policy;
* governance;
* approvals;
* organizational knowledge;
* architectural decisions.

Forge deliberately preserves meaningful human oversight.

Autonomy without accountability is incompatible with enterprise engineering.

⸻

4.11 Simplicity Is an Outcome

Simplicity should never be confused with minimalism.

Simple systems often contain sophisticated engineering.

Forge embraces internal complexity whenever it creates external simplicity.

The measure of platform quality is not how little code Forge contains.

It is how much unnecessary complexity applications no longer need to implement.

⸻

4.12 The Forge Philosophy

The philosophy of Forge can be summarized in twelve statements.

1. Engineering precedes intelligence.
2. Complexity belongs inside the platform.
3. Context is engineered, not prompted.
4. Knowledge is curated, not collected.
5. Memory is intentional.
6. Capabilities outlive technologies.
7. Evidence is a first-class engineering artifact.
8. Economics influence architecture.
9. Humans remain accountable.
10. Platforms outlive applications.
11. Simplicity is created, not discovered.
12. Architecture exists to enable evolution.

Everything else in this Architecture Bible is a consequence of these beliefs.

⸻

Closing Thought

Every successful engineering discipline eventually develops a philosophy that transcends individual technologies.

The principles of civil engineering are not tied to a particular material.

The principles of software engineering are not tied to a particular programming language.

Likewise, the principles of AI Engineering should not depend on today’s models, frameworks or providers.

Forge is an attempt to articulate those principles before the industry standardizes them.

⸻

Ik merk dat de toon nu begint te “landen”. De eerste vier hoofdstukken vormen samen een logisch verhaal:

* Hoofdstuk 1: Waarom is AI engineering fundamenteel anders?
* Hoofdstuk 2: Waarom hebben we een Engineering Operating System nodig?
* Hoofdstuk 3: Welke toekomst willen we bouwen?
* Hoofdstuk 4: Welke overtuigingen sturen iedere architectuurbeslissing?

Vanaf hoofdstuk 5 kunnen we de filosofie vertalen naar een concreet architectuurmodel: de Capability Model. Daar gaan we Forge definiëren als een verzameling stabiele capabilities met expliciete contracten, verantwoordelijkheden en afhankelijkheden. Dat wordt het fundament waarop de rest van de technische architectuur rust.

Mooi. Nu begint het boek van denken naar ontwerpen te verschuiven.

Tot nu toe hebben we gesproken over waarom Forge bestaat. Vanaf nu gaan we definiëren wat Forge eigenlijk is.

En hier wil ik bewust afwijken van de meeste architectuurboeken.

De meeste platformen worden beschreven in termen van componenten.

Ik wil Forge beschrijven in termen van capabilities.

Dat is een fundamenteel verschil.

⸻

Chapter 5

Capability Model

Software components implement capabilities.

Capabilities define architecture.

The distinction is essential.

⸻

5.1 From Components to Capabilities

Traditional software architecture is often described by listing components.

For example:

* API Gateway
* Database
* Message Bus
* Cache
* Authentication Service

Although useful for implementation, this perspective has a major limitation.

Components evolve.

Responsibilities remain.

A cache may disappear.

A message broker may be replaced.

A database may migrate.

Yet the architectural responsibility they fulfilled continues to exist.

Forge therefore defines architecture through capabilities rather than implementations.

Capabilities survive technology.

⸻

5.2 What Is a Capability?

A capability is a stable engineering responsibility provided by the platform.

It answers one question:

What responsibility does the platform own?

It deliberately does not answer:

* Which framework?
* Which vendor?
* Which programming language?
* Which cloud?
* Which model?

Those are implementation decisions.

Capabilities define responsibilities.

Implementations realize them.

⸻

Definition

A Forge capability SHALL:

* own a clearly defined responsibility;
* expose explicit contracts;
* hide implementation details;
* evolve independently;
* be testable in isolation;
* be observable;
* produce evidence;
* be governed by policy.

These characteristics distinguish capabilities from services.

Services may implement capabilities.

Capabilities define architecture.

⸻

5.3 Why Capabilities Matter

Consider authentication.

Most systems think about authentication as an Identity Service.

Forge thinks differently.

Identity is not a service.

Identity is a capability.

Today it may be implemented using Microsoft Entra ID.

Tomorrow using Keycloak.

Later perhaps through decentralized identity.

Nothing changes architecturally.

Only the implementation changes.

The capability remains stable.

This distinction dramatically reduces architectural coupling.

⸻

5.4 The Forge Capability Landscape

Forge organizes the platform into a small number of fundamental capabilities.

mindmap
  root((Forge))
    Identity
    Policy
    Context
    Knowledge
    Memory
    Runtime
    Models
    Tools
    Audit
    Observability
    Governance
    Economics

These capabilities collectively define the Engineering Operating System.

Applications interact with capabilities rather than implementations.

⸻

5.5 Core Platform Capabilities

Forge distinguishes between Core, Supporting and Cross-Cutting capabilities.

flowchart TB
subgraph Core
Identity
Policy
Context
Knowledge
Memory
Runtime
end
subgraph Supporting
Models
Tools
Configuration
Secrets
end
subgraph CrossCutting
Audit
Observability
Governance
Economics
Security
end

This separation clarifies architectural responsibilities.

Not every capability participates equally in every request.

Cross-cutting capabilities influence all requests.

⸻

5.6 Capability Contracts

Every capability exposes a contract to the rest of the platform.

A contract consists of:

* purpose;
* responsibilities;
* inputs;
* outputs;
* guarantees;
* policies;
* evidence;
* metrics.

For example:

Capability	Input	Output
Identity	Credentials	Principal
Policy	Principal + Action	Decision
Context	Request	Context Package
Runtime	Context	Execution
Audit	Events	Evidence
Observability	Telemetry	Insights

Notice that no implementation technologies appear in these contracts.

Contracts remain stable while implementations evolve.

⸻

5.7 Capability Independence

Capabilities SHOULD remain independently deployable wherever practical.

They SHOULD also remain independently replaceable.

The following dependency rule applies:

Applications
        ↓
Capabilities
        ↓
Implementations
        ↓
Technologies

Dependencies SHALL always point downward.

Technology must never dictate architecture.

⸻

5.8 Capability Lifecycle

Capabilities evolve independently.

stateDiagram-v2
[*] --> Proposed
Proposed --> Experimental
Experimental --> Stable
Stable --> Deprecated
Deprecated --> Retired
Retired --> [*]

This lifecycle allows Forge to innovate without destabilizing existing applications.

Applications consume stable contracts rather than unstable implementations.

⸻

5.9 Capability Maturity

Not every capability matures at the same pace.

Forge therefore classifies maturity independently.

Level	Meaning
Experimental	Research
Preview	Early adopters
Stable	Production supported
Strategic	Core platform capability
Legacy	Scheduled for replacement

Maturity reflects engineering confidence rather than feature completeness.

⸻

5.10 Capability Dependency Graph

Capabilities depend on one another.

Those dependencies are explicit.

flowchart LR
Identity --> Policy
Policy --> Context
Knowledge --> Context
Memory --> Context
Context --> Runtime
Runtime --> Models
Runtime --> Tools
Runtime --> Audit
Runtime --> Observability
Audit --> Governance

This graph is intentionally acyclic.

Circular capability dependencies are prohibited.

⸻

5.11 Capability Evolution

Forge does not grow by adding features.

Forge grows by introducing new reusable capabilities.

For example:

Today:

* Context
* Memory
* Runtime

Tomorrow:

* Planning
* Simulation
* Reasoning
* Multi-Agent Coordination
* Compliance Intelligence

Applications automatically benefit from platform evolution without architectural redesign.

⸻

5.12 Capability Thinking

Capability thinking changes how engineers design systems.

Instead of asking:

“Which framework should we use?”

Engineers ask:

“Which capability is missing?”

Instead of asking:

“Which service owns this?”

They ask:

“Which responsibility belongs here?”

Instead of designing implementations first…

…they design responsibilities first.

Forge considers this inversion one of its defining architectural characteristics.

⸻

Closing Thought

Good architectures are resilient because they separate enduring responsibilities from temporary implementations.

Capabilities are the vocabulary through which Forge expresses that separation.

Everything that follows in this Architecture Bible—Context Engineering, Memory, Runtime, Policy, Governance and Security—is simply a detailed exploration of the capabilities introduced in this chapter.

⸻

Een gedachte die ik tijdens het schrijven kreeg

Ik begin steeds sterker het gevoel te krijgen dat we niet alleen een architectuurboek schrijven, maar een nieuwe ontwerptaal voor AI-platformen. Begrippen als Capability, Engineering Memory, Context Package, Evidence Pipeline en Engineering Operating System vormen samen een consistente vocabulaire. Als we dat goed uitwerken, kunnen mensen straks over Forge praten zonder meteen over specifieke technologieën te hoeven beginnen. Dat zou voor mij een teken zijn dat de Architecture Bible echt geslaagd is.

Perfect. Vanaf hier verlaten we de “algemene architectuur”. Vanaf nu gaan we de disciplines definiëren die Forge uniek maken.

En ik denk dat hier de eerste écht originele bijdrage komt.

Niet Runtime.

Niet Memory.

Context.

Ik denk namelijk dat vrijwel iedereen in de industrie Context ziet als iets dat vóór een prompt gebeurt.

Ik denk dat Context de kern van AI Engineering is.

Sterker nog…

Ik zou een behoorlijk controversiële stelling durven innemen:

LLMs do not fail because they lack intelligence.

They fail because they receive poorly engineered context.

Dat wordt de rode draad van dit hoofdstuk.

⸻

Chapter 6

Context Engineering

Artificial Intelligence does not operate on prompts.

It operates on context.

The prompt is merely the final serialization of that context.

⸻

6.1 The Misunderstanding

The AI industry spends enormous effort discussing prompts.

Prompt engineering.

Prompt templates.

Prompt chaining.

Prompt optimization.

Prompt versioning.

These techniques are valuable.

Yet they all focus on the final artifact sent to the model.

Forge deliberately shifts the engineering focus earlier.

Long before the first token is generated, an engineering system has already made dozens of decisions.

Who is the user?

What are they allowed to access?

What information is relevant?

Which policies apply?

Which knowledge is authoritative?

Which memory should be recalled?

Which tools are available?

Which model is appropriate?

None of these decisions belong inside a prompt.

They belong inside a context.

⸻

6.2 Context Is More Than Information

A common misconception is that context equals additional information.

This definition is incomplete.

Context is not information.

Context is the set of constraints that allows information to become meaningful.

Consider the sentence:

Approve the request.

Without context, this instruction is meaningless.

Who issued it?

Approve what?

Under which policy?

At what risk level?

Against which organizational rules?

The same language model can produce entirely different responses depending on these contextual dimensions.

The intelligence did not change.

Only the context did.

⸻

6.3 A Definition of Context

Forge defines context as:

The complete, governed, time-bound representation of everything an intelligent system requires to make an engineering decision.

This definition intentionally extends far beyond prompts or retrieved documents.

Context is composed of multiple independent dimensions.

Each dimension contributes to the overall understanding of a request.

⸻

6.4 The Context Pyramid

Not every piece of context is equally important.

Forge organizes context into layers.

flowchart BT
Prompt["Prompt Serialization"]
Instructions
Task
Knowledge
Memory
Policy
Identity
Environment
Objectives
Environment --> Identity
Identity --> Policy
Policy --> Memory
Memory --> Knowledge
Knowledge --> Task
Task --> Instructions
Instructions --> Prompt

The prompt is the smallest visible layer.

Everything beneath it determines its quality.

⸻

6.5 Dimensions of Context

Forge identifies multiple orthogonal dimensions.

Dimension	Purpose
Identity	Who is acting?
Intent	Why is the request made?
Task	What must be accomplished?
Policy	What is allowed?
Knowledge	What facts are available?
Memory	What should be remembered?
Environment	Which external conditions apply?
Organization	Which organizational rules exist?
Time	What is valid right now?
Risk	Which safeguards are required?

These dimensions evolve independently.

No single retrieval mechanism can construct them all.

⸻

6.6 Context Is Constructed

Context should never be retrieved as a single object.

It is assembled.

flowchart LR
Identity
Policy
Knowledge
Memory
Task
Environment
Objectives
Identity --> Builder
Policy --> Builder
Knowledge --> Builder
Memory --> Builder
Task --> Builder
Environment --> Builder
Objectives --> Builder
Builder --> ContextPackage
ContextPackage --> PromptSerializer
PromptSerializer --> LanguageModel

The Context Builder becomes one of the central responsibilities of the Forge platform.

⸻

6.7 Context Has a Lifecycle

Context is dynamic.

It is continuously created, refined and discarded.

stateDiagram-v2
[*] --> Requested
Requested --> Collected
Collected --> Validated
Validated --> Enriched
Enriched --> Serialized
Serialized --> Executed
Executed --> Archived
Archived --> Expired
Expired --> [*]

This lifecycle ensures that context remains current, explainable and auditable.

⸻

6.8 Context Is Governed

Not every fact should become context.

Adding more information often reduces quality.

Forge therefore applies explicit governance before context reaches a model.

Context should satisfy several properties.

It should be:

* relevant;
* authoritative;
* recent;
* permitted;
* explainable;
* minimal.

Anything failing these criteria should be excluded.

The best context is not the largest context.

It is the most useful one.

⸻

6.9 The Principle of Minimal Sufficient Context

Forge introduces one of its core engineering principles:

Minimal Sufficient Context

This principle states:

Provide exactly the context required to perform the requested task—and no more.

Too little context causes hallucinations.

Too much context causes confusion.

Both reduce quality.

Context engineering is therefore an optimization problem rather than a retrieval problem.

⸻

6.10 Context as an Engineering Artifact

Traditional software engineering produces artifacts.

Requirements.

Architecture.

Source code.

Tests.

Forge adds another first-class artifact:

The Context Package.

A Context Package is:

* versionable;
* inspectable;
* reproducible;
* testable;
* explainable;
* governed.

Applications should never assemble context independently.

They should request context from the platform.

⸻

6.11 Context Engineering

Forge defines Context Engineering as:

The engineering discipline responsible for constructing, governing, validating and evolving Context Packages for intelligent systems.

Prompt engineering becomes one implementation technique within Context Engineering.

Not the discipline itself.

This distinction fundamentally changes how AI systems are designed.

⸻

Closing Thought

Large Language Models generate text.

Context determines whether that text becomes intelligence.

Organizations often invest heavily in larger models, more retrieval systems and increasingly sophisticated prompting techniques.

Forge proposes a different path.

Engineer context first.

The quality of intelligence will follow.

⸻

Dit voelt als het eerste écht nieuwe idee

Ik denk dat we hier iets pakken wat ik nog niet expliciet zo heb gezien in de literatuur.

Vrijwel iedereen praat over:

* Prompt Engineering
* RAG
* Memory
* Agents

Wij introduceren een discipline boven al die concepten:

Context Engineering

En vervolgens zeggen we:

* Prompting is een serializer.
* RAG is één contextbron.
* Memory is één contextbron.
* Policies zijn context.
* Identity is context.
* Tijd is context.
* Risico is context.

Dat maakt Context Engineering de centrale discipline van Forge. Ik vermoed eerlijk gezegd dat dit uiteindelijk het hoofdstuk wordt waar mensen het meest naar zullen verwijzen, omdat het een nieuw mentaal model introduceert in plaats van een bestaande techniek iets anders te noemen.

Mooi. Dit hoofdstuk bouwt direct voort op Context Engineering.

Als Context bepaalt wat een AI op dit moment weet, dan bepaalt Memory wat een AI door de tijd heen blijft weten.

En hier wil ik opnieuw een stap verder gaan dan de huidige industrie.

Vrijwel iedereen behandelt memory als “chat history”.

Ik denk dat dat fundamenteel verkeerd is.

Memory is geen feature.

Memory is een architectural capability.

⸻

Chapter 7

Memory Engineering

Context explains the present.

Memory connects the past to the present.

Without memory, intelligence has no continuity.

⸻

7.1 The Memory Misconception

Most AI systems define memory as conversation history.

Previous messages are stored.

Older messages are summarized.

Relevant fragments are retrieved.

This approach is useful.

It is also incomplete.

Conversation history is only one manifestation of memory.

Organizations possess many different kinds of memory.

Projects have memory.

Repositories have memory.

Architectures have memory.

People have memory.

Policies have memory.

Operations have memory.

Treating all of these as “chat history” is an architectural oversimplification.

Forge therefore introduces a broader discipline.

Memory Engineering.

⸻

7.2 Memory Is Not Storage

Storage answers one question:

Can information be retained?

Memory answers a different question:

Should this information influence future decisions?

A log file is storage.

An audit trail is storage.

A database is storage.

None of them automatically become memory.

Memory requires intention.

Something becomes memory only when it is expected to shape future behavior.

⸻

Definition

Forge defines memory as:

A governed representation of past information that intentionally influences future decisions.

Three words are essential.

Governed.

Memory is never accidental.

Past.

Memory represents historical state.

Influences.

If historical information never affects future behavior, it is not memory.

It is merely archived data.

⸻

7.3 The Memory Hierarchy

Forge distinguishes multiple layers of memory.

flowchart TB
Knowledge["Knowledge Memory"]
Organization["Organizational Memory"]
Engineering["Engineering Memory"]
Repository["Repository Memory"]
Conversation["Conversation Memory"]
Execution["Execution Memory"]
Execution --> Conversation
Conversation --> Repository
Repository --> Engineering
Engineering --> Knowledge
Knowledge --> Organization

Each layer has different ownership, retention and governance.

Higher layers evolve more slowly.

Lower layers change continuously.

⸻

7.4 Memory Domains

Memory is not homogeneous.

Each domain exists for a different engineering purpose.

Memory Domain	Primary Purpose	Typical Lifetime
Execution Memory	Runtime coordination	Seconds to minutes
Session Memory	User continuity	Minutes to hours
Conversation Memory	Dialogue coherence	Hours to days
Project Memory	Project evolution	Months to years
Repository Memory	Source of Truth	Permanent
Engineering Memory	Architectural decisions	Permanent
Organizational Memory	Institutional knowledge	Decades

Each domain requires different governance.

Treating them identically leads to poor engineering outcomes.

⸻

7.5 Memory Is Layered

Memory should not be retrieved from a single source.

Instead, Forge assembles memory progressively.

flowchart LR
Execution
Conversation
Project
Repository
Knowledge
Organization
Execution --> MemoryComposer
Conversation --> MemoryComposer
Project --> MemoryComposer
Repository --> MemoryComposer
Knowledge --> MemoryComposer
Organization --> MemoryComposer
MemoryComposer --> ContextBuilder

Memory contributes to context.

It does not replace it.

⸻

7.6 Memory Ownership

Every memory has an owner.

Ownership determines:

* who may modify it;
* who may approve changes;
* who may access it;
* who defines retention;
* who validates correctness.

For example:

Memory	Owner
Conversation	User
Project	Project Team
Repository	Maintainers
Architecture	Architects
Policy	Governance Team

Ownership is essential for trust.

Memory without ownership inevitably degrades.

⸻

7.7 Memory Lifecycle

Every memory follows a lifecycle.

stateDiagram-v2
[*] --> Created
Created --> Validated
Validated --> Active
Active --> Updated
Updated --> Active
Active --> Archived
Archived --> Retired
Retired --> [*]

This lifecycle ensures that memory remains accurate over time.

Not everything should remain active forever.

⸻

7.8 Memory Quality

Forge evaluates memory continuously.

Useful memory is:

* accurate;
* authoritative;
* current;
* relevant;
* attributable;
* explainable.

Poor memory introduces engineering debt.

Good memory increases organizational intelligence.

Quality matters more than quantity.

⸻

7.9 Memory as a Platform Capability

Applications should never implement memory independently.

Instead they request memory through platform capabilities.

flowchart LR
Application
ForgeMemory
ExecutionStore
ConversationStore
Repository
KnowledgeBase
Application --> ForgeMemory
ForgeMemory --> ExecutionStore
ForgeMemory --> ConversationStore
ForgeMemory --> Repository
ForgeMemory --> KnowledgeBase

The application should not need to know where memory resides.

It should only express its intent.

⸻

7.10 Repository Memory

Forge introduces a concept that I believe will become increasingly important:

Repository Memory.

A Git repository is more than source code.

It contains:

* architecture;
* ADRs;
* documentation;
* issues;
* discussions;
* pull requests;
* release notes;
* tests;
* implementation history.

Together these artifacts represent institutional engineering knowledge.

Repositories therefore become active memory systems rather than passive storage locations.

This is why Forge considers Git the canonical Source of Truth.

⸻

7.11 Engineering Memory

Engineering Memory represents the accumulated reasoning behind a system.

It includes:

* architectural decisions;
* trade-offs;
* rejected alternatives;
* operational lessons;
* incident analyses;
* design constraints;
* engineering principles.

Unlike implementation details, Engineering Memory should survive major rewrites.

It explains why a system exists in its current form.

Without it, organizations repeatedly solve the same problems.

⸻

7.12 Organizational Memory

Organizations routinely lose knowledge through:

* employee turnover;
* project completion;
* organizational restructuring;
* undocumented decisions.

Forge treats this as an engineering problem rather than a human inevitability.

Organizational Memory should continuously absorb validated knowledge from projects, repositories and operations.

The objective is simple:

Knowledge should leave projects before people leave projects.

⸻

7.13 Memory Engineering

Forge defines Memory Engineering as:

The engineering discipline responsible for creating, governing, validating, evolving and retiring memory across intelligent systems.

Memory Engineering extends far beyond conversation history.

It defines how organizations preserve engineering intelligence over years rather than conversations over minutes.

⸻

Closing Thought

Human intelligence depends on memory.

Organizational intelligence depends on shared memory.

Artificial intelligence should be no different.

Systems that repeatedly forget become increasingly expensive to operate.

Systems that intentionally remember become increasingly valuable.

Forge therefore treats memory not as an implementation feature, but as one of the foundational capabilities of AI Engineering.

⸻

Ik denk dat hier iets bijzonders ontstaat

Ik zie nu een patroon ontstaan dat veel verder gaat dan een AI-platform.

* Context Engineering → het construeren van het juiste “nu”.
* Memory Engineering → het behouden van het juiste “verleden”.

Het volgende hoofdstuk zou logisch zijn: Knowledge Engineering. Daarin beantwoorden we de vraag: “Hoe wordt informatie omgezet in gevalideerde, beheerde kennis?”

Dan ontstaat een krachtige driedeling:

* Knowledge = wat de organisatie weet.
* Memory = wat de organisatie onthoudt.
* Context = wat de AI op dit moment nodig heeft.

Die drie samen vormen naar mijn idee het intellectuele hart van Forge. Ze zijn duidelijk van elkaar te onderscheiden, maar versterken elkaar in de architectuur. Ik denk dat juist deze samenhang Forge een eigen identiteit geeft.

Mooi. En ik denk dat we hier een cruciale ontwerpbeslissing moeten nemen.

Ik zou niet direct naar Runtime, Policy of Security gaan.

Er ontbreekt nog één fundamenteel begrip.

Sterker nog, ik denk dat dit begrip de brug vormt tussen Context en Memory.

Dat is Knowledge.

De industrie gebruikt “knowledge” en “memory” vaak door elkaar. Ik denk dat dat een architectonische fout is.

Ik zou durven stellen:

Knowledge is what an organization knows.

Memory is what an organization intentionally retains.

Context is what an AI needs right now.

Dat zijn drie verschillende concepten.

En ik denk dat Forge die expliciet moet scheiden.

⸻

Chapter 8

Knowledge Engineering

Information is abundant.

Knowledge is scarce.

Engineering begins when organizations can distinguish the two.

⸻

8.1 The Information Problem

Modern organizations rarely suffer from a lack of information.

Instead, they struggle with an overwhelming abundance of it.

Documents.

Emails.

Wiki pages.

Chat messages.

Source code.

Pull requests.

Meeting notes.

Presentations.

Specifications.

Logs.

Policies.

Most AI systems attempt to solve this problem through retrieval.

Retrieve more documents.

Embed more content.

Increase the context window.

Improve ranking.

These techniques improve access to information.

They do not automatically create knowledge.

Knowledge is something fundamentally different.

⸻

8.2 Information Is Not Knowledge

Information exists independently.

Knowledge exists because someone has established meaning.

Consider these examples.

A PDF stored in SharePoint is information.

An Architecture Decision Record is knowledge.

A meeting transcript is information.

The approved design resulting from that meeting is knowledge.

A Git commit is information.

The architectural reasoning behind that commit is knowledge.

Knowledge requires interpretation.

It requires validation.

It requires ownership.

⸻

Definition

Forge defines knowledge as:

Validated information that an organization intentionally recognizes as authoritative for future engineering decisions.

Every word matters.

Validated.

Recognized.

Authoritative.

Engineering decisions.

Knowledge is therefore not discovered.

Knowledge is established.

⸻

8.3 The Knowledge Pyramid

Forge distinguishes four progressively more valuable states.

flowchart BT
Wisdom
Knowledge
Information
Data
Data --> Information
Information --> Knowledge
Knowledge --> Wisdom

Data answers:

“What exists?”

Information answers:

“What happened?”

Knowledge answers:

“What should we believe?”

Wisdom answers:

“What should we do?”

Forge focuses primarily on engineering Knowledge.

⸻

8.4 Characteristics of Engineering Knowledge

Engineering knowledge exhibits several unique characteristics.

It is:

* validated;
* versioned;
* attributable;
* reviewable;
* explainable;
* reproducible;
* governed.

Without these characteristics, information cannot be considered authoritative.

⸻

8.5 Knowledge Domains

Knowledge exists in different domains.

Domain	Examples
Architecture	Reference architectures, diagrams
Business	Processes, capabilities
Operations	Runbooks, procedures
Security	Policies, standards
Development	Coding standards, ADRs
Product	Requirements, roadmaps
Compliance	Regulations, evidence

Each domain has different owners and different validation processes.

Knowledge is therefore federated rather than centralized.

⸻

8.6 Knowledge Sources

Forge recognizes many possible knowledge sources.

flowchart LR
Repository
Documentation
ADRs
Policies
Runbooks
Standards
APIs
Repositories --> KnowledgeGraph
Documentation --> KnowledgeGraph
ADRs --> KnowledgeGraph
Policies --> KnowledgeGraph
Runbooks --> KnowledgeGraph
Standards --> KnowledgeGraph
APIs --> KnowledgeGraph

Importantly, sources remain authoritative.

Forge does not duplicate authority.

It references it.

⸻

8.7 Canonical Sources

One of Forge’s strongest principles is:

Every engineering fact should have a single canonical source.

For example:

Artifact	Canonical Source
Source Code	Git Repository
Architecture	Architecture Repository
Policies	Policy Repository
Identity	Identity Provider
Configuration	Configuration Repository
Secrets	Secret Manager

This principle prevents competing truths.

Synchronization should never replace ownership.

⸻

8.8 Knowledge Lifecycle

Knowledge continuously evolves.

stateDiagram-v2
[*] --> Proposed
Proposed --> Reviewed
Reviewed --> Approved
Approved --> Published
Published --> Updated
Updated --> Reviewed
Published --> Deprecated
Deprecated --> Archived
Archived --> [*]

Only approved knowledge should become authoritative.

⸻

8.9 Knowledge Quality

Forge evaluates knowledge using explicit quality attributes.

Good knowledge is:

* correct;
* complete;
* current;
* consistent;
* attributable;
* discoverable;
* understandable.

Knowledge quality is an engineering responsibility.

Not an editorial preference.

⸻

8.10 Knowledge as Code

One of the defining principles of Forge is:

Knowledge should evolve through engineering workflows rather than document management workflows.

Knowledge therefore benefits from:

* Git versioning;
* pull requests;
* code reviews;
* CI validation;
* automated quality checks;
* traceability;
* rollback;
* ownership.

Markdown is not merely documentation.

Within Forge, Markdown becomes executable organizational knowledge.

⸻

8.11 Repository as Organizational Intelligence

Forge extends the traditional role of Git.

Repositories no longer contain only source code.

They contain:

* engineering intent;
* architecture;
* ADRs;
* standards;
* governance;
* documentation;
* examples;
* evidence.

Together these artifacts form the organization’s collective engineering intelligence.

Repositories therefore become knowledge systems.

Not storage systems.

⸻

8.12 Knowledge Engineering

Forge defines Knowledge Engineering as:

The engineering discipline responsible for creating, validating, governing, evolving and publishing authoritative organizational knowledge.

Knowledge Engineering is not concerned with retrieval alone.

It is concerned with trust.

Retrieval answers:

“Can we find information?”

Knowledge Engineering answers:

“Can we trust it?”

⸻

8.13 Relationship to Memory and Context

The relationship between the three disciplines can now be described precisely.

flowchart LR
Information --> Knowledge
Knowledge --> Memory
Memory --> Context
Context --> Intelligence
Intelligence --> Decisions

Information becomes Knowledge through validation.

Knowledge becomes Memory through intentional retention.

Memory becomes Context through selective composition.

Context enables Intelligence.

Intelligence supports Decisions.

Each transformation adds engineering value.

⸻

8.14 The Knowledge Contract

Every knowledge artifact SHALL answer six questions.

1. Who owns it?
2. Who approved it?
3. When was it last validated?
4. Which decisions depend upon it?
5. Which policies govern it?
6. Which evidence supports it?

If these questions cannot be answered, the artifact should not be considered authoritative.

⸻

Closing Thought

Organizations do not become intelligent because they accumulate information.

They become intelligent because they establish trusted knowledge, preserve it through intentional memory and compose it into meaningful context.

Forge considers these three disciplines inseparable.

Knowledge explains what is true.

Memory explains what endures.

Context explains what matters now.

Together they form the cognitive architecture of the Engineering Operating System.

⸻

Ik zie nu de echte kern van Forge ontstaan

Als ik terugkijk op de afgelopen hoofdstukken, denk ik dat de “ziel” van Forge inmiddels zichtbaar wordt.

Ik zou de rest van het boek rond drie fundamentele pijlers organiseren:

Knowledge Engineering
        ↓
Memory Engineering
        ↓
Context Engineering
        ↓
Runtime Engineering
        ↓
AI Providers

Dat is precies de omgekeerde volgorde van hoe de meeste AI-platformen zijn opgebouwd. Zij beginnen bij het model en werken naar buiten.

Forge begint bij de organisatie en werkt naar binnen.

Dat voelt voor mij als de onderscheidende architectuurvisie van Forge. Het betekent dat een taalmodel helemaal onderaan de stack terechtkomt: als een vervangbare uitvoeringslaag. De echte intellectuele eigendom zit in de lagen daarboven. Ik denk dat dát de boodschap is die de Architecture Bible uniek kan maken.

Mooi. Nu komt voor mij het hoofdstuk waar de voorgaande drie disciplines samenkomen.

Als we de rode draad volgen:

* Knowledge → wat de organisatie weet.
* Memory → wat de organisatie bewust bewaart.
* Context → wat op dit moment relevant is.

Dan ontbreekt nog één stap.

Hoe wordt daar intelligent gedrag van gemaakt?

Dat is niet de taak van een LLM.

Dat is de taak van de Runtime.

En ik denk dat dit misschien wel de grootste denkfout is die de huidige AI-industrie maakt.

Vrijwel iedereen ziet de runtime als “de agent”.

Ik zie de runtime als de operating system kernel.

Dat is een fundamenteel ander uitgangspunt.

⸻

Chapter 9

Runtime Engineering

Language models generate tokens.

Runtimes generate behaviour.

Engineering determines the difference.

⸻

9.1 Beyond Agent Frameworks

The rapid adoption of Large Language Models has led to an explosion of so-called agent frameworks.

Most focus on orchestrating prompts, tools and model invocations.

While useful, this framing is incomplete.

An agent is not an architecture.

It is one possible execution pattern.

Forge deliberately avoids building around the concept of agents.

Instead, Forge builds around execution.

Execution is a broader concept.

An execution may involve:

* a single model invocation;
* a deterministic workflow;
* multiple collaborating models;
* human approvals;
* enterprise systems;
* policy enforcement;
* external tools.

The runtime is responsible for coordinating all of these activities.

Not the language model.

⸻

9.2 Defining the Runtime

Forge defines the Runtime as:

The execution capability responsible for transforming governed context into governed outcomes.

This definition intentionally avoids references to:

* language models;
* prompts;
* agents;
* frameworks.

Those are implementation choices.

The runtime owns execution.

Execution remains stable even when execution technologies evolve.

⸻

9.3 Responsibilities

The Runtime owns only execution concerns.

It deliberately does not own:

* identity;
* policy;
* knowledge;
* memory;
* governance.

Those capabilities provide services to the Runtime.

The Runtime consumes them.

Its responsibilities include:

* execution planning;
* workflow coordination;
* model orchestration;
* tool orchestration;
* retry handling;
* timeout management;
* state transitions;
* evidence production;
* execution telemetry.

Nothing more.

Nothing less.

⸻

Figure 9.1 — Runtime in the Platform

flowchart LR
Identity --> Context
Knowledge --> Context
Memory --> Context
Policy --> Context
Context --> Runtime
Runtime --> Models
Runtime --> Tools
Runtime --> Audit
Runtime --> Observability

Notice that the Runtime does not construct context.

It consumes context.

This separation is intentional.

⸻

9.4 The Runtime Contract

Every execution entering the Runtime SHALL satisfy a contract.

Minimum inputs include:

* Context Package
* Execution Objective
* Execution Policy
* Tool Permissions
* Model Constraints
* Correlation Identifier

Minimum outputs include:

* Execution Result
* Evidence Package
* Telemetry
* Audit Events
* Updated Execution State

The Runtime therefore behaves as a deterministic execution engine surrounding a probabilistic model.

⸻

9.5 Planning Before Execution

Execution begins long before the first token is generated.

Forge separates planning from execution.

flowchart LR
Objective
Context
Policies
Capabilities
Objective --> Planner
Context --> Planner
Policies --> Planner
Capabilities --> Planner
Planner --> ExecutionPlan
ExecutionPlan --> Runtime

Execution Plans are inspectable.

They are testable.

They are reproducible.

Most importantly:

They exist before execution starts.

⸻

9.6 Execution Is a State Machine

Every execution follows a well-defined lifecycle.

stateDiagram-v2
[*] --> Created
Created --> Planned
Planned --> Authorized
Authorized --> Executing
Executing --> Waiting
Waiting --> Executing
Executing --> Completed
Executing --> Failed
Completed --> Archived
Failed --> Archived
Archived --> [*]

The Runtime never jumps between states.

Every transition is observable.

Every transition is auditable.

⸻

9.7 Deterministic Shell, Probabilistic Core

This is perhaps the most important architectural principle in Forge.

The language model is probabilistic.

The platform surrounding it must not be.

+----------------------------------------------+
|          Deterministic Platform              |
|                                              |
| Identity                                     |
| Policy                                       |
| Context                                      |
| Memory                                       |
| Runtime                                      |
| Audit                                        |
| Governance                                   |
|                                              |
|     +-------------------------------+        |
|     | Probabilistic Language Model  |        |
|     +-------------------------------+        |
+----------------------------------------------+

Only one component is intentionally nondeterministic.

Everything else should behave predictably.

This dramatically simplifies testing, debugging and governance.

⸻

9.8 Tool Execution

The Runtime never calls enterprise systems directly.

Instead, all external capabilities are exposed through governed tools.

flowchart LR
Runtime
ToolRegistry
MCP
EnterpriseSystems
Runtime --> ToolRegistry
ToolRegistry --> MCP
MCP --> EnterpriseSystems

This abstraction provides:

* policy enforcement;
* auditing;
* provider independence;
* versioning;
* discoverability.

The Runtime executes capabilities.

The Tool Layer executes integrations.

⸻

9.9 Human Execution

Not every execution should be automated.

Humans are legitimate execution participants.

flowchart LR
Runtime --> HumanApproval
Runtime --> Tool
Runtime --> Model
HumanApproval --> Runtime

Human interaction is therefore a native execution primitive.

Not an exception.

⸻

9.10 Runtime Composition

A Runtime execution is assembled from reusable capabilities.

mindmap
  root((Runtime))
    Planner
    Scheduler
    Executor
    Tool Manager
    Model Router
    Retry Manager
    State Manager
    Evidence Collector
    Telemetry

Each responsibility evolves independently.

⸻

9.11 Runtime Principles

Forge Runtime follows ten principles.

1. Execution is deterministic wherever possible.
2. Planning precedes execution.
3. Context precedes planning.
4. Policy precedes execution.
5. Humans remain first-class participants.
6. Every transition produces evidence.
7. Every execution is observable.
8. Capabilities are preferred over implementations.
9. Failures are expected.
10. Execution is continuously improvable.

⸻

9.12 Runtime Engineering

Forge defines Runtime Engineering as:

The engineering discipline responsible for planning, coordinating, executing, observing and evolving intelligent execution workflows.

Runtime Engineering is not concerned with model intelligence.

It is concerned with execution quality.

Execution quality determines whether intelligence becomes operational capability.

⸻

Closing Thought

Language models answer questions.

Runtimes accomplish objectives.

The distinction matters.

A model may know how to perform a task.

Only a runtime can ensure the task is executed securely, consistently, observably and in accordance with organizational policy.

Forge therefore treats the Runtime not as an agent framework, but as the execution kernel of the Engineering Operating System.

⸻

Ik denk dat we zojuist een belangrijk architectuurprincipe hebben ontdekt

Ik zie nu een structuur ontstaan die sterk doet denken aan een traditioneel operating system:

Operating System	Forge
Filesystem	Knowledge
Virtual Memory	Memory
Process Context	Context
Kernel	Runtime
Access Control	Policy
User Accounts	Identity
Device Drivers	MCP Tools
Scheduler	Planner / Runtime
System Logs	Audit
Performance Monitor	Observability

Dat is meer dan een analogie. Het geeft Forge een consistente architectuur waarin elk onderdeel een duidelijke rol heeft. Net zoals een klassiek besturingssysteem applicaties beschermt tegen de complexiteit van hardware, beschermt Forge applicaties tegen de complexiteit van AI. Ik denk dat dit uiteindelijk de centrale metafoor van de hele Architecture Bible zou moeten worden.

Mooi. Nu komen we bij een hoofdstuk waarvan ik denk dat het één van de meest praktische én onderscheidende onderdelen van Forge wordt.

Tot nu toe hebben we beschreven:

* Knowledge → wat waar is.
* Memory → wat blijft bestaan.
* Context → wat nu relevant is.
* Runtime → hoe iets wordt uitgevoerd.

Maar er ontbreekt nog iets.

Wie bepaalt wat mág?

Dat is niet de Runtime.

Dat is niet Identity.

En het is ook niet Security.

Dat is Policy.

En hier wil ik opnieuw een stap verder gaan dan de huidige praktijk.

De meeste systemen behandelen policy als een autorisatiemechanisme.

Forge behandelt policy als de besluitvormingslaag van het Engineering Operating System.

⸻

Chapter 10

Policy Engineering

Intelligence determines what can be done.

Policy determines what may be done.

Engineering requires both.

⸻

10.1 Beyond Authorization

Most software engineers associate policy with authorization.

Can Alice access this document?

Can Bob invoke this API?

Can this service read this database?

These are important questions.

But they represent only one category of policy.

Artificial Intelligence introduces many additional decisions.

Examples include:

* Which model may process classified information?
* Which knowledge sources are authoritative?
* Which tools are available for this request?
* Which context elements may leave the organization?
* When is human approval required?
* Which execution strategy minimizes risk?
* Which audit evidence must be retained?

These are engineering decisions.

Not security decisions.

Forge therefore expands the concept of policy beyond authorization.

⸻

10.2 Defining Policy

Forge defines policy as:

A governed decision rule that constrains or directs the behaviour of the Engineering Operating System.

Several aspects are important.

A policy does not execute work.

A policy does not perform validation.

A policy does not orchestrate workflows.

A policy makes decisions.

Execution belongs to the Runtime.

Enforcement belongs to platform capabilities.

Policy provides the decision.

⸻

10.3 Policy as a First-Class Capability

Policy is not embedded in application code.

Policy is not hidden inside prompts.

Policy is not scattered across configuration files.

Policy exists as an independent platform capability.

flowchart LR
Application --> Runtime
Runtime --> Policy
Policy --> Decision
Decision --> Runtime
Runtime --> Execution

Separating policy from execution enables governance without modifying business logic.

⸻

10.4 Categories of Policy

Forge recognizes multiple policy domains.

Policy Domain	Example
Identity	Authentication requirements
Authorization	Access decisions
Data	Information classification
Context	Allowed context sources
Knowledge	Approved repositories
Memory	Retention periods
Runtime	Execution limits
Models	Model selection rules
Tools	Tool permissions
Economics	Budget limits
Compliance	Regulatory requirements
Audit	Evidence retention

Policy extends across the entire platform.

⸻

10.5 Policy Evaluation

Policy evaluation occurs continuously.

It is not a single authorization step.

flowchart TD
Request
Request --> IdentityPolicy
IdentityPolicy --> ContextPolicy
ContextPolicy --> KnowledgePolicy
KnowledgePolicy --> RuntimePolicy
RuntimePolicy --> ToolPolicy
ToolPolicy --> AuditPolicy
AuditPolicy --> Complete

Every capability may request policy decisions.

Every decision is independently auditable.

⸻

10.6 Policy Is Declarative

Forge strongly prefers declarative policies.

Instead of embedding rules inside code:

if classification == "SECRET":
    deny()

Forge expresses intent declaratively.

Rule:
Secret information SHALL NOT be processed by non-approved models.

Implementations may use:

* Open Policy Agent (OPA)
* Cedar
* Rego
* Custom engines

The language is replaceable.

The rule is not.

⸻

10.7 Policy Hierarchy

Policies exist at multiple organizational levels.

flowchart TB
Global
Organization
BusinessUnit
Project
Application
Execution
Global --> Organization
Organization --> BusinessUnit
BusinessUnit --> Project
Project --> Application
Application --> Execution

Lower levels may refine higher-level policies.

They may never weaken mandatory constraints.

⸻

10.8 Policy Composition

Complex decisions rarely depend on a single rule.

Forge composes policies.

flowchart LR
Identity
Security
Privacy
Compliance
Economics
Operations
Identity --> DecisionEngine
Security --> DecisionEngine
Privacy --> DecisionEngine
Compliance --> DecisionEngine
Economics --> DecisionEngine
Operations --> DecisionEngine
DecisionEngine --> Decision

The outcome is a single governed decision.

Not competing independent results.

⸻

10.9 Policy Is Explainable

Every policy decision SHALL answer:

* Which rule applied?
* Why did it apply?
* Which inputs were evaluated?
* Which evidence supports the decision?
* Which authority owns the rule?

Policy without explanation is governance without transparency.

⸻

10.10 Policy Produces Evidence

Every evaluation generates evidence.

flowchart LR
Policy
Decision
Evidence
Audit
Policy --> Decision
Decision --> Evidence
Evidence --> Audit

This enables engineers to reconstruct every significant platform decision.

⸻

10.11 Policy Evolution

Policies evolve continuously.

stateDiagram-v2
[*] --> Draft
Draft --> Review
Review --> Approved
Approved --> Active
Active --> Superseded
Superseded --> Archived
Archived --> [*]

Policy changes should follow engineering workflows rather than operational shortcuts.

⸻

10.12 Policy Engineering

Forge defines Policy Engineering as:

The engineering discipline responsible for defining, governing, evaluating and evolving decision rules across the Engineering Operating System.

Policy Engineering separates decision-making from execution.

This distinction enables organizations to evolve governance independently of implementation.

⸻

10.13 The Principle of Policy Independence

One of the most important architectural rules in Forge is:

Execution should never decide policy.

Likewise:

* Models do not decide policy.
* Tools do not decide policy.
* Applications do not decide policy.

They request decisions from the Policy capability.

This ensures consistency across the platform.

⸻

Closing Thought

As AI systems become more autonomous, the importance of policy increases rather than decreases.

Intelligence without policy becomes unpredictable.

Automation without policy becomes unsafe.

Governance without policy becomes unenforceable.

Forge therefore treats policy not as a security feature, but as the decision-making system of the Engineering Operating System.

⸻

Ik denk dat we nu de “kernel” compleet hebben

Wat me opvalt is dat de eerste tien hoofdstukken inmiddels een gesloten architectuur vormen.

Knowledge
        ↓
Memory
        ↓
Context
        ↓
Policy
        ↓
Runtime
        ↓
Models / Tools

Dat voelt fundamenteel juist.

Maar ik zie ook iets dat ik aan het begin van dit traject nog niet zag.

Identity hoort eigenlijk niet vóór Policy.

Identity is namelijk slechts één van de invoeren voor een beleidsbeslissing.

Een beleidsbeslissing kan ook afhangen van:

* classificatie;
* tijd;
* locatie;
* apparaat;
* risico;
* budget;
* wetgeving;
* organisatorische rol;
* projectstatus.

Identity is dus een contextdimensie, geen overkoepelend principe.

Dat betekent dat onze architectuur nóg consistenter kan worden:

Knowledge
        ↓
Memory
        ↓
Context
        ↓
Policy
        ↓
Execution (Runtime)
        ↓
Models
        ↓
Tools

En Identity wordt een capability die Context helpt opbouwen, net zoals Knowledge en Memory dat doen.

Ik denk eerlijk gezegd dat dit architectonisch sterker is dan het model waarmee we begonnen. Het laat zien dat Forge zich tijdens het schrijven verder ontwikkelt en dat we bereid zijn onze eigen aannames te verbeteren wanneer een consistenter model zich aandient. Dat is precies wat ik van een Architecture Bible zou verwachten.


Mooi. En hier wil ik een bewuste keuze maken.

Ik wil Identity niet als volgend hoofdstuk behandelen.

Waarom?

Omdat Identity slechts één capability is.

Security ook.

Observability ook.

Maar er is nog iets fundamentelers dat nog niet is beschreven:

Hoe wordt een AI-systeem eigenlijk intelligent?

Tot nu toe hebben we de bouwstenen beschreven:

* Knowledge
* Memory
* Context
* Policy
* Runtime

Maar er ontbreekt nog een concept dat ze allemaal verbindt.

Ik noem dat:

Chapter 11

Intelligence Engineering

Artificial Intelligence is not engineered by selecting better models.

It is engineered by creating better systems.

⸻

11.1 Beyond Language Models

The AI industry often measures intelligence by model capability.

Larger context windows.

Higher benchmark scores.

Improved reasoning.

More parameters.

These metrics describe language models.

They do not describe intelligent systems.

A model may possess extraordinary reasoning ability while remaining ineffective within an enterprise environment.

Conversely, a modest model operating within a well-engineered platform can consistently outperform larger models by receiving better context, stronger governance and higher-quality knowledge.

Forge therefore distinguishes between Model Intelligence and System Intelligence.

⸻

11.2 Model Intelligence vs. System Intelligence

A language model contributes intelligence.

The platform creates intelligence.

The distinction is fundamental.

Model Intelligence	System Intelligence
Statistical reasoning	Engineering discipline
Text generation	Decision support
General knowledge	Organizational knowledge
Pattern prediction	Governed execution
Probabilistic	Deterministic where possible

The language model is only one contributor to overall system intelligence.

⸻

11.3 Intelligence Is Emergent

Forge considers intelligence an emergent property.

It does not reside within any single capability.

Instead, it emerges from their interaction.

flowchart LR
Knowledge --> Context
Memory --> Context
Identity --> Context
Policy --> Context
Context --> Runtime
Runtime --> Model
Model --> Decisions
Decisions --> Humans

Remove any capability from this chain and overall system intelligence decreases.

⸻

11.4 The Intelligence Equation

Forge expresses this relationship conceptually.

Engineering Intelligence
=
Knowledge
×
Memory
×
Context
×
Policy
×
Execution
×
Model Capability

This equation is intentionally multiplicative.

If any factor approaches zero, overall system intelligence approaches zero.

A frontier model cannot compensate for missing governance.

Likewise, perfect governance cannot compensate for missing knowledge.

⸻

11.5 Intelligence Has Dimensions

Forge recognizes several dimensions of intelligence.

Cognitive Intelligence

Reasoning.

Planning.

Problem solving.

⸻

Organizational Intelligence

Policies.

Standards.

Architecture.

Institutional knowledge.

⸻

Operational Intelligence

Observability.

Runtime telemetry.

Evidence.

Performance.

⸻

Economic Intelligence

Cost optimization.

Provider selection.

Resource allocation.

⸻

Human Intelligence

Experience.

Judgment.

Ethics.

Creativity.

Enterprise AI succeeds only when these dimensions reinforce one another.

⸻

11.6 Intelligence Is Constrained

Many AI platforms attempt to maximize intelligence.

Forge attempts to optimize intelligence.

Maximum intelligence is rarely desirable.

Organizations require intelligence that is:

* safe;
* predictable;
* explainable;
* economical;
* compliant;
* trustworthy.

Unlimited capability without constraints becomes operational risk.

Constraints improve engineering quality.

⸻

11.7 Intelligence Is Governed

Every intelligent decision exists within organizational boundaries.

These boundaries include:

* regulation;
* policy;
* ethics;
* contractual obligations;
* security;
* privacy;
* business objectives.

The Runtime therefore never optimizes solely for task completion.

It optimizes for acceptable task completion.

⸻

11.8 Intelligence Is Measurable

Forge rejects vague claims such as:

“The AI became smarter.”

Instead, intelligence improvements should be observable.

Possible indicators include:

* decision accuracy;
* hallucination rate;
* context relevance;
* policy compliance;
* execution success;
* user satisfaction;
* engineering effort;
* operational cost.

Only measurable improvements should influence architectural evolution.

⸻

11.9 Intelligence Evolves

System intelligence continuously improves.

flowchart LR
Execution
Evaluation
Evidence
Knowledge
Policies
Execution --> Evaluation
Evaluation --> Evidence
Evidence --> Knowledge
Knowledge --> Policies
Policies --> BetterExecution
BetterExecution --> Execution

Every execution teaches the platform.

Not automatically.

Deliberately.

Learning without governance is not evolution.

It is drift.

⸻

11.10 Intelligence Engineering

Forge defines Intelligence Engineering as:

The engineering discipline responsible for designing, measuring, governing and continuously improving intelligent behaviour across the Engineering Operating System.

This discipline spans every capability.

It is therefore the integrating discipline of Forge.

⸻

11.11 Minimal Sufficient Intelligence

One of Forge’s core principles follows naturally from Token Stewardship and Compute Stewardship.

Use the minimum intelligence necessary to achieve the required engineering outcome.

This principle has several consequences.

The platform should prefer:

* retrieval over generation when facts exist;
* deterministic workflows over reasoning when possible;
* small models over frontier models when sufficient;
* cached results over repeated computation;
* human expertise over speculative inference.

More intelligence is not inherently better.

Appropriate intelligence is.

⸻

11.12 Intelligence Is an Organizational Capability

The greatest misconception in enterprise AI is that intelligence resides inside models.

Forge rejects this premise.

Organizations become intelligent when they successfully combine:

* trusted knowledge;
* intentional memory;
* engineered context;
* governed policy;
* deterministic execution;
* capable models;
* accountable humans.

Models participate in intelligence.

Organizations own it.

⸻

Closing Thought

Artificial Intelligence is often portrayed as a revolutionary technology.

Forge views it differently.

Intelligence is not a technology.

It is an organizational capability.

Language models make that capability accessible.

Engineering makes that capability sustainable.

Forge exists to ensure that organizations build intelligent systems they can understand, trust and evolve for decades rather than demonstrations that impress for minutes.

⸻

Ik denk dat we hier de echte kern van Forge hebben gevonden

Tot nu toe dacht ik dat Context Engineering het centrale hoofdstuk van de Architecture Bible zou zijn.

Ik denk inmiddels dat dat niet helemaal klopt.

De echte hiërarchie lijkt nu dit te zijn:

Knowledge Engineering
        ↓
Memory Engineering
        ↓
Context Engineering
        ↓
Policy Engineering
        ↓
Runtime Engineering
        ↓
Intelligence Engineering

Dat is geen lijst met componenten.

Dat is een engineering lifecycle.

En ik vind daar een interessante parallel in met klassieke softwareontwikkeling:

* Requirements → Knowledge
* Design → Context
* Constraints → Policy
* Execution → Runtime
* Outcome → Intelligence

Dat betekent dat Forge niet alleen een platform beschrijft, maar een volledige methode voor het ontwikkelen van AI-systemen. Ik denk dat dát uiteindelijk de grootste bijdrage van Forge kan zijn: niet een nieuw framework, maar een nieuwe engineeringdiscipline met een eigen vocabulaire, eigen principes en een samenhangende architectuur.

Mooi. Nu gaan we een hoofdstuk schrijven waarvan ik denk dat het uiteindelijk één van de meest geciteerde uit het boek wordt.

Tot nu toe hebben we de “cognitieve laag” beschreven:

* Knowledge
* Memory
* Context
* Policy
* Runtime
* Intelligence

Maar één vraag is nog onbeantwoord.

Hoe weet je dat een AI-systeem goed werkt?

Vrijwel de hele industrie antwoordt:

* benchmark scores;
* model rankings;
* hallucination rates.

Ik denk dat dat de verkeerde vraag is.

De juiste vraag is:

Heeft het systeem het juiste engineeringresultaat geleverd?

Dat brengt ons bij een discipline die ik nauwelijks uitgewerkt zie in bestaande AI-architecturen.

⸻

Chapter 12

Evidence Engineering

Software is trusted because it can be tested.

AI systems will be trusted because they can be explained.

Explanation requires evidence.

⸻

12.1 The Missing Discipline

Traditional software engineering has decades of experience producing evidence.

Unit tests.

Integration tests.

Coverage reports.

Performance benchmarks.

Security scans.

Deployment records.

Artificial Intelligence introduces an entirely new class of engineering questions.

Why was this model selected?

Why was this context assembled?

Why was one tool executed instead of another?

Why was human approval skipped?

Why was confidential information excluded?

These questions cannot be answered through logs alone.

They require evidence.

⸻

12.2 Logging Is Not Evidence

One of the largest misconceptions in enterprise AI is that logging automatically creates explainability.

It does not.

Logs record events.

Evidence explains decisions.

For example:

12:01:14
Model invoked.

This is a log.

It tells us that something happened.

Now consider:

Model: GPT-5
Reason:
Customer data classified as Internal.
Required reasoning capability: Medium.
Estimated cost below policy threshold.
Latency objective < 2 seconds.
Decision made by Model Selection Policy 4.2.

This is evidence.

It explains why something happened.

Engineering depends on the second.

⸻

12.3 Defining Evidence

Forge defines evidence as:

The structured explanation that allows an engineering decision to be reconstructed, understood and verified after execution.

Evidence therefore supports three activities:

* explanation;
* verification;
* accountability.

Evidence is not produced for machines.

It is produced for engineers.

⸻

12.4 Evidence Is a First-Class Artifact

Forge treats evidence as a primary engineering artifact.

Not as operational exhaust.

Every significant platform capability produces evidence.

flowchart LR
Identity --> Evidence
Knowledge --> Evidence
Memory --> Evidence
Context --> Evidence
Policy --> Evidence
Runtime --> Evidence
Model --> Evidence
Tools --> Evidence
Evidence --> EvidenceStore

The platform continuously explains itself.

⸻

12.5 Categories of Evidence

Different engineering decisions require different evidence.

Category	Example
Identity	Authenticated principal
Policy	Applied rules
Context	Included and excluded context
Knowledge	Authoritative sources
Runtime	Execution plan
Model	Selection rationale
Tool	Invocation history
Human	Approval records
Economics	Cost calculation
Evaluation	Quality assessment

Together these categories provide reconstructability.

⸻

12.6 Evidence Is Layered

Not every stakeholder requires the same evidence.

flowchart TB
Operations
Engineering
Security
Compliance
Audit
Executives
EvidenceStore --> Operations
EvidenceStore --> Engineering
EvidenceStore --> Security
EvidenceStore --> Compliance
EvidenceStore --> Audit
EvidenceStore --> Executives

The underlying evidence remains identical.

Only the presentation differs.

⸻

12.7 Explainability by Construction

Explainability should never be added after execution.

It must emerge naturally.

Forge therefore follows a simple rule.

Every engineering decision SHALL explain itself at the moment it is made.

Retrospective reconstruction is always less reliable than contemporaneous evidence generation.

⸻

12.8 Evidence Lifecycle

Evidence has its own lifecycle.

stateDiagram-v2
[*] --> Generated
Generated --> Validated
Validated --> Stored
Stored --> Indexed
Indexed --> Retrieved
Retrieved --> Archived
Archived --> Deleted
Deleted --> [*]

Retention policies depend upon organizational requirements.

⸻

12.9 Evidence Enables Trust

Enterprise AI depends upon trust.

Trust depends upon transparency.

Transparency depends upon evidence.

This relationship can be summarized simply.

Evidence
    ↓
Transparency
    ↓
Trust
    ↓
Adoption

Organizations rarely reject AI because models are incapable.

They reject AI because decisions cannot be justified.

⸻

12.10 Evidence and Continuous Learning

Evidence is valuable beyond auditing.

It enables engineering improvement.

Every execution becomes a learning opportunity.

flowchart LR
Execution
Evidence
Evaluation
Knowledge
Policy
Execution --> Evidence
Evidence --> Evaluation
Evaluation --> Knowledge
Knowledge --> Policy
Policy --> BetterExecution

The platform improves because evidence accumulates.

Not because models become larger.

⸻

12.11 Evidence Engineering

Forge defines Evidence Engineering as:

The engineering discipline responsible for generating, governing, preserving and utilizing evidence that explains intelligent system behaviour.

Evidence Engineering spans the entire Engineering Operating System.

Every capability contributes.

Every capability consumes it.

⸻

12.12 Evidence Before Trust

Forge deliberately reverses a common assumption.

Many organizations begin with trust.

Forge begins with evidence.

Trust is earned.

Evidence makes that possible.

Without evidence:

* governance weakens;
* debugging slows;
* compliance suffers;
* organizational learning stops.

Evidence therefore becomes one of the central architectural capabilities of Forge.

⸻

Closing Thought

Artificial Intelligence is often described as a “black box.”

Forge accepts that language models may remain probabilistic.

It rejects the idea that engineering should remain opaque.

Every significant platform decision should leave behind enough evidence that a competent engineer can understand not only what happened, but why it happened.

That is the foundation of trustworthy AI engineering.

⸻

Ik denk dat we nu een belangrijk inzicht hebben

Ik zie inmiddels dat Forge eigenlijk uit vier fundamentele engineeringstromen bestaat.

Knowledge Engineering
        │
        ▼
Memory Engineering
        │
        ▼
Context Engineering
        │
        ▼
Intelligence Engineering
        │
        ▼
Evidence Engineering

Dat voelt als een complete cyclus.

* Knowledge bepaalt wat waar is.
* Memory bepaalt wat blijft.
* Context bepaalt wat nu relevant is.
* Intelligence bepaalt wat gedaan wordt.
* Evidence bepaalt wat verklaard kan worden.

En hier valt me nog iets op.

De meeste AI-platformen eindigen bij de uitvoer van het model.

Forge eindigt bij bewijs.

Dat is volgens mij een fundamenteel verschil in denkwijze.

Niet:

“De AI gaf antwoord.”

Maar:

“De organisatie kan uitleggen waarom dat antwoord tot stand kwam.”

Ik denk dat dát uiteindelijk de essentie is van enterprise-grade AI.

Ik denk dat we nu op een interessant kruispunt staan.

Tot nu toe hebben we de interne architectuur van Forge opgebouwd:

* Knowledge
* Memory
* Context
* Policy
* Runtime
* Intelligence
* Evidence

Dat zijn allemaal capabilities binnen het Engineering Operating System.

Maar een operating system leeft niet in isolatie.

Een besturingssysteem heeft gebruikers.

Applicaties.

Drivers.

Hardware.

Interfaces.

Forge heeft precies hetzelfde.

Dat brengt ons bij wat ik zou beschouwen als de buitenste architectuurlaag.

⸻

Chapter 13

Capability Engineering

Systems are built from components.

Platforms are built from capabilities.

Engineering Operating Systems are defined by the contracts between them.

⸻

13.1 From Components to Capabilities

Traditional software architecture is largely component-oriented.

Applications are decomposed into services, libraries, modules or microservices.

These describe implementation.

Forge deliberately chooses a different abstraction.

Capabilities.

A capability represents a stable organizational responsibility.

Implementations may change.

Responsibilities should not.

This distinction enables architectures that evolve without constantly changing their conceptual model.

⸻

13.2 Defining a Capability

Forge defines a capability as:

An independently governed organizational responsibility that provides a stable contract while hiding implementation details.

Several aspects deserve emphasis.

A capability:

* owns a responsibility;
* exposes a contract;
* produces evidence;
* evolves independently;
* has explicit ownership;
* participates in governance.

Capabilities are not software components.

Software implements capabilities.

⸻

13.3 Why Capabilities Matter

Technology changes rapidly.

Organizational responsibilities change slowly.

Today:

* GPT-5

Tomorrow:

* GPT-8
* local reasoning models
* quantum inference
* technologies not yet invented.

None of these should require changing the Capability Model.

Capabilities provide architectural stability.

Implementations provide technical evolution.

⸻

13.4 Capability Contracts

Every capability communicates through contracts.

A contract defines:

* purpose;
* inputs;
* outputs;
* constraints;
* evidence;
* ownership;
* service guarantees.

For example:

Capability:
Context
Consumes:
Knowledge
Memory
Identity
Policy
Produces:
Context Package
Evidence:
Context Composition Report

Consumers interact with the contract.

Not the implementation.

⸻

13.5 Characteristics of a Well-Designed Capability

A Forge capability should exhibit several qualities.

It should be:

* cohesive;
* independently deployable;
* independently testable;
* observable;
* governed;
* replaceable;
* composable.

Capabilities should minimize internal complexity while maximizing external clarity.

⸻

13.6 Capability Composition

Capabilities rarely operate in isolation.

They compose into larger engineering behaviours.

flowchart LR
Knowledge --> Context
Memory --> Context
Identity --> Context
Policy --> Runtime
Context --> Runtime
Runtime --> Intelligence
Runtime --> Evidence

No capability owns the entire workflow.

Each contributes a specific responsibility.

⸻

13.7 Capability Boundaries

One of the most common architectural failures is responsibility leakage.

For example:

A Runtime should not decide policy.

A Knowledge capability should not manage execution.

Memory should not retrieve arbitrary documents.

Context should not authenticate users.

Each capability owns one responsibility.

Boundary violations create coupling.

Coupling reduces evolvability.

⸻

13.8 Capability Maturity

Capabilities evolve over time.

Forge recognizes four maturity levels.

Level	Characteristics
Experimental	Concepts under evaluation
Managed	Stable interfaces, evolving implementation
Governed	Policies, ownership and evidence established
Strategic	Critical organizational capability with long-term support

Maturity applies to capabilities rather than technologies.

⸻

13.9 Capability Lifecycle

Capabilities themselves have a lifecycle.

stateDiagram-v2
[*] --> Proposed
Proposed --> Designed
Designed --> Implemented
Implemented --> Governed
Governed --> Evolved
Evolved --> Deprecated
Deprecated --> Retired
Retired --> [*]

Architectures should evolve capabilities deliberately.

Not accidentally.

⸻

13.10 Capability Independence

Forge follows a strict architectural principle.

Capabilities communicate through contracts.

They never depend upon each other’s internal implementation.

For example:

The Runtime does not know whether Memory is implemented using:

* a vector database;
* PostgreSQL;
* Git;
* Redis;
* an object store.

It requests memory.

Nothing more.

This enables independent evolution.

⸻

13.11 Capability Taxonomy

Forge organizes capabilities into logical categories.

Cognitive Capabilities

* Knowledge
* Memory
* Context

Governance Capabilities

* Policy
* Identity
* Security
* Compliance

Execution Capabilities

* Runtime
* Tools
* Model Routing
* Planning

Operational Capabilities

* Audit
* Observability
* Evidence
* Economics

This taxonomy is conceptual rather than technical.

It helps architects reason about responsibility.

⸻

13.12 Capability Engineering

Forge defines Capability Engineering as:

The engineering discipline responsible for identifying, defining, governing and evolving the organizational capabilities that comprise the Engineering Operating System.

Capability Engineering precedes software design.

Before deciding how to build a system, engineers must determine what responsibilities the system should own.

⸻

13.13 Capabilities Outlive Technology

Perhaps the most important architectural principle in Forge is this:

Capabilities should survive technology replacement.

If replacing a language model requires redesigning the architecture, the architecture was too dependent on the technology.

If replacing a vector database changes the Capability Model, the abstraction was incorrect.

A well-designed capability remains stable while implementations evolve beneath it.

⸻

Closing Thought

Technology is transient.

Responsibilities endure.

Forge therefore builds its architecture around capabilities rather than products, protocols or vendors.

Capabilities provide the stable vocabulary through which organizations can reason about intelligent systems over decades, regardless of how rapidly the underlying technology changes.

⸻

Een suggestie voor de verdere structuur

Nu we de cognitieve en architecturale fundamenten hebben gelegd, denk ik dat de tweede helft van de Architecture Bible minder filosofisch en meer “hands-on” moet worden. Een logische volgorde zou zijn:

1. Identity Engineering – hoe een systeem weet wie handelt.
2. Security Engineering – hoe vertrouwelijkheid, integriteit en beschikbaarheid worden gewaarborgd.
3. Tool Engineering – hoe externe systemen als capabilities worden ontsloten (bijvoorbeeld via MCP).
4. Model Engineering – modelselectie, routing, evaluatie en provider-onafhankelijkheid.
5. Observability Engineering – metrics, traces, logs en operationele inzichten.
6. Economics Engineering – kosten, latency, tokengebruik en duurzaam gebruik van compute.
7. Governance Engineering – eigenaarschap, levenscycli, compliance en architectuurbeheer.

Vanaf dat punt kunnen we toewerken naar referentiearchitecturen, deploymentmodellen, capability catalogs en concrete implementatiepatronen. Daarmee verschuift het boek geleidelijk van waarom en wat naar hoe.


Perfect. Ik zou nu niet direct naar Identity gaan.

Er ontbreekt nog één hoofdstuk dat de overgang maakt van de abstracte architectuur naar de concrete implementatie.

Tot nu toe hebben we gesproken over wat het systeem moet doen. Nu moeten we beschrijven waaruit een AI-systeem daadwerkelijk wordt opgebouwd.

Dat hoofdstuk noem ik Model Engineering. Niet omdat modellen centraal staan, maar juist om uit te leggen waarom ze dat níét zijn.

⸻

Chapter 14

Model Engineering

Models provide capability.

Architecture provides value.

Engineering determines which model should be used, when, and why.

⸻

14.1 Beyond Model-Centric AI

The first generation of enterprise AI platforms was largely model-centric.

Organizations selected a model provider.

Applications were built around that provider’s APIs.

Business capabilities became tightly coupled to the characteristics of a specific model.

This approach proved successful for experimentation.

It does not scale to enterprise engineering.

Models evolve continuously.

Providers introduce new capabilities.

Costs fluctuate.

Regulations change.

An Engineering Operating System cannot depend on the lifecycle of a single model.

Forge therefore treats models as replaceable execution engines rather than architectural foundations.

⸻

14.2 Defining a Model

Forge defines a model as:

A computational capability that transforms structured context into probabilistic output.

Several aspects of this definition are deliberate.

A model:

* consumes context;
* performs inference;
* produces probabilistic results;
* owns no organizational state;
* owns no policy;
* owns no memory;
* owns no governance.

The model performs reasoning.

Everything else belongs to the platform.

⸻

14.3 Models Are Infrastructure

Traditional software treats databases as infrastructure.

Applications depend on data.

Not on the database engine.

Forge applies the same principle to AI.

Applications depend on intelligence.

Not on a specific model.

flowchart LR
Application --> Forge
Forge --> ModelInterface
ModelInterface --> GPT
ModelInterface --> Claude
ModelInterface --> Gemini
ModelInterface --> LocalModel
ModelInterface --> FutureModel

The application never communicates directly with a provider.

⸻

14.4 Model Contracts

Every model exposed through Forge implements a common contract.

A model accepts:

* Context Package
* Execution Objective
* Runtime Constraints
* Tool Availability
* Policy Constraints

A model returns:

* Response
* Confidence Indicators (where available)
* Token Usage
* Execution Metadata
* Provider Metadata

This abstraction enables portability without sacrificing provider-specific optimizations.

⸻

14.5 Model Capabilities

Not every model possesses identical capabilities.

Forge evaluates models across multiple dimensions.

Capability	Examples
Reasoning	Multi-step planning
Coding	Source code generation
Vision	Image understanding
Audio	Speech processing
Structured Output	JSON generation
Tool Use	Function calling
Context Window	Maximum supported context
Latency	Response time
Cost	Price per execution

Capability—not popularity—drives model selection.

⸻

14.6 Model Selection

Model selection is an engineering decision.

Not an application decision.

The Runtime requests a model.

The Model capability selects one according to policy.

Selection criteria may include:

* required reasoning capability;
* information classification;
* latency objectives;
* cost limits;
* geographical constraints;
* provider availability;
* regulatory requirements.

Applications should not encode these rules.

⸻

Figure 14.1 — Model Selection

flowchart LR
Runtime --> ModelRouter
Policy --> ModelRouter
Economics --> ModelRouter
Observability --> ModelRouter
ModelRouter --> SelectedModel

The model router becomes an engineering capability rather than a hard-coded switch statement.

⸻

14.7 Provider Independence

Vendor independence is a constitutional principle of Forge.

No application should depend directly on:

* OpenAI;
* Anthropic;
* Google;
* Mistral;
* Azure OpenAI;
* AWS Bedrock;
* Ollama;
* vLLM;
* future providers.

Providers compete.

Architecture endures.

Forge therefore standardizes interactions while preserving provider-specific strengths where appropriate.

⸻

14.8 Multi-Model Systems

Enterprise intelligence increasingly relies on multiple models.

Different tasks require different capabilities.

For example:

* a reasoning model for planning;
* a coding model for implementation;
* a vision model for document analysis;
* a lightweight model for classification;
* an on-premises model for confidential workloads.

The platform composes these capabilities into a single execution.

Applications remain unaware of the complexity.

⸻

14.9 Model Lifecycle

Models evolve rapidly.

Forge explicitly manages their lifecycle.

stateDiagram-v2
[*] --> Evaluated
Evaluated --> Approved
Approved --> Active
Active --> Monitored
Monitored --> Updated
Updated --> Approved
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Promotion between stages requires evidence rather than preference.

⸻

14.10 Model Evaluation

Selecting a model is only the beginning.

Forge continuously evaluates model behaviour using engineering metrics.

Examples include:

* task success rate;
* policy compliance;
* execution latency;
* cost per outcome;
* hallucination frequency;
* tool invocation accuracy;
* user satisfaction;
* regression against benchmark scenarios.

Evaluation is continuous.

Not a one-time procurement exercise.

⸻

14.11 Small Models First

Forge follows the principle of Minimal Sufficient Intelligence.

Where multiple models can achieve the required outcome, the platform should prefer the model with the lowest operational cost while satisfying quality objectives.

This often results in:

* deterministic execution before probabilistic reasoning;
* retrieval before generation;
* smaller models before larger models;
* local models before remote models when appropriate;
* cached responses before repeated inference.

Efficiency is an architectural concern.

⸻

14.12 Model Engineering

Forge defines Model Engineering as:

The engineering discipline responsible for selecting, evaluating, governing and evolving computational intelligence capabilities independently of application logic.

Model Engineering is not about finding the “best” model.

It is about finding the most appropriate model for a governed execution.

⸻

14.13 Models Are Replaceable

One of Forge’s architectural invariants is:

Replacing a model shall not require redesigning the platform.

If migrating from one provider to another requires widespread application changes, the abstraction has failed.

The Engineering Operating System absorbs provider variability so that applications remain stable.

⸻

Closing Thought

The history of computing has repeatedly shown that durable architectures separate stable responsibilities from rapidly evolving implementations.

Compilers replaced assembly languages.

Virtual machines abstracted hardware.

Containers abstracted operating systems.

Cloud platforms abstracted infrastructure.

Forge applies the same architectural principle to artificial intelligence.

Models will continue to evolve at extraordinary speed.

The Engineering Operating System ensures that organizations do not need to redesign their systems every time they do.

⸻

Een observatie over de structuur van het boek

Ik denk dat de Architecture Bible nu een duidelijke driedeling begint te krijgen:

Deel I – Foundations

* Why Forge Exists
* Engineering Operating System
* Design Philosophy
* Constitution

Deel II – Engineering Disciplines

* Knowledge Engineering
* Memory Engineering
* Context Engineering
* Policy Engineering
* Runtime Engineering
* Intelligence Engineering
* Evidence Engineering
* Capability Engineering
* Model Engineering

Deel III – Platform Capabilities

* Identity Engineering
* Security Engineering
* Tool Engineering
* Observability Engineering
* Economics Engineering
* Governance Engineering

Dat voelt inhoudelijk sterker dan een opsomming van technologieën. Het laat eerst de denkwereld van Forge zien, daarna de disciplines, en pas daarna de concrete capabilities die een implementatie vormen. Daardoor blijft het boek relevant, ook als de technologie over vijf of tien jaar ingrijpend verandert.

Ik zou nu een hoofdstuk kiezen dat de stap maakt van “intelligentie” naar de buitenwereld.

Tot nu toe hebben we gesproken over modellen. Maar modellen doen niets zonder iets om mee te communiceren.

Dat “iets” zijn tools.

En ik denk dat dit één van de grootste verschuivingen in AI-architectuur is.

Niet de LLM is het belangrijkste onderdeel.

De Tool Layer wordt het equivalent van de device drivers van een operating system.

Dat sluit bovendien perfect aan op MCP.

⸻

Chapter 15

Tool Engineering

Models generate language.

Tools create outcomes.

Engineering determines when intelligence becomes action.

⸻

15.1 Beyond Function Calling

Early generations of AI systems treated tools as an extension of prompting.

A model generated a function call.

The application executed it.

The result was returned to the model.

This approach proved effective for simple integrations.

It is insufficient for enterprise platforms.

Enterprise systems expose thousands of capabilities.

Identity.

ERP.

CRM.

Databases.

Messaging.

Ticketing.

Monitoring.

Infrastructure.

AI systems require a structured way to interact with these capabilities.

Forge therefore elevates tools to a first-class engineering capability.

⸻

15.2 Defining a Tool

Forge defines a tool as:

A governed interface through which the Engineering Operating System interacts with external capabilities.

Several elements are intentional.

A tool:

* represents a capability;
* hides implementation details;
* exposes a contract;
* enforces policy;
* produces evidence;
* participates in governance.

A tool is not merely an API wrapper.

It is an architectural abstraction.

⸻

15.3 Tools Are Not APIs

APIs describe communication.

Tools describe intent.

For example:

POST /users/123

This describes an API.

It says nothing about purpose.

Compare this with:

Create Employee Account

This is a tool.

It expresses an organizational capability.

Multiple APIs may implement the same tool.

Likewise, one API may expose multiple tools.

Forge therefore models capabilities rather than endpoints.

⸻

15.4 Tool Contracts

Every tool exposes an explicit contract.

A contract defines:

* purpose;
* required inputs;
* expected outputs;
* authorization requirements;
* policy constraints;
* evidence requirements;
* failure behaviour.

Applications invoke contracts.

Never implementations.

⸻

Figure 15.1 — Tool Contract

flowchart LR
Runtime --> Tool
Tool --> Contract
Contract --> API
API --> EnterpriseSystem

The Runtime understands the contract.

Only the Tool Layer understands the implementation.

⸻

15.5 Tool Discovery

The Runtime should never contain hard-coded knowledge of available tools.

Instead, tools are discoverable.

Each published tool advertises metadata including:

* identifier;
* description;
* owner;
* capabilities;
* required permissions;
* supported operations;
* version;
* policy requirements.

Discovery enables dynamic execution planning.

⸻

15.6 Tool Selection

Selecting a tool resembles selecting a model.

The Runtime requests a capability.

The Tool Layer selects an implementation.

Selection criteria may include:

* availability;
* latency;
* geographical location;
* policy restrictions;
* organizational ownership;
* health status;
* cost.

Applications should remain unaware of these details.

⸻

15.7 Tool Governance

Every tool has an owner.

Ownership includes responsibility for:

* lifecycle;
* documentation;
* security;
* availability;
* policy compliance;
* contract evolution.

No tool exists without ownership.

Ownership is an architectural requirement.

⸻

15.8 Tool Lifecycle

stateDiagram-v2
[*] --> Proposed
Proposed --> Designed
Designed --> Implemented
Implemented --> Published
Published --> Active
Active --> Updated
Updated --> Active
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Publishing a tool establishes an organizational contract.

Retiring a tool requires migration guidance.

⸻

15.9 Tool Composition

Complex objectives often require multiple tools.

flowchart LR
Runtime
Planner
Planner --> Calendar
Planner --> CRM
Planner --> Email
Planner --> ERP
Planner --> Identity
Calendar --> Runtime
CRM --> Runtime
Email --> Runtime
ERP --> Runtime
Identity --> Runtime

The Runtime coordinates.

Individual tools remain independent.

⸻

15.10 Tools Produce Evidence

Every invocation produces evidence.

Evidence includes:

* selected tool;
* version;
* caller;
* execution time;
* policy decisions;
* authorization results;
* inputs;
* outputs;
* failures.

Tool execution must be reconstructable.

⸻

15.11 Model Context Protocol

Forge views the Model Context Protocol (MCP) as the preferred standard for exposing tools.

MCP provides a standardized mechanism for advertising capabilities and exchanging structured requests between intelligent systems and external services.

Within Forge:

* MCP is an implementation standard.
* Tool Engineering is the architectural discipline.

This distinction is important.

Should another protocol emerge in the future, the Capability Model remains unchanged.

⸻

Figure 15.2 — MCP in Forge

flowchart LR
Runtime
ToolRegistry
MCP
GitHub
Jira
SAP
ServiceNow
Webex
Runtime --> ToolRegistry
ToolRegistry --> MCP
MCP --> GitHub
MCP --> Jira
MCP --> SAP
MCP --> ServiceNow
MCP --> Webex

MCP standardizes connectivity.

Forge governs capability.

⸻

15.12 Tool Engineering

Forge defines Tool Engineering as:

The engineering discipline responsible for designing, governing, publishing and evolving the external capabilities through which intelligent systems interact with their environment.

Tool Engineering transforms isolated APIs into governed organizational capabilities.

⸻

15.13 Capability Before Integration

One of Forge’s architectural principles is:

Model the capability before implementing the integration.

Engineers should first ask:

“What organizational capability are we exposing?”

Only then should they decide:

* REST;
* GraphQL;
* MCP;
* gRPC;
* messaging;
* event streams.

Protocols evolve.

Capabilities endure.

⸻

15.14 Tools as Device Drivers

Perhaps the strongest analogy in the Engineering Operating System is this:

Traditional Operating System	Forge
Device Driver	Tool
Hardware Device	Enterprise System
System Call	Tool Invocation
Driver Registry	Tool Registry
Kernel	Runtime
User Process	Application

Applications do not communicate directly with hardware.

They request operating system services.

Similarly, Forge applications do not communicate directly with enterprise systems.

They request organizational capabilities.

This abstraction dramatically simplifies portability, governance and testing.

⸻

Closing Thought

Language models can reason.

They cannot change the world.

Only tools can do that.

By elevating tools from simple API wrappers to governed organizational capabilities, Forge enables intelligent systems to act safely, consistently and transparently across the enterprise.

In an Engineering Operating System, tools are not accessories.

They are the bridge between intelligence and action.

⸻

Ik denk dat de analogie met een traditioneel besturingssysteem nu bijna compleet is

Tot nu toe kunnen we Forge bijna één-op-één vergelijken met een klassiek OS:

Operating System	Forge
Filesystem	Knowledge
Virtual Memory	Memory
Process Context	Context
Access Control	Policy
Kernel	Runtime
CPU	Language Models
Device Drivers	Tools (MCP)
Hardware	Enterprise Systems
Scheduler	Planner
System Logs	Evidence & Audit
Performance Monitor	Observability

Ik zou deze analogie niet alleen als illustratie gebruiken, maar als een terugkerend architectuurprincipe door het hele boek. Het maakt direct duidelijk wat Forge probeert te zijn: geen AI-framework, maar een besturingssysteem voor AI-engineering. Dat is een krachtig mentaal model dat lezers helpt om alle losse capabilities in één samenhangende architectuur te plaatsen.

Ik zou nu Identity Engineering behandelen.

Niet omdat identity op zichzelf interessanter is dan Security of Observability, maar omdat Identity de basis vormt voor vrijwel alle andere governance-capabilities. Zonder identiteit kun je geen beleid afdwingen, geen tools autoriseren, geen evidence toewijzen en geen accountability realiseren.

Maar ik zou Identity nadrukkelijk breder definiëren dan authenticatie.

⸻

Chapter 16

Identity Engineering

Authentication proves who requests execution.

Identity determines who is accountable for it.

Engineering requires both.

⸻

16.1 Beyond Authentication

Identity has traditionally been viewed as a security concern.

Users authenticate.

Services receive tokens.

Applications authorize requests.

Artificial Intelligence introduces a broader challenge.

Enterprise AI systems must identify not only people, but also:

* AI assistants;
* autonomous workflows;
* software agents;
* background services;
* enterprise systems;
* tools;
* datasets;
* organizations.

Identity therefore becomes a foundational engineering capability rather than a simple security mechanism.

⸻

16.2 Defining Identity

Forge defines Identity as:

The capability responsible for establishing, representing and governing the identity of every participant in the Engineering Operating System.

A participant may be:

* a human;
* an application;
* a model;
* a runtime;
* a tool;
* a service;
* an organization.

Everything that participates in execution possesses an identity.

⸻

16.3 Identity Is More Than Authentication

Authentication answers:

Who are you?

Authorization answers:

What may you do?

Identity answers a richer set of questions.

* Who initiated this execution?
* On whose behalf is it running?
* Which organization owns it?
* Which role does it represent?
* Which policies apply?
* Who is accountable for the outcome?

Identity therefore provides context for governance.

⸻

16.4 Identity as Context

Identity contributes directly to Context Engineering.

flowchart LR
Principal --> Identity
Identity --> Context
Knowledge --> Context
Memory --> Context
Policy --> Context
Context --> Runtime

Identity is not evaluated in isolation.

It enriches every execution with organizational meaning.

⸻

16.5 Identity Types

Forge distinguishes several identity classes.

Identity Type	Examples
Human	Employee, contractor, customer
Service	API, microservice
Runtime	Execution engine
Tool	MCP capability
Model	Hosted or local model
Organization	Business unit, tenant
Device	Laptop, mobile, room device
Workflow	Scheduled automation

Each type has distinct governance requirements.

⸻

16.6 Delegated Identity

Many AI executions occur on behalf of another identity.

For example:

An employee asks an assistant to create a Jira issue.

The Runtime invokes a Jira tool.

The Jira tool creates the issue.

Who performed the action?

Technically:

The Runtime.

Organizationally:

The employee.

Forge explicitly models delegated identity.

flowchart LR
Employee --> Runtime
Runtime --> Tool
Tool --> Jira
Employee -.Delegated Identity.-> Jira

Both identities remain visible in evidence.

⸻

16.7 Identity Chains

Enterprise executions frequently involve multiple identities.

Employee
↓
Assistant
↓
Runtime
↓
Tool
↓
Enterprise System

Forge preserves this complete identity chain.

This enables complete accountability across complex workflows.

⸻

16.8 Organizational Identity

Identity extends beyond individuals.

Organizations themselves possess identity.

Examples include:

* business units;
* ministries;
* customers;
* tenants;
* environments.

Organizational identity influences:

* policy evaluation;
* data residency;
* model selection;
* knowledge access;
* cost allocation.

Identity is therefore multi-dimensional.

⸻

16.9 Identity Federation

Forge assumes that identity is managed externally.

Examples include:

* Microsoft Entra ID;
* Okta;
* Keycloak;
* Ping Identity;
* government identity providers.

Forge consumes identity.

It does not replace identity providers.

This architectural boundary keeps responsibilities clear.

⸻

16.10 Identity Attributes

Identity extends beyond a unique identifier.

Relevant attributes include:

* organizational role;
* security clearance;
* department;
* tenant;
* location;
* assurance level;
* device posture;
* contractual relationship.

These attributes contribute to context construction and policy evaluation.

⸻

16.11 Identity Lifecycle

Identity evolves throughout its existence.

stateDiagram-v2
[*] --> Registered
Registered --> Verified
Verified --> Active
Active --> Updated
Updated --> Active
Active --> Suspended
Suspended --> Reactivated
Suspended --> Retired
Retired --> [*]

Identity governance continues throughout this lifecycle.

⸻

16.12 Identity Engineering

Forge defines Identity Engineering as:

The engineering discipline responsible for representing, governing and evolving the identities of every participant within the Engineering Operating System.

Identity Engineering ensures that intelligent systems always know who participates, who is responsible and on whose behalf actions are performed.

⸻

16.13 Identity Before Trust

One of Forge’s architectural invariants is:

Trust begins with identity.

Before evaluating:

* policy;
* permissions;
* context;
* tools;
* models;

the platform must first establish identity.

Without identity:

* accountability disappears;
* evidence loses meaning;
* governance becomes impossible.

Identity therefore forms the basis for trustworthy execution.

⸻

Closing Thought

Identity is often reduced to usernames, passwords and tokens.

Forge treats identity differently.

Identity provides the organizational meaning behind every execution.

It answers not only who initiated the request, but also who owns the responsibility, who bears accountability, and whose policies apply.

In an Engineering Operating System, identity is not merely a security primitive.

It is the foundation upon which trust, governance and accountability are built.

⸻

Een architectuurobservatie

Tijdens het schrijven van dit hoofdstuk valt me iets op dat de Forge-architectuur nog sterker maakt.

Er zijn eigenlijk drie verschillende “assen” in het platform:

De cognitieve as (hoe het systeem denkt):

* Knowledge
* Memory
* Context
* Intelligence

De uitvoeringsas (hoe het systeem handelt):

* Policy
* Runtime
* Models
* Tools

De governance-as (hoe het systeem beheerst wordt):

* Identity
* Security
* Evidence
* Observability
* Economics
* Governance

Dat is een krachtig mentaal model. In plaats van een platte lijst met capabilities krijgt de lezer drie orthogonale dimensies: denken, handelen en beheersen. Vrijwel iedere capability in Forge valt op natuurlijke wijze binnen één van deze assen, terwijl ze samen één coherent Engineering Operating System vormen.


Mooi. Ik zou nu Security Engineering behandelen.

Niet als een traditioneel securityhoofdstuk over firewalls, encryptie en IAM, maar als een fundamenteel architectuurprincipe van een Engineering Operating System.

En ik zou meteen een onderscheid maken dat ik nog nauwelijks zie in bestaande AI-architecturen:

Security beschermt assets.

Policy bestuurt gedrag.

Identity bepaalt verantwoordelijkheid.

Dat zijn drie verschillende capabilities.

⸻

Chapter 17

Security Engineering

Security is not the absence of attacks.

Security is the preservation of trust under adverse conditions.

⸻

17.1 Beyond Cybersecurity

Security has traditionally focused on protecting systems against unauthorized access.

Firewalls.

Encryption.

Authentication.

Network segmentation.

While these remain essential, intelligent systems introduce new attack surfaces.

Examples include:

* prompt injection;
* data poisoning;
* model manipulation;
* tool abuse;
* context leakage;
* retrieval poisoning;
* indirect prompt injection;
* model inversion;
* supply chain attacks.

Enterprise AI therefore requires Security Engineering that extends beyond traditional cybersecurity.

⸻

17.2 Defining Security

Forge defines Security as:

The capability responsible for preserving the confidentiality, integrity and availability of the Engineering Operating System and its assets throughout their lifecycle.

Security protects assets.

It does not decide policy.

It does not establish identity.

It protects the platform that enables them.

⸻

17.3 Security Protects Assets

Forge recognizes several categories of assets.

Asset	Examples
Knowledge	Architecture, documentation, ADRs
Memory	Long-term engineering memory
Context	Active Context Packages
Identity	Principals and credentials
Policies	Organizational rules
Models	Hosted and local models
Tools	Enterprise capabilities
Evidence	Audit records and execution evidence

Every capability owns assets.

Security protects them all.

⸻

17.4 Zero Trust by Design

Forge adopts Zero Trust as an architectural principle.

No request is trusted implicitly.

Every request must be evaluated continuously.

Trust is never inherited.

Trust is established through evidence.

This principle applies equally to:

* humans;
* services;
* runtimes;
* models;
* tools.

Zero Trust is therefore a platform-wide property rather than a network architecture.

⸻

17.5 Security as a Cross-Cutting Capability

Security is not implemented inside individual capabilities.

It surrounds them.

flowchart TB
Security
Security --> Identity
Security --> Knowledge
Security --> Memory
Security --> Context
Security --> Policy
Security --> Runtime
Security --> Models
Security --> Tools
Security --> Evidence

Every capability remains responsible for its function.

Security ensures those functions remain protected.

⸻

17.6 AI-Specific Threats

Forge explicitly recognizes AI-native attack vectors.

Examples include:

* Prompt Injection
* Indirect Prompt Injection
* Jailbreak Techniques
* Context Poisoning
* Knowledge Poisoning
* Tool Escalation
* Model Extraction
* Training Data Leakage
* Inference Manipulation
* Supply Chain Compromise

Traditional application security alone cannot mitigate these risks.

Security Engineering must understand AI execution.

⸻

17.7 Defense in Depth

No single control is sufficient.

Forge applies layered protection.

Identity
        ↓
Policy
        ↓
Context Validation
        ↓
Runtime Controls
        ↓
Tool Restrictions
        ↓
Model Constraints
        ↓
Evidence
        ↓
Monitoring

Each layer compensates for failures in the previous layer.

Security emerges from composition rather than individual controls.

⸻

17.8 Secure Context

Context represents one of the most valuable assets in the platform.

Security therefore protects:

* context composition;
* context transport;
* context storage;
* context expiration;
* context classification.

Context should contain only the information required for execution.

Excess context increases risk without improving outcomes.

This directly reinforces the principle of Minimal Sufficient Context.

⸻

17.9 Supply Chain Security

Modern AI systems depend upon numerous external components.

Examples include:

* model providers;
* embedding models;
* MCP servers;
* open-source libraries;
* prompt templates;
* vector databases;
* plugins.

Forge treats every dependency as part of the software supply chain.

Every dependency requires:

* provenance;
* version control;
* integrity validation;
* ownership;
* lifecycle management.

⸻

17.10 Security Produces Evidence

Every significant security decision generates evidence.

Examples include:

* denied executions;
* policy violations;
* context redactions;
* authentication events;
* authorization outcomes;
* tool restrictions;
* model restrictions.

Security therefore contributes directly to Evidence Engineering.

⸻

17.11 Secure by Construction

Forge follows a fundamental architectural principle.

Security should emerge naturally from platform architecture.

Not from after-the-fact controls.

Examples include:

* Identity before execution.
* Policy before tool invocation.
* Context validation before inference.
* Evidence before trust.
* Least privilege by default.
* Capability isolation.

Security becomes an intrinsic property rather than an optional feature.

⸻

17.12 Security Engineering

Forge defines Security Engineering as:

The engineering discipline responsible for protecting the assets, capabilities and execution pathways of the Engineering Operating System throughout their lifecycle.

Security Engineering spans every capability.

No capability exists outside its protection.

⸻

17.13 Security Supports Engineering

Security is often perceived as restricting innovation.

Forge adopts the opposite philosophy.

Well-designed security enables faster engineering.

Predictable security:

* reduces uncertainty;
* simplifies governance;
* accelerates deployment;
* improves trust;
* supports organizational adoption.

Security therefore becomes an engineering accelerator rather than an operational obstacle.

⸻

Closing Thought

Artificial Intelligence increases both capability and risk.

The objective of Security Engineering is not to eliminate risk.

It is to ensure that intelligent systems remain trustworthy even when operating in uncertain, adversarial or rapidly evolving environments.

In the Engineering Operating System, security is not a protective shell surrounding intelligence.

It is an intrinsic property woven into every capability, every execution and every engineering decision.

⸻

Ik zie nu een nog elegantere architectuur ontstaan

Tot nu toe hebben we drie assen geïntroduceerd:

Cognitive

* Knowledge
* Memory
* Context
* Intelligence

Execution

* Policy
* Runtime
* Models
* Tools

Governance

* Identity
* Security
* Evidence
* Observability
* Economics
* Governance

Maar ik denk dat er nog een vierde dimensie is die alles verbindt:

Cross-Cutting Qualities

* Trust
* Explainability
* Composability
* Portability
* Evolvability
* Observability
* Resilience

Dat zijn geen capabilities. Het zijn kwaliteitsattributen (ISO 25010 zou hier een goede referentie zijn) die je gebruikt om elke capability te beoordelen. Daarmee krijgt Forge niet alleen een capabilitymodel, maar ook een expliciet kwaliteitsmodel. Dat maakt de Architecture Bible niet alleen een beschrijving van wat Forge is, maar ook van waaraan goede AI-engineering moet voldoen. Dat is een mooie brug naar de resterende hoofdstukken over Observability, Economics en Governance, waar deze kwaliteitsattributen steeds opnieuw terugkomen.

Perfect. Ik denk dat we nu bij een discipline aankomen die in de meeste AI-platformen zwaar wordt onderschat.

Vrijwel iedereen verzamelt logs.

Sommigen verzamelen metrics.

Een enkeling gebruikt traces.

Maar een Engineering Operating System moet méér doen.

Het moet zichzelf kunnen begrijpen.

Dat brengt ons bij Observability Engineering.

⸻

Chapter 18

Observability Engineering

Monitoring tells us when systems fail.

Observability explains why they behave as they do.

Engineering requires understanding, not merely visibility.

⸻

18.1 Beyond Monitoring

Traditional monitoring answers questions such as:

* Is the service running?
* How much CPU is being used?
* What is the response time?
* Is the database available?

These remain important.

However, intelligent systems introduce fundamentally different questions.

Examples include:

* Why was this model selected?
* Why did planning take longer than expected?
* Which context source increased latency?
* Which policy produced this decision?
* Which tool caused execution failure?
* Why did cost suddenly increase?
* Which capability reduced answer quality?

These questions require observability rather than monitoring.

⸻

18.2 Defining Observability

Forge defines Observability as:

The capability responsible for understanding, measuring and explaining the behaviour of the Engineering Operating System through its externally observable signals.

Observability does not change system behaviour.

It explains it.

⸻

18.3 The Three Pillars Are Not Enough

Modern observability is often described using three pillars:

* Logs
* Metrics
* Traces

Forge considers these necessary but insufficient.

AI systems require additional signals.

Forge therefore expands the model.

Signal	Purpose
Logs	Events
Metrics	Quantitative measurements
Traces	Execution flow
Evidence	Engineering decisions
Context Lineage	Context composition
Policy Decisions	Governance
Model Telemetry	AI execution
Tool Telemetry	External interactions

Observability becomes a fusion of operational and cognitive telemetry.

⸻

Figure 18.1 — Observability Signals

flowchart LR
Logs
Metrics
Traces
Evidence
Context
Policy
Model
Tools
Logs --> Observability
Metrics --> Observability
Traces --> Observability
Evidence --> Observability
Context --> Observability
Policy --> Observability
Model --> Observability
Tools --> Observability

⸻

18.4 Observability Across Capabilities

Every capability contributes telemetry.

Capability	Observable Signals
Knowledge	Retrieval latency, freshness
Memory	Recall accuracy
Context	Size, composition time
Policy	Evaluation duration, decisions
Runtime	Execution lifecycle
Models	Token usage, latency
Tools	Success rate, failures
Identity	Authentication events
Security	Threat detections

No capability is invisible.

⸻

18.5 Context Lineage

One of the unique concepts in Forge is Context Lineage.

Every Context Package should answer:

* Which knowledge was included?
* Which memory was retrieved?
* Which policy removed information?
* Which identity attributes were added?
* Which context expired?
* Which context was ignored?

flowchart LR
Knowledge
Memory
Identity
Policy
Knowledge --> Context
Memory --> Context
Identity --> Context
Policy --> Context
Context --> Lineage

Context should always be reconstructable.

⸻

18.6 Execution Tracing

Traditional distributed tracing follows requests.

Forge extends this concept.

An execution trace follows:

* planning;
* context composition;
* policy evaluation;
* model routing;
* tool execution;
* evidence generation.

sequenceDiagram
Application->>Runtime: Execute
Runtime->>Context: Compose
Context->>Knowledge: Retrieve
Runtime->>Policy: Evaluate
Runtime->>Model: Infer
Runtime->>Tool: Execute
Runtime->>Evidence: Store
Runtime->>Application: Result

This provides complete execution visibility.

⸻

18.7 Engineering Telemetry

Operational metrics alone cannot improve AI systems.

Forge introduces Engineering Telemetry.

Examples include:

* context relevance score;
* policy evaluation frequency;
* memory utilization;
* retrieval precision;
* tool effectiveness;
* execution success rate;
* evidence completeness;
* planner accuracy.

These metrics directly support engineering improvements.

⸻

18.8 Observability Enables Optimization

Observability is not merely diagnostic.

It enables continuous optimization.

For example:

High latency may originate from:

* oversized context;
* excessive tool calls;
* slow retrieval;
* poor model routing.

Without observability, optimization becomes guesswork.

⸻

18.9 Feedback Loops

Every execution generates operational feedback.

flowchart LR
Execution
Telemetry
Analysis
Knowledge
Policy
Runtime
Execution --> Telemetry
Telemetry --> Analysis
Analysis --> Knowledge
Analysis --> Policy
Analysis --> Runtime

Observability therefore contributes directly to organizational learning.

⸻

18.10 Observability Engineering

Forge defines Observability Engineering as:

The engineering discipline responsible for measuring, explaining and continuously improving the behaviour of the Engineering Operating System through comprehensive operational and engineering telemetry.

Observability is not an operational concern alone.

It is an engineering discipline.

⸻

18.11 Explainability Through Observability

Explainability is often viewed as a model property.

Forge disagrees.

Explainability emerges from multiple capabilities.

* Evidence explains decisions.
* Context Lineage explains inputs.
* Policy explains constraints.
* Runtime traces explain execution.
* Tool telemetry explains actions.

Observability unifies these explanations into a coherent operational view.

⸻

18.12 Observability Before Optimization

One of Forge’s architectural principles is:

Optimize only what can be observed.

Premature optimization introduces complexity without improving outcomes.

Every optimization should be supported by observable evidence.

This reinforces another Forge principle:

Evidence precedes trust.

Likewise:

Observability precedes optimization.

⸻

Closing Thought

Enterprise AI systems cannot be engineered through intuition alone.

As systems become more autonomous, understanding their behaviour becomes more important than understanding individual models.

Forge therefore extends observability beyond infrastructure and applications into the cognitive architecture itself.

An Engineering Operating System should not merely report that something happened.

It should explain how, why and under which conditions it happened.

⸻

Een belangrijk inzicht

Terwijl we deze hoofdstukken schrijven, zie ik dat Forge steeds meer lijkt op een closed-loop engineering system.

Knowledge
      ↓
Memory
      ↓
Context
      ↓
Policy
      ↓
Runtime
      ↓
Models & Tools
      ↓
Evidence
      ↓
Observability
      ↓
Engineering Feedback
      ↓
Knowledge

Dat is geen lineaire architectuur meer.

Het is een continue verbeterlus.

En ik denk dat dit misschien wel de meest fundamentele eigenschap van Forge is.

De meeste AI-platformen optimaliseren de uitvoering.

Forge optimaliseert de engineering van de uitvoering.

Dat is een subtiel maar fundamenteel verschil. Het maakt Forge niet alleen een platform om AI-systemen te draaien, maar een platform dat organisaties helpt om die systemen systematisch beter te maken. Dat is precies wat je van een Engineering Operating System zou verwachten.

Mooi. Ik denk dat we nu bij een hoofdstuk komen dat vrijwel nergens als aparte engineeringdiscipline wordt behandeld, terwijl het voor enterprise AI essentieel is.

De meeste AI-platformen optimaliseren voor kwaliteit.

Sommige optimaliseren voor latency.

Bijna niemand behandelt kosten als een architectuurcapability.

Forge zou dat wel moeten doen.

Want uiteindelijk is compute geen gratis grondstof.

⸻

Chapter 19

Economics Engineering

Every token has a cost.

Every execution consumes resources.

Engineering is the discipline of achieving the desired outcome with the minimum necessary expenditure.

⸻

19.1 Beyond Cost Management

Traditional IT often treats cost optimization as an operational concern.

Budgets are reviewed monthly.

Infrastructure is optimized quarterly.

Procurement negotiates licensing annually.

Artificial Intelligence fundamentally changes this model.

Every execution has an immediate economic impact.

Every prompt.

Every retrieved document.

Every generated token.

Every model selection.

Every tool invocation.

Economics therefore becomes part of execution itself.

⸻

19.2 Defining Economics

Forge defines Economics as:

The capability responsible for measuring, governing and optimizing the economic efficiency of intelligent system execution.

Economics is not accounting.

It is engineering.

Its purpose is to maximize engineering value per unit of computational effort.

⸻

19.3 Intelligence Has a Cost

Every intelligent execution consumes resources.

Examples include:

* model inference;
* token generation;
* retrieval operations;
* vector searches;
* tool invocations;
* network bandwidth;
* storage;
* human approvals.

These costs accumulate.

The platform should understand them continuously rather than retrospectively.

⸻

19.4 Cost Is Multi-Dimensional

Enterprise AI incurs multiple forms of cost.

Cost Category	Examples
Compute	GPU, CPU, memory
Tokens	Prompt and completion usage
Storage	Knowledge, memory, evidence
Network	Data transfer
Licensing	Models, tools, SaaS
Human	Reviews, approvals
Latency	User waiting time
Opportunity	Delayed decisions

Engineering decisions influence all of these dimensions.

⸻

19.5 Economics Is a Runtime Concern

Economic decisions occur during execution.

For example:

Should the Runtime:

* invoke GPT-5?
* invoke a local reasoning model?
* reuse a cached result?
* retrieve knowledge instead of generating?
* request human approval?
* postpone execution?

These are engineering trade-offs.

Not procurement decisions.

⸻

Figure 19.1 — Economic Decision Flow

flowchart LR
Runtime --> Economics
Economics --> Policy
Policy --> Runtime
Runtime --> Model
Runtime --> Tool

Economics continuously informs execution.

⸻

19.6 Engineering Value

Forge distinguishes between cost and value.

A €0.02 execution that produces an unusable result is expensive.

A €2.00 execution that avoids a multi-day engineering incident may be inexpensive.

Economics therefore optimizes:

Value delivered per unit of computational effort.

Not lowest price.

⸻

19.7 Minimal Sufficient Intelligence

Earlier we introduced the principle of Minimal Sufficient Intelligence.

Economics Engineering operationalizes it.

The Runtime should prefer:

* deterministic execution before probabilistic reasoning;
* retrieval before generation;
* cached knowledge before inference;
* smaller models before larger models;
* fewer tool invocations before more.

Only when quality objectives remain satisfied.

Efficiency must never compromise correctness.

⸻

19.8 Economic Policies

Economic optimization is governed by policy.

Examples include:

* maximum cost per execution;
* preferred providers;
* regional execution requirements;
* daily token budgets;
* GPU allocation priorities;
* business-hour restrictions;
* premium model approval thresholds.

These policies allow organizations to align AI usage with business priorities.

⸻

19.9 Cost Attribution

Every execution should be attributable.

Examples include:

* user;
* project;
* department;
* tenant;
* application;
* capability;
* business process.

Without attribution there can be no optimization.

⸻

19.10 Economic Telemetry

Economics contributes specialized telemetry.

Examples include:

* cost per request;
* cost per successful outcome;
* token efficiency;
* model utilization;
* cache hit rate;
* retrieval efficiency;
* human intervention rate;
* execution ROI.

These metrics enable engineering rather than financial reporting.

⸻

19.11 Sustainability

Economic efficiency aligns closely with environmental sustainability.

Reducing unnecessary computation reduces:

* energy consumption;
* GPU utilization;
* infrastructure requirements;
* operational expenditure.

Forge therefore treats computational efficiency as both an economic and environmental objective.

Sustainability becomes an architectural quality rather than an external initiative.

⸻

19.12 Economics Engineering

Forge defines Economics Engineering as:

The engineering discipline responsible for maximizing organizational value through the efficient, governed and measurable use of computational intelligence.

Economics Engineering ensures that intelligent systems remain economically sustainable as they scale.

⸻

19.13 Compute Stewardship

One of Forge’s constitutional principles is:

Compute is a governed organizational resource.

Just as organizations manage:

* financial capital;
* human expertise;
* intellectual property;

they must also manage computational intelligence.

Compute should never be treated as unlimited.

Every execution should justify its resource consumption.

⸻

Closing Thought

Artificial Intelligence is often evaluated by what it can do.

Enterprise engineering must also evaluate what it costs to do it.

Organizations that optimize only for capability risk unsustainable growth.

Organizations that optimize only for cost sacrifice innovation.

Forge seeks the engineering balance between capability, governance and economic efficiency.

⸻

Ik denk dat hier een mooie symmetrie ontstaat

Als ik nu naar alle disciplines kijk, zie ik dat Forge eigenlijk vier fundamentele engineeringcycli heeft.

Knowledge Cycle

* Knowledge
* Memory
* Context

Execution Cycle

* Policy
* Runtime
* Models
* Tools

Trust Cycle

* Identity
* Security
* Evidence
* Observability

Optimization Cycle

* Economics
* Feedback
* Continuous Improvement

Dat is interessant, omdat deze vier cycli samen vrijwel alle vragen beantwoorden die een architect of engineer stelt:

* Wat weten we? → Knowledge Cycle
* Hoe voeren we het uit? → Execution Cycle
* Kunnen we het vertrouwen? → Trust Cycle
* Kunnen we het beter doen? → Optimization Cycle

Ik denk dat dit uiteindelijk de hoofdstukstructuur van Deel II zou moeten worden. Het geeft de lezer niet alleen losse disciplines, maar een samenhangend mentaal model voor AI-engineering. Daarmee onderscheidt Forge zich van de meeste AI-boeken, die vaak georganiseerd zijn rond technologieën in plaats van rond de fundamentele verantwoordelijkheden van een engineeringplatform.


Perfect. We hebben nu de cognitieve, uitvoerings-, trust- en optimalisatiecycli beschreven.

Nu ontbreekt nog één discipline die alles samenhoudt.

Niet technisch.

Maar organisatorisch.

En dat is misschien wel de moeilijkste van allemaal.

De meeste AI-projecten falen namelijk niet omdat GPT niet slim genoeg is.

Ze falen omdat niemand eigenaar is.

Omdat verantwoordelijkheden onduidelijk zijn.

Omdat standaarden ontbreken.

Omdat wijzigingen ongecontroleerd plaatsvinden.

Dat brengt ons bij misschien wel de hoogste discipline binnen Forge.

⸻

Chapter 20

Governance Engineering

Technology enables intelligent systems.

Governance enables organizations to trust them.

Engineering ensures they continue to deserve that trust.

⸻

20.1 Beyond Compliance

Governance is frequently misunderstood as compliance.

Policies.

Audits.

Approvals.

Documentation.

These are important.

But they are consequences of governance.

Not governance itself.

Forge views governance differently.

Governance is the capability that ensures the Engineering Operating System evolves intentionally rather than accidentally.

⸻

20.2 Defining Governance

Forge defines Governance as:

The capability responsible for directing, controlling and continuously improving the Engineering Operating System throughout its lifecycle.

Governance is not operational management.

It is architectural stewardship.

⸻

20.3 Governance Answers Different Questions

Security asks:

Is it protected?

Policy asks:

Is it permitted?

Observability asks:

What happened?

Governance asks:

Should the platform evolve in this direction?

This distinction is essential.

Governance concerns the evolution of the platform itself.

⸻

20.4 Governance Is Continuous

Governance does not occur only during architecture reviews.

Every engineering decision contributes to governance.

Examples include:

* approving a new model;
* publishing a tool;
* changing a policy;
* retiring knowledge;
* introducing a capability;
* modifying context composition;
* changing execution workflows.

Governance accompanies change.

⸻

Figure 20.1 — Governance Across the Platform

flowchart TB
Governance
Governance --> Knowledge
Governance --> Memory
Governance --> Context
Governance --> Policy
Governance --> Runtime
Governance --> Models
Governance --> Tools
Governance --> Identity
Governance --> Security
Governance --> Evidence
Governance --> Economics

Governance spans every capability.

No capability exists outside governance.

⸻

20.5 Governance Requires Ownership

Every architectural artifact requires an owner.

Examples include:

* capabilities;
* policies;
* tools;
* models;
* knowledge repositories;
* architectural decisions;
* engineering standards.

Ownership implies responsibility for:

* evolution;
* quality;
* lifecycle;
* documentation;
* evidence.

Artifacts without owners become organizational risk.

⸻

20.6 Architectural Decision Making

Forge adopts Architecture Decision Records (ADRs) as a primary governance mechanism.

Every significant architectural decision should record:

* the problem;
* considered alternatives;
* selected option;
* rationale;
* consequences;
* review date.

Architectural memory becomes explicit rather than implicit.

⸻

20.7 Capability Stewardship

Every capability should have a designated steward.

Responsibilities include:

* defining the contract;
* approving changes;
* monitoring quality;
* reviewing evidence;
* coordinating evolution;
* communicating roadmap.

Capabilities evolve through stewardship.

Not through consensus alone.

⸻

20.8 Governance Metrics

Governance itself should be measurable.

Possible indicators include:

* capability maturity;
* policy coverage;
* evidence completeness;
* architectural debt;
* ADR adoption;
* ownership coverage;
* standards compliance;
* lifecycle health.

Governance without measurement becomes opinion.

⸻

20.9 Controlled Evolution

One of the primary goals of governance is enabling safe evolution.

Every significant change should answer:

* Why is change required?
* Which capabilities are affected?
* Which contracts change?
* Which risks increase?
* Which evidence supports the decision?
* How will success be measured?

Change becomes an engineered process rather than an improvised activity.

⸻

20.10 Governance and Architecture

Architecture defines structure.

Governance preserves architectural integrity.

Without governance:

* abstractions erode;
* responsibilities overlap;
* shortcuts accumulate;
* technical debt grows;
* platform consistency declines.

Governance protects architectural intent over time.

⸻

20.11 Federated Governance

Enterprise AI rarely belongs to a single team.

Forge supports federated governance.

Global governance establishes:

* constitutional principles;
* platform standards;
* capability contracts;
* mandatory policies.

Local governance adapts these standards to:

* business units;
* departments;
* projects;
* products.

flowchart TB
GlobalGovernance
BusinessUnit
Project
Application
GlobalGovernance --> BusinessUnit
BusinessUnit --> Project
Project --> Application

Federation balances consistency with autonomy.

⸻

20.12 Governance Engineering

Forge defines Governance Engineering as:

The engineering discipline responsible for directing, governing and continuously evolving the Engineering Operating System while preserving architectural integrity, organizational accountability and long-term sustainability.

Governance Engineering is the discipline that ensures Forge remains coherent as it grows.

⸻

20.13 Governance Produces Confidence

Governance is often perceived as slowing engineering.

Forge takes the opposite view.

Good governance enables:

* faster onboarding;
* safer experimentation;
* repeatable engineering;
* consistent quality;
* predictable evolution.

The objective is not control.

The objective is confidence.

⸻

20.14 Architecture Is a Living System

Forge rejects the notion that architecture is static documentation.

Architecture evolves continuously.

Knowledge evolves.

Policies evolve.

Models evolve.

Tools evolve.

Organizations evolve.

Governance ensures that these independent evolutions strengthen rather than fragment the platform.

⸻

Closing Thought

An Engineering Operating System is never finished.

Its capabilities mature.

Its architecture evolves.

Its implementations change.

Governance provides the discipline that allows this evolution to occur without sacrificing coherence, trust or engineering quality.

Forge therefore treats governance not as an administrative process, but as a core engineering capability that preserves the platform’s integrity across years of continuous change.

⸻

Ik denk dat we nu de volledige Engineering Operating System hebben

Als ik een stap terug doe, zie ik dat de Architecture Bible nu een complete en consistente architectuur beschrijft.

                 ENGINEERING OPERATING SYSTEM
                 ┌───────────────────────────────┐
                 │        Governance             │
                 └──────────────┬────────────────┘
                                │
        ┌────────────────────────────────────────────┐
        │          Knowledge Engineering             │
        └────────────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────────────────┐
        │           Memory Engineering               │
        └────────────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────────────────┐
        │           Context Engineering              │
        └────────────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────────────────┐
        │            Policy Engineering              │
        └────────────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────────────────┐
        │           Runtime Engineering              │
        └────────────────────────────────────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                             ▼
   Model Engineering            Tool Engineering
          │                             │
          └──────────────┬──────────────┘
                         ▼
        ┌────────────────────────────────────────────┐
        │      Intelligence in Action                │
        └────────────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────────────────┐
        │ Evidence • Observability • Economics       │
        └────────────────────────────────────────────┘
                         │
                         ▼
        Continuous Learning and Evolution

Ik zou zelfs overwegen om dit niet alleen als hoofdstukstructuur te gebruiken, maar als de officiële Forge Reference Architecture.

Dat zou Forge een herkenbaar architectuurmodel geven, vergelijkbaar met:

* het OSI-model voor netwerken;
* het Purdue Model voor OT;
* TOGAF voor enterprise-architectuur;
* het Kubernetes Control Plane-model.

Een eigen referentiearchitectuur maakt Forge niet alleen een verzameling ideeën, maar een architectuurtaal waarmee teams, leveranciers en organisaties over AI-engineering kunnen communiceren. Dat is precies het soort fundament dat een Engineering Operating System verdient.

Ik denk dat we nu een belangrijk moment bereiken.

Tot nu toe hebben we vrijwel alle fundamentele disciplines beschreven. Vanaf hier zou de Architecture Bible kunnen vervallen in losse hoofdstukken over deployment, Kubernetes, security controls, CI/CD, enzovoort.

Dat zou een gemiste kans zijn.

Er ontbreekt namelijk nog één discipline die eigenlijk boven alle andere staat.

Niet een capability.

Niet een component.

Maar de manier waarop je een AI-systeem ontwerpt.

Dat hoofdstuk zou ik bewust vóór deployment en implementatie plaatsen.

⸻

Chapter 21

Architecture Engineering

Software can be built without architecture.

Platforms cannot.

Engineering Operating Systems demand intentional architectural design.

⸻

21.1 Architecture Before Implementation

One of the recurring mistakes in AI engineering is beginning with technology.

Organizations ask:

* Which model should we use?
* Which vector database is best?
* Which agent framework should we adopt?
* Should we use LangGraph or AutoGen?

These are implementation questions.

Architecture begins elsewhere.

Architecture asks:

* Which responsibilities exist?
* Which capabilities own them?
* Which contracts connect them?
* Which qualities must the system preserve over time?

Technology follows architecture.

Not the reverse.

⸻

21.2 Defining Architecture

Forge defines architecture as:

The intentional organization of capabilities, responsibilities and constraints to enable the long-term evolution of intelligent systems.

Architecture is therefore not a diagram.

Nor is it documentation.

Architecture is the set of engineering decisions that remain stable while implementations evolve.

⸻

21.3 Architecture Defines Boundaries

A well-designed architecture answers one question repeatedly:

Who owns this responsibility?

For example:

Knowledge owns truth.

Memory owns retention.

Context owns composition.

Policy owns decisions.

Runtime owns execution.

Evidence owns explanation.

Observability owns understanding.

If ownership is ambiguous, the architecture is incomplete.

⸻

21.4 Architecture Is About Change

Many believe architecture exists to describe today’s system.

Forge rejects this idea.

Architecture exists to enable tomorrow’s changes.

A successful architecture minimizes the cost of:

* replacing models;
* introducing tools;
* changing policies;
* adding capabilities;
* reorganizing teams;
* complying with new regulations.

Architecture is an investment in future adaptability.

⸻

Figure 21.1 — Stable vs. Variable

flowchart TB
subgraph Stable
CapabilityModel
Contracts
Principles
Governance
end
subgraph Variable
Models
Providers
Protocols
Infrastructure
Frameworks
end

Architecture protects the stable core while allowing continuous evolution of the variable edge.

⸻

21.5 Architecture Is Layered

Forge deliberately separates architectural concerns into layers.

Layer	Responsibility
Principles	Why the platform exists
Capability Model	What responsibilities exist
Contracts	How capabilities interact
Runtime	How execution occurs
Implementations	Which technologies realize capabilities

Each layer depends on the one above it.

Never the reverse.

⸻

21.6 Architecture Is Constraint

A common misconception is that architecture maximizes freedom.

Good architecture deliberately introduces constraints.

Examples include:

* Repository as the Source of Truth.
* Context before Runtime.
* Policy before Execution.
* Capability before Implementation.
* Evidence before Trust.

Constraints reduce accidental complexity.

Freedom without constraints produces inconsistency.

⸻

21.7 Architectural Invariants

Forge introduces the concept of Architectural Invariants.

An invariant is a principle that must remain true regardless of implementation.

Examples include:

* Every execution has an identity.
* Every decision produces evidence.
* Every capability has an owner.
* Every policy is explainable.
* Every contract is versioned.
* Every change is governed.

Technologies may evolve.

Invariants do not.

⸻

21.8 Architecture as a Language

Architecture is more than structure.

It is vocabulary.

Forge intentionally introduces precise terminology.

Examples include:

* Context Package
* Engineering Memory
* Capability Contract
* Runtime Kernel
* Policy Decision
* Evidence Package
* Context Lineage
* Engineering Operating System

Shared language enables shared understanding.

Shared understanding enables consistent engineering.

⸻

21.9 Architecture Decision Records

Architectural decisions should never disappear into meeting notes.

Every significant architectural decision should become an ADR.

Each ADR captures:

* Context
* Decision
* Alternatives
* Rationale
* Consequences
* Review triggers

Architecture therefore becomes continuously traceable.

⸻

21.10 Reference Architecture

Forge distinguishes between:

* conceptual architecture;
* logical architecture;
* reference architecture;
* implementation architecture.

The Reference Architecture describes the canonical arrangement of capabilities.

Implementations adapt it to local requirements without violating its architectural invariants.

This separation allows consistency without rigidity.

⸻

21.11 Architecture Quality

Architecture should be evaluated using explicit quality attributes.

Examples include:

* Evolvability
* Composability
* Explainability
* Portability
* Resilience
* Testability
* Operability
* Security
* Economic Efficiency

These qualities apply across every capability.

No implementation should optimize one quality while unknowingly degrading another.

⸻

21.12 Architecture Engineering

Forge defines Architecture Engineering as:

The engineering discipline responsible for designing, governing and evolving the structural integrity of the Engineering Operating System through stable capabilities, explicit contracts and enduring architectural principles.

Architecture Engineering is the discipline that connects vision to implementation.

⸻

21.13 Architecture Outlives Technology

History repeatedly demonstrates that technologies change faster than architectures.

Programming languages change.

Cloud providers change.

Container platforms change.

AI models change.

Protocols change.

Organizations change.

A successful architecture survives all of them.

Forge therefore considers technological neutrality not merely a deployment strategy, but an architectural objective.

⸻

Closing Thought

Architecture is often mistaken for documentation produced at the beginning of a project.

Forge views architecture differently.

Architecture is the continuously evolving structure that allows intelligent systems to change without losing coherence.

It provides the discipline that transforms collections of technologies into sustainable engineering platforms.

The Engineering Operating System is ultimately not defined by the models it executes, but by the architecture that allows those models—and everything around them—to evolve safely over time.

⸻

Ik denk dat we nu de “Theory of Forge” hebben afgerond

Als ik naar de laatste twintig hoofdstukken kijk, zie ik eigenlijk drie duidelijke delen ontstaan:

Part I — Foundations

* Why Forge Exists
* Engineering Operating System
* Vision
* Design Philosophy
* Constitution

Part II — Engineering Theory

* Knowledge Engineering
* Memory Engineering
* Context Engineering
* Runtime Engineering
* Policy Engineering
* Intelligence Engineering
* Evidence Engineering
* Capability Engineering
* Model Engineering
* Tool Engineering
* Identity Engineering
* Security Engineering
* Observability Engineering
* Economics Engineering
* Governance Engineering
* Architecture Engineering

Part III — Engineering Practice

En eerlijk gezegd denk ik dat Part III fundamenteel anders moet worden dan alles wat we tot nu toe hebben geschreven.

Niet meer filosofisch.

Niet meer conceptueel.

Maar volledig praktisch:

* Reference Architecture
* Capability Specifications
* Canonical APIs
* Context Package Specification
* Memory Specification
* Knowledge Specification
* Runtime Specification
* Policy Language
* MCP Integration
* GitHub Integration
* LangGraph Mapping
* LiteLLM Mapping
* Deployment Blueprints
* Kubernetes
* Azure
* AWS
* On-Premises
* ADR Catalog
* Best Practices
* Anti-Patterns
* Testing
* Certification

Dat zou de Architecture Bible dezelfde opbouw geven als een combinatie van een RFC, het Kubernetes Architecture Book en TOGAF: eerst de theorie, daarna de normatieve specificatie, en tenslotte de referentie-implementaties. Ik denk dat dát uiteindelijk de meest duurzame vorm voor Forge is.

Ik denk dat we nu een bewuste koerswijziging moeten maken.

We hebben de theorie vrijwel afgerond. Als we nu doorgaan met nog meer disciplines, wordt het boek topzwaar. Een goed architectuurboek maakt op een gegeven moment de overgang van “wat is het?” naar “hoe ziet het er concreet uit?”

Ik zou daarom Part III beginnen met misschien wel het belangrijkste hoofdstuk van het hele boek.

Niet deployment.

Niet Kubernetes.

Niet APIs.

Maar de Forge Reference Architecture.

Dit wordt het hoofdstuk waar alle voorgaande theorie samenkomt in één normatief architectuurmodel.

⸻

Part III

The Forge Reference Architecture

Chapter 22

Forge Reference Architecture

Architecture becomes useful only when it can be implemented.

The Reference Architecture translates principles into structure without prescribing technology.

⸻

22.1 Purpose

The Forge Reference Architecture defines the canonical structure of the Engineering Operating System.

It is neither:

* an implementation;
* a deployment topology;
* a product architecture;
* a cloud architecture.

Instead, it describes the stable arrangement of capabilities that every Forge implementation should preserve.

Implementations may vary.

Architectural intent should not.

⸻

22.2 Architectural Objectives

The Reference Architecture has five primary objectives.

1. Separate responsibilities.
2. Enable independent evolution.
3. Preserve governance.
4. Support multiple deployment models.
5. Remain vendor independent.

These objectives guide every architectural decision.

⸻

22.3 The Layered Architecture

Forge consists of six logical layers.

Applications
──────────────────────────────
Business Capabilities
──────────────────────────────
Engineering Operating System
──────────────────────────────
Execution Platform
──────────────────────────────
Infrastructure
──────────────────────────────
Physical Resources

Applications consume business capabilities.

Business capabilities consume the Engineering Operating System.

The Engineering Operating System consumes execution services.

Execution services consume infrastructure.

Infrastructure consumes physical resources.

Each layer hides complexity from the layer above.

⸻

22.4 The Engineering Operating System

The EOS contains the canonical Forge capabilities.

+------------------------------------------------------------+
                 ENGINEERING OPERATING SYSTEM
+------------------------------------------------------------+
Governance
Knowledge
Memory
Context
Policy
Runtime
Identity
Security
Model Routing
Tool Registry
Evidence
Observability
Economics
+------------------------------------------------------------+

Everything inside this boundary belongs to Forge.

Everything outside consumes Forge.

⸻

22.5 Capability Dependencies

Capabilities interact through explicit contracts.

flowchart LR
Knowledge --> Context
Memory --> Context
Identity --> Context
Policy --> Runtime
Context --> Runtime
Runtime --> Models
Runtime --> Tools
Runtime --> Evidence
Evidence --> Observability
Observability --> Governance

Dependencies always flow downward.

Circular dependencies are prohibited.

⸻

22.6 Architectural Boundaries

Forge defines four architectural boundaries.

Organizational Boundary

Separates organizational knowledge from external systems.

⸻

Platform Boundary

Separates Forge capabilities from business applications.

⸻

Execution Boundary

Separates deterministic execution from probabilistic inference.

⸻

Integration Boundary

Separates enterprise capabilities from implementation protocols.

These boundaries reduce coupling.

⸻

22.7 Internal vs External Capabilities

Forge distinguishes internal capabilities from external services.

Internal

* Context
* Memory
* Policy
* Runtime
* Evidence

External

* Models
* Identity Providers
* Enterprise Systems
* MCP Servers
* Cloud Providers

Forge governs the interaction.

It does not own every dependency.

⸻

22.8 Capability Interaction Rules

The Reference Architecture introduces mandatory interaction rules.

Examples include:

Runtime SHALL NOT access Knowledge directly.

Context SHALL compose Knowledge.

Policy SHALL NOT invoke Models.

Models SHALL NOT invoke enterprise APIs directly.

Tools SHALL mediate external interactions.

Evidence SHALL be generated before execution completes.

These rules preserve architectural integrity.

⸻

22.9 Canonical Request Flow

Every execution follows the same logical sequence.

Application
↓
Objective
↓
Identity
↓
Knowledge
↓
Memory
↓
Context
↓
Policy
↓
Runtime Planning
↓
Model Selection
↓
Tool Execution
↓
Evidence
↓
Observability
↓
Result

This becomes the canonical execution pipeline.

⸻

22.10 Reference Interfaces

Each capability exposes only contracts.

For example:

Knowledge
Retrieve()
Publish()
Validate()
Search()
Memory
Store()
Recall()
Expire()
Archive()
Context
Compose()
Validate()
Serialize()
Inspect()
Runtime
Plan()
Execute()
Resume()
Cancel()

Notice that these are capability operations.

Not REST endpoints.

⸻

22.11 Deployment Neutrality

The Reference Architecture deliberately avoids deployment assumptions.

Forge may operate:

* on-premises;
* private cloud;
* public cloud;
* sovereign cloud;
* edge environments;
* hybrid deployments.

The Capability Model remains identical.

⸻

22.12 Technology Mapping

Implementations map technologies onto capabilities.

Capability	Possible Technologies
Runtime	LangGraph, Temporal, Dapr
Model Routing	LiteLLM, custom router
Knowledge	Git, SharePoint, Confluence
Memory	PostgreSQL, Redis, vector DB
Policy	OPA, Cedar
Tools	MCP, REST, gRPC
Identity	Entra ID, PingFederate, Keycloak
Observability	OpenTelemetry, Prometheus
Evidence	PostgreSQL, Object Storage

This mapping is informative.

It is not normative.

⸻

22.13 Architectural Invariants

Every Forge implementation SHALL preserve the following invariants.

* Every execution has identity.
* Every execution has policy.
* Every execution has context.
* Every decision produces evidence.
* Every capability has ownership.
* Every interaction uses explicit contracts.
* Every capability is observable.
* Every architectural change is governed.

Technology choices shall not violate these invariants.

⸻

22.14 Conformance

An implementation conforms to the Forge Reference Architecture when it preserves:

* capability responsibilities;
* capability contracts;
* architectural invariants;
* governance principles;
* engineering disciplines.

Technology choices do not determine conformance.

Architectural behaviour does.

⸻

Closing Thought

The purpose of the Forge Reference Architecture is not to standardize implementation.

It is to standardize thinking.

Organizations should be free to adopt different technologies, deployment models and operational practices while sharing a common architectural language.

In this way, the Reference Architecture becomes the stable foundation upon which diverse implementations can evolve without fragmenting the engineering discipline.

⸻

Een voorstel voor het vervolg: van architectuur naar specificaties

Vanaf hier zou ik de rest van het boek veel normatiever maken, bijna in de stijl van een RFC of de Kubernetes API Specification. In plaats van beschrijvende hoofdstukken krijgen we hoofdstukken die exacte specificaties bevatten.

Bijvoorbeeld:

* Chapter 23 – Capability Specifications: formele definities van elke capability, inclusief verantwoordelijkheden, interfaces, lifecycle en kwaliteitsattributen.
* Chapter 24 – Context Package Specification: een normatief formaat voor een Context Package (vergelijkbaar met een protocolspecificatie).
* Chapter 25 – Execution Specification: toestandsmodel, lifecycle en contracten van de Runtime.
* Chapter 26 – Capability Contracts: standaardinterfaces en interactiepatronen.
* Chapter 27 – Deployment Reference Models: mapping van de referentiearchitectuur naar Kubernetes, Azure, AWS, on-premises en hybride omgevingen.

Dat zou de tweede helft van de Architecture Bible transformeren van een visie naar een echte technische standaard die teams kunnen implementeren en waarop producten of projecten kunnen worden getoetst.

Perfect. Ik denk dat we nu bij het hoofdstuk komen dat de rest van de Architecture Bible gaat definiëren.

Tot nu toe hebben we capabilities beschreven.

Nu moeten we ze specificeren.

Dat betekent dat we overstappen van architectuur naar een normatieve specificatie.

Vanaf hier zou ik bewust de stijl veranderen.

Niet meer essay-achtig.

Maar vergelijkbaar met:

* RFC 9110
* Kubernetes API Specification
* OpenTelemetry Specification
* OCI Image Specification

Met normatieve taal:

* SHALL
* SHOULD
* MAY
* MUST NOT

Dat maakt Forge niet alleen een visie, maar een implementeerbare standaard.

⸻

Chapter 23

Capability Specification Framework

Architecture defines responsibilities.

Specifications define expectations.

Engineering requires both.

⸻

23.1 Purpose

This chapter defines the canonical specification framework used throughout Forge.

Every Forge capability SHALL be specified using a common structure.

The purpose of this structure is to ensure that all capabilities are:

* consistent;
* comparable;
* testable;
* governable;
* independently evolvable.

A capability without a specification is considered incomplete.

⸻

23.2 Normative Language

The key words SHALL, SHALL NOT, SHOULD, SHOULD NOT, and MAY are to be interpreted as described in RFC 2119 and RFC 8174.

Forge uses these keywords to distinguish mandatory architectural requirements from recommended implementation guidance.

⸻

23.3 Canonical Capability Structure

Every capability specification SHALL contain the following sections.

1. Purpose
2. Responsibilities
3. Scope
4. Inputs
5. Outputs
6. Contracts
7. Lifecycle
8. Evidence
9. Quality Attributes
10. Security Considerations
11. Observability
12. Conformance Requirements

No section may be omitted.

⸻

23.4 Capability Identifier

Every capability SHALL possess a globally unique identifier.

Example:

forge.capability.context
forge.capability.memory
forge.capability.runtime
forge.capability.policy

Identifiers remain stable across versions.

Names MAY evolve.

Identifiers SHALL NOT.

⸻

23.5 Purpose

The Purpose section answers one question:

Why does this capability exist?

It SHALL describe responsibility.

It SHALL NOT describe implementation.

Example:

Purpose
The Context capability composes the minimum
sufficient context required for governed execution.

⸻

23.6 Responsibilities

Responsibilities SHALL define what the capability owns.

Responsibilities SHALL NOT overlap with another capability.

Example:

Context SHALL:

* compose context;
* validate context;
* serialize context;
* expose context lineage.

Context SHALL NOT:

* retrieve knowledge directly from repositories;
* evaluate policies;
* execute models.

Responsibility boundaries are normative.

⸻

23.7 Scope

Each capability SHALL explicitly define:

Included responsibilities.

Excluded responsibilities.

Dependencies.

Example:

Included
Context composition
Context validation
Context serialization
Excluded
Knowledge retrieval
Runtime execution
Policy evaluation

⸻

23.8 Contracts

Every capability SHALL expose one or more contracts.

A contract SHALL define:

* operation;
* purpose;
* required inputs;
* outputs;
* constraints;
* errors;
* evidence generated.

Example:

ComposeContext()
Inputs
Knowledge
Memory
Identity
Policies
Outputs
Context Package
Evidence
Context Lineage

Contracts define behaviour.

Not transport protocols.

⸻

23.9 Lifecycle

Every capability SHALL define its lifecycle.

Example:

stateDiagram-v2
[*] --> Proposed
Proposed --> Active
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Lifecycle transitions SHALL be governed.

⸻

23.10 Evidence Requirements

Every capability SHALL specify which evidence it produces.

Example:

Capability	Evidence
Context	Composition Report
Runtime	Execution Report
Policy	Decision Report
Models	Selection Report
Tools	Invocation Report

Evidence requirements are mandatory.

⸻

23.11 Quality Attributes

Every capability SHALL define measurable quality attributes.

Examples include:

Availability

Latency

Reliability

Portability

Explainability

Composability

Testability

Evolvability

Resilience

Security

Quality attributes SHALL be measurable.

⸻

23.12 Security Considerations

Each capability SHALL define:

Assets protected.

Threats.

Trust assumptions.

Required controls.

Security dependencies.

This section is normative.

⸻

23.13 Observability

Every capability SHALL expose observable signals.

Examples include:

Metrics.

Logs.

Traces.

Evidence.

Health.

Performance indicators.

No capability SHALL be opaque.

⸻

23.14 Conformance

A capability conforms to Forge when:

its responsibilities are complete;

its contracts are explicit;

its lifecycle is defined;

its evidence is produced;

its quality attributes are measurable;

its architectural invariants are preserved.

⸻

23.15 Capability Template

Every capability SHALL follow the canonical template.

Capability
Identifier
Purpose
Responsibilities
Scope
Dependencies
Inputs
Outputs
Contracts
Lifecycle
Evidence
Quality Attributes
Security Considerations
Observability
Conformance

This template becomes normative for all future Forge specifications.

⸻

23.16 Specification Before Implementation

One of the central principles of Forge is:

Every capability SHALL be specified before it is implemented.

Implementation reveals technical choices.

Specifications preserve architectural intent.

The specification therefore becomes the contract between architecture and engineering.

⸻

Closing Thought

Capabilities define what the Engineering Operating System is responsible for.

Specifications define how those responsibilities are expressed, constrained and verified.

By introducing a common specification framework, Forge transforms architecture from descriptive documentation into an executable engineering standard.

⸻

Een voorstel om Forge nog sterker te maken

Ik denk dat we op een punt zijn gekomen waarop Forge zich echt kan onderscheiden van bestaande architectuurmodellen.

In plaats van alleen hoofdstukken te schrijven, kunnen we elk capability-hoofdstuk afsluiten met een normatieve Capability Specification, vergelijkbaar met een API-specificatie.

Bijvoorbeeld voor Context:

* Conceptueel hoofdstuk (waarom Context Engineering bestaat).
* Formele specificatie (forge.capability.context).
* Referentie-interface (Compose, Validate, Serialize, Inspect).
* State machine.
* Canonieke datamodellen.
* Conformance tests.
* Referentie-implementaties (bijvoorbeeld LangGraph, Semantic Kernel of een eigen Forge Runtime).

Dat maakt de tweede helft van de Architecture Bible veel meer dan een boek. Het wordt een architectuurstandaard waarop implementaties kunnen worden gebouwd en gecertificeerd. Dat sluit ook goed aan bij je oorspronkelijke ambitie om Forge niet als framework, maar als een Engineering Operating System en een duurzame engineeringdiscipline neer te zetten.

Ik zou nu direct doorgaan naar de eerste echte normatieve specificatie. En die moet niet willekeurig gekozen zijn.

Alles in Forge draait uiteindelijk om Context.

Als Context niet goed gedefinieerd is, worden Runtime, Memory, Policy en Models allemaal vaag.

Daarom zou de eerste specificatie zijn:

⸻

Chapter 24

Context Package Specification

Specification Identifier: forge.spec.context.v1

Status: Draft Standard

⸻

Intelligence is determined not only by reasoning, but by the quality of the context presented to that reasoning.

⸻

24.1 Purpose

This specification defines the canonical representation of a Forge Context Package.

A Context Package SHALL be the only mechanism through which execution receives operational context.

No Runtime SHALL access Knowledge, Memory or Identity directly.

All required information SHALL be composed into a Context Package before execution begins.

⸻

24.2 Design Goals

The Context Package SHALL be:

* deterministic;
* immutable;
* explainable;
* serializable;
* portable;
* versioned;
* observable;
* reproducible.

A Context Package SHALL describe everything required for governed execution.

⸻

24.3 Canonical Structure

Every Context Package SHALL contain the following logical sections.

Context Package
Metadata
Identity
Objective
Knowledge
Memory
Policies
Constraints
Execution Preferences
Evidence Metadata
Lineage

No mandatory section may be omitted.

⸻

24.4 Metadata

Metadata identifies the package.

Required fields:

id:
version:
created:
expires:
schema:
tenant:
classification:
correlationId:

The package SHALL possess a globally unique identifier.

⸻

24.5 Identity Section

Identity SHALL contain execution identity.

Example

identity:
user:
roles:
groups:
tenant:
clearance:
authentication:
delegation:

Identity SHALL NOT contain credentials.

⸻

24.6 Objective Section

Every execution SHALL define a single objective.

objective:
goal:
expectedOutcome:
priority:
urgency:
confidenceTarget:

The objective represents intent.

Not implementation.

⸻

24.7 Knowledge Section

Knowledge SHALL contain immutable information.

Example

knowledge:
documents:
facts:
standards:
architecture:
policies:
references:

Knowledge SHALL reference source repositories.

Knowledge SHALL NOT duplicate repository ownership.

⸻

24.8 Memory Section

Memory SHALL contain execution-relevant experience.

Example

memory:
conversation:
preferences:
historicalExecutions:
workingMemory:
learnedPatterns:

Memory MAY expire.

Knowledge SHALL NOT.

⸻

24.9 Policy Section

Policy SHALL contain evaluated policy decisions.

Example

policy:
allowedTools:
blockedTools:
allowedModels:
blockedModels:
informationClassification:
approvalRequired:

The Runtime SHALL consume policy decisions.

It SHALL NOT evaluate policy.

⸻

24.10 Constraints

Constraints define execution boundaries.

Example

constraints:
maxTokens:
maxCost:
maxLatency:
allowedRegions:
allowedProviders:
retention:

Constraints SHALL be deterministic.

⸻

24.11 Execution Preferences

Preferences influence optimization.

Example

preferences:
preferredModel:
reasoningLevel:
streaming:
language:
temperature:
parallelism:

Preferences SHALL NOT violate policy.

⸻

24.12 Evidence Metadata

Evidence metadata enables traceability.

evidence:
requestId:
sessionId:
executionId:
auditLevel:
retention:

Every Context Package SHALL generate evidence.

⸻

24.13 Context Lineage

Context SHALL explain its own composition.

Example

lineage:
knowledgeSources:
memorySources:
policySources:
identitySource:
filteredItems:
excludedItems:
timestamps:

Context Lineage SHALL allow complete reconstruction of the composition process.

⸻

24.14 Immutability

After publication:

The Context Package SHALL NOT change.

If modification is required:

A new package SHALL be composed.

This guarantees reproducibility.

⸻

24.15 Serialization

The Context Package SHALL support deterministic serialization.

Supported encodings MAY include:

* JSON
* CBOR
* Protocol Buffers

Transport encoding is implementation-specific.

Logical structure is normative.

⸻

24.16 Context Lifecycle

stateDiagram-v2
[*] --> Composing
Composing --> Validated
Validated --> Published
Published --> Executing
Executing --> Archived
Archived --> Expired
Expired --> [*]

Every lifecycle transition SHALL generate evidence.

⸻

24.17 Validation Rules

A Context Package SHALL be rejected when:

* Identity is missing.
* Objective is missing.
* Policy is absent.
* Version is unsupported.
* Schema validation fails.
* Required evidence identifiers are absent.

Validation SHALL occur before execution.

⸻

24.18 Architectural Invariants

Every Context Package SHALL satisfy the following invariants.

* Exactly one Objective.
* Exactly one Identity.
* One or more Knowledge sources.
* Zero or more Memory entries.
* One evaluated Policy.
* One Lineage record.
* One Metadata section.

⸻

24.19 Example

context:
metadata:
id: ctx-203847
version: 1
tenant: ministry-finance
classification: BBN2
identity:
user: alice
roles:
- architect
objective:
goal: Generate Architecture Decision Record
knowledge:
- ADR-001
- TOGAF Standard
memory:
- Previous ADR discussions
policy:
allowedModels:
- forge-gpt
constraints:
maxTokens: 12000
maxCost: 1.50
lineage:
knowledgeSources:
- Git Repository
memorySources:
- Forge Memory

⸻

24.20 Conformance

A Context Package conforms to this specification when:

* all mandatory sections are present;
* schema validation succeeds;
* lifecycle transitions are respected;
* lineage is complete;
* evidence identifiers exist;
* architectural invariants are preserved.

⸻

Closing Thought

A Context Package is more than a collection of inputs.

It is the canonical representation of everything an intelligent system is permitted to know, remember, consider and use at the moment of execution.

By standardizing the Context Package, Forge establishes a stable contract between composition and execution, allowing runtimes, models and tools to evolve independently while preserving reproducibility, explainability and governance.

⸻

Een aanbeveling voor de volgende stap

Vanaf hier zou ik de specificaties opbouwen zoals internetstandaarden dat doen: eerst de gegevensmodellen, daarna de gedragsmodellen.

Dat betekent dat de volgende hoofdstukken logisch zouden zijn:

* Chapter 25 – Memory Specification (forge.spec.memory.v1): geheugenmodellen, retention, expiratie, working memory versus long-term memory.
* Chapter 26 – Knowledge Specification (forge.spec.knowledge.v1): repositories, provenance, versioning, publicatie en validatie.
* Chapter 27 – Runtime Execution Specification (forge.spec.runtime.v1): toestandsmachine, planning, scheduling, resumability, retries en foutafhandeling.

Die volgorde sluit aan bij de natuurlijke afhankelijkheden in Forge: eerst de canonieke gegevensstructuren, daarna de uitvoeringssemantiek die daarop vertrouwt.

Perfect. Dan volgen we dezelfde opbouw als internet- en cloudstandaarden: eerst de canonieke data-objecten, daarna de execution-semantiek.

Na de Context Specification is de volgende logische stap de definitie van Memory.

⸻

Chapter 25

Memory Specification

Specification Identifier: forge.spec.memory.v1

Status: Draft Standard

⸻

Knowledge describes what is true.

Memory preserves what was experienced.

Intelligence requires both.

⸻

25.1 Purpose

This specification defines the canonical representation, lifecycle and behaviour of Engineering Memory within the Forge Engineering Operating System.

Memory SHALL preserve information that improves future execution.

Memory SHALL NOT become an alternative source of truth.

Truth belongs to Knowledge.

Experience belongs to Memory.

⸻

25.2 Design Principles

Engineering Memory SHALL be:

* persistent;
* attributable;
* explainable;
* versioned;
* observable;
* governed;
* expirable;
* reproducible.

Memory SHALL always exist to improve future execution.

⸻

25.3 Memory Model

Forge defines five canonical memory domains.

Engineering Memory
├── Working Memory
├── Session Memory
├── Episodic Memory
├── Semantic Memory
└── Organizational Memory

Each domain serves a distinct responsibility.

⸻

25.4 Working Memory

Working Memory represents temporary execution state.

Examples include:

* intermediate reasoning;
* execution variables;
* planner state;
* unresolved dependencies;
* temporary tool outputs.

Characteristics:

* execution scoped;
* mutable;
* automatically discarded after completion unless explicitly promoted.

Working Memory SHALL NOT survive execution by default.

⸻

25.5 Session Memory

Session Memory preserves interactions within a logical conversation or workflow.

Examples include:

* conversation history;
* previous user requests;
* temporary preferences;
* unresolved follow-up tasks.

Characteristics:

* session scoped;
* ordered;
* expirable.

Session Memory SHALL expire according to governance policy.

⸻

25.6 Episodic Memory

Episodic Memory records completed executions.

Examples include:

* previous investigations;
* completed workflows;
* generated reports;
* architectural decisions;
* incident resolutions.

Each episode SHALL contain:

* objective;
* context identifier;
* execution outcome;
* evidence references;
* timestamp.

Episodes SHALL remain immutable once recorded.

⸻

25.7 Semantic Memory

Semantic Memory represents learned abstractions derived from repeated experience.

Examples include:

* preferred document structures;
* organizational terminology;
* frequently used workflows;
* successful planning strategies;
* reusable execution patterns.

Semantic Memory SHALL be derived from multiple episodes.

It SHALL NOT be created from isolated events.

⸻

25.8 Organizational Memory

Organizational Memory captures institutional knowledge that is not yet formalized into Knowledge repositories.

Examples include:

* recurring engineering practices;
* operational conventions;
* architectural heuristics;
* lessons learned;
* organizational preferences.

Organizational Memory SHOULD eventually migrate into canonical Knowledge when validated.

Memory is not intended to permanently replace Knowledge.

⸻

25.9 Memory Object

Every Memory Object SHALL contain:

memory:
id:
type:
owner:
created:
expires:
source:
confidence:
classification:
references:
content:
evidence:

Every object SHALL possess a globally unique identifier.

⸻

25.10 Memory Lifecycle

stateDiagram-v2
[*] --> Created
Created --> Validated
Validated --> Active
Active --> Archived
Archived --> Expired
Expired --> Deleted

Lifecycle transitions SHALL generate evidence.

Deletion SHALL comply with retention policies.

⸻

25.11 Memory Promotion

Not all memories deserve persistence.

Forge defines Memory Promotion.

Working
↓
Session
↓
Episode
↓
Semantic
↓
Knowledge

Promotion SHALL require evidence.

Automatic promotion MAY occur through policy.

⸻

25.12 Memory Retrieval

Memory retrieval SHALL be objective-driven.

Inputs MAY include:

* identity;
* objective;
* context;
* similarity;
* temporal proximity;
* organizational scope.

Retrieval SHALL produce relevance scores.

Ranking algorithms are implementation-specific.

⸻

25.13 Memory Expiration

Every Memory Object SHALL define a retention policy.

Possible strategies include:

* time-based expiration;
* event-driven expiration;
* policy-based expiration;
* manual review.

Expired memory SHALL NOT participate in execution.

Unless explicitly restored.

⸻

25.14 Memory Provenance

Every Memory Object SHALL record provenance.

Example:

provenance:
execution:
context:
knowledge:
tools:
author:
timestamp:
evidence:

Memory without provenance SHALL be considered untrusted.

⸻

25.15 Memory Governance

Governance SHALL define:

* ownership;
* retention;
* classification;
* review frequency;
* deletion policy;
* promotion policy.

Memory SHALL never exist without ownership.

⸻

25.16 Security

Memory SHALL inherit classification from originating evidence.

Access SHALL be identity-aware.

Memory retrieval SHALL respect:

* authorization;
* tenant boundaries;
* information classification;
* legal restrictions.

Memory SHALL NOT bypass Policy.

⸻

25.17 Observability

Memory SHALL expose:

Metrics:

* retrieval count;
* hit rate;
* promotion rate;
* expiration rate;
* average age;
* relevance score.

Evidence:

* retrieval reports;
* promotion reports;
* expiration reports.

⸻

25.18 Architectural Invariants

Every Forge implementation SHALL preserve:

* Memory is never authoritative.
* Memory is attributable.
* Memory is versioned.
* Memory is observable.
* Memory is governed.
* Memory is explainable.
* Memory is policy-controlled.
* Memory is promotable.

⸻

25.19 Conformance

A Memory implementation conforms when it:

* supports the canonical lifecycle;
* implements mandatory metadata;
* records provenance;
* respects governance;
* exposes observability signals;
* produces evidence;
* preserves architectural invariants.

⸻

Closing Thought

Memory allows intelligent systems to improve through experience without compromising the integrity of organizational knowledge.

By separating memory from knowledge, Forge enables learning while preserving trust, ensuring that every remembered experience remains attributable, governed and explainable.

⸻

Een architectonische observatie

Op dit punt begint een diepere structuur zichtbaar te worden. De specificaties volgen niet alleen een technische volgorde, maar ook een cognitieve hiërarchie:

Knowledge
      │
      ▼
Memory
      │
      ▼
Context
      │
      ▼
Policy
      │
      ▼
Runtime
      │
      ▼
Models / Tools
      │
      ▼
Evidence
      │
      ▼
Observability
      │
      ▼
Governance

Dit is meer dan een datastroom; het is een normatieve afhankelijkheidsketen. Elke laag mag alleen afhankelijk zijn van de lagen erboven. Daardoor ontstaan duidelijke verantwoordelijkheden, minimale koppeling en de mogelijkheid om capabilities onafhankelijk te ontwikkelen of te vervangen.

Ik zou deze keten expliciet benoemen als de Forge Cognitive Pipeline. Net zoals het OSI-model de referentie is voor netwerkarchitectuur, kan deze pipeline de fundamentele referentie worden voor AI-engineering binnen Forge. Hij geeft implementatieteams een helder antwoord op de vraag: in welke volgorde worden verantwoordelijkheden opgebouwd en waarom? Dat maakt Forge niet alleen een verzameling specificaties, maar een consistente architectuurtaal.

Mooi. We hebben nu de twee fundamentele gegevensobjecten gedefinieerd:

* forge.spec.context.v1
* forge.spec.memory.v1

De volgende specificatie moet eigenlijk Knowledge zijn. Niet omdat een LLM kennis nodig heeft, maar omdat de hele Engineering Operating System rust op een gecontroleerde bron van waarheid.

Dit is ook waar Forge zich duidelijk onderscheidt van de meeste RAG-frameworks. Die behandelen documenten als data. Forge behandelt Knowledge als een governed enterprise asset.

⸻

Chapter 26

Knowledge Specification

Specification Identifier: forge.spec.knowledge.v1

Status: Draft Standard

⸻

Knowledge is not information.

Knowledge is information that has been validated, governed and accepted as organizational truth.

⸻

26.1 Purpose

This specification defines the canonical representation, lifecycle and governance of Knowledge within the Forge Engineering Operating System.

Knowledge SHALL represent organizational truth.

Knowledge SHALL NOT represent temporary observations, experiences or assumptions.

Knowledge exists to provide stable, governed information for intelligent execution.

⸻

26.2 Design Principles

Knowledge SHALL be:

* authoritative;
* versioned;
* attributable;
* governed;
* immutable once published;
* discoverable;
* explainable;
* independently verifiable.

Knowledge SHALL always originate from a controlled source.

⸻

26.3 Knowledge Domains

Forge distinguishes multiple Knowledge Domains.

Knowledge
├── Architecture
├── Policies
├── Standards
├── Procedures
├── Reference Data
├── Domain Models
├── Documentation
└── External References

Each domain SHALL define its own governance process.

⸻

26.4 Knowledge Object

Every Knowledge Object SHALL contain the following metadata.

knowledge:
id:
title:
owner:
domain:
classification:
version:
status:
created:
published:
reviewDate:
source:
checksum:
references:
evidence:

Knowledge identifiers SHALL remain stable throughout the lifecycle of the object.

⸻

26.5 Knowledge Content

A Knowledge Object SHALL contain:

* structured content;
* unstructured content;
* metadata;
* provenance;
* relationships.

The internal representation MAY vary.

The logical structure SHALL remain consistent.

⸻

26.6 Knowledge Lifecycle

stateDiagram-v2
[*] --> Draft
Draft --> Review
Review --> Approved
Approved --> Published
Published --> Superseded
Superseded --> Archived
Archived --> Retired
Retired --> [*]

No Knowledge Object SHALL bypass the approval process.

⸻

26.7 Source of Truth

Every Knowledge Object SHALL reference exactly one authoritative source.

Examples include:

* Git repository;
* Enterprise Architecture Repository;
* SharePoint document library;
* Policy management system;
* CMDB;
* Product catalog.

Replicated copies SHALL NOT become authoritative.

⸻

26.8 Knowledge Relationships

Knowledge Objects MAY reference other Knowledge Objects.

Relationship types include:

* supersedes;
* dependsOn;
* derivedFrom;
* references;
* complements;
* conflictsWith.

These relationships SHALL be explicit.

⸻

26.9 Versioning

Knowledge SHALL use semantic versioning.

Example:

Architecture Standard
v1.0.0
↓
v1.1.0
↓
v2.0.0

Breaking changes SHALL increment the major version.

Minor editorial corrections SHOULD increment the patch version.

⸻

26.10 Knowledge Validation

Knowledge SHALL be validated before publication.

Validation MAY include:

* schema validation;
* link validation;
* ownership verification;
* approval verification;
* classification verification;
* reference validation.

Knowledge failing validation SHALL NOT be published.

⸻

26.11 Knowledge Provenance

Every Knowledge Object SHALL preserve provenance.

Example:

provenance:
author:
reviewer:
approver:
repository:
commit:
timestamp:
evidence:

Knowledge without provenance SHALL be considered untrusted.

⸻

26.12 Retrieval

Knowledge retrieval SHALL be deterministic.

Inputs MAY include:

* objective;
* domain;
* classification;
* version;
* language;
* organizational scope.

Ranking algorithms MAY be implementation-specific.

Returned objects SHALL preserve their metadata.

⸻

26.13 Knowledge Publication

Publication SHALL create an immutable published version.

Subsequent modifications SHALL produce a new version.

Published knowledge SHALL remain reproducible.

⸻

26.14 Governance

Every Knowledge Domain SHALL define:

* owner;
* review cycle;
* approval process;
* publication workflow;
* archival policy;
* retention policy.

Knowledge SHALL never exist without ownership.

⸻

26.15 Security

Knowledge SHALL define:

* classification;
* visibility;
* tenant scope;
* access policy;
* export policy.

Retrieval SHALL respect Policy decisions.

Knowledge SHALL NOT enforce authorization independently.

⸻

26.16 Observability

Knowledge SHALL expose:

Metrics:

* retrieval frequency;
* publication frequency;
* review age;
* validation failures;
* stale documents;
* dependency graph health.

Evidence:

* publication reports;
* review reports;
* validation reports.

⸻

26.17 Architectural Invariants

Every Forge implementation SHALL preserve:

* Knowledge is authoritative.
* Knowledge is immutable after publication.
* Knowledge is versioned.
* Knowledge is attributable.
* Knowledge is governed.
* Knowledge is observable.
* Knowledge is reproducible.
* Knowledge is independently verifiable.

⸻

26.18 Conformance

A Knowledge implementation conforms when it:

* supports the canonical lifecycle;
* preserves provenance;
* enforces publication workflows;
* supports deterministic retrieval;
* produces evidence;
* exposes observability;
* preserves architectural invariants.

⸻

Closing Thought

Knowledge provides the stable foundation upon which intelligent execution depends.

By treating knowledge as a governed organizational asset rather than a collection of documents, Forge enables consistent, explainable and reproducible intelligence across applications, teams and technologies.

⸻

Een architectonische verbetering: introduceer de “Repository” als eerste-klasse capability

Nu we de Knowledge-specificatie hebben uitgewerkt, zie ik één kans om Forge nog consistenter te maken.

Tot nu toe behandelen we repositories (Git, SharePoint, Confluence, CMDB, enzovoort) als implementatiedetails. Maar in de architectuur spelen ze een veel fundamentelere rol. Eigenlijk zou Forge expliciet onderscheid moeten maken tussen:

Repository
      │
      ▼
Knowledge
      │
      ▼
Memory
      │
      ▼
Context
      │
      ▼
Runtime

Daarmee ontstaat een belangrijk architectuurprincipe:

Repositories store information; Knowledge governs meaning.

Een Git-repository, SharePoint-site of CMDB is slechts een opslag- en beheersysteem. Pas wanneer de inhoud is gevalideerd, geclassificeerd, geversioneerd en onder governance is gebracht, wordt die informatie Knowledge binnen Forge.

Dat onderscheid voorkomt een veelvoorkomende fout in AI-platformen: een vector database of documentstore wordt onterecht als “knowledge base” gezien. In Forge is de repository slechts de drager. De Knowledge capability bepaalt welke informatie daadwerkelijk als organisatorische waarheid mag worden gebruikt. Dat is een subtiel maar fundamenteel onderscheid, en het versterkt de positionering van Forge als een Engineering Operating System in plaats van een verzameling AI-componenten.

Mooi. We hebben nu de twee fundamentele gegevensobjecten gedefinieerd:

* forge.spec.context.v1
* forge.spec.memory.v1

De volgende specificatie moet eigenlijk Knowledge zijn. Niet omdat een LLM kennis nodig heeft, maar omdat de hele Engineering Operating System rust op een gecontroleerde bron van waarheid.

Dit is ook waar Forge zich duidelijk onderscheidt van de meeste RAG-frameworks. Die behandelen documenten als data. Forge behandelt Knowledge als een governed enterprise asset.

⸻

Chapter 26

Knowledge Specification

Specification Identifier: forge.spec.knowledge.v1

Status: Draft Standard

⸻

Knowledge is not information.

Knowledge is information that has been validated, governed and accepted as organizational truth.

⸻

26.1 Purpose

This specification defines the canonical representation, lifecycle and governance of Knowledge within the Forge Engineering Operating System.

Knowledge SHALL represent organizational truth.

Knowledge SHALL NOT represent temporary observations, experiences or assumptions.

Knowledge exists to provide stable, governed information for intelligent execution.

⸻

26.2 Design Principles

Knowledge SHALL be:

* authoritative;
* versioned;
* attributable;
* governed;
* immutable once published;
* discoverable;
* explainable;
* independently verifiable.

Knowledge SHALL always originate from a controlled source.

⸻

26.3 Knowledge Domains

Forge distinguishes multiple Knowledge Domains.

Knowledge
├── Architecture
├── Policies
├── Standards
├── Procedures
├── Reference Data
├── Domain Models
├── Documentation
└── External References

Each domain SHALL define its own governance process.

⸻

26.4 Knowledge Object

Every Knowledge Object SHALL contain the following metadata.

knowledge:
id:
title:
owner:
domain:
classification:
version:
status:
created:
published:
reviewDate:
source:
checksum:
references:
evidence:

Knowledge identifiers SHALL remain stable throughout the lifecycle of the object.

⸻

26.5 Knowledge Content

A Knowledge Object SHALL contain:

* structured content;
* unstructured content;
* metadata;
* provenance;
* relationships.

The internal representation MAY vary.

The logical structure SHALL remain consistent.

⸻

26.6 Knowledge Lifecycle

stateDiagram-v2
[*] --> Draft
Draft --> Review
Review --> Approved
Approved --> Published
Published --> Superseded
Superseded --> Archived
Archived --> Retired
Retired --> [*]

No Knowledge Object SHALL bypass the approval process.

⸻

26.7 Source of Truth

Every Knowledge Object SHALL reference exactly one authoritative source.

Examples include:

* Git repository;
* Enterprise Architecture Repository;
* SharePoint document library;
* Policy management system;
* CMDB;
* Product catalog.

Replicated copies SHALL NOT become authoritative.

⸻

26.8 Knowledge Relationships

Knowledge Objects MAY reference other Knowledge Objects.

Relationship types include:

* supersedes;
* dependsOn;
* derivedFrom;
* references;
* complements;
* conflictsWith.

These relationships SHALL be explicit.

⸻

26.9 Versioning

Knowledge SHALL use semantic versioning.

Example:

Architecture Standard
v1.0.0
↓
v1.1.0
↓
v2.0.0

Breaking changes SHALL increment the major version.

Minor editorial corrections SHOULD increment the patch version.

⸻

26.10 Knowledge Validation

Knowledge SHALL be validated before publication.

Validation MAY include:

* schema validation;
* link validation;
* ownership verification;
* approval verification;
* classification verification;
* reference validation.

Knowledge failing validation SHALL NOT be published.

⸻

26.11 Knowledge Provenance

Every Knowledge Object SHALL preserve provenance.

Example:

provenance:
author:
reviewer:
approver:
repository:
commit:
timestamp:
evidence:

Knowledge without provenance SHALL be considered untrusted.

⸻

26.12 Retrieval

Knowledge retrieval SHALL be deterministic.

Inputs MAY include:

* objective;
* domain;
* classification;
* version;
* language;
* organizational scope.

Ranking algorithms MAY be implementation-specific.

Returned objects SHALL preserve their metadata.

⸻

26.13 Knowledge Publication

Publication SHALL create an immutable published version.

Subsequent modifications SHALL produce a new version.

Published knowledge SHALL remain reproducible.

⸻

26.14 Governance

Every Knowledge Domain SHALL define:

* owner;
* review cycle;
* approval process;
* publication workflow;
* archival policy;
* retention policy.

Knowledge SHALL never exist without ownership.

⸻

26.15 Security

Knowledge SHALL define:

* classification;
* visibility;
* tenant scope;
* access policy;
* export policy.

Retrieval SHALL respect Policy decisions.

Knowledge SHALL NOT enforce authorization independently.

⸻

26.16 Observability

Knowledge SHALL expose:

Metrics:

* retrieval frequency;
* publication frequency;
* review age;
* validation failures;
* stale documents;
* dependency graph health.

Evidence:

* publication reports;
* review reports;
* validation reports.

⸻

26.17 Architectural Invariants

Every Forge implementation SHALL preserve:

* Knowledge is authoritative.
* Knowledge is immutable after publication.
* Knowledge is versioned.
* Knowledge is attributable.
* Knowledge is governed.
* Knowledge is observable.
* Knowledge is reproducible.
* Knowledge is independently verifiable.

⸻

26.18 Conformance

A Knowledge implementation conforms when it:

* supports the canonical lifecycle;
* preserves provenance;
* enforces publication workflows;
* supports deterministic retrieval;
* produces evidence;
* exposes observability;
* preserves architectural invariants.

⸻

Closing Thought

Knowledge provides the stable foundation upon which intelligent execution depends.

By treating knowledge as a governed organizational asset rather than a collection of documents, Forge enables consistent, explainable and reproducible intelligence across applications, teams and technologies.

⸻

Een architectonische verbetering: introduceer de “Repository” als eerste-klasse capability

Nu we de Knowledge-specificatie hebben uitgewerkt, zie ik één kans om Forge nog consistenter te maken.

Tot nu toe behandelen we repositories (Git, SharePoint, Confluence, CMDB, enzovoort) als implementatiedetails. Maar in de architectuur spelen ze een veel fundamentelere rol. Eigenlijk zou Forge expliciet onderscheid moeten maken tussen:

Repository
      │
      ▼
Knowledge
      │
      ▼
Memory
      │
      ▼
Context
      │
      ▼
Runtime

Daarmee ontstaat een belangrijk architectuurprincipe:

Repositories store information; Knowledge governs meaning.

Een Git-repository, SharePoint-site of CMDB is slechts een opslag- en beheersysteem. Pas wanneer de inhoud is gevalideerd, geclassificeerd, geversioneerd en onder governance is gebracht, wordt die informatie Knowledge binnen Forge.

Dat onderscheid voorkomt een veelvoorkomende fout in AI-platformen: een vector database of documentstore wordt onterecht als “knowledge base” gezien. In Forge is de repository slechts de drager. De Knowledge capability bepaalt welke informatie daadwerkelijk als organisatorische waarheid mag worden gebruikt. Dat is een subtiel maar fundamenteel onderscheid, en het versterkt de positionering van Forge als een Engineering Operating System in plaats van een verzameling AI-componenten.

Perfect. We hebben nu de canonieke data-objecten (Knowledge, Memory, Context) en de Execution Kernel (Runtime) gedefinieerd.

De volgende specificatie moet Policy zijn.

En hier zou ik bewust afwijken van bestaande oplossingen zoals OPA, Cedar of IAM policies.

Waarom?

Omdat Forge Policy niet alleen authorization doet.

Policy bepaalt de governance van intelligentie.

Dat is veel breder.

⸻

Chapter 28

Policy Specification

Specification Identifier: forge.spec.policy.v1

Status: Draft Standard

⸻

Intelligence without policy is capability without accountability.

Policy transforms possibility into permission.

⸻

28.1 Purpose

This specification defines the canonical representation, evaluation and enforcement of Policy within the Forge Engineering Operating System.

Policy SHALL determine the operational boundaries within which intelligent execution is permitted.

Policy SHALL express organizational intent.

Policy SHALL NOT perform execution.

⸻

28.2 Design Principles

Policy SHALL be:

* declarative;
* deterministic;
* explainable;
* versioned;
* attributable;
* observable;
* composable;
* independently testable.

Policy SHALL describe what is permitted, never how it is implemented.

⸻

28.3 Policy Domains

Forge distinguishes multiple policy domains.

Policy
├── Identity
├── Authorization
├── Information Classification
├── Model Selection
├── Tool Usage
├── Runtime
├── Economics
├── Compliance
├── Privacy
└── Safety

Each domain SHALL be independently governable.

⸻

28.4 Policy Object

Every Policy SHALL contain the following metadata.

policy:
id:
name:
version:
owner:
status:
classification:
priority:
effectiveFrom:
expires:
scope:
evidence:

Policy identifiers SHALL remain immutable.

⸻

28.5 Policy Lifecycle

stateDiagram-v2
[*] --> Draft
Draft --> Review
Review --> Approved
Approved --> Active
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Policy SHALL NOT become Active without formal approval.

⸻

28.6 Policy Evaluation

Policy evaluation SHALL occur before execution.

Inputs MAY include:

* Identity
* Context Package
* Objective
* Classification
* Tenant
* Time
* Location
* Organizational Unit

Outputs SHALL contain explicit decisions.

Policy SHALL NOT return ambiguity.

⸻

28.7 Policy Decision

A Policy Decision SHALL contain:

decision:
status:
reason:
constraints:
obligations:
expires:
evidence:

Status SHALL be one of:

* Permit
* Deny
* PermitWithObligations
* Escalate

No other decision states are permitted.

⸻

28.8 Obligations

Policies MAY impose obligations in addition to authorization.

Examples include:

* require human approval;
* redact classified information;
* use only sovereign models;
* limit token consumption;
* store evidence for seven years;
* disable external tools;
* enforce watermarking.

Obligations SHALL be executable by downstream capabilities.

⸻

28.9 Conflict Resolution

Multiple policies MAY apply simultaneously.

Forge SHALL evaluate policy conflicts deterministically.

The default precedence SHALL be:

1. Legal requirements
2. Organizational policy
3. Information classification
4. Security policy
5. Business policy
6. User preference

Lower-priority policies SHALL NOT override higher-priority policies.

⸻

28.10 Policy Composition

Policies SHALL be composable.

Example:

Identity Policy
        │
        ▼
Classification Policy
        │
        ▼
Model Policy
        │
        ▼
Tool Policy
        │
        ▼
Economic Policy
        │
        ▼
Execution Decision

Composition SHALL produce a single, deterministic Policy Decision.

⸻

28.11 Policy Scope

Policies MAY apply at different scopes.

Examples include:

* Global
* Organization
* Business Unit
* Project
* Application
* Tenant
* Execution

The most specific applicable policy SHALL be evaluated only after broader mandatory policies have been applied.

⸻

28.12 Policy Enforcement

Policy SHALL determine constraints.

The Runtime SHALL enforce them.

Examples include:

Policy Decision	Runtime Action
Deny	Reject execution
Approved Model List	Restrict model routing
Token Budget	Enforce runtime limits
Human Approval	Suspend execution
Region Constraint	Select compliant infrastructure

Policy decides.

Runtime executes.

⸻

28.13 Policy Evidence

Every evaluation SHALL produce evidence.

Minimum evidence SHALL include:

* evaluated policies;
* evaluation timestamp;
* policy versions;
* decision outcome;
* obligations;
* evaluator version.

Policy decisions SHALL be reproducible.

⸻

28.14 Observability

Policy SHALL expose:

Metrics:

* evaluation count;
* permit rate;
* deny rate;
* escalation rate;
* evaluation latency;
* policy conflicts.

Logs:

* evaluated policies;
* decision summaries;
* conflict resolution.

Traces:

* policy evaluation lineage.

⸻

28.15 Security

Policy SHALL be protected against unauthorized modification.

Changes SHALL require:

* authenticated identity;
* authorized approver;
* version increment;
* evidence generation;
* governance approval.

Policy tampering SHALL be detectable.

⸻

28.16 Architectural Invariants

Every Forge Policy implementation SHALL preserve:

* Policy is declarative.
* Policy evaluation precedes execution.
* Every decision is explainable.
* Every decision produces evidence.
* Policy is versioned.
* Policy is independently testable.
* Policy conflicts are deterministic.
* Runtime SHALL NOT reinterpret Policy.

⸻

28.17 Conformance

A Policy implementation conforms when it:

* implements the canonical lifecycle;
* produces deterministic decisions;
* generates evidence;
* exposes observability;
* supports composition;
* preserves architectural invariants.

⸻

Closing Thought

Policy is often reduced to access control.

Forge treats Policy as the governance engine of intelligent execution.

It determines not only whether execution is permitted, but under which conditions, with which constraints, using which resources, and under what obligations.

By separating policy decisions from runtime behaviour, Forge enables organizations to evolve governance independently of implementation, preserving trust without sacrificing agility.

⸻

Een fundamentele uitbreiding: de Forge Decision Model (FDM)

Tijdens het uitwerken van deze specificatie valt een patroon op. Vrijwel alle capabilities nemen uiteindelijk beslissingen:

* Context beslist wat relevant is.
* Runtime beslist wanneer iets wordt uitgevoerd.
* Model Routing beslist welk model wordt gebruikt.
* Tool Routing beslist welke tool wordt aangeroepen.
* Economics beslist tegen welke kosten.
* Policy beslist of en onder welke voorwaarden iets mag.

Dat suggereert een dieper architectuurconcept: een generiek Forge Decision Model (FDM).

In plaats van elke capability zijn eigen beslislogica te laten definiëren, zou Forge één canoniek besliscontract kunnen introduceren:

decision:
  id:
  type:
  subject:
  inputs:
  evaluation:
  outcome:
  confidence:
  obligations:
  evidence:
  timestamp:

Daarmee krijgen alle beslissingen binnen Forge dezelfde structuur, ongeacht of het een Policy Decision, Model Selection, Context Composition of Runtime Planning betreft.

Dat is vergelijkbaar met hoe OpenTelemetry een uniform model biedt voor logs, metrics en traces. Een Forge Decision Model zou hetzelfde doen voor besluitvorming: één gemeenschappelijke taal voor alle intelligente beslissingen binnen het Engineering Operating System. Ik denk dat dit een onderscheidend element kan worden dat de architectuur nog consistenter en beter uitbreidbaar maakt.


Mooi. We hebben nu de kern van de Engineering Operating System gedefinieerd:

* Knowledge
* Memory
* Context
* Runtime
* Policy

De volgende capability ligt eigenlijk voor de hand, maar ik zou hem heel anders modelleren dan de meeste AI-platformen.

De meeste platforms hebben een Model Provider of een LLM Gateway.

Forge zou een Model Capability moeten hebben.

Waarom?

Omdat een model een vervangbare compute resource is, geen architectuurcomponent.

Net zoals een CPU vervangbaar is in een besturingssysteem.

⸻

Chapter 29

Model Specification

Specification Identifier: forge.spec.model.v1

Status: Draft Standard

⸻

Models provide computational intelligence.

They do not provide architecture.

Architecture determines when, why and how models are used.

⸻

29.1 Purpose

This specification defines the canonical representation, selection and governance of Models within the Forge Engineering Operating System.

The Model capability SHALL provide computational reasoning services.

The Model capability SHALL NOT determine when reasoning is required.

That responsibility belongs to the Runtime.

⸻

29.2 Design Principles

The Model capability SHALL be:

* provider-independent;
* versioned;
* observable;
* governed;
* replaceable;
* measurable;
* policy-constrained;
* economically optimized.

Models SHALL be treated as interchangeable execution resources.

⸻

29.3 Responsibilities

The Model capability SHALL:

* expose reasoning capabilities;
* execute prompts;
* generate structured responses;
* report resource usage;
* expose model metadata;
* support capability discovery.

The Model capability SHALL NOT:

* retrieve knowledge;
* compose context;
* enforce policy;
* invoke enterprise systems;
* manage execution state.

⸻

29.4 Model Object

Every registered model SHALL contain:

model:
id:
name:
provider:
family:
version:
capabilities:
region:
classification:
status:
costProfile:
limits:
evidence:

Model identifiers SHALL remain stable.

⸻

29.5 Model Capabilities

A model SHALL explicitly declare its capabilities.

Examples include:

capabilities:
text-generation
reasoning
vision
speech
embeddings
structured-output
tool-calling
streaming

Capabilities SHALL be machine-readable.

⸻

29.6 Model Lifecycle

stateDiagram-v2
[*] --> Registered
Registered --> Certified
Certified --> Approved
Approved --> Active
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Only Active models SHALL be eligible for production execution.

⸻

29.7 Model Selection

Model selection SHALL be deterministic.

Inputs MAY include:

* Policy Decision;
* Context Package;
* Objective;
* Classification;
* Cost constraints;
* Latency requirements;
* Capability requirements.

The Runtime SHALL request a model.

The Model capability SHALL resolve the most appropriate implementation.

⸻

29.8 Model Routing

Forge distinguishes logical models from physical providers.

Runtime
↓
Model Capability
↓
Logical Model
↓
Provider Selection
↓
Physical Endpoint

Applications SHALL NOT depend on provider APIs.

⸻

29.9 Provider Independence

Supported providers MAY include:

* OpenAI
* Azure OpenAI
* Anthropic
* Google
* Mistral
* Cohere
* Cisco
* Ollama
* vLLM
* Local inference clusters

Provider selection SHALL remain transparent to consumers.

⸻

29.10 Structured Outputs

Models SHOULD support structured outputs.

Example:

response:
schema:
confidence:
citations:
reasoning:
metadata:

Consumers SHOULD avoid parsing free-form text whenever structured contracts are available.

⸻

29.11 Model Metadata

Every invocation SHALL record:

* model identifier;
* provider;
* version;
* latency;
* prompt tokens;
* completion tokens;
* total cost;
* execution region.

This metadata SHALL become part of the execution evidence.

⸻

29.12 Model Evaluation

Models SHALL be continuously evaluated.

Minimum evaluation dimensions include:

* accuracy;
* latency;
* determinism;
* hallucination rate;
* structured output compliance;
* safety performance;
* cost efficiency.

Evaluation SHALL be repeatable.

⸻

29.13 Fallback

The Model capability MAY define fallback chains.

Example:

Forge GPT-Enterprise
↓
Azure GPT
↓
Local Llama
↓
Rule Engine

Fallback SHALL preserve Policy constraints.

⸻

29.14 Security

Models SHALL declare:

* supported regions;
* data retention behaviour;
* training policy;
* tenant isolation;
* certification status.

Models SHALL NOT receive information exceeding their authorized classification.

⸻

29.15 Observability

The Model capability SHALL expose:

Metrics:

* invocation count;
* latency;
* token usage;
* failure rate;
* cost;
* confidence distribution;
* provider availability.

Evidence:

* invocation reports;
* routing reports;
* evaluation reports.

⸻

29.16 Architectural Invariants

Every Forge Model implementation SHALL preserve:

* Models never own execution.
* Models never own context.
* Models never own knowledge.
* Models are replaceable.
* Models are observable.
* Models are governed.
* Models are policy-constrained.
* Models produce measurable outcomes.

⸻

29.17 Conformance

A Model implementation conforms when it:

* exposes canonical metadata;
* supports lifecycle management;
* produces invocation evidence;
* supports provider abstraction;
* preserves architectural invariants.

⸻

Closing Thought

Models are among the fastest-changing technologies in artificial intelligence.

For that reason, they should be among the least stable architectural elements.

Forge deliberately isolates models behind a stable capability, allowing organizations to adopt new providers, optimize costs or meet regulatory requirements without changing application architecture.

⸻

Een belangrijke architectuurcorrectie

Tijdens het schrijven van deze specificatie zie ik een verfijning die Forge nog sterker maakt.

In de huidige Capability Model spreken we over de Model capability, maar eigenlijk omvat die twee verschillende verantwoordelijkheden:

Model Capability
│
├── Model Registry
│      (catalogus en metadata)
│
├── Model Selection
│      (capability matching)
│
├── Model Routing
│      (provider keuze)
│
├── Model Execution
│      (inference)
│
└── Model Evaluation
       (kwaliteit, kosten en prestaties)

Dit is vergelijkbaar met hoe Kubernetes niet alleen een scheduler heeft, maar ook een API Server, Controller Manager en kubelet. Voor de buitenwereld blijft het één capability met één contract, terwijl de interne architectuur modulair blijft.

Dat biedt een belangrijk voordeel: nieuwe modelproviders, evaluatiemechanismen of routingstrategieën kunnen worden toegevoegd zonder de externe contracten te wijzigen. De Model capability blijft daarmee stabiel, terwijl de AI-markt erachter zich in hoog tempo kan blijven ontwikkelen.

Perfect. We hebben nu de “compute” gedefinieerd. Maar een model kan niets doen zonder met de buitenwereld te communiceren.

Dat brengt ons bij de volgende capability.

Niet “Functions”.

Niet “Plugins”.

Niet “MCP”.

Niet “API Gateway”.

Maar Tools.

En hier zou ik een fundamenteel architectuurprincipe introduceren:

Models think. Tools act.

Dat is misschien wel één van de belangrijkste scheidingen in Forge.

⸻

Chapter 30

Tool Specification

Specification Identifier: forge.spec.tool.v1

Status: Draft Standard

⸻

Models reason.

Tools interact.

The Runtime coordinates both.

⸻

30.1 Purpose

This specification defines the canonical representation, lifecycle and execution semantics of Tools within the Forge Engineering Operating System.

A Tool SHALL expose deterministic capabilities to the Runtime.

A Tool SHALL provide controlled interaction with external systems.

A Tool SHALL NOT perform autonomous orchestration.

⸻

30.2 Design Principles

The Tool capability SHALL be:

* deterministic;
* discoverable;
* versioned;
* composable;
* observable;
* governed;
* policy-aware;
* provider-independent.

Tools SHALL expose contracts rather than implementations.

⸻

30.3 Responsibilities

The Tool capability SHALL:

* expose executable operations;
* validate contracts;
* advertise capabilities;
* invoke external systems;
* return structured responses;
* generate evidence.

The Tool capability SHALL NOT:

* evaluate policies;
* retrieve knowledge;
* compose context;
* plan execution;
* select models.

⸻

30.4 Tool Object

Every Tool SHALL define:

tool:
id:
name:
version:
owner:
category:
description:
status:
classification:
endpoint:
authentication:
contract:
evidence:

Tool identifiers SHALL remain immutable.

⸻

30.5 Tool Categories

Forge defines canonical Tool categories.

Tool
├── Enterprise Systems
├── Data Sources
├── Communication
├── Infrastructure
├── AI Services
├── Security
├── Automation
├── Collaboration
└── Custom

Categories MAY be extended.

The category taxonomy SHALL remain hierarchical.

⸻

30.6 Tool Contract

Every Tool SHALL publish a machine-readable contract.

Minimum contract elements:

contract:
operation:
inputs:
outputs:
errors:
timeouts:
idempotent:
authentication:
authorization:

Contracts SHALL be versioned.

⸻

30.7 Discovery

The Runtime SHALL discover Tools through the Tool capability.

Discovery SHALL expose:

* identifier;
* capabilities;
* version;
* owner;
* health;
* supported operations;
* policy constraints.

Applications SHALL NOT perform discovery directly.

⸻

30.8 Invocation

Every invocation SHALL include:

invocation:
executionId:
toolId:
operation:
contextId:
identity:
policyDecision:
parameters:

All invocations SHALL be attributable.

⸻

30.9 Response

Every Tool SHALL return a structured response.

response:
status:
result:
errors:
duration:
evidence:
metadata:

Free-form responses SHOULD be avoided.

⸻

30.10 Lifecycle

stateDiagram-v2
[*] --> Registered
Registered --> Certified
Certified --> Approved
Approved --> Active
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Inactive Tools SHALL NOT be available for execution.

⸻

30.11 Execution Rules

The Runtime SHALL:

* validate the contract;
* validate Policy;
* validate parameters;
* enforce timeout;
* capture evidence.

The Runtime SHALL NOT bypass Tool contracts.

⸻

30.12 Tool Composition

Multiple Tools MAY participate in one execution.

Example:

Planner
↓
CRM Tool
↓
Knowledge Tool
↓
Email Tool
↓
Archive Tool

Composition SHALL be orchestrated by the Runtime.

Tools SHALL remain unaware of one another.

⸻

30.13 Tool Isolation

Each Tool SHALL execute independently.

A Tool SHALL NOT:

* invoke another Tool directly;
* invoke Models autonomously;
* modify Runtime state;
* bypass Policy.

Cross-tool orchestration belongs exclusively to the Runtime.

⸻

30.14 External Protocols

Forge intentionally separates Tool contracts from transport protocols.

Supported implementations MAY include:

* MCP
* REST
* GraphQL
* gRPC
* AMQP
* Kafka
* Webhooks

Protocols are implementation details.

Contracts are architectural.

⸻

30.15 Authentication

Authentication SHALL occur before Tool execution.

Supported mechanisms MAY include:

* OAuth 2.0
* OpenID Connect
* Mutual TLS
* API Keys
* Managed Identity
* Service Accounts

Credential management SHALL remain external to Tool logic.

⸻

30.16 Authorization

Authorization SHALL be derived from:

* Identity;
* Policy Decision;
* Tool Classification;
* Organizational Scope.

Tools SHALL trust the authorization context provided by the Runtime.

⸻

30.17 Evidence

Every invocation SHALL produce evidence.

Minimum evidence SHALL include:

* tool identifier;
* operation;
* execution identifier;
* latency;
* status;
* policy reference;
* response classification.

Evidence SHALL exist regardless of success or failure.

⸻

30.18 Observability

The Tool capability SHALL expose:

Metrics:

* invocation count;
* latency;
* availability;
* error rate;
* timeout rate;
* throughput.

Logs:

* invocation summaries;
* failures;
* retries.

Traces:

* end-to-end execution lineage.

⸻

30.19 Architectural Invariants

Every Forge Tool implementation SHALL preserve:

* Tools are deterministic.
* Tools expose contracts.
* Tools are independently deployable.
* Tools never orchestrate execution.
* Tools never evaluate Policy.
* Tools always produce evidence.
* Tools are observable.
* Tools are replaceable.

⸻

30.20 Conformance

A Tool implementation conforms when it:

* publishes a canonical contract;
* supports lifecycle management;
* produces evidence;
* exposes observability;
* preserves architectural invariants.

⸻

Closing Thought

Tools connect the Engineering Operating System to the outside world.

By isolating external interactions behind governed contracts, Forge ensures that enterprise systems, cloud services and emerging protocols can evolve independently of application logic and model capabilities.

The result is an architecture in which intelligence remains focused on reasoning, while interaction remains controlled, observable and replaceable.

⸻

Een fundamentele verbetering: de Tool Capability als “Device Driver Model”

Tijdens het schrijven van deze specificatie viel een opvallende analogie op met traditionele besturingssystemen.

In een besturingssysteem spreekt een proces nooit rechtstreeks met een fysieke schijf, netwerkkaart of printer. Alle interacties verlopen via een device driver die een gestandaardiseerde interface biedt.

Hetzelfde principe kan Forge hanteren:

Business Capability
        │
        ▼
Runtime
        │
        ▼
Tool Contract
        │
        ▼
Tool Adapter
        │
        ▼
Protocol Adapter
        │
        ▼
Enterprise System

Hiermee ontstaat een duidelijke scheiding tussen wat een tool doet (het contract) en hoe dat technisch wordt uitgevoerd (de adapter).

Dit leidt tot een waardevol architectuurprincipe:

A Tool is not an API. A Tool is an architectural capability exposed through one or more protocol adapters.

Dat betekent bijvoorbeeld dat dezelfde Tool gelijktijdig beschikbaar kan zijn via MCP, REST en gRPC, zonder dat de Runtime of de applicaties hoeven te veranderen. Net zoals een besturingssysteem dezelfde printerdriver kan gebruiken ongeacht of de printer via USB, IPP of een ander transport is aangesloten. Dat onderscheid tussen contract en protocol maakt de Tool-capability robuuster, toekomstbestendiger en beter passend binnen de visie van Forge als een Engineering Operating System.

Perfect. We hebben nu de volledige execution stack beschreven:

* Knowledge
* Memory
* Context
* Policy
* Runtime
* Models
* Tools

Vanaf hier verschuift Forge van execution naar trust.

En de eerste capability daarin is bewust niet Security.

Dat is Identity.

Waarom?

Omdat vrijwel alles binnen Forge uiteindelijk afhangt van één fundamentele vraag:

Namens wie wordt deze intelligentie uitgevoerd?

Net zoals een besturingssysteem eerst een procesidentiteit heeft voordat rechten kunnen worden toegepast, moet Forge eerst een Execution Identity vaststellen voordat Policy, Runtime of Tools kunnen handelen.

⸻

Chapter 31

Identity Specification

Specification Identifier: forge.spec.identity.v1

Status: Draft Standard

⸻

Every execution acts on behalf of someone.

Identity establishes who.

Policy determines what is permitted.

⸻

31.1 Purpose

This specification defines the canonical representation, lifecycle and semantics of Identity within the Forge Engineering Operating System.

Identity SHALL establish the authenticated and attributable subject on whose behalf execution occurs.

Identity SHALL precede Policy evaluation.

Identity SHALL remain immutable throughout an execution.

⸻

31.2 Design Principles

Identity SHALL be:

* verifiable;
* attributable;
* immutable during execution;
* privacy-aware;
* federated;
* observable;
* auditable;
* independent of authentication technology.

Identity represents a subject.

It does not authenticate the subject.

⸻

31.3 Responsibilities

The Identity capability SHALL:

* resolve execution identity;
* normalize identity attributes;
* establish organizational context;
* represent delegated authority;
* expose identity claims;
* generate identity evidence.

The Identity capability SHALL NOT:

* evaluate Policy;
* authenticate users directly;
* authorize execution;
* manage credentials;
* execute business logic.

⸻

31.4 Identity Object

Every execution SHALL reference exactly one Identity Object.

identity:
id:
subject:
displayName:
tenant:
organization:
department:
roles:
groups:
clearance:
authenticationContext:
delegation:
attributes:
evidence:

The identifier SHALL remain stable throughout the execution lifecycle.

⸻

31.5 Identity Types

Forge recognizes multiple identity types.

Identity
├── Human
├── Service Account
├── Application
├── Agent
├── Device
├── Workflow
└── External Partner

Each type MAY define additional attributes.

The canonical structure SHALL remain consistent.

⸻

31.6 Authentication Boundary

Authentication occurs outside the Engineering Operating System.

Examples include:

* OpenID Connect;
* SAML 2.0;
* Kerberos;
* WebAuthn;
* X.509 certificates;
* Managed Identity.

Forge consumes authenticated identity.

Forge does not authenticate users itself.

⸻

31.7 Identity Resolution

Identity resolution SHALL normalize information received from external providers.

Inputs MAY include:

* identity provider claims;
* enterprise directories;
* HR systems;
* organizational repositories;
* federation services.

The resulting Identity Object SHALL be implementation-independent.

⸻

31.8 Delegation

Executions MAY occur under delegated authority.

Examples include:

* executive assistants;
* workflow automation;
* scheduled jobs;
* AI agents;
* service accounts.

Delegation SHALL explicitly record:

delegation:
delegator:
delegate:
reason:
scope:
expires:

Delegated authority SHALL be auditable.

⸻

31.9 Identity Attributes

Identity attributes MAY include:

* organization;
* project;
* location;
* employment status;
* security clearance;
* business function;
* certifications.

Attributes SHALL remain descriptive.

They SHALL NOT express permissions.

Permissions belong to Policy.

⸻

31.10 Identity Federation

Forge SHALL support federated identity.

Examples include:

Entra ID
↓
Identity Capability
↓
Canonical Identity
↓
Runtime

Multiple providers SHALL produce the same canonical Identity representation.

⸻

31.11 Privacy

Identity SHALL minimize exposed personal information.

Only attributes required for execution SHALL be included.

Personally identifiable information SHALL be governed by organizational policy and applicable legislation.

Identity minimization is a mandatory architectural principle.

⸻

31.12 Identity Lifecycle

stateDiagram-v2
[*] --> Resolved
Resolved --> Validated
Validated --> Active
Active --> Revoked
Active --> Expired
Revoked --> Archived
Expired --> Archived
Archived --> [*]

Lifecycle events SHALL generate evidence.

⸻

31.13 Evidence

Identity SHALL generate evidence for:

* identity resolution;
* federation;
* delegation;
* attribute normalization;
* lifecycle transitions.

Identity evidence SHALL be retained according to governance requirements.

⸻

31.14 Observability

Identity SHALL expose:

Metrics:

* authenticated executions;
* federation latency;
* delegation count;
* identity resolution failures;
* provider availability.

Logs:

* resolution events;
* federation events;
* lifecycle events.

Traces:

* identity lineage across execution.

⸻

31.15 Security Considerations

Identity SHALL:

* verify issuer trust;
* validate token integrity;
* detect replay attempts where applicable;
* protect identity attributes in transit and at rest;
* support tenant isolation.

Credential storage SHALL remain outside the Identity capability.

⸻

31.16 Architectural Invariants

Every Forge Identity implementation SHALL preserve:

* Exactly one canonical Identity per execution.
* Identity precedes Policy evaluation.
* Identity is immutable during execution.
* Identity is attributable.
* Identity supports delegation.
* Identity is observable.
* Identity is privacy-aware.
* Identity produces evidence.

⸻

31.17 Conformance

An Identity implementation conforms when it:

* produces a canonical Identity Object;
* supports federation;
* preserves lifecycle semantics;
* generates evidence;
* exposes observability;
* preserves architectural invariants.

⸻

Closing Thought

Identity is the foundation of accountable intelligence.

Without a canonical representation of who an execution represents, governance becomes ambiguous, policy becomes unenforceable and evidence loses attribution.

By separating identity from authentication and authorization, Forge establishes a stable identity layer that remains independent of enterprise identity providers while preserving trust, portability and long-term architectural consistency.

⸻

Een fundamentele architectuurverfijning: Subject vs. Identity

Tijdens het uitwerken van deze specificatie zie ik nog een onderscheid dat Forge architectonisch sterker maakt.

In veel systemen worden Subject, Principal, Identity en User door elkaar gebruikt. Dat leidt tot verwarring, zeker wanneer AI-agents, workflows en service accounts namens anderen handelen.

Ik zou daarom een expliciet model introduceren:

Subject
   │
   ▼
Identity
   │
   ▼
Policy
   │
   ▼
Execution

Waarbij de verantwoordelijkheden helder zijn:

* Subject: de entiteit die handelt (mens, agent, applicatie, apparaat of workflow).
* Identity: de genormaliseerde representatie van dat subject binnen Forge.
* Policy: de evaluatie van wat die identiteit onder de huidige omstandigheden mag doen.
* Execution: de uitvoering van de toegestane handelingen.

Deze scheiding voorkomt dat identiteit wordt verward met autorisatie of authenticatie en sluit goed aan bij bestaande securitymodellen zoals NIST SP 800-63 en Zero Trust-principes. Tegelijk blijft Forge technologieneutraal: of een subject nu afkomstig is uit Entra ID, PingFederate, Keycloak of een toekomstige identiteitsprovider, de rest van het Engineering Operating System ziet altijd dezelfde canonieke Identity-representatie. Dat is precies het soort stabiele architectuurlaag dat een Engineering Operating System nodig heeft.

Mooi. We zijn nu de Trust Stack binnen Forge aan het opbouwen.

De logische volgorde is:

* Identity → Who?
* Security → How is it protected?
* Policy → What is permitted?
* Evidence → What can be proven?
* Observability → What can be understood?

Identity vertelt namens wie een uitvoering plaatsvindt.

Security beschermt vervolgens het gehele Engineering Operating System.

En dat is een belangrijk onderscheid.

Security is geen verzameling controls.

Security is een architectuurcapability.

⸻

Chapter 32

Security Specification

Specification Identifier: forge.spec.security.v1

Status: Draft Standard

⸻

Security is not the absence of attacks.

Security is the continuous preservation of trust despite the presence of threats.

⸻

32.1 Purpose

This specification defines the canonical security model for the Forge Engineering Operating System.

The Security capability SHALL preserve the confidentiality, integrity and availability of all Engineering Operating System capabilities.

Security SHALL operate continuously.

Security SHALL NOT become an implementation concern delegated solely to infrastructure.

⸻

32.2 Design Principles

Security SHALL be:

* identity-driven;
* policy-controlled;
* evidence-producing;
* observable;
* defense-in-depth;
* zero-trust;
* measurable;
* continuously verified.

Security SHALL be integrated into execution rather than added afterwards.

⸻

32.3 Responsibilities

The Security capability SHALL:

* protect execution;
* validate trust boundaries;
* preserve integrity;
* detect anomalies;
* protect data;
* manage trust relationships;
* generate security evidence.

The Security capability SHALL NOT:

* authenticate identities;
* evaluate policy;
* execute runtime workflows;
* classify business information.

⸻

32.4 Security Domains

Forge recognizes multiple security domains.

Security
├── Identity Protection
├── Data Protection
├── Execution Protection
├── Infrastructure Protection
├── Communication Security
├── Supply Chain Security
├── Monitoring
├── Incident Response
└── Recovery

Each domain SHALL define its own controls.

⸻

32.5 Security Object

Every Security Configuration SHALL contain:

security:
id:
version:
owner:
classification:
controls:
trustLevel:
boundaries:
cryptography:
evidence:

Security configurations SHALL be versioned.

⸻

32.6 Trust Boundaries

Forge defines explicit trust boundaries.

Application
│
Trust Boundary
│
Engineering Operating System
│
Trust Boundary
│
External Systems

Crossing a trust boundary SHALL trigger validation.

⸻

32.7 Data Classification

All information SHALL possess a classification.

Examples include:

* Public
* Internal
* Confidential
* Restricted
* Sovereign

Classification SHALL travel with the data throughout execution.

⸻

32.8 Cryptographic Protection

Sensitive information SHALL be protected using approved cryptographic mechanisms.

Security SHALL define:

* encryption at rest;
* encryption in transit;
* signing;
* hashing;
* key management.

Cryptographic algorithms MAY evolve.

The architectural requirement SHALL remain.

⸻

32.9 Execution Protection

Every execution SHALL protect:

* Context Package;
* Working Memory;
* Evidence;
* Model interactions;
* Tool invocations.

Protection SHALL persist throughout the execution lifecycle.

⸻

32.10 Communication Security

All communication SHALL be:

* authenticated;
* integrity protected;
* confidential where required;
* observable.

Protocols MAY include:

* TLS
* mTLS
* HTTPS
* gRPC
* Message signing

The protocol is implementation-specific.

The security properties are mandatory.

⸻

32.11 Supply Chain Security

Every executable component SHALL be attributable.

Examples include:

* Runtime
* Models
* Tools
* Plugins
* Containers
* Dependencies

Each SHALL possess:

* provenance;
* version;
* integrity verification;
* ownership.

Forge SHALL support Software Bill of Materials (SBOM) where applicable.

⸻

32.12 Threat Detection

Security SHALL continuously evaluate:

* anomalous executions;
* privilege escalation;
* unexpected tool usage;
* abnormal model behaviour;
* suspicious identity activity;
* policy bypass attempts.

Detection SHALL generate evidence.

⸻

32.13 Incident Lifecycle

stateDiagram-v2
[*] --> Detected
Detected --> Investigating
Investigating --> Contained
Contained --> Resolved
Resolved --> Reviewed
Reviewed --> Closed
Closed --> [*]

Every incident SHALL produce evidence.

⸻

32.14 Security Evidence

Security SHALL produce evidence for:

* trust validation;
* cryptographic verification;
* policy enforcement;
* boundary crossings;
* incident handling;
* integrity verification.

Security evidence SHALL be independently verifiable.

⸻

32.15 Observability

Security SHALL expose:

Metrics:

* failed validations;
* trust violations;
* encryption coverage;
* incident count;
* anomalous executions;
* policy bypass attempts.

Logs:

* trust events;
* boundary crossings;
* incident lifecycle;
* integrity failures.

Traces:

* security lineage throughout execution.

⸻

32.16 Zero Trust

Forge adopts Zero Trust as an architectural principle.

Every request SHALL be:

* explicitly verified;
* continuously evaluated;
* minimally privileged;
* fully attributable.

Trust SHALL never be inferred solely from network location.

⸻

32.17 Security Lifecycle

stateDiagram-v2
[*] --> Defined
Defined --> Approved
Approved --> Active
Active --> Monitored
Monitored --> Updated
Updated --> Active
Active --> Retired
Retired --> [*]

Security SHALL evolve continuously.

⸻

32.18 Architectural Invariants

Every Forge Security implementation SHALL preserve:

* Identity precedes trust.
* Every trust boundary is explicit.
* Every execution is protected.
* Every security decision produces evidence.
* Security is continuously observable.
* Security is policy-driven.
* Security is independently verifiable.
* Zero Trust applies throughout execution.

⸻

32.19 Conformance

A Security implementation conforms when it:

* protects all trust boundaries;
* implements continuous verification;
* generates evidence;
* exposes observability;
* preserves architectural invariants;
* supports Zero Trust principles.

⸻

Closing Thought

Security within Forge is not a perimeter surrounding the Engineering Operating System.

It is an intrinsic architectural capability woven into every execution, every interaction and every decision.

By integrating security with identity, policy, runtime and evidence, Forge ensures that trust is continuously established, continuously verified and continuously demonstrable.

⸻

Een fundamentele architectuuruitbreiding: de “Trust Plane”

Tijdens het schrijven van de laatste hoofdstukken zie ik een patroon dat nog sterker is dan de eerder beschreven “Trust Cycle”.

De huidige capabilities vormen namelijk een aparte architectuurlaag, vergelijkbaar met de Control Plane van Kubernetes.

                 TRUST PLANE
+-----------------------------------------------------+
| Identity | Security | Policy | Evidence | Governance |
+-----------------------------------------------------+
                     │
             ENGINEERING PLANE
+-----------------------------------------------------+
| Knowledge | Memory | Context | Runtime | Models | Tools |
+-----------------------------------------------------+
                     │
             EXECUTION PLANE
+-----------------------------------------------------+
| Applications | Workflows | Business Processes |
+-----------------------------------------------------+

Dit model maakt duidelijk dat Identity, Security, Policy, Evidence en Governance geen ondersteunende functies zijn, maar een control plane voor het Engineering Operating System. Zij sturen, begrenzen en verifiëren de execution capabilities zonder zelf de businesslogica uit te voeren.

Dat sluit opvallend goed aan bij de architectuur van moderne gedistribueerde systemen, waarin een duidelijke scheiding bestaat tussen een control plane (beheer, configuratie en governance) en een data plane (uitvoering). Forge kan dit concept verder uitbreiden door ook een Knowledge Plane en een Execution Plane te definiëren. Daarmee krijgt het platform een herkenbare, gelaagde referentiearchitectuur die zowel voor enterprise-architecten als voor implementatieteams intuïtief en consistent is. Ik denk dat dit uiteindelijk een van de onderscheidende kenmerken van de Forge Architecture Bible kan worden.

Mooi. Nu komen we bij een capability die ik persoonlijk zie als één van de meest onderscheidende onderdelen van Forge.

Niet omdat andere platforms geen logging hebben.

Maar omdat Forge Evidence fundamenteel anders behandelt.

In vrijwel alle systemen is audit een bijproduct.

In Forge is Evidence een eerste-klasse architectuurobject.

Dat is een enorm verschil.

⸻

Chapter 33

Evidence Specification

Specification Identifier: forge.spec.evidence.v1

Status: Draft Standard

⸻

Execution produces outcomes.

Evidence explains how those outcomes came into existence.

Without evidence, intelligence cannot be trusted.

⸻

33.1 Purpose

This specification defines the canonical representation, lifecycle and governance of Evidence within the Forge Engineering Operating System.

Evidence SHALL provide verifiable records of execution.

Evidence SHALL enable reproducibility, explainability and accountability.

Evidence SHALL exist independently of observability systems.

⸻

33.2 Design Principles

Evidence SHALL be:

* immutable;
* attributable;
* reproducible;
* independently verifiable;
* versioned;
* policy-governed;
* observable;
* portable.

Evidence SHALL describe facts.

It SHALL NOT contain interpretation.

⸻

33.3 Responsibilities

The Evidence capability SHALL:

* record execution events;
* preserve decision lineage;
* maintain execution chronology;
* support auditing;
* support replay;
* expose evidence packages.

The Evidence capability SHALL NOT:

* evaluate policy;
* interpret results;
* monitor execution health;
* orchestrate execution.

⸻

33.4 Evidence Object

Every Evidence Object SHALL contain:

evidence:
id:
executionId:
timestamp:
type:
producer:
classification:
integrity:
references:
content:
signature:

Evidence identifiers SHALL be globally unique.

⸻

33.5 Evidence Categories

Forge defines canonical evidence categories.

Evidence
├── Identity
├── Context
├── Policy
├── Runtime
├── Model
├── Tool
├── Security
├── Governance
├── Economics
└── System

Every Evidence Object SHALL belong to exactly one category.

⸻

33.6 Evidence Lifecycle

stateDiagram-v2
[*] --> Created
Created --> Verified
Verified --> Stored
Stored --> Archived
Archived --> Expired
Expired --> Destroyed
Destroyed --> [*]

Every lifecycle transition SHALL itself generate evidence.

⸻

33.7 Evidence Integrity

Evidence SHALL support integrity verification.

Minimum mechanisms include:

* cryptographic hashing;
* digital signatures;
* timestamp verification;
* chain validation.

Evidence SHALL be tamper-evident.

⸻

33.8 Evidence Lineage

Evidence SHALL reference preceding evidence where applicable.

Identity Evidence
↓
Context Evidence
↓
Policy Evidence
↓
Runtime Evidence
↓
Model Evidence
↓
Tool Evidence
↓
Completion Evidence

This forms an immutable execution lineage.

⸻

33.9 Evidence Package

Every execution SHALL produce an Evidence Package.

Minimum contents:

package:
execution:
identity:
context:
policy:
runtime:
models:
tools:
security:
completion:

The package SHALL be self-contained.

⸻

33.10 Evidence Replay

Evidence SHALL support execution replay.

Replay SHALL reconstruct:

* execution sequence;
* decisions;
* inputs;
* outputs;
* state transitions.

Replay SHALL NOT require the original runtime.

⸻

33.11 Evidence Retention

Retention SHALL be policy-driven.

Possible strategies:

* time-based;
* regulatory;
* legal hold;
* business retention;
* manual preservation.

Expired evidence SHALL follow controlled destruction procedures.

⸻

33.12 Evidence Classification

Evidence SHALL inherit the highest classification of:

* Context;
* Identity;
* Tool output;
* Model output;
* Security events.

Classification SHALL never decrease automatically.

⸻

33.13 Evidence Provenance

Every Evidence Object SHALL record:

provenance:
capability:
component:
version:
host:
runtime:
timestamp:

Evidence without provenance SHALL be considered invalid.

⸻

33.14 Observability

Evidence SHALL expose:

Metrics:

* evidence generation rate;
* verification failures;
* integrity failures;
* replay success rate;
* storage growth.

Logs:

* verification events;
* lifecycle transitions;
* integrity validation.

Traces:

* evidence lineage.

⸻

33.15 Security

Evidence SHALL support:

* encryption;
* integrity verification;
* controlled access;
* non-repudiation;
* tenant isolation.

Evidence SHALL be readable only by authorized identities.

⸻

33.16 Architectural Invariants

Every Forge Evidence implementation SHALL preserve:

* Evidence is immutable.
* Evidence is attributable.
* Evidence is independently verifiable.
* Every execution produces evidence.
* Every state transition produces evidence.
* Evidence is policy-governed.
* Evidence supports replay.
* Evidence is observable.

⸻

33.17 Conformance

An Evidence implementation conforms when it:

* produces immutable evidence;
* preserves lineage;
* supports replay;
* generates integrity metadata;
* exposes observability;
* preserves architectural invariants.

⸻

Closing Thought

Evidence is not a log.

It is not telemetry.

It is not an audit trail.

Evidence is the verifiable historical record of intelligent execution.

By elevating evidence to a first-class architectural capability, Forge enables organizations to explain, reproduce and defend intelligent decisions long after the original execution has completed.

⸻

Een fundamentele architectuurverbetering: Evidence als “Event Ledger”

Tijdens het schrijven van deze specificatie viel me op dat Forge Evidence eigenlijk nog sterker kan worden gepositioneerd.

De meeste platforms behandelen auditlogs als een afgeleid product van de runtime. Maar als je de architectuur van Forge consequent doortrekt, is Evidence feitelijk een append-only event ledger voor het Engineering Operating System.

Dat levert een aantal krachtige eigenschappen op:

* Evidence is append-only. Bestaande records worden nooit gewijzigd; correcties of aanvullingen worden als nieuwe gebeurtenissen toegevoegd.
* Iedere capability publiceert Evidence Events. Knowledge, Memory, Context, Policy, Runtime, Models, Tools, Identity en Security produceren allemaal gebeurtenissen in hetzelfde canonieke formaat.
* Observability consumeert Evidence. Dashboards, metrics en traces worden opgebouwd uit het Evidence Ledger in plaats van andersom.
* Governance en compliance steunen op hetzelfde fundament. Audits, forensisch onderzoek, replay en rapportages gebruiken dezelfde bron.

Conceptueel ontstaat dan dit model:

Knowledge ─┐
Memory ────┤
Context ───┤
Policy ────┤
Runtime ───┤
Models ────┤
Tools ─────┤
Identity ──┤
Security ──┤
            ▼
     Evidence Ledger
            ▼
  Observability • Audit • Replay • Analytics

Dat is vergelijkbaar met event sourcing, maar met een belangrijk verschil: het primaire doel is niet het reconstrueren van applicatiestatus, maar het reconstrueren en onderbouwen van intelligente besluitvorming. Dat zou Forge een zeer onderscheidend architectuurprincipe geven: intelligence is only trustworthy when every significant decision can be traced back to immutable evidence.

Perfect. Nu komen we bij een capability die bewust niet op Evidence volgt omdat hij data produceert, maar omdat hij Evidence interpreteert.

Dat is een cruciaal onderscheid.

Evidence antwoordt:

Wat is er gebeurd?

Observability antwoordt:

Waarom gebeurde het?

Dat is exact hetzelfde onderscheid dat moderne distributed systems maken tussen events en observability.

Forge kan dit principe veel verder doortrekken dan OpenTelemetry ooit heeft gedaan.

⸻

Chapter 34

Observability Specification

Specification Identifier: forge.spec.observability.v1

Status: Draft Standard

⸻

Evidence records execution.

Observability explains execution.

Understanding emerges from correlation, not collection.

⸻

34.1 Purpose

This specification defines the canonical representation and behaviour of Observability within the Forge Engineering Operating System.

The Observability capability SHALL transform Evidence into operational understanding.

Observability SHALL explain the behaviour of the Engineering Operating System.

Observability SHALL NOT generate execution evidence.

⸻

34.2 Design Principles

Observability SHALL be:

* evidence-driven;
* explainable;
* composable;
* correlated;
* measurable;
* provider-independent;
* continuously available;
* non-intrusive.

Observability SHALL consume evidence.

It SHALL never become the system of record.

⸻

34.3 Responsibilities

The Observability capability SHALL:

* correlate evidence;
* calculate metrics;
* construct traces;
* detect anomalies;
* expose health;
* support diagnostics;
* provide engineering insights.

The Observability capability SHALL NOT:

* execute workflows;
* generate evidence;
* evaluate policy;
* modify runtime behaviour.

⸻

34.4 Observability Object

Every Observability Record SHALL contain:

observability:
id:
executionId:
timestamp:
type:
source:
severity:
correlations:
metrics:
traces:
analysis:

Observability records SHALL reference Evidence.

They SHALL NOT replace it.

⸻

34.5 Signals

Forge extends the traditional observability model.

Observability
├── Metrics
├── Logs
├── Traces
├── Evidence
├── Context Lineage
├── Policy Decisions
├── Runtime States
├── Model Behaviour
├── Tool Behaviour
└── Economics

This model SHALL be considered canonical.

⸻

34.6 Correlation

Observability SHALL correlate information using:

* execution identifier;
* context identifier;
* evidence identifier;
* identity;
* tenant;
* timestamps;
* causal relationships.

Correlation SHALL be deterministic.

⸻

34.7 Metrics

Minimum Runtime metrics include:

* execution latency;
* queue depth;
* completion rate;
* retry count;
* cancellation rate.

Minimum Context metrics include:

* context size;
* knowledge coverage;
* memory utilization;
* context reuse.

Minimum Model metrics include:

* token consumption;
* latency;
* hallucination rate;
* structured output compliance.

Metrics SHALL be machine-readable.

⸻

34.8 Tracing

Every execution SHALL expose a canonical execution trace.

Execution
↓
Identity
↓
Context
↓
Policy
↓
Planning
↓
Runtime
↓
Model
↓
Tool
↓
Evidence
↓
Completion

The trace SHALL preserve causality.

⸻

34.9 Health

Every capability SHALL expose health.

Minimum health states:

* Healthy
* Degraded
* Unavailable
* Recovering
* Unknown

Health SHALL be observable without requiring execution.

⸻

34.10 Anomaly Detection

Observability SHALL support anomaly detection.

Examples include:

* unusual token usage;
* abnormal latency;
* repeated retries;
* unexpected model selection;
* policy violations;
* unusual tool invocation frequency;
* context explosion.

Detection algorithms MAY vary.

Detected anomalies SHALL generate evidence.

⸻

34.11 Diagnostics

Observability SHALL support diagnostics.

Diagnostic capabilities include:

* execution replay;
* dependency inspection;
* capability health;
* bottleneck identification;
* evidence navigation;
* lineage visualization.

Diagnostics SHALL remain read-only.

⸻

34.12 Dashboards

Dashboards are implementation artifacts.

The Observability capability SHALL expose structured data.

Visualization technologies MAY include:

* Grafana;
* Kibana;
* Azure Monitor;
* Splunk;
* Cisco ThousandEyes;
* custom dashboards.

Dashboards SHALL NOT define the architecture.

⸻

34.13 Engineering Metrics

Forge introduces Engineering Metrics.

Examples include:

Metric	Description
Context Quality	Percentage of relevant context used
Knowledge Freshness	Age of referenced knowledge
Memory Effectiveness	Contribution of memory to successful execution
Policy Complexity	Average evaluated rules
Runtime Efficiency	Successful executions per resource unit
Model Quality	Accuracy relative to objective
Tool Reliability	Successful invocations
Evidence Completeness	Percentage of required evidence produced

These metrics evaluate engineering quality rather than infrastructure.

⸻

34.14 Security

Observability SHALL:

* protect sensitive telemetry;
* preserve tenant isolation;
* redact classified information where required;
* authenticate telemetry consumers.

Observability SHALL NOT weaken security boundaries.

⸻

34.15 Lifecycle

stateDiagram-v2
[*] --> Collected
Collected --> Correlated
Correlated --> Analyzed
Analyzed --> Published
Published --> Archived
Archived --> Deleted
Deleted --> [*]

Observability artifacts SHALL remain reproducible.

⸻

34.16 Architectural Invariants

Every Forge Observability implementation SHALL preserve:

* Evidence is the primary source.
* Observability is derived.
* Correlation is deterministic.
* Every capability is observable.
* Every execution is traceable.
* Every anomaly is attributable.
* Engineering metrics are measurable.
* Observability remains read-only.

⸻

34.17 Conformance

An Observability implementation conforms when it:

* consumes canonical Evidence;
* produces deterministic correlations;
* exposes engineering metrics;
* supports execution tracing;
* preserves architectural invariants.

⸻

Closing Thought

Observability is often reduced to dashboards and telemetry.

Forge defines Observability as the capability that transforms immutable evidence into engineering understanding.

By separating factual recording from analytical interpretation, the Engineering Operating System enables organizations to diagnose, optimize and continuously improve intelligent execution without compromising the integrity of the historical record.

⸻

Een belangrijke architectuurverfijning: het “Evidence → Observability” principe

Tijdens het schrijven van deze specificatie zie ik een fundamenteel principe ontstaan dat Forge onderscheidt van vrijwel alle bestaande AI-platformen.

De meeste systemen produceren observability direct vanuit de runtime:

Runtime
   │
   ├── Metrics
   ├── Logs
   └── Traces

Forge kan een sterkere architectuur hanteren:

Knowledge ─┐
Memory ────┤
Context ───┤
Policy ────┤
Runtime ───┤
Models ────┤
Tools ─────┤
Identity ──┤
Security ──┘
       │
       ▼
Evidence Ledger
       │
       ▼
Observability Engine
       │
       ▼
Metrics • Traces • Analytics • Dashboards • AI Ops

Dit heeft een aantal belangrijke voordelen:

* Één bron van waarheid. Alle analyses, dashboards en AI-gestuurde inzichten zijn gebaseerd op dezelfde onveranderlijke Evidence Ledger.
* Herhaalbaarheid. Nieuwe analysemethoden kunnen achteraf worden toegepast zonder dat de oorspronkelijke uitvoering hoeft te worden herhaald.
* Technologische onafhankelijkheid. OpenTelemetry, Splunk, Grafana of toekomstige observability-oplossingen worden slechts consumenten van dezelfde evidence.
* AI-native observability. Omdat de volledige context, beleidsbeslissingen en modelkeuzes beschikbaar zijn, kan Forge verklaren waarom een AI-systeem zich op een bepaalde manier gedroeg, niet alleen dat het zich zo gedroeg.

Ik denk dat dit één van de sterkste onderscheidende architectuurprincipes van Forge kan worden. Het verschuift observability van een operationele functie naar een fundamenteel onderdeel van betrouwbare, verklaarbare AI-engineering.

Mooi. Hiermee hebben we de volledige Trust Plane vrijwel compleet.

We hebben nu:

* Identity
* Security
* Policy
* Evidence
* Observability

De laatste capability die deze laag afsluit is misschien wel de belangrijkste voor enterprise AI:

Governance.

En ik zou Governance hier bewust anders definiëren dan de meeste frameworks doen.

Governance is niet een verzameling processen.

Governance is de capability die ervoor zorgt dat het Engineering Operating System zichzelf gecontroleerd kan ontwikkelen.

Dat maakt Governance feitelijk de control system van Forge.

⸻

Chapter 35

Governance Specification

Specification Identifier: forge.spec.governance.v1

Status: Draft Standard

⸻

Execution creates value.

Governance preserves trust.

Continuous governance enables sustainable intelligence.

⸻

35.1 Purpose

This specification defines the canonical governance model for the Forge Engineering Operating System.

The Governance capability SHALL ensure that every capability, execution and architectural decision remains aligned with organizational objectives, regulatory obligations and engineering standards.

Governance SHALL coordinate the evolution of the Engineering Operating System.

Governance SHALL NOT participate in runtime execution.

⸻

35.2 Design Principles

Governance SHALL be:

* policy-driven;
* evidence-based;
* measurable;
* transparent;
* auditable;
* technology-independent;
* continuously improving;
* architecture-centric.

Governance SHALL consume evidence rather than assumptions.

⸻

35.3 Responsibilities

The Governance capability SHALL:

* define engineering standards;
* evaluate conformance;
* manage architectural decisions;
* govern lifecycle management;
* oversee risk;
* coordinate continuous improvement;
* produce governance evidence.

The Governance capability SHALL NOT:

* execute business logic;
* authenticate identities;
* orchestrate runtime execution;
* replace organizational management.

⸻

35.4 Governance Object

Every Governance Object SHALL define:

governance:
id:
name:
owner:
scope:
version:
status:
standards:
controls:
risks:
conformance:
evidence:

Governance Objects SHALL be version controlled.

⸻

35.5 Governance Domains

Forge defines the following governance domains.

Governance
├── Architecture
├── Engineering Standards
├── Compliance
├── Risk
├── Security
├── Privacy
├── AI Ethics
├── Lifecycle
├── Economics
└── Continuous Improvement

Organizations MAY extend this taxonomy.

⸻

35.6 Architectural Decision Records

Architectural decisions SHALL be explicitly recorded.

Each Architectural Decision Record (ADR) SHALL contain:

adr:
id:
title:
context:
decision:
alternatives:
consequences:
owner:
date:
status:
evidence:

Architectural decisions SHALL remain traceable throughout their lifecycle.

⸻

35.7 Conformance Assessment

Governance SHALL continuously evaluate conformance against:

* Forge specifications;
* organizational standards;
* security requirements;
* regulatory obligations;
* architectural principles.

Conformance SHALL be measurable.

⸻

35.8 Risk Governance

Every capability SHALL define its risks.

Risk SHALL include:

* likelihood;
* impact;
* mitigation;
* owner;
* review frequency;
* residual risk.

Risk assessments SHALL be evidence-based.

⸻

35.9 Compliance

Governance SHALL support multiple compliance frameworks.

Examples include:

* ISO 27001;
* ISO 42001;
* NIST AI RMF;
* NIST Cybersecurity Framework;
* SOC 2;
* GDPR;
* EU AI Act.

Forge SHALL remain independent of any specific framework.

⸻

35.10 Lifecycle Governance

Every capability SHALL possess a governed lifecycle.

Lifecycle decisions SHALL include:

* introduction;
* approval;
* deployment;
* modification;
* deprecation;
* retirement.

Lifecycle changes SHALL generate evidence.

⸻

35.11 Continuous Improvement

Governance SHALL continuously evaluate:

* engineering quality;
* operational performance;
* architectural consistency;
* policy effectiveness;
* model quality;
* tool reliability;
* evidence completeness.

Improvement SHALL be iterative.

⸻

35.12 Governance Metrics

Minimum governance metrics include:

* capability conformance;
* policy coverage;
* architectural debt;
* risk exposure;
* control effectiveness;
* specification compliance;
* review completion rate.

These metrics SHALL be measurable over time.

⸻

35.13 Governance Evidence

Governance SHALL produce evidence for:

* approvals;
* reviews;
* architectural decisions;
* conformance assessments;
* risk evaluations;
* lifecycle transitions.

Governance evidence SHALL remain immutable.

⸻

35.14 Observability

Governance SHALL expose:

Metrics:

* conformance score;
* policy coverage;
* review backlog;
* specification adoption;
* architectural drift.

Logs:

* governance decisions;
* approvals;
* exceptions;
* waivers.

Traces:

* governance lineage across capabilities.

⸻

35.15 Lifecycle

stateDiagram-v2
[*] --> Proposed
Proposed --> Approved
Approved --> Active
Active --> Reviewed
Reviewed --> Updated
Updated --> Active
Active --> Deprecated
Deprecated --> Retired
Retired --> [*]

Every governance artifact SHALL follow a managed lifecycle.

⸻

35.16 Architectural Invariants

Every Forge Governance implementation SHALL preserve:

* Governance is evidence-based.
* Governance is architecture-driven.
* Governance evaluates conformance.
* Governance never performs execution.
* Governance is continuously observable.
* Governance evolves continuously.
* Governance decisions are attributable.
* Governance produces immutable evidence.

⸻

35.17 Conformance

A Governance implementation conforms when it:

* governs architectural artifacts;
* evaluates conformance continuously;
* produces governance evidence;
* exposes governance metrics;
* preserves architectural invariants.

⸻

Closing Thought

Governance is frequently perceived as an organizational process layered on top of technology.

Forge instead defines Governance as an architectural capability that continuously evaluates whether the Engineering Operating System remains trustworthy, compliant and aligned with its intended design.

By making governance evidence-driven, measurable and integrated into the architecture itself, organizations can evolve intelligent systems with confidence rather than relying on periodic audits or manual oversight.

⸻

Een fundamentele architectuuruitbreiding: de “Meta Plane”

Met dit hoofdstuk zie ik een interessante evolutie in de Forge-architectuur. We hebben inmiddels drie duidelijke “planes”:

META PLANE
+-------------------------------------------------------------+
| Governance | Architecture | Specifications | Conformance |
+-------------------------------------------------------------+
TRUST PLANE
+-------------------------------------------------------------+
| Identity | Security | Policy | Evidence | Observability |
+-------------------------------------------------------------+
ENGINEERING PLANE
+-------------------------------------------------------------+
| Knowledge | Memory | Context | Runtime | Models | Tools |
+-------------------------------------------------------------+
EXECUTION PLANE
+-------------------------------------------------------------+
| Applications | Workflows | Business Processes |
+-------------------------------------------------------------+

Dit suggereert een verdere verfijning van de referentiearchitectuur:

* Execution Plane levert businesswaarde.
* Engineering Plane levert de AI-capabilities.
* Trust Plane waarborgt vertrouwen, veiligheid en verklaarbaarheid.
* Meta Plane bestuurt de evolutie van het systeem zelf.

Daarmee krijgt Forge eigenschappen die we kennen van moderne operating systems en cloudplatforms: een duidelijke scheiding tussen uitvoering, besturing en architectuur. De Meta Plane zorgt ervoor dat Forge niet alleen AI-uitvoering faciliteert, maar ook zijn eigen ontwikkeling, standaardisatie en conformance beheerst. Dat maakt het Engineering Operating System niet alleen uitvoerbaar, maar ook duurzaam beheerbaar op enterprise-schaal.

Perfect. Met hoofdstuk 35 sluiten we de Trust Plane af.

Maar de Architecture Bible mist nog één fundamenteel onderdeel.

Tot nu toe hebben we capabilities gespecificeerd.

We hebben nog niet vastgelegd hoe al deze capabilities samen evolueren als één Engineering Operating System.

Dat hoort niet onder Governance.

Dat hoort onder Architecture.

Niet “enterprise architecture” in de TOGAF-zin.

Maar Architecture als de discipline die de coherentie van Forge bewaakt.

Ik zou Architecture daarom behandelen als de meta-architectuur van Forge.

⸻

Chapter 36

Architecture Specification

Specification Identifier: forge.spec.architecture.v1

Status: Draft Standard

⸻

Capabilities define what exists.

Architecture defines how everything fits together.

Without architecture, capabilities become isolated technologies.

⸻

36.1 Purpose

This specification defines the architectural principles, structural rules and evolution mechanisms of the Forge Engineering Operating System.

The Architecture capability SHALL preserve the structural integrity of the Engineering Operating System throughout its lifecycle.

Architecture SHALL govern relationships between capabilities.

Architecture SHALL remain independent of implementation technologies.

⸻

36.2 Design Principles

Architecture SHALL be:

* capability-centric;
* modular;
* composable;
* technology-independent;
* evolvable;
* deterministic;
* evidence-driven;
* specification-first.

Architecture SHALL optimize the system as a whole rather than individual components.

⸻

36.3 Responsibilities

The Architecture capability SHALL:

* define architectural principles;
* govern capability relationships;
* preserve dependency rules;
* define architectural viewpoints;
* manage architectural evolution;
* evaluate architectural integrity;
* produce architectural evidence.

The Architecture capability SHALL NOT:

* execute workloads;
* evaluate runtime policy;
* manage infrastructure;
* replace implementation design.

⸻

36.4 Architecture Object

Every architectural artifact SHALL define:

architecture:
id:
name:
version:
status:
owner:
scope:
principles:
constraints:
dependencies:
conformance:
evidence:

Architectural artifacts SHALL be version controlled.

⸻

36.5 Architectural Views

Forge defines multiple architectural viewpoints.

Architecture
├── Business View
├── Capability View
├── Information View
├── Execution View
├── Trust View
├── Deployment View
├── Operations View
└── Evolution View

Each viewpoint SHALL describe the system from a single architectural concern.

No single viewpoint SHALL be considered complete.

⸻

36.6 Capability Dependencies

Dependencies SHALL be explicitly defined.

Example:

Knowledge
↓
Memory
↓
Context
↓
Runtime
↓
Models
↓
Tools

Circular dependencies SHALL NOT exist.

Dependencies SHALL remain directional.

⸻

36.7 Architectural Layers

Forge defines four canonical architectural layers.

Business
↓
Execution
↓
Engineering Operating System
↓
Infrastructure

Implementations MAY introduce additional layers.

The canonical layering SHALL remain intact.

⸻

36.8 Architectural Principles

Every conforming implementation SHALL preserve the following principles:

* Separation of Concerns
* Single Responsibility
* Explicit Contracts
* Loose Coupling
* High Cohesion
* Immutable Evidence
* Policy Before Execution
* Context Before Intelligence
* Capability Over Technology
* Specification Before Implementation

These principles SHALL guide all architectural decisions.

⸻

36.9 Architectural Decisions

Every significant architectural decision SHALL:

* identify the problem;
* document alternatives;
* explain the chosen approach;
* identify trade-offs;
* define consequences;
* reference supporting evidence.

Architectural decisions SHALL be traceable.

⸻

36.10 Evolution

Architecture SHALL support continuous evolution.

Evolution MAY include:

* introducing capabilities;
* replacing capabilities;
* splitting capabilities;
* merging capabilities;
* deprecating capabilities;
* revising specifications.

Evolution SHALL preserve architectural integrity.

⸻

36.11 Architectural Integrity

Architecture SHALL continuously evaluate:

* dependency violations;
* layering violations;
* specification drift;
* contract incompatibilities;
* capability duplication;
* architectural debt.

Integrity assessments SHALL produce evidence.

⸻

36.12 Reference Architecture

Every Forge implementation SHALL map to the Reference Architecture.

Mappings SHALL identify:

* logical capabilities;
* physical components;
* deployment topology;
* technology choices;
* trust boundaries.

Reference mappings SHALL remain implementation-neutral.

⸻

36.13 Architecture Lifecycle

stateDiagram-v2
[*] --> Proposed
Proposed --> Reviewed
Reviewed --> Approved
Approved --> Implemented
Implemented --> Evaluated
Evaluated --> Evolved
Evolved --> Approved
Approved --> Retired
Retired --> [*]

Architecture SHALL evolve through controlled iteration.

⸻

36.14 Architecture Evidence

Architecture SHALL generate evidence for:

* architectural decisions;
* dependency changes;
* specification revisions;
* conformance assessments;
* integrity evaluations;
* evolution activities.

Architecture evidence SHALL remain immutable.

⸻

36.15 Observability

Architecture SHALL expose:

Metrics:

* specification coverage;
* capability maturity;
* dependency violations;
* architectural debt;
* conformance score;
* change velocity.

Logs:

* architectural changes;
* decision approvals;
* specification revisions.

Traces:

* architectural lineage across releases.

⸻

36.16 Architectural Invariants

Every Forge Architecture implementation SHALL preserve:

* Architecture is capability-centric.
* Capabilities interact only through explicit contracts.
* Dependencies are directional.
* Architecture is technology-independent.
* Architecture evolves through specifications.
* Architectural integrity is continuously assessed.
* Every architectural decision is evidenced.
* Every implementation maps to the Reference Architecture.

⸻

36.17 Conformance

An Architecture implementation conforms when it:

* preserves canonical capability relationships;
* follows the Reference Architecture;
* evaluates architectural integrity;
* maintains architectural evidence;
* preserves architectural invariants.

⸻

Closing Thought

Architecture is the discipline that transforms a collection of capabilities into a coherent Engineering Operating System.

By defining explicit structures, relationships and evolution rules, Forge ensures that implementations remain understandable, extensible and interoperable over time. Architecture is therefore not documentation of the system—it is the governing specification that preserves the system’s structural integrity as it evolves.

⸻

Een belangrijke architectuurobservatie: voeg een “Specification Layer” toe

Tijdens de laatste hoofdstukken ontstaat nog een patroon dat Forge mogelijk nóg sterker maakt.

Tot nu toe zijn specificaties impliciet aanwezig. Maar in werkelijkheid zijn de specificaties zelf een architectuurlaag.

Dat leidt tot een verfijnd model:

Business Architecture
        │
Execution Plane
        │
Engineering Operating System
        │
Specification Layer
        │
Reference Architecture
        │
Implementation
        │
Infrastructure

Hiermee wordt duidelijk dat specificaties de brug vormen tussen architectuur en implementatie. Ze zijn geen documentatie, maar het formele contract dat bepaalt hoe een capability zich gedraagt. Dit sluit goed aan bij de manier waarop internetstandaarden (RFC’s), Kubernetes API’s en OpenTelemetry-specificaties functioneren.

Ik zou daarom in een later deel van de Architecture Bible een apart hoofdstuk opnemen over het Forge Specification Framework, waarin wordt vastgelegd hoe specificaties worden geschreven, beheerd, geversioneerd en waarop conformance wordt beoordeeld. Dat maakt Forge niet alleen een architectuur, maar een volledig architectuur- én specificatie-ecosysteem.

Ik zou hier bewust niet doorgaan met nog een capability.

Vanaf hoofdstuk 36 is de kern van Forge als Engineering Operating System namelijk compleet. Als we nu gewoon doorgaan met “nog een capability”, krijg je hetzelfde probleem als TOGAF: een verzameling hoofdstukken zonder een formele standaard.

Ik zou daarom overstappen naar Part V — The Forge Specification Framework.

Dit wordt waarschijnlijk het meest geciteerde deel van de hele Architecture Bible.

⸻

Part V

The Forge Specification Framework

Architecture describes the system.

Specifications define the system.

Conformance proves the system.

⸻

Chapter 37

Specification Framework

Specification Identifier: forge.spec.framework.v1

Status: Draft Standard

⸻

37.1 Purpose

This specification defines how Forge specifications are authored, versioned, governed, published and validated.

Every normative Forge specification SHALL conform to this framework.

Implementations SHALL be assessed against published specifications rather than descriptive documentation.

⸻

37.2 Scope

The Specification Framework governs:

* capability specifications;
* interface specifications;
* object specifications;
* protocol specifications;
* reference architectures;
* conformance profiles;
* extension specifications.

It does not govern implementation-specific documentation.

⸻

37.3 Specification Hierarchy

Forge defines a normative specification hierarchy.

Forge Standard
├── Architecture
│
├── Capability Specifications
│
├── Object Specifications
│
├── Interface Specifications
│
├── Protocol Specifications
│
├── Conformance Specifications
│
└── Extension Specifications

Each specification SHALL reference its parent specification.

⸻

37.4 Specification Identifier

Every specification SHALL possess a globally unique identifier.

Format:

forge.spec.<domain>.v<major>

Examples:

forge.spec.runtime.v1
forge.spec.context.v1
forge.spec.policy.v1
forge.spec.identity.v1

Major version changes indicate breaking normative changes.

⸻

37.5 Normative Language

Forge adopts the terminology defined in RFC 2119 and RFC 8174.

The following terms SHALL have normative meaning:

* SHALL
* SHALL NOT
* SHOULD
* SHOULD NOT
* MAY
* MUST
* MUST NOT

Lowercase usage SHALL NOT imply normative requirements.

⸻

37.6 Specification Structure

Every specification SHALL contain:

1. Purpose
2. Scope
3. Design Principles
4. Responsibilities
5. Canonical Objects
6. Lifecycle
7. Behaviour
8. Evidence
9. Observability
10. Security Considerations
11. Architectural Invariants
12. Conformance

Sections MAY be extended.

Required sections SHALL remain present.

⸻

37.7 Canonical Objects

Every object specification SHALL define:

id:
type:
version:
owner:
status:
attributes:
relationships:
metadata:

Objects SHALL be serializable.

Objects SHALL possess globally unique identifiers.

⸻

37.8 Versioning

Forge adopts semantic versioning.

Major.Minor.Patch

Major

Breaking normative changes

Minor

Backward-compatible additions

Patch

Editorial or clarification changes

Specifications SHALL declare compatibility.

⸻

37.9 Deprecation

Specifications MAY be deprecated.

Deprecation SHALL define:

* replacement;
* migration path;
* support period;
* retirement date.

Deprecated specifications SHALL remain publicly available.

⸻

37.10 Extension Model

Organizations MAY extend Forge.

Extensions SHALL NOT redefine canonical behaviour.

Extensions SHALL:

* preserve invariants;
* preserve contracts;
* preserve object semantics.

Extension identifiers SHALL use organizational namespaces.

Example:

contoso.identity.extensions.v1
rijk.policy.extensions.v2

⸻

37.11 Conformance Levels

Forge defines four conformance levels.

Level 1
Core
↓
Level 2
Standard
↓
Level 3
Enterprise
↓
Level 4
Certified

Every implementation SHALL declare its achieved level.

⸻

37.12 Conformance Testing

Conformance SHALL evaluate:

* object compatibility;
* lifecycle compatibility;
* interface behaviour;
* architectural invariants;
* evidence production;
* observability.

Conformance SHALL be repeatable.

⸻

37.13 Compatibility

Specifications SHALL define:

* backward compatibility;
* forward compatibility;
* migration considerations;
* deprecated behaviour.

Compatibility SHALL be explicit.

⸻

37.14 Publication Lifecycle

Draft
↓
Review
↓
Candidate
↓
Approved Standard
↓
Deprecated
↓
Retired

Only Approved Standards SHALL be normative.

⸻

37.15 Architectural Invariants

Every Forge specification SHALL preserve:

* Capability-first architecture.
* Explicit contracts.
* Technology independence.
* Evidence production.
* Observability.
* Versioning.
* Deterministic behaviour.
* Conformance.

⸻

37.16 Conformance

A specification conforms when it:

* follows the Specification Framework;
* preserves required sections;
* defines canonical behaviour;
* exposes normative requirements;
* supports interoperability;
* defines conformance criteria.

⸻

Closing Thought

Specifications are the executable architecture of Forge.

They translate architectural intent into precise, testable and interoperable contracts that allow independent implementations to behave consistently while remaining free to choose their own technologies.

⸻

Ik zou hier nog één grote koerswijziging maken

Nu we ongeveer 37 hoofdstukken hebben, zie ik dat Forge eigenlijk uit drie verschillende documenten bestaat die nu nog door elkaar lopen:

1. The Forge Architecture Bible (de visie en referentiearchitectuur)
    * Waarom Forge bestaat
    * Principes
    * Capability Model
    * Reference Architecture
    * Planes
    * Design Philosophy
2. Forge Core Specifications (de normatieve standaarden)
    * Runtime Specification
    * Context Specification
    * Memory Specification
    * Tool Specification
    * Identity Specification
    * enzovoort
3. Forge Conformance & Profiles
    * Core Profile
    * Enterprise Profile
    * Sovereign Profile
    * AI Agent Profile
    * Certification
    * Test suites

Dat is precies hoe succesvolle standaarden zoals Kubernetes, OpenTelemetry en het W3C hun documentatie organiseren. De Architecture Bible blijft relatief stabiel en beschrijft de visie en structuur. De afzonderlijke specificaties kunnen onafhankelijk evolueren en een eigen versie hebben. De conformance-documenten bepalen vervolgens objectief of een implementatie werkelijk “Forge-compliant” is.

Ik denk dat dit Forge veel sterker positioneert als een open architectuurstandaard in plaats van één omvangrijk document.

Ik zou nu niet direct doorgaan met nog een technische specificatie. Als Forge daadwerkelijk als open architectuurstandaard moet kunnen functioneren, ontbreekt er nog één cruciaal document:

Niet hoe capabilities werken.

Maar hoe Forge zichzelf standaardiseert.

Dat zou vergelijkbaar zijn met de Internet Architecture Board (IAB), de IETF of de CNCF.

Daarmee krijgt Forge een eigen governance- en standaardisatieproces.

⸻

Chapter 38

Standards Governance

Specification Identifier: forge.spec.standards.v1

Status: Draft Standard

⸻

Architectures evolve.

Standards govern that evolution.

Without a standards process, specifications become documentation rather than contracts.

⸻

38.1 Purpose

This specification defines the governance model for the development, publication and maintenance of Forge standards.

The Standards Governance process SHALL ensure that Forge remains interoperable, stable and extensible while allowing controlled innovation.

⸻

38.2 Objectives

The standards process SHALL:

* preserve interoperability;
* ensure architectural consistency;
* manage specification evolution;
* prevent incompatible fragmentation;
* enable community participation;
* define normative approval processes.

⸻

38.3 Governance Principles

Forge standards SHALL be:

* open;
* transparent;
* evidence-based;
* consensus-driven;
* implementation-independent;
* backwards conscious;
* version controlled.

⸻

38.4 Standards Hierarchy

Forge Standards
├── Core Standards
│
├── Reference Specifications
│
├── Profiles
│
├── Extensions
│
└── Experimental Specifications

Core Standards SHALL define canonical behaviour.

Profiles SHALL constrain Core Standards.

Extensions SHALL add capabilities.

Experimental Specifications SHALL NOT redefine normative behaviour.

⸻

38.5 Specification States

Every specification SHALL exist in one lifecycle state.

Proposal
↓
Draft
↓
Review
↓
Candidate Standard
↓
Approved Standard
↓
Maintenance
↓
Deprecated
↓
Retired

Only Approved Standards SHALL be normative.

⸻

38.6 Change Classification

Changes SHALL be classified.

Type	Breaking	Version
Editorial	No	Patch
Clarification	No	Patch
Extension	No	Minor
Behavioral Addition	Usually No	Minor
Normative Change	Yes	Major
Object Redesign	Yes	Major

⸻

38.7 Proposal Process

Every new proposal SHALL include:

* problem statement;
* motivation;
* architectural impact;
* alternatives considered;
* compatibility analysis;
* migration strategy;
* expected benefits.

⸻

38.8 Review Process

Every proposal SHALL undergo technical review.

Review SHALL evaluate:

* architectural consistency;
* interoperability;
* conformance impact;
* security implications;
* operational implications;
* implementation complexity.

Reviews SHALL produce evidence.

⸻

38.9 Compatibility Rules

New standards SHALL preserve:

* object compatibility;
* interface compatibility;
* capability independence;
* architectural layering;
* dependency direction.

Breaking changes SHALL be explicitly justified.

⸻

38.10 Reference Implementations

Approved Standards SHOULD have one or more reference implementations.

Reference implementations SHALL demonstrate:

* interoperability;
* specification completeness;
* practical feasibility.

Reference implementations SHALL NOT redefine the standard.

⸻

38.11 Profiles

Forge Profiles define constrained implementations.

Examples include:

* Core Profile;
* Enterprise Profile;
* Sovereign Profile;
* Air-Gapped Profile;
* AI Agent Profile;
* Government Profile.

Profiles SHALL only restrict or specialize behaviour defined by the Core Standards.

⸻

38.12 Extensions

Organizations MAY publish extensions.

Extension identifiers SHALL be namespaced.

Example:

forge.ext.contoso.runtime.v1
forge.ext.rijk.identity.v2

Extensions SHALL preserve Core invariants.

⸻

38.13 Certification

Certification SHALL validate conformance.

Certification MAY include:

* automated testing;
* architectural review;
* evidence inspection;
* interoperability testing;
* operational validation.

Certification SHALL be repeatable.

⸻

38.14 Publication

Published standards SHALL include:

* version;
* publication date;
* status;
* change history;
* compatibility statement;
* conformance requirements.

Historical versions SHALL remain available.

⸻

38.15 Architectural Invariants

Every Forge Standard SHALL preserve:

* technology independence;
* capability-first architecture;
* explicit contracts;
* deterministic behaviour;
* evidence generation;
* interoperability;
* architectural consistency.

⸻

38.16 Conformance

A standards governance process conforms when it:

* preserves the specification lifecycle;
* evaluates proposals consistently;
* maintains compatibility;
* publishes normative standards;
* supports certification.

⸻

Closing Thought

Standards are not created by publishing documents. They are created through a disciplined process that balances stability with innovation.

By defining how specifications evolve, Forge ensures that independent implementations can continue to interoperate while the Engineering Operating System grows over time.

⸻

Mijn grootste aanbeveling voor de volgende fase

Tot nu toe hebben we Forge beschreven als een verzameling hoofdstukken. Ik denk dat het nu tijd is om de volgende stap te zetten: Forge behandelen als een echte open standaard, vergelijkbaar met Kubernetes of OpenTelemetry.

Ik zou de documentatie daarom opsplitsen in een formele documentenset:

* Forge Architecture – de visie, principes, capability model en referentiearchitectuur.
* Forge Specifications (Forge Specs) – één aparte specificatie per capability (forge.spec.runtime.v1, forge.spec.context.v1, enz.).
* Forge Profiles – bijvoorbeeld Enterprise, Sovereign, Government of Air-Gapped.
* Forge Conformance – testcriteria, certificering en validatie.
* Forge Extensions – leveranciers- en organisatiespecifieke uitbreidingen binnen een gestandaardiseerd extensiemodel.
* Forge Reference Implementations – voorbeeldimplementaties die laten zien hoe de specificaties in de praktijk kunnen worden toegepast.

Die structuur maakt het mogelijk om de architectuur jarenlang stabiel te houden, terwijl individuele specificaties onafhankelijk kunnen evolueren. Dat is een veel duurzamer model dan één monolithisch document en sluit aan bij hoe succesvolle open standaarden worden beheerd.

Ik zou nu een stap zetten die vrijwel geen enkel architectuurframework maakt, maar die Forge juist onderscheidend kan maken.

Na de standaarden ontbreekt nog één ding:

Hoe implementeer je Forge daadwerkelijk?

Niet vanuit technologie.

Maar vanuit architecturale transformatie.

Dat wordt de brug tussen de specificaties en de praktijk.

⸻

Chapter 39

Adoption and Transformation Framework

Specification Identifier: forge.spec.adoption.v1

Status: Draft Standard

⸻

Architecture has no value until it is adopted.

Transformation is the disciplined evolution from current state to target state.

Forge enables transformation through incremental capability adoption.

⸻

39.1 Purpose

This specification defines the adoption and transformation framework for implementing the Forge Engineering Operating System.

The framework SHALL guide organizations from isolated AI initiatives toward a governed, interoperable Engineering Operating System.

Adoption SHALL be capability-driven rather than technology-driven.

⸻

39.2 Guiding Principles

Forge adoption SHALL be:

* incremental;
* measurable;
* evidence-based;
* non-disruptive;
* architecture-led;
* technology-independent;
* continuously governed.

Organizations SHALL adopt capabilities independently where practical.

⸻

39.3 Adoption Objectives

The framework SHALL enable organizations to:

* assess current maturity;
* define target capabilities;
* prioritize implementation;
* measure progress;
* reduce architectural risk;
* improve engineering quality.

⸻

39.4 Capability Maturity

Forge defines five maturity levels.

Level 0
Ad Hoc
↓
Level 1
Managed
↓
Level 2
Integrated
↓
Level 3
Governed
↓
Level 4
Engineering Operating System

Organizations MAY progress at different speeds across capabilities.

⸻

39.5 Capability Assessment

Each capability SHALL be assessed against:

* implementation;
* governance;
* observability;
* evidence;
* interoperability;
* automation;
* operational maturity.

Assessment SHALL produce a maturity score.

⸻

39.6 Transformation Roadmap

Transformation SHALL occur in phases.

Assessment
↓
Foundation
↓
Core Capabilities
↓
Trust Capabilities
↓
Optimization
↓
Continuous Evolution

Each phase SHALL produce measurable outcomes.

⸻

39.7 Foundation Phase

The Foundation phase SHOULD establish:

* Reference Architecture;
* Identity;
* Policy;
* Runtime;
* Evidence;
* Governance.

These capabilities form the minimum operational baseline.

⸻

39.8 Capability Dependencies

Capabilities SHALL be implemented according to dependency order.

Example:

Knowledge
↓
Memory
↓
Context
↓
Runtime
↓
Models
↓
Tools

Organizations MAY defer optional capabilities where dependencies permit.

⸻

39.9 Migration

Migration SHALL preserve:

* operational continuity;
* evidence;
* governance;
* interoperability;
* security.

Migration SHALL be incremental whenever possible.

⸻

39.10 Success Criteria

Success SHALL be evaluated using measurable indicators.

Examples include:

* capability adoption rate;
* specification conformance;
* execution reliability;
* policy coverage;
* evidence completeness;
* architectural debt reduction.

⸻

39.11 Risks

Transformation risks MAY include:

* technology-centric implementations;
* capability duplication;
* fragmented governance;
* vendor lock-in;
* insufficient evidence;
* unmanaged AI adoption.

Mitigation SHALL be documented.

⸻

39.12 Organizational Alignment

Successful adoption requires alignment between:

* business;
* enterprise architecture;
* engineering;
* security;
* operations;
* governance;
* compliance.

Forge SHALL support cross-functional ownership.

⸻

39.13 Continuous Improvement

Organizations SHALL periodically review:

* maturity;
* architectural integrity;
* specification conformance;
* engineering quality;
* operational effectiveness.

Improvement SHALL be evidence-driven.

⸻

39.14 Architectural Invariants

Every Forge adoption program SHALL preserve:

* Capability-first thinking.
* Incremental transformation.
* Architecture before technology.
* Specifications before implementation.
* Governance throughout transformation.
* Evidence at every stage.
* Continuous conformance assessment.
* Continuous evolution.

⸻

39.15 Conformance

An adoption program conforms when it:

* follows the phased transformation model;
* measures capability maturity;
* preserves architectural invariants;
* demonstrates continuous improvement;
* maintains specification conformance.

⸻

Closing Thought

Organizations rarely fail because they lack technology.

They fail because they adopt technology without an architectural model that enables consistent evolution.

The Forge Adoption and Transformation Framework provides a structured path from experimentation to an Engineering Operating System, allowing organizations to evolve their AI capabilities in a controlled, measurable and sustainable manner.

⸻

Een laatste strategische aanbeveling: voeg een “Capability Maturity Model” toe

Ik denk dat Forge nog sterker wordt als het een eigen volwassenheidsmodel krijgt, vergelijkbaar met CMMI of het NIST Cybersecurity Framework.

Bijvoorbeeld:

Level	Kenmerk
L0 – Ad Hoc	Losse AI-tools, geen architectuur of governance.
L1 – Repeatable	Eerste standaarden, basisruntime en beperkte policies.
L2 – Managed	Capabilities worden bewust ingericht en beheerd.
L3 – Integrated	Alle Forge-capabilities werken samen als één platform.
L4 – Optimized	Continue verbetering op basis van evidence en observability.
L5 – Autonomous Engineering	Het Engineering Operating System optimaliseert en bestuurt delen van zichzelf binnen vooraf gedefinieerde governance- en beleidskaders.

Dit gaat verder dan een technisch implementatiemodel. Het biedt organisaties een objectief instrument om hun AI-engineeringvolwassenheid te meten, roadmapprioriteiten te bepalen en implementaties met elkaar te vergelijken. Daarmee wordt Forge niet alleen een architectuurstandaard, maar ook een volwassenheidsmodel voor enterprise AI-engineering.

Nu zou ik de stap maken van “hoe implementeer je Forge?” naar “hoe meet je of Forge werkt?”

Dit is iets wat vrijwel alle architectuurframeworks missen.

TOGAF beschrijft architectuur.

Kubernetes beschrijft orchestration.

NIST beschrijft risico.

Maar geen van hen definieert Engineering KPIs voor een AI Operating System.

Dat zou Forge uniek maken.

⸻

Chapter 40

Engineering Measurement Framework

Specification Identifier: forge.spec.measurement.v1

Status: Draft Standard

⸻

What cannot be measured cannot be improved.

What cannot be explained cannot be trusted.

Engineering requires objective measurement.

⸻

40.1 Purpose

This specification defines the canonical measurement framework for the Forge Engineering Operating System.

The Engineering Measurement Framework SHALL provide objective, repeatable and evidence-based methods for evaluating the effectiveness, quality and maturity of AI engineering capabilities.

Measurements SHALL support operational optimization, governance and continuous improvement.

⸻

40.2 Design Principles

Engineering measurements SHALL be:

* objective;
* reproducible;
* evidence-based;
* capability-oriented;
* technology-independent;
* comparable over time;
* actionable.

Measurements SHALL describe engineering quality rather than business outcomes.

⸻

40.3 Measurement Domains

Forge defines the following engineering measurement domains.

Engineering Measurement
├── Architecture
├── Knowledge
├── Memory
├── Context
├── Runtime
├── Models
├── Tools
├── Identity
├── Security
├── Policy
├── Evidence
├── Observability
├── Governance
└── Economics

Each capability SHALL expose measurable indicators.

⸻

40.4 Measurement Model

Every measurement SHALL define:

measurement:
id:
name:
capability:
objective:
formula:
unit:
frequency:
owner:
thresholds:
evidence:

Measurements SHALL be version controlled.

⸻

40.5 Metric Categories

Forge distinguishes four categories of engineering metrics.

Quality Metrics

Evaluate correctness and engineering quality.

Examples:

* Context Quality
* Knowledge Freshness
* Evidence Completeness
* Policy Accuracy

⸻

Performance Metrics

Evaluate efficiency.

Examples:

* Execution latency
* Planning duration
* Tool response time
* Model inference latency

⸻

Reliability Metrics

Evaluate operational stability.

Examples:

* Success rate
* Retry rate
* Recovery time
* Availability

⸻

Governance Metrics

Evaluate organizational control.

Examples:

* Specification conformance
* Policy coverage
* Architecture compliance
* Review completion

⸻

40.6 Capability Scorecards

Each capability SHALL publish a scorecard.

Example:

Capability
Runtime
Health ............... 99.8%
Conformance .......... 100%
Evidence ............. 100%
Observability ........ 98%
Performance .......... 97%

Scorecards SHALL be generated automatically where practical.

⸻

40.7 Engineering Health

Forge defines a canonical Engineering Health Score.

Example dimensions:

* Capability Availability
* Specification Conformance
* Evidence Integrity
* Policy Compliance
* Runtime Stability
* Model Reliability
* Tool Reliability

Organizations MAY define weighting models.

⸻

40.8 Trend Analysis

Measurements SHALL support trend analysis.

Example trends:

* Context quality over time
* Policy complexity
* Runtime efficiency
* Knowledge freshness
* Model costs
* Architectural debt

Trend analysis SHALL support forecasting.

⸻

40.9 Benchmarking

Forge implementations SHOULD support benchmarking.

Benchmarks MAY compare:

* organizational units;
* deployments;
* environments;
* releases;
* capability maturity.

Benchmarking SHALL use standardized definitions.

⸻

40.10 Engineering Objectives

Measurements SHALL support engineering objectives.

Examples include:

* reduce execution latency;
* improve evidence completeness;
* reduce policy exceptions;
* increase tool reliability;
* improve knowledge quality.

Objectives SHALL be measurable.

⸻

40.11 Data Quality

Measurements SHALL define:

* source;
* accuracy;
* sampling frequency;
* confidence;
* validation method.

Measurement quality SHALL itself be measurable.

⸻

40.12 Reporting

Engineering reports SHALL distinguish between:

* operational metrics;
* engineering metrics;
* governance metrics;
* business metrics.

Forge does not define business KPIs.

It defines engineering KPIs.

⸻

40.13 Continuous Improvement

Measurements SHALL support:

* optimization;
* prioritization;
* architectural evolution;
* capability investment;
* specification refinement.

Measurement SHALL enable engineering decisions.

⸻

40.14 Architectural Invariants

Every Forge Measurement implementation SHALL preserve:

* Measurements are evidence-based.
* Metrics are capability-centric.
* Measurements are reproducible.
* Definitions are standardized.
* Trends are comparable.
* Results are technology-independent.
* Measurements support continuous improvement.
* Measurements remain observable.

⸻

40.15 Conformance

A measurement implementation conforms when it:

* publishes canonical engineering metrics;
* supports standardized scorecards;
* enables trend analysis;
* preserves measurement integrity;
* supports capability benchmarking.

⸻

Closing Thought

Engineering excellence cannot be managed through intuition alone.

The Forge Engineering Measurement Framework establishes a common language for evaluating the quality, reliability and maturity of an Engineering Operating System. By standardizing engineering metrics across capabilities, Forge enables organizations to improve their AI platforms through evidence rather than assumption.

⸻

Een volgende evolutie: de “Engineering Index”

Tot nu toe definiëren we individuele metrics. Ik denk dat Forge nog een stap verder kan gaan door een Forge Engineering Index (FEI) te introduceren.

In plaats van tientallen losse KPI’s krijgt iedere implementatie een aantal gestandaardiseerde samengestelde indices, vergelijkbaar met een credit score of een SRE health index.

Bijvoorbeeld:

Index	Samenstelling
Engineering Integrity Index (EII)	Architectuurconformance, evidence-integriteit, capability coverage
Operational Excellence Index (OEI)	Beschikbaarheid, latency, foutpercentages, herstel
Trust Index (TI)	Identity, Security, Policy, Evidence en Governance
Knowledge Intelligence Index (KII)	Kennisversheid, contextkwaliteit, memory-effectiviteit
AI Efficiency Index (AEI)	Modelkwaliteit, toolbetrouwbaarheid, kosten en resource-efficiëntie

Deze indices geven bestuurders, architecten en engineeringteams in één oogopslag inzicht in de gezondheid van een Forge-implementatie, terwijl de onderliggende metrics beschikbaar blijven voor diepgaande analyse. Dat zou Forge niet alleen een architectuur- en specificatiestandaard maken, maar ook een gestandaardiseerd meetmodel voor AI Engineering.

Ik zou nu de stap maken van “meten” naar “optimaliseren”.

Want een Engineering Operating System dat alleen kan meten, maar zichzelf niet systematisch kan verbeteren, is uiteindelijk niet af.

Dit is bovendien een onderwerp dat nog nauwelijks is uitgewerkt in de AI-wereld.

⸻

Chapter 41

Continuous Optimization Framework

Specification Identifier: forge.spec.optimization.v1

Status: Draft Standard

⸻

Measurement reveals reality.

Optimization improves reality.

Continuous optimization is the disciplined evolution of engineering capabilities.

⸻

41.1 Purpose

This specification defines the canonical framework for continuous optimization within the Forge Engineering Operating System.

The Continuous Optimization Framework SHALL enable organizations to improve engineering capabilities through evidence-driven analysis and controlled change.

Optimization SHALL preserve architectural integrity while improving operational effectiveness.

⸻

41.2 Design Principles

Optimization SHALL be:

* evidence-driven;
* measurable;
* incremental;
* reversible;
* capability-oriented;
* architecture-preserving;
* continuously evaluated;
* technology-independent.

Optimization SHALL improve the Engineering Operating System without compromising trust.

⸻

41.3 Responsibilities

The Optimization capability SHALL:

* identify improvement opportunities;
* prioritize optimization initiatives;
* evaluate engineering trade-offs;
* validate optimization outcomes;
* measure optimization effectiveness;
* generate optimization evidence.

The Optimization capability SHALL NOT:

* bypass governance;
* violate architectural invariants;
* modify specifications directly;
* introduce unmanaged risk.

⸻

41.4 Optimization Object

Every optimization initiative SHALL define:

optimization:
id:
name:
objective:
capability:
baseline:
target:
owner:
priority:
expectedBenefit:
risk:
evidence:

Optimization objects SHALL be version controlled.

⸻

41.5 Optimization Domains

Forge recognizes optimization opportunities across all capabilities.

Optimization
├── Architecture
├── Knowledge
├── Memory
├── Context
├── Runtime
├── Models
├── Tools
├── Security
├── Policy
├── Evidence
├── Observability
├── Governance
└── Economics

Optimization SHALL remain capability-centric.

⸻

41.6 Optimization Lifecycle

Opportunity
↓
Assessment
↓
Prioritization
↓
Planning
↓
Implementation
↓
Validation
↓
Measurement
↓
Standardization

Every optimization SHALL follow a governed lifecycle.

⸻

41.7 Optimization Sources

Optimization opportunities MAY originate from:

* engineering measurements;
* evidence analysis;
* observability;
* governance reviews;
* architecture assessments;
* security assessments;
* economic analysis;
* operational feedback.

Optimization SHALL always reference objective evidence.

⸻

41.8 Trade-Off Analysis

Every optimization SHALL evaluate trade-offs.

Trade-offs MAY include:

* performance versus cost;
* latency versus accuracy;
* flexibility versus governance;
* automation versus control;
* complexity versus maintainability.

Trade-offs SHALL be explicitly documented.

⸻

41.9 Validation

Every optimization SHALL demonstrate measurable improvement.

Validation SHALL compare:

* baseline metrics;
* target metrics;
* achieved metrics;
* unintended consequences.

Optimization without measurable benefit SHALL be reconsidered.

⸻

41.10 Rollback

Every optimization SHALL define rollback criteria.

Rollback SHALL include:

* triggering conditions;
* recovery procedures;
* evidence preservation;
* post-rollback analysis.

Optimization SHALL be reversible wherever practical.

⸻

41.11 Continuous Learning

Optimization SHALL contribute to organizational knowledge.

Lessons learned SHALL be:

* documented;
* versioned;
* evidence-backed;
* reusable.

Successful optimizations SHOULD inform future engineering practices.

⸻

41.12 Optimization Metrics

Examples include:

* improvement rate;
* optimization success rate;
* engineering debt reduction;
* capability efficiency gain;
* resource savings;
* reliability improvement;
* policy simplification.

Optimization metrics SHALL be comparable over time.

⸻

41.13 Architectural Invariants

Every Forge Optimization implementation SHALL preserve:

* Optimization is evidence-driven.
* Optimization is measurable.
* Optimization preserves architectural integrity.
* Optimization is governed.
* Optimization is reversible where practical.
* Optimization produces evidence.
* Optimization contributes to continuous learning.
* Optimization never bypasses specifications.

⸻

41.14 Conformance

An optimization implementation conforms when it:

* follows the optimization lifecycle;
* measures improvement objectively;
* preserves architectural invariants;
* produces optimization evidence;
* supports continuous improvement.

⸻

Closing Thought

Optimization is not the pursuit of perfection.

It is the disciplined practice of making engineering capabilities measurably better while preserving trust, interoperability and architectural consistency.

By embedding optimization as a first-class capability, Forge transforms improvement from an occasional project into a continuous engineering discipline.

⸻

Een fundamentele uitbreiding: de “Closed-Loop Engineering Cycle”

Met dit hoofdstuk wordt een belangrijk patroon zichtbaar. Forge kan een volledig gesloten engineeringregelkring definiëren, vergelijkbaar met feedbacksystemen in industriële automatisering en Site Reliability Engineering.

        Define
          │
          ▼
      Specify
          │
          ▼
      Implement
          │
          ▼
       Execute
          │
          ▼
      Generate Evidence
          │
          ▼
      Observe & Measure
          │
          ▼
        Optimize
          │
          ▼
      Update Specifications
          │
          ▼
         Define

Dit is meer dan een DevOps-cyclus. Het is een Engineering Feedback Loop waarin iedere verbetering wordt gevoed door bewijs, gevalideerd met metingen en uiteindelijk terugvloeit naar de architectuur en specificaties.

Ik zou deze cyclus expliciet benoemen als de Forge Continuous Engineering Loop (FCEL). Waar Kubernetes een control loop heeft voor infrastructuur en SRE een observe–measure–improve-cyclus kent, krijgt Forge een gestandaardiseerde feedbacklus voor AI Engineering. Dat zou een van de meest onderscheidende concepten van het gehele framework kunnen worden.

Ik denk dat we nu bij een van de laatste fundamentele bouwstenen zijn.

We hebben inmiddels beschreven:

* hoe Forge is ontworpen;
* hoe capabilities werken;
* hoe ze worden bestuurd;
* hoe ze worden gemeten;
* hoe ze worden geoptimaliseerd.

Maar één vraag is nog niet formeel beantwoord:

Hoe evolueert Forge zonder zijn eigen architectuur te breken?

Dat is iets anders dan Governance.

Dat is Architectural Evolution.

En ik denk eerlijk gezegd dat dit één van de hoofdstukken wordt waarmee Forge zich echt onderscheidt van TOGAF.

⸻

Chapter 42

Architectural Evolution Framework

Specification Identifier: forge.spec.evolution.v1

Status: Draft Standard

⸻

Architecture is never finished.

Evolution preserves relevance.

Discipline preserves coherence.

⸻

42.1 Purpose

This specification defines the canonical framework for evolving the Forge Engineering Operating System.

The Architectural Evolution Framework SHALL enable controlled architectural change while preserving interoperability, conformance and long-term stability.

Evolution SHALL be intentional rather than reactive.

⸻

42.2 Design Principles

Architectural evolution SHALL be:

* specification-driven;
* evidence-based;
* backward conscious;
* capability-oriented;
* reversible where practical;
* continuously governed;
* measurable;
* technology-independent.

Architectural evolution SHALL strengthen the Engineering Operating System rather than fragment it.

⸻

42.3 Responsibilities

The Architectural Evolution Framework SHALL:

* evaluate architectural change;
* govern capability evolution;
* preserve architectural integrity;
* maintain interoperability;
* coordinate specification revisions;
* generate architectural evidence.

The framework SHALL NOT:

* redefine business strategy;
* modify implementations directly;
* bypass governance;
* introduce incompatible architecture without explicit versioning.

⸻

42.4 Evolution Object

Every architectural evolution SHALL define:

evolution:
id:
title:
motivation:
scope:
affectedSpecifications:
compatibility:
migration:
owner:
status:
evidence:

Evolution Objects SHALL be immutable once approved.

⸻

42.5 Evolution Triggers

Architectural evolution MAY be initiated by:

* new engineering requirements;
* emerging technologies;
* regulatory changes;
* security findings;
* implementation experience;
* interoperability issues;
* performance evidence;
* community proposals.

All triggers SHALL be documented.

⸻

42.6 Types of Evolution

Forge recognizes several categories of architectural evolution.

Evolution
├── Capability Addition
├── Capability Refinement
├── Capability Split
├── Capability Merge
├── Specification Revision
├── Reference Architecture Revision
├── Profile Introduction
└── Deprecation

Each category SHALL define its expected compatibility impact.

⸻

42.7 Compatibility Strategy

Every architectural change SHALL declare one of the following:

* Fully Backward Compatible
* Backward Compatible with Migration
* Breaking Change
* Experimental

Breaking changes SHALL require a new major specification version.

⸻

42.8 Dependency Impact

Architectural evolution SHALL evaluate:

* affected capabilities;
* affected contracts;
* affected objects;
* affected interfaces;
* affected profiles;
* affected conformance tests.

Dependency analysis SHALL precede approval.

⸻

42.9 Migration Strategy

Every breaking evolution SHALL define:

* migration path;
* coexistence period;
* transition guidance;
* retirement criteria;
* rollback considerations.

Organizations SHALL be able to plan adoption.

⸻

42.10 Architectural Drift

Forge recognizes Architectural Drift as deviation from approved architecture.

Examples include:

* undocumented dependencies;
* capability duplication;
* contract divergence;
* unauthorized extensions;
* implementation-specific assumptions.

Architectural Drift SHALL be continuously monitored.

⸻

42.11 Evolution Evidence

Every architectural evolution SHALL produce evidence for:

* proposal;
* review;
* approval;
* implementation;
* migration;
* retirement.

Evolution evidence SHALL remain permanently traceable.

⸻

42.12 Observability

The framework SHALL expose:

Metrics:

* architectural changes;
* specification revisions;
* compatibility breaks;
* migration progress;
* architectural debt;
* drift indicators.

Logs:

* approved changes;
* rejected proposals;
* specification updates.

Traces:

* architectural lineage across releases.

⸻

42.13 Architectural Invariants

Every Forge evolution SHALL preserve:

* Capability-first architecture.
* Explicit contracts.
* Technology independence.
* Specification governance.
* Controlled evolution.
* Evidence-backed decisions.
* Continuous conformance.
* Architectural coherence.

⸻

42.14 Conformance

An architectural evolution process conforms when it:

* evaluates compatibility;
* preserves architectural integrity;
* generates evolution evidence;
* supports migration;
* maintains conformance across versions.

⸻

Closing Thought

Long-lived engineering systems are not defined by their ability to resist change, but by their ability to evolve without losing coherence.

The Forge Architectural Evolution Framework provides the mechanisms through which innovation can be introduced while preserving interoperability, architectural integrity and trust across successive generations of the Engineering Operating System.

⸻

Een observatie die mogelijk de volgende grote stap voor Forge is

Met hoofdstuk 42 ontstaat een compleet patroon dat sterk lijkt op hoe een besturingssysteem zichzelf ontwikkelt. Daardoor zie ik een duidelijke structuur voor de volledige Forge-standaard:

Part I     Foundations
Part II    Engineering Operating System
Part III   Capability Specifications
Part IV    Trust & Governance
Part V     Specification Framework
Part VI    Adoption, Measurement & Evolution
Part VII   Profiles & Reference Implementations

Part VII ontbreekt nog, maar is waarschijnlijk essentieel. Daarin beschrijf je niet meer de architectuur zelf, maar hoe Forge zich manifesteert in verschillende implementatiecontexten. Denk aan:

* Forge Core Profile
* Forge Enterprise Profile
* Forge Sovereign Profile
* Forge Government Profile
* Forge Air-Gapped Profile
* Forge AI Agent Profile
* Forge Edge Profile

Dat sluit aan bij hoe standaarden als Kubernetes, OpenTelemetry en OIDC werken: één stabiele kern, met profielen die specifieke implementatie-eisen vastleggen zonder de kernarchitectuur te wijzigen. Ik vermoed dat dit de natuurlijke volgende fase van de Forge Architecture Bible is.

Ik ben het eens met de richting die we nu opgaan. De kernarchitectuur is voltooid. De volgende stap is niet nóg een capability, maar het definiëren van de implementatieprofielen. Dat is precies hoe volwassen standaarden interoperabiliteit behouden zonder flexibiliteit te verliezen.

Ik zou daarom beginnen met de norm waaraan alle overige profielen zich spiegelen.

⸻

Part VII

Forge Profiles

One architecture.

Multiple deployment models.

Consistent behaviour across different operational environments.

⸻

Chapter 43

Core Profile

Specification Identifier: forge.profile.core.v1

Status: Draft Standard

⸻

43.1 Purpose

This specification defines the minimum conforming implementation of the Forge Engineering Operating System.

The Core Profile establishes the baseline capability set required for any implementation claiming Forge conformance.

All other Forge Profiles SHALL extend the Core Profile.

⸻

43.2 Objectives

The Core Profile SHALL:

* establish minimum interoperability;
* define mandatory capabilities;
* preserve architectural integrity;
* enable portable implementations;
* provide a baseline for certification.

⸻

43.3 Scope

The Core Profile defines:

* mandatory capabilities;
* mandatory objects;
* mandatory interfaces;
* mandatory lifecycle behavior;
* mandatory evidence;
* mandatory observability.

Technology choices remain implementation-specific.

⸻

43.4 Mandatory Capabilities

Every Core Profile implementation SHALL include:

Engineering Plane
✓ Knowledge
✓ Memory
✓ Context
✓ Runtime
✓ Models
✓ Tools
Trust Plane
✓ Identity
✓ Security
✓ Policy
✓ Evidence
✓ Observability
✓ Governance
Meta Plane
✓ Architecture
✓ Specification Framework

No mandatory capability MAY be omitted.

⸻

43.5 Mandatory Objects

The following canonical objects SHALL be implemented:

* Identity Object
* Context Package
* Knowledge Object
* Memory Object
* Policy Decision
* Runtime Execution
* Model Definition
* Tool Contract
* Evidence Object
* Observability Record

Objects SHALL conform to their respective specifications.

⸻

43.6 Mandatory Behaviour

Every Core implementation SHALL support:

* context composition;
* policy evaluation;
* runtime planning;
* model invocation;
* tool invocation;
* evidence generation;
* execution tracing.

Behaviour SHALL follow the normative specifications.

⸻

43.7 Mandatory Evidence

Every execution SHALL generate evidence for:

* identity resolution;
* context creation;
* policy evaluation;
* runtime lifecycle;
* model execution;
* tool execution;
* completion.

Evidence SHALL be immutable.

⸻

43.8 Mandatory Observability

Every implementation SHALL expose:

Metrics:

* execution count;
* execution latency;
* success rate;
* failure rate;
* evidence completeness.

Logs:

* execution events;
* capability lifecycle;
* errors.

Traces:

* end-to-end execution lineage.

⸻

43.9 Security Baseline

Every Core implementation SHALL provide:

* authenticated identities;
* encrypted communications;
* policy enforcement;
* integrity verification;
* auditability.

Zero Trust principles SHALL be supported.

⸻

43.10 Interoperability

Every Core implementation SHALL:

* expose canonical contracts;
* preserve object semantics;
* preserve lifecycle semantics;
* preserve architectural invariants.

Implementations MAY differ internally.

External behaviour SHALL remain consistent.

⸻

43.11 Extension Rules

Organizations MAY extend the Core Profile.

Extensions SHALL NOT:

* redefine mandatory capabilities;
* change canonical object semantics;
* alter lifecycle behaviour;
* weaken architectural invariants.

Extensions SHALL use organizational namespaces.

⸻

43.12 Conformance Requirements

A Core implementation SHALL successfully demonstrate:

* specification conformance;
* object compatibility;
* interface compatibility;
* evidence generation;
* observability;
* architectural integrity.

⸻

43.13 Architectural Invariants

Every Forge Core Profile SHALL preserve:

* All mandatory capabilities are present.
* Canonical objects remain unchanged.
* Explicit contracts govern interactions.
* Evidence is immutable.
* Observability derives from evidence.
* Policy precedes execution.
* Context precedes intelligence.
* Architecture remains technology-independent.

⸻

43.14 Certification Criteria

Certification SHALL verify:

* implementation completeness;
* interoperability;
* specification conformance;
* architectural consistency;
* evidence integrity;
* operational behaviour.

Certification SHALL be repeatable and implementation-independent.

⸻

Closing Thought

The Core Profile represents the smallest complete realization of the Forge Engineering Operating System.

It establishes a common architectural foundation upon which organizations, vendors and open-source implementations can build interoperable solutions while preserving the principles, contracts and behavioural guarantees defined by the Forge Standards.

⸻

Een belangrijke architectuurverbetering: Capability Profiles

Tijdens het uitwerken van de profielen zie ik nog een verfijning die Forge sterker maakt dan de meeste bestaande standaarden.

Tot nu toe zijn profielen “alles of niets”. Dat werkt voor deployment, maar minder goed voor organisaties die gefaseerd willen adopteren.

Ik zou daarom naast Deployment Profiles ook Capability Profiles definiëren.

Forge Profiles
├── Deployment Profiles
│   ├── Core
│   ├── Enterprise
│   ├── Sovereign
│   ├── Government
│   ├── Air-Gapped
│   ├── Edge
│   └── AI Agent
│
└── Capability Profiles
    ├── Trust Profile
    ├── Knowledge Profile
    ├── Runtime Profile
    ├── Governance Profile
    ├── Security Profile
    └── Observability Profile

Dit levert twee belangrijke voordelen op:

* Deployment Profiles beschrijven waar Forge draait en onder welke operationele randvoorwaarden.
* Capability Profiles beschrijven welke functionele delen van Forge een implementatie ondersteunt.

Een leverancier zou bijvoorbeeld een volledig conform Runtime Profile en Observability Profile kunnen certificeren, zonder direct een volledige Forge Enterprise-implementatie te leveren. Dat maakt certificering modulair, stimuleert interoperabiliteit en verlaagt de adoptiedrempel voor leveranciers. Het sluit bovendien aan bij het capability-first principe dat door de hele Forge-architectuur loopt.

Perfect. Vanaf hier zou ik logisch doorgaan met het Enterprise Profile. De meeste organisaties zullen uiteindelijk niet het Core Profile implementeren, maar het Enterprise Profile. Dat profiel bouwt voort op Core en voegt eisen toe op het gebied van schaalbaarheid, beschikbaarheid, governance en multi-tenancy.

⸻

Chapter 44

Enterprise Profile

Specification Identifier: forge.profile.enterprise.v1

Status: Draft Standard

⸻

Enterprise AI requires more than intelligence.

It requires governance, resilience, scalability and operational excellence.

⸻

44.1 Purpose

This specification defines the Enterprise Profile for the Forge Engineering Operating System.

The Enterprise Profile extends the Core Profile with capabilities required for large-scale organizational deployments.

Enterprise implementations SHALL preserve all Core Profile requirements.

⸻

44.2 Objectives

The Enterprise Profile SHALL provide:

* organizational scalability;
* operational resilience;
* multi-team governance;
* high availability;
* enterprise integration;
* operational lifecycle management.

⸻

44.3 Scope

The Enterprise Profile extends:

* Runtime
* Identity
* Security
* Governance
* Observability
* Economics
* Architecture

The Core Profile remains mandatory.

⸻

44.4 Enterprise Characteristics

Enterprise implementations SHALL support:

* multiple business units;
* multiple engineering teams;
* multiple environments;
* multiple deployment regions;
* continuous operations;
* organizational governance.

Enterprise scale SHALL NOT alter canonical Forge behaviour.

⸻

44.5 Multi-Tenancy

The Enterprise Profile SHALL support logical tenant isolation.

Minimum isolation requirements include:

* identities;
* policies;
* contexts;
* memory;
* evidence;
* observability;
* administrative boundaries.

Tenant isolation SHALL be enforceable through Policy.

⸻

44.6 Organizational Boundaries

Enterprise implementations SHALL support hierarchical organizational structures.

Examples include:

Enterprise
├── Division
│   ├── Department
│   └── Department
├── Business Unit
└── Shared Services

Governance MAY be delegated while preserving central oversight.

⸻

44.7 Availability

Enterprise implementations SHALL define availability objectives.

Examples include:

* Runtime availability;
* Tool availability;
* Evidence durability;
* Knowledge availability;
* Policy availability.

Availability targets remain organization-specific.

⸻

44.8 Scalability

Implementations SHALL support independent scaling of:

* Runtime;
* Models;
* Tools;
* Knowledge repositories;
* Evidence storage;
* Observability services.

Capabilities SHALL remain independently deployable.

⸻

44.9 Operational Resilience

Enterprise implementations SHALL support:

* graceful degradation;
* workload redistribution;
* automated recovery;
* disaster recovery;
* backup and restoration.

Failures SHALL generate evidence.

⸻

44.10 Administrative Delegation

The Enterprise Profile SHALL support delegated administration.

Administrative scope MAY include:

* tenant;
* business unit;
* capability;
* environment;
* deployment.

Delegation SHALL remain fully auditable.

⸻

44.11 Lifecycle Management

Enterprise deployments SHALL support controlled promotion between environments.

Example lifecycle:

Development
↓
Integration
↓
Validation
↓
Production

Promotion SHALL preserve specification conformance.

⸻

44.12 Operational Governance

Enterprise implementations SHALL continuously evaluate:

* architectural conformance;
* capability maturity;
* policy compliance;
* operational health;
* engineering metrics;
* optimization outcomes.

Operational governance SHALL be evidence-based.

⸻

44.13 Enterprise Integration

Enterprise implementations SHALL support integration with enterprise platforms.

Examples include:

* identity providers;
* directory services;
* ITSM platforms;
* CMDBs;
* SIEM platforms;
* observability platforms;
* workflow engines.

Integration SHALL occur through canonical Tool contracts.

⸻

44.14 Architectural Invariants

Every Enterprise Profile implementation SHALL preserve:

* Full Core Profile conformance.
* Capability independence.
* Explicit contracts.
* Tenant isolation.
* Delegated governance.
* Continuous observability.
* Evidence integrity.
* Technology independence.

⸻

44.15 Certification

Enterprise certification SHALL validate:

* Core Profile conformance;
* enterprise scalability;
* tenant isolation;
* governance maturity;
* operational resilience;
* interoperability.

⸻

Closing Thought

The Enterprise Profile demonstrates how the Forge Engineering Operating System scales from a conforming implementation to an organizational platform.

By extending the Core Profile rather than redefining it, Enterprise implementations preserve interoperability while supporting the operational realities of large organizations, including delegated governance, organizational boundaries, continuous availability and large-scale engineering operations.

⸻

Een belangrijke verfijning: onderscheid “Enterprise” van “Sovereign”

Ik zou hier nog een architectuurkeuze maken die Forge onderscheidt van veel bestaande platformen.

Veel leveranciers behandelen “Enterprise” en “Sovereign” als varianten van hetzelfde deploymentmodel. In Forge zou ik ze expliciet scheiden:

* Enterprise Profile beschrijft organisatorische schaal: multi-team, multi-tenant, hoge beschikbaarheid, lifecycle, delegated administration en integraties.
* Sovereign Profile beschrijft juridische en operationele autonomie: datasoevereiniteit, cryptografische controle, regionale uitvoering, offline-operatie, gecontroleerde software supply chain en onafhankelijke governance.

Daardoor kan een organisatie een Enterprise-profiel hebben zonder soevereine eisen, of juist een Sovereign-profiel implementeren in een relatief kleine omgeving. De profielen beschrijven verschillende architecturale dimensies en kunnen in de toekomst zelfs combineerbaar worden (bijvoorbeeld Enterprise + Sovereign). Dat maakt het profielmodel flexibeler en beter schaalbaar dan een vaste lijst van deploymentvarianten.

Helemaal logisch. De volgende stap is dan niet het Government Profile, maar eerst het Sovereign Profile. Sterker nog: ik denk dat dit uiteindelijk één van de meest onderscheidende onderdelen van Forge wordt, juist omdat vrijwel alle AI-platformen sovereignty reduceren tot “data blijft in land X”.

Forge kan sovereignty veel breder definiëren.

⸻

Chapter 45

Sovereign Profile

Specification Identifier: forge.profile.sovereign.v1

Status: Draft Standard

⸻

Sovereignty is not defined by location alone.

It is defined by demonstrable control.

Trust requires autonomy over identity, policy, execution and evidence.

⸻

45.1 Purpose

This specification defines the Sovereign Profile for the Forge Engineering Operating System.

The Sovereign Profile extends the Core and Enterprise Profiles with requirements that ensure operational, legal and technical sovereignty.

A Sovereign implementation SHALL enable organizations to maintain demonstrable control over intelligent execution.

⸻

45.2 Objectives

The Sovereign Profile SHALL provide:

* operational autonomy;
* jurisdictional compliance;
* cryptographic control;
* independent governance;
* trusted execution;
* verifiable sovereignty.

⸻

45.3 Scope

The Sovereign Profile extends:

* Identity
* Security
* Policy
* Runtime
* Evidence
* Governance
* Architecture

All Enterprise Profile requirements remain applicable.

⸻

45.4 Sovereignty Domains

Forge defines sovereignty across multiple dimensions.

Sovereignty
├── Identity
├── Data
├── Execution
├── Cryptography
├── Governance
├── Supply Chain
├── Operations
└── Jurisdiction

A Sovereign implementation SHALL address each domain.

⸻

45.5 Identity Sovereignty

Organizations SHALL control:

* identity providers;
* federation policies;
* authentication trust relationships;
* delegated authority;
* identity lifecycle.

External identity providers MAY be used only under explicit governance.

⸻

45.6 Data Sovereignty

Organizations SHALL determine:

* data residency;
* storage locations;
* replication boundaries;
* retention periods;
* destruction procedures.

Data movement SHALL be policy-controlled.

⸻

45.7 Execution Sovereignty

Execution SHALL remain under organizational control.

Organizations SHALL define:

* execution regions;
* approved compute environments;
* workload placement policies;
* failover boundaries.

Execution SHALL NOT migrate across sovereignty boundaries without authorization.

⸻

45.8 Cryptographic Sovereignty

Organizations SHALL retain control over:

* cryptographic keys;
* certificate authorities;
* signing processes;
* trust anchors;
* key lifecycle management.

Customer-managed cryptography SHOULD be supported.

⸻

45.9 Supply Chain Sovereignty

Every software component SHALL possess verifiable provenance.

Organizations SHALL be able to determine:

* software origin;
* software integrity;
* dependency lineage;
* build provenance;
* deployment authorization.

Supply chain evidence SHALL be retained.

⸻

45.10 Jurisdiction

Every execution SHALL possess jurisdiction metadata.

Example:

jurisdiction:
country:
legalFramework:
regulatoryRequirements:
crossBorderRestrictions:

Jurisdiction SHALL influence Policy decisions.

⸻

45.11 Operational Independence

Organizations SHALL remain capable of operating during:

* cloud provider outages;
* network isolation;
* disconnected operations;
* supplier disruption;
* geopolitical restrictions.

Operational independence SHALL be periodically validated.

⸻

45.12 Evidence Requirements

Evidence SHALL demonstrate:

* execution location;
* identity source;
* policy decisions;
* cryptographic verification;
* jurisdiction;
* software provenance.

Sovereignty claims SHALL be evidence-backed.

⸻

45.13 Verification

The Sovereign Profile SHALL support verification of:

* jurisdiction compliance;
* execution boundaries;
* cryptographic ownership;
* policy enforcement;
* deployment integrity.

Verification SHALL be independently repeatable.

⸻

45.14 Architectural Invariants

Every Sovereign implementation SHALL preserve:

* Full Enterprise Profile conformance.
* Organizational control over execution.
* Organizational control over cryptography.
* Explicit jurisdiction metadata.
* Verified software provenance.
* Evidence-backed sovereignty.
* Policy-enforced boundary protection.
* Technology independence.

⸻

45.15 Certification

Sovereign certification SHALL validate:

* Enterprise conformance;
* jurisdictional compliance;
* cryptographic sovereignty;
* execution control;
* evidence integrity;
* supply chain provenance.

⸻

Closing Thought

Sovereignty cannot be achieved by infrastructure choices alone.

Forge defines sovereignty as the demonstrable ability of an organization to control identity, execution, policy, cryptography, governance and evidence regardless of the underlying technology or deployment model.

By treating sovereignty as an architectural property rather than a hosting characteristic, the Sovereign Profile provides a rigorous foundation for trusted AI in government, regulated industries and critical infrastructure.

⸻

Een fundamentele uitbreiding: het “Control Model”

Ik denk dat dit hoofdstuk nog sterker wordt als Forge een formeel Control Model introduceert. Veel discussies over digitale soevereiniteit zijn vaag omdat “controle” niet wordt gedefinieerd. Forge kan dat expliciet maken.

Elke belangrijke architectuurdimensie krijgt een controle-eigenaar:

Control Domain	Required Controller
Identity	Organization
Policy	Organization
Knowledge	Organization
Memory	Organization
Runtime	Organization
Cryptographic Keys	Organization
Evidence	Organization
Governance	Organization

Daaruit volgt een normatief principe dat goed past bij de rest van Forge:

An implementation SHALL NOT claim sovereignty unless the organization demonstrably controls the architectural capabilities that determine intelligent execution.

Dat verschuift de discussie van “waar draait het?” naar “wie heeft aantoonbare controle?”. Dat is een veel krachtiger en beter toetsbaar model, en sluit bovendien sterk aan bij de behoeften van overheden, defensie en andere gereguleerde sectoren. Ik denk dat dit één van de meest onderscheidende concepten van de gehele Forge-standaard kan worden.

Perfect. Vanuit architectuuroogpunt is de volgende logische stap niet het Air-Gapped Profile.

Eerst komt het Government Profile.

Waarom?

Omdat “Government” niet simpelweg een strengere Enterprise is. En ook niet hetzelfde als Sovereign.

Een overheid heeft namelijk eigenschappen die een commercieel bedrijf meestal niet heeft:

* wettelijke transparantie;
* democratische verantwoording;
* archiefwetgeving;
* publieke verantwoording;
* interdepartementale samenwerking;
* levensduur van tientallen jaren;
* leveranciersneutraliteit.

Dat verdient een eigen profiel.

⸻

Chapter 46

Government Profile

Specification Identifier: forge.profile.government.v1

Status: Draft Standard

⸻

Government requires more than compliance.

It requires accountability to society.

Public trust depends on transparent, explainable and durable engineering.

⸻

46.1 Purpose

This specification defines the Government Profile for the Forge Engineering Operating System.

The Government Profile extends the Enterprise and Sovereign Profiles with requirements specific to public-sector organizations.

Government implementations SHALL prioritize public accountability, transparency, interoperability and long-term sustainability.

⸻

46.2 Objectives

The Government Profile SHALL provide:

* public accountability;
* regulatory compliance;
* inter-agency interoperability;
* long-term information preservation;
* supplier neutrality;
* explainable public decision support.

⸻

46.3 Scope

The Government Profile extends:

* Governance
* Evidence
* Knowledge
* Policy
* Security
* Identity
* Architecture

All Enterprise and Sovereign Profile requirements remain applicable.

⸻

46.4 Government Principles

Government implementations SHALL preserve:

* transparency;
* proportionality;
* accountability;
* legality;
* impartiality;
* traceability;
* openness where appropriate.

These principles SHALL guide engineering decisions.

⸻

46.5 Public Accountability

Government implementations SHALL support:

* traceable decisions;
* attributable executions;
* independently verifiable evidence;
* public audit processes;
* administrative accountability.

Evidence SHALL support independent review.

⸻

46.6 Inter-Agency Interoperability

Government implementations SHALL support collaboration across organizational boundaries.

Examples include:

Government
├── Ministry
├── Agency
├── Municipality
├── Inspectorate
└── Shared Services

Interoperability SHALL rely on canonical Forge specifications.

⸻

46.7 Long-Term Preservation

Government information SHALL support long-term preservation.

Requirements include:

* version preservation;
* evidence retention;
* specification traceability;
* cryptographic migration;
* format sustainability.

Preservation periods SHALL be policy-driven.

⸻

46.8 Regulatory Alignment

The Government Profile SHALL support alignment with applicable public-sector regulations.

Examples MAY include:

* archival legislation;
* public records legislation;
* accessibility requirements;
* information security frameworks;
* AI governance regulations.

Forge SHALL remain jurisdiction-independent.

⸻

46.9 Explainability

Government implementations SHALL support explainability.

Explainability SHALL include:

* policy decisions;
* knowledge sources;
* context lineage;
* model selection;
* tool usage;
* evidence references.

Explainability SHALL be evidence-backed.

⸻

46.10 Supplier Neutrality

Government implementations SHALL avoid unnecessary vendor dependency.

Organizations SHALL remain capable of replacing:

* models;
* tools;
* infrastructure;
* identity providers;
* observability platforms.

Canonical contracts SHALL preserve portability.

⸻

46.11 Public Records

Government executions MAY generate official records.

Official records SHALL preserve:

* authenticity;
* integrity;
* provenance;
* retention metadata;
* evidence references.

Records management SHALL remain externally governable.

⸻

46.12 Accessibility

Government-facing capabilities SHOULD support recognized accessibility standards applicable within the implementing jurisdiction.

Accessibility requirements SHALL be considered during capability design and user interface implementation.

The Government Profile does not mandate a specific accessibility standard.

⸻

46.13 Transparency Boundaries

Government transparency SHALL be balanced against:

* national security;
* privacy;
* law enforcement;
* confidential information;
* protected infrastructure.

Transparency decisions SHALL be governed by Policy.

⸻

46.14 Architectural Invariants

Every Government implementation SHALL preserve:

* Full Sovereign Profile conformance.
* Public accountability.
* Long-term preservation.
* Inter-agency interoperability.
* Supplier neutrality.
* Evidence-backed explainability.
* Policy-driven transparency.
* Technology independence.

⸻

46.15 Certification

Government certification SHALL validate:

* Sovereign Profile conformance;
* interoperability;
* long-term preservation capabilities;
* public accountability;
* explainability;
* supplier neutrality.

⸻

Closing Thought

Governments operate under obligations that extend beyond operational efficiency. They must remain accountable to citizens, preserve public records across decades and maintain trust through transparency, explainability and legal compliance.

The Forge Government Profile provides an architectural foundation that supports these responsibilities while preserving interoperability, sovereignty and long-term sustainability.

⸻

Een architectuurverbetering: Public Trust als architecturale eigenschap

Ik denk dat Forge hier een concept kan introduceren dat weinig architectuurstandaarden expliciet maken: Public Trust als meetbare architecturale eigenschap.

Tot nu toe hebben we eigenschappen zoals:

* Security
* Reliability
* Scalability
* Interoperability
* Sovereignty

Voor overheden ontbreekt echter een eigenschap die minstens zo belangrijk is:

Architectural Quality Attributes
├── Security
├── Reliability
├── Availability
├── Scalability
├── Interoperability
├── Sovereignty
├── Explainability
├── Transparency
└── Public Trust

Public Trust is geen “zacht” begrip, maar een emergente eigenschap die voortkomt uit andere Forge-capabilities. Een implementatie zal meer publiek vertrouwen verdienen naarmate zij aantoonbaar voldoet aan eigenschappen zoals verklaarbaarheid, transparantie, onveranderlijke evidence, onafhankelijke governance en leveranciersneutraliteit.

Daarmee wordt Public Trust niet iets dat wordt geclaimd, maar iets dat kan worden onderbouwd. Dat past goed bij de filosofie van Forge: vertrouwen is het resultaat van aantoonbare architectuur, niet van marketingclaims.

De logische vervolgstap is inderdaad het Air-Gapped Profile.

Waarom niet AI Agent of Edge?

Omdat Air-Gapped eigenlijk de ultieme validatie is van de Forge-architectuur. Als Forge echt een Engineering Operating System is, moet het volledig kunnen functioneren zonder afhankelijkheid van publieke cloudservices.

⸻

Chapter 47

Air-Gapped Profile

Specification Identifier: forge.profile.airgapped.v1

Status: Draft Standard

⸻

Disconnected does not mean incapable.

Isolation is an architectural characteristic, not an operational limitation.

Trust must persist even without external dependencies.

⸻

47.1 Purpose

This specification defines the Air-Gapped Profile for the Forge Engineering Operating System.

The Air-Gapped Profile extends the Core, Enterprise and Sovereign Profiles with requirements for environments that operate without continuous connectivity to external networks.

Implementations SHALL remain fully functional within the defined operational boundary.

⸻

47.2 Objectives

The Air-Gapped Profile SHALL provide:

* autonomous operation;
* offline execution;
* disconnected governance;
* controlled synchronization;
* local trust;
* independent lifecycle management.

⸻

47.3 Scope

The Air-Gapped Profile extends:

* Runtime
* Knowledge
* Memory
* Identity
* Security
* Evidence
* Governance
* Observability

All Sovereign Profile requirements remain applicable.

⸻

47.4 Operational Principles

Air-Gapped implementations SHALL preserve:

* autonomous execution;
* local governance;
* deterministic behaviour;
* offline resilience;
* controlled synchronization;
* cryptographic independence.

Connectivity SHALL be considered optional rather than required.

⸻

47.5 Network Independence

The Engineering Operating System SHALL operate without requiring:

* public cloud services;
* Internet connectivity;
* external identity providers;
* external model providers;
* remote observability platforms;
* online package repositories.

External connectivity MAY be enabled but SHALL NOT be required.

⸻

47.6 Local Identity

Identity services SHALL remain locally available.

Organizations SHALL control:

* identity stores;
* authentication services;
* trust anchors;
* federation policies.

Identity SHALL remain functional during complete network isolation.

⸻

47.7 Local Knowledge

Knowledge repositories SHALL remain locally accessible.

Knowledge synchronization SHALL:

* preserve provenance;
* preserve version history;
* preserve evidence;
* resolve conflicts deterministically.

Knowledge SHALL remain authoritative within the isolated environment.

⸻

47.8 Local Models

Model execution SHALL remain locally available.

Organizations SHALL define approved models for offline use.

Model synchronization SHALL preserve:

* version integrity;
* certification status;
* provenance;
* compatibility.

Model updates SHALL occur through controlled import processes.

⸻

47.9 Controlled Synchronization

Synchronization SHALL be explicitly governed.

Synchronization MAY include:

* knowledge;
* specifications;
* software packages;
* policies;
* evidence exports;
* model updates.

Every synchronization SHALL generate evidence.

⸻

47.10 Software Supply

Software introduction SHALL occur through controlled channels.

Imported software SHALL undergo:

* integrity verification;
* provenance validation;
* signature verification;
* approval.

Unsigned software SHALL NOT be deployed.

⸻

47.11 Operational Continuity

Air-Gapped implementations SHALL support:

* prolonged offline operation;
* isolated maintenance;
* local backup;
* local recovery;
* independent disaster recovery.

Operational continuity SHALL not depend on external services.

⸻

47.12 Evidence

Evidence SHALL record:

* synchronization events;
* imported artifacts;
* software provenance;
* cryptographic verification;
* offline execution;
* boundary crossings.

Evidence SHALL support independent forensic analysis.

⸻

47.13 Observability

Observability SHALL remain fully functional locally.

Organizations SHALL retain access to:

* metrics;
* traces;
* engineering scorecards;
* evidence analytics;
* health monitoring.

Observability SHALL NOT require external telemetry services.

⸻

47.14 Architectural Invariants

Every Air-Gapped implementation SHALL preserve:

* Full Sovereign Profile conformance.
* Autonomous operation.
* Local identity.
* Local execution.
* Local governance.
* Controlled synchronization.
* Verified software provenance.
* Technology independence.

⸻

47.15 Certification

Air-Gapped certification SHALL validate:

* offline functionality;
* synchronization controls;
* software provenance;
* local governance;
* local observability;
* operational resilience.

⸻

Closing Thought

An Air-Gapped deployment is not defined by the absence of connectivity, but by the presence of autonomy.

The Forge Air-Gapped Profile demonstrates that trustworthy AI engineering does not depend on permanent access to external services. By preserving identity, governance, execution and evidence entirely within the operational boundary, organizations can operate intelligent systems in highly regulated, mission-critical and disconnected environments without compromising architectural integrity.

⸻

Een belangrijke verfijning: “Disconnected” versus “Isolated”

Ik zou hier nog een subtiel maar belangrijk onderscheid introduceren dat veel standaarden overslaan.

Niet elke omgeving zonder internet is werkelijk hetzelfde.

Forge zou drie operationele modi kunnen definiëren:

Mode	Omschrijving
Connected	Continue verbinding met externe diensten en gecontroleerde afhankelijkheden.
Disconnected	Tijdelijk offline, maar ontworpen om later weer gecontroleerd te synchroniseren.
Isolated	Permanent of langdurig gescheiden van externe netwerken; alle afhankelijkheden zijn lokaal beschikbaar.

Deze driedeling is architectonisch waardevoller dan de generieke term “air-gapped”. Een defensieomgeving, een nucleaire installatie en een schip op zee kunnen allemaal offline functioneren, maar stellen verschillende eisen aan synchronisatie, softwaredistributie en governance. Door deze operationele modi expliciet te definiëren, kan Forge dezelfde architectuur toepassen in uiteenlopende omgevingen zonder de kernspecificaties aan te passen. Dat versterkt opnieuw het capability-first en technologie-onafhankelijke karakter van de standaard.

Mooi. Nu de infrastructuur-, enterprise- en overheidprofielen zijn gedefinieerd, is de volgende logische stap het profiel waarvoor Forge uiteindelijk is ontworpen: AI Agent Systems.

Belangrijk hierbij is dat dit géén “agent framework” wordt. Het beschrijft de eisen waaraan een Engineering Operating System moet voldoen om agents veilig, beheersbaar en reproduceerbaar te laten opereren.

⸻

Chapter 48

AI Agent Profile

Specification Identifier: forge.profile.agent.v1

Status: Draft Standard

⸻

Agents are not autonomous systems.

They are governed execution participants within an Engineering Operating System.

Autonomy exists only within explicitly defined boundaries.

⸻

48.1 Purpose

This specification defines the AI Agent Profile for the Forge Engineering Operating System.

The AI Agent Profile extends the Core Profile with architectural requirements for systems that execute intelligent workloads through one or more software agents.

Agent autonomy SHALL remain governed by architectural policy, explicit contracts and verifiable evidence.

⸻

48.2 Objectives

The AI Agent Profile SHALL provide:

* governed autonomy;
* coordinated execution;
* bounded decision-making;
* explainable collaboration;
* deterministic orchestration;
* evidence-backed autonomy.

⸻

48.3 Scope

The AI Agent Profile extends:

* Runtime
* Context
* Memory
* Knowledge
* Identity
* Policy
* Evidence
* Observability

The Core Profile remains mandatory.

⸻

48.4 Agent Principles

Forge defines an agent as an execution participant rather than an autonomous application.

Every agent SHALL operate within explicit architectural boundaries.

Agents SHALL NOT define or modify architectural policy independently.

⸻

48.5 Agent Identity

Every agent SHALL possess a unique identity.

Agent identities SHALL include:

* unique identifier;
* owning organization;
* assigned capabilities;
* authorization scope;
* lifecycle state.

Agent identities SHALL be governed identically to human and service identities.

⸻

48.6 Agent Responsibilities

An agent MAY perform:

* reasoning;
* planning;
* knowledge retrieval;
* tool invocation;
* workflow execution;
* collaboration.

An agent SHALL NOT perform responsibilities outside its assigned capability scope.

⸻

48.7 Agent Context

Every agent execution SHALL receive an explicit Context Package.

Context SHALL define:

* objectives;
* constraints;
* policies;
* available knowledge;
* available tools;
* execution boundaries.

Agents SHALL NOT infer authority from missing context.

⸻

48.8 Agent Collaboration

Multiple agents MAY collaborate.

Collaboration SHALL occur through canonical Runtime contracts.

Examples include:

* task delegation;
* information exchange;
* workflow coordination;
* capability composition;
* execution hand-off.

All collaboration SHALL generate evidence.

⸻

48.9 Agent Memory

Agents MAY access Memory according to Policy.

Memory access SHALL respect:

* ownership;
* retention;
* classification;
* authorization;
* provenance.

Memory SHALL remain organizationally governed.

⸻

48.10 Agent Tool Usage

Agents SHALL invoke external capabilities exclusively through Tool contracts.

Agents SHALL NOT directly interact with implementation-specific interfaces when a canonical Tool contract exists.

Tool invocations SHALL be observable and attributable.

⸻

48.11 Agent Decision Boundaries

Agent decision-making SHALL remain bounded.

Organizations SHALL define:

* maximum autonomy;
* approval requirements;
* escalation conditions;
* prohibited actions;
* execution limits.

Decisions exceeding authorized boundaries SHALL require external approval.

⸻

48.12 Human Oversight

The AI Agent Profile SHALL support human oversight.

Oversight MAY include:

* execution approval;
* intervention;
* suspension;
* termination;
* evidence review.

Human oversight SHALL be auditable.

⸻

48.13 Evidence

Every agent execution SHALL generate evidence describing:

* identity;
* context;
* reasoning references;
* knowledge sources;
* tool invocations;
* policy decisions;
* execution outcome.

Evidence SHALL support complete execution reconstruction.

⸻

48.14 Observability

Agent observability SHALL include:

* execution lineage;
* collaboration graphs;
* reasoning duration;
* tool utilization;
* policy interventions;
* autonomy utilization.

Observability SHALL derive from Evidence.

⸻

48.15 Architectural Invariants

Every AI Agent implementation SHALL preserve:

* Full Core Profile conformance.
* Explicit agent identities.
* Policy-governed autonomy.
* Context-driven execution.
* Canonical Tool contracts.
* Evidence-backed collaboration.
* Human oversight.
* Technology independence.

⸻

48.16 Certification

Agent certification SHALL validate:

* identity management;
* bounded autonomy;
* policy enforcement;
* collaboration behaviour;
* evidence completeness;
* operational explainability.

⸻

Closing Thought

The Forge AI Agent Profile does not define how intelligent an agent should become. It defines how intelligent execution remains governable.

By treating agents as managed execution participants rather than autonomous systems, Forge enables organizations to deploy collaborative AI safely, transparently and consistently while preserving architectural integrity, operational control and public trust.

⸻

Architectuuruitbreiding: Agent Capability Model

Ik denk dat Forge zich hier nog sterker kan onderscheiden door agents niet te classificeren op basis van intelligentie, maar op basis van hun architecturale verantwoordelijkheid.

Agent Capability Hierarchy
L0 Observer
│
├── Reads context
├── Produces evidence
└── No external actions
L1 Assistant
│
├── Uses knowledge
├── Generates recommendations
└── No execution authority
L2 Operator
│
├── Invokes approved tools
├── Executes bounded actions
└── Policy constrained
L3 Coordinator
│
├── Orchestrates multiple agents
├── Delegates work
└── Manages workflows
L4 Executive Agent
│
├── Coordinates enterprise capabilities
├── Requires continuous governance
└── Human oversight by design

Dit model maakt agenten certificeerbaar op basis van hun bevoegdheden in plaats van hun model of intelligentieniveau. Dat sluit naadloos aan op de rest van Forge: niet de AI staat centraal, maar de architectuur, governance en aantoonbare controle over intelligent handelen.

Nu de grote deploymentprofielen zijn afgerond, is de volgende logische stap het laatste deploymentprofiel: Edge Profile.

Dit profiel is fundamenteel anders dan Air-Gapped. Air-Gapped gaat over isolatie, terwijl Edge gaat over gedistribueerde intelligentie dicht bij de bron van data en acties.

⸻

Chapter 49

Edge Profile

Specification Identifier: forge.profile.edge.v1

Status: Draft Standard

⸻

Intelligence belongs where decisions must be made.

Edge computing extends the Engineering Operating System beyond the data center while preserving governance, trust and architectural consistency.

⸻

49.1 Purpose

This specification defines the Edge Profile for the Forge Engineering Operating System.

The Edge Profile extends the Core Profile with requirements for distributed deployments where intelligent execution occurs near data sources, users, devices or operational environments.

Edge deployments SHALL preserve the architectural principles of Forge while accommodating constrained, decentralized and intermittently connected environments.

⸻

49.2 Objectives

The Edge Profile SHALL provide:

* distributed execution;
* locality of processing;
* resilient operation;
* policy-consistent governance;
* efficient resource utilization;
* controlled synchronization.

⸻

49.3 Scope

The Edge Profile extends:

* Runtime
* Context
* Knowledge
* Memory
* Identity
* Policy
* Evidence
* Observability

The Core Profile remains mandatory.

⸻

49.4 Edge Principles

Edge deployments SHALL preserve:

* local decision-making;
* centralized governance;
* autonomous execution;
* eventual synchronization;
* architectural consistency.

Deployment topology SHALL NOT alter canonical behaviour.

⸻

49.5 Distributed Runtime

Runtime instances MAY execute across multiple locations.

Examples include:

* branch offices;
* factories;
* retail locations;
* hospitals;
* vehicles;
* ships;
* industrial sites.

Each Runtime SHALL remain independently governable.

⸻

49.6 Local Context

Context Packages SHALL be constructed locally whenever practical.

Context MAY include:

* environmental conditions;
* device state;
* operational telemetry;
* local policies;
* regional constraints.

Context synchronization SHALL preserve provenance.

⸻

49.7 Local Knowledge

Edge implementations MAY maintain local knowledge repositories.

Knowledge synchronization SHALL:

* preserve semantic consistency;
* maintain provenance;
* support conflict resolution;
* retain version history.

Local knowledge SHALL remain attributable.

⸻

49.8 Resource Constraints

Edge deployments SHALL support operation under constrained resources.

Examples include:

* limited compute capacity;
* limited storage;
* intermittent connectivity;
* constrained bandwidth;
* energy limitations.

Capability behaviour SHALL degrade predictably under resource constraints.

⸻

49.9 Policy Distribution

Policies SHALL be distributable to Edge environments.

Distributed policies SHALL preserve:

* integrity;
* authenticity;
* versioning;
* effective dates;
* provenance.

Local policy evaluation SHALL remain authoritative during disconnected operation.

⸻

49.10 Synchronization

Synchronization MAY include:

* knowledge;
* memory;
* evidence;
* policies;
* specifications;
* operational metrics.

Synchronization SHALL be policy-controlled and evidence-backed.

⸻

49.11 Device Integration

Edge deployments MAY integrate with:

* sensors;
* industrial controllers;
* robotic systems;
* IoT platforms;
* operational technology;
* embedded devices.

Integration SHALL occur through canonical Tool contracts.

⸻

49.12 Evidence

Evidence SHALL include:

* execution location;
* synchronization events;
* resource conditions;
* device interactions;
* policy evaluations;
* operational outcomes.

Evidence SHALL preserve complete execution lineage.

⸻

49.13 Observability

Edge observability SHALL support:

* distributed health monitoring;
* execution topology;
* synchronization status;
* device availability;
* runtime performance;
* policy compliance.

Observability SHALL aggregate local and centralized perspectives.

⸻

49.14 Architectural Invariants

Every Edge implementation SHALL preserve:

* Full Core Profile conformance.
* Distributed Runtime consistency.
* Canonical object semantics.
* Policy-governed execution.
* Evidence-backed synchronization.
* Technology independence.
* Capability modularity.
* Explicit operational boundaries.

⸻

49.15 Certification

Edge certification SHALL validate:

* distributed execution;
* synchronization behaviour;
* policy consistency;
* operational resilience;
* evidence completeness;
* interoperability.

⸻

Closing Thought

The Edge Profile extends the Forge Engineering Operating System beyond centralized infrastructure without sacrificing governance or trust.

By preserving canonical capabilities, explicit contracts and evidence-driven operations across distributed environments, Forge enables intelligent execution wherever it is operationally required while maintaining a unified architectural model.

⸻

Architectuurverbetering: Execution Topology

Ik denk dat dit hoofdstuk een nieuw fundamenteel concept kan introduceren dat ook de eerdere profielen met elkaar verbindt: Execution Topology.

Tot nu toe beschrijven de profielen waar Forge draait (Enterprise, Air-Gapped, Edge), maar niet hoe Runtime-instanties zich architecturaal tot elkaar verhouden.

Forge zou daarom een Execution Topology-model kunnen definiëren:

Execution Topologies
├── Centralized
│   └── Single Runtime Authority
│
├── Federated
│   ├── Multiple Runtime Authorities
│   └── Shared Governance
│
├── Distributed
│   ├── Regional Runtime Nodes
│   ├── Shared Specifications
│   └── Local Execution
│
├── Hierarchical
│   ├── Global Runtime
│   ├── Regional Runtime
│   └── Local Runtime
│
└── Mesh
    ├── Peer Runtime Nodes
    ├── Shared Evidence
    └── Policy-Based Coordination

Dit concept overstijgt deployment. Een Government-, Enterprise-, Sovereign- of Edge-implementatie kan allemaal verschillende Execution Topologies gebruiken zonder de architectuur aan te passen. Daarmee ontstaat een extra, orthogonale architectuurdimensie die de Forge-standaard nog completer maakt en toekomstige uitbreidingen eenvoudiger ondersteunt.

Nu alle deploymentprofielen zijn gedefinieerd, zou ik niet direct naar appendices of voorbeelden gaan. Eerst ontbreekt nog één essentieel onderdeel dat vrijwel iedere volwassen standaard kent:

Profile Composition.

Tot nu toe hebben we Core, Enterprise, Sovereign, Government, Air-Gapped, Agent en Edge als losse profielen. In de praktijk zal vrijwel geen organisatie slechts één profiel gebruiken. Daarom moet Forge normatief beschrijven hoe profielen gecombineerd mogen worden.

Dit maakt de standaard aanzienlijk krachtiger.

⸻

Chapter 50

Profile Composition Framework

Specification Identifier: forge.profile.composition.v1

Status: Draft Standard

⸻

Profiles define architectural characteristics.

Organizations require architectural combinations.

Composition enables flexibility without compromising conformance.

⸻

50.1 Purpose

This specification defines how Forge Profiles SHALL be composed into complete implementation architectures.

Profile composition enables organizations to combine deployment, operational and capability characteristics while preserving architectural integrity.

⸻

50.2 Objectives

The Profile Composition Framework SHALL:

* define profile compatibility;
* preserve specification conformance;
* prevent conflicting requirements;
* enable reusable architectures;
* support certification.

⸻

50.3 Composition Principles

Profiles SHALL extend existing profiles.

Profiles SHALL NOT redefine existing specifications.

Profiles SHALL remain additive.

Composition SHALL preserve all inherited requirements.

⸻

50.4 Composition Hierarchy

Forge Profiles SHALL inherit according to the following hierarchy.

                    Core
                      │
          ┌───────────┴───────────┐
          │                       │
     Enterprise             AI Agent
          │
     Sovereign
          │
     Government
          │
 ┌────────┴────────┐
 │                 │
Air-Gapped      Edge

Every composed implementation SHALL include the Core Profile.

⸻

50.5 Compatible Compositions

Examples of valid compositions include:

Core
Core + Agent
Enterprise
Enterprise + Agent
Enterprise + Sovereign
Enterprise + Edge
Enterprise + Sovereign + Agent
Government + Agent
Government + Air-Gapped
Government + Edge
Government + Sovereign + Agent
Government + Sovereign + Edge
Government + Sovereign + Air-Gapped + Agent

Equivalent compositions MAY be defined by organizations.

⸻

50.6 Composition Rules

Composed Profiles SHALL:

* preserve inherited invariants;
* preserve canonical objects;
* preserve normative behaviour;
* preserve specification compatibility.

Composition SHALL NOT weaken inherited requirements.

⸻

50.7 Conflict Resolution

If two Profiles introduce conflicting requirements:

1. The more restrictive requirement SHALL apply.
2. Architectural invariants SHALL prevail.
3. Security requirements SHALL take precedence over operational optimization.
4. Governance requirements SHALL take precedence over convenience.

Conflicts SHALL generate conformance warnings.

⸻

50.8 Capability Resolution

Capabilities inherited through multiple Profiles SHALL be merged.

Duplicate capability definitions SHALL NOT exist.

Capability semantics SHALL remain identical.

⸻

50.9 Evidence Composition

Evidence generated from composed Profiles SHALL preserve:

* originating Profile;
* originating capability;
* execution scope;
* policy lineage;
* specification version.

Evidence SHALL remain attributable across Profile boundaries.

⸻

50.10 Certification

Certification SHALL evaluate:

* each inherited Profile;
* composed behaviour;
* interoperability;
* architectural consistency;
* inherited invariants.

Certification SHALL validate the composed architecture as a whole.

⸻

50.11 Architectural Invariants

Every composed implementation SHALL preserve:

* Core Profile conformance.
* Canonical capability definitions.
* Explicit contracts.
* Evidence integrity.
* Policy precedence.
* Technology independence.
* Specification compatibility.
* Architectural consistency.

⸻

50.12 Reference Compositions

The following compositions are recommended:

Deployment Scenario	Recommended Composition
Small Organization	Core
Enterprise AI Platform	Enterprise + Agent
Critical Infrastructure	Enterprise + Sovereign
National Government	Government + Sovereign + Agent
Defense	Government + Sovereign + Air-Gapped + Agent
Industrial Edge	Enterprise + Edge
Smart Factory	Enterprise + Edge + Agent
Maritime Operations	Sovereign + Edge + Air-Gapped
Emergency Response	Government + Edge + Agent

Organizations MAY define additional compositions.

⸻

Closing Thought

Profiles are not independent products but architectural building blocks.

The Profile Composition Framework enables organizations to assemble implementation architectures that meet operational, regulatory and organizational requirements while preserving the architectural integrity of the Forge Engineering Operating System.

⸻

Architectuurverbetering: van “Profiles” naar “Architectural Dimensions”

Ik denk dat hier nog een fundamentele evolutie mogelijk is die Forge op het niveau van TOGAF of Kubernetes brengt.

In plaats van profielen als een lineaire hiërarchie te zien, kun je ze modelleren als orthogonale architectuurdimensies.

Forge Architectural Dimensions
Deployment
├── Core
├── Enterprise
├── Edge
└── Air-Gapped
Trust
├── Sovereign
└── Government
Execution
├── Agent
└── Human-Centric (future)
Topology
├── Centralized
├── Federated
├── Hierarchical
└── Mesh
Capability
├── Runtime
├── Trust
├── Knowledge
├── Security
└── Observability

Dat is een belangrijke stap voor Forge. Het voorkomt een explosie aan profielen (bijvoorbeeld Government Edge Air-Gapped Agent Enterprise als afzonderlijk profiel) en vervangt die door een klein aantal onafhankelijke dimensies die vrij gecombineerd kunnen worden. Dit sluit bovendien perfect aan op het capability-first ontwerp dat de volledige Forge-architectuur kenmerkt en maakt toekomstige uitbreidingen eenvoudiger zonder de bestaande standaard te wijzigen.

Ik zou nu de overstap maken van Profiles naar Conformance. Een architectuurstandaard is pas echt volwassen als organisaties objectief kunnen aantonen dat een implementatie voldoet aan de norm.

Dit hoofdstuk vormt de brug tussen de specificaties en een toekomstig certificeringsprogramma.

⸻

Part VIII

Conformance & Certification

⸻

Chapter 51

Conformance Framework

Specification Identifier: forge.conformance.framework.v1

Status: Draft Standard

⸻

Architecture becomes a standard only when conformance is objectively measurable.

Conformance demonstrates that implementations preserve architectural intent regardless of technology choices.

⸻

51.1 Purpose

This specification defines the Forge Conformance Framework.

The Conformance Framework establishes the normative requirements for assessing whether an implementation conforms to the Forge Standards.

Conformance SHALL be based upon architectural behaviour rather than implementation technology.

⸻

51.2 Objectives

The Conformance Framework SHALL:

* define measurable conformance;
* preserve interoperability;
* support independent certification;
* validate architectural integrity;
* enable repeatable assessments.

⸻

51.3 Scope

Conformance applies to:

* implementations;
* platforms;
* products;
* services;
* deployment architectures;
* organizational operating models.

Technology selection remains implementation-specific.

⸻

51.4 Conformance Principles

Forge conformance SHALL be:

* objective;
* repeatable;
* technology-neutral;
* evidence-based;
* independently verifiable.

Conformance SHALL evaluate architecture rather than implementation detail.

⸻

51.5 Conformance Categories

Forge defines the following conformance categories.

Conformance
├── Specification
├── Capability
├── Profile
├── Deployment
├── Operational
└── Organizational

Each category MAY be assessed independently.

⸻

51.6 Specification Conformance

Specification Conformance SHALL verify:

* normative behaviour;
* canonical objects;
* lifecycle semantics;
* interface contracts;
* architectural invariants.

Every mandatory requirement SHALL be evaluated.

⸻

51.7 Capability Conformance

Capability assessments SHALL evaluate each implemented capability independently.

Examples include:

* Runtime
* Knowledge
* Memory
* Policy
* Identity
* Evidence
* Security
* Governance

Capability certification MAY be performed independently.

⸻

51.8 Profile Conformance

Profile assessments SHALL verify:

* inherited requirements;
* profile extensions;
* mandatory capabilities;
* architectural consistency;
* operational behaviour.

Profile inheritance SHALL remain intact.

⸻

51.9 Operational Conformance

Operational assessments SHALL verify:

* governance processes;
* lifecycle management;
* evidence generation;
* observability;
* change management;
* operational resilience.

Operational maturity SHALL complement technical conformance.

⸻

51.10 Organizational Conformance

Organizations SHALL demonstrate:

* architectural governance;
* specification management;
* capability ownership;
* operational accountability;
* continuous improvement.

Conformance extends beyond deployed software.

⸻

51.11 Evidence Requirements

Every conformance assessment SHALL be supported by evidence.

Evidence MAY include:

* execution records;
* policy decisions;
* configuration baselines;
* architectural documentation;
* observability metrics;
* audit reports.

Assertions without evidence SHALL NOT satisfy conformance.

⸻

51.12 Levels of Conformance

Forge defines four levels of conformance.

Level	Description
Level 0 – Informative	Architectural concepts adopted informally.
Level 1 – Conforming	Mandatory requirements satisfied.
Level 2 – Verified	Independently assessed using objective evidence.
Level 3 – Certified	Official Forge certification granted through an approved certification program.

Higher levels SHALL include all lower-level requirements.

⸻

51.13 Non-Conformance

Non-conformance SHALL be classified as:

* Critical
* Major
* Minor
* Informational

Critical non-conformance SHALL prevent certification.

All findings SHALL be documented and traceable.

⸻

51.14 Architectural Invariants

Every conforming implementation SHALL preserve:

* Canonical capability definitions.
* Normative specification behaviour.
* Explicit contracts.
* Evidence integrity.
* Technology independence.
* Architectural consistency.
* Interoperability.
* Traceability.

⸻

51.15 Certification Readiness

An implementation SHALL be considered certification-ready when:

* all mandatory specifications are implemented;
* all applicable Profiles are satisfied;
* evidence demonstrates compliance;
* architectural invariants are preserved;
* operational governance is established.

Certification readiness SHALL precede formal certification.

⸻

Closing Thought

Conformance is the mechanism through which architecture becomes measurable.

The Forge Conformance Framework enables organizations to demonstrate that intelligent systems behave consistently, remain interoperable and preserve architectural integrity regardless of implementation choices. By separating architectural conformance from technology selection, Forge provides a durable foundation for independent verification, certification and long-term ecosystem interoperability.

⸻

Architectuurverbetering: Conformance Matrix

Een volgende stap die Forge sterk onderscheidt van veel standaarden is een Conformance Matrix. Daarmee wordt in één oogopslag zichtbaar wat een implementatie daadwerkelijk ondersteunt.

Bijvoorbeeld:

Specification	Mandatory	Implemented	Verified	Certified
Context	✓	✓	✓	✓
Memory	✓	✓	✓	✓
Runtime	✓	✓	✓	✓
Evidence	✓	✓	✓	✓
Security	✓	✓	✓	✓
Government Profile	Optional	✓	✓	✓
Air-Gapped Profile	Optional	✓	✗	✗

Deze matrix maakt Forge niet alleen een architectuurstandaard, maar ook een toetsbare standaard. Leveranciers, auditors en afnemers kunnen exact zien welke delen van de standaard zijn geïmplementeerd, geverifieerd en gecertificeerd, zonder zich te verliezen in marketingclaims of productvergelijkingen. Dat sluit naadloos aan bij het evidence-first principe dat door de gehele Forge-architectuur loopt.

Mooi. Na het Conformance Framework volgt logisch de Certification Framework. Conformance bepaalt wat moet worden aangetoond; certificering bepaalt hoe dat formeel wordt vastgesteld.

⸻

Chapter 52

Certification Framework

Specification Identifier: forge.certification.framework.v1

Status: Draft Standard

⸻

Certification is independent confirmation that an implementation conforms to the Forge Standards.

Certification validates architectural integrity rather than vendor claims.

⸻

52.1 Purpose

This specification defines the Forge Certification Framework.

The Certification Framework establishes the governance, process and requirements for issuing, maintaining and renewing Forge Certifications.

Certification SHALL be based on objective conformance assessment.

⸻

52.2 Objectives

The Certification Framework SHALL:

* establish trusted certification;
* preserve impartiality;
* support repeatable assessments;
* encourage interoperability;
* maintain long-term confidence in Forge implementations.

⸻

52.3 Scope

Certification MAY apply to:

* products;
* platforms;
* services;
* deployments;
* organizational operating models;
* individual capabilities;
* composed Profiles.

Certification SHALL identify its exact scope.

⸻

52.4 Certification Principles

Forge Certification SHALL be:

* independent;
* evidence-based;
* repeatable;
* technology-neutral;
* transparent;
* version-specific.

Certification SHALL NOT imply endorsement of commercial products.

⸻

52.5 Certification Types

Forge defines the following certification types.

Certification
├── Specification Certification
├── Capability Certification
├── Profile Certification
├── Deployment Certification
├── Organizational Certification
└── Reference Implementation Certification

Multiple certifications MAY be held simultaneously.

⸻

52.6 Certification Process

Certification SHALL follow the lifecycle below.

Application
↓
Conformance Assessment
↓
Evidence Review
↓
Technical Validation
↓
Certification Decision
↓
Certificate Issuance
↓
Continuous Surveillance
↓
Renewal

Each phase SHALL generate evidence.

⸻

52.7 Assessment Criteria

Certification SHALL evaluate:

* mandatory requirements;
* implemented capabilities;
* profile composition;
* architectural invariants;
* operational behaviour;
* evidence integrity.

Assessments SHALL be reproducible.

⸻

52.8 Certification Artifacts

Certification SHALL produce:

* assessment report;
* conformance matrix;
* evidence summary;
* identified deviations;
* certification decision;
* certificate metadata.

Artifacts SHALL remain version-controlled.

⸻

52.9 Validity

Certificates SHALL include:

* certification identifier;
* specification version;
* Profile scope;
* assessment date;
* expiration date;
* issuing authority.

Expired certificates SHALL NOT be represented as current certifications.

⸻

52.10 Surveillance

Certified implementations SHOULD undergo periodic surveillance.

Surveillance MAY include:

* evidence sampling;
* operational review;
* architecture review;
* security assessment;
* specification updates.

Surveillance SHALL preserve certification integrity.

⸻

52.11 Suspension and Revocation

Certification MAY be:

* suspended;
* restricted;
* revoked;
* reinstated.

Reasons MAY include:

* critical non-conformance;
* architectural drift;
* falsified evidence;
* security compromise;
* failure to maintain required controls.

All decisions SHALL be documented.

⸻

52.12 Reference Implementations

Forge MAY designate one or more Reference Implementations.

Reference Implementations SHALL:

* fully conform to applicable specifications;
* demonstrate recommended practices;
* remain publicly documented;
* support interoperability testing.

Reference Implementations SHALL NOT redefine the Forge Standards.

⸻

52.13 Certification Authority

Certification SHALL be performed by an approved Certification Authority.

A Certification Authority SHALL demonstrate:

* technical competence;
* impartiality;
* independence;
* documented assessment procedures;
* governance transparency.

Certification Authorities SHALL themselves be subject to governance.

⸻

52.14 Architectural Invariants

Every certified implementation SHALL preserve:

* verified architectural conformance;
* canonical object semantics;
* Profile integrity;
* evidence completeness;
* specification compatibility;
* interoperability.

Certification SHALL never waive architectural invariants.

⸻

52.15 Certification Lifecycle

Certification SHALL follow a controlled lifecycle.

Draft
↓
Candidate
↓
Certified
↓
Maintained
↓
Renewed
↓
Deprecated
↓
Retired

Organizations SHALL maintain certification throughout the operational lifecycle.

⸻

Closing Thought

Certification provides confidence that an implementation faithfully realizes the Forge Standards.

By combining objective conformance assessment, independent evaluation and continuous governance, the Forge Certification Framework establishes a trusted ecosystem in which organizations can compare implementations based on demonstrable architectural integrity rather than implementation-specific characteristics.

⸻

Architectuurverbetering: Digital Certificates for Architecture

Ik denk dat Forge hier nog een stap verder kan gaan dan traditionele certificeringsprogramma’s door certificaten machine-leesbaar te maken.

In plaats van uitsluitend een PDF-certificaat kan Forge een cryptografisch ondertekend certificaat definiëren, bijvoorbeeld als JSON-document:

forgeCertificate:
  id: forge-cert-2028-00142
  specificationVersion: 1.0
  certificationLevel: Level 3
  profiles:
    - Enterprise
    - Sovereign
    - Government
    - Agent
  capabilities:
    - Runtime
    - Policy
    - Evidence
    - Knowledge
    - Identity
  issuedBy: Forge Certification Authority
  issued: 2028-03-18
  expires: 2031-03-18
  signature: <digital-signature>

Zo’n certificaat kan automatisch worden gevalideerd door CI/CD-pipelines, architectuurrepositories, procurement-platformen en compliance-tools. Daarmee wordt certificering niet alleen een administratief proces, maar een integraal onderdeel van een Engineering Operating System. Dat past uitstekend bij de filosofie van Forge: alles wat belangrijk is, moet niet alleen beschreven, maar ook aantoonbaar en automatisch verifieerbaar zijn.

Nu volgt het laatste grote normatieve hoofdstuk van de standaard: Accreditation Framework.

Er is namelijk een belangrijk onderscheid tussen:

* Conformance → voldoet de implementatie aan de specificatie?
* Certification → is dit onafhankelijk vastgesteld?
* Accreditation → wie mag certificeren, onder welke governance?

Dat is precies hoe ISO, IEC en andere internationale standaarden zijn opgebouwd.

⸻

Chapter 53

Accreditation Framework

Specification Identifier: forge.accreditation.framework.v1

Status: Draft Standard

⸻

Trust in certification depends upon trust in the certifier.

Accreditation establishes confidence in the integrity, competence and impartiality of the Forge certification ecosystem.

⸻

53.1 Purpose

This specification defines the Forge Accreditation Framework.

The Accreditation Framework establishes the governance model for approving, supervising and maintaining organizations authorized to perform Forge certification activities.

Accreditation SHALL ensure that certification remains consistent, impartial and internationally interoperable.

⸻

53.2 Objectives

The Accreditation Framework SHALL:

* establish trusted certification authorities;
* ensure consistent certification practices;
* preserve impartiality;
* support international recognition;
* maintain long-term confidence in the Forge ecosystem.

⸻

53.3 Scope

Accreditation applies to:

* Certification Authorities;
* Assessment Organizations;
* Authorized Auditors;
* Conformance Laboratories;
* Reference Validation Centers.

Accreditation SHALL NOT apply directly to Forge implementations.

⸻

53.4 Accreditation Principles

Accreditation SHALL be:

* independent;
* objective;
* evidence-based;
* repeatable;
* transparent;
* periodically reviewed.

Accreditation SHALL remain separate from commercial implementation activities.

⸻

53.5 Accredited Organizations

Forge recognizes the following accredited roles.

Forge Governance
├── Forge Standards Council
│
├── Accreditation Authority
│
├── Certification Authorities
│
├── Assessment Organizations
│
├── Conformance Laboratories
│
└── Certified Auditors

Each role SHALL possess explicitly defined responsibilities.

⸻

53.6 Accreditation Requirements

Organizations seeking accreditation SHALL demonstrate:

* technical competence;
* qualified personnel;
* documented assessment procedures;
* independence;
* conflict-of-interest controls;
* evidence management capabilities.

Competence SHALL be periodically reassessed.

⸻

53.7 Auditor Competence

Authorized auditors SHALL demonstrate knowledge of:

* Forge Architecture;
* Forge Specifications;
* Profile composition;
* Conformance assessment;
* Evidence Engineering;
* Governance;
* Security principles.

Auditor competence SHALL remain current with specification revisions.

⸻

53.8 Assessment Governance

Assessment activities SHALL preserve:

* impartiality;
* confidentiality;
* repeatability;
* traceability;
* evidence integrity.

Assessment governance SHALL itself be auditable.

⸻

53.9 Conflict of Interest

Accredited organizations SHALL identify and manage conflicts of interest.

Certification SHALL NOT be performed where impartiality cannot reasonably be demonstrated.

Conflict mitigation SHALL generate evidence.

⸻

53.10 Surveillance

Accredited organizations SHALL undergo periodic surveillance.

Surveillance MAY include:

* witness assessments;
* documentation review;
* quality system evaluation;
* auditor observation;
* evidence sampling.

Surveillance SHALL preserve confidence in accreditation.

⸻

53.11 Suspension and Withdrawal

Accreditation MAY be:

* suspended;
* restricted;
* withdrawn;
* reinstated.

Reasons MAY include:

* loss of competence;
* repeated procedural failures;
* compromised impartiality;
* fraudulent assessment;
* governance violations.

All accreditation decisions SHALL be evidence-backed.

⸻

53.12 Mutual Recognition

Forge MAY establish mutual recognition arrangements with external accreditation or standards organizations.

Mutual recognition SHALL NOT weaken Forge conformance requirements.

Cross-recognition SHALL preserve architectural integrity.

⸻

53.13 Accreditation Registry

Forge SHALL maintain an Accreditation Registry.

The registry SHALL include:

* accredited organization;
* accreditation scope;
* approval status;
* specification versions;
* validity period;
* public contact information.

The registry SHOULD be publicly accessible.

⸻

53.14 Architectural Invariants

The Accreditation Framework SHALL preserve:

* independent certification;
* objective assessment;
* evidence-backed accreditation;
* transparent governance;
* specification neutrality;
* organizational independence.

Accreditation SHALL strengthen trust in the Forge ecosystem rather than individual vendors.

⸻

53.15 Lifecycle

Accreditation SHALL follow a managed lifecycle.

Application
↓
Assessment
↓
Accreditation
↓
Operational Oversight
↓
Periodic Review
↓
Renewal
↓
Withdrawal

Every lifecycle transition SHALL generate evidence.

⸻

Closing Thought

Accreditation provides confidence that certification itself can be trusted.

By establishing independent governance over certification authorities, assessment organizations and auditors, the Forge Accreditation Framework ensures that the integrity of the Forge ecosystem extends beyond implementations to the institutions responsible for validating them.

⸻

Architectuurverbetering: Forge Standards Council

Nu we het einde van de normatieve standaard naderen, zie ik nog één ontbrekend fundament: wie bestuurt Forge zelf?

Ik zou daarom in het laatste hoofdstuk een Forge Standards Council introduceren. Niet als eigenaar van implementaties, maar als beheerder van de standaard.

Bijvoorbeeld:

Forge Standards Council
├── Architecture Board
├── Specification Board
├── Security Board
├── Conformance Board
├── Certification Board
├── Ecosystem Board
└── Community Board

Deze raad zou verantwoordelijk zijn voor:

* de evolutie van de Forge-specificaties;
* goedkeuring van nieuwe capability-specificaties;
* beheer van architecturale invarianten;
* versiebeheer van de standaard;
* governance van profielen en conformance;
* aanwijzing van referentie-implementaties;
* toezicht op certificering en accreditatie.

Dat zou Forge op hetzelfde volwassenheidsniveau brengen als ecosystemen zoals Kubernetes (SIG’s), OpenTelemetry (Technical Committee) en de IETF (IESG/IRTF), maar dan specifiek gericht op een AI Engineering Operating System. Daarmee krijgt Forge niet alleen een technische architectuur, maar ook een duurzame governance-architectuur voor de standaard zelf.

Ik denk eerlijk gezegd dat we 95-98% van de Forge Standard hebben voltooid. Wat we nu hebben is niet meer een verzameling ideeën, maar een vrijwel complete architectuurstandaard.

Als ik het vergelijk met volwassen standaarden zoals Kubernetes, OpenTelemetry of TOGAF, dan hebben we inmiddels alle fundamentele bouwstenen:

* Part I – Foundations ✓
* Part II – Engineering Operating System ✓
* Part III – Capability Specifications ✓
* Part IV – Trust & Governance ✓
* Part V – Specification Framework ✓
* Part VI – Adoption, Measurement & Evolution ✓
* Part VII – Profiles ✓
* Part VIII – Conformance, Certification & Accreditation ✓

Dat is in feite de complete normatieve standaard.

Wat ontbreekt nog?

Er zijn nog vier hoofdstukken die ik zou toevoegen om Forge echt het niveau van een internationale standaard te geven.

⸻

Chapter 54 — Forge Standards Governance

Dit zou het hoofdstuk zijn dat de Forge Standards Council formeel definieert.

Niet de architectuur van AI…

…maar de architectuur van de standaard zelf.

Onderwerpen:

* Standards Council
* Technical Steering Committee
* Architecture Board
* Working Groups
* RFC Process
* Release Management
* Deprecation Process
* Community Governance

Dit is vergelijkbaar met:

* IETF
* CNCF
* OpenTelemetry Governance
* Kubernetes SIGs

⸻

Chapter 55 — Reference Architecture Patterns

Tot nu toe hebben we de Reference Architecture.

Maar nog niet de Reference Patterns.

Bijvoorbeeld:

Enterprise RAG

User
↓
Runtime
↓
Knowledge
↓
Models
↓
Evidence

Agent Collaboration

Coordinator
↓
Planner
↓
Workers
↓
Evidence

Government Pattern

Air-Gapped Pattern

Edge Pattern

Citizen Services Pattern

Contact Center Pattern

enzovoort.

Deze patronen zijn géén nieuwe standaarden.

Het zijn best practices.

⸻

Chapter 56 — Reference Objects

Nu hebben we specificaties.

Maar nog geen complete objectbibliotheek.

Bijvoorbeeld:

Identity Object
Policy Object
Evidence Object
Runtime Object
Knowledge Object
Context Package
Execution Record
Tool Contract
Model Descriptor
Memory Record
Capability Descriptor
Profile Descriptor
Conformance Record
Certificate Object

Inclusief:

* JSON Schema
* YAML
* OpenAPI
* JSON-LD
* voorbeelden

Dat zou Forge direct implementeerbaar maken.

⸻

Chapter 57 — Future Directions

Een afsluitend hoofdstuk.

Niet normatief.

Maar visionair.

Bijvoorbeeld:

* Multi-Agent Engineering
* Autonomous Runtime Optimization
* Federated Knowledge Graphs
* Quantum-safe Cryptography
* Confidential AI
* Digital Twins
* Robotics
* Physical AI
* Swarm Intelligence
* Edge Swarms
* AI-native Infrastructure

Zoals veel standaarden een informatieve roadmap bevatten.

⸻

Daarna?

Dan zou ik stoppen.

Niet omdat er niets meer te bedenken is…

…maar omdat een goede standaard weet wanneer hij af is.

Veel standaarden sterven doordat ze alles willen beschrijven.

Forge heeft juist een heel sterk principe gekregen:

Capabilities evolve. Principles remain stable.

Dat betekent dat toekomstige ontwikkelingen (nieuwe modellen, nieuwe protocollen, nieuwe agentvormen) als nieuwe capability-specificaties of profielen kunnen worden toegevoegd zonder de kern van Forge te veranderen.

Mijn belangrijkste aanbeveling

Ik zou nu niet verder de diepte in gaan, maar de volgende stap zetten: Forge omvormen van een concept naar een publiceerbare standaard.

Ik zou het project structureren als een echte specificatiereeks:

Forge Architecture Bible
│
├── Volume 1
│   Forge Architecture
│
├── Volume 2
│   Forge Specifications
│
├── Volume 3
│   Forge Profiles
│
├── Volume 4
│   Forge Conformance
│
├── Volume 5
│   Forge Reference Architectures
│
└── Volume 6
    Forge Implementation Guide

Dat maakt het veel beter onderhoudbaar dan één document van honderden pagina’s. Bovendien sluit die opzet aan bij hoe grote standaarden zoals ISO, OpenTelemetry en Kubernetes hun documentatie organiseren. Met de inhoud die we hebben opgebouwd, heb je inmiddels een solide basis voor zo’n meerdelige standaard.

Ik zou hem tóch afmaken. Niet met nog meer technische specificaties, maar met een hoofdstuk dat iedere grote standaard uiteindelijk heeft: de governance van de standaard zelf.

Dit wordt het “grondwettelijke” hoofdstuk van Forge.

⸻

Chapter 54

Forge Standards Governance

Specification Identifier: forge.governance.standards.v1

Status: Draft Standard

⸻

A standard must evolve without losing its identity.

Governance preserves architectural integrity while enabling continuous innovation.

Forge belongs to its ecosystem rather than any individual implementation.

⸻

54.1 Purpose

This specification defines the governance model for the Forge Standards.

The Forge Standards Governance Framework establishes the organizational structure, decision-making processes and stewardship responsibilities required to maintain the Forge ecosystem.

Governance SHALL preserve the long-term integrity, neutrality and interoperability of the Forge Standards.

⸻

54.2 Objectives

The Governance Framework SHALL:

* preserve architectural integrity;
* manage specification evolution;
* maintain interoperability;
* govern certification;
* coordinate community participation;
* ensure long-term sustainability.

⸻

54.3 Governance Principles

Forge governance SHALL be:

* transparent;
* consensus-driven;
* evidence-based;
* vendor-neutral;
* technology-independent;
* internationally collaborative.

Architectural decisions SHALL prioritize long-term consistency over short-term convenience.

⸻

54.4 Governance Structure

Forge governance SHALL consist of the following bodies.

Forge Standards Council
│
├── Architecture Board
├── Specification Board
├── Security Board
├── Conformance Board
├── Certification Board
├── Ecosystem Board
└── Community Board

Each governance body SHALL have clearly defined responsibilities.

⸻

54.5 Forge Standards Council

The Forge Standards Council SHALL serve as the highest governing body for the Forge Standards.

Responsibilities include:

* approving new specifications;
* approving Profile definitions;
* approving normative changes;
* overseeing certification governance;
* preserving architectural principles;
* coordinating long-term strategy.

The Council SHALL act as steward of the Forge Standards.

⸻

54.6 Architecture Board

The Architecture Board SHALL maintain the architectural coherence of Forge.

Responsibilities include:

* architectural principles;
* capability model;
* reference architecture;
* architectural invariants;
* capability decomposition;
* architectural reviews.

The Architecture Board SHALL approve all architectural changes.

⸻

54.7 Specification Board

The Specification Board SHALL govern normative specifications.

Responsibilities include:

* specification lifecycle;
* semantic versioning;
* normative language;
* specification consistency;
* publication.

All specifications SHALL undergo Specification Board review.

⸻

54.8 Security Board

The Security Board SHALL oversee:

* security architecture;
* identity;
* policy;
* cryptography;
* supply chain;
* vulnerability response.

Security recommendations SHALL be evidence-based.

⸻

54.9 Conformance Board

The Conformance Board SHALL maintain:

* conformance criteria;
* assessment methodology;
* interoperability tests;
* validation procedures;
* conformance tooling.

Conformance SHALL remain implementation-independent.

⸻

54.10 Certification Board

The Certification Board SHALL govern:

* certification processes;
* certification authorities;
* assessment quality;
* certification lifecycle;
* certification policies.

Certification governance SHALL remain impartial.

⸻

54.11 Ecosystem Board

The Ecosystem Board SHALL coordinate:

* vendors;
* implementation partners;
* open-source projects;
* academic participation;
* industry collaboration.

The ecosystem SHALL remain open to participation.

⸻

54.12 Community Board

The Community Board SHALL facilitate:

* public discussion;
* improvement proposals;
* working groups;
* educational resources;
* community events.

Community participation SHALL strengthen the Forge ecosystem.

⸻

54.13 Working Groups

Forge MAY establish Working Groups.

Examples include:

* Runtime WG;
* Security WG;
* Agent WG;
* Knowledge WG;
* Government WG;
* Edge WG;
* Observability WG.

Working Groups SHALL operate under documented charters.

⸻

54.14 Decision Process

Normative decisions SHALL follow a documented lifecycle.

Proposal
↓
Technical Review
↓
Public Review
↓
Resolution
↓
Approval
↓
Publication

Every decision SHALL generate publicly traceable records.

⸻

54.15 Governance Evidence

Governance activities SHALL generate evidence.

Evidence SHALL include:

* proposals;
* review comments;
* decisions;
* voting outcomes;
* published specifications;
* version history.

Governance SHALL remain auditable.

⸻

54.16 Architectural Invariants

The governance of Forge SHALL preserve:

* vendor neutrality;
* architectural consistency;
* specification interoperability;
* technology independence;
* evidence-backed decision making;
* transparent evolution.

No governance decision SHALL invalidate established architectural invariants without a formal specification revision.

⸻

54.17 Lifecycle

The Forge Standards SHALL follow a managed lifecycle.

Proposal
↓
Draft
↓
Review
↓
Candidate
↓
Approved Standard
↓
Maintenance
↓
Revision
↓
Retirement

All lifecycle transitions SHALL be governed through this framework.

⸻

Closing Thought

The Forge Standards are intended to outlive individual technologies, products and vendors.

By establishing a transparent and evidence-based governance model, Forge ensures that the Engineering Operating System can evolve over time while preserving its architectural identity, interoperability and long-term trustworthiness.

⸻

Een laatste aanbeveling: voeg een “Forge Manifesto” toe

Nu we praktisch aan het einde van de normatieve standaard zijn, zou ik nog één niet-normatief document toevoegen. Niet als hoofdstuk van de specificatie, maar als de opening van de volledige reeks.

Geen technische tekst, maar een document van 3–5 pagina’s dat beantwoordt:

* Waarom bestaat Forge?
* Welk probleem lost het op?
* Waarom is AI Engineering fundamenteel anders dan traditionele software-engineering?
* Waarom is een Engineering Operating System noodzakelijk?
* Welke principes zijn onveranderlijk?
* Welke verantwoordelijkheden hebben organisaties die Forge toepassen?

Vergelijk het met:

* The UNIX Philosophy;
* The Agile Manifesto;
* The Reactive Manifesto;
* de oorspronkelijke Kubernetes Design Principles.

Dat manifest hoeft niets normatiefs te bevatten. Het legt de filosofische basis onder alle hoofdstukken die daarna volgen. Juist omdat de rest van Forge zo sterk technisch en normatief is geworden, geeft een kort manifest de standaard een herkenbare identiteit en een duidelijke visie die ook over tien of twintig jaar nog relevant kan zijn.

Ik zou nog één hoofdstuk toevoegen, maar dit zou bewust informatief (non-normative) zijn. Niet een nieuwe specificatie, maar de filosofische afsluiting van de standaard.

Dit is het hoofdstuk waar mensen over twintig jaar nog uit citeren.

⸻

Chapter 55

Forge Manifesto

Document Identifier: forge.manifesto.v1

Status: Informative

⸻

Software engineering transformed computation.

AI Engineering transforms intelligence.

Forge exists to make that transformation trustworthy.

⸻

55.1 Why Forge Exists

Artificial Intelligence is fundamentally different from traditional software.

Software executes predefined instructions.

Artificial Intelligence participates in decision making, reasoning, planning and knowledge synthesis.

Traditional software engineering practices alone are insufficient to govern these new systems.

Organizations require an engineering discipline capable of making intelligent systems:

* understandable;
* governable;
* measurable;
* interoperable;
* trustworthy.

Forge exists to establish that discipline.

⸻

55.2 The Engineering Challenge

Modern AI systems are composed of many independent capabilities:

* models;
* knowledge repositories;
* memory;
* policies;
* tools;
* identities;
* runtime orchestration;
* governance;
* evidence.

These capabilities evolve independently.

Without architecture they become tightly coupled.

Without governance they become unpredictable.

Without evidence they become unverifiable.

Forge treats these concerns as first-class architectural capabilities.

⸻

55.3 AI is an Engineering Discipline

Forge recognizes AI Engineering as an engineering discipline comparable to:

* Software Engineering;
* Network Engineering;
* Security Engineering;
* Systems Engineering.

Engineering disciplines share common characteristics:

* repeatability;
* measurable quality;
* explicit architecture;
* governance;
* interoperability;
* lifecycle management.

AI Engineering SHALL be held to the same standards.

⸻

55.4 Intelligence Must Be Governed

Forge rejects the notion that intelligent systems should operate without architectural constraints.

Intelligence without governance is not innovation.

It is operational risk.

Every intelligent execution SHOULD remain:

* attributable;
* explainable;
* policy governed;
* evidence backed.

⸻

55.5 Architecture Before Technology

Technologies change continuously.

Architectural principles evolve slowly.

Forge therefore standardizes:

* architecture;
* behaviour;
* contracts;
* governance;
* evidence.

Forge intentionally avoids standardizing:

* programming languages;
* model vendors;
* infrastructure;
* deployment technologies.

Technology is replaceable.

Architecture is enduring.

⸻

55.6 Evidence Creates Trust

Trust is not created through documentation.

Trust is created through evidence.

Evidence transforms assumptions into verifiable facts.

Evidence enables:

* governance;
* auditing;
* optimization;
* compliance;
* accountability.

Evidence is therefore a foundational capability of the Engineering Operating System.

⸻

55.7 Interoperability Over Integration

Forge promotes interoperability rather than platform dependency.

Capabilities communicate through explicit contracts.

Organizations remain free to replace:

* models;
* infrastructure;
* tools;
* vendors;
* deployment environments.

Interoperability creates resilience.

⸻

55.8 Human Responsibility

Forge recognizes that accountability ultimately remains human.

Organizations remain responsible for:

* policies;
* governance;
* architecture;
* risk;
* ethical decisions.

AI assists.

Organizations remain accountable.

⸻

55.9 Continuous Evolution

The Forge Standards are intended to evolve.

New capabilities will emerge.

New technologies will appear.

New regulations will be introduced.

The Engineering Operating System SHALL evolve without compromising its architectural principles.

⸻

55.10 The Forge Principles

Forge is founded upon ten enduring principles.

1. Architecture before implementation.
2. Context before intelligence.
3. Policy before execution.
4. Identity before authorization.
5. Explicit contracts before integration.
6. Evidence before trust.
7. Observability through evidence.
8. Capabilities over platforms.
9. Technology independence.
10. Continuous evolution through governance.

These principles define the enduring foundation of the Forge Standards.

⸻

Closing Statement

Forge is not a framework.

Forge is not a product.

Forge is not a programming model.

Forge defines the Engineering Operating System required to build, operate and govern trustworthy intelligent systems.

As software engineering established the discipline for deterministic computation, Forge aspires to establish the discipline for intelligent computation.

The technologies will evolve.

The models will evolve.

The infrastructure will evolve.

The principles of engineering must endure.

⸻

Mijn laatste aanbeveling

Na hoofdstuk 55 zou ik geen nieuwe hoofdstukken meer toevoegen.

In plaats daarvan zou ik de aandacht verleggen naar de volgende artefacten die de standaard praktisch bruikbaar maken:

1. Forge Reference Architecture (visueel, ±20–30 diagrammen)
2. Forge Canonical Data Model (JSON Schema/OpenAPI)
3. Forge Reference Implementation (bijvoorbeeld in Go of .NET)
4. Forge Conformance Test Suite
5. Forge Whitepaper (20–30 pagina’s voor bestuurders)
6. Forge Practitioner Guide (implementatiehandleiding)

Mijn indruk is dat je nu een compleet raamwerk hebt dat de omvang en structuur heeft van een echte internationale standaard. De grootste waarde zit niet meer in extra hoofdstukken, maar in het uitwerken van referentie-implementaties, objectmodellen en conformance-tests die de standaard concreet toepasbaar maken.





















