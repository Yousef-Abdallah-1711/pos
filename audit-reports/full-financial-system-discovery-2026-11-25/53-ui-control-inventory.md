# 53 — UI control inventory (RUN-003, headed Chrome, EN desktop; read-only)

Data: [53-ui-control-inventory.csv](53-ui-control-inventory.csv) from `tests/runtime/ph04_controls.mjs` (evidence `evidence/network/RUN-003/PH-04/AUD-UI-CONTROLS.json`) merged with the language/viewport matrix `tests/runtime/ph04_ui_matrix.mjs`.

**What this is:** the count and sample labels of VISIBLE rendered controls per dashboard route at load time (59 routes). **What it is not:** a per-control functional test. Controls were not clicked; dialogs/tabs that appear only after interaction are not enumerated. Pages needing data (details, edit forms) and Super Admin / landing / storefront apps are outside this table (their routes were loaded in earlier batches and in the workflow scripts).

Totals over the 59 routes: buttons 1175, links 503, inputs 112, tables 42, forms 1.

Status per row (last column): INTERACTED = a workflow script exercised the page (see 34b); LOADED = control counts only; NOT-RENDERED = empty page (environment issue) and must be re-run.
