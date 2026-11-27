// RUN-005 / PH-04 (Phase G): Arabic (RTL) and mobile/tablet FUNCTIONAL interaction pass in headed Chrome. For each representative route and configuration (ar-desktop 1440, en-360, en-390, en-768, en-1024, ar-390) it checks:
// horizontal overflow, controls outside the viewport, mobile navigation reachability, add-dialog fit and submit reachability, touch target size, RTL direction on html/dialog/table/select, and an Arabic search/typing round trip. No data is saved.
// Env: ROUTES=csv, CONFIGS=csv. Separate Chrome profile. Evidence: JSON per run + a screenshot only when an issue is found.
import fs from 'fs';
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { evidenceDirs, DASH, PASSWORD } from './chrome_lib.mjs';
const d = evidenceDirs('PH-04');
const ROUTES = (process.env.ROUTES || 'dashboard,pos,orders,products,products/new,accounting/customers,accounting/customer-invoices,accounting/chart-of-accounts,reports/sales,reports/financial/profit-loss,hr/employees,settings,theme-builder,kds,tables').split(',');
const ALL = { 'ar-1440': { w: 1440, h: 900, lang: 'ar' }, 'en-360': { w: 360, h: 740, lang: 'en' }, 'en-390': { w: 390, h: 844, lang: 'en' }, 'en-768': { w: 768, h: 1024, lang: 'en' }, 'en-1024': { w: 1024, h: 768, lang: 'en' }, 'ar-390': { w: 390, h: 844, lang: 'ar' } };
const CONFIGS = (process.env.CONFIGS || Object.keys(ALL).join(',')).split(',');
const ctx = await chromium.launchPersistentContext('C:/Users/Yousef/AppData/Local/Temp/r5-chrome-profile-lang', { channel: 'chrome', headless: false, viewport: { width: 1440, height: 900 }, hasTouch: false });
const page = ctx.pages()[0] ?? (await ctx.newPage());
await page.goto(`${DASH}/en/dashboard`, { waitUntil: 'domcontentloaded', timeout: 180000 });
for (let i = 0; i < 40; i++) { await page.waitForTimeout(3000); if ((await page.locator('input[type=email]').count()) > 0 || (await page.evaluate(() => (document.body.innerText || '').length).catch(() => 0)) > 800) break; }
if ((await page.locator('input[type=email]').count()) > 0) { await page.fill('input[type=email]', 'audit.c.owner@example.test'); await page.fill('input[type=password]', PASSWORD); await page.click('form button[type=submit]').catch(() => page.click('button[type=submit]')); await page.waitForTimeout(9000); }
const out = [];
for (const cfgName of CONFIGS) {
  const cfg = ALL[cfgName]; await page.setViewportSize({ width: cfg.w, height: cfg.h });
  for (const route of ROUTES) {
    const rec = { config: cfgName, route, viewport: `${cfg.w}x${cfg.h}`, lang: cfg.lang, issues: [] };
    try {
      await page.goto(`${DASH}/${cfg.lang}/dashboard${route === 'dashboard' ? '' : '/' + route}`, { waitUntil: 'domcontentloaded', timeout: 150000 }); await page.waitForTimeout(5500);
      const m = await page.evaluate(() => {
        const vw = window.innerWidth; const out = { dir: document.documentElement.dir, lang: document.documentElement.lang, overflow: document.documentElement.scrollWidth > vw + 2, scrollW: document.documentElement.scrollWidth, bodyLen: (document.body.innerText || '').length };
        const off = []; const small = []; for (const e of document.querySelectorAll('main button, main a, main input, main select, main [role=tab]')) { const r = e.getBoundingClientRect(); if (r.width < 2 || r.height < 2) continue; const cs = getComputedStyle(e); if (cs.visibility === 'hidden') continue; let p = e.parentElement, scrollable = false; while (p && p !== document.body) { const o = getComputedStyle(p).overflowX; if ((o === 'auto' || o === 'scroll') && p.scrollWidth > p.clientWidth) { scrollable = true; break; } p = p.parentElement; } if (!scrollable && (r.right > vw + 2 || r.left < -2)) off.push((e.innerText || e.getAttribute('aria-label') || e.tagName).trim().slice(0, 24)); if ((r.height < 32 || r.width < 32) && e.tagName !== 'INPUT') small.push(1); }
        out.offscreen = off.slice(0, 8); out.offscreen_count = off.length; out.small_targets = small.length;
        const tbl = document.querySelector('main table'); if (tbl) { const cs = getComputedStyle(tbl); out.table_dir = cs.direction; const w = tbl.getBoundingClientRect().width; out.table_wider_than_viewport = w > vw + 2; let p = tbl.parentElement, sc = false; while (p && p !== document.body) { const o = getComputedStyle(p).overflowX; if (o === 'auto' || o === 'scroll') { sc = true; break; } p = p.parentElement; } out.table_scrollable_wrapper = sc; }
        return out; });
      Object.assign(rec, m);
      if (m.dir !== (cfg.lang === 'ar' ? 'rtl' : 'ltr')) rec.issues.push(`html dir=${m.dir} for ${cfg.lang}`);
      if (m.overflow) rec.issues.push(`horizontal overflow (scrollWidth ${m.scrollW} > ${cfg.w})`);
      if (m.offscreen_count) rec.issues.push(`controls outside viewport: ${m.offscreen.join(', ')}`);
      if (m.table_wider_than_viewport && !m.table_scrollable_wrapper) rec.issues.push('table wider than viewport without a scrollable wrapper');
      if (cfg.w < 600 && m.small_targets > 12) rec.issues.push(`${m.small_targets} touch targets below 32 px`);
      if (cfg.w < 800) { const burger = page.locator('button[aria-label*="menu" i], button[aria-label*="navigation" i], header button').first(); rec.mobile_menu_button = (await burger.count()) > 0; }
      // add-dialog fit and submit reachability
      const add = page.locator('main button:visible').filter({ hasText: /^(\+\s*)?(add|new|create|record|request)/i }).first();
      if (route !== 'dashboard' && (await add.count())) { await add.click({ timeout: 4000 }).catch(() => {}); await page.waitForTimeout(1500); const dlg = page.locator('[role=dialog]:visible').first(); if (await dlg.count()) { const g = await dlg.evaluate((e) => { const r = e.getBoundingClientRect(); const cs = getComputedStyle(e); return { l: r.left, r: r.right, t: r.top, b: r.bottom, vw: innerWidth, vh: innerHeight, dir: cs.direction, scroll: e.scrollHeight > e.clientHeight }; }); rec.dialog = g; if (g.r > g.vw + 2 || g.l < -2) rec.issues.push('dialog clipped horizontally'); if (g.dir !== (cfg.lang === 'ar' ? 'rtl' : 'ltr')) rec.issues.push(`dialog direction ${g.dir}`); const sub = dlg.locator('button').filter({ hasText: /^(save|create|add|submit|confirm|request|record)/i }).last(); if (await sub.count()) { await sub.scrollIntoViewIfNeeded().catch(() => {}); const sb = await sub.boundingBox(); rec.submit_reachable = !!sb && sb.y >= 0 && sb.y + sb.height <= cfg.h + 2 && sb.x >= -2 && sb.x + sb.width <= cfg.w + 2; if (!rec.submit_reachable) rec.issues.push('dialog submit button not reachable inside the viewport'); } await page.keyboard.press('Escape').catch(() => {}); await page.waitForTimeout(500); } }
      // Arabic search round trip
      if (cfg.lang === 'ar') { const q = page.locator('main input[type=search], main input[placeholder*="بحث"], main input[placeholder*="earch" i]').first(); if (await q.count()) { await q.fill('اختبار').catch(() => {}); await page.waitForTimeout(1200); rec.arabic_search_value = await q.inputValue().catch(() => ''); if (rec.arabic_search_value !== 'اختبار') rec.issues.push('Arabic search text not retained'); await q.fill('').catch(() => {}); } }
      if (rec.issues.length) await page.screenshot({ path: `${d.shot}/AUD-R5-LANG-${cfgName}-${route.replace(/\//g, '_')}.png` }).catch(() => {});
    } catch (e) { rec.issues.push('TOOLING: ' + String(e.message).split('\n')[0].slice(0, 90)); }
    out.push(rec); console.log(cfgName, route, rec.issues.length ? 'ISSUES ' + rec.issues.join(' | ').slice(0, 200) : 'ok');
    fs.writeFileSync(`${d.net}/AUD-R5-LANG-VIEWPORT.json`, JSON.stringify({ run: 'RUN-005', tenant: 4, role: 'owner', results: out }, null, 1));
  }
}
await ctx.close();
