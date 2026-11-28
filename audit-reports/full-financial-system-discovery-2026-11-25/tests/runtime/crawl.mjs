// usage: node crawl.mjs <app> <baseUrl> <stateFile|none> <label> [lang]
// Loads every route of an app from routes.json, records load outcome, API failures, page errors, control counts, key text markers.
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import fs from 'fs';
const [app, base, state, label, lang = 'en'] = process.argv.slice(2);
const ROOT = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25';
const routes = JSON.parse(fs.readFileSync(`${ROOT}/tests/runtime/routes.json`, 'utf8'))[app];
const b = await chromium.launch({ headless: true });
const ctx = await b.newContext({ viewport: { width: 1440, height: 900 }, ...(state !== 'none' ? { storageState: `${ROOT}/evidence/${state}` } : {}) });
const out = [];
for (const r of routes) {
  const url = r.replace('[locale]', lang).replace('[id]', '1').replace('[orderNumber]', 'ORD-X').replace('[tableId]', '1').replace('[slug]', 'menu-items');
  const p = await ctx.newPage(); const apiFail = []; const errs = []; let nav = 0;
  p.on('response', rs => { const u = rs.url(); if (u.includes(':18082') && rs.status() >= 400 && rs.request().method() !== 'OPTIONS') apiFail.push(`${rs.request().method()} ${u.replace('http://localhost:18082/api/v1', '').slice(0, 70)} ${rs.status()}`); });
  p.on('pageerror', e => errs.push(e.message.slice(0, 100)));
  try { const rs = await p.goto(base + url, { waitUntil: 'domcontentloaded', timeout: 45000 }); nav = rs ? rs.status() : 0; } catch (e) { nav = -1; }
  await p.waitForTimeout(3000);
  const info = await p.evaluate(() => ({ inputs: document.querySelectorAll('input,select,textarea').length, buttons: document.querySelectorAll('button').length, tables: document.querySelectorAll('table').length, links: document.querySelectorAll('a[href]').length, forms: document.querySelectorAll('form').length, dir: document.documentElement.dir, text: (document.body.innerText || '').replace(/\s+/g, ' ').slice(0, 1500) })).catch(() => ({ inputs: 0, buttons: 0, tables: 0, links: 0, forms: 0, dir: '', text: '' }));
  const t = info.text;
  const marker = /could not be found|404/i.test(t) ? '404' : /access denied|do not have permission/i.test(t) ? 'ACCESS_DENIED' : /unable to load|something went wrong|application error|failed to/i.test(t) ? 'ERROR_TEXT' : /redirecting|checking your sign-in/i.test(t) ? 'REDIRECT' : 'OK';
  out.push({ route: r, url, status: nav, finalUrl: p.url().replace(base, ''), marker, inputs: info.inputs, buttons: info.buttons, tables: info.tables, forms: info.forms, apiFail: apiFail.slice(0, 5), pageErrors: errs.slice(0, 3), snippet: t.slice(0, 160) });
  await p.close();
}
fs.writeFileSync(`${ROOT}/evidence/crawl-${label}.json`, JSON.stringify(out, null, 1));
const c = {}; out.forEach(o => c[o.marker] = (c[o.marker] || 0) + 1);
console.log(label, 'routes', out.length, JSON.stringify(c), 'withApiFail', out.filter(o => o.apiFail.length).length, 'withPageErrors', out.filter(o => o.pageErrors.length).length);
await b.close();
