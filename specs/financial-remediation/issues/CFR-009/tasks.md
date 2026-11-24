# CFR-009 tasks — recovery draft

- [x] T001 Trace paid mutation guards, locked payment and shared posting transaction; verify #40 merge and historical test inventory.
- [x] T002 Seven CFR-009 PostgreSQL scenarios passed at candidate `244c44c8f7c53fc81adbd83453edf2875acce9fb` in the [disposable PostgreSQL job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663), including failure, concurrency and scope. Revalidate on final main.
- [ ] T003 Reconcile synthetic payroll reports and GL independently, including authorized deductions and tenant/restaurant isolation.
- [ ] T004 Obtain correction, reversal and statutory decisions; design audit-preserving migration and rollback if needed.
- [ ] T005 Implement approved correction and test one-time concurrent payment, reversal, authorization, API and browser flows.
- [ ] T006 Record CI, independent review, merge and exact-main acceptance before closure.

Source tracing does not establish correction or runtime acceptance.
