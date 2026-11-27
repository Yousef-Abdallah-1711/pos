# 30d - Metric verification classes (RUN-005)

Built by `tests/build_metric_verification_classes.py`; the class is stored per row in the new `verification_class` column of [30c](30c-metric-coverage-reconciliation.csv).

Total rows: 119. Class counts: {'1 independently verified and correct': 61, '2 independently verified and incorrect': 23, '3 measured without sufficient independent verification': 17, '4 attempted but blocked': 12, '5 policy dependent': 6}.

Relationship to the RUN-004 numbers (coverage_class column): {'MEASURED vs independent value': 26, 'MEASURED vs independent value (RUN-003)': 56, 'MEASURED vs independent value (RUN-004, headed Chrome + SQL)': 6, 'MEASURED vs independent value (RUN-003, partial: boundaries untested)': 11, 'ATTEMPTED - BLOCKED BY DEFECT F-009 (HTTP 500)': 12, 'MEASURED (vacuous: no data; flow pending)': 3, 'MEASURED (not in 30 inventory: added in 43)': 5}. In RUN-004 terms: 99 rows were MEASURED vs an independent value (56 + 26 + 11 + 6), 8 were measured without one (5 added after the inventory + 3 vacuous) and 12 were blocked by F-009 (menu profitability and break-even): 99 + 8 + 12 = 119. "Measured vs independent value" is NOT the same as "correct": the correct/incorrect split above comes from the runtime status in 43 (FAIL, INCONSISTENT and MISLEADING map to class 2).

Class definitions: 1 independently verified and correct; 2 independently verified and incorrect (a finding exists); 3 measured without sufficient independent verification (vacuous data, untested boundaries, one populated bucket, or not verified); 4 attempted but blocked by a defect; 5 policy dependent (the expected value depends on an unapproved decision in 56); 6 missing capability and 7 not applicable (none among the 119 rows).

Matching the application's own SQL is not accepted as proof of accounting correctness; rows whose only comparison is the application query are class 3.
