#!/usr/bin/env python3
import http.cookiejar,json,urllib.request,urllib.error
from pathlib import Path
from acceptance_target import api_base, safe_opener, ui_origin
BASE=api_base()+'/super-admin';ORIGIN=ui_origin('admin')
jar=http.cookiejar.CookieJar();op=safe_opener(urllib.request.HTTPCookieProcessor(jar));rows=[]
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
    request=urllib.request.Request(BASE+path,data=None if data is None else json.dumps(data).encode(),headers={'Origin':ORIGIN,'Content-Type':'application/json'},method=method)
    try:response=op.open(request,timeout=20)
    except urllib.error.HTTPError as e:response=e
    response.read();code=response.status
    rows.append(dict(check=label,actual=code,expected=expected,status='PASS' if code==expected else 'FAIL'))
(Path(__file__).resolve().parents[2]/'ops/aws/evidence/admin-smoke.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows:print(row['status']+' '+row['check'])
raise SystemExit(any(r['status']=='FAIL' for r in rows))
