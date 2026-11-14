#!/usr/bin/env python3
"""Provision a private staging-only gateway without revealing its credential."""
import json,secrets,subprocess
from pathlib import Path
root=Path('/opt/servio/ops/aws');private=root/'evidence/staging-gateway.json'
if private.exists():credential=json.loads(private.read_text())
else:
    credential={'username':'staging','password':secrets.token_urlsafe(32)}
    private.write_text(json.dumps(credential));private.chmod(0o600)
result=subprocess.run(['sudo','docker','run','--rm','-i','--network','none','--memory','64m','--security-opt','no-new-privileges:true','--cap-drop','ALL','caddy@sha256:78bd35025c7ba85961dd149ae670315129954e77abe2383ac47c89106b35ade9','caddy','hash-password'],input=credential['password']+'\n',text=True,capture_output=True,check=True)
hashed=result.stdout.strip();assert hashed.startswith('$2')
env=root/'.env';lines=[line for line in env.read_text().splitlines() if not line.startswith('STAGING_GATEWAY_PASSWORD_HASH=')];lines.append("STAGING_GATEWAY_PASSWORD_HASH='"+hashed+"'");env.write_text('\n'.join(lines)+'\n');env.chmod(0o600)
p=root/'compose.json';d=json.loads(p.read_text());d['services']['proxy']['environment']={'STAGING_GATEWAY_PASSWORD_HASH':'${STAGING_GATEWAY_PASSWORD_HASH:?provision staging gateway}'};p.write_text(json.dumps(d,indent=2)+'\n')
p=root/'Caddyfile';s=p.read_text();marker='(common) {'
if '@staging_frontend' not in s:
    s=s.replace(marker,marker+'''\n\t@staging_frontend {
\t\tnot host api.pos.rakmyat.online
\t}
\tbasic_auth @staging_frontend {
\t\tstaging {$STAGING_GATEWAY_PASSWORD_HASH}
\t}
\t@image_optimizer path /_next/image /_next/image/*
\trespond @image_optimizer "Image optimization is unavailable on staging" 503
''')
p.write_text(s);(root/'Caddyfile.https-ready').write_text(s)
print('Gateway provisioned; credential remains in a mode-0600 private file.')
