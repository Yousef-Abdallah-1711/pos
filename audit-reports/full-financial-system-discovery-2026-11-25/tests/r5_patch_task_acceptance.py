import os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
p = 'tests/build_finding_packages.py'
s = open(p, encoding='utf-8').read()
marker = '# RUN-005 tailored acceptance'
if marker not in s:
    old = "            n[0] += 1\n            tid = 'REM-%s-%02d' % (fid, n[0])\n"
    new = old + """            # RUN-005 tailored acceptance: replace generic phase criteria with finding-specific, measurable ones
            _exp = (f['expected'] or '')[:150]
            if phase.startswith('P4'):
                acc = (ov.get('acceptance') or [None])[0] or ('UI, API and stored values agree for the reproduction dataset: ' + _exp)
            elif phase.startswith('P5'):
                acc = 'Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (' + _exp + ') and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged'
            elif phase.startswith('P6'):
                acc = 'Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for ' + fid + '; no cross-tenant data and no raw database text in responses'
            elif phase.startswith('P7'):
                acc = 'Read-only exception report lists the affected historical records (' + (hist or '')[:110].replace('\\n', ' ') + '); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals'
            elif phase.startswith('P8'):
                acc = 'Every REG-' + fid + '-* test in regression-test-spec.md passes on the isolated stack, including the expected result: ' + _exp + '; evidence stored under the release evidence folder'
            elif phase.startswith('P9'):
                acc = 'On staging the reproduction steps of ' + fid + ' no longer reproduce (expected: ' + _exp + '); the rollback in remediation-spec section 19 is rehearsed and its duration recorded'
"""
    assert old in s
    s = s.replace(old, new, 1)
    open(p, 'w', encoding='utf-8').write(s)
print('patched')
