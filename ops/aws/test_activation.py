"""Activation cannot approve a different release promoted by its provenance recheck."""
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

spec=importlib.util.spec_from_file_location('activate_agent',Path(__file__).with_name('activate-pull-agent.py'))
activate=importlib.util.module_from_spec(spec);spec.loader.exec_module(activate)
SHA='a'*40

class ActivationRaceTests(unittest.TestCase):
    def test_changed_release_cannot_write_marker_or_enable_timer(self):
        with tempfile.TemporaryDirectory() as directory:
            root=Path(directory);evidence=root/'ops/aws/evidence';evidence.mkdir(parents=True)
            release=root/'releases'/SHA/'ops/aws/evidence';(release/'browser').mkdir(parents=True)
            state={'sha':SHA,'candidate':False,'compose':str(release.parent/'compose.json')}
            state_path=evidence/'deployed-release.json';state_path.write_text(json.dumps(state))
            names=['database-role-smoke.json','api-smoke.json','refresh-smoke.json','admin-smoke.json','pos-kds-smoke.json','completed-sale-inventory.json','browser/results.json','storefront-browser-flow.json','upload-browser.json','journal-api-acceptance.json','auth-form-validation.json']
            for name in names:(release/name).write_text(json.dumps([{'status':'PASS'}]))
            (evidence/'full-browser').mkdir()
            (evidence/'full-browser/run.json').write_text(json.dumps({'scope':'full','sha':SHA,'observed_sha':SHA,'completed_utc':'now','failed':0,'passed':1,'expected_checks':1,'not_tested':0}))
            calls=[]
            def recheck(command):
                calls.append(command)
                if command[0].endswith('servio-pull-agent'):state_path.write_text(json.dumps(dict(state,sha='b'*40)))
            marker=root/'accepted'
            with patch.object(activate,'ROOT',root),patch.object(activate,'LOCK',root/'lock'),patch.object(activate,'MARKER',marker),patch.object(activate,'checked',side_effect=recheck):
                with self.assertRaisesRegex(RuntimeError,'Release changed'):activate.activate()
            self.assertFalse(marker.exists())
            self.assertFalse(any(command[0]=='systemctl' for command in calls))
