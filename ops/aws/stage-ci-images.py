#!/usr/bin/env python3
"""Archive browser-tested main images without registry or signing credentials."""

import gzip
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
import tarfile
import urllib.request
from release_image_archive import uncompressed_layer_bytes


ROOT = Path(__file__).resolve().parents[2]
IMAGE_NAMES = ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']
MAX_ARCHIVE_SIZE = 2 * 1024**3
MAX_SOURCE_SIZE = 150 * 1024**2
RELEASE_API_HOST = b'api.pos.rakmyat.online'


def verify_frontend_static(name: str, static_directory: Path) -> None:
    scripts = list(static_directory.rglob('*.js'))
    if not scripts:
        raise RuntimeError('Release frontend has no compiled JavaScript: ' + name)
    contains_release_api = False
    for script in scripts:
        content = script.read_bytes()
        if b'servio.test' in content:
            raise RuntimeError('Disposable CI host leaked into release frontend: ' + name)
        contains_release_api = contains_release_api or RELEASE_API_HOST in content
    if not contains_release_api:
        raise RuntimeError('Release frontend bundle lacks its staging API host: ' + name)


def verify_frontend_image(name: str, image_tag: str) -> None:
    container = subprocess.check_output(['docker', 'create', image_tag], text=True).strip()
    try:
        with tempfile.TemporaryDirectory(prefix='servio-release-frontend-') as directory:
            subprocess.run(
                ['docker', 'cp', f'{container}:/app/apps/{name}/.next/static', directory],
                check=True,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
            verify_frontend_static(name, Path(directory) / 'static')
    finally:
        subprocess.run(['docker', 'rm', container], check=True, stdout=subprocess.DEVNULL)


def digest(path: Path) -> str:
    with path.open('rb') as stream:
        return hashlib.file_digest(stream, 'sha256').hexdigest()


def archive_image(bundle: Path, name: str, sha: str) -> tuple[Path, str, int]:
    image_tag = f'servio-release-{name}:{sha}'
    config = json.loads(subprocess.check_output(['docker', 'image', 'inspect', image_tag], text=True))[0]
    labels = config.get('Config', {}).get('Labels') or {}
    if labels.get('org.opencontainers.image.revision') != sha:
        raise RuntimeError('Image revision label mismatch: ' + name)
    if name != 'backend':
        verify_frontend_image(name, image_tag)
    archive = bundle / (name + '.docker.tar.gz')
    with gzip.open(archive, 'wb', compresslevel=3) as output:
        process = subprocess.Popen(['docker', 'save', image_tag], stdout=subprocess.PIPE)
        if process.stdout is None:
            raise RuntimeError('Docker image stream unavailable')
        with process.stdout:
            shutil.copyfileobj(process.stdout, output)
        if process.wait() != 0:
            raise RuntimeError('Image archive failed: ' + name)
    if not 0 < archive.stat().st_size <= MAX_ARCHIVE_SIZE:
        raise RuntimeError('Image archive exceeds the deployment bound: ' + name)
    with tarfile.open(archive, 'r:gz') as saved:
        entries = json.load(saved.extractfile('manifest.json'))
        if len(entries) != 1 or entries[0].get('RepoTags') != [image_tag]:
            raise RuntimeError('Image archive tag mismatch: ' + name)
        raw = saved.extractfile(entries[0]['Config']).read(2 * 1024**2 + 1)
        if len(raw) > 2 * 1024**2:
            raise RuntimeError('Image config exceeds the size bound: ' + name)
    image_id = 'sha256:' + hashlib.sha256(raw).hexdigest()
    if config.get('Id') != image_id:
        raise RuntimeError('Archived image config does not match the tested image: ' + name)
    image_size = uncompressed_layer_bytes(archive)
    if not isinstance(image_size, int) or not 0 < image_size <= 5 * 1024**3:
        raise RuntimeError('Image size is missing or exceeds the deployment bound: ' + name)
    return archive, image_id, image_size


def stage_bundle(bundle: Path, sha: str, run_id: int) -> dict:
    bundle.mkdir()
    manifest = {
        'schema_version': 3,
        'environment': 'staging',
        'sha': sha,
        'ci_run_id': run_id,
        'images': {},
        'image_ids': {},
        'image_sizes': {},
        'archives': {},
    }
    for name in IMAGE_NAMES:
        archive, image_id, image_size = archive_image(bundle, name, sha)
        manifest['image_ids'][name] = image_id
        manifest['image_sizes'][name] = image_size
        manifest['archives'][name] = {
            'file': archive.name,
            'sha256': digest(archive),
            'size': archive.stat().st_size,
        }
    source = bundle / 'source.tar.gz'
    urllib.request.urlretrieve(
        'https://codeload.github.com/Yousef-Abdallah-1711/pos/tar.gz/' + sha,
        source,
    )
    if not 0 < source.stat().st_size <= MAX_SOURCE_SIZE:
        raise RuntimeError('Source archive exceeds the deployment bound')
    manifest['source_sha256'] = digest(source)
    manifest['source_size'] = source.stat().st_size
    (bundle / 'release-manifest.json').write_text(
        json.dumps(manifest, sort_keys=True, separators=(',', ':')) + '\n',
        encoding='utf-8',
    )
    return manifest


def main() -> int:
    sha = os.environ.get('GITHUB_SHA', '')
    run_id = os.environ.get('GITHUB_RUN_ID', '')
    if (
        os.environ.get('GITHUB_ACTIONS') != 'true'
        or os.environ.get('GITHUB_REPOSITORY') != 'Yousef-Abdallah-1711/pos'
        or os.environ.get('GITHUB_EVENT_NAME') != 'push'
        or os.environ.get('GITHUB_REF') != 'refs/heads/main'
        or not re.fullmatch('[0-9a-f]{40}', sha)
        or not re.fullmatch('[1-9][0-9]*', run_id)
    ):
        raise SystemExit('Refusing to stage release inputs outside an exact main push')
    manifest = stage_bundle(Path('release-bundle'), sha, int(run_id))
    print('PASS unprivileged release inputs staged for ' + sha + ' (CI run ' + str(manifest['ci_run_id']) + ')')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
