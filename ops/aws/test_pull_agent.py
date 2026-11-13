"""Cryptographic and provenance regression tests without network or live infrastructure."""
import hashlib
import importlib.util
import io
import json
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch

spec=importlib.util.spec_from_file_location('servio_agent',Path(__file__).with_name('pull-agent.py'))
agent=importlib.util.module_from_spec(spec)
spec.loader.exec_module(agent)
SHA='a'*40

class PullAgentTrustTests(unittest.TestCase):
    def setUp(self):
        self.temporary=tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root=Path(self.temporary.name)
        self.private=self.root/'key.pem';self.public=self.root/'public.pem'
        subprocess.run(['openssl','genpkey','-algorithm','ED25519','-out',str(self.private)],check=True,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
        subprocess.run(['openssl','pkey','-in',str(self.private),'-pubout','-out',str(self.public)],check=True,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
        self.manifest={'sha':SHA,'environment':'staging','ci_run_id':123,'schema_version':1,'workflow_ref':agent.REPO+'/.github/workflows/production-ci.yml@refs/heads/main','workflow_sha':SHA}
        self.raw=json.dumps(self.manifest).encode()
        source=self.root/'manifest';source.write_bytes(self.raw)
        signature=self.root/'signature'
        subprocess.run(['openssl','pkeyutl','-sign','-rawin','-inkey',str(self.private),'-in',str(source),'-out',str(signature)],check=True)
        self.signature=signature.read_bytes()
        self.release={'target_commitish':SHA,'tag_name':'staging-sha-'+SHA,'draft':False,'assets':[{'name':name,'browser_download_url':'https://github.com/'+agent.REPO+'/releases/download/staging-sha-'+SHA+'/'+name} for name in ['release-manifest.json','release-manifest.sig']]}
        self.run={'workflow_id':1234,'head_sha':SHA,'head_branch':'main','event':'push','conclusion':'success','status':'completed','path':'.github/workflows/production-ci.yml'}
        (self.root/'workflow-id').write_text('1234')
        for name,value in [('WORKFLOW_ID',self.root/'workflow-id'),('BASE',self.root/'state'),('KEY',self.public),('STATE',self.root/'deployed.json')]:
            patcher=patch.object(agent,name,value);patcher.start();self.addCleanup(patcher.stop)

    def verify(self,raw=None,run=None):
        with patch.object(agent,'github',side_effect=[{'object':{'sha':SHA}},run or self.run]),patch.object(agent,'request',side_effect=[raw or self.raw,self.signature]):
            return agent.verify_release(self.release)

    def test_accepts_real_signature_for_exact_successful_main(self):
        manifest,_=self.verify()
        self.assertEqual(manifest['sha'],SHA)

    def test_rejects_manifest_tampering_before_artifact_use(self):
        with self.assertRaisesRegex(RuntimeError,'signature'):
            self.verify(raw=self.raw+b' ')

    def test_rejects_failed_wrong_branch_and_untrusted_workflow(self):
        for key,value in [('conclusion','failure'),('head_branch','feature'),('event','pull_request'),('path','.github/workflows/other.yml')]:
            with self.subTest(key=key),self.assertRaisesRegex(RuntimeError,'provenance'):
                self.verify(run=dict(self.run,**{key:value}))

    def test_pending_ci_never_promotes(self):
        with self.assertRaises(agent.PendingRelease):
            self.verify(run=dict(self.run,status='in_progress',conclusion=None))

    def test_stale_release_does_not_even_download_manifest(self):
        with patch.object(agent,'github',return_value={'object':{'sha':'b'*40}}),patch.object(agent,'request') as download:
            with self.assertRaises(agent.PendingRelease):agent.verify_release(self.release)
            download.assert_not_called()

    def test_checksum_failure_never_installs_download(self):
        release={'tag_name':'staging-sha-'+SHA,'assets':[{'name':'backend.docker.tar.gz','browser_download_url':'https://github.com/'+agent.REPO+'/releases/download/staging-sha-'+SHA+'/backend.docker.tar.gz'}]}
        target=self.root/'backend.docker.tar.gz'
        with patch.object(agent.urllib.request,'urlopen',return_value=io.BytesIO(b'ac')):
            with self.assertRaisesRegex(RuntimeError,'checksum'):
                agent.download(release,'backend.docker.tar.gz',target,2,hashlib.sha256(b'ab').hexdigest())
        self.assertFalse(target.exists())

    def test_bounded_download_refuses_extra_bytes(self):
        release={'tag_name':'staging-sha-'+SHA,'assets':[{'name':'backend.docker.tar.gz','browser_download_url':'https://github.com/'+agent.REPO+'/releases/download/staging-sha-'+SHA+'/backend.docker.tar.gz'}]}
        with patch.object(agent.urllib.request,'urlopen',return_value=io.BytesIO(b'abc')):
            with self.assertRaisesRegex(RuntimeError,'exceeds'):
                agent.download(release,'backend.docker.tar.gz',self.root/'archive',2,hashlib.sha256(b'ab').hexdigest())

    def test_release_cannot_redirect_to_another_artifact_origin(self):
        self.release['assets'][0]['browser_download_url']='https://attacker.example/manifest'
        with self.assertRaisesRegex(RuntimeError,'untrusted'):
            agent.asset_url(self.release,'release-manifest.json')

    def test_retry_limit_stops_additional_deployment_attempts(self):
        agent.BASE.mkdir()
        (agent.BASE/'attempts.json').write_text(json.dumps({SHA:{'count':3}}))
        with patch.object(agent,'github',return_value=[self.release]),patch.object(agent,'verify_release') as verify:
            with self.assertRaisesRegex(RuntimeError,'Retry limit'):agent.poll()
            verify.assert_not_called()

if __name__=='__main__':unittest.main()
