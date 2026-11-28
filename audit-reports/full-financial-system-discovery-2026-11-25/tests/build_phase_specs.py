"""Generates phases/phase-NN/{spec.md,tasks.md} (references only, no copying of reports). Audit tooling; idempotent.
Task status lives in 51-phased-audit-task-register.csv and progress.md; tasks.md is an execution instruction file that links to them."""
import csv, glob, os, re

AUDIT = r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
os.chdir(AUDIT)

P = {
 '00': ('Existing audit takeover and file verification', 'COMPLETE', ['49-existing-audit-takeover-and-verification.md', '55-finding-evidence-and-root-cause-verification.csv'], ['22', '34b'], ['Baseline counts, SHA, stashes verified; discrepancies DISC-001..004 recorded; evidence paths verified (55: all RESOLVED).'], []),
 '01': ('Build the complete remaining execution backlog', 'COMPLETE', ['50-remaining-audit-execution-plan.md', '51-phased-audit-task-register.csv'], ['51'], ['Task register generated and extended; new work is recorded in 34b/progress.md and tasks are updated by tests/build_task_register.py.'], []),
 '02': ('Restore the isolated environment and open real Google Chrome', 'COMPLETE', ['47-environment-restoration-report.md', 'evidence/runtime/RUN-003/47-isolated-gate-and-restore-check.txt'], ['40'], ['Isolated stack on persistent volume; real headed Chrome 150 with dedicated profile; checksummed dumps (start, mid and final).', 'Environment notes: API recreate breaks frontends sharing its network namespace (recreate them too); storefront needs wildcard CORS origin for tenant subdomains.'], []),
 '03': ('Complete system architecture and financial surface discovery', 'COMPLETE-WITH-LIMITS', ['03-complete-system-architecture.md', '27b-financial-lineage-chains.md', '52-complete-money-operation-trace.md'], ['08', '11'], ['Money-operation trace (52) covers 16 operations with file:line; F-001 root cause CONFIRMED.', 'Limit: 119 vs 117 route reconciliation (DISC-001) still explained, not re-run.'], ['F-001']),
 '04': ('Exhaustive UI inventory and functional testing', 'PARTIAL', ['06-complete-page-inventory.md', '07-field-level-inventory.csv', '34-complete-manual-ui-test-plan.csv'], ['06', '07', '34'], ['Dashboard route matrix in headed Chrome (EN, AR/RTL, mobile) recorded in evidence/network/RUN-003/PH-04.', 'Limit: the 1,491-row field matrices and every dialog/filter were NOT executed individually; field-level validation was exercised through API batteries (Phases 16-17).'], ['F-027', 'F-032']),
 '05': ('Registration, onboarding, authentication, settings and roles', 'PARTIAL', ['41-fresh-tenant-onboarding-results.md'], ['41'], ['F-001 re-validated in headed Chrome (tenant 5 preserved); subdomain observation F-038; settings forms F-037; auth: login/origin/rate limit/session tamper executed.', 'Limit: roles/users/access-control/billing settings forms and password recovery not executed in the UI.'], ['F-001', 'F-037', 'F-038']),
 '06': ('Products, categories, inventory and suppliers', 'PARTIAL', ['26-business-workflow-catalog.md'], ['33'], ['Controlled stock-vs-ledger experiment (OBS-INV-001 not reproduced on a fresh tracked product); supplier invoice/payment/concurrency matrix.', 'Limit: product/category/supplier flows were driven by API-diagnostic scripts, not by UI forms.'], ['F-011', 'F-034', 'F-040', 'F-041']),
 '07': ('POS, cashier, shifts, orders and payments', 'PARTIAL', ['10-financial-transaction-lifecycles.md', '52-complete-money-operation-trace.md'], ['35b'], ['Cash POS order in headed Chrome with legible UI totals; duplicates, concurrency, cancel/refund matrix executed.', 'BLOCKED: card/wallet settlement (F-005/F-026); no shift/cash-drawer module was found in the product (capability matrix 33).'], ['F-004', 'F-005', 'F-017', 'F-023', 'F-026', 'F-028', 'F-049', 'F-050', 'F-052']),
 '08': ('Storefront, theme builder, KDS, tables and delivery', 'PARTIAL', ['26-business-workflow-catalog.md'], ['35b'], ['Storefront order placed end to end in headed Chrome (website order, KDS entry on confirm, paid posting); KDS board shows cancelled orders (F-044).', 'Limit: theme builder, tables layout, drivers/delivery and QR not exercised beyond page loads.'], ['F-022', 'F-044']),
 '09': ('Customers, customer invoices and receivables', 'COMPLETE-WITH-LIMITS', ['44-transaction-level-reconciliation.csv'], ['35b'], ['Invoice creation, tax/discount, receipts, overpayment, concurrency and AR residual trace; AR 1,800 residual traced to F-006.', 'Limit: credit notes not supported/tested.'], ['F-006', 'F-007', 'F-008', 'F-024', 'F-052']),
 '10': ('Complete HR and payroll review', 'PARTIAL', ['43-report-metric-runtime-results.csv'], ['35b'], ['Attendance, leave, salary arithmetic, mark-paid idempotency and posting verified through API matrix.', 'BLOCKED: UI employee creation (F-014); training/insurance/contracts/resignations/performance not exercised.'], ['F-014', 'F-015', 'F-016', 'F-053', 'F-054']),
 '11': ('Accounting, banking and advanced accounting review', 'PARTIAL', ['11-financial-posting-matrix.csv', '13-accounting-event-architecture.md'], ['35b'], ['Account CRUD, journal validation/post/void, capital/withdrawal, asset purchase, recurring, budget, FX, bank opening via API matrix.', 'BLOCKED: period close (no API, F-013). Bank reconciliation, depreciation and import/export flows not executed.'], ['F-013', 'F-029', 'F-047', 'F-048']),
 '12': ('Exhaustive reporting and metric verification', 'COMPLETE-WITH-LIMITS', ['30-report-metric-inventory.csv', '30c-metric-coverage-reconciliation.csv', '43-report-metric-runtime-results.csv'], ['30', '43'], ['113 of 119 coverage rows have a runtime comparison (API + SQL oracle; Chrome rendering for statement pages); MENU/BE blocked by F-009; revenue/cash/VAT bridges explained.', 'Limit: exports/print and multi-dataset boundary variants not executed for every metric; 6 rows still NOT EXECUTED (see 30c).'], ['F-002', 'F-003', 'F-009', 'F-010', 'F-012', 'F-042', 'F-043', 'F-045', 'F-046']),
 '13': ('SaaS billing, subscriptions and Super Admin', 'PARTIAL', ['28-security-and-reliability-architecture.md'], ['35b'], ['Plans, trial, disable/enable, extend-trial validation, mock provider and metrics executed.', 'BLOCKED-CAPABILITY: paid checkout disabled (503) so successful payment, webhook settlement and refund success path are not executable.'], ['F-055', 'F-036']),
 '14': ('Realistic restaurant operating-month simulation', 'PARTIAL', ['42-complete-operating-month-reconciliation.md'], ['35'], ['Compressed month on fresh tenant: independent ledger = system GL on every account; F-002/F-003 reproduced with exact residuals.', 'BLOCKED: payroll (F-014), dated sales (orders cannot be back-dated), period close (F-013), inventory purchase accounting (F-011).'], ['F-002', 'F-003']),
 '15': ('Complete transaction-level financial forensics', 'COMPLETE-WITH-LIMITS', ['44-transaction-level-reconciliation.csv', '20-financial-reconciliation-report.md'], ['44'], ['AR 1,800, cash-flow, revenue, VAT, KPI, inventory residuals each decomposed term by term (evidence/database/RUN-003/PH-09, PH-12).'], ['F-003', 'F-006', 'F-010', 'F-011']),
 '16': ('Negative, boundary and concurrency testing', 'COMPLETE-WITH-LIMITS', ['18-negative-boundary-test-results.md', '35-financial-scenario-test-specification.md'], ['36'], ['Receipt/supplier/stock/redrive/recurring/mark-paid concurrency; validation matrices for journals, accounts, HR, FX.', 'Limit: no k6 100-VU profile (see Phase 18).'], ['F-006', 'F-017', 'F-047', 'F-052']),
 '17': ('Security, RBAC and tenant isolation', 'COMPLETE-WITH-LIMITS', ['19-security-concurrency-recovery-results.md', '36-security-reliability-test-plan.md'], ['36'], ['Stub login, cookie tamper, cross-tenant read/write, mass assignment, role matrix, origin/rate limit, injection probes executed.', 'Observation: session cookie lacks Secure over HTTP (verify under TLS). SEC-013 staging-vs-production comparison not executable locally.'], ['F-020', 'F-033', 'F-031']),
 '18': ('Performance, reliability and recovery', 'PARTIAL', ['19-security-concurrency-recovery-results.md'], ['36'], ['Durable outbox across API restart (CFR-004), concurrent redrive, stock/order concurrency, report latency on small data.', 'NOT EXECUTED: 100-VU load, 50k-line report scale benchmark, Redis outage re-test.'], ['F-018']),
 '19': ('Historical regression revalidation', 'COMPLETE-WITH-LIMITS', ['21-existing-issue-regression-matrix.csv'], ['21'], ['HC-04, HC-05, FIN-016, FIN-020, CFR-004, FIN-019 (partial), CFR-006 (small data only), HC-10 re-run with evidence.'], ['F-049']),
 '20': ('Complete evidence and root-cause verification', 'COMPLETE-WITH-LIMITS', ['55-finding-evidence-and-root-cause-verification.csv', 'findings/'], ['22', '55'], ['54/54 findings have resolvable evidence; root cause CONFIRMED to file:line for most defects (see counts in 61); remaining ones carry candidate causes and next investigations.'], []),
 '21': ('Implementation-ready remediation specifications and tasks', 'COMPLETE', ['23-remediation-roadmap.md', '56-accounting-policy-decision-package.md', '57-future-remediation-task-register.csv', '58-remediation-specification-index.md', '59-post-remediation-regression-test-register.csv'], ['22'], ['findings/F-xxx-*/ with 10 files each; 57 task register and 59 regression register generated by tests/build_finding_packages.py; all tasks NOT STARTED.'], []),
 '22': ('Final audit closure, coverage and handoff', 'COMPLETE-WITH-LIMITS', ['61-final-exhaustive-audit-closure-report.md', '62-next-model-handoff.md', '60-evidence-manifest.csv'], ['45'], ['Closure report, manifest, handoff, development-environment restoration and fingerprint comparison.'], []),
}

def tasks_for(ph):
    rows = list(csv.DictReader(open('51-phased-audit-task-register.csv', encoding='utf-8')))
    return [r for r in rows if r['phase'] == ph]

def build():
    os.makedirs('phases', exist_ok=True)
    reg = {r['finding_id']: r for r in csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8'))}
    ex = list(csv.DictReader(open('34b-executed-test-results.csv', encoding='utf-8')))
    for ph, (title, status, reports, invs, results, fids) in P.items():
        d = os.path.join('phases', 'phase-' + ph)
        os.makedirs(d, exist_ok=True)
        scripts = sorted(os.path.basename(p) for p in glob.glob('tests/runtime/ph%s_*.mjs' % ph))
        evid = [p.replace('\\', '/') for p in glob.glob('evidence/*/RUN-00*/PH-%s' % ph)]
        tests = [r['exec_id'] for r in ex if re.search(r'AUD-%s|PH-?%s\b' % ('|'.join({'04': 'UI', '05': 'ONB|SET', '06': 'INV|SUP|PUR', '07': 'POS|REFUND', '08': 'STORE|KDS', '09': 'CINV|AR|CUST', '10': 'ATT|LEV|SAL|HRC|HR', '11': 'ACC', '12': 'MET|REV|CF', '13': 'SAAS', '14': 'MONTH', '16': 'REL', '17': 'SEC', '18': 'REL|CFR004|CFR006', '19': 'REG|HC|CFR|FIN019', '20': 'EVD'}.get(ph, 'NONE-%s' % ph)), ph), r['exec_id'])]
        fl = ['- [%s](../../findings/%s) — %s (%s)' % (f, [x for x in os.listdir('findings') if x.startswith(f + '-')][0] + '/finding.md', reg[f]['title'][:90], reg[f]['severity'].split(' ')[0]) for f in fids if f in reg and os.path.isdir('findings')]
        tk = tasks_for(ph)
        import collections as _c
        _cnt = _c.Counter(t['status'] for t in tk)
        status = ('DERIVED FROM REGISTER 51 (RUN-004): ' + (', '.join('%d %s' % (v, k) for k, v in sorted(_cnt.items())) or 'no tasks') + '; legacy label from RUN-003: ' + status)
        open(os.path.join(d, 'spec.md'), 'w', encoding='utf-8').write(f"""# Phase {ph} — {title}

**Closure status:** {status}

## Objective and scope
See the master prompt section for Phase {ph} (`writing-block.md`) — this file only references existing artefacts.

## Authoritative reports
{chr(10).join('- [%s](../../%s)' % (r, r) for r in reports)}

## Inventories
{', '.join('file ' + i for i in invs)} (see the audit root README for names).

## Existing test scripts
{chr(10).join('- `tests/runtime/%s`' % s for s in scripts) or '- none specific to this phase (shared helpers in `tests/runtime/chrome_lib.mjs`)'}

## Evidence
{chr(10).join('- `%s`' % e for e in evid) or '- see evidence/*/RUN-003 for related folders'}

## Existing test IDs (34b)
{', '.join(tests[:60]) or 'none mapped; see 34b and 51'}{' …' if len(tests) > 60 else ''}

## Actual results
{chr(10).join('- ' + r for r in results)}

## Findings
{chr(10).join(fl) or '- none specific'}

## Remaining blockers and unexecuted work
Recorded per task in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and in [48-blocked-vs-unexecuted-and-prerequisites.md](../../48-blocked-vs-unexecuted-and-prerequisites.md).

## Required new tests (if the phase is re-opened)
Any row of 51 with status PLANNED/BLOCKED for this phase; plus the regression tests in [59](../../59-post-remediation-regression-test-register.csv) after remediation.

## Required deliverables
Evidence under `evidence/{{screenshots,network,database}}/RUN-003/PH-{ph}/`, updated registers (22, 34b, 43, 44), progress.md entry.

## Completion criteria
Every executable task in 51 for this phase is EXECUTED-* or has a documented blocker; findings recorded in 22 with evidence; no product code changed.
""")
        open(os.path.join(d, 'tasks.md'), 'w', encoding='utf-8').write(f"""# Phase {ph} — tasks (execution instructions; NOT a progress tracker)

Status of each task is authoritative in [51-phased-audit-task-register.csv](../../51-phased-audit-task-register.csv) and [progress.md](../../progress.md). This file lists what the register holds for the phase at generation time.

{chr(10).join('- **%s** [%s] %s — module %s — expected: %s' % (t['task_id'], t['status'], t['objective'], t['module'], (t['expected'] or '')[:120]) for t in tk) or '- no phase-specific register rows; see spec.md for the work performed'}

## How to execute
1. Bring up the isolated stack (see 62 handoff). 2. Use real headed Chrome for UI steps (`tests/runtime/chrome_lib.mjs`). 3. Capture screenshot, network and SQL evidence under the RUN-003 folders. 4. Update 22/34b/43/44 and progress.md; never modify product code.
""")
    print('phases written', len(P))

if __name__ == '__main__':
    build()
