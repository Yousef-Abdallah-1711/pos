# Landing performance and accessibility evidence

Measured 2026-11-27 against the local Next.js production preview at `http://localhost:3001` with Lighthouse 13.6.0 and installed Chrome. These are simulated lab results, not field Core Web Vitals or physical device measurements. Raw reports are `artifacts/lighthouse-*-final.json`. The brief targets LCP at or below 2.5 s and CLS below 0.1.

| Route / profile | Performance | Accessibility | FCP | LCP | TBT | CLS |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `/en`, mobile | 90 | 100 | 1.37 s | 3.31 s | 167 ms | 0 |
| `/ar`, mobile | 78 | 100 | 1.52 s | 4.01 s | 358 ms | 0 |
| `/en`, desktop | 100 | 100 | 0.37 s | 0.75 s | 0 ms | 0.004 |

The hero photograph is the mobile LCP element in both locales. It was initially a CSS background and the first English mobile measurement was LCP 6.1 s / performance 69. A prioritized `next/image`, quality 50, and removal of an unnecessary below-fold preload brought the final English measurement to 3.31 s / 90. The image remains subdued beneath the intended dark overlay. The solid Coral setup-card mark resolved the desktop contrast finding, and reserved mobile hero title height removed the measured Arabic layout shift.

**Verdict:** the measured CLS and accessibility goals pass. Mobile LCP still misses the 2.5 s target in both locales. INP cannot be inferred from these short lab runs. The next performance pass should profile hydration and the eagerly loaded GSAP client bundle, then retest on representative hardware. These scores must not be presented as real-user performance.

The final rebuild after these measurements changed only the desktop notice layout for a temporarily unavailable auth API. At 1200 px, a forced 503 check placed the notice below the navigation at `y=76`, kept the document width at 1200 px, and preserved login and trial controls. The Lighthouse reports were collected before that scoped notice-only CSS change.
