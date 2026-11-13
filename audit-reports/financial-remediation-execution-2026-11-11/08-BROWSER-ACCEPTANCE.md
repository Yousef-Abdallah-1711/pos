# 08 — Browser acceptance

## Follow-up Arabic RTL viewport verification

The closure follow-up captured a separate real-browser run at Arabic 390, 768, and 1440px and English 390px. The financial report's locale root, main, navigation, warning and table computed the expected RTL/LTR direction; no tested viewport overflowed horizontally. The document element remained `lang=en`, without `dir`, and computed LTR even for Arabic. The earlier `ltr` observation sampled the document root rather than the localized page. See `closure/03-ARABIC-RTL-INVESTIGATION.md`, `closure/evidence/rtl-viewports.json` and screenshots. Page-level RTL passes; document metadata remains open.

## Closure verification

Browser checks were run against an isolated local production build using synthetic API fixtures. All seven scenarios passed: unreconciled status, balanced status, unknown integrity status, 390px Arabic RTL layout, loading state, API 500 state, and unauthorized report response. For the balanced fixture, PDF download passed; PDF export controls were absent for unreconciled and unknown status. No page-level exceptions were recorded.

The run recorded 801 requests because unrelated dashboard shell components repeatedly requested user settings; this is a mock-browser request count and does not establish production request volume. Browser console included the deliberate fixture 500 and 401 responses plus a 404 for a missing resource. The 500/401 produced the expected visible API error and login redirect with figures withheld. The initial mobile probe sampled document direction; follow-up page-level RTL verification and the root metadata discrepancy are documented above.

Evidence: `test-results/browser-acceptance.json`, `screenshots/p00-balanced-en-desktop.png`, `screenshots/p00-unreconciled-en-desktop.png`, `screenshots/p00-unreconciled-ar-mobile.png`, and `test-results/synthetic-balanced-balance-sheet.pdf`. These checks use synthetic API responses and do not verify authenticated backend behavior, tenant isolation, or live financial data.
