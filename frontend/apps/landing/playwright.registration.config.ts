import { defineConfig, devices } from '@playwright/test'

const landingURL = process.env.REG_E2E_LANDING_URL ?? 'http://localhost:14211'

export default defineConfig({
  testDir: './tests/registration-e2e',
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: 'list',
  outputDir: process.env.REG_E2E_OUTPUT_DIR ?? './test-results/registration-e2e',
  timeout: 45_000,
  expect: { timeout: 10_000 },
  use: {
    ...devices['Desktop Chrome'],
    baseURL: landingURL,
    browserName: 'chromium',
    headless: true,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    viewport: { width: 1440, height: 1000 },
  },
})
