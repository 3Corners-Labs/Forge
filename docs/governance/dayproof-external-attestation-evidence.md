# Day-proof: digest-bound approval through external attestation

Status: Draft
Owner: Forge Engineering and Architecture
Verified: 2026-07-25

## Question

Market research recorded that no shipped agent tool binds a human approval to an
artifact digest, and that roughly 70–85% of Forge's trust chain already exists as
free standard tooling. This proof tests the practical consequence:

> Can `cosign` bind a Forge Approval Record to the exact digest Forge already
> computes, fully offline, and does that make any Forge module redundant?

Fully offline means: local key, no transparency log, no Fulcio, no OIDC, no
timestamp authority, no network. The target customer runs air-gapped.

## Method

Subject: commit `08097f9` ("rebaseline Forge as a bounded means"), a real change
of 11 files.

1. Compute the candidate tree digest by importing Forge's own
   `workercontract.CandidateTreeDigest` — the production algorithm, not a
   reimplementation.
2. Derive a `HumanApproval` predicate from the existing
   `schemas/forge/trust/v1/approval-record-v2.schema.json`. No new format.
3. Sign with `cosign attest-blob` using a locally generated key.
4. Verify offline.
5. Attempt to break it two ways.

Tooling: cosign v3.1.1, Go 1.26.5. Keys and artefacts were written to a
scratchpad outside the repository; no key material was committed.

## Result: the mechanism works

**Tree digest** (Forge's own algorithm):

~~~text
sha256:033f2ecda3f4302e292a6f2d59bc275ad813294366b99d5917a8baa43d066bef
~~~

Independently reproduced from the file list in Python, byte-identical. The
algorithm is not Go-specific and can be implemented by any verifier.

**Attestation**: a Sigstore bundle v0.3 containing a DSSE envelope, payload type
`application/vnd.in-toto+json`, in-toto Statement v0.1, predicate type
`https://3corners.nl/attestations/HumanApproval/v1`, carrying the Approval
Record fields — `action`, `subjectDigest`, `decision`, `approverId`,
`environment`, `riskClass`, `decidedAt`, `expiresAt`, and the 11 file digests.

**Verification**: `Verified OK`, offline.

### Break test A — one byte of source

Flipping one character in the content digest of one file changes the tree digest
to `sha256:d9cd86e0…`. The approved digest no longer matches the actual tree, so
publication must be refused. The approval is bound to exact content, not to a
branch, a commit message, or a reviewer's intent.

### Break test B — forging the approval

Editing the signed payload to name a different approver and lower the risk class,
leaving the signature untouched, produces:

~~~text
Error: failed to verify signature: could not verify envelope:
accepted signatures do not match threshold, Found: 0, Expected 1
~~~

The unmodified bundle still verifies. An approval cannot be rewritten after the
fact without the signing key.

## What this does and does not replace

This is where the earlier estimate was too generous. Being precise:

| Forge module | Lines | Disposition |
|---|---|---|
| `internal/contracts/trust` integrity binding | part of 1257 | **Can delegate** the signing and verification mechanics to DSSE/cosign; the schemas stay |
| `internal/evidence` | 529 | **Partly.** Per-record signing and tamper detection are covered. The append-only stream, fork detection, and artefact promotion are not |
| `internal/verifier` | 160 | **Partly**, on the integrity-checking side |
| `internal/verification` | 451 | **Keep.** cosign verifies signatures, not whether the change is correct. Independent assurance-plan derivation and execution is a different problem and remains unsolved by the market |
| `internal/policy` | 2172 | **Keep the rules, replace the engine** with OPA later. Not tested here |

Honest total directly displaced by this proof: on the order of **500–700 lines**,
not the thousands claimed earlier in the session. The value is not deletion. It
is that the retained code now sits behind a standard, externally verifiable
format instead of a bespoke one.

## Air-gap friction found

Every step that needed a workaround, because this list is the beginning of the
air-gapped bundle the target customer needs and nobody publishes:

1. **`--tlog-upload=false` is deprecated in cosign 3.x** and errors out when a
   signing config is in play.
2. **The error message's suggested fix requires the internet.** It instructs the
   operator to `curl` a signing config from `raw.githubusercontent.com` and strip
   the Rekor URLs with `jq`. In an air-gapped environment that instruction is
   unusable.
3. **The working offline route is undocumented as such:**
   `cosign signing-config create --no-default-fulcio --no-default-rekor
   --no-default-oidc --no-default-tsa --out signing-config.json` produces a local
   empty config that permits offline signing. Nothing points an operator here.
4. **`--output-attestation` was removed** in favour of `--bundle` with the new
   bundle format. Existing instructions and scripts break silently on upgrade.
5. **Offline verification requires `--insecure-ignore-tlog`**, which prints
   `WARNING: Skipping tlog verification is an insecure practice…`. Technically
   accurate, but presenting that warning to a public-sector auditor undermines
   the evidence it accompanies.

Item 5 is not cosmetic. The proper resolutions are a self-hosted Rekor v2
(`rekor-server-posix`, documented as having no cloud dependencies) or an
eIDAS-qualified European RFC 3161 timestamp authority. The latter is worth
evaluating on its own merits: a qualified EU timestamp carries legal evidentiary
weight in the Netherlands that a US transparency log does not.

## Follow-on: external RFC 3161 timestamping

The proof was extended to test whether an external timestamp authority can
replace the Sigstore transparency log, as the route to an eIDAS-qualified
timestamp. A public, **non-qualified** authority was used to prove the mechanism;
it carries no Article 41 presumption.

**The mechanism works.** Signing with `--signing-config` carrying a `--tsa`
entry, and verifying against a `--trusted-root` containing only that authority,
produces a bundle with one RFC 3161 token and **zero transparency-log entries**,
verifying offline. Tampering with the approver in the timestamped bundle still
fails signature verification.

Two results matter more than the success:

1. **Only a hash leaves the environment.** The timestamp request is 70 bytes
   containing the SHA-256 and a nonce — no filename, no content, no metadata.
   This is the argument that permits timestamping across an air-gap boundary.
2. **The `--insecure-ignore-tlog` warning cannot be removed.** An earlier claim
   in this session that a timestamp would eliminate it was wrong. Verifying
   without the flag fails with `not enough verified log entries from transparency
   log: 0 < 1`, **even when the trusted root declares no transparency logs at
   all**. cosign does not treat an empty `tlogs` set as "no log required". A
   valid qualified timestamp therefore still produces a warning about missing
   transparency — the opposite of the intended effect in front of an auditor.

Item 2 is a cosign limitation, not a design constraint of the approach. It is
worth reporting upstream; the trusted root already carries enough information to
infer that no log is expected.

Additional air-gap friction, extending the list above:

6. `--timestamp-server-url` and `--rfc3161-timestamp` were removed in cosign 3.x.
   The published Sigstore documentation still describes them.
7. The timestamp authority's certificate chain must be supplied at **signing**
   time as well as verification time, because cosign verifies the token as it
   signs. An air-gapped signer needs the chain provisioned in advance.

### Qualified provider selection

No Netherlands-established qualified provider sells a documented, publicly priced
qualified RFC 3161 endpoint, so the service must be procured from another Member
State. Article 41(3) provides cross-border recognition. Candidates recorded for
evaluation: Aruba PEC (IT, self-service, 50 timestamps for EUR 13.50), SK ID
Solutions (EE, contract plus IP allowlist, EUR 0.036 per timestamp falling to
EUR 0.006 at volume), Camerfirma (ES), InfoCert (IT).

Qualified status for every candidate is provider-published and SHALL be confirmed
against the authoritative EU Trusted List before contracting.

## Conclusion

The one capability the market does not ship — a human approval cryptographically
bound to an exact artifact digest — is achievable today, offline, with standard
tooling, in an afternoon. It is a missing bundle, not a missing technology.

That confirms the direction: assemble on standard tooling, keep the contract and
the parts of Forge that genuinely have no market equivalent, and treat the
air-gapped, verified bundle as a deliverable in its own right.

## Not established

- Independent verification the agent cannot influence. Break test B proves an
  approval cannot be forged; it says nothing about whether the reviewer was
  independent of the author. That is the second gap and it is untested.
- Behaviour with a hardware or KMS-held key.
- Self-hosted Rekor, and whether an eIDAS timestamp authority integrates cleanly.
- Whether OPA can express the existing `internal/policy` rules without loss.
- Long-term key custody, rotation, and revocation.

## Reproduction

The proof harness is `dist/dayproof/main.go`, which is deliberately inside the
gitignored `dist/` directory: it is throwaway evidence code, not a Forge
capability. It imports `internal/contracts/worker` so the digest under test is
the production one.
