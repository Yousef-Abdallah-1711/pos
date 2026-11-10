import { expect, test } from '@playwright/test'
import { apiURL, dashboardURL, fillRegistration, landingURL, openRegistration, submitRegistration, testAccount, waitForAvailable } from './fixtures'

test.describe.configure({ mode: 'serial' })

test('registers a tenant owner, persists an HttpOnly session, logs out, denies access, and logs in again', async ({ page }) => {
  await openRegistration(page)
  await fillRegistration(page)
  await waitForAvailable(page)
  await submitRegistration(page)

  await expect(page).toHaveURL(`${dashboardURL}/en/dashboard`, { timeout: 60_000 })
  const cookies = await page.context().cookies(apiURL)
  const session = cookies.find((cookie) => cookie.name === 'poos_session')
  expect(session?.httpOnly).toBe(true)
  expect(session?.secure).toBe(false)
  const browserAuth = await page.evaluate(() => {
    const state = localStorage.getItem('auth-storage')
    return state ? JSON.parse(state) as { state?: { token?: string } } : null
  })
  expect(browserAuth?.state?.token ?? '').toBe('')
  expect(new URL(page.url()).search).toBe('')

  const sessionResponse = await page.request.get(`${apiURL}/auth/session`)
  expect(sessionResponse.ok()).toBe(true)
  const sessionUser = await sessionResponse.json() as { email: string; tenant_id: number; restaurant_id: number; role: string }
  expect(sessionUser.email).toBe(testAccount.email)
  expect(sessionUser.tenant_id).toBeGreaterThan(0)
  expect(sessionUser.restaurant_id).toBeGreaterThan(0)
  expect(sessionUser.role).toBe('owner')
  if (process.env.REG_E2E_FIXTURE_FILE) {
    const { writeFileSync } = await import('node:fs')
    writeFileSync(process.env.REG_E2E_FIXTURE_FILE, JSON.stringify({ email: testAccount.email, slug: testAccount.slug }), 'utf8')
  }

  const logout = await page.request.post(`${apiURL}/auth/logout`, { headers: { Origin: landingURL } })
  expect(logout.status()).toBe(204)
  await expect.poll(async () => (await page.context().cookies(apiURL)).some((cookie) => cookie.name === 'poos_session')).toBe(false)
  await page.goto(`${dashboardURL}/en/dashboard`, { waitUntil: 'domcontentloaded' })
  await expect(page).toHaveURL(/\/en\/auth\/login/, { timeout: 30_000 })

  await page.goto(`${landingURL}/en/auth/login`)
  await page.locator('#email').fill(testAccount.email)
  await page.locator('#password').fill(testAccount.password)
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(`${dashboardURL}/en/dashboard`, { timeout: 60_000 })
  expect((await page.context().cookies(apiURL)).find((cookie) => cookie.name === 'poos_session')?.httpOnly).toBe(true)
})

test('rejects duplicate email with a field error and preserves the form values', async ({ page }) => {
  await openRegistration(page)
  await fillRegistration(page, { email: testAccount.email, slug: `${testAccount.slug}-email-retry`, restaurantName: `${testAccount.restaurantName} Retry` })
  await waitForAvailable(page)
  await submitRegistration(page)
  await expect(page.locator('#email')).toHaveAttribute('aria-invalid', 'true')
  await expect(page.locator('#email-error')).toContainText(/already exists|already in use/i)
  await expect(page.locator('#name')).toHaveValue(testAccount.name)
})

test('rejects a duplicate subdomain from the server after an advisory available result', async ({ page }) => {
  await page.route('**/api/v1/auth/check-subdomain?**', (route) => route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify({ available: true, slug: testAccount.slug }),
  }))
  await openRegistration(page)
  await fillRegistration(page, { email: `registration-slug-conflict-${Date.now()}@example.test`, slug: testAccount.slug, restaurantName: `Duplicate Slug ${Date.now()}` })
  await waitForAvailable(page)
  await submitRegistration(page)
  await expect(page.locator('#restaurantSlug')).toHaveAttribute('aria-invalid', 'true')
  await expect(page.locator('#restaurantSlug-error')).toContainText(/already|use/i)
})

test('rejects an overlong phone before sending registration', async ({ page }) => {
  let registrationRequests = 0
  await page.on('request', (request) => {
    if (request.url().includes('/api/v1/auth/register')) registrationRequests += 1
  })
  await openRegistration(page)
  await fillRegistration(page, { email: `registration-long-phone-${Date.now()}@example.test`, slug: `long-phone-${Date.now()}`, restaurantName: `Long Phone ${Date.now()}`, phone: '+1234567890123456' })
  await submitRegistration(page)
  await expect(page.locator('#phone')).toHaveAttribute('aria-invalid', 'true')
  await expect(page.locator('#phone-error')).toBeVisible()
  expect(registrationRequests).toBe(0)
})
