# Shared Global Requests and Duplication

| Surface | Static request/dependency | Trigger | Duplication or blocking concern | Verification |
|---|---|---|---|---|
| Dashboard | Auth bootstrap; optional impersonation validation | Provider mount / token state | Children stay hidden until auth and hydration finish | STATIC_ONLY; anonymous redirect observed only |
| Dashboard | User preferences | Shared dashboard layout mount | May gate or compete with first route requests; exact await behavior needs trace | STATIC_ONLY |
| Dashboard | Route preloading | Dashboard layout | Can compete for bandwidth/compile before click; source mentions common route preload | STATIC_ONLY |
| Dashboard | Notifications | Bell/hook polling | Background polling; source includes visibility handling in hook | STATIC_ONLY |
| Dashboard | KDS + reports | KDS/reports routes | Polling can duplicate refresh work; visibility policy needs per hook verification | STATIC_ONLY |
| Restaurant Website | Tenant status | Middleware on tenant subdomain | Additional request before route continuation; failure is fail-open | STATIC_ONLY |
| Restaurant Website | Theme | Metadata and page/server fetch | May repeat logically equivalent fetch; Next request memoization/runtime cache hit not captured | STATIC_ONLY |
| Restaurant Website | Order tracking | Tracking hooks | 10–15s refetch; hidden-tab pause and key scope need review | STATIC_ONLY |
| Super Admin | Auth bootstrap | Root provider mount | Wraps login and protected paths | STATIC_ONLY |

No HAR exists. The only measured public landing page navigation recorded zero `/api/` resource entries. That does not rule out same-origin or server-side data requests.
