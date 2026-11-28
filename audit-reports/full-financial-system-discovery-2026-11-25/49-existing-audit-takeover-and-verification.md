# 49 — Existing audit takeover and verification (Master Prompt 05, RUN-003, 2026-11-26)

Authoritative progress tracker remains [progress.md](progress.md). This file records what was independently re-verified at takeover and what can be reused.

## 1. Baseline verification (measured 2026-11-26 16:22Z)
| Item | Prompt / historical claim | Measured now | Verdict |
|---|---|---|---|
| Git HEAD | `133260db66b1143c254d79476daa7650b41a33f0` | `133260db66b1143c254d79476daa7650b41a33f0` | UNCHANGED — historical evidence tied to the code is still applicable |
| Repository | `Yousef-Abdallah-1711/pos` | local `origin` = `https://github.com/yousefabdallah171/POSS.git` | DISCREPANCY (already noted in progress.md: local `origin` is a different/obsolete remote). Not changed (no git/remote edits allowed). |
| Files under audit dir | "≈350" | **331** at takeover (60 top-level, not 63; 177 evidence files confirmed; 27 diagrams; 4 oracle; 13 SQL; 42 runtime scripts + 8 other tests) | PROMPT COUNT WAS APPROXIMATE; record real counts |
| Register sizes | 627 fields / 374 APIs / 1,491 tests / 36 findings / 112 metrics / 36 runtime metric rows / 936 executed-result rows | 627 / 374 / 1,491 / 36 / 112 / 36 / 936 | CONFIRMED (CSV row counts) |
| Findings by severity (file 22) | "1 Critical, 11 High…" in progress.md batch 1 | 36 rows: Critical 1, High 15, Medium 12, Low 8 | progress.md batch-1 numbers were for 26 findings; current register is 36 |
| Evidence references | — | 24 distinct `evidence/…` paths cited by the md/csv registers; **0 missing** | OK (does not prove every register row cites evidence — see §3) |
| Working tree | 6 tooling-modified frontend files | same 6 (`tsconfig.json` ×4, `next-env.d.ts` ×2) — written by Next dev servers in containers via bind mount, not by the audit | UNCHANGED; **no product source modified by this run** |
| Stashes | 2 | 2 (`stash@{0}`, `stash@{1}`) | intact |

## 2. Environment verification (RUN-003)
- Dev DB fingerprint before changes: `evidence/runtime/RUN-003/46-dev-db-fingerprint-before.txt` (db `pos_saas`: tenants 2, users 6, orders 0, journals 0, COA 62, migrations 145) — identical to the previous fingerprints.
- Dev app containers (`pos-landing/dashboard/restaurant-website/super-admin/backend`) stopped (not removed); `pos-postgres` untouched.
- Isolated stack `pos-p2a-ws05-e2e` started on the **persistent** volume `pos-p2a-ws05-e2e_p2a_audit_pgdata` — **no restore was needed**; row counts equal the batch-2 snapshot (tenants 4, orders 34, journals 54, GL lines 115, migrations 145, RLS tables 20); roles `pos_app` / `pos_accounting` non-superuser, no BYPASSRLS; `ENV=phase2a-e2e`, `RBAC_ENFORCEMENT=enforce`. Evidence: `evidence/runtime/RUN-003/47-isolated-gate-and-restore-check.txt`.
- **Runbook defect found and corrected (audit tooling only):** `tests/runtime/docker-compose.audit-persistent.yml` uses `./backend` / `./frontend` relative bind mounts. Because the file lives in `tests/runtime/`, `docker compose -f …` resolves them to the wrong folder and `init-roles` fails with `/fixture/create_runtime_roles.sql: No such file`. **Always pass `--project-directory <repo root>`.** (The repo compose file is untouched.)
- Super Admin app serves at `/` and `/login` (no `/en` prefix); `/en` returns 404 — expected, not a defect.
- Real Chrome: installed Google Chrome 150.0.7871.124, headed, dedicated profile outside the evidence tree, logged in through the real UI (`tests/runtime/chrome_headed_smoke.mjs`; evidence `evidence/screenshots/RUN-003/PH-02/AUD-ENV-001-*.png`, `evidence/network/RUN-003/PH-02/AUD-ENV-001.json`). Earlier batches used Playwright **headless** Chromium (`lib.mjs`: `chromium.launch({headless:true})`) — they did **not** meet this prompt's "headed Chrome" requirement; headed Chrome is used from RUN-003 on.
- Observation (not yet a finding): `GET /auth/session` returned 401 on the login page before authentication. Expected for an anonymous probe; recorded so a later security pass can confirm it is the only 401.

## 3. Classification of previous results
| Class | Content | Reuse decision |
|---|---|---|
| Reusable as-is (code unchanged, raw evidence exists) | Static architecture/lineage docs 02–13, 25–28, 33b; DB schema facts (102 tables/287 FK/47 triggers/20 RLS); oracle files (30 self-checks — *not* app tests); batch-1/2 runtime evidence 00–45 for F-001…F-036 as reported; API permission sweep (366 routes) | Reuse; do not rerun |
| Reusable but weak (page load / permission-only) | 117 page loads (≈19% real interaction), 366 anonymous + 366 Employee API checks (permission status only, not business logic) | Do NOT count as functional coverage (stated in 45) |
| Used SQL fixtures, not UI | Tenant C COA/posting rules copied from tenant 1; employee created by API diagnostic fixture | Valid only as diagnostic fixtures; do not present as end-to-end UI success |
| Requires revalidation | Historical regressions: 9 REPRODUCED and 6 PARTIAL (per 45, from single batch runs); 8 historical issues never run (CFR-004, CFR-006, FIN-016, FIN-019, FIN-020, HC-04, HC-05, HC-10); report metrics (43 holds 36 measured rows that are not 1:1 with the 112 inventory rows, so the exact number of inventory metrics still lacking a runtime comparison is UNKNOWN until 30↔43 is reconciled — T-P12-001) | Scheduled in 51 |
| Never executed | Operating month (SCN-24); 3 scenarios; 8 security/reliability tests; HR sub-modules, theme builder, settings forms, advanced accounting console, import/export, responsive/RTL matrix, 1,491-row field matrices | Scheduled in 51 |
| Blocked by product defect/capability | See [48](48-blocked-vs-unexecuted-and-prerequisites.md) §A | Remain BLOCKED; verify the blocker is still reproducible before relying on it |

## 4. Known data-quality caveats carried forward
1. Page "PASS" in 34b = loaded without crash only.
2. 34-complete-manual-ui-test-plan.csv rows are templated; they are plan rows, not results.
3. 22's `evidence_paths` was validated row by row (T-P00-003, `tests/build_finding_evidence_check.py`, result in [55](55-finding-evidence-and-root-cause-verification.csv)): 30 findings RESOLVED, 2 cite only rows of file 44 (F-008, F-015), and **4 have no resolvable evidence file (F-010, F-011, F-016, F-023)** — their cells held directory-level hints. F-014 was repaired in RUN-003 with precise paths. The remaining four need precise raw evidence (T-P20-001). Correction found: F-014's title says "raw DB error shown"; in the UI only a generic "Request failed with status code 500" toast is shown — the raw DB error is in the API response body (information exposure, see F-031).
4. Sensitive files `evidence/state-C.json`, `state-K.json`, `state-SA.json` (Playwright session state) exist in the shareable evidence tree. They were **not opened**. They must be excluded from any commit/upload; new session state is stored outside `evidence/`.

## 5. Open discrepancy list (not silently normalised)
| ID | Discrepancy | Status |
|---|---|---|
| DISC-001 | 119 discovered pages/route handlers vs 117 page routes loaded | Explained in 45: 2 are route handlers (N/A). Still to confirm which two by re-running route discovery (T-P03-002). |
| DISC-002 | Prompt says ≈350 files; actual 331 | Counted; no files are missing from the documented categories except the prompt's estimate |
| DISC-003 | progress.md batch-1 severity totals (26 findings) vs register (36) | progress.md superseded by batch-2 lines; update done in this run |
| DISC-004 | Prompt names snapshot `isolated-audit-batch2-final-2026-11-26.dump`; the persistent volume (not the dump) was used | Volume row counts equal the dump's recorded counts; the dump's restorability is validated by `evidence/runtime/27` from batch 2 |
