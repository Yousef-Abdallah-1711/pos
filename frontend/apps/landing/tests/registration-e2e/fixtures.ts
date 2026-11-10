import { expect, type Page } from '@playwright/test'

export const landingURL = (process.env.REG_E2E_LANDING_URL ?? 'http://localhost:14211').replace(/\/$/, '')
export const dashboardURL = (process.env.REG_E2E_DASHBOARD_URL ?? 'http://localhost:14212').replace(/\/$/, '')
export const apiURL = (process.env.REG_E2E_API_URL ?? 'http://localhost:19482/api/v1').replace(/\/$/, '')
export const suffix = `${Date.now()}-${Math.floor(Math.random() * 1_000_000)}`
export const testAccount = {
  name: 'Registration E2E Owner',
  email: `registration-e2e-${suffix}@example.test`,
  password: 'E2E-Disposable-Password-2026!',
  restaurantName: `Registration E2E Restaurant ${suffix}`,
  slug: `registration-e2e-${suffix}`,
  phone: '+12025550123',
}

export async function openRegistration(page: Page, locale = 'en') {
  await page.goto(`${landingURL}/${locale}/auth/register`, { waitUntil: 'domcontentloaded' })
  await expect(page.locator('form')).toBeVisible()
}

export async function fillRegistration(page: Page, overrides: Partial<typeof testAccount> = {}) {
  const account = { ...testAccount, ...overrides }
  await page.locator('#name').fill(account.name)
  await page.locator('#email').fill(account.email)
  await page.locator('#restaurantName').fill(account.restaurantName)
  await page.locator('#restaurantSlug').fill(account.slug)
  if (account.phone) await page.locator('#phone').fill(account.phone)
  await page.locator('#password').fill(account.password)
  await page.locator('#confirmPassword').fill(account.password)
}

export async function waitForAvailable(page: Page) {
  await expect(page.getByRole('status').filter({ hasText: /available|متاح/i })).toBeVisible()
}

export async function submitRegistration(page: Page) {
  await page.getByRole('button', { name: /create account|إنشاء حساب/i }).click()
}

export async function fulfillAvailableSlug(page: Page) {
  await page.route('**/api/v1/auth/check-subdomain?**', async (route) => {
    const slug = new URL(route.request().url()).searchParams.get('slug') ?? ''
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ available: true, slug }) })
  })
}

export async function expectInlineError(page: Page, inputId: string, message: RegExp) {
  const input = page.locator(`#${inputId}`)
  await expect(input).toHaveAttribute('aria-invalid', 'true')
  await expect(page.locator(`#${inputId}-error`)).toContainText(message)
  await expect(input).toHaveAttribute('aria-describedby', new RegExp(`${inputId}-error`))
}
