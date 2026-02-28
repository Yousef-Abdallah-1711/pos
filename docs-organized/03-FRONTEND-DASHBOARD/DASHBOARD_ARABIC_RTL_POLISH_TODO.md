# Dashboard Arabic + RTL Polish TODO (Reports / Orders / Products / Theme Builder)

**Status:** In Progress  
**Scope:** `frontend/apps/dashboard` + reports export backend (`backend/internal/handler/http/report_handler.go`)  
**Goal:** Make Arabic UI/RTL and Arabic CSV/PDF exports production-ready.

## 1. Reports (UI + Export)

- [x] Reports pages use real backend data (no mock data for `overview`, `sales`, `profitability`, `inventory`)
- [x] Arabic UI strings added for reports pages/sidebar
- [x] Arabic CSV export labels for reports (`sales`, `profitability`, `inventory`)
- [x] Arabic PDF export font support added (UTF-8 Arabic font embedded)
- [x] PDF layout improved (cards/tables/spacing/page footer)
- [x] PDF section/table labels localized for Arabic
- [x] PDF no-data text localized (`لا توجد بيانات متاحة`)
- [x] PDF metadata labels fixed (`الفترة`, `التجميع`, `الحد الأعلى`, `تم الإنشاء`, `صفحة`, `إلى`)
- [ ] Visual QA pass in browser/PDF viewer for all 3 report PDFs in Arabic (not text extraction only)
- [ ] RTL table column alignment review for Arabic PDF (optional polish)

## 1.5 Dashboard Overview + Shared Layout (`/ar/dashboard`)

- [x] Dashboard overview KPI/action/status strings localized (fixed broken `????` Arabic values)
- [x] Dashboard overview mixed English subtitle/status chips localized via `ar.json`
- [x] Header action cluster order/alignment improved for RTL
- [x] Sidebar nested section padding switched to RTL padding (`pr-*`) instead of margin offset
- [ ] Visual QA for header/sidebar spacing on desktop + mobile in Arabic

## 2. Orders Page (`/ar/dashboard/orders`)

- [x] Pagination/footer summary string localized (`Showing X to Y of Z orders`)
- [x] Status/payment labels localized
- [x] CSV export localized for Arabic (headers + values + BOM for Excel)
- [x] Common action/button/toast strings localized
- [ ] Full Arabic page QA pass (filters, export, pagination, dialogs)
- [ ] Confirm no remaining English hardcoded strings in visible path

## 3. Products Pages (`/ar/dashboard/products`, `/ar/dashboard/products/new`, edit form)

- [x] Product form/new page major labels/placeholders localized
- [x] Fixed raw key display (`products.addProductDescription`)
- [x] RTL spacing fixes for checkbox/toggle rows
- [ ] Products list page Arabic audit (toolbar/actions/status badges/table labels)
- [ ] Product edit page visual QA in Arabic
- [ ] Validate all product form validation messages appear localized

## 4. Theme Builder (`/ar/dashboard/theme-builder`)

- [x] Top-level page/list strings localized (title/subtitle/actions/filters/card actions)
- [x] Debug info panel hidden in production unless `NEXT_PUBLIC_DEBUG_LOGS=true`
- [x] Debug panel forced `dir="ltr"` for readability when visible
- [x] Theme Builder editor tabs / color palettes / contrast checker strings localized
- [x] Theme Builder editor preview defaults to Arabic on `/ar/...` routes
- [x] Theme Builder editor page header/loading/error strings localized (`/editor/[id]`)
- [x] Theme Builder nested editor core panels localized (HeaderEditor / FooterEditor / FontPicker / ColorPicker)
- [ ] Arabic audit for nested Theme Builder dialogs/pages (import/export/template flows)
- [ ] Localize remaining deeper editor sub-components (`ContactEditor`, `ProductGridEditor`, `CTAEditor`, `GenericEditor` tips)
- [ ] Verify no remaining hardcoded English labels in common user flows

## 5. Dashboard Layout RTL (Header + Sidebar)

- [x] Header action cluster alignment improved for RTL/LTR
- [x] Profile text alignment fixed in RTL
- [x] Sidebar `Theme Builder` label localized
- [x] Sidebar collapse chevron direction respects RTL
- [x] Sidebar nested group indentation/padding mirrors correctly in RTL
- [ ] Visual QA for header spacing/order in Arabic on mobile + desktop
- [ ] Visual QA for sidebar spacing/icons on collapsed + expanded modes in RTL

## 6. QA / Production Readiness Checklist (Arabic/RTL)

- [ ] `/ar/dashboard/reports/overview` renders with no `????` UI text
- [ ] `/ar/dashboard/reports/sales` UI + CSV + PDF verified
- [ ] `/ar/dashboard/reports/profitability` UI + CSV + PDF verified
- [ ] `/ar/dashboard/reports/inventory` UI + CSV + PDF verified
- [ ] `/ar/dashboard/orders` UI + CSV verified
- [ ] `/ar/dashboard/drivers` UI verified (table alignment/status badges/actions)
- [x] `/ar/dashboard/drivers` rating symbol mojibake fixed
- [ ] `/ar/dashboard/products` UI verified
- [ ] `/ar/dashboard/products/new` UI/form validation verified
- [ ] `/ar/dashboard/theme-builder` top-level list page verified
- [ ] `/ar/dashboard/theme-builder/editor/[id]` colors/header/footer/settings/history tabs verified
- [ ] `/ar/dashboard/settings` + account/restaurant/billing/access-control pages verified in Arabic
- [ ] Browser hard refresh performed after i18n updates (`Ctrl+F5`)
- [ ] Capture screenshots/PDF samples for signoff

## Notes

- Arabic PDF **text extraction** from some tools may still appear reversed/shaped incorrectly even when the visible PDF renders correctly. Judge by PDF viewer rendering, not extraction output alone.
- Reports data and orders data are real DB/API data (not mock) in the dashboard pages covered above.
