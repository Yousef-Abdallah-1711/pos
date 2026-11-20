"""Local, network-free checks for staging acceptance's release boundary."""
import importlib.util
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

SPEC = importlib.util.spec_from_file_location(
    'servio_staging_acceptance', Path(__file__).with_name('staging-safe-acceptance.py'))
suite = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(suite)


class StagingReceiptTests(unittest.TestCase):
    def setUp(self):
        suite.EXPECTED_SHA = 'a' * 40
        suite.RECEIPT_STATUS = 'candidate'
        suite.CHECKS.clear()

    def test_accepts_only_the_expected_candidate_release_receipt(self):
        receipt = {'sha': suite.EXPECTED_SHA, 'status': 'candidate', 'environment': 'staging'}
        with patch.object(suite, 'request', return_value=(200, receipt, {})):
            suite.verify_receipt('candidate')
        self.assertEqual(suite.CHECKS[-1]['status'], 'PASS')

    def test_rejects_receipt_sha_status_or_environment_mismatch(self):
        variants = [
            {'sha': 'b' * 40, 'status': 'candidate', 'environment': 'staging'},
            {'sha': suite.EXPECTED_SHA, 'status': 'verified', 'environment': 'staging'},
            {'sha': suite.EXPECTED_SHA, 'status': 'candidate', 'environment': 'production'},
        ]
        for receipt in variants:
            with self.subTest(receipt=receipt), patch.object(suite, 'request', return_value=(200, receipt, {})):
                with self.assertRaisesRegex(RuntimeError, 'Release receipt'):
                    suite.verify_receipt('candidate')

    def test_rejects_fixtures_with_group_or_world_permissions(self):
        with tempfile.TemporaryDirectory() as directory:
            fixture = Path(directory) / 'browser-admin.json'
            fixture.write_text('{"email":"qa@example.invalid","password":"synthetic"}')
            fixture.chmod(0o644)
            with patch.object(suite, 'FIXTURE_DIR', Path(directory)):
                with self.assertRaisesRegex(RuntimeError, 'Protected synthetic'):
                    suite.read_fixture('browser-admin.json')


if __name__ == '__main__':
    unittest.main()
