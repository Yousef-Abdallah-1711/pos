// RUN-003 / PH-06: category creation through the real UI form (HEADED Chrome): name only, then name + display order. Records the mutating request body and status.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, DASH, PASSWORD } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
const browser = await chromium.launch({ channel: 'chrome', headless: false });
const page = await (await browser.newContext({ viewport: { width: 1440, height: 900 } })).newPage();
const calls = []; page.on('response', async (r) => { if (r.url().includes(':18082') && r.request().method() === 'POST' && r.url().includes('/categories')) { let b = ''; try { b = (await r.text()).slice(0, 200); } catch {} calls.push({ u: r.url().replace('http://localhost:18082/api/v1', ''), s: r.status(), req: (r.request().postData() || '').slice(0, 300), b }); } });
await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded', timeout: 180000 });
for (let i = 0; i < 40; i++) { await page.waitForTimeout(3000); if ((await page.locator('input[type=email]:not([disabled])').count()) > 0) break; }
await page.fill('input[type=email]', 'audit.c.owner@example.test'); await page.fill('input[type=password]', PASSWORD); await page.click('button[type=submit]'); await page.waitForTimeout(8000);
const out = [];
for (const variant of ['name-only', 'name+display-order']) {
  await page.goto(`${DASH}/en/dashboard/categories/new`, { waitUntil: 'domcontentloaded', timeout: 180000 }); await page.waitForTimeout(6000);
  await page.fill('[name=name]', 'UI Cat ' + variant); 
  if (variant === 'name+display-order') await page.fill('[name=display_order]', '7');
  const before = calls.length;
  await page.locator('button:visible').filter({ hasText: /^(save|create|add|submit)/i }).last().click({ timeout: 10000 }).catch((e) => console.log('click', String(e).slice(0, 80)));
  await page.waitForTimeout(5000);
  await page.screenshot({ path: `${d.shot}/AUD-UIFORM-category-${variant}.png`, timeout: 30000 }).catch(() => {});
  out.push({ variant, calls: calls.slice(before), displayOrderDefault: await page.locator('[name=display_order]').inputValue().catch(() => 'n/a') });
  console.log(variant, JSON.stringify(calls.slice(before)));
}
fs.writeFileSync(`${d.net}/AUD-UIFORM-CATEGORY.json`, JSON.stringify({ run: 'RUN-003', out }, null, 1));
await browser.close();
