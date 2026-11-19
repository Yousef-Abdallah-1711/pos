# SERVIO Landing Header UX Specification

**Status:** Proposed behavior. The current Landing Navbar does not implement these authenticated states.
**Evidence baseline:** `frontend/apps/landing/src/components/Navbar.tsx:16-88`; current public deployment observed 2026-11-17 UTC.

## Current observation

**VERIFIED CURRENT BEHAVIOR:** The Navbar renders marketing links, locale switch, theme toggle, Login, Trial CTA and mobile menu. It does not subscribe to the landing auth store, query `/auth/session`, show a profile, or offer logout. Landing login does establish an HttpOnly API-host cookie, but the UI cannot read it; a full navigation to Dashboard checks the session successfully in source flow.

**OBSERVED DEFECT:** A user with a valid dashboard session still sees guest actions on Landing. This is an absent session-aware header, not evidence that the API cookie is inaccessible to credentialed requests. Current landing store is memory-only and must not be treated as cross-app auth proof.

## Header states

| State | Desktop | Mobile | Rules |
|---|---|---|---|
| Guest | Brand, marketing navigation, language/theme, Login, prominent Start Free Trial | Brand, language/theme and menu; Login and trial are visible in the opened drawer and reachable without scrolling | Public shell can be statically rendered and cached. |
| Session checking | Keep stable header dimensions; show neutral account placeholder or retain guest actions only if no misleading login flash; mark account area `aria-busy=true` | Same neutral account slot, no blocking overlay | One bounded check on first hydration. Do not block hero or navigation. |
| Authenticated restaurant user | Avatar/initials button, display name, current restaurant if authorized, Account menu with Dashboard, Settings and Sign out; marketing links remain | Drawer exposes Dashboard and account links; focus management and Escape behavior | Data comes from no-store session DTO. Never infer role/tenant from localStorage. |
| Multiple restaurant memberships | Account menu shows current restaurant and “Switch restaurant” link/selector only after API returns authorized choices | Accessible list/radio dialog with current item and confirmation | Do not accept tenant IDs without server authorization; selection triggers server endpoint and then reloads session. |
| Super Admin session only | Do not expose ordinary staff account as signed in; optionally show a neutral “Admin portal” destination only when separately validated for that host/origin | Same, with no tenant/admin actions embedded in marketing drawer | Landing must not infer admin status from a staff cookie or expose privileged control. Default: keep guest state on landing. |
| Expired/revoked | After authenticated state gets 401, clear local view, return to guest actions without modal; next protected navigation sends to login | Same | 401 is distinct from 5xx/network outage. |
| API unavailable | Keep current stable header; show no stale authenticated identity as authoritative; optional non-blocking retry after navigation/focus | Same | Do not redirect from marketing page due to temporary API failure. |

## Profile menu content and interaction

- Button accessible name: `Account menu for {display name}`; expose `aria-haspopup="menu"`, `aria-expanded`, and keyboard arrow/escape behavior.
- Menu order: display name and email (not tenant secrets); current restaurant/tenant; Dashboard; Account settings; Sign out.
- Use server-generated safe paths/origins. Dashboard link targets only the canonical allowlisted `DASHBOARD_URL`; never use arbitrary `return_to` URL.
- Sign out calls API `POST /api/v1/auth/logout` with credentials and Origin/CSRF protections. On success clear landing view state; optionally BroadcastChannel to other same-origin landing tabs; dashboard/admin state is refreshed from API on focus/navigation.
- For API/network failure, do not claim server logout succeeded. Keep a retry message/action and clear only transient menu UI, or provide an explicit local-only sign-out state with clear wording if business accepts that policy.
- No profile avatar URL should be fetched from arbitrary external host; use safe initials fallback and bounded image source policy.

## Recommended rendering model

**PROPOSED DESIGN:** Keep the marketing shell public and cacheable; add one small client-only session island in the account-action area. The island calls `GET /api/v1/auth/session` once after hydration with `credentials: include` and `cache: no-store`, then uses an in-memory DTO. This fits the current host-only cookie because the request goes to `api.pos.rakmyat.online` with exact CORS credentials. It avoids personalized HTML and per-request landing SSR. Use a short in-memory freshness window and recheck on focus/navigation after a defined idle period. No repeated request per component render.

An app-specific BFF is an alternative if browser CORS/Origin complexity or API origin exposure becomes undesirable. The BFF must pass the cookie securely, set private/no-store for personalized responses, and never cache by a shared public key. Do not set `Domain=.rakmyat.online` for convenience.

| Model | Advantages | Trade-offs | Recommendation |
|---|---|---|---|
| Fully personalized SSR | Correct profile on first HTML; no client auth flash | Public HTML becomes user-specific, can bypass CDN or risk cache leakage; every page render hits session store; server must forward cookie safely | Avoid for all marketing routes. Consider only a private app shell if measured need. |
| Cached public shell + session island | Preserves CDN caching, low API traffic (one small request per visit), simple integration with current API cookie | Brief neutral/loading account slot; client JS needed for profile | **Recommended initial landing approach.** |
| BFF session proxy | Avoids browser CORS handling; can return server-rendered session on a private route | Adds server code, cookie forwarding and cache boundary; introduces another runtime hop | Use if CORS or API coupling proves operationally costly. |

## Layout, responsive and accessibility requirements

- Keep header height fixed between guest, checking and authenticated states to prevent layout shift.
- Preserve current EN LTR and AR RTL direction. Align profile menu to logical `start/end`, not hard-coded left/right.
- On small screens use a single drawer; do not add a separate profile bar that pushes content below the fold.
- Interactive targets at least 44×44 CSS px; visible focus; semantic `<button>` for menu and `<nav>` labels.
- Menu must work with keyboard and screen reader. Return focus to trigger on Escape/close; trap focus only if rendered as modal dialog.
- Respect light/dark theme, reduced motion, high contrast and user zoom. Do not encode status only through color.
- RTL display name, email, punctuation and mixed-direction identifiers should be isolated with `dir="auto"` / `<bdi>` where needed.

## Acceptance criteria

1. Guest render has login/trial and no profile API leak in HTML or shared cache.
2. Login on landing, then returning to landing in same browser shows profile after one successful session call.
3. Dashboard session already active, new landing tab loads profile without entering credentials.
4. API 401 returns to guest state; API 503/network failure never masquerades as logged out or causes loop.
5. Sign out revokes server session and a new dashboard request receives 401.
6. Staff user sees only authorized tenant list, safe account settings and dashboard destinations.
7. Super Admin identity is not presented as restaurant user; admin action is never exposed on landing.
8. No personalized response has shared CDN cacheability; anonymous `/en` and `/ar` remain cacheable.
9. EN/AR, RTL/LTR, dark/light and 320/390/768/1024/1440 viewport checks pass with no overflow, focus trap or hydration mismatch.

**UNVERIFIED ASSUMPTION:** Marketing pages currently use a stable static/ISR strategy in all deployment modes; live response shows `s-maxage=300`, but exact Next rendering/cdn behavior per route needs a browser/network trace.

**REQUIRED BUSINESS DECISION:** Should the header show the current restaurant name? Should it expose a restaurant switcher directly or link to Dashboard? Should Super Admin see an admin portal shortcut on marketing pages? Recommended safe default is no admin identity leakage and a Dashboard link only for staff sessions.

See [session architecture](SERVIO_MULTI_PORTAL_SESSION_DESIGN.md) and [acceptance matrix](SERVIO_AUTH_E2E_ACCEPTANCE_MATRIX.md).

## Implementation update — 2026-11-17

This report was reconciled against `origin/main` at `c052c817e114df4d060d8428e6e236f6011876f8`. The current implementation slice, P0/P1 traceability, test evidence, and outstanding work are recorded in [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md); the refreshed repository/CI baseline is in [BASELINE.md](BASELINE.md). Historical audit conclusions remain intact. This update does not mark unimplemented Redis replay/rate-limit work, the full E2E matrix, Super Admin/storefront acceptance, cache-boundary testing, or 10K measurements as complete.
