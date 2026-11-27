// RUN-005 / PH-09: customer invoice form through the real UI (headed Chrome, separate profile). Creates a draft invoice: 1 x 5,000 with 10% discount and 14% VAT, captures the exact request payload (F-024 float artefacts, F-007 tax base),
// then reads the stored invoice and the journal state (F-008 posts while draft). Tenant 4 fixture, synthetic.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, DASH, PASSWORD, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-09');
const profile = 'C:/Users/Yousef/AppData/Local/Temp/r5-chrome-profile-inv';
const ctx = await chromium.launchPersistentContext(profile, { channel: 'chrome', headless: false, viewport: { width: 1440, height: 900 } });
const page = ctx.pages()[0] ?? (await ctx.newPage());
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && !['GET', 'OPTIONS'].includes(r.request().method())) { let b = ''; try { b = (await r.text()).slice(0, 300); } catch {} calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: r.request().postData() || '', b }); } });
await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded', timeout: 180000 });
for (let i = 0; i < 40; i++) { await page.waitForTimeout(3000); if ((await page.locator('input[type=email]').count()) > 0 || (await page.evaluate(() => (document.body.innerText || '').length).catch(() => 0)) > 800) break; }
if ((await page.locator('input[type=email]').count()) > 0) { await page.fill('input[type=email]', 'audit.c.owner@example.test'); await page.fill('input[type=password]', PASSWORD); await page.click('main button[type=submit], form button[type=submit]').catch(() => page.click('button[type=submit]')); await page.waitForTimeout(9000); }
// create a valid customer through the real customer form (receivable account chosen explicitly: 1120)
await page.goto(`${DASH}/en/dashboard/accounting/customers`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(5000);
await page.locator('button:visible').filter({ hasText: /^(\+\s*)?(add|new|create)/i }).first().click({ timeout: 8000 }).catch(() => {}); await page.waitForTimeout(1800);
const CN = '000 R5 Invoice Customer ' + Date.now().toString().slice(-5);
const cf = page.locator('[role=dialog]:visible input:not([type=search]):not([placeholder*="earch" i])'); await cf.nth(0).fill(CN).catch(() => {});
const racc = page.locator('[role=dialog]:visible select').filter({ has: page.locator('option', { hasText: 'Select receivable account' }) }).first(); const rl = await racc.locator('option').allInnerTexts(); const ri = rl.findIndex((t) => /1120/.test(t)); console.log('receivable options', rl.length, 'chosen', ri, rl[ri], '| first options', JSON.stringify(rl.slice(0, 6))); if (ri > 0) await racc.selectOption({ index: ri });
const cb = calls.length; await page.locator('button:visible').filter({ hasText: /^(save|create|add)/i }).last().click({ timeout: 8000 }).catch(() => {}); await page.waitForTimeout(4000);
console.log('customer create', calls.slice(cb).map((c) => c.m + ' ' + c.u + ' ' + c.s).join(','), 'db', sql(`select id||'|'||receivable_account_id from customers_accounting where customer_name='${CN}'`));
await page.goto(`${DASH}/en/dashboard/accounting/customer-invoices`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(6000);
await page.locator('button:visible').filter({ hasText: /new customer invoice/i }).first().click({ timeout: 10000 }).catch((e) => console.log('open failed', e.message.slice(0, 60))); await page.waitForTimeout(3000);
const dump = await page.evaluate(() => ({ inputs: [...document.querySelectorAll('[role=dialog] input, [role=dialog] select, [role=dialog] textarea, form input, form select')].filter((e) => e.offsetParent).map((e, i) => `${i}|${e.tagName}|${e.type}|${e.name || e.id}|${e.placeholder || ''}|${(e.labels?.[0]?.textContent || '').trim().slice(0, 30)}`), buttons: [...document.querySelectorAll('[role=dialog] button, form button')].map((b) => (b.innerText || '').trim()).filter(Boolean).slice(0, 14) }));
fs.writeFileSync(`${d.net}/AUD-R5-INVOICE-FORM-DUMP.json`, JSON.stringify(dump, null, 1));
const form = page.locator('[role=dialog]:visible, form').first(); const inp = form.locator('input, select, textarea').filter({ visible: true });
await inp.nth(0).fill('2026-11-27').catch(() => {}); await inp.nth(1).fill('2026-12-27').catch(() => {});
const custLabels = await inp.nth(2).locator('option').allInnerTexts(); const ci = custLabels.findIndex((t) => /000 R5 Invoice Customer/.test(t)); console.log('customer options count', custLabels.length, 'customers in DB', sql("select count(*) from customers_accounting where tenant_id=4 and is_active"), 'first', custLabels.slice(0, 4).join(' / '), 'chosen', ci); await inp.nth(2).selectOption({ index: ci > 0 ? ci : 1 });
await inp.nth(6).fill('R5 float probe line').catch(() => {}); await inp.nth(7).fill('بند اختبار').catch(() => {});
await inp.nth(9).fill('1').catch(() => {}); await inp.nth(10).fill('5000').catch(() => {}); await inp.nth(11).fill('14').catch(() => {}); await inp.nth(12).fill('10').catch(() => {});
const accOpts = await inp.nth(8).locator('option').allInnerTexts().catch(() => []); const accIdx = accOpts.findIndex((t) => /4100|food sales|sales/i.test(t)); if (accIdx > 0) await inp.nth(8).selectOption({ index: accIdx }).catch(() => {}); console.log('account options', accOpts.length, 'chosen', accIdx, accOpts[accIdx]);
await page.waitForTimeout(800); await page.screenshot({ path: `${d.shot}/AUD-R5-INVOICE-FORM-filled.png` }).catch(() => {});
const shown = await page.locator('[role=dialog], form').first().innerText().catch(() => ''); const m = shown.split(String.fromCharCode(10)).filter((l) => /tax|total|subtotal|vat/i.test(l));
const before = calls.length; await page.locator('button:visible').filter({ hasText: /create customer invoice/i }).last().click({ timeout: 8000 }).catch((e) => console.log('submit failed', e.message.slice(0, 60))); await page.waitForTimeout(6000);
const c = calls.slice(before)[0]; let payload = null; try { payload = JSON.parse(c.req); } catch { /* n/a */ }
const line = payload?.lines?.[0] || payload?.items?.[0] || null;
const out = { ui_summary_lines: m.slice(0, 6), request: c ? `${c.m} ${c.u} ${c.s}` : 'none', response_head: c?.b?.slice(0, 200), line_payload: line, payload_top: payload ? Object.fromEntries(Object.entries(payload).filter(([k]) => k !== 'lines' && k !== 'items')) : null, raw_has_float_artefact: c ? /\d\.\d{8,}/.test(c.req) : null };
const invNo = (c?.b || '').match(/INV-[0-9-]+/)?.[0]; if (invNo) { out.db = sql(`select status||'|sub '||subtotal||'|tax '||tax_amount||'|total '||total_amount||'|je '||coalesce(journal_entry_id::text,'-') from invoices_customer where invoice_number='${invNo}'`); out.journal = sql(`select string_agg(a.account_code||' D'||l.debit_amount||' C'||l.credit_amount,' ; ') from invoices_customer i join journal_entry_lines l on l.journal_entry_id=i.journal_entry_id join chart_of_accounts a on a.id=l.account_id where i.invoice_number='${invNo}'`); }
console.log(JSON.stringify(out, null, 1)); fs.writeFileSync(`${d.net}/AUD-R5-INVOICE-FORM-RESULT.json`, JSON.stringify({ run: 'RUN-005', findings: ['F-024', 'F-007', 'F-008'], out, calls: calls.slice(before) }, null, 1));
await ctx.close();
