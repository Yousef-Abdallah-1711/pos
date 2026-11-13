// Authenticated synthetic journals only; no provider payments or real business tenant.
const {request}=require('playwright');const fs=require('fs');const path=require('path');
const evidence=process.env.SERVIO_BROWSER_EVIDENCE||path.join(__dirname,'evidence');
const users=process.env.SERVIO_UPLOAD_USERS||path.join(__dirname,'evidence/browser-fixtures.json');
const user=JSON.parse(fs.readFileSync(users))[0];
if(!user.email.endsWith('@example.invalid'))throw new Error('Synthetic accounting fixture required');
const base='https://api.pos.rakmyat.online/api/v1';const headers={Origin:'https://pos.app.rakmyat.online'};const checks=[];
function record(check,actual,expected){checks.push({check,actual,expected,status:actual===expected?'PASS':'FAIL'});}
async function object(response){const raw=await response.json();return raw.data||raw;}
async function accounts(context){
 const rows=[];const stamp=Date.now().toString();
 for(const [index,type] of ['asset','equity'].entries()){
  const response=await context.post(base+'/accounting/accounts',{headers,data:{account_code:'QA'+stamp+index,account_name:'Synthetic SERVIO acceptance '+type,account_type:type,normal_balance:type==='asset'?'debit':'credit',currency:'EGP',is_active:true}});
  record('Synthetic '+type+' account creation',response.status(),201);
  if(response.status()!==201)throw new Error('Synthetic account creation failed');
  rows.push(await object(response));
 }
 return rows;
}
async function ledger(context,rows,entry){
 const responses=await Promise.all(rows.map(account=>context.get(base+'/accounting/ledger?account_id='+account.id+'&limit=100')));
 if(responses.some(response=>response.status()!==200))throw new Error('Scoped ledger retrieval failed');
 const data=await Promise.all(responses.map(object));
 return {count:data.flatMap(item=>item.lines||[]).filter(line=>line.journal_entry_id===entry).length,balanced:data.every(item=>Number(item.balance)===0)};
}
async function journal(context,rows,state){
 const payload={entry_date:new Date().toISOString().slice(0,10),description:'Synthetic staging acceptance. Not a business transaction.',entry_type:'manual',reference_type:'staging_acceptance',notes:'Synthetic tenant only; reverse through audited void.',lines:[{account_id:rows[0].id,debit_amount:1,credit_amount:0,line_order:1},{account_id:rows[1].id,debit_amount:0,credit_amount:1,line_order:2}]};
 const create=await context.post(base+'/accounting/journal-entries',{headers,data:payload});record('Balanced synthetic journal draft',create.status(),201);
 if(create.status()!==201)throw new Error('Synthetic draft creation failed');
 const draft=await object(create);state.id=draft.id;
 record('Draft totals match its lines',Number(draft.total_debit)===1&&Number(draft.total_credit)===1&&draft.status==='draft',true);
 const posted=await context.post(base+'/accounting/journal-entries/'+state.id+'/post',{headers,data:{}});record('Authorized journal posting',posted.status(),200);state.posted=posted.status()===200;
 if(!state.posted)throw new Error('Synthetic posting failed');
 record('Duplicate posting rejected',(await context.post(base+'/accounting/journal-entries/'+state.id+'/post',{headers,data:{}})).status(),400);
 record('One ledger effect per journal line',(await ledger(context,rows,state.id)).count,2);
 record('Posted journal edit rejected',(await context.put(base+'/accounting/journal-entries/'+state.id,{headers,data:{...payload,description:'Invalid historical edit'}})).status(),400);
 const reversed=await context.post(base+'/accounting/journal-entries/'+state.id+'/void',{headers,data:{reason:'Audited reversal of synthetic staging acceptance'}});record('Audited journal reversal',reversed.status(),200);state.voided=reversed.status()===200;
 record('Duplicate reversal rejected',(await context.post(base+'/accounting/journal-entries/'+state.id+'/void',{headers,data:{reason:'Duplicate synthetic retry'}})).status(),400);
 const final=await ledger(context,rows,state.id);record('Original and inverse ledger effects retained',final.count,4);record('Synthetic account balances return to zero',final.balanced,true);
}
(async()=>{
 const context=await request.newContext({ignoreHTTPSErrors:process.env.SERVIO_DISPOSABLE_BROWSER==='1'});const state={id:null,posted:false,voided:false};
 try{
  const login=await context.post(base+'/auth/login',{headers,data:{email:user.email,password:user.password}});
  if(login.status()!==200)throw new Error('Synthetic accounting login failed');
  const session=await object(await context.get(base+'/auth/session'));
  if(session.email!==user.email||session.tenant_id<=1)throw new Error('Synthetic accounting tenant scope mismatch');
  await journal(context,await accounts(context),state);
 }catch(error){checks.push({check:'Authenticated journal API acceptance',status:'FAIL',error:String(error.message).replaceAll(user.password,'[redacted]').slice(0,300)});}
 finally{
  if(state.posted&&!state.voided){const result=await context.post(base+'/accounting/journal-entries/'+state.id+'/void',{headers,data:{reason:'Emergency audited cleanup of synthetic acceptance'}});record('Synthetic journal cleanup reversal',result.status(),200);}
  await context.post(base+'/auth/logout',{headers,data:{}});
  record('Logged-out journal access denied',(await context.get(base+'/accounting/journal-entries')).status(),401);await context.dispose();
 }
 fs.writeFileSync(path.join(evidence,'journal-api-acceptance.json'),JSON.stringify({checks},null,2),{mode:0o600});
 for(const row of checks)console.log(row.status+' '+row.check);
 process.exitCode=checks.some(row=>row.status==='FAIL')?1:0;
})().catch(error=>{console.error(String(error.message).replaceAll(user.password,'[redacted]'));process.exitCode=1;});
