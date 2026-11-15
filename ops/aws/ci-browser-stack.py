#!/usr/bin/env python3
"""Verify built images in a separate, synthetic GitHub Actions Compose stack."""
import crypt
import json
import os
from pathlib import Path
import re
import secrets
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
hosts = ['pos.servio.test', 'api.servio.test', 'pos.app.servio.test', 'admin.pos.servio.test', 'restaurant-1.pos.servio.test']
host_replacements = {
    'https://api.pos.rakmyat.online': 'https://api.servio.test',
    'https://pos.app.rakmyat.online': 'https://pos.app.servio.test',
    'https://admin.pos.rakmyat.online': 'https://admin.pos.servio.test',
    'https://restaurant-1.pos.rakmyat.online': 'https://restaurant-1.pos.servio.test',
    'https://pos.rakmyat.online': 'https://pos.servio.test',
    'api.pos.rakmyat.online': 'api.servio.test',
    'pos.app.rakmyat.online': 'pos.app.servio.test',
    'admin.pos.rakmyat.online': 'admin.pos.servio.test',
    'restaurant-1.pos.rakmyat.online': 'restaurant-1.pos.servio.test',
    'pos.rakmyat.online': 'pos.servio.test',
}
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
        service.pop('ports', None)
        service['networks']['edge'] = {'aliases': hosts}
    if name == 'backend':
        environment = service.setdefault('environment', {})
        environment['ENV'] = 'test'
        for key, value in list(environment.items()):
            if isinstance(value, str):
                for old, new in host_replacements.items():
                    value = value.replace(old, new)
                environment[key] = value
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
config['networks']['edge']['internal'] = True
# CI uses its own local CA only; this does not modify or relax public HTTPS.
caddy = (OPS / 'Caddyfile.public').read_text()
caddy = caddy.replace('rakmyat.online', 'servio.test')
# CI's API TLS route uses the same short host as the browser target guard.
caddy = caddy.replace('api.pos.servio.test', 'api.servio.test')
caddy = re.sub(r'\s*@staging_frontend\s*\{[^}]*\}\s*basic_auth @staging_frontend\s*\{[^\n]*\n\s*staging [^\n]*\n\s*\}', '', caddy)
(PRIVATE / 'ci-Caddyfile').write_text('{\n local_certs\n}\n' + caddy)
path = PRIVATE / 'ci-compose.json'
path.write_text(json.dumps(config, indent=2))
compose = ['docker', 'compose', '--env-file', str(OPS / '.env'), '-f', str(path)]

def sql(text):
    return subprocess.run(compose + ['exec', '-T', 'postgres', 'psql', '-U', 'postgres', '-d', 'servio_ci_test', '-qAt', '-v', 'ON_ERROR_STOP=1'], input=text, text=True, stdout=subprocess.PIPE, check=True).stdout

def verify_disposable_stack():
    for name in ['postgres', 'redis', 'backend', 'proxy']:
        container = json.loads(subprocess.check_output(['docker', 'inspect', f'servio-ci-{name}-1']))[0]
        if container['Config']['Labels'].get('com.servio.disposable') != 'true':
            raise RuntimeError(f'{name} is not labeled as disposable')
        if container['HostConfig'].get('PortBindings'):
            raise RuntimeError(f'{name} unexpectedly publishes a host port')
        if name == 'postgres' and 'POSTGRES_DB=servio_ci_test' not in container['Config']['Env']:
            raise RuntimeError('Disposable PostgreSQL database identity mismatch')
        if name == 'backend':
            environment = container['Config']['Env']
            if 'ENV=test' not in environment or 'DB_NAME=servio_ci_test' not in environment:
                raise RuntimeError('API database/environment identity mismatch')
            if 'PAYMOB_MODE=mock' not in environment:
                raise RuntimeError('External payment integration is not sandboxed')
            if any(key.startswith(('PAYMOB_API_KEY=', 'SMTP_PASSWORD=')) for key in environment):
                raise RuntimeError('External production integration credentials are present')
    for name in ['servio-ci_private', 'servio-ci_app', 'servio-ci_edge']:
        network = json.loads(subprocess.check_output(['docker', 'network', 'inspect', name]))[0]
        if network.get('Internal') is not True:
            raise RuntimeError(f'{name} permits external egress')

def register(suffix):
    email = 'ci-responsive-overflow-' + suffix + '-' + ('x' * 72) + '@example.invalid' if suffix == 'a' else 'ci-' + suffix + '@example.invalid'
    name = 'CI responsive overflow ' + ('x' * 72) if suffix == 'a' else 'CI synthetic ' + suffix
    user = {'name': name, 'email': email, 'password': secrets.token_urlsafe(32), 'restaurant_name': 'CI ' + suffix, 'restaurant_slug': 'ci-' + suffix, 'phone': ''}
    # Seed through backend-local loopback; the GitHub runner has no app-network route.
    result = subprocess.run(
        compose + ['exec', '-T', 'backend', 'curl', '--silent', '--show-error',
                   '--max-time', '30', '--output', '-', '--write-out', '\n%{http_code}',
                   '--request', 'POST', '--header', 'Host: api.servio.test',
                   '--header', 'Origin: https://pos.servio.test',
                   '--header', 'Content-Type: application/json', '--data-binary', '@-',
                   'http://127.0.0.1:8080/api/v1/auth/register'],
        input=json.dumps(user), text=True, stdout=subprocess.PIPE, check=True,
    )
    status = result.stdout.rsplit('\n', 1)[-1]
    if status != '201':
        raise RuntimeError('Disposable API registration failed: HTTP ' + status)
    return user

try:
    subprocess.run(compose + ['up', '-d', '--wait', '--wait-timeout', '180'], check=True)
    verify_disposable_stack()
    users = [register('a'), register('b')]
    (PRIVATE / 'browser-fixtures.json').write_text(json.dumps(users))
    admin = {'email': 'superadmin@pos-saas.com', 'password': secrets.token_urlsafe(32)}
    hashed = crypt.crypt(admin['password'], crypt.mksalt(crypt.METHOD_BLOWFISH))
    sql("UPDATE users SET password_hash = '" + hashed + "' WHERE email = 'superadmin@pos-saas.com' AND password_hash = '!seed-required';")
    storefront_sql = (OPS / 'align-seeded-storefront.sql').read_text().replace('rakmyat.online', 'servio.test')
    sql(storefront_sql)
    fixture=sql(runpy.run_path(str(OPS/'seed-storefront-qa.py'))['SQL'])
    (PRIVATE/'storefront-qa-fixture.json').write_text(fixture)
    (PRIVATE/'storefront-qa-fixture.json').chmod(0o600)
    (PRIVATE / 'browser-admin.json').write_text(json.dumps(admin))
    (PRIVATE / 'staging-gateway.json').write_text(json.dumps({'username': 'unused-ci', 'password': secrets.token_urlsafe(32)}))
    for p in ['browser-fixtures.json', 'browser-admin.json', 'staging-gateway.json']:
        (PRIVATE / p).chmod(0o600)
    output = PRIVATE / 'browser'
    output.mkdir(exist_ok=True)
    metadata = {'commit': SHA, 'run_id': os.environ.get('GITHUB_RUN_ID'), 'environment': 'isolated-github-actions', 'database': 'servio_ci_test', 'api_target': 'https://api.servio.test/api/v1'}
    (output / 'run-metadata.json').write_text(json.dumps(metadata, indent=2) + '\n')
    image = 'mcr.microsoft.com/playwright@sha256:ee15712262f017b8157421b212b12099146645f44b1bc5c0514cbbe6a229b1d0'
    browser_modules = os.environ.get('SERVIO_CI_BROWSER_NODE_MODULES')
    if not browser_modules or not (Path(browser_modules) / 'playwright').exists():
        raise RuntimeError('Install CI browser dependencies on the runner before starting the isolated stack')
    docker_args = [
        'docker', 'run', '--rm', '--user', f'{os.getuid()}:{os.getgid()}',
        '--security-opt', 'no-new-privileges:true', '--cap-drop', 'ALL',
        '--network', 'servio-ci_edge', '--memory', '1500m', '--shm-size', '256m',
        '-e', 'SERVIO_ACCEPTANCE_ENV=ci', '-e', 'SERVIO_DISPOSABLE_BROWSER=1',
        '-e', 'SERVIO_ACCEPTANCE_LANDING_URL=https://pos.servio.test',
        '-e', 'SERVIO_ACCEPTANCE_DASHBOARD_URL=https://pos.app.servio.test',
        '-e', 'SERVIO_ACCEPTANCE_ADMIN_URL=https://admin.pos.servio.test',
        '-e', 'SERVIO_ACCEPTANCE_STOREFRONT_URL=https://restaurant-1.pos.servio.test',
        '-e', 'SERVIO_ACCEPTANCE_API_BASE_URL=https://api.servio.test/api/v1',
        '-e', 'SERVIO_MARKETING_URL=https://pos.servio.test',
        '-e', 'POOS_ACCEPTANCE_URL=https://pos.servio.test',
        '-e', 'POOS_ACCEPTANCE_EVIDENCE=/evidence/landing-regression',
        '-e', 'SERVIO_MARKETING_API_URL=https://api.servio.test/api/v1',
        '-e', 'SERVIO_SYNTHETIC_TRIAL=1', '-e', 'SERVIO_BROWSER_EVIDENCE=/evidence',
        '-e', 'SERVIO_STORE_FIXTURE=/fixtures/storefront.json',
        '-e', 'SERVIO_UPLOAD_USERS=/fixtures/users.json',
        '-v', f'{browser_modules}:/tmp/browser/node_modules:ro',
        '-v', f'{OPS / "run-ci-browser-checks.sh"}:/work/run-ci-browser-checks.sh:ro',
        '-v', f'{ROOT / "frontend/apps/landing/scripts/capture-landing-acceptance.mjs"}:/work/capture-landing-acceptance.mjs:ro',
        '-v', f'{ROOT / "frontend/apps/landing/scripts/lighthouse-acceptance.mjs"}:/work/lighthouse-acceptance.mjs:ro',
        '-v', f'{ROOT / "frontend/apps/landing/scripts/trial-registration.mjs"}:/work/trial-registration.mjs:ro',
        '-v', f'{ROOT / "frontend/apps/landing/scripts/marketing-acceptance.mjs"}:/work/marketing-acceptance.mjs:ro',
        '-v', f'{OPS / "browser-upload.cjs"}:/work/browser-upload.cjs:ro',
        '-v', f'{OPS / "journal-api-acceptance.cjs"}:/work/journal-api-acceptance.cjs:ro',
        '-v', f'{OPS / "acceptance-target.cjs"}:/work/acceptance-target.cjs:ro',
        '-v', f'{OPS / "browser-auth-validation.cjs"}:/work/browser-auth-validation.cjs:ro',
        '-v', f'{OPS / "browser-storefront-flow.cjs"}:/work/storefront-flow.cjs:ro',
        '-v', f'{OPS / "browser-smoke.cjs"}:/work/browser-smoke.cjs:ro',
        '-v', f'{PRIVATE / "storefront-qa-fixture.json"}:/fixtures/storefront.json:ro',
        '-v', f'{PRIVATE / "browser-fixtures.json"}:/fixtures/users.json:ro',
        '-v', f'{PRIVATE / "browser-admin.json"}:/fixtures/admin.json:ro',
        '-v', f'{output}:/evidence', image, 'sh', '/work/run-ci-browser-checks.sh',
    ]
    subprocess.run(docker_args, check=True)
finally:
    subprocess.run(compose + ['stop'], check=True)
