#!/usr/bin/env bash
set -euo pipefail
[[ "${GITHUB_SHA:-}" =~ ^[0-9a-f]{40}$ ]] || { echo "GITHUB_SHA must be a full commit SHA" >&2; exit 2; }
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
REGISTRY="${IMAGE_REGISTRY:-ghcr.io/yousef-abdallah-1711/pos}"
TAG="sha-${GITHUB_SHA}"
MODE="${1:---push}"
case "$MODE" in --push|--load) ;; *) echo "usage: $0 [--push|--load]" >&2; exit 2;; esac
[[ "$MODE" != --push || -n "${GHCR_USERNAME:-}" && -n "${GHCR_TOKEN:-}" ]] || { echo "GHCR credentials are required for push" >&2; exit 1; }
if [[ "$MODE" == --push ]]; then printf '%s' "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USERNAME" --password-stdin; fi
build() {
  local service="$1" dockerfile="$2" context="$3"
  local action=--push
  [[ "$MODE" == --push ]] || action=--load
  local args=(--build-arg "APP=$service")
  if [[ "$context" == frontend ]]; then
    args+=(--build-arg NEXT_PUBLIC_SITE_URL=https://pos.rakmyat.online \
      --build-arg NEXT_PUBLIC_API_URL=https://api.pos.rakmyat.online/api/v1 \
      --build-arg NEXT_PUBLIC_API_BASE_URL=https://api.pos.rakmyat.online/api/v1 \
      --build-arg NEXT_PUBLIC_SUPER_ADMIN_API_URL=https://api.pos.rakmyat.online/api/v1/super-admin \
      --build-arg NEXT_PUBLIC_LANDING_URL=https://pos.rakmyat.online \
      --build-arg NEXT_PUBLIC_DASHBOARD_URL=https://pos.app.rakmyat.online \
      --build-arg NEXT_PUBLIC_RESTAURANT_DOMAIN=restaurants.pos.rakmyat.online \
      --build-arg NEXT_PUBLIC_DEBUG_LOGS=false --build-arg NEXT_PUBLIC_DEBUG=false)
  fi
  docker buildx build --platform linux/amd64 "$action" "${args[@]}" \
    --tag "$REGISTRY/$service:$TAG" --label "org.opencontainers.image.revision=$GITHUB_SHA" \
    --label org.opencontainers.image.source=https://github.com/Yousef-Abdallah-1711/pos \
    -f "$ROOT/$dockerfile" "$ROOT/$context"
}
build api backend/Dockerfile backend
for app in landing dashboard restaurant-website super-admin; do
  build "$app" frontend/Dockerfile.prod frontend
done
if [[ "$MODE" == --push ]]; then docker logout ghcr.io >/dev/null; fi
echo "built immutable images tagged $TAG"
