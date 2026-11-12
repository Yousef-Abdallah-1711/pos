# Page Usability and Loading Timeline

## Measured public page

Landing `/en` began navigation in isolated Chromium. Resource Timing reports response start at 10.259 s on first isolated-context access and 2.277 s on repeat; DOMContentLoaded at 11.279 s / 2.378 s; load at 12.100 s / 3.075 s. These times are document events, not meaningful usability. The page title and initial marketing page content were present at completion; no explicit primary-action readiness condition was instrumented.

## Blocked page

Anonymous `/en/dashboard` navigation waited 60 seconds for DOMContentLoaded then timed out. Browser ended up at Landing `/en/auth/login`. This supports an auth redirect path but provides no reliable dashboard readiness or latency number.

## Readiness condition not met

A page-specific condition (for example dashboard table visible and primary action enabled, storefront menu rendered with synthetic menu data) was not instrumented for any page. All other 109 pages are UNMEASURED. No authenticated synthetic accounts were exercised. No 3G/slow API/failure emulation was run.
