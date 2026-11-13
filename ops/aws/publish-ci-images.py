#!/usr/bin/env python3
"""Publish already smoke-tested CI images without exposing registry credentials."""
import json
import os
from pathlib import Path
import re
import subprocess

sha = os.environ['GITHUB_SHA']
if os.environ.get('GITHUB_ACTIONS') != 'true' or not re.fullmatch('[0-9a-f]{40}', sha):
    raise SystemExit('Trusted immutable CI context required')
actor = os.environ['GITHUB_ACTOR']
subprocess.run(['docker', 'login', 'ghcr.io', '-u', actor, '--password-stdin'], input=os.environ['GHCR_TOKEN'], text=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=True)
manifest = {'sha': sha, 'ci_run_id': int(os.environ['GITHUB_RUN_ID']), 'images': {}}
try:
    for name in ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']:
        target = f'ghcr.io/yousef-abdallah-1711/pos/{name}:sha-{sha}'
        subprocess.run(['docker', 'tag', f'servio-ci-{name}:{sha}', target], check=True)
        subprocess.run(['docker', 'push', target], check=True)
        metadata = json.loads(subprocess.check_output(['docker', 'image', 'inspect', target]))[0]
        manifest['images'][name] = next(d for d in metadata['RepoDigests'] if d.startswith(target.split(':')[0] + '@sha256:'))
    Path('release-manifest.json').write_text(json.dumps(manifest, indent=2) + '\n')
finally:
    subprocess.run(['docker', 'logout', 'ghcr.io'], stdout=subprocess.DEVNULL, check=True)
