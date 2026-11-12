# Performance and reliability

- No Lighthouse, Core Web Vitals, authenticated production-build metrics, k6/load stage, API latency distribution, database slow-query review, or saturation test was run.
- Production builds fail before a production server can be launched, so production CLS/LCP are NOT RUN.
- POS UI diagnostic timings displayed product-fetch around 131ms, page-load around 132ms, and order-submit around 395ms in this local run. These are single-run UI diagnostics, not a benchmark or acceptance result.
- The KDS showed a disconnected live connection indicator despite its initial API order load and successful status update. Realtime reliability is unverified.

