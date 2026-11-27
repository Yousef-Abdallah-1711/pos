// RUN-003 / PH-12 / AUD-MET-AGE: receivables & payables aging in HEADED Chrome, tenant 4 owner, read-only.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
const { ctx, page } = await openChrome();
const bodies = {}; page.on('response', async (r) => { if (/aging/.test(r.url()) && r.url().includes(':18082')) { try { bodies[r.url().replace('http://localhost:18082/api/v1', '')] = await r.text(); } catch {} } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/reports/financial/receivables-aging`);
await page.waitForTimeout(7000); await page.screenshot({ path: `${d.shot}/AUD-MET-AGE-receivables.png`, fullPage: true });
const ar = (await page.locator('main').innerText().catch(() => '')).slice(0, 1500);
await page.goto(`${DASH}/en/dashboard/reports/financial/payables-aging`, { waitUntil: 'domcontentloaded', timeout: 120000 }); await page.waitForTimeout(7000);
await page.screenshot({ path: `${d.shot}/AUD-MET-AGE-payables.png`, fullPage: true });
const ap = (await page.locator('main').innerText().catch(() => '')).slice(0, 1500);
fs.writeFileSync(`${d.net}/AUD-MET-AGE.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, ar, ap, bodies }, null, 1));
console.log('AR UI:\n' + ar.slice(0, 700)); console.log('AP UI:\n' + ap.slice(0, 700)); console.log(Object.entries(bodies).map(([k, v]) => k + ' => ' + v.slice(0, 600)).join('\n'));
await ctx.close();
