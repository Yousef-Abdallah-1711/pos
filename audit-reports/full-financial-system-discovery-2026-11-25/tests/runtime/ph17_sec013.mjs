// RUN-004 / PH-17: SEC-013 staging vs main isolated API (ENV=phase2a-e2e, RBAC_ENFORCEMENT=enforce) vs a second API container started with ENV=staging and RBAC_ENFORCEMENT unset. API-DIAGNOSTIC, read-only requests.
import fs from 'fs';
import { PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-17'); const O = 'http://localhost:13112';
const TARGETS = { main_enforce: 'http://127.0.0.1:18082/api/v1', staging_default: 'http://127.0.0.1:18085/api/v1' };
const out = {};
for (const [name, base] of Object.entries(TARGETS)) {
  const o = (out[name] = { base, checks: {} });
  const login = await fetch(base + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'cashier.c@example.test', password: PASSWORD }) });
  const setc = login.headers.getSetCookie();
  o.login_status = login.status;
  o.cookie_flags = setc.map((c) => c.replace(/=[^;]*/, '=<redacted>'));
  const ck = setc.map((c) => c.split(';')[0]).join('; ');
  const g = async (p, extra = {}) => { const r = await fetch(base + p, { headers: { Origin: O, Cookie: ck, ...extra } }); return r.status; };
  for (const p of ['/accounting/journal-entries', '/accounting/accounts', '/hr/salaries', '/hr/employees', '/reports/financial/profit-loss?date_from=2026-11-01&date_to=2026-11-30', '/reports/financial/balance-sheet?as_of=2026-11-30', '/users', '/settings/restaurant', '/products?page=1&page_size=1']) o.checks[p] = await g(p);
  o.metrics_no_token = (await fetch(base + '/health/metrics')).status;
  o.pos_write_as_cashier = (await fetch(base + '/accounting/accounts', { method: 'POST', headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: JSON.stringify({}) })).status;
  console.log(name, JSON.stringify(o));
}
fs.writeFileSync(`${d.net}/AUD-R4-SEC013.json`, JSON.stringify({ run: 'RUN-004', kind: 'API-DIAGNOSTIC', note: 'staging_default = ENV=staging, RBAC_ENFORCEMENT unset, same DB/Redis, isolated; main_enforce = isolated main API', out }, null, 1));
