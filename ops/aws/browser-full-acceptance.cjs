// Full route discovery and read-only browser acceptance. Evidence contains no passwords.
const {chromium}=require('playwright');
const fs=require('fs');
const path=require('path');
const root=path.resolve(__dirname,'../..');
const evidence=path.join(root,'ops/aws/evidence/full-browser');
fs.mkdirSync(evidence,{recursive:true,mode:0o700});
const rows=[];
const hosts={landing:'https://pos.rakmyat.online',dashboard:'https://pos.app.rakmyat.online','super-admin':'https://admin.pos.rakmyat.online','restaurant-website':'https://restaurant-1.pos.rakmyat.online'};
function discover(directory){
 return fs.readdirSync(directory,{withFileTypes:true}).flatMap(entry=>entry.isDirectory()?discover(path.join(directory,entry.name)):entry.name==='page.tsx'?[path.join(directory,entry.name)]:[]);
}
function routes(app){
 const base=path.join(root,'frontend/apps',app,app==='restaurant-website'?'app':'src/app');
 return [...new Set(discover(base).map(file=>'/'+path.relative(base,file).replace(/\/page.tsx$|^page.tsx$/g,'')))].sort();
}
function safeURL(url){const parsed=new URL(url);return parsed.origin+parsed.pathname;}
async function inspect(page,app,route,width){
 const row={app,route,width,status:'FAIL',console:[],network:[],errors:[]};
 const error=event=>row.errors.push(event.message.slice(0,300));
 const log=message=>{if(message.type()==='error')row.console.push(message.text().slice(0,300));};
 const response=result=>{if(result.status()>=400&&!/\/auth\/session$/.test(new URL(result.url()).pathname))row.network.push({url:safeURL(result.url()),status:result.status()});};
 page.on('pageerror',error);page.on('console',log);page.on('response',response);
 try{
  await page.setViewportSize({width,height:900});
  const result=await page.goto(hosts[app]+route,{waitUntil:'domcontentloaded',timeout:45000});
  await page.waitForTimeout(1200);
  row.http=result.status();row.finalURL=safeURL(page.url());
  row.layout=await page.evaluate(()=>({lang:document.documentElement.lang,dir:document.documentElement.dir,overflow:document.documentElement.scrollWidth>innerWidth+2,bodyLength:document.body.innerText.trim().length,applicationError:/Application error|Internal Server Error/.test(document.body.innerText),brokenImages:[...document.images].filter(i=>i.complete&&i.currentSrc&&i.naturalWidth===0).map(i=>({alt:i.alt,src:new URL(i.currentSrc).pathname})),forms:[...document.forms].map(f=>({fields:f.elements.length,invalid:f.querySelectorAll(':invalid').length})),unlabelledFields:[...document.querySelectorAll('input:not([type=hidden]),select,textarea')].filter(i=>!i.labels?.length&&!i.getAttribute('aria-label')&&!i.getAttribute('aria-labelledby')).map(i=>({type:i.type,name:i.name})),headings:document.querySelectorAll('h1').length}));
  await page.keyboard.press('Tab');
  row.keyboardFocus=await page.evaluate(()=>({tag:document.activeElement.tagName,visible:!!document.activeElement.getClientRects().length}));
  const locale=route.startsWith('/ar')?'ar':route.startsWith('/en')?'en':null;
  const localeOK=!locale||(row.layout.lang===locale&&(locale!=='ar'||row.layout.dir==='rtl'));
  row.status=row.http<400&&row.layout.bodyLength>20&&!row.layout.applicationError&&!row.layout.overflow&&!row.errors.length&&!row.layout.brokenImages.length&&!row.network.length&&localeOK?'PASS':'FAIL';
  if(row.status==='FAIL'||route==='/ar'||route==='/en'){
   row.screenshot=`${app}-${route.replaceAll('/','_')||'root'}-${width}.png`;
   await page.screenshot({path:path.join(evidence,row.screenshot),fullPage:true});
  }
 }catch(error){row.error=error.message.slice(0,300);}
 page.off('pageerror',error);page.off('console',log);page.off('response',response);
 rows.push(row);fs.writeFileSync(path.join(evidence,'results.json'),JSON.stringify(rows,null,2));
 console.log(`${row.status} ${app} ${route} ${width}`);
}
async function login(context,admin=false){
 const file=admin?'browser-admin.json':'browser-fixtures.json';
 const fixture=JSON.parse(fs.readFileSync(path.join(root,'ops/aws/evidence',file)));
 const user=admin?fixture:fixture[0];const page=await context.newPage();
 await page.goto(admin?hosts['super-admin']+'/login':hosts.landing+'/en/auth/login');
 await page.waitForFunction(()=>{const i=document.querySelector('input[type=email]');return i&&Object.keys(i).some(k=>k.startsWith('__reactProps')&&typeof i[k]?.onChange==='function');});
 await page.locator('input[type=email]').fill(user.email);await page.locator('input[type=password]').fill(user.password);
 await Promise.all([page.waitForURL(admin?'**/dashboard':'**/en/dashboard**'),page.locator('button[type=submit]').click()]);
 await page.close();
}
async function verifyLogout(context,admin){
 const prefix=admin?'/super-admin':'';
 const result=await context.request.post('https://api.pos.rakmyat.online/api/v1'+prefix+'/auth/logout',{headers:{Origin:hosts[admin?'super-admin':'landing']},data:{}});
 const session=await context.request.get('https://api.pos.rakmyat.online/api/v1'+prefix+'/auth/session');
 rows.push({app:admin?'super-admin':'dashboard',route:'logout and revoked session',http:result.status(),session:session.status(),status:result.ok()&&session.status()===401?'PASS':'FAIL'});
}
(async()=>{
 const browser=await chromium.launch({headless:true});
 const inventory=Object.fromEntries(Object.keys(hosts).map(app=>[app,routes(app)]));
 fs.writeFileSync(path.join(evidence,'discovered-routes.json'),JSON.stringify(inventory,null,2));
 for(const app of Object.keys(hosts)){
  const context=await browser.newContext({viewport:{width:1440,height:900}});
  const authenticated=app==='dashboard'||app==='super-admin';
  if(authenticated)await login(context,app==='super-admin');
  const page=await context.newPage();
  for(const raw of inventory[app]){
   if(raw.includes('[id]')||raw.includes('[tableId]')||raw.includes('[orderNumber]')||raw.includes('/test-url')||raw.includes('/callback')){
    rows.push({app,route:raw,status:'NOT_TESTED',reason:'Requires a synthetic resource or one-time callback; generic traversal does not invent identifiers.'});continue;
   }
   if(app==='dashboard'&&!raw.startsWith('/[locale]/dashboard'))continue;
   const localized=raw.includes('[locale]');
   for(const locale of localized?['en','ar']:[''])for(const width of [390,768,1440])await inspect(page,app,raw.replace('[locale]',locale),width);
  }
  if(authenticated)await verifyLogout(context,app==='super-admin');
  await context.close();
 }
 await browser.close();
 fs.writeFileSync(path.join(evidence,'results.json'),JSON.stringify(rows,null,2));
 process.exitCode=rows.some(row=>row.status==='FAIL')?1:0;
})().catch(error=>{console.error(error.message);process.exitCode=1;});
