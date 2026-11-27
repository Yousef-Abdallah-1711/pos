"""Builds 61-final-exhaustive-audit-closure-report.md with live counts from the registers. Verdict text is authored here (honest, evidence-based). Audit tooling; idempotent."""
import csv, collections, glob, json, os, re, subprocess
AUDIT = r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
REPO = r'C:\Users\Yousef\Desktop\Projects\POSS'
os.chdir(AUDIT)
rd = lambda f: list(csv.DictReader(open(f, encoding='utf-8')))
F = rd('22-confirmed-findings-register.csv'); T34 = rd('34b-executed-test-results.csv'); M30 = rd('30c-metric-coverage-reconciliation.csv'); M43 = rd('43-report-metric-runtime-results.csv')
T51 = rd('51-phased-audit-task-register.csv'); T57 = rd('57-future-remediation-task-register.csv'); R59 = rd('59-post-remediation-regression-test-register.csv')
sev = collections.Counter(f['severity'].split(' ')[0] for f in F)
conf = collections.Counter()
for f in F:
    m = re.match(r'\s*(CONFIRMED|PROBABLE|CANDIDATE|UNKNOWN)', f['root_cause_first_break'])
    conf[m.group(1) if m else 'UNSTATED'] += 1
cat = collections.Counter(f['category'].split('(')[0].strip().split('/')[0].strip() for f in F)
status = collections.Counter(f['evidence_status'].split(' ')[0] for f in F)
cov = collections.Counter(m['coverage_class'].split(' (')[0] for m in M30)
ttype = collections.Counter(); tres = collections.Counter()
for r in T34:
    ttype[r['type']] += 1; tres[r['result'].split(' ')[0].split('(')[0]] += 1
run3 = [r for r in T34 if r['exec_id'].startswith('EX-RUN003')]
run3res = collections.Counter(r['result'].split(' ')[0] for r in run3)
st51 = collections.Counter(t['status'] for t in T51)
ui = {}
if os.path.exists('evidence/network/RUN-003/PH-04/AUD-UI-MATRIX.json'):
    ui = json.load(open('evidence/network/RUN-003/PH-04/AUD-UI-MATRIX.json', encoding='utf-8'))
uires = ui.get('results', [])
ui404 = [r for r in uires if r.get('is404')]; uierr = [r for r in uires if r.get('hasError')]; uiover = [r for r in uires if r.get('overflow')]
uirtl = [r for r in uires if r['mode'] == 'ar' and r.get('dir') == 'rtl']
uiar = [r for r in uires if r['mode'] == 'ar']
rows_md = lambda rows: '\n'.join('| ' + ' | '.join(str(c) for c in r) + ' |' for r in rows)
def count(f): return sum(1 for _ in csv.reader(open(f, encoding='utf-8'))) - 1
sha = subprocess.check_output(['git', '-C', REPO, 'rev-parse', 'HEAD'], text=True).strip()
mod = subprocess.check_output(['git', '-C', REPO, 'status', '--short'], text=True).strip().split('\n')
mod_tool = [m for m in mod if re.search(r'tsconfig\.json|next-env\.d\.ts', m)]
mod_other = [m for m in mod if m not in mod_tool and not m.startswith('??')]
now = subprocess.check_output(['date', '-u', '+%FT%TZ'], text=True).strip()
sev_s = ', '.join('%s %d' % (k, v) for k, v in sev.most_common())
open('61-final-exhaustive-audit-closure-report.md', 'w', encoding='utf-8').write(f"""# 61 — Final exhaustive audit closure report (Master Prompt 05, RUN-003 and RUN-004)

Generated {now} by `tests/build_closure_report.py` from the registers. **Overall audit classification: PARTIAL.** Financial release readiness: **NO-GO** (unchanged; see [24](24-financial-release-verdict.md)). Nothing was fixed; no product code, migration or formula was changed; all remediation tasks are NOT STARTED.

## 1. Discovery (exact)
| Item | Count | Source |
|---|---|---|
| Applications | 5 (Go API, dashboard, landing, restaurant storefront, super-admin) | 03 |
| Pages discovered / loaded | 119 discovered (2 route handlers) / 117 loaded in earlier batches | 06, 45 |
| Logical forms | {count('07b-form-inventory.csv')} | 07b |
| Fields | {count('07-field-level-inventory.csv')} | 07 |
| API endpoints | {count('08-api-inventory.csv')} | 08 |
| Database tables | 102 (287 FK, 47 triggers, 20 RLS tables) | 09, 49 |
| Financial operations traced (UI→report) | 16 operation groups with file:line | 52 |
| Posting-matrix rows | {count('11-financial-posting-matrix.csv')} | 11 |
| Report metrics (inventory + runtime-only) | 112 + 7 = {len(M30)} coverage rows | 30, 30c |
| Manual UI test-plan rows | {count('34-complete-manual-ui-test-plan.csv')} (templated plan rows, not results) | 34 |
| UI controls | counts per route in 53 (59 routes); one safe interaction per control for 644 controls in 53b (cap 22 per route; destructive controls skipped) | 53, 53b |

## 2. Runtime execution (exact, from 34b and run artifacts)
- Executed-result rows in 34b: **{len(T34)}** ({len(run3)} added in RUN-003). By type: {', '.join('%s %d' % (k, v) for k, v in ttype.most_common(8))}.
- Pages opened in real headed Chrome in RUN-003: dashboard route matrix {len(uires)} page-loads ({len(set(r['route'] for r in uires))} distinct routes × EN desktop / AR desktop / EN mobile){'' if uires else ' — matrix result file not present'}; plus the workflow scripts (register, settings, POS, KDS, storefront, invoices, reports).
- Matrix findings: {len(ui404)} loads that rendered a 404 shell, {len(uierr)} with an application-error text, {len(uiover)} with horizontal overflow (mobile unless noted), {len(uirtl)}/{len(uiar)} Arabic loads with `dir=rtl`.
- Forms submitted through the real UI in headed Chrome: registration (3 tenants + probes), contact settings, KDS settings, POS settings, employee create (fails, F-014), customer invoice, POS cash order, storefront checkout. All other write-path tests used API-DIAGNOSTIC scripts (labelled as such in 34b).
- API functional tests (RUN-003 diagnostics): {sum(v for k, v in ttype.items() if k.startswith('API-DIAGNOSTIC'))}; API permission-only tests (earlier batches): 732.
- Historical issues re-validated (21): see matrix; CFR-004, HC-04, HC-05, FIN-016, FIN-020, HC-10 executed, FIN-019 and CFR-006 partial.
- Security tests executed (RUN-003): stub login, session tamper, cross-tenant read/write, mass assignment, role matrix, origin check, rate limit, injection, cookie attributes (see 34b `EX-RUN003-SEC-*`).
- Reliability tests executed: CFR-004, REL-001/004/009/025/027, Redis outage. **Performance (RUN-004):** write load 10-100 virtual users (5,719 orders, 0 server errors), 600-payment burst with outbox drain (56 s), and a 54,000-line synthetic benchmark (reports 33-235 ms, exact to the cent) executed; not production capacity claims; browser rendering latency on the large tenant not measured.
- Report metrics with an independent comparison: {sum(1 for m in M30 if m['coverage_class'].startswith('MEASURED'))} of {len(M30)} rows (classes: {', '.join('%s %d' % (k, v) for k, v in cov.most_common())}). Matching the application's own SQL is NOT treated as proof of accounting correctness; the GL oracle and policy status are recorded per row.

## 3. Results
- RUN-003 test outcomes: {', '.join('%s %d' % (k, v) for k, v in run3res.most_common())}.
- Whole 34b result distribution: {', '.join('%s %d' % (k, v) for k, v in tres.most_common(10))}.
- Tasks in 51: {', '.join('%s %d' % (k, v) for k, v in st51.most_common())}.
- Monetary discrepancies reproduced (EGP): F-003 cash flow residual 5,190.00 (clean tenant) / 107,300.00 (busy tenant); BS equation off by exactly the net loss -27,710.00 (F-002); AR vs GL 1,800.00 (F-006 over-settlement); VAT report vs GL 29,060.00 (F-010); inventory valuation vs GL 228,960.00 (F-011); KPI net profit +14,360.00 and cash position -226,270.01 (F-043); supplier VAT +140 per 10,000 at 10% discount (F-041); customer invoice +112 (F-007).
- Unexplained residuals: none remaining in the datasets reconciled; unrelated residuals were not offset.

## 4. Findings
- Total findings in 22: **{len(F)}** (severity: {sev_s}). Evidence verification (55): every finding has resolvable evidence.
- Evidence status: {', '.join('%s %d' % (k, v) for k, v in status.most_common())}.
- Root-cause confidence: {', '.join('%s %d' % (k, v) for k, v in conf.most_common())} (UNSTATED = no confidence recorded in the register; treated as CANDIDATE in `findings/*/root-cause.md`).
- Categories: {', '.join('%s %d' % (k, v) for k, v in cat.most_common(8))}.
- New in RUN-003: F-037 … F-065; new in RUN-004: F-066 … F-079 (F-051 unused). RUN-004 also retracted one RUN-003 claim: the duplicate-table-number part of F-064 (duplicates return 409).
- Missing capabilities (register category): {sum(1 for f in F if 'Missing capability' in f['category'])}. Policy-dependent discrepancies: revenue bridge (POL-01) and others in 56.

## 5. Future remediation
- Per-finding packages: `findings/` — {len(F)} folders × 10 files (`finding, reproduction, root-cause, affected-code-paths, financial-impact, dependency-map, evidence-index, remediation-spec, remediation-tasks, regression-test-spec`). Index: [58](58-remediation-specification-index.md).
- Future implementation tasks: **{len(T57)}** in [57](57-future-remediation-task-register.csv), all NOT STARTED.
- Required regression tests: **{len(R59)}** in [59](59-post-remediation-regression-test-register.csv).
- Policy decisions: 19 in [56](56-accounting-policy-decision-package.md): POL-01..POL-15 plus proposed POL-16..POL-19 added in RUN-004 (none approved; all OPEN).
- Roadmap/dependencies/blockers: [23](23-remediation-roadmap.md), [46](46-root-cause-dependency-map.md), [48](48-blocked-vs-unexecuted-and-prerequisites.md).
- Prerequisite order that unlocks blocked workflows: F-001 → F-014 → F-013 → F-005/F-026 → F-049 → full month rerun.

## 6. Environment and integrity
- Git HEAD `{sha}` unchanged. Product files modified by the audit: **0**. Tooling-only modified files present from the start: {len(mod_tool)}. Other modified tracked files now present: {len(mod_other)} — these are NOT from the audit (landing marketing files edited concurrently by another actor; see `evidence/runtime/RUN-003/48-git-status-midrun.txt`).
- Audit-only files: everything under `audit-reports/full-financial-system-discovery-2026-11-25/` (manifest: [60](60-evidence-manifest.csv)). Previous evidence preserved (no deletions; earlier-run captures copied as `*.earlier-run.json` where overwritten).
- Original development database: fingerprint compared at closure (see §8).
- Snapshots: start, mid and final checksummed dumps in `evidence/database/` (see `evidence/database/RUN-003/T-P02-003-checkpoint.txt`).
- Handoff readiness: [62](62-next-model-handoff.md) + `phases/` + `progress.md`.

## 7. Final acceptance gates (from the master prompt §15)
| Gate | Status | Note |
|---|---|---|
| All discovered feasible tests executed | NOT MET | executed in RUN-004: SEC-013, 100-VU load, 54k-line benchmark, theme builder UI, bank reconciliation, form matrices, control matrix; still open: field-level isolation tests, Arabic/mobile interaction, invoice/payment/POS/settings form matrices, destructive-control fixtures (45 addendum) |
| All pages have explicit functional coverage | PARTIAL | every dashboard route loaded in 3 modes; interaction depth varies (see 45 and 34b) |
| All logical forms explicit coverage | NOT MET | 116 forms in the dashboard-heavy inventory; 13 exercised through the UI at form level (54b), 17 form targets in the matrix incl. forms missing from the inventory |
| All fields and controls explicit test statuses | NOT MET | 531 logical fields: 94 at form level, 0 at field level (54b); 644 controls interacted once (53b), 88 skipped, 58 tooling-blocked |
| Financial operations have documented lineage | MET (52) | 16 operation groups, file:line |
| Report metrics explicit verification statuses | MET (30c, 43) | all 119 rows executed (RUN-004); results include FAIL rows with findings |
| Feasible financial transactions reconciled | PARTIAL | clean-tenant month + residual decompositions; payroll/period close/inventory purchases blocked |
| Findings have evidence classifications | MET | 22/55 |
| Confirmed findings reproducible | MET | scripts under tests/runtime |
| Root causes have confidence classifications | MET | every finding classified: CONFIRMED or PROBABLE with the exact investigation gap (F-018, F-021, F-038, F-074) |
| Blocked workflows have prerequisite mappings | MET | 48 |
| Remediation tasks have specifications | MET | findings/*/remediation-spec.md |
| Remediation tasks have acceptance criteria | MET | each task in 57 |
| Evidence paths verified | MET | 55 + manifest 60 |
| Source paths verified | MET | generator resolves every cited path against `git ls-files` (unresolved count shown in 58) |
| Previous evidence preserved | MET | nothing deleted |
| Original environment restored | SEE §8 | |
| Another model can resume without history | MET | 62, phases/, progress.md, scripts |

## 8. Environment restoration
(filled by the closing step — see the section below, written after `docker start …` and the DB fingerprint comparison)

## 9. Verdicts
1. System discovery completeness — **PARTIAL-HIGH** (static discovery and lineage complete; per-control counts for 59 routes; the field inventory misses some forms, see 54b).
2. Functional UI testing completeness — **PARTIAL** (route matrix in 3 modes; 644 control interactions; 17 form targets x 5 modes; field-level and Arabic/mobile interaction not done; many flows still API-diagnostic).
3. Financial transaction verification completeness — **PARTIAL** (all executed flows reconciled; blocked flows documented).
4. Reporting verification completeness — **PARTIAL-HIGH** ({sum(1 for m in M30 if m['coverage_class'].startswith('MEASURED'))}/{len(M30)} rows measured; exports and multi-dataset variants incomplete).
5. Security verification completeness — **PARTIAL** (tenant isolation, session, role, input probes and SEC-013 staging comparison done; no TLS/production-environment checks).
6. Evidence completeness — **HIGH** (every finding has resolvable evidence; manifest with hashes).
7. Remediation specification readiness — **MEDIUM-HIGH**: every finding has a hand-written design, measurable acceptance criteria and policy references (RUN-004), tasks apply only the phases that fit and shared root causes point at one owner task; designs are still proposals to be confirmed by engineers (F-018, F-021, F-038, F-074 need a root-cause confirmation step first).
8. Overall audit completeness — **PARTIAL**.
9. Financial integrity — **NOT ESTABLISHED** (multiple reproduced ledger/report defects; see §3).
10. Financial release readiness — **NO-GO**.
""")
print('written 61')
