// RUN-003 / PH-17 / AUD-XSS-001: stored script text (order 71 customer_name) rendered by orders page and KDS in HEADED Chrome (separate non-persistent context). Read-only.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, DASH, PASSWORD } from './chrome_lib.mjs';
const d = evidenceDirs('PH-17');
const browser = await chromium.launch({ channel: 'chrome', headless: false });
const page = await (await browser.newContext({ viewport: { width: 1440, height: 900 } })).newPage();
const dialogs = []; page.on('dialog', async (dlg) => { dialogs.push(dlg.message()); await dlg.dismiss(); });
await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded', timeout: 180000 });
for (let i = 0; i < 40; i++) { await page.waitForTimeout(3000); if ((await page.locator('input[type=email]').count()) > 0) break; }
await page.fill('input[type=email]', 'audit.c.owner@example.test'); await page.fill('input[type=password]', PASSWORD); await page.click('button[type=submit]'); await page.waitForTimeout(8000);
const out = {};
for (const [tag, path] of [['orders', '/en/dashboard/orders'], ['kds', '/en/dashboard/kds']]) {
  await page.goto(DASH + path, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(9000);
  const found = await page.evaluate(() => { const t = document.body.innerText; const h = document.body.innerHTML; return { textShowsRawTag: t.includes('<script>alert(1)</script>'), htmlHasLiveScriptNode: !!document.querySelector('script') && Array.from(document.querySelectorAll('script')).some((s) => (s.textContent || '').includes('alert(1)')), htmlEscaped: h.includes('&lt;script&gt;alert(1)') }; });
  await page.screenshot({ path: `${d.shot}/AUD-XSS-001-${tag}.png`, timeout: 30000 }).catch(() => {});
  out[tag] = found;
}
out.dialogsFired = dialogs;
fs.writeFileSync(`${d.net}/AUD-XSS-001.json`, JSON.stringify({ run: 'RUN-003', order: 71, out }, null, 1));
console.log(JSON.stringify(out));
await browser.close();
