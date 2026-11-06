import type { Page } from '@playwright/test'

const candidateApiOrigin = new URL(
  process.env.WS05_API_URL || 'http://localhost:18081/api/v1',
).origin

export async function routeContainerApiToCandidate(page: Page) {
  const loopbackServicePorts = new Map([
    ['13101', '3001'],
    ['13102', '3002'],
    ['13103', '3003'],
    ['13104', '3004'],
  ])
  await page.route(/^http:\/\/localhost:(13101|13102|13103|13104)\//, async (route) => {
    const requestedURL = new URL(route.request().url())
    const servicePort = loopbackServicePorts.get(requestedURL.port)
    if (!servicePort) return route.continue()
    requestedURL.port = servicePort
    const response = await route.fetch({ url: requestedURL.toString() })
    await route.fulfill({ response })
  })
  await page.route('http://localhost:18080/**', async (route) => {
    const requestedUrl = new URL(route.request().url())
    const response = await route.fetch({
      url: `${candidateApiOrigin}${requestedUrl.pathname}${requestedUrl.search}`,
    })
    if (requestedUrl.pathname.endsWith('/auth/session')) {
      console.log(`E2E_SESSION_FORWARD ${response.status()}`)
    }
    await route.fulfill({ response })
  })
}
