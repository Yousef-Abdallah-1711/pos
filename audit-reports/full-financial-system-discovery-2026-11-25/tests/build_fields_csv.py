import csv, re
rows = []
for line in open('tests/fields_src.txt', encoding='utf-8'):
    line = line.rstrip('\n')
    if not line.strip(): continue
    p = line.split('|')
    app, route, form, field, typ, val, maps = p[:7]
    v = val.lower()
    if re.search(r'\brequired\b|>=1 non-blank|zod min|server: >0|must be >0', v): req = 'yes'
    elif v.startswith('optional') or 'optional' in v: req = 'no'
    else: req = 'unspecified'
    rows.append([app, route, form, field, typ, req, val, maps, 'NOT TESTED'])
with open('07-field-level-inventory.csv', 'w', newline='', encoding='utf-8') as f:
    w = csv.writer(f)
    w.writerow(['app', 'page_route', 'form', 'field_name', 'type', 'required', 'validation', 'maps_to_api_field', 'tested_status'])
    w.writerows(rows)
print(len(rows))
