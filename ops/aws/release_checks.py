"""Fail-closed main publication gate for every protected-branch check identity."""
import json
import re
import subprocess
import time

REPO='Yousef-Abdallah-1711/pos'
EXPECTED={
    'required':(415727583,'.github/workflows/production-ci.yml'),
    'runtime-localhost-guardrail':(404030479,'.github/workflows/runtime-localhost-guardrail.yml'),
    'url-domain-config-smoke':(404030480,'.github/workflows/url-domain-config-smoke.yml'),
}
POLLS=60
DELAY=10

def api(path):
    result=subprocess.run(['gh','api','repos/'+REPO+'/'+path],capture_output=True,text=True,timeout=30)
    if result.returncode or len(result.stdout)>2*1024**2:
        raise RuntimeError('Required-check metadata unavailable or exceeds bounds')
    return json.loads(result.stdout)

def run_ready(check,sha,main_run):
    pattern=r'https://github.com/'+re.escape(REPO)+r'/actions/runs/(\d+)/job/\d+'
    match=re.fullmatch(pattern,check.get('details_url',''))
    if not match:raise RuntimeError('Untrusted required-check run URL')
    run_id=int(match.group(1));run=api('actions/runs/'+str(run_id))
    if (run.get('workflow_id'),run.get('path'))!=EXPECTED[check['name']]:
        raise RuntimeError('Required-check workflow identity mismatch')
    if (run.get('head_sha'),run.get('head_branch'),run.get('event'))!=(sha,'main','push'):
        raise RuntimeError('Required check is not an exact main push')
    if check['name']=='required':
        if run_id!=main_run:raise RuntimeError('Required aggregate belongs to a different CI run')
        if run.get('status') not in ['in_progress','completed'] or run.get('conclusion') not in [None,'success']:
            raise RuntimeError('Publication CI is cancelled or unsuccessful')
        return True  # This workflow is still running its publication job.
    if run.get('status')!='completed':return False
    if run.get('conclusion')!='success':raise RuntimeError('Required guardrail workflow failed')
    return True

def checks_ready(checks,sha,main_run):
    ready=True
    for name in EXPECTED:
        matches=[c for c in checks if c.get('name')==name and c.get('app',{}).get('slug')=='github-actions']
        if not matches:ready=False;continue
        check=max(matches,key=lambda c:c.get('id',0))
        if check.get('head_sha')!=sha:raise RuntimeError('Required-check commit mismatch')
        if check.get('status')!='completed':ready=False;continue
        if check.get('conclusion')!='success':raise RuntimeError('A protected-branch required check failed')
        if not run_ready(check,sha,main_run):ready=False
    return ready

def wait_checks(sha,main_run):
    if not re.fullmatch('[0-9a-f]{40}',sha) or not isinstance(main_run,int) or main_run<=0:
        raise RuntimeError('Invalid main publication identity')
    for attempt in range(POLLS):
        if api('git/ref/heads/main')['object']['sha']!=sha:
            raise RuntimeError('Main advanced before release publication')
        checks=api('commits/'+sha+'/check-runs?per_page=100&filter=latest')['check_runs']
        if checks_ready(checks,sha,main_run):
            print('PASS all protected-branch checks for exact successful main workflows');return
        if attempt+1<POLLS:time.sleep(DELAY)
    raise RuntimeError('Required main checks did not complete within the publication deadline')
