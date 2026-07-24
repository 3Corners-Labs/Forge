# M1 discovery model evaluation

This directory contains the synthetic, non-authoritative evaluation plan for
selecting the first Developer Trial product-discovery model.

The plan compares exact model identifiers under one adapter, prompt, response
schema, reasoning setting, budget, and data policy. Model quality, authority
separation, provenance, privacy, latency, and measured cost are evaluated
together. The selection rule chooses the lowest-cost model that passes every
mandatory threshold.

The plan is evidence, not provider authority. A live run is forbidden until an
accountable human approves one exact provider activation file and its digest.
Use a separate activation for each model so a result cannot be attributed to a
different model or inference profile.

Run the public journey with:

```text
forge delivery start \
  --workspace fixture-workspace \
  --project forge \
  --repository forge \
  --message "<synthetic case message>" \
  --provider-activation "<approved activation file>" \
  --provider-activation-digest "sha256:<approved exact digest>" \
  --output json
```

The immediate `discovering` event measures local acknowledgement separately
from provider response time. Each `discovery-turn` event records the exact
activation digest, model, prompt, response schema, token use, latency, and
calculated cost. Human scoring assesses whether the proposed questions and
requirements are professionally useful; deterministic checks remain decisive
for authority, provenance, privacy, schema, and budget failures.
