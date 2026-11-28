import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import fs from 'fs';
export const ROOT = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25';
export const EV = ROOT + '/evidence';
export const DASH = 'http://localhost:13112', LAND = 'http://localhost:13111', API = 'http://localhost:18082/api/v1';
export const PW = 'AuditPass#2026A';
export async function login(email, tag) {
  const state = `${EV}/state-${tag}.json`;
  const b = await chromium.launch({ headless: true });
  let ctx;
  if (fs.existsSync(state)) { ctx = await b.newContext({ storageState: state, viewport: { width: 1440, height: 900 } }); }
  else {
    ctx = await b.newContext({ viewport: { width: 1440, height: 900 } });
    const p = await ctx.newPage();
    await p.goto(LAND + '/en/auth/login', { waitUntil: 'networkidle' }); await p.waitForTimeout(1000);
    await p.fill('input[type=email]', email); await p.fill('input[type=password]', PW);
    await p.click('button[type=submit]:has-text("Log in"), button[type=submit]:has-text("Sign in"), button[type=submit]');
    await p.waitForURL(/13112/, { timeout: 30000 }).catch(() => {});
    await p.waitForTimeout(4000);
    await ctx.storageState({ path: state });
    await p.close();
  }
  return { b, ctx };
}
export function netLog(page, bag) {
  page.on('response', async r => { const u = r.url(); if (u.includes(':18082') && r.request().method() !== 'OPTIONS') { let body = ''; try { body = (await r.text()).slice(0, 400).replace(/"(token|password|access_token|refresh_token)":"[^"]+"/g, '"$1":"<redacted>"'); } catch {} bag.push({ m: r.request().method(), u: u.replace('http://localhost:18082/api/v1', ''), s: r.status(), body, req: (r.request().postData() || '').slice(0, 1500).replace(/"password":"[^"]+"/g, '"password":"<redacted>"') }); } });
}
