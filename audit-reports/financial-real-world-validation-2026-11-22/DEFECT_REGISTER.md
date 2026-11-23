# Defect register

| ID | Finding | Classification | Disposition |
| --- | --- | --- | --- |
| F-01 | Synthetic order Balance Sheet report has assets 5,500 cents, liabilities 500, equity 0; 5,000-cent residual; report marks unreconciled. | Confirmed report-level unreconciled condition; accounting-policy implications unresolved | No adjustment posted. OPEN-01/EARN-01 and presentation approval remain pending. Requires approved policy and further diagnosis. |
| F-02 | Existing #110 10k comparison shows large-scope warm end-to-end p50 regression of 3.8% despite improved query plan and no temp I/O. | Confirmed measurement result; no defect classification without approved budget | Performance acceptance pending; repeat controlled measurements and obtain budget. |
| F-03 | Initial #108 CI failure / intermittent worker test mechanism | Unresolved | Isolated retry passed, but insufficient evidence to establish flake or root cause. |
| F-04 | Initial 10k local benchmark did not measure due missing disposable PostgreSQL `pg_stat_statements` preload | Confirmed test-environment configuration blocker | Disposable DB was corrected; 1k benchmark then ran. Not an application defect. |

No confirmed application defect was remediated in this validation run; no focused fix PR was created. No accounting policy was changed.
