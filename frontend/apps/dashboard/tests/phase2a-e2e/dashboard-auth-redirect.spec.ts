import { expect, test } from '@playwright/test'

const dashboardUrl = process.env.WS05_DASHBOARD_URL || 'http://localhost:13102'
const landingUrl = process.env.WS05_LANDING_URL || 'http://localhost:13101'
const landingServiceUrl = process.env.WS05_LANDING_SERVICE_URL || 'http://localhost:3001'

test('unauthenticated dashboard visit reaches the working landing login page', async ({ page }) => {
  // The browser runner shares the API container network namespace, so localhost
  // host-published ports are not reachable from it. Keep the real redirect URL
  // while routing its requests to the actual landing service on the Compose network.
  let requestedPublishedLandingURL = false
  await page.route(`${landingUrl}/**`, async (route) => {
    const redirectedURL = new URL(route.request().url())
    requestedPublishedLandingURL = true
    const response = await route.fetch({
      url: `${landingServiceUrl}${redirectedURL.pathname}${redirectedURL.search}`,
    })
    await route.fulfill({ response })
  })

  await page.goto(`${dashboardUrl}/en/dashboard`, { waitUntil: 'domcontentloaded' })

  await expect(page).toHaveURL(`${landingUrl}/en/auth/login`, { timeout: 30_000 })
  expect(requestedPublishedLandingURL).toBe(true)
  await expect(page.getByRole('textbox', { name: /email/i })).toBeVisible()
  await expect(page.getByRole('textbox', { name: /password/i })).toBeVisible()
  await expect(page.getByRole('button', { name: /sign in/i })).toBeVisible()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})
