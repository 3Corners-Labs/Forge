#!/usr/bin/env python3
"""Dependency-free architecture-as-code checks for the Forge repository."""

from __future__ import annotations

import hashlib
import json
import re
import sys
from pathlib import Path
from urllib.parse import unquote, urlparse

from validate_autonomy_grant import (
    GRANT_DIRECTORY,
    MANDATORY_HUMAN_GATES,
    validate_grant_document,
)


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


def check_m1_authorization_governance() -> None:
    manifest_path = ROOT / "docs/governance/f0-acceptance-manifest.md"
    adr_path = ROOT / "docs/decisions/0010-authorize-bounded-m1-implementation.md"
    manifest = manifest_path.read_text(encoding="utf-8")
    adr = adr_path.read_text(encoding="utf-8")
    register = (ROOT / "docs/governance/architecture-readiness.md").read_text(encoding="utf-8")

    def document_status(text: str, path: Path) -> str:
        match = re.search(r"^Status: (.+)$", text, re.MULTILINE)
        if not match:
            fail(f"{path.relative_to(ROOT)} has no status")
        return match.group(1)

    manifest_status = document_status(manifest, manifest_path)
    adr_status = document_status(adr, adr_path)
    if manifest_status != adr_status or manifest_status not in {"Proposed", "Accepted"}:
        fail("F0 manifest and ADR-0010 must have matching Proposed or Accepted status")

    baseline_commit = "93efcfcae8e946310ebd07294d744f38a3083b0e"
    baseline_tree = "eeea9e5e0d95534cebc8ac253672d00b818027c8"
    for text, path in ((manifest, manifest_path), (adr, adr_path)):
        if baseline_commit not in text or baseline_tree not in text:
            fail(f"{path.relative_to(ROOT)} does not bind the immutable F0 baseline")

    for identifier in range(11):
        if f"| G{identifier} |" not in manifest:
            fail(f"F0 acceptance manifest is missing G{identifier}")
    required_manifest_sections = (
        "## Immutable F0 baseline",
        "## Gate evidence manifest",
        "## Verification evidence",
        "## Authorized M1 envelope",
        "## Explicit exclusions",
        "## Deferred evidence",
        "## Residual risks",
        "## Accountable review record",
    )
    missing = [section for section in required_manifest_sections if section not in manifest]
    if missing:
        fail("F0 acceptance manifest lacks sections: " + ", ".join(missing))

    required_adr_content = (
        "forge delivery start",
        "Developer Trial profile only",
        "## Explicitly unauthorized",
        "## Delivery controls",
        "## Stop conditions",
        "## M1 exit evidence",
        "## Rollback and recovery",
        "active or autonomous offensive testing",
        "merge a pull request",
        "Workplace-specific",
        "Supported",
    )
    missing = [item for item in required_adr_content if item not in adr]
    if missing:
        fail("ADR-0010 lacks required authorization boundaries: " + ", ".join(missing))

    if manifest_status == "Proposed":
        if "M1 functional implementation remains denied" not in register:
            fail("Proposed M1 authorization must leave implementation denied")
        if "| Pending | — |" not in manifest or "| Pending | — |" not in adr:
            fail("Proposed authorization documents must preserve pending human review")
    else:
        if "explicitly authorizes M1 functional implementation" not in register:
            fail("Accepted ADR-0010 must be reflected in the readiness register")
        if "| Pending | — |" in manifest or "| Pending | — |" in adr:
            fail("Accepted authorization documents cannot retain pending review outcomes")

    roadmap = (ROOT / "docs/roadmap.md").read_text(encoding="utf-8")
    if "governance/f0-acceptance-manifest.md" not in roadmap or "decisions/0010-authorize-bounded-m1-implementation.md" not in roadmap:
        fail("roadmap must link the final F0 manifest and ADR-0010")


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
        "approval-request.schema.json",
        "approval-record-v2.schema.json",
        "evidence-record.schema.json",
        "worker-request.schema.json",
        "candidate-bundle.schema.json",
        "worker-result.schema.json",
        "assurance-plan.schema.json",
        "verification-result.schema.json",
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


def check_common_contract_schemas() -> None:
    path = ROOT / "schemas/forge/common/v1/normalized-error.schema.json"
    document = json.loads(path.read_text(encoding="utf-8"))
    if document.get("type") != "object" or document.get("additionalProperties") is not False:
        fail("normalized error must be a closed top-level object schema")
    if document.get("properties", {}).get("schemaVersion", {}).get("const") != "forge.normalized-error.v1":
        fail("normalized error has an unsupported schema version")
    expected_codes = {
        "invalid_request",
        "unauthenticated",
        "forbidden",
        "indeterminate",
        "conflict",
        "expired",
        "budget_exhausted",
        "dependency_unavailable",
        "verification_failed",
        "evidence_unavailable",
        "manual_recovery_required",
    }
    actual_codes = set(document.get("properties", {}).get("code", {}).get("enum", []))
    if actual_codes != expected_codes:
        fail("normalized error code set changed")
    expected_retry = {"none", "reauthenticate", "reconcile", "new-authority", "bounded", "manual"}
    actual_retry = set(document.get("properties", {}).get("retry", {}).get("enum", []))
    if actual_retry != expected_retry:
        fail("normalized error retry set changed")


def check_intelligence_schemas() -> None:
    path = ROOT / "schemas/forge/intelligence/v1/provider-activation.schema.json"
    document = json.loads(path.read_text(encoding="utf-8"))
    for field in ("$schema", "$id", "title", "type"):
        if field not in document:
            fail(f"{path.relative_to(ROOT)} is missing {field}")
    if document.get("type") != "object" or document.get("additionalProperties") is not False:
        fail("model-provider activation must be a closed top-level object schema")
    properties = document.get("properties", {})
    if properties.get("schemaVersion", {}).get("const") != "forge.model-provider-activation.v1":
        fail("model-provider activation has an unsupported schema version")
    if properties.get("storeProviderState", {}).get("const") is not False:
        fail("Developer Trial model-provider activation must disable provider state storage")
    if properties.get("trainingUse", {}).get("const") != "disabled":
        fail("Developer Trial model-provider activation must disable provider training use")
    if properties.get("allowedPurposes", {}).get("items", {}).get("const") != "forge.product.discovery":
        fail("M1 model-provider activation must be limited to product discovery")
    if set(properties.get("allowedClassifications", {}).get("items", {}).get("enum", [])) != {
        "Public",
        "Internal",
    }:
        fail("M1 model-provider activation has unsupported data classifications")
    if set(document.get("required", [])) != set(properties):
        fail("model-provider activation must require every declared property")


def check_contributor_autonomy_grants() -> None:
    schema_path = ROOT / "schemas/forge/contributor/v1/autonomy-grant.schema.json"
    schema = json.loads(schema_path.read_text(encoding="utf-8"))
    for field in ("$schema", "$id", "title", "type"):
        if field not in schema:
            fail(f"{schema_path.relative_to(ROOT)} is missing {field}")
    if schema.get("type") != "object" or schema.get("additionalProperties") is not False:
        fail("contributor autonomy grant schema must be a closed top-level object")
    if schema.get("properties", {}).get("schemaVersion", {}).get("const") != "forge.contributor-autonomy-grant.v1":
        fail("contributor autonomy grant schema has an unsupported version")

    grant_directory = ROOT / GRANT_DIRECTORY
    grant_paths = sorted(grant_directory.glob("*.json"))
    if not grant_paths:
        fail("canonical contributor autonomy grant directory is empty")
    identifiers: set[str] = set()
    for path in grant_paths:
        grant = validate_grant_document(json.loads(path.read_text(encoding="utf-8")))
        identifier = grant["id"]
        if identifier in identifiers:
            fail(f"duplicate contributor autonomy grant ID: {identifier}")
        identifiers.add(identifier)
        if path.name != f"{identifier}.json":
            fail(f"{path.relative_to(ROOT)} does not match its grant ID")
        if set(grant["mandatoryHumanGates"]) != MANDATORY_HUMAN_GATES:
            fail(f"{path.relative_to(ROOT)} weakens mandatory human gates")

    expected = "forge.autonomy-grant.m1-bounded-delivery-v1"
    if expected not in identifiers:
        fail("the reviewed inactive-by-default M1 contributor grant is missing")
    m1_grant_path = grant_directory / f"{expected}.json"
    m1_grant = validate_grant_document(json.loads(m1_grant_path.read_text(encoding="utf-8")))
    if "architecture/module-boundaries.json" not in m1_grant["scope"]["pathPatterns"]:
        fail("the M1 grant cannot register accepted capability modules")
    if "architecture/**" in m1_grant["scope"]["pathPatterns"]:
        fail("the M1 grant must not receive broad architecture scope")

    contract = (ROOT / "AGENTS.md").read_text(encoding="utf-8")
    required_contract_markers = (
        "## Bounded autonomous delivery",
        GRANT_DIRECTORY,
        "candidate branch SHALL NOT activate",
        "default-branch revision",
        "pull-request-merge-ready",
    )
    missing = [marker for marker in required_contract_markers if marker not in contract]
    if missing:
        fail("operating contract lacks bounded-autonomy controls: " + ", ".join(missing))

    guide = (ROOT / "docs/governance/bounded-autonomous-delivery.md").read_text(encoding="utf-8")
    required_guide_markers = (
        "forge.autonomy-grant.m1-bounded-delivery-v1",
        "scripts/validate_autonomy_grant.py",
        "eligible-for-human-activation",
        "protected-default-branch-landing",
        "## Revocation and safe recovery",
    )
    missing = [marker for marker in required_guide_markers if marker not in guide]
    if missing:
        fail("bounded-autonomy guide is incomplete: " + ", ".join(missing))


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


def check_journey_contracts() -> None:
    scenario_directory = ROOT / "assurance/journeys/m1"
    expected_scenarios = {
        "forge.journey.m1.allowed-intent-capture",
        "forge.journey.m1.denied-intent-capture",
    }
    schema = json.loads(
        (ROOT / "schemas/forge/assurance/v1/journey-scenario.schema.json").read_text(
            encoding="utf-8"
        )
    )
    required_fields = set(schema["required"])
    scenarios: dict[str, dict] = {}
    for path in sorted(scenario_directory.glob("*.scenario.json")):
        scenario = json.loads(path.read_text(encoding="utf-8"))
        missing = required_fields - set(scenario)
        if missing:
            fail(f"{path.relative_to(ROOT)} lacks required fields: {sorted(missing)}")
        identifier = scenario.get("id")
        if identifier in scenarios:
            fail(f"duplicate journey scenario ID: {identifier}")
        scenarios[identifier] = scenario
        if scenario.get("owner") == "Forge Engineering":
            fail(f"{path.relative_to(ROOT)} oracle must be owned independently from implementation")
        if scenario.get("lifecycle") not in {"reviewed", "accepted"}:
            fail(f"{path.relative_to(ROOT)} must be reviewed before G9 review")
        interface = scenario.get("interface", {})
        if interface.get("type") != "cli" or interface.get("publicEntryPointReference") != "forge delivery start":
            fail(f"{path.relative_to(ROOT)} must use the ADR-0009 public CLI boundary")
        if scenario.get("environmentProfile") != "developer-trial":
            fail(f"{path.relative_to(ROOT)} must use the Developer Trial profile")
        for fixture_reference in scenario.get("fixtureReferences", []):
            if not (ROOT / fixture_reference).is_file():
                fail(f"{path.relative_to(ROOT)} references missing fixture {fixture_reference}")
        expected_outcomes = scenario.get("expectedOutcomes", [])
        if not expected_outcomes or any(item.get("oracle") != "deterministic" for item in expected_outcomes):
            fail(f"{path.relative_to(ROOT)} must use deterministic M1 oracles")
        if not scenario.get("forbiddenOutcomes"):
            fail(f"{path.relative_to(ROOT)} must declare forbidden outcomes")
        if scenario.get("retryPolicy", {}).get("maximumAttempts") != 1:
            fail(f"{path.relative_to(ROOT)} must not mask M1 failures with retries")
        quarantine = scenario.get("quarantinePolicy", {})
        if not quarantine.get("owner") or not quarantine.get("releaseConditions"):
            fail(f"{path.relative_to(ROOT)} lacks accountable quarantine governance")
        evidence = scenario.get("evidencePolicy", {})
        if evidence.get("redactionRequired") is not True:
            fail(f"{path.relative_to(ROOT)} must require evidence redaction")
        actions = " ".join(step.get("action", "") for step in scenario.get("steps", []))
        if "forge delivery start" not in actions:
            fail(f"{path.relative_to(ROOT)} does not execute the public CLI command")

    if set(scenarios) != expected_scenarios:
        fail(f"M1 journey set differs: expected {sorted(expected_scenarios)}, found {sorted(scenarios)}")
    allowed = scenarios["forge.journey.m1.allowed-intent-capture"]
    denied = scenarios["forge.journey.m1.denied-intent-capture"]
    if allowed["persona"]["role"] == denied["persona"]["role"]:
        fail("allowed and denied journeys must exercise different roles")
    denied_assertions = " ".join(item["assertion"] for item in denied["expectedOutcomes"])
    if "forbidden" not in denied_assertions or "before repository discovery" not in denied_assertions:
        fail("denied journey must fail closed before repository discovery")
    if allowed.get("version") != "2.0.0":
        fail("allowed M1 journey must use the accepted 2.0.0 oracle")
    allowed_assertions = " ".join(item["assertion"] for item in allowed["expectedOutcomes"])
    for required in (
        "Requirements Proposal",
        "plan digest",
        "independent Assurance Plan",
        "candidate digest",
        "contract-double",
    ):
        if required not in allowed_assertions:
            fail(f"allowed M1 journey must assert {required}")
    allowed_forbidden = " ".join(allowed["forbiddenOutcomes"])
    for prohibited in ("Live GitHub", "merge", "release", "deployment", "production"):
        if prohibited.lower() not in allowed_forbidden.lower():
            fail(f"allowed M1 journey must forbid {prohibited}")

    plan_path = scenario_directory / "trial-execution-plan.json"
    plan = json.loads(plan_path.read_text(encoding="utf-8"))
    if plan.get("publicEntryPoint") != "forge delivery start":
        fail("M1 execution plan must bind the ADR-0009 public CLI boundary")
    if plan.get("owner") == plan.get("implementationOwner"):
        fail("M1 execution-plan owner must be independent from implementation")
    if plan.get("executionState") != "executable-m1-candidate":
        fail("M1 journey runner must be an executable M1 candidate")
    if plan.get("resultSchema") != "schemas/forge/assurance/v1/journey-result.schema.json":
        fail("M1 execution plan must bind the journey-result schema")
    if not plan.get("applicationRevisionBinding", {}).get("required"):
        fail("M1 execution plan must bind an immutable application revision")
    runner = plan.get("runnerBinding", {})
    if not runner.get("required") or runner.get("candidateControlled") is not False:
        fail("M1 execution plan must bind an independent immutable runner")
    bindings = plan.get("scenarioBindings", [])
    if len(bindings) != 2:
        fail("M1 execution plan must bind exactly two scenarios")
    for binding in bindings:
        fixture_path = ROOT / binding["fixture"]
        actual_digest = hashlib.sha256(fixture_path.read_bytes()).hexdigest()
        if binding.get("fixtureSha256") != actual_digest:
            fail(f"fixture digest mismatch for {fixture_path.relative_to(ROOT)}")
        fixture = json.loads(fixture_path.read_text(encoding="utf-8"))
        if binding["identityProfileId"] != fixture["identityProfileId"]:
            fail(f"identity binding mismatch for {binding['scenario']}")


def check_offensive_assurance_safety() -> None:
    schema_path = ROOT / "schemas/forge/security/v1/rules-of-engagement.schema.json"
    schema = json.loads(schema_path.read_text(encoding="utf-8"))
    for field in ("$schema", "$id", "title", "type"):
        if field not in schema:
            fail(f"{schema_path.relative_to(ROOT)} is missing {field}")
    if schema.get("type") != "object" or schema.get("additionalProperties") is not False:
        fail("rules-of-engagement schema must be a closed top-level object")

    roe_path = ROOT / "security/assurance/m1.rules-of-engagement.json"
    roe = json.loads(roe_path.read_text(encoding="utf-8"))
    missing = set(schema["required"]) - set(roe)
    if missing:
        fail(f"{roe_path.relative_to(ROOT)} lacks required fields: {sorted(missing)}")
    target_path = ROOT / "security/assurance/m1-target-set.json"
    runner_path = ROOT / "security/assurance/m1-runner-config.json"
    expected_target_digest = "sha256:" + hashlib.sha256(target_path.read_bytes()).hexdigest()
    expected_runner_digest = "sha256:" + hashlib.sha256(runner_path.read_bytes()).hexdigest()
    authorization = roe["authorization"]
    if authorization.get("targetSetDigest") != expected_target_digest:
        fail("rules of engagement do not bind the exact target set")
    if authorization.get("runnerConfigurationDigest") != expected_runner_digest:
        fail("rules of engagement do not bind the exact runner configuration")
    target_set = json.loads(target_path.read_text(encoding="utf-8"))
    runner_config = json.loads(runner_path.read_text(encoding="utf-8"))
    allowed_targets = set(target_set["targets"])

    def violations(candidate: dict) -> set[str]:
        result: set[str] = set()
        environment = candidate.get("authorization", {}).get("environment")
        scope = candidate.get("scope", {})
        runner = candidate.get("runner", {})
        stop = candidate.get("stopControls", {})
        limits = candidate.get("limits", {})
        targets = set(scope.get("targets", []))
        if environment == "production" or scope.get("productionTargetsAllowed") is not False:
            result.add("production-target")
        if scope.get("thirdPartyTargetsAllowed") is not False or any(
            not target.startswith("fixture://") for target in targets
        ):
            result.add("third-party-target")
        if targets != allowed_targets:
            result.add("scope-expansion")
        if runner.get("enabled") is not False or runner.get("generatedContentMayConfigure") is not False:
            result.add("generated-activation")
        if runner.get("autonomous") is not False or runner.get("mode") != "contract-only":
            result.add("active-runner")
        if runner.get("scopeEnforcement") != "external-deterministic-controller":
            result.add("scope-enforcement")
        if stop.get("killSwitchEngaged") is not True or stop.get("failClosed") is not True:
            result.add("kill-switch-override")
        if any(limits.get(field) != 0 for field in ("activeNetworkRequests", "requestsPerSecond", "costUnits", "payloadBytes")):
            result.add("active-impact")
        return result

    base_violations = violations(roe)
    if base_violations:
        fail(f"M1 rules of engagement are unsafe: {sorted(base_violations)}")
    if runner_config != {
        "runnerConfigurationVersion": "forge.security-runner-config.v1",
        "enabled": False,
        "mode": "contract-only",
        "activeNetworkRequests": 0,
        "generatedContentMayConfigure": False,
        "scopeEnforcement": "external-deterministic-controller",
        "killSwitchEngaged": True,
    }:
        fail("M1 offensive runner configuration must remain exactly disabled and contract-only")

    def merged(base: dict, override: dict) -> dict:
        result = json.loads(json.dumps(base))
        stack = [(result, override)]
        while stack:
            destination, source = stack.pop()
            for key, value in source.items():
                if isinstance(value, dict) and isinstance(destination.get(key), dict):
                    stack.append((destination[key], value))
                else:
                    destination[key] = value
        return result

    fixture_directory = ROOT / "security/assurance/testdata"
    expected_fixtures = {
        "generated-activation.json",
        "kill-switch-override.json",
        "production-target.json",
        "scope-expansion.json",
        "third-party-target.json",
    }
    fixture_paths = sorted(fixture_directory.glob("*.json"))
    if {path.name for path in fixture_paths} != expected_fixtures:
        fail("offensive-safety negative fixture set changed")
    for path in fixture_paths:
        fixture = json.loads(path.read_text(encoding="utf-8"))
        expected = fixture.get("expectedViolation")
        actual = violations(merged(roe, fixture.get("override", {})))
        if expected not in actual:
            fail(f"{path.relative_to(ROOT)} did not fail with {expected}; found {sorted(actual)}")


def check_module_boundaries() -> None:
    manifest_path = ROOT / "architecture/module-boundaries.json"
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    if manifest.get("schemaVersion") != "forge.module-boundaries.v1":
        fail("module-boundary manifest has an unsupported schema version")
    modules = manifest.get("modules", [])
    by_id = {module["id"]: module for module in modules}
    if len(by_id) != len(modules):
        fail("module-boundary manifest contains duplicate module IDs")
    by_path = {module["path"]: module for module in modules}
    if len(by_path) != len(modules):
        fail("module-boundary manifest contains duplicate module paths")
    expected_identity_module = {
        "id": "identity",
        "path": "internal/identity",
        "deploymentUnit": "forge-control-plane",
        "maturity": "Trial",
        "allowedInternalDependencies": ["failure-contract", "trust-contract"],
    }
    if by_id.get("identity") != expected_identity_module:
        fail("M1 Identity must depend only on failure and shared trust contracts")
    expected_policy_module = {
        "id": "policy",
        "path": "internal/policy",
        "deploymentUnit": "forge-control-plane",
        "maturity": "Trial",
        "allowedInternalDependencies": [
            "delivery",
            "failure-contract",
            "identity",
            "tenant",
            "trust-contract",
        ],
    }
    if by_id.get("policy") != expected_policy_module:
        fail("M1 policy must remain a Trial control-plane module with its accepted dependencies")
    expected_delivery_module = {
        "id": "delivery",
        "path": "internal/delivery",
        "deploymentUnit": "forge-control-plane",
        "maturity": "Trial",
        "allowedInternalDependencies": [
            "failure-contract",
            "identity",
            "tenant",
            "trust-contract",
            "verification-contracts",
            "worker-contracts",
        ],
    }
    if by_id.get("delivery") != expected_delivery_module:
        fail("M1 Delivery must own planning without importing Approval or Policy runtime state")
    expected_context_module = {
        "id": "context",
        "path": "internal/context",
        "deploymentUnit": "forge-control-plane",
        "maturity": "Trial",
        "allowedInternalDependencies": [
            "delivery",
            "failure-contract",
            "identity",
            "tenant",
            "trust-contract",
        ],
    }
    if by_id.get("context") != expected_context_module:
        fail("M1 context must remain a Trial enforcement module without Policy runtime access")
    expected_approval_module = {
        "id": "approval",
        "path": "internal/approval",
        "deploymentUnit": "forge-control-plane",
        "maturity": "Trial",
        "allowedInternalDependencies": [
            "failure-contract",
            "identity",
            "trust-contract",
        ],
    }
    if by_id.get("approval") != expected_approval_module:
        fail("M1 Approval must own approval state without importing Delivery or Policy runtime state")
    expected_m1_b_modules = {
        "execution": ("forge-control-plane", {"delivery", "failure-contract", "trust-contract", "worker-contracts"}),
        "verification": ("forge-control-plane", {"delivery", "failure-contract", "verification-contracts", "worker-contracts"}),
        "evidence": ("forge-control-plane", {"failure-contract", "trust-contract"}),
        "worker": ("forge-worker", {"worker-contracts"}),
        "verifier": ("forge-verifier", {"verification-contracts", "worker-contracts"}),
    }
    for module_id, (unit, dependencies) in expected_m1_b_modules.items():
        module = by_id.get(module_id)
        if module is None or module.get("deploymentUnit") != unit or module.get("maturity") != "Trial" or set(module.get("allowedInternalDependencies", [])) != dependencies:
            fail(f"M1-B module {module_id} does not preserve its accepted capability and deployment boundary")
    allowed_maturities = {"Research", "Trial", "Supported", "Deprecated"}
    for module in modules:
        if module.get("maturity") not in allowed_maturities:
            fail(f"module {module['id']} has invalid maturity")
        for dependency in module.get("allowedInternalDependencies", []):
            if dependency not in by_id:
                fail(f"module {module['id']} allows unknown dependency {dependency}")
        path_exists = (ROOT / module["path"]).is_dir()
        if module.get("planned") is True and path_exists:
            fail(f"planned module {module['id']} exists but is still represented as unimplemented")
        if module.get("planned") is not True and not path_exists:
            fail(f"implemented module {module['id']} has no source directory")

    isolation_rules: dict[str, set[str]] = {}
    for rule in manifest.get("isolationRules", []):
        identifier = rule["id"]
        if identifier in isolation_rules:
            fail(f"duplicate module isolation rule {identifier}")
        isolation_rules[identifier] = set(rule["deniedToUnits"])
    expected_rules = {
        "control-plane-cannot-import-worker",
        "worker-cannot-import-control-plane",
        "control-plane-cannot-import-verifier",
        "verifier-cannot-import-control-plane",
        "shared-contract-cannot-import-runtime",
    }
    if set(isolation_rules) != expected_rules:
        fail("module isolation rule set changed")

    module_prefix = "github.com/3corners/forge/"

    def target_module(import_path: str) -> dict | None:
        if not import_path.startswith(module_prefix):
            return None
        relative = import_path.removeprefix(module_prefix)
        candidates = [module for path, module in by_path.items() if relative == path or relative.startswith(path + "/")]
        if not candidates:
            return {"id": "unknown", "path": relative, "deploymentUnit": "unknown"}
        return max(candidates, key=lambda module: len(module["path"]))

    def dependency_violations(importer: dict, imports: list[str]) -> set[str]:
        result: set[str] = set()
        allowed = set(importer.get("allowedInternalDependencies", []))
        for import_path in imports:
            target = target_module(import_path)
            if target is None:
                continue
            if target["id"] == "unknown":
                result.add("unknown-internal-module")
                continue
            if target["id"] != importer["id"] and target["id"] not in allowed:
                result.add("undeclared-internal-dependency")
            source_unit = importer["deploymentUnit"]
            target_unit = target["deploymentUnit"]
            for rule in manifest["isolationRules"]:
                if rule["fromUnit"] == source_unit and target_unit in rule["deniedToUnits"]:
                    result.add(rule["id"])
        return result

    go_import_pattern = re.compile(r'"(github\.com/3corners/forge/[^"\s]+)"')
    for module in modules:
        if module.get("planned") is True:
            continue
        imports: list[str] = []
        for path in (ROOT / module["path"]).rglob("*.go"):
            if path.name.endswith("_test.go"):
                continue
            imports.extend(go_import_pattern.findall(path.read_text(encoding="utf-8")))
        violations = dependency_violations(module, imports)
        if violations:
            fail(f"module {module['id']} violates boundaries: {sorted(violations)}")

    delivery_source = "\n".join(
        path.read_text(encoding="utf-8")
        for path in (ROOT / "internal/delivery").rglob("*.go")
        if not path.name.endswith("_test.go")
    )
    if "trust.ExecutionEnvelope{" in delivery_source or "BindExecutionEnvelopeIntegrity(" in delivery_source:
        fail("Delivery must not construct or integrity-bind Execution Envelopes")
    if "trust.ApprovalRecordV2{" in delivery_source or "BindApprovalRecordV2Integrity(" in delivery_source:
        fail("Delivery must not construct or integrity-bind Approval Records")

    fixture_directory = ROOT / "architecture/testdata"
    expected_fixture_names = {
        "forbidden-control-to-worker.json",
        "forbidden-undeclared-module.json",
        "forbidden-worker-to-control.json",
        "m1-b-boundaries/forbidden-control-to-verifier.json",
        "m1-b-boundaries/forbidden-verifier-to-control.json",
        "m1-c-boundaries/forbidden-experience-to-delivery.json",
        "m1-c-boundaries/forbidden-source-control-to-approval.json",
    }
    fixture_paths = sorted(fixture_directory.rglob("*.json"))
    if {str(path.relative_to(fixture_directory)) for path in fixture_paths} != expected_fixture_names:
        fail("module-boundary negative fixture set changed")
    for path in fixture_paths:
        fixture = json.loads(path.read_text(encoding="utf-8"))
        importer = by_id[fixture["importerModule"]]
        actual = dependency_violations(importer, fixture["imports"])
        expected = set(fixture["expectedViolations"])
        if actual != expected:
            fail(f"{path.relative_to(ROOT)} expected {sorted(expected)}, found {sorted(actual)}")


def check_required_ci() -> None:
    architecture_workflow = (ROOT / ".github/workflows/architecture.yml").read_text(encoding="utf-8")
    if "pull_request:" not in architecture_workflow:
        fail("architecture workflow must run for every pull request")
    if re.search(r"^\s+paths(?:-ignore)?:", architecture_workflow, re.MULTILINE):
        fail("architecture workflow must not be path-filtered")
    required_architecture_markers = (
        "permissions:\n  contents: read",
        "timeout-minutes: 10",
        "run: python3 scripts/check_architecture.py",
        "run: python3 -m unittest scripts/test_autonomy_grant.py",
        "mkdocs build --strict",
    )
    missing = [marker for marker in required_architecture_markers if marker not in architecture_workflow]
    if missing:
        fail("architecture workflow lacks required controls: " + ", ".join(missing))

    ci_workflow = (ROOT / ".github/workflows/ci.yml").read_text(encoding="utf-8")
    required_ci_markers = (
        'test -z "$(gofmt -l cmd internal)"',
        "go vet ./...",
        "go test ./... -race -coverprofile=coverage.out",
        "go build -trimpath -o dist/forge ./cmd/forge",
    )
    missing = [marker for marker in required_ci_markers if marker not in ci_workflow]
    if missing:
        fail("CI workflow lacks required verification: " + ", ".join(missing))

    codeowners = (ROOT / ".github/CODEOWNERS").read_text(encoding="utf-8")
    required_owners = (
        "* @cyberbeuk @3c-unified-git",
        "/architecture/ @cyberbeuk @3c-unified-git",
        "/assurance/ @cyberbeuk @3c-unified-git",
        "/security/ @cyberbeuk @3c-unified-git",
        "/scripts/check_architecture.py @cyberbeuk @3c-unified-git",
        "/scripts/validate_autonomy_grant.py @cyberbeuk @3c-unified-git",
        "/scripts/test_autonomy_grant.py @cyberbeuk @3c-unified-git",
        "/.forge/autonomy/ @cyberbeuk @3c-unified-git",
        "/schemas/forge/contributor/ @cyberbeuk @3c-unified-git",
        "/AGENTS.md @cyberbeuk @3c-unified-git",
        "/.github/workflows/architecture.yml @cyberbeuk @3c-unified-git",
    )
    missing = [marker for marker in required_owners if marker not in codeowners]
    if missing:
        fail("architecture enforcement ownership is incomplete: " + ", ".join(missing))
    if "@3corners/" in codeowners:
        fail("CODEOWNERS must not rely on unverified organization-team aliases")


def check_threat_privacy_readiness() -> None:
    plan_path = ROOT / "security/threat-controls/m1.json"
    plan = json.loads(plan_path.read_text(encoding="utf-8"))
    if plan.get("schemaVersion") != "forge.threat-control-plan.v1" or plan.get("status") != "reviewed":
        fail("M1 threat-control plan must be a reviewed v1 contract")
    expected_threats = {
        "approval-replay",
        "candidate-self-attestation",
        "confused-deputy",
        "cross-tenant-leakage",
        "denial-of-wallet",
        "evidence-tampering",
        "malicious-repository",
        "offensive-testing-scope-escape",
        "operator-misuse",
        "poisoned-knowledge",
        "prompt-injection",
        "provider-leakage",
        "scenario-oracle-capture",
        "supply-chain-compromise",
        "synthetic-user-data-leakage",
        "tool-abuse",
        "untrusted-scanner-template",
    }
    mappings = plan.get("mappings", [])
    by_id = {mapping["id"]: mapping for mapping in mappings}
    if len(by_id) != len(mappings) or set(by_id) != expected_threats:
        fail("M1 threat-control mapping does not exactly cover the accepted priority threats")
    phases: set[str] = set()
    test_ids: set[str] = set()
    for mapping in mappings:
        for field in ("threat", "owner", "controls", "enforcementPoints", "negativeTests", "evidenceObligations"):
            if not mapping.get(field):
                fail(f"threat {mapping['id']} lacks {field}")
        for test in mapping["negativeTests"]:
            identifier = test.get("id")
            if not identifier or identifier in test_ids:
                fail(f"duplicate or missing threat negative-test ID: {identifier}")
            test_ids.add(identifier)
            phase = test.get("phase")
            if phase not in {"before-m1", "m1-exit", "before-trial"}:
                fail(f"threat test {identifier} has invalid evidence phase")
            phases.add(phase)
            fixture = test.get("fixture")
            if not fixture or not (ROOT / fixture).is_file():
                fail(f"threat test {identifier} references missing fixture {fixture}")
            if not test.get("expected") or test.get("disclosureSafe") is not True:
                fail(f"threat test {identifier} lacks a disclosure-safe expected result")
    if phases != {"before-m1", "m1-exit", "before-trial"}:
        fail("threat tests must assign evidence across all three phases")

    privacy_path = ROOT / "security/privacy/developer-trial-decisions.json"
    privacy = json.loads(privacy_path.read_text(encoding="utf-8"))
    if privacy.get("profile") != "developer-trial" or privacy.get("status") != "reviewed":
        fail("Developer Trial security/privacy decisions must be reviewed and profile-specific")
    expected_decisions = {
        "cancellation",
        "context_classification",
        "evidence_integrity",
        "identity_assurance",
        "repository_credentials",
        "retention",
        "revocation",
        "source_control_destination",
        "telemetry_minimization",
        "tenant_topology",
        "worker_isolation",
    }
    decisions = privacy.get("decisions", {})
    if set(decisions) != expected_decisions:
        fail("Developer Trial security/privacy decision set is incomplete")
    for topic, decision in decisions.items():
        for field in ("decision", "enforcement", "m1Evidence", "beforeTrialEvidence"):
            if not decision.get(field):
                fail(f"Developer Trial decision {topic} lacks {field}")
    unresolved = privacy.get("unresolvedProductionConstraints", [])
    if len(unresolved) < len(expected_decisions):
        fail("production-shaped security/privacy constraints are not explicit")

    seeded_path = ROOT / "security/fixtures/seeded-sensitive-data.json"
    seeded = json.loads(seeded_path.read_text(encoding="utf-8"))
    if seeded.get("synthetic") is not True or seeded.get("classification") != "Restricted":
        fail("seeded data fixture must be synthetic and Restricted")
    if "DO_NOT_RECORD" not in seeded.get("secretMarker", ""):
        fail("seeded secret fixture lacks a deterministic leak marker")
    personal = seeded.get("personalData", {})
    if ".invalid" not in personal.get("email", ""):
        fail("seeded personal-data fixture must use a reserved invalid domain")
    required_destinations = {
        "model-provider",
        "pull-request-description",
        "routine-evidence",
        "telemetry",
        "tool-output",
    }
    if set(seeded.get("forbiddenDestinations", [])) != required_destinations:
        fail("seeded data fixture does not cover every forbidden destination")


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


def check_m1_planning_lifecycle() -> None:
    decision = (ROOT / "docs/decisions/0015-context-informed-delivery-plan-review.md").read_text(
        encoding="utf-8"
    )
    runtime = (ROOT / "docs/architecture/runtime-and-workflow.md").read_text(encoding="utf-8")
    contracts = (ROOT / "docs/architecture/first-slice-capability-contracts.md").read_text(
        encoding="utf-8"
    )
    vertical_slice = (ROOT / "docs/architecture/vertical-slice.md").read_text(encoding="utf-8")
    data_flow = (ROOT / "docs/architecture/first-slice-data-flow.md").read_text(encoding="utf-8")
    roadmap = (ROOT / "docs/roadmap.md").read_text(encoding="utf-8")

    if "Status: Accepted" not in decision:
        fail("ADR-0015 must remain Accepted after accountable review")
    required_runtime_order = (
        "Captured -> ContextAuthorizing -> ContextReady -> Planned",
        "-> AwaitingPlanApproval -> PlanApproved",
        "-> ExecutionAuthorizing -> Executing -> Verifying",
        "-> AwaitingPublicationApproval -> Publishing -> Completed",
    )
    if any(marker not in runtime for marker in required_runtime_order):
        fail("runtime workflow does not preserve the accepted context-informed planning order")
    forbidden_lifecycle_markers = (
        "Captured -> Planned -> Authorized -> ContextReady",
        "Verifying -> AwaitingApproval -> Publishing",
    )
    lifecycle_documents = (runtime, contracts, vertical_slice, data_flow, roadmap)
    if any(
        marker in document
        for marker in forbidden_lifecycle_markers
        for document in lifecycle_documents
    ):
        fail("an ambiguous pre-context planning or approval state was reintroduced")

    required_contract_markers = (
        "RecordContext(runId, expectedVersion, contextPackageId, manifestDigest, evidenceRefs)",
        "forge.delivery.execute-plan",
        "Publication approval is a separate request",
    )
    if any(marker not in contracts for marker in required_contract_markers):
        fail("first-slice contracts do not preserve separate context, plan, and publication gates")
    if "PlanApproval --> ExecutionPolicy" not in vertical_slice:
        fail("vertical slice bypasses plan approval before execution policy")
    if "PlanApproval --> ExecutionAuthority" not in data_flow:
        fail("first-slice data flow bypasses plan approval before execution authority")
    if "Human plan review binds that exact plan digest but does not grant worker authority" not in roadmap:
        fail("roadmap does not preserve the ADR-0015 authority boundary")


def check_m1_remaining_runway() -> None:
    decisions = {
        identifier: (ROOT / path).read_text(encoding="utf-8")
        for identifier, path in {
            "ADR-0016": "docs/decisions/0016-scope-bound-approval-record-v2.md",
            "ADR-0017": "docs/decisions/0017-bounded-autonomous-worker-and-verification.md",
            "ADR-0018": "docs/decisions/0018-developer-trial-durable-state-and-evidence.md",
            "ADR-0019": "docs/decisions/0019-github-developer-trial-publication-adapter.md",
        }.items()
    }
    if any("Status: Accepted" not in text for text in decisions.values()):
        fail("the accepted M1 architecture-runway packet is incomplete")

    approval = decisions["ADR-0016"]
    for marker in (
        "forge.approval-record.v2",
        "forge.approval-request.v1",
        "V1 may validate its own historical structure but SHALL be rejected",
    ):
        if marker not in approval:
            fail(f"ADR-0016 lost required Approval v2 boundary: {marker}")

    worker = decisions["ADR-0017"]
    for marker in (
        "After one valid human approval of the exact Delivery Plan",
        "It does not request additional human approval for:",
        "one deterministic change-preparation worker",
        "one independent verifier",
        "forge.worker.prepare-change",
    ):
        if marker not in worker:
            fail(f"ADR-0017 lost bounded-autonomy or independence rule: {marker}")

    persistence = decisions["ADR-0018"]
    for marker in (
        "Use PostgreSQL as the Developer Trial adapter",
        "Use MinIO as the Developer Trial adapter",
        "Modules do not read another capability's tables",
        "Evidence Records remain distinct from telemetry",
    ):
        if marker not in persistence:
            fail(f"ADR-0018 lost durable-state or evidence boundary: {marker}")

    publication = decisions["ADR-0019"]
    for marker in (
        "dedicated GitHub App installation identity",
        "pull_requests:write",
        "adapter endpoint allowlist containing no merge",
        "repository rules on the protected default branch that deny updates by the GitHub App",
        "It SHALL NOT automatically choose",
    ):
        if marker not in publication:
            fail(f"ADR-0019 lost GitHub publication or no-merge boundary: {marker}")

    roadmap = (ROOT / "docs/roadmap.md").read_text(encoding="utf-8")
    for marker in (
        "M1-A Context to approved plan",
        "M1-B Approved plan to verified candidate",
        "M1-C Verified candidate to reviewable pull request",
        "do not introduce micro-approval gates",
    ):
        if marker not in roadmap:
            fail(f"roadmap lost the accelerated M1 tranche plan: {marker}")


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
        check_m1_authorization_governance,
        check_capability_ids,
        check_trust_schemas,
        check_common_contract_schemas,
        check_intelligence_schemas,
        check_contributor_autonomy_grants,
        check_deployment_profiles,
        check_assurance_schemas,
        check_journey_contracts,
        check_offensive_assurance_safety,
        check_module_boundaries,
        check_required_ci,
        check_threat_privacy_readiness,
        check_m1_planning_lifecycle,
        check_m1_remaining_runway,
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
