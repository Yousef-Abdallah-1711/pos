// RUN-005 / PH-06 (F-074 root-cause confirmation): for the customer, supplier, account and contract forms, submit an over-long value through the real UI and record (a) toast text, (b) inline field messages,
// (c) unhandled promise rejections / page errors (listener on pageerror + window unhandledrejection), (d) the Next.js dev error overlay text. Dev build (isolated stack): the overlay is a dev-only artefact; the unhandled rejection is not.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
const { ctx, page } = await openChrome();
const errs = []; page.on('pageerror', (e) => errs.push('pageerror: ' + String(e.message).slice(0, 100))); page.on('console', (m) => { if (m.type() === 'error') errs.push('console: ' + m.text().slice(0, 100)); });
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && !['GET', 'OPTIONS'].includes(r.request().method())) { let b = ''; try { b = (await r.text()).slice(0, 140); } catch {} calls.push(`${r.request().method()} ${r.url().replace('http://localhost:18082/api/v1', '')} ${r.status()} ${b}`); } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard`);
const TARGETS = [['customers', 'accounting/customers', true], ['suppliers', 'accounting/suppliers', true], ['accounts', 'accounting/chart-of-accounts', true], ['contracts', 'hr/contracts', true]];
const BASE = 'input:visible:not([type=hidden]):not([type=checkbox]):not([type=radio]):not([type=file]):not([type=search]):not([placeholder*="earch" i]), textarea:visible';
const out = [];
for (const [key, route, add] of TARGETS) {
  errs.length = 0; calls.length = 0;
  await page.goto(`${DASH}/en/dashboard/${route}`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(4500);
  await page.evaluate(() => { window.__unhandled = []; window.addEventListener('unhandledrejection', (e) => window.__unhandled.push(String(e.reason?.message || e.reason).slice(0, 100))); });
  if (add) { const b = page.locator('button:visible').filter({ hasText: /^(\+\s*)?(add|new|create|record|invite|request)/i }).first(); if (await b.count()) { await b.click().catch(() => {}); await page.waitForTimeout(1800); } }
  const scope = (await page.locator('[role=dialog]:visible').count()) > 0 ? BASE.split(', ').map((x) => `[role=dialog] ${x}`).join(', ') : BASE; const loc = page.locator(scope); const n = await loc.count();
  const tag = Date.now().toString().slice(-6);
  for (let i = 0; i < n; i++) { const e = loc.nth(i); const ty = await e.getAttribute('type'); const ph = (await e.getAttribute('placeholder')) || ''; let v = /@/.test(ph) || ty === 'email' ? `r5.${tag}@example.test` : ty === 'tel' || /\+20/.test(ph) ? '+201000000999' : ty === 'number' ? '12' : ty === 'date' ? '2026-11-20' : (i === 0 || i === 1) ? 'X'.repeat(300) : `R5 ${tag}`; await e.fill(v, { timeout: 3000 }).catch(() => {}); }
  const sel = await page.locator('select:visible').count(); for (let k = 0; k < sel; k++) await page.locator('select:visible').nth(k).selectOption({ index: 1 }).catch(() => {});
  await page.locator('button:visible').filter({ hasText: /^(save|create|add|submit|confirm)/i }).last().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(1200);
  const toast = (await page.locator('[data-sonner-toast], [role=status], li[data-type]').allInnerTexts().catch(() => [])).join(' | ').slice(0, 220);
  await page.waitForTimeout(2500); const overlay = await page.evaluate(() => { const o = document.querySelector('nextjs-portal'); const t = o?.shadowRoot?.textContent || ''; return t.slice(0, 140); }).catch(() => '');
  const unhandled = await page.evaluate(() => window.__unhandled || []).catch(() => []);
  await page.screenshot({ path: `${d.shot}/AUD-R5-F074-${key}.png` }).catch(() => {});
  out.push({ form: key, request: calls.slice(0, 2), toast, inlineAlerts: (await page.locator('[role=alert]').allInnerTexts().catch(() => [])).slice(0, 3), devOverlay: overlay, unhandledRejections: unhandled, pageErrors: errs.slice(0, 4) });
  console.log(JSON.stringify(out[out.length - 1]));
}
fs.writeFileSync(`${d.net}/AUD-R5-F074-UNHANDLED.json`, JSON.stringify({ run: 'RUN-005', finding: 'F-074', build: 'Next dev (isolated stack)', out }, null, 1));
await ctx.close();
