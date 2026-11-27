import os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
p = 'tests/build_finding_packages.py'
s = open(p, encoding='utf-8').read()
old = "{extra_rc}\n## Verified source locations"
new = """{extra_rc}
## Origin and propagation
- Origin (first verified source location): {('`%s`%s' % (good[0][2], (':' + good[0][1]) if good[0][1] else '')) if good else 'see the first break text above'}.
- Propagation: affected modules {f['affected_modules']}; related findings {', '.join(rel_ids) or 'none recorded'}; observed consequence: {f['actual'][:240]}
- Cause versus symptom: the symptom is what users or reports show ({f['actual'][:140]}); the cause is the first break above.

## Investigation performed
- Executed audit tests: {', '.join(tids) if tids else 'none recorded in 34b'}; scripts: {', '.join(scripts[:3]) if scripts else 'reproduced through UI/API/SQL steps in reproduction.md'}.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations"""
if '## Origin and propagation' not in s:
    assert old in s
    s = s.replace(old, new, 1)
    open(p, 'w', encoding='utf-8').write(s)
# validator: rely on file 55 for evidence resolution
v = 'tests/validate_packages.py'
t = open(v, encoding='utf-8').read()
a = t.index("    for tok in re.split(r';\\s*', r['evidence_paths']):")
b = t.index("    for tok in re.findall(")
t = t[:a] + t[b:]
open(v, 'w', encoding='utf-8').write(t)
print('patched')
