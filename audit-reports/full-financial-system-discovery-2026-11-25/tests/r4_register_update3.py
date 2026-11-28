"""RUN-004 register update 3 (idempotent): F-072, F-073 (bank reconciliation) + 34b rows. Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
DB = 'evidence/database/RUN-004/PH-11/'
F = {
    'F-072': dict(
        title='Bank reconciliation is not tied to the ledger: book balance is the sum of client-supplied line amounts, so a session can be "completed" with invented numbers and stamps the bank account as reconciled',
        category='Software defect (accounting control)', severity='High', evidence_status='REPRODUCED',
        affected_modules='Accounting > Bank reconciliation, bank accounts',
        reproduction='Tenant 4 bank account 2 (ledger 12,000.00, GL 0): create a session with statement balance 777,777.77, add one line referencing a real journal that has no bank line with book_amount 777,777.77, complete the session (API-DIAGNOSTIC)',
        expected='Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts agree with the statement and the ledger',
        actual='Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the ledger balance is 12,000.00 and the journal has 0 bank-account lines. A session with difference 123,456.78 or 500.00 is also stamped completed_at although its status stays in_progress, and the complete call returns HTTP 200',
        monetary_impact_EGP='Control failure: a reconciled flag can be produced without any real matching (illustrated with 777,777.77 against a 12,000.00 ledger balance)',
        root_cause_first_break='CONFIRMED (source + runtime): refreshReconciliationSummaryTx sums bank_reconciliation_lines.book_amount supplied by the client (accounting_advanced_repo.go:2126-2149); AddBankReconciliationLine inserts without checking that the journal touches the bank GL account or that amounts equal the journal line (:113-166); CompleteBankReconciliationSession writes last_reconciled_* from the summed client amounts and sets completed_at regardless of status (:169-230)',
        source_refs='backend/internal/repository/accounting_advanced_repo.go:113; backend/internal/repository/accounting_advanced_repo.go:169; backend/internal/repository/accounting_advanced_repo.go:2126; backend/internal/handler/http/accounting_advanced_handler.go:58',
        related_static_ids='F-029; F-043', evidence_paths=DB + 'AUD-R4-BANK-RECON.json; ' + DB + 'AUD-R4-BANK-RECON2.json; tests/runtime/ph11_r4_bank_recon.mjs; tests/runtime/ph11_r4_bank_recon2.mjs',
        remediation_direction='Compute book balance from the GL for the bank account at the statement date; require lines to reference journal lines on that account and take amounts from them; allow completion only when difference is zero (or an explicit adjustment journal is recorded); do not set completed_at otherwise',
        regression_tests_needed='Session with invented amounts cannot complete; matched real lines complete and update last_reconciled values; non-zero difference returns 409'),
    'F-073': dict(
        title='Bank reconciliation lines accept a journal_entry_id that belongs to another tenant (only a foreign key is checked)',
        category='Security (tenant isolation integrity)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='Accounting > Bank reconciliation',
        reproduction='Tenant 4 owner adds a reconciliation line with journal_entry_id 33, a journal of tenant 2 (API-DIAGNOSTIC)',
        expected='400/404 because the journal is not visible to the tenant',
        actual='HTTP 200; the line is stored in the tenant-4 session referencing tenant 2 journal 33; a nonexistent id returns 400 only because of the foreign key, with the raw constraint text in the error body',
        monetary_impact_EGP='0 direct; cross-tenant reference integrity and error-text leakage (F-025 class)',
        root_cause_first_break='CONFIRMED (source + runtime): AddBankReconciliationLine verifies only that the session exists for the tenant, then inserts journal_entry_id as supplied (accounting_advanced_repo.go:113-166); the table has a plain foreign key and the table has no row-level security (F-019)',
        source_refs='backend/internal/repository/accounting_advanced_repo.go:113; backend/internal/handler/http/accounting_advanced_handler.go:58',
        related_static_ids='F-019; F-025; F-072', evidence_paths=DB + 'AUD-R4-BANK-RECON.json; tests/runtime/ph11_r4_bank_recon.mjs',
        remediation_direction='Validate that the journal (and its line) belongs to the same tenant and restaurant and touches the reconciled bank account; map FK errors to 400 with a generic message; add RLS on the reconciliation tables',
        regression_tests_needed='Foreign-tenant and nonexistent journal ids return 4xx with no stored line and no raw database text'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        v = dict(v, finding_id=fid)
        rows.append([v.get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN004-BANKRECON', 'API-DIAGNOSTIC', 'Bank reconciliation: invented amounts, cross-tenant journal, completion with difference (BR-01..BR-10)', 'owner', 'FAIL', 'invented 777,777.77 line completes the session and stamps last_reconciled_balance (F-072); other-tenant journal accepted (F-073); FK errors leak constraint text'),
    ('EX-RUN004-DEPRECIATION', 'SOURCE+SQL', 'Depreciation / fixed assets capability', 'owner', 'BLOCKED-CAPABILITY', 'fixed_assets table (migration 145) has 0 rows; no route or UI; PostDepreciationEntries has no caller; not implemented (see 33 N23)'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    w = csv.writer(h)
    for r in new:
        if r[0] not in ex:
            w.writerow(r)
print(len(rows) - 1, 'findings')
