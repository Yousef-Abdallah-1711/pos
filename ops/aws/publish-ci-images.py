#!/usr/bin/env python3
"""Publish exact-main CI image inputs to GHCR and record immutable digests."""

import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
from release_image_archive import uncompressed_layer_bytes

from release_checks import wait_checks


IMAGE_NAMES = ['backend', 'landing', 'dashboard', 'restaurant-website', 'super-admin']
REGISTRY = 'ghcr.io/yousef-abdallah-1711/pos'


def inspect_image(name: str, tag: str) -> dict:
    result = subprocess.check_output(['docker', 'image', 'inspect', tag], text=True)
    images = json.loads(result)
    if len(images) != 1:
        raise RuntimeError('Expected exactly one image: ' + name)
    return images[0]


def publish_images(bundle: Path, sha: str, actor: str, token: str) -> dict:
    manifest_path = bundle / 'release-manifest.json'
    manifest = json.loads(manifest_path.read_text(encoding='utf-8'))
    expected_files = {'release-manifest.json', 'source.tar.gz'} | {
        name + '.docker.tar.gz' for name in IMAGE_NAMES
    }
    if {path.name for path in bundle.iterdir()} != expected_files:
        raise RuntimeError('Unexpected staged release bundle contents')
    if manifest.get('sha') != sha or manifest.get('environment') != 'staging':
        raise RuntimeError('Staged release identity mismatch')
    if not isinstance(manifest.get('ci_run_id'), int) or manifest['ci_run_id'] <= 0:
        raise RuntimeError('Invalid CI run identity')
    if set(manifest.get('archives', {})) != set(IMAGE_NAMES):
        raise RuntimeError('Incomplete image archive metadata')
    if set(manifest.get('image_ids', {})) != set(IMAGE_NAMES):
        raise RuntimeError('Incomplete immutable image identities')
    if set(manifest.get('image_sizes', {})) != set(IMAGE_NAMES):
        raise RuntimeError('Incomplete immutable image size metadata')
    if not re.fullmatch('[0-9a-f]{64}', manifest.get('source_sha256', '')):
        raise RuntimeError('Invalid source archive identity')
    source = bundle / 'source.tar.gz'
    if not 0 < source.stat().st_size <= 150 * 1024**2:
        raise RuntimeError('Source archive exceeds the deployment bound')
    if manifest.get('source_size') != source.stat().st_size:
        raise RuntimeError('Source archive size metadata mismatch')
    with source.open('rb') as stream:
        if hashlib.file_digest(stream, 'sha256').hexdigest() != manifest['source_sha256']:
            raise RuntimeError('Source archive checksum mismatch')
    for name in IMAGE_NAMES:
        archive = bundle / (name + '.docker.tar.gz')
        recorded = manifest['archives'][name]
        if recorded.get('file') != archive.name:
            raise RuntimeError('Unexpected image archive name: ' + name)
        if (
            not isinstance(recorded.get('size'), int)
            or not 0 < recorded['size'] <= 2 * 1024**3
            or recorded.get('size') != archive.stat().st_size
        ):
            raise RuntimeError('Image archive size mismatch: ' + name)
        if not re.fullmatch('[0-9a-f]{64}', recorded.get('sha256', '')):
            raise RuntimeError('Invalid image archive checksum: ' + name)
        with archive.open('rb') as stream:
            actual_digest = hashlib.file_digest(stream, 'sha256').hexdigest()
        if actual_digest != recorded['sha256']:
            raise RuntimeError('Image archive checksum mismatch: ' + name)

    subprocess.run(
        ['docker', 'login', 'ghcr.io', '-u', actor, '--password-stdin'],
        input=token,
        text=True,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        check=True,
    )
    try:
        manifest['images'] = {}
        for name in IMAGE_NAMES:
            archive = bundle / (name + '.docker.tar.gz')
            subprocess.run(['docker', 'load', '--input', str(archive)], check=True)
            local_tag = f'servio-release-{name}:{sha}'
            image = inspect_image(name, local_tag)
            if image.get('Id') != manifest['image_ids'][name]:
                raise RuntimeError('Loaded image identity mismatch: ' + name)
            if uncompressed_layer_bytes(archive) != manifest['image_sizes'][name]:
                raise RuntimeError('Archived image size metadata mismatch: ' + name)
            labels = image.get('Config', {}).get('Labels') or {}
            if labels.get('org.opencontainers.image.revision') != sha:
                raise RuntimeError('Loaded image revision mismatch: ' + name)

            target = f'{REGISTRY}/{name}:sha-{sha}'
            subprocess.run(['docker', 'tag', local_tag, target], check=True)
            subprocess.run(['docker', 'push', target], check=True)
            published = inspect_image(name, target)
            repo_prefix = f'{REGISTRY}/{name}@sha256:'
            digests = [value for value in published.get('RepoDigests', []) if value.startswith(repo_prefix)]
            if len(digests) != 1 or not re.fullmatch(re.escape(repo_prefix) + '[0-9a-f]{64}', digests[0]):
                raise RuntimeError('Registry did not return one immutable image digest: ' + name)
            manifest['images'][name] = digests[0]
    finally:
        subprocess.run(['docker', 'logout', 'ghcr.io'], stdout=subprocess.DEVNULL, check=True)

    manifest_path.write_text(
        json.dumps(manifest, sort_keys=True, separators=(',', ':')) + '\n',
        encoding='utf-8',
    )
    return manifest


def main() -> int:
    sha = os.environ.get('GITHUB_SHA', '')
    if (
        os.environ.get('GITHUB_ACTIONS') != 'true'
        or os.environ.get('GITHUB_REPOSITORY') != 'Yousef-Abdallah-1711/pos'
        or os.environ.get('GITHUB_EVENT_NAME') != 'push'
        or os.environ.get('GITHUB_REF') != 'refs/heads/main'
        or not re.fullmatch('[0-9a-f]{40}', sha)
    ):
        raise SystemExit('Refusing image publication outside an exact main push')
    token = os.environ.get('GHCR_TOKEN', '')
    if not token:
        raise SystemExit('A scoped GHCR publishing token is required')
    wait_checks(sha, int(os.environ['GITHUB_RUN_ID']))
    manifest = publish_images(Path('release-bundle'), sha, os.environ['GITHUB_ACTOR'], token)
    print('PASS immutable main image digests recorded for ' + sha + ' (CI run ' + str(manifest['ci_run_id']) + ')')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
