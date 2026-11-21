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

    def test_marketing_console_errors_are_serialized_before_route_navigation(self):
        script = (ROOT / 'frontend/apps/landing/scripts/marketing-acceptance.mjs').read_text()
        create_page = script.split('async function createPage(browser, viewport) {', maxsplit=1)[1].split('\nasync function ', maxsplit=1)[0]
        self.assertLess(create_page.index('await page.addInitScript(installConsoleErrorSerialization)'), create_page.index('return page'))
        self.assertIn('const serializedErrors = parseSerializedConsoleError(text)', script)
        self.assertIn('response = await page.goto(`${base}/${locale}${route}`', script)

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

    def test_staging_deployer_uses_its_safe_suite_and_candidate_receipt(self):
        deployer = (ROOT / 'ops/aws/deploy-release.py').read_text()
        staging_acceptance = deployer.split('def acceptance(', maxsplit=1)[1].split('\ndef rollback_apps', maxsplit=1)[0]
        for forbidden in ['smoke-api.py', 'smoke-refresh.py', 'smoke-admin.py', 'smoke-pos-kds.py', 'journal-api-acceptance.cjs']:
            with self.subTest(helper=forbidden):
                self.assertNotIn(forbidden, staging_acceptance)
        self.assertIn('staging-safe-acceptance.py', staging_acceptance)
        self.assertIn('SERVIO_STAGING_RECEIPT_STATUS=candidate', staging_acceptance)
        self.assertIn('SERVIO_STAGING_FIXTURE_DIR=', staging_acceptance)

        deploy = deployer.split('def deploy(', maxsplit=1)[1].split("\nif __name__", maxsplit=1)[0]
        self.assertLess(deploy.index('publish_candidate_receipt(manifest)'), deploy.index('acceptance(destination, sha)'))
        self.assertIn('restore_public_receipt(previous_receipt)', deploy)

    def test_disposable_ci_guard_is_not_relaxed_for_staging(self):
        target = (ROOT / 'ops/aws/acceptance_target.py').read_text()
        suite = (ROOT / 'ops/aws/staging-safe-acceptance.py').read_text()
        self.assertIn("os.environ.get('SERVIO_ACCEPTANCE_ENV') != 'ci'", target)
        self.assertIn('Mutating acceptance is allowed only in disposable CI', target)
        self.assertNotIn('SERVIO_ACCEPTANCE_ENV=ci', suite)
        self.assertNotIn('SERVIO_DISPOSABLE_BROWSER=1', suite)
        self.assertLess(suite.index('verify_receipt(RECEIPT_STATUS)'), suite.index("read_fixture('browser-fixtures.json')"))
        self.assertNotIn('/accounting/journal-entries', suite)
        self.assertNotIn('/public/orders', suite)

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

    def test_release_image_publisher_survives_only_unrelated_skips(self):
        workflow = (ROOT / '.github/workflows/production-ci.yml').read_text()
        publisher = workflow.split('  publish-release-images:', maxsplit=1)[1].split('\n  required:', maxsplit=1)[0]
        self.assertIn('if: ${{ !cancelled()', publisher)
        self.assertIn("github.event_name == 'push'", publisher)
        self.assertIn("github.ref == 'refs/heads/main'", publisher)
        self.assertIn("needs.application-images.result == 'success'", publisher)
        self.assertIn("needs.required.result == 'success'", publisher)
        self.assertIn('needs: [application-images, required]', publisher)


if __name__ == '__main__':
    unittest.main()
