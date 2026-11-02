# Current architecture diagrams

Diagrams describe observed source/runtime boundaries, not a desired target.

## A. System context

```mermaid
flowchart TD
  U[Restaurant staff / platform admin / customers] --> L[Landing Next.js :3001]
  U --> D[Dashboard Next.js :3002]
  U --> S[Restaurant storefront Next.js :3003]
  U --> A[Super-admin Next.js :3004]
  L --> API[Go API :8080]
  D --> API
  S --> API
  A --> API
  API --> PG[(PostgreSQL 15)]
  API --> R[(Redis)]
  API --> P[Paymob live or mock adapter]
  API --> WS[KDS WebSocket]
```

## B. Backend dependency shape

```mermaid
flowchart LR
  Main[cmd/api main.go: composition + 363 route registrations] --> MW[Middleware wrappers]
  Main --> HTTP[HTTP handlers]
  Main --> UC[Use cases]
  Main --> Repo[Repositories]
  HTTP --> UC
  HTTP --> Repo
  HTTP --> SQL[Inline SQL / *sql.DB]
  UC --> Repo
  Repo --> SQL
  SQL --> PG[(PostgreSQL)]
  UC --> Paymob[Paymob interface]
  Main --> Redis[(Redis/cache)]
```

Direct handler→repository and SQL paths remain; target layering is not universal.

## C. Frontend workspace

```mermaid
flowchart TD
  L[Landing] --> P[Workspace packages]
  D[Dashboard] --> P
  S[Storefront] --> P
  A[Super Admin] --> P
  P --> T[Theme / UI packages and layers]
  P -. API clients are duplicated / no canonical shared client .-> API[Go API]
  L --> API
  D --> API
  S --> API
  A --> API
```

## D. Identity and tenant flow

```mermaid
sequenceDiagram
  participant C as Client
  participant M as Auth middleware
  participant T as Tenant context
  participant S as Tenant status
  participant R as Permission middleware
  participant H as Handler / repository
  C->>M: Bearer JWT (or session-bound impersonation)
  M->>T: authenticated identity / tenant claim
  T->>T: validate restaurant context against tenant
  T->>S: tenant and restaurant context
  S->>R: reject blocked tenant except allowed status/billing routes
  R->>H: compare requested module/permission
```

Resource ownership remains a separate check and is not globally verified.

## E. Database tenancy

```mermaid
erDiagram
  TENANTS ||--o{ RESTAURANTS : owns
  TENANTS ||--o{ USERS : scopes
  RESTAURANTS ||--o{ PRODUCTS : scopes
  RESTAURANTS ||--o{ ORDERS : scopes
  ORDERS ||--o{ ORDER_ITEMS : contains
  RESTAURANTS ||--o{ DRIVERS : scopes
```

This is representative. Existing schema has 91 tables, many derived by foreign keys. RLS is enabled on zero live tables; the relation diagram is not a complete schema map.

## F. Testing architecture

```mermaid
flowchart LR
  Go[Go unit/default suite: 680 prior verified] --> CI[Current CI]
  FE[Frontend tests: dashboard red; landing/storefront red; super-admin 3 pass] --> CI
  Smoke[Disposable API + PostgreSQL security smoke: 73/73 prior verified] --> CI
  E2E[Full browser / four-app E2E: missing] -. not wired .-> CI
  CI --> Guard[Two config guardrail workflows only]
```

CI currently lacks Go/frontend regression and disposable security jobs. Counts are from current checkpoint reports except the fresh landing/storefront runs recorded in `09_TEST_INFRASTRUCTURE_AUDIT.md`.
