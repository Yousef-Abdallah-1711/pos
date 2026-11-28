import { login, DASH } from './lib.mjs';
const { b, ctx } = await login('audit.c.owner@example.test', 'C');
const p = await ctx.newPage();
await p.goto(DASH + '/en/dashboard/hr/payroll', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
await p.click('button:has-text("Create Salary")'); await p.waitForTimeout(1000);
for (const t of ['General', 'Earnings', 'Deductions']) {
  await p.click(`button:has-text("${t}")`); await p.waitForTimeout(500);
  const f = await p.$$eval('label,input,select', els => els.filter(e => e.offsetParent !== null).map(e => e.tagName === 'LABEL' ? 'L:' + e.innerText.trim().slice(0, 40) : `${e.tagName}[${e.type}] name=${e.name}`));
  console.log('== ' + t + '\n' + f.slice(10).join(' | '));
}
await b.close();
