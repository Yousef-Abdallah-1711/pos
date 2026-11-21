#!/bin/sh
# Every suite runs for useful failure evidence; any failure rejects the image gate.
set -eu
if [ ! -d /tmp/browser/node_modules/playwright ]; then
  echo 'Runner-installed Playwright dependencies are required.' >&2
  exit 1
fi
failed=0
for script in browser-smoke.cjs storefront-flow.cjs browser-upload.cjs \
  journal-api-acceptance.cjs browser-auth-validation.cjs marketing-acceptance.mjs \
  marketing-console-diagnostics.test.mjs trial-registration.mjs capture-landing-acceptance.mjs lighthouse-acceptance.mjs; do
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
