# Landing performance and accessibility evidence

## CI failure and font-loading experiment (2026-11-27)

PR #137 application-images CI failed its Lighthouse performance threshold (90): EN scored 88 with LCP 3.15 s and AR scored 86 with LCP 3.13 s. Accessibility was 100 and CLS was zero in both. The gate was not weakened.

The prior local audit showed eight font files preloaded at high priority for both locales, including IBM Plex Arabic files on the English page. A candidate disables automatic Plex preloading while retaining `font-display: swap` and the same font families and weights; the browser requests the applicable font when used. In sequential local mobile Chrome Lighthouse runs, the prior production container scored EN 74 / AR 76 with LCP 4.36 / 4.37 s; the candidate scored EN 85 / AR 85 with LCP 3.75 / 3.54 s. Both had accessibility 100 and CLS zero. English font requests fell from eight to two; Arabic continued to load eight fonts when needed. The Arabic hero reported IBM Plex Sans Arabic loaded after `document.fonts.ready`.

These are individual lab runs, not field results; the baseline build predates the small POS illustration, so this is not a perfectly isolated comparison. They support the change but do not prove that CI will pass or that the 2.5 s LCP goal has been met. Repeat the image gate and test actual devices and assistive technology before a final release verdict.

The first CI run after the font change passed Lighthouse (EN 96 / AR 93, LCP 2.40 / 2.86 s, accessibility 100, CLS 0 / 0.025) but failed Firefox dark marketing acceptance because Firefox reported canceled Arabic font requests during navigation. Locally, deliberately delaying font responses and immediately navigating reproduced Firefox status `2152398850`; waiting for `document.fonts.ready` loaded six Arabic font faces before navigating and produced zero font errors. The acceptance script now waits for fonts before its rapid route changes. It does not ignore console errors or relax thresholds; the full image gate must pass again. Arabic CI LCP remains above 2.5 s in that run.

## Final-source follow-up (2026-11-27)

One local Lighthouse 13.6.0 mobile run against the final production container returned performance 57, accessibility 100, LCP 5.41 s and CLS 0. The LCP element remained `.cv-hero-photo`, and the audit reported the image was discoverable, eager and high priority. Root-document response was about 20 ms and total blocking time about 1.04 s. These are single-run local lab measurements under different runtime conditions from the earlier run, so they do not establish a controlled regression. The 2.5 s target remains unmet; no unmeasured image or animation change was made. Repeat controlled EN/AR runs and profile main-thread work before claiming a fix.

Measured 2026-11-27 against the local Next.js production preview at `http://localhost:3001` with Lighthouse 13.6.0 and installed Chrome. These are simulated lab results, not field Core Web Vitals or physical device measurements. Raw reports are `artifacts/lighthouse-*-final.json`. The brief targets LCP at or below 2.5 s and CLS below 0.1.

| Route / profile | Performance | Accessibility | FCP | LCP | TBT | CLS |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `/en`, mobile | 90 | 100 | 1.37 s | 3.31 s | 167 ms | 0 |
| `/ar`, mobile | 78 | 100 | 1.52 s | 4.01 s | 358 ms | 0 |
| `/en`, desktop | 100 | 100 | 0.37 s | 0.75 s | 0 ms | 0.004 |

The hero photograph is the mobile LCP element in both locales. It was initially a CSS background and the first English mobile measurement was LCP 6.1 s / performance 69. A prioritized `next/image`, quality 50, and removal of an unnecessary below-fold preload brought the final English measurement to 3.31 s / 90. The image remains subdued beneath the intended dark overlay. The solid Coral setup-card mark resolved the desktop contrast finding, and reserved mobile hero title height removed the measured Arabic layout shift.

**Verdict:** the measured CLS and accessibility goals pass. Mobile LCP still misses the 2.5 s target in both locales. INP cannot be inferred from these short lab runs. The next performance pass should profile hydration and the eagerly loaded GSAP client bundle, then retest on representative hardware. These scores must not be presented as real-user performance.

The final rebuild after these measurements changed only the desktop notice layout for a temporarily unavailable auth API. At 1200 px, a forced 503 check placed the notice below the navigation at `y=76`, kept the document width at 1200 px, and preserved login and trial controls. The Lighthouse reports were collected before that scoped notice-only CSS change.
