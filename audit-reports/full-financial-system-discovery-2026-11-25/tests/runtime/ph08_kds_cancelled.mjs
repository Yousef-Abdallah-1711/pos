// RUN-003 / PH-08 / AUD-KDS-001: do cancelled orders still show on the KDS board? HEADED Chrome, read-only.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, sql, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08');
const { ctx, page } = await openChrome();
const cancelled = sql("select o.id||'|'||o.order_number||'|'||k.status from orders o join kds_queue k on k.order_id=o.id where o.tenant_id=4 and o.status='cancelled' order by o.id").split('\n');
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/kds`);
await page.waitForTimeout(8000);
const text = await page.locator('body').innerText();
const shown = cancelled.map((c) => { const [id, num, ks] = c.split('|'); const short = num.slice(-8); return { orderId: id, orderNumber: num, kdsStatus: ks, visibleOnBoard: text.includes(num) || text.includes(num.slice(0, 22)) || text.includes(short) }; });
await page.screenshot({ path: `${d.shot}/AUD-KDS-001-board.png`, fullPage: true, timeout: 40000 }).catch(() => {});
fs.writeFileSync(`${d.net}/AUD-KDS-001.json`, JSON.stringify({ run: 'RUN-003', test: 'AUD-KDS-001', cancelledOrdersAndKdsRows: shown, boardTextSample: text.slice(0, 1800) }, null, 1));
console.log(JSON.stringify(shown, null, 1)); console.log(text.slice(0, 700));
await ctx.close();
