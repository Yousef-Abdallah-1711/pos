# PR change consolidation matrix

| PR | Disposition in candidate | Rationale |
|---|---|---|
| #14 | Partially/full source candidate | Conflict resolved by preserving current main theme/render behavior and server-side API configuration; related tests/config retained. |
| #18 | Source/test portion integrated | Generated browser screenshots/results excluded to avoid duplicating evidence; original PR retained. |
| #19 | Focused commits already represented; broad snapshot excluded | Five focused patches applied as empty/equivalent against current source. Snapshot is broad, conflicts with newer code, and has no full regression evidence. |
| #20 | Integrated with manual release trigger | Infrastructure and tests retained; workflow cannot deploy automatically on CI. Nothing was deployed. |
| #21 | Partially integrated | E2E stack/fixture changes were already covered by current candidate; large Phase 2A snapshot remains on original PR pending file-level reconciliation and tests. |
| #22 | Preserved by original PR | Research-only material; no production implementation change. |
| #23 | Preserved by original PR | Brand/UI evidence and screenshots remain on original PR; no brand publication. |

Generated browser evidence from #18 and bulk research/screenshots from #22/#23 were not copied into this candidate. No project runtime change was silently discarded from a dirty worktree.
