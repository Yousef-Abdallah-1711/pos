#!/usr/bin/env python3
"""Fail the release gate unless every required workflow job succeeded."""

from __future__ import annotations

import json
import os
import sys
from collections.abc import Mapping


REQUIRED_JOBS = (
    "backend",
    "frontend",
    "security-and-config",
    "application-images",
    "financial-postgres-acceptance",
)
OPTIONAL_BROWSER_JOB = "auth-browser-acceptance"


def failed_required_jobs(results: Mapping[str, Mapping[str, str]], require_browser: bool) -> dict[str, str]:
    """Return required jobs that are absent, skipped, or otherwise unsuccessful."""
    required_names = list(REQUIRED_JOBS)
    if require_browser:
        required_names.append(OPTIONAL_BROWSER_JOB)
    return {
        name: results.get(name, {}).get("result", "missing")
        for name in required_names
        if results.get(name, {}).get("result") != "success"
    }


def main() -> int:
    try:
        results = json.loads(os.environ["RESULTS"])
    except (KeyError, json.JSONDecodeError) as error:
        print(f"invalid required CI job results: {error}", file=sys.stderr)
        return 2
    if not isinstance(results, dict):
        print("required CI job results must be a JSON object", file=sys.stderr)
        return 2

    require_browser = os.environ.get("REQUIRE_AUTH_BROWSER") == "true"
    failed = failed_required_jobs(results, require_browser)
    if failed:
        print(f"required CI jobs did not succeed: {failed}", file=sys.stderr)
        return 1
    print("all required CI jobs succeeded")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
