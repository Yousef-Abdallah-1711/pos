# Registration and login validation contract

The shared form package at `frontend/packages/form-foundation` is the frontend source for registration/login schemas, form semantics, safe API error parsing, and debounced async availability. New forms can reuse `FormField`, `FormErrorSummary`, `PasswordField`, `SubmitButton`, `fieldAccessibility`, `parseApiError`, and `applyApiFieldErrors`; add `@pos-saas/form-foundation: workspace:*` and include it in the app's Next.js `transpilePackages` list.

The backend is authoritative and validates the decoded DTO before any registration transaction. PostgreSQL remains the final authority for unique tenant slug/domain and per-tenant user email constraints. Migration `179_index_users_email_lower.sql` adds a non-unique functional index for case-normalized identity lookups; rollback is `DROP INDEX IF EXISTS idx_users_email_lower`. No global user-email constraint is added because existing login explicitly supports the same email in multiple tenants. Tenant display name remains globally unique in the current schema; this change preserves that behavior pending a product decision.

| Field | Normalization and limits | Persistence rule |
| --- | --- | --- |
| Owner name | Trim; 2–255 Unicode code points | `users.name VARCHAR(255)` |
| Email | Trim and lowercase; valid address; at most 255 Unicode code points | Tenant scoped unique constraint remains authoritative |
| Restaurant name | Trim; 2–255 Unicode code points | `tenants.name VARCHAR(255)` remains unique |
| Website subdomain | Trim/lowercase; 2–63 ASCII lowercase letters, digits, and single hyphens; alphanumeric ends; reserved names rejected | Unique slug/domain constraint; availability is advisory |
| Phone | Optional; blank becomes absent; accept explicit international `+` input, validate with libphonenumber metadata, store E.164 | At most 15 digits plus `+`, within existing `VARCHAR(20)` |
| Password | At least 6 Unicode code points; at most 72 UTF-8 bytes | Existing bcrypt hashing; confirmation is frontend-only |

Login trims and lowercases email, validates address syntax, requires a password, and enforces bcrypt's 72-byte input limit. It does not apply registration's six-character minimum to existing accounts. Credential mismatches retain one generic 401 response.

The availability hook exposes idle, invalid, checking, available, taken, reserved, and network-error states. A response only applies to its still-current request; registration also checks that the verified slug equals the submitted canonical slug. The registration endpoint checks uniqueness again inside its transaction.

API auth errors keep the legacy `error` property and add stable `code`, `message`, optional `fields`, and request correlation ID. The client selects localized safe messages by code/status and does not render unknown server strings.
