import { openChrome, ensureLoggedIn, DASH } from './chrome_lib.mjs';
const { ctx, page } = await openChrome();
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard`);
for (const r of ['hr/employees', 'hr/employees/new', 'hr/leaves', 'settings/roles', 'rbac/roles']) {
  await page.goto(`${DASH}/en/dashboard/${r}`, { waitUntil: 'domcontentloaded', timeout: 120000 }).catch(() => {}); await page.waitForTimeout(6000);
  const i = await page.evaluate(() => ({ h: (document.querySelector('h1,h2')?.innerText || '').slice(0, 50), n404: /could not be found/i.test(document.body.innerText), buttons: [...document.querySelectorAll('main button, main a[href]')].map((b) => (b.innerText || '').trim().slice(0, 28)).filter(Boolean).slice(0, 14), inputs: document.querySelectorAll('main input:not([type=hidden])').length }));
  console.log(r, JSON.stringify(i));
}
await ctx.close();
