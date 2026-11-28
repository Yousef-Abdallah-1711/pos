// RUN-005 / PH-18 (F-018 root-cause confirmation): with the isolated Redis stopped, measures the TIME TO FAILURE (not just "hangs") of an authenticated request, login and /health with a 200 s client timeout, then restarts Redis and measures recovery.
// Isolated stack only (container pos-p2a-ws05-e2e-redis-1). API-DIAGNOSTIC.
import fs from 'fs';
import { execSync } from 'child_process';
import { API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112'; const RED = 'pos-p2a-ws05-e2e-redis-1';
const t = async (label, fn) => { const s = performance.now(); let r; try { r = await fn(); } catch (e) { r = { error: e.name + ': ' + String(e.message).slice(0, 60) }; } const ms = Math.round(performance.now() - s); const o = { label, ms, ...r }; console.log(JSON.stringify(o)); return o; };
const req = (m, p, ck, b, to = 200000) => async () => { const x = await fetch(API + p, { method: m, headers: { Origin: O, ...(ck ? { Cookie: ck } : {}), 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined, signal: AbortSignal.timeout(to) }); const body = (await x.text()).slice(0, 100); return { status: x.status, retryAfter: x.headers.get('retry-after'), body, cookie: x.headers.getSetCookie?.().map((c) => c.split(';')[0]).join('; ') }; };
const out = [];
const login = await t('login (Redis up)', req('POST', '/auth/login', null, { email: 'audit.c.owner@example.test', password: PASSWORD })); out.push({ ...login, cookie: undefined }); const ck = login.cookie;
out.push(await t('GET /products (Redis up)', req('GET', '/products?page=1&page_size=1', ck)));
execSync(`docker stop ${RED}`); out.push({ label: 'REDIS STOPPED', at: new Date().toISOString() });
out.push(await t('GET /health (Redis down)', req('GET', '/health', null, null, 20000)));
out.push(await t('GET /products with valid cookie (Redis down)', req('GET', '/products?page=1&page_size=1', ck)));
out.push(await t('POST /auth/login (Redis down)', req('POST', '/auth/login', null, { email: 'audit.c.owner@example.test', password: PASSWORD })));
out.push(await t('GET public menu (Redis down)', req('GET', '/public/restaurants/audit-c-diag/theme', null, null, 60000)));
execSync(`docker start ${RED}`); out.push({ label: 'REDIS STARTED', at: new Date().toISOString() });
await new Promise((r) => setTimeout(r, 8000));
out.push(await t('GET /products with the OLD cookie after Redis restart', req('GET', '/products?page=1&page_size=1', ck, null, 60000)));
const l2 = await t('POST /auth/login after Redis restart', req('POST', '/auth/login', null, { email: 'audit.c.owner@example.test', password: PASSWORD }, 60000)); out.push({ ...l2, cookie: undefined });
out.push(await t('GET /products with a NEW cookie', req('GET', '/products?page=1&page_size=1', l2.cookie, null, 60000)));
fs.writeFileSync(`${d.net}/AUD-R5-F018-REDIS-TIMING.json`, JSON.stringify({ run: 'RUN-005', finding: 'F-018', out: out.map((o) => ({ ...o, cookie: undefined })) }, null, 1));
