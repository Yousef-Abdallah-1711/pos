# All Pages Inventory (summary)
Detail: `pages-dashboard.md`, `pages-others.md`, machine-readable `evidence/routes-dashboard.json` (87 files) and `evidence/routes-others.json`.

| App | Route files | Notes |
|---|---|---|
| Landing (3001) | 10 page files, 2 layouts, 1 middleware | Only 4 effective pages (`/[locale]`, `auth/login`, `auth/register`, `auth/super-admin-login`); 6 redirect stubs are shadowed by middleware. Bare `/login` & `/register` end at 404 (browser-verified). |
| Dashboard (3002) | 87 route files = 79 pages + 8 layouts; no `route.ts`, no middleware | All auth/locale handling client-side. 283 API call sites: 274 match `main.go`, 8 orphan, 1 method mismatch. |
| Storefront (3003) | 12 pages, 1 route handler, 2 layouts, 1 middleware | `GET /public/orders` and `/public/themes*` called but not registered in main.go (FE-OTH-015). |
| Super admin (3004) | 9 pages, 1 route handler, 2 layouts, no middleware | Auth is client-side. |

Total discovered: 110 page files (10 landing + 79 dashboard + 12 storefront + 9 super-admin). Status of each: see FINAL report §C — only anonymous/public pages were opened in a browser.
