import sys
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2] / "scripts"))

from summarize_k6 import SummaryError, build_report  # noqa: E402


class SummarizeK6Tests(unittest.TestCase):
    def setUp(self):
        self.summary = {
            "state": {"testRunDurationMs": 60000},
            "metrics": {
                "http_reqs": {"count": 120, "rate": 2},
                "http_req_failed": {"value": 0.025},
                "http_req_duration": {
                    "avg": 84,
                    "med": 70,
                    "max": 240,
                    "p(90)": 140,
                    "p(95)": 170,
                    "p(99)": 220,
                },
                "vus_max": {"max": 10},
            },
        }

    def test_reports_only_measured_aggregate_metrics(self):
        self.assertEqual(
            build_report(self.summary),
            {
                "schema": "servio.k6.measured-summary.v1",
                "requests": {"count": 120.0, "per_second": 2.0},
                "failed_request_rate": 0.025,
                "http_request_duration_ms": {
                    "average": 84.0,
                    "median": 70.0,
                    "maximum": 240.0,
                    "p90": 140.0,
                    "p95": 170.0,
                    "p99": 220.0,
                },
                "test_duration_ms": 60000.0,
                "max_vus": 10.0,
            },
        )

    def test_missing_measurement_fails_instead_of_inventing_a_value(self):
        del self.summary["metrics"]["http_req_duration"]["p(99)"]
        with self.assertRaisesRegex(SummaryError, r"p\(99\)"):
            build_report(self.summary)

    def test_out_of_range_error_rate_is_rejected(self):
        self.summary["metrics"]["http_req_failed"]["value"] = 1.1
        with self.assertRaisesRegex(SummaryError, "out-of-range"):
            build_report(self.summary)

    def test_includes_measured_dropped_iterations_when_k6_reports_them(self):
        self.summary["metrics"]["dropped_iterations"] = {"count": 3}
        self.assertEqual(build_report(self.summary)["dropped_iterations"], 3.0)


if __name__ == "__main__":
    unittest.main()
