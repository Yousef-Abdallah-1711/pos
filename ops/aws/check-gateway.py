#!/usr/bin/env python3
import base64,json,urllib.request,urllib.error
from pathlib import Path
root=Path('/opt/servio/ops/aws/evidence');cred=json.loads((root/'staging-gateway.json').read_text());auth='Basic '+base64.b64encode((cred['username']+':'+cred['password']).encode()).decode();rows=[]
for host in ['pos.rakmyat.online','pos.app.rakmyat.online','admin.pos.rakmyat.online','restaurant-1.pos.rakmyat.online']:
    for headers,label,path,expected in [({},'unauthenticated','/',401),({'Authorization':auth},'authorized','/en' if host!='admin.pos.rakmyat.online' else '/login',200),({'Authorization':auth},'optimizer-blocked','/_next/image?url=https://example.invalid/test.avif&w=640&q=75',503)]:
        try:r=urllib.request.urlopen(urllib.request.Request('https://'+host+path,headers=headers),timeout=15)
        except urllib.error.HTTPError as e:r=e
        rows.append(dict(host=host,check=label,actual=r.status,expected=expected,status='PASS' if r.status==expected else 'FAIL'));r.close()
r=urllib.request.urlopen('https://api.pos.rakmyat.online/api/v1/health',timeout=15);rows.append(dict(host='api.pos.rakmyat.online',check='API health',actual=r.status,expected=200,status='PASS' if r.status==200 else 'FAIL'));r.close()
(root/'gateway-checks.json').write_text(json.dumps(rows,indent=2)+'\n')
for row in rows:print(row['status'],row['host'],row['check'],row['actual'])
raise SystemExit(any(r['status']=='FAIL' for r in rows))
