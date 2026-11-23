# Security Acceptance

- The report request is authenticated as an accountant in the synthetic browser fixture.
- SQL and report calculations scope tenant and restaurant identifiers. Disposable PostgreSQL integration tests add a second tenant and a sibling restaurant and verify their ledger rows do not change the original report.
- The preceding authenticated API acceptance verified a second-tenant denial and duplicate paid-callback handling. No private credentials or fixture JSON are included in this evidence directory.
- This patch changes only report cache usage; it adds no new input, SQL, route, or authorization surface.

These results establish the tested synthetic scopes only. They are not a production security certification.
