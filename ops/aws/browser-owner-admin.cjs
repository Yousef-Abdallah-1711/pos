// Verify the owner's provisioned administrator without printing its credentials.
const {chromium}=require('playwright');
const fs=require('fs');
const path=require('path');
const root=path.resolve(__dirname,'../..');
const evidence=path.join(root,'ops/aws/evidence');
const owner=JSON.parse(fs.readFileSync(path.join(evidence,'owner-admin.json')));
const checks=[];
function record(check,actual,expected){checks.push({check,actual,expected,status:actual===expected?'PASS':'FAIL'});}
(async()=>{
 const browser=await chromium.launch();const context=await browser.newContext();const page=await context.newPage();
 try{
  const response=await page.goto('https://admin.pos.rakmyat.online/login');record('Public owner login without gateway',response.status(),200);
  await page.waitForFunction(()=>{const i=document.querySelector('input[type=email]');return i&&Object.keys(i).some(k=>k.startsWith('__reactProps')&&typeof i[k]?.onChange==='function');});
  await page.locator('input[type=email]').fill(owner.email);await page.locator('input[type=password]').fill(owner.password);
  await Promise.all([page.waitForURL('**/dashboard'),page.locator('button[type=submit]').click()]);
  const statuses=await page.evaluate(async()=>{
   const base='https://api.pos.rakmyat.online/api/v1/super-admin';
   const request=route=>fetch(base+route,{credentials:'include'}).then(r=>r.status);
   return {session:await request('/auth/session'),tenants:await request('/tenants?page=1&limit=1')};
  });
  record('Owner authenticated session',statuses.session,200);record('Owner authorized tenant administration',statuses.tenants,200);
  const cookies=await context.cookies('https://api.pos.rakmyat.online');
  record('Owner session has Secure HttpOnly cookie',cookies.some(c=>c.secure&&c.httpOnly),true);
  await page.screenshot({path:path.join(evidence,'owner-admin-dashboard.png'),fullPage:true});
  await page.getByRole('link',{name:'Logout',exact:true}).or(page.getByRole('button',{name:'Logout',exact:true})).click();
  await page.waitForTimeout(1200);
  const result=await page.evaluate(async()=>{
   const base='https://api.pos.rakmyat.online/api/v1/super-admin/auth';
   const session=await fetch(base+'/session',{credentials:'include'});
   if(session.status===200)await fetch(base+'/logout',{method:'POST',credentials:'include',headers:{'Content-Type':'application/json'},body:'{}'});
   return {session:session.status};
  });
  record('Owner UI logout revokes session',result.session,401);
  await page.goto('https://admin.pos.rakmyat.online/dashboard/tenants');await page.waitForURL('**/login',{timeout:30000});
  record('Logged-out owner cannot view protected tenants',new URL(page.url()).pathname,'/login');
 }catch(error){checks.push({check:'Owner administrator browser journey',status:'FAIL',error:String(error.message).replaceAll(owner.password,'[redacted]')});}
 finally{await browser.close();}
 const target=path.join(evidence,'owner-admin-browser-results.json');fs.writeFileSync(target,JSON.stringify(checks,null,2),{mode:0o600});
 for(const row of checks)console.log(row.status+' '+row.check);
 process.exitCode=checks.some(row=>row.status==='FAIL')?1:0;
})().catch(error=>{console.error(String(error.message).replaceAll(owner.password,'[redacted]'));process.exitCode=1;});
