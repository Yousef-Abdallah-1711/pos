#!/usr/bin/env python3
"""Run bounded Go checks with positively identified disposable dependencies only."""
import json
import os
from pathlib import Path
import subprocess
import sys
import time

ROOT = Path(__file__).resolve().parents[2]
EVIDENCE = ROOT / 'ops/aws/evidence'
values = dict(line.split('=', 1) for line in (EVIDENCE / 'acceptance.env').read_text().splitlines() if '=' in line)
env = dict(os.environ, GOMAXPROCS='1', GOFLAGS='-p=1', GOTOOLCHAIN='local', ENV='test')

def address(service):
    name = 'servio-acceptance-' + service + '-1'
    data = json.loads(subprocess.check_output(['sudo', 'docker', 'inspect', name]))[0]
    if data['Config']['Labels'].get('com.servio.disposable') != 'true':
        raise RuntimeError('Refusing dependencies without disposable label')
    if service == 'postgres' and 'POSTGRES_DB=servio_test' not in data['Config']['Env']:
        raise RuntimeError('Refusing unexpected database')
    return data['NetworkSettings']['Networks']['servio-acceptance_isolated']['IPAddress']

host = address('postgres')
db_name = 'servio_test'
if '--fresh' in sys.argv:
    db_name = 'servio_remediation_test_' + str(int(time.time()))
    bootstrap = f'''CREATE DATABASE {db_name} OWNER pos_migrator;
\\connect {db_name}
ALTER SCHEMA public OWNER TO pos_migrator;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO pos_app;
REVOKE ALL ON DATABASE {db_name} FROM PUBLIC;
GRANT CONNECT ON DATABASE {db_name} TO pos_app, pos_migrator;
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO pos_app;
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO pos_app;
ALTER DEFAULT PRIVILEGES FOR ROLE pos_migrator IN SCHEMA public GRANT EXECUTE ON FUNCTIONS TO pos_app;
'''
    subprocess.run(['sudo','docker','exec','-i','servio-acceptance-postgres-1','psql','-U','postgres','-v','ON_ERROR_STOP=1'],input=bootstrap,text=True,stdout=subprocess.DEVNULL,check=True)
    (EVIDENCE / 'remediation-test-database.txt').write_text(db_name + '\n')
for key in ['POS_TEST_DATABASE_URL', 'POS_RLS_ADMIN_URL']:
    env[key] = f"postgres://postgres:{values['TEST_POSTGRES_PASSWORD']}@{host}:5432/{db_name}?sslmode=disable"
for key in ['POS_APP_DATABASE_URL', 'POS_RLS_APP_URL']:
    env[key] = f"postgres://pos_app:{values['TEST_APP_PASSWORD']}@{host}:5432/{db_name}?sslmode=disable"
env['POOS_SESSION_REDIS_URL'] = f"redis://:{values['TEST_REDIS_PASSWORD']}@{address('redis')}:6379/0"
env.update(DB_HOST=host, DB_PORT='5432', DB_NAME=db_name, DB_USER='pos_migrator', DB_PASSWORD=values['TEST_MIGRATOR_PASSWORD'], DB_SSLMODE='disable')

def run(command, filename):
    with (EVIDENCE / filename).open('w') as output:
        process = subprocess.Popen(command, cwd=ROOT / 'backend', env=env, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, errors='replace')
        for line in process.stdout:
            for value in values.values():
                line = line.replace(value, '[redacted]')
            output.write(line)
            output.flush()
        code = process.wait()
    print(('PASS' if code == 0 else 'FAIL') + ' ' + filename, flush=True)
    return code

if __name__ == '__main__':
    if '--targeted' in sys.argv:
        raise SystemExit(run(['go', 'test', '-tags', 'integration', '-count=1', '-timeout', '10m', '-run', 'TestOrderStock', './internal/repository'], 'remediation-stock-tests.log'))
    codes = []
    for command, filename in [
        (['go', 'mod', 'verify'], 'remediation-go-mod-verify.log'),
        (['go', 'run', './cmd/migrate'], 'remediation-migrations.log'),
        (['go', 'vet', './...'], 'remediation-go-vet.log'),
        (['go', 'test', '-tags', 'integration', '-count=1', '-timeout', '10m', '-json', './...'], 'remediation-go-tests.log'),
    ]:
        code = run(command, filename)
        codes.append(code)
        if code and 'migrations' in filename:
            break
    raise SystemExit(any(codes))
