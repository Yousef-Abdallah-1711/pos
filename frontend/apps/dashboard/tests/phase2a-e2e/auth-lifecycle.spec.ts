import { mkdir } from 'node:fs/promises'
import path from 'node:path'
import { expect, test, type Page } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const landingUrl = (process.env.WS05_LANDING_URL || 'http://localhost:13111').replace(/\/$/, '')
const landingRedirectUrl = (process.env.WS05_LANDING_REDIRECT_URL || 'http://localhost:13111').replace(/\/$/, '')
const dashboardUrl = (process.env.WS05_DASHBOARD_URL || 'http://localhost:13112').replace(/\/$/, '')
const apiOrigin = new URL(process.env.WS05_API_URL || 'http://localhost:18082/api/v1').origin
const ownerEmail = 'phase2a-owner-a@example.test'
const ownerName = 'WS05 A Owner'
const ownerPassword = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'

async function captureScreenshot(page: Page, name: string) {
  const evidenceDirectory = process.env.WS05_AUTH_EVIDENCE_DIR
  if (!evidenceDirectory) return
  await mkdir(evidenceDirectory, { recursive: true })
  await page.screenshot({ path: path.join(evidenceDirectory, `${name}.png`), fullPage: true })
}

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
  }).catch(() => ({ pathname: '', hasUser: false, hasStoredToken: false, hasCredentialQuery: false })), { timeout: 90_000 })
    .toEqual({ pathname: '/en/dashboard', hasUser: true, hasStoredToken: false, hasCredentialQuery: false })
  const sessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie?.httpOnly).toBe(true)
}

async function signIn(page: Page, password: string) {
  await routeContainerApiToCandidate(page)
  const landingSessionCheck = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/session'
  }, { timeout: 30_000 })
  let resolveLoginResponse!: (status: number) => void
  const loginResponse = new Promise<number>((resolve) => {
    resolveLoginResponse = resolve
  })
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
    if (pathname === '/api/v1/auth/login' || pathname === '/api/v1/auth/session' || pathname === '/api/v1/auth/logout') {
      console.log(`BROWSER_AUTH_RESPONSE ${pathname} ${response.status()}`)
    }
    if (pathname === '/api/v1/auth/login') resolveLoginResponse(response.status())
  })
  await page.goto(`${landingUrl}/en/auth/login`, { waitUntil: 'domcontentloaded' })
  await expect(page.locator('form')).toHaveAttribute('method', 'post')
  await landingSessionCheck
  await page.getByRole('textbox', { name: /email/i }).fill(ownerEmail)
  await page.getByRole('textbox', { name: /password/i }).fill(password)
  await page.getByRole('button', { name: /^sign in$/i }).click()
  return loginResponse
}

test('guest can open Landing and follow the canonical Login link', async ({ page }) => {
  await routeContainerApiToCandidate(page)
  const guestSessionCheck = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/session'
  }, { timeout: 30_000 })
  const response = await page.goto(`${landingUrl}/en`, { waitUntil: 'domcontentloaded' })
  expect(response?.status()).toBe(200)
  const html = await response?.text()
  expect(html).not.toContain(ownerEmail)
  expect(html).not.toContain(ownerName)
  await guestSessionCheck
  const loginLink = page.getByRole('link', { name: /log in|login/i }).first()
  await expect(loginLink).toBeVisible()
  await captureScreenshot(page, 'landing-guest')
  await loginLink.click()
  await expect(page).toHaveURL(`${landingUrl}/en/auth/login`)
  await expect(page.locator('form')).toBeVisible()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('invalid credentials stay on login and show an error', async ({ page }) => {
  const loginResponse = await signIn(page, 'Wrong-WS05-Pass!')

  expect(loginResponse).toBe(401)
  await expect(page).toHaveURL(`${landingUrl}/en/auth/login`)
  expect(new URL(page.url()).searchParams.has('password')).toBe(false)
  await expect(page.getByText(/invalid|failed|credentials/i).last()).toBeVisible()
  await captureScreenshot(page, 'landing-invalid-login')
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('logout clears access and denies the next protected dashboard visit', async ({ page }) => {
  test.setTimeout(240_000)
  const loginResponse = await signIn(page, ownerPassword)
  expect(loginResponse).toBe(200)
  expect(new URL(page.url()).searchParams.has('password')).toBe(false)
  await expectSignedInDashboard(page)
  await captureScreenshot(page, 'dashboard-authenticated')

  const landingPage = await page.context().newPage()
  await landingPage.goto(`${landingUrl}/en`, { waitUntil: 'domcontentloaded' })
  const landingAccountButton = landingPage.getByRole('button', { name: new RegExp(ownerName) })
  await expect(landingAccountButton).toBeVisible({ timeout: 30_000 })
  await landingAccountButton.click()
  const accountPanel = landingPage.locator('#sv-account-panel')
  await expect(accountPanel.getByText('WS05 A Restaurant')).toBeVisible()
  expect(await accountPanel.innerText()).not.toMatch(/phone|user id|tenant id|session|token/i)
  await expect(accountPanel.getByRole('link', { name: /dashboard/i })).toHaveAttribute('href', `${dashboardUrl}/en/dashboard`)
  await expect(accountPanel.getByRole('link', { name: /settings/i })).toHaveAttribute('href', `${dashboardUrl}/en/dashboard/settings/account`)
  await captureScreenshot(landingPage, 'landing-authenticated-account-menu')

  await landingPage.keyboard.press('Escape')
  await expect(accountPanel).toBeHidden()
  await expect(landingAccountButton).toBeFocused()
  await landingAccountButton.press('Enter')
  await expect(accountPanel).toBeVisible()

  for (const width of [320, 390, 768, 1024, 1440]) {
    await landingPage.setViewportSize({ width, height: 900 })
    const viewport = await landingPage.evaluate(() => ({
      clientWidth: document.documentElement.clientWidth,
      scrollWidth: document.documentElement.scrollWidth,
    }))
    expect(viewport.scrollWidth).toBeLessThanOrEqual(viewport.clientWidth)
  }
  await landingPage.setViewportSize({ width: 1280, height: 720 })

  const themeToggle = landingPage.locator('.sv-theme-toggle')
  const initialDarkTheme = await landingPage.evaluate(() => document.documentElement.classList.contains('dark'))
  await themeToggle.click()
  await expect.poll(() => landingPage.evaluate(() => document.documentElement.classList.contains('dark'))).toBe(!initialDarkTheme)
  await captureScreenshot(landingPage, 'landing-authenticated-theme-toggle')
  await themeToggle.click()
  await expect.poll(() => landingPage.evaluate(() => document.documentElement.classList.contains('dark'))).toBe(initialDarkTheme)

  await landingPage.goto(`${landingUrl}/ar`, { waitUntil: 'domcontentloaded' })
  await expect(landingPage.getByRole('button', { name: new RegExp(ownerName) })).toBeVisible({ timeout: 30_000 })
  const arabicDocument = await landingPage.evaluate(() => ({ language: document.documentElement.lang, direction: document.documentElement.dir }))
  expect(arabicDocument).toEqual({ language: 'ar', direction: 'rtl' })
  for (const width of [320, 390, 768, 1024, 1440]) {
    await landingPage.setViewportSize({ width, height: 800 })
    await expect.poll(() => landingPage.evaluate(() => document.documentElement.scrollWidth <= document.documentElement.clientWidth)).toBe(true)
  }
  await landingPage.setViewportSize({ width: 320, height: 800 })
  const mobileMenuButton = landingPage.locator('.sv-menu-button')
  const mobileDrawer = landingPage.locator('#sv-mobile-navigation')
  await mobileMenuButton.click()
  await expect(mobileDrawer).toBeVisible()
  const mobileAccount = mobileDrawer.locator('.sv-mobile-account')
  await expect(mobileAccount.getByText('WS05 A Restaurant')).toBeVisible()
  await expect(mobileAccount.locator('a').nth(0)).toHaveAttribute('href', `${dashboardUrl}/ar/dashboard`)
  await expect(mobileAccount.locator('a').nth(1)).toHaveAttribute('href', `${dashboardUrl}/ar/dashboard/settings/account`)
  const signOutBounds = await mobileAccount.locator('.sv-mobile-account__signout').boundingBox()
  expect(signOutBounds?.height ?? 0).toBeGreaterThanOrEqual(44)
  await captureScreenshot(landingPage, 'landing-authenticated-arabic-320')
  await landingPage.keyboard.press('Escape')
  await expect(mobileDrawer).toBeHidden()
  await expect(mobileMenuButton).toBeFocused()
  await mobileMenuButton.click()
  await expect(mobileDrawer).toBeVisible()

  const firstSessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(firstSessionCookie?.httpOnly).toBe(true)
  expect(firstSessionCookie).toBeDefined()
  const firstRevokedCookie = `${firstSessionCookie!.name}=${firstSessionCookie!.value}`
  const landingLogoutResponse = landingPage.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/logout'
  }, { timeout: 30_000 })
  await mobileAccount.locator('.sv-mobile-account__signout').click()
  expect((await landingLogoutResponse).status()).toBe(204)
  await mobileMenuButton.click()
  await expect(mobileDrawer.getByRole('link', { name: /تسجيل الدخول/i })).toBeVisible()
  await captureScreenshot(landingPage, 'landing-after-logout')

  const landingLogoutReplay = await page.request.get(`${apiOrigin}/api/v1/auth/session`, {
    headers: { Cookie: firstRevokedCookie, Origin: landingUrl },
  })
  expect(landingLogoutReplay.status()).toBe(401)
  await landingPage.close()

  expect(new URL(page.url()).search).toBe('')
  await page.goto(`${dashboardUrl}/en/dashboard/inventory`, { waitUntil: 'domcontentloaded' })
  await expectPathname(page, '/en/auth/login')

  await page.unrouteAll({ behavior: 'ignoreErrors' })
  const secondLoginResponse = await signIn(page, ownerPassword)
  expect(secondLoginResponse).toBe(200)
  await expectSignedInDashboard(page)
  await page.reload({ waitUntil: 'domcontentloaded' })
  await expectSignedInDashboard(page)

  const secondSessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(secondSessionCookie?.httpOnly).toBe(true)
  expect(secondSessionCookie?.sameSite).toBe('Lax')
  expect(secondSessionCookie).toBeDefined()
  const secondRevokedCookie = `${secondSessionCookie!.name}=${secondSessionCookie!.value}`

  await page.getByRole('button', { name: new RegExp(ownerName) }).click()
  const dashboardLogoutResponse = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/logout'
  }, { timeout: 30_000 })
  await page.getByRole('button', { name: /logout/i }).click()
  expect((await dashboardLogoutResponse).status()).toBe(204)
  await expect(page).toHaveURL(`${landingRedirectUrl}/en/auth/login`)
  await expect(page.getByRole('textbox', { name: /email/i })).toBeVisible()
  await expect.poll(async () => (await page.context().cookies(apiOrigin)).some((cookie) => cookie.name === 'poos_session')).toBe(false)
  await captureScreenshot(page, 'landing-after-dashboard-logout')

  const dashboardLogoutReplay = await page.request.get(`${apiOrigin}/api/v1/auth/session`, {
    headers: { Cookie: secondRevokedCookie, Origin: landingUrl },
  })
  expect(dashboardLogoutReplay.status()).toBe(401)

  await page.goto(`${dashboardUrl}/en/dashboard/inventory`, { waitUntil: 'domcontentloaded' })
  await expectPathname(page, '/en/auth/login')
  await expect(page.getByRole('textbox', { name: /password/i })).toBeVisible()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})
