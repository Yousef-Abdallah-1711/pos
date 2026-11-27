// RUN-004 / PH-00: read-only smoke of the updated (uncommitted) landing app in HEADED Chrome: routes EN/AR, language switch, nav links, register and login pages, console errors, horizontal overflow at desktop and 390px. No form is submitted.
import fs from 'fs';
import { openChrome, evidenceDirs, LAND } from './chrome_lib.mjs';
const d = evidenceDirs('PH-00');
const { ctx, page } = await openChrome();
let errs = []; page.on('pageerror', (e) => errs.push('pageerror: ' + String(e.message).slice(0, 120))); page.on('console', (m) => { if (m.type() === 'error') errs.push('console: ' + m.text().slice(0, 140)); });
const paths = ['/en', '/ar', '/en/auth/login', '/en/auth/register', '/ar/auth/login', '/ar/auth/register', '/en/pricing', '/en/features', '/en/demo', '/en/how-it-works', '/en/contact', '/en/does-not-exist'];
const out = [];
for (const vp of [{ w: 1440, h: 900, n: 'desktop' }, { w: 390, h: 844, n: 'mobile' }]) {
  await page.setViewportSize({ width: vp.w, height: vp.h });
  for (const p of paths) {
    errs = [];
    let status = null; try { const r = await page.goto(LAND + p, { waitUntil: 'domcontentloaded', timeout: 120000 }); status = r?.status(); } catch (e) { status = 'ERR ' + String(e.message).slice(0, 60); }
    await page.waitForTimeout(2500);
    const info = await page.evaluate(() => ({ title: document.title, h1: (document.querySelector('h1')?.innerText || '').slice(0, 70), dir: document.documentElement.dir, lang: document.documentElement.lang, overflow: document.documentElement.scrollWidth > window.innerWidth + 2, len: (document.body.innerText || '').length, links: [...document.querySelectorAll('a[href^="/"]')].length })).catch(() => ({}));
    out.push({ viewport: vp.n, path: p, status, ...info, errors: errs.slice(0, 3) });
    if (vp.n === 'desktop' && ['/en', '/ar'].includes(p)) await page.screenshot({ path: `${d.shot}/AUD-R4-LANDING-${p.slice(1)}.png`, fullPage: false }).catch(() => {});
  }
}
// language switch on /en
await page.setViewportSize({ width: 1440, height: 900 }); await page.goto(LAND + '/en', { waitUntil: 'domcontentloaded' }); await page.waitForTimeout(2500);
const sw = page.locator('a[href^="/ar"], button:has-text("العربية"), a:has-text("العربية"), button:has-text("AR")').first();
let langSwitch = 'no switch control found';
if (await sw.count()) { await sw.click({ timeout: 5000 }).catch(() => {}); await page.waitForTimeout(3000); langSwitch = `after click url=${page.url().replace(LAND, '')} dir=${await page.evaluate(() => document.documentElement.dir)}`; }
fs.writeFileSync(`${d.net}/AUD-R4-LANDING-SMOKE.json`, JSON.stringify({ run: 'RUN-004', note: 'working tree contains uncommitted landing changes', langSwitch, results: out }, null, 1));
for (const r of out) console.log(r.viewport, r.path, r.status, 'len', r.len, 'dir', r.dir, 'overflow', r.overflow, 'errs', r.errors?.length);
console.log('lang switch:', langSwitch);
await ctx.close();
