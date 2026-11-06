import { expect, test } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const dashboardUrl = process.env.WS05_DASHBOARD_URL || 'http://localhost:13102'
const landingUrl = process.env.WS05_LANDING_URL || 'http://localhost:13101'
const landingRedirectUrl = process.env.WS05_LANDING_REDIRECT_URL || 'http://localhost:13101'

test('unauthenticated dashboard visit reaches the working landing login page', async ({ page }) => {
  await routeContainerApiToCandidate(page)
  await page.goto(`${dashboardUrl}/en/dashboard`, { waitUntil: 'domcontentloaded' })

  await expect(page).toHaveURL(`${landingRedirectUrl}/en/auth/login`, { timeout: 30_000 })
  await expect(page.getByRole('textbox', { name: /email/i })).toBeVisible()
  await expect(page.getByRole('textbox', { name: /password/i })).toBeVisible()
  await expect(page.getByRole('button', { name: /sign in/i })).toBeVisible()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})
