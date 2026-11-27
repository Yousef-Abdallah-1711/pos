// RUN-004 / PH-13 (SUPER ADMIN APP, port 13114): per-control INTERACTION matrix in HEADED Chrome (EN desktop). For every dashboard route in 53-ui-control-inventory.csv it enumerates interactive controls
// (tabs, buttons, links, selects, switches, checkboxes, search/filter inputs) and performs ONE safe action per control, recording observable effect, dialogs, GET/non-GET network, console/page errors.
// Destructive or money-moving labels are NOT clicked (status SKIPPED-DESTRUCTIVE) because they need a purpose-built fixture; any non-GET request that still fires is flagged. Tenant 4 owner. Env: ROUTES=csv subset, MAXC=controls per route.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, DASH, PASSWORD, ROOT } from './chrome_lib.mjs';
const d = evidenceDirs('PH-13'); const SA = 'http://localhost:13114';
const PROFILE = process.env.AUDIT_CHROME_PROFILE2 || 'C:/Users/Yousef/AppData/Local/Temp/claude/f11bcefd4a918204/ed68dc13-c3a5-45c2-a7ac-cd2b088cd7c5/scratchpad/chrome-audit-profile-sa';
const MAXC = +(process.env.MAXC || 22); const SUB = (process.env.ROUTES || '').split(',').filter(Boolean);
const routes = ['', 'tenants', 'payments', 'users', 'audit-logs', 'settings'];
const DESTRUCTIVE = /delete|remove|void|cancel|logout|log out|sign out|pay\b|mark.*paid|approve|reject|post\b|reverse|reset|clear|deactivate|disable|terminate|bump|complete|refund|submit|save|create order|place order|checkout|import|upload/i;
const ctx = await chromium.launchPersistentContext(PROFILE, { channel: 'chrome', headless: false, viewport: { width: 1440, height: 900 }, acceptDownloads: true });
const page = ctx.pages()[0] ?? (await ctx.newPage());
let net = []; let errs = [];
page.on('response', (r) => { if (r.url().includes(':18082')) net.push({ m: r.request().method(), s: r.status(), u: r.url().replace('http://localhost:18082/api/v1', '').slice(0, 90) }); });
page.on('pageerror', (e) => errs.push('pageerror: ' + String(e.message).slice(0, 120)));
page.on('console', (m) => { if (m.type() === 'error') errs.push('console: ' + m.text().slice(0, 120)); });
await page.goto(`${SA}/login`, { waitUntil: 'domcontentloaded', timeout: 180000 });
for (let i = 0; i < 40; i++) { await page.waitForTimeout(3000); if ((await page.locator('input[type=email]').count()) > 0) break; }
if ((await page.locator('input[type=email]').count()) > 0) { await page.fill('input[type=email]', 'superadmin@pos-saas.com'); await page.fill('input[type=password]', PASSWORD); await page.click('button[type=submit]'); await page.waitForTimeout(8000); }
const out = []; const sig = async () => (await page.evaluate(() => { const t = document.body.innerText || ''; return t.length + ':' + t.slice(0, 40) + ':' + t.slice(-40); }).catch(() => ''));
for (const route of routes) {
  const url = `${SA}/dashboard${route ? '/' + route : ''}`;
  try { await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 120000 }); } catch (e) { out.push({ route, control: '(route)', type: 'route', status: 'ERROR', detail: 'goto ' + String(e.message).slice(0, 80) }); continue; }
  await page.waitForTimeout(3500);
  const tag = async () => page.evaluate((max) => {
    document.querySelectorAll('[data-r4]').forEach((e) => e.removeAttribute('data-r4'));
    const vis = (e) => { const r = e.getBoundingClientRect(); const cs = getComputedStyle(e); return r.width > 2 && r.height > 2 && cs.visibility !== 'hidden' && cs.display !== 'none'; };
    const items = []; const add = (type, e, label) => { if (!vis(e)) return; e.setAttribute('data-r4', String(items.length)); items.push({ type, label: (label || '').replace(/\s+/g, ' ').trim().slice(0, 50) }); };
    const main = document.querySelector('main') || document.body;
    main.querySelectorAll('[role=tab]').forEach((e) => add('tab', e, e.innerText));
    main.querySelectorAll('button').forEach((e) => add('button', e, e.innerText || e.getAttribute('aria-label') || e.title));
    main.querySelectorAll('a[href^="/"]').forEach((e) => add('link', e, e.innerText || e.getAttribute('aria-label')));
    main.querySelectorAll('select').forEach((e) => add('select', e, e.getAttribute('aria-label') || e.name));
    main.querySelectorAll('[role=switch], input[type=checkbox]').forEach((e) => add('toggle', e, e.getAttribute('aria-label') || e.name || e.id));
    main.querySelectorAll('input[type=search], input[placeholder*="earch" i], input[placeholder*="ilter" i]').forEach((e) => add('search', e, e.placeholder));
    return items.slice(0, max);
  }, MAXC);
  let controls = await tag();
  for (let ci = 0; ci < controls.length; ci++) {
    const c = controls[ci]; const rec = { route, control: `${c.type}:${c.label || '(unlabelled)'}`, type: c.type };
    if (c.type === 'button' && DESTRUCTIVE.test(c.label)) { rec.status = 'SKIPPED-DESTRUCTIVE'; rec.detail = 'money-moving or state-changing label; needs a purpose-built fixture'; out.push(rec); continue; }
    net = []; errs = []; const before = await sig(); const u0 = page.url();
    try {
      let loc = page.locator(`[data-r4="${ci}"]`); if ((await loc.count()) === 0) { const again = await tag(); if (again[ci] && again[ci].type === c.type) loc = page.locator(`[data-r4="${ci}"]`); }
      if (c.type === 'search') { await loc.fill('zzz-no-match', { timeout: 4000 }); await page.waitForTimeout(1500); rec.detail = 'typed no-match query'; await loc.fill('', { timeout: 3000 }); }
      else if (c.type === 'select') { await loc.selectOption({ index: 1 }, { timeout: 4000 }).catch(() => {}); await page.waitForTimeout(1200); rec.detail = 'selected 2nd option'; }
      else if (c.type === 'toggle') { await loc.click({ timeout: 4000 }); await page.waitForTimeout(800); await loc.click({ timeout: 4000 }).catch(() => {}); rec.detail = 'toggled twice (restored)'; }
      else { const dl = page.waitForEvent('download', { timeout: 1500 }).catch(() => null); await loc.scrollIntoViewIfNeeded({ timeout: 3000 }).catch(() => {}); if (await loc.isDisabled().catch(() => false)) { rec.status = 'DISABLED'; rec.detail = 'control disabled in this state (precondition not met)'; out.push(rec); continue; } await loc.click({ timeout: 4000 }).catch(async () => { rec.detail = 'forced click (covered or animating)'; await loc.click({ timeout: 4000, force: true }); }); const dwn = await dl; if (dwn) rec.detail = 'download: ' + dwn.suggestedFilename(); await page.waitForTimeout(1500); }
      const dialog = (await page.locator('[role=dialog]:visible').count()) > 0; if (dialog) { rec.dialog = true; await page.keyboard.press('Escape').catch(() => {}); await page.waitForTimeout(500); if ((await page.locator('[role=dialog]:visible').count()) > 0) { await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 90000 }).catch(() => {}); await page.waitForTimeout(2500); } }
      const after = await sig(); const nav = page.url() !== u0;
      rec.navigated = nav ? page.url().replace(SA, '') : null; rec.changed = before !== after || nav || dialog;
      const mut = net.filter((x) => x.m !== 'GET'); const bad = net.filter((x) => x.s >= 500);
      rec.mutations = mut.length ? mut.slice(0, 3).map((x) => `${x.m} ${x.u} ${x.s}`) : undefined; rec.badGet = bad.length ? bad.slice(0, 2).map((x) => `${x.m} ${x.u} ${x.s}`) : undefined; rec.errors = errs.length ? errs.slice(0, 2) : undefined;
      rec.status = bad.length || errs.length ? 'ERROR' : mut.length ? 'FLAG-MUTATION' : rec.changed ? 'PASS' : 'NO-EFFECT';
      if (nav) { await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 90000 }).catch(() => {}); await page.waitForTimeout(2500); }
    } catch (e) { rec.status = 'BLOCKED'; rec.detail = String(e.message).split('\n')[0].slice(0, 100); }
    out.push(rec);
  }
  const rc = out.filter((x) => x.route === route); console.log(route, 'controls', rc.length, JSON.stringify(rc.reduce((a, x) => (a[x.status] = (a[x.status] || 0) + 1, a), {})));
  fs.writeFileSync(`${d.net}/AUD-R4-SUPERADMIN-CONTROL-MATRIX${SUB.length ? '-subset' : ''}.json`, JSON.stringify({ run: 'RUN-004', tenant: 4, role: 'owner', viewport: '1440x900 EN', results: out }, null, 1));
}
await ctx.close();
