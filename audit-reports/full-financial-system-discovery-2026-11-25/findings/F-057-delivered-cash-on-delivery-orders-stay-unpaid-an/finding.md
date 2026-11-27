# F-057 — Delivered cash-on-delivery orders stay unpaid and no settlement step records the cash collected by the driver

- **Severity:** Medium  ·  **Category:** Missing capability  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Delivery; Drivers; Orders; Accounting
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-DRV-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: A COD settlement action (driver cash handover) marks the order paid and posts revenue/COGS, or delivery completion prompts for payment

Observed: status delivered with payment_status pending and 0 journals; revenue stays unrecognised until someone manually marks paid; no driver cash reconciliation exists | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): delivered COD order remained payment_status pending with 0 journals (REPRODUCED). 

## Monetary / operational impact
All unsettled COD sales are missing from GL revenue/VAT/cash

## Where it was reproduced
Create delivery order (cash), assign driver, move status preparing -> ready -> out_for_delivery -> delivered, then read payment_status and journals

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
