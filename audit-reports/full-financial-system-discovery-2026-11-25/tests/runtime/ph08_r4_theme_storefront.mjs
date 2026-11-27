// RUN-004 / PH-08: activate theme 13 (primary #ff0000 in config) through the UI, load the real storefront in Chrome, read the applied primary colour; restore Default Theme.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08');
const { ctx, page } = await openChrome();
const log = []; const rec = (k, desc, extra = {}) => { log.push({ k, desc, ...extra }); console.log(k, desc, JSON.stringify(extra).slice(0, 420)); };
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/theme-builder`); await page.waitForTimeout(7000);
await page.locator('xpath=//h3[normalize-space()="Default Theme (Copy)"]/ancestor::div[.//button[contains(.,"Activate Theme")]][1]//button[contains(.,"Activate Theme")]').first().click({ timeout: 6000 }).catch((e) => rec('A0', 'activate click failed', { e: String(e.message).slice(0, 80) })); await page.waitForTimeout(5000);
rec('A1', 'active themes', { active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
const sp = await ctx.newPage();
await sp.goto('http://audit-c-diag.localhost:13113/en', { waitUntil: 'domcontentloaded', timeout: 120000 }); await sp.waitForTimeout(10000);
const info = await sp.evaluate(() => { const cs = getComputedStyle(document.documentElement); const out = {}; for (const v of ['--primary', '--theme-primary', '--color-primary', '--primary-color', '--brand']) out[v] = cs.getPropertyValue(v).trim(); const sheets = [...document.styleSheets].length; const html = document.documentElement.outerHTML; return { vars: out, sheets, mentionsFF0000: /ff0000|255,\s*0,\s*0|0 100% 50%/i.test(html), title: document.title, len: document.body.innerText.length }; });
rec('A2', 'storefront applied theme', info);
await sp.screenshot({ path: `${d.shot}/AUD-R4-THEME-storefront-theme13.png` }).catch(() => {}); await sp.close();
await page.goto(`${DASH}/en/dashboard/theme-builder`, { waitUntil: 'domcontentloaded' }); await page.waitForTimeout(6000);
await page.locator('xpath=//h3[normalize-space()="Default Theme"]/ancestor::div[.//button[contains(.,"Activate Theme")]][1]//button[contains(.,"Activate Theme")]').first().click({ timeout: 6000 }).catch(() => {}); await page.waitForTimeout(4000);
rec('R1', 'restored', { active: sql("select string_agg(id||':'||is_active,' ') from themes_v2 where tenant_id=4") });
fs.writeFileSync(`${d.net}/AUD-R4-THEME-STOREFRONT.json`, JSON.stringify({ run: 'RUN-004', log }, null, 1));
await ctx.close();
