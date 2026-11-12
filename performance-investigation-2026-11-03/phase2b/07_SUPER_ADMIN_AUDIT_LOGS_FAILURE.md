# Super Admin Audit Logs Functional Failure

Classification: REPRODUCIBLE FUNCTIONAL FAILURE. This is not a performance bottleneck finding.

In the authenticated Super Admin session, /dashboard/audit-logs returned HTTP 200 and the API returned 9 synthetic audit rows. Two rows had details: null; 7 had object values. The frontend at frontend/apps/super-admin/src/app/dashboard/audit-logs/page.tsx:162 calls Object.keys(log.details).length while rendering rows. Null raises TypeError: Cannot convert undefined or null to object, and the browser displays its generic application error. The first authenticated navigation reproduced this crash. A later navigation after the session expired redirected to login, so a second valid-session run was not captured.

The minimal remediation proposal is to make rendering null-safe and decide whether absent details should hide the block or display {}. Alternatively normalize API output only if that contract is intended. Add a regression case with null details and an authenticated page smoke test. No source changes were made. Expected impact: restore route availability; performance impact is not estimated. Tenant isolation implications are none for the null guard. Rollback: revert the isolated UI change if needed.

Sanitized exception evidence is in 16_BROWSER_TRACES/audit-logs-exception-sanitized.log and audit-logs-exception-snapshot-sanitized.yml. The original browser-generated files are preserved in .playwright-mcp; no credentials or row contents are included in these report excerpts.
