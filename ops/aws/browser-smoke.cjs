const { chromium } = require('playwright');
const { assertCompleteTargetConfig, apiBase, guardContext } = require('./acceptance-target.cjs');
assertCompleteTargetConfig();
const api = apiBase();
const fs = require('fs');
const evidence=process.env.SERVIO_BROWSER_EVIDENCE||'/evidence';
const fixtures=process.env.SERVIO_BROWSER_FIXTURES||'/fixtures';
const frameEvidence=`${evidence}/device-frames`;
fs.mkdirSync(frameEvidence,{recursive:true});
const rows=[];
const responsiveViewports=[
 {width:320,height:568},{width:360,height:800},{width:375,height:812},{width:390,height:844},
 {width:430,height:932},{width:600,height:960},{width:768,height:1024},{width:820,height:1180},
 {width:1024,height:768},{width:1180,height:820},{width:1280,height:800},{width:1440,height:900},{width:1920,height:1080},
];
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
async function snapshot(page,label,width,url,height=900) {
 const errors=[];const onError=e=>errors.push(e.message);page.on('pageerror',onError);
 try {
  await page.setViewportSize({width,height});
  const response=await page.goto(url,{waitUntil:'domcontentloaded',timeout:45000});
  await page.waitForTimeout(1200);
  const layout=await page.evaluate(()=>{const root=document.documentElement;const width=innerWidth;const overflow=root.scrollWidth>width+2;const elements=overflow?[...document.querySelectorAll('body *')].map(element=>{const rect=element.getBoundingClientRect();return{tag:element.tagName,id:element.id,className:typeof element.className==='string'?element.className.slice(0,100):'',right:Math.round(rect.right),width:Math.round(rect.width)}}).filter(element=>element.right>width+2&&element.width>0).sort((a,b)=>b.right-a.right).slice(0,8):[];return{lang:root.lang,dir:root.dir,overflow,viewportWidth:width,scrollWidth:root.scrollWidth,overflowElements:elements,marketingHeader:!!document.querySelector('.sv-header'),marketingFooter:!!document.querySelector('.sv-footer'),applicationError:/Application error|Internal Server Error/.test(document.body.innerText)}});
  await page.screenshot({path:`${evidence}/${label}-${width}x${height}.png`,fullPage:true});
  if((label.startsWith('authenticated-')||label.startsWith('admin-')||label.startsWith('landing-'))&&[320,390,768,1440].includes(width)) await captureDeviceFrame(page,label,width,height);
  const expectsMarketingShell=['landing-en','landing-ar','login','register'].includes(label);
  rows.push({label,width,height,url:page.url(),http:response.status(),layout,errors,status:response.status()<400&&!layout.applicationError&&!errors.length&&!layout.overflow&&(!new URL(url).pathname.startsWith('/ar')||(layout.lang==='ar'&&layout.dir==='rtl'))&&(!expectsMarketingShell||(layout.marketingHeader&&layout.marketingFooter))?'PASS':'FAIL'});
 } catch(e) { rows.push({label,width,status:'FAIL',error:e.message}); }
 page.off('pageerror',onError);
}
async function captureDeviceFrame(page,label,width,height) {
 const rawDirectory=`${evidence}/raw-viewport`;
 fs.mkdirSync(rawDirectory,{recursive:true});
 const rawPath=`${rawDirectory}/${label}-${width}x${height}.png`;
 const bytes=await page.screenshot({path:rawPath,fullPage:false});
 const {locale,direction,theme}=await page.evaluate(()=>({locale:document.documentElement.lang||'en',direction:document.documentElement.dir||'ltr',theme:document.documentElement.classList.contains('dark')?'dark':'light'}));
 const device=width<=430?'phone':width<=820?'tablet':width<1440?'laptop':'desktop';
 const frame=await page.context().newPage();
 const outer={width:width+100,height:height+156};
 await frame.setViewportSize(outer);
 const chrome=device==='laptop'||device==='desktop'?`<div class="chrome">Genuine browser capture · ${width} × ${height}</div>`:'';
 const encodedLabel=label.replace(/[&<>"']/g,char=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[char]));
 const base64=bytes.toString('base64');
 await frame.setContent(`<!doctype html><html lang="${locale}" dir="${direction}"><meta charset="utf-8"><style>
 *{box-sizing:border-box}html,body{margin:0;min-height:100%;background:#eeeaf3;color:#241d31;font-family:Arial,sans-serif}body{display:grid;place-items:center;padding:20px}.board{display:grid;justify-items:center;gap:10px}.caption{font-size:12px;font-weight:700;text-align:center}.device{padding:10px;background:#17131f;border:1px solid #625a6d;box-shadow:0 20px 50px #241d3130}.device.phone{padding:9px;border:8px solid #17131f;border-radius:38px}.device.tablet{padding:10px;border:8px solid #17131f;border-radius:24px}.device.laptop{border-radius:11px 11px 4px 4px}.device.desktop{border-radius:14px}.screen{display:block;width:${width}px;height:${height}px;object-fit:fill;background:#fff}.chrome{height:18px;margin:-2px 0 7px;border-radius:5px;background:#302a39;color:#ded5e8;font-size:9px;padding:3px 7px}.source{font-size:10px;color:#625b70;text-align:center}
 </style><body><main class="board"><div class="caption">SERVIO · ${locale.toUpperCase()} ${direction.toUpperCase()} · ${theme.toUpperCase()} · ${encodedLabel}</div><div class="device ${device}">${chrome}<img class="screen" width="${width}" height="${height}" src="data:image/png;base64,${base64}" alt="Unmodified SERVIO browser viewport"></div><div class="source">Presentation frame only · raw viewport: ${encodedLabel}-${width}x${height}.png</div></main></body></html>`);
 await frame.screenshot({path:`${frameEvidence}/${label}-${locale}-${theme}-${width}x${height}.png`,fullPage:true});
 await frame.close();
}
async function loginJourney(browser,admin=false) {
 const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width:1440,height:900}});
 await guardContext(context);

 const page=await context.newPage();const label=admin?'admin-browser-login':'tenant-browser-login';
 const fixture=JSON.parse(fs.readFileSync(admin?fixtures+'/admin.json':fixtures+'/users.json','utf8'));
 const user=admin?fixture:fixture[0];
 const adminBase='https://admin.pos.servio.test';
 const authenticateAdmin=async()=>{
  await context.clearCookies();
  await page.goto(`${adminBase}/login`,{waitUntil:'domcontentloaded',timeout:45000});
  await page.evaluate(()=>localStorage.removeItem('super-admin-auth-store'));
  await page.reload({waitUntil:'domcontentloaded'});
  await page.waitForFunction(()=>{const input=document.querySelector('input[type=email]');return input&&Object.keys(input).some(k=>k.startsWith('__reactProps')&&typeof input[k]?.onChange==='function');},{},{timeout:30000});
  await page.locator('input[type=email]').fill(user.email);
  await page.locator('input[type=password]').fill(user.password);
  await Promise.all([page.waitForURL('**/dashboard',{timeout:45000}),page.locator('button[type=submit]').click()]);
  await page.waitForFunction(async apiBase=>{
   const response=await fetch(`${apiBase}/super-admin/auth/session`,{credentials:'include'});
   return response.status===200;
  },api,{timeout:20000});
 };
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
  if(!admin) {
   const combinations=[
    {locale:'en',theme:'light',dir:'ltr'},
    {locale:'en',theme:'dark',dir:'ltr'},
    {locale:'ar',theme:'light',dir:'rtl'},
    {locale:'ar',theme:'dark',dir:'rtl'},
   ];
   for(const {locale,theme,dir} of combinations){
    const currentLocale=page.url().includes('/ar/')?'ar':'en';
    await page.goto(`https://pos.app.servio.test/${currentLocale}/dashboard`,{waitUntil:'domcontentloaded',timeout:45000});
    await page.locator('main h1').first().waitFor({state:'visible',timeout:30000});
    if(currentLocale!==locale){
     await page.getByRole('button',{name:currentLocale==='ar'?'تغيير اللغة':'Change language'}).click();
     await page.getByRole('button',{name:locale==='ar'?'عربي':'English'}).click();
     await page.waitForURL(`**/${locale}/dashboard**`,{timeout:30000});
    }
    const darkNow=await page.evaluate(()=>document.documentElement.classList.contains('dark'));
    const shouldBeDark=theme==='dark';
    if(darkNow!==shouldBeDark){
     await page.getByRole('button',{name:locale==='ar'?(shouldBeDark?'التبديل إلى الوضع الداكن':'التبديل إلى الوضع الفاتح'):`Switch to ${shouldBeDark?'dark':'light'} mode`}).click();
     await page.waitForFunction(expected=>document.documentElement.classList.contains('dark')===expected,shouldBeDark);
    }
    await page.reload({waitUntil:'domcontentloaded'});
    await page.waitForFunction(expected=>document.documentElement.lang===expected.locale&&document.documentElement.dir===expected.dir&&document.documentElement.classList.contains('dark')===(expected.theme==='dark'),{locale,dir,theme});
    for(const route of ['dashboard','dashboard/pos','dashboard/kds']) for(const viewport of responsiveViewports){
     const url=`https://pos.app.servio.test/${locale}/${route}`;
     await snapshot(page,`authenticated-${locale}-${theme}-${route.replaceAll('/','-')}`,viewport.width,url,viewport.height);
    }
   }
  }
  else {
   const localeThemes=[
    {locale:'en',dir:'ltr',theme:'light'},
    {locale:'en',dir:'ltr',theme:'dark'},
    {locale:'ar',dir:'rtl',theme:'light'},
    {locale:'ar',dir:'rtl',theme:'dark'},
   ];
   const adminPageErrors=[];
   const adminConsoleErrors=[];
   const adminAssetFailures=[];
   const onAdminPageError=error=>adminPageErrors.push(error.message);
   page.on('pageerror',onAdminPageError);
   page.on('console',message=>{if(message.type()==='error')adminConsoleErrors.push(message.text())});
   page.on('response',response=>{if(response.status()>=400&&response.request().resourceType()==='image')adminAssetFailures.push(`${response.status()} ${new URL(response.url()).pathname}`)});
   const adminRoutes=['/dashboard','/dashboard/tenants','/dashboard/users','/dashboard/audit-logs','/dashboard/payments','/dashboard/settings'];
   const switchTo=async(locale,theme)=>{
    const current=await page.evaluate(()=>({locale:document.documentElement.lang,theme:document.documentElement.classList.contains('dark')}));
    if(current.locale!==locale) await page.getByRole('button',{name:current.locale==='en'?'Switch to Arabic':'التغيير إلى الإنجليزية'}).click();
    await page.waitForFunction(expected=>document.documentElement.lang===expected&&document.documentElement.dir===(expected==='ar'?'rtl':'ltr'),locale);
    const dark=current.theme;
    const desired=theme==='dark';
    if(dark!==desired) await page.getByRole('button',{name:locale==='ar'?(desired?'تفعيل الوضع الداكن':'تفعيل الوضع الفاتح'):`Switch to ${desired?'dark':'light'} mode`}).click();
    await page.waitForFunction(expected=>document.documentElement.classList.contains('dark')===expected,desired);
   };
   for(const {locale,dir,theme} of localeThemes){
    await authenticateAdmin();
    await switchTo(locale,theme);
    await page.reload({waitUntil:'domcontentloaded'});
    await page.waitForFunction(expected=>document.documentElement.lang===expected.locale&&document.documentElement.dir===(expected.locale==='ar'?'rtl':'ltr')&&document.documentElement.classList.contains('dark')===(expected.theme==='dark'),{locale,theme});
    for(const route of adminRoutes){
     let response=await page.goto(adminBase+route,{waitUntil:'domcontentloaded',timeout:45000});
     await page.waitForFunction(()=>window.location.pathname==='/login'||document.querySelector('main h1'),null,{timeout:30000});
     if(new URL(page.url()).pathname==='/login'){
      await authenticateAdmin();
      response=await page.goto(adminBase+route,{waitUntil:'domcontentloaded',timeout:45000});
      await page.locator('main h1').first().waitFor({state:'visible',timeout:30000});
     } else {
      await page.locator('main h1').first().waitFor({state:'visible',timeout:30000});
     }
     await page.waitForTimeout(350);
     for(const {width,height} of responsiveViewports){
      await page.setViewportSize({width,height});
      const layout=await page.evaluate(()=>({lang:document.documentElement.lang,dir:document.documentElement.dir,dark:document.documentElement.classList.contains('dark'),overflow:document.documentElement.scrollWidth>innerWidth+2,scrollWidth:document.documentElement.scrollWidth,viewportWidth:innerWidth}));
      const screenshotName=`admin-${locale}-${theme}-${route.replaceAll('/','-')||'home'}-${width}x${height}`;
      await page.screenshot({path:`${evidence}/${screenshotName}.png`,fullPage:true});
      if([320,390,768,1440].includes(width)) await captureDeviceFrame(page,screenshotName,width,height);
      const heading=await page.locator('main h1').first().innerText();
      const localeContentOkay=locale!=='ar'||route!=='/dashboard/tenants'||heading==='الجهات';
      const expectedSurface=theme==='dark'?'rgb(31, 41, 55)':'rgb(255, 255, 255)';
      const surface=await page.locator('.bg-white').first().evaluate(element=>getComputedStyle(element).backgroundColor);
      const primary=route==='/dashboard/settings'?await page.locator('.bg-blue-600').first().evaluate(element=>getComputedStyle(element).backgroundColor).catch(()=>null):null;
      const primaryBrandOkay=route!=='/dashboard/settings'||primary==='rgb(124, 58, 237)';
      const okay=response?.status()===200&&layout.lang===locale&&layout.dir===dir&&layout.dark===(theme==='dark')&&!layout.overflow&&localeContentOkay&&surface===expectedSurface&&primaryBrandOkay&&!adminPageErrors.length&&!adminConsoleErrors.length&&!adminAssetFailures.length;
      rows.push({label:screenshotName,url:page.url(),width,height,http:response?.status(),heading,layout,surface,primaryBrandColor:primary,pageErrors:[...adminPageErrors],consoleErrors:[...adminConsoleErrors],assetFailures:[...adminAssetFailures],status:okay?'PASS':'FAIL'});
     }
     if(route==='/dashboard/tenants'){
      const viewLink=page.getByRole('link',{name:locale==='ar'?'عرض':'View'}).first();
      await viewLink.waitFor({state:'visible',timeout:30000});
      await viewLink.click();
      await page.waitForURL(/\/dashboard\/tenants\/\d+$/, {timeout:30000});
      await page.locator('main h1').first().waitFor({state:'visible',timeout:30000});
      for(const {width,height} of responsiveViewports){
       await page.setViewportSize({width,height});
       const layout=await page.evaluate(()=>({lang:document.documentElement.lang,dir:document.documentElement.dir,dark:document.documentElement.classList.contains('dark'),overflow:document.documentElement.scrollWidth>innerWidth+2,scrollWidth:document.documentElement.scrollWidth,viewportWidth:innerWidth}));
       const screenshotName=`admin-${locale}-${theme}-tenant-detail-${width}x${height}`;
       await page.screenshot({path:`${evidence}/${screenshotName}.png`,fullPage:true});
       if([320,390,768,1440].includes(width)) await captureDeviceFrame(page,screenshotName,width,height);
       rows.push({label:screenshotName,url:page.url(),width,height,layout,status:layout.lang===locale&&layout.dir===dir&&layout.dark===(theme==='dark')&&!layout.overflow&&!adminPageErrors.length?'PASS':'FAIL'});
      }
      await page.setViewportSize({width:390,height:844});
      const extendLabel=locale==='ar'?'تمديد الفترة التجريبية':'Extend Trial';
      const extendButton=page.getByRole('button',{name:extendLabel}).first();
      if(await extendButton.isEnabled()){
       await extendButton.click();
       await page.getByRole('heading',{name:locale==='ar'?'تمديد الفترة التجريبية':'Extend Trial Period'}).waitFor({state:'visible'});
       const modalOverflow=await page.evaluate(()=>document.documentElement.scrollWidth>innerWidth+2);
       const modalPath=`${evidence}/admin-${locale}-${theme}-extend-trial-modal-390.png`;
       await page.screenshot({path:modalPath,fullPage:true});
       await captureDeviceFrame(page,`admin-${locale}-${theme}-extend-trial-modal-390`,390,844);
       rows.push({label:`admin-${locale}-${theme}-extend-trial-modal-390`,url:page.url(),width:390,modalOverflow,status:modalOverflow?'FAIL':'PASS'});
       await page.getByRole('button',{name:locale==='ar'?'إلغاء':'Cancel'}).last().click();
      }
     }
    }
   }
   page.off('pageerror',onAdminPageError);
  }
 } catch(e) {
  rows.push({label,status:'FAIL',error:String(e.message).replaceAll(user.password,'[redacted]'),url:page.url()});
  await page.screenshot({path:`${evidence}/${label}-failure.png`,fullPage:true}).catch(()=>{});
 }
 await context.close();
}
(async()=>{
 const browser=await chromium.launch({headless:true,args:['--disable-dev-shm-usage']});
 for(const {width,height} of responsiveViewports) {
  const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width,height}});await guardContext(context);const page=await context.newPage();

  for(const [label,url] of anonymous) await snapshot(page,label,width,url,height);
  await context.close();
 }
 await loginJourney(browser);await loginJourney(browser,true);
 await browser.close();
 fs.writeFileSync(evidence+'/results.json',JSON.stringify(rows,null,2));
 for(const row of rows)console.log(row.status+' '+row.label+(row.width?' '+row.width:''));
 process.exitCode=rows.some(r=>r.status==='FAIL')?1:0;
})().catch(e=>{console.error(e.message);process.exit(1)});
