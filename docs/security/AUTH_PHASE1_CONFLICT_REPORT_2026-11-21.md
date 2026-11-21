# Phase 1 Synthetic Identity Conflict Report

**Purpose:** Exercise the proposed legacy-user backfill/linking rules against representative conflicts.
**Data source:** Fully synthetic fixture specification. No staging or production rows, email addresses, password hashes, or database values were queried.
**Baseline schema:** `users` in migration `003_create_users.sql`, exact-case `(tenant_id,email)` uniqueness, normalized case-insensitive login lookup in the auth repository, and `user_roles` in `103_create_user_roles_table.sql`.

## Fixture Summary

Eight synthetic `users` rows cover three duplicate normalized-email groups plus two individual edge cases. The labels below are fixture-only identifiers; no actual hashes or customer data are represented.

| Fixture group | Rows | Current-schema condition | Safe backfill result | Conflict action |
|---|---:|---|---|---|
| `SYN-A` | 2 | Same normalized email in tenants 101 and 202; distinct credential markers | Two unverified identities and two memberships | No merge. A password for one row must not authorize the other. Require proof of each account plus verified email before linking. |
| `SYN-B` | 2 | Same normalized email in tenants 303 and 404; same credential marker | Two unverified identities and two memberships | No merge. Password equality is not proof that records are the same person. |
| `SYN-C` | 2 | Same tenant 505; `Case.Manager@example.test` and `case.manager@example.test`; exact-case unique key permits both | Two unverified identities and two memberships | Flag a normalized-email collision. Keep separately scoped and block canonical linking until reviewed. |
| `SYN-D` | 1 | Active user with `restaurant_id IS NULL` | One unverified identity and one unresolved membership | No restaurant grant until the legacy null meaning is approved and source behavior is verified. |
| `SYN-E` | 1 | Inactive user row with historical actor/audit references | One unverified identity and one suspended membership | Preserve row and FKs; do not grant login access. |

A separate synthetic employee record uses an email matching `SYN-A` but has no `users` foreign key. It remains an HR record and is not linked to either identity by email.

A separate role fixture gives one legacy user `users.role='manager'` while its `user_roles` mapping resolves to `cashier`. It is marked as a permission-source conflict; backfill must not choose the more privileged result automatically.

## Expected Aggregate Output

For this fixture only:

- Legacy user rows: **8**
- New identity rows after one-to-one backfill: **8**
- New membership rows after one-to-one backfill: **8**
- Automatic identity merges: **0**
- Automatically linked cross-tenant memberships: **0**
- Normalized-email duplicate groups: **3**
- Duplicate groups with divergent credential markers: **1**
- Case-only same-tenant collision groups: **1**
- Unresolved null restaurant scopes: **1**
- Role-source conflicts requiring manual mapping: **1**
- Employee-to-user links inferred from email alone: **0**

These counts are an expected fixture result, not production or staging counts. The feature flag remains off; legacy authentication continues while the new membership rows are compared in isolation.

## Migration Invariants

1. Every legacy user row maps to exactly one new identity and one membership on the first backfill pass.
2. Re-running the backfill produces the same row counts and mapping, using `legacy_user_id` uniqueness as the idempotency key.
3. No two legacy identities merge because their email or password markers match.
4. No new membership is active until the identity email is verified and its account-link proof is accepted.
5. A membership with unknown restaurant scope or unresolved role conflict is not enabled in the new login path.
6. Every `users.id`, employee reference, role reference, and audit actor reference remains unchanged.
7. A password proof for one legacy row produces no selection proof for a different row with a distinct credential.

## Implementation Gate

Implement this fixture in disposable PostgreSQL integration tests before any production backfill. The test must verify row counts, rerun idempotency, preserved foreign keys, deny-by-default RLS, and rollback to the old auth read path. Independent security review must approve the conflict policy and its database constraints before enabling feature-flagged reads.
