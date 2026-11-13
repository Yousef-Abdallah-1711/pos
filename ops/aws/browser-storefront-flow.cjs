// Synthetic cash checkout only. No provider payments, real contacts, or token output.
const {chromium}=require('playwright');const fs=require('fs');const path=require('path');
const root=path.resolve(__dirname,'../..');const evidence=process.env.SERVIO_BROWSER_EVIDENCE||path.join(root,'ops/aws/evidence');
const fixture=JSON.parse(fs.readFileSync(process.env.SERVIO_STORE_FIXTURE||path.join(evidence,'storefront-qa-fixture.json')));
if(!fixture.synthetic||fixture.tenant_id!==1||fixture.restaurant_id!==1)throw new Error('Guarded synthetic storefront fixture required');
const checks=[];const errors=[];let created=null;let submitted=0;
function record(check,actual,expected){checks.push({check,actual,expected,status:actual===expected?'PASS':'FAIL'});}
async function controlContrast(locator){
 return locator.evaluate(element=>{
  const style=getComputedStyle(element);const canvas=document.createElement('canvas');
  canvas.width=canvas.height=1;const context=canvas.getContext('2d');
  const luminance=color=>{
   context.clearRect(0,0,1,1);context.fillStyle=color;context.fillRect(0,0,1,1);
   const pixel=context.getImageData(0,0,1,1).data;
   const values=[...pixel].slice(0,3).map(value=>{const s=value/255;return s<=0.04045?s/12.92:((s+0.055)/1.055)**2.4;});
   return values[0]*0.2126+values[1]*0.7152+values[2]*0.0722;
  };
  const foreground=luminance(style.color);const background=luminance(style.backgroundColor);
  return (Math.max(foreground,background)+0.05)/(Math.min(foreground,background)+0.05);
 });
}

(async()=>{
 const browser=await chromium.launch();const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width:390,height:900}});const page=await context.newPage();
 page.on('pageerror',e=>errors.push(e.message.slice(0,200)));
 page.on('request',r=>{if(r.method()==='POST'&&new URL(r.url()).pathname==='/api/v1/public/orders')submitted++;});
 try{
  const status=await context.request.get('https://api.pos.rakmyat.online/api/v1/public/tenant-status/restaurant-1');
  const statusBody=await status.json();
  if((statusBody.data||statusBody).tenant_name!=='Test Tenant')throw new Error('Refusing checkout outside the exact synthetic demo tenant');
  await page.goto('https://restaurant-1.pos.rakmyat.online/en/menu');
  await page.getByText('SERVIO staging QA item',{exact:true}).waitFor();
  for(const label of ['Search products','Change language'])record(label+' has visible contrast',await controlContrast(page.getByRole('button',{name:label,exact:true}))>=4.5,true);
  record('Synthetic product and inline details visible',await page.getByText('Synthetic staging acceptance. No real sale.',{exact:true}).isVisible(),true);
  await page.getByRole('button',{name:'Add',exact:true}).click();
  await page.goto('https://restaurant-1.pos.rakmyat.online/en/cart');
  await page.getByText('SERVIO staging QA item',{exact:true}).waitFor();
  record('Cart shows the selected synthetic product',true,true);
  await page.screenshot({path:path.join(evidence,'storefront-cart-390.png'),fullPage:true});
  await page.goto('https://restaurant-1.pos.rakmyat.online/en/checkout');
  await page.locator('button[type=submit]').click();await page.waitForTimeout(300);
  record('Empty form does not submit an order',submitted,0);
  record('Checkout supports cash only',await page.locator('input[name=paymentMethod]').inputValue(),'cash');
  await page.locator('input[name=customerName]').fill('SERVIO synthetic browser QA');
  await page.locator('input[name=customerEmail]').fill('servio-browser-qa@example.invalid');
  await page.locator('input[name=customerPhone]').fill('+12025550123');
  await page.locator('textarea[name=deliveryAddress]').fill('Synthetic staging address. No delivery.');
  await page.locator('textarea[name=specialInstructions]').fill('Synthetic staging QA only. Do not fulfill or notify.');
  const quote=Number((await page.locator('.checkout-grand-total').innerText()).replace(/[^0-9.]/g,''));
  const responsePromise=page.waitForResponse(r=>r.request().method()==='POST'&&new URL(r.url()).pathname==='/api/v1/public/orders');
  await page.locator('button[type=submit]').click();const response=await responsePromise;
  record('Synthetic cash order accepted',response.status(),201);
  const raw=await response.json();created=raw.data||raw;
  fixture.order_id=created.id;fixture.order_number=created.order_number;fixture.stored_total=created.total_amount;fixture.displayed_total=quote;
  record('Displayed checkout total matches authoritative order',quote,Number(created.total_amount));
  await page.getByText(created.order_number,{exact:true}).waitFor();
  record('Confirmation displays returned order number',true,true);
  await page.screenshot({path:path.join(evidence,'storefront-confirmation-390.png'),fullPage:true});
  await page.goto('https://restaurant-1.pos.rakmyat.online/en/orders');
  await page.getByText(created.order_number,{exact:false}).first().waitFor();
  record('Device-held history displays its authorized order',true,true);
  await page.goto('https://restaurant-1.pos.rakmyat.online/en/orders/'+created.id);
  await page.getByText(created.order_number,{exact:false}).first().waitFor({timeout:10000});
  record('Authorized order detail renders',true,true);
 }catch(error){checks.push({check:'Storefront browser flow',status:'FAIL',error:error.message.slice(0,300)});await page.screenshot({path:path.join(evidence,'storefront-flow-failure.png'),fullPage:true});}
 finally{
  if(created){
   const response=await context.request.get('https://api.pos.rakmyat.online/api/v1/public/orders/'+created.id,{headers:{'X-Tenant-ID':'1','X-Restaurant-ID':'1'}});
   record('Order identifier alone does not authorize access',response.status(),404);
  }
  record('No browser runtime exceptions',errors.length,0);
  await context.storageState({path:path.join(evidence,'storefront-qa-state.json')});
  fs.chmodSync(path.join(evidence,'storefront-qa-state.json'),0o600);
  fs.writeFileSync(path.join(evidence,'storefront-qa-fixture.json'),JSON.stringify(fixture,null,2),{mode:0o600});
  await browser.close();
 }
 fs.writeFileSync(path.join(evidence,'storefront-browser-flow.json'),JSON.stringify({checks,errors},null,2),{mode:0o600});
 for(const row of checks)console.log(row.status+' '+row.check);
 process.exitCode=checks.some(r=>r.status==='FAIL')?1:0;
})().catch(error=>{console.error(error.message);process.exitCode=1;});
