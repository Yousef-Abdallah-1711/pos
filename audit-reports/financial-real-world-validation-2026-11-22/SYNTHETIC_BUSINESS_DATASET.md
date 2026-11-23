# Synthetic business dataset

## Scope and provenance

All exercised records were generated in disposable local infrastructure with synthetic identities and values. The authenticated acceptance scenarios used isolated tenants/restaurants and deterministic mock payment behavior. The separate performance run used database `servio_perf_test` and a seeded workload; the API acceptance database `servio_test` also contained unrelated repository-test fixtures and must not be treated as one clean restaurant dataset.

## Counts and financial scale

The broad `servio_test` inventory observed during the run included 316 tenants, 519 restaurants, 100 users, 144 products, 2,078 orders, 10,159 journal entries, and 50,332 general-ledger rows. These totals include many test fixtures and do not represent the six API acceptance scenarios alone. Exact source IDs and cents-level expected/actual values for those scenarios are in `financial-live-api-tests.jsonl`.

The benchmark dataset comprised 2 tenants, 1,102 journal entries (tenant counts 1,077 and 25), 5,229 journal lines, 5,229 GL rows, 1,000 synthetic sales orders, and 2,000 sales-order items. Seeded journal and GL debit/credit totals were each 11,032,700 cents in the structured benchmark report. A separate seed reconciliation log recorded 11,000,000 cents for its specific invariant check; these are distinct reported checks and should not be combined as if they were one total.

## Limits

This run does not establish that the minimum requested three-tenant, multi-restaurant role matrix or every product, ingredient, shift, bank, inventory, and tax configuration was constructed. Unsupported or unexercised scenarios remain untested, not implicitly accepted. No customer PII or actual payments were used.
