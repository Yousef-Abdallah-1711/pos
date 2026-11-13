#!/usr/bin/env python3
"""Fail closed when the protected production environment is incomplete."""

from __future__ import annotations

import re
import sys
from pathlib import Path


def parse_env(path: Path) -> dict[str, str]:
    values: dict[str, str] = {}
    for line_number, line in enumerate(path.read_text(encoding="utf-8").splitlines(), 1):
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue
        if "=" not in stripped:
            raise ValueError(f"invalid environment entry at line {line_number}")
        key, value = stripped.split("=", 1)
        key, value = key.strip(), value.strip()
        if not re.fullmatch(r"[A-Z][A-Z0-9_]*", key):
            raise ValueError(f"invalid key at line {line_number}")
        if key in values:
            raise ValueError(f"duplicate key {key}")
        values[key] = value
    return values


def fail(message: str) -> int:
    print(f"production environment rejected: {message}", file=sys.stderr)
    return 1


def main() -> int:
    if len(sys.argv) != 2:
        return fail("usage: validate-env.py /path/to/production.env")
    path = Path(sys.argv[1])
    if not path.is_file():
        return fail("environment file is missing")
    values = parse_env(path)

    required = (
        "DB_BOOTSTRAP_PASSWORD",
        "DB_PASSWORD",
        "MIGRATE_DB_PASSWORD",
        "REDIS_PASSWORD",
        "JWT_SECRET",
    )
    for key in required:
        value = values.get(key, "")
        if len(value) < 32 or any(marker in value.lower() for marker in ("change_me", "generate_", "replace_")):
            return fail(f"{key} is missing, too short, or still a placeholder")
    if len({values[key] for key in required}) != len(required):
        return fail("database, Redis, and JWT secrets must be distinct")

    if values.get("IMAGE_REGISTRY") != "ghcr.io/yousef-abdallah-1711/pos":
        return fail("IMAGE_REGISTRY must point to the reviewed Rakmyat POS package")
    if not re.fullmatch(r"sha-[0-9a-f]{40}", values.get("IMAGE_TAG", "")):
        return fail("IMAGE_TAG must be an immutable sha-<40 lowercase hex> release tag")
    expected = {
        "PUBLIC_API_URL": "https://api.pos.rakmyat.online",
        "LANDING_URL": "https://pos.rakmyat.online",
        "DASHBOARD_URL": "https://pos.app.rakmyat.online",
        "SUPER_ADMIN_URL": "https://admin.pos.rakmyat.online",
        "RESTAURANT_DOMAIN": "restaurants.pos.rakmyat.online",
        "PAYMOB_MODE": "live",
    }
    for key, value in expected.items():
        if values.get(key) != value:
            return fail(f"{key} does not match the approved production hostname/mode")

    for key in ("PAYMOB_API_KEY", "PAYMOB_INTEGRATION_ID", "PAYMOB_IFRAME_ID", "PAYMOB_HMAC_SECRET"):
        if not values.get(key, "").strip():
            return fail(f"{key} is required; configure it via the protected server environment")
    for key in ("PAYMOB_INTEGRATION_ID", "PAYMOB_IFRAME_ID"):
        if not values[key].isdigit() or int(values[key]) <= 0:
            return fail(f"{key} must be a positive integer")
    if len(values["PAYMOB_API_KEY"]) < 16 or len(values["PAYMOB_HMAC_SECRET"]) < 16:
        return fail("Paymob API/HMAC credentials are too short")
    if values.get("PRODUCTION_DEPLOYMENT_ENABLED") != "true":
        return fail("PRODUCTION_DEPLOYMENT_ENABLED is not enabled")
    if values.get("PRODUCTION_READINESS_APPROVED") != "true":
        return fail("PRODUCTION_READINESS_APPROVED is not enabled")
    if values.get("MIGRATIONS_BACKWARD_COMPATIBLE") != "true":
        return fail("release has no backward-compatible migration approval")
    print("production environment checks passed")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (OSError, ValueError) as error:
        raise SystemExit(fail(str(error)))
