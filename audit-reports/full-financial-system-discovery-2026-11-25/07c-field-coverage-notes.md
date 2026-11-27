# 07c — Field inventory coverage notes (honest)

`07-field-level-inventory.csv` now has **627 rows** (original manual rows retained and unchanged except for added columns; original backed up at `evidence/07-field-level-inventory.original.csv`). Earlier packages said "200 rows"; the original file actually held 201 data rows (header excluded) — the earlier count was off by one.

| Source | Rows |
|---|---|
| Original manual rows (high fidelity: form name, validation, API mapping) | 201 |
| Added by a **regex extraction** pass over JSX (`<input|select|textarea|Checkbox|Switch|DatePicker|TextField>`) in files not already represented | 426 |
| By app (total) | dashboard 596 · restaurant-website 24 · super-admin 5 · landing 2 |
| By kind | user-editable 619 · read-only/disabled 8 · hidden 0 detected |

`07b-form-inventory.csv`: **128 input-bearing component/page files** (dashboard 108, super-admin 9, restaurant-website 7, landing 4); 74 of them were not represented in the original rows. These are **file-level containers, not logical forms** — the true number of forms is not established.

## What the added rows do NOT provide
Label text, default values, validation rules, conditional visibility, dropdown data source, API payload key, backend validation, DB destination: columns are blank or `none visible in tag`. The extraction captures `name`/`id`/`placeholder`/`type`/`required`/`min`/`max`/`pattern`/`maxLength` attributes only; calculated and hidden values are not derived.

## Known misses (coverage NOT achieved)
- Custom wrapper components (e.g. `FormField`, `FormInput`, controlled fields that spread `register(...)`, shadcn `Input` wrappers) are only caught where the tag name matches; **landing (2 rows) and super-admin (5 rows) are almost certainly under-counted** — login/registration/password-recovery forms must be re-read manually.
- UI primitive definitions (`components/ui/input.tsx`, `textarea.tsx`) appear as `UNNAMED` rows; they are not forms.
- Dialogs/drawers/tabs/multi-step flows are detected only if the input tag sits in the scanned file; inputs rendered via data-driven field arrays are missed.
- No browser crawl was done (that would need form interaction / navigation beyond read-only inspection and is deferred to execution).

Per the acceptance criterion, **field-level discovery is NOT exhaustive**. A manual component-by-component read of EmployeeForm, HR forms, auth, RBAC, super-admin and settings forms remains required (tracked as KG-ORC-001 PARTIALLY-RESOLVED).
