// RUN-003 / PH-12: batch capture of remaining report endpoints (API-DIAGNOSTIC, tenant 4 owner session via real login). UI/Chrome rendering is captured separately.
import fs from 'fs';
import { API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12');
let cookie = '';
const get = async (path, method = 'GET', body) => { const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 300); } return { status: r.status, body: j }; };
await get('/auth/login', 'POST', { email: 'audit.c.owner@example.test', password: PASSWORD });
const Q = 'date_from=2026-10-28&date_to=2026-11-26';
const eps = { vat: `/accounting/tax/vat-report?${Q}`, invval: '/reports/financial/inventory-valuation', kpis: `/reports/financial/kpis?${Q}`, menu: `/reports/financial/menu-profitability?${Q}`, be: `/reports/financial/break-even?${Q}`, pcf: '/reports/financial/predictive-cash-flow', hrcost: `/reports/hr/cost-analysis?${Q}`, hrhead: `/reports/hr/headcount?${Q}`, hrturn: `/reports/hr/turnover?${Q}`, hrabs: `/reports/hr/absenteeism?${Q}`, hrot: `/reports/hr/overtime?${Q}`, hrleave: `/reports/hr/leave-balances`, hrtrain: `/reports/hr/training-compliance`, ordstats: '/admin/orders/stats', kds: '/kds/analytics', tables: '/tables/stats', bs: `/reports/financial/balance-sheet?as_of=2026-11-26`, ledger: '/accounting/ledger?limit=500', tb: '/accounting/trial-balance' };
const out = {};
for (const [k, p] of Object.entries(eps)) { const r = await get(p); out[k] = { path: p, status: r.status, body: r.body }; console.log(k, r.status, JSON.stringify(r.body).slice(0, 260)); }
fs.writeFileSync(`${d.net}/AUD-MET-BATCH.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, kind: 'API-DIAGNOSTIC', out }, null, 1));
