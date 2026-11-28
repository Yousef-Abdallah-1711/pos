import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, nameEn, price, cost, qty] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/en/dashboard/products/new', { waitUntil: 'networkidle' }); await p.waitForTimeout(2500);
await p.fill('[name=name_en]', nameEn); await p.fill('[name=price]', price); await p.fill('[name=cost]', cost);
// Track inventory toggle: checkbox following the "Track Inventory" label
const tr = p.locator('text=Track Inventory').first();
const cb = p.locator('xpath=//*[contains(text(),"Track Inventory")]/ancestor::div[1]//input[@type="checkbox"] | //label[contains(.,"Track Inventory")]//input').first();
try { if (await cb.count()) { if (!(await cb.isChecked())) await cb.check({ force: true }); } } catch (e) { console.log('track toggle', e.message.slice(0, 80)); }
await p.waitForTimeout(500);
await p.fill('[name=quantity_in_stock]', qty).catch(e => console.log('qty', e.message.slice(0, 80)));
await p.screenshot({ path: `${EV}/screenshots/prod-${tag}-${nameEn.replace(/\W+/g, '_')}-filled.png` });
await p.click('button:has-text("Create Product")'); await p.waitForTimeout(4000);
console.log('URL', p.url());
console.log((await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 300));
console.log(net.filter(n => n.m !== 'GET').map(n => `${n.m} ${n.u} ${n.s} req=${n.req} resp=${n.body.slice(0, 250)}`).join('\n'));
await b.close();
