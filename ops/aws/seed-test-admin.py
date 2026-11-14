#!/usr/bin/env python3
import crypt,json,secrets,subprocess,time
from pathlib import Path
root=Path('/opt/servio/ops/aws/evidence')
p=root/'browser-admin.json'
if p.exists():raise SystemExit('Synthetic admin fixture already exists; preserved.')
owner=json.loads((root/'browser-fixtures.json').read_text())[1]
admin={'email':'aws-admin-'+str(int(time.time()))+'@example.invalid','password':secrets.token_urlsafe(32)}
hashed=crypt.crypt(admin['password'],crypt.mksalt(crypt.METHOD_BLOWFISH))
sql=r"""BEGIN;
INSERT INTO users (tenant_id,restaurant_id,email,password_hash,name,role,status,is_super_admin)
SELECT tenant_id,restaurant_id,:'admin_email',:'admin_hash','AWS staging synthetic super admin','super_admin','active',true
FROM users WHERE email=:'owner_email' AND tenant_id IN (SELECT id FROM tenants WHERE slug LIKE 'aws-%') RETURNING id AS fixture_admin_id \gset
INSERT INTO super_admins (user_id,is_active) VALUES (:fixture_admin_id,true);
COMMIT;
"""
subprocess.run(['/opt/servio/ops/aws/servio','exec','-T','postgres','psql','-U','postgres','-d','servio_staging','-v','ON_ERROR_STOP=1','-v','owner_email='+owner['email'],'-v','admin_email='+admin['email'],'-v','admin_hash='+hashed],input=sql.encode(),check=True,stdout=subprocess.DEVNULL)
p.write_text(json.dumps(admin));p.chmod(0o600)
print('Synthetic admin provisioned; credentials stored privately and not printed.')
