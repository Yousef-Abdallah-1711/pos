import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import fs from 'fs';
const EV = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25/evidence';
const b = await chromium.launch({ headless: true }); const p = await (await b.newContext({ viewport: { width: 1280, height: 900 } })).newPage();
const net = []; p.on('response', async r => { const u = r.url(); if (u.includes(':18082') && r.request().method() !== 'GET' && r.request().method() !== 'OPTIONS') { let body = ''; try { body = (await r.text()).slice(0, 300); } catch {} net.push(`${r.request().method()} ${u.replace('http://localhost:18082/api/v1', '')} ${r.status()} req=${(r.request().postData() || '').slice(0, 300)} resp=${body}`); } });
await p.goto('http://audit-c-diag.localhost:13113/en/menu', { waitUntil: 'networkidle' }); await p.waitForTimeout(3000);
console.log('MENU:', (await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 500));
const btns = await p.$$eval('button', els => els.map(e => (e.innerText || e.getAttribute('aria-label') || '').trim().slice(0, 30)).filter(Boolean)); console.log('BTN', btns.join(' || '));
await p.screenshot({ path: `${EV}/screenshots/store-C-menu.png` });
const add = p.locator('button:has-text("Add")').first(); if (await add.count()) { await add.click(); await p.waitForTimeout(1500); }
await p.goto('http://audit-c-diag.localhost:13113/en/cart', { waitUntil: 'networkidle' }).catch(() => {}); await p.waitForTimeout(2500);
console.log('CART:', (await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 700));
const f = await p.$$eval('input,select,textarea,button', els => els.filter(e => e.offsetParent !== null).map(e => `${e.tagName}[${e.type}] name=${e.name} ph=${e.placeholder} txt=${(e.innerText || '').slice(0, 30)}`)); console.log(f.join('\n'));
await p.screenshot({ path: `${EV}/screenshots/store-C-cart.png` });
console.log(net.join('\n'));
await b.close();
