// RUN-003 / PH-08: storefront (F-022 revalidation) in HEADED Chrome on the tenant subdomain. Read-only.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08');
const browser = await chromium.launch({ channel: 'chrome', headless: false });
const page = await (await browser.newContext({ viewport: { width: 1280, height: 900 } })).newPage();
const calls = []; page.on('response', (r) => { if (r.url().includes(':18082')) calls.push(`${r.request().method()} ${r.url().replace('http://localhost:18082/api/v1', '')} ${r.status()}`); });
const out = {};
for (const [tag, host] of [['C', 'audit-c-diag'], ['E', 'audit-e-outbox']]) {
  await page.goto(`http://${host}.localhost:13113/en`, { waitUntil: 'domcontentloaded', timeout: 120000 }).catch((e) => { out[tag + '-goto'] = String(e).slice(0, 100); });
  await page.waitForTimeout(9000);
  const text = (await page.locator('body').innerText()).slice(0, 1500);
  await page.screenshot({ path: `${d.shot}/AUD-STORE-${tag}-home.png`, fullPage: true, timeout: 40000 }).catch(() => {});
  out[tag] = { url: page.url(), text: text.replace(/\n+/g, ' | ').slice(0, 700), products: sql(`select string_agg(name_en||'(cat='||coalesce(category_id::text,'NULL')||',avail='||is_available||')',', ') from products where tenant_id=${tag === 'C' ? 4 : 8}`) };
}
out.apiCalls = calls.slice(0, 40);
fs.writeFileSync(`${d.net}/AUD-STORE-001.json`, JSON.stringify({ run: 'RUN-003', out }, null, 1));
console.log(JSON.stringify(out, null, 1).slice(0, 3000));
await browser.close();
