#!/usr/bin/env python3
import copy,http.cookiejar,json,urllib.request,urllib.error
from pathlib import Path
jar=http.cookiejar.CookieJar();op=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar));rows=[]
def request(client,method,path,payload=None):
    req=urllib.request.Request('https://api.pos.rakmyat.online/api/v1'+path,data=None if payload is None else json.dumps(payload).encode(),headers={'Origin':'https://pos.app.rakmyat.online','Content-Type':'application/json'},method=method)
    try:r=client.open(req,timeout=20)
    except urllib.error.HTTPError as e:r=e
    r.read();return r.status

def check(name,actual,expected):rows.append(dict(check=name,actual=actual,expected=expected,status='PASS' if actual==expected else 'FAIL'))
user=json.loads((Path(__file__).resolve().parents[2]/'ops/aws/evidence/browser-fixtures.json').read_text())[1]
check('Refresh fixture login',request(op,'POST','/auth/login',{'email':user['email'],'password':user['password']}),200)
old_jar=http.cookiejar.CookieJar()
for cookie in jar:old_jar.set_cookie(copy.copy(cookie))
old_values=[c.value for c in old_jar]
old_op=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(old_jar))
check('Session refresh',request(op,'POST','/auth/refresh',{}),204)
check('Session identifier rotated',[c.value for c in jar]!=old_values,True)
check('Refreshed session accepted',request(op,'GET','/auth/session'),200)
check('Previous session revoked',request(old_op,'GET','/auth/session'),401)
check('Refreshed session logout',request(op,'POST','/auth/logout',{}),204)
(Path(__file__).resolve().parents[2]/'ops/aws/evidence/refresh-smoke.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows:print(row['status']+' '+row['check'])
raise SystemExit(any(r['status']=='FAIL' for r in rows))
