#!/usr/bin/env bash
set -euo pipefail
[[ $# -ge 1 && $# -le 2 ]] || { echo "usage: $0 /path/to/encrypted.dump.gpg [/path/to/encrypted-uploads.tar.gz.gpg]" >&2; exit 2; }
backup="$1"
uploads="${2:-}"
key="${POS_BACKUP_KEY_FILE:-/etc/rakmyat-pos/backup.gpg.key}"
[[ -r "$backup" && -r "$key" ]] || { echo "backup or key unavailable" >&2; exit 1; }
[[ -r "$backup.sha256" ]] || { echo "backup checksum sidecar is required" >&2; exit 1; }
(cd "$(dirname "$backup")" && sha256sum --check "$(basename "$backup").sha256")
if [[ -n "$uploads" ]]; then
  [[ -r "$uploads.sha256" ]] || { echo "upload backup checksum sidecar is required" >&2; exit 1; }
  (cd "$(dirname "$uploads")" && sha256sum --check "$(basename "$uploads").sha256")
fi
command -v docker >/dev/null || { echo "Docker is required for an isolated restore test" >&2; exit 1; }
[[ "$(stat -c '%a' "$key")" == 600 ]] || { echo "backup key must be mode 0600" >&2; exit 1; }
suffix="$(od -An -N8 -tx1 /dev/urandom | tr -d ' \n')"
name="pos-restore-test-$suffix"
volume="${name}-uploads"
container_created=false
volume_created=false
cleanup() {
  if [[ "$container_created" == true ]]; then docker rm -f "$name" >/dev/null 2>&1 || true; fi
  if [[ "$volume_created" == true ]]; then docker volume rm "$volume" >/dev/null 2>&1 || true; fi
}
trap cleanup EXIT
if docker container inspect "$name" >/dev/null 2>&1 || docker volume inspect "$volume" >/dev/null 2>&1; then
  echo "isolated restore-test resource name collision" >&2
  exit 1
fi
docker run --detach --rm --name "$name" --network none --read-only \
  --tmpfs /var/lib/postgresql/data:rw,noexec,nosuid,size=2g \
  -e POSTGRES_PASSWORD=throwaway-restore-test-only -e POSTGRES_DB=pos_restore_test \
  postgres:15.19-alpine3.24@sha256:eab009ce9447a891d9482702d370bbf132c18563bb87cfdcc2b815c6d3b0372f >/dev/null
container_created=true
for _ in $(seq 1 60); do docker exec "$name" pg_isready -U postgres >/dev/null 2>&1 && break; sleep 2; done
docker exec "$name" pg_isready -U postgres >/dev/null
docker exec "$name" psql -v ON_ERROR_STOP=1 -U postgres -d pos_restore_test -c \
  'CREATE ROLE pos_bootstrap; CREATE ROLE pos_migrator; CREATE ROLE pos_app'
gpg --batch --quiet --pinentry-mode loopback --passphrase-file "$key" --decrypt "$backup" \
  | docker exec -i "$name" pg_restore --exit-on-error --no-owner -U postgres -d pos_restore_test
docker exec "$name" psql -U postgres -d pos_restore_test -Atc \
  "SELECT to_regclass('public.schema_migrations') IS NOT NULL" | grep -qx t
if [[ -n "$uploads" ]]; then
  docker volume create "$volume" >/dev/null
  volume_created=true
  gpg --batch --quiet --pinentry-mode loopback --passphrase-file "$key" --decrypt "$uploads" \
    | docker run --rm -i --network none --volume "$volume:/restore" \
      alpine:3.23@sha256:469c9ea410b9215eaee7976ba6c0639039b444be8e5d638f7abbada189971eb2 \
      sh -ec 'tar -xzf - -C /restore'
fi
echo "isolated restore test passed; disposable container removed on exit"
