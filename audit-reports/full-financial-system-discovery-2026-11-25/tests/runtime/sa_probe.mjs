import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
const EV = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25/evidence';
const b = await chromium.launch({ headless: true }); const ctx = await b.newContext({ storageState: `${EV}/state-SA.json`, viewport: { width: 1440, height: 900 } }); const p = await ctx.newPage();
const net = []; p.on('response', r => { if (r.url().includes(':18082') && r.request().method() !== 'GET' && r.request().method() !== 'OPTIONS') net.push(`${r.request().method()} ${r.url().replace('http://localhost:18082/api/v1', '')} ${r.status()}`); });
for (const path of ['dashboard/tenants/1', 'dashboard/tenants/3']) {
  await p.goto('http://localhost:13114/' + path, { waitUntil: 'domcontentloaded', timeout: 60000 }); await p.waitForTimeout(2500);
  console.log('==', path, (await p.locator('body').innerText()).split(String.fromCharCode(10)).join(' | ').slice(0, 520));
  console.log('BTN', (await p.$$eval('button', els => els.map(e => e.innerText.trim().slice(0, 25)).filter(Boolean))).join(' || '));
}
console.log('writes', net.join('; '));
await b.close();
