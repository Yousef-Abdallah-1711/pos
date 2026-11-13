#!/usr/bin/env python3
"""Sign data-only CI artifacts and publish a verified staging release from protected main."""
import json
import hashlib
import os
from pathlib import Path
import subprocess
import tempfile

sha=os.environ['GITHUB_SHA']
if os.environ.get('GITHUB_ACTIONS')!='true' or os.environ.get('GITHUB_REF')!='refs/heads/main' or os.environ.get('GITHUB_EVENT_NAME')!='push':
    raise SystemExit('Protected main push context required')
root=Path('release-bundle')
manifest=json.loads((root/'release-manifest.json').read_text())
if manifest['sha']!=sha or manifest['ci_run_id']!=int(os.environ['GITHUB_RUN_ID']):
    raise SystemExit('Artifact provenance identity mismatch')
expected={'release-manifest.json','source.tar.gz',*[name+'.docker.tar.gz' for name in ['backend','landing','dashboard','restaurant-website','super-admin']]}
if {p.name for p in root.iterdir()}!=expected:
    raise SystemExit('Unexpected release artifact; refusing publication')
for name,asset in manifest['archives'].items():
    with (root/asset['file']).open('rb') as source:
        if hashlib.file_digest(source,'sha256').hexdigest()!=asset['sha256']:raise SystemExit('Image artifact integrity failure')
with (root/'source.tar.gz').open('rb') as source:
    if hashlib.file_digest(source,'sha256').hexdigest()!=manifest['source_sha256']:raise SystemExit('Source artifact integrity failure')
manifest['workflow_ref']=os.environ['GITHUB_WORKFLOW_REF']
manifest['workflow_sha']=os.environ['GITHUB_WORKFLOW_SHA']
if manifest['workflow_sha']!=sha:raise SystemExit('Workflow source SHA mismatch')
(root/'release-manifest.json').write_text(json.dumps(manifest,sort_keys=True,separators=(',',':'))+'\n')
with tempfile.TemporaryDirectory(prefix='servio-sign-') as private:
    key=Path(private)/'signing.pem'
    key.write_text(os.environ['SERVIO_RELEASE_SIGNING_KEY'])
    key.chmod(0o600)
    subprocess.run(['openssl','pkeyutl','-sign','-rawin','-inkey',str(key),'-in',str(root/'release-manifest.json'),'-out',str(root/'release-manifest.sig')],check=True)
tag='staging-sha-'+sha
# Draft until every immutable asset is uploaded. Never overwrite an existing release.
subprocess.run(['gh','release','create',tag,'--target',sha,'--title','Verified SERVIO staging '+sha[:12],'--notes','Staging only. Immutable images and source verified by SERVIO CI; automatic production promotion is disabled.','--draft',*[str(p) for p in sorted(root.iterdir())]],check=True)
subprocess.run(['gh','release','edit',tag,'--draft=false','--latest'],check=True)
print('PASS signed staging release '+sha)
