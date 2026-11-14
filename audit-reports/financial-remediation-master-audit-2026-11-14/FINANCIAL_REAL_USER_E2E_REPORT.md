# Financial Real-User E2E Report
> **Phase 2 execution update (2026-11-14):** `origin/main` was fetched and verified at `165837c7ad2e506dbf20b1a45c9a123476f2f8e9`. Focused CFR-010 A–D PostgreSQL fixtures pass against a newly provisioned schema-only disposable database; independent SQL confirms Fixture A balanced opening (assets/equity 1,000/1,000), Fixture B off-ledger opening reproduces the historical -100 EGP residual, Fixture C expense and Fixture D two-period/close cases reconcile to their independent figures. Migrations `186_accounting_event_redrive_audit.sql` (PR #37), `186_servio_trial_offer.sql` (main), and `187_expand_journal_entry_type_to_32.sql` (PR #44) applied together in fresh and upgraded disposable schemas; long `auto_supplier_*` values persisted, and unsafe rollback refused. This does **not** execute authenticated supplier auto-posting: the only production references to those methods/constants are definitions, with tests as the only callers. Manual browser acceptance remains blocked; no credentialed synthetic browser session is available. PR #46 adds the regression and its CI is green on all required jobs. Financial production remains **NO-GO**.

> **Execution update (2026-11-14):** A later audit pass ran a verified local disposable stack: 1 authenticated supplier invoice/payment API workflow, 6 automated API-backed report views, 1 unauthenticated authorization-negative request, and 0 manual financial browser workflows. Independent SQL reconciliation found source/payment/GL/TB/P&L/AP/AR differences of 0.00 for the narrow fixture, but a -100.00 EGP Balance Sheet residual because opening cash is absent from GL. Local isolation PASS for observed controls; existing staging isolation UNKNOWN. Financial release remains **NO-GO**. Detailed scope, values, commands, and limitations: [disposable execution evidence](evidence/disposable-execution-2026-11-14.md). This update supersedes earlier statements of no local fixture where applicable.

**Targets:** public staging for read-only release identity and login-page observation only, plus a separate local disposable Compose stack for synthetic authenticated API and automated report execution. No staging financial writes were made. Local isolation passed for the observed disposable configuration; staging isolation remains unknown.

## Environment preflight

| Check | Observed | Acceptance status |
|---|---|---|
| Current main SHA | `018682beb9448eea6dd72fd2418876b439a94dbf` | Recorded. |
| Staging release identity | Public `/_ops/release` returned `environment=staging`, `status=candidate`, SHA `f733efb5e7f5730fea5a56a65d66d9c339e6c040`. | Does not match current main and is not `verified`. |
| API health | HTTP 200; `status=ok`, `db_healthy=true`. | API and DB connectivity observed only. |
| Dashboard route | Opening `/en/dashboard` redirected to the sign-in form at `/en/auth/login`. | No authenticated test session. |
| PostgreSQL isolation | Staging runtime config was not exposed; separate local disposable PG 15.19 used tmpfs database `pos_phase2a_e2e`, loopback port only. | Staging **UNKNOWN**; local observed controls **PASS**. |
| Redis isolation | Staging Redis unverified; local stack used a dedicated Redis container. | Staging **UNKNOWN**; local observed controls **PASS**. |
| Payment/email/other integrations | Staging config unverified; local Paymob mock mode, relevant external credential variables unset. | Staging **UNKNOWN**; local observed controls **PASS**. |
| Synthetic-only account and record guarantee | No evidence available. | **Unverified.** |
| Production database/credentials unreachable | No configuration or network-isolation evidence available to this audit. | **Unverified.** |

No SSH session, environment secret, connection string, session cookie, or password was inspected or disclosed. The local app inventory showed no pre-authenticated staging browser session. No approved secret manager/credential-broker connector was available in this task’s tools. The login form was not submitted.

## Manual scenario record

“Manual browser” means visible human-style UI interaction; it is not a Playwright assertion. Only the read-only route/load observation was performed. No screenshots or traces of financial actions exist because no financial action was attempted.

| Scenario | Actions performed | Expected | Observed | Verdict |
|---|---|---|---|---|
| Preflight / sign-in | Opened `https://pos.app.rakmyat.online/en/dashboard` in a fresh visible browser tab. | Correct staging build opens; authorized synthetic session is available. | Redirected to `https://pos.rakmyat.online/en/auth/login`; sign-in form visible; no session available. Release endpoint reports older `candidate` SHA. | **Blocked** |
| 1 — Restaurant owner | None beyond preflight; no login, sale, payment, report, or export. | Persisted synthetic sale ties through journal, GL, P&L, BS and reconciliation. | No authorized isolated data or login. | **Blocked** |
| 2 — Accountant | None. | Valid journal posts; unbalanced rejected; posted mutation denied; TB and export tie. | No session; no journal entered. | **Blocked** |
| 3 — Payroll | None. | Unpaid salary edits; paid state has one posting; repeats and edits/deletes rejected. | No session; no salary entered. | **Blocked** |
| 4 — Supplier accounting | None. | Invoice and payment persist with journal/GL and AP tie. | No session; this mainline path also has the CFR-011 schema mismatch. | **Blocked** |
| 5 — Security | None. | Unauthorized/cross-tenant/forbidden export denied at API and UI. | No test roles or verified isolated tenant data. | **Blocked** |
| 6 — Recovery | None. | Durable event survives controlled failure/restart, replays once and is diagnosable. | No isolated worker/DB controls available. | **Blocked** |
| 7 — Browser quality | Sign-in route only. | Arabic/English, RTL/LTR, mobile, validation, empty/error/loading states, console/network and export checked. | No financial page or authenticated locale dashboard reached. | **Blocked** |

## Automated and API evidence separation

### Later local disposable execution

One authenticated synthetic supplier invoice/payment API workflow persisted 100.00 EGP. Six Playwright report routes loaded from the real API/backend: English Balance Sheet, P&L, Cash Flow, AP aging, AR aging, and Arabic Balance Sheet. Console errors/warnings were 0/0. One unauthenticated Balance Sheet API GET returned 401. At 375×812, document width measured 422px, an observed 47px overflow from the fixed off-canvas sidebar. These are automated browser/API observations, not manual browser acceptance. Full amounts and isolation evidence are in the linked execution record.


- GitHub Actions PR checks are automated. The CI workflow uses disposable PostgreSQL/Redis and runs full tagged Go integration tests, frontend tests for four apps, security/config checks, and an image-based Chromium smoke. This audit did not treat that smoke as manual acceptance.
- The staging health check is a real backend/DB health observation, not a financial API transaction.
- One authenticated supplier invoice/payment workflow and six API-backed report views were run in the later local disposable stack. The initial staging browser observation remained read-only. No manual financial browser workflow, export download, or staging financial write was performed.

## Required access to resume

Provide either an already authenticated, synthetic staging test session or an approved secret manager/credential broker path, plus auditable confirmation that the exact staging release uses a dedicated disposable PostgreSQL database, synthetic accounts/records, isolated Redis, and sandboxed/disabled external integrations with no route to production data or payment credentials. Do not send passwords or tokens in chat. Then verify the deployment SHA is current and `verified` before beginning write-side scenarios.


## Phase 2 latest-main execution boundary

No new authenticated user workflow ran in this continuation. New execution consists of two focused database integration tests (CFR-010 A-D fixtures), independent SQL reconciliation, and migration compatibility checks. These counts are database fixtures, not API users or browser scenarios. Earlier candidate-stack report views/workflow remain clearly identified as historical candidate evidence and are not promoted to latest-main passes. The local service identity check reconfirmed tmpfs-backed disposable PostgreSQL, dedicated Redis, mock Paymob, restricted app DB role, and unset inspected external-provider credentials; generic network egress remains unproven. Dashboard browser probe returned the landing app at port 13111 and 404s on the other web ports; it is not acceptance evidence. Manual financial browser workflow count remains **0 / BLOCKED**.


