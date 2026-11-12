# Critical User Journey Performance (Static / Partial Runtime)

All journeys are source-mapped only at a broad app/page/endpoint-family level unless noted. No authenticated account was used and no order/payment/inventory mutation was issued.

| Journey | Exists in source? | Main shared dependencies | Runtime state | Candidate concern |
|---|---|---|---|---|
| Login → Dashboard | Yes | Landing auth; dashboard auth callback/provider | Not exercised; anonymous dashboard redirect path observed but no valid readiness number | Token handoff + dashboard bootstrap |
| Dashboard → POS | Yes | Dashboard layout/auth, POS data hooks | Not measured | Provider/auth gating; endpoint dependencies not joined |
| POS → Create Order | Yes | POS API, order handler/usecase/repository | Mutation not exercised | Do not benchmark writes without isolated synthetic lifecycle |
| POS → Inventory | Yes | Dashboard navigation, inventory route | Not measured | Global shell + missing route skeleton coverage |
| Inventory → Adjustment | Route source exists; current branch has dirty handler addition | Inventory adjustment endpoint/handler | Mutation not exercised | Pagination/query and audit trail need trace |
| Dashboard → Orders | Yes | Admin/orders pages and list endpoints | Not measured | Static N+1 at admin order list |
| Orders → KDS | Yes | Order API, KDS polling/WebSocket paths | Not measured | Poll visibility and duplicate request behavior |
| Restaurant Website → Menu | Yes | Host slug middleware, tenant status, public menu APIs | Not measured; existing prior audit saw missing demo tenant, but not reused as current runtime evidence | Middleware gate and data cache scope |
| Restaurant Website → Checkout | Source route exists | Public order/customer/payment APIs | Not exercised | Never trigger live payment; synthetic end-to-end fixture needed |
| Super Admin → Tenant List | Yes | Super Admin provider + list API | Not measured | Provider wraps routes; pagination and response size unknown |
| Super Admin → Tenant Details | Yes | Tenant details/subscription/history APIs | Not measured | Multiple sequential fetches need browser waterfall |
| Settings → Theme Editor | Yes | Dashboard auth/layout, theme APIs, dynamic editor components | Not measured | Large editor/client dependencies and query provider wiring |
| HR navigation | Yes | Dashboard shell, HR endpoints | Not measured | Page count and shared provider effects |
| Accounting navigation | Yes | Dashboard shell, financial report cache/endpoints | Not measured | Cache hit/miss and report query cost unknown |

Only the Landing homepage was measured in browser. One safe health endpoint was sampled. No end-to-end business journey timing exists.
