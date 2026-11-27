# 53b - UI control interaction results (RUN-004, headed Chrome)

Data: [53b-ui-control-interaction-results.csv](53b-ui-control-interaction-results.csv). Built by `tests/build_control_interaction_results.py` from the matrices produced by `tests/runtime/ph04_r4_control_matrix.mjs` (dashboard) and `tests/runtime/ph13_r4_superadmin_ui.mjs` (Super Admin).

**What this is:** one safe interaction per control (tab, button, link, select, toggle, search/filter), up to 22 controls per route, on 59 of 59 dashboard routes and 4 Super Admin routes; 644 controls in total. **What it is not:** an exhaustive functional test. Destructive or money-moving labels were not clicked (SKIPPED-DESTRUCTIVE: need a fixture and an oracle); the 22-per-route cap leaves lower-priority controls on long pages untested; Arabic and mobile interaction was not repeated (only the load matrix of RUN-003).

Result counts: {'PASS': 274, 'NO-EFFECT': 187, 'FLAG-MUTATION': 18, 'SKIPPED-DESTRUCTIVE': 88, 'ERROR': 9, 'DISABLED': 10, 'BLOCKED': 58}. By app: {'dashboard': 610, 'super-admin': 34}.

Definitions: PASS = observable change without errors; NO-EFFECT = no observable change in the page signature (may still be correct, for example a toggle that re-renders identically); FLAG-MUTATION = the click caused a non-GET request (reviewed individually: exports use POST, Renew/Acknowledge/Mark-all-read/Create Session/Create Sample Template change fixture data, Appearance colours overwrite the active public theme - F-077); ERROR = console error or 5xx (F-009, F-074); BLOCKED = tooling failure (click timeout, Chrome tab crash), not a product verdict; DISABLED = control disabled in the fixture state.

Findings linked: F-009 (menu-profitability and break-even 500 still reproduced through the UI Refresh button), F-074 (unhandled 400 on Create Schedule), F-076 (Super Admin user View links 404), F-077 (Appearance Theme Colors rewrites the active storefront theme).
