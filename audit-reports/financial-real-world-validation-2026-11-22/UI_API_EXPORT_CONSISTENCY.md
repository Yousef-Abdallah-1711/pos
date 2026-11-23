# UI, API, and export consistency

**Browser acceptance:** 11/11 passed in the authenticated Playwright suite (`financial-browser-tests.txt`, `financial-browser-junit.xml`). The journeys included English/Arabic, desktop/mobile, accounting navigation, report rendering, CSV/PDF download flows, RBAC denial, and tenant isolation.

The API order-to-report/export check found source order gross 5,500 cents; sales export and CSV gross 5,500; PDF gross 5,500 for one order; P&L and GL revenue 5,000; trial balance 5,500 per side. These values matched their distinct expected measures. Arabic mobile and English desktop Balance Sheet screenshots are diagnostic and visibly correspond to an unreconciled report; they do not prove accounting acceptance.

Evidence screenshots: `acceptance-en-balance-sheet.png`, `acceptance-ar-mobile-balance-sheet.png`. The suite exercised dashboard, not every storefront, Super Admin, tablet, or role requested in the original matrix. No claim is made for untested screens or exports.
