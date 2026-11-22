# CFR-002 ordered tasks

- [x] T001 Identify merged PR #34, migration 180, historical evidence and current-main source.
- [ ] T002 Create a disposable PostgreSQL database and runtime roles; record version, migration state and exact code SHA.
- [ ] T003 Test debit-heavy, credit-heavy and zero-line direct SQL posting; verify rejection leaves the journal draft and no ledger effects.
- [ ] T004 Verify balanced posting, immutable header/lines/ledger, duplicate effects, cross-tenant and restaurant scope changes under appropriate roles.
- [ ] T005 Race six postings of one draft; independently assert one original ledger effect per line, then verify void and duplicate void behavior.
- [ ] T006 Run synthetic preflight, invalid-data migration rejection, documented downgrade and forward migration without altering financial rows.
- [ ] T007 Obtain recorded accounting decisions for historical correction and linked/closed-period reversal; coordinate CFR-004, CFR-005 and CFR-013.
- [ ] T008 Record current-SHA results, role boundaries, independent review and outstanding financial release gates before changing the issue status.
