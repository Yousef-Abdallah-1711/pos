from __future__ import annotations

import argparse
import json
import ssl
import sys
import urllib.error
import urllib.request
from dataclasses import dataclass
from typing import Iterable


SUCCESS_HTTP_STATUSES = {200, 201, 204, 301, 302, 307, 308}


@dataclass
class CheckResult:
    name: str
    url: str
    ok: bool
    status_code: int | None
    message: str


def normalize_base_url(url: str) -> str:
    return url.rstrip("/")


def ensure_api_health_url(api_url: str) -> str:
    normalized = normalize_base_url(api_url)
    if normalized.endswith("/api/v1/health"):
        return normalized
    if normalized.endswith("/api/v1"):
        return f"{normalized}/health"
    return f"{normalized}/api/v1/health"


def fetch_url(url: str, timeout: float, insecure: bool) -> tuple[int, bytes]:
    request = urllib.request.Request(
        url,
        method="GET",
        headers={
            "User-Agent": "pos-saas-deployment-verifier/1.0",
            "Accept": "application/json,text/html,*/*",
        },
    )
    context = None
    if insecure:
        context = ssl.create_default_context()
        context.check_hostname = False
        context.verify_mode = ssl.CERT_NONE

    with urllib.request.urlopen(request, timeout=timeout, context=context) as response:
        return int(response.status), response.read()


def check_http_endpoint(name: str, url: str, timeout: float, insecure: bool) -> CheckResult:
    try:
        status_code, _ = fetch_url(url, timeout, insecure)
        ok = status_code in SUCCESS_HTTP_STATUSES
        message = "OK" if ok else f"Unexpected status {status_code}"
        return CheckResult(name=name, url=url, ok=ok, status_code=status_code, message=message)
    except urllib.error.HTTPError as exc:
        status_code = int(exc.code)
        ok = status_code in SUCCESS_HTTP_STATUSES
        message = "OK" if ok else f"HTTPError {status_code}"
        return CheckResult(name=name, url=url, ok=ok, status_code=status_code, message=message)
    except Exception as exc:  # pragma: no cover
        return CheckResult(name=name, url=url, ok=False, status_code=None, message=str(exc))


def check_api_health(api_url: str, timeout: float, insecure: bool) -> CheckResult:
    health_url = ensure_api_health_url(api_url)
    try:
        status_code, body = fetch_url(health_url, timeout, insecure)
        if status_code not in SUCCESS_HTTP_STATUSES:
            return CheckResult("API Health", health_url, False, status_code, f"Unexpected status {status_code}")

        body_text = body.decode("utf-8", errors="ignore").strip()
        if body_text:
            try:
                parsed = json.loads(body_text)
                if isinstance(parsed, dict):
                    # Best-effort check only; API schemas can vary.
                    summary = ", ".join(sorted(parsed.keys())[:5]) or "json"
                    return CheckResult("API Health", health_url, True, status_code, f"OK (JSON keys: {summary})")
            except json.JSONDecodeError:
                pass

        return CheckResult("API Health", health_url, True, status_code, "OK")
    except urllib.error.HTTPError as exc:
        return CheckResult("API Health", health_url, False, int(exc.code), f"HTTPError {int(exc.code)}")
    except Exception as exc:  # pragma: no cover
        return CheckResult("API Health", health_url, False, None, str(exc))


def print_results(results: Iterable[CheckResult]) -> int:
    failed = False
    print("\nDeployment endpoint verification results:\n")
    for result in results:
        status = "PASS" if result.ok else "FAIL"
        code_part = f" [{result.status_code}]" if result.status_code is not None else ""
        print(f"- {status} {result.name}{code_part}: {result.url} :: {result.message}")
        if not result.ok:
            failed = True

    return 1 if failed else 0


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Verify staging/production public endpoints after deployment."
    )
    parser.add_argument("--api-url", required=True, help="Public API base URL (e.g. https://api.yourdomain.com)")
    parser.add_argument("--landing-url", help="Landing site URL (e.g. https://www.yourdomain.com)")
    parser.add_argument("--dashboard-url", help="Dashboard URL (e.g. https://dashboard.yourdomain.com)")
    parser.add_argument("--admin-url", help="Super admin URL (e.g. https://admin.yourdomain.com)")
    parser.add_argument(
        "--restaurant-url",
        help="Sample tenant restaurant URL (e.g. https://demo.yourdomain.com or https://slug.yourdomain.com/en)",
    )
    parser.add_argument("--timeout", type=float, default=10.0, help="HTTP timeout in seconds (default: 10)")
    parser.add_argument(
        "--insecure",
        action="store_true",
        help="Disable TLS verification (use only for temporary staging/self-signed cert checks).",
    )
    args = parser.parse_args()

    results: list[CheckResult] = []
    results.append(check_api_health(args.api_url, args.timeout, args.insecure))

    optional_checks = [
        ("Landing", args.landing_url),
        ("Dashboard", args.dashboard_url),
        ("Super Admin", args.admin_url),
        ("Restaurant Website", args.restaurant_url),
    ]
    for name, url in optional_checks:
        if url:
            results.append(check_http_endpoint(name, normalize_base_url(url), args.timeout, args.insecure))

    if len(results) == 1:
        print("Only API health was checked (no frontend URLs provided).")

    return print_results(results)


if __name__ == "__main__":
    sys.exit(main())
