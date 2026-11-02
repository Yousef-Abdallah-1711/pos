# Route Inventory: landing, restaurant-website, super-admin

Machine-readable copy: `evidence/routes-others.json`. Backend column refers to `backend/cmd/api/main.go` registrations.

## Landing (3001) - 10 page files (4 effective), 2 layouts, 1 middleware
| URL | File | Notes |
|---|---|---|
| `/` | src/app/page.tsx | redirect('/en'); shadowed by middleware |
| `/login` | src/app/login/page.tsx | redirect('/en/auth/login'); middleware sends to `/en/login` -> 404 |
| `/register` | src/app/register/page.tsx | same defect |
| `/auth/login` | src/app/auth/login/page.tsx | stub; middleware redirects first |
| `/auth/register` | src/app/auth/register/page.tsx | stub |
| `/auth/super-admin-login` | src/app/auth/super-admin-login/page.tsx | stub |
| `/{en,ar}` | src/app/[locale]/page.tsx | marketing page: Navbar, Hero, Features, HowItWorks, Pricing, Testimonials, CTA, Footer |
| `/{locale}/auth/login` | src/app/[locale]/auth/login/page.tsx | POST /auth/login, POST /auth/login/confirm |
| `/{locale}/auth/register` | src/app/[locale]/auth/register/page.tsx | GET /auth/check-subdomain, POST /auth/register |
| `/{locale}/auth/super-admin-login` | src/app/[locale]/auth/super-admin-login/page.tsx | POST /auth/token-exchange |
| layouts | src/app/layout.tsx, src/app/[locale]/layout.tsx | |
| middleware | src/middleware.ts | locale prefix |
Dead link targets: `/{locale}/auth/forgot-password`, `/demo`, `/about`, `/blog`, `/help`, `/docs`, `/api`, `/privacy`, `/terms`, `#about`, `#contact`.

## Restaurant website (3003) - 12 pages, 1 route handler, 2 layouts, error, not-found, middleware
| URL | File | Backend |
|---|---|---|
| `/` | app/page.tsx | none (dev placeholder; middleware redirects) |
| `/maintenance` | app/maintenance/page.tsx | none (query params) |
| `/api/metrics` | app/api/metrics/route.ts | Next handler only |
| `/{locale}` | app/[locale]/page.tsx | homepage, theme, products (exist) |
| `/{locale}/menu` | app/[locale]/menu/page.tsx | categories, products, search (exist) |
| `/{locale}/cart` | app/[locale]/cart/page.tsx | none |
| `/{locale}/checkout` | app/[locale]/checkout/page.tsx | POST /public/orders |
| `/{locale}/orders` | app/[locale]/orders/page.tsx | GET /public/orders - NOT REGISTERED |
| `/{locale}/orders/[id]` | app/[locale]/orders/[id]/page.tsx | GET /public/orders/{id}, /status |
| `/{locale}/orders/track` | app/[locale]/orders/track/page.tsx | lookup + track |
| `/{locale}/orders/track/[orderNumber]` | .../track/[orderNumber]/page.tsx | lookup + track |
| `/{locale}/settings` | app/[locale]/settings/page.tsx | none (local preferences store) |
| `/{locale}/table/[tableId]/menu` | app/[locale]/table/[tableId]/menu/page.tsx | GET /public/table/{id}/info |
| layouts / error / not-found | app/layout.tsx, app/[locale]/layout.tsx, error.tsx, not-found.tsx | |
| middleware | middleware.ts | GET /public/tenant-status/{slug} |

## Super admin (3004) - 9 pages, 1 route handler, 2 layouts, no middleware
| URL | File | Backend |
|---|---|---|
| `/` | src/app/page.tsx | none |
| `/login` | src/app/login/page.tsx | POST /super-admin/auth/login |
| `/dashboard` | src/app/dashboard/page.tsx | GET dashboard/stats |
| `/dashboard/tenants` | .../tenants/page.tsx | GET tenants |
| `/dashboard/tenants/[id]` | .../tenants/[id]/page.tsx | tenant, users, subscription, history, extend-trial, disable, enable, delete, dashboard-login, refund, impersonate |
| `/dashboard/users` | .../users/page.tsx | GET users (links to missing `/dashboard/users/{id}`) |
| `/dashboard/payments` | .../payments/page.tsx | payment-metrics, webhook-health |
| `/dashboard/audit-logs` | .../audit-logs/page.tsx | GET audit-logs |
| `/dashboard/settings` | .../settings/page.tsx | PUT settings/ip-whitelist |
| `/favicon.ico` | src/app/favicon.ico/route.ts | |
| layouts | src/app/layout.tsx, src/app/dashboard/layout.tsx | |

Totals: landing 4 effective pages (+6 stubs), restaurant-website 12 pages, super-admin 9 pages.
