# Infrastructure and database

**Isolated stack:** seven services from main SHA `a6a2c6700a5fe883256deb1e5a3893108bd9865c`: PostgreSQL, Redis, API, landing, dashboard, restaurant website, super-admin.

- PostgreSQL and Redis reported healthy; API reported healthy.
- API health endpoint returned HTTP 200 with database status `ok`.
- API startup applied migrations through the known application migration sequence; the database migration ledger's latest row is `999_theme_migration_legacy.sql` at 2026-11-11 16:16 UTC. The ledger is filename-based, not a numeric version column.
- `pos_app` and `pos_migrator` exist, can log in, and are not superusers, database creators, or role creators.
- Direct frontend configuration used `http://localhost:18080`, while the isolated API was published at `18082`. The first browser login request failed with connection refused. A temporary isolated TCP forward allowed browser/API acceptance; it was stopped after testing. No application configuration was changed.
- The isolated database is disposable and uses tmpfs. Persistence across a PostgreSQL restart was NOT RUN to avoid destroying that database or affecting any existing database.
- The earlier `pos_migrator` failure was reproduced as an initialization-order issue in the prior tmpfs demo stack: PostgreSQL restarted empty, the one-shot role initializer had exited, and direct API startup bypassed it. Rerunning the intended initializer on that fresh disposable database restored the restricted roles and API health. Persistent databases were not altered.
- All pre-existing Docker stacks remain running; no volumes were removed.

