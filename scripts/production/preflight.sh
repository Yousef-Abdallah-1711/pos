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
echo "host preflight passed"
