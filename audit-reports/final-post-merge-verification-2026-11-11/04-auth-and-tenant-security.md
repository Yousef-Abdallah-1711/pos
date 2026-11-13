# Authentication and tenant security

Tests used disposable accounts and the isolated tmpfs database.

- Registration via the real landing form: API returned 201; the user session was created and the browser navigated to the dashboard.
- Logout via the dashboard: API returned 204 and redirected to login.
- Login using the newly created disposable account: API returned 200; the dashboard loaded the authenticated account.
- A made-up account login returned 401 with the expected invalid-credentials response.
- Authenticated dashboard calls for session, user settings, RBAC permissions, orders, products, categories, employees, drivers, notifications, and tenant status returned 200.
- Unauthenticated requests to session, user settings, products, and admin orders returned 401.
- Two disposable tenants were created. The second tenant reading user 2's settings returned 403; reading its own user 3 settings returned 200.
- CORS preflight from `http://localhost:13111` to the isolated API returned 200 and the matching allow-origin header. Direct browser testing was possible only with the temporary port forward described in report 02.
- These are bounded route checks, not full authentication acceptance. Session expiration, refresh rotation, password recovery, CSRF behavior, arbitrary tenant/resource IDs, role mutation, suspended tenant, super-admin, and broader cookie policy checks remain NOT RUN.
- No valid account credentials are included in this report.

