#!/usr/bin/env bash
# Phase 2A security smoke test against the DISPOSABLE test stack (docker-compose.pos-test.yml).
# Refuses to run unless the target is positively identified as disposable.
set -u
API="${API:-http://127.0.0.1:18080/api/v1}"
PG=pos-test-postgres
BE=pos-test-backend
PASS=0; FAIL=0

die() { echo "REFUSING TO RUN: $*" >&2; exit 2; }

# ---- isolation guards (fail closed) -------------------------------------------------
[ "$(docker inspect -f '{{ index .Config.Labels "com.pos.disposable" }}' "$PG" 2>/dev/null)" = "true" ] || die "$PG is not labelled com.pos.disposable=true"
[ "$(docker inspect -f '{{ index .Config.Labels "com.pos.disposable" }}' "$BE" 2>/dev/null)" = "true" ] || die "$BE is not labelled com.pos.disposable=true"
[ "$(docker exec "$PG" printenv POSTGRES_DB)" = "pos_test_disposable" ] || die "database name is not pos_test_disposable"
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$BE" | grep -qx 'POS_DISPOSABLE_TEST_STACK=1' || die "backend is not flagged POS_DISPOSABLE_TEST_STACK=1"
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$BE" | grep -qx 'ENV=test' || die "backend ENV is not test"
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$BE" | grep -qx 'PAYMOB_MODE=mock' || die "payments are not in mock mode"
case "$API" in http://127.0.0.1:*|http://localhost:*) ;; *) die "API target $API is not local" ;; esac
echo "isolation guards passed (disposable stack: $PG / $BE)"

psql_q() { docker exec "$PG" psql -U postgres -d pos_test_disposable -At -c "$1"; }
json() { python -c "import sys,json; d=json.load(sys.stdin); print(eval(sys.argv[1]))" "$1" 2>/dev/null; }
code() { curl -s -o /tmp/smoke_body.txt -w '%{http_code}' "$@"; }
check() { # name expected actual
  if [ "$2" = "$3" ]; then PASS=$((PASS+1)); echo "PASS  $1 (got $3)"; else FAIL=$((FAIL+1)); echo "FAIL  $1 (want $2, got $3)  body=$(head -c 160 /tmp/smoke_body.txt)"; fi
}

# ---- fixtures: two tenants, each with one restaurant and an owner ------------------------
SUFFIX=$(date +%s)
register() { # slug email -> prints token
  curl -s -X POST "$API/auth/register" -H 'Content-Type: application/json' \
    -d "{\"name\":\"Owner $1\",\"email\":\"$2\",\"password\":\"Test-Only-Passw0rd!\",\"restaurant_name\":\"Rest $1\",\"restaurant_slug\":\"$1\",\"phone\":\"+201000000000\"}"
}
RA=$(register "smoke-a-$SUFFIX" "owner-a-$SUFFIX@example.test")
RB=$(register "smoke-b-$SUFFIX" "owner-b-$SUFFIX@example.test")
TOKEN_A=$(echo "$RA" | json 'd["token"]')
TOKEN_B=$(echo "$RB" | json 'd["token"]')
TENANT_A=$(echo "$RA" | json 'd["user"]["tenant_id"]')
TENANT_B=$(echo "$RB" | json 'd["user"]["tenant_id"]')
REST_A=$(echo "$RA" | json 'd["user"]["restaurant_id"]')
REST_B=$(echo "$RB" | json 'd["user"]["restaurant_id"]')
[ -n "$TOKEN_A" ] && [ -n "$TOKEN_B" ] || die "fixture registration failed: A=$(echo "$RA" | head -c 200) B=$(echo "$RB" | head -c 200)"
echo "fixtures: tenant A=$TENANT_A/rest $REST_A, tenant B=$TENANT_B/rest $REST_B"
authA=(-H "Authorization: Bearer $TOKEN_A")
authB=(-H "Authorization: Bearer $TOKEN_B")

echo "--- SEC-A impersonation"
check "no credentials -> 401" 401 "$(code "$API/products")"
check "bare X-Impersonation-Token -> 401" 401 "$(code -H 'X-Impersonation-Token: x' -H "X-Tenant-ID: $TENANT_A" -H "X-Restaurant-ID: $REST_A" "$API/products")"
check "bare impersonation header + valid bearer -> 401 (no fallback)" 401 "$(code -H 'X-Impersonation-Token: x' "${authA[@]}" "$API/products")"

echo "--- SEC-B tenant context"
check "own token lists products -> 200" 200 "$(code "${authA[@]}" "$API/products")"
check "token A + tenant B headers still serves tenant A" 200 "$(code "${authA[@]}" -H "X-Tenant-ID: $TENANT_B" -H "X-Restaurant-ID: $REST_B" "$API/products")"
CAT_CODE=$(code -X POST "${authA[@]}" -H "X-Tenant-ID: $TENANT_B" -H "X-Restaurant-ID: $REST_B" -H 'Content-Type: application/json' -d "{\"name\":\"SpoofProbe-$SUFFIX\",\"name_ar\":\"SpoofProbe\",\"is_active\":true}" "$API/categories")
OWN_TENANT=$(psql_q "SELECT tenant_id FROM categories WHERE name='SpoofProbe-$SUFFIX'")
check "category created with spoofed tenant-B headers lands in tenant A (DB check)" "$TENANT_A" "$OWN_TENANT"
check "public route: mismatched tenant/restaurant pair -> 403" 403 "$(code -H "X-Tenant-ID: $TENANT_A" -H "X-Restaurant-ID: $REST_B" "$API/public/orders/1")"
check "public route: consistent pair reaches handler (404 no such order)" 404 "$(code -H "X-Tenant-ID: $TENANT_A" -H "X-Restaurant-ID: $REST_A" "$API/public/orders/999999")"

echo "--- SEC-C login confirm"
check "confirm without proof -> 401" 401 "$(code -X POST "$API/auth/login/confirm" -H 'Content-Type: application/json' -d "{\"email\":\"owner-a-$SUFFIX@example.test\",\"tenant_id\":$TENANT_A}")"
check "confirm with forged proof -> 401" 401 "$(code -X POST "$API/auth/login/confirm" -H 'Content-Type: application/json' -d "{\"email\":\"owner-a-$SUFFIX@example.test\",\"tenant_id\":$TENANT_A,\"selection_token\":\"forged\"}")"

echo "--- SEC-F theme isolation"
THEME=$(curl -s -X POST "$API/admin/themes" "${authA[@]}" -H 'Content-Type: application/json' -d "{\"name\":\"Smoke $SUFFIX\",\"slug\":\"smoke-$SUFFIX\",\"config\":{}}")
THEME_ID=$(echo "$THEME" | json 'd.get("id") or d.get("data",{}).get("id")')
if [ -n "$THEME_ID" ] && [ "$THEME_ID" != "None" ]; then
  check "owner A reads own theme -> 200" 200 "$(code "${authA[@]}" "$API/admin/themes/$THEME_ID")"
  check "owner B reads A's theme -> 404" 404 "$(code "${authB[@]}" "$API/admin/themes/$THEME_ID")"
  check "owner B updates A's theme -> 404" 404 "$(code -X PUT "${authB[@]}" -H 'Content-Type: application/json' -d '{"name":"pwned"}' "$API/admin/themes/$THEME_ID")"
  check "owner B activates A's theme -> 404" 404 "$(code -X POST "${authB[@]}" "$API/admin/themes/$THEME_ID/activate")"
  check "owner B deletes A's theme -> 404" 404 "$(code -X DELETE "${authB[@]}" "$API/admin/themes/$THEME_ID")"
  check "A's theme survives -> 200" 200 "$(code "${authA[@]}" "$API/admin/themes/$THEME_ID")"
else
  FAIL=$((FAIL+1)); echo "FAIL  could not create fixture theme: $(echo "$THEME" | head -c 200)"
fi

echo "--- SEC-D RBAC (enforce mode)"
check "owner has ADMIN-equivalent access to HR employees -> 200" 200 "$(code "${authA[@]}" "$API/hr/employees")"
# employee fixture: clone owner A's credentials into a plain 'employee' row (no role grants)
psql_q "INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, phone, avatar_url, role, status) SELECT tenant_id, restaurant_id, 'emp-a-$SUFFIX@example.test', password_hash, 'Employee A', '', '', 'employee', status FROM users WHERE email='owner-a-$SUFFIX@example.test'" >/dev/null
EMP=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"email\":\"emp-a-$SUFFIX@example.test\",\"password\":\"Test-Only-Passw0rd!\"}")
TOKEN_E=$(echo "$EMP" | json 'd["token"]')
authE=(-H "Authorization: Bearer $TOKEN_E")
if [ -n "$TOKEN_E" ] && [ "$TOKEN_E" != "None" ]; then
  check "employee without grants denied: HR salaries -> 403" 403 "$(code "${authE[@]}" "$API/hr/salaries")"
  check "employee without grants denied: create product -> 403" 403 "$(code -X POST "${authE[@]}" -H 'Content-Type: application/json' -d '{}' "$API/products")"
  check "owner still allowed: HR salaries -> 200" 200 "$(code "${authA[@]}" "$API/hr/salaries")"
else
  FAIL=$((FAIL+1)); echo "FAIL  could not log in employee fixture: $(echo "$EMP" | head -c 200)"
fi

echo "--- SEC-H subscription status"
check "active tenant business route -> 200" 200 "$(code "${authB[@]}" "$API/products")"
psql_q "UPDATE tenants SET is_disabled=true WHERE id=$TENANT_B" >/dev/null
sleep 31   # tenant access decisions are cached for 30s
check "disabled tenant business route -> 403" 403 "$(code "${authB[@]}" "$API/products")"
check "disabled tenant can still read /tenant/status -> 200" 200 "$(code "${authB[@]}" "$API/tenant/status")"
check "disabled tenant can still reach billing /subscription/current -> not 403" "ok" "$([ "$(code "${authB[@]}" "$API/subscription/current")" != 403 ] && echo ok || echo blocked)"
check "other tenant unaffected -> 200" 200 "$(code "${authA[@]}" "$API/products")"

echo "--- SEC-E public order object-level access"
KNOWN_TOKEN="smoke-known-token-$SUFFIX"
KNOWN_HASH=$(printf '%s' "$KNOWN_TOKEN" | sha256sum | cut -d' ' -f1)
PHONE_A='+20 100 555 0101'
psql_q "INSERT INTO orders (tenant_id, restaurant_id, order_number, customer_name, customer_phone, subtotal, total_amount, public_token_hash) VALUES ($TENANT_A, $REST_A, 'SMOKE-T-$SUFFIX', 'Token Customer', '+201005550101', 10, 10, '$KNOWN_HASH'), ($TENANT_A, $REST_A, 'SMOKE-L-$SUFFIX', 'Legacy Customer', '+201005550102', 10, 10, NULL)" >/dev/null
OID_T=$(psql_q "SELECT id FROM orders WHERE order_number='SMOKE-T-$SUFFIX'")
OID_L=$(psql_q "SELECT id FROM orders WHERE order_number='SMOKE-L-$SUFFIX'")
hA=(-H "X-Tenant-ID: $TENANT_A" -H "X-Restaurant-ID: $REST_A")
hB=(-H "X-Tenant-ID: $TENANT_B" -H "X-Restaurant-ID: $REST_B")
check "order by id, no proof -> 404" 404 "$(code "${hA[@]}" "$API/public/orders/$OID_T")"
check "order by id, wrong token -> 404" 404 "$(code "${hA[@]}" -H 'X-Order-Token: wrong' "$API/public/orders/$OID_T")"
check "order by id, correct token -> 200" 200 "$(code "${hA[@]}" -H "X-Order-Token: $KNOWN_TOKEN" "$API/public/orders/$OID_T")"
check "order by id, correct token but other tenant headers -> 404" 404 "$(code "${hB[@]}" -H "X-Order-Token: $KNOWN_TOKEN" "$API/public/orders/$OID_T")"
check "legacy order (no stored token) with any token -> 404" 404 "$(code "${hA[@]}" -H "X-Order-Token: $KNOWN_TOKEN" "$API/public/orders/$OID_L")"
check "number lookup, no proof -> 404" 404 "$(code "${hA[@]}" "$API/public/orders/number/SMOKE-T-$SUFFIX/lookup")"
check "number lookup, wrong phone -> 404" 404 "$(code "${hA[@]}" -H 'X-Customer-Phone: +201005559999' "$API/public/orders/number/SMOKE-T-$SUFFIX/lookup")"
check "number lookup, matching phone (formatted differently) -> 200" 200 "$(code "${hA[@]}" -H "X-Customer-Phone: $PHONE_A" "$API/public/orders/number/SMOKE-T-$SUFFIX/lookup")"
check "status by id with phone -> 200" 200 "$(code "${hA[@]}" -H "X-Customer-Phone: $PHONE_A" "$API/public/orders/$OID_T/status")"
check "status by id, no proof -> 404" 404 "$(code "${hA[@]}" "$API/public/orders/$OID_T/status")"
check "track by id, no proof -> 404" 404 "$(code "${hA[@]}" "$API/public/orders/$OID_T/track")"
check "cancel with phone only -> 404 (token required)" 404 "$(code -X DELETE "${hA[@]}" -H "X-Customer-Phone: $PHONE_A" "$API/public/orders/$OID_T")"
check "cancel, no proof -> 404" 404 "$(code -X DELETE "${hA[@]}" "$API/public/orders/$OID_T")"
check "non-existent order identical to forbidden (404)" 404 "$(code "${hA[@]}" "$API/public/orders/99999999")"
check "cancel with token -> not 404" "ok" "$([ "$(code -X DELETE "${hA[@]}" -H "X-Order-Token: $KNOWN_TOKEN" "$API/public/orders/$OID_T")" != 404 ] && echo ok || echo no)"

echo "--- C003 formerly ungated routes (employee without grants, enforce mode)"
if [ -n "${TOKEN_E:-}" ] && [ "$TOKEN_E" != "None" ]; then
  check "employee: GET /categories -> 403" 403 "$(code "${authE[@]}" "$API/categories")"
  check "employee: GET /users -> 403" 403 "$(code "${authE[@]}" "$API/users")"
  check "employee: GET /admin/drivers -> 403" 403 "$(code "${authE[@]}" "$API/admin/drivers")"
  check "employee: GET /admin/themes -> 403" 403 "$(code "${authE[@]}" "$API/admin/themes")"
  check "employee: GET /notifications -> 403" 403 "$(code "${authE[@]}" "$API/notifications")"
  check "employee: GET /rbac/roles -> 403" 403 "$(code "${authE[@]}" "$API/rbac/roles")"
  check "employee: POST /rbac/permissions/assign (privilege escalation) -> 403" 403 "$(code -X POST "${authE[@]}" -H 'Content-Type: application/json' -d '{"role_id":1,"module_id":1,"permission_level":"ADMIN"}' "$API/rbac/permissions/assign")"
  check "employee: PUT /tenant/contact-info -> 403" 403 "$(code -X PUT "${authE[@]}" -H 'Content-Type: application/json' -d '{}' "$API/tenant/contact-info")"
  check "employee: own /me still works -> 200" 200 "$(code "${authE[@]}" "$API/me")"
fi
check "owner: GET /categories -> 200" 200 "$(code "${authA[@]}" "$API/categories")"
check "owner: GET /users -> 200" 200 "$(code "${authA[@]}" "$API/users")"
check "owner: GET /admin/themes -> 200" 200 "$(code "${authA[@]}" "$API/admin/themes")"
check "owner: GET /notifications -> 200" 200 "$(code "${authA[@]}" "$API/notifications")"
check "owner: GET /rbac/roles -> 200" 200 "$(code "${authA[@]}" "$API/rbac/roles")"

echo "--- C007 unsafe defaults"
check "translations upload without credentials -> 401" 401 "$(code -X POST "$API/translations/en/upload" -H 'Content-Type: application/json' -d '{}')"
check "translations upload with a tenant owner token -> 401" 401 "$(code -X POST "${authA[@]}" "$API/translations/en/upload" -H 'Content-Type: application/json' -d '{}')"
check "metrics open in non-production without METRICS_TOKEN -> 200" 200 "$(code "$API/health/metrics")"
RL=0; for i in $(seq 1 8); do c=$(code -X POST "$API/auth/forgot-password" -H 'Content-Type: application/json' -d '{"email":"nobody@example.test"}'); [ "$c" = 429 ] && RL=1; done
check "forgot-password is rate limited (429 within 8 rapid calls)" 1 "$RL"

echo "--- C010 restricted database runtime role"
ATTRS=$(psql_q "SELECT rolsuper::text||','||rolbypassrls::text||','||rolcreatedb::text||','||rolcreaterole::text||','||rolreplication::text FROM pg_roles WHERE rolname='pos_app'")
check "pos_app has no superuser/bypassrls/createdb/createrole/replication" "false,false,false,false,false" "$ATTRS"
MATTRS=$(psql_q "SELECT rolsuper::text||','||rolcreaterole::text FROM pg_roles WHERE rolname='pos_migrator'")
check "pos_migrator is not a superuser and cannot create roles" "false,false" "$MATTRS"
APPCONN=$(psql_q "SELECT count(*) FROM pg_stat_activity WHERE datname='pos_test_disposable' AND client_addr IS NOT NULL AND usename='pos_app'")
SUCONN=$(psql_q "SELECT count(*) FROM pg_stat_activity WHERE datname='pos_test_disposable' AND client_addr IS NOT NULL AND usename='postgres'")
check "API holds connections as pos_app" "ok" "$([ "${APPCONN:-0}" -gt 0 ] && echo ok || echo none)"
check "API holds no network connection as the postgres superuser" "0" "$SUCONN"
app_sql() { docker exec -e PGPASSWORD=test-only-app-password "$PG" psql -h 127.0.0.1 -U pos_app -d pos_test_disposable -At -c "$1" >/tmp/app_sql.txt 2>&1; echo $?; }
check "pos_app can read business tables" 0 "$(app_sql 'SELECT count(*) FROM users')"
for stmt in "CREATE TABLE smoke_ddl_probe (id int)" "DROP TABLE users" "ALTER TABLE users ADD COLUMN smoke_probe int" "TRUNCATE users" "CREATE ROLE smoke_probe_role" "DROP SCHEMA public CASCADE" "CREATE FUNCTION smoke_f() RETURNS int LANGUAGE sql AS 'SELECT 1'" "SELECT * FROM schema_migrations" "UPDATE schema_migrations SET version='x'"; do
  rc=$(app_sql "$stmt")
  check "pos_app denied: $stmt" "denied" "$([ "$rc" != 0 ] && grep -qi 'permission denied\|must be owner\|permission' /tmp/app_sql.txt && echo denied || echo "allowed($rc)")"
done
check "pos_app cannot COPY to a program" "denied" "$(rc=$(app_sql "COPY (SELECT 1) TO PROGRAM 'id'"); [ "$rc" != 0 ] && echo denied || echo allowed)"

echo; echo "RESULT: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ]
