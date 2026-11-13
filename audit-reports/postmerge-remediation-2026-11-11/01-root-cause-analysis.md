# Root cause analysis

## Baseline and policy

- Baseline: `a6a2c6700a5fe883256deb1e5a3893108bd9865c` (`origin/main` at fetch).
- `C:\Users\OPT\Desktop\POS\CODING_STANDARDS.md` is absent. The tracked repository policy `CODING_STANDARDS.md` was read in full and used instead; update the stale absolute path in the local `AGENTS.md` policy when its owner permits.
- The prior verification recorded an `Html` prerender error for all four Next.js apps. A source search found no `next/document` or `Html` import under frontend application source.

## Reproduction result

The four clean application builds passed on Next.js 15.5.9 when supplied with the required production URL build variables. A landing build with those variables omitted failed earlier at `/robots.txt` because `NEXT_PUBLIC_LANDING_URL` is intentionally required for production SEO URLs. The former shared `/404` `Html` message did not reproduce under the declared production Compose build configuration, so its exact historical cause remains unproven and must not be described as a confirmed source defect.

The dashboard settings route was a client component that returned `null` while scheduling a router effect. It is now a server route redirect. The route is present in the successful build as static locale routes and no longer triggers the previously recorded settings prerender failure.

## Confirmed source defects fixed

- Landing hero areas shared the same `#product` key; the key now combines the stable href and label.
- The E2E Compose app bundle targeted API port `18080` while host browsers reach the API on published port `18082`; both now derive from `P2A_API_HOST_PORT`.
- Cookie-backed auth uses a non-bearer marker. The KDS socket hook incorrectly required a bearer token and therefore never opened a socket for these sessions.
- The dashboard flex content lacked `min-width: 0`, allowing children to widen the main flex item on small screens.
