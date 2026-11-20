import re
import unittest
from pathlib import Path

from check_required_ci_gates import failed_required_jobs


class RequiredCIGateTests(unittest.TestCase):
    def test_required_job_checks_out_repository_before_running_gate_script(self):
        workflow = (Path(__file__).parents[2] / ".github/workflows/production-ci.yml").read_text(encoding="utf-8")
        required_job = re.search(r"(?ms)^  required:\n(.*?)(?=^  [A-Za-z0-9_-]+:\n|\Z)", workflow)
        self.assertIsNotNone(required_job)
        job_steps = required_job.group(1)
        checkout = job_steps.find("actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1")
        gate_script = job_steps.find("python3 ops/aws/check_required_ci_gates.py")
        self.assertGreaterEqual(checkout, 0, "required aggregator must check out the repository")
        self.assertGreater(gate_script, checkout, "checkout must precede the repository-local gate script")

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
