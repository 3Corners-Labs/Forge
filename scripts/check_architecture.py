#!/usr/bin/env python3
"""Dependency-free architecture-as-code checks for the Forge repository."""

from __future__ import annotations

import hashlib
import json
import re
import sys
from pathlib import Path


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
PRINCIPLES_SOURCE = ROOT / "docs/source-material/3corners-forge-principles-v0.1.md"
PRINCIPLES_SHA256 = "7805ca4e13f1254e81f15e06a8a1d787862a4909aed276302f78f41354c3cee8"
PRINCIPLES_ADDENDUM = ROOT / "docs/source-material/3corners-forge-principles-addendum-v0.1.md"
PRINCIPLES_ADDENDUM_SHA256 = "9a3e68eb0588a3bdff230167f4c442ac6b57f7774c7d6450ca009a8fcffdeec5"
PRINCIPLES_ADDENDUM_V02 = ROOT / "docs/source-material/3corners-forge-principles-addendum-v0.2.md"
PRINCIPLES_ADDENDUM_V02_SHA256 = "9e80a6ba904e3eae209a81794a5621be46c0da1c13144b1ec65ae5f6a8a1e676"


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
    if "x-forge-profile: developer-only" not in text:
        fail("Compose bootstrap must declare developer-only profile")
    port_blocks = re.findall(r"^\s*ports:\s*\[([^\]]+)\]", text, re.MULTILINE)
    published_ports = []
    for block in port_blocks:
        published_ports.extend(re.findall(r'"([^"]+)"', block))
    unsafe = [port for port in published_ports if not port.startswith("127.0.0.1:")]
    if unsafe:
        fail(f"Compose ports must bind loopback only: {', '.join(unsafe)}")


def main() -> int:
    checks = (
        check_operating_contract,
        check_document_metadata,
        check_capability_ids,
        check_trust_schemas,
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
