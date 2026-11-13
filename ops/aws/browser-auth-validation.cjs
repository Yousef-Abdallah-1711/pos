// Negative synthetic form acceptance: never creates an account or sends notifications.
const {chromium}=require('playwright');
const fs=require('fs');const path=require('path');
const evidence=process.env.SERVIO_BROWSER_EVIDENCE||path.join(__dirname,'evidence');
const checks=[];
function record(check,actual,expected){checks.push({check,actual,expected,status:actual===expected?'PASS':'FAIL'});}
async function hydrated(page){
 await page.waitForFunction(()=>{const i=document.querySelector('input[type=email]');return i&&Object.keys(i).some(k=>k.startsWith('__reactProps')&&typeof i[k]?.onChange==='function');});
}
async function validateLanding(browser,locale,form){
 const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1'});
 const page=await context.newPage();let mutations=0;
 page.on('request',request=>{if(request.method()==='POST'&&new URL(request.url()).pathname.startsWith('/api/v1/auth/'))mutations++;});
 try{
  await page.goto(`https://pos.rakmyat.online/${locale}/auth/${form}`);await hydrated(page);
  await page.locator('button[type=submit]').click();
  await page.locator('input[aria-invalid=true]').first().waitFor();
  record(`${locale} ${form} rejects empty required fields`,await page.locator('input[aria-invalid=true]').count()>0,true);
  record(`${locale} ${form} invalid submission performs no API mutation`,mutations,0);
 }finally{await context.close();}
}
async function validateAdmin(browser){
 const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1'});
 const page=await context.newPage();let mutations=0;
 page.on('request',request=>{if(request.method()==='POST'&&new URL(request.url()).pathname.endsWith('/super-admin/auth/login'))mutations++;});
 try{
  await page.goto('https://admin.pos.rakmyat.online/login');await hydrated(page);
  await page.locator('input[type=email]').fill('');await page.locator('button[type=submit]').click();
  await page.getByText('Please enter both email and password',{exact:true}).waitFor();
  record('Administrator empty fields show validation',true,true);
  record('Administrator empty fields perform no login mutation',mutations,0);
  await page.locator('input[type=email]').fill('servio-invalid-auth@example.invalid');
  await page.locator('input[type=password]').fill('Synthetic-invalid-password-only-42!');
  const [response]=await Promise.all([page.waitForResponse(r=>r.url().endsWith('/super-admin/auth/login')),page.locator('button[type=submit]').click()]);
  record('Invalid administrator credentials are rejected',response.status(),401);
  await page.locator('.text-red-700').waitFor();
  record('Administrator rejection remains visible after request',await page.locator('.text-red-700').isVisible(),true);
  record('Administrator login form remains usable',await page.locator('button[type=submit]').isEnabled(),true);
 }finally{await context.close();}
}
(async()=>{
 const browser=await chromium.launch();
 try{
  for(const locale of ['en','ar'])for(const form of ['login','register'])await validateLanding(browser,locale,form);
  await validateAdmin(browser);
 }catch(error){checks.push({check:'Negative authentication form journey',status:'FAIL',error:String(error.message).slice(0,300)});}
 finally{await browser.close();}
 fs.mkdirSync(evidence,{recursive:true});fs.writeFileSync(path.join(evidence,'auth-form-validation.json'),JSON.stringify({checks},null,2),{mode:0o600});
 for(const row of checks)console.log(row.status+' '+row.check);
 process.exitCode=checks.some(row=>row.status==='FAIL')?1:0;
})().catch(error=>{console.error(error.message);process.exitCode=1;});
