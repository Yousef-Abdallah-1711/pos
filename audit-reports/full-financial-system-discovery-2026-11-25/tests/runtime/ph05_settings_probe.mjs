// RUN-003 / PH-05 / AUD-SET-000: read-only enumeration of every settings page control in HEADED Chrome (tenant 4 owner). No submits.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-05');
const routes = ['settings', 'settings/restaurant', 'settings/account', 'settings/appearance', 'settings/billing', 'settings/access-control', 'settings/roles', 'settings/roles/new', 'settings/users', 'pos/settings', 'kds/settings'];
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/settings`);
const calls = []; page.on('response', (r) => { if (r.url().includes(':18082') && r.request().method() === 'GET') calls.push(`${r.status()} ${r.url().replace('http://localhost:18082/api/v1', '')}`); });
const out = [];
for (const rt of routes) {
  calls.length = 0;
  await page.goto(`${DASH}/en/dashboard/${rt}`, { waitUntil: 'domcontentloaded', timeout: 120000 });
  await page.waitForTimeout(5000);
  const controls = await page.evaluate(() => [...document.querySelectorAll('input,select,textarea,button,[role=switch],[role=combobox],[role=tab]')].map((e) => {
    const lab = (e.labels && e.labels[0]?.innerText) || e.getAttribute('aria-label') || e.closest('label')?.innerText || '';
    return { tag: e.tagName.toLowerCase(), type: e.type || e.getAttribute('role') || '', name: e.name || e.id || '', label: lab.trim().slice(0, 60), text: (e.tagName === 'BUTTON' ? e.innerText.trim().slice(0, 40) : ''), value: e.type === 'password' ? '<hidden>' : String(e.value ?? '').slice(0, 60), disabled: !!e.disabled };
  }));
  await page.screenshot({ path: `${d.shot}/AUD-SET-000-${rt.replace(/\//g, '_')}.png`, fullPage: true });
  out.push({ route: rt, finalUrl: page.url(), h1: (await page.locator('h1').first().innerText().catch(() => '')), controls, apiGets: [...calls] });
  console.log(rt, '->', page.url().replace(DASH, ''), 'controls', controls.length);
}
fs.writeFileSync(`${d.net}/AUD-SET-000-controls.json`, JSON.stringify(out, null, 1));
await ctx.close();
