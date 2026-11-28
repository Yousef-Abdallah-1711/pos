// RUN-004 / PH-12: MET-HRR-002 (turnover) and MET-HRR-007 (social-insurance CSV). API-DIAGNOSTIC against isolated tenant 4 (diagnostic fixture) + SQL oracle. Rendering of the HR report page is captured by ph12_r4_metrics2.mjs (headed Chrome).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b, raw) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = null; } return { status: x.status, body: j, text: raw ? t : t.slice(0, 300) }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 420)); };
const from = '2026-10-29', to = '2026-11-27';
// ---- HRR-007 social insurance CSV
let r = await A('GET', '/reports/hr/social-insurance?year=2026&month=11', null, true); fs.writeFileSync(`${d.net}/AUD-R4-HRR007-before.csv`, r.text);
rec('HRR007-01', 'social insurance CSV before', { http: r.status, lines: r.text.trim().split('\n').length, head: r.text.split('\n').slice(0, 3).join(' / ') });
rec('HRR007-02', 'SQL oracle: salary rows for 2026-11', { salaries: sql(`select count(*) from salaries where tenant_id=4 and year=2026 and month=11`), active_policies_emp4: sql(`select count(*) from insurance_tracking where tenant_id=4 and employee_id=4 and status='active'`) });
r = await A('POST', '/hr/insurance', { employee_id: 4, policy_type: 'life', provider_name: 'Audit Insurer 2', policy_number: 'POL-R4-2', coverage_start: '2026-11-01', coverage_end: '2027-10-31', monthly_premium: 100 }); rec('HRR007-03', 'add second active policy for the same employee', { http: r.status, resp: r.text });
r = await A('GET', '/reports/hr/social-insurance?year=2026&month=11', null, true); fs.writeFileSync(`${d.net}/AUD-R4-HRR007-after-second-policy.csv`, r.text);
rec('HRR007-04', 'CSV after second active policy (oracle: one row per employee salary row)', { http: r.status, lines: r.text.trim().split('\n').length, sql_salary_rows: sql(`select count(*) from salaries where tenant_id=4 and year=2026 and month=11`), body: r.text.split('\n').slice(0, 4).join(' / ') });
// ---- HRR-002 turnover
const turn = async () => (await A('GET', `/reports/hr/turnover?date_from=${from}&date_to=${to}`)).body?.data?.summary;
rec('HRR002-01', 'turnover before', { summary: JSON.stringify(await turn()), employees: sql(`select count(*) from employees where tenant_id=4`) });
const code = 'EMP-R4-' + Date.now();
r = await A('POST', '/hr/employees', { employee_code: code, first_name: 'Audit', last_name: 'Leaver', email: `leaver.${Date.now()}@example.test`, hire_date: '2026-10-01', employment_type: 'full_time', employment_status: 'active', position: 'Cook', department: 'Kitchen', base_salary: 4000, salary_currency: 'EGP', payment_frequency: 'monthly', working_hours_per_week: 40, shift_type: 'day' });
const eid = r.body?.data?.id ?? r.body?.id; rec('HRR002-02', 'create employee with backend-native keys (F-014 workaround, diagnostic)', { http: r.status, id: eid, resp: r.text });
if (eid) {
  r = await A('POST', '/hr/resignations', { employee_id: eid, resignation_date: '2026-11-20', last_working_date: '2026-11-25', reason: 'RUN-004 turnover probe' }); const rid = r.body?.data?.id ?? r.body?.id; rec('HRR002-03', 'create resignation', { http: r.status, id: rid });
  r = await A('POST', `/hr/resignations/${rid}/acknowledge`, {}); rec('HRR002-04', 'acknowledge', { http: r.status });
  r = await A('POST', `/hr/resignations/${rid}/complete`, { settlement_amount: 0 }); rec('HRR002-05', 'complete', { http: r.status, resp: r.text });
  rec('HRR002-06', 'employee row after completed resignation (oracle: terminated + termination_date 2026-11-25)', { row: sql(`select employment_status,is_active,termination_date from employees where id=${eid}`) });
  rec('HRR002-07', 'turnover after completed resignation (oracle: departures 1)', { summary: JSON.stringify(await turn()) });
  r = await A('DELETE', `/hr/employees/${eid}`); rec('HRR002-08', 'DELETE employee (the only path that sets termination_date)', { http: r.status });
  const after = await turn(); rec('HRR002-09', 'turnover after delete', { summary: JSON.stringify(after), row: sql(`select employment_status,termination_date from employees where id=${eid}`), headcount_start_sql: sql(`select count(*) from employees where tenant_id=4 and hire_date<='${from}' and (termination_date is null or termination_date>'${from}')`), headcount_end_sql: sql(`select count(*) from employees where tenant_id=4 and hire_date<='${to}' and (termination_date is null or termination_date>'${to}')`), departures_sql: sql(`select count(*) from employees where tenant_id=4 and termination_date between '${from}' and '${to}'`) });
}
fs.writeFileSync(`${d.db}/AUD-R4-HR.json`, JSON.stringify({ run: 'RUN-004', kind: 'API-DIAGNOSTIC', tenant: 4, log }, null, 1));
