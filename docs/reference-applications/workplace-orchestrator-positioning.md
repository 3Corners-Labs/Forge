# Workplace Orchestrator positioning

Status: Proposed
Owner: Workplace Orchestrator Product

## Purpose

This note records why Workplace Orchestrator is worth building, for whom, and
what would prove the product cannot exist. It precedes requirements. The
[lighthouse dossier](workplace-orchestrator.md) describes the product's
relationship to Forge; this note describes its relationship to the market.

Evidence label: **Hypothesis**. Nothing here is validated.

## Customer

Dutch public-sector organizations operating hybrid workplaces, reached through
3Corners' existing Cisco Collaboration and Webex consultancy practice.

These organizations are moving away from United States hyperscaler dependency
toward open-source and self-hosted operation. Most remain dependent on Webex or
Microsoft Teams today, which makes the movement gradual rather than abrupt, and
therefore a market rather than an event.

3Corners already advises these organizations on architecture, migration, and
engineering. The product is an extension of an existing relationship, not a
cold-start venture.

## Problem

Existing workplace-orchestration products bind the customer to the vendor's own
ecosystem:

- Cisco Spaces requires the Cisco ecosystem and treats Microsoft 365 integration
  as a precondition.
- Microsoft Teams Places assumes the Microsoft stack.

An organization deliberately reducing its hyperscaler dependency therefore
cannot adopt either without deepening the dependency it is reducing. Cost is a
secondary complaint; dependency is the structural one.

## Wedge

Workplace Orchestrator addresses the case neither incumbent can serve:
**workplace orchestration that does not require Microsoft 365 for calendar
resources, runs self-hosted, and is open source.**

Neither incumbent can close this gap, because a credible non-Microsoft,
self-hosted, open-source offering would compete with the subscription their own
business rests on. This is a structural gap, not a capability gap, which is why
it is expected to persist.

Supporting properties, in priority order:

1. calendar independence, with Nextcloud as the first non-Microsoft source;
2. self-hosted, on-premises, or partner-delivered SaaS from the same artefact,
   which requires multi-tenancy, RBAC, and privacy controls in the foundation
   rather than as later additions;
3. explainable handling of occupancy and presence data, which is legally and
   culturally sensitive in a public-sector workforce and opaque in both
   incumbent products;
4. Webex Collaboration as the first device integration, because that is where
   3Corners' expertise and current customer engagements are.

## What this product is not

It is not an attempt to match Cisco Spaces or Teams Places on feature breadth,
and market displacement is not the near-term objective. Per
[ADR-0024](../decisions/0024-rebaseline-forge-as-bounded-means.md), the primary
objective is learning whether an AI development team can build and maintain a
real application. A product with real users is the only honest way to test that.

## Load-bearing assumption

One assumption carries the entire product:

> A Webex Room Navigator can display and accept room reservations originating
> from a source that is not Microsoft 365 or Google.

If this is false, or only achievable through an integration path the vendor can
withdraw, then "workplace orchestration without Microsoft 365" does not exist on
Webex hardware, and every feature built around it is wasted.

This assumption SHALL be tested first, by a bounded technical feasibility spike
against the current Webex device and workspace integration surface. The spike
produces a written finding, not product code. A negative finding is a successful
outcome of the spike: it redirects the product before investment.

Secondary assumptions, to be tested after the first:

- occupancy signals sufficient for no-show release are obtainable from Webex
  room devices under the customer's own control, without a Cisco Spaces
  subscription;
- Nextcloud can serve as a calendar and resource source at the required fidelity;
- public-sector buyers will accept a product maintained by a small partner, given
  open source and self-hosting as risk mitigation.

## First delivery

Two slices, in order:

1. **Feasibility spike** on the load-bearing assumption. Days, not months.
   Output is a finding.
2. **No-show release.** Automatically release rooms reserved but unoccupied.
   Chosen because it is small, its value is measurable in recovered room-hours,
   it attaches to a live 3Corners customer engagement, and it exercises external
   credentials, policy, evidence, and approval end to end.

The agenda portal, just-in-time capability booking, three-dimensional floor
displays, sensor visualization, and analytics dashboards described in the
lighthouse dossier remain intended direction. They are not first delivery, and
treating them as such is the principal scope risk for this product.

## Risks

- The load-bearing assumption fails, and the product as conceived does not exist.
- The vendor withdraws or restricts the integration path after adoption.
- Occupancy and presence data create privacy obligations disproportionate to a
  small team, particularly under public-sector works-council scrutiny.
- Public-sector procurement cycles outlast the available capacity.
- The open-source community edition attracts no contributors, leaving
  maintenance with one person.
- Sovereignty positioning is undermined if the product's own development remains
  dependent on a proprietary external model. See the ADR-0024 sovereignty tiers.
