"""Unit checks for credential-separated, exact-main release image handling."""

import hashlib
import importlib.util
import io
import json
from pathlib import Path
import sys
import tempfile
import tarfile
import unittest
from unittest.mock import patch


ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / 'ops/aws'))


def load_script(name: str, module_name: str):
    spec = importlib.util.spec_from_file_location(module_name, ROOT / 'ops/aws' / name)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


stager = load_script('stage-ci-images.py', 'servio_stage_ci_images')
publisher = load_script('publish-ci-images.py', 'servio_publish_ci_images')


def write_docker_archive(path: Path, name: str, sha: str, image_id: str) -> int:
    config = json.dumps({
        'architecture': 'amd64',
        'os': 'linux',
        'config': {'Labels': {'org.opencontainers.image.revision': sha}},
        'rootfs': {'type': 'layers', 'diff_ids': []},
    }).encode()
    config_name = image_id.removeprefix('sha256:') + '.json'
    layer_name = 'layer/layer.tar'
    manifest = json.dumps([{'Config': config_name, 'RepoTags': [f'servio-ci-{name}:{sha}'], 'Layers': [layer_name]}]).encode()
    with tarfile.open(path, 'w:gz') as archive:
        for member_name, data in [('manifest.json', manifest), (config_name, config), (layer_name, b'x' * 8192)]:
            info = tarfile.TarInfo(member_name)
            info.size = len(data)
            archive.addfile(info, io.BytesIO(data))
    return 8192


class ReleaseImagePipelineTests(unittest.TestCase):
    def test_size_budget_uses_stable_uncompressed_archive_layer_bytes(self):
        with tempfile.TemporaryDirectory() as directory:
            archive = Path(directory) / 'backend.docker.tar.gz'
            write_docker_archive(archive, 'backend', 'd' * 40, 'sha256:' + 'a' * 64)
            self.assertEqual(publisher.uncompressed_layer_bytes(archive), 8192)

    def test_stager_creates_fixed_set_and_exact_ci_identity(self):
        sha = 'a' * 40
        with tempfile.TemporaryDirectory() as directory:
            bundle = Path(directory) / 'release-bundle'

            def fake_archive(target: Path, name: str, actual_sha: str):
                archive = target / (name + '.docker.tar.gz')
                archive.write_bytes(('archive-' + name).encode())
                return archive, 'sha256:' + name.encode().hex().ljust(64, '0')[:64], 123456

            def fake_source(_url: str, destination: Path):
                Path(destination).write_bytes(b'synthetic source archive')
                return str(destination)

            with patch.object(stager, 'archive_image', side_effect=fake_archive), patch.object(
                stager.urllib.request, 'urlretrieve', side_effect=fake_source
            ):
                manifest = stager.stage_bundle(bundle, sha, 12345)

            self.assertEqual(manifest['sha'], sha)
            self.assertEqual(manifest['ci_run_id'], 12345)
            self.assertEqual(manifest['schema_version'], 3)
            self.assertEqual(set(manifest['image_sizes']), set(stager.IMAGE_NAMES))
            self.assertGreater(manifest['source_size'], 0)
            self.assertEqual(set(manifest['archives']), set(stager.IMAGE_NAMES))
            self.assertEqual(set(path.name for path in bundle.iterdir()), {
                'release-manifest.json', 'source.tar.gz',
                *(name + '.docker.tar.gz' for name in stager.IMAGE_NAMES),
            })

    def test_publisher_loads_only_verified_image_archives_and_records_digests(self):
        sha = 'b' * 40
        image_ids = {name: 'sha256:' + str(index) * 64 for index, name in enumerate(publisher.IMAGE_NAMES, 1)}
        registry_digests = {
            name: f'{publisher.REGISTRY}/{name}@sha256:' + format(index + 5, 'x') * 64
            for index, name in enumerate(publisher.IMAGE_NAMES, 1)
        }
        with tempfile.TemporaryDirectory() as directory:
            bundle = Path(directory)
            archives = {}
            for name in publisher.IMAGE_NAMES:
                archive = bundle / (name + '.docker.tar.gz')
                write_docker_archive(archive, name, sha, image_ids[name])
                archives[name] = {
                    'file': archive.name,
                    'sha256': hashlib.sha256(archive.read_bytes()).hexdigest(),
                    'size': archive.stat().st_size,
                }
            source = bundle / 'source.tar.gz'
            source.write_bytes(b'source')
            manifest = {
                'schema_version': 3,
                'environment': 'staging',
                'sha': sha,
                'ci_run_id': 12345,
                'images': {},
                'image_ids': image_ids,
                'image_sizes': {name: 8192 for name in publisher.IMAGE_NAMES},
                'archives': archives,
                'source_sha256': hashlib.sha256(source.read_bytes()).hexdigest(),
                'source_size': source.stat().st_size,
            }
            (bundle / 'release-manifest.json').write_text(json.dumps(manifest))

            def fake_inspect(command, text):
                tag = command[-1]
                name = tag.removeprefix('servio-ci-').split('/')[-1].split(':')[0]
                image = {
                    'Id': image_ids[name],
                    # Docker Engine's reported Size is intentionally allowed to differ.
                    'Size': 9999,
                    'Config': {'Labels': {'org.opencontainers.image.revision': sha}},
                }
                if '/' in tag:
                    image['RepoDigests'] = [registry_digests[name]]
                return json.dumps([image])

            with patch.object(publisher.subprocess, 'check_output', side_effect=fake_inspect), patch.object(
                publisher.subprocess, 'run'
            ) as docker_run:
                result = publisher.publish_images(bundle, sha, 'servio-ci', 'synthetic-token')

            self.assertEqual(result['images'], registry_digests)
            docker_run.assert_any_call(['docker', 'push', f'{publisher.REGISTRY}/backend:sha-{sha}'], check=True)

    def test_publisher_rejects_modified_archive_before_registry_login(self):
        sha = 'c' * 40
        with tempfile.TemporaryDirectory() as directory:
            bundle = Path(directory)
            source = bundle / 'source.tar.gz'
            source.write_bytes(b'source')
            archives = {}
            for name in publisher.IMAGE_NAMES:
                archive = bundle / (name + '.docker.tar.gz')
                archive.write_bytes(b'changed')
                archives[name] = {'file': archive.name, 'sha256': '0' * 64, 'size': archive.stat().st_size}
            manifest = {
                'environment': 'staging', 'sha': sha, 'ci_run_id': 12345,
                'image_ids': {name: 'sha256:' + '1' * 64 for name in publisher.IMAGE_NAMES},
                'image_sizes': {name: 123456 for name in publisher.IMAGE_NAMES},
                'archives': archives,
                'source_sha256': hashlib.sha256(source.read_bytes()).hexdigest(),
                'source_size': source.stat().st_size,
            }
            (bundle / 'release-manifest.json').write_text(json.dumps(manifest))
            with patch.object(publisher.subprocess, 'run') as docker_run:
                with self.assertRaisesRegex(RuntimeError, 'archive checksum mismatch'):
                    publisher.publish_images(bundle, sha, 'servio-ci', 'synthetic-token')
            docker_run.assert_not_called()


if __name__ == '__main__':
    unittest.main()
