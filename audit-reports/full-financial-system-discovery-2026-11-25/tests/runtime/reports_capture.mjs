import { login, EV, DASH } from './lib.mjs';
import fs from 'fs';
const [tag, email, suffix, ...pages] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
for (const pg of pages) {
  const p = await ctx.newPage(); const calls = [];
  p.on('response', async r => { const u = r.url(); if (u.includes(':18082') && r.request().method() === 'GET' && /reports|trial-balance|ledger|accounting\/(tax|accounts)|kpi/.test(u)) { try { calls.push({ u: u.replace('http://localhost:18082/api/v1', ''), s: r.status(), body: await r.text() }); } catch {} } });
  await p.goto(DASH + '/en/dashboard/' + pg, { waitUntil: 'networkidle' }).catch(() => {}); await p.waitForTimeout(4000);
  const name = pg.replace(/\W+/g, '_');
  fs.writeFileSync(`${EV}/reports/${tag}-${suffix}-${name}.json`, JSON.stringify({ page: pg, text: (await p.locator('main, body').first().innerText()).slice(0, 4000), calls }, null, 1));
  await p.screenshot({ path: `${EV}/screenshots/rep-${tag}-${suffix}-${name}.png`, fullPage: true });
  console.log(pg, '| api:', calls.map(c => `${c.s} ${c.u.slice(0, 70)}`).join(' ; '));
  await p.close();
}
await b.close();
