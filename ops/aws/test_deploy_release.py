"""Deployment trust-boundary tests; no Docker, SSH, database or network access."""
import copy
import importlib.util
from pathlib import Path
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('servio_deploy', Path(__file__).with_name('deploy-release.py'))
deploy = importlib.util.module_from_spec(spec)
spec.loader.exec_module(deploy)
SHA = 'a' * 40

class ReleaseTrustTests(unittest.TestCase):
    def setUp(self):
        self.manifest = {'sha': SHA, 'ci_run_id': 123, 'images': {name: 'ghcr.io/yousef-abdallah-1711/pos/' + name + '@sha256:' + 'b' * 64 for name in deploy.SERVICES}}
        self.run = {'head_sha': SHA, 'head_branch': 'main', 'event': 'push', 'conclusion': 'success', 'path': '.github/workflows/production-ci.yml'}

    def verify(self, manifest=None, run=None, main=SHA):
        with patch.object(deploy, 'github', side_effect=[{'object': {'sha': main}}, run or self.run]):
            return deploy.verify_manifest(manifest or self.manifest, SHA)

    def test_accepts_exact_current_main_with_successful_ci_and_all_digests(self):
        self.assertEqual(self.verify(), SHA)

    def test_rejects_stale_main_before_deploy(self):
        with self.assertRaisesRegex(RuntimeError, 'current main'):
            self.verify(main='c' * 40)

    def test_rejects_failed_or_untrusted_workflow_runs(self):
        for key, value in [('conclusion', 'failure'), ('event', 'pull_request'), ('head_branch', 'other'), ('path', '.github/workflows/untrusted.yml'), ('head_sha', 'c' * 40)]:
            with self.subTest(key=key):
                run = dict(self.run, **{key: value})
                with self.assertRaisesRegex(RuntimeError, 'acceptance'):
                    self.verify(run=run)

    def test_rejects_mutable_or_foreign_image_references(self):
        for image in ['ghcr.io/yousef-abdallah-1711/pos/backend:latest', 'attacker.example/backend@sha256:' + 'b' * 64]:
            with self.subTest(image=image):
                manifest = copy.deepcopy(self.manifest)
                manifest['images']['backend'] = image
                with self.assertRaisesRegex(RuntimeError, 'registry'):
                    self.verify(manifest=manifest)

    def test_rejects_missing_application_image(self):
        manifest = copy.deepcopy(self.manifest)
        del manifest['images']['backend']
        with self.assertRaisesRegex(RuntimeError, 'Incomplete'):
            self.verify(manifest=manifest)

    def test_rejects_requested_sha_mismatch_before_network(self):
        with patch.object(deploy, 'github') as network:
            with self.assertRaisesRegex(RuntimeError, 'identity'):
                deploy.verify_manifest(self.manifest, 'd' * 40)
            network.assert_not_called()

if __name__ == '__main__':
    unittest.main()
