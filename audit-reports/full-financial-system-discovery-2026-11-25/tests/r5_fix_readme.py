"""RUN-005: corrects stale counts in README.md and appends the RUN-005 index (idempotent). Historical text is kept; only wrong counts are corrected."""
import csv, os, re
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
n_find = len(list(csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8'))))
s = open('README.md', encoding='utf-8').read()
s = s.replace('[22 (36 findings)]', f'[22 ({n_find} findings)]')
s = s.replace('| 22–24 | [22](22-confirmed-findings-register.csv)', '| 22–24 | [22](22-confirmed-findings-register.csv)')
M = '## RUN-005 index (2026-11-27)'
add = f"""
{M}
Current authoritative counts are in [61](61-final-exhaustive-audit-closure-report.md) (generated) and [progress.md](progress.md); the numbers quoted in the historical table above (for example 36 findings, 128 forms, 627 fields) belong to earlier runs. Findings now: {n_find}.
New or extended in RUN-005: [64 baseline and continuity](64-run005-baseline.md), [06b route and API reconciliation](06b-route-and-api-reconciliation.md), [30d metric verification classes](30d-metric-verification-classes.md), [53b control interaction results](53b-ui-control-interaction-results.md), [54b field denominator and form coverage](54b-field-denominator-and-form-coverage.md), [54c isolated field-level results](54c-isolated-field-results.md), [56b policy decision detail](56b-policy-decision-detail.md), [58b package validation](58b-package-validation.md), [63 RUN-004 baseline](63-resume-baseline-run004.md). Evidence: `evidence/*/RUN-005/`. Handoff: [62](62-next-model-handoff.md).
"""
if M not in s:
    s = s.rstrip('\n') + '\n' + add
open('README.md', 'w', encoding='utf-8').write(s)
print('ok', n_find)
