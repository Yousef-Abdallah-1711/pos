#!/usr/bin/env python3
"""Outbound-only staging release verifier. Uses anonymous, read-only public artifacts."""
import fcntl
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import re
import signal
import shutil
import subprocess
import sys
import time
import urllib.request

BASE = Path('/var/lib/servio-deployer')
KEY = Path('/etc/servio/release-signing.pub.pem')
WORKFLOW_ID=Path('/etc/servio/approved-workflow-id')
REPO = 'Yousef-Abdallah-1711/pos'
SERVICES = ['backend','landing','dashboard','restaurant-website','super-admin']
STATE = Path('/opt/servio/ops/aws/evidence/deployed-release.json')
MAX_ATTEMPTS = 3

class PendingRelease(Exception):
    pass

def request(url, limit=1024**2):
    req=urllib.request.Request(url,headers={'User-Agent':'servio-outbound-staging-agent','Accept':'application/vnd.github+json'})
    with urllib.request.urlopen(req,timeout=30) as response:
        body=response.read(limit+1)
    if len(body)>limit: raise RuntimeError('Release metadata exceeds size bound')
    return body

def github(path):
    return json.loads(request('https://api.github.com/repos/'+REPO+'/'+path))

def atomic_json(path, value):
    path.parent.mkdir(parents=True,exist_ok=True,mode=0o700)
    temporary=path.with_suffix('.tmp')
    temporary.write_text(json.dumps(value,indent=2)+'\n')
    temporary.chmod(0o600)
    temporary.replace(path)

def digest(path):
    with path.open('rb') as stream: return hashlib.file_digest(stream,'sha256').hexdigest()

def asset_url(release, name):
    asset=next((a for a in release['assets'] if a['name']==name),None)
    expected='https://github.com/'+REPO+'/releases/download/'+release['tag_name']+'/'
    if not asset or asset['browser_download_url']!=expected+name:
        raise RuntimeError('Missing or untrusted release asset: '+name)
    return asset['browser_download_url']

def download(release, name, target, expected):
    expected_size,expected_digest=expected
    if target.exists() and target.stat().st_size==expected_size and digest(target)==expected_digest: return
    temporary=target.with_suffix(target.suffix+'.part')
    count=0
    with urllib.request.urlopen(asset_url(release,name),timeout=60) as response,temporary.open('wb') as output:
        while chunk:=response.read(1024**2):
            count+=len(chunk)
            if count>expected_size: raise RuntimeError('Release asset exceeds declared size')
            output.write(chunk)
    if count!=expected_size or digest(temporary)!=expected_digest:
        raise RuntimeError('Release asset checksum mismatch: '+name)
    temporary.replace(target)

def verify_release(release):
    match=re.fullmatch(r'staging-sha-([0-9a-f]{40})',release['tag_name'])
    if not match or release['draft']: raise RuntimeError('Not a published immutable staging release')
    sha=match.group(1)
    if release.get('target_commitish')!=sha:raise RuntimeError('Release target SHA mismatch')
    if github('git/ref/heads/main')['object']['sha']!=sha: raise PendingRelease('Release is not current main')
    directory=BASE/'cache'/sha
    directory.mkdir(parents=True,exist_ok=True,mode=0o700)
    raw=request(asset_url(release,'release-manifest.json'),65536)
    signature=request(asset_url(release,'release-manifest.sig'),64)
    (directory/'release-manifest.json').write_bytes(raw)
    (directory/'release-manifest.sig').write_bytes(signature)
    result=subprocess.run(['openssl','pkeyutl','-verify','-rawin','-pubin','-inkey',str(KEY),'-in',str(directory/'release-manifest.json'),'-sigfile',str(directory/'release-manifest.sig')],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,timeout=10)
    if result.returncode: raise RuntimeError('Release signature does not match the pinned CI public key')
    manifest=json.loads(raw)
    if manifest.get('schema_version')!=1 or manifest.get('workflow_ref')!=REPO+'/.github/workflows/production-ci.yml@refs/heads/main' or manifest.get('workflow_sha')!=sha:raise RuntimeError('Signed workflow identity mismatch')
    if manifest.get('sha')!=sha or manifest.get('environment')!='staging': raise RuntimeError('Signed environment or SHA mismatch')
    run=github('actions/runs/'+str(manifest['ci_run_id']))
    if run.get('workflow_id')!=int(WORKFLOW_ID.read_text().strip()):raise RuntimeError('Unapproved workflow ID')
    if run['status']!='completed': raise PendingRelease('Release CI has not completed')
    expected=(sha,'main','push','success','.github/workflows/production-ci.yml')
    if (run['head_sha'],run['head_branch'],run['event'],run['conclusion'],run['path'])!=expected:
        raise RuntimeError('Release lacks successful exact-main CI provenance')
    return manifest,directory

def fetch_artifacts(release, manifest, directory):
    if set(manifest.get('images',{}))!=set(SERVICES) or set(manifest.get('archives',{}))!=set(SERVICES):
        raise RuntimeError('Incomplete release artifact set')
    if set(manifest.get('image_ids',{}))!=set(SERVICES):raise RuntimeError('Incomplete immutable image identities')
    for name in SERVICES:
        if not re.fullmatch('sha256:[0-9a-f]{64}',manifest['image_ids'][name]):raise RuntimeError('Invalid immutable image identity')
        if not re.fullmatch('ghcr.io/yousef-abdallah-1711/pos/'+re.escape(name)+'@sha256:[0-9a-f]{64}',manifest['images'][name]):raise RuntimeError('Unapproved immutable registry reference')
        asset=manifest['archives'][name]
        if asset['file']!=name+'.docker.tar.gz' or not isinstance(asset['size'],int) or not 0<asset['size']<=2*1024**3 or not re.fullmatch('[0-9a-f]{64}',asset['sha256']):
            raise RuntimeError('Invalid bounded archive metadata')
        download(release,asset['file'],directory/asset['file'],(asset['size'],asset['sha256']))
    source=next(a for a in release['assets'] if a['name']=='source.tar.gz')
    if not 0<source['size']<=150*1024**2 or not re.fullmatch('[0-9a-f]{64}',manifest['source_sha256']):
        raise RuntimeError('Invalid source artifact metadata')
    download(release,'source.tar.gz',directory/'source.tar.gz',(source['size'],manifest['source_sha256']))

def deployer():
    spec=importlib.util.spec_from_file_location('servio_installed_deployer','/usr/local/lib/servio/deploy_release.py')
    module=importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module

def poll(check_only=False):
    releases=github('releases?per_page=10')
    current=github('git/ref/heads/main')['object']['sha']
    release=next((r for r in releases if r['tag_name']=='staging-sha-'+current and not r['draft']),None)
    if not release: print('NOT READY: no published staging release'); return
    sha=release['tag_name'].removeprefix('staging-sha-')
    attempts_path=BASE/'attempts.json'
    attempts=json.loads(attempts_path.read_text()) if attempts_path.exists() else {}
    if attempts.get(sha,{}).get('count',0)>=MAX_ATTEMPTS and attempts.get(sha,{}).get('result')!='success': raise RuntimeError('Retry limit reached for '+sha+'; operator inspection required')
    try:
        manifest,directory=verify_release(release)
    except PendingRelease: raise
    except Exception:
        if not check_only:
            attempts[sha]={'count':attempts.get(sha,{}).get('count',0)+1,'result':'verification_failed'}
            atomic_json(attempts_path,attempts)
        raise
    if check_only: print('PASS signature and exact successful-main CI provenance '+sha); return
    if STATE.exists():
        deployed=json.loads(STATE.read_text())
        if deployed.get('sha')==sha and deployed.get('image_ids')==manifest.get('image_ids') and not deployed.get('candidate',False):
            for name in SERVICES:
                actual=subprocess.check_output(['docker','inspect','--format','{{.Image}}','servio-staging-'+name+'-1'],text=True).strip()
                if actual!=deployed.get('runtime_image_ids',{}).get(name):raise RuntimeError('Runtime image drift: '+name)
            print('PASS already deployed and immutable identities match '+sha); return
    attempts[sha]={'count':attempts.get(sha,{}).get('count',0)+1,'started_utc':time.strftime('%Y-%m-%dT%H:%M:%SZ',time.gmtime())}
    atomic_json(attempts_path,attempts)
    required=sum(asset['size'] for asset in manifest['archives'].values())+4*1024**3
    if shutil.disk_usage(BASE).free<required: raise RuntimeError('Insufficient disk headroom for verified download and deployment')
    fetch_artifacts(release,manifest,directory)
    module=deployer()
    module.deploy(manifest,sha,directory)
    # Only disposable download copies are removed; every loaded image, source,
    # volume, backup, private environment and signed provenance receipt remains.
    for name in SERVICES: (directory/manifest['archives'][name]['file']).unlink()
    attempts[sha]['result']='success'
    atomic_json(attempts_path,attempts)
    print('PASS outbound staging deployment '+sha)

def interrupted(signum, frame):
    raise RuntimeError('Deployment interrupted; stop promotion and restore application images')

if __name__=='__main__':
    os.umask(0o077)
    if os.geteuid()!=0: raise SystemExit('Root-managed staging agent only')
    signal.signal(signal.SIGTERM,interrupted)
    BASE.mkdir(parents=True,exist_ok=True,mode=0o700)
    lock=open('/run/lock/servio-staging-release.lock','w')
    try:
        fcntl.flock(lock,fcntl.LOCK_EX|fcntl.LOCK_NB)
        poll('--check-only' in sys.argv)
    except PendingRelease as error:
        print('NOT READY: '+str(error))
    except BlockingIOError:
        print('NOT READY: another staging deployment owns the lock')
    except Exception as error:
        print('BLOCKED: '+type(error).__name__+': '+str(error),file=sys.stderr)
        raise SystemExit(1)
