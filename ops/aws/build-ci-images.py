#!/usr/bin/env python3
"""Build the exact five application images on a trusted CI runner."""
import json
import os
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[2]
SHA = os.environ.get('GITHUB_SHA', '')
if os.environ.get('GITHUB_ACTIONS') != 'true' or not re.fullmatch('[0-9a-f]{40}', SHA):
    raise SystemExit('Refusing CI builds outside GitHub Actions or without an immutable SHA')
config = json.loads((ROOT / 'ops/aws/compose.json').read_text())
for name in ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']:
    build = config['services'][name]['build']
    context = (ROOT / 'ops/aws' / build['context']).resolve()
    dockerfile = (context / build['dockerfile']).resolve()
    command = ['docker', 'build', '--label', 'org.opencontainers.image.revision=' + SHA, '--label', 'org.opencontainers.image.source=https://github.com/Yousef-Abdallah-1711/pos', '-t', f'servio-ci-{name}:{SHA}', '-f', str(dockerfile)]
    for key, value in build.get('args', {}).items():
        command.extend(['--build-arg', key + '=' + value])
    subprocess.run(command + [str(context)], check=True)
