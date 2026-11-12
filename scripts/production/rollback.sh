#!/usr/bin/env bash
set -euo pipefail
[[ $# -eq 1 && "$1" =~ ^[0-9a-f]{40}$ ]] || { echo "usage: $0 <previous-40-character-release-sha>" >&2; exit 2; }
[[ "${POS_ROLLBACK_APPROVED:-false}" == true && "${POS_MIGRATION_BACKWARD_COMPATIBLE:-false}" == true ]] || {
  echo "rollback requires explicit approval and backward-compatible migrations" >&2; exit 1;
}
ROOT="${POS_RELEASE_ROOT:-/opt/rakmyat-pos/current}"
ENV_FILE="${POS_ENV_FILE:-/etc/rakmyat-pos/production.env}"
exec 9>/run/lock/rakmyat-pos-deploy.lock
flock -n 9 || { echo "another release holds the deployment lock" >&2; exit 1; }
export IMAGE_TAG="sha-$1"
docker compose --project-name rakmyat-pos-production --env-file "$ENV_FILE" -f "$ROOT/docker-compose.prod.yml" pull
docker compose --project-name rakmyat-pos-production --env-file "$ENV_FILE" -f "$ROOT/docker-compose.prod.yml" up -d --no-deps backend landing dashboard restaurant-website super-admin
"$ROOT/scripts/production/healthcheck.sh"
printf '%s\n' "$1" > /var/lib/rakmyat-pos/active-release
echo "application release rolled back to $1; database was not changed"
