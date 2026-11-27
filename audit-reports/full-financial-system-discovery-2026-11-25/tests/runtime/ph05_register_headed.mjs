// RUN-003 / PH-05 / AUD-ONB-001: revalidate F-001 (fresh tenant has no accounting bootstrap) in HEADED Chrome.
// Registers a NEW tenant (Audit D) through the real UI, with NO seeding, then measures bootstrap state by read-only SQL.
// Target: isolated stack only. Side effect: creates tenant/user/restaurant rows via the product's own registration flow.
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import fs from 'fs';
import { evidenceDirs, sql, PASSWORD } from './chrome_lib.mjs';

const d = evidenceDirs('PH-05');
const EMAIL = 'audit.d.owner@example.test';
const SLUG = 'audit-d-onboarding';

const measure = (tenantId) => (tenantId
  ? sql(`select 'coa',count(*) from chart_of_accounts where tenant_id=${tenantId} union all select 'auto_posting_rules',count(*) from auto_posting_rules where tenant_id=${tenantId} union all select 'fiscal_periods',count(*) from fiscal_periods where tenant_id=${tenantId} union all select 'tax_periods',count(*) from tax_periods where tenant_id=${tenantId} union all select 'users',count(*) from users where tenant_id=${tenantId} union all select 'restaurants',count(*) from restaurants where tenant_id=${tenantId}`)
  : 'tenant not yet created');
const tenantsBefore = sql('select count(*) from tenants');

// Fresh, non-persistent Chrome context: registration must start unauthenticated.
const browser = await chromium.launch({ channel: 'chrome', headless: false });
const page = await (await browser.newContext({ viewport: { width: 1440, height: 900 } })).newPage();
const net = [];
page.on('response', async (r) => {
  if (!r.url().includes(':18082') || r.request().method() === 'OPTIONS') return;
  let body = '';
  try { body = (await r.text()).slice(0, 600); } catch { /* unavailable */ }
  net.push({ method: r.request().method(), url: r.url().replace('http://localhost:18082/api/v1', ''), status: r.status(), response: body.replace(/"(token|password|access_token|refresh_token)":"[^"]+"/g, '"$1":"<redacted>"') });
});

await page.goto('http://localhost:13111/en/auth/register', { waitUntil: 'domcontentloaded', timeout: 120000 });
await page.waitForSelector('#name', { timeout: 120000 });
await page.fill('#name', 'Audit Owner D');
await page.fill('#email', EMAIL);
await page.fill('#restaurantName', 'Audit D Onboarding');
await page.fill('#restaurantSlug', SLUG);
await page.fill('#phone', '+201000000004');
await page.fill('#password', PASSWORD);
await page.fill('#confirmPassword', PASSWORD);
await page.screenshot({ path: `${d.shot}/AUD-ONB-001-before-submit.png` });
// Probe 1 (observation): submit immediately after typing, before the debounced subdomain check returns.
await page.click('button[type=submit]:has-text("Create account")');
await page.waitForTimeout(6000);
await page.screenshot({ path: `${d.shot}/AUD-ONB-001-probe1-fast-submit.png` });
const probe1Posts = net.filter((c) => c.method === 'POST').length;
const probe1Text = (await page.locator('body').innerText()).includes('Please choose an available subdomain');
// Probe 2 (the test): wait for the UI to confirm availability, then submit like a normal user.
await page.getByText('Subdomain is available').first().waitFor({ timeout: 30000 });
await page.waitForTimeout(1500);
await page.click('button[type=submit]:has-text("Create account")');
await page.waitForTimeout(10000);
await page.screenshot({ path: `${d.shot}/AUD-ONB-001-after-submit.png` });
console.log(JSON.stringify({ probe1: { postsSent: probe1Posts, staleSubdomainErrorShown: probe1Text } }));
const finalUrl = page.url();

const tenantId = sql(`select tenant_id from users where email='${EMAIL}'`);
const after = measure(tenantId);
const tenantsAfter = sql('select count(*) from tenants');
fs.writeFileSync(`${d.net}/AUD-ONB-001.json`, JSON.stringify({ run: 'RUN-003', test: 'AUD-ONB-001', finding: 'F-001', finalUrl, calls: net }, null, 2));
fs.writeFileSync(`${d.db}/AUD-ONB-001-before-after.txt`, `tenants before=${tenantsBefore} after=${tenantsAfter}\nnew tenant_id=${tenantId || 'NONE'}\nBOOTSTRAP STATE AFTER REGISTRATION (no seeding performed)\n${after}\n`);
console.log(JSON.stringify({ finalUrl, tenantsBefore, tenantsAfter, tenantId, bootstrap: after.split('\n'), calls: net.map((c) => `${c.method} ${c.url} ${c.status}`) }, null, 2));
await browser.close();
