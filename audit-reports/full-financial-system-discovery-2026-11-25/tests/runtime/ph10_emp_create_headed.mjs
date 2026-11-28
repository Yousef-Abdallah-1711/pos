// RUN-003 / PH-10 / AUD-HR-001: reproduce F-014 (employee create via UI fails) in HEADED Chrome with before/after SQL.
// Target: isolated stack, tenant 4 (Audit C), owner role. Expected side effect if product is correct: 1 employee row.
// If product defect reproduces: 0 employee rows (and we record any orphan rows e.g. users).
import { chromium } from 'file:///C:/Users/Yousef/Desktop/Projects/POSS/node_modules/playwright/index.mjs';
import { execSync } from 'child_process';
import fs from 'fs';
import { ensureLoggedIn } from './chrome_lib.mjs';

const ROOT = 'C:/Users/Yousef/Desktop/Projects/POSS/audit-reports/full-financial-system-discovery-2026-11-25';
const SHOT = `${ROOT}/evidence/screenshots/RUN-003/PH-10`;
const NET = `${ROOT}/evidence/network/RUN-003/PH-10`;
const DB = `${ROOT}/evidence/database/RUN-003/PH-10`;
const PROFILE = process.env.AUDIT_CHROME_PROFILE;
if (!PROFILE) throw new Error('AUDIT_CHROME_PROFILE required');
[SHOT, NET, DB].forEach((d) => fs.mkdirSync(d, { recursive: true }));

const sql = (q) => execSync(`docker exec pos-p2a-ws05-e2e-postgres-1 psql -U postgres -d pos_phase2a_e2e -Atc "${q}"`, { encoding: 'utf8' }).trim();
const snapshot = () => sql("select 'employees_t4',count(*) from employees where tenant_id=4 union all select 'users_t4',count(*) from users where tenant_id=4 union all select 'journals_t4',count(*) from journal_entries where tenant_id=4");
const before = snapshot();

const ctx = await chromium.launchPersistentContext(PROFILE, { channel: 'chrome', headless: false, viewport: { width: 1440, height: 900 } });
const page = ctx.pages()[0] ?? (await ctx.newPage());
const calls = [];
page.on('response', async (r) => {
  if (r.url().includes(':18082') && r.request().method() !== 'GET' && r.request().method() !== 'OPTIONS') {
    let body = '';
    try { body = (await r.text()).slice(0, 600); } catch { /* body unavailable */ }
    calls.push({ method: r.request().method(), url: r.url().replace('http://localhost:18082/api/v1', ''), status: r.status(), request: (r.request().postData() || '').slice(0, 1500), response: body });
  }
});

await ensureLoggedIn(page, 'audit.c.owner@example.test', 'http://localhost:13112/en/dashboard/hr/employees');
await page.waitForSelector('button:has-text("Add Employee")', { timeout: 120000 });
await page.waitForTimeout(2000);
await page.click('button:has-text("Add Employee")');
await page.waitForTimeout(1000);
await page.fill('[name=first_name]', 'Headed');
await page.fill('[name=last_name]', 'Repro');
await page.fill('[name=email]', 'headed.repro@example.test');
await page.fill('[name=phone]', '+201000000199');
await page.click('button:has-text("Employment")');
await page.fill('[name=hire_date]', '2026-11-01');
await page.fill('[name=position]', 'Chef');
await page.fill('[name=department]', 'Kitchen');
await page.click('button:has-text("Compensation")');
await page.fill('[name=monthly_salary]', '9000');
await page.screenshot({ path: `${SHOT}/AUD-HR-001-before-submit.png` });
await page.click('button:has-text("Create Employee")');
await page.waitForTimeout(4000);
await page.screenshot({ path: `${SHOT}/AUD-HR-001-after-submit.png` });
const visibleText = (await page.locator('body').innerText()).slice(0, 1500);

const after = snapshot();
fs.writeFileSync(`${NET}/AUD-HR-001.json`, JSON.stringify({ run: 'RUN-003', test: 'AUD-HR-001', finding: 'F-014', tenant: 4, role: 'owner', calls }, null, 2));
fs.writeFileSync(`${DB}/AUD-HR-001-before-after.txt`, `BEFORE\n${before}\n\nAFTER\n${after}\n\nVISIBLE_TEXT_AFTER_SUBMIT (first 1500 chars)\n${visibleText}\n`);
console.log(JSON.stringify({ before: before.split('\n'), after: after.split('\n'), calls: calls.map((c) => `${c.method} ${c.url} ${c.status}`) }, null, 2));
await ctx.close();
