from __future__ import annotations

import argparse
import copy
import json
import subprocess
import tempfile
import unittest
from datetime import datetime, timezone
from pathlib import Path

from scripts.validate_autonomy_grant import (
    GrantError,
    ZERO_DIGEST,
    canonical_digest,
    load_protected_grant,
    validate_activation,
    validate_grant_document,
)


ROOT = Path(__file__).resolve().parents[1]
CANONICAL_GRANT = ROOT / ".forge/autonomy/grants/forge.autonomy-grant.m1-bounded-delivery-v1.json"


def signed(grant: dict) -> dict:
    candidate = copy.deepcopy(grant)
    candidate["integrity"]["canonicalDigest"] = ZERO_DIGEST
    candidate["integrity"]["canonicalDigest"] = canonical_digest(candidate)
    return candidate


class AutonomyGrantTests(unittest.TestCase):
    def setUp(self) -> None:
        self.grant = json.loads(CANONICAL_GRANT.read_text(encoding="utf-8"))

    def test_canonical_grant_is_closed_and_integrity_bound(self) -> None:
        self.assertEqual(validate_grant_document(self.grant)["status"], "authorized")
        altered = copy.deepcopy(self.grant)
        altered["limits"]["maxConcurrency"] = 2
        with self.assertRaisesRegex(GrantError, "concurrency"):
            validate_grant_document(signed(altered))
        altered = copy.deepcopy(self.grant)
        altered["unexpected"] = True
        with self.assertRaisesRegex(GrantError, "unknown fields"):
            validate_grant_document(altered)

    def test_mandatory_gates_and_authority_paths_cannot_be_removed(self) -> None:
        altered = copy.deepcopy(self.grant)
        altered["mandatoryHumanGates"].remove("protected-default-branch-landing")
        with self.assertRaisesRegex(GrantError, "complete ADR-0011 gate set"):
            validate_grant_document(signed(altered))
        altered = copy.deepcopy(self.grant)
        altered["scope"]["pathPatterns"].append("AGENTS.md")
        with self.assertRaisesRegex(GrantError, "authority-defining paths"):
            validate_grant_document(signed(altered))

        self.assertIn("architecture/module-boundaries.json", self.grant["scope"]["pathPatterns"])
        altered = copy.deepcopy(self.grant)
        altered["scope"]["pathPatterns"].append("architecture/**")
        with self.assertRaisesRegex(GrantError, "broadly include architecture"):
            validate_grant_document(signed(altered))

    def test_integrity_change_fails_closed(self) -> None:
        altered = copy.deepcopy(self.grant)
        altered["limits"]["maxChangedFiles"] += 1
        with self.assertRaisesRegex(GrantError, "integrity digest"):
            validate_grant_document(altered)


class ProtectedActivationTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary = tempfile.TemporaryDirectory()
        self.repository = Path(self.temporary.name)
        self._git("init", "-b", "main")
        self._git("config", "user.name", "Forge Test")
        self._git("config", "user.email", "forge-test@example.invalid")
        (self.repository / "README.md").write_text("base\n", encoding="utf-8")
        self._git("add", "README.md")
        self._git("commit", "-m", "base")
        base = self._git("rev-parse", "HEAD").stdout.strip()

        grant = json.loads(CANONICAL_GRANT.read_text(encoding="utf-8"))
        grant["repository"]["protectedBaseRevision"] = base
        grant = signed(grant)
        directory = self.repository / ".forge/autonomy/grants"
        directory.mkdir(parents=True)
        self.path = directory / f"{grant['id']}.json"
        self.path.write_text(json.dumps(grant, indent=2) + "\n", encoding="utf-8")
        self._git("add", str(self.path.relative_to(self.repository)))
        self._git("commit", "-m", "protected grant")
        self.protected_revision = self._git("rev-parse", "HEAD").stdout.strip()
        self._git("remote", "add", "origin", str(self.repository))
        self._git("update-ref", "refs/remotes/origin/main", self.protected_revision)
        self.protected_grant = grant
        self.args = argparse.Namespace(
            repo_root=self.repository,
            grant_id=grant["id"],
            protected_revision=self.protected_revision,
            protected_ref="refs/remotes/origin/main",
            repository="3Corners-Labs/Forge",
            task="m1-iteration-4-policy-decision",
            branch="codex/m1-iteration-4-policy-decision",
            capability=["forge.capability.policy"],
            path=["internal/delivery/policy.go"],
            action=["run-verification", "edit-task-scope"],
            terminal_state="pull-request-merge-ready",
            elapsed_minutes=1,
            external_cost_minor_units=0,
            attempts_per_failure=0,
            concurrency=1,
            changed_files=1,
            changed_lines=10,
        )
        self.now = datetime(2026, 8, 1, tzinfo=timezone.utc)

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def _git(self, *args: str) -> subprocess.CompletedProcess[str]:
        return subprocess.run(
            ["git", "-C", str(self.repository), *args],
            text=True,
            capture_output=True,
            check=True,
        )

    def test_candidate_worktree_cannot_widen_protected_grant(self) -> None:
        self._git("switch", "-c", "candidate")
        candidate = copy.deepcopy(self.protected_grant)
        candidate["allowedActions"].append("merge-protected-branch")
        self.path.write_text(json.dumps(signed(candidate), indent=2) + "\n", encoding="utf-8")

        loaded = load_protected_grant(
            self.repository,
            self.protected_grant["id"],
            self.protected_revision,
            "refs/remotes/origin/main",
        )

        self.assertNotIn("merge-protected-branch", loaded["allowedActions"])
        validate_activation(loaded, self.args, self.now)

    def test_valid_activation_is_eligible(self) -> None:
        loaded = load_protected_grant(
            self.repository,
            self.protected_grant["id"],
            self.protected_revision,
            "refs/remotes/origin/main",
        )
        validate_activation(loaded, self.args, self.now)

    def test_wrong_scope_lifecycle_and_budget_fail_closed(self) -> None:
        cases = {
            "repository": ("repository", "another/Repository"),
            "task": ("task", "m2-release"),
            "branch": ("branch", "main"),
            "capability": ("capability", ["forge.capability.release"]),
            "path": ("path", ["AGENTS.md"]),
            "action": ("action", ["merge-protected-branch"]),
            "terminal": ("terminal_state", "merged"),
            "budget": ("changed_files", 101),
        }
        for label, (field, value) in cases.items():
            with self.subTest(label=label):
                args = copy.deepcopy(self.args)
                setattr(args, field, value)
                with self.assertRaises(GrantError):
                    validate_activation(self.protected_grant, args, self.now)

        with self.assertRaisesRegex(GrantError, "validity window"):
            validate_activation(self.protected_grant, self.args, datetime(2026, 10, 1, tzinfo=timezone.utc))
        revoked = copy.deepcopy(self.protected_grant)
        revoked["status"] = "revoked"
        revoked["revocation"] = {
            "revoked": True,
            "revokedAt": "2026-07-31T00:00:00Z",
            "reason": "test revocation",
        }
        revoked = validate_grant_document(signed(revoked))
        with self.assertRaisesRegex(GrantError, "not authorized"):
            validate_activation(revoked, self.args, self.now)

    def test_unknown_or_unreachable_protected_revision_fails_closed(self) -> None:
        with self.assertRaises(GrantError):
            load_protected_grant(
                self.repository,
                self.protected_grant["id"],
                "f" * 40,
                "refs/remotes/origin/main",
            )

    def test_stale_or_substituted_local_protected_ref_fails_closed(self) -> None:
        self._git("update-ref", "refs/remotes/origin/main", self.protected_grant["repository"]["protectedBaseRevision"])
        with self.assertRaisesRegex(GrantError, "does not match remote main"):
            load_protected_grant(
                self.repository,
                self.protected_grant["id"],
                self.protected_revision,
                "refs/remotes/origin/main",
            )

    def test_origin_identity_must_match_activated_github_repository(self) -> None:
        with self.assertRaisesRegex(GrantError, "origin does not match"):
            load_protected_grant(
                self.repository,
                self.protected_grant["id"],
                self.protected_revision,
                "refs/remotes/origin/main",
                "3Corners-Labs/Forge",
            )


if __name__ == "__main__":
    unittest.main()
