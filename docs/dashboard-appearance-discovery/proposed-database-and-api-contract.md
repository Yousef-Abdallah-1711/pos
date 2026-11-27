# Proposed database and API contract (future; not implemented)

Proposed resource: a tenant-keyed workspace appearance document with schema version, published revision, validated override payload, updated-at and actor metadata; add audit history with transactionally consistent updates. Enforce forced RLS using the `app.tenant_id` context established by `backend/internal/pkg/database/tenant_tx.go`; derive tenant identity from authenticated middleware, never from request JSON. Grant application roles narrowly and test cross-tenant reads/writes with a real database. No migration number or production migration has been chosen.

Candidate operations: authenticated GET resolved configuration; tenant-admin PUT with `expectedRevision` and validated overrides; tenant-admin reset with expected revision. Return conflict on stale revision; invalid keys, colors, fonts, sizes or contrast fail validation. Distinguish database failure from missing overrides so an outage never silently overwrites data. Decide whether drafts need persistence based on actual collaboration requirements; local preview alone is simpler but cannot replace a published DB record. Include audit actor, before/after revision, safe rollback and cache invalidation scoped to tenant. Existing `/api/v1/users/{userId}/settings/*` remains personal; `/api/v1/admin/themes/*` remains storefront scope unless an independently reviewed migration explicitly changes that boundary.

## Candidate wire format for implementation review

These examples are proposals only; route names, permission names and status mappings must be reconciled with existing backend conventions before implementation.

```json
{
  "schemaVersion": 1,
  "revision": 7,
  "source": "tenant",
  "overrides": {
    "colors": {"primary": {"light": "#A13224", "dark": "#E17861"}},
    "typography": {"bodyFont": "manrope"}
  }
}
```

Read derives tenant from authenticated server context and returns either the tenant document or an explicit default source with its revision. A proposed write sends `{ "expectedRevision": 7, "overrides": { ... } }`; the server validates the entire candidate and atomically updates only if the revision matches, increments it, and records an audit entry in the same tenant-scoped transaction. An absent initial row needs an explicit creation revision rule. Reset is a distinct authorized action with the same concurrency check; never accept tenant ID, CSS, URLs or arbitrary font names from the body. A conflict should return the current revision without applying any part of the request.

Before coding, choose whether a GET returns raw overrides, fully resolved tokens, or both; specify ETag/cache headers and whether override values may differ per mode. Define missing/invalid legacy settings, schema migrations, API timeouts, authorization errors, contrast validation against inherited defaults, and restore behavior. Never deploy the proposed schema or migrate stored customer values without an explicit implementation review and recovery plan.
