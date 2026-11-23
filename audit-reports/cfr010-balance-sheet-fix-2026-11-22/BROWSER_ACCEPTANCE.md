# Browser Acceptance

Authenticated Playwright Balance Sheet test run against the disposable API (`localhost:8180`) and dashboard (`localhost:3002`):

`npx playwright test --config playwright.financial.config.ts --grep 'Balance Sheet values match the SQL oracle' --trace on --output test-results/cfr010-bs-rootcause-cache-fix-verified`

Result: Balance Sheet test passed; two additional targeted tests passed for duplicate journal-post rejection and second-tenant isolation. It asserted API totals against the synthetic SQL oracle, English and Arabic formatting, LTR/RTL direction, 1440×900 desktop and 375×812 mobile viewports, no mobile overflow, correct negative contra-asset formatting, visible red alert while unreconciled, and absence of the PDF export button while residual is nonzero. The underlying warning was not hidden.

Screenshots: `browser-artifacts/en-desktop-balance-sheet.png`, `en-mobile-balance-sheet.png`, `ar-desktop-balance-sheet.png`, and `ar-mobile-balance-sheet.png`. Playwright trace was captured in the task-local Playwright output. A content scan found synthetic login identity, password, and cookie/auth headers in the archive, so it is deliberately excluded from the PR. The local trace remains available for private inspection and must be scrubbed before sharing.

The prior run also reproduced stale `RECONCILED` status after a posted journal. The post-fix run returns 55.00 assets / 5.00 liabilities / 0.00 equity and keeps the warning/PDF restriction. Browser coverage here is scoped to Balance Sheet; it is not a full acceptance of unrelated financial flows.
