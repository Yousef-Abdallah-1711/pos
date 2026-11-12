# Test Coverage Matrix (static)

| Module | Backend tests | Frontend tests | E2E | Gap |
|---|---|---|---|---|
| Auth/JWT/RBAC | jwt_test, rbac_test (live DB), rbac_integration (build-tagged) | landing login test | none | no login/401 E2E |
| Multi-tenancy/RLS | phase6_tenant_isolation (13), rls tests, tenant_isolation_test | - | none | live verification depends on DB |
| Orders/POS | order_usecase_table_integration | posApi.test, pos component | none | order_repo not directly tested |
| KDS | repo/usecase/handler + kds_flow_e2e | 3 component tests | none | relatively good |
| Tables | repo/usecase/handler | 6 component tests | none | good |
| Accounting | edge cases, auto-posting, phase6 cycle | none | none | accounting_advanced_repo (2133 lines) / handler (1346) no direct tests |
| HR | employee_handler (tagged), salary auto-posting | none | none | hr_report_handler (938 lines) untested |
| Reports | financial_report cache/usecase | reportsApi, reports tests | none | report_handler (2630 lines) |
| Payments/Paymob | payment_handler (tagged), Newman json | none | none | webhook paths unverified |
| Subscription/Trial | none found | none | none | tenant_status middleware untested |
| Themes | admin_theme_handler* (tagged) | 6 theme-builder tests | none | |
| Super-admin | none found | 1 test | none | impersonation untested |
| Restaurant website | - | 8 jest | none | |
| Scaling/Sharding | 7 unit | - | - | code not wired to server |
