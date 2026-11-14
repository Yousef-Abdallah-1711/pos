# Visual QA

## Verified candidate

[Verified CI run](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43577385555), PR head `60f7499689927fd3439370c1a46bef34527fce93`, includes the current main finance fixes. This is built-image acceptance in a disposable synthetic stack, not proof of AWS deployment. [Machine-readable summary](evidence/candidate-verification.json).

Manual visual review covered English and Arabic at 1440px and 390px. The evidence shows loaded product captures, readable hierarchy, independent mobile compositions, RTL ordering, distinct section rhythms, and a clear monthly offer. Hero captures explicitly disclose English application views and synthetic data.

| Composition | Evidence |
|---|---|
| English desktop | [1440px](evidence/en-1440.png) |
| Arabic desktop | [1440px](evidence/ar-1440.png) |
| English mobile | [390px](evidence/en-390.png) |
| Arabic mobile | [390px](evidence/ar-390.png) |

Before: [English desktop](../audits/product-landing-2026-11-14/evidence/en-desktop.png), [Arabic mobile](../audits/product-landing-2026-11-14/evidence/ar-mobile.png).

## Responsive and interaction evidence

Every public route passed at 360, 390, 768, 1024, 1440 and 1920px in both languages: 168 route/viewport samples. The screenshots retained here are byte-identical to the manually reviewed candidate. The retained regression suite also passed 60 homepage/login/registration samples at ten widths and four mobile-menu checks. Annual mobile pricing was also captured and verified at 390px. The final release gate additionally measures visible, unobscured focus indicators with at least 3:1 contrast on light, dark, final-action and footer surfaces. Its exact result is in the latest [PR validation record](https://github.com/Yousef-Abdallah-1711/pos/pull/43) and `chromium-smoke-results` artifact. At 1920px the marketing image is a viewport capture; the page is still scrolled for screenshot loading and checked for document overflow.

No failed public network responses, browser exceptions, or console errors were recorded in the successful candidate. Native FAQ disclosure, keyboard walkthrough tabs, annual registration routing, locale switching, and menu Escape behavior passed.

## Accessibility and performance

The candidate's axe WCAG A/AA scans had zero violations. Lighthouse mobile accessibility scored 100 in both languages. Main palette contrast ratios: body 15.82:1; muted on the light violet surface 5.31:1; white action text 6.76:1; dark-section labels 8.24:1; final-section body 6.67:1; footer 8.92:1. Keyboard focus and reduced-motion behavior are implemented; controls use 48px target dimensions.

| Locale | Lighthouse performance | Accessibility | LCP | CLS |
|---|---:|---:|---:|---:|
| EN | 93 | 100 | 0.913s | 0 |
| AR | 96 | 100 | 0.987s | 0 |

These are standard Lighthouse mobile simulations in CI, not field conversion data or real-device measurements. NVDA/VoiceOver and physical Android/iOS testing were not performed. No certification or conversion-uplift claim is made.

SEO scored 69 solely because the staging edge deliberately sends `X-Robots-Tag: noindex, nofollow`. All other scored SEO audits passed. Keep this header for staging.
