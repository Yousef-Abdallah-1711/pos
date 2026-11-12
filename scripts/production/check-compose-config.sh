#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TEMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TEMP_DIR"' EXIT
ENV_FILE="$TEMP_DIR/compose-check.env"
CONFIG_FILE="$TEMP_DIR/compose-config.json"
REDIS_SECRET="$TEMP_DIR/redis-password"
printf '%s' 'compose-validation-only-redis-secret' >"$REDIS_SECRET"
chmod 600 "$REDIS_SECRET"

python3 - "$ENV_FILE" "$REDIS_SECRET" <<'PY'
import sys
from pathlib import Path

Path(sys.argv[1]).write_text("""IMAGE_REGISTRY=ghcr.io/yousef-abdallah-1711/pos
IMAGE_TAG=sha-0000000000000000000000000000000000000000
PUBLIC_API_URL=https://api.pos.rakmyat.online
LANDING_URL=https://pos.rakmyat.online
DASHBOARD_URL=https://pos.app.rakmyat.online
SUPER_ADMIN_URL=https://admin.pos.rakmyat.online
RESTAURANT_DOMAIN=restaurants.pos.rakmyat.online
CORS_ORIGINS=https://pos.rakmyat.online,https://pos.app.rakmyat.online,https://admin.pos.rakmyat.online,https://*.restaurants.pos.rakmyat.online
DB_NAME=pos_saas
DB_BOOTSTRAP_USER=pos_bootstrap
DB_BOOTSTRAP_PASSWORD=compose-validation-only-bootstrap-secret
DB_PASSWORD=compose-validation-only-runtime-secret
MIGRATE_DB_PASSWORD=compose-validation-only-migrator-secret
REDIS_PASSWORD=compose-validation-only-redis-secret
REDIS_PASSWORD_FILE={redis_file}
JWT_SECRET=compose-validation-only-jwt-secret-012345678901
PAYMOB_MODE=live
PAYMOB_API_KEY=
PAYMOB_INTEGRATION_ID=
PAYMOB_IFRAME_ID=
PAYMOB_HMAC_SECRET=
LANDING_HOST_PORT=13401
DASHBOARD_HOST_PORT=13402
RESTAURANT_HOST_PORT=13403
SUPER_ADMIN_HOST_PORT=13404
API_HOST_PORT=13405
""".replace("{redis_file}", sys.argv[2]), encoding="utf-8")
PY
chmod 600 "$ENV_FILE"

docker compose --env-file "$ENV_FILE" -f "$REPO_ROOT/docker-compose.prod.yml" config --format json >"$CONFIG_FILE"
python3 - "$CONFIG_FILE" <<'PY'
import json
import sys
from pathlib import Path

compose = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
services = compose["services"]
required = {"postgres", "redis", "migrate", "backend", "landing", "dashboard", "restaurant-website", "super-admin"}
if set(services) != required:
    raise SystemExit("production Compose service set is incomplete or unexpected")
for service in ("postgres", "redis", "migrate"):
    if services[service].get("ports"):
        raise SystemExit(f"{service} must not publish a host port")
for name, service in services.items():
    for port in service.get("ports", []):
        if not port.get("published", "").isdigit() or not port.get("host_ip", "").startswith("127."):
            raise SystemExit(f"{name} publishes a non-loopback port")
for name in ("postgres", "redis"):
    if not services[name].get("volumes"):
        raise SystemExit(f"{name} must use persistent storage")
for name in required:
    if name in {"postgres", "redis"}:
        continue
    if "sha-0000000000000000000000000000000000000000" not in services[name].get("image", ""):
        raise SystemExit(f"{name} does not use the immutable test release tag")
PY
printf 'production Compose structure validated\n'
