import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/financial-browser',
  testMatch: '**/*.browser.js',
  fullyParallel: false,
  workers: 1,
  retries: 0,
  timeout: 90_000,
  expect: { timeout: 15_000 },
  outputDir: './test-results/financial-browser',
  reporter: [
    ['list'],
    ['junit', { outputFile: './test-results/financial-browser-junit.xml' }],
    ['html', { outputFolder: './playwright-report/financial-browser', open: 'never' }],
  ],
  use: {
    ...devices['Desktop Chrome'],
    baseURL: process.env.SERVIO_FINANCIAL_DASHBOARD_URL,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
})
