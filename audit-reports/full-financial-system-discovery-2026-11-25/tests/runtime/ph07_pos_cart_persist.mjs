// RUN-003 / PH-07 / OBS-POS-001: does the POS cart survive after a successful order + browser restart? Read-only (no clicks on products, no order placed).
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-07');
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/pos`);
await page.waitForSelector('text=Audit Burger', { timeout: 120000 }); await page.waitForTimeout(3000);
const body = await page.locator('body').innerText();
const ls = await page.evaluate(() => Object.fromEntries(Object.entries(localStorage).filter(([k]) => /cart|pos/i.test(k)).map(([k, v]) => [k, v.slice(0, 400)])));
const cartLines = body.split('\n').filter((l) => /Audit Burger|Subtotal|Tax|Total|x\s*\d|^\d+$/.test(l)).slice(0, 15);
try { await page.screenshot({ path: `${d.shot}/OBS-POS-001-cart-on-open.png`, timeout: 20000 }); } catch {}
fs.writeFileSync(`${d.net}/OBS-POS-001.json`, JSON.stringify({ run: 'RUN-003', cartLines, localStorageKeys: Object.keys(ls), localStorage: ls }, null, 1));
console.log(JSON.stringify({ cartLines, ls }, null, 1));
await ctx.close();
