#!/usr/bin/env python3
"""Extract measured aggregate metrics from a k6 --summary-export JSON file."""

from __future__ import annotations

import argparse
import json
import math
import sys
from pathlib import Path
from typing import Any


class SummaryError(ValueError):
    pass


def metric_values(summary: dict[str, Any], metric_name: str) -> dict[str, Any]:
    metrics = summary.get("metrics")
    if not isinstance(metrics, dict):
        raise SummaryError("summary must contain a metrics object")
    metric = metrics.get(metric_name)
    if not isinstance(metric, dict):
        raise SummaryError(f"summary is missing metric values for {metric_name}")
    values = metric.get("values")
    return values if isinstance(values, dict) else metric


def measured_number(values: dict[str, Any], key: str, metric_name: str) -> float:
    value = values.get(key)
    if isinstance(value, bool) or not isinstance(value, (int, float)) or not math.isfinite(value):
        raise SummaryError(f"summary is missing a finite {metric_name}.{key} measurement")
    return float(value)


def build_report(summary: dict[str, Any]) -> dict[str, Any]:
    request_values = metric_values(summary, "http_reqs")
    failure_values = metric_values(summary, "http_req_failed")
    duration_values = metric_values(summary, "http_req_duration")

    request_count = measured_number(request_values, "count", "http_reqs")
    request_rate = measured_number(request_values, "rate", "http_reqs")
    failure_key = "rate" if "rate" in failure_values else "value"
    failure_rate = measured_number(failure_values, failure_key, "http_req_failed")
    if request_count < 0 or request_rate < 0 or not 0 <= failure_rate <= 1:
        raise SummaryError("summary contains out-of-range request measurements")

    duration_keys = {
        "average": "avg",
        "median": "med",
        "maximum": "max",
        "p90": "p(90)",
        "p95": "p(95)",
        "p99": "p(99)",
    }
    duration = {
        label: measured_number(duration_values, key, "http_req_duration")
        for label, key in duration_keys.items()
    }
    if any(value < 0 for value in duration.values()):
        raise SummaryError("summary contains negative request durations")

    report: dict[str, Any] = {
        "schema": "servio.k6.measured-summary.v1",
        "requests": {"count": request_count, "per_second": request_rate},
        "failed_request_rate": failure_rate,
        "http_request_duration_ms": duration,
    }

    state = summary.get("state")
    if isinstance(state, dict) and "testRunDurationMs" in state:
        report["test_duration_ms"] = measured_number(state, "testRunDurationMs", "state")

    vus_values = metric_values(summary, "vus_max") if "vus_max" in summary["metrics"] else None
    if vus_values is not None:
        report["max_vus"] = measured_number(vus_values, "max", "vus_max")

    metrics = summary["metrics"]
    if "dropped_iterations" in metrics:
        dropped_values = metric_values(summary, "dropped_iterations")
        dropped = measured_number(dropped_values, "count", "dropped_iterations")
        if dropped < 0:
            raise SummaryError("summary contains a negative dropped-iteration count")
        report["dropped_iterations"] = dropped

    return report


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("summary", type=Path, help="k6 --summary-export JSON file")
    parser.add_argument("--output", type=Path, help="write the measured summary JSON to this file")
    args = parser.parse_args(argv)

    try:
        with args.summary.open(encoding="utf-8") as source:
            summary = json.load(source)
        if not isinstance(summary, dict):
            raise SummaryError("summary root must be a JSON object")
        report = build_report(summary)
        rendered = json.dumps(report, indent=2, sort_keys=True) + "\n"
        if args.output:
            args.output.write_text(rendered, encoding="utf-8")
        else:
            sys.stdout.write(rendered)
    except (OSError, json.JSONDecodeError, SummaryError) as error:
        parser.error(str(error))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
