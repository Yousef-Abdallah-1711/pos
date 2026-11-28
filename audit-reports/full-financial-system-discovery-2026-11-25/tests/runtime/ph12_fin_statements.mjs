// RUN-003 / PH-12: trial balance, ledger, P&L, cash flow in HEADED Chrome, tenant 4 owner, read-only.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
const { ctx, page } = await openChrome();
const bodies = {}; page.on('response', async (r) => { const u = r.url(); if (u.includes(':18082') && r.request().method() === 'GET' && /trial-balance|profit-loss|cash-flow|\/ledger/.test(u)) { try { bodies[u.replace('http://localhost:18082/api/v1', '')] = (await r.text()).slice(0, 12000); } catch {} } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/accounting`);
const vis = {};
for (const [id, rt] of [['tb', 'accounting'], ['ledger', 'accounting/ledger'], ['pnl', 'reports/financial/profit-loss'], ['cf', 'reports/financial/cash-flow']]) {
  await page.goto(`${DASH}/en/dashboard/${rt}`, { waitUntil: 'domcontentloaded', timeout: 120000 }); await page.waitForTimeout(9000);
  await page.screenshot({ path: `${d.shot}/AUD-MET-FIN-${id}.png`, fullPage: true });
  vis[id] = (await page.locator('body').innerText()).slice(0, 2500);
}
fs.writeFileSync(`${d.net}/AUD-MET-FIN.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, visible: vis, bodies }, null, 1));
for (const [k, v] of Object.entries(bodies)) console.log(k, '=>', v.slice(0, 450).replace(/\n/g, ' '));
await ctx.close();
