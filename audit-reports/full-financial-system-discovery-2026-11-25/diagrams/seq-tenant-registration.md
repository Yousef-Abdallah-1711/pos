# Sequence: tenant registration / restaurant initialization

Source: `backend/internal/handler/http/auth_handler.go:47`, `backend/internal/usecase/auth_usecase.go:99`, `backend/internal/repository/auth_repo.go:42-210`. Label: CONFIRMED-BY-SOURCE. The "not created" notes are confirmed by grep (no accounting seeding in the registration path) and by catalog queries, not by running registration (REQUIRES-EXECUTION).

```mermaid
sequenceDiagram
    autonumber
    participant UI as Landing app
    participant H as AuthHandler.Register
    participant UC as AuthUseCase.RegisterContext
    participant R as AuthRepository.CreateTenantWithUserContext
    participant DB as PostgreSQL
    UI->>H: POST /api/v1/auth/register
    H->>UC: validate (registration_validation.go)
    UC->>R: create
    R->>DB: INSERT tenants (slug retry loop, trial fields)
    R->>DB: INSERT users (owner, password hash)
    R->>DB: INSERT restaurants
    R->>DB: INSERT themes_v2
    R->>DB: INSERT roles + role_permissions + user_roles (RBAC seed)
    R-->>UC: user
    UC-->>H: tokens / session
    Note over DB: NOT created: chart_of_accounts, auto_posting_rules, fiscal_periods, bank_accounts_ledger, opening balances, tax periods
    Note over DB: Runtime catalog: tenant 2 has 0 accounts, 0 rules, 0 fiscal periods; tenant 1 (migration seeded) has 62, 6, 0
```
