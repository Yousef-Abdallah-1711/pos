// RUN-003 / PH-10 part 2: leave approval over allowance, leave validation, clock-in API, attendance hours. API-DIAGNOSTIC tenant 4.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-10'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 160); } return { status: x.status, body: j }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 400)); };
const lid = Number(sql("select id from leaves where employee_id=4 order by id limit 1") || 0);
rec('LEV-ROW', 'existing 30-day leave', { row: sql(`select id||'|'||status||'|'||total_days||'|'||leave_type from leaves where id=${lid}`) });
let r = await A('POST', `/hr/leaves/${lid}/approve`, {}); rec('LEV-02', 'approve 30-day annual leave over allowance 21', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200), status: sql(`select status from leaves where id=${lid}`) });
const bal = await A('GET', '/reports/hr/leave-balances'); rec('LEV-BAL', 'leave balance after approval', { item: JSON.stringify(bal.body?.data?.items?.[0]) });
r = await A('POST', '/hr/leaves', { employee_id: 4, start_date: '2027-01-20', end_date: '2027-01-10', total_days: 1, leave_type: 'annual', leave_category: 'paid', reason: 'end before start' }); rec('LEV-03b', 'leave end before start (no overlap)', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200) });
r = await A('POST', '/hr/leaves', { employee_id: 4, start_date: '2027-01-10', end_date: '2027-01-12', total_days: 40, leave_type: 'annual', leave_category: 'paid', reason: 'total_days mismatch' }); rec('LEV-05', 'total_days 40 for a 3-day range', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200), stored: sql("select total_days from leaves where reason='total_days mismatch'") });
r = await A('POST', '/hr/attendance/clock-in', { employee_id: 4 }); rec('ATT-05', 'clock-in API', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200) });
r = await A('POST', '/hr/attendance', { employee_id: 4, attendance_date: '2026-11-23', clock_in: '09:00', clock_out: '19:30', status: 'present' }); rec('ATT-06', 'attendance 09:00-19:30 hours computed?', { http: r.status, row: sql("select clock_in::text||'|'||clock_out::text||'|'||regular_hours||'|'||overtime_hours from attendance where employee_id=4 and attendance_date='2026-11-23'") });
const ot = await A('GET', '/reports/hr/overtime?date_from=2026-11-01&date_to=2026-11-30'); rec('ATT-OT', 'overtime report for November', { summary: JSON.stringify(ot.body?.data?.summary) });
fs.writeFileSync(`${d.db}/AUD-HR-MATRIX2.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
