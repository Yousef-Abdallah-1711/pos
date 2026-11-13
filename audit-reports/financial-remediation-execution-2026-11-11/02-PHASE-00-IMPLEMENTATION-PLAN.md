# 02 — Phase 00 implementation plan

Task PH00-001: add API integrity status and equation difference based only on existing equationBalanced result; preserve response figures and legacy field. UI displays accessible warning while preserving numbers and suppresses PDF export for anything other than explicit RECONCILED plus equation_balanced true. Task PH00-002: make existing paid-order posting failures safer to identify only if existing log permits; never claim durable retry. Task PH00-003: deterministic usecase/component/export-gate tests. Task PH00-004: operator runbook.

No migration planned. Policy blocker: no approved currency precision; use existing legacy round2 rule solely for continuity, disclose basis and non-certification.