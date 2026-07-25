# AI delivery experiment measurement

Status: Accepted
Owner: Forge Product and Architecture

## Purpose

[ADR-0024](../decisions/0024-rebaseline-forge-as-bounded-means.md) states that
Forge exists to answer one question:

> Can an AI development team, under human accountability, build **and maintain** a
> real application without a large human engineering organization?

This document makes that question falsifiable. Without declared criteria the
exercise produces an impression rather than a result, and a favourable
impression is the most likely outcome of any effort its author is invested in.

Measurement begins when Workplace delivery begins. Figures recorded before then
describe Forge building Forge, which is a different and easier problem.

## What is being measured

| Measure | Definition | Why it matters |
|---|---|---|
| Human hours per delivered feature | Accountable-owner time from intent to merged, working feature | The whole claim is that this number is small |
| Correction rate | Proportion of agent proposals requiring material human correction before acceptance | Distinguishes assistance from delegation |
| Escaped defects | Defects found after merge, per delivered feature | Speed that produces rework is not speed |
| Maintenance load | Human hours per month on dependency updates, breakage, and regression repair, once the application runs | Building is the easy half; the claim is about maintaining |
| Cost per feature | Model spend plus infrastructure, per delivered feature | Determines whether the approach survives without subsidy |
| Model class sufficiency | Which model class proved sufficient per task type: deterministic, small local, large local, cloud | Required by the ADR-0024 sovereignty tiers |
| Context re-explanation | How often the same context must be supplied again across sessions | Measures whether the repository genuinely holds the knowledge |

Lines of code, number of agents, commit count, and delivered document volume are
not measures. They were the most visible signals during the Forge self-build and
none of them indicated product value.

## Recording

Each delivered feature records the measures above in the Workplace repository,
alongside the change. Recording is part of delivery, not a separate reporting
activity; a measure that requires its own workflow will not be collected.

Estimates are acceptable where exact figures are impractical, provided they are
labelled as estimates. A missing measure is recorded as missing rather than
omitted silently.

## Review

The measures are reviewed after the first three delivered features, and then
quarterly. Review is a written judgement, not a dashboard reading.

## Falsification

The approach SHALL be judged unsuccessful, and the accountable owner SHALL
reconsider it, if after the first three delivered features any of the following
holds:

- human hours per feature are not materially below a conventional
  human-engineering estimate for the same feature;
- the correction rate is high enough that reviewing proposals costs more
  attention than writing the change would have;
- maintenance load grows with application size at a rate that a single person
  cannot absorb;
- cost per feature exceeds what the intended public-sector market would pay;
- no delivery path works on a locally hosted model, making community
  maintenance impossible; or
- delivery depends on context the accountable owner holds personally and that
  the repository does not capture.

An unsuccessful result is a valid and useful outcome of this experiment. It is
recorded, not reframed.

## Known measurement risks

- The accountable owner is also the sole assessor. Favourable interpretation is
  the default failure mode.
- The first features are the easiest; early figures will flatter the approach.
- Feature size is not comparable between deliveries, so per-feature figures are
  indicative rather than exact.
- Measuring an approach one has invested in is not independent assessment. These
  figures inform judgement; they do not substitute for it.
