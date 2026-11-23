# CFR-004 tasks — recovery draft

- [x] T001 Trace current atomic journal/GL entry point and paid-order intent transaction; verify merges #33, #35 and #37.
- [x] T002 Run existing failure-boundary, duplicate and scope tests in clean disposable PostgreSQL on a candidate; [PR #109 job](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/45157868007/job/137635542663) passed them at `244c44c8f7c53fc81adbd83453edf2875acce9fb`. Revalidate after integration to main.
- [ ] T003 Run paid transition rollback, retry, expired lease, reconciliation and redrive scenarios; verify a real worker restart and one journal/GL effect with an independent oracle.
- [ ] T004 Exercise scoped operator API and failure visibility, including terminal retry state; define and verify external alerts if required.
- [ ] T005 Reconcile provider settlement scope and policy decisions; implement only authorized missing engineering behavior and rerun relevant checks.
- [ ] T006 Record exact candidate SHA, CI, independent review, merge SHA and accounting approvals; verify on main before closure.

The completed source trace is not a runtime acceptance result.
