# Remaining blockers

| ID | Impact | Owner | Required action / evidence | Independent work that can continue |
|---|---|---|---|---|
| CFR010-DB-01 | Journal/GL integration behavior is not verified in PostgreSQL. | Engineering / CI owner | Provide an isolated disposable PostgreSQL service with credentials for a database whose name includes `test`; rerun the targeted tests and complete acceptance on exact SHA. | Code review, source recovery through authorized bundles, policy-independent design. |
| CFR010-SRC-01 | Original discrepancy cannot be explained or closed. | Finance data owner | Supply provenance-preserving original statements, signed opening TB, scoped journals/GL, report payload/export and fiscal close state. | Read-only integrity tooling and synthetic-only tests. |
| CFR010-POLICY-01 | Earnings/opening/close production behavior cannot be changed safely. | Qualified accountant/controller; fiscal calendar owner for close | Record signed OPEN-01, EARN-01 and CLOSE-01 decisions with scope and effective date. | Detection, test infrastructure and implementation contracts that do not alter posted results. |
| CFR010-ACCEPT-01 | Financial release cannot be accepted. | Independent finance/engineering reviewers and CI/browser owners | Exact-SHA database, authenticated browser, security, export parity, performance, CI and independent review evidence. | Continue other independent engineering work; production remains NO-GO. |
