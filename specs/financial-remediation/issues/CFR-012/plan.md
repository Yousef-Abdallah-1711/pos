# CFR-012 implementation plan — recovery draft

## Verified candidate paths

- `backend/internal/usecase/financial_report_usecase.go`
- `backend/internal/handler/http/financial_report_handler.go`
- PR #114 `backend/internal/usecase/cfr010_cash_flow_void_integration_test.go`

## Approach

1. Preserve the existing difference evidence and review #114's movement-date correction separately from cash-flow methodology.
2. Obtain CF-01 approval for method, cash equivalents, mappings and treatment of tax/noncash transactions; define a versioned mapping and migration only if required.
3. Build synthetic cash and noncash scenarios with an independent GL/bank numerical oracle, including void and closed-period boundaries; prove opening plus net movement equals closing cash.
4. Implement only approved methodology; validate scope, failure cases, export/comparatives and performance on a reviewed candidate and final main.

## Gaps

No direct/indirect method, cash-equivalent policy or category mapping is approved. A passing void-date test does not close this issue.
