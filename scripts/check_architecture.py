#!/usr/bin/env python3
"""Dependency-free architecture-as-code checks for the Forge repository."""

from __future__ import annotations

import hashlib
import json
import re
import sys
from pathlib import Path
from urllib.parse import unquote, urlparse


ROOT = Path(__file__).resolve().parents[1]
ALLOWED_DOCUMENT_STATUSES = {
    "Draft",
    "Proposed",
    "Accepted",
    "Superseded",
    "Deprecated",
    "Archive",
}
METADATA_DIRECTORIES = (
    "docs/foundation",
    "docs/architecture",
    "docs/governance",
    "docs/decisions",
    "docs/reference",
    "docs/reference-applications",
    "docs/context",
)
NAVIGATION_DIRECTORIES = METADATA_DIRECTORIES
PRINCIPLES_SOURCE = ROOT / "docs/source-material/3corners-forge-principles-v0.1.md"
PRINCIPLES_SHA256 = "7805ca4e13f1254e81f15e06a8a1d787862a4909aed276302f78f41354c3cee8"
PRINCIPLES_ADDENDUM = ROOT / "docs/source-material/3corners-forge-principles-addendum-v0.1.md"
PRINCIPLES_ADDENDUM_SHA256 = "9a3e68eb0588a3bdff230167f4c442ac6b57f7774c7d6450ca009a8fcffdeec5"
PRINCIPLES_ADDENDUM_V02 = ROOT / "docs/source-material/3corners-forge-principles-addendum-v0.2.md"
PRINCIPLES_ADDENDUM_V02_SHA256 = "9e80a6ba904e3eae209a81794a5621be46c0da1c13144b1ec65ae5f6a8a1e676"
ACCEPTED_DOCUMENTATION_BASELINE = (
    "docs/README.md",
    "docs/architecture/agentic-delivery.md",
    "docs/architecture/capabilities.md",
    "docs/architecture/continuous-security-testing.md",
    "docs/architecture/data-and-privacy.md",
    "docs/architecture/deployment.md",
    "docs/architecture/evidence-and-observability.md",
    "docs/architecture/experience.md",
    "docs/architecture/extensibility.md",
    "docs/architecture/identity-and-tenancy.md",
    "docs/architecture/intelligence.md",
    "docs/architecture/knowledge.md",
    "docs/architecture/overview.md",
    "docs/architecture/policy-and-authorization.md",
    "docs/architecture/quality-and-operations.md",
    "docs/architecture/runtime-and-workflow.md",
    "docs/architecture/security.md",
    "docs/architecture/software-supply-chain.md",
    "docs/architecture/user-journey-assurance.md",
    "docs/context/README.md",
    "docs/decisions/README.md",
    "docs/foundation/masterplan.md",
    "docs/foundation/principles.md",
    "docs/foundation/vision.md",
    "docs/governance/architecture-readiness.md",
    "docs/governance/content-migration.md",
    "docs/governance/delivery.md",
    "docs/governance/implementation-baseline.md",
    "docs/governance/workplace-baseline-assessment.md",
    "docs/reference-applications/workplace-orchestrator.md",
    "docs/reference/glossary.md",
    "docs/roadmap.md",
)


def fail(message: str) -> None:
    raise AssertionError(message)


def check_operating_contract() -> None:
    text = (ROOT / "AGENTS.md").read_text(encoding="utf-8")
    if "Status: Active" not in text:
        fail("AGENTS.md must be the active repository operating contract")
    if "docs/context/README.md" not in text:
        fail("AGENTS.md must route task context through docs/context/README.md")


def check_document_metadata() -> None:
    for relative_directory in METADATA_DIRECTORIES:
        for path in (ROOT / relative_directory).rglob("*.md"):
            text = path.read_text(encoding="utf-8")
            status = re.search(r"^Status: (.+)$", text, re.MULTILINE)
            owner = re.search(r"^Owner: (.+)$", text, re.MULTILINE)
            if not status or status.group(1) not in ALLOWED_DOCUMENT_STATUSES:
                fail(f"{path.relative_to(ROOT)} has missing or invalid Status metadata")
            if not owner or not owner.group(1).strip():
                fail(f"{path.relative_to(ROOT)} has missing Owner metadata")


def check_documentation_acceptance_baseline() -> None:
    expected_statuses = {
        **{path: "Accepted" for path in ACCEPTED_DOCUMENTATION_BASELINE},
        "docs/architecture/cryptography.md": "Proposed",
        "docs/governance/minimum-architecture-acceptance.md": "Superseded",
    }
    mismatches: list[str] = []
    for relative_path, expected in expected_statuses.items():
        text = (ROOT / relative_path).read_text(encoding="utf-8")
        match = re.search(r"^Status: (.+)$", text, re.MULTILINE)
        actual = match.group(1) if match else "missing"
        if actual != expected:
            mismatches.append(f"{relative_path}: expected={expected}, actual={actual}")
    if mismatches:
        fail("documentation acceptance baseline changed: " + "; ".join(mismatches))


def check_local_markdown_links() -> None:
    excluded = {PRINCIPLES_SOURCE}
    link_pattern = re.compile(r"\[[^\]]*\]\(([^)]+)\)")
    failures: list[str] = []
    for path in (ROOT / "docs").rglob("*.md"):
        if path in excluded:
            continue
        text = path.read_text(encoding="utf-8")
        for raw_target in link_pattern.findall(text):
            target = raw_target.strip().strip("<>").split()[0]
            parsed = urlparse(target)
            if parsed.scheme or target.startswith("#"):
                continue
            decoded_path = unquote(parsed.path)
            if not decoded_path:
                continue
            if decoded_path.startswith("/"):
                resolved = ROOT / decoded_path.lstrip("/")
            else:
                resolved = path.parent / decoded_path
            if not resolved.exists():
                failures.append(
                    f"{path.relative_to(ROOT)} -> {target}"
                )
    if failures:
        fail("broken local Markdown links: " + "; ".join(sorted(failures)))


def check_navigation_coverage() -> None:
    navigation = (ROOT / "mkdocs.yml").read_text(encoding="utf-8")
    missing: list[str] = []
    for relative_directory in NAVIGATION_DIRECTORIES:
        directory = ROOT / relative_directory
        for path in directory.rglob("*.md"):
            relative_to_docs = path.relative_to(ROOT / "docs").as_posix()
            if relative_to_docs not in navigation:
                missing.append(relative_to_docs)
    if missing:
        fail("canonical documents missing from MkDocs navigation: " + ", ".join(sorted(missing)))


def check_decision_register() -> None:
    decision_directory = ROOT / "docs/decisions"
    register = (decision_directory / "README.md").read_text(encoding="utf-8")
    entries = {
        filename: status
        for filename, status in re.findall(
            r"\| \[ADR-\d+\]\((\d{4}-[^)]+\.md)\) \| [^|]+ \| ([A-Za-z]+) \|",
            register,
        )
    }
    decision_files = {
        path.name for path in decision_directory.glob("[0-9][0-9][0-9][0-9]-*.md")
    }
    if set(entries) != decision_files:
        fail(
            "decision register differs from decision files: "
            f"registered {sorted(entries)}, found {sorted(decision_files)}"
        )
    mismatches: list[str] = []
    for filename, registered_status in entries.items():
        text = (decision_directory / filename).read_text(encoding="utf-8")
        match = re.search(r"^Status: (.+)$", text, re.MULTILINE)
        if not match or match.group(1) != registered_status:
            actual = match.group(1) if match else "missing"
            mismatches.append(f"{filename}: register={registered_status}, document={actual}")
    if mismatches:
        fail("decision status mismatch: " + "; ".join(mismatches))


def check_roadmap_structure() -> None:
    path = ROOT / "docs/roadmap.md"
    text = path.read_text(encoding="utf-8")
    required_foundation_intent = (
        "## F0 — Foundation readiness gate",
        "Accept the constitution, minimum capability boundary, trust model, and foundational ADRs",
        "Close the mandatory readiness gates",
        "Reconcile document status, capability maturity, and evidence labels",
        "Define non-functional requirements and deployment profiles",
        "Complete source-content migration without loss",
        "Replace development defaults that could be mistaken for production posture",
        "decisions/0007-sequence-self-building-before-full-advisory.md",
    )
    missing_intent = [item for item in required_foundation_intent if item not in text]
    if missing_intent:
        fail("roadmap lost required foundation intent: " + ", ".join(missing_intent))

    milestone_matches = list(
        re.finditer(r"^## M([1-9]) — .+$", text, re.MULTILINE)
    )
    found_numbers = [int(match.group(1)) for match in milestone_matches]
    if found_numbers != list(range(1, 10)):
        fail(f"roadmap milestones must be M1 through M9 in order, found {found_numbers}")
    required_sections = (
        "### Objective",
        "### Deliverables",
        "### Dependencies",
        "### Acceptance criteria",
        "### Estimated complexity",
        "### Risks",
    )
    for index, match in enumerate(milestone_matches):
        end = (
            milestone_matches[index + 1].start()
            if index + 1 < len(milestone_matches)
            else text.find("\n## Cross-milestone success measures", match.end())
        )
        section = text[match.start():end if end != -1 else len(text)]
        missing_sections = [item for item in required_sections if item not in section]
        if missing_sections:
            fail(
                f"roadmap M{match.group(1)} is missing sections: "
                + ", ".join(missing_sections)
            )


def check_f0_readiness_governance() -> None:
    register_path = ROOT / "docs/governance/architecture-readiness.md"
    register = register_path.read_text(encoding="utf-8")
    plan_path = ROOT / "docs/governance/f0-readiness-plan.md"
    plan = plan_path.read_text(encoding="utf-8")

    allowed_statuses = {"Open", "In progress", "Ready for review", "Accepted for F0"}
    gate_rows = re.findall(
        r"^\| G(\d+) [^|]+\|[^|]+\|[^|]+\| ([^|]+) \|$",
        register,
        re.MULTILINE,
    )
    gate_statuses = {int(identifier): status.strip() for identifier, status in gate_rows}
    expected_gates = set(range(11))
    if len(gate_rows) != len(expected_gates) or set(gate_statuses) != expected_gates:
        fail(
            "F0 readiness register must contain G0 through G10 exactly once: "
            f"found rows={len(gate_rows)}, gates={sorted(gate_statuses)}"
        )
    invalid = {
        identifier: status
        for identifier, status in gate_statuses.items()
        if status not in allowed_statuses
    }
    if invalid:
        fail(f"F0 readiness register uses invalid gate statuses: {invalid}")

    required_plan_content = (
        "## Gate status vocabulary",
        "## Evidence phases",
        "## Ownership model",
        "## Gate evidence assignments",
        "## Gate acceptance record",
        "## F0 final acceptance",
        "Before M1",
        "M1 exit",
        "Before Trial",
    )
    missing = [item for item in required_plan_content if item not in plan]
    if missing:
        fail("F0 readiness plan is missing required content: " + ", ".join(missing))
    for identifier in range(11):
        if f"| G{identifier} |" not in plan:
            fail(f"F0 readiness plan is missing ownership for G{identifier}")
    if "governance/f0-readiness-plan.md" not in (ROOT / "docs/roadmap.md").read_text(encoding="utf-8"):
        fail("roadmap must link the F0 readiness plan")


def check_capability_ids() -> None:
    path = ROOT / "docs/architecture/capabilities.md"
    text = path.read_text(encoding="utf-8")
    identifiers = re.findall(r"`(forge\.capability\.[a-z0-9-]+)`", text)
    if not identifiers:
        fail("capability catalogue contains no stable identifiers")
    duplicates = sorted({item for item in identifiers if identifiers.count(item) > 1})
    if duplicates:
        fail(f"duplicate capability identifiers: {', '.join(duplicates)}")


def check_trust_schemas() -> None:
    schema_directory = ROOT / "schemas/forge/trust/v1"
    required = {
        "execution-envelope.schema.json",
        "context-package.schema.json",
        "policy-decision.schema.json",
        "approval-record.schema.json",
        "evidence-record.schema.json",
    }
    present = {path.name for path in schema_directory.glob("*.schema.json")}
    if present != required:
        fail(f"trust schema set differs: expected {sorted(required)}, found {sorted(present)}")
    identifiers: set[str] = set()
    for path in schema_directory.glob("*.schema.json"):
        document = json.loads(path.read_text(encoding="utf-8"))
        for field in ("$schema", "$id", "title", "type"):
            if field not in document:
                fail(f"{path.relative_to(ROOT)} is missing {field}")
        if document["$id"] in identifiers:
            fail(f"duplicate schema identifier: {document['$id']}")
        identifiers.add(document["$id"])
        if document.get("type") != "object" or document.get("additionalProperties") is not False:
            fail(f"{path.relative_to(ROOT)} must be a closed top-level object schema")


def check_deployment_profiles() -> None:
    schema_path = ROOT / "schemas/forge/deployment/v1/profile.schema.json"
    schema = json.loads(schema_path.read_text(encoding="utf-8"))
    for field in ("$schema", "$id", "title", "type"):
        if field not in schema:
            fail(f"{schema_path.relative_to(ROOT)} is missing {field}")
    if schema.get("type") != "object" or schema.get("additionalProperties") is not False:
        fail("deployment profile schema must be a closed top-level object schema")

    profile_paths = sorted((ROOT / "profiles").glob("*.profile.json"))
    expected_profiles = {
        "developer-trial": "Trial",
        "connected-self-hosted-reference": "Research",
    }
    required_controls = {
        "identity",
        "storage",
        "repository",
        "worker_isolation",
        "network",
        "secrets",
        "evidence",
        "telemetry",
        "backup_restore",
        "upgrade_rollback",
        "capacity_recovery",
    }
    found_profiles: dict[str, str] = {}
    test_ids: set[str] = set()
    for path in profile_paths:
        profile = json.loads(path.read_text(encoding="utf-8"))
        profile_id = profile.get("profile_id")
        maturity = profile.get("maturity")
        if not isinstance(profile_id, str) or not isinstance(maturity, str):
            fail(f"{path.relative_to(ROOT)} is missing profile_id or maturity")
        if profile_id in found_profiles:
            fail(f"duplicate deployment profile ID: {profile_id}")
        found_profiles[profile_id] = maturity
        dimensions = profile.get("dimensions", {})
        required_dimensions = {
            "connectivity",
            "tenancy",
            "operations",
            "assurance",
            "provider_ownership",
        }
        if set(dimensions) != required_dimensions:
            fail(f"{path.relative_to(ROOT)} does not resolve every profile dimension")
        controls = profile.get("controls", {})
        if set(controls) != required_controls:
            fail(f"{path.relative_to(ROOT)} has incomplete control ownership")
        for control_name, control in controls.items():
            if not control.get("owner") or not control.get("assumption"):
                fail(f"{path.relative_to(ROOT)} control {control_name} lacks owner or assumption")
            for phase in ("before_m1_tests", "m1_exit_tests", "before_trial_tests"):
                tests = control.get(phase)
                if not isinstance(tests, list) or not tests:
                    fail(f"{path.relative_to(ROOT)} control {control_name} lacks {phase}")
        phases = set()
        for test in profile.get("conformance_tests", []):
            identifier = test.get("id")
            if not identifier or identifier in test_ids:
                fail(f"duplicate or missing deployment conformance test ID: {identifier}")
            test_ids.add(identifier)
            phases.add(test.get("phase"))
            if not test.get("command") or not test.get("expected"):
                fail(f"{path.relative_to(ROOT)} test {identifier} lacks command or expected result")
        if phases != {"before-m1", "m1-exit", "before-trial"}:
            fail(f"{path.relative_to(ROOT)} must define tests for all evidence phases")
        unsupported = profile.get("unsupported_claims", [])
        if not unsupported:
            fail(f"{path.relative_to(ROOT)} must declare unsupported claims")

    if found_profiles != expected_profiles:
        fail(
            "deployment profile set or maturity differs: "
            f"expected {expected_profiles}, found {found_profiles}"
        )


def check_assurance_schemas() -> None:
    schema_directory = ROOT / "schemas/forge/assurance/v1"
    required = {
        "journey-scenario.schema.json",
        "journey-result.schema.json",
    }
    present = {path.name for path in schema_directory.glob("*.schema.json")}
    if present != required:
        fail(f"assurance schema set differs: expected {sorted(required)}, found {sorted(present)}")
    identifiers: set[str] = set()
    for path in schema_directory.glob("*.schema.json"):
        document = json.loads(path.read_text(encoding="utf-8"))
        for field in ("$schema", "$id", "title", "type"):
            if field not in document:
                fail(f"{path.relative_to(ROOT)} is missing {field}")
        if document["$id"] in identifiers:
            fail(f"duplicate schema identifier: {document['$id']}")
        identifiers.add(document["$id"])
        if document.get("type") != "object" or document.get("additionalProperties") is not False:
            fail(f"{path.relative_to(ROOT)} must be a closed top-level object schema")


def check_principles_integrity() -> None:
    digest = hashlib.sha256(PRINCIPLES_SOURCE.read_bytes()).hexdigest()
    if digest != PRINCIPLES_SHA256:
        fail(
            "the immutable 3Corners Forge Principles snapshot changed; "
            "create a versioned addendum or new source version instead"
        )
    addendum_digest = hashlib.sha256(PRINCIPLES_ADDENDUM.read_bytes()).hexdigest()
    if addendum_digest != PRINCIPLES_ADDENDUM_SHA256:
        fail(
            "the versioned Forge Principles addendum changed; "
            "create a new addendum version instead"
        )
    addendum_v02_digest = hashlib.sha256(PRINCIPLES_ADDENDUM_V02.read_bytes()).hexdigest()
    if addendum_v02_digest != PRINCIPLES_ADDENDUM_V02_SHA256:
        fail(
            "the versioned Forge Principles addendum v0.2 changed; "
            "create a new addendum version instead"
        )


def check_developer_compose() -> None:
    path = ROOT / "deploy/compose/docker-compose.yml"
    text = path.read_text(encoding="utf-8")
    required_compose_markers = (
        "name: forge-development",
        "x-forge-profile: developer-only",
        "POSTGRES_PASSWORD: forge-dev-only",
        'command: ["start-dev"]',
        "KC_BOOTSTRAP_ADMIN_PASSWORD: admin-dev-only",
        "-dev-root-token-id=dev-only-token",
        "MINIO_ROOT_PASSWORD: minio-dev-only",
        "GF_SECURITY_ADMIN_PASSWORD: admin-dev-only",
        "postgres-data:",
        "minio-data:",
        "grafana-data:",
    )
    missing_markers = [marker for marker in required_compose_markers if marker not in text]
    if missing_markers:
        fail("developer Compose controls are missing: " + ", ".join(missing_markers))
    port_blocks = re.findall(r"^\s*ports:\s*\[([^\]]+)\]", text, re.MULTILINE)
    published_ports = []
    for block in port_blocks:
        published_ports.extend(re.findall(r'"([^"]+)"', block))
    unsafe = [port for port in published_ports if not port.startswith("127.0.0.1:")]
    if unsafe:
        fail(f"Compose ports must bind loopback only: {', '.join(unsafe)}")
    telemetry = (ROOT / "observability/otel-collector/config.yaml").read_text(encoding="utf-8")
    if "debug:" not in telemetry or "exporters: [debug]" not in telemetry:
        fail("Developer Trial telemetry must keep its debug exporter visible to conformance checks")
    negative_fixture = ROOT / "internal/cli/testdata/production-shaped.yaml"
    fixture_text = negative_fixture.read_text(encoding="utf-8")
    if "profile: connected-self-hosted" not in fixture_text:
        fail("G0 requires a production-shaped negative CLI fixture")


def main() -> int:
    checks = (
        check_operating_contract,
        check_document_metadata,
        check_documentation_acceptance_baseline,
        check_local_markdown_links,
        check_navigation_coverage,
        check_decision_register,
        check_roadmap_structure,
        check_f0_readiness_governance,
        check_capability_ids,
        check_trust_schemas,
        check_deployment_profiles,
        check_assurance_schemas,
        check_principles_integrity,
        check_developer_compose,
    )
    for check in checks:
        check()
        print(f"PASS {check.__name__}")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (AssertionError, json.JSONDecodeError) as error:
        print(f"FAIL {error}", file=sys.stderr)
        raise SystemExit(1)
