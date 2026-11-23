# CFR-011 tasks — recovery draft

- [x] T001 Verify migration 187, protected rollback, #44 merge and current absence of production callers for the two hook names.
- [x] T002 Verify invoice/payment constructed-entry persistence and GL tests passed in [disposable PostgreSQL](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) at candidate `244c44c8f7c53fc81adbd83453edf2875acce9fb`; revalidate final main.
- [ ] T003 Test migration upgrade on existing rows and rollback refusal with long values on a disposable database.
- [ ] T004 Establish approved supplier hook contract, source lifecycle and account mapping; confirm support or supersession.
- [ ] T005 If supported, test and reconcile actual invoice/payment, failures, reversals, tenant scope, report filters and exports.
- [ ] T006 Record final main, CI, review and accounting evidence before closure.

The passing constructed-entry tests do not establish supplier lifecycle acceptance.
