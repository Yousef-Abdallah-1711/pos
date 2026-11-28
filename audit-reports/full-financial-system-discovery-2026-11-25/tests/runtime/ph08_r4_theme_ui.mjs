// RUN-004 / PH-08: Theme Builder through the REAL UI (headed Chrome, tenant 4 fixture): export, copy, edit, activate, storefront effect, restore. Records every theme API call.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08');
const { ctx, page } = await openChrome(); await page.context().grantPermissions([]).catch(() => {});
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && /theme/i.test(r.url()) && r.request().method() !== 'GET') { let b = ''; try { b = (await r.text()).slice(0, 200); } catch {} calls.push({ m: r.request().method(), u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').slice(0, 200), b }); } });
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 360)); };
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/theme-builder`); await page.waitForTimeout(8000);
const themesBefore = sql("select id||':'||name||':'||is_active from themes_v2 where tenant_id=4 order by id").replace(/\n/g, ' ; ');
rec('TH-00', 'themes before', { themesBefore });
// export (two-step dialog: Export -> Download)
await page.locator('button:has-text("Export")').first().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(1500);
const dl = page.waitForEvent('download', { timeout: 8000 }).catch(() => null);
await page.locator('[role=dialog] button:has-text("Download"), button:has-text("Download")').first().click({ timeout: 6000 }).catch(() => {}); const dwn = await dl; rec('TH-01', 'Export dialog -> Download', { download: dwn ? dwn.suggestedFilename() : 'no download' });
if (dwn) { const p = `${d.net}/AUD-R4-THEME-export.json`; await dwn.saveAs(p).catch(() => {}); try { rec('TH-01b', 'export contents', { bytes: fs.statSync(p).size, head: fs.readFileSync(p, 'utf8').slice(0, 200) }); } catch {} }
await page.keyboard.press('Escape'); await page.waitForTimeout(800); rec('TH-01c', 'Escape closes the export dialog?', { stillOpen: (await page.locator('[role=dialog]:visible').count()) > 0 }); await page.locator('[role=dialog] button:has-text("Close")').first().click({ timeout: 4000 }).catch(() => {}); await page.waitForTimeout(800);
// copy (fresh page so no stale dialog intercepts the click)
await page.reload({ waitUntil: 'domcontentloaded' }); await page.waitForTimeout(6000);
let before = calls.length; await page.locator('main button:has-text("Copy")').first().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(5000);
await page.screenshot({ path: `${d.shot}/AUD-R4-THEME-after-copy.png`, fullPage: true }).catch(() => {});
rec('TH-02', 'Copy the default theme', { calls: calls.slice(before).map((c) => `${c.m} ${c.u} ${c.s}`), themes: sql("select id||':'||name||':'||is_active from themes_v2 where tenant_id=4 order by id").replace(/\n/g, ' ; ') });
const newId = sql("select id from themes_v2 where tenant_id=4 order by id desc limit 1");
// open editor for the newest copy (last Edit button)
await page.reload({ waitUntil: 'domcontentloaded' }); await page.waitForTimeout(6000);
await page.locator('button:has-text("Edit")').last().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(8000);
await page.screenshot({ path: `${d.shot}/AUD-R4-THEME-editor.png`, fullPage: true }).catch(() => {});
const editor = await page.evaluate(() => ({ url: location.pathname, buttons: [...document.querySelectorAll('button')].map((b) => (b.innerText || b.getAttribute('aria-label') || '').trim()).filter(Boolean).slice(0, 40), colorInputs: document.querySelectorAll('input[type=color]').length, textInputs: [...document.querySelectorAll('input[type=text]')].map((i) => i.name || i.placeholder || i.value).slice(0, 20), tabs: [...document.querySelectorAll('[role=tab]')].map((t) => t.innerText.trim()) }));
rec('TH-03', 'Editor opened', editor);
fs.writeFileSync(`${d.net}/AUD-R4-THEME-UI.json`, JSON.stringify({ run: 'RUN-004', tenant: 4, themeId: newId, log, calls }, null, 1));
await ctx.close();
