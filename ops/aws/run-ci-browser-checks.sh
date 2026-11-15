#!/bin/sh
# Every suite runs for useful failure evidence; any failure rejects the image gate.
set -eu
if ! npm install --prefix /tmp/browser --ignore-scripts --no-audit --no-fund \
  playwright@1.64.1 axe-core@4.11.0 lighthouse@13.6.0 >/tmp/install.log 2>&1; then
  tail -30 /tmp/install.log
  exit 1
fi
failed=0
for script in browser-smoke.cjs storefront-flow.cjs browser-upload.cjs \
  journal-api-acceptance.cjs browser-auth-validation.cjs marketing-acceptance.mjs \
  trial-registration.mjs capture-landing-acceptance.mjs lighthouse-acceptance.mjs; do
  echo "Running $script"
  if ! NODE_PATH=/tmp/browser/node_modules node "/work/$script"; then
    failed=1
  fi
  if [ "$script" = marketing-acceptance.mjs ]; then
    for combination in 'chromium light' 'firefox light' 'firefox dark'; do
      set -- $combination
      browser=$1
      theme=$2
      echo "Running $browser $theme theme marketing acceptance"
      if ! SERVIO_MARKETING_BROWSER="$browser" SERVIO_THEME="$theme" SERVIO_COLOR_SCHEME="$theme" \
        SERVIO_MARKETING_EVIDENCE="${SERVIO_BROWSER_EVIDENCE:-/evidence}/marketing-$browser-$theme" \
        NODE_PATH=/tmp/browser/node_modules node /work/marketing-acceptance.mjs; then
        failed=1
      fi
    done
  fi
done
exit "$failed"
