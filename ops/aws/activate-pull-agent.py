#!/usr/bin/env python3
"""Acceptance-gated activation; installation alone never enables automatic deployment."""
import fcntl
import json
import os
from pathlib import Path
import subprocess

ROOT=Path('/opt/servio')
LOCK=Path('/run/lock/servio-staging-release.lock')
MARKER=Path('/etc/servio/staging-agent.accepted')

def checked(command):subprocess.run(command,check=True,timeout=120)

def activate():
    state=json.loads((ROOT/'ops/aws/evidence/deployed-release.json').read_text())
    sha=state['sha']
    if state.get('recovery_required'):raise RuntimeError('Incomplete recovery cannot activate automation')
    if state.get('candidate',True):raise RuntimeError('Only an accepted main release can activate automation')
    release=ROOT/'releases'/sha
    if Path(state['compose'])!=release/'ops/aws/compose.json':raise RuntimeError('Unexpected deployment directory')
    evidence=release/'ops/aws/evidence'
    for name in ['database-role-smoke.json','api-smoke.json','refresh-smoke.json','admin-smoke.json','pos-kds-smoke.json','completed-sale-inventory.json','browser/results.json','storefront-browser-flow.json','upload-browser.json','journal-api-acceptance.json','auth-form-validation.json']:
        data=json.loads((evidence/name).read_text());rows=data['checks'] if isinstance(data,dict) else data
        if not rows or any(row.get('status')!='PASS' for row in rows):raise RuntimeError('Acceptance failed: '+name)
    full=json.loads((ROOT/'ops/aws/evidence/full-browser/run.json').read_text())
    if full.get('scope')!='full' or full.get('expected_checks')!=full.get('passed',0)+full.get('failed',0) or full.get('sha')!=sha or full.get('observed_sha')!=sha or not full.get('completed_utc') or full.get('failed')!=0 or full.get('passed',0)<1 or full.get('not_tested')!=0:
        raise RuntimeError('Complete successful browser traversal of this main SHA is required')
    # Normal polling also rechecks pinned signature, current successful main CI and
    # every runtime image identity. It cannot accept a feature-branch candidate.
    checked(['/usr/local/sbin/servio-pull-agent'])
    with LOCK.open('w') as lock:
        fcntl.flock(lock,fcntl.LOCK_EX|fcntl.LOCK_NB)
        current=json.loads((ROOT/'ops/aws/evidence/deployed-release.json').read_text())
        if current!=state:raise RuntimeError('Release changed during activation; verify the new SHA before retrying')
        checked([str(ROOT/'ops/aws/health')])
        MARKER.write_text(sha+'\n');MARKER.chmod(0o600)
        checked(['systemctl','enable','--now','servio-deploy.timer'])
    print('PASS staging-only recurring deployment enabled after acceptance: '+sha)

if __name__=='__main__':
    if os.geteuid()!=0:raise SystemExit('Root operator command only')
    activate()
