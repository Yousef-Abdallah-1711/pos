# Main guardrail and deployment workflow

- `python scripts/check_runtime_localhost.py`: PASS — no unallowlisted hardcoded localhost URLs in runtime code.
- Checkout framing policy in `backend/internal/middleware/cache.go` remains environment-configured and fails closed for missing/invalid origins; exact loopback allowlisting is development-only.
- Security regression tests from the remediation candidate remain present.
- PR #20 release workflow was changed to `workflow_dispatch`. CI completion no longer starts image publishing or production deployment automatically. No workflow was dispatched.
- `git diff --check`: PASS at source checkpoint.
