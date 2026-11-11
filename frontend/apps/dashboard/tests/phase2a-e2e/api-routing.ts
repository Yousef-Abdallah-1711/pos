import type { Page } from '@playwright/test'

const apiHostPort = process.env.WS05_API_HOST_PORT || '18082'
const containerApiOrigin = new URL(
  process.env.WS05_API_INTERNAL_URL || 'http://localhost:18080/api/v1',
).origin

export async function routeContainerApiToCandidate(page: Page) {
  const loopbackServicePorts = new Map([
    ['13111', '3001'],
    ['13112', '3002'],
    ['13113', '3003'],
    ['13114', '3004'],
  ])
  await page.route(/^http:\/\/localhost:(13111|13112|13113|13114)\//, async (route) => {
    if (process.env.WS05_TEST_RUNNER_ON_HOST) return route.continue()
    const requestedURL = new URL(route.request().url())
    const servicePort = loopbackServicePorts.get(requestedURL.port)
    if (!servicePort) return route.continue()
    // Dockerized runners reach container ports on the compose network. A host runner
    // instead reaches the same services through their published localhost ports.
    requestedURL.port = servicePort
    if (process.env.WS05_MAP_LOCALHOST) requestedURL.hostname = process.env.WS05_MAP_LOCALHOST
    const response = await route.fetch({ url: requestedURL.toString(), timeout: 30_000 })
    await route.fulfill({ response })
  })

  if (process.env.WS05_TEST_RUNNER_ON_HOST) return

  // A browser inside the compose network uses the API container port. Host-run
  // browsers reach the same API directly through the published loopback port.
  await page.route(`http://localhost:${apiHostPort}/**`, async (route) => {
    const requestedURL = new URL(route.request().url())
    const response = await route.fetch({
      url: `${containerApiOrigin}${requestedURL.pathname}${requestedURL.search}`,
      timeout: 30_000,
    })
    await route.fulfill({ response })
  })
}
