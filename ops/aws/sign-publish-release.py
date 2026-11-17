#!/usr/bin/env python3
"""Sign data-only CI artifacts and publish a verified staging release from protected main."""
import json
import hashlib
import os
from pathlib import Path
import re
import subprocess
import tempfile

from release_checks import EXPECTED, api, wait_checks

sha=os.environ['GITHUB_SHA']
ci_run_id=int(os.environ['SERVIO_CI_RUN_ID'])
if os.environ.get('GITHUB_ACTIONS')!='true' or os.environ.get('GITHUB_REPOSITORY')!='Yousef-Abdallah-1711/pos' or os.environ.get('GITHUB_REF')!='refs/heads/main' or os.environ.get('GITHUB_EVENT_NAME')!='workflow_dispatch':
    raise SystemExit('Owner-dispatched protected main publication context required')
if os.environ.get('GITHUB_ACTOR')!=os.environ.get('GITHUB_REPOSITORY_OWNER') or os.environ.get('SERVIO_RELEASE_CONFIRMATION')!='PUBLISH-STAGING':
    raise SystemExit('Repository-owner staging publication confirmation required')
if not re.fullmatch('[0-9a-f]{40}',sha) or not os.environ.get('SERVIO_RELEASE_SIGNING_KEY'):
    raise SystemExit('Exact main SHA and protected signing key are required')
wait_checks(sha,ci_run_id)
ci_run=api('actions/runs/'+str(ci_run_id))
if (ci_run.get('workflow_id'),ci_run.get('path'),ci_run.get('head_sha'),ci_run.get('head_branch'),ci_run.get('event'),ci_run.get('status'),ci_run.get('conclusion'))!=(EXPECTED['required'][0],EXPECTED['required'][1],sha,'main','push','completed','success'):
    raise SystemExit('The exact-main CI run is not fully successful')
ci_jobs=api('actions/runs/'+str(ci_run_id)+'/jobs?per_page=100')['jobs']
publish_jobs=[job for job in ci_jobs if job.get('name')=='publish-release-images']
if len(publish_jobs)!=1 or publish_jobs[0].get('status')!='completed' or publish_jobs[0].get('conclusion')!='success':
    raise SystemExit('Exact-main immutable image publication did not complete successfully')
root=Path('release-bundle')
manifest=json.loads((root/'release-manifest.json').read_text())
if manifest.get('schema_version')!=3 or manifest.get('environment')!='staging' or manifest.get('sha')!=sha or manifest.get('ci_run_id')!=ci_run_id:
    raise SystemExit('Artifact provenance identity mismatch')
names=['backend','landing','dashboard','restaurant-website','super-admin']
expected={'release-manifest.json','source.tar.gz',*[name+'.docker.tar.gz' for name in names]}
if {p.name for p in root.iterdir()}!=expected:
    raise SystemExit('Unexpected release artifact; refusing publication')
if set(manifest.get('archives',{}))!=set(names) or set(manifest.get('images',{}))!=set(names) or set(manifest.get('image_ids',{}))!=set(names) or set(manifest.get('image_sizes',{}))!=set(names):
    raise SystemExit('Incomplete signed image identity set')
for name,asset in manifest['archives'].items():
    if asset.get('file')!=name+'.docker.tar.gz' or not isinstance(asset.get('size'),int) or not 0<asset['size']<=2*1024**3 or not re.fullmatch('[0-9a-f]{64}',asset.get('sha256','')):
        raise SystemExit('Invalid image artifact metadata')
    artifact=root/asset['file']
    if artifact.stat().st_size!=asset['size']:
        raise SystemExit('Image artifact size mismatch')
    with artifact.open('rb') as source:
        if hashlib.file_digest(source,'sha256').hexdigest()!=asset['sha256']:raise SystemExit('Image artifact integrity failure')
with (root/'source.tar.gz').open('rb') as source:
    source_size=os.fstat(source.fileno()).st_size
    if not 0<source_size<=150*1024**2 or manifest.get('source_size')!=source_size or hashlib.file_digest(source,'sha256').hexdigest()!=manifest['source_sha256']:raise SystemExit('Source artifact integrity failure')
for name in names:
    if not re.fullmatch('sha256:[0-9a-f]{64}',manifest['image_ids'][name]):
        raise SystemExit('Invalid immutable image identity')
    if not isinstance(manifest['image_sizes'][name],int) or not 0<manifest['image_sizes'][name]<=5*1024**3:
        raise SystemExit('Invalid immutable image size')
    if not re.fullmatch('ghcr.io/yousef-abdallah-1711/pos/'+re.escape(name)+r'@sha256:[0-9a-f]{64}',manifest['images'][name]):
        raise SystemExit('Invalid immutable registry image reference')
manifest['workflow_ref']='Yousef-Abdallah-1711/pos/.github/workflows/production-ci.yml@refs/heads/main'
manifest['workflow_sha']=sha
if manifest['workflow_sha']!=sha:raise SystemExit('Workflow source SHA mismatch')
(root/'release-manifest.json').write_text(json.dumps(manifest,sort_keys=True,separators=(',',':'))+'\n')
with tempfile.TemporaryDirectory(prefix='servio-sign-') as private:
    key=Path(private)/'signing.pem'
    key.write_text(os.environ['SERVIO_RELEASE_SIGNING_KEY'])
    key.chmod(0o600)
    subprocess.run(['openssl','pkeyutl','-sign','-rawin','-inkey',str(key),'-in',str(root/'release-manifest.json'),'-out',str(root/'release-manifest.sig')],check=True)
if api('git/ref/heads/main')['object']['sha']!=sha:
    raise SystemExit('Main advanced during signing; refusing to publish a stale release')
tag='staging-sha-'+sha
# Draft until every immutable asset is uploaded. Never overwrite an existing release.
subprocess.run(['gh','release','create',tag,'--target',sha,'--title','Verified SERVIO staging '+sha[:12],'--notes','Staging only. Immutable images and source verified by SERVIO CI; automatic production promotion is disabled.','--draft',*[str(p) for p in sorted(root.iterdir())]],check=True)
subprocess.run(['gh','release','edit',tag,'--draft=false','--latest'],check=True)
print('PASS signed staging release '+sha)
