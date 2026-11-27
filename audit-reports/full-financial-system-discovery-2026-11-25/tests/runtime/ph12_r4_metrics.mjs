// RUN-004 / PH-12: headed-Chrome capture of PROF-007 (reports), DASH-004 (dashboard home), KDS-003 (kds analytics), HRR-002/007 (hr reports). Tenant 4 owner, read-only GETs. Evidence JSON holds every report API body seen by the page.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
const { ctx, page } = await openChrome();
const bodies = {};
page.on('response', async (r) => { const u = r.url(); if (u.includes(':18082') && r.request().method() === 'GET' && /report|kds|admin\/orders|products|hr\/|tables\/stats|dashboard/.test(u)) { try { bodies[u.replace('http://localhost:18082/api/v1', '')] = (await r.text()).slice(0, 20000); } catch {} } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard`);
const vis = {};
for (const [id, rt] of [['dash', ''], ['reports', 'reports'], ['kds', 'kds'], ['hrreports', 'hr/reports']]) {
  await page.goto(`${DASH}/en/dashboard${rt ? '/' + rt : ''}`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(12000);
  await page.screenshot({ path: `${d.shot}/AUD-R4-MET-${id}.png`, fullPage: true }).catch(() => {});
  vis[id] = (await page.locator('body').innerText()).slice(0, 4000);
}
fs.writeFileSync(`${d.net}/AUD-R4-MET-CAPTURE.json`, JSON.stringify({ run: 'RUN-004', tenant: 4, visible: vis, bodies }, null, 1));
for (const [k, v] of Object.entries(bodies)) console.log(k, '=>', v.slice(0, 200).replace(/\n/g, ' '));
await ctx.close();
