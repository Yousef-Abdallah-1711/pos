#!/usr/bin/env python3
"""
Capture cold vs warm timings for dashboard reports endpoints.

Supports either:
- `--token <jwt>`
- or local JWT generation via `--jwt-secret` (+ tenant/user/restaurant claims)

Examples:
  python scripts/capture_reports_cache_timings.py --token "$env:REPORTS_TEST_TOKEN"

  python scripts/capture_reports_cache_timings.py `
    --jwt-secret "your-super-secret-jwt-key-change-in-production-2024" `
    --tenant-id 8 --restaurant-id 8 --user-id 7 --email "owner@example.com"
"""

from __future__ import annotations

import argparse
import base64
import datetime as dt
import hmac
import hashlib
import json
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from typing import Any, Dict, List, Tuple


def b64url(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode("ascii")


def sign_hs256(payload: Dict[str, Any], secret: str) -> str:
    header = {"alg": "HS256", "typ": "JWT"}
    header_b64 = b64url(json.dumps(header, separators=(",", ":"), ensure_ascii=False).encode("utf-8"))
    payload_b64 = b64url(json.dumps(payload, separators=(",", ":"), ensure_ascii=False).encode("utf-8"))
    signing_input = f"{header_b64}.{payload_b64}".encode("ascii")
    sig = hmac.new(secret.encode("utf-8"), signing_input, hashlib.sha256).digest()
    return f"{header_b64}.{payload_b64}.{b64url(sig)}"


def build_local_jwt(
    secret: str,
    user_id: int,
    tenant_id: int,
    restaurant_id: int,
    email: str,
    role: str,
    ttl_hours: int,
) -> str:
    now = int(time.time())
    payload = {
        "user_id": user_id,
        "email": email,
        "tenant_id": tenant_id,
        "restaurant_id": restaurant_id,
        "role": role,
        "iat": now,
        "exp": now + (ttl_hours * 3600),
    }
    return sign_hs256(payload, secret)


def request_json(url: str, token: str, timeout: float = 20.0) -> Tuple[int, bytes, Dict[str, Any], float]:
    req = urllib.request.Request(
        url,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/json",
        },
        method="GET",
    )
    started = time.perf_counter()
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            body = resp.read()
            elapsed_ms = (time.perf_counter() - started) * 1000.0
            payload = json.loads(body.decode("utf-8"))
            return resp.status, body, payload, elapsed_ms
    except urllib.error.HTTPError as e:
        body = e.read() if e.fp else b""
        elapsed_ms = (time.perf_counter() - started) * 1000.0
        try:
            payload = json.loads(body.decode("utf-8")) if body else {}
        except Exception:
            payload = {"raw": body.decode("utf-8", errors="replace")}
        return e.code, body, payload, elapsed_ms


def generated_at_from_payload(payload: Dict[str, Any]) -> str | None:
    if not isinstance(payload, dict):
        return None
    data = payload.get("data")
    if isinstance(data, dict):
        value = data.get("generated_at")
        return value if isinstance(value, str) else None
    return None


def build_endpoint_urls(base_url: str, date_from: str, date_to: str) -> List[Tuple[str, str]]:
    q = urllib.parse.urlencode({"date_from": date_from, "date_to": date_to})
    sales_q = urllib.parse.urlencode({"date_from": date_from, "date_to": date_to, "granularity": "day"})
    return [
        ("overview", f"{base_url}/reports/overview?{q}"),
        ("sales", f"{base_url}/reports/sales?{sales_q}"),
        ("profitability", f"{base_url}/reports/profitability?{q}"),
        ("inventory", f"{base_url}/reports/inventory?{q}"),
    ]


def main() -> int:
    parser = argparse.ArgumentParser(description="Capture cold vs warm timings for reports endpoints")
    parser.add_argument("--api-base", default="http://localhost:8080/api/v1", help="API base URL (default: local backend)")
    parser.add_argument("--date-from", default=(dt.date.today() - dt.timedelta(days=29)).isoformat())
    parser.add_argument("--date-to", default=dt.date.today().isoformat())
    parser.add_argument("--token", default="", help="Existing JWT bearer token")
    parser.add_argument("--jwt-secret", default="", help="Generate a local JWT using this secret (HS256)")
    parser.add_argument("--user-id", type=int, default=7)
    parser.add_argument("--tenant-id", type=int, default=8)
    parser.add_argument("--restaurant-id", type=int, default=8)
    parser.add_argument("--email", default="reports-benchmark@local.test")
    parser.add_argument("--role", default="owner")
    parser.add_argument("--ttl-hours", type=int, default=24)
    parser.add_argument("--repeat", type=int, default=2, help="Number of requests per endpoint (min 2 recommended)")
    args = parser.parse_args()

    token = args.token.strip()
    if not token:
        if not args.jwt_secret:
            print("ERROR: provide --token or --jwt-secret", file=sys.stderr)
            return 2
        token = build_local_jwt(
            secret=args.jwt_secret,
            user_id=args.user_id,
            tenant_id=args.tenant_id,
            restaurant_id=args.restaurant_id,
            email=args.email,
            role=args.role,
            ttl_hours=args.ttl_hours,
        )

    base_url = args.api_base.rstrip("/")
    endpoints = build_endpoint_urls(base_url, args.date_from, args.date_to)

    print("# Reports Endpoint Timing Capture")
    print(f"- API Base: `{base_url}`")
    print(f"- Date Range: `{args.date_from}` -> `{args.date_to}`")
    print(f"- Repeat Count: `{args.repeat}`")
    print("")
    print("| Endpoint | Run | Status | ms | Payload Bytes | generated_at | Cache-Hit Inference |")
    print("|---|---:|---:|---:|---:|---|---|")

    failures = 0
    for endpoint_name, url in endpoints:
        prev_generated = None
        for i in range(1, max(2, args.repeat) + 1):
            status, body, payload, elapsed_ms = request_json(url, token)
            generated_at = generated_at_from_payload(payload)
            cache_inference = ""
            if i == 1:
                cache_inference = "cold (expected miss)"
            else:
                if generated_at and prev_generated and generated_at == prev_generated:
                    cache_inference = "likely HIT (generated_at unchanged)"
                elif generated_at and prev_generated and generated_at != prev_generated:
                    cache_inference = "miss/rebuild (generated_at changed)"
                else:
                    cache_inference = "unknown"
            prev_generated = generated_at or prev_generated
            print(
                f"| {endpoint_name} | {i} | {status} | {elapsed_ms:.1f} | {len(body)} | "
                f"{generated_at or '-'} | {cache_inference} |"
            )
            if status != 200:
                failures += 1

    print("")
    if failures:
        print(f"FAILURES: {failures} non-200 responses")
        return 1
    print("OK: all requests returned 200")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

