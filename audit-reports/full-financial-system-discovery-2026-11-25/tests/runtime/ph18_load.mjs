// RUN-003 / PH-18 / REL-011-lite: incremental load (10, 25, 50 concurrent virtual users, 15 s each) on read endpoints with one authenticated session; isolated stack, RBAC enforce. Not a production capacity claim.
import fs from 'fs';
import { API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const EPS = ['/products', '/admin/orders?limit=20', '/accounting/trial-balance', '/reports/financial/profit-loss?date_from=2026-10-28&date_to=2026-11-26', '/reports/sales?date_from=2026-10-28&date_to=2026-11-26&granularity=day', '/kds/queue'];
const out = {};
for (const vu of [10, 25, 50]) {
  const lat = []; let err = 0, n = 0; const stop = Date.now() + 15000;
  await Promise.all(Array.from({ length: vu }, async (_, w) => { let i = w; while (Date.now() < stop) { const e = EPS[i++ % EPS.length]; const s = performance.now(); try { const x = await fetch(API + e, { headers: { Origin: O, Cookie: ck }, signal: AbortSignal.timeout(20000) }); await x.arrayBuffer(); if (x.status >= 400) err++; } catch { err++; } lat.push(performance.now() - s); n++; } }));
  lat.sort((a, b) => a - b); out[vu] = { requests: n, rps: Math.round(n / 15), errors: err, errorRate: +(err / n * 100).toFixed(2), p50: Math.round(lat[Math.floor(lat.length * 0.5)]), p95: Math.round(lat[Math.floor(lat.length * 0.95)]), max: Math.round(lat[lat.length - 1]) };
  console.log(vu, 'VU', JSON.stringify(out[vu]));
}
fs.writeFileSync(`${d.net}/AUD-LOAD.json`, JSON.stringify({ run: 'RUN-003', note: 'read endpoints only; concurrent headed-Chrome sessions were running on the same host', out }, null, 1));
