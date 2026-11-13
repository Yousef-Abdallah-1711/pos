# Open PR inventory

Latest GitHub state fetched 2026-11-13. Open PRs inspected: #14, #18, #19, #20, #21, #22, #23.

- #14 — `fix: remediate verified post-merge issues`, draft, 51 files, conflicts with main. Its candidate implementation was merged selectively into this branch; storefront conflict retained current rendering/theme behavior and the server API base helper.
- #18 — WIP Phase 00, draft, 30 files. Source/test commit `b1bfdc0` integrated; its separate generated browser evidence commit was left on the original PR.
- #19 — WIP WS02 security snapshot, draft, 75 files, conflict state DIRTY. Five focused source/test commits were compared; each patch was already present in the candidate/main-equivalent source. Its broad snapshot commit was not imported.
- #20 — gated deployment infrastructure, draft, 28 files, checks reported green at inspection. Seven commits imported. The release workflow was changed from automatic `workflow_run` to manual `workflow_dispatch`; no release or deployment was run.
- #21 — WIP Phase 2A, draft, 100 files, conflict state DIRTY. Its E2E Compose/seed changes were already represented by the more complete candidate. The unreviewed broad snapshot was not imported wholesale.
- #22 — financial architecture evidence, draft, documentation/research only; left preserved on its original PR.
- #23 — brand/UI research and screenshots, draft, evidence only; left preserved on its original PR. ORDERLY was not published.

No original PR was closed, force-pushed, or rewritten. Independent approval is not recorded for the open PRs.
