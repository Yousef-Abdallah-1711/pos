# UI, responsive, and accessibility observations

Screenshots captured by Playwright and copied into this directory's `evidence/` folder:

- `landing-en-1440.png`
- `landing-ar-390.png`
- `login-en-1440.png`
- `login-en-390.png`
- `dashboard-authenticated-en-1440.png`
- `dashboard-authenticated-en-mobile.png`

Observed:

- Landing and auth pages rendered in English LTR and Arabic RTL.
- The shared website header remained visible on login in both language checks; the register page rendered with its mobile menu control.
- At the requested 390px dashboard resize, the browser reported viewport width 487 CSS px and document width 603 px, indicating horizontal overflow. The Playwright host applies a scale factor, so record the tool's reported dimensions when comparing.
- Dashboard rendered and loaded authenticated data on desktop and mobile.
- Landing browser console reports duplicate React key `#product` in the hero area. This was also the user's previously reported console defect.
- Dashboard and KDS had no observed console errors during the latest route check; warnings remain.
- The screenshots are evidence of rendering only, not accessibility certification. Keyboard-only flows, contrast, screen-reader semantics, 768/820/1024/1920 layouts, all route-level visual review, and hydration coverage remain NOT RUN.

