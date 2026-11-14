"""Publication cannot outrun separate required guardrails or accept spoofed contexts."""
import copy
import importlib.util
from pathlib import Path
import unittest
from unittest.mock import patch

spec=importlib.util.spec_from_file_location('release_checks',Path(__file__).with_name('release_checks.py'))
gate=importlib.util.module_from_spec(spec);spec.loader.exec_module(gate)
SHA='a'*40

class RequiredPublicationTests(unittest.TestCase):
    def setUp(self):
        self.checks=[];self.runs={}
        for run_id,(name,(workflow_id,path)) in enumerate(gate.EXPECTED.items(),start=1):
            self.checks.append({'id':run_id,'name':name,'head_sha':SHA,'app':{'slug':'github-actions'},'status':'completed','conclusion':'success','details_url':f'https://github.com/{gate.REPO}/actions/runs/{run_id}/job/99'})
            self.runs[run_id]={'workflow_id':workflow_id,'path':path,'head_sha':SHA,'head_branch':'main','event':'push','status':'in_progress' if run_id==1 else 'completed','conclusion':None if run_id==1 else 'success'}
        self.mock=patch.object(gate,'api',side_effect=lambda path:self.runs[int(path.rsplit('/',1)[1])]);self.mock.start();self.addCleanup(self.mock.stop)

    def ready(self):return gate.checks_ready(self.checks,SHA,1)

    def test_accepts_exact_main_and_all_completed_guardrail_workflows(self):
        self.assertTrue(self.ready())

    def test_missing_required_context_waits_without_publication(self):
        self.checks.pop();self.assertFalse(self.ready())

    def test_running_required_context_waits_without_publication(self):
        self.checks[1]['status']='in_progress';self.assertFalse(self.ready())

    def test_failed_guardrail_is_fatal(self):
        self.checks[2]['conclusion']='failure'
        with self.assertRaisesRegex(RuntimeError,'required check failed'):self.ready()

    def test_new_failed_result_cannot_be_hidden_by_old_success(self):
        newer=copy.deepcopy(self.checks[1]);newer.update(id=100,conclusion='failure');self.checks.append(newer)
        with self.assertRaisesRegex(RuntimeError,'required check failed'):self.ready()

    def test_wrong_application_cannot_supply_required_context(self):
        self.checks[1]['app']['slug']='untrusted-app';self.assertFalse(self.ready())

    def test_different_commit_is_rejected(self):
        self.checks[1]['head_sha']='b'*40
        with self.assertRaisesRegex(RuntimeError,'commit mismatch'):self.ready()

    def test_spoofed_workflow_name_or_id_is_rejected(self):
        self.runs[2]['workflow_id']=999
        with self.assertRaisesRegex(RuntimeError,'workflow identity'):self.ready()

    def test_dispatch_or_nonmain_run_cannot_publish(self):
        for field,value in [('event','workflow_dispatch'),('head_branch','feature')]:
            with self.subTest(field=field):
                old=self.runs[2][field];self.runs[2][field]=value
                with self.assertRaisesRegex(RuntimeError,'exact main push'):self.ready()
                self.runs[2][field]=old

    def test_different_primary_ci_run_is_rejected(self):
        with self.assertRaisesRegex(RuntimeError,'different CI run'):gate.checks_ready(self.checks,SHA,123)

    def test_incomplete_guardrail_workflow_waits_after_job_success(self):
        self.runs[2]['status']='in_progress';self.runs[2]['conclusion']=None
        self.assertFalse(self.ready())

    def test_main_advance_stops_publication(self):
        with patch.object(gate,'api',return_value={'object':{'sha':'b'*40}}):
            with self.assertRaisesRegex(RuntimeError,'Main advanced'):gate.wait_checks(SHA,1)

    def test_polling_has_a_fixed_deadline(self):
        responses=[{'object':{'sha':SHA}},{'check_runs':[]}]*2
        with patch.object(gate,'api',side_effect=responses),patch.object(gate,'POLLS',2),patch.object(gate.time,'sleep') as sleep:
            with self.assertRaisesRegex(RuntimeError,'publication deadline'):gate.wait_checks(SHA,1)
            sleep.assert_called_once_with(gate.DELAY)

    def test_cancelled_publication_ci_is_rejected(self):
        self.runs[1].update(status='completed',conclusion='cancelled')
        with self.assertRaisesRegex(RuntimeError,'cancelled or unsuccessful'):self.ready()
