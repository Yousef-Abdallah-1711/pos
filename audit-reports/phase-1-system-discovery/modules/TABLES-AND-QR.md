# Tables and QR (dashboard) - static audit

Status: **Verified implemented** (table CRUD, status, floor-plan layout, QR generation); QR design is single-tenant.

## Pages
`tables/page.tsx` (459; gate TABLES 7 refs; imports `qrcode`), `tables/[id]/page.tsx`, `tables/[id]/edit`, `tables/[id]/qr`, `tables/layout/page.tsx` (floor-plan editor, 461 lines). Components: `FloorPlanEditor`, `FloorPlanTable`, `QRStudio`, `QRPreview`, `AddTableForm`, `TableList`, `TableStatusBadge`; `QRBuilderModal.tsx` is unreferenced by pages (heuristic). API: `lib/api/tablesApi.ts`, hook `hooks/useTables.ts`. 6 vitest files in `tables/components`.

## Contract (all match main.go:1305-1314, RBAC module 14)
POST/GET `/tables`, GET `/tables/stats`, GET/POST `/tables/layouts`, PUT `/tables/positions` (floor-plan save), GET/PUT/DELETE `/tables/{id}`, PUT `/tables/{id}/status`.

## QR
- `lib/table-qr.ts:19-40`: URL = `{NEXT_PUBLIC_RESTAURANT_DOMAIN}/{locale}/table/{tableID}/menu`. Domain comes from one deploy-time env (`config.ts:16,109`, default `localhost:3003`), table id is the integer PK. Customer page calls public `GET /public/table/{id}/info` (main.go:465, IP rate-limited).
- QR rendered client-side with `qrcode` (+ `qrcode.react`), downloadable PNG/SVG, file-name sanitised (`sanitizeQrFileName`). Also `theme-builder/qr` page and `WebsiteQRStudio` for site QR.
- Tests: `lib/table-qr.test.ts`, `website-qr.test.ts`.
- Gap: in multi-tenant production every tenant's QR would point to the same domain unless the env is per deployment; ids are guessable (FE-DASH-025).
- `qrcode` has no `@types/qrcode` devDependency (TS error source, FE-DASH-009).

## Notes
Dashboard overview calls `GET /tables/stats` and backs off on 404 (`dashboard/page.tsx:350-356`). POS table selector uses `/tables`. Table `status` set manually; no auto-link to open orders verified.

## Findings
FE-DASH-025, FE-DASH-009, FE-DASH-022.
