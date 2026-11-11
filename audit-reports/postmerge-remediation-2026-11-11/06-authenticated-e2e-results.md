# Authenticated E2E results

## New verification

- KDS WebSocket hook unit tests: 2 passed, including cookie-backed session connection eligibility and unauthenticated non-connection.
- Browser loaded the landing page at `/en` successfully with no browser-console errors. This was a UI render check, not an authentication E2E.

## Prior bounded evidence retained

The prior isolated run reported passing disposable checks for signup (201), logout (204), valid login (200), invalid login (401), unauthenticated reads (401), own settings (200), cross-user settings denial (403), category and product creation (201), cash order creation (201), payment status update (200), and one KDS status transition. Those results apply only to their exact paths and prior disposable stack.

## Not run in this remediation environment

Registration/login/session refresh and expiration, protected route lifecycle, tenant selection and cross-tenant isolation, inventory effects, storefront order flow, super-admin access controls, payment recovery, full KDS WebSocket event delivery, and suspended-tenant behavior. Docker and PostgreSQL were unavailable; no real authenticated acceptance is claimed for this candidate.
