#!/usr/bin/env bash
set -euo pipefail
[[ "${1:-}" =~ ^[0-9a-f]{40}$ ]] || { echo "usage: $0 <40-character-release-sha>" >&2; exit 2; }
[[ "${POS_PRODUCTION_READINESS_APPROVED:-false}" == true && "${POS_BACKUP_RESTORE_VERIFIED:-false}" == true ]] || {
  echo "deployment gates are not approved" >&2; exit 1;
}
ROOT="${POS_RELEASE_ROOT:-/opt/rakmyat-pos/current}"
ENV_FILE="${POS_ENV_FILE:-/etc/rakmyat-pos/production.env}"
COMPOSE=(docker compose --project-name rakmyat-pos-production --env-file "$ENV_FILE" -f "$ROOT/docker-compose.prod.yml")
export IMAGE_TAG="sha-$1"
"$ROOT/scripts/production/preflight.sh"
exec 9>/run/lock/rakmyat-pos-deploy.lock
flock -n 9 || { echo "another release holds the deployment lock" >&2; exit 1; }
"$ROOT/scripts/production/backup-postgres.sh"
"${COMPOSE[@]}" pull
"${COMPOSE[@]}" up -d postgres redis
"${COMPOSE[@]}" run --rm migrate
"${COMPOSE[@]}" up -d --remove-orphans backend landing dashboard restaurant-website super-admin
"$ROOT/scripts/production/healthcheck.sh"
printf '%s\n' "$1" > /var/lib/rakmyat-pos/active-release
echo "release $1 deployed and healthy"
