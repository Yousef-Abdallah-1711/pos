#!/usr/bin/env bash
set -euo pipefail
BASE="${1:-https://pos.rakmyat.online}"
API="${2:-https://api.pos.rakmyat.online}"
for url in "$BASE" "https://pos.app.rakmyat.online" "https://admin.pos.rakmyat.online" "https://healthcheck.restaurants.pos.rakmyat.online" "$API/api/v1/health"; do
  code="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' --max-time 12 "$url")"
  [[ "$code" =~ ^[23][0-9][0-9]$ ]] || { echo "health check failed: $url returned $code" >&2; exit 1; }
done
echo "public endpoint checks passed"
