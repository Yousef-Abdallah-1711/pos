# Performance smoke

- A development browser render of `/en` completed. No Lighthouse run or production Core Web Vitals measurement was performed.
- Landing layout was sampled at browser resize settings 320, 390, 768, and 1440; the Playwright browser reported effective viewport widths 400, 487, 960, and 1800 respectively (device scale factor 1.25). No horizontal document overflow was measured at those effective widths. These measurements do not verify the exact requested CSS widths.
- Dashboard mobile overflow was not re-measured because authentication and the isolated API were unavailable. The `min-w-0` flex correction is implemented but not presented as a verified fix.
- LCP, CLS, API latency, bundle analysis, database query plans, Redis health, and k6 load are NOT RUN. The browser console showed no errors; its accumulated history includes the initial pre-fix Next.js LCP priority warning. The hero image now sets `priority`.
