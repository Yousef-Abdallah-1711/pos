# Authenticated synthetic visual baseline

Captured 2026-11-27 using the actual local dashboard in a Playwright-controlled browser and a fresh disposable `pos-test` PostgreSQL/backend stack. The tenant and account were created solely for this audit; displayed owner, restaurant and data are synthetic. Dashboard ran in development mode from the local `servio-landing-qa:20261127` image; exact image source commit was not independently recorded. Most captures show default empty data and a trial owner, not representative populated orders or restricted permissions. The Next.js development badge is visible. No cookies, authentication credentials or database files are stored here.

| File | Route / state | Viewport |
| --- | --- | --- |
| [dashboard-en-light-expanded.png](dashboard-en-light-expanded.png) | `/en/dashboard`, light, expanded | 1440 × 900 CSS pixels |
| [dashboard-en-dark-expanded.png](dashboard-en-dark-expanded.png) | `/en/dashboard`, dark, expanded | 1440 × 900 CSS pixels |
| [dashboard-en-dark-collapsed.png](dashboard-en-dark-collapsed.png) | `/en/dashboard`, dark, collapsed | 1440 × 900 CSS pixels |
| [dashboard-ar-light-expanded.png](dashboard-ar-light-expanded.png) | `/ar/dashboard`, light, expanded | 1440 × 900 CSS pixels |
| [dashboard-ar-dark-expanded.png](dashboard-ar-dark-expanded.png) | `/ar/dashboard`, dark, expanded | 1440 × 900 CSS pixels |
| [dashboard-ar-light-mobile-menu.png](dashboard-ar-light-mobile-menu.png) | `/ar/dashboard`, light, open mobile drawer | 390 × 844 CSS pixels |
| [dashboard-en-pos-light.png](dashboard-en-pos-light.png) | `/en/dashboard/pos`, light, empty cart | 1440 × 900 CSS pixels |
| [dashboard-en-pos-synthetic-cart-light.png](dashboard-en-pos-synthetic-cart-light.png) | `/en/dashboard/pos`, light, one synthetic product in cart; developer performance overlay obscures part of the submit button | 1440 × 900 CSS pixels |
| [dashboard-en-kds-light.png](dashboard-en-kds-light.png) | `/en/dashboard/kds`, light, empty queue | 1440 × 900 CSS pixels |
| [dashboard-en-orders-light.png](dashboard-en-orders-light.png) | `/en/dashboard/orders`, light, empty list | 1440 × 900 CSS pixels |
| [dashboard-en-products-light.png](dashboard-en-products-light.png) | `/en/dashboard/products`, light, empty list | 1440 × 900 CSS pixels |
| [dashboard-en-inventory-light.png](dashboard-en-inventory-light.png) | `/en/dashboard/inventory`, light, empty stock | 1440 × 900 CSS pixels |
| [dashboard-en-hr-employees-light.png](dashboard-en-hr-employees-light.png) | `/en/dashboard/hr/employees`, light, empty list | 1440 × 900 CSS pixels |
| [dashboard-en-reports-light.png](dashboard-en-reports-light.png) | `/en/dashboard/reports`, light, reports center | 1440 × 900 CSS pixels |
| [dashboard-en-settings-light.png](dashboard-en-settings-light.png) | `/en/dashboard/settings`, light, account settings | 1440 × 900 CSS pixels |
| [dashboard-en-restricted-light.png](dashboard-en-restricted-light.png) | `/en/dashboard`, light, synthetic staff user with no `user_roles` assignments | 1440 × 900 CSS pixels |
| [dashboard-en-appearance-light.png](dashboard-en-appearance-light.png) | `/en/dashboard/settings/appearance`, light, current personal color settings | 1440 × 900 CSS pixels |
| [dashboard-en-site-theme-builder-light.png](dashboard-en-site-theme-builder-light.png) | `/en/dashboard/theme-builder`, light, synthetic storefront theme | 1440 × 900 CSS pixels |
| [baseline-ar-pos-light.png](baseline-ar-pos-light.png) | `/ar/dashboard/pos`, light, synthetic product | 1440 × 900 CSS pixels |
| [baseline-ar-kds-light.png](baseline-ar-kds-light.png) | `/ar/dashboard/kds`, light, empty queue | 1440 × 900 CSS pixels |
| [baseline-ar-pos-dark.png](baseline-ar-pos-dark.png) | `/ar/dashboard/pos`, dark, synthetic product | 1440 × 900 CSS pixels |
| [baseline-ar-kds-dark.png](baseline-ar-kds-dark.png) | `/ar/dashboard/kds`, dark, empty queue | 1440 × 900 CSS pixels |
| [baseline-en-pos-dark.png](baseline-en-pos-dark.png) | `/en/dashboard/pos`, dark, synthetic product | 1440 × 900 CSS pixels |
| [baseline-en-kds-dark.png](baseline-en-kds-dark.png) | `/en/dashboard/kds`, dark, empty queue | 1440 × 900 CSS pixels |
| [baseline-en-manager-light.png](baseline-en-manager-light.png) | `/en/dashboard`, light, synthetic Manager role | 1440 × 900 CSS pixels |
| [baseline-en-employee-light.png](baseline-en-employee-light.png) | `/en/dashboard`, light, synthetic Employee role | 1440 × 900 CSS pixels |

The captures are visual baseline evidence, not proof that every control, permission, locale, dark mode, or screen reader works. Capture further routes, populated realistic synthetic data, restricted roles, and remaining viewport/theme combinations before approving a future rebrand. Physical devices were not used.

For the populated POS example, a synthetic category and product were inserted only into the fresh disposable test database; the product was selected in the actual POS UI. No order was submitted. This image is discovery evidence and has not replaced any marketing image.

The separate synthetic staff account was inserted into the same disposable tenant with zero `user_roles` rows and logged in through the API in a separate browser context. Its observed sidebar contained Dashboard and Settings, while the dashboard displayed a partial-data warning. The Theme Builder shortcut was still visible in the dashboard content; its authorization outcome was not tested. This single artificial role state does not prove that protected routes or APIs enforce access correctly.

Two additional synthetic accounts were assigned the existing tenant-scoped Manager and Employee system roles and used separate browser contexts. The manager sidebar showed POS, KDS, orders, products, HR, financial reports, reports, theme builder and settings; the employee sidebar showed KDS, orders, products, HR and settings. These are observations of this fresh fixture's seeded permissions, not a complete role/action authorization test.

In the Arabic KDS capture, the Completed filter remained English; verify localization coverage in the future implementation. Attempts to switch dark mode by editing the browser's local theme key before navigation returned to light after preferences loaded. The dark POS/KDS captures were instead taken after switching the real POS control; dark mode remained active across the subsequent POS/KDS and locale navigations during that session. This observation does not cover a fresh login or all theme settings.
