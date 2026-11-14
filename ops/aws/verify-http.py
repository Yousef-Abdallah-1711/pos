#!/usr/bin/env python3
"""Check the approved public routes and redirects without recording credentials."""
import base64,json,urllib.request,urllib.error
from pathlib import Path
class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self,req,fp,code,msg,headers,newurl):return None
op=urllib.request.build_opener(NoRedirect);rows=[]
gateway=json.loads(Path("/opt/servio/ops/aws/evidence/staging-gateway.json").read_text())
authorization="Basic "+base64.b64encode((gateway["username"]+":"+gateway["password"]).encode()).decode()
for label,host,path in [
 ('landing','pos.rakmyat.online','/en'),
 ('api','api.pos.rakmyat.online','/api/v1/health'),
 ('dashboard','pos.app.rakmyat.online','/en/dashboard'),
 ('super-admin','admin.pos.rakmyat.online','/login'),
 ('storefront','restaurant-1.pos.rakmyat.online','/en')]:
    for scheme,expected in [('https',200),('http',308)]:
        url=scheme+'://'+host+path
        try:
            try:
                request=urllib.request.Request(url)
                if scheme=="https" and label!="api":request.add_header("Authorization",authorization)
                r=op.open(request,timeout=15)
            except urllib.error.HTTPError as e:r=e
            location=r.headers.get('Location');status=r.status;r.close()
            correct=status==expected and (scheme=='https' or location=='https://'+host+path)
            rows.append(dict(service=label,url=url,http=status,location=location,status='PASS' if correct else 'FAIL'))
        except Exception as e:rows.append(dict(service=label,url=url,status='FAIL',error=str(e)))
Path('/opt/servio/ops/aws/evidence/http-checks.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows:print(row['status']+' '+row['url'])
raise SystemExit(any(r['status']=='FAIL' for r in rows))
