// RUN-003 / PH-12: operational report metrics (sales, profitability, inventory turnover, overview, order stats) in HEADED Chrome, tenant 4, read-only.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
const { ctx, page } = await openChrome();
const bodies = {}; page.on('response', async (r) => { const u = r.url(); if (u.includes(':18082') && r.request().method() === 'GET' && /\/reports\/|orders\/stats|kds\/analytics/.test(u)) { try { bodies[u.replace('http://localhost:18082/api/v1', '')] = await r.text(); } catch {} } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/reports/overview`);
const visible = {};
for (const rt of ['overview', 'sales', 'profitability', 'inventory']) {
  await page.goto(`${DASH}/en/dashboard/reports/${rt}`, { waitUntil: 'domcontentloaded', timeout: 120000 }); await page.waitForTimeout(9000);
  await page.screenshot({ path: `${d.shot}/AUD-MET-OPS-${rt}.png`, fullPage: true });
  visible[rt] = (await page.locator('main').innerText()).slice(0, 2500);
}
await page.goto(`${DASH}/en/dashboard/orders`, { waitUntil: 'domcontentloaded', timeout: 120000 }); await page.waitForTimeout(6000);
await page.screenshot({ path: `${d.shot}/AUD-MET-OPS-orders.png`, fullPage: true });
const tabs = [], text = '';
fs.writeFileSync(`${d.net}/AUD-MET-OPS.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, visible, bodies }, null, 1));
console.log(Object.entries(bodies).map(([k, v]) => k + ' => ' + v.slice(0, 500)).join('\n\n'));
await ctx.close();
