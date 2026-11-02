# Products and Categories (dashboard) - static audit

Status: **Verified implemented** (CRUD present); categories pages ungated client-side and backend auth-only.

## Pages
- Products: `products/page.tsx` (209; gate PRODUCTS, 8 refs), `products/new/page.tsx`, `products/[id]/edit/page.tsx` (gate PRODUCTS). Form: `components/products/ProductFormPage.tsx` (used) - `ProductForm.tsx` is an unreferenced duplicate (heuristic).
- Categories: `categories/page.tsx` (174), `new`, `[id]/edit` - **no useHasPermission/MODULE_IDS in any of the 3 pages**; shared `components/categories/CategoryFormPage.tsx`.

## Contract (all match main.go)
GET/POST `/products`, GET/PUT/DELETE `/products/{id}` (RBAC module 1 via wrapWithPermission, main.go:663-668); GET `/products/alerts/low-stock` (no UI caller). GET/POST `/categories`, GET/PUT/DELETE `/categories/{id}` use `wrapProtected` only (main.go:671-675) i.e. authentication, no RBAC module 2 (FE-DASH-014).

## Behaviour
- Product create/update submits `multipart/form-data` with JSON in field `data` plus optional `image` file (`ProductFormPage.tsx:84-120`); image upload handled by backend `product_handler.go`. Client coerces numbers (price, cost, discount, stock, calories). Bilingual fields name/name_ar, description_en/ar.
- List page size / pagination handled server side (`products/page.tsx` uses GET with params); POS caps at 500 (FE-DASH-017).
- Delete uses confirm + DELETE; no soft-delete indicator in UI.
- Raw form values logged via `debugLog` (`ProductFormPage.tsx:113-114`; debug gated by lib/debug).
- Validation: react-hook-form; schemas in `lib/validation.ts` (unreferenced per heuristic scan) so validation is mostly inline.

## Findings
FE-DASH-014 (categories ungated), FE-DASH-026 (error strings), FE-DASH-017 (500 cap).

## Unverified
Image size/type validation (client or server), duplicate SKU handling, category-with-products delete behaviour.
