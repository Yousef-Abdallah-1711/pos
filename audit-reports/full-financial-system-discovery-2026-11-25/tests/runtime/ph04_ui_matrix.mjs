// RUN-003 / PH-04: dashboard route matrix in HEADED Chrome: EN desktop, AR (RTL) desktop, EN mobile (390x844). Read-only.
// Records per route: final URL, load error text (404/500 shells), console errors, horizontal overflow, <html dir>, visible text length, page title. Screenshots for AR and mobile.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-04');
const ROUTES = ['', 'accounting/advanced', 'accounting/bank-accounts', 'accounting/chart-of-accounts', 'accounting/customer-invoices', 'accounting/customers', 'accounting/general-ledger', 'accounting/journal-entries', 'accounting/supplier-invoices', 'accounting/suppliers', 'categories', 'drivers', 'hr/attendance', 'hr/contracts', 'hr/employees', 'hr/insurance', 'hr/leaves', 'hr/payroll', 'hr/performance', 'hr/resignations', 'hr/training', 'inventory', 'kds', 'kds/settings', 'notifications', 'orders', 'pos', 'pos/settings', 'products', 'products/new', 'rbac/roles', 'reports', 'reports/financial', 'reports/financial/advanced', 'reports/financial/balance-sheet', 'reports/financial/cash-flow', 'reports/financial/inventory-valuation', 'reports/financial/payables-aging', 'reports/financial/profit-loss', 'reports/financial/receivables-aging', 'reports/financial/vat', 'reports/hr', 'reports/inventory', 'reports/overview', 'reports/profitability', 'reports/sales', 'settings', 'settings/access-control', 'settings/account', 'settings/appearance', 'settings/billing', 'settings/restaurant', 'settings/roles', 'settings/users', 'tables', 'tables/layout', 'theme-builder', 'theme-builder/presets', 'theme-builder/qr'];
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard`);
const results = []; let relogins = 0;
const modes = [['en', { width: 1440, height: 900 }, false], ['ar', { width: 1440, height: 900 }, true], ['en-mobile', { width: 390, height: 844 }, true]];
const ONLY = process.env.MODES ? process.env.MODES.split(',') : null;
for (const [mode, vp, shot] of modes) {
  if (ONLY && !ONLY.includes(mode)) continue;
  await page.setViewportSize(vp);
  const lang = mode.startsWith('ar') ? 'ar' : 'en';
  for (const route of ROUTES) {
    const cons = []; const onc = (m) => { if (m.type() === 'error') cons.push(m.text().slice(0, 120)); }; page.on('console', onc);
    const url = `${DASH}/${lang}/dashboard${route ? '/' + route : ''}`;
    let status = 'ok';
    try { await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 120000 }); } catch (e) { status = 'goto-error ' + String(e).slice(0, 80); }
    await page.waitForTimeout(3500);
    if ((await page.locator('input[type=email]').count()) > 0 || page.url().includes('/auth/login')) { relogins++; await ensureLoggedIn(page, 'audit.c.owner@example.test', url); await page.waitForTimeout(2500); }
    const info = await page.evaluate(() => ({ url: location.href, dir: document.documentElement.dir, lang: document.documentElement.lang, overflow: document.documentElement.scrollWidth > window.innerWidth + 2, sw: document.documentElement.scrollWidth, iw: window.innerWidth, len: (document.body.innerText || '').length, is404: /This page could not be found|404/.test(document.body.innerText.slice(0, 200)), hasError: /Application error|Internal Server Error|Something went wrong/i.test(document.body.innerText), title: document.title.slice(0, 60) })).catch(() => ({ url: 'eval-error' }));
    page.off('console', onc);
    const rec = { mode, route: route || '(home)', status, ...info, consoleErrors: cons.length, consoleSample: cons.slice(0, 2) };
    if (shot) await page.screenshot({ path: `${d.shot}/AUD-UI-${mode}-${(route || 'home').replace(/\//g, '_')}.png`, timeout: 30000 }).catch(() => {});
    results.push(rec); console.log(mode, route || '(home)', info.dir, 'overflow=' + info.overflow, 'len=' + info.len, info.is404 ? '404' : '', info.hasError ? 'ERR' : '', 'cons=' + cons.length);
  }
}
fs.writeFileSync(`${d.net}/AUD-UI-MATRIX${process.env.OUT_SUFFIX || ''}.json`, JSON.stringify({ run: 'RUN-003', routes: ROUTES.length, relogins, results }, null, 1));
await ctx.close();
