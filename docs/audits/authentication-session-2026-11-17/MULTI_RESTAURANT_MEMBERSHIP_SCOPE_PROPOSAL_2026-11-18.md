# Same-Tenant Multi-Restaurant Membership Scope Proposal

**Status: proposed, pending product-owner decision. No schema change is authorized.**

## Current model and verified scope

`backend/migrations/003_create_users.sql` stores `restaurant_id` on `users` and enforces `UNIQUE (tenant_id, email)`. A single tenant user row can therefore point to one restaurant. The existing login-time choice is for a user with memberships in multiple tenants; it does not select among multiple restaurants inside one tenant. The real-service browser scenario covers tenant selection and session scope, but it does not prove same-tenant multi-restaurant membership.

The repository's implementation roadmap says in-session tenant switching is optional until the identity and membership model is decided. The subscription UI also lists “Multi-Restaurant Management” for the Pro plan. The repository does not say whether that commercial feature must be available in the first authentication release. That conflict needs a product decision; code and tests cannot determine the release promise.

## Recommended release-1 scope

Defer same-tenant multi-restaurant membership from this authentication release. Keep one restaurant association per `(tenant_id, email)` membership, preserve the tested login-time multi-tenant choice, and keep in-session switching disabled. This recommendation does not cancel or alter a plan entitlement. Product must confirm that multi-restaurant access is not a release-1 acceptance requirement or reconcile the entitlement before release approval.

Until that decision is recorded, scenario 16 stays incomplete and the release decision remains blocked. Do not treat absence of a schema migration as evidence that the feature is unnecessary.

## Secure behavior supported now

- Login credentials are verified before multi-tenant choices are returned.
- A selection proof is required to confirm a tenant; a session is not issued while selection is pending.
- The selected tenant and default restaurant come from server-authorized records, not client-supplied IDs.
- The accepted tenant is recorded in the session/JWT context and scope checks reject forged tenant or restaurant headers.
- Logout and re-login use server-side session revocation/rotation behavior covered by the authentication acceptance suite.

The successful browser check covers multiple tenants, not multiple restaurants in one tenant.

## If Product requires the feature

Do not implement it as an extension of the current `restaurant_id` field. Prepare an additive identity and membership migration with the following design constraints:

1. **Identity:** Introduce a tenant-scoped identity with one normalized email and password credential per tenant, retaining a stable identity ID.
2. **Membership:** Add a `user_restaurant_memberships` relation with identity, tenant, restaurant, role, status, permission-set reference, timestamps, and uniqueness constraints. Add composite foreign keys so restaurant and membership tenant IDs must agree.
3. **Roles:** Store authorization role and permissions on each membership; define whether “owner” is tenant-wide, restaurant-wide, or both before migration.
4. **Authorization:** Derive the active restaurant only from an active membership loaded by the server. Recheck membership status on protected requests; ignore client tenant/restaurant selectors without membership validation.
5. **Session rotation:** Selecting or changing membership rotates the opaque session ID and updates the server-side session atomically. Revoke the prior session, rotate CSRF state as applicable, and write an audit event.
6. **Audit:** Record actor, previous and selected membership IDs, tenant, request ID, timestamp, and outcome. Do not log cookies, tokens, or proofs.
7. **Cache invalidation:** Clear client query caches and tenant/restaurant scoped API caches on membership change; all shared-cache keys must remain independent of authenticated identity and private routes must bypass shared caching.
8. **Backward compatibility:** Add and backfill membership rows from existing users first, keep old reads during a verification window, then move reads/writes behind feature-gated code. Do not remove `users.restaurant_id` until consistency checks and a rollback window complete.
9. **Rollback:** Keep migration additive and reversible before any membership selection is written. Provide forward/backfill validation, a rollback script, and a restore path from the prior single-restaurant state. Do not delete memberships on rollback.
10. **Acceptance:** Test authorized sibling restaurant selection, unauthorized sibling and foreign tenant denial, role differences, revocation, concurrent switch/refresh, session rotation, audit records, and cache invalidation with isolated synthetic fixtures.

## Owner action

Product owner: confirm one of these outcomes in the release decision record:

- `DEFERRED FROM RELEASE 1`: confirm the current single-restaurant membership model is acceptable for the first release, and reconcile the Pro plan promise; or
- `REQUIRED FOR RELEASE 1`: approve the additive design and migration proposal, then schedule implementation and re-run scenario 16.

No code or migration change is authorized by this proposal.
