// RUN-003 / PH-12: batch capture of remaining report endpoints (API-DIAGNOSTIC, tenant 4 owner session via real login). UI/Chrome rendering is captured separately.
import fs from 'fs';
import { API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
let cookie = '';
const get = async (path, method = 'GET', body) => { const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 300); } return { status: r.status, body: j }; };
await get('/auth/login', 'POST', { email: 'audit.c.owner@example.test', password: PASSWORD });
const Q = 'date_from=2026-10-28&date_to=2026-11-26';
const eps = { pnl: `/reports/financial/profit-loss?${Q}&compare_date_from=2026-09-28&compare_date_to=2026-10-27`, tb: '/accounting/trial-balance', bs: '/reports/financial/balance-sheet?as_of=2026-11-26', ledger70: '/accounting/ledger?account_id=70&limit=500', ledger71: '/accounting/ledger?account_id=71&limit=500', turn: `/reports/inventory?${Q}&granularity=day`, ovr: `/reports/overview?${Q}`, hrcost2: '/reports/hr/cost-analysis?date_from=2026-11-01&date_to=2026-11-30', sales: `/reports/sales?${Q}&granularity=day`, prof: `/reports/profitability?${Q}&granularity=day`, taxperiods: '/accounting/tax/periods', orders: '/admin/orders?limit=100', products: '/products', categories: '/categories', employees: '/hr/employees', drivers: '/admin/drivers', tables: '/tables/stats', kdsstats: '/kds/stats', cf: `/reports/financial/cash-flow?${Q}`, hrturn: `/reports/hr/turnover?${Q}` };
const out = {};
for (const [k, p] of Object.entries(eps)) { const r = await get(p); out[k] = { path: p, status: r.status, body: r.body }; console.log(k, r.status, JSON.stringify(r.body).slice(0, 260)); }
fs.writeFileSync(`${d.net}/AUD-MET-BATCH2.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, kind: 'API-DIAGNOSTIC', out }, null, 1));
