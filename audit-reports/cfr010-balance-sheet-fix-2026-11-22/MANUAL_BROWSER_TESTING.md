# Manual Browser Testing

Status: **NOT RUN** for this continuation.

The earlier `BROWSER_ACCEPTANCE.md` records authenticated automated Playwright journeys. They are not manual browser interaction and must not be described as manual testing. The earlier task's disposable PostgreSQL, Redis, API, and dashboard services are no longer running; no new manual transaction or report journey was completed in this continuation.

To complete this acceptance, start a new disposable environment and synthetic accountant fixture, then interact with the dashboard to create and post a synthetic transaction, compare the Trial Balance, P&L, and Balance Sheet with exact SQL results, post another transaction, and repeat the reports and PDF gate in English and Arabic at desktop and mobile sizes. Capture safe screenshots and a scrubbed trace, and record each amount and date filter. Do not use production credentials or operating financial records.
