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
ACTIVE_RELEASE_FILE=/var/lib/rakmyat-pos/active-release
previous_release=""
if [[ -r "$ACTIVE_RELEASE_FILE" ]]; then previous_release="$(cat "$ACTIVE_RELEASE_FILE")"; fi

wait_for_backend() {
  for _ in $(seq 1 60); do
    if curl --fail --silent --show-error --max-time 3 http://127.0.0.1:13405/api/v1/health >/dev/null 2>&1; then return 0; fi
    sleep 2
  done
  echo "local API health check timed out" >&2
  return 1
}

rollback_on_failure() {
  rc=$?
  trap - ERR
  set +e
  if [[ "$previous_release" =~ ^[0-9a-f]{40}$ ]]; then
    echo "release failed; restoring application images from $previous_release" >&2
    export IMAGE_TAG="sha-$previous_release"
    "${COMPOSE[@]}" pull backend landing dashboard restaurant-website super-admin
    "${COMPOSE[@]}" up -d --no-deps backend
    wait_for_backend
    "${COMPOSE[@]}" up -d --no-deps landing dashboard restaurant-website super-admin
    "$ROOT/scripts/production/healthcheck.sh"
    if [[ $? -eq 0 ]]; then printf '%s\n' "$previous_release" > "$ACTIVE_RELEASE_FILE"; fi
  else
    echo "first deployment failed; stopping only POS application containers" >&2
    "${COMPOSE[@]}" stop backend landing dashboard restaurant-website super-admin
  fi
  exit "$rc"
}
trap rollback_on_failure ERR
"$ROOT/scripts/production/backup-postgres.sh"
"${COMPOSE[@]}" pull
"${COMPOSE[@]}" up -d postgres redis
"${COMPOSE[@]}" run --rm migrate
"${COMPOSE[@]}" up -d --no-deps backend
wait_for_backend
"${COMPOSE[@]}" up -d --no-deps landing dashboard restaurant-website super-admin
"$ROOT/scripts/production/healthcheck.sh"
printf '%s\n' "$1" > "$ACTIVE_RELEASE_FILE"
trap - ERR
echo "release $1 deployed and healthy"
