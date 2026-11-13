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
PUBLIC_STATUS=Path('/var/lib/servio-deployer/public/release.json')

def github(path):
    request = urllib.request.Request('https://api.github.com/repos/' + REPO + '/' + path, headers={'User-Agent': 'servio-staging-deployer', 'Accept': 'application/vnd.github+json'})
    with urllib.request.urlopen(request, timeout=20) as response:
        return json.load(response)

def checked(command, **kwargs):
    kwargs.setdefault("timeout",900)
    return subprocess.run(command, check=True, **kwargs)

def verify_manifest(manifest, requested, candidate=False):
    sha = manifest.get('sha', '')
    if not re.fullmatch('[0-9a-f]{40}', sha) or sha != requested:
        raise RuntimeError('Immutable release identity mismatch')
    main=github('git/ref/heads/main')['object']['sha']
    if not candidate and main != sha:
        raise RuntimeError('Release is not the current main commit')
    run_id = manifest.get('ci_run_id')
    if not isinstance(run_id, int) or run_id <= 0:
        raise RuntimeError('Missing CI run identity')
    run = github('actions/runs/' + str(run_id))
    if candidate:
        verify_candidate(run,sha,main)
    elif (run['head_sha'], run['head_branch'], run['event'], run['conclusion'], run['path']) != (sha, 'main', 'push', 'success', '.github/workflows/production-ci.yml'):
        raise RuntimeError('Exact main CI acceptance is missing')
    if set(manifest.get('images', {})) != set(SERVICES):
        raise RuntimeError('Incomplete application image set')
    for name, image in manifest['images'].items():
        if not re.fullmatch('ghcr.io/yousef-abdallah-1711/pos/' + re.escape(name) + '@sha256:[0-9a-f]{64}', image):
            raise RuntimeError('Unapproved registry or mutable image reference')
    return sha

def verify_candidate(run,sha,main):
    expected=(sha,'workflow_dispatch','success','.github/workflows/production-ci.yml','completed')
    if tuple(run.get(k) for k in ['head_sha','event','conclusion','path','status'])!=expected:
        raise RuntimeError('Candidate requires completed successful owner CI')
    if run.get('actor',{}).get('login')!='Yousef-Abdallah-1711' or run.get('workflow_id')!=415727583:
        raise RuntimeError('Unapproved candidate workflow identity')
    branch=run['head_branch']
    if github('git/ref/heads/'+branch)['object']['sha']!=sha:
        raise RuntimeError('Candidate branch has advanced')
    comparison=github('compare/'+main+'...'+sha)
    if comparison.get('status') not in ['ahead','identical']:
        raise RuntimeError('Candidate does not include current main')

def extract_source(sha, destination, artifacts):
    archive = artifacts / 'source.tar.gz'
    with tarfile.open(archive) as source:
        for member in source.getmembers():
            relative = Path(*Path(member.name).parts[1:])
            if not relative.parts:
                continue
            member.name = str(relative)
            source.extract(member, destination, filter='data')
    # Python's data filter intentionally ignores directory modes. Restore public
    # source traversal under the root agent's restrictive umask, before secrets exist.
    for file in destination.rglob('*'):
        if file.is_symlink():continue
        file.chmod(0o755 if file.is_dir() or file.stat().st_mode&0o111 else 0o644)

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

def verify_migrations(destination, previous_path):
    policy=Path('/etc/servio/approved-migrations.json')
    approved=json.loads(policy.read_text()) if policy.exists() else {}
    result=compose(previous_path,'exec','-T','postgres','psql','-U','postgres','-d','servio_staging','-At','-v','ON_ERROR_STOP=1','-c','SELECT version FROM schema_migrations ORDER BY version',stdout=subprocess.PIPE,text=True)
    applied=set(result.stdout.splitlines())
    for file in sorted((destination/'backend/migrations').glob('*.sql')):
        if file.name in applied:continue
        observed=hashlib.sha256(file.read_bytes()).hexdigest()
        if approved.get(file.name)!=observed:
            raise RuntimeError('Migration requires operator compatibility review: '+file.name)


def archive_member(source,name,limit=2*1024**2):
    stream=source.extractfile(name)
    if stream is None:raise RuntimeError('Missing image metadata')
    raw=stream.read(limit+1)
    if len(raw)>limit:raise RuntimeError('Image metadata exceeds size bound')
    return raw

def archive_identities(path,config_id):
    # Classic Docker identifies configs; the containerd store identifies manifests.
    # Both identities must be derived from the already checksum-verified archive.
    with tarfile.open(path) as source:
        entries=json.loads(archive_member(source,'manifest.json'))
        if len(entries)!=1:raise RuntimeError('Expected one immutable archived image')
        config=archive_member(source,entries[0]['Config'])
        if 'sha256:'+hashlib.sha256(config).hexdigest()!=config_id:
            raise RuntimeError('Archived config digest mismatch')
        identities={config_id}
        if 'index.json' in source.getnames():
            index=json.loads(archive_member(source,'index.json'))
            if len(index['manifests'])!=1:raise RuntimeError('Expected one OCI image manifest')
            descriptor=index['manifests'][0];digest=descriptor['digest']
            if not re.fullmatch('sha256:[0-9a-f]{64}',digest):raise RuntimeError('Invalid OCI identity')
            raw=archive_member(source,'blobs/sha256/'+digest.removeprefix('sha256:'))
            if 'sha256:'+hashlib.sha256(raw).hexdigest()!=digest or json.loads(raw)['config']['digest']!=config_id:
                raise RuntimeError('OCI manifest/config integrity mismatch')
            identities.add(digest)
        return identities

def load_verified_images(manifest,artifacts):
    runtime={}
    for name in SERVICES:
        archive=artifacts/manifest['archives'][name]['file']
        identities=archive_identities(archive,manifest['image_ids'][name])
        checked(['docker','load','-i',str(archive)],stdout=subprocess.DEVNULL)
        tag=manifest['images'][name].split('@')[0]+':sha-'+manifest['sha']
        metadata=json.loads(subprocess.check_output(['docker','image','inspect',tag]))[0]
        if metadata.get('Os')!='linux' or metadata.get('Architecture')!='amd64' or metadata['Id'] not in identities or metadata['Config'].get('Labels',{}).get('org.opencontainers.image.revision')!=manifest['sha']:
            raise RuntimeError('Loaded image identity mismatch: '+name)
        runtime[name]=metadata['Id']
    return runtime

def configure_runtime_images(path,runtime):
    config=json.loads(path.read_text())
    for name in SERVICES+['migrate']:
        config['services'][name]['image']=runtime['backend' if name=='migrate' else name]
    path.write_text(json.dumps(config,indent=2)+'\n')


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
    admin=json.loads((OPS/'evidence/browser-admin.json').read_text())
    if not admin.get('email','').endswith('@example.invalid'):raise RuntimeError('Synthetic acceptance admin required')
    storefront=json.loads((OPS/'evidence/storefront-qa-fixture.json').read_text())
    if not storefront.get('synthetic') or storefront.get('tenant_id')!=1 or storefront.get('restaurant_id')!=1:raise RuntimeError('Guarded storefront fixture required')
    shutil.copyfile(OPS / 'evidence/browser-admin.json', evidence / 'browser-admin.json')
    shutil.copyfile(OPS/'evidence/storefront-qa-fixture.json',evidence/'storefront-qa-fixture.json')
    # Synthetic acceptance admin only; never use the owner's interactive account.
    for file in evidence.iterdir(): file.chmod(0o600)
    (evidence/'.docker').mkdir(mode=0o700)
    checked(['chown', '-R', 'ubuntu:ubuntu', str(evidence)])
    commands = [['python3', str(destination / 'ops/aws' / name)] for name in ['smoke-api.py', 'smoke-refresh.py', 'smoke-admin.py', 'smoke-pos-kds.py', 'check-completed-sale.py']]
    commands.append([str(destination / 'ops/aws/run-browser-checks')])
    commands.append(['/usr/local/lib/servio/browser/node',str(destination/'ops/aws/browser-storefront-flow.cjs')])
    commands.append(['/usr/local/lib/servio/browser/node',str(destination/'ops/aws/browser-upload.cjs')])
    for command in commands:
        checked(['runuser','-u','ubuntu','--','env','DOCKER_CONFIG='+str(evidence/'.docker'),'NODE_PATH=/usr/local/lib/servio/browser/node_modules','PLAYWRIGHT_BROWSERS_PATH=/usr/local/lib/servio/browser/browsers',*command])

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

def record_release(manifest,path,previous,candidate=False):
    sha=manifest['sha']
    state={'candidate':candidate,'sha':sha,'compose':str(path),'images':manifest['images'],'image_ids':manifest['image_ids'],'runtime_image_ids':manifest['runtime_image_ids'],'ci_run_id':manifest['ci_run_id'],'previous':{key:previous[key] for key in ['sha','compose','images','image_ids','runtime_image_ids'] if key in previous},'deployed_utc':time.strftime('%Y-%m-%dT%H:%M:%SZ',time.gmtime())}
    status=PUBLIC_STATUS
    receipt={'sha':sha,'status':'candidate' if candidate else 'verified','environment':'staging','ci_run_id':manifest['ci_run_id']}
    snapshots={file:file.read_bytes() if file.exists() else None for file in [STATE,status]}
    try:
        replace_metadata(STATE,(json.dumps(state,indent=2)+'\n').encode(),0o600,True)
        replace_metadata(status,(json.dumps(receipt)+'\n').encode(),0o644)
    except Exception:
        for file,raw in snapshots.items():
            if raw is None:file.unlink(missing_ok=True)
            else:replace_metadata(file,raw,0o600 if file==STATE else 0o644,file==STATE)
        raise


def prepare_release(manifest,sha,artifacts,previous):
    destination=ROOT/'releases'/sha
    if destination.exists():
        if Path(previous['compose']).is_relative_to(destination):
            raise RuntimeError('Refusing to overwrite an active release')
        # Failed attempt source/evidence is retained, never deleted or overwritten.
        destination.rename(destination.with_name(sha+'-failed-'+str(time.time_ns())))
    destination.mkdir(parents=True,mode=0o755)
    destination.parent.chmod(0o755);destination.chmod(0o755)
    extract_source(sha,destination,artifacts)
    path=destination/'ops/aws/compose.json'
    config=json.loads(path.read_text())
    previous_config=json.loads(Path(previous['compose']).read_text())
    for name in ['postgres','redis','proxy']:
        if config['services'][name]['image']!=previous_config['services'][name]['image']:
            raise RuntimeError('Infrastructure image changes require a separate reviewed operation')
    for name in SERVICES+['migrate']:
        service=config['services'][name]
        service['image']=manifest['image_ids']['backend' if name=='migrate' else name]
        service['pull_policy']='never';service.pop('build',None)
    proxy=config['services']['proxy'];proxy['environment']={}
    proxy['volumes'].append('/var/lib/servio-deployer/public:/srv/servio-release:ro')
    proxy['volumes']=[v.replace('./Caddyfile:','./Caddyfile.public:') for v in proxy['volumes']]
    path.write_text(json.dumps(config,indent=2)+'\n')
    compose(path,'config','-q')
    return destination,path

def deploy(manifest,requested,artifacts,candidate=False):
    sha=verify_manifest(manifest,requested,candidate=candidate)
    verify_artifacts(manifest,artifacts)
    if shutil.disk_usage(ROOT).free<4*1024**3:
        raise RuntimeError('Less than 4 GiB disk headroom; refusing image load')
    previous=json.loads(STATE.read_text()) if STATE.exists() else {'compose':str(OPS/'compose.json')}
    previous_path=Path(previous['compose'])
    destination,path=prepare_release(manifest,sha,artifacts,previous)
    verify_migrations(destination,previous_path)
    runtime=load_verified_images(manifest,artifacts)
    configure_runtime_images(path,runtime)
    backup(previous_path)
    # Never reverse SQL automatically. Recheck provenance before promotion.
    compose(path,'run','--rm','--no-deps','migrate')
    verify_manifest(manifest,requested,candidate=candidate)
    try:
        for name in SERVICES+['proxy']:
            compose(path,'up','-d','--no-deps','--wait','--wait-timeout','150',name)
        public_checks()
        acceptance(destination)
        verify_manifest(manifest,requested,candidate=candidate)
        record_release(dict(manifest,runtime_image_ids=runtime),path,previous,candidate=candidate)
    except Exception:
        print('Deployment failed; restoring prior application configuration. SQL/volumes are preserved.',file=sys.stderr)
        rollback_apps(previous_path)
        raise
    print('PASS deployed immutable '+('staging candidate ' if candidate else 'main ')+sha)

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
