// RUN-005 / PH-06 (Phase E): ISOLATED field-level tests through the real UI in headed Chrome. For every discovered field of a form: fill ALL fields with a valid baseline, change ONLY the target field to a case value,
// submit, and record frontend validation, the real request/response, whether a row was persisted and WHAT was persisted (detects silent trimming/truncation). Tenant 4 diagnostic fixture; synthetic data only.
// Env: ONLY=<form keys>, KINDS=<csv kinds>, MAXF=<fields per form>, CASES=<csv case names>.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
const ONLY = (process.env.ONLY || '').split(',').filter(Boolean); const MAXF = +(process.env.MAXF || 99); const CASESEL = (process.env.CASES || '').split(',').filter(Boolean);
const TARGETS = [
  { key: 'customers', route: 'accounting/customers', add: true, table: 'customers_accounting' },
  { key: 'suppliers', route: 'accounting/suppliers', add: true, table: 'suppliers' },
  { key: 'accounts', route: 'accounting/chart-of-accounts', add: true, table: 'chart_of_accounts' },
  { key: 'drivers', route: 'drivers/new', add: false, table: 'drivers' },
  { key: 'products', route: 'products/new', add: false, table: 'products' },
  { key: 'categories', route: 'categories/new', add: false, table: 'categories' },
  { key: 'roles', route: 'rbac/roles', add: true, table: 'roles' },
  { key: 'users', route: 'settings/users', add: true, table: 'users' },
  { key: 'tables', route: 'tables', add: true, table: 'restaurant_tables' },
  { key: 'contracts', route: 'hr/contracts', add: true, table: 'employment_contracts' },
  { key: 'attendance', route: 'hr/attendance', add: true, table: 'attendance' },
].filter((t) => !ONLY.length || ONLY.includes(t.key));
const BASE = 'input:visible:not([type=hidden]):not([type=checkbox]):not([type=radio]):not([type=file]):not([type=search]):not([type=submit]):not([placeholder*="earch" i]), textarea:visible';
const AR = /[؀-ۿ]/;
const kindOf = (f) => { const key = `${f.name} ${f.label} ${f.ph}`.toLowerCase(); if (f.type === 'email' || /e-?mail|@/.test(key)) return 'email'; if (f.type === 'tel' || /phone|mobile|\+20/.test(key)) return 'phone'; if (f.type === 'password' || /password/.test(key)) return 'password'; if (f.type === 'date') return 'date'; if (f.type === 'datetime-local') return 'datetime'; if (f.type === 'time') return 'time'; if (f.type === 'number' || /^\d+(\.\d+)?$/.test(f.ph.trim()) || /(price|cost|amount|salary|rate|quantity|stock|days|hours|capacity|limit|balance|order)/.test(key)) return 'number'; if (/^[a-z]{3}$/i.test(f.ph.trim())) return 'currency'; if (/\d{3}-\d{3}/.test(f.ph) || /(tax|iban|account.?number|code|sku|national|number)/.test(key)) return 'code'; return f.type === 'textarea' ? 'textarea' : 'text'; };
const baseline = (f, kind, tag) => { const short = tag.slice(-8); switch (kind) { case 'email': return `r5.${short}@example.test`; case 'phone': return '+201000000888'; case 'password': return 'AuditPass#2026A'; case 'date': return '2026-11-20'; case 'time': return '09:00'; case 'number': return '12'; case 'currency': return f.ph.trim().toUpperCase(); case 'code': return `R5${short}`.slice(0, 12); default: return AR.test(f.ph) ? `اختبار ${short.slice(-4)}` : `R5 ${short}`; } };

const FORMOV = {
  tables: [[/table_number/i, () => String(300 + Math.floor(Math.random() * 600))], [/capacity/i, () => '4']],
  categories: [[/display_order/i, () => '']],
  contracts: [[/employee id/i, () => '4'], [/contract type/i, () => 'permanent'], [/contract number/i, (t) => `R5C${t.slice(-8)}`], [/currency/i, () => 'EGP'], [/salary/i, () => '5000']],
  attendance: [[/^clock_in$/i, () => '2026-11-20T09:00'], [/^clock_out$/i, () => '2026-11-20T17:00'], [/scheduled_clock_in/i, () => '09:00'], [/scheduled_clock_out/i, () => '17:00'], [/break/i, () => '30']],
  users: [[/john doe|^name$/i, (t) => `R5 User ${t.slice(-6)}`], [/@/, (t) => `r5.${t.slice(-8)}@example.test`], [/555|phone/i, () => '+201000000888'], [/•|password/i, () => 'AuditPass#2026A']],
};
const base2 = (form, g, kind, tag) => { for (const [re, fn] of FORMOV[form] || []) { if (re.test(`${g.name} ${g.label} ${g.ph}`)) return fn(tag); } return baseline(g, kind, tag); };
const CASES = { text: [['empty', ''], ['whitespace', '   '], ['long300', 'X'.repeat(300)], ['arabic', 'اختبار عربي ٧٧'], ['special', `<b>x</b>'; -- &amp;`]], textarea: [['empty', ''], ['long300', 'X'.repeat(300)], ['arabic', 'اختبار عربي ٧٧'], ['special', `<b>x</b>'; --`]],
  email: [['empty', ''], ['invalid', 'not-an-email'], ['long260', 'a'.repeat(250) + '@example.test'], ['uppercase', 'R5.UPPER@EXAMPLE.TEST']], phone: [['empty', ''], ['letters', 'abc'], ['long30', '+2010000008881234567890123456'], ['arabic-digits', '٠١٠٠٠٠٠٠٨٨٨']],
  number: [['empty', ''], ['negative', '-5'], ['zero', '0'], ['decimal6', '12.345678'], ['huge', '99999999999'], ['tiny', '0.001']], date: [['empty', ''], ['past1900', '1900-01-01'], ['future2999', '2999-12-31']], code: [['empty', ''], ['long40', 'C'.repeat(40)], ['special', `'; --`], ['arabic', 'رمز٧٧']],
  datetime: [['empty', ''], ['before-clock-in', '2026-11-20T08:00'], ['future2999', '2999-12-31T09:00']], currency: [['empty', ''], ['lowercase', 'egp'], ['invalid', 'ZZZZ']], password: [['empty', ''], ['short1', '1']], time: [['empty', '']] };
const { ctx, page } = await openChrome();
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && !['GET', 'OPTIONS'].includes(r.request().method())) { let b = ''; try { b = (await r.text()).slice(0, 200); } catch { /* n/a */ } calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').replace(/"password":"[^"]*"/g, '"password":"<redacted>"').slice(0, 300), b }); } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard`);
const open = async (t) => { await page.goto(`${DASH}/en/dashboard/${t.route}`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(2800); if (t.add) { const add = page.locator('button:visible').filter({ hasText: /^(\+\s*)?(add|new|create|record|invite|request)/i }).first(); if (await add.count()) { await add.click({ timeout: 8000 }).catch(() => {}); await page.waitForTimeout(1200); } } const scope = (await page.locator('[role=dialog]:visible').count()) > 0 ? BASE.split(', ').map((x) => `[role=dialog] ${x}`).join(', ') : BASE; const fields = await page.locator(scope).evaluateAll((els) => els.map((e, i) => ({ i, type: e.type || 'textarea', name: e.getAttribute('name') || e.id || '', ph: e.getAttribute('placeholder') || '', label: (e.labels?.[0]?.textContent || e.getAttribute('aria-label') || '').trim().slice(0, 40), maxlength: e.maxLength > 0 ? e.maxLength : null }))); return { scope, fields }; };
for (const t of TARGETS) {
  const results = []; let probe; try { probe = await open(t); } catch (e) { console.log(t.key, 'open failed', String(e.message).slice(0, 80)); continue; }
  const fields = probe.fields.slice(0, MAXF); console.log(t.key, 'fields', fields.length);
  for (const f of fields) {
    const kind = kindOf(f); const cases = (CASES[kind] || CASES.text).filter((c) => !CASESEL.length || CASESEL.includes(c[0]));
    for (const [cname, cval] of cases) {
      const tag = `r5${Date.now()}`; const rec = { form: t.key, field_index: f.i, field: f.label || f.name || f.ph || f.type, placeholder: f.ph, input_type: f.type, maxlength_attr: f.maxlength, kind, case: cname, value_sent: cval.slice(0, 40) + (cval.length > 40 ? '...(' + cval.length + ')' : '') };
      try {
        const pr = await open(t); const loc = page.locator(pr.scope); const sel = await page.locator('select:visible').count(); for (let k = 0; k < sel; k++) await page.locator('select:visible').nth(k).selectOption({ index: 1 }).catch(() => {});
        for (const g of pr.fields) { const v = g.i === f.i ? cval : base2(t.key, g, kindOf(g), tag); await loc.nth(g.i).fill(v, { timeout: 4000 }).catch(() => {}); }
        const before = calls.length; let cnt0 = null; try { cnt0 = +sql(`select count(*) from ${t.table}`); } catch { /* n/a */ }
        const submit = page.locator('button:visible').filter({ hasText: /^(save|create|add|submit|confirm|invite|record|post|request)/i }).last();
        rec.submit_found = (await submit.count()) > 0; rec.submit_disabled = rec.submit_found ? await submit.isDisabled().catch(() => null) : null;
        if (rec.submit_found && !rec.submit_disabled) await submit.click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(2300);
        rec.html5_invalid = await page.locator('input:invalid, textarea:invalid').count().catch(() => 0); rec.html5_message = await page.evaluate(() => { const e = document.querySelector('input:invalid, textarea:invalid'); return e ? ((e.labels?.[0]?.textContent || e.placeholder || e.name || '').trim().slice(0, 30) + ': ' + e.validationMessage.slice(0, 80)) : ''; }).catch(() => '');
        rec.messages = [...new Set(await page.locator('[role=alert], p, span, div').evaluateAll((els) => els.filter((e) => e.children.length === 0 && e.innerText && e.innerText.length < 110 && /required|invalid|must|cannot|too long|at least|already|unique|failed|error|status code|مطلوب|غير صالح/i.test(e.innerText)).map((e) => e.innerText.trim())).catch(() => []))].slice(0, 4);
        const mut = calls.slice(before); rec.request = mut.map((c) => `${c.m} ${c.u} ${c.s}`).slice(0, 3); rec.response = (mut[0]?.b || '').slice(0, 150); rec.request_body = (mut[0]?.req || '').slice(0, 200);
        let cnt1 = null; if (mut.some((c) => c.s >= 200 && c.s < 300)) { try { cnt1 = +sql(`select count(*) from ${t.table}`); } catch { /* n/a */ } } else cnt1 = cnt0;
        rec.rows_created = cnt0 === null || cnt1 === null ? 'n/a' : cnt1 - cnt0;
        if (rec.rows_created > 0) { try { rec.persisted = sql(`select row_to_json(x) from (select * from ${t.table} order by id desc limit 1) x`).slice(0, 420); } catch { rec.persisted = 'n/a'; } }
        rec.outcome = rec.rows_created > 0 ? 'ACCEPTED-PERSISTED' : mut.length ? `REJECTED-BY-SERVER(${mut[0].s})` : (rec.html5_invalid > 0 || rec.messages.length ? 'REJECTED-BY-CLIENT' : 'NO-REQUEST-NO-MESSAGE');
      } catch (e) { rec.outcome = 'TOOLING-ERROR'; rec.error = String(e.message).split('\n')[0].slice(0, 120); }
      results.push(rec);
    }
    console.log(`  ${t.key}/${f.label || f.ph || f.i} (${kind}): ${results.filter((r) => r.field_index === f.i).map((r) => r.case + '=' + r.outcome.replace('REJECTED-BY-', 'rej-').replace('ACCEPTED-PERSISTED', 'ok')).join(' ')}`);
    fs.writeFileSync(`${d.net}/AUD-R5-FIELD-${t.key}.json`, JSON.stringify({ run: 'RUN-005', form: t.key, tenant: 4, results }, null, 1));
  }
}
await ctx.close();
