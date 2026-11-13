// Synthetic catalog image acceptance: no real contacts, payments or notifications.
const {chromium}=require('playwright');
const fs=require('fs');const path=require('path');
const evidence=process.env.SERVIO_BROWSER_EVIDENCE||path.join(__dirname,'evidence');
const usersFile=process.env.SERVIO_UPLOAD_USERS||path.join(__dirname,'evidence/browser-fixtures.json');
const user=JSON.parse(fs.readFileSync(usersFile))[0];
if(!user.email.endsWith('@example.invalid'))throw new Error('Synthetic tenant fixture required');
const checks=[];const errors=[];const base='https://api.pos.rakmyat.online/api/v1';
const headers={Origin:'https://pos.app.rakmyat.online'};
const png=Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGNw6TjzHwAFRAKYH64YhwAAAABJRU5ErkJggg==','base64');
function record(check,actual,expected){checks.push({check,actual,expected,status:actual===expected?'PASS':'FAIL'});}
async function object(response){const raw=await response.json();return raw.data||raw;}
(async()=>{
 const browser=await chromium.launch();
 const context=await browser.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1',viewport:{width:390,height:900}});
 const page=await context.newPage();page.on('pageerror',error=>errors.push(error.message.slice(0,200)));
 try{
  const login=await context.request.post(base+'/auth/login',{headers,data:{email:user.email,password:user.password}});
  if(login.status()!==200)throw new Error('Synthetic tenant login failed');
  const session=await object(await context.request.get(base+'/auth/session'));
  if(session.email!==user.email||session.tenant_id<=1)throw new Error('Synthetic tenant scope mismatch');
  const categoryResponse=await context.request.post(base+'/categories',{headers,data:{name:'SERVIO synthetic upload QA '+Date.now(),is_active:true}});
  if(categoryResponse.status()!==201)throw new Error('Synthetic category creation failed');
  const category=await object(categoryResponse);
  const product={name_en:'SERVIO synthetic image QA',name_ar:'اختبار الصورة',price:1,category_id:category.id,sku:'qa-image-'+Date.now(),track_inventory:false,is_available:true,status:'active'};
  const create=await context.request.post(base+'/products',{headers,multipart:{data:JSON.stringify(product),image:{name:'untrusted-name.html',mimeType:'image/png',buffer:png}}});
  record('Synthetic image upload accepted',create.status(),201);
  if(create.status()!==201)throw new Error('Synthetic image creation failed');
  const created=await object(create);const saved=await object(await context.request.get(base+'/products/'+created.id));
  const initialURL=saved.main_image_url;
  record('Image URL persisted on creation',typeof initialURL==='string'&&initialURL.endsWith('.png'),true);
  if(!initialURL)throw new Error('Image URL missing after creation');
  record('Client filename is replaced',initialURL.includes('untrusted-name'),false);
  const image=await context.request.get(initialURL);
  record('Image available over HTTPS',image.status(),200);
  record('Image response has safe raster type',image.headers()['content-type'],'image/png');
  record('Image response disables MIME sniffing',image.headers()['x-content-type-options'],'nosniff');
  record('Stored image bytes match',Buffer.compare(await image.body(),png),0);
  const update=await context.request.put(base+'/products/'+created.id,{headers,multipart:{data:JSON.stringify({name_en:product.name_en}),image:{name:'replacement.png',mimeType:'image/png',buffer:png}}});
  record('Replacement image accepted',update.status(),200);
  const updated=await object(await context.request.get(base+'/products/'+created.id));
  record('Replacement URL persisted',typeof updated.main_image_url==='string'&&updated.main_image_url!==initialURL,true);
  record('Previous uploaded image retained',(await context.request.get(initialURL)).status(),200);
  const invalid=await context.request.post(base+'/products',{headers,multipart:{data:JSON.stringify({...product,sku:product.sku+'-invalid'}),image:{name:'spoof.png',mimeType:'image/png',buffer:Buffer.from('<html><script>invalid()</script></html>')}}});
  record('Spoofed executable upload rejected',invalid.status(),400);
  record('Upload directories are not listed',(await context.request.get('https://api.pos.rakmyat.online/uploads/products/')).status(),404);
  await page.goto('https://pos.app.rakmyat.online/en/dashboard/products/'+created.id+'/edit');
  await page.waitForFunction(url=>[...document.images].some(image=>{
   const value=image.currentSrc||image.src;if(!value)return false;const src=new URL(value,location.href);return (src.href===url||src.searchParams.get('url')===url)&&image.complete&&image.naturalWidth>0;
  }),updated.main_image_url,{timeout:20000});
  record('Dashboard renders persisted uploaded image',true,true);
  await page.screenshot({path:path.join(evidence,'uploaded-product-390.png'),fullPage:true});
 }catch(error){checks.push({check:'Synthetic image browser acceptance',status:'FAIL',error:String(error.message).replaceAll(user.password,'[redacted]').slice(0,300)});}
 finally{
  record('No upload browser runtime exceptions',errors.length,0);
  await context.request.post(base+'/auth/logout',{headers,data:{}});await browser.close();
 }
 fs.writeFileSync(path.join(evidence,'upload-browser.json'),JSON.stringify({checks,errors},null,2),{mode:0o600});
 for(const row of checks)console.log(row.status+' '+row.check);
 process.exitCode=checks.some(row=>row.status==='FAIL')?1:0;
})().catch(error=>{console.error(String(error.message).replaceAll(user.password,'[redacted]'));process.exitCode=1;});
