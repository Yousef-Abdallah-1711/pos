"""One-off codemod (Phase 2A C003): wrap authenticated-only business routes with permission checks.

Reads backend/cmd/api/main.go, rewrites `mux.Handle("METHOD /path", wrapProtected|wrapAuthOnly(...))`
into wrapWithPermission(..., MODULE, "LEVEL") / wrapAuthPermission(..., MODULE, "LEVEL") according to POLICY.
Prints every rewrite. Safe to re-run (already-wrapped routes are skipped).
"""
import re
import sys

PATH = "backend/cmd/api/main.go"

# (method, path) -> (module, level). Module ids come from the seeded `modules` table.
# 2 CATEGORIES, 4 NOTIFICATIONS, 5 SETTINGS, 6 USERS, 7 ROLES, 9 THEMES, 10 ORDERS (drivers live under orders).
POLICY = {
    ("GET", "/api/v1/categories"): (2, "READ"),
    ("GET", "/api/v1/categories/{id}"): (2, "READ"),
    ("POST", "/api/v1/categories"): (2, "WRITE"),
    ("PUT", "/api/v1/categories/{id}"): (2, "WRITE"),
    ("DELETE", "/api/v1/categories/{id}"): (2, "DELETE"),
    ("GET", "/api/v1/tenant/contact-info"): (5, "READ"),
    ("PUT", "/api/v1/tenant/contact-info"): (5, "WRITE"),
    ("GET", "/api/v1/users"): (6, "READ"),
    ("POST", "/api/v1/users"): (6, "WRITE"),
    ("PUT", "/api/v1/users/{id}"): (6, "WRITE"),
    ("DELETE", "/api/v1/users/{id}"): (6, "DELETE"),
    ("POST", "/api/v1/admin/drivers"): (10, "WRITE"),
    ("GET", "/api/v1/admin/drivers"): (10, "READ"),
    ("GET", "/api/v1/admin/drivers/available"): (10, "READ"),
    ("GET", "/api/v1/admin/drivers/{id}"): (10, "READ"),
    ("PUT", "/api/v1/admin/drivers/{id}"): (10, "WRITE"),
    ("DELETE", "/api/v1/admin/drivers/{id}"): (10, "DELETE"),
    ("PUT", "/api/v1/admin/drivers/{id}/status"): (10, "WRITE"),
    ("GET", "/api/v1/admin/drivers/{id}/stats"): (10, "READ"),
    # drivers report their own location with the DRIVER role, which only holds ORDERS:READ
    ("POST", "/api/v1/admin/drivers/{id}/location"): (10, "READ"),
    ("GET", "/api/v1/admin/drivers/{id}/location"): (10, "READ"),
    ("GET", "/api/v1/notifications"): (4, "READ"),
    ("GET", "/api/v1/notifications/stats"): (4, "READ"),
    ("GET", "/api/v1/notifications/unread-count"): (4, "READ"),
    ("POST", "/api/v1/notifications"): (4, "WRITE"),
    ("POST", "/api/v1/notifications/mark-all-read"): (4, "READ"),
    ("POST", "/api/v1/notifications/{id}/read"): (4, "READ"),
    ("POST", "/api/v1/notifications/{id}/unread"): (4, "READ"),
    ("GET", "/api/v1/notifications/{id}"): (4, "READ"),
    ("DELETE", "/api/v1/notifications/{id}"): (4, "DELETE"),
    ("GET", "/api/v1/admin/themes"): (9, "READ"),
    ("POST", "/api/v1/admin/themes"): (9, "WRITE"),
    ("GET", "/api/v1/admin/themes/{id}"): (9, "READ"),
    ("PUT", "/api/v1/admin/themes/{id}"): (9, "WRITE"),
    ("DELETE", "/api/v1/admin/themes/{id}"): (9, "DELETE"),
    ("POST", "/api/v1/admin/themes/{id}/activate"): (9, "WRITE"),
    ("POST", "/api/v1/admin/themes/{id}/duplicate"): (9, "WRITE"),
    ("GET", "/api/v1/admin/themes/{id}/history"): (9, "READ"),
    ("GET", "/api/v1/admin/theme-presets"): (9, "READ"),
    ("GET", "/api/v1/rbac/roles"): (7, "READ"),
    ("GET", "/api/v1/rbac/roles/{id}"): (7, "READ"),
    ("POST", "/api/v1/rbac/roles"): (7, "WRITE"),
    ("PUT", "/api/v1/rbac/roles/{id}"): (7, "WRITE"),
    ("DELETE", "/api/v1/rbac/roles/{id}"): (7, "DELETE"),
    ("POST", "/api/v1/rbac/permissions/assign"): (7, "ADMIN"),
    ("DELETE", "/api/v1/rbac/permissions/{roleId}/{moduleId}"): (7, "ADMIN"),
    ("GET", "/api/v1/rbac/roles/{roleId}/permissions"): (7, "READ"),
    ("POST", "/api/v1/rbac/users/{userId}/roles"): (7, "ADMIN"),
    ("DELETE", "/api/v1/rbac/users/{userId}/roles/{roleId}"): (7, "ADMIN"),
    ("GET", "/api/v1/rbac/users/{userId}/roles"): (7, "READ"),
    ("POST", "/api/v1/rbac/users-by-email/{email}/roles"): (7, "ADMIN"),
    ("DELETE", "/api/v1/rbac/users-by-email/{email}/roles/{roleId}"): (7, "ADMIN"),
    ("GET", "/api/v1/rbac/users-by-email/{email}/roles"): (7, "READ"),
}
# Deliberately left authenticated-only (self-service / status / billing): /me, /tenant/status, /subscription/*,
# /rbac/me/*, /rbac/check-permission, /users/{userId}/settings*, /users/{userId}/profile (handlers enforce self-only).

SRC = open(PATH, encoding="utf-8").read()


def match_paren(text: str, open_idx: int) -> int:
    """Index of the ')' that closes the '(' at open_idx; skips strings, runes and comments."""
    depth = 0
    i = open_idx
    n = len(text)
    while i < n:
        c = text[i]
        if c == '"':
            i += 1
            while text[i] != '"':
                i += 2 if text[i] == "\\" else 1
        elif c == "`":
            i = text.index("`", i + 1)
        elif c == "'":
            i += 1
            while text[i] != "'":
                i += 2 if text[i] == "\\" else 1
        elif text.startswith("//", i):
            i = text.index("\n", i)
        elif text.startswith("/*", i):
            i = text.index("*/", i) + 1
        elif c == "(":
            depth += 1
        elif c == ")":
            depth -= 1
            if depth == 0:
                return i
        i += 1
    raise ValueError("unbalanced parentheses")


pattern = re.compile(r'mux\.Handle\(\s*"([A-Z]+) ([^"]+)",\s*(wrapProtected|wrapAuthOnly)\(')
out = []
pos = 0
changed = []
for m in pattern.finditer(SRC):
    method, path, wrapper = m.group(1), m.group(2), m.group(3)
    policy = POLICY.get((method, path))
    if policy is None:
        continue
    open_idx = m.end() - 1  # the '(' of wrapX(
    close_idx = match_paren(SRC, open_idx)
    module, level = policy
    new_wrapper = "wrapWithPermission" if wrapper == "wrapProtected" else "wrapAuthPermission"
    out.append(SRC[pos:m.start(3)])
    out.append(new_wrapper)
    out.append(SRC[m.end(3):close_idx])
    out.append(f', {module}, "{level}"')
    pos = close_idx
    changed.append((method, path, wrapper, module, level))
out.append(SRC[pos:])
NEW = "".join(out)
open(PATH, "w", encoding="utf-8", newline="").write(NEW)
for c in changed:
    print(*c)
print(f"rewrote {len(changed)} registrations", file=sys.stderr)
