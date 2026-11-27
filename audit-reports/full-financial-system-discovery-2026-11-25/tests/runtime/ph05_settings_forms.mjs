// RUN-003 / PH-05 / AUD-SET-001..: settings form behaviour in HEADED Chrome, tenant 4 owner, isolated stack. Mutates tenant 4 contact info + KDS settings only.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, recordMutations, sql, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-05');
const { ctx, page } = await openChrome();
const calls = []; recordMutations(page, calls);
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/settings/restaurant`);
const log = []; const rec = (id, o) => { log.push({ id, ...o }); console.log(id, JSON.stringify(o)); };
const shot = (n) => page.screenshot({ path: `${d.shot}/${n}.png` });
const mark = () => calls.length;
const since = (m) => calls.slice(m).map((c) => `${c.method} ${c.url} ${c.status} ${c.response.slice(0, 160)}`);
const tenantRow = () => sql("select email,phone from tenants where id=4");

// ---- Restaurant contact settings
await page.waitForSelector('button:has-text("Save Changes")', { timeout: 120000 });
const email = page.locator('input[type=email]'), phone = page.locator('input[type=tel]');
rec('SET-REST-before', { db: tenantRow(), email: await email.inputValue(), phone: await phone.inputValue() });
for (const [id, e, p] of [['AUD-SET-001-invalid-email', 'not-an-email', '+201000000444'], ['AUD-SET-002-invalid-phone', 'ok4@example.test', 'abc-###'], ['AUD-SET-003-empty', '', ''], ['AUD-SET-004-valid', 'audit.c.contact@example.test', '+201000000444']]) {
  await email.fill(e); await phone.fill(p); await shot(`${id}-before`);
  const m = mark(); await page.click('button:has-text("Save Changes")'); await page.waitForTimeout(3500); await shot(`${id}-after`);
  rec(id, { input: { e, p }, calls: since(m), dbAfter: tenantRow(), nativeInvalid: await email.evaluate((x) => !x.validity.valid) });
}
await page.reload({ waitUntil: 'domcontentloaded' }); await page.waitForSelector('input[type=email]', { timeout: 60000 }); await page.waitForTimeout(2500);
rec('AUD-SET-005-persist-reload', { emailShown: await page.locator('input[type=email]').inputValue(), phoneShown: await page.locator('input[type=tel]').inputValue(), db: tenantRow() });

// ---- KDS settings
await page.goto(`${DASH}/en/dashboard/kds/settings`, { waitUntil: 'domcontentloaded', timeout: 120000 });
await page.waitForSelector('button:has-text("Save settings")', { timeout: 120000 });
const nums = page.locator('input[type=number]');
const kds = () => sql("select * from kds_settings where tenant_id=4").slice(0, 300);
let kdsTable = true; try { kds(); } catch { kdsTable = false; }
for (const [id, w, dg] of [['AUD-SET-006-kds-inverted', '30', '5'], ['AUD-SET-007-kds-negative', '-5', '-10'], ['AUD-SET-008-kds-huge', '999999', '9999999'], ['AUD-SET-009-kds-valid', '12', '25']]) {
  await nums.nth(0).fill(w); await nums.nth(1).fill(dg); await shot(`${id}-before`);
  await page.waitForTimeout(800);
  const saveBtn = page.locator('button:has-text("Save settings")'); const enabled = await saveBtn.isEnabled();
  const m = mark(); if (enabled) { await saveBtn.click(); await page.waitForTimeout(3000); } await shot(`${id}-after`);
  const txt = await page.locator('body').innerText();
  rec(id, { input: { warn: w, danger: dg }, saveEnabled: enabled, calls: since(m), validationHint: txt.split('\n').find((x) => /must|invalid|greater|less than|error/i.test(x)) ?? null });
}
await page.reload({ waitUntil: 'domcontentloaded' }); await page.waitForSelector('input[type=number]', { timeout: 60000 }); await page.waitForTimeout(2500);
rec('AUD-SET-010-kds-persist', { warn: await nums.nth(0).inputValue(), danger: await nums.nth(1).inputValue(), kdsTable });

// ---- POS settings persistence (server or browser?)
await page.goto(`${DASH}/en/dashboard/pos/settings`, { waitUntil: 'domcontentloaded', timeout: 120000 });
await page.waitForSelector('input[type=checkbox]', { timeout: 120000 }); await page.waitForTimeout(2000);
const m = mark(); const box = page.locator('input[type=checkbox]').nth(3);
const was = await box.isChecked(); await box.click({ force: true }); await page.click('button:has-text("Save")'); await page.waitForTimeout(2500);
rec('AUD-SET-011-pos-autoprint-toggle', { wasChecked: was, mutatingApiCalls: since(m) });
fs.writeFileSync(`${d.net}/AUD-SET-001-011.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, role: 'owner', log, allMutations: calls }, null, 1));
await ctx.close();
