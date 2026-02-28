from __future__ import annotations

import re
import sys
import os
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]
if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")

SCAN_ROOTS = [
    REPO_ROOT / "backend" / "cmd" / "api",
    REPO_ROOT / "backend" / "internal",
    REPO_ROOT / "frontend" / "apps" / "restaurant-website",
    REPO_ROOT / "frontend" / "apps" / "dashboard",
    REPO_ROOT / "frontend" / "apps" / "landing",
    REPO_ROOT / "frontend" / "apps" / "super-admin",
]

SOURCE_EXTENSIONS = {".go", ".ts", ".tsx", ".js", ".mjs", ".cjs"}
EXCLUDED_DIR_NAMES = {"node_modules", ".next", "dist", "coverage", ".turbo", "out", "build"}
EXCLUDED_FILE_SUFFIXES = (".test.ts", ".test.tsx", ".spec.ts", ".spec.tsx", "_test.go")
EXCLUDED_FILE_NAMES = {"README.md"}
EXCLUDED_PATH_PARTS = {"__tests__"}

# Intentional dev-only localhost/default behavior or local-development support.
ALLOWLIST_FILES = {
    Path("backend/cmd/api/main.go"),
    Path("backend/internal/config/runtime_env.go"),
    Path("backend/internal/middleware/cors.go"),
    Path("backend/internal/service/webhook_retry_worker.go"),
    Path("backend/internal/handler/http/mock_paymob_handler.go"),
    Path("backend/internal/handler/http/super_admin_handler.go"),
    Path("frontend/apps/restaurant-website/lib/config/api.ts"),
    Path("frontend/apps/restaurant-website/lib/subdomain.ts"),
    Path("frontend/apps/restaurant-website/middleware.ts"),
    Path("frontend/apps/dashboard/src/lib/config.ts"),
    Path("frontend/apps/landing/src/lib/config.ts"),
    Path("frontend/apps/super-admin/src/lib/config.ts"),
    Path("frontend/apps/dashboard/next.config.js"),
    Path("frontend/apps/landing/next.config.js"),
    Path("frontend/apps/restaurant-website/next.config.js"),
}

LOCALHOST_URL_RE = re.compile(
    r"(?i)\b(?:https?|wss?|ws)://(?:localhost|127\.0\.0\.1)(?::\d+)?(?:[/?#][^\s\"'`]*)?"
)


def should_skip_file(path: Path) -> bool:
    if path.suffix not in SOURCE_EXTENSIONS:
        return True
    if path.name in EXCLUDED_FILE_NAMES:
        return True
    if path.name.endswith(EXCLUDED_FILE_SUFFIXES):
        return True
    rel = path.relative_to(REPO_ROOT)
    if any(part in EXCLUDED_DIR_NAMES for part in rel.parts):
        return True
    if any(part in EXCLUDED_PATH_PARTS for part in rel.parts):
        return True
    if "docs-organized" in rel.parts:
        return True
    return False


def is_legacy_build_tag_file(text: str) -> bool:
    return "//go:build legacy" in "\n".join(text.splitlines()[:3])


def iter_source_files():
    for root in SCAN_ROOTS:
        if not root.exists():
            continue
        for dirpath, dirnames, filenames in os.walk(root):
            # Prune heavy/generated directories early for speed.
            dirnames[:] = [d for d in dirnames if d not in EXCLUDED_DIR_NAMES and d not in EXCLUDED_PATH_PARTS]
            current_dir = Path(dirpath)
            for filename in filenames:
                path = current_dir / filename
                if should_skip_file(path):
                    continue
                yield path


def main() -> int:
    violations: list[tuple[str, int, str]] = []

    for path in iter_source_files():
        rel = path.relative_to(REPO_ROOT)
        text = path.read_text(encoding="utf-8", errors="ignore")

        if is_legacy_build_tag_file(text):
            continue
        if rel in ALLOWLIST_FILES:
            continue

        for lineno, line in enumerate(text.splitlines(), start=1):
            stripped = line.strip()
            if stripped.startswith(("//", "*", "/*", "#")):
                continue
            if LOCALHOST_URL_RE.search(line):
                violations.append((str(rel).replace("\\", "/"), lineno, line.strip()))

    if violations:
        print("Runtime localhost URL guardrail failed. Found hardcoded localhost URL(s) in runtime code:\n")
        for rel, lineno, line in violations:
            safe_line = line.encode("utf-8", errors="replace").decode("utf-8", errors="replace")
            print(f"- {rel}:{lineno}: {safe_line}")
        print(
            "\nFix the code to use env-driven config with dev-only fallbacks, or add a narrow allowlist entry if intentional."
        )
        return 1

    print("Runtime localhost URL guardrail passed (no unallowlisted hardcoded localhost URLs in runtime code).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
