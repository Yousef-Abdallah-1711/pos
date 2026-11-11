# Performance smoke

- A development browser render of `/en` completed. The isolated authenticated dashboard CLS Playwright check reported **0.047** on the final run, below the 0.1 threshold for that specific dev-stack route and run.
- Landing layout was sampled at browser resize settings 320, 390, 768, and 1440; the Playwright browser reported effective viewport widths 400, 487, 960, and 1800 respectively (device scale factor 1.25). No horizontal document overflow was measured at those effective widths. These measurements do not verify the exact requested CSS widths.
- Dashboard mobile overflow was not re-measured at the exact 320/390/768/1440 CSS widths. The `min-w-0` correction remains implemented; the CLS result does not verify overflow behavior.
- Lighthouse, LCP, API latency distribution, bundle analysis, database query plans, and k6 load are NOT RUN. Redis health passed in the isolated Compose stack. The hero image now sets `priority`; the earlier pre-fix warning was not part of the final successful Playwright assertions.
