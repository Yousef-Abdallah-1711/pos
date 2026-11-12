#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TEMPLATE="$ROOT/deploy/cloudpanel/pos-vhosts.conf.template"
python3 - "$TEMPLATE" <<'PY'
from pathlib import Path
import re, sys
data = Path(sys.argv[1]).read_text()
servers = re.findall(r'listen 443 ssl http2;.*?server_name\s+([^;]+);', data, re.S)
expected = {
    'pos.rakmyat.online', 'pos.app.rakmyat.online', 'admin.pos.rakmyat.online',
    'api.pos.rakmyat.online', '*.restaurants.pos.rakmyat.online',
}
if set(servers) != expected or len(servers) != 5:
    raise SystemExit('expected five exact host patterns')
if re.search(r'proxy_pass\s+http://(?!127\.0\.0\.1:1340[1-5]\b)', data):
    raise SystemExit('an upstream is not one of the planned loopback ports')
if data.count('ssl_certificate {{') != 5 or data.count('ssl_certificate_key {{') != 5:
    raise SystemExit('all five hosts need a CloudPanel certificate placeholder')
if re.search(r'\b(?:cache|proxy_cache)\s+', data):
    raise SystemExit('authenticated traffic must not be cached by this vhost')
print('CloudPanel host and upstream mapping validated; certificate placeholders remain')
PY
if command -v nginx >/dev/null; then
  temp="$(mktemp -d)"
  trap 'rm -rf "$temp"' EXIT
  rendered="$temp/pos-vhosts.conf"
  openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
    -keyout "$temp/test.key" -out "$temp/test.crt" -subj /CN=localhost >/dev/null 2>&1
  cp "$TEMPLATE" "$rendered"
  for name in POS_APEX POS_DASHBOARD POS_ADMIN POS_API POS_RESTAURANT_WILDCARD; do
    sed -i "s|{{${name}_CERT}}|$temp/test.crt|g; s|{{${name}_KEY}}|$temp/test.key|g" "$rendered"
  done
  cat >"$temp/nginx.conf" <<EOF
pid $temp/nginx.pid;
error_log stderr;
events {}
http {
    include /etc/nginx/mime.types;
    include $rendered;
}
EOF
  nginx -t -p "$temp/" -c "$temp/nginx.conf"
fi
