import { expect, test, type Page } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const landingUrl = (process.env.WS05_LANDING_URL || 'http://localhost:13101').replace(/\/$/, '')
const dashboardUrl = (process.env.WS05_DASHBOARD_URL || 'http://localhost:13102').replace(/\/$/, '')
const ownerEmail = 'phase2a-owner-a@example.test'
const ownerPassword = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'

async function expectPathname(page: Page, pathname: string) {
  await expect.poll(() => new URL(page.url()).pathname, { timeout: 30_000 }).toBe(pathname)
}

async function expectSignedInDashboard(page: Page) {
  await expect.poll(async () => page.evaluate(() => {
    const rawAuthState = localStorage.getItem('auth-storage')
    const authState = rawAuthState ? JSON.parse(rawAuthState) as { state?: { user?: unknown; token?: string } } : null
    const query = new URLSearchParams(location.search)
    const hasCredentialQuery = ['token', 'jwt', 'user'].some((name) => query.has(name))
    return { pathname: location.pathname, hasUser: Boolean(authState?.state?.user), hasStoredToken: Boolean(authState?.state?.token), hasCredentialQuery }
  }).catch(() => ({ pathname: '', hasUser: false, hasStoredToken: false, hasCredentialQuery: false })), { timeout: 30_000 })
    .toEqual({ pathname: '/en/dashboard', hasUser: true, hasStoredToken: false, hasCredentialQuery: false })
  const sessionCookie = (await page.context().cookies('http://localhost:18081')).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie?.httpOnly).toBe(true)
}

async function signIn(page: Page, password: string) {
  await routeContainerApiToCandidate(page)
  page.on('framenavigated', (frame) => {
    if (frame === page.mainFrame()) console.log(`BROWSER_NAVIGATION ${new URL(frame.url()).pathname}`)
  })
  page.on('pageerror', (error) => console.log(`BROWSER_PAGE_ERROR ${error.name}: ${error.message}`))
  page.on('console', (message) => {
    if (message.type() === 'error') console.log(`BROWSER_CONSOLE_ERROR ${message.text()}`)
  })
  page.on('requestfailed', (request) => console.log(`BROWSER_REQUEST_FAILED ${request.method()} ${new URL(request.url()).pathname} ${request.failure()?.errorText ?? 'unknown'}`))
  page.on('response', (response) => {
    const pathname = new URL(response.url()).pathname
    if (pathname === '/api/v1/auth/login' || pathname === '/api/v1/auth/session') {
      console.log(`BROWSER_AUTH_RESPONSE ${pathname} ${response.status()}`)
    }
  })
  await page.goto(`${landingUrl}/en/auth/login`, { waitUntil: 'domcontentloaded' })
  await expect(page.locator('form')).toHaveAttribute('method', 'post')
  await page.waitForLoadState('networkidle')
  await page.getByRole('textbox', { name: /email/i }).fill(ownerEmail)
  await page.getByRole('textbox', { name: /password/i }).fill(password)
  await page.getByRole('button', { name: /^sign in$/i }).click()
}

test('invalid credentials stay on login and show an error', async ({ page }) => {
  await signIn(page, 'Wrong-WS05-Pass!')

  await expect(page).toHaveURL(`${landingUrl}/en/auth/login`)
  expect(new URL(page.url()).searchParams.has('password')).toBe(false)
  await expect(page.getByText(/invalid|failed|credentials/i).last()).toBeVisible()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('logout clears access and denies the next protected dashboard visit', async ({ page }) => {
  await signIn(page, ownerPassword)
  expect(new URL(page.url()).searchParams.has('password')).toBe(false)
  await expectSignedInDashboard(page)
  expect(new URL(page.url()).search).toBe('')
  await page.reload({ waitUntil: 'domcontentloaded' })
  await expectSignedInDashboard(page)

  await page.getByRole('button', { name: /phase2a-owner-a@example\.test/i }).click()
  await page.getByRole('button', { name: /logout/i }).click()
  await expect(page).toHaveURL(`${landingUrl}/en/auth/login`)
  await expect(page.getByRole('textbox', { name: /email/i })).toBeVisible()

  await page.goto(`${dashboardUrl}/en/dashboard/inventory`, { waitUntil: 'domcontentloaded' })
  await expectPathname(page, '/en/auth/login')
  await expect(page.getByRole('textbox', { name: /password/i })).toBeVisible()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})
