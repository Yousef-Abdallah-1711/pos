#!/usr/bin/env python3
import http.cookiejar,json,urllib.request,urllib.error
from pathlib import Path
jar=http.cookiejar.CookieJar();op=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar));rows=[]
fixture=json.loads((Path(__file__).resolve().parents[2]/'ops/aws/evidence/browser-admin.json').read_text())
for label,method,path,data,expected in [
 ('Super-admin login','POST','/auth/login',fixture,200),
 ('Super-admin session','GET','/auth/session',None,200),
 ('Super-admin dashboard','GET','/dashboard/stats',None,200),
 ('Super-admin tenant list','GET','/tenants',None,200),
 ('Super-admin audit log','GET','/audit-logs',None,200),
 ('Super-admin logout','POST','/auth/logout',{},200),
 ('Super-admin session revocation','GET','/auth/session',None,401),
]:
    request=urllib.request.Request('https://api.pos.rakmyat.online/api/v1/super-admin'+path,data=None if data is None else json.dumps(data).encode(),headers={'Origin':'https://admin.pos.rakmyat.online','Content-Type':'application/json'},method=method)
    try:response=op.open(request,timeout=20)
    except urllib.error.HTTPError as e:response=e
    response.read();code=response.status
    rows.append(dict(check=label,actual=code,expected=expected,status='PASS' if code==expected else 'FAIL'))
(Path(__file__).resolve().parents[2]/'ops/aws/evidence/admin-smoke.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows:print(row['status']+' '+row['check'])
raise SystemExit(any(r['status']=='FAIL' for r in rows))
