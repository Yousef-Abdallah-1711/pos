#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="${POS_ENV_FILE:-/etc/rakmyat-pos/production.env}"
MIN_FREE_KB="${POS_MIN_FREE_KB:-26214400}"
MIN_MEM_KB="${POS_MIN_MEM_KB:-8388608}"

[[ -r "$ENV_FILE" ]] || { echo "missing protected environment file: $ENV_FILE" >&2; exit 1; }
MODE="$(stat -c '%a' "$ENV_FILE")"
[[ "$MODE" == 600 || "$MODE" == 640 ]] || { echo "environment file must have mode 0600 or 0640" >&2; exit 1; }
python3 "$ROOT/scripts/production/validate-env.py" "$ENV_FILE"
FREE_KB="$(df -Pk / | awk 'NR==2 {print $4}')"
MEM_KB="$(awk '/MemAvailable:/ {print $2}' /proc/meminfo)"
(( FREE_KB >= MIN_FREE_KB )) || { echo "insufficient free disk space" >&2; exit 1; }
(( MEM_KB >= MIN_MEM_KB )) || { echo "insufficient available memory" >&2; exit 1; }
docker info >/dev/null
docker compose version >/dev/null
for port in 13401 13402 13403 13404 13405; do
  if ss -H -ltn "sport = :$port" | grep -q .; then
    echo "planned POS loopback port $port is already occupied" >&2
    exit 1
  fi
done
EXPECTED_IPV4="${POS_EXPECTED_IPV4:-187.127.229.66}"
for host in pos.rakmyat.online pos.app.rakmyat.online admin.pos.rakmyat.online api.pos.rakmyat.online healthcheck.restaurants.pos.rakmyat.online; do
  if ! getent ahostsv4 "$host" | awk '{print $1}' | grep -Fxq "$EXPECTED_IPV4"; then
    echo "DNS for $host does not resolve to the configured VPS address" >&2
    exit 1
  fi
  tls="$(timeout 12 openssl s_client -connect "$host:443" -servername "$host" \
    -verify 5 -verify_hostname "$host" -verify_return_error -brief </dev/null 2>&1)" || {
    echo "valid HTTPS certificate check failed for $host" >&2; exit 1;
  }
  grep -q 'Verification: OK' <<<"$tls" || { echo "certificate verification failed for $host" >&2; exit 1; }
done
echo "host preflight passed"
