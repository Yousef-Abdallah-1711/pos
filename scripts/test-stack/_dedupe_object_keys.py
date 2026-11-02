"""Phase 2A frontend codemod: resolve TS1117 (duplicate keys in an object literal) without changing behaviour.

In `{ key: default, ...spread, key: normalized }` the later `key` always wins, so the earlier one is dead code.
For every flagged (later) property this deletes the nearest preceding property with the same key at the same
indentation, including multi-line values.

usage: _dedupe_object_keys.py <tsc-output.txt> <relative-source-file-as-printed-by-tsc>
"""
import re
import sys

tsc, target = sys.argv[1], sys.argv[2]
ROOT = "frontend/apps/dashboard/"

flagged = sorted(
    {int(m.group(1)) for line in open(tsc, encoding="utf-8", errors="ignore")
     if line.startswith(target + "(") and (m := re.match(re.escape(target) + r"\((\d+),\d+\): error TS1117", line))},
)
path = ROOT + target
raw = open(path, encoding="utf-8").read()
crlf = "\r\n" in raw
lines = raw.replace("\r\n", "\n").split("\n")


def delete_property(start: int) -> int:
    """Remove the property starting at lines[start]; returns number of lines removed."""
    depth = 0
    end = start
    while end < len(lines):
        text = re.sub(r"'(?:\\.|[^'\\])*'|\"(?:\\.|[^\"\\])*\"|`(?:\\.|[^`\\])*`", "", lines[end])
        depth += sum(text.count(c) for c in "([{") - sum(text.count(c) for c in ")]}")
        if depth <= 0 and lines[end].rstrip().endswith(","):
            break
        end += 1
    removed = end - start + 1
    del lines[start:end + 1]
    return removed


removed_total = 0
shift = 0  # lines deleted so far (always above the next flagged line), so later targets move up
for ln in flagged:
    ln -= shift
    later = lines[ln - 1]
    m = re.match(r"^(\s*)([A-Za-z_][A-Za-z0-9_]*)\s*:", later)
    if not m:
        print("skip (not a plain key)", ln, later.strip()[:60])
        continue
    indent, key = m.group(1), m.group(2)
    for i in range(ln - 2, -1, -1):
        cur = lines[i]
        if cur.strip() and not cur.startswith(indent) and not cur.strip().startswith(("}", "]", ")")):
            break  # left the literal
        if re.match(rf"^{re.escape(indent)}{re.escape(key)}\s*:", cur):
            removed = delete_property(i)
            removed_total += removed
            shift += removed
            break
    else:
        print("no earlier duplicate found for", key, "at", ln)

text = "\n".join(lines)
if crlf:
    text = text.replace("\n", "\r\n")
open(path, "w", encoding="utf-8", newline="").write(text)
print("removed", removed_total, "dead lines for", len(flagged), "duplicate keys in", target)
