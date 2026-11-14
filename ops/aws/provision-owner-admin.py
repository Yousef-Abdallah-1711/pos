#!/usr/bin/env python3
"""Provision the disabled staging seed administrator through local database access."""
import crypt
import json
import os
from pathlib import Path
import secrets
import subprocess

ROOT = Path(__file__).resolve().parents[2]
PRIVATE = ROOT / 'ops/aws/evidence/owner-admin.json'
EMAIL = 'superadmin@pos-saas.com'
if PRIVATE.exists() and json.loads(PRIVATE.read_text()).get('provisioned'):
    raise SystemExit('Owner credential file already exists; preserved. Verify or rotate through a reviewed maintenance procedure.')
credential = {'email': EMAIL, 'password': secrets.token_urlsafe(36), 'rotation_required_before_production': True, 'provisioned': False}
if PRIVATE.exists():
    credential = json.loads(PRIVATE.read_text())
else:
    fd = os.open(PRIVATE, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
    with os.fdopen(fd, 'w') as stream:
        json.dump(credential, stream)
hashed = crypt.crypt(credential['password'], crypt.mksalt(crypt.METHOD_BLOWFISH, rounds=4096))
# Only the intentionally disabled, already privileged seed may be provisioned.
# No runtime endpoint, tenant role conversion, or existing password overwrite exists.
sql = """BEGIN;
UPDATE users SET password_hash = :'hash', updated_at = CURRENT_TIMESTAMP
WHERE email = 'superadmin@pos-saas.com' AND password_hash = '!seed-required'
  AND is_super_admin AND role = 'super_admin' AND status = 'active'
  AND EXISTS (SELECT 1 FROM super_admins s WHERE s.user_id = users.id AND s.is_active)
RETURNING id AS provisioned_id \\gset
INSERT INTO super_admin_audit_logs(super_admin_id,action,target_type,details,ip_address)
SELECT id,'UPDATE_SETTINGS','SYSTEM','{"reason":"authorized staging owner bootstrap; rotate before production"}'::jsonb,'127.0.0.1'
FROM super_admins WHERE user_id = :provisioned_id;
COMMIT;
"""
result = subprocess.run([str(ROOT / 'ops/aws/servio'), 'exec', '-T', 'postgres', 'psql', '-U', 'postgres', '-d', 'servio_staging', '-v', 'ON_ERROR_STOP=1', '-v', 'hash=' + hashed], input=sql, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
if result.returncode:
    raise SystemExit('Bootstrap refused or failed; private credential retained with provisioned=false. Existing accounts were not overwritten.')
credential['provisioned'] = True
PRIVATE.write_text(json.dumps(credential, indent=2) + '\n')
PRIVATE.chmod(0o600)
print('Owner staging administrator provisioned: ' + EMAIL + '; credential remains private and requires rotation before production.')
