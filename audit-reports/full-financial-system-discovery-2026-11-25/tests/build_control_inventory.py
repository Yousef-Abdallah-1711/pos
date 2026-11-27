"""Builds 53-ui-control-inventory.csv/.md from evidence/network/RUN-003/PH-04/AUD-UI-CONTROLS.json (rendered controls per dashboard route, EN desktop, headed Chrome)
and merges the language/mobile matrix (AUD-UI-MATRIX*.json). Audit tooling; idempotent. Counts are of VISIBLE rendered controls at load time (no clicking)."""
import csv, json, os
AUDIT = r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
os.chdir(AUDIT)
ctl = json.load(open('evidence/network/RUN-003/PH-04/AUD-UI-CONTROLS.json', encoding='utf-8'))['rows']
mat = {}
for f in ('AUD-UI-MATRIX.run1-en-valid.json', 'AUD-UI-MATRIX.ar-mobile.json'):
    p = 'evidence/network/RUN-003/PH-04/' + f
    if os.path.exists(p):
        for r in json.load(open(p, encoding='utf-8'))['results']:
            mat[(r['mode'], r['route'])] = r
cols = ['route', 'buttons', 'links', 'inputs', 'input_types', 'selects_comboboxes', 'textareas', 'tabs', 'tables', 'forms', 'checkboxes', 'switches', 'dialog_triggers', 'sample_button_labels', 'en_loaded', 'ar_dir', 'ar_len', 'mobile_overflow', 'mobile_len', 'functional_test_status']
out = []
def status(route, c):
    if c.get('error') or c.get('len', 0) < 200: return 'NOT-RENDERED (page empty or environment error) - re-run'
    if route in ('pos', 'products/new', 'accounting/customer-invoices', 'kds', 'reports/financial/profit-loss', 'reports/financial/cash-flow', 'reports/financial/balance-sheet', 'hr/payroll', 'settings/restaurant', 'pos/settings', 'kds/settings'): return 'INTERACTED (see 34b / PH-05..PH-12)'
    return 'LOADED + control counts only (controls not individually exercised)'
for c in ctl:
    r = c['route']; en = mat.get(('en', r), {}); ar = mat.get(('ar', r), {}); mb = mat.get(('en-mobile', r), {})
    out.append([r, c.get('buttons'), c.get('links'), c.get('inputs'), c.get('inputTypes'), c.get('selects'), c.get('textareas'), c.get('tabs'), c.get('tables'), c.get('forms'), c.get('checkboxes'), c.get('switches'), c.get('dialogTriggers'), ' | '.join(c.get('buttonLabels', [])[:12]), 'yes' if en.get('len', 0) > 200 else ('n/a' if not en else 'short'), ar.get('dir', 'n/a'), ar.get('len', 'n/a'), mb.get('overflow', 'n/a'), mb.get('len', 'n/a'), status(r, c)])
with open('53-ui-control-inventory.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.writer(h); w.writerow(cols); w.writerows(out)
tot = lambda i: sum(int(r[i] or 0) for r in out if str(r[i]).isdigit() or isinstance(r[i], int))
open('53-ui-control-inventory.md', 'w', encoding='utf-8').write(f"""# 53 — UI control inventory (RUN-003, headed Chrome, EN desktop; read-only)

Data: [53-ui-control-inventory.csv](53-ui-control-inventory.csv) from `tests/runtime/ph04_controls.mjs` (evidence `evidence/network/RUN-003/PH-04/AUD-UI-CONTROLS.json`) merged with the language/viewport matrix `tests/runtime/ph04_ui_matrix.mjs`.

**What this is:** the count and sample labels of VISIBLE rendered controls per dashboard route at load time ({len(out)} routes). **What it is not:** a per-control functional test. Controls were not clicked; dialogs/tabs that appear only after interaction are not enumerated. Pages needing data (details, edit forms) and Super Admin / landing / storefront apps are outside this table (their routes were loaded in earlier batches and in the workflow scripts).

Totals over the {len(out)} routes: buttons {tot(1)}, links {tot(2)}, inputs {tot(3)}, tables {tot(8)}, forms {tot(9)}.

Status per row (last column): INTERACTED = a workflow script exercised the page (see 34b); LOADED = control counts only; NOT-RENDERED = empty page (environment issue) and must be re-run.
""")
print(len(out), 'routes')
