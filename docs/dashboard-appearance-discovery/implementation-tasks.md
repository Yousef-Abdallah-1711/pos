# Future implementation tasks

All tasks are NOT STARTED. Complete source and runtime discovery before estimating implementation.

| ID | Objective and existing files | Depends on | Acceptance and tests | Risk / complexity |
| --- | --- | --- | --- | --- |
| A01 | Finalize token contract from `globals.css`, `tailwind.config.js`, `DESIGN.md` | Discovery | Document defaults, accessibility and consumers | Medium / high |
| A02 | Design migration and compatibility for personal colors and storefront sync in appearance page | A01 | Existing tenant and storefront data preserved | High / high |
| A03 | Add validated tenant persistence, RLS and audit using `backend/migrations`, `tenant_tx.go` | A01–A02 | Real DB isolation, rollback and concurrency tests | High / high |
| A04 | Add authenticated and authorized API using `backend/cmd/api/main.go` and middleware | A03 | RBAC, invalid input and cross-tenant tests | High / high |
| A05 | Implement safe resolver and scoped preview in dashboard shell and preferences store | A04 | Refresh, mode, tenant switching and fallback tests | High / high |
| A06 | Create workspace settings UI with locale-aware controls and accessibility feedback | A05 | Keyboard, RTL, save/reset, errors and unsaved-state tests | High / high |
| A07 | Migrate sidebar, shared controls, routes, POS/KDS and portals incrementally | A05 | Visual and business workflow regression | High / very high |
| A08 | Recapture synthetic product matrix and integrate marketing assets | A07 | Provenance, all variants, geometry and link tests | Medium / high |
| A09 | Run browser, security, accessibility, performance and full functional gates | A03–A08 | Evidence-based release verdict | High / high |
