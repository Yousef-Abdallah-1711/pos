"""RUN-005: F-082 (P&L report overflows at 768 px) + 34b rows for the Arabic/mobile/tablet pass (idempotent)."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N = 'evidence/network/RUN-005/PH-04/AUD-R5-LANG-VIEWPORT.json'
S = 'evidence/screenshots/RUN-005/PH-04/AUD-R5-LANG-en-768-reports_financial_profit-loss.png'
F = {'F-082': dict(
    title='The Profit and Loss report overflows horizontally at tablet width (768 px): content is 803 px wide and the page scrolls sideways',
    category='UX defect (responsive layout)', severity='Low', evidence_status='REPRODUCED',
    affected_modules='Reports > Financial > Profit and Loss',
    reproduction='Headed Chrome, viewport 768x1024, English, tenant 4 owner: open /en/dashboard/reports/financial/profit-loss and compare document scrollWidth with the viewport width',
    expected='No horizontal page scroll; wide tables scroll inside their own container',
    actual='document.documentElement.scrollWidth is 803 against a 768 px viewport; the same report fits at 1024 px, 390 px and 360 px and in Arabic at 1440 px',
    monetary_impact_EGP='0', root_cause_first_break='PROBABLE (measured layout fact; the element that sets the 803 px width was not isolated): a fixed-width element or table in the report page without an overflow wrapper at the md breakpoint',
    source_refs='frontend/apps/dashboard/src/app/[locale]/dashboard/reports/financial/profit-loss/page.tsx', related_static_ids='',
    evidence_paths=N + '; ' + S + '; tests/runtime/ph04_r5_lang_viewport.mjs',
    remediation_direction='Wrap the wide table in an overflow-x container or let the layout reflow at the md breakpoint',
    regression_tests_needed='At 768, 1024 and 1440 px the page has no horizontal scroll (scrollWidth <= innerWidth)')}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        rows.append([dict(v, finding_id=fid).get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [('EX-RUN005-ARABIC-RTL', 'UI-HEADED-CHROME', 'Arabic functional pass at 1440 px and 390 px on 15 representative routes: html direction, dialog direction and fit, submit reachability, overflow, Arabic search round trip (no data saved)', 'owner', 'PASS',
        'html dir=rtl on 30 of 30 Arabic checks (the earlier "0/59" in report 61 was a generator bug, fixed); no overflow, dialogs fit, Arabic text retained; evidence ' + N),
       ('EX-RUN005-MOBILE-TABLET', 'UI-HEADED-CHROME', 'Mobile and tablet pass at 360, 390, 768 and 1024 px on 15 routes: menu button, overflow, off-screen controls, dialog fit and submit reachability, touch targets', 'owner', 'PARTIAL',
        'mobile menu button present on all 60 mobile/tablet checks; 16 dialogs measured, none clipped, 0 unreachable submit buttons; horizontal overflow only on P&L at 768 px (F-082); orders and chart of accounts have more than 12 touch targets under 32 px at 360 and 390 px (OBSERVATION OBS-A11Y-001, not registered as a finding); 15 routes only, not the full route set; evidence ' + N)]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    for r in new:
        if r[0] not in ex:
            csv.writer(h).writerow(r)
print(len(rows) - 1)
