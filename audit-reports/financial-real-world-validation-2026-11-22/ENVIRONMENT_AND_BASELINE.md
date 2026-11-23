# Environment and baseline

**Run date:** 2026-11-22  
**Validation status:** synthetic local acceptance only; not production acceptance.

## Source version

- Repository: `Yousef-Abdallah-1711/pos`
- Tested unmerged integration candidate: `98aa00951756a2c89e15678888eeb67b2d66dd5e`
- Candidate base at assembly: `4a630b35f90f9d475d3d0ba3472877058a0e2952`
- Remote `origin/main` refreshed after testing: `c9a57c1db92ce24acdd3ced0c55650a3ad031d90` (documentation-only PR #106). Application code tested was unchanged by this later main commit.
- The candidate contains local integration of PR heads #108, #109, and #110 for testing. It was not pushed or deployed.

## Disposable environment

The isolated Compose project `cfr010-realworld-20261122` ran PostgreSQL, Redis, Go API, dashboard, restaurant storefront, and Super Admin. The backend used mock payment mode and enforced RBAC. The test API used the disposable `servio_test` database. PostgreSQL `pg_stat_statements` was enabled for the benchmark after its initial absence prevented the first 10,000-journal run from reaching measurements. No staging or production database or payment credentials were used.

Artifacts: `compose.realworld.override.yml`, `financial-live-api-tests.jsonl`, `financial-acceptance-tests.jsonl`, `financial-browser-tests.txt`, `financial-browser-junit.xml`, and `financial-benchmark-results-1000.json` in this directory.

## Acceptance boundary

Only synthetic workflows described in the evidence set passed. Original operating-data reconciliation is **NOT ESTABLISHED**. OPEN-01, EARN-01, and CLOSE-01 remain **PENDING**. CFR-010 remains **OPEN** and financial production remains **NO-GO**. This environment proves no deployment readiness or real-world capacity.
