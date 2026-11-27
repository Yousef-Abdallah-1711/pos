// RUN-003 / PH-12: general-ledger page (tabs, trial balance) in HEADED Chrome, tenant 4 owner, read-only.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
const { ctx, page } = await openChrome();
const bodies = {}; page.on('response', async (r) => { const u = r.url(); if (u.includes(':18082') && r.request().method() === 'GET' && /trial|ledger|account/.test(u)) { try { bodies[u.replace('http://localhost:18082/api/v1', '')] = (await r.text()).slice(0, 5000); } catch {} } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/accounting/general-ledger`);
await page.waitForTimeout(8000);
const tabs = await page.locator('button,[role=tab]').evaluateAll((e) => e.map((x) => x.innerText.trim()).filter((t) => t && t.length < 40));
await page.screenshot({ path: `${d.shot}/AUD-MET-GL-page.png`, fullPage: true });
for (const t of ['Trial Balance', 'Trial balance']) { const b = page.getByText(t, { exact: true }).first(); if (await b.count()) { await b.click(); await page.waitForTimeout(5000); await page.screenshot({ path: `${d.shot}/AUD-MET-TB-tab.png`, fullPage: true }); break; } }
const txt = (await page.locator('body').innerText()).slice(0, 3000);
fs.writeFileSync(`${d.net}/AUD-MET-GL-TB.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, tabs, txt, bodies }, null, 1));
console.log(tabs.slice(0, 30)); console.log(txt.slice(-900)); console.log(Object.keys(bodies));
for (const [k, v] of Object.entries(bodies)) if (/trial/.test(k)) console.log(k, v.slice(0, 600));
await ctx.close();
