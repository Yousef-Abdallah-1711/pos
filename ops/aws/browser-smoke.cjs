const { chromium } = require('playwright');
const fs = require('fs');
const rows=[];
const gateway=JSON.parse(fs.readFileSync('/fixtures/gateway.json','utf8'));
const gatewayHeader='Basic '+Buffer.from(gateway.username+':'+gateway.password).toString('base64');
const frontHosts=new Set(['pos.rakmyat.online','pos.app.rakmyat.online','admin.pos.rakmyat.online','restaurant-1.pos.rakmyat.online']);
async function protect(context){await context.route('**/*',async route=>{const u=new URL(route.request().url());const headers={...route.request().headers()};if(u.protocol==='https:'&&frontHosts.has(u.hostname))headers.authorization=gatewayHeader;else if(headers.authorization===gatewayHeader)delete headers.authorization;await route.continue({headers});});}

const anonymous=[
 ['landing-en','https://pos.rakmyat.online/en'],
 ['landing-ar','https://pos.rakmyat.online/ar'],
 ['login','https://pos.rakmyat.online/en/auth/login'],
 ['register','https://pos.rakmyat.online/en/auth/register'],
 ['dashboard-anonymous','https://pos.app.rakmyat.online/en/dashboard'],
 ['super-admin-login','https://admin.pos.rakmyat.online/login'],
 ['storefront-en','https://restaurant-1.pos.rakmyat.online/en'],
 ['storefront-ar','https://restaurant-1.pos.rakmyat.online/ar'],
 ['storefront-menu','https://restaurant-1.pos.rakmyat.online/en/menu'],
];
async function snapshot(page,label,width,url) {
 const errors=[];const onError=e=>errors.push(e.message);page.on('pageerror',onError);
 try {
  const response=await page.goto(url,{waitUntil:'domcontentloaded',timeout:45000});
  await page.waitForTimeout(1200);
  const layout=await page.evaluate(()=>({lang:document.documentElement.lang,dir:document.documentElement.dir,overflow:document.documentElement.scrollWidth>innerWidth+2,applicationError:/Application error|Internal Server Error/.test(document.body.innerText)}));
  await page.screenshot({path:`/evidence/${label}-${width}.png`,fullPage:true});
  rows.push({label,width,url:page.url(),http:response.status(),layout,errors,status:response.status()<400&&!layout.applicationError&&!errors.length?'PASS':'FAIL'});
 } catch(e) { rows.push({label,width,status:'FAIL',error:e.message}); }
 page.off('pageerror',onError);
}
async function loginJourney(browser,admin=false) {
 const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width:1440,height:900}});
 await protect(context);
 const page=await context.newPage();const label=admin?'admin-browser-login':'tenant-browser-login';
 const fixture=JSON.parse(fs.readFileSync(admin?'/fixtures/admin.json':'/fixtures/users.json','utf8'));
 const user=admin?fixture:fixture[0];
 try {
  await page.goto(admin?'https://admin.pos.rakmyat.online/login':'https://pos.rakmyat.online/en/auth/login',{waitUntil:'domcontentloaded',timeout:45000});
  // Server-rendered inputs exist before React attaches their controlled handlers.
  await page.waitForFunction(()=>{const input=document.querySelector('input[type=email]');return input&&Object.keys(input).some(k=>k.startsWith('__reactProps')&&typeof input[k]?.onChange==='function');},{},{timeout:30000});
  await page.locator('input[type=email]').fill(user.email);
  await page.locator('input[type=password]').fill(user.password);
  await Promise.all([page.waitForURL(admin?'**/dashboard':'**/en/dashboard**',{timeout:45000}),page.locator('button[type=submit]').click()]);
  await page.waitForTimeout(2000);
  const session=await page.evaluate(async(admin)=>{
   const r=await fetch('https://api.pos.rakmyat.online/api/v1/'+(admin?'super-admin/':'')+'auth/session',{credentials:'include'});return r.status;
  },admin);
  const cookies=(await context.cookies('https://api.pos.rakmyat.online')).map(({name,secure,httpOnly,sameSite,domain})=>({name,secure,httpOnly,sameSite,domain}));
  const jwtInStorage=await page.evaluate(()=>Object.values(localStorage).some(v=>/eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+/.test(v)));
  rows.push({label,url:page.url(),session,cookies,jwtInStorage,status:session===200&&!jwtInStorage&&cookies.some(c=>c.secure&&c.httpOnly)?'PASS':'FAIL'});
  await page.screenshot({path:`/evidence/${label}.png`,fullPage:true});
  if(!admin) for(const route of ['dashboard','dashboard/pos','dashboard/kds']) await snapshot(page,'authenticated-'+route.replaceAll('/','-'),1440,'https://pos.app.rakmyat.online/en/'+route);
  else await snapshot(page,'admin-authenticated-dashboard',1440,'https://admin.pos.rakmyat.online/dashboard');
 } catch(e) {
  rows.push({label,status:'FAIL',error:String(e.message).replaceAll(user.password,'[redacted]'),url:page.url()});
  await page.screenshot({path:`/evidence/${label}-failure.png`,fullPage:true}).catch(()=>{});
 }
 await context.close();
}
(async()=>{
 const browser=await chromium.launch({headless:true,args:['--disable-dev-shm-usage']});
 for(const width of [1440,390]) {
  const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width,height:900}});const page=await context.newPage();
  await protect(context);
  for(const [label,url] of anonymous) await snapshot(page,label,width,url);
  await context.close();
 }
 await loginJourney(browser);await loginJourney(browser,true);
 await browser.close();
 fs.writeFileSync('/evidence/results.json',JSON.stringify(rows,null,2));
 for(const row of rows)console.log(row.status+' '+row.label+(row.width?' '+row.width:''));
 process.exitCode=rows.some(r=>r.status==='FAIL')?1:0;
})().catch(e=>{console.error(e.message);process.exit(1)});
