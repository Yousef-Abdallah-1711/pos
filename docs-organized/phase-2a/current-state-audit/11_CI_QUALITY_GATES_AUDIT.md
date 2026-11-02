# CI quality gates audit

## Workflows found

`.github/workflows/runtime-localhost-guardrail.yml` runs `python scripts/check_runtime_localhost.py` on PR and pushes to main/master. `url-domain-config-smoke.yml` installs Python 3.11, Go from backend/go.mod, pnpm 9/Node 20 and frozen frontend dependencies, then runs `scripts/smoke_check_url_domain_config.py --skip-guardrail`.

| Gate | CI evidence | Status |
|---|---|---|
| Go build | None | Missing |
| Go tests | None | Missing |
| Go vet | None | Missing |
| Frontend typecheck | None | Missing |
| Frontend unit tests | None | Missing |
| Frontend production builds | None | Missing |
| Security smoke | None | Missing |
| Migration validation | None | Missing |
| Docker build | None | Missing |
| Dependency vulnerability scan | None | Missing |
| Secret scan | None | Missing |
| URL/runtime config guardrail | Two workflows | Present |

No evidence in checked-in workflows shows required branch protection or mandatory status checks; GitHub settings were not queried. Existing scripts do not prove the repository requires passing CI.

## Minimum gate proposal

Add deterministic Go build/vet/test job; install frozen frontend workspace and execute typechecks plus test suites; restore broken suites before making them blocking; build each app with no bypass; run isolated disposable security smoke and migration forward/rollback validation; add dependency/secret scanning. Publish reports as artifacts. Set required checks only after the suite is stable; do not conceal failures with allow-failure or ignored tests.
