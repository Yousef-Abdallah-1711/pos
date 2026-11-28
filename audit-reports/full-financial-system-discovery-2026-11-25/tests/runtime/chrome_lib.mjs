// Shared helpers for RUN-003+ headed-Chrome audit scripts. Isolated stack only. No product code touched.
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { execSync } from 'child_process';
import fs from 'fs';

export const ROOT = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25';
export const DASH = 'http://localhost:13112';
export const LAND = 'http://localhost:13111';
export const API = 'http://localhost:18082/api/v1';
export const PASSWORD = 'AuditPass#2026A'; // synthetic audit-only credential

export const RUN_ID = process.env.RUN_ID || 'RUN-005'; // RUN-003 scripts keep their evidence by running with RUN_ID=RUN-003
export function evidenceDirs(phase) {
  const d = { shot: `${ROOT}/evidence/screenshots/${RUN_ID}/${phase}`, net: `${ROOT}/evidence/network/${RUN_ID}/${phase}`, db: `${ROOT}/evidence/database/${RUN_ID}/${phase}` };
  Object.values(d).forEach((p) => fs.mkdirSync(p, { recursive: true }));
  return d;
}

// Read-only/diagnostic SQL against the ISOLATED audit DB only.
export const sql = (q) => execSync(`docker exec pos-p2a-ws05-e2e-postgres-1 psql -U postgres -d pos_phase2a_e2e -Atc "${q}"`, { encoding: 'utf8' }).trim();

export async function openChrome() {
  const profile = process.env.AUDIT_CHROME_PROFILE || 'C:/Users/Yousef/AppData/Local/Temp/claude/f11bcefd4a918204/ed68dc13-c3a5-45c2-a7ac-cd2b088cd7c5/scratchpad/chrome-audit-profile'; // dedicated profile outside evidence/, never the owner's
  if (!profile) throw new Error('AUDIT_CHROME_PROFILE required');
  const ctx = await chromium.launchPersistentContext(profile, { channel: 'chrome', headless: false, viewport: { width: 1440, height: 900 } });
  const page = ctx.pages()[0] ?? (await ctx.newPage());
  return { ctx, page };
}

// Authenticates through the real login UI only when the app redirects to it.
export async function ensureLoggedIn(page, email, targetUrl) {
  // Always authenticate against a known-good protected route (a 404 target never redirects to login).
  for (let attempt = 0; attempt < 3; attempt++) {
    await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded', timeout: 180000 });
    let needsLogin = false;
    for (let i = 0; i < 40; i++) { // first Next.js compile of the login page can take > 60 s
      await page.waitForTimeout(3000);
      if (page.url().includes('/auth/login') || (await page.locator('input[type=email]').count()) > 0) { needsLogin = true; break; }
      const len = await page.evaluate(() => (document.body.innerText || '').length).catch(() => 0);
      if (page.url().startsWith(DASH) && len > 400 && i >= 2) break; // dashboard rendered
    }
    if (!needsLogin) break;
    await page.fill('input[type=email]', email);
    await page.fill('input[type=password]', PASSWORD);
    await page.click('button[type=submit]');
    await page.waitForTimeout(8000);
  }
  await page.goto(targetUrl, { waitUntil: 'domcontentloaded', timeout: 180000 });
  await page.waitForTimeout(3000);
}

// Records every mutating API call (method/url/status/request/response) with secrets redacted.
export function recordMutations(page, bag) {
  page.on('response', async (r) => {
    if (!r.url().includes(':18082') || ['GET', 'OPTIONS'].includes(r.request().method())) return;
    let body = '';
    try { body = (await r.text()).slice(0, 600); } catch { /* body unavailable */ }
    const redact = (s) => s.replace(/"(password|token|access_token|refresh_token)":"[^"]+"/g, '"$1":"<redacted>"');
    bag.push({ method: r.request().method(), url: r.url().replace(API, ''), status: r.status(), request: redact((r.request().postData() || '').slice(0, 1500)), response: redact(body) });
  });
}
