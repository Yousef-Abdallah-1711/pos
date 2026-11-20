# Restaurant storefront credential boundary

## Status

This note records the source review and code change in the
`codex/auth-boundary-review-20261117` branch. It is not a deployment receipt.

## Source trace

The restaurant website had a shared Axios interceptor that read `auth-token`
from `localStorage`, attached it as `Authorization: Bearer`, removed it on a
401, and redirected visitors to `/login`. A repository-wide search found no
caller of `apiClient.setAuthToken`, `clearAuthToken`, or `getAxiosInstance` in
the storefront. The only credential use came from that unused legacy path.

The interceptor now removes stale `auth-token` storage when possible and
strips Authorization headers before public requests. A 401 remains a request
error and does not redirect anonymous visitors to a staff login flow.
Development request logs now omit headers and response bodies, so order
capabilities and API response credentials are not written to those logs.

## Preserved public capabilities

- Anonymous menu browsing and order submission remain available.
- The API-issued per-order tracking capability remains in its separate
  `order-access-tokens` store and is sent only as `X-Order-Token` for matching
  public-order URLs.
- Tenant and restaurant context headers remain in place for public tenant
  routing. They are context hints, not authentication credentials; the API
  must continue validating the resource and tenant relationship.
- Cross-device order lookup continues to use the existing order-number and
  checkout-phone flow.

No customer account or staff session was added. Backend authorization,
Super Admin isolation, cookies, and order-capability validation are unchanged
by this slice.

## Verification

The new API-client test asserts that neither stale local storage nor a caller
provided Authorization header is sent, while the per-order capability and
tenant context headers are preserved. Browser checkout, cross-device lookup,
and production cache behavior still require their existing acceptance runs.

Backend middleware tests also exercise staff/Super Admin cookie cross-access
denial and trusted-Origin requirements for unsafe requests using each portal's
cookie. No backend authorization middleware changed in this slice.
