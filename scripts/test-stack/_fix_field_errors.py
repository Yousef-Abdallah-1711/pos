"""Phase 2A frontend codemod: render react-hook-form field errors through errorText().

Reads the tsc output of the dashboard, finds TS2322 'FieldError ... is not assignable to type ReactNode'
sites and rewrites `errors.<field>.message` / `errors.<field>?.message` on exactly those lines.
"""
import re
import sys
from collections import defaultdict

TSC = sys.argv[1]
ROOT = "frontend/apps/dashboard/"
sites = defaultdict(set)
for line in open(TSC, encoding="utf-8", errors="ignore"):
    m = re.match(r"(.+?)\((\d+),\d+\): error TS2322: Type 'string \| FieldError", line)
    if m:
        sites[m.group(1)].add(int(m.group(2)))

pattern = re.compile(r"errors\.([A-Za-z0-9_]+)\??\.message")
total = 0
for path, lines_to_fix in sites.items():
    full = ROOT + path
    src = open(full, encoding="utf-8").read().split("\n")
    changed = False
    for ln in lines_to_fix:
        new = pattern.sub(lambda m: f"errorText(errors.{m.group(1)})", src[ln - 1])
        if new != src[ln - 1]:
            src[ln - 1] = new
            changed = True
            total += 1
    if changed:
        text = "\n".join(src)
        if "errorText" in text and "from '@/lib/form-errors'" not in text:
            # add the import after the last top-level import line
            idx = max(i for i, l in enumerate(text.split("\n")) if l.startswith("import "))
            parts = text.split("\n")
            parts.insert(idx + 1, "import { errorText } from '@/lib/form-errors'")
            text = "\n".join(parts)
        open(full, "w", encoding="utf-8", newline="").write(text)
print("rewrote", total, "sites in", len(sites), "files")
