#!/usr/bin/env python3
"""Publish already smoke-tested CI images without exposing registry credentials."""
import json
import gzip
import hashlib
import shutil
import urllib.request
import os
from pathlib import Path
import re
import subprocess
import tarfile

sha = os.environ['GITHUB_SHA']
if os.environ.get('GITHUB_ACTIONS') != 'true' or not re.fullmatch('[0-9a-f]{40}', sha):
    raise SystemExit('Trusted immutable CI context required')
actor = os.environ['GITHUB_ACTOR']
subprocess.run(['docker', 'login', 'ghcr.io', '-u', actor, '--password-stdin'], input=os.environ['GHCR_TOKEN'], text=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=True)
bundle=Path('release-bundle');bundle.mkdir()
manifest = {'schema_version':1,'environment':'staging','sha':sha,'ci_run_id':int(os.environ['GITHUB_RUN_ID']),'images':{},'image_ids':{},'archives':{}}
try:
    for name in ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']:
        target = f'ghcr.io/yousef-abdallah-1711/pos/{name}:sha-{sha}'
        subprocess.run(['docker', 'tag', f'servio-ci-{name}:{sha}', target], check=True)
        subprocess.run(['docker', 'push', target], check=True)
        metadata = json.loads(subprocess.check_output(['docker', 'image', 'inspect', target]))[0]
        archive=bundle/(name+'.docker.tar.gz')
        with gzip.open(archive,'wb',compresslevel=3) as output:
            process=subprocess.Popen(['docker','save',target],stdout=subprocess.PIPE)
            shutil.copyfileobj(process.stdout,output)
            if process.wait()!=0:raise RuntimeError('Image archive failed')
        with tarfile.open(archive) as saved:
            entries=json.load(saved.extractfile('manifest.json'))
            if len(entries)!=1:raise RuntimeError('Expected one archived image')
            raw=saved.extractfile(entries[0]['Config']).read(2*1024**2+1)
            if len(raw)>2*1024**2:raise RuntimeError('Image config exceeds size bound')
            manifest['image_ids'][name]='sha256:'+hashlib.sha256(raw).hexdigest()
        manifest['archives'][name]={'file':archive.name,'sha256':hashlib.file_digest(archive.open('rb'),'sha256').hexdigest(),'size':archive.stat().st_size}
        manifest['images'][name] = next(d for d in metadata['RepoDigests'] if d.startswith(target.split(':')[0] + '@sha256:'))
    source=bundle/'source.tar.gz'
    urllib.request.urlretrieve('https://codeload.github.com/Yousef-Abdallah-1711/pos/tar.gz/'+sha,source)
    manifest['source_sha256']=hashlib.file_digest(source.open('rb'),'sha256').hexdigest()
    (bundle/'release-manifest.json').write_text(json.dumps(manifest,sort_keys=True,separators=(',',':'))+'\n')
finally:
    subprocess.run(['docker', 'logout', 'ghcr.io'], stdout=subprocess.DEVNULL, check=True)
