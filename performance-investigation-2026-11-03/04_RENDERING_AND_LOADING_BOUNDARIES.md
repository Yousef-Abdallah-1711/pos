# Rendering and Loading Boundaries

## Page-file and boundary inventory

| App | Page files | `loading.tsx` files | `error.tsx` files | Static assessment |
|---|---:|---:|---:|---|
| Dashboard | 79 | 3 | Not fully inventoried | Loading boundaries exist only in selected route segments; dashboard-wide auth provider has its own full-screen loading state. |
| Landing | 10 | 0 | Not fully inventoried | No route-level loading boundary found. |
| Restaurant Website | 12 | 0 | 1 (`[locale]/error.tsx`) | No route-level loading boundary found; server page waits for its data promises. |
| Super Admin | 9 | 0 | Not fully inventoried | No route-level loading boundary found; provider/layout have loading behavior. |
| **Total** | **110** | **3** | **at least 1** | Counts are source file counts, not runtime coverage. |

Exact page file paths are in `02_COMPLETE_FRONTEND_PAGE_INVENTORY.csv`. The three dashboard loading files are `frontend/apps/dashboard/src/app/[locale]/dashboard/reports/inventory/loading.tsx`, `.../reports/profitability/loading.tsx`, and `.../reports/sales/loading.tsx`. The Restaurant Website error boundary is `frontend/apps/restaurant-website/app/[locale]/error.tsx`.

## Findings

- Dashboard auth hydration blocks child route content behind a full-screen spinner (`frontend/apps/dashboard/src/components/providers/DashboardAuthProvider.tsx:33–38, 58, 121–139, 232–244`; wrapper at `src/app/[locale]/dashboard/layout.tsx:1–17`). The fixed 100–300 ms waits are explicit code delays; their user-visible cost was not isolated from network/render time.
- Dashboard locale layout adds global client providers and a Suspense fallback described in source as an empty full-height shell. Determine which of these are client boundaries and their bundle cost with a production build before recommending removal.
- Restaurant Website root layout globally mounts providers; tenant-status middleware adds an awaited server-side API dependency on tenant subdomains (`frontend/apps/restaurant-website/middleware.ts:11–86`). Homepage server fetches are parallel (`app/[locale]/page.tsx:120–156`) but the returned page awaits all of them. No streaming / component-level progressive readiness was runtime verified.
- Loading boundary counts do not prove a skeleton is good or bad: layout/provider/component-level loading may exist. Visual shape, layout shift, keyboard responsiveness and readiness were not systematically tested.
- Runtime browser test reached Landing `/en` and a second Landing `/en` navigation. Anonymous dashboard navigation timed out in the browser tool while development compilation/redirect occurred; it ultimately landed on Landing `/en/auth/login`. No authenticated route readiness was measured.

## Page-level matrix

`19_SKELETON_AND_PERCEIVED_PERFORMANCE_MATRIX.csv` maps each page file to its app-level loading boundary status. `NOT_VERIFIED` means no visual/runtime confirmation, not that no component skeleton can exist.

