// RUN-003 / PH-04/06: generic real-UI form submission probe in HEADED Chrome (separate non-persistent context). For each route: open, click an Add/New/Create button when no form is visible,
// fill visible text/number/email/tel/date inputs with synthetic values, submit with a Save/Create/Add/Submit button, and record every mutating network call + visible result. Tenant 4 (diagnostic). Synthetic data only.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, DASH, PASSWORD, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
const TARGETS = [['categories/new', 'category'], ['products/new', 'product'], ['drivers/new', 'driver'], ['accounting/customers', 'customer'], ['accounting/suppliers', 'supplier'], ['accounting/bank-accounts', 'bank account'], ['tables', 'table']];
const browser = await chromium.launch({ channel: 'chrome', headless: false });
const page = await (await browser.newContext({ viewport: { width: 1440, height: 900 } })).newPage();
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && !['GET', 'OPTIONS'].includes(r.request().method())) { let b = ''; try { b = (await r.text()).slice(0, 220); } catch { /* n/a */ } calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').replace(/"password":"[^"]*"/g, '"password":"<redacted>"').slice(0, 300), b }); } });
await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded', timeout: 180000 });
for (let i = 0; i < 40; i++) { await page.waitForTimeout(3000); if ((await page.locator('input[type=email]').count()) > 0) break; }
await page.fill('input[type=email]', 'audit.c.owner@example.test'); await page.fill('input[type=password]', PASSWORD); await page.click('button[type=submit]'); await page.waitForTimeout(8000);
const results = [];
for (const [route, label] of TARGETS) {
  const before = calls.length;
  await page.goto(`${DASH}/en/dashboard/${route}`, { waitUntil: 'domcontentloaded', timeout: 180000 }).catch(() => {}); await page.waitForTimeout(6000);
  const visibleInputs = async () => page.locator('input:visible:not([type=hidden]):not([type=checkbox]):not([type=radio]):not([type=file]):not([type=search]), textarea:visible').count();
  if ((await visibleInputs()) === 0) { const add = page.locator('button:visible').filter({ hasText: /^(\+\s*)?(add|new|create)/i }).first(); if (await add.count()) { await add.click({ timeout: 8000 }).catch(() => {}); await page.waitForTimeout(2500); } }
  const fields = await page.locator('input:visible:not([type=hidden]):not([type=checkbox]):not([type=radio]):not([type=file]):not([type=search]), textarea:visible').evaluateAll((els) => els.map((e, i) => ({ i, type: e.type || 'textarea', name: e.getAttribute('name') || e.id || e.getAttribute('placeholder') || '', label: (e.labels?.[0]?.textContent || e.getAttribute('aria-label') || '').trim().slice(0, 30) })));
  const filled = [];
  const loc = page.locator('input:visible:not([type=hidden]):not([type=checkbox]):not([type=radio]):not([type=file]):not([type=search]), textarea:visible');
  for (const f of fields) {
    const key = (f.name + ' ' + f.label).toLowerCase(); let v = 'Audit UI ' + label;
    if (f.type === 'email' || /email/.test(key)) v = `ui.${label.replace(/\W/g, '')}@example.test`; else if (f.type === 'tel' || /phone|mobile/.test(key)) v = '+201000000555'; else if (f.type === 'number' || /price|cost|amount|stock|quantity|capacity|balance|salary|rate/.test(key)) v = '12'; else if (f.type === 'date') v = '2026-11-26'; else if (/iban|account.?number|code|sku|number/.test(key)) v = 'UI' + Date.now().toString().slice(-6);
    await loc.nth(f.i).fill(v, { timeout: 5000 }).then(() => filled.push(`${f.name || f.label || f.type}=${v}`)).catch(() => {});
  }
  await page.screenshot({ path: `${d.shot}/AUD-UIFORM-${route.replace(/\//g, '_')}-filled.png`, timeout: 30000 }).catch(() => {});
  const submit = page.locator('button:visible').filter({ hasText: /^(save|create|add|submit|confirm)/i }).last();
  let clicked = false; if (await submit.count()) { clicked = await submit.click({ timeout: 8000 }).then(() => true).catch(() => false); }
  await page.waitForTimeout(5000);
  await page.screenshot({ path: `${d.shot}/AUD-UIFORM-${route.replace(/\//g, '_')}-after.png`, timeout: 30000 }).catch(() => {});
  const toast = (await page.locator('[role=status], [data-sonner-toast], [role=alert]').allInnerTexts().catch(() => [])).join(' | ').slice(0, 200);
  results.push({ route, label, fieldCount: fields.length, filled: filled.slice(0, 14), submitClicked: clicked, toast, mutations: calls.slice(before) });
  console.log(route, 'fields', fields.length, 'clicked', clicked, 'toast', toast.slice(0, 80), 'mutations', calls.slice(before).map((c) => `${c.m} ${c.u} ${c.s}`).join(', '));
}
fs.writeFileSync(`${d.net}/AUD-UIFORMS.json`, JSON.stringify({ run: 'RUN-003', results }, null, 1));
await browser.close();
