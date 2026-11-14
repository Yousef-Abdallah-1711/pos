#!/usr/bin/env python3
"""Run repository tests only against an isolated, positively identified test stack."""
import os,subprocess
from pathlib import Path
root=Path('/opt/servio');os.chdir(root)
secret_file=root/'ops/aws/evidence/acceptance.env'
secret_values=[line.split('=',1)[1] for line in secret_file.read_text().splitlines() if '=' in line]
def run_log(command,name):
    path=root/'ops/aws/evidence'/name
    with path.open('w') as output:
        process=subprocess.Popen(command,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True,errors='replace')
        for line in process.stdout:
            for value in secret_values:line=line.replace(value,'[redacted]')
            output.write(line);output.flush()
        return process.wait()

def verify():
    result=subprocess.run(['sudo','docker','inspect','--format','{{index .Config.Labels "com.servio.disposable"}}','servio-acceptance-postgres-1'],capture_output=True,text=True,check=True)
    if result.stdout.strip()!='true':raise RuntimeError('Disposable database label is missing; refusing tests')
    result=subprocess.run(['ops/aws/acceptance','exec','-T','postgres','printenv','POSTGRES_DB'],capture_output=True,text=True,check=True)
    if result.stdout.strip()!='servio_test':raise RuntimeError('Unexpected database; refusing tests')

try:
    code=run_log(['sudo','docker','buildx','build','--load','--target','builder','-f','ops/aws/Dockerfile.backend-check','-t','servio-go-acceptance:208bcae-aws1','backend'],'go-tools-build.log')
    if code:raise RuntimeError('Go tool image build failed')
    if run_log(['ops/aws/acceptance','up','-d','--wait','--wait-timeout','120','postgres','redis'],'go-acceptance-start.log'):raise RuntimeError('Isolated dependencies failed readiness')
    verify()
    if run_log(['ops/aws/acceptance','run','--rm','migrate'],'go-acceptance-migrations.log'):raise RuntimeError('Isolated test migrations failed')
    verify()
    vet=run_log(['ops/aws/acceptance','run','--rm','--no-deps','tools','vet','-p','1','./...'],'go-vet.log')
    tests=run_log(['ops/aws/acceptance','run','--rm','--no-deps','tools'],'go-test-integration.log')
    print(('PASS' if vet==0 else 'FAIL')+' go vet ./...')
    print(('PASS' if tests==0 else 'FAIL')+' go test -tags integration -p 1 ./...')
    Path('ops/aws/evidence/go-exit-codes.txt').write_text(f'vet={vet}\ntests={tests}\n')
    raise SystemExit(bool(vet or tests))
finally:
    # Preserve volumes and test evidence; release test-process RAM.
    run_log(['ops/aws/acceptance','stop'],'go-acceptance-stop.log')
