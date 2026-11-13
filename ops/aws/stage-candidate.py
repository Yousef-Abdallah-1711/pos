#!/usr/bin/env python3
"""Operator-only staging acceptance of an owner-dispatched CI artifact, never automation."""
import fcntl
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import re
import sys
import urllib.request
import zipfile

REPO='Yousef-Abdallah-1711/pos'

def github(path):
    request=urllib.request.Request('https://api.github.com/repos/'+REPO+'/'+path,headers={'User-Agent':'servio-staging-candidate'})
    with urllib.request.urlopen(request,timeout=30) as response:return json.load(response)

def extract_verified(archive,artifact,destination):
    expected=artifact.get('digest','')
    with archive.open('rb') as stream:observed='sha256:'+hashlib.file_digest(stream,'sha256').hexdigest()
    if not re.fullmatch('sha256:[0-9a-f]{64}',expected) or observed!=expected:
        raise RuntimeError('GitHub artifact ZIP integrity mismatch')
    names={'release-manifest.json','source.tar.gz',*[name+'.docker.tar.gz' for name in ['backend','landing','dashboard','restaurant-website','super-admin']]}
    with zipfile.ZipFile(archive) as source:
        entries=source.infolist()
        if {entry.filename for entry in entries}!=names or len(entries)!=len(names):
            raise RuntimeError('Unexpected candidate artifact members')
        if any(entry.file_size>2*1024**3 for entry in entries) or sum(entry.file_size for entry in entries)>3*1024**3:
            raise RuntimeError('Candidate artifact exceeds extraction bounds')
        destination.mkdir(mode=0o700,parents=True)
        source.extractall(destination)

def main(run_id,archive):
    artifact=next(a for a in github('actions/runs/'+str(run_id)+'/artifacts')['artifacts'] if a['name']=='release-bundle' and not a['expired'])
    spec=importlib.util.spec_from_file_location('installed_deployer','/usr/local/lib/servio/deploy_release.py')
    module=importlib.util.module_from_spec(spec);spec.loader.exec_module(module)
    run=github('actions/runs/'+str(run_id))
    sha=run['head_sha'];main_sha=github('git/ref/heads/main')['object']['sha']
    module.verify_candidate(run,sha,main_sha)
    destination=Path('/var/lib/servio-deployer/candidates')/sha
    extract_verified(archive,artifact,destination)
    manifest=json.loads((destination/'release-manifest.json').read_text())
    if manifest.get('ci_run_id')!=run_id:raise RuntimeError('Candidate artifact run mismatch')
    module.deploy(manifest,sha,destination,candidate=True)

if __name__=='__main__':
    if os.geteuid()!=0 or len(sys.argv)!=3 or not sys.argv[1].isdigit():
        raise SystemExit('Operator root command: stage-candidate.py CI_RUN_ID VERIFIED_ARTIFACT_ZIP')
    os.umask(0o077)
    lock=open('/run/lock/servio-staging-release.lock','w')
    fcntl.flock(lock,fcntl.LOCK_EX|fcntl.LOCK_NB)
    main(int(sys.argv[1]),Path(sys.argv[2]))
