import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './tests/phase2a-e2e',
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: 'list',
  timeout: 120_000,
  use: {
    baseURL: process.env.WS05_DASHBOARD_URL || 'http://localhost:13102',
    headless: true,
    // Chrome 142+ gates local-site requests to loopback APIs behind LNA permission.
    // The disposable browser suite exercises this local development topology.
    permissions: ['local-network-access'],
    launchOptions: {
      args: process.env.WS05_MAP_LOCALHOST ? [
        `--host-resolver-rules=MAP localhost ${process.env.WS05_MAP_LOCALHOST}, MAP phase2a-restaurant-a.localhost ${process.env.WS05_MAP_LOCALHOST}`,
      ] : [],
    },
  },
})
