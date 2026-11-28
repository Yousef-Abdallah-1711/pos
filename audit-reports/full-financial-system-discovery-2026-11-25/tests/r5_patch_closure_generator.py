import os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
p = 'tests/build_closure_report.py'
s = open(p, encoding='utf-8').read()
old = "uires = ui.get('results', [])"
new = """uires = list(ui.get('results', []))
for _f in ('AUD-UI-MATRIX.run1-en-valid.json', 'AUD-UI-MATRIX.ar-mobile.json'):  # RUN-005: the AR/mobile rows live in these files (the base file has none), which made the RTL count read 0/59
    _p = 'evidence/network/RUN-003/PH-04/' + _f
    if os.path.exists(_p):
        uires += json.load(open(_p, encoding='utf-8')).get('results', [])"""
if 'RUN-005: the AR/mobile rows' not in s:
    assert old in s
    s = s.replace(old, new, 1)
open(p, 'w', encoding='utf-8').write(s)
print('patched')
