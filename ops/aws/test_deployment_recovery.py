"""Breaking guard migrations must not restore an incompatible journal writer."""
import importlib.util
import json
from pathlib import Path
import tempfile
from types import SimpleNamespace
import unittest
from unittest.mock import Mock,patch

import deployment_recovery as recovery

spec=importlib.util.spec_from_file_location('recovery_deployer',Path(__file__).with_name('deploy-release.py'))
deploy=importlib.util.module_from_spec(spec);spec.loader.exec_module(deploy)
SHA='a'*40
OLD='sha256:'+'b'*64
NEW='sha256:'+'c'*64

class MigrationRecoveryTests(unittest.TestCase):
    def policy(self,mode=0o100644,owner=0):
        policy=Mock();policy.exists.return_value=True
        policy.stat.return_value=SimpleNamespace(st_uid=owner,st_mode=mode)
        policy.read_text.return_value=json.dumps({'guard.sql':'keep_backend'})
        return policy

    def test_only_reviewed_pending_contracts_need_forward_backend(self):
        with patch.object(recovery,'POLICY',self.policy()):
            self.assertTrue(recovery.requires_compatible_backend(['guard.sql'],{}))
            self.assertFalse(recovery.requires_compatible_backend(['ordinary.sql'],{}))
            self.assertTrue(recovery.requires_compatible_backend([] ,{'recovery_required':True}))

    def test_untrusted_or_writable_policy_is_rejected(self):
        for policy in [self.policy(owner=1000),self.policy(mode=0o100666)]:
            with patch.object(recovery,'POLICY',policy),self.assertRaisesRegex(RuntimeError,'root-owned'):
                recovery.requires_compatible_backend(['guard.sql'],{})

    def fixture(self,root):
        old=root/'compose.json'
        config={'services':{name:{'image':OLD,'volumes':['./Caddyfile.public:/etc/caddy/Caddyfile:ro']} for name in ['backend','migrate','landing']}}
        old.write_text(json.dumps(config))
        previous={'sha':'d'*40,'compose':str(old),'images':{'backend':'old-registry'},'image_ids':{'backend':OLD},'runtime_image_ids':{'backend':OLD,'landing':OLD}}
        attempt={'manifest':{'sha':SHA,'images':{'backend':'new-registry'},'image_ids':{'backend':NEW}},'runtime':{'backend':NEW},'backend_ready':True}
        return previous,attempt,config

    def test_recovery_keeps_compatible_backend_and_previous_frontend(self):
        with tempfile.TemporaryDirectory() as directory:
            previous,attempt,original=self.fixture(Path(directory))
            state,fallback=recovery.compatible_recovery(previous,attempt)
            actual=json.loads(fallback.read_text())
            self.assertEqual(actual['services']['backend']['image'],NEW)
            self.assertEqual(actual['services']['landing']['image'],OLD)
            self.assertEqual(fallback.parent,Path(previous['compose']).parent)
            self.assertEqual(json.loads(Path(previous['compose']).read_text()),original)
            self.assertTrue(state['recovery_required'])
            self.assertEqual(state['runtime_image_ids']['backend'],NEW)

    def test_failed_new_backend_restores_last_compatible_recovery_backend(self):
        with tempfile.TemporaryDirectory() as directory:
            previous,attempt,_=self.fixture(Path(directory))
            previous['recovery_required']=True;attempt['backend_ready']=False
            state,fallback=recovery.compatible_recovery(previous,attempt)
            self.assertEqual(json.loads(fallback.read_text())['services']['backend']['image'],OLD)
            self.assertEqual(state['runtime_image_ids']['backend'],OLD)

    def test_mixed_recovery_never_advertises_an_accepted_commit(self):
        with tempfile.TemporaryDirectory() as directory:
            root=Path(directory);previous,attempt,_=self.fixture(root)
            with patch.object(deploy,'STATE',root/'state.json'),patch.object(deploy,'PUBLIC_STATUS',root/'public.json'),patch.object(deploy,'checked'),patch.object(deploy,'rollback_apps') as rollback:
                deploy.recover_failed_deployment(previous,attempt,True)
                receipt=json.loads((root/'public.json').read_text())
                self.assertIsNone(receipt['sha'])
                self.assertEqual(receipt['status'],'recovery_required')
                self.assertTrue(json.loads((root/'state.json').read_text())['recovery_required'])
                rollback.assert_called_once()

    def test_contract_writer_starts_before_migration_and_failure_preserves_it(self):
        from contextlib import ExitStack
        with tempfile.TemporaryDirectory() as directory:
            root=Path(directory);ops=root/'ops/aws';ops.mkdir(parents=True)
            previous,attempt,config=self.fixture(ops)
            config['services']={name:{'image':OLD,'volumes':[]} for name in deploy.SERVICES+['postgres','redis','proxy','migrate']}
            Path(previous['compose']).write_text(json.dumps(config))
            state=ops/'state.json';state.write_text(json.dumps(previous))
            manifest={'sha':SHA,'image_ids':{name:NEW for name in deploy.SERVICES}}
            calls=[]
            def extract(sha,destination,artifacts):
                target=destination/'ops/aws';target.mkdir(parents=True);(target/'compose.json').write_text(json.dumps(config))
            patches=[
                patch.object(deploy,'ROOT',root),
                patch.object(deploy,'OPS',ops),
                patch.object(deploy,'STATE',state),
                patch.object(deploy,'verify_manifest',return_value=SHA),
                patch.object(deploy,'verify_artifacts'),
                patch.object(deploy,'extract_source',side_effect=extract),
                patch.object(deploy,'verify_migrations',return_value=['guard.sql']),
                patch.object(deploy,'requires_compatible_backend',return_value=True),
                patch.object(deploy,'load_verified_images',return_value=manifest['image_ids']),
                patch.object(deploy,'backup'),
                patch.object(deploy,'compose',side_effect=lambda path,*args,**kwargs:calls.append(args)),
                patch.object(deploy,'public_checks',side_effect=RuntimeError('readiness failed')),
                patch.object(deploy,'recover_failed_deployment'),
                patch.object(deploy.shutil,'disk_usage',return_value=SimpleNamespace(free=20*1024**3))
            ]
            with ExitStack() as stack:
                mocks=[stack.enter_context(item) for item in patches]
                with self.assertRaisesRegex(RuntimeError,'readiness'):deploy.deploy(manifest,SHA,root)
                self.assertEqual(calls[1][0],'up')
                self.assertEqual(calls[1][-1],'backend')
                self.assertEqual(calls[2][0],'run')
                recovery_call=mocks[-2].call_args
                self.assertTrue(recovery_call.args[1]['backend_ready'])
                self.assertTrue(recovery_call.args[2])
