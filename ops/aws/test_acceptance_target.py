import os
import unittest
from unittest.mock import patch

from acceptance_target import api_base, require_acceptance_url, ui_origin


class AcceptanceTargetTests(unittest.TestCase):
    def test_allows_exact_disposable_api_and_ui_targets(self):
        env = {
            'SERVIO_ACCEPTANCE_ENV': 'ci',
            'SERVIO_ACCEPTANCE_API_BASE_URL': 'https://api.servio.test/api/v1',
            'SERVIO_ACCEPTANCE_LANDING_URL': 'https://pos.servio.test',
        }
        with patch.dict(os.environ, env, clear=False):
            self.assertEqual(api_base(), env['SERVIO_ACCEPTANCE_API_BASE_URL'])
            self.assertEqual(ui_origin('landing'), env['SERVIO_ACCEPTANCE_LANDING_URL'])

    def test_rejects_missing_target_and_non_ci_execution(self):
        with patch.dict(os.environ, {}, clear=True):
            with self.assertRaisesRegex(RuntimeError, 'disposable CI'):
                api_base()
        with patch.dict(os.environ, {'SERVIO_ACCEPTANCE_ENV': 'staging'}, clear=True):
            with self.assertRaisesRegex(RuntimeError, 'disposable CI'):
                require_acceptance_url('https://api.servio.test/api/v1')

    def test_rejects_public_unknown_and_redirect_targets(self):
        unsafe = [
            'https://api.pos.rakmyat.online/api/v1',
            'https://api.example.com/api/v1',
            'http://api.servio.test/api/v1',
            'https://user:pass@api.servio.test/api/v1',
            'https://api.servio.test:444/api/v1',
            'https://api.servio.test/api/v1?redirect=https://example.com',
            'https://api.servio.test/api/v1#fragment',
        ]
        with patch.dict(os.environ, {'SERVIO_ACCEPTANCE_ENV': 'ci'}, clear=True):
            for target in unsafe:
                with self.subTest(target=target):
                    with self.assertRaises(RuntimeError):
                        require_acceptance_url(target)


if __name__ == '__main__':
    unittest.main()
