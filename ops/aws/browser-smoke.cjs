const { chromium } = require('playwright');
const { assertCompleteTargetConfig, apiBase, guardContext } = require('./acceptance-target.cjs');
assertCompleteTargetConfig();
const api = apiBase();
const fs = require('fs');
const evidence=process.env.SERVIO_BROWSER_EVIDENCE||'/evidence';
const fixtures=process.env.SERVIO_BROWSER_FIXTURES||'/fixtures';
const rows=[];
const anonymous=[
 ['landing-en','https://pos.servio.test/en'],
 ['landing-ar','https://pos.servio.test/ar'],
 ['login','https://pos.servio.test/en/auth/login'],
 ['register','https://pos.servio.test/en/auth/register'],
 ['dashboard-anonymous','https://pos.app.servio.test/en/dashboard'],
 ['super-admin-login','https://admin.pos.servio.test/login'],
 ['storefront-en','https://restaurant-1.pos.servio.test/en'],
 ['storefront-ar','https://restaurant-1.pos.servio.test/ar'],
 ['storefront-menu','https://restaurant-1.pos.servio.test/en/menu'],
];
async function snapshot(page,label,width,url) {
 const errors=[];const onError=e=>errors.push(e.message);page.on('pageerror',onError);
 try {
  await page.setViewportSize({width,height:900});
  const response=await page.goto(url,{waitUntil:'domcontentloaded',timeout:45000});
  await page.waitForTimeout(1200);
  const layout=await page.evaluate(()=>{const root=document.documentElement;const width=innerWidth;const overflow=root.scrollWidth>width+2;const elements=overflow?[...document.querySelectorAll('body *')].map(element=>{const rect=element.getBoundingClientRect();return{tag:element.tagName,id:element.id,className:typeof element.className==='string'?element.className.slice(0,100):'',right:Math.round(rect.right),width:Math.round(rect.width)}}).filter(element=>element.right>width+2&&element.width>0).sort((a,b)=>b.right-a.right).slice(0,8):[];return{lang:root.lang,dir:root.dir,overflow,viewportWidth:width,scrollWidth:root.scrollWidth,overflowElements:elements,applicationError:/Application error|Internal Server Error/.test(document.body.innerText)}});
  await page.screenshot({path:`${evidence}/${label}-${width}.png`,fullPage:true});
  rows.push({label,width,url:page.url(),http:response.status(),layout,errors,status:response.status()<400&&!layout.applicationError&&!errors.length&&!layout.overflow&&(!new URL(url).pathname.startsWith('/ar')||(layout.lang==='ar'&&layout.dir==='rtl'))?'PASS':'FAIL'});
 } catch(e) { rows.push({label,width,status:'FAIL',error:e.message}); }
 page.off('pageerror',onError);
}
async function loginJourney(browser,admin=false) {
 const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width:1440,height:900}});
 await guardContext(context);

 const page=await context.newPage();const label=admin?'admin-browser-login':'tenant-browser-login';
 const fixture=JSON.parse(fs.readFileSync(admin?fixtures+'/admin.json':fixtures+'/users.json','utf8'));
 const user=admin?fixture:fixture[0];
 try {
  await page.goto(admin?'https://admin.pos.servio.test/login':'https://pos.servio.test/en/auth/login',{waitUntil:'domcontentloaded',timeout:45000});
  // Server-rendered inputs exist before React attaches their controlled handlers.
  await page.waitForFunction(()=>{const input=document.querySelector('input[type=email]');return input&&Object.keys(input).some(k=>k.startsWith('__reactProps')&&typeof input[k]?.onChange==='function');},{},{timeout:30000});
  await page.locator('input[type=email]').fill(user.email);
  await page.locator('input[type=password]').fill(user.password);
  await Promise.all([page.waitForURL(admin?'**/dashboard':'**/en/dashboard**',{timeout:45000}),page.locator('button[type=submit]').click()]);
  await page.waitForTimeout(2000);
  const session=await page.evaluate(async({admin,api})=>{
   const r=await fetch(api+'/'+(admin?'super-admin/':'')+'auth/session',{credentials:'include'});return r.status;
  },{admin,api});
  const cookies=(await context.cookies(new URL(api).origin)).map(({name,secure,httpOnly,sameSite,domain})=>({name,secure,httpOnly,sameSite,domain}));
  const jwtInStorage=await page.evaluate(()=>Object.values(localStorage).some(v=>/eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+/.test(v)));
  rows.push({label,url:page.url(),session,cookies,jwtInStorage,status:session===200&&!jwtInStorage&&cookies.some(c=>c.secure&&c.httpOnly)?'PASS':'FAIL'});
  await page.screenshot({path:`${evidence}/${label}.png`,fullPage:true});
  if(!admin) for(const route of ['dashboard','dashboard/pos','dashboard/kds']) for(const width of [320,375,390,768,1440]) await snapshot(page,'authenticated-'+route.replaceAll('/','-'),width,'https://pos.app.servio.test/en/'+route);
  else for(const width of [320,375,390,768,1440]) await snapshot(page,'admin-authenticated-dashboard',width,'https://admin.pos.servio.test/dashboard');
 } catch(e) {
  rows.push({label,status:'FAIL',error:String(e.message).replaceAll(user.password,'[redacted]'),url:page.url()});
  await page.screenshot({path:`${evidence}/${label}-failure.png`,fullPage:true}).catch(()=>{});
 }
 await context.close();
}
(async()=>{
 const browser=await chromium.launch({headless:true,args:['--disable-dev-shm-usage']});
 for(const width of [1440,768,390,375]) {
  const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width,height:900}});await guardContext(context);const page=await context.newPage();

  for(const [label,url] of anonymous) await snapshot(page,label,width,url);
  await context.close();
 }
 await loginJourney(browser);await loginJourney(browser,true);
 await browser.close();
 fs.writeFileSync(evidence+'/results.json',JSON.stringify(rows,null,2));
 for(const row of rows)console.log(row.status+' '+row.label+(row.width?' '+row.width:''));
 process.exitCode=rows.some(r=>r.status==='FAIL')?1:0;
})().catch(e=>{console.error(e.message);process.exit(1)});
