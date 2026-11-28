// RUN-003 / PH-05 / OBS-ONB-001: human-paced probes of the registration subdomain-availability race. HEADED Chrome, fresh context.
// Every scenario submits BEFORE availability is confirmed, so no tenant is created (POST count asserted = 0).
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import fs from 'fs';
import { evidenceDirs, sql, PASSWORD } from './chrome_lib.mjs';
const d = evidenceDirs('PH-05');
const results = [];
const browser = await chromium.launch({ channel: 'chrome', headless: false });
async function scenario(id, { delayAfterTypeMs, checkDelayMs = 0, slug, then }) {
  const ctx = await browser.newContext({ viewport: { width: 1440, height: 1000 } });
  const page = await ctx.newPage(); const posts = [];
  page.on('request', (r) => { if (r.method() === 'POST' && r.url().includes(':18082')) posts.push(r.url()); });
  if (checkDelayMs) await page.route('**/auth/check-subdomain*', async (route) => { await new Promise((r) => setTimeout(r, checkDelayMs)); await route.continue(); });
  await page.goto('http://localhost:13111/en/auth/register', { waitUntil: 'domcontentloaded', timeout: 120000 });
  await page.waitForSelector('#name', { timeout: 120000 });
  await page.fill('#name', 'Probe User'); await page.fill('#email', `probe.${id.toLowerCase()}@example.test`); await page.fill('#restaurantName', 'Probe ' + id);
  await page.fill('#password', PASSWORD); await page.fill('#confirmPassword', PASSWORD);
  await page.locator('#restaurantSlug').pressSequentially(slug, { delay: 90 }); // human typing speed
  await page.waitForTimeout(delayAfterTypeMs);
  if (then) await then(page);
  await page.click('button[type=submit]:has-text("Create account")');
  const snaps = [];
  for (const t of [300, 3000, 8000]) {
    await page.waitForTimeout(t === 300 ? 300 : t - (t === 3000 ? 300 : 3000));
    const txt = await page.locator('body').innerText();
    snaps.push({ afterMs: t, availableShown: /Subdomain is available/.test(txt), staleRequiredErrorShown: /Please choose an available subdomain/.test(txt) });
  }
  await page.screenshot({ path: `${d.shot}/OBS-ONB-001-${id}.png`, fullPage: true });
  results.push({ id, delayAfterTypeMs, checkDelayMs, slug, postsSent: posts.length, snaps }); console.log(JSON.stringify(results.at(-1)));
  await ctx.close();
}
await scenario('S1-click-200ms-after-typing', { delayAfterTypeMs: 200, slug: 'obs-s1-probe' });
await scenario('S2-click-450ms-after-typing', { delayAfterTypeMs: 450, slug: 'obs-s2-probe' });
await scenario('S3-click-1500ms-wait-for-check', { delayAfterTypeMs: 1500, slug: 'obs-s3-probe', then: async () => {} });
await scenario('S4-slow-network-3s-check', { delayAfterTypeMs: 1200, checkDelayMs: 3000, slug: 'obs-s4-probe' });
await scenario('S5-change-slug-after-validation', { delayAfterTypeMs: 1500, slug: 'obs-s5-probe', then: async (p) => { await p.locator('#restaurantSlug').pressSequentially('x', { delay: 90 }); } });
fs.writeFileSync(`${d.net}/OBS-ONB-001-probes.json`, JSON.stringify({ run: 'RUN-003', results, tenantsAfter: sql('select count(*) from tenants') }, null, 1));
console.log('tenants now', sql('select count(*) from tenants'));
await browser.close();
