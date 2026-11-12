# Admin Orders N+1 Analysis

## Finding

The source contains a clear N+1 query pattern in `GET /api/v1/admin/orders`. `AdminOrderHandler.ListOrders` first calls `OrderUseCase.ListOrders`, which reaches `OrderRepository.ListOrders`; that repository performs one count query and one paginated list query. The handler then calls `GetOrder` once for every returned order. `GetOrderByID` performs an order/table query and then calls `GetOrderItems`, which performs another query.

Ignoring driver-context work and errors, the source-derived count is therefore approximately `2 + 2N` SQL statements for a page of N orders. This is a source estimate, **not a runtime query-count measurement**. The route also has a second `GetOrder`-per-row pattern in the order export handler.

## Exact proposed implementation batch

Once Phase 2A ownership of these files ends, add a list-specific repository path for the admin list:

1. Keep the existing count query, filters, pagination, and descending `created_at` ordering.
2. Have the paginated order query return all fields required by `adminOrderSummary` from the order row and tenant/restaurant-scoped table join, rather than reloading each order.
3. Add a batch item query for the page IDs, joining `order_items` to `orders` and constraining both `orders.tenant_id = $tenant` and `orders.restaurant_id = $restaurant`; use the page IDs as an additional predicate and order by `order_items.created_at`. Group items by `order_id` and calculate `ItemsCount` from those grouped results.
4. Preserve the repository result order by mapping the batch results by ID and projecting in the original `response.Orders` order. Do not change response keys, pagination, authorization, filters, or error behavior without a separate contract decision.
5. Keep the single-order detail path unchanged. Avoid broad changes to shared `ListOrders` behavior until callers are audited.

With a detail-complete page query and one batch item query, the expected list path is count + page query + item batch: three SQL statements independent of N. Confirm the actual count after implementation, including driver context and any middleware/database instrumentation.

## Required verification

- Inspect current Phase 2A diff before editing and coordinate exact file ownership.
- Add repository/handler regression tests for response parity, stable order, pagination, empty pages, and item counts.
- Add cross-tenant and cross-restaurant negative tests proving the batch item lookup cannot return foreign rows.
- Benchmark pages 1, 10, 25, and 50 with the same synthetic fixtures and build conditions before and after.
- Record SQL count, per-order query count, SQL duration, handler duration, serialization duration, end-to-end latency, response bytes, and repeated-sample distributions.
- Use `EXPLAIN (ANALYZE, BUFFERS)` only on the disposable database and verified read-only queries.

No before/after counts or latency data is available in Phase 2C. Phase 2B order latency figures are historical and must not be used as the baseline for this source snapshot.
