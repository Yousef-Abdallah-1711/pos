"""Deployment trust-boundary tests; no Docker, SSH, database or network access."""
import copy
import importlib.util
import json
from pathlib import Path
import unittest
from types import SimpleNamespace
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
    def test_staging_headroom_rechecks_signed_image_budget(self):
        manifest={'image_sizes':{name:100*1024**2 for name in deploy.SERVICES}}
        with patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=7*1024**3)):
            deploy.verify_disk_headroom(manifest)
        with patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=5*1024**3)):
            with self.assertRaisesRegex(RuntimeError,'Insufficient deployment disk headroom'):
                deploy.verify_disk_headroom(manifest)
        with patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=5*1024**3)):
            deploy.verify_disk_headroom(manifest,after_images=True)
        with patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=4*1024**3)):
            with self.assertRaisesRegex(RuntimeError,'Insufficient deployment disk headroom'):
                deploy.verify_disk_headroom(manifest,after_images=True)

    def test_fresh_backup_must_restore_before_it_is_accepted(self):
        import tempfile
        with tempfile.TemporaryDirectory() as directory:
            root=Path(directory);ops=root/'ops/aws';(ops/'backups').mkdir(parents=True)
            calls=[]
            def compose(path,*args,**kwargs):
                calls.append(args)
                if args[0:3]==('exec','-T','postgres') and args[3]=='pg_dump':
                    kwargs['stdout'].write(b'verified synthetic dump')
                if args[0:3]==('exec','-T','postgres') and args[3]=='psql':
                    return SimpleNamespace(stdout='3\n')
                return SimpleNamespace(stdout='')
            with patch.object(deploy,'OPS',ops),patch.object(deploy,'compose',side_effect=compose),patch.object(deploy,'checked'):
                deploy.backup(Path('/tmp/compose.json'))
            dumps=list((ops/'backups').glob('*.dump'))
            self.assertEqual(len(dumps),1)
            self.assertTrue(dumps[0].with_suffix('.dump.sha256').exists())
            receipt=json.loads(dumps[0].with_suffix('.restore.json').read_text())
            self.assertEqual(receipt['status'],'PASS')
            self.assertEqual(receipt['restored_migration_receipts'],3)
            self.assertTrue(any(call[-1].startswith('servio_restore_verify_') for call in calls if call and call[0]=='exec'))

    def test_failed_post_deployment_check_restores_apps_without_reversing_sql(self):
        import tempfile
        from types import SimpleNamespace
        with tempfile.TemporaryDirectory() as temporary:
            root=Path(temporary);ops=root/'ops/aws';ops.mkdir(parents=True)
            services={name:{'image':'old-'+name,'volumes':[]} for name in deploy.SERVICES+['postgres','redis','proxy','migrate']}
            old={'services':services}
            previous=ops/'compose.json';previous.write_text(__import__('json').dumps(old))
            status=ops/'public/release.json';status.parent.mkdir();status.write_text('{"sha":"'+'c'*40+'","status":"verified"}\n')
            manifest={'sha':SHA,'ci_run_id':123,'images':{},'image_ids':{name:'sha256:'+'b'*64 for name in deploy.SERVICES},'image_sizes':{name:100*1024**2 for name in deploy.SERVICES}}
            calls=[]
            def extract(sha,destination,artifacts):
                target=destination/'ops/aws';target.mkdir(parents=True)
                (target/'compose.json').write_text(__import__('json').dumps(old))
            def compose(path,*args,**kwargs):calls.append((path,args))
            patches=[patch.object(deploy,'ROOT',root),patch.object(deploy,'OPS',ops),patch.object(deploy,'STATE',ops/'state.json'),patch.object(deploy,'PUBLIC_STATUS',status),patch.object(deploy,'verify_manifest',return_value=SHA),patch.object(deploy,'verify_artifacts'),patch.object(deploy,'extract_source',side_effect=extract),patch.object(deploy,'load_verified_images',return_value=manifest['image_ids']),patch.object(deploy,'verify_migrations'),patch.object(deploy,'backup'),patch.object(deploy,'compose',side_effect=compose),patch.object(deploy,'public_checks',side_effect=RuntimeError('simulated readiness failure')),patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=20*1024**3))]
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

class CandidateTrustTests(unittest.TestCase):
    def test_operator_candidate_requires_successful_owner_workflow_and_current_main(self):
        run={'head_sha':SHA,'event':'workflow_dispatch','conclusion':'success','path':'.github/workflows/production-ci.yml','status':'completed','actor':{'login':'Yousef-Abdallah-1711'},'workflow_id':415727583,'head_branch':'candidate'}
        with patch.object(deploy,'github',side_effect=[{'object':{'sha':SHA}},{'status':'ahead'}]):
            deploy.verify_candidate(run,SHA,'b'*40)
        for key,value in [('event','pull_request'),('conclusion','failure'),('actor',{'login':'foreign'}),('workflow_id',99)]:
            with self.subTest(key=key),self.assertRaises(RuntimeError):
                deploy.verify_candidate(dict(run,**{key:value}),SHA,'b'*40)
        with patch.object(deploy,'github',side_effect=[{'object':{'sha':SHA}},{'status':'diverged'}]):
            with self.assertRaisesRegex(RuntimeError,'include current main'):
                deploy.verify_candidate(run,SHA,'b'*40)

class MetadataAtomicityTests(unittest.TestCase):
    def test_failed_public_receipt_restores_previous_deployment_identity(self):
        import tempfile
        with tempfile.TemporaryDirectory() as directory:
            state=Path(directory)/'state.json';status=Path(directory)/'release.json'
            state.write_bytes(b'previous private state');status.write_bytes(b'previous public receipt')
            original=deploy.replace_metadata;failed=False
            def replace(file,raw,mode,owner=False):
                nonlocal failed
                if file==status and not failed:
                    failed=True;raise OSError('simulated metadata write failure')
                original(file,raw,mode,owner)
            manifest={'sha':SHA,'images':{},'image_ids':{},'runtime_image_ids':{},'ci_run_id':123}
            with patch.object(deploy,'STATE',state),patch.object(deploy,'PUBLIC_STATUS',status),patch.object(deploy,'checked'),patch.object(deploy,'replace_metadata',side_effect=replace):
                with self.assertRaisesRegex(OSError,'metadata write'):
                    deploy.record_release(manifest,Path('/release.json'),{})
            self.assertEqual(state.read_bytes(),b'previous private state')
            self.assertEqual(status.read_bytes(),b'previous public receipt')

class CandidateReceiptTests(unittest.TestCase):
    def test_candidate_receipt_is_published_and_previous_receipt_can_be_restored(self):
        import tempfile
        with tempfile.TemporaryDirectory() as directory:
            status = Path(directory) / 'release.json'
            previous = b'{"sha":"' + b'c' * 40 + b'","status":"verified"}\n'
            status.write_bytes(previous)
            manifest = {'sha': SHA, 'ci_run_id': 456}
            with patch.object(deploy, 'PUBLIC_STATUS', status):
                deploy.publish_candidate_receipt(manifest)
                self.assertEqual(json.loads(status.read_text()), {
                    'sha': SHA, 'status': 'candidate', 'environment': 'staging', 'ci_run_id': 456,
                })
                deploy.restore_public_receipt(previous)
            self.assertEqual(status.read_bytes(), previous)

    def test_failed_staging_acceptance_restores_previous_receipt_after_app_rollback(self):
        import tempfile
        from contextlib import ExitStack
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            ops = root / 'ops/aws'
            ops.mkdir(parents=True)
            previous_compose = ops / 'compose.json'
            services = {name: {'image': 'old-' + name, 'volumes': []}
                        for name in deploy.SERVICES + ['postgres', 'redis', 'proxy', 'migrate']}
            previous_compose.write_text(json.dumps({'services': services}))
            state = ops / 'state.json'
            state.write_text(json.dumps({'sha': 'c' * 40, 'compose': str(previous_compose)}))
            status = ops / 'public/release.json'
            status.parent.mkdir()
            prior_receipt = b'{"sha":"' + b'c' * 40 + b'","status":"verified"}\n'
            status.write_bytes(prior_receipt)
            manifest = {
                'sha': SHA, 'ci_run_id': 456,
                'image_ids': {name: 'sha256:' + 'b' * 64 for name in deploy.SERVICES},
                'image_sizes': {name: 100 * 1024**2 for name in deploy.SERVICES},
            }

            def extract(sha, destination, artifacts):
                target = destination / 'ops/aws'
                target.mkdir(parents=True)
                (target / 'compose.json').write_text(json.dumps({'services': services}))

            def fail_acceptance(destination, expected_sha):
                self.assertEqual(json.loads(status.read_text())['sha'], expected_sha)
                self.assertEqual(json.loads(status.read_text())['status'], 'candidate')
                raise RuntimeError('simulated staging-safe acceptance failure')

            patches = [
                patch.object(deploy, 'ROOT', root), patch.object(deploy, 'OPS', ops),
                patch.object(deploy, 'STATE', state), patch.object(deploy, 'PUBLIC_STATUS', status),
                patch.object(deploy, 'verify_manifest', return_value=SHA),
                patch.object(deploy, 'verify_artifacts'), patch.object(deploy, 'extract_source', side_effect=extract),
                patch.object(deploy, 'verify_migrations', return_value=[]),
                patch.object(deploy, 'requires_compatible_backend', return_value=False),
                patch.object(deploy, 'load_verified_images', return_value=manifest['image_ids']),
                patch.object(deploy, 'backup'), patch.object(deploy, 'public_checks'),
                patch.object(deploy, 'acceptance', side_effect=fail_acceptance),
                patch.object(deploy, 'compose'),
                patch.object(deploy.shutil, 'disk_usage', return_value=SimpleNamespace(free=20 * 1024**3)),
            ]
            with ExitStack() as stack:
                for mocked in patches:
                    stack.enter_context(mocked)
                with self.assertRaisesRegex(RuntimeError, 'simulated staging-safe acceptance failure'):
                    deploy.deploy(manifest, SHA, root)
            self.assertEqual(status.read_bytes(), prior_receipt)
            self.assertEqual(json.loads(state.read_text())['sha'], 'c' * 40)

if __name__ == '__main__':unittest.main()
