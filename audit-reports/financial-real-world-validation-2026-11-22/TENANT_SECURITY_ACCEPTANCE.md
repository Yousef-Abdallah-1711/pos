# Tenant and authorization acceptance

**Synthetic acceptance:** PASS for exercised checks.

- Cashier financial-screen access and journal posting were denied; a server-side API check returned HTTP 403 for posting.
- Accountant posting was allowed for the balanced test journal; duplicate posting was rejected.
- Cross-tenant report/account access returned not found/denied in API and browser journeys. The order acceptance's foreign tenant read returned HTTP 404. Browser test 11 explicitly verified a second tenant could not read or export the first tenant account.
- AR and payroll acceptance logs recorded foreign-tenant denial.

Evidence: `financial-live-api-tests.jsonl`, `financial-browser-tests.txt`, `financial-browser-junit.xml`. This is bounded scenario evidence, not a complete security audit, penetration test, or proof for every endpoint and role. No production credentials or data were used.
