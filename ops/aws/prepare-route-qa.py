#!/usr/bin/env python3
"""Create bounded, authenticated resources only in the existing synthetic QA tenant."""
import http.cookiejar
import json
from pathlib import Path
import time
import urllib.error
import urllib.request

ROOT=Path(__file__).resolve().parents[2]
BASE='https://api.pos.rakmyat.online/api/v1'

class Client:
    def __init__(self):
        self.opener=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
    def call(self,method,path,data=None):
        request=urllib.request.Request(BASE+path,data=None if data is None else json.dumps(data).encode(),method=method,headers={'Origin':'https://pos.app.rakmyat.online','Content-Type':'application/json'})
        try:response=self.opener.open(request,timeout=20)
        except urllib.error.HTTPError as error:response=error
        raw=response.read()
        if response.status>=400:raise RuntimeError(method+' '+path+' HTTP '+str(response.status))
        body=json.loads(raw) if raw else {}
        return body.get('data',body) if isinstance(body,dict) else body

def collection(value,key):
    if isinstance(value,list):return value
    if isinstance(value,dict):return value.get(key,[])
    raise RuntimeError('Unexpected scoped collection response: '+key)

def journal_fixture(client,tenant):
    entries=collection(client.call('GET','/accounting/journal-entries'),'entries')
    if entries:return entries[0]['id']
    accounts=collection(client.call('GET','/accounting/accounts'),'accounts')
    accounts=[a for a in accounts if not a.get('is_header') and a.get('is_active',True)]
    for suffix,kind,balance in [('A','asset','debit'),('E','equity','credit')]:
        if len(accounts)>=2:break
        account=client.call('POST','/accounting/accounts',{'account_code':'QA'+str(tenant)+suffix,'account_name':'SERVIO synthetic QA '+kind,'account_type':kind,'normal_balance':balance,'currency':'EGP','is_active':True})
        accounts.append(account)
    entry=client.call('POST','/accounting/journal-entries',{'entry_date':time.strftime('%Y-%m-%d',time.gmtime()),'description':'SERVIO synthetic QA draft','entry_type':'manual','reference_type':'manual','notes':'Synthetic tenant only; draft, no real payment','lines':[{'account_id':accounts[0]['id'],'debit_amount':1,'credit_amount':0},{'account_id':accounts[1]['id'],'debit_amount':0,'credit_amount':1}]})
    return entry['id']

def prepare():
    fixture=json.loads((ROOT/'ops/aws/evidence/browser-fixtures.json').read_text())[0]
    if not fixture['email'].endswith('@example.invalid'):raise RuntimeError('Synthetic account required')
    client=Client();client.call('POST','/auth/login',{'email':fixture['email'],'password':fixture['password']})
    session=client.call('GET','/auth/session')
    if session.get('email')!=fixture['email'] or session.get('tenant_id',0)<=1:raise RuntimeError('Synthetic tenant identity mismatch')
    tenant=session['tenant_id']
    result={'tenant_id':tenant,'restaurant_id':session['restaurant_id'],'synthetic':True}
    for resource,key in [('products','products'),('categories','categories'),('rbac/roles','roles'),('admin/themes','themes'),('admin/orders?limit=20','orders')]:
        items=collection(client.call('GET','/'+resource),key)
        if items:result[key]=items[0]['id']
    tables=collection(client.call('GET','/tables'),'tables')
    result['tables']=(tables[0] if tables else client.call('POST','/tables',{'table_number':901,'name':'SERVIO synthetic QA','capacity':4,'location':'indoor','shape':'square'}))['id']
    drivers=collection(client.call('GET','/admin/drivers'),'drivers')
    driver=(drivers[0] if drivers else client.call('POST','/admin/drivers',{'first_name':'SERVIO','last_name':'Synthetic QA','email':'qa-driver-'+str(tenant)+'@example.invalid','phone_number':'+12025550123','license_number':'QA-'+str(tenant),'vehicle_type':'motorcycle'}))
    result['drivers']=driver.get('driver',driver)['id']
    result['journal_entries']=journal_fixture(client,tenant)
    client.call('POST','/auth/logout',{})
    target=ROOT/'ops/aws/evidence/route-qa-fixtures.json';target.write_text(json.dumps(result,indent=2)+'\n');target.chmod(0o600)
    print('PASS authenticated synthetic resource preparation; no real payments or privileged role changes')

if __name__=='__main__':prepare()
