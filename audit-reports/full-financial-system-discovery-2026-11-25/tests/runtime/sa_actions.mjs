import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
const EV = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25/evidence';
const b = await chromium.launch({ headless: true }); const ctx = await b.newContext({ storageState: `${EV}/state-SA.json`, viewport: { width: 1440, height: 900 } }); const p = await ctx.newPage();
const net = []; p.on('response', async r => { if (r.url().includes(':18082') && r.request().method() !== 'GET' && r.request().method() !== 'OPTIONS') { let t = ''; try { t = (await r.text()).slice(0, 160).replace(/"(token|password|access_token)":"[^"]+"/g, '"$1":"<redacted>"'); } catch {} net.push(`${r.request().method()} ${r.url().replace('http://localhost:18082/api/v1', '')} ${r.status()} ${t}`); } });
p.on('dialog', async d => { console.log('dialog:', d.type(), d.message().slice(0, 80)); await d.accept(d.type() === 'prompt' ? '7' : undefined); });
await p.goto('http://localhost:13114/dashboard/tenants/3', { waitUntil: 'domcontentloaded', timeout: 60000 }); await p.waitForTimeout(3000);
await p.click('button:has-text("Extend Trial")'); await p.waitForTimeout(2000);
const inputs = await p.$$eval('input:visible', els => els.map(e => `${e.type} ph=${e.placeholder} val=${e.value}`)); console.log('extend dialog inputs', inputs.join(' | '));
const sub = p.locator('button:has-text("Extend"), button:has-text("Confirm"), button:has-text("Save")').last(); if (await sub.count()) { await sub.click().catch(() => {}); await p.waitForTimeout(2500); }
console.log('after extend:', (await p.locator('body').innerText()).split(String.fromCharCode(10)).join(' | ').match(/Trial Days Remaining[^|]*\|[^|]*\|[^|]*/)?.[0]);
await p.click('button:has-text("Impersonate")'); await p.waitForTimeout(4000);
console.log('after impersonate URL', p.url().slice(0, 80));
console.log(net.join('\n'));
await b.close();
