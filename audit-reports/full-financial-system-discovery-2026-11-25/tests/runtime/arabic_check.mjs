import { login, EV, DASH } from './lib.mjs';
const { b, ctx } = await login('audit.c.owner@example.test', 'C');
const out = [];
for (const path of ['ar/dashboard', 'ar/dashboard/pos', 'ar/dashboard/accounting/journal-entries', 'ar/dashboard/reports/financial/balance-sheet', 'ar/dashboard/hr/employees', 'ar/dashboard/products']) {
  const p = await ctx.newPage(); const errs = []; p.on('pageerror', e => errs.push(e.message.slice(0, 80)));
  await p.goto(DASH + '/' + path, { waitUntil: 'networkidle' }).catch(() => {}); await p.waitForTimeout(2500);
  const dir = await p.evaluate(() => document.documentElement.dir + '/' + document.documentElement.lang);
  const txt = (await p.locator('body').innerText()); const arabic = (txt.match(/[؀-ۿ]/g) || []).length; const latinHeader = /Journal Entries|Balance Sheet|Cashier POS/.test(txt);
  out.push(`${path} dir/lang=${dir} arabicChars=${arabic} englishLeftover=${latinHeader} pageerrors=${errs.length}`);
  await p.screenshot({ path: `${EV}/screenshots/ar-${path.replace(/\W+/g, '_')}.png` }); await p.close();
}
console.log(out.join('\n')); await b.close();
