#!/usr/bin/env node
import { readFile, writeFile } from "node:fs/promises";

const HELP = [
  "Usage: node scripts/performance/measure-api.mjs --routes <safe-get-routes.json>",
  "",
  "Environment:",
  "  P2C_API_BASE_URL   API base URL",
  "  P2C_BEARER_TOKEN   Optional synthetic token; never printed",
  "  P2C_TENANT_ID      Optional tenant header",
  "  P2C_RESTAURANT_ID  Optional restaurant header",
  "  P2C_ROLE           Non-secret role label for evidence (default synthetic-admin)",
  "  P2C_SOURCE_COMMIT  Source commit label for evidence",
  "  P2C_SOURCE_FINGERPRINT  Source change-set fingerprint for evidence",
  "  P2C_API_SAMPLES    Warm samples per route (default 25, range 3-100)",
  "  P2C_API_OUTPUT     Optional CSV output path (otherwise stdout)",
  "",
  "Route JSON: an array of { name, path }; only same-origin GET requests are sent.",
  "Response bodies and headers are never written."
].join("\n");

function percentile(values, quantile) {
  const sorted = [...values].sort((left, right) => left - right);
  const index = Math.max(0, Math.ceil(quantile * sorted.length) - 1);
  return sorted[index];
}

function countRecords(value) {
  if (!value || typeof value !== "object") return "";
  for (const key of ["orders", "products", "items", "tenants", "tables", "results"]) {
    if (Array.isArray(value[key])) return value[key].length;
  }
  if (Array.isArray(value.data)) return value.data.length;
  if (value.data && typeof value.data === "object") return countRecords(value.data);
  return "";
}

function escapeCsv(value) {
  const text = String(value ?? "");
  return /[",\r\n]/.test(text) ? '"' + text.replaceAll('"', '""') + '"' : text;
}

async function main() {
  const args = process.argv.slice(2);
  if (args.includes("--help") || args.includes("-h")) {
    process.stdout.write(HELP + "\n");
    return;
  }
  const routeIndex = args.indexOf("--routes");
  const routeFile = routeIndex >= 0 ? args[routeIndex + 1] : "";
  const baseValue = process.env.P2C_API_BASE_URL;
  if (!routeFile || !baseValue) throw new Error("Provide --routes and P2C_API_BASE_URL.");
  const base = new URL(baseValue);
  const tenantId = process.env.P2C_TENANT_ID ?? "";
  const restaurantId = process.env.P2C_RESTAURANT_ID ?? "";
  const role = process.env.P2C_ROLE ?? "synthetic-admin";
  const sourceCommit = process.env.P2C_SOURCE_COMMIT ?? "";
  const sourceFingerprint = process.env.P2C_SOURCE_FINGERPRINT ?? "";
  const routes = JSON.parse(await readFile(routeFile, "utf8"));
  if (!Array.isArray(routes) || routes.length === 0) throw new Error("Route file must be a non-empty JSON array.");

  const samples = Number(process.env.P2C_API_SAMPLES ?? 25);
  if (!Number.isInteger(samples) || samples < 3 || samples > 100) {
    throw new Error("P2C_API_SAMPLES must be an integer between 3 and 100.");
  }

  const headers = new Headers({ Accept: "application/json" });
  if (process.env.P2C_BEARER_TOKEN) headers.set("Authorization", "Bearer " + process.env.P2C_BEARER_TOKEN);
  if (process.env.P2C_TENANT_ID) headers.set("X-Tenant-ID", process.env.P2C_TENANT_ID);
  if (process.env.P2C_RESTAURANT_ID) headers.set("X-Restaurant-ID", process.env.P2C_RESTAURANT_ID);

  const rows = [];
  for (const route of routes) {
    if (!route || typeof route.path !== "string" || !route.path.startsWith("/") || route.path.startsWith("//")) {
      throw new Error("Every route must use a same-origin absolute path.");
    }
    const url = new URL(route.path, base);
    if (url.origin !== base.origin) throw new Error("Cross-origin routes are not allowed.");
    const latencies = [];
    const statuses = new Map();
    let responseBytes = "";
    let records = "";
    let firstSampleMs = "";
    let errors = 0;

    for (let sample = 0; sample < samples + 1; sample += 1) {
      const startedAt = performance.now();
      try {
        const response = await fetch(url, {
          method: "GET",
          headers,
          signal: AbortSignal.timeout(15000),
        });
        const body = await response.text();
        const elapsedMs = performance.now() - startedAt;
        statuses.set(response.status, (statuses.get(response.status) ?? 0) + 1);
        if (sample === 0) {
          firstSampleMs = elapsedMs.toFixed(2);
        } else {
          latencies.push(elapsedMs);
          responseBytes = Buffer.byteLength(body);
          try {
            records = countRecords(JSON.parse(body));
          } catch {
            records = "";
          }
        }
      } catch {
        errors += 1;
        if (sample > 0) latencies.push(performance.now() - startedAt);
      }
    }

    rows.push({
      method: "GET",
      route: route.name || url.pathname,
      role,
      tenant_id: tenantId,
      restaurant_id: restaurantId,
      source_commit: sourceCommit,
      source_fingerprint: sourceFingerprint,
      samples: latencies.length,
      status_counts: [...statuses.entries()].map(([status, count]) => status + ":" + count).join("|"),
      errors,
      response_bytes: responseBytes,
      records,
      first_sample_ms: firstSampleMs,
      warm_p50_ms: percentile(latencies, 0.5)?.toFixed(2) ?? "",
      warm_p75_ms: percentile(latencies, 0.75)?.toFixed(2) ?? "",
      warm_p95_ms: percentile(latencies, 0.95)?.toFixed(2) ?? "",
      warm_p99_ms: percentile(latencies, 0.99)?.toFixed(2) ?? "",
    });
  }

  const columns = Object.keys(rows[0]);
  const csv = [columns.join(","), ...rows.map((row) => columns.map((column) => escapeCsv(row[column])).join(","))].join("\n") + "\n";
  if (process.env.P2C_API_OUTPUT) {
    await writeFile(process.env.P2C_API_OUTPUT, csv, "utf8");
  } else {
    process.stdout.write(csv);
  }
}

main().catch((error) => {
  process.stderr.write(error.message + "\n");
  process.exitCode = 1;
});

export { countRecords, percentile };
