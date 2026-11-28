// RUN-005 / PH-07: F-023 (POS currency symbol) and F-004 context in headed Chrome: open POS, add a product, read every currency-looking token and compare with the dashboard (EGP) and the stored currency. No order is submitted.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-07');
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/pos`); await page.waitForTimeout(9000);
const card = page.locator('main [class*=card], main button').filter({ hasText: /Audit|MD|Probe/ }).first();
await card.click({ timeout: 8000 }).catch(() => {}); await page.waitForTimeout(2500);
const text = await page.locator('body').innerText();
const tokens = [...new Set((text.match(/(?:[$€£]|EGP|USD|E£|ج\.م)\s?[\d,.]+|[\d,.]+\s?(?:EGP|USD)/g) || []))].slice(0, 20);
await page.screenshot({ path: `${d.shot}/AUD-R5-POS-currency.png` }).catch(() => {});
await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded' }); await page.waitForTimeout(8000);
const dashTokens = [...new Set(((await page.locator('body').innerText()).match(/(?:[$€£]|EGP|USD)\s?[\d,.]+/g) || []))].slice(0, 6);
const out = { pos_tokens: tokens, dashboard_tokens: dashTokens, tenant_currency_column: (() => { try { return sql("select string_agg(distinct currency,',') from suppliers where tenant_id=4"); } catch { return 'n/a'; } })() };
console.log(JSON.stringify(out)); fs.writeFileSync(`${d.net}/AUD-R5-POS-CURRENCY.json`, JSON.stringify({ run: 'RUN-005', finding: 'F-023', ...out }, null, 1));
await ctx.close();
