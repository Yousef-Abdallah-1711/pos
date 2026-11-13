// Full route discovery and read-only browser acceptance. Evidence contains no passwords.
const {chromium}=require('playwright');
const fs=require('fs');
const path=require('path');
const root=path.resolve(__dirname,'../..');
const evidence=path.join(root,'ops/aws/evidence/full-browser');
fs.mkdirSync(evidence,{recursive:true,mode:0o700});
const rows=[];
const dynamicOnly=process.env.SERVIO_DYNAMIC_ONLY==='1';
const run={scope:dynamicOnly?'dynamic-diagnostic':'full',started_utc:new Date().toISOString(),sha:process.env.SERVIO_RELEASE_SHA||null};
const hosts={landing:'https://pos.rakmyat.online',dashboard:'https://pos.app.rakmyat.online','super-admin':'https://admin.pos.rakmyat.online','restaurant-website':'https://restaurant-1.pos.rakmyat.online'};
function discover(directory){
 return fs.readdirSync(directory,{withFileTypes:true}).flatMap(entry=>entry.isDirectory()?discover(path.join(directory,entry.name)):entry.name==='page.tsx'?[path.join(directory,entry.name)]:[]);
}
function routes(app){
 const base=path.join(root,'frontend/apps',app,app==='restaurant-website'?'app':'src/app');
 return [...new Set(discover(base).map(file=>'/'+path.relative(base,file).replace(/\/page.tsx$|^page.tsx$/g,'')))].sort();
}
function resolveFixture(app,raw){
 const file=app==='restaurant-website'?'storefront-qa-fixture.json':'route-qa-fixtures.json';
 const fixture=JSON.parse(fs.readFileSync(path.join(root,'ops/aws/evidence',file)));
 if(!fixture.synthetic)throw new Error('Synthetic route fixtures required');
 if(raw.includes('[tableId]'))return raw.replace('[tableId]',fixture.table_id);
 if(raw.includes('[orderNumber]'))return raw.replace('[orderNumber]',fixture.order_number);
 if(!raw.includes('[id]'))return raw;
 let id;
 if(app==='restaurant-website')id=fixture.order_id;
 else if(app==='super-admin')id=fixture.tenant_id;
 else {
  const resources={'categories':'categories','products':'products','roles':'roles','theme-builder':'themes','orders':'orders','tables':'tables','drivers':'drivers','journal-entries':'journal_entries'};
  const resource=Object.keys(resources).find(name=>raw.split('/').includes(name));
  id=resource&&fixture[resources[resource]];
 }
 if(!Number.isInteger(id)||id<=0)throw new Error('Missing synthetic identifier for '+raw);
 return raw.replace('[id]',id);
}
function safeURL(url){const parsed=new URL(url);return parsed.origin+parsed.pathname;}
async function inspect(page,app,route,width){
 const row={app,route,width,status:'FAIL',console:[],network:[],errors:[]};
 const error=event=>row.errors.push(event.message.slice(0,300));
 const log=message=>{if(message.type()==='error'&&!/Failed to load resource: the server responded with a status of 401/.test(message.text()))row.console.push(message.text().slice(0,300));};
 const response=result=>{if(result.status()>=400&&!(result.status()===401&&/\/auth\/session$/.test(new URL(result.url()).pathname)))row.network.push({url:safeURL(result.url()),status:result.status()});};
 page.on('pageerror',error);page.on('console',log);page.on('response',response);
 try{
  await page.setViewportSize({width,height:900});
  const result=await page.goto(hosts[app]+route,{waitUntil:'domcontentloaded',timeout:45000});
  await page.waitForTimeout(1200);
  row.http=result.status();row.finalURL=safeURL(page.url());
  row.layout=await page.evaluate(()=>({lang:document.documentElement.lang,dir:document.documentElement.dir,overflow:document.documentElement.scrollWidth>innerWidth+2,bodyLength:document.body.innerText.trim().length,applicationError:/Application error|Internal Server Error|Something went wrong|حدث خطأ غير متوقع/.test(document.body.innerText),brokenImages:[...document.images].filter(i=>i.complete&&i.currentSrc&&i.naturalWidth===0).map(i=>({alt:i.alt,src:new URL(i.currentSrc).pathname})),forms:[...document.forms].map(f=>({fields:f.elements.length,invalid:f.querySelectorAll(':invalid').length})),unlabelledFields:[...document.querySelectorAll('input:not([type=hidden]),select,textarea')].filter(i=>!i.labels?.length&&!i.getAttribute('aria-label')&&!i.getAttribute('aria-labelledby')).map(i=>({type:i.type,name:i.name})),headings:document.querySelectorAll('h1').length}));
  await page.keyboard.press('Tab');
  row.keyboardFocus=await page.evaluate(()=>({tag:document.activeElement.tagName,visible:!!document.activeElement.getClientRects().length}));
  const locale=route.startsWith('/ar')?'ar':route.startsWith('/en')?'en':null;
  const localeOK=!locale||(row.layout.lang===locale&&(locale!=='ar'||row.layout.dir==='rtl'));
  row.status=row.http<400&&row.layout.bodyLength>20&&!row.layout.applicationError&&!row.layout.overflow&&!row.errors.length&&!row.console.length&&!row.layout.brokenImages.length&&!row.network.length&&localeOK?'PASS':'FAIL';
  if(row.status==='FAIL'||route==='/ar'||route==='/en'){
   row.screenshot=`${app}-${route.replaceAll('/','_')||'root'}-${width}.png`;
   await page.screenshot({path:path.join(evidence,row.screenshot),fullPage:true});
  }
 }catch(error){row.error=error.message.slice(0,300);}
 page.off('pageerror',error);page.off('console',log);page.off('response',response);
 rows.push(row);fs.writeFileSync(path.join(evidence,'results.json'),JSON.stringify(rows,null,2));
 console.log(`${row.status} ${app} ${route} ${width}`);
}
async function verifyAnonymousCallbacks(browser){
 const context=await browser.newContext();const page=await context.newPage();
 for(const route of ['/auth/callback','/en/auth/callback','/ar/auth/callback'])for(const width of [390]){
  const row={app:'dashboard',route:route+' anonymous rejection',width,status:'FAIL'};
  try{
   await page.setViewportSize({width,height:900});
   await page.goto(hosts.dashboard+route+'?token=invalid-synthetic-token');
   await page.waitForURL(url=>url.origin===hosts.landing&&url.pathname.endsWith('/auth/login'));
   const session=await context.request.get('https://api.pos.rakmyat.online/api/v1/auth/session');
   row.http=session.status();row.status=row.http===401?'PASS':'FAIL';
  }catch(error){row.error=error.message.slice(0,200);}
  rows.push(row);console.log(row.status+' '+row.route+' '+width);
 }
 await context.close();
}
async function login(context,admin=false){
 const file=admin?'browser-admin.json':'browser-fixtures.json';
 const fixture=JSON.parse(fs.readFileSync(path.join(root,'ops/aws/evidence',file)));
 const user=admin?fixture:fixture[0];const page=await context.newPage();
 try {
 await page.goto(admin?hosts['super-admin']+'/login':hosts.landing+'/en/auth/login');
 await page.waitForFunction(()=>{const i=document.querySelector('input[type=email]');return i&&Object.keys(i).some(k=>k.startsWith('__reactProps')&&typeof i[k]?.onChange==='function');});
 await page.locator('input[type=email]').fill(user.email);await page.locator('input[type=password]').fill(user.password);
 await Promise.all([page.waitForURL(admin?'**/dashboard':'**/en/dashboard**'),page.locator('button[type=submit]').click()]);
 }catch(error){throw new Error(String(error.message).replaceAll(user.password,'[redacted]'));}
 finally{await page.close();}
}
async function verifyLogout(context,admin){
 const prefix=admin?'/super-admin':'';
 const result=await context.request.post('https://api.pos.rakmyat.online/api/v1'+prefix+'/auth/logout',{headers:{Origin:hosts[admin?'super-admin':'landing']},data:{}});
 let session;let rateLimitRetries=0;
 for(let attempt=0;attempt<3;attempt++){
  session=await context.request.get('https://api.pos.rakmyat.online/api/v1'+prefix+'/auth/session');
  if(session.status()!==429)break;
  rateLimitRetries++;const delay=Math.max(1,Math.min(60,Number(session.headers()['retry-after'])||60));
  await new Promise(resolve=>setTimeout(resolve,delay*1000));
 }
 rows.push({app:admin?'super-admin':'dashboard',route:'logout and revoked session',http:result.status(),session:session.status(),rateLimitRetries,status:result.ok()&&session.status()===401?'PASS':'FAIL'});
}
(async()=>{
 const browser=await chromium.launch({headless:true});
 const proof=await (await browser.newContext()).request.get('https://api.pos.rakmyat.online/_ops/release');
 run.observed_sha=(await proof.json()).sha;
 if(!run.sha||run.sha!==run.observed_sha)throw new Error('Full acceptance requires the exact deployed SHA');
 fs.writeFileSync(path.join(evidence,'run.json'),JSON.stringify(run,null,2));
 if(!dynamicOnly)await verifyAnonymousCallbacks(browser);
 const inventory=Object.fromEntries(Object.keys(hosts).map(app=>[app,routes(app)]));
 fs.writeFileSync(path.join(evidence,'discovered-routes.json'),JSON.stringify(inventory,null,2));
 run.expected_checks=(dynamicOnly?0:3)+2+Object.entries(inventory).reduce((total,[app,list])=>total+list.filter(raw=>!raw.includes('/callback')&&!raw.includes('/test-url')&&(app!=='dashboard'||raw.startsWith('/[locale]/dashboard'))&&(!dynamicOnly||/\[(id|tableId|orderNumber)\]/.test(raw))).reduce((sum,raw)=>sum+(raw.includes('[locale]')?6:3),0),0);
 for(const app of Object.keys(hosts)){
  const state=path.join(root,'ops/aws/evidence/storefront-qa-state.json');
  const context=await browser.newContext({viewport:{width:1440,height:900},...(app==='restaurant-website'?{storageState:state}:{})});
  const authenticated=app==='dashboard'||app==='super-admin';
  if(authenticated)await login(context,app==='super-admin');
  const page=await context.newPage();
  for(const raw of inventory[app]){
   if(raw.includes('/callback')||raw.includes('/test-url')){
    rows.push({app,route:raw,status:'NOT_APPLICABLE',reason:raw.includes('/callback')?'One-time authentication transport; full mode checks anonymous callback rejection separately.':'Developer URL diagnostic; outside authenticated business route inventory.'});continue;
   }
   if(app==='dashboard'&&!raw.startsWith('/[locale]/dashboard'))continue;
   if(dynamicOnly&&!/\[(id|tableId|orderNumber)\]/.test(raw))continue;
   const localized=raw.includes('[locale]');
   for(const locale of localized?['en','ar']:[''])for(const width of [390,768,1440])await inspect(page,app,resolveFixture(app,raw).replace('[locale]',locale),width);
  }
  if(authenticated)await verifyLogout(context,app==='super-admin');
  await context.close();
 }
 const finalProof=await (await browser.newContext()).request.get('https://api.pos.rakmyat.online/_ops/release');
 if((await finalProof.json()).sha!==run.sha)throw new Error('Deployed release changed during traversal');
 await browser.close();
 fs.writeFileSync(path.join(evidence,'results.json'),JSON.stringify(rows,null,2));
 run.completed_utc=new Date().toISOString();run.passed=rows.filter(r=>r.status==='PASS').length;run.failed=rows.filter(r=>r.status==='FAIL').length;run.not_tested=rows.filter(r=>r.status==='NOT_TESTED').length;
 if(run.passed+run.failed!==run.expected_checks)throw new Error('Discovered route coverage does not match actual checks');
 fs.writeFileSync(path.join(evidence,'run.json'),JSON.stringify(run,null,2));
 process.exitCode=run.failed?1:0;
})().catch(error=>{console.error(error.message);process.exitCode=1;});
