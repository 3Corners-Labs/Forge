# ADR-0026: License Forge under Apache-2.0 and accept contributions under the DCO

Status: Accepted
Owner: Forge Architecture
Date: 2026-10-29

## Context

The Forge repository became public on 2026-10-29. Its `LICENSE` file named the
Apache License 2.0 but held only the short header notice, not the licence text,
so GitHub and licence scanners reported it as "Other". The repository had no
`NOTICE` file, no licence field in the Workbench package manifest, and no rule
for the origin of outside contributions.

Every current dependency is permissively licensed (MIT, BSD, or Apache-2.0):
the Go modules in `go.mod` and the Workbench web client in
`internal/experience/workbench/ui/package.json`. The Developer Trial MinIO image
(AGPL-3.0) runs only as a separate container in development and CI; Forge does
not link to or distribute it.

## Decision

Forge is licensed under the Apache License, Version 2.0, for all code,
schemas, and documentation in this repository:

- `LICENSE` holds the unmodified Apache-2.0 text;
- `NOTICE` holds the project's copyright notice, as Apache-2.0 §4(d) expects;
- package manifests declare the SPDX identifier `Apache-2.0`.

Contributions are accepted inbound under the same licence (Apache-2.0 §5) and
SHALL carry a Developer Certificate of Origin 1.1 sign-off (`DCO`): every
commit by a human author has a `Signed-off-by:` trailer with the author's
email address (`git commit -s`). The `DCO` workflow checks this on every pull
request. Bot accounts (`*[bot]`) cannot certify origin and are skipped; their
changes are reviewed as usual.

## Alternatives considered

- **EUPL-1.2.** The preferred licence for Dutch and EU public-sector code. Its
  weak copyleft would oblige anyone who modifies and distributes or hosts Forge
  to share the changes. Not selected: the repository already declared
  Apache-2.0, and Forge is a 3Corners means, not public-sector code.
- **MPL-2.0.** File-level copyleft; uncommon in the Go ecosystem. Not selected.
- **Contributor Licence Agreement.** Would give 3Corners wider rights, such as
  relicensing, at the cost of paperwork and a signing service, and it deters
  contributors. Not selected.
- **Inbound=outbound without sign-off.** Least effort, but leaves no per-commit
  record that the contributor had the right to submit the work. Not selected.

## Consequences

- GitHub and licence scanners detect Apache-2.0.
- Contributors must sign off their commits; the DCO check fails otherwise and
  names the commit to amend.
- A sign-off records the contributor's name and email address in the public
  history indefinitely. DCO clause (d) states this, and contributors choose
  which address to use (a GitHub no-reply address is sufficient). No other
  personal data is collected.
- Relicensing later needs agreement from every contributor, because no CLA
  grants 3Corners that right.
- Binary distributions that embed the Workbench client SHALL carry the
  third-party notices of its bundled dependencies; release packaging does not
  yet generate them.

## Validation

- GitHub shows "Apache-2.0" for the repository.
- The `DCO` workflow fails a pull request containing an unsigned human commit.
- A new dependency with a licence other than MIT, BSD, ISC, or Apache-2.0
  requires review against this ADR before it is added.

## Supersession

None.
