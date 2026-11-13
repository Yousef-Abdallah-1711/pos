# 05 — Financial integrity evidence

Historical audit independently observed FIN-001 on four fixtures: Balance Sheet equation gaps 100, 60, 60 and 60; P&L net profit matched 100/60/60/60. Phase 00 does not fill gaps.

Current Phase 00 test: balanced sale journal Dr cash 100 / Cr sales revenue 100, stored in isolated PostgreSQL. Since Balance Sheet includes only asset/liability/equity, API usecase reports assets 100, liabilities 0, equity 0, difference 100 and UNRECONCILED; amounts remain unchanged. SQLMock covers balanced 100/100 assets/equity and reports RECONCILED. The balance status uses legacy `round2` only; it is not accountant-approved and is explicitly identified in integrity_check_basis. This is arithmetic detection, not a certification of full statement correctness.