#!/usr/bin/env bash
set -euo pipefail
[[ "$(id -u)" == 0 ]] || { echo "run this on the target host as root" >&2; exit 1; }
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEST="${POS_CONFIG_DIR:-/etc/rakmyat-pos}"
[[ ! -e "$DEST/production.env" && ! -e "$DEST/backup.gpg.key" && ! -e "$DEST/redis-password" ]] || {
  echo "protected POS configuration already exists; refusing to overwrite" >&2; exit 1;
}
install -d -o root -g root -m 0700 "$DEST"
if getent group 19990 >/dev/null; then
  [[ "$(getent group 19990 | cut -d: -f1)" == rakmyat-pos-secrets ]] || {
    echo "GID 19990 is already assigned; refusing to reuse an unrelated group" >&2; exit 1;
  }
else
  groupadd --system --gid 19990 rakmyat-pos-secrets
fi
umask 077
python3 - "$ROOT/.env.prod.example" "$DEST/production.env" "$DEST/redis-password" <<'PY'
from pathlib import Path
import secrets, sys
source, target, redis_file = map(Path, sys.argv[1:])
values = {
    "DB_BOOTSTRAP_PASSWORD": secrets.token_hex(32),
    "DB_PASSWORD": secrets.token_hex(32),
    "MIGRATE_DB_PASSWORD": secrets.token_hex(32),
    "REDIS_PASSWORD": secrets.token_hex(32),
    "JWT_SECRET": secrets.token_hex(64),
}
lines = source.read_text().splitlines()
result = []
for line in lines:
    key = line.partition("=")[0]
    if key in values:
        line = f"{key}={values[key]}"
    if key == "IMAGE_TAG":
        line = "IMAGE_TAG=sha-" + "0" * 40
    result.append(line)
target.write_text("\n".join(result) + "\n")
redis_file.write_text(values["REDIS_PASSWORD"] + "\n")
PY
openssl rand -hex 48 >"$DEST/backup.gpg.key"
chmod 600 "$DEST/production.env" "$DEST/redis-password" "$DEST/backup.gpg.key"
chown root:19990 "$DEST/redis-password"
chmod 640 "$DEST/redis-password"
echo "generated distinct local secrets and a protected environment template; payment credentials and release gates remain unset"
