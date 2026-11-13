# 06 — Database safety evidence

- Container: financial-phase00-disposable-20261111.
- PostgreSQL image: postgres:15-alpine.
- Label: codex.disposable=true; codex.task=financial-phase00.
- Data: tmpfs /var/lib/postgresql/data; published only to 127.0.0.1:15444.
- Test database: pos_phase00_test_disposable; connection used in Go container via host.docker.internal:15444.
- Migration history: 126 SQL files applied through 999_theme_migration_legacy.sql.
- Schema inspected: journal debit/credit NUMERIC(15,2), GL running balance NUMERIC(15,2), attachments JSONB, entry_type VARCHAR(20). No core financial table RLS policies in this migration snapshot.
- One seeded test tenant and restaurant. Integration inserted a balanced test journal and two ledger lines; cleanup confirmed zero journal and ledger rows.
- No primary/shared/protected database was queried or modified. No application migration was added.