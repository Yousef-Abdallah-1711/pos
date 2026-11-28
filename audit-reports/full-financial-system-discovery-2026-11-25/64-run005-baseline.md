# 64 - RUN-005 baseline and continuity record (started 2026-11-27)

Authorization: Master Prompt 06 (complete remaining feasible audit work; no remediation; do not touch the landing/marketing workstream; no commit/push/merge/deploy of the audit into the product repository).

## Storage and backup (Phase A.1)
- Disk C: 238 GB, ~18 GB free (93% used) at start. Audit package ~76 MB on disk (16.5 MB are ignored files: dumps and session-state). Docker: images 4.6 GB, volumes 14.7 GB (13.7 GB reclaimable, includes the persistent audit DB volume and node/go caches - NOT pruned), build cache 2.5 GB.
- **Continuity incident:** the audit folder had been committed to a local/remote branch `audit/financial-discovery-run004` (user request at the end of RUN-004). A later branch switch by the other workstream (`feature/servio-bento-showcase-20261127`) removed the tracked audit files from the working tree (17 ignored files remained). All 1,825 tracked files were restored non-destructively with `git restore --source=audit/financial-discovery-run004 --worktree -- <audit dir>` (no branch switch, no staging); 1,842 files present. The branch on origin is an additional off-machine copy of everything except session-state files and dumps.
- Per Master Prompt 06 the audit is NOT committed to the product repository from RUN-005 on. Do not switch branches blindly: the other workstream's branch checkout deletes tracked audit files again; this working copy is now untracked on the current branch (verify with `git status`).
- Verified backup: `C:\Users\Yousef\audit-backups\audit-full-financial-system-discovery-RUN005-start-20261127T114510Z.tar.gz` (55.1 MB, sha256 `6eba4e5424825dd7dadd0fd41d5057c77918ae9b0bd90dc3114c2195b38fa7fb`, checksum file alongside, 2,052 entries, archive readable, **no state-*.json inside**). Limitation: same physical drive (C: is the only drive); copy it off-machine when possible.
- Nothing deleted; no `docker system prune`, no volume deletion, stashes intact (2).

## Repository state (Phase A.2, read-only)
- Branch at start of RUN-005: `feature/servio-bento-showcase-20261127` (the other workstream), HEAD `35100f76efd12268077e1791689d1fd7523c7866` = `origin/main` tip (merge of PR #136). Working tree: modified `docs/landing-handoff.md` and `docs/landing-interactive-experience/progress.md` (other workstream), tooling tsconfig/next-env files, untracked `.playwright-mcp/*` and the restored audit folder. Remote `origin = https://github.com/Yousef-Abdallah-1711/pos.git`. `origin/main` is protected (pull request with 3 required checks); the RUN-004 push went to the audit branch.
- Changes since RUN-004 closure: PR #136 merged (landing). Diff of backend since the audit baseline `133260d`: none (verified in RUN-004; re-verified below in the evidence record).
- Historical SHAs: `133260d...` (RUN-001..003), `0de839b...` (RUN-004 closure). Evidence keeps its original attribution; RUN-005 evidence is attributed to the SHA recorded in `evidence/runtime/RUN-005/01-baseline.txt`.

## Environment at start (Phase B)
- The development containers named in RUN-004 (`pos-landing`, `pos-dashboard`, `pos-restaurant-website`, `pos-super-admin`) and the two landing CI containers **no longer exist** (removed outside this audit, apparently while freeing disk space); the isolated-stack containers and their images were removed as well (volumes kept: `pos-p2a-ws05-e2e_p2a_audit_pgdata` and 17 others). Only `pos-backend` (dev, in a restart loop waiting for `postgres:5432`) and `pos-postgres` run. Dev DB `pos_saas` fingerprint re-verified identical to `evidence/runtime/RUN-003/46-dev-db-fingerprint-before.txt` (tenants 2, users 6, orders 0, journals 0, COA 62, migrations 145).
- Because those containers are absent, `docker start ...` is not possible; recreating the dev stack (compose project `pos-dev`) is the user's action and is not performed by the audit.
- Isolated stack: rebuilt from the audit compose file on the persistent volume (see evidence/runtime/RUN-005).


## RUN-005 final state (2026-11-27)
- Final isolated DB snapshot: `evidence/database/isolated-audit-run005-final-2026-11-27.dump`, sha256 `0ca052179a615ab4f5c9915a19ad47212920fe5bccce4d86628fc38d3e0ced46`.
- Final backup: `C:\Users\Yousef\audit-backups\audit-full-financial-system-discovery-RUN005-final-20261127T1345Z.tar.gz` (68 MB, sha256 `95878f5f50096120ed98c7aa6d50c4b5c8197b9359a4e8b9423900b987d83c02`, session-state files excluded).
- Protected tenants 5, 6, 7: chart of accounts rows 0 each (unseeded). Git stashes: 2 intact. Product code: not modified by the audit (the six tsconfig/next-env changes in git status pre-date this session and belong to another workstream).
- Dev DB `pos_saas` fingerprint unchanged (tenants 2, users 6, orders 0, journals 0). Dev frontend containers still absent; recreating the dev stack is the user's action.
- Regeneration order that works: build_finding_evidence_check, build_finding_packages, r4_task_overlay, build_phase_specs, build_field_denominator, build_field_results, build_control_interaction_results, r5_record_destructive, r5_record_destructive_b, validate_packages, build_policy_detail, build_metric_verification_classes, build_evidence_manifest, build_closure_report, r5_checkpoint (build_control_interaction_results resets 53b, so both recorders must follow it).
