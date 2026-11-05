import type { Page } from '@playwright/test'

const candidateApiOrigin = new URL(
  process.env.WS05_API_URL || 'http://localhost:18081/api/v1',
).origin

export async function routeContainerApiToCandidate(page: Page) {
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
