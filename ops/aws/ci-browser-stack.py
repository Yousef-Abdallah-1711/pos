#!/usr/bin/env python3
"""Verify built images in a separate, synthetic GitHub Actions Compose stack."""
import crypt
import http.client
import json
import os
from pathlib import Path
import re
import secrets
import ssl
import socket
import subprocess
import runpy

ROOT = Path(__file__).resolve().parents[2]
OPS = ROOT / 'ops/aws'
PRIVATE = OPS / 'evidence'
SHA = os.environ.get('GITHUB_SHA', '')
if os.environ.get('GITHUB_ACTIONS') != 'true' or not re.fullmatch('[0-9a-f]{40}', SHA):
    raise SystemExit('CI-only command: live EC2 execution refused')
PRIVATE.mkdir(exist_ok=True, mode=0o700)
subprocess.run(['python3', str(OPS / 'init-secrets.py')], check=True)
config = json.loads((OPS / 'compose.json').read_text())
config['name'] = 'servio-ci'
for name, service in config['services'].items():
    service.pop('build', None)
    service['labels'] = {'com.servio.disposable': 'true'}
    if name in ['backend', 'migrate', 'landing', 'dashboard', 'restaurant-website', 'super-admin']:
        service['image'] = f'servio-ci-{"backend" if name == "migrate" else name}:{SHA}'
    for key in ['POSTGRES_DB', 'DB_NAME']:
        if key in service.get('environment', {}):
            service['environment'][key] = 'servio_ci_test'
    if name == 'postgres':
        service['healthcheck']['test'] = ['CMD-SHELL', 'pg_isready -U postgres -d servio_ci_test']
    if name == 'proxy':
        service['environment'] = {}
        service['ports'] = ['127.0.0.1:80:80', '127.0.0.1:443:443']
        service['networks']['edge'] = {'aliases': ['pos.rakmyat.online', 'api.pos.rakmyat.online', 'pos.app.rakmyat.online', 'admin.pos.rakmyat.online', 'restaurant-1.pos.rakmyat.online']}
    if name in ['landing', 'dashboard', 'restaurant-website', 'super-admin']:
        service['environment']['API_INTERNAL_URL'] = 'http://backend:8080/api/v1'
        service.pop('ports', None)
    volumes = []
    for volume in service.get('volumes', []):
        source, *rest = volume.split(':')
        if source.startswith('.'):
            source = str((OPS / source).resolve())
        if name == 'proxy' and source.endswith('/Caddyfile'):
            source = str(PRIVATE / 'ci-Caddyfile')
        volumes.append(':'.join([source, *rest]))
    service['volumes'] = volumes
# CI uses its own local CA only; this does not modify or relax public HTTPS.
caddy = (OPS / 'Caddyfile.public').read_text()
caddy = re.sub(r'\s*@staging_frontend\s*\{[^}]*\}\s*basic_auth @staging_frontend\s*\{[^\n]*\n\s*staging [^\n]*\n\s*\}', '', caddy)
(PRIVATE / 'ci-Caddyfile').write_text('{\n local_certs\n}\n' + caddy)
path = PRIVATE / 'ci-compose.json'
path.write_text(json.dumps(config, indent=2))
compose = ['docker', 'compose', '--env-file', str(OPS / '.env'), '-f', str(path)]

def sql(text):
    return subprocess.run(compose + ['exec', '-T', 'postgres', 'psql', '-U', 'postgres', '-d', 'servio_ci_test', '-qAt', '-v', 'ON_ERROR_STOP=1'], input=text, text=True, stdout=subprocess.PIPE, check=True).stdout

class LoopbackTLS(http.client.HTTPSConnection):
    def connect(self):
        # Connect only to this disposable runner while sending the correct TLS SNI.
        raw=socket.create_connection(('127.0.0.1',443),timeout=self.timeout)
        self.sock=self._context.wrap_socket(raw,server_hostname='api.pos.rakmyat.online')

def register(suffix):
    user = {'name': 'CI synthetic ' + suffix, 'email': 'ci-' + suffix + '@example.invalid', 'password': secrets.token_urlsafe(32), 'restaurant_name': 'CI ' + suffix, 'restaurant_slug': 'ci-' + suffix, 'phone': ''}
    # Fixed loopback destination: no production DNS lookup or live registration.
    connection = LoopbackTLS('127.0.0.1', 443, context=ssl._create_unverified_context(), timeout=30)
    connection.request('POST', '/api/v1/auth/register', json.dumps(user), {'Host': 'api.pos.rakmyat.online', 'Origin': 'https://pos.rakmyat.online', 'Content-Type': 'application/json'})
    response = connection.getresponse()
    response.read()
    if response.status != 201:
        raise RuntimeError('CI registration failed: HTTP ' + str(response.status))
    connection.close()
    return user

try:
    subprocess.run(compose + ['up', '-d', '--wait', '--wait-timeout', '180'], check=True)
    data = json.loads(subprocess.check_output(['docker', 'inspect', 'servio-ci-postgres-1']))[0]
    if data['Config']['Labels'].get('com.servio.disposable') != 'true' or 'POSTGRES_DB=servio_ci_test' not in data['Config']['Env']:
        raise RuntimeError('Disposable database identity mismatch')
    users = [register('a'), register('b')]
    (PRIVATE / 'browser-fixtures.json').write_text(json.dumps(users))
    admin = {'email': 'superadmin@pos-saas.com', 'password': secrets.token_urlsafe(32)}
    hashed = crypt.crypt(admin['password'], crypt.mksalt(crypt.METHOD_BLOWFISH))
    sql("UPDATE users SET password_hash = '" + hashed + "' WHERE email = 'superadmin@pos-saas.com' AND password_hash = '!seed-required';")
    sql((OPS / 'align-seeded-storefront.sql').read_text())
    fixture=sql(runpy.run_path(str(OPS/'seed-storefront-qa.py'))['SQL'])
    (PRIVATE/'storefront-qa-fixture.json').write_text(fixture)
    (PRIVATE/'storefront-qa-fixture.json').chmod(0o600)
    (PRIVATE / 'browser-admin.json').write_text(json.dumps(admin))
    (PRIVATE / 'staging-gateway.json').write_text(json.dumps({'username': 'unused-ci', 'password': secrets.token_urlsafe(32)}))
    for p in ['browser-fixtures.json', 'browser-admin.json', 'staging-gateway.json']:
        (PRIVATE / p).chmod(0o600)
    output = PRIVATE / 'browser'
    output.mkdir(exist_ok=True)
    image = 'mcr.microsoft.com/playwright@sha256:ee15712262f017b8157421b212b12099146645f44b1bc5c0514cbbe6a229b1d0'
    subprocess.run(['docker', 'run', '--rm','--user',str(os.getuid())+':'+str(os.getgid()),'--security-opt','no-new-privileges:true','--cap-drop','ALL','-e','NPM_CONFIG_CACHE=/tmp/servio-npm-cache', '--network', 'servio-ci_edge', '--memory', '1500m', '--shm-size', '256m', '-e', 'SERVIO_DISPOSABLE_BROWSER=1', '-e','SERVIO_BROWSER_EVIDENCE=/evidence','-e','SERVIO_STORE_FIXTURE=/fixtures/storefront.json','-e','SERVIO_UPLOAD_USERS=/fixtures/users.json','-v',str(OPS/'browser-upload.cjs')+':/work/browser-upload.cjs:ro','-v',str(OPS/'journal-api-acceptance.cjs')+':/work/journal-api-acceptance.cjs:ro','-v',str(OPS/'browser-auth-validation.cjs')+':/work/browser-auth-validation.cjs:ro','-v',str(OPS/'browser-storefront-flow.cjs')+':/work/storefront-flow.cjs:ro','-v',str(PRIVATE/'storefront-qa-fixture.json')+':/fixtures/storefront.json:ro', '-v', str(OPS / 'browser-smoke.cjs') + ':/work/browser-smoke.cjs:ro', '-v', str(PRIVATE / 'browser-fixtures.json') + ':/fixtures/users.json:ro', '-v', str(PRIVATE / 'browser-admin.json') + ':/fixtures/admin.json:ro', '-v', str(output) + ':/evidence', image, 'sh', '-c', 'npm install --prefix /tmp/browser --ignore-scripts --no-audit --no-fund playwright@1.64.1 >/tmp/install.log 2>&1 && NODE_PATH=/tmp/browser/node_modules node /work/browser-smoke.cjs && NODE_PATH=/tmp/browser/node_modules node /work/storefront-flow.cjs && NODE_PATH=/tmp/browser/node_modules node /work/browser-upload.cjs && NODE_PATH=/tmp/browser/node_modules node /work/journal-api-acceptance.cjs && NODE_PATH=/tmp/browser/node_modules node /work/browser-auth-validation.cjs'], check=True)
finally:
    subprocess.run(compose + ['stop'], check=True)
