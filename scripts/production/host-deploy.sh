#!/usr/bin/env bash
set -euo pipefail
[[ "${1:-}" =~ ^[0-9a-f]{40}$ ]] || { echo "usage: $0 <release-sha>" >&2; exit 2; }
[[ "$(id -u)" == 0 ]] || { echo "run through the root-owned deployment wrapper" >&2; exit 1; }
release="$1"
base=/opt/rakmyat-pos/releases
target="$base/$release"
install -d -o root -g root -m 0755 "$base" /var/lib/rakmyat-pos
if [[ ! -d "$target" ]]; then
  temp="$base/.${release}.tmp"
  trap 'rm -rf "$temp"' EXIT
  install -d -o root -g root -m 0755 "$temp"
  git -C "$temp" init -q
  git -C "$temp" remote add origin https://github.com/Yousef-Abdallah-1711/pos.git
  git -C "$temp" fetch --depth=1 origin "$release"
  git -C "$temp" checkout --detach FETCH_HEAD
  [[ "$(git -C "$temp" rev-parse HEAD)" == "$release" ]] || { echo "release SHA mismatch" >&2; exit 1; }
  rm -rf "$temp/.git"
  mv "$temp" "$target"
fi
ln -sfn "$target" /opt/rakmyat-pos/current.new
mv -Tf /opt/rakmyat-pos/current.new /opt/rakmyat-pos/current
exec /opt/rakmyat-pos/current/scripts/production/deploy.sh "$release"
