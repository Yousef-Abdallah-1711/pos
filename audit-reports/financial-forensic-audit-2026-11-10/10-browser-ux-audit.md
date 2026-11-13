# Browser UX investigation

**Status: PARTIAL.** Playwright MCP ran in isolated contexts, not claimed as production browser coverage.

- Synthetic registration returned HTTP 201 and separate login succeeded; credentials omitted.
- Requested accounting and financial report pages were opened, plus suppliers and VAT. Most began empty.
- Journal form displayed EGP 100 debit, EGP 100 credit and Balanced. Save failed HTTP 400 because attachments was absent. Authenticated browser request with attachments: [] succeeded; UI then posted it.
- P&L rendered revenue 100 and net profit 100. Export PDF created profit-loss-2026-10-12-to-2026-11-10.pdf.
- Arabic locale and mobile viewport 390x844 were captured in screenshots folder.
- P&L page had zero console errors at one snapshot; journal interaction logged API 400 errors. Browser session accumulated warnings.
- One route measurement pass ranged 2.7–10.6 seconds; invoice content needed extra wait in one observation.

Evidence: screenshots, snapshots, registration-result.yml, downloaded PDF. Accessibility, keyboard-only operation, all filters, back/forward, multiple tabs and large tables remain unverified.
