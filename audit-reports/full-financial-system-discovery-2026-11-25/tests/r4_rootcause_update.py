"""RUN-004 root-cause reconciliation (idempotent). Upgrades root_cause_first_break only where source/runtime review in RUN-004 supports it;
keeps PROBABLE with the exact gap where it does not; corrects F-064 (duplicate-table claim retracted). Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
U = {
 'F-012': ('CONFIRMED (RUN-004 source): three independent queries define revenue differently: Sales = SUM(o.total_amount) tax-inclusive (report_handler.go:2177), Profitability = SUM(oi.total_price) net of tax for non-cancelled orders (report_handler.go:2367), P&L = GL revenue-type accounts (financial_report_usecase.go:208); no shared definition or label; numeric bridge in evidence/network/RUN-003/PH-12/AUD-REV-TRACE.txt',
           'backend/internal/handler/http/report_handler.go:2177; backend/internal/handler/http/report_handler.go:2367; backend/internal/usecase/financial_report_usecase.go:208'),
 'F-019': ('CONFIRMED (RUN-004 runtime SQL): only 13 of the 70 public tables with a tenant_id column have row-level security (RLS migrations 164-171 are the only ones); the application role pos_app has UPDATE on orders and payments and no immutability trigger exists on payments (orders carry only stock-transition triggers); journal_entry_lines is not updatable by pos_app',
           'backend/migrations/164_rls_user_auth_audit_logs.sql; backend/migrations/171_rls_e_invoice_submissions.sql'),
 'F-024': ('CONFIRMED (RUN-004 source): InvoiceForm computes discountAmount, taxAmount and lineTotal with raw floating-point multiplication and sends them unrounded (frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:620-650)',
           'frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:623; frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:637'),
 'F-025': ('CONFIRMED (RUN-004 source): handlers pass database errors to clients, e.g. backend/cmd/api/main.go:1302 (fmt.Sprintf("Failed to create user: %v", err)) and respondError(..., err.Error()) at admin_theme_handler_v2.go:80,120,154,361; 128 response sites in internal/handler/http use err.Error()',
           'backend/cmd/api/main.go:1302; backend/internal/handler/http/admin_theme_handler_v2.go:80'),
 'F-026': ('CONFIRMED (RUN-004 source): card payments resolve the debit account to 1113 Card Payment Clearing (auto_posting_engine.go:615-616; domain/auto_posting.go:37; seeded in migration 156:35); no code or migration references 1113 anywhere else, so no settlement/clearing flow exists',
           'backend/internal/usecase/auto_posting_engine.go:615; backend/internal/domain/auto_posting.go:37; backend/migrations/156_seed_chart_of_accounts.sql:35'),
 'F-027': ('CONFIRMED (RUN-004 source, partial scope): only 6 dashboard page files reference a permission hook; e.g. settings/roles/page.tsx has no guard reference, so the shell renders for any signed-in role while the API returns 403; other pages guard inconsistently',
           'frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx'),
 'F-030': ('CONFIRMED (RUN-004 source): POST /users handler decodes name/email/phone/password with no validation, hashes the (possibly empty) password and calls CreateUserInTenant, then assigns the EMPLOYEE role with the error discarded (_ =) in separate non-atomic steps (backend/cmd/api/main.go:1261-1335)',
           'backend/cmd/api/main.go:1261; backend/cmd/api/main.go:1300'),
 'F-031': ('CONFIRMED (RUN-004 source): error mapping is done by passing err text to the client (see F-025 sites) and returning 500 for constraint violations (e.g. main.go:1302; hr handlers rely on DB check constraints)',
           'backend/cmd/api/main.go:1302'),
 'F-039': ('CONFIRMED (RUN-004 source + RUN-003 runtime): CreateJournalEntryRequest.Attachments and JournalEntry.Attachments are json.RawMessage passthroughs with no size or schema limit (domain/accounting.go:95,139); a 1 MB attachment was stored (HC-05)',
           'backend/internal/domain/accounting.go:95; backend/internal/domain/accounting.go:139'),
 'F-044': ('CONFIRMED (RUN-004 SQL): no trigger or code path updates kds_queue when an order is cancelled (pg_trigger on orders lists only order_stock_transition/order_stock_items); RUN-004 shows 3 cancelled orders still active in the queue (AUD-R4-KDS-after-read.json)',
           'backend/internal/repository/kds_repo.go:307'),
 'F-048': ('CONFIRMED (RUN-004 source): the handler only requires non-empty codes and a positive rate (accounting_advanced_handler.go:475-482); there is no base != quote check',
           'backend/internal/handler/http/accounting_advanced_handler.go:481'),
 'F-053': ('CONFIRMED (RUN-004 source): clock_in/clock_out are parsed only as RFC3339 and a failed parse is silently ignored (hr_attendance_handler.go:235-243), so "09:00" is dropped and the row stores NULL times',
           'backend/internal/handler/http/hr_attendance_handler.go:235'),
 'F-054': ('CONFIRMED (RUN-004 source): total_days is taken from the client when non-zero (hr_leave_handler.go:207-210; update :328-329) and no allowance or date-consistency check exists',
           'backend/internal/handler/http/hr_leave_handler.go:207; backend/internal/handler/http/hr_leave_handler.go:328'),
 'F-055': ('CONFIRMED (RUN-004 source): ExtendTenantTrial validates only days > 0 (subscription_handler.go:128-138); there is no upper bound',
           'backend/internal/handler/http/subscription_handler.go:137'),
 'F-056': ('CONFIRMED (RUN-004 source): contract create/update validate only the date format (contract_handler.go:165-170, 221-226); no end >= start rule; resignation create has the same gap (resignation_repo.go:117-143)',
           'backend/internal/handler/http/contract_handler.go:165; backend/internal/repository/resignation_repo.go:117'),
 'F-057': ('CONFIRMED (RUN-004 source): CompleteDelivery only changes the order status to delivered (order_usecase.go:433-447) and records no collection or payment',
           'backend/internal/usecase/order_usecase.go:433'),
 'F-058': ('CONFIRMED (RUN-004 source): driver create requires non-empty fields only (admin_driver_handler.go:146-150); no phone format validation',
           'backend/internal/handler/http/admin_driver_handler.go:146'),
 'F-059': ('CONFIRMED (RUN-004 source absence + RUN-003 runtime): the theme handler/service persist config JSON without colour or value validation (no validate call in the theme service; url(javascript:) was accepted and served)',
           'backend/internal/handler/http/admin_theme_handler_v2.go:293'),
 'F-060': ('CONFIRMED (RUN-004 source): GetThemeHistory returns 501 Not Implemented (admin_theme_handler.go:68)',
           'backend/internal/handler/http/admin_theme_handler.go:68'),
}
STAY = {
 'F-018': ' [RUN-004: kept PROBABLE. Source shows default go-redis options (cache.go:232-237, ping timeout 2 s at startup only); the 30 s hang implies repeated per-request calls but the middleware call site is not yet identified. Required next step: trace session lookup call chain and set a Redis-down fault injection with per-call timing.]',
 'F-021': ' [RUN-004: kept PROBABLE. The stored-column write path is confirmed (supplier_repo.go:44,88,116,189) but the exact sequence that produced -22,000 versus AP 800 has not been replayed step by step. Required next step: replay supplier fixture transactions with before/after current_balance.]',
 'F-038': ' [RUN-004: kept PROBABLE. Static and runtime consistent (landing register page.tsx:59-63) but not reproduced at human pace. Required next step: paced headed-Chrome repro with network timeline.]',
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
ix = {k: i for i, k in enumerate(hd)}
for r in rows[1:]:
    f = r[0]
    if f in U and not r[ix['root_cause_first_break']].startswith('CONFIRMED (RUN-004'):
        r[ix['root_cause_first_break']] = U[f][0]
        r[ix['source_refs']] = (r[ix['source_refs']] + '; ' if r[ix['source_refs']] else '') + U[f][1]
    if f in STAY and '[RUN-004: kept PROBABLE' not in r[ix['root_cause_first_break']]:
        r[ix['root_cause_first_break']] += STAY[f]
    if f == 'F-064' and not r[ix['title']].startswith('Tables API returns 500'):
        r[ix['title']] = 'Tables API returns 500 and a generic message for an over-long name (the duplicate-table-number claim was RETRACTED in RUN-004)'
        r[ix['severity']] = 'Low'
        r[ix['evidence_status']] = 'REPRODUCED (name length); duplicate claim RETRACTED'
        r[ix['actual']] = ('RUN-004 correction: duplicate table numbers are rejected with 409 "table number already exists for this restaurant" (constraint uq_restaurant_tables_number, table_repo.go:81-83). '
                           'The RUN-003 result (TBL-01 201) was a misread: table id 1 had been renumbered to 4 at 20:53, so number 1 did not exist when id 2 was created at 21:10. '
                           'Still reproduced: a 1,800-char name -> 500 "Failed to create table" (column name varchar(100); the DB error is not mapped to 400).')
        r[ix['root_cause_first_break']] = ('CONFIRMED (RUN-004): the name length is not validated in the handler and the varchar(100) overflow surfaces as a generic 500 (table_repo.go:78-86 wraps the error as "failed to create table"; handler table_handler.go:42 maps it to 500). '
                                           'Uniqueness IS enforced by uq_restaurant_tables_number.')
        r[ix['source_refs']] = 'backend/internal/repository/table_repo.go:81; backend/internal/handler/http/table_handler.go:42'
        r[ix['evidence_paths']] += '; evidence/database/RUN-004/PH-19/AUD-R4-F064-dup-table.json; tests/runtime/ph19_r4_tables_dup.mjs'
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
import collections, re
c = collections.Counter()
for r in rows[1:]:
    m = re.match(r'(CONFIRMED|PROBABLE|CANDIDATE|UNKNOWN)', r[ix['root_cause_first_break']].upper())
    c[m.group(1) if m else 'UNSTATED'] += 1
print(len(rows) - 1, dict(c))
