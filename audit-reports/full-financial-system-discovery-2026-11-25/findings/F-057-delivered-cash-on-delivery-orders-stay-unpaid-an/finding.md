# F-057 — Delivered cash-on-delivery orders stay unpaid and no settlement step records the cash collected by the driver

- **Severity:** Medium  ·  **Category:** Missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Delivery; Drivers; Orders; Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-DRV-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment

Observed: status delivered with payment_status pending and 0 journals; revenue stays unrecognised until someone manually marks paid; no driver cash reconciliation exists | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): delivered COD order remained payment_status pending with 0 journals (REPRODUCED). 

## Severity rationale
Rated **Medium** (category: Missing capability; evidence: REPRODUCED). Basis recorded in the register: All unsettled COD sales are missing from GL revenue/VAT/cash. Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Delivery; Drivers; Orders; Accounting. All unsettled COD sales are missing from GL revenue/VAT/cash

## Technical impact
CONFIRMED (RUN-004 source): CompleteDelivery only changes the order status to delivered (order_usecase.go:433-447) and records no collection or payment

Data integrity: stored monetary or ledger data can become inconsistent (see financial-impact.md). Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
All unsettled COD sales are missing from GL revenue/VAT/cash

## Where it was reproduced
Create delivery order (cash), assign driver, move status preparing -> ready -> out_for_delivery -> delivered, then read payment_status and journals

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
