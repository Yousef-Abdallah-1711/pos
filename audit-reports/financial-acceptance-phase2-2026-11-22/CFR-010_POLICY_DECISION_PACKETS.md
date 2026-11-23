# CFR-010 Accounting Policy Decision Packets

**Status:** Draft for qualified review; no approval recorded
**Prepared:** 2026-11-22
**Required disposition:** Keep CFR-010 OPEN and financial production NO-GO until original-source reconciliation, qualified accounting approval, and the applicable acceptance gates are complete.
**Original operating root cause:** NOT ESTABLISHED

This packet requests decisions that engineering must not infer. Synthetic tests and reports establish software behavior against fixtures; they do not establish SERVIO's original balances or accounting policy. Each decision below is unapproved until an authorized person completes and signs its approval record.

## Evidence and limits reviewed

- PR [#104](https://github.com/Yousef-Abdallah-1711/pos/pull/104) merged as `32f79c55f1e1643aafe09b6d9d93be8e62fb4b96`.
- PR [#108](https://github.com/Yousef-Abdallah-1711/pos/pull/108), head `1a9a949ae59a303444c02670cb0209ab61b046cf`: read-only opening-balance and cross-report diagnostics. Financial PostgreSQL acceptance and application-image smoke passed in workflow run [44984726238](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/44984726238) after retry. Auth-browser acceptance was skipped.
- PR [#109](https://github.com/Yousef-Abdallah-1711/pos/pull/109), head `8b4831576a664a1ee82df5ac5ed59a2b904f9726`: read-only current-earnings diagnostic. CI and financial acceptance passed; auth-browser acceptance was skipped.
- PR [#110](https://github.com/Yousef-Abdallah-1711/pos/pull/110), head `3592dd1a50837d3827d8677f50480c3566243c46`: reconciliation benchmark evidence. Benchmark, financial acceptance, and application-image smoke passed; auth-browser acceptance was skipped.
- Synthetic benchmark workload: 10,000 journals and 50,000 journal-line/GL rows. Candidate warm p50/p95/max: 180.726/182.565/182.660 ms; baseline: 174.044/175.441/175.524 ms. The candidate plan improved and removed baseline temp I/O, while measured end-to-end latency was slower. No approved latency budget was identified; this is not performance acceptance.
- Synthetic report fixtures include the 5,000-cent residual scenario. No original tenant/restaurant identifiers, signed opening Trial Balance, source bank statements, original journals/GL, original reports, fiscal-close state, or original producing environment/SHA were available in the evidence reviewed for this packet.

The PRs remain engineering candidates. Their checks and independent code reviews do not constitute accounting approval. No independent GitHub review is recorded on PRs #108–#110 as of this packet.

## OPEN-01 — Opening balance source, date, and offset

### Decision requested

Approve the authoritative source documents, cutover date, account mapping, and treatment of opening balances for the specifically identified tenant, restaurant, currency, and fiscal period. Confirm whether opening entries may be posted from those sources, whether balances may only be imported for reporting, or whether the matter must remain deferred.

### Options for the approver

- [ ] **A — Source-backed opening entries.** Approve only the opening amounts and accounts reconciled to the signed source Trial Balance and supporting statements listed in the decision record. Specify the journal date, offset account(s), references, and review controls.
- [ ] **B — Source-backed reporting import, no journal.** Approve a clearly identified reporting-only opening source, its accounts and cutover date, and how it reconciles to the GL. Do not present it as a posted journal.
- [ ] **C — Defer unsupported balances (recommended with current evidence).** Keep the diagnostic warning and report residual visible. Do not synthesize opening entries or a balancing plug. Reconsider when original evidence is supplied.

### Implications and guardrails

The selected source, date, account mapping, and debit/credit treatment determine whether an opening-balance difference is a real omission or an unsupported inference. A decision cannot validate amounts that are absent from the approved source documents. Engineering will not post or rewrite balances as a result of this packet.

### Scope and evidence to attach

- Tenant ID: **not supplied**
- Restaurant ID(s): **not supplied**
- Currency and fiscal period: **not supplied**
- Cutover date: **not established**
- Signed opening Trial Balance and version/date: **not supplied**
- Source cash/bank statements and reconciliation: **not supplied**
- Approved chart-of-accounts mapping and offset account(s): **not supplied**
- Original journal/GL references and producing environment/SHA: **not supplied**

### Approval record — OPEN-01

- Selected option and any conditions: ____________________________________
- Approved source documents / immutable references: ______________________
- Tenant / restaurant / currency / period scope: _________________________
- Cutover date and approved account/offset mapping: ______________________
- Effective date (`YYYY-MM-DD`): _________________________________________
- Qualified accountant/controller name and title: ________________________
- Signature or recorded approval reference: _____________________________
- Approval date (`YYYY-MM-DD`): _________________________________________

## EARN-01 — Presentation of unclosed current and prior earnings

### Decision requested

For each approved fiscal-period boundary, decide whether candidate unclosed earnings are displayed separately, included in equity presentation without a journal, or transferred through an approved close journal. Specify treatment of prior-period earnings, existing retained-earnings movements, existing close journals, and reversals.

### Options for the approver

- [ ] **A — Present candidate earnings separately, without posting (recommended interim choice).** Keep the calculation visibly distinct from posted equity and state that it is not a close entry.
- [ ] **B — Include in Balance Sheet equity presentation without posting.** Define how the report labels the amount, applies period boundaries, and prevents double counting with existing retained-earnings or close movements.
- [ ] **C — Transfer only through an approved period-close journal.** Specify the accounts, date, idempotency, reversal/correction treatment, and the approved close process under CLOSE-01.
- [ ] **D — Defer the presentation decision.** Preserve the policy-pending status and avoid changing production Balance Sheet treatment.

### Implications and guardrails

The read-only diagnostic derives candidate net earnings from scoped posted GL activity. A candidate amount is not proof that a close is due, that an existing close is absent, or that an amount belongs in retained earnings. No computed amount will be injected into the production Balance Sheet or posted to the GL without an approved policy and separate accepted implementation.

### Scope and evidence to attach

- Tenant / restaurant / currency: **not supplied**
- Fiscal year and inclusive/exclusive period boundaries: **not approved**
- Revenue and expense account classifications: **existing data only; policy confirmation pending**
- Existing retained-earnings and close journal references: **not reconciled to original records**
- Reversal and closed-period treatment: **not approved**
- Synthetic reference case: **5,000 cents; not original operating evidence**

### Approval record — EARN-01

- Selected option and any conditions: ____________________________________
- Period boundary and account-classification scope: ______________________
- Treatment of existing retained earnings / close journals: ______________
- Reversal, prior-period, and double-counting rules: _____________________
- Effective date (`YYYY-MM-DD`): _________________________________________
- Qualified accountant/controller name and title: ________________________
- Signature or recorded approval reference: _____________________________
- Approval date (`YYYY-MM-DD`): _________________________________________

## CLOSE-01 — Fiscal close, adjustment periods, reopen, and correction

### Decision requested

Define the fiscal calendar, close states, adjustment-period rules, authorized actors, reopen conditions, and immutable correction lifecycle. Confirm whether close/reopen behavior may be implemented and for which tenant and fiscal periods.

### Options for the approver

- [ ] **A — Controlled hard close with a defined adjustment period.** Set the close date, allowed adjustment dates, required approvals, retained-earnings treatment, and correction path.
- [ ] **B — Soft close followed by controlled finalization.** Define which reports are provisional, when finalization occurs, and how later corrections are recorded and disclosed.
- [ ] **C — Defer close/reopen activation (recommended with current evidence).** Keep posting/period behavior unchanged; do not lock, reopen, or create closing journals until the lifecycle and original fiscal state are established.

### Implications and guardrails

Close policy controls posting eligibility, retained-earnings roll-forward, corrections, and audit history. A reopen decision must name authorized roles and approval requirements. Corrections must preserve posted history and source references; no in-place journal edits or unapproved fiscal locks are authorized by this packet.

### Scope and evidence to attach

- Tenant / restaurant / fiscal year: **not supplied**
- Fiscal calendar and close date: **not supplied**
- Current original close state and adjustment-period records: **not supplied**
- Roles allowed to close, adjust, and reopen: **not approved**
- Approval count, audit events, and correction journal rules: **not approved**
- Effective date and any prospective/retroactive application: **not specified**

### Approval record — CLOSE-01

- Selected option and any conditions: ____________________________________
- Fiscal calendar / close date / adjustment periods: _____________________
- Authorized roles and required approvals: ______________________________
- Reopen conditions and correction lifecycle: ___________________________
- Retained-earnings treatment and journal rules: ________________________
- Tenant / restaurant / fiscal-year scope: ______________________________
- Effective date (`YYYY-MM-DD`): _________________________________________
- Qualified accountant/controller name and title: ________________________
- Product owner name and title (required for system lifecycle decisions): _
- Signatures or recorded approval references: ___________________________
- Approval date (`YYYY-MM-DD`): _________________________________________

## Interim operating disposition

Until each applicable approval is completed and its evidence is attached:

1. Keep CFR-010 OPEN and financial production NO-GO.
2. Preserve `ORIGINAL OPERATING ROOT CAUSE: NOT ESTABLISHED`.
3. Keep opening-balance and earnings diagnostics read-only and policy-neutral.
4. Do not create opening or closing journals, balancing plugs, fiscal locks, or reopen transitions.
5. Do not treat synthetic tests, benchmark results, or passing CI as accounting approval.
6. Record each final approval with approver identity, authority, scope, effective date, source-document references, and the implementation/acceptance evidence required before release.

**Decision register status:** OPEN-01 pending; EARN-01 pending; CLOSE-01 pending. No approval signature or effective date is present in this draft.
