# Backend and database acceptance

- Go unit tests, vet, build, PostgreSQL integration, migration forward/rollback, RLS, payment, concurrency, and tenant isolation: NOT RUN locally because Go and a disposable PostgreSQL/Docker Engine are unavailable.
- Hosted check state observed for source PRs is not equivalent to full integrated regression on this candidate.
- No database was modified.
- Required next verification: hosted CI against the exact consolidated head and a disposable PostgreSQL runner.
