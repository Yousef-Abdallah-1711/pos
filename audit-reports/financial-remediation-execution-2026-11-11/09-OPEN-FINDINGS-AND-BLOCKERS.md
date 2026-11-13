# 09 — Open findings and blockers

- FIN-001 remains unresolved accounting behavior. Phase 00 detects, warns and suppresses PDF export; it does not add earnings to equity.
- Current equation status uses the existing two-decimal comparison, not an approved currency precision policy; no accountant sign-off.
- Authenticated two-tenant report/export tests now pass for the tested owner and cashier scenarios using actual cookie sessions and RBAC enforcement. This is scoped endpoint evidence, not proof of product-wide tenant isolation; comprehensive RLS coverage remains unestablished.
- Dashboard Vitest suite completed 48 files / 656 tests. The full test typecheck still fails on pre-existing duplicate JSX attributes in `TableList.test.tsx` and readonly `NODE_ENV` writes in `config.test.ts`; neither file changed from the Phase 00 base.
- Arabic page-level RTL passes at 390, 768 and 1440px; English LTR passes at 390px. The document element still has `lang=en` and no `dir`, an app-wide metadata issue.
- FIN-002/003/008/009/010 remediation remains deferred to later phases except safe paid-order diagnostic fields; no durable recovery is added.
- Phase 01 remains NOT READY pending the frontend typecheck gate, FIN-001 accounting decision and correction, owner/accountant policy approvals, and broader release gates.
