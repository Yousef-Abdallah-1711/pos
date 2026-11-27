"""RUN-005 Phase E: builds 54c-isolated-field-results.csv/.md from evidence/network/RUN-005/PH-06/AUD-R5-FIELD-*.json and marks the isolated coverage in 07 (tested_status) and 54b. Idempotent. Audit tooling only."""
import csv, json, glob, os, re, collections
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
files = [f for f in glob.glob('evidence/network/RUN-005/PH-06/AUD-R5-FIELD-*.json') if 'partial' not in f]
rows = []
per_form = collections.OrderedDict()
for f in sorted(files):
    j = json.load(open(f, encoding='utf-8'))
    form = j['form']
    for r in j['results']:
        m = re.match(r'REJECTED-BY-SERVER\((2\d\d)\)', r['outcome'])
        if m:  # 2xx response but the row-count probe saw no new row: classify by the response, not as a rejection
            r['outcome'] = 'ACCEPTED-2xx(' + m.group(1) + ', persisted row not verified by the count probe)'
        stored = ''
        if r.get('persisted'):
            stored = r['persisted'][:160]
        rows.append([form, r['field_index'], r['field'], r.get('placeholder', ''), r['kind'], r['case'], r['value_sent'], r['outcome'], '; '.join(r.get('request') or []), (r.get('response') or '')[:110], '; '.join(r.get('messages') or [])[:110], stored, os.path.basename(f)])
        d = per_form.setdefault(form, collections.Counter())
        d[r['outcome'].split('(')[0]] += 1
# baseline validity: a form with no accepted case at all gives no isolated verdict
inconclusive = {fm for fm, c in per_form.items() if not any(k.startswith('ACCEPTED') for k in c)}
KNOWN = {'categories': 'F-065 (display_order is sent as a string, every submit returns 400 Invalid request body)'}
for r in rows:
    if r[0] in inconclusive:
        r[7] = r[7] + ' [INCONCLUSIVE: no accepted submit for this form' + (' - ' + KNOWN[r[0]] if r[0] in KNOWN else '') + ']'
cols = ['form', 'field_index', 'field', 'placeholder', 'inferred_kind', 'case', 'value_sent', 'outcome', 'request', 'response_head', 'ui_messages', 'persisted_row_head', 'evidence_file']
with open('54c-isolated-field-results.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.writer(h); w.writerow(cols); w.writerows(rows)
fields = {(r[0], r[1]) for r in rows}
# notable validation gaps derived mechanically
gap = collections.Counter()
for r in rows:
    form, idx, fld, ph, kind, case, val, out = r[:8]
    if kind == 'phone' and case in ('letters', 'arabic-digits') and out.startswith('ACCEPTED'):
        gap['phone accepts letters or Arabic-Indic digits'] += 1
    if kind == 'currency' and case == 'lowercase' and out.startswith('ACCEPTED'):
        gap['currency code accepted in lowercase'] += 1
    if kind == 'email' and case == 'uppercase' and out.startswith('ACCEPTED'):
        gap['email stored without case normalisation'] += 1
    if case in ('long300', 'long40', 'long260', 'long30') and out.startswith('REJECTED-BY-SERVER'):
        gap['over-length value rejected only by the database (raw pq text returned)'] += 1
    if case in ('long300', 'long40', 'long260', 'long30') and out.startswith('REJECTED-BY-SERVER(500'):
        gap['over-length value returns HTTP 500'] += 1
    if kind == 'number' and case == 'huge' and out.startswith('ACCEPTED'):
        gap['very large number (99,999,999,999) accepted'] += 1
    if case == 'special' and out.startswith('ACCEPTED'):
        gap['HTML/SQL characters stored verbatim (data); escaping is a rendering concern, verified PASS in RUN-003'] += 1
    if out == 'NO-REQUEST-NO-MESSAGE':
        gap['submit produced no request and no visible message'] += 1
    if out == 'TOOLING-ERROR':
        gap['tooling error (not a product verdict)'] += 1
tot = collections.Counter(r[7].split('(')[0] for r in rows)
open('54c-isolated-field-results.md', 'w', encoding='utf-8').write(f"""# 54c - Isolated field-level results (RUN-005, headed Chrome)

Data: [54c-isolated-field-results.csv](54c-isolated-field-results.csv) (one row per field and case). Runner: `tests/runtime/ph06_r5_field_isolation.mjs`. Evidence JSON: `evidence/network/RUN-005/PH-06/AUD-R5-FIELD-<form>.json`.

**Method:** for each discovered field of a form, all other fields receive a valid baseline value, ONLY the target field is changed to a case value chosen by its inferred kind (text, email, phone, number, date, code, currency, password), the real UI form is submitted, and the actual request, response, visible messages and the persisted row are recorded. Applicability: text-like fields get empty, whitespace, 300 characters, Arabic and special characters; email, phone, number, date, code and currency fields get kind-specific cases. Select, checkbox, radio, file and rich controls are NOT covered by this runner.

**Coverage reached so far:** {len(fields)} fields in {len(per_form)} forms ({', '.join(per_form)}), {len(rows)} executed cases. Outcome totals: {dict(tot)}. Per form: {{{', '.join(f'{k}: {dict(v)}' for k, v in per_form.items())}}}.

**Derived observations (mechanical counts over the rows):** {dict(gap)}.

**Inconclusive forms (no case was accepted, so a field-level verdict is not possible):** {sorted(inconclusive) or 'none'}; known cause: {KNOWN}.

**Honest limits:** fields of forms not listed above remain NOT TESTED at field level (see 54b); numeric min/max rules were probed with fixed values (negative, zero, 6 decimals, 99,999,999,999, 0.001), not with each column's exact bounds; inferred kinds come from labels, placeholders and input types and can be wrong for unusual fields; file upload, date-pair ordering and cross-field rules are not covered; financial side effects were verified for forms that create ledger-relevant rows only where noted in 34b.
""")
# mark isolated coverage in 07 (form-level text kept; never claims PASS)
iso_forms = set(per_form)
MAP = {'customers': 'customerform', 'suppliers': 'supplierform', 'accounts': 'accountform', 'products': 'productform', 'categories': 'categoryformpage', 'drivers': 'driver', 'tables': 'addtableform', 'contracts': 'contractform', 'attendance': 'attendanceform', 'roles': 'role', 'users': 'user'}
r07 = list(csv.DictReader(open('07-field-level-inventory.csv', encoding='utf-8')))
hd = list(r07[0].keys()); n = 0
for r in r07:
    base = os.path.splitext(os.path.basename(r['component_path'] or ''))[0].lower(); fm = r['form'].lower().split(' (')[0].strip()
    for k, v in MAP.items():
        if k in iso_forms and (fm == v or base == v or (k == 'products' and base == 'productformpage')):
            if k in inconclusive:
                r['tested_status'] = 'ISOLATED FIELD RUN INCONCLUSIVE RUN-005 (no submit accepted' + (', ' + KNOWN[k] if k in KNOWN else '') + '); see 54c'
                n += 1
            elif not r['tested_status'].startswith('ISOLATED'):
                r['tested_status'] = 'ISOLATED FIELD CASES EXECUTED RUN-005 for the form (see 54c); per-field rows in 54c; selects/files not covered'
                n += 1
with open('07-field-level-inventory.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.DictWriter(h, fieldnames=hd); w.writeheader(); w.writerows(r07)
print(len(fields), 'fields,', len(rows), 'cases;', n, '07 rows marked;', dict(tot))
