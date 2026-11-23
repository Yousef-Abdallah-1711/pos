# CFR-005 tasks — recovery draft

- [x] T001 Trace current paid transition and intent, confirm the unsupported refund guard and merge of #39.
- [ ] T002 Transition, scope and pay/cancel race tests passed in [PR #109 PostgreSQL job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) at `244c44c8f7c53fc81adbd83453edf2875acce9fb`; independently reconcile paid orders to events and GL and revalidate main.
- [ ] T003 Inventory all order channels and provider capture/refund contracts, including durable identity and uncertain-outcome recovery.
- [ ] T004 Obtain accounting decisions for recognition, allocation, tax and closed-period reversal; specify migration, rollback and linked journal model.
- [ ] T005 Implement approved refund lifecycle and test full/partial, duplicate, concurrent, timeout, restart and cross-scope cases with a numerical oracle.
- [ ] T006 Record exact candidate SHA, CI, independent review, merge and end-to-end acceptance evidence.

The completed source trace does not establish runtime or refund acceptance.
