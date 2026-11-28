import os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
p = 'tests/build_finding_packages.py'
s = open(p, encoding='utf-8').read()
old = "## Monetary / operational impact\n{f['monetary_impact_EGP'] or 'Not quantified in the register.'}\n"
new = """## Severity rationale
Rated **{f['severity']}** (category: {f['category']}; evidence: {f['evidence_status']}). Basis recorded in the register: {f['monetary_impact_EGP'] or 'no monetary effect recorded'}. Root-cause confidence: {c}. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: {f['affected_modules']}. {f['monetary_impact_EGP'] or 'The register records no direct monetary effect; the impact is on correctness, usability or control quality of the listed modules.'}

## Technical impact
{f['root_cause_first_break'][:420]}

Data integrity: {'stored monetary or ledger data can become inconsistent (see financial-impact.md)' if (('accounting' in K or 'concurrency' in K) and not re.search(NON_MONEY_CATEGORY, f['category'], re.I)) else 'no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access'}. Security relevance: {'yes, see category' if 'security' in K else 'none recorded'}. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
{f['monetary_impact_EGP'] or 'Not quantified in the register.'}
"""
if '## Severity rationale' not in s:
    assert old in s
    s = s.replace(old, new, 1)
    open(p, 'w', encoding='utf-8').write(s)
print('patched')
