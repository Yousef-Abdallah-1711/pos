# Fiscal Period Boundaries

The Balance Sheet query includes GL movements through the supplied `as_of` date (`gl.entry_date <= as_of`). The P&L query uses inclusive `entry_date >= from` and `entry_date <= to` boundaries. A disposable PostgreSQL test verifies an entry dated after an earlier `as_of` date is excluded from the earlier report and included in the later report.

The synthetic earnings fixture spans 2025 activity, an explicit 2026 close, and new 2026 activity. These tests validate date filtering and explicit close behavior; they do not establish an organization-approved fiscal calendar or closing policy.

Specific first/last-day fiscal-year behavior, timezone cutoff policy, and authorized fiscal-period definitions remain unverified and policy-dependent. No production fiscal-period configuration was changed.
