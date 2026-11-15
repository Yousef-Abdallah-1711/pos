"""Static guardrails for the disposable CI browser runner and release triggers."""
import json
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[2]


class CIIsolationTests(unittest.TestCase):
    def test_generated_browser_stack_has_no_external_route(self):
        runner = (ROOT / 'ops/aws/ci-browser-stack.py').read_text()
        self.assertIn("config['networks']['edge']['internal'] = True", runner)
        self.assertIn("'SERVIO_ACCEPTANCE_ENV=ci'", runner)
        self.assertIn("SERVIO_ACCEPTANCE_API_BASE_URL=https://api.servio.test/api/v1", runner)
        self.assertIn('/work/acceptance-target.cjs:ro', runner)
        self.assertIn("replace('rakmyat.online', 'servio.test')", runner)
        self.assertIn("service.pop('ports', None)", runner)
        self.assertNotIn("'--network', 'host'", runner)

    def test_compose_targets_and_external_services_are_safe(self):
        compose = json.loads((ROOT / 'ops/aws/compose.json').read_text())
        self.assertTrue(compose['networks']['app']['internal'])
        self.assertEqual(compose['services']['backend']['environment']['PAYMOB_MODE'], 'mock')
        self.assertNotIn('PAYMOB_API_KEY', compose['services']['backend']['environment'])

    def test_mutating_acceptance_scripts_have_no_public_production_api_literals(self):
        scripts = [
            'ops/aws/browser-smoke.cjs',
            'ops/aws/browser-storefront-flow.cjs',
            'ops/aws/browser-upload.cjs',
            'ops/aws/browser-auth-validation.cjs',
            'ops/aws/journal-api-acceptance.cjs',
            'ops/aws/browser-tenant-lifecycle.cjs',
            'ops/aws/smoke-api.py',
            'ops/aws/smoke-refresh.py',
            'ops/aws/smoke-admin.py',
            'ops/aws/smoke-pos-kds.py',
            'ops/aws/check-completed-sale.py',
            'ops/aws/prepare-route-qa.py',
            'frontend/apps/landing/scripts/trial-registration.mjs',
        ]
        for name in scripts:
            with self.subTest(script=name):
                self.assertNotIn('rakmyat.online', (ROOT / name).read_text(encoding='utf-8'))

    def test_staging_publish_and_deploy_are_not_automatic(self):
        ci = (ROOT / '.github/workflows/production-ci.yml').read_text()
        release = (ROOT / '.github/workflows/production-release.yml').read_text()
        self.assertNotIn('publish-main-release:', ci)
        self.assertNotIn('publish-ci-images.py', ci)
        self.assertIn('workflow_dispatch:', release)
        self.assertNotIn('workflow_run:', release)


if __name__ == '__main__':
    unittest.main()
