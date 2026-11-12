# Redis and Tenant Isolation Results

- Cross-tenant GET verification passed for the two synthetic tenants: each returned 21 products, with no product-ID overlap or opposite-tenant fixture item. Tenant-status responses differed between tenant contexts.
- A mixed 20-request GET probe increased Redis global counters by 19 hits and 1 miss. Counters are global and cannot be assigned to an endpoint or tenant.
- Observed key prefix: tenant_access; source TTL: 30 seconds. Endpoint-specific hit ratio, database reduction, cold/warm latency delta, memory footprint, and invalidation were not measured.
- Tested product and tenant-status reads showed no cross-tenant cache leakage. This is narrow synthetic evidence, not a full cache-isolation audit.
- No mutations, cache flushes, or cache configuration changes were made.
