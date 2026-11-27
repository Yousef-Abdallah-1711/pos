# 54b - Field denominator and form coverage (RUN-004)

Source: [07-field-level-inventory.csv](07-field-level-inventory.csv) (627 rows). Detail: [54b-field-denominator-and-form-coverage.csv](54b-field-denominator-and-form-coverage.csv). Built by `tests/build_field_denominator.py`.

## Corrected denominator
- Inventory rows: 627
- Logical user-facing fields: **531** (dashboard 506, storefront 19, super-admin 5, landing 1)
- Duplicate rows (same app, route, form and field name): 88
- Read-only or disabled: 8; generated: 0; hidden: 0
- Distinct forms: 116

## What was and was not executed
- RUN-004 real-UI form matrix (headed Chrome, empty/invalid/Arabic/long/valid modes) exercised **13 of 116 forms** covering **94 of 531 logical fields at form level** (results in 34b EX-RUN004-FORMMATRIX and FORMMATRIX2, evidence in evidence/network/RUN-004/PH-06).
- **Field-level verification was NOT performed**: the driver fills every visible field of a form with a mode value and observes the form outcome; it does not isolate one field at a time, test min/max per numeric field, date ordering per pair, file types or sizes, or decimal precision. These rows are therefore labelled FORM-LEVEL EXERCISED, never PASS.
- 437 logical fields belong to forms not exercised (theme builder, settings sub-forms, storefront checkout, super-admin, POS, accounting invoices, reports filters and others) and stay NOT TESTED.
- Per-field boundary results that exist come from the API batteries summarised in 54-field-validation-and-boundary-results.md (RUN-003).

Forms exercised by the matrix that are ABSENT from the inventory (inventory incompleteness): driver create page (12 fields), settings users create dialog (4 fields), RBAC role create dialog (2 fields); the employee form exposes 6 fields in the UI against 1 inventoried. Add them before the next inventory refresh.

Findings from the matrix: F-074, F-075, F-079 (plus reconfirmation of F-014, F-031, F-062, F-064, F-065).
