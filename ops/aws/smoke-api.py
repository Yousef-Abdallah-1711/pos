#!/usr/bin/env python3
"""Non-destructive checks using unique, synthetic staging tenants only."""
import http.cookiejar, json, secrets, time, urllib.request, urllib.error
from pathlib import Path
BASE='https://api.pos.rakmyat.online/api/v1'
ORIGIN='https://pos.rakmyat.online'
results=[]
fixtures=[]
def client():
    jar=http.cookiejar.CookieJar()
    return urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar)),jar

def call(opener,method,path,payload=None,headers=None):
    hdr={'Origin':ORIGIN,'Content-Type':'application/json',**(headers or {})}
    req=urllib.request.Request(BASE+path,data=None if payload is None else json.dumps(payload).encode(),headers=hdr,method=method)
    try:
        resp=opener.open(req,timeout=20)
    except urllib.error.HTTPError as exc:
        resp=exc
    raw=resp.read()
    try:
        body=json.loads(raw)
    except (ValueError,UnicodeDecodeError):
        body={}
    return resp.status,body,resp.headers

def check(name,actual,expected):
    expected=expected if isinstance(expected,list) else [expected]
    results.append({'check':name,'actual':actual,'expected':expected,'status':'PASS' if actual in expected else 'FAIL'})

def run():
    anon,_=client()
    code,body,hdr=call(anon,'GET','/health');check('API health',code,200)
    code,_,_=call(anon,'GET','/me');check('Anonymous protected API denied',code,401)
    code,_,_=call(anon,'POST','/auth/login',{}, {'Origin':'https://untrusted.invalid'});check('Untrusted login origin denied',code,403)
    code,_,hdr=call(anon,'OPTIONS','/auth/login',headers={'Access-Control-Request-Method':'POST'});check('Allowed CORS preflight',code,200);check('Exact allowed CORS origin',hdr.get('Access-Control-Allow-Origin'),ORIGIN)
    code,_,_=call(anon,'OPTIONS','/auth/login',headers={'Origin':'https://untrusted.invalid','Access-Control-Request-Method':'POST'});check('Untrusted CORS denied',code,403)
    stamp=str(int(time.time()))
    tenants=[]
    for suffix in ['a','b']:
        op,jar=client();password=secrets.token_urlsafe(28)
        payload={'name':'AWS Staging Test '+suffix,'email':f'aws-{stamp}-{suffix}@example.invalid','password':password,'restaurant_name':'AWS Staging Test '+stamp+suffix,'restaurant_slug':'aws-'+stamp+suffix,'phone':''}
        code,body,hdr=call(op,'POST','/auth/register',payload);check('Register synthetic tenant '+suffix,code,201)
        if code!=201: continue
        fixtures.append({'email':payload['email'],'password':password,'slug':payload['restaurant_slug']})
        check('Registration session established '+suffix,body.get('session_established'),True)
        cookie=hdr.get('Set-Cookie','');check('HttpOnly session cookie '+suffix,'httponly' in cookie.lower(),True);check('Secure session cookie '+suffix,'secure' in cookie.lower(),True)
        code,session,_=call(op,'GET','/auth/session');check('Session bootstrap '+suffix,code,200)
        user=session.get('data',session)
        code,_,_=call(op,'POST','/auth/login',{'email':payload['email'],'password':password});check('Password login '+suffix,code,200)
        tenants.append((op,jar,user))
    if len(tenants)!=2:return
    a,jar,ua=tenants[0];b,_,ub=tenants[1]
    code,body,_=call(a,'GET','/auth/session',headers={'X-Tenant-ID':str(ub.get('tenant_id')),'X-Restaurant-ID':str(ub.get('restaurant_id'))})
    check('Forged tenant headers cannot switch identity',body.get('data',body).get('tenant_id'),ua.get('tenant_id'))
    code,_,_=call(a,'GET','/super-admin/tenants');check('Tenant owner cannot access super-admin',code,[401,403])
    code,body,_=call(a,'POST','/categories',{'name':'AWS synthetic isolation '+stamp,'is_active':True});check('Create own category',code,201)
    category=body.get('data',body)
    if code==201 and category.get('id'):
        path='/categories/'+str(category['id'])
        code,_,_=call(b,'GET',path);check('Cross-tenant category read denied',code,[403,404])
        code,_,_=call(b,'PUT',path,{'name':'forbidden'});check('Cross-tenant category update denied',code,[403,404])
        code,_,_=call(b,'DELETE',path);check('Cross-tenant category delete denied',code,[403,404])
        code,after,_=call(a,'GET',path);check('Original category remains accessible',code,200);check('Original category unchanged',after.get('data',after).get('name'),category.get('name'))
    code,_,_=call(a,'POST','/categories',{'name':'untrusted'}, {'Origin':'https://untrusted.invalid'});check('Cookie mutation rejects untrusted Origin',code,403)
    code,_,_=call(a,'GET','/mock/paymob/iframe');check('Public mock payments blocked',code,503)
    code,_,_=call(a,'POST','/auth/logout',{});check('Logout',code,[200,204])
    code,_,_=call(a,'GET','/auth/session');check('Session revoked after logout',code,401)

try:
    run()
except Exception as exc:
    results.append({'check':'smoke execution','status':'FAIL','error':str(exc)})
fixture_path=(Path(__file__).resolve().parents[2]/'ops/aws/evidence/browser-fixtures.json')
fixture_path.write_text(json.dumps(fixtures))
fixture_path.chmod(0o600)
(Path(__file__).resolve().parents[2]/'ops/aws/evidence/api-smoke.json').write_text(json.dumps(results,indent=2)+'\n')
for row in results:print(row['status']+' '+row['check'])
raise SystemExit(any(r['status']=='FAIL' for r in results))
