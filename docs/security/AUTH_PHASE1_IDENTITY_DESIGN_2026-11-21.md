# SERVIO Phase 1 Identity and Membership Design Proposal

**Status:** Review draft. This is a design proposal, not an applied schema change.
**Baseline:** `origin/main` `c527e08bb47ec54e4d2f6fd14ecbd9b30e7cadc7`.
**Related remediation:** PR #105, `fix(auth): bind tenant choices to verified accounts` (owner-directed self-review; CI pending).
**Data boundary:** No production or staging data was queried for this proposal. Conflict cases below are synthetic.

## 1. Phase 1 Objective

Give one verified person identity independently authorized memberships in multiple tenants. Login must list only eligible memberships after credential verification, then create a session bound to one membership, tenant, and restaurant scope. A membership’s roles and branch grants must never carry into another membership.

The current task authorizes multiple memberships, login-time tenant selection, and a minimal email verification path needed to establish one verified identity across tenants. In-session tenant switching, invitations, account recovery, and MFA remain out of scope. Preserve the existing staff session policy (session-only cookie, 30-minute idle limit, eight-hour absolute limit), Super Admin identity/session boundary, legacy user IDs, and existing audit references.

No migration or runtime feature should be enabled until the owner decisions, independent review, exact-head CI, disposable-database migration tests, and rollback rehearsal below are complete.

## 2. Current Schema and Authorization Facts

- `users` (`003_create_users.sql`) has a tenant-scoped integer ID, optional `restaurant_id`, email, password hash, one legacy role string, JSON permissions, status, and exact-case `UNIQUE(tenant_id,email)`. It is not a global person identity. `179_index_users_email_lower.sql` adds case-insensitive lookup, not global uniqueness.
- `employees` (`020_create_employees_table.sql`) is a separate HR record. It has tenant and restaurant scope and tenant-local email uniqueness; its creator/updater references point to legacy `users.id`.
- `user_roles` (`103_create_user_roles_table.sql`) maps tenant/user/role IDs, but still references tenant-scoped users. Roles in `021_create_roles_table.sql` are themselves tenant- and restaurant-scoped; `role_permissions` (`104`) maps module levels.
- RLS exists for selected slices, including auth audit logs (`164`), role permissions (`176`), and user roles (`177`). These policies use transaction-local `app.tenant_id` and `app.restaurant_id`. The full business schema is not uniformly RLS protected, so repository/query scoping remains required.
- Permission evaluation currently combines role assignments, role permissions, KDS/TABLES fallbacks, and legacy `users.role` owner/admin elevation. Migration must preserve the effective result or deny access; it must not broaden grants.
- Legacy user IDs are used by many foreign keys and audit actor fields. Preserve each row and ID through the compatibility period. Do not rewrite actor references to a new identity ID in Phase 1.
- Super Admin is a distinct account and session class; it must not become a tenant membership.

## 3. Threat Model and Immediate Finding

`AUTH-11` is a confirmed code-path flaw conditional on duplicate same-email rows with different password hashes: current login checks the first row, then signs a tenant-selection proof covering the other rows. A holder of the first row's password can therefore select a different tenant row without proving its password.

PR #105 narrows login to rows whose password hash matches and binds the signed selection proof to their exact legacy user/tenant pairs. Confirmation re-reads the exact row, and same-tenant duplicate matches fail closed. This is an interim guard, not canonical identity linking. Equal email, equal password, matching hashes, or a staging aggregate do not establish that two legacy rows belong to the same person. Production duplicate state remains unknown.

Phase 1 must also prevent:

- A pending, suspended, or revoked membership from appearing in tenant selection or receiving a session.
- A client-supplied tenant, restaurant, role, user, or membership ID from overriding verified server-side scope.
- A valid role in tenant A from authorizing an operation in tenant B or an ungranted restaurant in the same tenant.
- A stale session from retaining access after membership suspension, role revocation, or identity security change.
- A dual-read/backfill path from silently increasing the effective permissions of an existing account.

## 4. Proposed Additive Data Model

Keep existing `users`, `user_roles`, role/permission tables, and their foreign keys during the compatibility and rollback window. Add new tables without dropping or repurposing legacy IDs.

### `auth_identities`

- `id UUID PRIMARY KEY`.
- `normalized_email TEXT NOT NULL`, using the existing lower(trim(email)) behavior initially; do not add provider-specific alias rewriting.
- `email_verified_at TIMESTAMPTZ NULL`; legacy rows are unverified until a verification flow proves control.
- `password_hash TEXT NOT NULL`, copied per legacy row for one-to-one backfill; a later explicit linking flow may replace it with one chosen credential.
- `status` and monotonic `security_version`, plus timestamps.
- A unique index on `normalized_email` for verified identities only. Keep unverified duplicates representable so the migration can preserve conflicting legacy accounts without merging them.
- Add a `identity_email_challenges` table with identity/email binding, purpose, one-way token hash, expiry, attempt count, consumed timestamp, and creation metadata. Store no plaintext code or link token. Rate-limit issuance and verification, and consume atomically.

### `tenant_memberships`

- `id UUID PRIMARY KEY`, `identity_id UUID REFERENCES auth_identities`, `tenant_id` referencing the existing tenant, and `legacy_user_id` referencing `users(id)` with a unique constraint while compatibility is active.
- Explicit state (`pending`, `active`, `suspended`, `revoked`), `accepted_at`, and audit metadata. Only `active` memberships may enter login selection.
- A uniqueness rule for one identity/tenant after same-tenant duplicate conflicts are resolved. Do not merge multiple legacy rows in one tenant automatically.
- A composite key including `(id, tenant_id)` so child role/scope tables can enforce tenant consistency with composite foreign keys.

### Membership roles and restaurant grants

- Add membership-scoped role assignments keyed by membership, tenant, restaurant, and role. Reference roles within the same tenant and restaurant; record grantor and status.
- Add explicit restaurant grants for the membership. Treat access as an allowlist. Do not interpret a nullable legacy `users.restaurant_id` as “all current and future restaurants” until the owner confirms that semantic and the source paths are reviewed.
- A tenant-wide grant, if approved, must be an explicit scope value with defined behavior for future restaurants. It must not be inferred from `NULL` in a new identity row.

Exact FK/index DDL depends on verifying existing composite keys for `tenants`, `restaurants`, and `roles`. Add uniqueness support in additive migrations before child foreign keys. New filtered columns require indexes. Each migration must include a rollback script and be exercised forward and rollback in a disposable PostgreSQL database.

## 5. Login, Membership Linking, and Session Contract

1. Normalize the supplied email using the existing normalization rule and authenticate an identity. Do not enumerate memberships before credential verification. A minimal single-use email challenge verifies mailbox control in Phase 1; invitations, recovery, and MFA are later phases.
2. List only active memberships for that authenticated identity whose email is verified and whose tenant is active. Return a short-lived signed proof covering the exact membership IDs, not a client-editable tenant list.
3. If multiple memberships qualify, the user chooses one at login. Any restaurant choice must be constrained to explicit active grants for that membership. If a scope is unresolved, fail closed and return a safe generic error.
4. Re-read identity, membership, tenant, restaurant grants, and effective role state when consuming the proof. Reject expired/replayed proofs and any membership that changed since issuance.
5. Issue the current server-managed opaque session with claims bound to identity ID, membership ID, legacy user ID, tenant ID, restaurant scope, and security version. Keep the legacy user ID for existing ownership/audit lookups until all consumers have a reviewed migration path.
6. All tenant/restaurant/role values come from validated membership state. Headers and URL parameters may select only a resource that the current membership is already authorized to use.
7. Do not add in-session tenant switching in this phase. A new login may choose another eligible membership; the issued session remains bound to its selected membership.

### Safe linking of existing accounts

Backfill one identity per existing user row first. Do not merge rows because email strings or password hashes match. Linking multiple legacy rows to one verified identity requires:

- In Phase 1, prove mailbox control with a short-lived, one-time email challenge whose token is stored only as a hash. The challenge does not prove control of any legacy tenant account by itself.
- Proof of control of each legacy membership being linked (current credential or an explicitly accepted invitation/recovery path), unless the owner approves a documented exception with independent security review.
- Explicit user confirmation of the memberships and roles being linked.
- Transactional conflict checks for tenant, role, branch, suspension, and audit state.

If credentials diverge, never choose the first row’s password. Require an explicit credential reset or verified link procedure. If same-tenant rows or role/scope conflicts cannot be reconciled without privilege changes, preserve them separately and leave linking disabled for that identity.

The challenge provider must use a disposable fake in automated tests and an owner-approved sandbox/allowlist in staging. Do not send test messages to real users. Keep an identity and its memberships pending until mailbox verification, account proof, and explicit membership confirmation succeed.

## 6. Backfill and Conflict Handling

### Inventory on a disposable sanitized database

Before writing a migration, produce aggregate counts (no email values, hashes, names, or IDs) for:

- Duplicate `lower(trim(email))` groups across tenants and within a tenant.
- Groups with distinct password hashes, active/inactive mixtures, missing or invalid hashes, and case-only email collisions.
- Users with null restaurant scope, restaurant/tenant mismatches, inactive tenants, and conflicting role sources (`users.role`, `user_roles`, role permissions, JSON permissions).
- Employee records, audit rows, auth logs, and other foreign keys referencing each legacy user ID.

The report must be generated from disposable/sanitized data or a privacy-approved read-only aggregate query. No production/staging identity values or password hashes belong in reports.

### Synthetic conflict matrix

| Synthetic case | Backfill/link result |
|---|---|
| Same normalized email, two tenants, different passwords | Two separate unverified identities and memberships. No cross-tenant proof; link only after both accounts are proven. |
| Same normalized email and same password in two tenants | Still two identities initially. Password equality is not a merge signal. |
| Case-only duplicate email in one tenant | Conflict report; keep separate and inactive for canonical linking until reviewed. |
| User with `restaurant_id = NULL` | Preserve old row; do not infer branch grants. Resolve the meaning through code review and owner decision. |
| Legacy role grants more than role table | Freeze and document the currently effective grant set; fail if new calculation broadens it. |
| Pending/revoked invitation or membership | No login choice, session, or permission. |
| Employee record shares email but has no user account relationship | Do not attach it to an identity by email; HR matching is a separate reviewed operation. |

### Idempotent migration sequence

1. Add identity, membership, role-grant, and restaurant-grant tables plus RLS and least-privilege grants. Existing reads/writes continue unchanged.
2. Backfill one identity and one membership per legacy user row in bounded batches. Use `legacy_user_id` uniqueness to make reruns idempotent; preserve each `users.id` and all audit/FK references.
3. Copy current password hashes only to the one-to-one identity mapping; leave verification null. Copy roles/scopes only after an effective-permission comparison report; unresolved rows remain inactive in the new path.
4. Run dual-read comparisons behind a default-off feature flag on synthetic fixtures. Record only counts and mismatch classes, never credential values.
5. Enable reads for a small synthetic/canary cohort only after independent review, owner approval, and exact-head CI. Do not dual-write privilege grants until transactional consistency is proven.
6. Keep legacy columns/tables and old auth reads through the rollback horizon. Do not drop any legacy user key or role source in Phase 1.

## 7. Authorization, RBAC, and RLS Requirements

- Update `PermissionEvaluator` to resolve role assignments through the selected membership and current restaurant, while preserving deny-by-default behavior and the established permission levels.
- Remove legacy `users.role` elevation only after a parity report proves no user loses or gains permissions unexpectedly and the owner approves the compatibility window. Never allow role resolution to omit membership/tenant/restaurant predicates.
- Add RLS to each new membership/role/scope table using transaction-local validated identity, membership, tenant, and restaurant context. Missing context returns no rows; the migration role remains separate from the runtime role.
- Protect identity and email-challenge rows with a narrowly scoped auth lookup/verification boundary. The ordinary application role must not gain unrestricted cross-tenant reads to resolve global identities; challenge hashes and verified addresses must not be exposed through tenant APIs.
- Do not grant the ordinary application role unrestricted cross-tenant access just to resolve identities at login. Select and review a least-privilege auth lookup boundary (dedicated auth DB role or a narrowly scoped database function) before implementation.
- Keep parameterized queries, tenant and restaurant filters in repository methods, and consistent authorization at handler, use-case, repository, and RLS layers.
- Audit membership creation, acceptance, suspension, role grant/revoke, restaurant grant/revoke, and identity-link events without logging passwords, selection proofs, cookies, or verification tokens.

## 8. Rollback and Compatibility

- Feature flag defaults off; no membership table is authoritative until rollout approval.
- Before cutover, rollback is a flag change back to legacy reads. Additive tables may remain unused.
- After any identity is linked, rollback must map each membership back to its preserved `legacy_user_id`; never split or delete account history based on email. Rehearse this on a disposable database with multi-tenant synthetic identities.
- Keep existing API response compatibility for single-tenant login. For a multi-tenant result, only return server-verified memberships. Do not change the approved session cookie lifetime or switch Super Admin to staff sessions.
- No deployment occurs until migrations are backward-compatible with the previous binary, rollback is rehearsed, exact signed artifact is available, fresh backup is verified, and staging acceptance is run on the exact deployed SHA.

## 9. Phase 1 Verification Plan

All tests use synthetic identities and disposable databases. Never edit staging/production data to simulate a duplicate identity.

### Unit

- Identity/email normalization and verified/unverified state transitions.
- Link authorization: email proof alone cannot attach an unproven legacy membership.
- Login proof includes exact membership IDs; expired, replayed, tampered, pending, suspended, and revoked memberships fail.
- Role/restaurant decision functions are deny-by-default and do not trust request-supplied tenant IDs.
- Permission parity fixtures cover owner, manager, cashier, KDS, branch-limited user, legacy JSON fallback, and conflicting role sources.

### PostgreSQL integration and migration

- Fresh migration, rerun/idempotency, backfill batches, forward rollback, restore from snapshot.
- Two tenant memberships for one verified synthetic identity with distinct roles and restaurant grants.
- Distinct legacy credentials never produce a cross-membership selection proof.
- Existing user IDs, actor FKs, auth audit records, employee links, and role assignments remain intact.
- `pos_app` RLS tests deny tenant A access to tenant B membership and deny restaurant B access to restaurant A grants, including missing/mismatched transaction context.
- Concurrent identity linking/registration attempts produce one canonical verified email and no partial memberships.

### Browser and security acceptance

- Real Chromium login shows only active verified memberships after password verification; selecting one binds the session to that membership and restaurant scope.
- Synthetic browser flow completes the email challenge, explicitly links only proven memberships, and rejects expired, replayed, rate-limited, or email-only linking attempts.
- Tampered tenant, membership, restaurant, role, and user parameters fail in UI/API; forged headers do not change session scope.
- A role or membership revocation denies the next authorized request within the approved invalidation bound.
- Cross-tenant reads/writes fail through API, repository, RLS, and browser flows. Verify no response/cache contains another tenant's profile or role.
- Keyboard and screen-reader checks cover the tenant and restaurant selector; EN/AR and mobile RTL are included.
- Independent security review checks authentication, RLS, data migration, rollback, legacy handler reachability, and the CI target boundary. CI green is necessary but not a substitute for this review.

## 10. Decisions Required Before Schema Cutover

1. **Legacy account linking:** require the email owner to prove every existing account with its password, allow an approved support-assisted exception, or use a different explicit recovery process. Default proposal: prove each membership; no email-only merge.
2. **Null restaurant scope:** decide whether legacy `users.restaurant_id IS NULL` means all existing branches, a tenant-wide grant including future branches, or unresolved scope. Default proposal: unresolved until source semantics are verified; no implicit expansion.
3. **Email delivery and pending activation:** Phase 1 includes minimal email verification because verified identity is part of the requested scope. Default behavior keeps identity/membership pending until verification and account-link proof; do not send real messages until an owner-approved provider sandbox and recipient allowlist are configured. Broader OTP recovery remains a later phase.
4. **Email normalization:** confirm the current lower/trim rule and reserve verified addresses even when an identity is disabled. Do not add plus-address or provider-specific equivalence rules.
5. **Review and repository protection:** name an independent security reviewer and reconcile the live branch rule (zero required reviews) with `ops/aws/RUNBOOK.md` (one approval). Keep existing status checks and force-push/deletion restrictions.

## 11. Review Request and Current Gate

Please review the threat model, additive schema, per-membership role/restaurant representation, email-linking proof, legacy permission parity, RLS design, and rollback before an implementation migration is opened. The first implementation PR should contain additive schema and tests only, with all new identity reads default-off. A separate reviewed change can then wire authentication behind the flag.

PR #105 addresses the immediate divergent-password selection flaw, but remains open with CI running and no reviewer assigned. The currently configured GitHub branch rule requires zero approving reviews, and the only write-capable collaborator found is the repository owner/author. No independent review can be claimed until a reviewer is supplied and records approval.

Do not restore automatic staging deployment verification as part of this design. The CI isolation incident involved acceptance scripts reaching public staging endpoints; current main contains target guards, an internal-only CI network, mocked payments, and isolation tests. Historical synthetic chart records have unresolved cleanup status. Keep the signed owner-dispatched publisher and exact-SHA verification intact until an independent reviewer signs off on that boundary and the incident disposition is documented.
