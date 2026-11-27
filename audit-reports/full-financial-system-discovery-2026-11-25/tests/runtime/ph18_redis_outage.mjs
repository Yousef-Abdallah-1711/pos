// RUN-003 / PH-18 / REL (F-018 revalidation): stop ONLY the isolated Redis, measure authenticated request and /health behaviour, then restore Redis.
import fs from 'fs';
import { execSync } from 'child_process';
import { API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const login = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = login.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const timed = async (p, withCookie, ms) => { const s = performance.now(); try { const r = await fetch(API + p, { headers: { Origin: O, ...(withCookie ? { Cookie: ck } : {}) }, signal: AbortSignal.timeout(ms) }); return { status: r.status, ms: Math.round(performance.now() - s) }; } catch (e) { return { error: String(e.name), ms: Math.round(performance.now() - s) }; } };
const out = { before: { auth: await timed('/products', true, 15000), health: await timed('/health', false, 15000) } };
execSync('docker stop pos-p2a-ws05-e2e-redis-1', { stdio: 'ignore' });
await new Promise((x) => setTimeout(x, 3000));
out.during = { auth: await timed('/products', true, 30000), health: await timed('/health', false, 15000), publicMenu: await timed('/public/restaurants/audit-c-diag/products', false, 15000), loginAttempt: await (async () => { const s = performance.now(); try { const r = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }), signal: AbortSignal.timeout(30000) }); return { status: r.status, ms: Math.round(performance.now() - s) }; } catch (e) { return { error: String(e.name), ms: Math.round(performance.now() - s) }; } })() };
execSync('docker start pos-p2a-ws05-e2e-redis-1', { stdio: 'ignore' });
await new Promise((x) => setTimeout(x, 8000));
out.after = { auth: await timed('/products', true, 15000), health: await timed('/health', false, 15000) };
fs.writeFileSync(`${d.net}/AUD-REL-REDIS.json`, JSON.stringify({ run: 'RUN-003', out }, null, 1));
console.log(JSON.stringify(out));
