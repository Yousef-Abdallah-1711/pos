import unittest

from check_required_ci_gates import failed_required_jobs


class RequiredCIGateTests(unittest.TestCase):
    def test_all_required_jobs_must_succeed(self):
        results = {name: {"result": "success"} for name in (
            "backend",
            "frontend",
            "security-and-config",
            "application-images",
            "financial-postgres-acceptance",
        )}
        self.assertEqual(failed_required_jobs(results, require_browser=False), {})

    def test_financial_acceptance_failure_blocks_release(self):
        results = {name: {"result": "success"} for name in (
            "backend",
            "frontend",
            "security-and-config",
            "application-images",
        )}
        results["financial-postgres-acceptance"] = {"result": "failure"}
        self.assertEqual(
            failed_required_jobs(results, require_browser=False),
            {"financial-postgres-acceptance": "failure"},
        )

    def test_skipped_financial_acceptance_blocks_release(self):
        results = {name: {"result": "success"} for name in (
            "backend",
            "frontend",
            "security-and-config",
            "application-images",
        )}
        results["financial-postgres-acceptance"] = {"result": "skipped"}
        self.assertEqual(
            failed_required_jobs(results, require_browser=False),
            {"financial-postgres-acceptance": "skipped"},
        )

    def test_required_browser_skip_blocks_release(self):
        results = {name: {"result": "success"} for name in (
            "backend",
            "frontend",
            "security-and-config",
            "application-images",
            "financial-postgres-acceptance",
        )}
        results["auth-browser-acceptance"] = {"result": "skipped"}
        self.assertEqual(
            failed_required_jobs(results, require_browser=True),
            {"auth-browser-acceptance": "skipped"},
        )


if __name__ == "__main__":
    unittest.main()
