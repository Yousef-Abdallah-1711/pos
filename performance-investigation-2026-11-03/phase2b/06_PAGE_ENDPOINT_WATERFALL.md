# Page-to-Endpoint Waterfall Summary

## Authenticated Dashboard home

Observed 23 successful API requests: 8 permission-check calls; 2 orders; 2 notifications; 2 notification-statistics; 2 tenant-status; 7 calls not attributed in the retained aggregate. Total resources: 73. These counts identify repeated request families; initiators and critical-path ordering were not retained.

## Storefront menu

/en/menu rendered 21 synthetic products and issued two API requests after the isolated CORS allowlist was corrected. Before correction, storefront requests were blocked by the stale allowlist. /en/cart and /en/checkout showed empty-cart states; no business mutation was used to populate a journey.

## Super Admin Audit Logs

/dashboard/audit-logs issued GET /api/v1/super-admin/audit-logs?page=1&page_size=50, returned HTTP 200 and 9 rows. Two rows contained null details; this caused the frontend exception described in 07_SUPER_ADMIN_AUDIT_LOGS_FAILURE.md.

Other page-to-endpoint waterfalls were not retained at exact endpoint-level granularity. The pre-existing dependency matrix remains static source evidence, not runtime evidence.
