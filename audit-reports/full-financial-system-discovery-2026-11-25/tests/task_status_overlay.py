"""Applies RUN-003 execution results to 51-phased-audit-task-register.csv (run after build_task_register.py). Idempotent audit tooling."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
U = {
 'T-P02-003': ('EXECUTED-PASS', 'pg_dump start and mid dumps with sha256 in evidence/database/RUN-003/T-P02-003-checkpoint.txt; final dump added at closure', ''),
 'T-P02-004': ('EXECUTED-PASS', 'outbox retry/lease/redrive verified live (CFR-004, REL-004, REL-027)', ''),
 'T-P03-001': ('EXECUTED-PASS', 'F-001 root cause CONFIRMED (auth_repo.go:42-212); 52-complete-money-operation-trace.md written (16 operations, file:line)', 'F-001'),
 'T-P03-002': ('PLANNED', 'DISC-001 (119 vs 117) explained in 45, route discovery not re-run', ''),
 'T-P03-003': ('EXECUTED-PARTIAL', 'hard-coded values found: posStore TAX_RATE 0.14, POS currency symbol, CSP; duplicate revenue definitions (F-012); not a full dead-code sweep', 'F-004,F-012,F-023'),
 'T-P04-001': ('PLANNED', 'control-level inventory (file 53) not created', ''),
 'T-P04-002': ('PLANNED', 'field matrices not executed per field; validation covered by API batteries (Phases 16-17)', ''),
 'T-P04-003': ('EXECUTED-PARTIAL', 'route matrix in headed Chrome (EN, AR, mobile) in evidence/network/RUN-003/PH-04/AUD-UI-MATRIX.json', 'F-027,F-032'),
 'T-P05-003': ('PLANNED', 'staff/role creation via UI not executed', ''),
 'T-P05-004': ('EXECUTED-PARTIAL', 'login via real UI, origin check, rate limit, cookie tamper executed; logout/password recovery not executed', ''),
 'T-P06-001': ('EXECUTED-PARTIAL', 'controlled fresh-product stock/ledger experiment PASS; F-040; API-diagnostic (not UI forms)', 'F-040'),
 'T-P06-002': ('EXECUTED-PARTIAL', 'supplier create/invoice/approve/partial/full/overpay/concurrent/void executed via API; F-041', 'F-041'),
 'T-P06-003': ('EXECUTED-PARTIAL', 'AP aging ties to GL 2101; supplier current_balance competing source F-021 not re-measured', 'F-021'),
 'T-P07-001': ('EXECUTED-PARTIAL', 'cash POS order in headed Chrome with legible totals (AUD-POS-001/002); dine-in/takeaway/delivery variants not separately executed', 'F-004,F-023'),
 'T-P07-003': ('EXECUTED-PARTIAL', 'API duplicate/idempotency tests REL-001/REL-025 (F-017/F-052); UI double-click not executed', 'F-017,F-052'),
 'T-P07-004': ('BLOCKED-CAPABILITY', 'no cashier shift/cash-drawer module found in the product (33 capability matrix)', ''),
 'T-P08-001': ('EXECUTED-PARTIAL', 'storefront menu and online order in headed Chrome after allowing wildcard CORS origin; default isolated CORS reproduces F-022 UI symptom (environment configuration)', 'F-022'),
 'T-P08-002': ('PLANNED', 'theme builder controls not exercised', ''),
 'T-P08-003': ('EXECUTED-PARTIAL', 'KDS board in headed Chrome; cancelled orders remain (F-044); website order enters KDS on confirm', 'F-044'),
 'T-P08-004': ('PLANNED', 'drivers/delivery not exercised (driver fixture exists)', ''),
 'T-P09-002': ('EXECUTED-FAIL', 'F-006 re-reproduced on a fresh invoice (5/5 concurrent receipts accepted)', 'F-006,F-052'),
 'T-P10-002': ('EXECUTED-PARTIAL', 'attendance, leaves executed via API (F-053, F-054); contracts/insurance/performance/training/resignations not exercised', 'F-053,F-054'),
 'T-P10-003': ('EXECUTED-PARTIAL', 'salary arithmetic and mark-paid idempotency PASS; paid-without-approval F-015; accrual/reversal not available', 'F-015'),
 'T-P11-001': ('EXECUTED-PARTIAL', 'journal validation/post/void/reversal executed via API; attachments/approval not exercised (F-039 earlier)', 'F-013'),
 'T-P11-002': ('EXECUTED-FAIL', 'bank opening balance posts nothing (F-029 re-reproduced on a fresh account); KPI cash uses bank ledger (F-043)', 'F-029,F-043'),
 'T-P11-003': ('EXECUTED-PARTIAL', 'budgets, FX, recurring journals executed (F-047, F-048); fixed assets/bank reconciliation/import-export not executed', 'F-047,F-048'),
 'T-P12-002': ('EXECUTED-PARTIAL', '113 of 119 coverage rows measured (API + SQL oracle; Chrome rendering for statements); exports/print and multi-dataset variants not executed', 'F-002,F-003,F-009,F-010,F-042,F-043,F-045,F-046'),
 'T-P12-003': ('EXECUTED-PARTIAL', 'HR/KDS/dashboard counts measured; several vacuous because the dataset has no attendance/overtime/training', ''),
 'T-P13-001': ('BLOCKED-CAPABILITY', 'paid checkout disabled in the environment (503); trial/plans/disable/enable/extend executed', 'F-055'),
 'T-P13-002': ('EXECUTED-PARTIAL', 'Super Admin stats/tenant list match SQL; payment metrics vacuous (no payments)', ''),
 'T-P14-001': ('EXECUTED-PARTIAL', 'compressed operating month on fresh tenant F reconciles to independent ledger; payroll blocked (F-014)', 'F-002,F-003'),
 'T-P15-001': ('EXECUTED-PARTIAL', 'residual decompositions AR 1,800, cash flow, revenue bridge, VAT, KPI recorded in evidence PH-09/PH-12', 'F-003,F-006,F-010'),
 'T-P15-002': ('EXECUTED-PASS', '56-accounting-policy-decision-package.md: 15 decisions, none approved', ''),
 'T-P16-001': ('EXECUTED-PARTIAL', 'concurrency and boundary batteries executed (receipts, supplier, stock, redrive, recurring, HR, accounting validation)', 'F-006,F-047,F-052'),
 'T-P17-001': ('EXECUTED-PARTIAL', 'SEC-014, SEC-011, SEC-012, IDOR, mass assignment, injection executed; SEC-013 staging comparison not executable locally', 'F-020'),
 'T-P17-002': ('EXECUTED-PARTIAL', 'cashier matrix executed (F-020/F-033 reproduced); full module x level matrix not executed', 'F-020,F-033'),
 'T-P18-001': ('EXECUTED-PARTIAL', 'REL-004/009 PASS, REL-025 FAIL (F-052), REL-027 partial', 'F-052'),
 'T-P18-002': ('PLANNED', '100-VU load and 50k-line report benchmark not run', ''),
 'T-P19-001': ('EXECUTED-PARTIAL', 'HC-04, HC-05, FIN-016, FIN-020, CFR-004, HC-10 run; FIN-019 and CFR-006 partial', 'F-049'),
 'T-P20-001': ('EXECUTED-PARTIAL', 'findings/ folders for all 54 findings (10 files each); root cause CONFIRMED for most; remaining carry candidate causes', ''),
 'T-P21-001': ('EXECUTED-PASS', '57 task register (all NOT STARTED), 58 index, 59 regression register, 56 policy package generated', ''),
}
rows = list(csv.reader(open('51-phased-audit-task-register.csv', encoding='utf-8')))
h = rows[0]; ix = {k: i for i, k in enumerate(h)}
n = 0
for r in rows[1:]:
    if r[0] in U:
        st, res, fids = U[r[0]]
        r[ix['status']] = st; r[ix['actual_result']] = res
        if fids: r[ix['finding_ids']] = fids
        n += 1
csv.writer(open('51-phased-audit-task-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
print('updated', n)
