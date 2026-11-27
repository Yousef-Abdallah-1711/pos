// RUN-003 / PH-07 / AUD-POS-002: read the POS cart's displayed subtotal/tax/total/currency in HEADED Chrome with the perf overlay hidden via browser-side style (no product change). No order is placed.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-07');
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/pos`);
await page.waitForSelector('text=Audit Burger', { timeout: 120000 }); await page.waitForTimeout(2000);
await page.locator('text=Audit Burger').first().click(); await page.waitForTimeout(500);
await page.locator('text=Audit Burger').first().click(); await page.waitForTimeout(800);
const hidden = await page.evaluate(() => { let n = 0; document.querySelectorAll('*').forEach((el) => { if (el.children.length < 8 && /POS Perf/i.test(el.textContent || '') && getComputedStyle(el).position === 'fixed') { el.style.display = 'none'; n++; } }); return n; });
await page.waitForTimeout(500);
const lines = (await page.locator('body').innerText()).split('\n').map((l) => l.trim()).filter(Boolean);
const idx = lines.findIndex((l) => /^subtotal/i.test(l));
const ui = lines.filter((l, i) => /subtotal|tax|total|vat|\$|EGP|E£|ج\.م/i.test(l) || (idx >= 0 && Math.abs(i - idx) <= 6)).slice(0, 40);
await page.screenshot({ path: `${d.shot}/AUD-POS-002-cart-totals-overlay-hidden.png`, timeout: 40000 }).catch((e) => console.log('shot fail', String(e).slice(0, 80)));
// expected for 2 x 10,000: subtotal 20,000; server 10% -> 2,000 / 22,000; UI label 14% -> 2,800 / 22,800
fs.writeFileSync(`${d.net}/AUD-POS-002.json`, JSON.stringify({ run: 'RUN-003', test: 'AUD-POS-002', overlayElementsHidden: hidden, ui }, null, 1));
console.log({ hidden, ui });
await page.evaluate(() => localStorage.removeItem('pos-storage')); // leave cart empty
await ctx.close();
