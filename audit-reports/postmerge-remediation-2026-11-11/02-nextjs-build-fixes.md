# Next.js production build results

All apps were built independently on Next.js 15.5.9 with production URL variables matching `docker-compose.prod.yml` (`NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_API_BASE_URL`, `NEXT_PUBLIC_SUPER_ADMIN_API_URL`, `NEXT_PUBLIC_LANDING_URL`, `NEXT_PUBLIC_DASHBOARD_URL`, and `NEXT_PUBLIC_RESTAURANT_DOMAIN`). No build-error bypass was enabled.

| App | Result | Notes |
|---|---|---|
| Landing | PASS | Static locale pages, `/robots.txt`, and `/sitemap.xml` generated. |
| Dashboard | PASS | 139 static pages generated; `/en/dashboard/settings` and `/ar/dashboard/settings` generated. |
| Restaurant website | PASS | Build completed; one existing `fetchOrder` effect-dependency lint warning remains. |
| Super Admin | PASS | Build completed; existing effect-dependency warnings remain in tenant detail and subscription history pages. |

The dashboard settings index now uses `redirect()` on the server to route to `/{locale}/dashboard/settings/account`. A clean standalone rebuild was not possible because Docker Desktop is stopped; these were host Next.js builds after a frozen-lockfile workspace install.
