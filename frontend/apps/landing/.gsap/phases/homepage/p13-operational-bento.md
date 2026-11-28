# P13 — Operational bento

Status: implemented, browser QA in progress (2026-11-27).

The previous `.cv-problem-grid` was replaced by `OperationalBento.tsx`. The composition uses a large POS card and two distinct supporting cards with verified locale/theme captures. Each card has three local explanatory states and keyboard/touch controls. Desktop ScrollTriggers start brief GSAP timelines on entry; the section is never pinned. Mobile and reduced-motion users see the content immediately and select steps manually. The bridge links to `#how-orders-move`. The explanatory steps are explicitly identified as illustrative and make no product API calls.

Evidence: `docs/landing-interactive-experience/artifacts/problems-before.png`, `artifacts/operational-revision/`, and `operational-bento-showcase-report.md`.
