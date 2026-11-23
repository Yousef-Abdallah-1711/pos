# CFR-007: Binary floating point is used across monetary domain and APIs

Status: HISTORICAL — CURRENT DISPOSITION REQUIRES VERIFICATION. Original severity: P2. This is a source-backed recovery draft, not acceptance.

## Original issue text

### CFR-007 — Binary floating point is used across monetary domain and APIs

- **Severity/status:** P2 — SOURCE-CONFIRMED; precision behavior not broadly property-tested.
- **Problem / impact:** Repeated decimal arithmetic, allocation and serialization can produce fractional-cent differences and inconsistent totals across Go, PostgreSQL, JavaScript and exports.
- **Current / expected:** Accounting/order/payroll models and repositories use `float64`, while DB uses `NUMERIC/DECIMAL`; JSON is numeric. Expected: one canonical currency-aware exact representation at boundaries, without silent precision loss.
- **Root cause:** `backend/internal/domain/accounting.go`, `domain/order.go`, `backend/internal/usecase/order_usecase.go:127-182`, `backend/internal/handler/http/hr_salary_handler.go:122-141`; numeric schema e.g. migrations `137`, `141`, `024`.
- **Reproduction/evidence:** Source-confirmed type mismatch; no property test or cross-layer precision experiment in this run.
- **Baseline scope:** Present main/current.
- **Recommended solution / components:** Define integer minor units or arbitrary-precision decimal contract, exact JSON string/decimal policy, currency exponent, supported rounding modes, allocation residual and FX handling. Migrate incrementally with compatibility converters and reconciliation.
- **Migration:** Potential staged type/API migrations; preserve old records and verify sums before any conversion.
- **Dependencies/policy:** Supported currencies and rounding policy require owner/accountant approval.
- **Acceptance:** property tests for fractional line sums, tax inclusive/exclusive, allocations, split/refund/FX, overflow/negative amounts; API→DB→report/export equality at approved precision.
- **Checklist:** OPEN.



## Current evidence and open questions

Read the later dated sections of the [current tracker](../../../../FINANCIAL-MODULE-COMPLETE-ISSUES-AND-REMEDIATION-PLAN.md). Confirm each current behavior, policy, and source path against the final candidate; historical text above is not a current verdict.

## Definition of done

Confirm root cause and approved behavior; demonstrate the failing regression and corrected result; run relevant PostgreSQL, isolation, API and browser checks; obtain required review and merge; reconcile independent financial evidence and document any required policy approval. Mark non-applicable checks with evidence.
