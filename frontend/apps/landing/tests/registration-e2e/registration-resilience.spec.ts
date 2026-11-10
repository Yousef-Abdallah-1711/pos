import { expect, test } from '@playwright/test'
import { fillRegistration, fulfillAvailableSlug, openRegistration, submitRegistration, waitForAvailable } from './fixtures'

test('shows a safe network error when registration cannot reach the API', async ({ page }) => {
  await fulfillAvailableSlug(page)
  await page.route('**/api/v1/auth/register', (route) => route.abort('failed'))
  await openRegistration(page)
  await fillRegistration(page, { email: `network-${Date.now()}@example.test`, slug: `network-${Date.now()}` })
  await waitForAvailable(page)
  await submitRegistration(page)
  await expect(page.locator('form [role="alert"]').filter({ hasText: /connect|service|try again/i })).toBeVisible()
})

test('maps a structured server field error to the email input and focuses it', async ({ page }) => {
  await fulfillAvailableSlug(page)
  await page.route('**/api/v1/auth/register', (route) => route.fulfill({
    status: 422,
    contentType: 'application/json',
    body: JSON.stringify({ code: 'VALIDATION_ERROR', error: 'Review fields', message: 'Review fields', fields: { email: 'This email is not allowed.' }, requestId: 'e2e-request-1' }),
  }))
  await openRegistration(page)
  await fillRegistration(page, { email: `field-error-${Date.now()}@example.test`, slug: `field-error-${Date.now()}` })
  await waitForAvailable(page)
  await submitRegistration(page)
  await expect(page.locator('#email-error')).toHaveText('This email is not allowed.')
  await expect(page.locator('#email')).toBeFocused()
})

test('prevents a second submission while the first registration request is pending', async ({ page }) => {
  await fulfillAvailableSlug(page)
  let requests = 0
  let releaseRequest: (() => void) | undefined
  const pendingResponse = new Promise<void>((resolve) => { releaseRequest = resolve })
  await page.route('**/api/v1/auth/register', async (route) => {
    requests += 1
    await pendingResponse
    await route.fulfill({ status: 201, contentType: 'application/json', body: JSON.stringify({ user: {}, session_established: false }) })
  })
  await openRegistration(page)
  await fillRegistration(page, { email: `double-${Date.now()}@example.test`, slug: `double-${Date.now()}` })
  await waitForAvailable(page)
  const submit = page.locator('form button[type="submit"]')
  await submit.click()
  await expect(submit).toHaveAttribute('aria-busy', 'true')
  await expect(submit.locator('[aria-live="polite"]')).toHaveText(/creating account/i)
  await submit.click({ force: true })
  expect(requests).toBe(1)
  releaseRequest?.()
  await expect(page).toHaveURL(/\/en\/auth\/login\?created=1/)
})

test('recovers to sign in when the server committed registration but could not establish a session', async ({ page }) => {
  await fulfillAvailableSlug(page)
  await page.route('**/api/v1/auth/register', (route) => route.fulfill({
    status: 201,
    contentType: 'application/json',
    body: JSON.stringify({ user: { id: 9 }, session_established: false, message: 'Account created' }),
  }))
  await openRegistration(page)
  await fillRegistration(page, { email: `recovery-${Date.now()}@example.test`, slug: `recovery-${Date.now()}` })
  await waitForAvailable(page)
  await submitRegistration(page)
  await expect(page).toHaveURL(/\/en\/auth\/login\?created=1/)
  await expect(page.locator('#email')).toBeVisible()
})

test('ignores an out-of-order subdomain response and uses only the latest value', async ({ page }) => {
  let firstResponseStarted: (() => void) | undefined
  const firstRequestStarted = new Promise<void>((resolve) => { firstResponseStarted = resolve })
  await page.route('**/api/v1/auth/check-subdomain?**', async (route) => {
    const slug = new URL(route.request().url()).searchParams.get('slug')
    if (slug === 'race-old') {
      firstResponseStarted?.()
      await new Promise((resolve) => setTimeout(resolve, 1200))
    }
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ available: slug === 'race-new', slug }) })
  })
  await openRegistration(page)
  const slugInput = page.getByLabel(/website subdomain/i)
  await slugInput.fill('race-old')
  await firstRequestStarted
  await slugInput.fill('race-new')
  await expect(page.getByRole('status').filter({ hasText: /available/i })).toBeVisible()
  await expect(page.getByRole('status').filter({ hasText: /not available/i })).toHaveCount(0)
})
