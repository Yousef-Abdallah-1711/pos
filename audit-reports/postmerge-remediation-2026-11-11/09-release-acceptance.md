# Release acceptance

| Gate | Result | Evidence / blocker |
|---|---|---|
| Four Next.js production builds | PASS | Independent host builds with required production URL variables. |
| Production Docker image builds | PASS (build only) | Go API and all four Next.js images built at candidate `ec0959d`; runtime stack was not started because placeholder `.invalid` URLs are build-only and live payment/database secrets were not used. |
| API without manual port-forward | PASS (isolated) | Disposable API health and frontend/API communication passed on the published ports in the seven-service test stack. Production deployment wiring is not inferred. |
| Authenticated E2E | PASS (bounded) | 16/16 Phase 2A tests and 12/12 registration Playwright tests passed; Phase 2A now includes live KDS WebSocket delivery. See report 06 for covered paths and exclusions. |
| Dashboard responsive layout / CLS | PASS (bounded) | No overflow at 320/390/768/1440; latest authenticated dev-stack CLS was 0.031. Production Lighthouse remains unrun. |
| Tenant isolation negative matrix | PARTIAL | 27 repository-level PostgreSQL tests plus restricted-role RLS and account/settings/payment coverage pass. Full table and HTTP role matrix remains incomplete; RLS is enabled/forced on 15 of 92 tables. |
| POS/KDS end-to-end | PASS (bounded) | POS order creation was observed over the authenticated KDS tenant WebSocket; inventory flows also passed. Broader KDS lifecycle acceptance remains open. |
| No unresolved P0 security defects | UNKNOWN | No confirmed exploit in prior bounded checks; full RLS and tenant matrix not run. |
| No unresolved P1 release blockers | NO | API/E2E/security/performance acceptance remains incomplete. |
| Hosted smoke checks | PASS | `runtime-localhost-guardrail` and `url-domain-config-smoke` passed on PR #14. These are not the full required regression suite. |
| Independent review | PENDING | No review approval is recorded on PR #14. |
| Full integrated CI | INCOMPLETE | Go tests/vet/build, four production builds (restaurant build rerun after its latest change), frontend suites, fresh migration, targeted rollback, payment, retry, repository DB suites, and bounded Playwright pass. Full RLS/route matrix, production Docker images, and performance gates remain unrun. |

**Decision: NO-GO for production.** This remediation is limited to verified application defects and build/runtime configuration. Public branding remains POOS; ORDERLY is not included because commercial and trademark clearance is unverified.
