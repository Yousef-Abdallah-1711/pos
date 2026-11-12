# Module Readiness Matrix
Status from static audit + limited runtime; "Implemented" means code exists and wiring was read, not that it was exercised.

| Module | Frontend | Backend | Key blockers |
|---|---|---|---|
| Registration/Auth | Implemented | Implemented | Clusters A, C; trial dates never set |
| Products/Categories | Implemented | Implemented | Categories lack permission wrapper |
| Inventory | Partial (adjust is a stub) | Partial | No adjust route; no decrement on order |
| Orders | Implemented | Partial | Non-atomic, number collisions, public IDOR |
| POS | Partial | Uses public order route | Offline duplicates; SW sync wrong origin |
| KDS | Implemented (code) | Implemented | WS auth unverified |
| Tables + QR | Implemented | Implemented | Not exercised |
| Drivers | Partial | Implemented | Auth-only, 18/23 repo methods unscoped |
| HR | Partial (advanced) | Implemented | Wrong module ID on routes; payroll→accounting link absent |
| Accounting | Partial (advanced) | Implemented | FIN-001/003/005 |
| Financial reports | Implemented | Implemented | Client-side calculations unverified |
| RBAC | Implemented | Not enforced | Cluster D |
| Theme builder | Partial (rollback stub) | Implemented | IDOR (F) |
| Billing/Paymob | Implemented | Mock by default | Cluster G |
| Super admin | Implemented | Implemented | Client-side auth only; rate limiter flaw (FE-OTH-006) |
| Storefront | Implemented | Partial | Missing `/public/orders`, `/public/themes*`; locale on maintenance |
