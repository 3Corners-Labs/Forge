#!/usr/bin/env python3
"""Fail-closed validator for protected Forge contributor autonomy grants."""

from __future__ import annotations

import argparse
import copy
import fnmatch
import hashlib
import json
import re
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Iterable


SCHEMA_VERSION = "forge.contributor-autonomy-grant.v1"
GRANT_DIRECTORY = ".forge/autonomy/grants"
ZERO_DIGEST = "sha256:" + ("0" * 64)
ALLOWED_ACTIONS = {
    "load-context",
    "create-feature-branch",
    "edit-task-scope",
    "run-verification",
    "repair-in-scope-failures",
    "commit-task-changes",
    "push-feature-branch",
    "open-draft-pull-request",
    "update-pull-request",
    "address-review-feedback",
    "mark-pull-request-ready",
    "monitor-merge-readiness",
}
MANDATORY_HUMAN_GATES = {
    "adr-or-architecture-decision",
    "accepted-boundary-change",
    "scope-expansion",
    "protected-default-branch-landing",
    "release-deploy-or-production-change",
    "production-credentials-or-data",
    "destructive-or-difficult-recovery",
    "control-bypass-or-weakening",
    "accepted-outcome-or-oracle-change",
    "material-product-ambiguity",
    "budget-overrun",
    "indeterminate-authority-or-evidence",
}
REQUIRED_CHECKS = {"validate", "verify", "govulncheck", "codeql"}
PROTECTED_PATHS = (
    "AGENTS.md",
    ".forge/autonomy/**",
    "docs/decisions/**",
    "docs/foundation/**",
    "docs/architecture/**",
)


class GrantError(ValueError):
    """A disclosure-safe grant validation failure."""


def _fail(message: str) -> None:
    raise GrantError(message)


def _exact_keys(value: Any, required: Iterable[str], optional: Iterable[str], label: str) -> dict[str, Any]:
    if not isinstance(value, dict):
        _fail(f"{label} must be an object")
    required_set = set(required)
    allowed = required_set | set(optional)
    keys = set(value)
    if keys - allowed:
        _fail(f"{label} contains unknown fields")
    if required_set - keys:
        _fail(f"{label} is missing required fields")
    return value


def _strings(value: Any, label: str, *, nonempty: bool = True) -> list[str]:
    if not isinstance(value, list) or (nonempty and not value):
        _fail(f"{label} must be a non-empty array")
    if any(not isinstance(item, str) or not item for item in value):
        _fail(f"{label} must contain non-empty strings")
    if len(value) != len(set(value)):
        _fail(f"{label} must contain unique values")
    return value


def _timestamp(value: Any, label: str) -> datetime:
    if not isinstance(value, str) or not value.endswith("Z"):
        _fail(f"{label} must be an RFC3339 UTC timestamp")
    try:
        parsed = datetime.fromisoformat(value[:-1] + "+00:00")
    except ValueError:
        _fail(f"{label} must be an RFC3339 UTC timestamp")
    if parsed.tzinfo is None:
        _fail(f"{label} must include a timezone")
    return parsed.astimezone(timezone.utc)


def canonical_digest(grant: dict[str, Any]) -> str:
    candidate = copy.deepcopy(grant)
    integrity = candidate.get("integrity")
    if not isinstance(integrity, dict):
        _fail("integrity must be an object")
    integrity["canonicalDigest"] = ZERO_DIGEST
    payload = json.dumps(candidate, sort_keys=True, separators=(",", ":")).encode("utf-8")
    return "sha256:" + hashlib.sha256(payload).hexdigest()


def validate_grant_document(grant: Any) -> dict[str, Any]:
    grant = _exact_keys(
        grant,
        {
            "schemaVersion", "id", "version", "status", "issuer", "accountableOwner",
            "issuedAt", "expiresAt", "revocation", "repository", "scope",
            "allowedActions", "mandatoryHumanGates", "destinations", "limits",
            "verification", "stop", "activation", "integrity",
        },
        set(),
        "grant",
    )
    if grant["schemaVersion"] != SCHEMA_VERSION:
        _fail("unsupported grant schema version")
    if not isinstance(grant["id"], str) or not re.fullmatch(r"forge\.autonomy-grant\.[a-z0-9.-]+", grant["id"]):
        _fail("grant id is invalid")
    if not isinstance(grant["version"], str) or not grant["version"]:
        _fail("grant version is invalid")
    if grant["status"] not in {"draft", "reviewed", "authorized", "expired", "revoked"}:
        _fail("grant status is invalid")
    for field in ("issuer", "accountableOwner"):
        if not isinstance(grant[field], str) or not grant[field]:
            _fail(f"{field} is invalid")
    issued_at = _timestamp(grant["issuedAt"], "issuedAt")
    expires_at = _timestamp(grant["expiresAt"], "expiresAt")
    if expires_at <= issued_at:
        _fail("expiresAt must be after issuedAt")

    revocation = _exact_keys(grant["revocation"], {"revoked"}, {"revokedAt", "reason"}, "revocation")
    if not isinstance(revocation["revoked"], bool):
        _fail("revocation.revoked must be boolean")
    if revocation["revoked"]:
        _timestamp(revocation.get("revokedAt"), "revocation.revokedAt")
        if not isinstance(revocation.get("reason"), str) or not revocation["reason"]:
            _fail("a revoked grant requires a reason")
    elif "revokedAt" in revocation or "reason" in revocation:
        _fail("an unrevoked grant cannot contain revocation details")

    repository = _exact_keys(grant["repository"], {"slug", "protectedRef", "protectedBaseRevision"}, set(), "repository")
    if not isinstance(repository["slug"], str) or not re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", repository["slug"]):
        _fail("repository.slug is invalid")
    if repository["protectedRef"] != "refs/remotes/origin/main":
        _fail("repository.protectedRef must be refs/remotes/origin/main")
    if not isinstance(repository["protectedBaseRevision"], str) or not re.fullmatch(r"[a-f0-9]{40}", repository["protectedBaseRevision"]):
        _fail("repository.protectedBaseRevision must be a full commit SHA")

    scope = _exact_keys(grant["scope"], {"milestones", "taskPatterns", "capabilityIds", "pathPatterns", "branchPatterns"}, set(), "scope")
    for field in ("milestones", "taskPatterns", "capabilityIds", "pathPatterns", "branchPatterns"):
        _strings(scope[field], f"scope.{field}")
    if any(pattern in PROTECTED_PATHS for pattern in scope["pathPatterns"]):
        _fail("grant scope cannot include authority-defining paths")
    if any(pattern in {"main", "master", "*"} for pattern in scope["branchPatterns"]):
        _fail("grant scope cannot include a default or unrestricted branch")

    allowed_actions = set(_strings(grant["allowedActions"], "allowedActions"))
    if not allowed_actions <= ALLOWED_ACTIONS:
        _fail("allowedActions contains an unsupported action")
    gates = set(_strings(grant["mandatoryHumanGates"], "mandatoryHumanGates"))
    if gates != MANDATORY_HUMAN_GATES:
        _fail("mandatoryHumanGates must preserve the complete ADR-0011 gate set")

    destinations = grant["destinations"]
    if not isinstance(destinations, list) or not destinations:
        _fail("destinations must be a non-empty array")
    for index, destination in enumerate(destinations):
        item = _exact_keys(destination, {"type", "value"}, set(), f"destinations[{index}]")
        if item["type"] != "github-repository" or item["value"] != repository["slug"]:
            _fail("destinations may contain only the scoped GitHub repository")

    limits = _exact_keys(
        grant["limits"],
        {"maxDurationMinutes", "maxExternalCostMinorUnits", "maxAttemptsPerFailure", "maxConcurrency", "maxChangedFiles", "maxChangedLines"},
        set(),
        "limits",
    )
    limit_minimums = {
        "maxDurationMinutes": 1,
        "maxExternalCostMinorUnits": 0,
        "maxAttemptsPerFailure": 1,
        "maxConcurrency": 1,
        "maxChangedFiles": 1,
        "maxChangedLines": 1,
    }
    for field, minimum in limit_minimums.items():
        if not isinstance(limits[field], int) or isinstance(limits[field], bool) or limits[field] < minimum:
            _fail(f"limits.{field} is invalid")
    if limits["maxConcurrency"] != 1:
        _fail("initial grants must limit concurrency to one contributor")

    verification = _exact_keys(
        grant["verification"],
        {"requiredChecks", "requiredEvidence", "updateCadence", "terminalState"},
        set(),
        "verification",
    )
    if not REQUIRED_CHECKS <= set(_strings(verification["requiredChecks"], "verification.requiredChecks")):
        _fail("verification.requiredChecks omits a required repository check")
    evidence = set(_strings(verification["requiredEvidence"], "verification.requiredEvidence"))
    if not {"grant-id", "protected-revision", "intent", "capabilities", "verification", "rollback", "human-gates"} <= evidence:
        _fail("verification.requiredEvidence is incomplete")
    if verification["updateCadence"] != "material-progress":
        _fail("verification.updateCadence must be material-progress")
    if verification["terminalState"] != "pull-request-merge-ready":
        _fail("verification.terminalState must be pull-request-merge-ready")

    stop = _exact_keys(grant["stop"], {"conditions", "safeRecovery"}, set(), "stop")
    conditions = set(_strings(stop["conditions"], "stop.conditions"))
    if not {"grant-expired", "grant-revoked", "scope-uncertain", "authority-indeterminate", "evidence-unavailable", "budget-exhausted", "terminal-state-reached"} <= conditions:
        _fail("stop.conditions is incomplete")
    if not isinstance(stop["safeRecovery"], str) or not stop["safeRecovery"]:
        _fail("stop.safeRecovery is invalid")

    activation = _exact_keys(
        grant["activation"],
        {"required", "authenticatedHumanRequired", "protectedRevisionRequired", "candidateMayActivate"},
        set(),
        "activation",
    )
    if activation != {
        "required": True,
        "authenticatedHumanRequired": True,
        "protectedRevisionRequired": True,
        "candidateMayActivate": False,
    }:
        _fail("activation controls must remain fail-closed")

    integrity = _exact_keys(grant["integrity"], {"algorithm", "canonicalDigest"}, set(), "integrity")
    if integrity["algorithm"] != "sha256" or not isinstance(integrity["canonicalDigest"], str):
        _fail("integrity metadata is invalid")
    if not re.fullmatch(r"sha256:[a-f0-9]{64}", integrity["canonicalDigest"]):
        _fail("integrity.canonicalDigest is invalid")
    if integrity["canonicalDigest"] != canonical_digest(grant):
        _fail("grant integrity digest does not match canonical content")
    return grant


def _git(repository: Path, *args: str, check: bool = True) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(
        ["git", "-C", str(repository), *args],
        text=True,
        capture_output=True,
        check=False,
    )
    if check and result.returncode != 0:
        _fail("protected repository state could not be verified")
    return result


def _github_slug(remote_url: str) -> str | None:
    match = re.fullmatch(
        r"(?:https://github\.com/|git@github\.com:)([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+?)(?:\.git)?",
        remote_url.strip(),
    )
    return f"{match.group(1)}/{match.group(2)}" if match else None


def load_protected_grant(
    repository: Path,
    grant_id: str,
    protected_revision: str,
    protected_ref: str,
    expected_repository: str | None = None,
) -> dict[str, Any]:
    if not re.fullmatch(r"forge\.autonomy-grant\.[a-z0-9.-]+", grant_id):
        _fail("grant id is invalid")
    if not re.fullmatch(r"[a-f0-9]{40}", protected_revision):
        _fail("protected revision must be a full commit SHA")
    resolved = _git(repository, "rev-parse", "--verify", f"{protected_revision}^{{commit}}").stdout.strip()
    if resolved != protected_revision:
        _fail("protected revision did not resolve exactly")
    remote_result = _git(repository, "ls-remote", "--exit-code", "origin", "refs/heads/main")
    remote_rows = [line.split() for line in remote_result.stdout.splitlines() if line.strip()]
    if len(remote_rows) != 1 or len(remote_rows[0]) != 2 or not re.fullmatch(r"[a-f0-9]{40}", remote_rows[0][0]):
        _fail("protected remote main could not be resolved exactly")
    remote_tip = remote_rows[0][0]
    local_tip = _git(repository, "rev-parse", "--verify", protected_ref).stdout.strip()
    if local_tip != remote_tip:
        _fail("local protected ref does not match remote main")
    if expected_repository is not None:
        remote_url = _git(repository, "remote", "get-url", "origin").stdout.strip()
        actual_repository = _github_slug(remote_url)
        if actual_repository is None or actual_repository.casefold() != expected_repository.casefold():
            _fail("origin does not match the activated GitHub repository")
    if _git(repository, "merge-base", "--is-ancestor", protected_revision, protected_ref, check=False).returncode != 0:
        _fail("protected revision is not reachable from the protected ref")
    relative_path = f"{GRANT_DIRECTORY}/{grant_id}.json"
    raw = _git(repository, "show", f"{protected_revision}:{relative_path}").stdout
    try:
        grant = json.loads(raw)
    except json.JSONDecodeError:
        _fail("protected grant is not valid JSON")
    if grant.get("id") != grant_id:
        _fail("protected grant id does not match its canonical path")
    return validate_grant_document(grant)


def _matches(value: str, patterns: list[str]) -> bool:
    return any(fnmatch.fnmatchcase(value, pattern) for pattern in patterns)


def validate_activation(grant: dict[str, Any], args: argparse.Namespace, now: datetime) -> None:
    if grant["status"] != "authorized":
        _fail("grant is not authorized")
    if grant["revocation"]["revoked"]:
        _fail("grant is revoked")
    if now < _timestamp(grant["issuedAt"], "issuedAt") or now >= _timestamp(grant["expiresAt"], "expiresAt"):
        _fail("grant is not within its validity window")
    repository = grant["repository"]
    if args.repository != repository["slug"] or args.protected_ref != repository["protectedRef"]:
        _fail("activation repository does not match the grant")
    base = repository["protectedBaseRevision"]
    if _git(args.repo_root, "merge-base", "--is-ancestor", base, args.protected_revision, check=False).returncode != 0:
        _fail("grant base is not an ancestor of the protected revision")
    scope = grant["scope"]
    if not _matches(args.task, scope["taskPatterns"]):
        _fail("activation task is outside grant scope")
    if not _matches(args.branch, scope["branchPatterns"]):
        _fail("activation branch is outside grant scope")
    if not args.capability or any(item not in scope["capabilityIds"] for item in args.capability):
        _fail("activation capabilities are outside grant scope")
    if not args.path or any(not _matches(item, scope["pathPatterns"]) for item in args.path):
        _fail("activation paths are outside grant scope")
    if not args.action or any(item not in grant["allowedActions"] for item in args.action):
        _fail("activation actions are outside grant scope")
    if args.terminal_state != grant["verification"]["terminalState"]:
        _fail("activation terminal state does not match the grant")
    observed = {
        "maxDurationMinutes": args.elapsed_minutes,
        "maxExternalCostMinorUnits": args.external_cost_minor_units,
        "maxAttemptsPerFailure": args.attempts_per_failure,
        "maxConcurrency": args.concurrency,
        "maxChangedFiles": args.changed_files,
        "maxChangedLines": args.changed_lines,
    }
    for field, value in observed.items():
        if value < 0 or value > grant["limits"][field]:
            _fail(f"activation exceeds {field}")


def _parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="Validate one task-scoped bounded delivery activation.")
    parser.add_argument("--repo-root", type=Path, default=Path.cwd())
    parser.add_argument("--grant-id", required=True)
    parser.add_argument("--protected-revision", required=True)
    parser.add_argument("--protected-ref", default="refs/remotes/origin/main")
    parser.add_argument("--repository", required=True)
    parser.add_argument("--task", required=True)
    parser.add_argument("--branch", required=True)
    parser.add_argument("--capability", action="append", default=[])
    parser.add_argument("--path", action="append", default=[])
    parser.add_argument("--action", action="append", default=[])
    parser.add_argument("--terminal-state", default="pull-request-merge-ready")
    parser.add_argument("--at", help="RFC3339 UTC validation time; defaults to now")
    parser.add_argument("--elapsed-minutes", type=int, default=0)
    parser.add_argument("--external-cost-minor-units", type=int, default=0)
    parser.add_argument("--attempts-per-failure", type=int, default=0)
    parser.add_argument("--concurrency", type=int, default=1)
    parser.add_argument("--changed-files", type=int, default=0)
    parser.add_argument("--changed-lines", type=int, default=0)
    return parser


def main() -> int:
    args = _parser().parse_args()
    try:
        now = _timestamp(args.at, "at") if args.at else datetime.now(timezone.utc)
        grant = load_protected_grant(
            args.repo_root.resolve(),
            args.grant_id,
            args.protected_revision,
            args.protected_ref,
            args.repository,
        )
        validate_activation(grant, args, now)
    except GrantError as error:
        print(json.dumps({"schemaVersion": 1, "state": "denied", "reason": str(error)}, sort_keys=True))
        return 20
    print(json.dumps({
        "schemaVersion": 1,
        "state": "eligible-for-human-activation",
        "grantId": grant["id"],
        "grantDigest": grant["integrity"]["canonicalDigest"],
        "protectedRevision": args.protected_revision,
        "task": args.task,
        "branch": args.branch,
        "terminalState": args.terminal_state,
    }, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
