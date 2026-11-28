"""RUN-004: corrects the field denominator of 07-field-level-inventory.csv and records form-level coverage honestly. Writes 54b-field-denominator-and-form-coverage.csv/.md
and updates the tested_status column of 07 (never claims field-level verification). Audit tooling; idempotent."""
import csv, os, re, collections
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
rows = list(csv.DictReader(open('07-field-level-inventory.csv', encoding='utf-8')))
fields = list(rows[0].keys())
GEN = {'id', 'created_at', 'updated_at', 'tenant_id', 'restaurant_id', 'created_by', 'updated_by', 'deleted_at'}
# forms exercised by the RUN-004 real-UI matrix (keyword in form name or route) -> modes
EXER = {'customer': 'empty,invalid,arabic,long,valid', 'supplier': 'empty,invalid,arabic,long,valid', 'accountform': 'empty,invalid,arabic,long,valid', 'chart': 'empty,invalid,arabic,long,valid', 'journal': 'empty,invalid,arabic,long,valid',
        'employee': 'empty,invalid,arabic,long,valid', 'leave': 'empty,invalid,arabic,long,valid', 'attendance': 'empty,invalid,arabic,long,valid', 'contract': 'empty,invalid,arabic,long,valid', 'user': 'empty,invalid,arabic,long,valid',
        'role': 'empty,invalid,arabic,long,valid', 'driver': 'empty,invalid,arabic,long,valid', 'table': 'empty,invalid,arabic,long,valid', 'categor': 'empty,invalid,arabic,long,valid', 'product': 'empty,invalid,arabic,long,valid'}
EXACT = {'customerform', 'supplierform', 'accountform', 'journalentryform', 'employeeform', 'leaveform', 'attendanceform', 'contractform', 'categoryformpage', 'productform', 'productformpage', 'addtableform'}
seen = {}; cls = collections.Counter(); out = []
for i, r in enumerate(rows, start=2):
    key = (r['app'], r['page_route'], r['form'], r['field_name'])
    t = (r['type'] or '').lower(); n = r['field_name'].lower()
    if r['editable_kind'] == 'read-only/disabled':
        c = 'read-only/disabled (not a user input)'
    elif n in GEN:
        c = 'generated (not user input)'
    elif t == 'hidden':
        c = 'hidden (non-interactive)'
    elif key in seen:
        c = f'duplicate row of line {seen[key]}'
    else:
        c = 'logical user-facing field'; seen[key] = i
    base = os.path.splitext(os.path.basename(r['component_path'] or ''))[0].lower()
    fm = r['form'].lower().split(' (')[0].strip()
    modes = 'empty,invalid,arabic,long,valid' if (fm in EXACT or base in EXACT or os.path.splitext(os.path.basename(fm))[0] in EXACT) else ''
    if c != 'logical user-facing field':
        status = 'EXCLUDED FROM DENOMINATOR (' + c.split(' (')[0] + ')'
    elif r['app'] == 'dashboard' and modes:
        status = 'FORM-LEVEL EXERCISED RUN-004 (modes ' + modes + '); FIELD-LEVEL NOT VERIFIED'
    else:
        status = 'NOT TESTED'
    cls[c.split(' of ')[0].split(' (')[0]] += 1
    out.append([i, r['app'], r['page_route'], r['form'], r['field_name'], t, c, status])
    r['tested_status'] = status
with open('07-field-level-inventory.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.DictWriter(h, fieldnames=fields); w.writeheader(); w.writerows(rows)
with open('54b-field-denominator-and-form-coverage.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.writer(h); w.writerow(['inventory_line', 'app', 'page_route', 'form', 'field_name', 'type', 'classification', 'coverage_status']); w.writerows(out)
logical = [x for x in out if x[6] == 'logical user-facing field']
ex = [x for x in logical if x[7].startswith('FORM-LEVEL')]
forms = {(x[1], x[3]) for x in logical}; forms_ex = {(x[1], x[3]) for x in ex}
open('54b-field-denominator-and-form-coverage.md', 'w', encoding='utf-8').write(f"""# 54b - Field denominator and form coverage (RUN-004)

Source: [07-field-level-inventory.csv](07-field-level-inventory.csv) ({len(rows)} rows). Detail: [54b-field-denominator-and-form-coverage.csv](54b-field-denominator-and-form-coverage.csv). Built by `tests/build_field_denominator.py`.

## Corrected denominator
- Inventory rows: {len(rows)}
- Logical user-facing fields: **{len(logical)}** (dashboard {sum(1 for x in logical if x[1] == 'dashboard')}, storefront {sum(1 for x in logical if x[1] == 'restaurant-website')}, super-admin {sum(1 for x in logical if x[1] == 'super-admin')}, landing {sum(1 for x in logical if x[1] == 'landing')})
- Duplicate rows (same app, route, form and field name): {cls.get('duplicate row', 0)}
- Read-only or disabled: {cls.get('read-only/disabled', 0)}; generated: {cls.get('generated', 0)}; hidden: {cls.get('hidden', 0)}
- Distinct forms: {len(forms)}

## What was and was not executed
- RUN-004 real-UI form matrix (headed Chrome, empty/invalid/Arabic/long/valid modes) exercised **{len(forms_ex)} of {len(forms)} forms** covering **{len(ex)} of {len(logical)} logical fields at form level** (results in 34b EX-RUN004-FORMMATRIX and FORMMATRIX2, evidence in evidence/network/RUN-004/PH-06).
- **Field-level verification was NOT performed**: the driver fills every visible field of a form with a mode value and observes the form outcome; it does not isolate one field at a time, test min/max per numeric field, date ordering per pair, file types or sizes, or decimal precision. These rows are therefore labelled FORM-LEVEL EXERCISED, never PASS.
- {len(logical) - len(ex)} logical fields belong to forms not exercised (theme builder, settings sub-forms, storefront checkout, super-admin, POS, accounting invoices, reports filters and others) and stay NOT TESTED.
- Per-field boundary results that exist come from the API batteries summarised in 54-field-validation-and-boundary-results.md (RUN-003).

Forms exercised by the matrix that are ABSENT from the inventory (inventory incompleteness): driver create page (12 fields), settings users create dialog (4 fields), RBAC role create dialog (2 fields); the employee form exposes 6 fields in the UI against 1 inventoried. Add them before the next inventory refresh.

Findings from the matrix: F-074, F-075, F-079 (plus reconfirmation of F-014, F-031, F-062, F-064, F-065).
""")
print(len(rows), 'rows;', len(logical), 'logical;', len(ex), 'form-level exercised;', len(forms_ex), '/', len(forms), 'forms')
