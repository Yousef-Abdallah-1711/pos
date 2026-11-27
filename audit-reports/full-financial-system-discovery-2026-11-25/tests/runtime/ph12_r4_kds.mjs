// RUN-004 / PH-12: MET-KDS-003 (KDS analytics + stats) in HEADED Chrome. Enables the analytics panel (diagnostic tenant 4 setting), bumps one in-progress order to completed via the KDS UI, captures /kds/stats and /kds/analytics, screenshots the panel.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, DASH, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
const { ctx, page } = await openChrome();
const calls = []; page.on('response', async (r) => { const u = r.url(); if (u.includes(':18082') && /kds/.test(u)) { let b = ''; try { b = (await r.text()).slice(0, 6000); } catch {} calls.push({ m: r.request().method(), u: u.replace('http://localhost:18082/api/v1', ''), s: r.status(), b }); } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/kds/settings`);
await page.waitForTimeout(6000);
const toggle = page.getByText('Show kitchen analytics panel').first();
const sw = page.locator('[role=switch], input[type=checkbox]');
console.log('switches', await sw.count());
await page.screenshot({ path: `${d.shot}/AUD-R4-KDS-settings-before.png`, fullPage: true }).catch(() => {});
// click the last switch (analytics panel is the 2nd toggle after alert sound)
const n = await sw.count(); const states = []; for (let i = 0; i < n; i++) states.push(await sw.nth(i).getAttribute('aria-checked').catch(() => null));
console.log('states', states);
await sw.nth(n - 1).click({ timeout: 8000 }).catch((e) => console.log('toggle click failed', e.message.slice(0, 80)));
await page.getByRole('button', { name: /save settings/i }).click({ timeout: 8000 }).catch((e) => console.log('save failed', e.message.slice(0, 80)));
await page.waitForTimeout(4000);
console.log('db setting:', sql(`select show_analytics_panel from kds_settings where tenant_id=4`));
await page.goto(`${DASH}/en/dashboard/kds`, { waitUntil: 'domcontentloaded', timeout: 120000 }); await page.waitForTimeout(10000);
await page.screenshot({ path: `${d.shot}/AUD-R4-KDS-before-bump.png`, fullPage: true }).catch(() => {});
const before = { stats: sql(`select status,count(*) from kds_queue where tenant_id=4 group by 1 order by 1`) };
// bump the in-progress order -> ready -> completed through the UI
for (let step = 0; step < 2; step++) {
  const card = page.locator('div,article').filter({ hasText: step === 0 ? 'In Progress -> Ready' : 'Ready -> Completed' }).last();
  const btn = page.getByRole('button', { name: /bump|complete/i }).filter({ hasText: /bump/i });
  const target = page.locator('button', { hasText: 'Bump' }).filter({ has: page.locator(':scope') });
  const label = step === 0 ? 'In Progress -> Ready' : 'Ready -> Completed';
  const b = page.locator(`button:has-text("Bump"):near(:text("${label}"))`).first();
  const ok = await b.click({ timeout: 8000 }).then(() => true).catch(() => false);
  console.log('bump step', step, label, ok); await page.waitForTimeout(4000);
}
await page.waitForTimeout(4000);
await page.screenshot({ path: `${d.shot}/AUD-R4-KDS-after-bump.png`, fullPage: true }).catch(() => {});
const vis = (await page.locator('body').innerText()).slice(0, 3500);
const after = { stats: sql(`select status,count(*) from kds_queue where tenant_id=4 group by 1 order by 1`), completed_rows: sql(`select id,order_id,received_at,completed_at, round(extract(epoch from (completed_at-received_at))/60.0,2) from kds_queue where tenant_id=4 and status='completed' order by id`), avg_all: sql(`select round(avg(extract(epoch from (completed_at-received_at))/60.0)::numeric,4) from kds_queue where tenant_id=4 and status='completed' and completed_at is not null`), completed_today: sql(`select count(*) from kds_queue where tenant_id=4 and status='completed' and completed_at>=current_date`), hour_buckets: sql(`select extract(hour from received_at)::int h,count(*) from kds_queue where tenant_id=4 and received_at>=current_date-interval '6 days' group by 1 order by 1`) };
fs.writeFileSync(`${d.net}/AUD-R4-KDS.json`, JSON.stringify({ run: 'RUN-004', before, after, visible: vis, calls }, null, 1));
console.log(JSON.stringify(after, null, 1).slice(0, 1200)); console.log(calls.filter((c) => /analytics|stats/.test(c.u)).slice(-3).map((c) => c.u + ' ' + c.s + ' ' + c.b.slice(0, 500)).join('\n'));
await ctx.close();
