# F-043 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): backend/internal/repository/accounting_advanced_repo.go:838-880 revenue=SUM(orders.total_amount paid, created_at), expenses=GL expense lines, cash_position=SUM(bank_accounts_ledger.current_balance)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/accounting_advanced_repo.go`:838-880.
- Propagation: affected modules Financial KPI dashboard; related findings F-029; observed consequence: KPI revenue 319,660 (paid orders incl tax) vs P&L revenue 305,300; net_profit 181,540 vs P&L 167,180 (+14,360); cash_position 50,000 (bank_accounts_ledger) vs GL cash+bank 276,270.01 | RUN-003: KPI net profit 181,540 vs P&L 167,180; cash_po
- Cause versus symptom: the symptom is what users or reports show (KPI revenue 319,660 (paid orders incl tax) vs P&L revenue 305,300; net_profit 181,540 vs P&L 167,180 (+14,360); cash_position 50,000 (bank_a); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: tests/runtime/ph12_batch_apis.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/accounting_advanced_repo.go | 838-880 | backend/internal/repository/accounting_advanced_repo.go | exact | yes |

