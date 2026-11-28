"""RUN-004: records load/benchmark/reconciliation-status results in 34b, 21 and the task overlay (idempotent). Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N = 'evidence/network/RUN-004/PH-18/'
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN004-WRITELOAD', 'API-DIAGNOSTIC load', 'Incremental write load 10/25/50/100 virtual users x 20 s, POS orders (untracked and tracked product) mixed with reads, tenant 4 fixture', 'owner', 'PASS',
     '5,719 orders created, 0 server errors or network errors at every stage; p95 184/420/634/2,454 ms; ~117-154 req/s; peak 16 PG connections; GET 429 rate limiting 12%-5%; tracked stock 12 -> 0 equals 12 accepted tracked orders; 0 duplicate order numbers; no orphan rows; evidence ' + N + 'AUD-R4-WRITE-LOAD.json. Client and server on the same host; not a capacity claim'),
    ('EX-RUN004-PAIDBURST', 'API-DIAGNOSTIC load', '600 concurrent payment confirmations (50 clients) then outbox drain', 'owner', 'PASS',
     'all 600 returned 200 in 2.0 s (p95 311 ms); outbox drained in 56 s (~10.6 events/s, ledger lags up to a minute for a 600-event burst); one sales journal per order, 13 extra COGS journals for tracked orders (not duplicates), revenue delta exactly 60,000.00, 0 unbalanced journals; evidence ' + N + 'AUD-R4-PAID-BURST.json'),
    ('EX-RUN004-BENCH-50K', 'API-DIAGNOSTIC + SYNTHETIC FIXTURE', '54,000 posted journal lines (27,000 journals via CSV import auto_post, tenant 10) report benchmark and correctness vs integer-cents oracle', 'owner', 'PASS',
     'P&L revenue 10,941,426.00 and expenses 2,703,789.00 equal the oracle; TB balanced; cash 1110 8,237,637.00; all report APIs 33-235 ms (ledger page 100 rows 209-235 ms; CSV export 337 KB 115 ms); P&L aggregation scans journal_entry_lines sequentially (planner choice at this size). Import itself: 1,500 journals took 58 s at 3k journals and 104 s at 24k and the client connection closed at ~60 s while the server completed (F-063). Browser rendering latency not measured here; evidence ' + N + 'AUD-R4-BENCH-MEASURE.json'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    w = csv.writer(h)
    for r in new:
        if r[0] not in ex:
            w.writerow(r)
rows = list(csv.reader(open('21-existing-issue-regression-matrix.csv', encoding='utf-8')))
hd = rows[0]; ix = {k: i for i, k in enumerate(hd)}
upd = {
    'CFR-006': ('RUN-004: FIXED at measured scale (54,000 posted lines: all reports 33-235 ms, correct to the cent)', N + 'AUD-R4-BENCH-MEASURE.json'),
    'CFR-001': ('RUN-004: NOT REPRODUCED; source review: auto_posting_engine.go contains no balancing/normalisation code (0 matches) and the deferred DB trigger rejects unbalanced journals', 'evidence/database/RUN-003/PH-19 (BND-01) + source grep'),
    'HC-06': ('RUN-004: NOT REPRODUCED (same evidence as CFR-001)', 'see CFR-001'),
    'HC-15': ('RUN-004: PARTIALLY FIXED - automatic recovery PASS and manual redrive executed (REL-004: 202 + 409, one audit row)', 'EX-RUN003-REL-004'),
}
for r in rows[1:]:
    if r[ix['issue_id']] in upd and 'RUN-004' not in r[ix['runtime_status_batch2']]:
        r[ix['runtime_status_batch2']], r[ix['runtime_evidence_batch2']] = upd[r[ix['issue_id']]]
csv.writer(open('21-existing-issue-regression-matrix.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
p = 'tests/r4_task_overlay.py'
s = open(p, encoding='utf-8').read()
add = "    'T-P18-002': ('EXECUTED-PASS', 'write load 10-100 VU (5,719 orders, 0 server errors), 600-payment burst with outbox drain, and 54,000-line report benchmark executed and reconciled; browser rendering latency for the 54k tenant not measured', ''),\n"
if "'T-P18-002'" not in s:
    s = s.replace("    'T-P20-001':", add + "    'T-P20-001':", 1)
open(p, 'w', encoding='utf-8').write(s)
print('ok')
