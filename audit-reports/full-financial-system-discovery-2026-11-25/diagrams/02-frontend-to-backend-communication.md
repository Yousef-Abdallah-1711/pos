# Diagram 02 - Frontend to backend communication

Legend: green = verified; red dashed = missing or broken connection.

```mermaid
sequenceDiagram
  autonumber
  participant UI as Page/component (dashboard)
  participant Hook as hook / lib/api/*.ts
  participant AX as axios api.ts (withCredentials, 6s timeout)
  participant MW as Go middleware (auth, tenantContext, statusGuard, requirePermission)
  participant H as Handler
  participant R as Repository
  participant DB as Postgres
  UI->>Hook: user action (e.g. post journal)
  Hook->>AX: api.post /accounting/journal-entries/{id}/post
  AX->>MW: HTTPS + HttpOnly session cookie (unsafe method requires trusted Origin)
  MW->>MW: session resolve (Redis) to claims; tenant from claims; tenant status; RBAC level
  MW->>H: request with tenant/restaurant/user in context
  H->>R: PostEntry(tenant, restaurant, id, user)
  R->>DB: BEGIN; SELECT FOR UPDATE; INSERT general_ledger; UPDATE status; COMMIT
  DB-->>H: ok
  H-->>AX: success envelope
  AX-->>UI: data (extractData)
  Note over UI,AX: No Idempotency-Key header on any money POST (verified). Client timeout 6000ms may abort while server commits (consequence UNVERIFIED).
  Note over UI: Client checkPermission only hides UI. Financial report pages have no client check; backend module check is the control.
```

```mermaid
flowchart TB
  classDef ok fill:#e6f4ea,stroke:#1e8e3e
  classDef bad fill:#fce8e6,stroke:#d93025,stroke-dasharray:5 3
  POS[POS page]:::ok -->|1 POST /pos/orders payment pending| B[(API)]:::ok
  POS -->|2 PUT payment-status paid - cash only, failure swallowed| B
  SITE[Storefront checkout]:::ok -->|quote then POST /public/orders, cash only| B
  SITE -.->|no UI path to paid| X[ledger]:::bad
  OFF[POS offline queue]:::ok -->|replay createOrder + payment-status, no idempotency| B
  ADV[Accounting advanced page]:::ok -->|CSV import, recurring run, bank rec| B
  PAY[Payroll page]:::ok -->|POST mark-paid, reference PAY-Date.now| B
  TB[Theme builder]:::bad -.->|admin/themes components, public/component-library: no route| B
  REDRIVE[Outbox redrive]:::bad -.->|no UI, API only| B
```
