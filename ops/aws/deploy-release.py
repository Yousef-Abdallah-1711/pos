#!/usr/bin/env python3
"""Root-installed staging deployer for integrity-verified, immutable image archives."""
import fcntl
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tarfile
import tempfile
import time
import urllib.request

ROOT = Path('/opt/servio')
OPS = ROOT / 'ops/aws'
REPO = 'Yousef-Abdallah-1711/pos'
SERVICES = ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']
STATE = OPS / 'evidence/deployed-release.json'

def github(path):
    request = urllib.request.Request('https://api.github.com/repos/' + REPO + '/' + path, headers={'User-Agent': 'servio-staging-deployer', 'Accept': 'application/vnd.github+json'})
    with urllib.request.urlopen(request, timeout=20) as response:
        return json.load(response)

def checked(command, **kwargs):
    kwargs.setdefault("timeout",900)
    return subprocess.run(command, check=True, **kwargs)

def verify_manifest(manifest, requested):
    sha = manifest.get('sha', '')
    if not re.fullmatch('[0-9a-f]{40}', sha) or sha != requested:
        raise RuntimeError('Immutable release identity mismatch')
    if github('git/ref/heads/main')['object']['sha'] != sha:
        raise RuntimeError('Release is not the current main commit')
    run_id = manifest.get('ci_run_id')
    if not isinstance(run_id, int) or run_id <= 0:
        raise RuntimeError('Missing CI run identity')
    run = github('actions/runs/' + str(run_id))
    if (run['head_sha'], run['head_branch'], run['event'], run['conclusion'], run['path']) != (sha, 'main', 'push', 'success', '.github/workflows/production-ci.yml'):
        raise RuntimeError('Exact main CI acceptance is missing')
    if set(manifest.get('images', {})) != set(SERVICES):
        raise RuntimeError('Incomplete application image set')
    for name, image in manifest['images'].items():
        if not re.fullmatch('ghcr.io/yousef-abdallah-1711/pos/' + re.escape(name) + '@sha256:[0-9a-f]{64}', image):
            raise RuntimeError('Unapproved registry or mutable image reference')
    return sha

def extract_source(sha, destination, artifacts):
    archive = artifacts / 'source.tar.gz'
    with tarfile.open(archive) as source:
        for member in source.getmembers():
            relative = Path(*Path(member.name).parts[1:])
            if not relative.parts:
                continue
            member.name = str(relative)
            source.extract(member, destination, filter='data')

def verify_artifacts(manifest, artifacts):
    source=artifacts/'source.tar.gz'
    with source.open('rb') as stream:
        if hashlib.file_digest(stream,'sha256').hexdigest()!=manifest['source_sha256']:
            raise RuntimeError('Source artifact integrity failure')
    for name in SERVICES:
        asset=manifest['archives'][name]
        if asset['file']!=name+'.docker.tar.gz' or not isinstance(asset['size'],int) or not 0<asset['size']<=2*1024**3:
            raise RuntimeError('Unapproved image artifact')
        path=artifacts/asset['file']
        with path.open('rb') as stream:
            if path.stat().st_size!=asset['size'] or hashlib.file_digest(stream,'sha256').hexdigest()!=asset['sha256']:
                raise RuntimeError('Image archive integrity failure: '+name)
        if not re.fullmatch('sha256:[0-9a-f]{64}',manifest['image_ids'][name]):
            raise RuntimeError('Missing immutable image content identity')

def load_verified_images(manifest, artifacts):
    for name in SERVICES:
        checked(['docker','load','-i',str(artifacts/manifest['archives'][name]['file'])],stdout=subprocess.DEVNULL)
        metadata=json.loads(subprocess.check_output(['docker','image','inspect',manifest['image_ids'][name]]))[0]
        if metadata['Id']!=manifest['image_ids'][name] or metadata['Config'].get('Labels',{}).get('org.opencontainers.image.revision')!=manifest['sha']:
            raise RuntimeError('Loaded image identity mismatch: '+name)


def compose(path, *arguments, **kwargs):
    return checked(['docker', 'compose', '-p', 'servio-staging', '--env-file', str(OPS / '.env'), '-f', str(path), *arguments], **kwargs)

def backup(path):
    target = OPS / 'backups' / ('pre-release-' + time.strftime('%Y%m%dT%H%M%SZ', time.gmtime()) + '.dump')
    fd = os.open(target, os.O_CREAT | os.O_EXCL | os.O_WRONLY, 0o600)
    with os.fdopen(fd, 'wb') as output:
        compose(path, 'exec', '-T', 'postgres', 'pg_dump', '-U', 'postgres', '-d', 'servio_staging', '-Fc', stdout=output)
    with target.open('rb') as source:
        compose(path, 'exec', '-T', 'postgres', 'pg_restore', '--list', stdin=source, stdout=subprocess.DEVNULL)
    sidecar = target.with_suffix('.dump.sha256')
    with target.open('rb') as stream:
        sidecar.write_text(hashlib.file_digest(stream,'sha256').hexdigest() + '  ' + target.name + '\n')
    checked(['chown', 'ubuntu:ubuntu', str(target), str(sidecar)])

def public_checks():
    paths = {'pos.rakmyat.online': '/en', 'api.pos.rakmyat.online': '/api/v1/health', 'pos.app.rakmyat.online': '/en/dashboard', 'admin.pos.rakmyat.online': '/login', 'restaurant-1.pos.rakmyat.online': '/en'}
    for host, path in paths.items():
        with urllib.request.urlopen('https://' + host + path, timeout=30) as response:
            if response.status != 200:
                raise RuntimeError('Public route failed: ' + host)
            if host.startswith('api.') and not json.load(response).get('db_healthy'):
                raise RuntimeError('Database health failed')

def acceptance(destination):
    evidence = destination / 'ops/aws/evidence'
    evidence.mkdir(mode=0o700)
    shutil.copyfile(OPS / 'evidence/browser-admin.json', evidence / 'browser-admin.json')
    # Synthetic acceptance admin only; never use the owner's interactive account.
    for file in evidence.iterdir(): file.chmod(0o600)
    (evidence/'.docker').mkdir(mode=0o700)
    checked(['chown', '-R', 'ubuntu:ubuntu', str(evidence)])
    commands = [['python3', str(destination / 'ops/aws' / name)] for name in ['smoke-api.py', 'smoke-refresh.py', 'smoke-admin.py', 'smoke-pos-kds.py', 'check-completed-sale.py']]
    commands.append([str(destination / 'ops/aws/run-browser-checks')])
    for command in commands:
        checked(['runuser','-u','ubuntu','--','env','DOCKER_CONFIG='+str(evidence/'.docker'),*command])

def rollback_apps(previous_path):
    failures=[]
    for name in SERVICES+['proxy']:
        try:
            compose(previous_path,'up','-d','--no-deps','--wait','--wait-timeout','150',name)
        except Exception:
            failures.append(name)
    if failures:
        raise RuntimeError('Application rollback incomplete: '+', '.join(failures))

def replace_metadata(path, raw, mode, owner=False):
    path.parent.mkdir(parents=True,exist_ok=True)
    temporary=path.with_suffix('.tmp')
    temporary.write_bytes(raw)
    temporary.chmod(mode)
    if owner:checked(['chown','ubuntu:ubuntu',str(temporary)])
    temporary.replace(path)

def record_release(manifest,sha,path,previous):
    state={'sha':sha,'compose':str(path),'images':manifest['images'],'image_ids':manifest['image_ids'],'ci_run_id':manifest['ci_run_id'],'previous':{key:previous[key] for key in ['sha','compose','images','image_ids'] if key in previous},'deployed_utc':time.strftime('%Y-%m-%dT%H:%M:%SZ',time.gmtime())}
    status=Path('/var/lib/servio-deployer/public/release.json')
    receipt={'sha':sha,'status':'verified','environment':'staging','ci_run_id':manifest['ci_run_id']}
    snapshots={file:file.read_bytes() if file.exists() else None for file in [STATE,status]}
    try:
        replace_metadata(STATE,(json.dumps(state,indent=2)+'\n').encode(),0o600,True)
        replace_metadata(status,(json.dumps(receipt)+'\n').encode(),0o644)
    except Exception:
        for file,raw in snapshots.items():
            if raw is None:file.unlink(missing_ok=True)
            else:replace_metadata(file,raw,0o600 if file==STATE else 0o644,file==STATE)
        raise


def deploy(manifest, requested, artifacts):
    sha = verify_manifest(manifest, requested)
    verify_artifacts(manifest, artifacts)
    if shutil.disk_usage(ROOT).free < 6 * 1024**3:
        raise RuntimeError('Less than 6 GiB disk headroom; refusing image pull')
    destination = ROOT / 'releases' / sha
    if destination.exists():
        raise RuntimeError('Release directory already exists; inspect prior attempt before retry')
    destination.mkdir(parents=True, mode=0o755)
    (ROOT / 'releases').chmod(0o755)
    destination.chmod(0o755)
    extract_source(sha, destination, artifacts)
    path = destination / 'ops/aws/compose.json'
    config = json.loads(path.read_text())
    previous = json.loads(STATE.read_text()) if STATE.exists() else {'compose': str(OPS / 'compose.json')}
    previous_path = Path(previous['compose'])
    previous_config = json.loads(previous_path.read_text())
    for name in ['postgres', 'redis', 'proxy']:
        if config['services'][name]['image'] != previous_config['services'][name]['image']:
            raise RuntimeError('Infrastructure image changes require a separate reviewed operation')
    for name in SERVICES + ['migrate']:
        service = config['services'][name]
        service['image'] = manifest['image_ids']['backend' if name == 'migrate' else name]
        service['pull_policy']='never'
        service.pop('build', None)
    config['services']['proxy']['environment'] = {}
    config['services']['proxy']['volumes'].append('/var/lib/servio-deployer/public:/srv/servio-release:ro')
    config['services']['proxy']['volumes'] = [v.replace('./Caddyfile:', './Caddyfile.public:') for v in config['services']['proxy']['volumes']]
    path.write_text(json.dumps(config, indent=2) + '\n')
    compose(path, 'config', '-q')
    load_verified_images(manifest, artifacts)
    backup(previous_path)
    # Never reverse SQL automatically. Recheck main immediately before promotion.
    compose(path, 'run', '--rm', '--no-deps', 'migrate')
    verify_manifest(manifest, requested)
    try:
        for name in SERVICES + ['proxy']:
            compose(path, 'up', '-d', '--no-deps', '--wait', '--wait-timeout', '150', name)
        public_checks()
        acceptance(destination)
        record_release(manifest,sha,path,previous)
    except Exception:
        print('Deployment failed; restoring prior application configuration. SQL/volumes are preserved.', file=sys.stderr)
        rollback_apps(previous_path)
        raise
    print('PASS deployed immutable main ' + sha)

if __name__ == '__main__':
    os.umask(0o077)
    if os.geteuid() != 0 or len(sys.argv) != 3 or not re.fullmatch('[0-9a-f]{40}', sys.argv[1]):
        raise SystemExit('Use the installed restricted staging entrypoint with a full main SHA')
    lock = open('/run/lock/servio-staging-release.lock', 'w')
    fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
    raw = sys.stdin.read(65537)
    if len(raw) > 65536:
        raise SystemExit('Release manifest exceeds limit')
    try:
        deploy(json.loads(raw), sys.argv[1], Path(sys.argv[2]))
    except Exception as error:
        # Exception class only: never reflect credential-bearing input or subprocess arguments.
        message = str(error)
        try:
            token = json.loads(raw).get('registry_token', '')
            if token: message = message.replace(token, '[redacted]')
        except ValueError: pass
        print('BLOCKED staging release: ' + type(error).__name__ + ': ' + message, file=sys.stderr)
        raise SystemExit(1)
