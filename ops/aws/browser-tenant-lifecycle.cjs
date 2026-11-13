// Guarded demo-tenant suspension. Always restore its original active state.
const{chromium}=require('playwright');const fs=require('fs');const path=require('path');
const evidence=path.resolve(__dirname,'evidence');const admin=JSON.parse(fs.readFileSync(path.join(evidence,'browser-admin.json')));const fixture=JSON.parse(fs.readFileSync(path.join(evidence,'storefront-qa-fixture.json')));
const checks=[];const base='https://api.pos.rakmyat.online/api/v1';let disabled=false;
function record(check,actual,expected){checks.push({check,actual,expected,status:actual===expected?'PASS':'FAIL'});}
(async()=>{
 const browser=await chromium.launch();const owner=await browser.newContext();const guest=await browser.newContext({viewport:{width:390,height:900}});const page=await guest.newPage();
 const adminHeaders={Origin:'https://admin.pos.rakmyat.online'};
 try{
  const login=await owner.request.post(base+'/super-admin/auth/login',{headers:adminHeaders,data:{email:admin.email,password:admin.password}});if(login.status()!==200)throw new Error('Synthetic administrator login failed');
  const response=await owner.request.get(base+'/super-admin/tenants/1');const raw=await response.json();const tenant=raw.data||raw;
  if(tenant.name!=='Test Tenant'||tenant.slug!=='restaurant-1'||tenant.is_active!==true)throw new Error('Refusing to alter a tenant outside the active seeded demo');
  const entitlement=async()=>{const raw=await(await owner.request.get(base+'/super-admin/tenants/1/subscription')).json();const s=raw.data||raw;return JSON.stringify([s.subscription_status,s.subscription_plan,s.trial_start_date,s.trial_end_date,s.subscription_expires_at,s.last_payment_at]);};
  const before=await entitlement();
  const suspend=await owner.request.put(base+'/super-admin/tenants/1/disable',{headers:adminHeaders,data:{reason:'Synthetic storefront suspension acceptance'}});
  if(suspend.status()!==200)throw new Error('Authorized demo suspension failed');disabled=true;
  const status=await guest.request.get(base+'/public/tenant-status/restaurant-1');const statusRaw=await status.json();record('Suspended demo tenant cannot access storefront',(statusRaw.data||statusRaw).can_access,false);
  await page.goto('https://restaurant-1.pos.rakmyat.online/en');record('Suspended storefront renders maintenance',new URL(page.url()).pathname,'/maintenance');
  await page.screenshot({path:path.join(evidence,'suspended-storefront.png'),fullPage:true});
  const create=await guest.request.post(base+'/public/orders',{headers:{Origin:'https://restaurant-1.pos.rakmyat.online','X-Tenant-ID':'1','X-Restaurant-ID':'1'},data:{customer_name:'SERVIO disabled-tenant QA',customer_email:'disabled-qa@example.invalid',customer_phone:'+12025550123',delivery_address:'Synthetic only. Do not deliver.',payment_method:'cash',order_source:'website',notes:'Negative staging acceptance. Do not fulfill.',items:[{product_id:fixture.product_id,quantity:1}]}});
  record('Suspended tenant rejects new anonymous orders',[401,402,403,423].includes(create.status()),true);
  if(create.status()===201){
   const body=await create.json();const order=body.data||body;
   if(order.tracking_token){await guest.request.delete(base+'/public/orders/'+order.id,{headers:{Origin:'https://restaurant-1.pos.rakmyat.online','X-Tenant-ID':'1','X-Restaurant-ID':'1','X-Order-Token':order.tracking_token}});}
  }
  const missing=await guest.request.get(base+'/public/tenant-status/servio-does-not-exist');record('Missing tenant status fails closed',missing.status(),404);
  const enable=await owner.request.put(base+'/super-admin/tenants/1/enable',{headers:adminHeaders,data:{}});if(enable.status()!==200)throw new Error('Demo restoration failed');disabled=false;
  const afterRaw=await(await owner.request.get(base+'/super-admin/tenants/1')).json();const after=afterRaw.data||afterRaw;
  record('Demo tenant restored to original active state',after.is_active,true);record('Subscription and trial state unchanged',await entitlement(),before);
 }catch(error){checks.push({check:'Tenant lifecycle browser acceptance',status:'FAIL',error:String(error.message).replaceAll(admin.password,'[redacted]')});}
 finally{
  if(disabled){const restore=await owner.request.put(base+'/super-admin/tenants/1/enable',{headers:adminHeaders,data:{}});record('Emergency demo restoration succeeded',restore.status(),200);}
  await owner.request.post(base+'/super-admin/auth/logout',{headers:adminHeaders,data:{}});await browser.close();
 }
 fs.writeFileSync(path.join(evidence,'tenant-lifecycle-browser.json'),JSON.stringify(checks,null,2),{mode:0o600});for(const row of checks)console.log(row.status+' '+row.check);process.exitCode=checks.some(r=>r.status==='FAIL')?1:0;
})().catch(error=>{console.error(String(error.message).replaceAll(admin.password,'[redacted]'));process.exitCode=1;});
