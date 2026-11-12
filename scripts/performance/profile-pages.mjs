#!/usr/bin/env node
import path from "node:path";
import { createHash } from "node:crypto";
import { createRequire } from "node:module";
import { readFile, writeFile } from "node:fs/promises";

const HELP = [
  "Usage: node scripts/performance/profile-pages.mjs --routes <local-routes.json>",
  "",
  "Environment:",
  "  P2C_BASE_URL              Local app origin",
  "  P2C_STORAGE_STATE         Optional synthetic Playwright state path",
  "  P2C_PAGE_SAMPLES          Hard navigations per route (default 5, range 1-30)",
  "  P2C_PAGE_OUTPUT           Optional CSV output path (otherwise stdout)",
  "  P2C_WATERFALL_OUTPUT      Optional sanitized API waterfall CSV path",
  "  P2C_API_ORIGIN            API origin to include in the waterfall",
  "  P2C_ROLE / P2C_TENANT_ID / P2C_RESTAURANT_ID  Evidence labels",
  "  P2C_SOURCE_COMMIT / P2C_SOURCE_FINGERPRINT    Evidence labels",
  "  P2C_READY_TIMEOUT_MS      Visible-content selector timeout (default 15000)",
  "",
  "Route JSON: an array of { path, readySelector }.",
  "The script stays on the configured origin and never saves HAR, headers, cookies, or response bodies."
].join("\n");

function escapeCsv(value) {
  const text = String(value ?? "");
  return /[",\r\n]/.test(text) ? '"' + text.replaceAll('"', '""') + '"' : text;
}

function safePath(value) {
  return new URL(value).pathname;
}

async function main() {
  const args = process.argv.slice(2);
  if (args.includes("--help") || args.includes("-h")) {
    process.stdout.write(HELP + "\n");
    return;
  }
  const routeIndex = args.indexOf("--routes");
  const routeFile = routeIndex >= 0 ? args[routeIndex + 1] : "";
  if (!routeFile || !process.env.P2C_BASE_URL) throw new Error("Provide --routes and P2C_BASE_URL.");
  const base = new URL(process.env.P2C_BASE_URL);
  const apiOrigin = new URL(process.env.P2C_API_ORIGIN ?? base.origin).origin;
  const routes = JSON.parse(await readFile(routeFile, "utf8"));
  const samples = Number(process.env.P2C_PAGE_SAMPLES ?? 5);
  const readyTimeoutMs = Number(process.env.P2C_READY_TIMEOUT_MS ?? 15000);
  if (!Array.isArray(routes) || routes.length === 0) throw new Error("Route file must contain a non-empty JSON array.");
  if (!Number.isInteger(samples) || samples < 1 || samples > 30) throw new Error("P2C_PAGE_SAMPLES must be between 1 and 30.");

  const dashboardRequire = createRequire(path.resolve("frontend/apps/dashboard/package.json"));
  const { chromium } = dashboardRequire("@playwright/test");
  const browser = await chromium.launch({ headless: true });
  const rows = [];
  const waterfallRows = [];

  try {
    for (const route of routes) {
      if (!route || typeof route.path !== "string" || !route.path.startsWith("/") || route.path.startsWith("//")) {
        throw new Error("Every route must use a same-origin absolute path.");
      }
      const target = new URL(route.path, base);
      if (target.origin !== base.origin) throw new Error("Cross-origin routes are not allowed.");
      const contextOptions = process.env.P2C_STORAGE_STATE
        ? { storageState: process.env.P2C_STORAGE_STATE }
        : {};
      const context = await browser.newContext(contextOptions);
      const page = await context.newPage();
      let consoleErrors = 0;
      let javascriptErrors = 0;
      let failedRequests = 0;
      let activeSample = 0;
      let sampleStartedAt = 0;
      let sampleApiRequests = [];
      const pendingApiRequests = new Map();
      page.on("console", (message) => {
        if (message.type() === "error") consoleErrors += 1;
      });
      page.on("pageerror", () => {
        javascriptErrors += 1;
      });
      page.on("requestfailed", () => {
        failedRequests += 1;
      });
      page.on("request", (request) => {
        const requestUrl = new URL(request.url());
        if (requestUrl.origin !== apiOrigin) return;
        const requestPath = requestUrl.pathname
          .replace(/\/[0-9]+(?=\/|$)/g, "/:id")
          .replace(/\/[0-9a-f-]{24,}(?=\/|$)/gi, "/:id");
        const queryFingerprint = requestUrl.search
          ? createHash("sha256").update(requestUrl.search).digest("hex").slice(0, 12)
          : "";
        const record = {
          route: safePath(target.href),
          sample: activeSample,
          request_path: requestPath,
          query_fingerprint: queryFingerprint,
          method: request.method(),
          start_offset_ms: Math.max(0, performance.now() - sampleStartedAt),
          status: "",
          duration_ms: "",
          content_length_bytes: "",
          failure: "",
        };
        sampleApiRequests.push(record);
        pendingApiRequests.set(request, record);
      });
      page.on("response", (response) => {
        const record = pendingApiRequests.get(response.request());
        if (!record) return;
        record.status = response.status();
        record.duration_ms = Math.max(0, performance.now() - sampleStartedAt - record.start_offset_ms).toFixed(2);
        record.content_length_bytes = response.headers()["content-length"] ?? "";
        pendingApiRequests.delete(response.request());
      });
      page.on("requestfailed", (request) => {
        const record = pendingApiRequests.get(request);
        if (!record) return;
        record.failure = request.failure()?.errorText ?? "request failed";
        record.duration_ms = Math.max(0, performance.now() - sampleStartedAt - record.start_offset_ms).toFixed(2);
        pendingApiRequests.delete(request);
      });
      await page.addInitScript(() => {
        const metrics = { lcp: 0, cls: 0, longTasks: 0 };
        Object.defineProperty(window, "__p2cMetrics", { value: metrics });
        try {
          new PerformanceObserver((entries) => {
            for (const entry of entries.getEntries()) metrics.lcp = entry.startTime;
          }).observe({ type: "largest-contentful-paint", buffered: true });
          new PerformanceObserver((entries) => {
            for (const entry of entries.getEntries()) {
              if (!entry.hadRecentInput) metrics.cls += entry.value;
            }
          }).observe({ type: "layout-shift", buffered: true });
          new PerformanceObserver((entries) => {
            metrics.longTasks += entries.getEntries().length;
          }).observe({ type: "longtask", buffered: true });
        } catch {
          metrics.observersUnavailable = true;
        }
      });

      for (let sample = 1; sample <= samples; sample += 1) {
        consoleErrors = 0;
        javascriptErrors = 0;
        failedRequests = 0;
        const startedAt = performance.now();
        sampleStartedAt = startedAt;
        activeSample = sample;
        sampleApiRequests = [];
        let status = "";
        let readyMs = "";
        let readinessStatus = route.readySelector ? "pending" : "not-configured";
        try {
          const response = sample === 1
            ? await page.goto(target.href, { waitUntil: "domcontentloaded", timeout: readyTimeoutMs })
            : await page.reload({ waitUntil: "domcontentloaded", timeout: readyTimeoutMs });
          status = response?.status() ?? "";
          if (route.readySelector) {
            await page.locator(route.readySelector).first().waitFor({
              state: "visible",
              timeout: readyTimeoutMs,
            });
            readinessStatus = "ready";
          }
          readyMs = (performance.now() - startedAt).toFixed(2);
        } catch {
          readinessStatus = status ? "selector-timeout" : "navigation-error";
          status = status || "TIMEOUT";
        }

        // Capture post-hydration API activity without treating network idleness as readiness.
        await page.waitForTimeout(1000);
        const metrics = await page.evaluate(() => {
          const nav = performance.getEntriesByType("navigation")[0];
          const resources = performance.getEntriesByType("resource");
          const p2c = window.__p2cMetrics || { lcp: 0, cls: 0, longTasks: 0 };
          return {
            ttfb: nav ? nav.responseStart : "",
            domContentLoaded: nav ? nav.domContentLoadedEventEnd : "",
            load: nav ? nav.loadEventEnd : "",
            fcp: performance.getEntriesByName("first-contentful-paint")[0]?.startTime ?? "",
            lcp: p2c.lcp,
            cls: p2c.cls,
            longTasks: p2c.longTasks,
            requestCount: resources.length,
            transferredBytes: resources.reduce((sum, resource) => sum + (resource.transferSize || 0), 0),
          };
        });
        const apiCounts = new Map();
        for (const apiRequest of sampleApiRequests) {
          const key = `${apiRequest.method}:${apiRequest.request_path}:${apiRequest.query_fingerprint}`;
          apiCounts.set(key, (apiCounts.get(key) ?? 0) + 1);
          waterfallRows.push(apiRequest);
        }
        const duplicateApiRequests = [...apiCounts.values()].reduce(
          (duplicates, count) => duplicates + Math.max(0, count - 1),
          0
        );
        rows.push({
          route: safePath(target.href),
          sample,
          cache_sample: sample === 1 ? "first-in-context" : "repeat-in-context",
          status,
          readiness_status: readinessStatus,
          ttfb_ms: metrics.ttfb,
          dom_content_loaded_ms: metrics.domContentLoaded,
          load_ms: metrics.load,
          fcp_ms: metrics.fcp,
          lcp_ms: metrics.lcp,
          cls: metrics.cls,
          meaningful_content_ms: readyMs,
          request_count: metrics.requestCount,
          api_request_count: sampleApiRequests.length,
          duplicate_api_requests: duplicateApiRequests,
          transferred_bytes: metrics.transferredBytes,
          long_tasks: metrics.longTasks,
          console_errors: consoleErrors,
          javascript_errors: javascriptErrors,
          failed_requests: failedRequests,
          final_path: safePath(page.url()),
          role: process.env.P2C_ROLE ?? "synthetic-admin",
          tenant_id: process.env.P2C_TENANT_ID ?? "",
          restaurant_id: process.env.P2C_RESTAURANT_ID ?? "",
          source_commit: process.env.P2C_SOURCE_COMMIT ?? "",
          source_fingerprint: process.env.P2C_SOURCE_FINGERPRINT ?? "",
        });
      }
      await context.close();
    }
  } finally {
    await browser.close();
  }

  const columns = Object.keys(rows[0]);
  const output = [columns.join(","), ...rows.map((row) => columns.map((column) => escapeCsv(row[column])).join(","))].join("\n") + "\n";
  if (process.env.P2C_PAGE_OUTPUT) {
    await writeFile(process.env.P2C_PAGE_OUTPUT, output, "utf8");
  } else {
    process.stdout.write(output);
  }
  if (process.env.P2C_WATERFALL_OUTPUT) {
    const waterfallColumns = [
      "route", "sample", "request_path", "query_fingerprint", "method", "status", "start_offset_ms",
      "duration_ms", "content_length_bytes", "failure",
    ];
    const waterfall = [
      waterfallColumns.join(","),
      ...waterfallRows.map((row) => waterfallColumns.map((column) => escapeCsv(row[column])).join(",")),
    ].join("\n") + "\n";
    await writeFile(process.env.P2C_WATERFALL_OUTPUT, waterfall, "utf8");
  }
}

main().catch((error) => {
  process.stderr.write(error.message + "\n");
  process.exitCode = 1;
});
