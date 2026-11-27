# F-061 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Sell a product with price 0 (product creation allows it), mark the order paid, read accounting_event_outbox

## Expected result
Zero-value sale completes without a journal (or posts a zero-line memo) and the event is completed

## Actual result (measured)
order 73 total 0.00 paid: outbox status retryable attempts 1/8 last_error "at least two journal lines are required"; it will exhaust retries and appear as a failed event needing manual redrive; no journal exists | RUN-003: zero-total paid order 73 -> outbox retryable, error at least two journal lines are required (CONFIRMED auto_posting_engine.go:535). 

## Re-runnable scripts / evidence
- `tests/runtime/ph16_masterdata.mjs`
- `tests/runtime/ph16_negprice.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
