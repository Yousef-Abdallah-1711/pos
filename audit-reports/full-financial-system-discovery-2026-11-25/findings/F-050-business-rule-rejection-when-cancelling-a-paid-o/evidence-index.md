# F-050 — Evidence index

- Register evidence status: **REPRODUCED**
- Evidence verification (file 55): see 55-finding-evidence-and-root-cause-verification.csv
- Executed tests: EX-RUN003-AUD-HC10, EX-RUN003-AUD-SAAS-006, EX-RUN003-AUD-OT-003, EX-RUN005-DESTR-ORDER-CANCEL

## Evidence files (relative to the audit folder)
| path | status |
|---|---|
| evidence/database/RUN-003/PH-19/AUD-REFUND-01-paid-order-cancel.json | exists |
| evidence/database/RUN-003/PH-19/AUD-REG-HC10-FIN019-CFR006.json | exists |
| tests/runtime/ph19_paid_cancel.mjs | exists |
| tests/runtime/ph19_regs2.mjs | exists |
| evidence/database/RUN-003/PH-07/AUD-ORDER-TYPES.json | exists |
| tests/runtime/ph07_order_types.mjs | exists |
| evidence/network/RUN-005/PH-04/AUD-R5-DESTRUCTIVE-A.run1.json | exists |
| tests/runtime/ph04_r5_destructive_a.mjs | exists |

State files `evidence/state-*.json` are session cookies and are intentionally not referenced.
