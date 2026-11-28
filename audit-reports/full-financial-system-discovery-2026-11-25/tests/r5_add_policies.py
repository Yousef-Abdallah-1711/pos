import os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
p = '56-accounting-policy-decision-package.md'
s = open(p, encoding='utf-8').read()
rows = [
    "| POL-20 | Reporting time zone and day boundary (PROPOSED in RUN-005) | UTC day (current) vs tenant time zone (Africa/Cairo default) for all date-bounded reports, KDS and exports | F-070: KDS hours are UTC buckets (Cairo buckets differ); F-042 delivered_today and report cut-offs use CURRENT_DATE | F-070, F-042 | tenant time zone setting, default Africa/Cairo; store timestamps with time zone |",
    "| POL-21 | Bank reconciliation completion rule (PROPOSED in RUN-005) | zero tolerance vs 0.01 vs approved adjustment journal; who may complete | F-072: session completed with an invented 777,777.77 line; F-073 cross-tenant journal reference | F-072, F-073, F-029 | book balance from the GL, completion only at zero difference (or with a recorded adjustment journal), finance permission to complete |",
    "| POL-22 | Subscription trial extension limits (PROPOSED in RUN-005) | unbounded (current) vs capped with reason and history | F-055: 100,000 days accepted (end year 2300) | F-055 | cap 90 days, reason required, history row |",
]
if 'POL-20' not in s:
    i = s.index('| POL-19 |'); j = s.index('\n', i)
    s = s[:j + 1] + '\n'.join(rows) + '\n' + s[j + 1:]
    open(p, 'w', encoding='utf-8').write(s)
print('POL-22' in s)
