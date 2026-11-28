// usage: node generic_create.mjs <tag> <email> <path> <openButtonText> <submitButtonText> <id>
// Opens a create form, fills every visible control with synthetic values by type, submits, records network writes + UI text.
import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, path, openBtn, submitBtn, id] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/' + path, { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
await p.click(`button:has-text("${openBtn}")`); await p.waitForTimeout(1000);
const form = p.locator('form:has(button:has-text("' + submitBtn + '"))').last();
const controls = await form.locator('input:visible, select:visible, textarea:visible').all();
const filled = [];
let n = 0;
for (const c of controls) {
  const tag2 = await c.evaluate(e => e.tagName.toLowerCase()); const type = await c.evaluate(e => e.type || '');
  n++;
  try {
    if (tag2 === 'select') { const v = await c.locator('option').evaluateAll(os => (os.find(o => o.value && !/^select/i.test(o.textContent.trim())) || {}).value); if (v) await c.selectOption(v); filled.push('select'); }
    else if (type === 'checkbox' || type === 'radio' || type === 'file' || type === 'hidden') continue;
    else if (type === 'number') { await c.fill(String(n % 2 ? 4 : 1)); filled.push('number'); }
    else if (type === 'date') { await c.fill('2026-11-26'); filled.push('date'); }
    else if (type === 'email') { await c.fill(`audit${Date.now() % 100000}@example.test`); filled.push('email'); }
    else if (type === 'tel') { await c.fill('+201000000999'); filled.push('tel'); }
    else if (type === 'password') { await c.fill('AuditPass#2026A'); filled.push('password'); }
    else if (type === 'time') { await c.fill('09:00'); filled.push('time'); }
    else { await c.fill(`Audit ${id} ${n}`); filled.push('text'); }
  } catch (e) { filled.push('ERR:' + e.message.slice(0, 40)); }
}
await p.screenshot({ path: `${EV}/screenshots/${id}-form.png` });
await form.locator(`button:has-text("${submitBtn}")`).click(); await p.waitForTimeout(3500);
const txt = (await p.locator('body').innerText()).split(String.fromCharCode(10)).join(' | ');
console.log(id, 'controls filled:', filled.length, filled.join(','));
console.log('writes:', JSON.stringify(net.filter(x => x.m !== 'GET').map(x => `${x.m} ${x.u} ${x.s} resp=${x.body.slice(0, 220)}`), null, 1));
console.log('ui tail:', txt.slice(-260));
await b.close();
