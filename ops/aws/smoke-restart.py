#!/usr/bin/env python3
import http.cookiejar,json,os,sys,time,urllib.request,urllib.error
from pathlib import Path
p=Path('/opt/servio/ops/aws/evidence/restart-cookies.private')
jar=http.cookiejar.LWPCookieJar(str(p))
op=urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))
def request(method,path,body=None):
    req=urllib.request.Request('https://api.pos.rakmyat.online/api/v1'+path,data=None if body is None else json.dumps(body).encode(),headers={'Origin':'https://pos.app.rakmyat.online','Content-Type':'application/json'},method=method)
    try:r=op.open(req,timeout=10)
    except urllib.error.HTTPError as e:r=e
    r.read();return r.status
if sys.argv[1]=='before':
    user=json.loads(Path('/opt/servio/ops/aws/evidence/browser-fixtures.json').read_text())[0]
    assert request('POST','/auth/login',{'email':user['email'],'password':user['password']})==200
    assert request('GET','/auth/session')==200
    jar.save(ignore_discard=True,ignore_expires=True);p.chmod(0o600)
    print('PASS Session established before restart')
else:
    jar.load(ignore_discard=True,ignore_expires=True)
    for attempt in range(10):
        try:
            if request('GET','/health')==200:break
        except Exception:pass
        time.sleep(2)
    assert request('GET','/auth/session')==200
    assert request('POST','/auth/logout',{})==204
    assert request('GET','/auth/session')==401
    print('PASS Redis-backed session survived '+(sys.argv[2] if len(sys.argv)>2 else 'API recreation')+' and was revoked on logout')
