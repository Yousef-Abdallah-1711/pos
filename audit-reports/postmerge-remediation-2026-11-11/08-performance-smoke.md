# Performance smoke

- A development browser render of `/en` completed. Authenticated dashboard CLS was **0.047** in the full Phase 2A run and **0.051** in a focused repeat (both below 0.1). A separate dev-stack browser timing sample measured dashboard LCP at **1.07 s**, TTFB at **357 ms**, and CLS at **0.0035**; these are one-run development measurements, not production Lighthouse results.
- Landing layout was sampled at browser resize settings 320, 390, 768, and 1440; the Playwright browser reported effective viewport widths 400, 487, 960, and 1800 respectively (device scale factor 1.25). No horizontal document overflow was measured at those effective widths. These measurements do not verify the exact requested CSS widths.
- Authenticated dashboard viewport checks passed at exact CSS widths 320, 390, 768, and 1440 with no document overflow after the language switcher was compacted on small screens.
- Lighthouse, LCP, API latency distribution, bundle analysis, database query plans, and k6 load are NOT RUN. Redis health passed in the isolated Compose stack. The hero image now sets `priority`; the earlier pre-fix warning was not part of the final successful Playwright assertions.
