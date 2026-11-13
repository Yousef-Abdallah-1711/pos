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

class RollbackSafetyTests(unittest.TestCase):
    def test_failed_post_deployment_check_restores_apps_without_reversing_sql(self):
        import tempfile
        from types import SimpleNamespace
        with tempfile.TemporaryDirectory() as temporary:
            root=Path(temporary);ops=root/'ops/aws';ops.mkdir(parents=True)
            services={name:{'image':'old-'+name,'volumes':[]} for name in deploy.SERVICES+['postgres','redis','proxy','migrate']}
            old={'services':services}
            previous=ops/'compose.json';previous.write_text(__import__('json').dumps(old))
            manifest={'sha':SHA,'ci_run_id':123,'images':{},'image_ids':{name:'sha256:'+'b'*64 for name in deploy.SERVICES}}
            calls=[]
            def extract(sha,destination,artifacts):
                target=destination/'ops/aws';target.mkdir(parents=True)
                (target/'compose.json').write_text(__import__('json').dumps(old))
            def compose(path,*args,**kwargs):calls.append((path,args))
            patches=[patch.object(deploy,'ROOT',root),patch.object(deploy,'OPS',ops),patch.object(deploy,'STATE',ops/'state.json'),patch.object(deploy,'verify_manifest',return_value=SHA),patch.object(deploy,'verify_artifacts'),patch.object(deploy,'extract_source',side_effect=extract),patch.object(deploy,'load_verified_images'),patch.object(deploy,'backup'),patch.object(deploy,'compose',side_effect=compose),patch.object(deploy,'public_checks',side_effect=RuntimeError('simulated readiness failure')),patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=20*1024**3))]
            with __import__('contextlib').ExitStack() as stack:
                for mocked in patches:stack.enter_context(mocked)
                with self.assertRaisesRegex(RuntimeError,'simulated'):deploy.deploy(manifest,SHA,root)
            rollback=[args[-1] for path,args in calls if path==previous]
            self.assertEqual(rollback,deploy.SERVICES+['proxy'])
            self.assertEqual(sum(args[:1]==('run',) for _,args in calls),1)
            self.assertFalse(any('down' in args for _,args in calls))
            self.assertEqual(__import__('json').loads(previous.read_text()),old)

class RollbackContinuationTests(unittest.TestCase):
    def test_one_failed_restore_does_not_prevent_other_services_restoring(self):
        calls=[]
        def restore(path,*args):
            calls.append(args[-1])
            if args[-1]=='backend':raise RuntimeError('unhealthy')
        with patch.object(deploy,'compose',side_effect=restore):
            with self.assertRaisesRegex(RuntimeError,'rollback incomplete: backend'):
                deploy.rollback_apps(Path('/previous.json'))
        self.assertEqual(calls,deploy.SERVICES+['proxy'])

if __name__ == '__main__':unittest.main()
