# F-006 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
5 parallel receipts of 1000 on balance 3200

## Expected result
At most 3 accepted; no over-collection

## Actual result (measured)
5 accepted; amount_paid 7000 on 5200 invoice; AR GL -1800 | RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment.  | RUN-003 FRESH-INVOICE RE-REPRODUCTION (AUD-CINV-FLOW): invoice 3,000, 5 concurrent receipts of 1,000 -> all 5 accepted (201), amount_paid 6,000, AR credit -3,000; single sequential overpayment 5,000 is rejected (guard exists but is not serialized). RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment. 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
