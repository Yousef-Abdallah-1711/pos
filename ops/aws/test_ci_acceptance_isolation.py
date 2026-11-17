"""Static guardrails for the disposable CI browser runner and release triggers."""
import json
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[2]


class CIIsolationTests(unittest.TestCase):
    def test_shared_postgres_backend_packages_run_serially(self):
        workflow = (ROOT / '.github/workflows/production-ci.yml').read_text()
        backend = workflow.split('  backend:', maxsplit=1)[1].split('\n  frontend:', maxsplit=1)[0]
        self.assertIn('GOFLAGS: -p=1', backend)

    def test_generated_browser_stack_has_no_external_route(self):
        runner = (ROOT / 'ops/aws/ci-browser-stack.py').read_text()
        self.assertIn("config['networks']['edge']['internal'] = True", runner)
        self.assertIn("'SERVIO_ACCEPTANCE_ENV=ci'", runner)
        self.assertIn("SERVIO_ACCEPTANCE_API_BASE_URL=https://api.servio.test/api/v1", runner)
        self.assertIn('/work/acceptance-target.cjs:ro', runner)
        self.assertIn("replace('rakmyat.online', 'servio.test')", runner)
        self.assertIn("replace('api.pos.servio.test', 'api.servio.test')", runner)
        self.assertIn("service.pop('ports', None)", runner)
        self.assertNotIn("'--network', 'host'", runner)

    def test_upload_acceptance_closes_browser_when_session_logout_fails(self):
        script = (ROOT / 'ops/aws/browser-upload.cjs').read_text()
        self.assertIn("timeout:10000", script)
        self.assertIn("await browser.close();", script)
        self.assertIn("record('Synthetic tenant session logout',logout.status(),204);", script)
        self.assertIn("catch(error){\n   checks.push({check:'Synthetic tenant session logout'", script)

    def test_compose_targets_and_external_services_are_safe(self):
        compose = json.loads((ROOT / 'ops/aws/compose.json').read_text())
        self.assertTrue(compose['networks']['app']['internal'])
        self.assertEqual(compose['services']['backend']['environment']['PAYMOB_MODE'], 'mock')
        self.assertNotIn('PAYMOB_API_KEY', compose['services']['backend']['environment'])

    def test_browser_smoke_measures_the_reported_375px_mobile_width(self):
        script = (ROOT / 'ops/aws/browser-smoke.cjs').read_text()
        self.assertIn('const responsiveViewports=[', script)
        for width in (320, 375, 390, 768, 1440):
            with self.subTest(width=width):
                self.assertIn(f'{{width:{width},height:', script)
        self.assertIn('await page.setViewportSize({width,height});', script)

    def test_marketing_console_errors_are_read_before_route_navigation(self):
        script = (ROOT / 'frontend/apps/landing/scripts/marketing-acceptance.mjs').read_text()
        load = script.split('async function load(page, locale, route) {', maxsplit=1)[1].split('\nasync function ', maxsplit=1)[0]
        self.assertLess(load.index('await flushPendingReactErrorDetails(page)'), load.index('await page.goto('))
        self.assertIn('await flushPendingReactErrorDetails(page)\n}', load)
        self.assertIn('pendingReactErrorDetails.set(page, pending)', script)

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

    def test_staging_release_requires_owner_dispatch_and_exact_main_ci(self):
        ci = (ROOT / '.github/workflows/production-ci.yml').read_text()
        release = (ROOT / '.github/workflows/production-release.yml').read_text()
        publisher = (ROOT / 'ops/aws/sign-publish-release.py').read_text()
        image_publisher = (ROOT / 'ops/aws/publish-ci-images.py').read_text()
        image_stager = (ROOT / 'ops/aws/stage-ci-images.py').read_text()
        self.assertIn("if: github.event_name == 'push' && github.ref == 'refs/heads/main'", ci)
        self.assertIn('ops/aws/stage-ci-images.py', ci)
        self.assertIn('ops/aws/publish-ci-images.py', ci)
        self.assertIn('name: release-bundle', ci)
        self.assertNotIn('gh release create', ci)
        unprivileged_jobs = ci.split('publish-release-images:', maxsplit=1)[0]
        self.assertNotIn('packages: write', unprivileged_jobs)
        self.assertIn('needs: [application-images, required]', ci)
        self.assertIn('wait_checks(sha, int(os.environ[\'GITHUB_RUN_ID\']))', image_publisher)
        self.assertIn("os.environ.get('GITHUB_EVENT_NAME') != 'push'", image_stager)
        self.assertNotIn('GHCR_TOKEN', image_stager)
        self.assertIn('workflow_dispatch:', release)
        self.assertNotIn('workflow_run:', release)
        self.assertIn("github.actor == github.repository_owner", release)
        self.assertIn('environment: staging', release)
        self.assertIn('PUBLISH-STAGING', release)
        self.assertIn('SERVIO_CI_RUN_ID', publisher)
        self.assertIn('SERVIO_RELEASE_CONFIRMATION', publisher)
        self.assertIn('wait_checks(sha,ci_run_id)', publisher)
        self.assertIn("ci_jobs=api('actions/runs/'+str(ci_run_id)+'/jobs?per_page=100')['jobs']", publisher)
        self.assertIn("job.get('name')=='publish-release-images'", publisher)
        self.assertIn("os.environ.get('GITHUB_EVENT_NAME')!='workflow_dispatch'", publisher)


if __name__ == '__main__':
    unittest.main()
