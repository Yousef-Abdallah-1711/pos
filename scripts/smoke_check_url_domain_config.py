from __future__ import annotations

import argparse
import os
import subprocess
import sys
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]

PROD_ENV_FILES = [
    REPO_ROOT / "backend" / ".env.production.example",
    REPO_ROOT / "frontend" / "apps" / "dashboard" / ".env.production.example",
    REPO_ROOT / "frontend" / "apps" / "landing" / ".env.production.example",
    REPO_ROOT / "frontend" / "apps" / "restaurant-website" / ".env.production.example",
    REPO_ROOT / "frontend" / "apps" / "super-admin" / ".env.production.example",
]

LOCAL_ENV_FILES = [
    REPO_ROOT / "backend" / ".env.example",
    REPO_ROOT / "frontend" / "apps" / "dashboard" / ".env.example",
    REPO_ROOT / "frontend" / "apps" / "landing" / ".env.example",
    REPO_ROOT / "frontend" / "apps" / "restaurant-website" / ".env.example",
    REPO_ROOT / "frontend" / "apps" / "super-admin" / ".env.example",
]


def run_step(name: str, cmd: list[str], cwd: Path | None = None) -> None:
    if os.name == "nt" and cmd and cmd[0] == "pnpm":
        cmd = ["pnpm.cmd", *cmd[1:]]
    print(f"\n==> {name}")
    print(" ".join(cmd))
    subprocess.run(cmd, cwd=str(cwd or REPO_ROOT), check=True)


def verify_env_templates() -> None:
    print("\n==> Static env template checks")
    for path in LOCAL_ENV_FILES:
        if not path.exists():
            raise RuntimeError(f"Missing local env template: {path}")
        content = path.read_text(encoding="utf-8", errors="ignore")
        active_lines = [
            line.strip()
            for line in content.splitlines()
            if line.strip() and not line.lstrip().startswith("#")
        ]
        active_blob = "\n".join(active_lines).lower()
        if "localhost" not in active_blob and "127.0.0.1" not in active_blob:
            raise RuntimeError(f"Local env template does not appear to contain local defaults: {path}")
        print(f"PASS local template: {path.relative_to(REPO_ROOT)}")

    for path in PROD_ENV_FILES:
        if not path.exists():
            raise RuntimeError(f"Missing production env template: {path}")
        content = path.read_text(encoding="utf-8", errors="ignore")
        active_lines = [
            line.strip()
            for line in content.splitlines()
            if line.strip() and not line.lstrip().startswith("#")
        ]
        active_blob = "\n".join(active_lines).lower()
        if "localhost" in active_blob or "127.0.0.1" in active_blob:
            raise RuntimeError(f"Production env template contains localhost: {path}")
        print(f"PASS production template: {path.relative_to(REPO_ROOT)}")


def main() -> int:
    parser = argparse.ArgumentParser(description="Smoke-check URL/domain resolver and config hardening across apps.")
    parser.add_argument("--skip-guardrail", action="store_true", help="Skip runtime localhost guardrail script")
    parser.add_argument("--skip-tests", action="store_true", help="Skip backend/frontend unit test commands")
    args = parser.parse_args()

    if not args.skip_guardrail:
        run_step("Runtime localhost guardrail", [sys.executable, "scripts/check_runtime_localhost.py"])

    verify_env_templates()
    if not args.skip_tests:
        run_step(
            "Backend config validation tests",
            ["go", "test", "./internal/config", "-run", "Test(CurrentEnv|ResolveRuntime|ValidateProductionURLConfig)"],
            cwd=REPO_ROOT / "backend",
        )

        run_step(
            "Landing resolver unit tests",
            [
                "pnpm",
                "--dir",
                "frontend/apps/landing",
                "test:config-resolver",
            ],
        )

        run_step(
            "Dashboard resolver unit tests",
            [
                "pnpm",
                "--dir",
                "frontend/apps/dashboard",
                "exec",
                "vitest",
                "run",
                "--config",
                "vitest.config.unit.mts",
                "src/lib/config.test.ts",
            ],
        )

        run_step(
            "Super-admin resolver unit tests",
            ["pnpm", "--dir", "frontend/apps/super-admin", "exec", "vitest", "run", "src/lib/config.test.ts"],
        )

        run_step(
            "Restaurant website resolver unit tests",
            ["pnpm", "--dir", "frontend/apps/restaurant-website", "exec", "jest", "lib/config/api.test.ts", "--runInBand"],
        )

    print("\nURL/domain configuration smoke check passed.")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except subprocess.CalledProcessError as exc:
        print(f"\nSmoke check failed: command exited with {exc.returncode}", file=sys.stderr)
        raise SystemExit(exc.returncode)
    except Exception as exc:  # pragma: no cover
        print(f"\nSmoke check failed: {exc}", file=sys.stderr)
        raise SystemExit(1)
