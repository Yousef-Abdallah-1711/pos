// RUN-004 / PH-04+06: mode-based real-UI form matrix in HEADED Chrome. For each target form and each mode (empty, invalid, arabic, long, valid) it opens the route, opens the add dialog if needed,
// fills visible fields with mode-specific values chosen by field type/name, submits, and records validation messages, dialog state, every mutating API call and a DB check for the marker.
// Tenant 4 (diagnostic fixture), synthetic data only. Env: ONLY=<comma list of route keys> to restrict, MODES=<comma list>.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
const ONLY = (process.env.ONLY || '').split(',').filter(Boolean); const MODES = (process.env.MODES || 'empty,invalid,arabic,long,valid').split(',');
const TARGETS = [
  { key: 'customers', route: 'accounting/customers', label: 'customer', table: 'customers_accounting', col: 'customer_name' },
  { key: 'suppliers', route: 'accounting/suppliers', label: 'supplier', table: 'suppliers', col: 'name' },
  { key: 'accounts', route: 'accounting/chart-of-accounts', label: 'account', table: 'chart_of_accounts', col: 'account_name' },
  { key: 'journal', route: 'accounting/journal-entries', label: 'journal', table: 'journal_entries', col: 'description' },
  { key: 'employees', route: 'hr/employees', label: 'employee', table: 'employees', col: 'first_name' },
  { key: 'leaves', route: 'hr/leaves', label: 'leave', table: 'leaves', col: 'reason' },
  { key: 'attendance', route: 'hr/attendance', label: 'attendance', table: 'attendance', col: 'notes' },
  { key: 'contracts', route: 'hr/contracts', label: 'contract', table: 'employment_contracts', col: 'contract_number' },
  { key: 'users', route: 'settings/users', label: 'user', table: 'users', col: 'name' },
  { key: 'roles', route: 'rbac/roles', label: 'role', table: 'roles', col: 'name' },
  { key: 'drivers', route: 'drivers/new', label: 'driver', table: 'drivers', col: 'first_name' },
  { key: 'tables', route: 'tables', label: 'table', table: 'restaurant_tables', col: 'name' },
  { key: 'categories', route: 'categories/new', label: 'category', table: 'categories', col: 'name' },
  { key: 'products', route: 'products/new', label: 'product', table: 'products', col: 'name_en' },
].filter((t) => !ONLY.length || ONLY.includes(t.key));
const BASE = 'input:visible:not([type=hidden]):not([type=checkbox]):not([type=radio]):not([type=file]):not([type=search]):not([type=submit]):not([placeholder*="Search" i]):not([placeholder*="search" i]), textarea:visible';
let FIELDS = BASE;
const { ctx, page } = await openChrome();
const calls = [];
page.on('response', async (r) => { if (r.url().includes(':18082') && !['GET', 'OPTIONS'].includes(r.request().method())) { let b = ''; try { b = (await r.text()).slice(0, 240); } catch { /* n/a */ } calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').replace(/"password":"[^"]*"/g, '"password":"<redacted>"').slice(0, 260), b }); } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard`);
const results = []; let n = 0;
const AR = /[؀-ۿ]/;
const valueFor = (f, mode, tag) => {
  const short = tag.replace(/[^a-z0-9]/gi, '').slice(-10);
  const key = `${f.name} ${f.label} ${f.ph}`.toLowerCase();
  const isAr = AR.test(f.ph) || /_ar|arabic/.test(key);
  if (f.type === 'email' || /e-?mail|@/.test(key)) return mode === 'invalid' ? 'not-an-email' : `r4.${short}@example.test`;
  if (f.type === 'tel' || /phone|mobile|\+20/.test(key)) return mode === 'invalid' ? 'abc' : '+201000000777';
  if (f.type === 'password' || /password/.test(key)) return mode === 'invalid' ? '1' : 'AuditPass#2026A';
  if (f.type === 'date') return '2026-11-20';
  if (f.type === 'time') return '09:00';
  if (f.type === 'number' || /^\d+(\.\d+)?$/.test(f.ph.trim()) || /price|cost|amount|salary|rate|quantity|stock|days|hours|capacity|limit|balance|debit|credit/.test(key)) return mode === 'invalid' ? '-5' : '12';
  if (/^[a-z]{3}$/i.test(f.ph.trim())) return f.ph.trim().toUpperCase();
  if (/\d{3}-\d{3}/.test(f.ph) || /tax|iban|account.?number|code|sku|national/.test(key)) return mode === 'invalid' ? `'; --` : `R4${short}`.slice(0, 12);
  if (mode === 'arabic') return (isAr ? 'اختبار ' : 'R4 ') + 'عربي ' + short.slice(-4);
  if (mode === 'long') return 'X'.repeat(300);
  if (mode === 'invalid') return `<b>x</b>'; --`;
  return isAr ? `اختبار ${short.slice(-4)}` : `R4UI ${short}`;
};
for (const t of TARGETS) {
  for (const mode of MODES) {
    n++; const tag = `${t.key}${mode}${Date.now().toString().slice(-6)}`; const before = calls.length;
    const rec = { target: t.key, route: t.route, mode, tag };
    let cntBefore = null; try { cntBefore = +sql(`select count(*) from ${t.table}`); } catch { /* table name unknown */ }
    try {
      await page.goto(`${DASH}/en/dashboard/${t.route}`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(5000);
      if (!t.route.endsWith('/new')) { const add = page.locator('button:visible').filter({ hasText: /^(\+\s*)?(add|new|create|record|invite|request)/i }).first(); if (await add.count()) { await add.click({ timeout: 8000 }).catch(() => {}); await page.waitForTimeout(2000); } }
      FIELDS = (await page.locator('[role=dialog]:visible').count()) > 0 ? BASE.split(', ').map((x) => `[role=dialog] ${x}`).join(', ') : BASE;
      const fields = await page.locator(FIELDS).evaluateAll((els) => els.map((e, i) => ({ i, type: e.type || 'textarea', name: e.getAttribute('name') || e.id || '', ph: e.getAttribute('placeholder') || '', label: (e.labels?.[0]?.textContent || e.getAttribute('aria-label') || '').trim().slice(0, 40), req: e.required || e.getAttribute('aria-required') === 'true', maxlength: e.maxLength > 0 ? e.maxLength : null })));
      rec.fieldCount = fields.length; rec.fields = fields.map((f) => `${f.label || f.name || f.ph || f.type}${f.req ? '*' : ''}${f.maxlength ? '(max' + f.maxlength + ')' : ''}`).slice(0, 30);
      const sel = await page.locator('select:visible').count(); rec.selects = sel;
      for (let k = 0; k < sel; k++) { await page.locator('select:visible').nth(k).selectOption({ index: 1 }).catch(() => {}); }
      const loc = page.locator(FIELDS); const filled = [];
      if (mode !== 'empty') for (const f of fields) { const v = valueFor(f, mode, tag); await loc.nth(f.i).fill(v, { timeout: 4000 }).then(() => filled.push(`${f.label || f.name || f.type}=${v.slice(0, 24)}`)).catch(() => {}); }
      rec.filled = filled.slice(0, 16);
      await page.screenshot({ path: `${d.shot}/AUD-R4-FORM-${t.key}-${mode}-filled.png`, timeout: 20000 }).catch(() => {});
      const submit = page.locator('button:visible').filter({ hasText: /^(save|create|add|submit|confirm|invite|record|post)/i }).last();
      rec.submitFound = (await submit.count()) > 0; rec.submitDisabled = rec.submitFound ? await submit.isDisabled().catch(() => null) : null;
      if (rec.submitFound && !rec.submitDisabled) { rec.submitClicked = await submit.click({ timeout: 6000 }).then(() => true).catch(() => false); }
      await page.waitForTimeout(4500);
      rec.validationMessages = [...new Set(await page.locator('[role=alert], [aria-invalid=true] ~ *, p, span, div').evaluateAll((els) => els.filter((e) => e.children.length === 0 && e.innerText && e.innerText.length < 110 && /required|invalid|must|cannot|too long|at least|already|unique|failed|error|مطلوب|غير صالح/i.test(e.innerText)).map((e) => e.innerText.trim())).catch(() => []))].slice(0, 8);
      rec.toast = (await page.locator('[role=status], [data-sonner-toast]').allInnerTexts().catch(() => [])).join(' | ').slice(0, 160);
      rec.dialogOpen = (await page.locator('[role=dialog]:visible').count()) > 0; rec.url = page.url().replace(DASH, '');
      rec.invalidFields = await page.locator('[aria-invalid=true]').count();
      rec.mutations = calls.slice(before).map((c) => ({ m: c.m, u: c.u, s: c.s, req: c.req.slice(0, 200), b: c.b.slice(0, 160) }));
      await page.screenshot({ path: `${d.shot}/AUD-R4-FORM-${t.key}-${mode}-after.png`, timeout: 20000 }).catch(() => {});
      try { rec.dbRowDelta = cntBefore === null ? 'n/a' : +sql(`select count(*) from ${t.table}`) - cntBefore; } catch { rec.dbRowDelta = 'n/a'; }
    } catch (e) { rec.error = String(e.message).slice(0, 160); }
    results.push(rec);
    console.log(`${n} ${t.key}/${mode}: fields=${rec.fieldCount} submit=${rec.submitClicked} inv=${rec.invalidFields} msgs=${(rec.validationMessages || []).length} mut=${(rec.mutations || []).map((c) => c.m + ' ' + c.u + ' ' + c.s).join(',')} dbDelta=${rec.dbRowDelta ?? '-'}${rec.error ? ' ERR ' + rec.error : ''}`);
    fs.writeFileSync(`${d.net}/AUD-R4-FORM-MATRIX${ONLY.length ? '-' + ONLY.join('_') : ''}.json`, JSON.stringify({ run: 'RUN-004', tenant: 4, results }, null, 1));
  }
}
await ctx.close();
