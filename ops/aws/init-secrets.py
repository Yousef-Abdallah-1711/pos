#!/usr/bin/env python3
"""Create private per-environment secrets without replacing existing credentials."""
import os,secrets,sys
from pathlib import Path
root=Path(__file__).resolve().parent
for name in ['evidence','backups']:
    (root/name).mkdir(exist_ok=True)
    (root/name).chmod(0o700)
acceptance='--acceptance' in sys.argv[1:]
p=root/'evidence/acceptance.env' if acceptance else root/'.env'
keys=['TEST_POSTGRES_PASSWORD','TEST_APP_PASSWORD','TEST_MIGRATOR_PASSWORD','TEST_REDIS_PASSWORD'] if acceptance else ['DB_BOOTSTRAP_PASSWORD','DB_PASSWORD','MIGRATE_DB_PASSWORD','JWT_SECRET','REDIS_PASSWORD']
if p.exists():
    p.chmod(0o600)
    print('Existing private environment preserved.')
else:
    fd=os.open(p,os.O_CREAT|os.O_EXCL|os.O_WRONLY,0o600)
    with os.fdopen(fd,'w') as stream:
        for key in keys:
            stream.write(key+'='+secrets.token_hex(32)+'\n')
    print('Private environment created; no values printed.')
