# Arabic RTL Investigation

## Result

**Tested page layout: PASS. Document metadata: OPEN ISSUE.** Browser evidence is in [`rtl-viewports.json`](evidence/rtl-viewports.json), screenshots in `evidence/rtl-screenshots/`, and the browser probe in [`phase00-rtl-viewports.mjs`](evidence/phase00-rtl-viewports.mjs).

| Locale / width | Page root | Main | Navigation | Warning | Table | Overflow |
|---|---|---|---|---|---|---|
| Arabic / 390px | RTL | RTL | RTL | RTL, right aligned | RTL, header right aligned | None |
| Arabic / 768px | RTL | RTL | RTL | RTL, right aligned | RTL, header right aligned | None |
| Arabic / 1440px | RTL | RTL | RTL | RTL, right aligned | RTL, header right aligned | None |
| English / 390px | LTR | LTR | LTR | LTR, left aligned | LTR, header left aligned | None |

## Root cause and scope

The earlier reported `ltr` reading sampled `document` and `body`, not the localized page container. The Arabic locale root supplies `lang=ar` and `dir=rtl`, and the relevant page elements inherit RTL. The document element still has `lang=en`, no `dir`, and computes LTR. This is an existing app-wide document metadata discrepancy observed during the Phase 00 page check; it is not a demonstrated Phase 00 layout regression and was not modified here.

This browser run used an isolated local production build and synthetic API fixtures for rendering. It proves responsive direction and alignment for the sampled financial report UI, not full product-wide RTL accessibility or real-data behavior.
