# 08 — Transaction atomicity and recovery (PROPOSED)

Inspection indicates GL insertion and posted state share a DB transaction while draft creation and posting are separate. Failure injection was NOT TESTED. No distributed exactly-once guarantee is claimed.

Commit source event and outbox row in one PostgreSQL transaction. Deliver at least once. Worker uses unique event identity/durable inbox. In one DB transaction write journal, lines, GL and posting outcome. Never hold a DB transaction across provider calls. Use provider idempotency where available and reconcile after timeout.

Recovery: before source commit, caller retries; after commit before publish, outbox republishes; worker crash before commit rolls back; crash after commit before ack replays into existing effect; provider success/local timeout queries provider by stable reference before retry; permanent validation errors become visible dead letters; orphan automatic drafts are aged, alerted and retried only after checking no posted effect. Fixes create new versioned events, not silent edits.

Specify lease, retry/backoff, DLQ owner, retention and recovery objectives before implementation. Test failure at each boundary and prove one effect plus visible status.
