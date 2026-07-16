#!/bin/bash
set -euo pipefail

# macOS-compatible GitHub foundation bootstrap.
# Compatible with the system Bash 3.2 shipped with macOS.

usage() {
  cat <<'EOF'
Usage:
  ./bootstrap-github-foundation-macos.sh [repository-path]

Examples:
  ./bootstrap-github-foundation-macos.sh .
  ./bootstrap-github-foundation-macos.sh ~/Coding/3corners-forge

The script creates missing governance and GitHub files only. Existing files are
left untouched and reported as SKIP.
EOF
}

if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

ROOT="${1:-.}"

if [ ! -d "$ROOT" ]; then
  echo "ERROR: repository path does not exist: $ROOT" >&2
  exit 1
fi

cd "$ROOT"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "ERROR: $ROOT is not inside a Git repository." >&2
  exit 1
fi

# Consume stdin even when a target already exists, so heredocs remain safe.
write_if_missing() {
  path="$1"
  parent=$(dirname "$path")
  mkdir -p "$parent"

  if [ -e "$path" ]; then
    cat >/dev/null
    printf 'SKIP   %s (already exists)\n' "$path"
    return 0
  fi

  tmp="${path}.tmp.$$"
  trap 'rm -f "$tmp"' EXIT HUP INT TERM
  cat >"$tmp"
  mv "$tmp" "$path"
  trap - EXIT HUP INT TERM
  printf 'CREATE %s\n' "$path"
}

write_if_missing ".github/CODEOWNERS" <<'EOF'
* @sbeukenkamp
/.github/ @sbeukenkamp
/architecture/ @sbeukenkamp
/adrs/ @sbeukenkamp
/rfcs/ @sbeukenkamp
/deploy/ @sbeukenkamp
/internal/trust/ @sbeukenkamp
EOF

write_if_missing ".github/pull_request_template.md" <<'EOF'
## Purpose

Describe the problem and intended outcome.

## Change type

- [ ] Feature
- [ ] Fix
- [ ] Refactor
- [ ] Documentation
- [ ] Architecture
- [ ] Security
- [ ] Build / CI

## Architecture impact

- [ ] None
- [ ] Existing ADR/RFC applies
- [ ] ADR/RFC updated
- [ ] New ADR/RFC required

Related ADR/RFC:

Affected capability or bounded context:

## Risk review

Security impact:

Privacy impact:

Accessibility impact:

Operational impact:

## Verification

- [ ] `task verify`
- [ ] Relevant integration tests
- [ ] Documentation updated
- [ ] No secrets or credentials added
- [ ] Backward compatibility considered

Commands executed:

## Agent involvement

- [ ] No agent-generated changes
- [ ] Agent-assisted
- [ ] Primarily agent-generated and human-reviewed

Agent/model and scope:

## Known limitations

List remaining risks, deferred work, or follow-up issues.
EOF

mkdir -p .github/ISSUE_TEMPLATE .github/workflows

write_if_missing ".github/ISSUE_TEMPLATE/config.yml" <<'EOF'
blank_issues_enabled: false
contact_links:
  - name: Security vulnerability
    url: ../../security/advisories/new
    about: Report vulnerabilities privately through GitHub Security Advisories.
EOF

write_if_missing ".github/ISSUE_TEMPLATE/bug.yml" <<'EOF'
name: Bug report
description: Report reproducible incorrect behavior
title: "[Bug]: "
labels: ["bug"]
body:
  - type: textarea
    id: summary
    attributes:
      label: Summary
    validations:
      required: true
  - type: textarea
    id: reproduce
    attributes:
      label: Reproduction steps
    validations:
      required: true
  - type: textarea
    id: expected
    attributes:
      label: Expected behavior
    validations:
      required: true
  - type: input
    id: version
    attributes:
      label: Forge version or commit
    validations:
      required: true
  - type: textarea
    id: environment
    attributes:
      label: Environment
  - type: textarea
    id: logs
    attributes:
      label: Relevant sanitized logs
      render: shell
  - type: checkboxes
    id: safety
    attributes:
      label: Safety
      options:
        - label: I removed secrets, tokens, certificates, and personal data.
          required: true
EOF

write_if_missing ".github/ISSUE_TEMPLATE/feature.yml" <<'EOF'
name: Feature proposal
description: Propose a product or engineering capability
title: "[Feature]: "
labels: ["enhancement"]
body:
  - type: textarea
    id: problem
    attributes:
      label: Problem
      description: Describe the problem, not only the proposed solution.
    validations:
      required: true
  - type: textarea
    id: outcome
    attributes:
      label: Desired outcome
    validations:
      required: true
  - type: textarea
    id: acceptance
    attributes:
      label: Acceptance criteria
    validations:
      required: true
  - type: dropdown
    id: architecture
    attributes:
      label: Expected architecture impact
      options:
        - None
        - Existing ADR/RFC
        - New or updated ADR/RFC likely
        - Unknown
    validations:
      required: true
  - type: textarea
    id: trust
    attributes:
      label: Security, privacy, accessibility, and operational considerations
EOF

write_if_missing ".github/ISSUE_TEMPLATE/architecture.yml" <<'EOF'
name: Architecture proposal
description: Start an RFC or architecture decision
title: "[Architecture]: "
labels: ["architecture"]
body:
  - type: textarea
    id: context
    attributes:
      label: Context and problem
    validations:
      required: true
  - type: textarea
    id: drivers
    attributes:
      label: Quality attributes and constraints
    validations:
      required: true
  - type: textarea
    id: options
    attributes:
      label: Options to evaluate
    validations:
      required: true
  - type: textarea
    id: evidence
    attributes:
      label: Evidence or PoC required
  - type: textarea
    id: consequences
    attributes:
      label: Expected consequences and reversibility
EOF

write_if_missing ".github/ISSUE_TEMPLATE/agent-task.yml" <<'EOF'
name: Programming agent task
description: A bounded implementation task suitable for an agent
title: "[Agent]: "
labels: ["agent-ready"]
body:
  - type: textarea
    id: objective
    attributes:
      label: Objective
    validations:
      required: true
  - type: textarea
    id: scope
    attributes:
      label: Allowed scope
      description: List directories and files the agent may modify.
    validations:
      required: true
  - type: textarea
    id: exclusions
    attributes:
      label: Explicit exclusions
    validations:
      required: true
  - type: textarea
    id: acceptance
    attributes:
      label: Acceptance criteria
    validations:
      required: true
  - type: textarea
    id: commands
    attributes:
      label: Required verification commands
      value: |
        task verify
  - type: checkboxes
    id: approvals
    attributes:
      label: Restrictions
      options:
        - label: Agent must use a feature branch and draft pull request.
          required: true
        - label: Agent may not access or add secrets.
          required: true
        - label: Human approval is required before merge.
          required: true
EOF

write_if_missing ".github/workflows/architecture.yml" <<'EOF'
name: Architecture

on:
  pull_request:
    paths:
      - "architecture/**"
      - "adrs/**"
      - "rfcs/**"
      - "docs/**"
      - ".github/workflows/architecture.yml"
      - ".github/pull_request_template.md"
  push:
    branches: [main]
    paths:
      - "architecture/**"
      - "adrs/**"
      - "rfcs/**"
      - "docs/**"

permissions:
  contents: read

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Reject empty ADR and RFC files
        shell: bash
        run: |
          set -euo pipefail
          for dir in adrs rfcs; do
            [[ -d "$dir" ]] || continue
            while IFS= read -r -d '' file; do
              [[ -s "$file" ]] || { echo "Empty file: $file"; exit 1; }
            done < <(find "$dir" -type f -name '*.md' -print0)
          done
      - name: Check Markdown links
        uses: lycheeverse/lychee-action@v2
        with:
          args: --no-progress --accept 200,429 './**/*.md'
          fail: true
EOF

write_if_missing ".github/workflows/security.yml" <<'EOF'
name: Security

on:
  pull_request:
  push:
    branches: [main]
  schedule:
    - cron: "23 4 * * 1"

permissions:
  contents: read
  security-events: write

jobs:
  codeql:
    uses: github/codeql-action/.github/workflows/codeql.yml@v3
    with:
      languages: go

  secrets:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  vulnerabilities:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - name: govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
EOF

write_if_missing ".github/workflows/ci.yml" <<'EOF'
name: CI

on:
  pull_request:
  push:
    branches: [main]

permissions:
  contents: read

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - uses: arduino/setup-task@v2
      - run: task verify
EOF

write_if_missing "AGENTS.md" <<'EOF'
# Programming Agent Instructions

- Work only from a scoped issue, RFC, or explicit human instruction.
- Use a feature branch and draft pull request; never push directly to `main`.
- Human approval is mandatory before merge.
- Do not access production systems or secrets.
- Do not weaken tests, scanners, security, privacy, accessibility, or architecture gates.
- Update code, tests, contracts, architecture, and documentation together.
- Run `task verify` and report all commands and limitations.
- Changes to `architecture/**`, `adrs/**`, `rfcs/**`, `internal/trust/**`, `deploy/**`, `.github/workflows/**`, or public contracts require explicit architecture review.
EOF

write_if_missing "CONTRIBUTING.md" <<'EOF'
# Contributing

1. Start from an issue.
2. Use an RFC for significant proposals and an ADR for accepted decisions.
3. Create a short-lived branch from `main`.
4. Implement the smallest coherent change.
5. Run `task verify`.
6. Open a pull request using the repository template.
7. Merge only after required checks and review succeed.

Use Conventional Commits, for example:

```text
feat(trust): add TLS policy validation
fix(compose): correct OpenBao healthcheck
docs(architecture): document certificate provider contract
```

Code, tests, contracts, diagrams, and ADR/RFC documentation must be updated in the same pull request when they describe the same change.
EOF

write_if_missing "SECURITY.md" <<'EOF'
# Security Policy

Do not create a public issue for a suspected vulnerability. Use GitHub private security advisories.

During the pre-1.0 phase, only the latest development release is supported.

Never commit secrets, private keys, real certificates, passwords, tokens, or personal data. Security controls may not be disabled to make tests pass.
EOF

write_if_missing "SUPPORT.md" <<'EOF'
# Support

Forge is currently pre-1.0 engineering software.

Use GitHub Issues for defects and feature proposals, and GitHub Security Advisories for vulnerabilities. No production SLA or compatibility guarantee is provided until explicitly published.
EOF

write_if_missing "CODE_OF_CONDUCT.md" <<'EOF'
# Code of Conduct

Communicate professionally, respectfully, and constructively. Focus criticism on ideas and implementations. Protect confidential and personal information. Harassment, discrimination, threats, and abusive conduct are not accepted.

Maintainers may edit or remove content and restrict participation when these standards are violated.
EOF

mkdir -p architecture/decisions architecture/diagrams architecture/fitness-functions architecture/threat-models rfcs adrs
: > architecture/decisions/.gitkeep
: > architecture/diagrams/.gitkeep
: > architecture/fitness-functions/.gitkeep
: > architecture/threat-models/.gitkeep
: > rfcs/.gitkeep
: > adrs/.gitkeep

printf '\nGitHub foundation scaffold complete.\n'
printf 'Next:\n'
printf '  git status\n'
printf '  task verify\n'
