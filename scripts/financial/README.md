# Local financial browser acceptance capture

`capture-local-balance-sheet.mjs` captures English desktop and Arabic mobile balance-sheet pages from a disposable local stack. Run it from the repository root and use only synthetic accounts and data.

Set `SERVIO_FINANCIAL_FIXTURE_PATH` to a synthetic JSON fixture kept outside the repository. It must contain `accountant.email` and `accountant.password`. Set `SERVIO_FINANCIAL_API_URL` and `SERVIO_FINANCIAL_DASHBOARD_URL` to loopback origins. The script rejects non-loopback origins, does not follow API redirects, and blocks browser HTTP(S) requests outside those two origins. Screenshots go to `SERVIO_FINANCIAL_SCREENSHOT_DIR`, or to a temporary directory when unset. Do not commit fixture files or captured screenshots.
