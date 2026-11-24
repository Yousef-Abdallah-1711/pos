#!/usr/bin/env node
import { mkdir, readFile } from "node:fs/promises";
import { createRequire } from "node:module";
import os from "node:os";
import path from "node:path";
import { requireLocalOrigin } from "../performance/local-origin.mjs";

const fixturePath = process.env.SERVIO_FINANCIAL_FIXTURE_PATH;
const apiOrigin = process.env.SERVIO_FINANCIAL_API_URL;
const dashboardOrigin = process.env.SERVIO_FINANCIAL_DASHBOARD_URL;
const repositoryRoot = path.resolve(process.cwd());
const fixtureFile = fixturePath ? path.resolve(fixturePath) : "";
const outputDirectory = path.resolve(
  process.env.SERVIO_FINANCIAL_SCREENSHOT_DIR ?? path.join(os.tmpdir(), "servio-financial-captures")
);

if (!fixturePath) throw new Error("Set SERVIO_FINANCIAL_FIXTURE_PATH to an out-of-repository synthetic fixture file.");
function assertOutsideRepository(candidate, label) {
  const relative = path.relative(repositoryRoot, candidate);
  if (relative === "" || (!path.isAbsolute(relative) && relative !== ".." && !relative.startsWith(`..${path.sep}`))) {
    throw new Error(`${label} must be outside the repository.`);
  }
}
assertOutsideRepository(fixtureFile, "SERVIO_FINANCIAL_FIXTURE_PATH");
assertOutsideRepository(outputDirectory, "SERVIO_FINANCIAL_SCREENSHOT_DIR");
const api = requireLocalOrigin(apiOrigin, "SERVIO_FINANCIAL_API_URL");
const dashboard = requireLocalOrigin(dashboardOrigin, "SERVIO_FINANCIAL_DASHBOARD_URL");
const fixture = JSON.parse(await readFile(fixtureFile, "utf8"));
if (!fixture?.accountant?.email || !fixture?.accountant?.password) {
  throw new Error("The synthetic fixture must contain accountant.email and accountant.password.");
}

const requireFromWorkingDirectory = createRequire(path.resolve(process.cwd(), "package.json"));
const { chromium } = requireFromWorkingDirectory("@playwright/test");
await mkdir(outputDirectory, { recursive: true });
const browser = await chromium.launch();
try {
  const context = await browser.newContext({ viewport: { width: 1440, height: 1000 }, deviceScaleFactor: 1 });
  const page = await context.newPage();
  const allowedOrigins = new Set([api.origin, dashboard.origin]);
  await page.route("**/*", (route) => {
    const requestUrl = new URL(route.request().url());
    if (!["http:", "https:"].includes(requestUrl.protocol) || allowedOrigins.has(requestUrl.origin)) {
      return route.continue();
    }
    return route.abort("blockedbyclient");
  });

  const login = await page.request.post(new URL("/api/v1/auth/login", api), {
    headers: { Origin: dashboard.origin },
    data: { email: fixture.accountant.email, password: fixture.accountant.password },
    maxRedirects: 0,
  });
  if (login.status() !== 200) throw new Error(`Synthetic accountant login failed with HTTP ${login.status()}.`);
  const session = await page.request.get(new URL("/api/v1/auth/session", api), { maxRedirects: 0 });
  if (session.status() !== 200) throw new Error(`Synthetic accountant session failed with HTTP ${session.status()}.`);

  const screenshots = [];
  for (const capture of [
    { locale: "en", viewport: { width: 1440, height: 1000 }, filename: "acceptance-en-balance-sheet.png" },
    { locale: "ar", viewport: { width: 390, height: 844 }, filename: "acceptance-ar-mobile-balance-sheet.png" },
  ]) {
    await page.setViewportSize(capture.viewport);
    await page.goto(new URL(`/${capture.locale}/dashboard/reports/financial/balance-sheet`, dashboard), {
      waitUntil: "domcontentloaded",
      timeout: 30_000,
    });
    await page.locator("main").waitFor({ state: "visible", timeout: 30_000 });
    await page.screenshot({ path: path.join(outputDirectory, capture.filename), fullPage: true });
    screenshots.push(capture.filename);
  }

  console.log(JSON.stringify({ login_status: login.status(), session_status: session.status(), screenshots }));
  await context.close();
} finally {
  await browser.close();
}
