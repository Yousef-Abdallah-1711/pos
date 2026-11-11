# UI and KDS fixes

- Landing hero list keys now uniquely combine each stable destination and label. Browser navigation to `/en` returned no console errors; the prior duplicate-key error was not observed. Existing browser console history contains an LCP priority warning from the first pre-fix render; the hero image now declares `priority`.
- Dashboard main flex content now has `min-w-0`, preventing the flex item from expanding to the width of oversized descendants. The code change is in place; the prior authenticated dashboard overflow could not be re-measured because no API or authenticated session was available.
- The KDS socket now treats the cookie-backed authenticated session as eligible to connect. It no longer requires a usable bearer token, since browser WebSockets use the API's httpOnly session cookie.
- Regression coverage verifies that the cookie-session marker opens the KDS socket and unauthenticated state does not. Both tests pass.

The live WebSocket handshake and KDS status indicator remain unverified without the isolated API/database stack.
