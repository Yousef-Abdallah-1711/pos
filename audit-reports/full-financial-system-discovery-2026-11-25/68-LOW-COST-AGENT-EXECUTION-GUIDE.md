# 68 - Low-cost agent execution guide

## Initial context (read only these)
1. 65-AUDIT-MASTER-OVERVIEW.md  2. this file  3. last block of progress.md  4. the next NOT-STARTED row in 66  5. that phase's `phases/phase-NN/README.md`  6. the row's required files. Only then open source or evidence for that task.

## Per-task protocol
1. Read the task row (66) and any per-item list (66b/66c/66d).  2. Verify prerequisites (dependency tasks done; fixtures).  3. Verify the environment (`docker ps`, health 200, protected tenants unchanged).  4. Prepare the fixture through the API with a unique prefix `000R6<random>`; never reuse tenants 5/6/7; tenant 10 is read only; tenant 4 is a mutated diagnostic fixture.  5. Execute the named script in headed Chrome (`AUDIT_CHROME_PROFILE` env gives each concurrent runner its own profile; `RUN_ID=RUN-006`).  6. Capture evidence (screenshot, network request+response, SQL before/after) under `evidence/<type>/RUN-006/PH-<phase>/`.  7. Compare expected and actual; a product failure is a finding, not a failed task.  8. Update: append a row to 34b through a small recorder script in `tests/` (copy `r5_record_destructive_b.py`), add findings through a `tests/r6_register_update*.py` script modelled on `r5_register_update5.py`, set the task status in `66-queue-status.csv`.  9. Append a short line to the RUN-006 block of progress.md.  10. Pick the next row in 66 whose dependencies are done.

## Rules learned the hard way
- Bash heredocs with apostrophes fail: write scripts with the file tool. Use `python -X utf8` (Arabic output).
- Mutating API calls need `Origin: http://localhost:13112` and cookie `poos_session`; login is rate limited to 20/min per IP (a burst of bad logins locks you out ~1 min; run the rate-limit test last).
- Lists show 20 rows sorted by name: prefix fixtures with `000` so they appear first. Icon-only buttons need aria-label selectors. The leaves page hides future dates by default.
- Memory: 7.7 GiB Docker VM. Keep restaurant-website and super-admin containers stopped unless the task needs them; start one at a time.
- Never switch git branches with the audit folder untracked; never `docker system prune`; never delete volumes, evidence, dumps or the 2 git stashes; never open `evidence/state-*.json`.
- A field batch result is meaningless until a clean baseline submit is ACCEPTED (RUN-005 lost time to five invalid baselines).

## Context efficiency
Search with Grep instead of reading whole files; read only the finding package of the related finding; reuse runners and fixtures; do not regenerate all documents after each test, only at checkpoints (order in 67); keep command output short (`| tail`); put durable results in the audit package, not in chat.

## Escalation (do not stop the audit; skip to the next independent row)
Escalate to a stronger model only for: a root cause that cannot be traced after reading the cited source paths; an accounting policy decision; a destructive action the owner has not authorised; an environment that differs materially from 65; evidence that contradicts the documented architecture. Capabilities in 66: ROUTINE-EXECUTION and BROWSER-AUTOMATION need only a low-cost model; DATA-VERIFICATION and SOURCE-TRACING need a mid-level model; FINANCIAL-RECONCILIATION, COMPLEX-INVESTIGATION and POLICY-REVIEW need the strongest available model.

## Closing run (Q-ENV-002)
Final pg_dump with full sha256 (not abbreviated), backup tar.gz excluding `evidence/state-*.json`, tenants 5/6/7 COA rows 0, 2 stashes, `git status` shows no product changes by the audit, dev DB `pos_saas` equals `evidence/runtime/RUN-003/46-dev-db-fingerprint-before.txt`, record all in 64 and 62.
