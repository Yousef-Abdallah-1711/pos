# Test coverage matrix

Legend: `Some` means tests exist but coverage is not complete; `Not verified` means no current evidence adequate for the layer. No percentages are inferred.

| Module | Unit | Integration | Security | E2E | Main gaps |
|---|---|---|---|---|---|
| Authentication | Some | Some | Some | Not verified | reset, refresh/logout lifecycle, browser cross-app session/token handling, lockout/distributed limits |
| RBAC | Some | Selected routes | Selected negatives | Not verified | complete role×permission matrix, resource-level constraints |
| Tenants | Some | Selected smoke | Selected negatives | Not verified | full cross-tenant CRUD across repositories |
| Restaurants | Some | Partial | Not verified | Not verified | two restaurants per tenant and cross-restaurant isolation |
| Products | Some | Limited | Not verified | Not verified | foreign IDs and inventory side effects |
| Inventory | Some | Limited | Not verified | Not verified | atomic stock changes and scope tests |
| Orders | Some | Selected | Public-token cases | Not verified | POS and public create end-to-end, transaction integrity |
| POS | Limited | Not verified | Not verified | Not verified | full cashier workflow |
| KDS | Some | Partial | Limited | Not verified | WebSocket handshake/tenant auth and order lifecycle |
| Payments | Some | Limited | Some webhook negatives | Not verified | success transaction, concurrency, out-of-order, retries |
| Accounting | Some | Limited | Not verified | Not verified | multi-table transaction and tenant isolation |
| HR | Some | Limited | Not verified | Not verified | employee/restaurant isolation and end-to-end workflows |
| Storefront | Some | Not verified | Public token tests | Not verified | all browser flows and cross-tenant resource disclosure |
| Themes | Some | Some smoke | Some isolation | Not verified | all repository paths, migration between theme models |
| Super Admin | Some | Not verified | Selected middleware tests | Not verified | end-to-end impersonation/audit/expiry/revocation |

Basis: Go suite reports 680 passing tests; selected disposable security smoke 73/73; frontend test outcomes detailed in `09_TEST_INFRASTRUCTURE_AUDIT.md`. These counts are not coverage percentages or proof of production correctness.
