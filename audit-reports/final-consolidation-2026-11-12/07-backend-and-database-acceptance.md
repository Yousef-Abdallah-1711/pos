# Backend and database acceptance

- Go unit tests, vet, build, PostgreSQL integration, migration forward/rollback, RLS, payment, concurrency, and tenant isolation: NOT RUN locally because Go and a disposable PostgreSQL/Docker Engine are unavailable.
- Hosted check state observed for source PRs is not equivalent to full integrated regression on this candidate.
- No database was modified.
- Required next verification: hosted CI against the exact consolidated head and a disposable PostgreSQL runner.

Hosted CI initially exposed two consolidation defects: duplicate string-length helper declarations and seven SQL-mock tests that had not been updated for pagination count queries. The duplicate implementation was consolidated into one Unicode-rune validator with invalid-rule rejection and unit tests. The seven SQL mocks now expect scoped COUNT and paged list queries. These fixes are queued for the next hosted CI run.
