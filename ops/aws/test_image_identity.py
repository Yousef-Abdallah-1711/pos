"""Docker classic/containerd image identities must stay bound to verified content."""
import hashlib
import importlib.util
import io
import json
from pathlib import Path
import tarfile
import tempfile
import unittest
from unittest.mock import patch

spec=importlib.util.spec_from_file_location('image_deployer',Path(__file__).with_name('deploy-release.py'))
deploy=importlib.util.module_from_spec(spec);spec.loader.exec_module(deploy)

class ImageIdentityTests(unittest.TestCase):
    def archive(self,oci=True,tamper=False):
        temporary=tempfile.TemporaryDirectory();self.addCleanup(temporary.cleanup)
        path=Path(temporary.name)/'image.tar.gz'
        config=b'{"architecture":"amd64","config":{}}'
        config_id='sha256:'+hashlib.sha256(config).hexdigest()
        config_path='blobs/sha256/'+config_id[7:]
        manifest=json.dumps({'config':{'digest':config_id},'layers':[]}).encode()
        manifest_id='sha256:'+hashlib.sha256(manifest).hexdigest()
        members={'manifest.json':json.dumps([{'Config':config_path}]).encode(),config_path:config+b' ' if tamper else config}
        if oci:
            members['index.json']=json.dumps({'manifests':[{'digest':manifest_id}]}).encode()
            members['blobs/sha256/'+manifest_id[7:]]=manifest
        with tarfile.open(path,'w:gz') as output:
            for name,data in members.items():
                member=tarfile.TarInfo(name);member.size=len(data);output.addfile(member,io.BytesIO(data))
        return path,config_id,manifest_id

    def test_classic_store_accepts_only_verified_config_identity(self):
        path,config,_=self.archive(oci=False)
        self.assertEqual(deploy.archive_identities(path,config),{config})

    def test_containerd_manifest_is_verified_and_linked_to_expected_config(self):
        path,config,manifest=self.archive()
        self.assertEqual(deploy.archive_identities(path,config),{config,manifest})

    def test_modified_config_cannot_borrow_a_valid_oci_identity(self):
        path,config,_=self.archive(tamper=True)
        with self.assertRaisesRegex(RuntimeError,'config digest mismatch'):
            deploy.archive_identities(path,config)

    def test_loaded_image_outside_verified_identity_set_is_rejected(self):
        manifest={'sha':'a'*40,'image_ids':{'backend':'sha256:'+'b'*64},'images':{'backend':'ghcr.io/yousef-abdallah-1711/pos/backend@sha256:'+'c'*64},'archives':{'backend':{'file':'backend.docker.tar.gz'}}}
        loaded=[{'Id':'sha256:'+'d'*64,'Config':{'Labels':{'org.opencontainers.image.revision':'a'*40}}}]
        with patch.object(deploy,'SERVICES',['backend']),patch.object(deploy,'archive_identities',return_value={'sha256:'+'b'*64}),patch.object(deploy,'checked'),patch.object(deploy.subprocess,'check_output',return_value=json.dumps(loaded)):
            with self.assertRaisesRegex(RuntimeError,'Loaded image identity mismatch'):
                deploy.load_verified_images(manifest,Path('/unused'))

class ExtractedSourcePermissionTests(unittest.TestCase):
    def test_private_agent_umask_does_not_hide_public_source_from_acceptance_user(self):
        import os
        with tempfile.TemporaryDirectory() as directory:
            root=Path(directory);destination=root/'release';destination.mkdir()
            with tarfile.open(root/'source.tar.gz','w:gz') as output:
                member=tarfile.TarInfo('repo/ops/aws/smoke.py');member.size=4;member.mode=0o644
                output.addfile(member,io.BytesIO(b'pass'))
            previous=os.umask(0o077)
            try:deploy.extract_source('a'*40,destination,root)
            finally:os.umask(previous)
            self.assertEqual((destination/'ops').stat().st_mode&0o777,0o755)
            self.assertEqual((destination/'ops/aws').stat().st_mode&0o777,0o755)
            self.assertEqual((destination/'ops/aws/smoke.py').stat().st_mode&0o777,0o644)

if __name__=='__main__':unittest.main()
