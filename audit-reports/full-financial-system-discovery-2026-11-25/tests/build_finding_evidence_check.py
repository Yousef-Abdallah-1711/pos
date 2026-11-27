"""RUN-003 audit-only script. Re-runnable.
1. Checks every evidence token in 22-confirmed-findings-register.csv resolves to a real file (prefix-glob aware,
   because the register uses shorthand like 'evidence/runtime/05' for the file '05-*').
2. Writes 55-finding-evidence-and-root-cause-verification.csv (one row per finding).
3. Appends RUN-003 results for F-014 (idempotent) to 22 (evidence_paths/actual) and 34b.
Run from the audit directory: python -I tests/build_finding_evidence_check.py
"""
import csv, glob, os, re

os.chdir(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
IGNORE = {"evidence/runtime", "runtime", "screenshots", "reports", "database", "network", "evidence", "explore-A", "json"}
ROW_REF = re.compile(r"^TX-\d+$")          # references to rows in 44, not files
TOP_LEVEL = re.compile(r"^\d{2}[a-z]?-.*\.(csv|md)$")  # references to top-level audit files


def resolve(tok):
    base = tok.strip().rstrip(".,")
    if os.path.exists(base):  # literal repo-relative path inside the audit dir (e.g. tests/runtime/x.mjs)
        return [base]
    if base.startswith("evidence/"):
        base = base[len("evidence/"):]
    if re.fullmatch(r"(runtime/)?\d{2}[a-z]?", base):
        n = base.split("/")[-1]
        return glob.glob(f"evidence/runtime/{n}-*") + glob.glob(f"evidence/runtime/{n}.*")
    if TOP_LEVEL.match(base):
        return glob.glob(base)
    found = []
    for d in ("", "runtime/", "reports/", "screenshots/", "network/", "database/"):
        found += glob.glob(f"evidence/{d}{base}*")
    return found


with open("22-confirmed-findings-register.csv", encoding="utf-8", newline="") as h:
    rd = csv.DictReader(h)
    cols = rd.fieldnames
    rows = list(rd)

NOTE = ("RUN-003 (headed Chrome 150, tenant 4 owner): POST /hr/employees -> 500 reproduced; DB before=after (no orphan rows); "
        "UI shows only generic 'Request failed with status code 500' - raw DB error is in the API response body (see evidence). ")
F14_EVIDENCE = ("evidence/screenshots/RUN-003/PH-10/AUD-HR-001-before-submit.png; evidence/screenshots/RUN-003/PH-10/AUD-HR-001-after-submit.png; "
                "evidence/network/RUN-003/PH-10/AUD-HR-001.json; evidence/database/RUN-003/PH-10/AUD-HR-001-before-after.txt; tests/runtime/ph10_emp_create_headed.mjs")
F001_EVIDENCE = ("evidence/screenshots/RUN-003/PH-05/AUD-ONB-001-after-submit.png; evidence/network/RUN-003/PH-05/AUD-ONB-001.json; "
                 "evidence/database/RUN-003/PH-05/AUD-ONB-001-before-after.txt; tests/runtime/ph05_register_headed.mjs")
F001_NOTE = ("RUN-003 headed-Chrome revalidation: new tenant 5 registered via real UI (POST /auth/register 201, no seeding) -> "
             "chart_of_accounts 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0 (STILL REPRODUCIBLE at SHA 133260d). ")
F001_NOTE += ("ROOT CAUSE (T-P03-001, CONFIRMED): handler auth_handler.go:47 Register -> usecase auth_usecase.go:99 RegisterContext -> "
              "repository auth_repo.go:42 CreateTenantWithUserContext (single tx, commit ~l.210) inserts tenants/users/restaurants/themes_v2/roles(seedDefaultRolesAndPermissions)/user_roles only; "
              "no chart_of_accounts/auto_posting_rules/fiscal_periods/tax_periods call. COA seed exists only in migration 156 (one-time, 'for existing restaurants') and posting rules in migration 152. "
              "Other COA writers (account_repo.go:136 CreateAccount, accounting_advanced_repo.go:~1833) are manual/advanced paths, not wired to registration. ")
RUN3 = {"F-014": (F14_EVIDENCE, NOTE), "F-001": (F001_EVIDENCE, F001_NOTE)}
RUN3["F-010"] = ("evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json; tests/runtime/ph20_evidence_gaps.mjs; evidence/screenshots/RUN-003/PH-20/AUD-EVD-F010-vat-report.png", "RUN-003 (AUD-EVD-F010, tenant 4): VAT report API output_vat=700 vs SQL sum(orders.tax_amount, paid)=14,000 (+700 customer-invoice VAT) and GL VAT Payable (acct 87) credit balance 14,700; residual 14,000 = POS order VAT missing from report (STILL REPRODUCIBLE). ")
RUN3["F-011"] = ("evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json; tests/runtime/ph20_evidence_gaps.mjs; evidence/screenshots/RUN-003/PH-20/AUD-EVD-F011-inventory-valuation.png", "RUN-003 (AUD-EVD-F011): inventory valuation report 192,000 (48 x 4,000) = SQL sum(quantity_in_stock*cost); GL Inventory (acct 1130) balance -36,000; residual 228,000; opening stock has no journal (STILL REPRODUCIBLE). ")
RUN3["F-016"] = ("evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json; tests/runtime/ph20_evidence_gaps.mjs; evidence/screenshots/RUN-003/PH-20/AUD-EVD-F016-payroll-default.png", "RUN-003 (AUD-EVD-F016): default window start_date=2026-10-27,end_date=2026-11-26 returns total 0 and UI says 'Create your first salary record', while salaries row id 1 (tenant 4, period 2026-11-01..2026-11-30, paid) exists (STILL REPRODUCIBLE). ")
RUN3["F-023"] = ("evidence/network/RUN-003/PH-20/AUD-EVD-gaps.json; tests/runtime/ph20_evidence_gaps.mjs; evidence/screenshots/RUN-003/PH-20/AUD-EVD-F023-pos.png", "RUN-003 (AUD-EVD-F023): POS shows '$10,000.00' for an EGP product; no EGP label on page; AUD-POS-002 shows every amount with '$' (cart subtotal $20,000.00) (STILL REPRODUCIBLE). ")
RUN3["F-008"] = ("evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-form-before-submit.png; evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-after-submit.png; evidence/network/RUN-003/PH-09/AUD-CINV-001.json; evidence/database/RUN-003/PH-09/AUD-CINV-001-before-after.txt; tests/runtime/ph09_custinv_draft_headed.mjs", "RUN-003 FULL REPRODUCTION (AUD-CINV-001, headed Chrome): new invoice id 3 status=draft; journal JE-2026-00050 posted at creation Dr 1120 8,320 / Cr 2120 1,120 / Cr 4100 7,200 (STILL REPRODUCIBLE). ")
RUN3["F-007"] = ("evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-form-before-submit.png; evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-after-submit.png; evidence/network/RUN-003/PH-09/AUD-CINV-001.json; evidence/database/RUN-003/PH-09/AUD-CINV-001-before-after.txt; tests/runtime/ph09_custinv_draft_headed.mjs", "RUN-003 (AUD-CINV-001): UI Tax 1,008 / Total 8,208 (matches discount-first oracle) vs stored tax 1,120 / total 8,320 (VAT on pre-discount subtotal): +112 tax, +112 AR (STILL REPRODUCIBLE). ")
RUN3["F-039"] = ("evidence/network/RUN-003/PH-19/AUD-REG-HC04-HC05.json; evidence/database/RUN-003/PH-19/AUD-REG-HC05-before-after.txt; tests/runtime/ph19_regressions.mjs", "")
RUN3["F-003"] = ("evidence/network/RUN-003/PH-12/AUD-MET-FIN.json; evidence/network/RUN-003/PH-12/AUD-MET-FIN.earlier-run.json; evidence/database/RUN-003/PH-12/AUD-CF-TRACE.txt; evidence/screenshots/RUN-003/PH-12/AUD-MET-FIN-cf.png; evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json; evidence/database/RUN-003/PH-14/AUD-MONTH.json; tests/runtime/ph12_fin_statements.mjs", "RUN-003 earlier run: cash-flow 82,710.01 vs GL cash+bank 115,010.01 (residual 32,300.00). Re-measured after more activity: report 172,710.01 vs GL 1110+1111 = 280,010.01, residual 107,300.00. ROOT CAUSE CONFIRMED (AUD-CF-TRACE): usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment that never reads cash/bank accounts and omits every working-capital movement (AR -6,520; inventory +96,000; input VAT -2,800; card clearing -11,000; AP +800; VAT payable +30,820 = +107,300); investing = '12%' debit-credit has inverted sign for asset purchases; financing excludes 21xx/2120 liabilities. CLEAN-TENANT RE-DERIVATION (AUD-MONTH): report 69,290.00 vs cash change 74,480.00, residual 5,190.00 = AR -1,850 + inventory +1,000 + input VAT -2,899.40 + AP +7,779.40 + VAT payable +1,160. DIRECT TEST (AUD-ACC-FLOW CF-INV): posting a 20,000 fixed-asset purchase (cash fell 20,000) INCREASED investing_cash_flow by +20,000 and net cash flow from 167,190.01 to 229,190.01. ")
RUN3["F-004"] = ("evidence/screenshots/RUN-003/PH-07/AUD-POS-001-before-submit.png; evidence/network/RUN-003/PH-07/AUD-POS-001.json; evidence/database/RUN-003/PH-07/AUD-POS-001-before-after.txt; evidence/screenshots/RUN-003/PH-07/AUD-POS-002-cart-totals-overlay-hidden.png; evidence/network/RUN-003/PH-07/AUD-POS-002.json; tests/runtime/ph07_pos_cash_headed.mjs; tests/runtime/ph07_pos_ui_totals.mjs", "RUN-003 (AUD-POS-001, order 38): POS label Tax (14%) vs stored tax 6,000 on subtotal 60,000 (10%); UI value captured with overlay hidden by browser-side style (AUD-POS-002): cart of 2 x 10,000 displays Subtotal $20,000.00, Tax (14%) $2,800.00, Total $22,800.00 whereas server stores 10% (2,000 / 22,000) - STILL REPRODUCIBLE (+800 shown to customer vs stored). ")
RUN3["F-015"] = ("evidence/database/RUN-003/PH-20/AUD-EVD-F008-F015-sql.txt", "RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are status=paid with is_approved=false; journals JE-2026-00007/00049 are Dr 5210 5,000 / Cr 1111 5,000 posted at payment only; no accrual journal exists (STILL REPRODUCIBLE; UI re-run pending T-P10-003). ")
RUN3["F-038"] = ("evidence/network/RUN-003/PH-05/OBS-ONB-001-probes.json; evidence/screenshots/RUN-003/PH-05/OBS-ONB-001-S1-click-200ms-after-typing.png; tests/runtime/ph05_subdomain_probe.mjs", "")
RUN3["F-037"] = ("evidence/network/RUN-003/PH-05/AUD-SET-001-011.json; tests/runtime/ph05_settings_forms.mjs", "")
RUN3["F-006"] = ("evidence/database/RUN-003/PH-09/AUD-AR-RESIDUAL-trace.txt; evidence/database/RUN-003/PH-09/AUD-CINV-FLOW.json; tests/runtime/ph09_customer_flow.mjs", "RUN-003 FRESH-INVOICE RE-REPRODUCTION (AUD-CINV-FLOW): invoice 3,000, 5 concurrent receipts of 1,000 -> all 5 accepted (201), amount_paid 6,000, AR credit -3,000; single sequential overpayment 5,000 is rejected (guard exists but is not serialized). RUN-003 TRANSACTION-LEVEL TRACE (AUD-AR-RESIDUAL): GL 1120 = Dr 5,200 (INV-2026-00001) + Dr 8,320 (INV-2026-00002 draft, F-008) - Cr 7,000 (6 receipts: 2,000 + 5x1,000) = 6,520; aging 8,320; residual 1,800 = receipts 7,000 minus invoice 5,200 = over-settlement of invoice 2 (amount_paid 7,000 > total 5,200, balance clamped to 0). ROOT CAUSE CONFIRMED (source): backend/internal/repository/payment_acct_repo.go:445-449 customer invoice balance SELECT has no FOR UPDATE (supplier twin at :499-503 has FOR UPDATE OF inv), and applyCustomerInvoicePaymentTx :569-582 updates unconditionally with GREATEST(...,0) clamp that hides the overpayment. ")
RUN3["F-040"] = ("evidence/database/RUN-003/PH-06/AUD-INV-COGS-UNTRACKED.txt; evidence/database/RUN-003/PH-06/OBS-INV-001-controlled-A-untracked.json; tests/runtime/ph06_inventory_controlled.mjs", "RUN-003 controlled experiment: untracked product order 39 posts no COGS while Profitability report shows cogs 120; tracked twin order 41 posts COGS 120 (CONFIRMED, auto_posting_engine.go:173). ")
RUN3["F-041"] = ("evidence/database/RUN-003/PH-06/AUD-SUP-FLOW.json; tests/runtime/ph06_supplier_flow.mjs", "RUN-003 supplier flow (API-DIAGNOSTIC): invoice 10x1000 disc 10% VAT 14% => tax 1,400 total 10,400 vs oracle 1,260/10,260; JE Dr 5110 9,000 Dr 1150 1,400 Cr 2101 10,400 (CONFIRMED trigger 142:75-79). ")
RUN3["F-009"] = ("evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json; evidence/database/RUN-003/PH-12/AUD-F009-ROOT-CAUSE.txt; tests/runtime/ph12_batch_apis.mjs", "RUN-003: GET /reports/financial/menu-profitability and /break-even still 500. ROOT CAUSE CONFIRMED: repository/accounting_advanced_repo.go:976 and :1045 select p.cost_price, a column that does not exist on products (the column is cost) -> SQL error mapped to generic 500 (accounting_advanced_handler.go:782). ")
RUN3["F-042"] = ("evidence/database/RUN-003/PH-12/AUD-MET-ORD-KPI.txt; evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json; tests/runtime/ph12_batch_apis.mjs", "RUN-003: order stats total_revenue 529,100 = all 41 orders incl 3 cancelled; delivered_today 30 vs delivered_count 0 (CONFIRMED order_repo.go:859-870). ")
RUN3["F-043"] = ("evidence/database/RUN-003/PH-12/AUD-MET-ORD-KPI.txt; evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json; tests/runtime/ph12_batch_apis.mjs", "RUN-003: KPI net profit 181,540 vs P&L 167,180; cash_position 50,000 vs GL cash+bank 276,270.01 (CONFIRMED accounting_advanced_repo.go:838-880). ")
RUN3["F-044"] = ("evidence/database/RUN-003/PH-12/AUD-MET-ORD-KPI.txt; evidence/network/RUN-003/PH-08/AUD-KDS-001.json; evidence/screenshots/RUN-003/PH-08/AUD-KDS-001-board.png; tests/runtime/ph08_kds_cancelled.mjs", "RUN-003 headed Chrome: 3 cancelled orders visible on KDS board (ACTIVE 25 / NEW 24). ")
RUN3["F-045"] = ("evidence/database/RUN-003/PH-12/AUD-MET-PCF-VAT-INV.txt; evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json; tests/runtime/ph12_batch_apis.mjs", "RUN-003: PCF Nov cash_in 13,000/cash_out 32,400 vs paid POS orders 319,660; forecast negative (CONFIRMED accounting_advanced_repo.go:1105-1175). ")
RUN3["F-046"] = ("evidence/network/RUN-003/PH-12/AUD-MET-BATCH2.json; evidence/database/RUN-003/PH-12/AUD-MET-BATCH2-compare.json; tests/runtime/ph12_batch_apis2.mjs", "RUN-003: turnover opening_stock_qty == current stock (33/21); waste_qty 0 vs movements waste 1 (REPRODUCED). ")
RUN3["F-047"] = ("evidence/database/RUN-003/PH-11/AUD-ACC-FLOW2.json; tests/runtime/ph11_accounting_flow2.mjs", "RUN-003: recurring run twice -> journals 98 and 99, rent +700 each (CONFIRMED accounting_advanced_repo.go:396-420). ")
RUN3["F-048"] = ("evidence/database/RUN-003/PH-11/AUD-ACC-FLOW2.json; tests/runtime/ph11_accounting_flow2.mjs", "RUN-003: USD/USD pair accepted (id 2). ")
RUN3["F-013"] = ("evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json; tests/runtime/ph11_accounting_flow.mjs", "RUN-003 (AUD-ACC-FLOW JE-13/JE-14): journals dated 2019-01-15 and 2030-06-01 post successfully and auto-create fiscal_period 2019-01/2030-06; GET /accounting/fiscal-periods is 404 (no period API). ")
RUN3["F-029"] = ("evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json; tests/runtime/ph11_accounting_flow.mjs", "RUN-003 (AUD-ACC-FLOW BNK-01): new bank account opening_balance 12,000 linked to GL 1112 creates no journal (GL 1112 stays 0) (REPRODUCED). ")
RUN3["F-031"] = ("evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json; tests/runtime/ph11_accounting_flow.mjs", "RUN-003 (AUD-ACC-FLOW ACC-02): duplicate account code returns raw DB text pq: duplicate key value violates unique constraint chart_of_accounts_tenant_id_restaurant_id_account_code_key (REPRODUCED). RUN-003 HR matrix: duplicate attendance/salary, invalid pay period, negative salary and overlapping leave return HTTP 500 with raw pq/constraint text (unique_employee_attendance_date, unique_employee_salary_period, valid_pay_period, positive_amounts, valid_leave_period) (REPRODUCED). ")
RUN3["F-049"] = ("evidence/database/RUN-003/PH-19/AUD-REFUND-01-paid-order-cancel.json; evidence/database/RUN-003/PH-19/AUD-REG-HC10-FIN019-CFR006.json; tests/runtime/ph19_paid_cancel.mjs; tests/runtime/ph19_regs2.mjs; evidence/database/RUN-003/PH-07/AUD-ORDER-TYPES.json", "RUN-003: paid order cannot be cancelled/refunded/voided through any endpoint; explicit guard message (REPRODUCED). ")
RUN3["F-050"] = ("evidence/database/RUN-003/PH-19/AUD-REFUND-01-paid-order-cancel.json; evidence/database/RUN-003/PH-19/AUD-REG-HC10-FIN019-CFR006.json; tests/runtime/ph19_paid_cancel.mjs; tests/runtime/ph19_regs2.mjs", "RUN-003: PUT status cancelled on paid order -> 500. Same class: unsigned Paymob webhook -> 500 and refund of unknown payment -> 500 (AUD-SAAS-MATRIX PAY-02/PAY-03). Also POST /pos/orders with expected_total mismatch and with an unknown product_id return 500 \"Failed to create order\" instead of 4xx (AUD-ORDER-TYPES OT-04/OT-08). ")
RUN3["F-020"] = ("evidence/network/RUN-003/PH-17/AUD-SEC-BATTERY.json; tests/runtime/ph17_security.mjs", "RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default Employee role: HR READ, SALARIES NONE) GET /hr/employees -> 200 with base_salary 5000 and contact data; accounting/report/admin routes 403/401 (REPRODUCED). ")
RUN3["F-052"] = ("evidence/database/RUN-003/PH-18/AUD-REL-BATTERY.json; tests/runtime/ph18_reliability.mjs", "RUN-003 REL-025: identical Idempotency-Key replay creates 2 receipts (sum 800) (REPRODUCED). ")
RUN3["F-017"] = ("evidence/database/RUN-003/PH-18/AUD-REL-BATTERY.json; tests/runtime/ph18_reliability.mjs", "RUN-003 REL-001: same order POST twice with identical Idempotency-Key creates 2 orders (REPRODUCED). ")
RUN3["F-053"] = ("evidence/database/RUN-003/PH-10/AUD-HR-MATRIX.json; evidence/database/RUN-003/PH-10/AUD-HR-MATRIX2.json; tests/runtime/ph10_hr_matrix.mjs; tests/runtime/ph10_hr_matrix2.mjs", "RUN-003: attendance with 09:00/19:30 stored NULL times and 0 hours (201); RFC3339 computes 8+2.5 (REPRODUCED). ")
RUN3["F-054"] = ("evidence/database/RUN-003/PH-10/AUD-HR-MATRIX.json; evidence/database/RUN-003/PH-10/AUD-HR-MATRIX2.json; tests/runtime/ph10_hr_matrix.mjs; tests/runtime/ph10_hr_matrix2.mjs", "RUN-003: 30-day annual leave approved over allowance 21; total_days 40 stored for a 3-day range (REPRODUCED). ")
RUN3["F-055"] = ("evidence/database/RUN-003/PH-13/AUD-SAAS-MATRIX.json; tests/runtime/ph13_saas.mjs", "RUN-003: extend-trial 100000 days accepted -> 2300-10-02 (REPRODUCED; tenant 8 restored by fixture). ")
RUN3["F-002"] = ("evidence/database/RUN-003/PH-14/AUD-MONTH.json; tests/runtime/ph14_month.mjs", "RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=false, equation_difference -27,710.00 = net loss exactly (REPRODUCED). ")
RUN3["F-056"] = ("evidence/database/RUN-003/PH-10/AUD-HR-MODULES.json; tests/runtime/ph10_hr_modules.mjs", "RUN-003: contract end 2026-10-01 < start 2026-11-01 stored (201); resignation last day before resignation date stored (201) (REPRODUCED). ")
RUN3["F-018"] = ("evidence/network/RUN-003/PH-18/AUD-REL-REDIS.json; tests/runtime/ph18_redis_outage.mjs", "RUN-003: with the isolated Redis stopped, authenticated GET /products and POST /auth/login both hang until the 30 s client timeout while /health stays 200 in 6 ms and the public menu answers in 11 ms; after Redis returns the previous session is rejected (401) (REPRODUCED). ")
RUN3["F-057"] = ("evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json; evidence/database/RUN-003/PH-08/AUD-THEME2.json; tests/runtime/ph08_theme_driver.mjs; tests/runtime/ph08_theme2.mjs", "RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): delivered COD order remained payment_status pending with 0 journals (REPRODUCED). ")
RUN3["F-058"] = ("evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json; evidence/database/RUN-003/PH-08/AUD-THEME2.json; tests/runtime/ph08_theme_driver.mjs; tests/runtime/ph08_theme2.mjs", "RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): driver phone abc accepted (201) (REPRODUCED). ")
RUN3["F-059"] = ("evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json; evidence/database/RUN-003/PH-08/AUD-THEME2.json; tests/runtime/ph08_theme_driver.mjs; tests/runtime/ph08_theme2.mjs", "RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): invalid colors and url(javascript:) stored and served by public theme endpoint (REPRODUCED). ")
RUN3["F-060"] = ("evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json; evidence/database/RUN-003/PH-08/AUD-THEME2.json; tests/runtime/ph08_theme_driver.mjs; tests/runtime/ph08_theme2.mjs", "RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): history 501; active theme deleted with 204 (REPRODUCED). ")
RUN3["F-061"] = ("evidence/database/RUN-003/PH-16/AUD-MASTERDATA.json; evidence/database/RUN-003/PH-16/AUD-MD-NEGPRICE.json; tests/runtime/ph16_masterdata.mjs; tests/runtime/ph16_negprice.mjs", "RUN-003: zero-total paid order 73 -> outbox retryable, error at least two journal lines are required (CONFIRMED auto_posting_engine.go:535). ")
RUN3["F-062"] = ("evidence/database/RUN-003/PH-16/AUD-MASTERDATA.json; evidence/database/RUN-003/PH-16/AUD-MD-NEGPRICE.json; tests/runtime/ph16_masterdata.mjs; tests/runtime/ph16_negprice.mjs", "RUN-003 master-data battery: negative/zero price, negative cost/stock, empty name, discount above price, 5-decimal price, duplicate/malformed supplier/customer accepted (CONFIRMED product_usecase.go:294-306). ")
RUN3["F-063"] = ("evidence/database/RUN-003/PH-11/AUD-IMPORT.json; tests/runtime/ph11_import.mjs; evidence/database/RUN-003/PH-18/AUD-JOURNAL-SCALE.json", "RUN-003 import battery: re-upload duplicates, all-failed imports answer 200 success, 6,000-journal import completed server-side after the client socket closed (REPRODUCED). ")
RUN3["F-064"] = ("evidence/database/RUN-003/PH-16/AUD-TABLES.json; tests/runtime/ph16_tables.mjs", "RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). ")
RUN3["F-065"] = ("evidence/network/RUN-003/PH-06/AUD-UIFORM-CATEGORY.json; evidence/screenshots/RUN-003/PH-06/AUD-UIFORM-category-name+display-order.png; tests/runtime/ph06_ui_category.mjs; evidence/network/RUN-003/PH-06/AUD-UIFORMS.json", "RUN-003 real-UI reproduction: display order typed -> string payload -> 400 Invalid request body; name only -> 201 (CONFIRMED CategoryFormPage.tsx:27,129 vs ProductForm.tsx:100). ")
RUN3["F-022"] = ("evidence/network/RUN-003/PH-06/AUD-UIFORM-F022-public-menu-500.txt; evidence/network/RUN-003/PH-06/AUD-UIFORMS.json; evidence/network/RUN-003/PH-08/AUD-STORE-002-public-api.txt; evidence/network/RUN-003/PH-08/AUD-STORE-003.json; evidence/network/RUN-003/PH-08/AUD-STORE-004.json; tests/runtime/ph06_ui_forms.mjs", "RUN-003: (1) REAL-UI reproduction: product form without a category sends category_id 0 -> POST /products 201 -> public menu/products API 500 for the whole tenant; (2) storefront UI zero items with default isolated CORS (environment), 3 items after the wildcard origin was allowed. ")
RUN3["F-025"] = ("evidence/network/RUN-003/PH-06/AUD-UIFORMS2.json; evidence/network/RUN-003/PH-06/AUD-UIFORMS.json; tests/runtime/ph06_ui_forms2.mjs", "RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-account forms have no client-side length limits; over-long input reaches the API and the UI user receives raw text such as pq: value too long for type character varying(3) (REPRODUCED; the filler used the same long synthetic value in every text field, so the inputs are artificial but the exposure is real). ")
for r in rows:
    if r["finding_id"] in RUN3:
        ev, note = RUN3[r["finding_id"]]
        for path in [x.strip() for x in ev.split(";") if x.strip()]:
            if path not in r["evidence_paths"]:
                r["evidence_paths"] = r["evidence_paths"] + " | " + path
        if note and note[:60] not in r["actual"]:
            r["actual"] = r["actual"] + " | " + note
with open("22-confirmed-findings-register.csv", "w", encoding="utf-8", newline="") as h:
    w = csv.DictWriter(h, fieldnames=cols)
    w.writeheader()
    w.writerows(rows)

out = []
for r in rows:
    toks = [x for x in re.split(r"[;,\s|]+", r["evidence_paths"]) if x and x not in IGNORE]
    toks = [x for x in toks if re.search(r"\d|\.", x)]
    file_toks = [x for x in toks if not ROW_REF.match(x)]
    unresolved = [x for x in file_toks if not resolve(x)]
    precise = [x for x in file_toks if resolve(x)]
    tx_refs = [x for x in toks if ROW_REF.match(x)]
    if not file_toks and tx_refs:
        state = "ROW-REFERENCE-ONLY (file 44; no raw evidence file cited)"
    elif not file_toks:
        state = "NO-RESOLVABLE-EVIDENCE"
    elif unresolved:
        state = "PARTIAL"
    else:
        state = "RESOLVED"
    run3 = {"F-014": "AUD-HR-001", "F-001": "AUD-ONB-001", "F-010": "AUD-EVD-F010", "F-004": "AUD-POS-001", "F-003": "AUD-MET-FIN", "F-039": "AUD-REG-HC05", "F-037": "AUD-SET-001..004", "F-038": "OBS-ONB-001 S1-S5", "F-008": "AUD-EVD-F008; AUD-CINV-001", "F-007": "AUD-CINV-001", "F-015": "AUD-EVD-F015", "F-011": "AUD-EVD-F011", "F-016": "AUD-EVD-F016", "F-023": "AUD-EVD-F023"}.get(r["finding_id"], "")
    out.append({
        "finding_id": r["finding_id"], "severity": r["severity"], "register_evidence_status": r["evidence_status"],
        "evidence_path_check": state, "resolved_file_tokens": len(precise), "unresolved_tokens": ";".join(unresolved),
        "tx_row_refs_in_44": ";".join(x for x in toks if ROW_REF.match(x)),
        "run003_revalidation_test": run3,
        "run003_result": {"F-014": "REPRODUCED (500, no DB change)", "F-001": "REPRODUCED (new tenant 5: COA/rules/fiscal/tax = 0)", "F-010": "REPRODUCED (700 vs 14,000)", "F-004": "REPRODUCED (label 14% vs stored 10%)", "F-003": "REPRODUCED (-32,300)", "F-039": "REPRODUCED (5/5 variants accepted)", "F-037": "REPRODUCED", "F-038": "REPRODUCED (S1,S4,S5)", "F-015": "REPRODUCED (SQL state)", "F-008": "REPRODUCED (draft + journal posted)", "F-007": "REPRODUCED (+112)", "F-011": "REPRODUCED (192,000 vs GL -36,000)", "F-016": "REPRODUCED (0 shown vs 1 paid salary)", "F-023": "REPRODUCED ($ shown)"}.get(r["finding_id"], "NOT REVALIDATED IN RUN-003"),
        "root_cause_confidence": {"F-014": "CONFIRMED (runtime error names constraint; source lines match)", "F-037": "CONFIRMED (handler main.go:845-873 read)", "F-038": "PROBABLE", "F-039": "PROBABLE", "F-015": "PROBABLE (source refs salary_repo.go:344-400, hr_salary_handler.go:470-558 not re-read yet)", "F-008": "PROBABLE (invoice repo not re-read yet)", "F-010": "PROBABLE (report formula not traced yet; T-P12)", "F-011": "PROBABLE (valuation from products table vs GL; opening stock unposted; trace T-P03)", "F-016": "PROBABLE (date filter on pay_period_end vs UI default; trace T-P10)", "F-023": "PROBABLE (hardcoded $ in POS component; trace T-P03)", "F-001": "CONFIRMED (T-P03-001: registration tx auth_repo.go:42-212 has no accounting init; migration 156/152 are one-time seeds for pre-existing tenants; runtime 0 rows on tenant 5)"}.get(r["finding_id"], "NOT YET ASSESSED (Phase 20)"),
        "note": "RUN-003 revalidated; precise evidence paths present" if run3 else ("needs precise evidence paths" if state != "RESOLVED" else ""),
    })
with open("55-finding-evidence-and-root-cause-verification.csv", "w", encoding="utf-8", newline="") as h:
    w = csv.DictWriter(h, fieldnames=list(out[0].keys()))
    w.writeheader()
    w.writerows(out)

# 34b: append RUN-003 executed rows once
new34b = [
    ["EX-RUN003-AUD-ONB-001", "UI-HEADED-CHROME", "Register new tenant 5 via UI, no seeding (F-001 revalidation)", "anonymous->owner", "FAIL",
     "POST /auth/register 201; coa 0, auto_posting_rules 0, fiscal_periods 0, tax_periods 0; evidence/*/RUN-003/PH-05/AUD-ONB-001*"],
    ["EX-RUN003-OBS-ONB-001", "UI-HEADED-CHROME", "Register form: submit before debounced subdomain check returns", "anonymous", "OBSERVED",
     "0 POST sent; red 'Please choose an available subdomain' shown while the same field shows 'Subdomain is available' (stale validation); seen once under automation, needs human-paced repro; evidence/screenshots/RUN-003/PH-05/AUD-ONB-001-probe1-fast-submit.png"],
    ["EX-RUN003-AUD-ENV-001", "UI-HEADED-CHROME", "login + /en/dashboard/reports/financial (isolated stack)", "owner", "PASS",
     "Chrome 150.0.7871.124 headed; login via real UI; financial hub rendered; 60 API calls, only anonymous /auth/session 401; evidence/screenshots/RUN-003/PH-02/"],
    ["EX-RUN003-AUD-HR-001", "UI-HEADED-CHROME", "Create Employee form submit (F-014 revalidation)", "owner", "FAIL",
     "POST /hr/employees 500 (positive_working_hours); DB unchanged; evidence/network/RUN-003/PH-10/AUD-HR-001.json"],
]
with open("34b-executed-test-results.csv", encoding="utf-8", newline="") as h:
    existing = {row[0] for row in csv.reader(h)}
with open("34b-executed-test-results.csv", "a", encoding="utf-8", newline="") as h:
    w = csv.writer(h)
    for row in new34b:
        if row[0] not in existing:
            w.writerow(row)

import collections
print(collections.Counter(o["evidence_path_check"] for o in out))
print("NO-RESOLVABLE:", [o["finding_id"] for o in out if o["evidence_path_check"] == "NO-RESOLVABLE-EVIDENCE"])
print("ROW-REF-ONLY:", [o["finding_id"] for o in out if o["evidence_path_check"].startswith("ROW-REFERENCE")])
print("PARTIAL:", [(o["finding_id"], o["unresolved_tokens"]) for o in out if o["evidence_path_check"] == "PARTIAL"])
