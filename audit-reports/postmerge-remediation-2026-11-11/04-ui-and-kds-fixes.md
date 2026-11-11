# UI and KDS fixes

- Landing hero list keys now uniquely combine each stable destination and label. Browser navigation to `/en` returned no console errors; the prior duplicate-key error was not observed. Existing browser console history contains an LCP priority warning from the first pre-fix render; the hero image now declares `priority`.
- Dashboard main flex content now has `min-w-0`, preventing the flex item from expanding to the width of oversized descendants. The authenticated dashboard now runs in the isolated stack and measured CLS is 0.047; exact-width mobile overflow remeasurement remains open.
- The KDS socket now treats the cookie-backed authenticated session as eligible to connect. It no longer requires a usable bearer token, since browser WebSockets use the API's httpOnly session cookie.
- Regression coverage verifies that the cookie-session marker opens the KDS socket and unauthenticated state does not. Both tests pass.

The live WebSocket handshake and KDS status indicator remain unverified; the current Phase 2A suite does not exercise live KDS events.
