# F-003 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Open Cash Flow after S13

## Expected result
Net cash change equals change in cash accounts

## Actual result (measured)
Operating = net profit (-2000); net cash flow 98000 vs GL cash+bank 93000 | RUN-003 (AUD-MET-FIN): cash-flow net 82,710.01 vs GL cash+bank 115,010.01, residual -32,300.00 (STILL REPRODUCIBLE).  | RUN-003 earlier run: cash-flow 82,710.01 vs GL cash+bank 115,010.01 (residual 32,300.00). Re-measured after more activity: report 172,710.01 vs GL 1110+1111 = 280,010.01, residual 107,300.00. ROOT CAUSE CONFIRMED (AUD-CF-TRACE): usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment that never reads cash/bank accounts and omits every working-capital movement (AR -6,520; inventory +96,000; input VAT -2,800; card clearing -11,000; AP +800; VAT payable +30,820 = +107,300); investing = '12%' debit-credit has inverted sign for asset purchases; financing excludes 21xx/2120 liabilities. 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
