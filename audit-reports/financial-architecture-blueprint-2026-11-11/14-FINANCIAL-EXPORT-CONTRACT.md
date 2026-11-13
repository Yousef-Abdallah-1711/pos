# 14 — Financial export contract (PROPOSED)

Audit evidence: PDF partial; CSV/XLSX untested; native XLSX is not established. CSV is not XLSX. Export shares report cutoff, filters, precision, authorization and watermark. Manifest includes generated time, tenant/restaurant, report, period, currency, policy version, watermark, filters and row count. Amounts are exact decimal strings.

Use streaming, row limits, cancellation, memory limits, request/download authorization, scoped short-lived links and audit provenance. Prevent formula injection for spreadsheet-leading =,+,-,@, tab/newline without altering canonical data. Native XLSX requires real workbook and reader tests. Compare totals independently to screen/ledger; measure peak memory. PDF tests signs, precision, RTL, pagination/totals.
