# UI and KDS fixes

- Landing hero list keys now uniquely combine each stable destination and label. Browser navigation to `/en` returned no console errors; the prior duplicate-key error was not observed. Existing browser console history contains an LCP priority warning from the first pre-fix render; the hero image now declares `priority`.
- Dashboard main flex content has `min-w-0`; a responsive check then exposed the language switcher text pushing the mobile header past the viewport. The switcher now uses an accessible globe-only control below the `sm` breakpoint. The authenticated browser assertion passes at CSS widths 320, 390, 768, and 1440 with no horizontal overflow; CLS measured 0.047 in the full run and 0.051 in the focused repeat, both below 0.1.
- The KDS socket now treats the cookie-backed authenticated session as eligible to connect. It no longer requires a usable bearer token, since browser WebSockets use the API's httpOnly session cookie.
- Regression coverage verifies that the cookie-session marker opens the KDS socket and unauthenticated state does not. Both tests pass.
- The authenticated Phase 2A browser flow now opens the KDS page, verifies the live connection, creates a POS order, and confirms the matching order arrives over the tenant WebSocket.
- The latest 16-test Phase 2A run passed, including live KDS delivery; latest authenticated dashboard CLS measurement was 0.031.
