# Progressive Rendering Recommendations

Recommendations are source-driven hypotheses pending user-facing authenticated trace validation:

1. Separate the immediately useful dashboard shell from auth-dependent business data, while preserving a strict redirect/permission gate before exposing tenant content.
2. Replace provider-wide full-screen waiting with a stable shell only if session bootstrap can safely show non-sensitive chrome and the loading state does not flash unauthorized data.
3. Add route-level loading boundaries to currently uncovered app segments where real route content has a predictable shape; avoid generic spinner pages.
4. In Restaurant Website, consider progressive independent rendering of theme/homepage/products if product content can render without holding the whole server response; preserve SEO and tenant scope.
5. Verify host middleware's status-check latency and freshness before relocating it; it currently gates tenant-host navigation.
6. Keep background refetch and polling out of initial readiness. Confirm hidden-tab pause and stable query keys.
7. Measure CLS and focus/keyboard behavior on skeletons before adoption.

No change was implemented. Expected benefits are unquantified until measured against production-mode baseline; see `18_FRONTEND_PERFORMANCE_OPPORTUNITIES.csv`.
