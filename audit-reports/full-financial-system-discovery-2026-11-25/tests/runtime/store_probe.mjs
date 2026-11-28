import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
const EV = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25/evidence';
const b = await chromium.launch({ headless: true }); const p = await (await b.newContext({ viewport: { width: 1280, height: 900 } })).newPage();
const net = []; p.on('response', async r => { if (r.url().includes(':18082') && r.request().method() !== 'GET') net.push(`${r.request().method()} ${r.url().replace('http://localhost:18082/api/v1', '')} ${r.status()}`); });
await p.goto('http://audit-c-diag.localhost:13113/en', { waitUntil: 'networkidle' }); await p.waitForTimeout(3000);
console.log((await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 700));
const btns = await p.$$eval('button, a', els => els.map(e => (e.innerText || e.getAttribute('aria-label') || '').trim().slice(0, 30)).filter(Boolean)); console.log('CLICKABLES', btns.slice(0, 40).join(' || '));
await p.screenshot({ path: `${EV}/screenshots/store-C-home.png` });
await b.close();
