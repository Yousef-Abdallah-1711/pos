# Phase 2A task reconciliation

The original Phase 2A plan contains 32 C-items (plus session 1 tasks), while this audit checklist contains 30 AUD-items. Status below uses current evidence and reopens claims where required.

| Original tasks | Implementation | Verification | Remaining work | Status |
|---|---|---|---|---|
| C001–C002 baseline/reconciliation | Present | Current branch/HEAD/status observed | Historical baseline and old findings not all replayed | PARTIAL |
| C003–C004 route RBAC/enforcement | Middleware and route wrappers | 274 routes + 73 disposable smoke; current local override enforce | Full route×role×resource coverage missing; generic non-prod default audit | PARTIAL |
| C005 public order authorization | Token/hash design and migration 161 | Selected unit + disposable smoke | Browser and real public order create/POS path | PARTIAL |
| C006 resource isolation | Some theme/resource checks | Selected smoke | Complete cross-tenant and cross-restaurant tests | PARTIAL |
| C007 unsafe auth defaults | Startup guards, rate limiting, metrics/internal token controls | Selected Go/security tests | Password reset, account lockout and distributed rate limiter review | PARTIAL |
| C008 webhook security | HMAC/state checks and transactional callback processing | Isolated worker test covers two concurrent claims, exactly-once delivery per fixture, persisted success and attempts | Backoff/recovery, provider/local-commit reconciliation, and broader provider acceptance | PARTIAL |
| C009 impersonation integration | Middleware/usecase components | Middleware tests only | Full login→impersonate→expiry/revoke/audit flow | NOT STARTED |
| C010 restricted DB role | SQL setup + isolated compose | Disposable role/security smoke | Dev and production rollout/rehearsal | PARTIAL |
| C011 tenancy classification | 91-table generated inventory plus explicit review of `webhook_retry_queue` | Schema/FK review and worker query review | Per-family RLS policy selection remains in C012 | COMPLETE (classification only) |
| C012 staged RLS | transaction helper and pilot | Unit + isolated pilot | Move all module paths to scoped transactions; no live policies | PARTIAL |
| C013 Go failure repair | Changes in Phase 2A code | 680 current-checkpoint tests passed | Maintain repeatable CI | VERIFIED COMPLETE (baseline) |
| C014 Go build/static check | Build and vet | Passed at checkpoint | CI gate absent | PARTIAL |
| C015–C018 frontend config/type/build | App typecheck and builds reported green | Current landing/storefront tests fail; dashboard red; super-admin 3 pass | Repair test config/failures, test type errors and rerun | PARTIAL |
| C019 fixture completeness | Basic tenant/restaurant/user/order smoke fixtures | Disposable test stack | More restaurants/products/inventory/roles/super-admin/storefront data | PARTIAL |
| C020 four-app test stack | Dev stack has apps, isolated test stack does not | Runtime HTTP 200 only | Isolated Redis + all four apps | NOT STARTED |
| C021 integration smoke | API/Postgres smoke | 73/73 selected checks | Full app and browser integration | PARTIAL |
| C022 architecture definition | Existing backend architecture report | Source inventory | Module boundaries not fully implemented | PARTIAL |
| C023–C024 backend refactor | No broad extraction | Not applicable | Route/composition and module refactor | NOT STARTED |
| C025 frontend architecture | Existing report | Static inventory | Type/package/API boundaries unresolved | PARTIAL |
| C026–C027 frontend refactor | Some type cleanup only | App typecheck prior pass | Shared package and feature boundaries | NOT STARTED |
| C028 error foundations | Existing logging/errors | Partial source review | Typed errors, common mapper, request IDs and redaction | NOT STARTED |
| C029 backend regression | 680 default tests | Pass at checkpoint | Wire into CI | PARTIAL |
| C030 frontend regression | Partial | 109 dashboard fail; fresh landing/storefront fail | Stabilize all app suites | FAILED (gate) |
| C031 cross-app smoke | No full journey | Not run | login → dashboard → order → storefront/admin journeys | NOT STARTED |
| C032 compatibility verification | Route set comparison and payload notes exist | Prior script evidence | Re-run on future changes; storefront contract test | PARTIAL |

No task is elevated to complete solely from a historical report. Source artifacts and detailed verification limits are linked in the individual audit reports.
