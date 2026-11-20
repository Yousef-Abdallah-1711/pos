import os
import unittest
from unittest.mock import MagicMock, patch

from acceptance_target import (
    STAGING_API_RELEASE_URL,
    _verify_live_staging,
    api_base,
    require_acceptance_url,
    ui_origin,
)


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
            with self.assertRaisesRegex(RuntimeError, 'disposable CI or verified staging'):
                api_base()
        with patch.dict(os.environ, {'SERVIO_ACCEPTANCE_ENV': 'staging'}, clear=True):
            with self.assertRaisesRegex(RuntimeError, 'approved HTTPS host'):
                require_acceptance_url('https://api.servio.test/api/v1')

    def test_allows_only_explicit_staging_targets_after_verifying_live_environment(self):
        env = {
            'SERVIO_ACCEPTANCE_ENV': 'staging',
            'SERVIO_ACCEPTANCE_API_BASE_URL': 'https://api.pos.rakmyat.online/api/v1',
            'SERVIO_ACCEPTANCE_DASHBOARD_URL': 'https://pos.app.rakmyat.online',
        }
        with patch.dict(os.environ, env, clear=True), \
                patch('acceptance_target._verify_live_staging') as verify:
            self.assertEqual(api_base(), env['SERVIO_ACCEPTANCE_API_BASE_URL'])
            self.assertEqual(ui_origin('dashboard'), env['SERVIO_ACCEPTANCE_DASHBOARD_URL'])
            self.assertEqual(verify.call_count, 2)

    def test_live_staging_verification_requires_proxy_header_and_receipt(self):
        response = MagicMock()
        response.status = 200
        response.headers.get.return_value = 'staging'
        response.read.return_value = b'{"environment":"staging","status":"verified","sha":"0000000000000000000000000000000000000000"}'
        response.__enter__.return_value = response
        opener = MagicMock()
        opener.open.return_value = response
        with patch('acceptance_target.safe_opener', return_value=opener):
            _verify_live_staging()
        request = opener.open.call_args.args[0]
        self.assertEqual(request.full_url, STAGING_API_RELEASE_URL)

        response.headers.get.return_value = 'production'
        with patch('acceptance_target.safe_opener', return_value=opener):
            with self.assertRaisesRegex(RuntimeError, 'identify itself as staging'):
                _verify_live_staging()

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

    def test_staging_rejects_unknown_hosts_and_redirect_targets(self):
        unsafe = [
            'https://api.example.com/api/v1',
            'https://api.pos.rakmyat.online.evil.invalid/api/v1',
            'https://api.pos.rakmyat.online:444/api/v1',
            'https://api.pos.rakmyat.online/api/v1?next=https://example.com',
        ]
        with patch.dict(os.environ, {'SERVIO_ACCEPTANCE_ENV': 'staging'}, clear=True):
            for target in unsafe:
                with self.subTest(target=target):
                    with self.assertRaises(RuntimeError):
                        require_acceptance_url(target)


if __name__ == '__main__':
    unittest.main()
