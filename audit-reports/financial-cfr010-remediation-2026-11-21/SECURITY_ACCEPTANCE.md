# Security and tenant isolation acceptance

**Result: PARTIAL for the candidate.** The repository reconciliation method is read-only and accepts an explicit positive tenant and restaurant scope. It adds no route or public API surface.

## Verified in code and PostgreSQL

- The method rejects missing/non-positive tenant or restaurant identifiers before querying.
- Journal rows and counts are scoped by both `tenant_id` and `restaurant_id`; account, journal-line, and GL relationships are checked for scope consistency.
- Database work uses a `REPEATABLE READ`, read-only transaction. All SQL values are bound parameters. Amount/date comparisons are performed by PostgreSQL without float conversion.
- The hosted workflow provisioned a disposable database and restricted migrator, application, and accounting roles. Its full `go test -tags integration` acceptance step passed.
- The targeted reconciliation integration suite passed on exact SHA `e68d70d45ccf507ee0474449fbf5e3f19a96f76e` in run [44953908598](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44953908598). It includes tenant and sibling-restaurant isolation, orphan/invalid link, scope and anomaly cases.

## Not claimed

- No new handler exists, so authorization behavior for a future diagnostic API is not accepted here. If exposed later, require the existing financial read permission and test unauthorized, missing tenant context, invalid restaurant context, and cross-tenant/sibling-restaurant denial through the live API.
- This run does not establish a complete application-wide RLS or security audit. It verifies the exercised disposable integration suite only.
- No production database or financial record was accessed or modified.
