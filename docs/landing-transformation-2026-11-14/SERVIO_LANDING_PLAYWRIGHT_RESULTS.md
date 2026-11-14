# Browser and test results

Local landing TypeScript: passed. Local landing component/unit suite: 64 tests passed in 11 files. Local landing lint: no errors, ten existing authentication/navigation/library warnings. Test TypeScript and expanded browser verification are in progress.

`frontend/apps/landing/scripts/marketing-acceptance.mjs` is the reproducible public browser check. It covers both locales, the six requested widths, all dedicated routes, annual signup navigation, walkthrough keyboard behavior and mobile menu Escape/focus behavior. Run with `SERVIO_MARKETING_URL` and `SERVIO_MARKETING_EVIDENCE`; optional `SERVIO_BROWSER_EXECUTABLE` supports the staging host's installed browser. Exact final counts and live deployment checks remain pending.

Backend tests include policy-copy consistency, billing preference validation, exact expiry boundary, registration duration/preference, and disabled checkout rejection. Full integration, race, build, and image/browser gates run in the existing verified CI workflow. No CI result is assumed from local checks.
