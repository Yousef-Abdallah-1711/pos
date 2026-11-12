# ADR-003 — Ledger immutability

**Status: PROPOSED — not accepted.**

## Context
The historical audit demonstrates incomplete or unsafe financial paths and leaves multiple workflows untested. Existing source and the partial architecture blueprint are evidence inputs, not policy approval.

## Options
A. Keep current implicit behavior and accept reconciliation ambiguity.
B. Adopt the explicit invariant and workflow contract described below.
C. Defer pending evidence or owner/accountant choice.

## Recommendation
Ledger immutability should follow a versioned, auditable contract with explicit tenant scope, provenance and correction behavior. Prefer immutable posted effects, idempotent processing and a reconciliation trail. This recommendation is provisional; policy-dependent choices require owner/accountant approval.

## Consequences and risks
Requires migration compatibility, backfill/reconciliation, tests and operational ownership. Potential costs include more schema, worker complexity and visible eventual consistency. Never rewrite posted history or claim exactly-once network delivery.

## Dependencies and approval
Depends on currency, tax, recognition, close, retention and SLO decisions where relevant. Approval required from product owner, qualified accountant, security or SRE according to decision register.

## Validation plan
Create independent fixtures; verify real PostgreSQL constraints/rollback, duplicates, tenant boundaries and report reconciliation as applicable. No validation is claimed complete.
