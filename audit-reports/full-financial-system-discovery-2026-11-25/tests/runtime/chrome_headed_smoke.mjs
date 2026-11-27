// RUN-003 / PH-02: proves real, installed Google Chrome (headed) can drive the ISOLATED audit stack.
// Target: isolated stack only (landing 13111, dashboard 13112, API 18082). Side effects: none (login + read-only navigation).
// Profile + session cookies live OUTSIDE the shareable evidence tree (sensitive).
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import fs from 'fs';

const ROOT = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25';
const OUT = `${ROOT}/evidence/screenshots/RUN-003/PH-02`;
const NET = `${ROOT}/evidence/network/RUN-003/PH-02`;
const PROFILE = process.env.AUDIT_CHROME_PROFILE; // dedicated profile dir, never the owner's profile
const EMAIL = 'audit.c.owner@example.test';
const PASSWORD = 'AuditPass#2026A'; // synthetic audit-only credential (documented in lib.mjs)

if (!PROFILE) throw new Error('AUDIT_CHROME_PROFILE env var is required');
fs.mkdirSync(OUT, { recursive: true });
fs.mkdirSync(NET, { recursive: true });

const ctx = await chromium.launchPersistentContext(PROFILE, {
  channel: 'chrome', // installed Google Chrome, not bundled Chromium
  headless: false,
  viewport: { width: 1440, height: 900 },
});
const page = ctx.pages()[0] ?? (await ctx.newPage());
const calls = [];
page.on('response', (r) => {
  const url = r.url();
  if (url.includes(':18082') && r.request().method() !== 'OPTIONS') {
    calls.push({ method: r.request().method(), url: url.replace('http://localhost:18082/api/v1', ''), status: r.status() });
  }
});

const browserVersion = ctx.browser()?.version() ?? 'persistent-context';
const userAgent = await page.evaluate(() => navigator.userAgent);

await page.goto('http://localhost:13111/en/auth/login', { waitUntil: 'networkidle' });
await page.screenshot({ path: `${OUT}/AUD-ENV-001-login-before.png` });
await page.fill('input[type=email]', EMAIL);
await page.fill('input[type=password]', PASSWORD);
await page.click('button[type=submit]');
await page.waitForURL(/13112/, { timeout: 45000 }).catch(() => {});
await page.waitForTimeout(4000);
await page.screenshot({ path: `${OUT}/AUD-ENV-001-after-login.png` });

await page.goto('http://localhost:13112/en/dashboard/reports/financial', { waitUntil: 'networkidle' });
await page.waitForTimeout(3000);
await page.screenshot({ path: `${OUT}/AUD-ENV-001-financial-report.png` });

const result = { browserVersion, userAgent, headed: true, finalUrl: page.url(), apiCalls: calls.length, failedCalls: calls.filter((c) => c.status >= 400) };
fs.writeFileSync(`${NET}/AUD-ENV-001.json`, JSON.stringify({ result, calls }, null, 2));
console.log(JSON.stringify(result, null, 2));
await ctx.close();
