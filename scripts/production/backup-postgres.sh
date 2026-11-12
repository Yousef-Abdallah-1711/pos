#!/usr/bin/env bash
set -euo pipefail
ENV_FILE="${POS_ENV_FILE:-/etc/rakmyat-pos/production.env}"
BACKUP_DIR="${POS_BACKUP_DIR:-/var/backups/rakmyat-pos}"
KEY_FILE="${POS_BACKUP_KEY_FILE:-/etc/rakmyat-pos/backup.gpg.key}"
ROOT="${POS_RELEASE_ROOT:-/opt/rakmyat-pos/current}"
COMPOSE=(docker compose --project-name rakmyat-pos-production --env-file "$ENV_FILE" -f "$ROOT/docker-compose.prod.yml")
[[ -r "$ENV_FILE" && -r "$KEY_FILE" ]] || { echo "environment or backup key is unavailable" >&2; exit 1; }
[[ "$(stat -c '%a' "$ENV_FILE")" == 600 && "$(stat -c '%a' "$KEY_FILE")" == 600 ]] || { echo "secret files must be mode 0600" >&2; exit 1; }
mkdir -p "$BACKUP_DIR"
chmod 700 "$BACKUP_DIR"
free_kb="$(df -Pk "$BACKUP_DIR" | awk 'NR==2 {print $4}')"
(( free_kb >= 10485760 )) || { echo "less than 10 GiB free at backup destination" >&2; exit 1; }
stamp="$(date -u +%Y%m%dT%H%M%SZ)-$$"
tmp="$BACKUP_DIR/.pos-${stamp}.dump.gpg.tmp"
out="$BACKUP_DIR/pos-${stamp}.dump.gpg"
uploads_tmp="$BACKUP_DIR/.pos-uploads-${stamp}.tar.gz.gpg.tmp"
uploads_out="$BACKUP_DIR/pos-uploads-${stamp}.tar.gz.gpg"
[[ ! -e "$out" && ! -e "$out.sha256" && ! -e "$uploads_out" && ! -e "$uploads_out.sha256" ]] || {
  echo "backup destination already exists; refusing to overwrite" >&2; exit 1;
}
trap 'rm -f "$tmp" "$uploads_tmp"' EXIT
"${COMPOSE[@]}" exec -T postgres sh -ec 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' \
  | gpg --batch --yes --pinentry-mode loopback --passphrase-file "$KEY_FILE" --symmetric --cipher-algo AES256 --output "$tmp"
gpg --batch --quiet --pinentry-mode loopback --passphrase-file "$KEY_FILE" --decrypt "$tmp" \
  | "${COMPOSE[@]}" exec -T postgres pg_restore --list >/dev/null
chmod 600 "$tmp"
mv -n "$tmp" "$out"
[[ ! -e "$tmp" ]] || { echo "database backup destination collision" >&2; exit 1; }
sha256sum "$out" >"$out.sha256"
chmod 600 "$out.sha256"
docker run --rm --network none --read-only \
  --volume rakmyat-pos-production-uploads:/source:ro \
  alpine:3.23@sha256:469c9ea410b9215eaee7976ba6c0639039b444be8e5d638f7abbada189971eb2 \
  tar -C /source -czf - . \
  | gpg --batch --yes --pinentry-mode loopback --passphrase-file "$KEY_FILE" --symmetric --cipher-algo AES256 --output "$uploads_tmp"
gpg --batch --quiet --pinentry-mode loopback --passphrase-file "$KEY_FILE" --decrypt "$uploads_tmp" \
  | docker run --rm -i --network none alpine:3.23@sha256:469c9ea410b9215eaee7976ba6c0639039b444be8e5d638f7abbada189971eb2 tar -tzf - >/dev/null
chmod 600 "$uploads_tmp"
mv -n "$uploads_tmp" "$uploads_out"
[[ ! -e "$uploads_tmp" ]] || { echo "uploads backup destination collision" >&2; exit 1; }
sha256sum "$uploads_out" >"$uploads_out.sha256"
chmod 600 "$uploads_out.sha256"
echo "encrypted database and upload backups verified: $out and $uploads_out"
