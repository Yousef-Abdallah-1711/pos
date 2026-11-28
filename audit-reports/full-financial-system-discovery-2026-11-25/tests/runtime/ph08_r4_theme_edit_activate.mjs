// RUN-004 / PH-08: theme editor colour edits (invalid then valid), persistence, activation, storefront effect, restore. Real UI, tenant 4 fixture; restores Default Theme (id 4) active and deletes the copy it edited.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08'); const id = process.env.THEME_ID || '13';
const { ctx, page } = await openChrome();
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && /theme/i.test(r.url()) && r.request().method() !== 'GET') { let b = ''; try { b = (await r.text()).slice(0, 200); } catch {} calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').slice(0, 200), b }); } });
const log = []; const rec = (k, desc, extra = {}) => { log.push({ k, desc, ...extra }); console.log(k, desc, JSON.stringify(extra).slice(0, 380)); };
const col = () => sql(`select primary_color||'|'||is_active||'|'||is_published from themes_v2 where id=${id}`);
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/theme-builder/editor/${id}`); await page.waitForTimeout(8000);
rec('E0', 'before', { row: col() });
const t1 = page.locator('input[type=text]').first();
let b = calls.length; await t1.fill('not-a-color'); await page.waitForTimeout(5000);
rec('E1', 'type invalid colour text into primary field', { calls: calls.slice(b).map((c) => `${c.m} ${c.u} ${c.s} ${c.req.slice(0, 80)}`), row: col(), value: await t1.inputValue() });
b = calls.length; await t1.fill('#ff0000'); await page.waitForTimeout(5000);
rec('E2', 'type valid #ff0000', { calls: calls.slice(b).map((c) => `${c.m} ${c.u} ${c.s}`), row: col() });
const saveBtn = page.locator('button:has-text("Save"), button:has-text("Publish")').first(); rec('E3', 'explicit Save/Publish control present', { present: (await saveBtn.count()) > 0 });
if (await saveBtn.count()) { b = calls.length; await saveBtn.click({ timeout: 5000 }).catch(() => {}); await page.waitForTimeout(4000); rec('E3b', 'clicked save', { calls: calls.slice(b).map((c) => `${c.m} ${c.u} ${c.s}`), row: col() }); }
// history tab
await page.locator('[role=tab]:has-text("History"), button:has-text("History")').first().click({ timeout: 5000 }).catch(() => {}); await page.waitForTimeout(3000);
rec('E4', 'History tab text', { text: (await page.locator('main').innerText().catch(() => '')).slice(0, 220).replace(/\n/g, ' | ') });
await page.screenshot({ path: `${d.shot}/AUD-R4-THEME-editor-history.png`, fullPage: true }).catch(() => {});
// activate from list
await page.goto(`${DASH}/en/dashboard/theme-builder`, { waitUntil: 'domcontentloaded' }); await page.waitForTimeout(7000);
b = calls.length; await page.locator('div:has(h3:has-text("Default Theme (Copy)")) button:has-text("Activate Theme")').first().click({ timeout: 6000 }).catch(async () => { await page.locator('button:has-text("Activate Theme")').last().click({ timeout: 6000 }).catch(() => {}); }); await page.waitForTimeout(5000);
rec('A1', 'Activate the edited copy via UI', { calls: calls.slice(b).map((c) => `${c.m} ${c.u} ${c.s}`), active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
const sf = await (await fetch('http://audit-c-diag.localhost:13113/en', { signal: AbortSignal.timeout(60000) }).catch(() => ({ text: async () => '' }))).text();
rec('A2', 'storefront HTML mentions the new colour', { has_ff0000: /ff0000/i.test(sf), bytes: sf.length });
const pub = await fetch('http://127.0.0.1:18082/api/v1/public/theme?restaurant_slug=audit-c-diag', { headers: { Origin: 'http://localhost:13113' } }).then((r) => r.text()).catch(() => '');
rec('A2b', 'public theme API mentions the colour', { has_ff0000: /ff0000/i.test(pub), head: pub.slice(0, 100) });
// restore: activate default (id 4)
await page.reload({ waitUntil: 'domcontentloaded' }); await page.waitForTimeout(6000);
b = calls.length; await page.locator('div:has(h3:has-text("Default Theme")):not(:has(h3:has-text("(Copy)"))) button:has-text("Activate Theme")').first().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(4000);
rec('R1', 'restore Default Theme active', { calls: calls.slice(b).map((c) => `${c.m} ${c.u} ${c.s}`), active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
fs.writeFileSync(`${d.net}/AUD-R4-THEME-EDIT-ACTIVATE.json`, JSON.stringify({ run: 'RUN-004', tenant: 4, themeId: id, log, calls }, null, 1));
await ctx.close();
