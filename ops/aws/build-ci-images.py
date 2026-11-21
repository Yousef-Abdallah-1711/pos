#!/usr/bin/env python3
"""Build isolated smoke images and exact release-configured main images."""
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
release_push = (
    os.environ.get('GITHUB_REPOSITORY') == 'Yousef-Abdallah-1711/pos'
    and os.environ.get('GITHUB_EVENT_NAME') == 'push'
    and os.environ.get('GITHUB_REF') == 'refs/heads/main'
)
test_hosts = {
    'api.pos.rakmyat.online': 'api.servio.test',
    'pos.app.rakmyat.online': 'pos.app.servio.test',
    'admin.pos.rakmyat.online': 'admin.pos.servio.test',
    'restaurant-1.pos.rakmyat.online': 'restaurant-1.pos.servio.test',
    'pos.rakmyat.online': 'pos.servio.test',
}
def build_image(name, profile, build, rewrite_hosts=False):
    context = (ROOT / 'ops/aws' / build['context']).resolve()
    dockerfile = (context / build['dockerfile']).resolve()
    command = [
        'docker', 'build',
        '--label', 'org.opencontainers.image.revision=' + SHA,
        '--label', 'org.opencontainers.image.source=https://github.com/Yousef-Abdallah-1711/pos',
        '-t', f'servio-{profile}-{name}:{SHA}', '-f', str(dockerfile),
    ]
    for key, original_value in build.get('args', {}).items():
        value = original_value
        if rewrite_hosts:
            for production_host, test_host in test_hosts.items():
                value = value.replace(production_host, test_host)
        command.extend(['--build-arg', key + '=' + value])
    subprocess.run(command + [str(context)], check=True)


for name in ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']:
    build = config['services'][name]['build']
    # These disposable images are the only images consumed by isolated browser smoke.
    build_image(name, 'ci', build, rewrite_hosts=True)
    if release_push:
        # Release artifacts must retain the exact origins from the signed staging compose.
        build_image(name, 'release', build)
