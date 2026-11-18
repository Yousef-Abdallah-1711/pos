import { execFileSync } from 'node:child_process'
import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { expect, test, type Page } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const landingUrl = (process.env.WS05_LANDING_URL || 'http://localhost:13111').replace(/\/$/, '')
const landingRedirectUrl = (process.env.WS05_LANDING_REDIRECT_URL || 'http://localhost:13111').replace(/\/$/, '')
const dashboardUrl = (process.env.WS05_DASHBOARD_URL || 'http://localhost:13112').replace(/\/$/, '')
const apiOrigin = new URL(process.env.WS05_API_URL || 'http://localhost:18082/api/v1').origin
const ownerEmail = 'phase2a-owner-a@example.test'
const ownerName = 'WS05 Tenant A Owner'
const ownerPassword = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'
const composeProject = process.env.WS05_COMPOSE_PROJECT || ''

type AuthTimelineEvent = {
  elapsedMs: number
  kind: 'request' | 'response' | 'navigation' | 'request-failed'
  method?: string
  path: string
  status?: number
  resourceType?: string
  error?: string
}

type AuthTimeline = { startedAt: number; events: AuthTimelineEvent[] }
const authTimelines = new WeakMap<Page, AuthTimeline>()
let authTimelinePages: Page[] = []

function startAuthTimeline(page: Page) {
  if (authTimelines.has(page)) return
  const timeline: AuthTimeline = { startedAt: Date.now(), events: [] }
  authTimelines.set(page, timeline)
  authTimelinePages.push(page)
  const elapsedMs = () => Date.now() - timeline.startedAt
  const isTimelinePath = (pathname: string) => pathname.startsWith('/api/v1/auth/') || pathname.startsWith('/api/v1/products/')
  const sanitizedPath = (pathname: string) => pathname.replace(/^\/api\/v1\/products\/\d+(?=\/|$)/, '/api/v1/products/:id')

  page.on('request', (request) => {
    const url = new URL(request.url())
    if (!isTimelinePath(url.pathname)) return
    timeline.events.push({
      elapsedMs: elapsedMs(),
      kind: 'request',
      method: request.method(),
      path: sanitizedPath(url.pathname),
      resourceType: request.resourceType(),
    })
  })
  page.on('response', (response) => {
    const url = new URL(response.url())
    if (!isTimelinePath(url.pathname)) return
    timeline.events.push({
      elapsedMs: elapsedMs(),
      kind: 'response',
      method: response.request().method(),
      path: sanitizedPath(url.pathname),
      status: response.status(),
      resourceType: response.request().resourceType(),
    })
  })
  page.on('requestfailed', (request) => {
    const url = new URL(request.url())
    if (!isTimelinePath(url.pathname)) return
    timeline.events.push({
      elapsedMs: elapsedMs(),
      kind: 'request-failed',
      method: request.method(),
      path: sanitizedPath(url.pathname),
      resourceType: request.resourceType(),
      error: request.failure()?.errorText ?? 'unknown',
    })
  })
  page.on('framenavigated', (frame) => {
    if (frame !== page.mainFrame()) return
    let pathname: string
    try {
      pathname = new URL(frame.url()).pathname
    } catch {
      return
    }
    timeline.events.push({ elapsedMs: elapsedMs(), kind: 'navigation', path: pathname })
  })
}

test.beforeEach(async ({ page }) => {
  authTimelinePages = []
  startAuthTimeline(page)
})

test.afterEach(async ({ page }, testInfo) => {
  startAuthTimeline(page)
  for (const [index, timelinePage] of authTimelinePages.entries()) {
    const timeline = authTimelines.get(timelinePage)
    if (!timeline) continue
    const requestCounts: Record<string, number> = {}
    const responseCounts: Record<string, number> = {}
    for (const event of timeline.events) {
      if (event.kind === 'request') {
        const key = `${event.method} ${event.path}`
        requestCounts[key] = (requestCounts[key] ?? 0) + 1
      } else if (event.kind === 'response') {
        const key = `${event.method} ${event.path} ${event.status}`
        responseCounts[key] = (responseCounts[key] ?? 0) + 1
      }
    }
    const evidence = JSON.stringify({ requestCounts, responseCounts, events: timeline.events }, null, 2)
    await testInfo.attach(`auth-browser-timeline-${index + 1}.json`, {
      body: evidence,
      contentType: 'application/json',
    })
    const evidenceDirectory = process.env.WS05_AUTH_EVIDENCE_DIR
    if (evidenceDirectory) {
      await mkdir(evidenceDirectory, { recursive: true })
      const safeTitle = testInfo.title.replace(/[^a-z0-9]+/gi, '-').replace(/^-|-$/g, '').toLowerCase()
      const fileName = `${safeTitle}-retry-${testInfo.retry}-page-${index + 1}.json`
      await writeFile(path.join(evidenceDirectory, fileName), evidence, { mode: 0o600 })
    }
  }
  authTimelinePages = []
})

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
  const loginResponse = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/login'
  }, { timeout: 30_000 })
  await page.goto(`${landingUrl}/en/auth/login`, { waitUntil: 'domcontentloaded' })
  await expect(page.locator('form')).toHaveAttribute('method', 'post')
  await landingSessionCheck
  await page.getByRole('textbox', { name: /email/i }).fill(ownerEmail)
  await page.getByRole('textbox', { name: /password/i }).fill(password)
  await page.getByRole('button', { name: /^sign in$/i }).click()
  return (await loginResponse).status()
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

test('forged tenant headers cannot disclose another tenant product', async ({ page }, testInfo) => {
  test.skip(!composeProject, 'requires the isolated GitHub Actions authentication Compose project')
  expect(composeProject).toMatch(/^servio-auth-gha-\d+$/)

  const postgresContainerIDs = execFileSync('docker', [
    'ps',
    '--filter', `label=com.docker.compose.project=${composeProject}`,
    '--filter', 'label=com.docker.compose.service=postgres',
    '--filter', 'label=com.pos.disposable=true',
    '--filter', 'label=com.pos.workstream=WS05',
    '--format', '{{.ID}}',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  expect(postgresContainerIDs).toHaveLength(1)
  const postgresContainerID = postgresContainerIDs[0]
  const postgresLabels = JSON.parse(execFileSync('docker', [
    'inspect', '--format', '{{json .Config.Labels}}', postgresContainerID,
  ], { encoding: 'utf8' })) as Record<string, string>
  expect(postgresLabels['com.docker.compose.project']).toBe(composeProject)
  expect(postgresLabels['com.docker.compose.service']).toBe('postgres')
  expect(postgresLabels['com.pos.disposable']).toBe('true')
  expect(postgresLabels['com.pos.workstream']).toBe('WS05')

  const targetProduct = execFileSync('docker', [
    'exec', postgresContainerID, 'psql', '-U', 'postgres', '-d', 'pos_phase2a_e2e', '-AtF', '|',
    '-c', "SELECT id, tenant_id, restaurant_id FROM products WHERE name_en = 'WS05 E2E Burger B' LIMIT 1",
  ], { encoding: 'utf8' }).trim().split('|').map(Number)
  expect(targetProduct).toHaveLength(3)
  expect(targetProduct.every(Number.isSafeInteger)).toBe(true)

  const ownerScope = execFileSync('docker', [
    'exec', postgresContainerID, 'psql', '-U', 'postgres', '-d', 'pos_phase2a_e2e', '-AtF', '|',
    '-c', "SELECT tenant_id, restaurant_id FROM products WHERE name_en = 'WS05 E2E Burger A' LIMIT 1",
  ], { encoding: 'utf8' }).trim().split('|').map(Number)
  expect(ownerScope).toHaveLength(2)
  expect(ownerScope.every(Number.isSafeInteger)).toBe(true)

  const scopedRowCount = execFileSync('docker', [
    'exec', '-e', 'PGPASSWORD=ws05-local-only-app', postgresContainerID,
    'psql', '-U', 'pos_app', '-d', 'pos_phase2a_e2e', '-At',
    '-c', `SELECT count(*) FROM products WHERE id = ${targetProduct[0]} AND tenant_id = ${ownerScope[0]} AND restaurant_id = ${ownerScope[1]}`,
  ], { encoding: 'utf8' }).trim()
  expect(scopedRowCount).toBe('0')

  const rlsEnabled = execFileSync('docker', [
    'exec', postgresContainerID, 'psql', '-U', 'postgres', '-d', 'pos_phase2a_e2e', '-At',
    '-c', "SELECT relrowsecurity FROM pg_class WHERE oid = 'public.products'::regclass",
  ], { encoding: 'utf8' }).trim()

  const loginResponse = await signIn(page, ownerPassword)
  expect(loginResponse).toBe(200)
  await expectSignedInDashboard(page)
  const sessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie?.httpOnly).toBe(true)
  const response = await page.evaluate(async ({ endpoint, tenantID, restaurantID }) => {
    const apiResponse = await fetch(endpoint, {
      credentials: 'include',
      headers: { 'X-Tenant-ID': String(tenantID), 'X-Restaurant-ID': String(restaurantID) },
    })
    return { status: apiResponse.status, body: await apiResponse.text() }
  }, {
    endpoint: `${apiOrigin}/api/v1/products/${targetProduct[0]}`,
    tenantID: targetProduct[1],
    restaurantID: targetProduct[2],
  })
  expect(response.status).toBe(404)
  expect(response.body).not.toContain('WS05 E2E Burger B')
  expect(response.body).not.toContain('WS05 Tenant B Restaurant')
  expect(response.body).not.toContain('phase2a-restaurant-b')

  await testInfo.attach('cross-tenant-zero-disclosure.json', {
    body: JSON.stringify({
      responseStatus: response.status,
      responseBodyContainsTargetProduct: response.body.includes('WS05 E2E Burger B'),
      responseBodyContainsTargetRestaurant: response.body.includes('WS05 Tenant B Restaurant'),
      forgedTenantAndRestaurantHeaders: true,
      ownerScopedProductRowsFoundByPOSApp: Number(scopedRowCount),
      productsTableRowLevelSecurityEnabled: rlsEnabled === 't',
    }, null, 2),
    contentType: 'application/json',
  })
})

test('multi-tenant login selects an authorized tenant before establishing the session', async ({ page }) => {
  await routeContainerApiToCandidate(page)
  await page.goto(`${landingUrl}/en/auth/login`, { waitUntil: 'domcontentloaded' })
  await expect(page.locator('form')).toBeVisible()
  await page.getByRole('textbox', { name: /email/i }).fill('phase2a-multi-owner@example.test')
  await page.getByRole('textbox', { name: /password/i }).fill(ownerPassword)
  const loginResponse = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/login'
  })
  await page.getByRole('button', { name: /^sign in$/i }).click()
  const response = await loginResponse
  expect(response.status()).toBe(200)
  const loginBody = await response.json() as {
    multiple_tenants: boolean
    tenants: Array<{ tenant_id: number; tenant_name: string }>
  }
  expect(loginBody.multiple_tenants).toBe(true)
  expect(loginBody.tenants.length).toBeGreaterThanOrEqual(2)
  expect((await page.context().cookies(apiOrigin)).some((cookie) => cookie.name === 'poos_session')).toBe(false)

  const selectedTenant = loginBody.tenants[1]
  const confirmationResponse = page.waitForResponse((candidate) => {
    return new URL(candidate.url()).pathname === '/api/v1/auth/login/confirm'
  })
  await page.getByRole('button', { name: new RegExp(selectedTenant.tenant_name) }).click()
  expect((await confirmationResponse).status()).toBe(200)
  await expectSignedInDashboard(page)
  const selectedTenantID = await page.evaluate(() => {
    const rawAuthState = localStorage.getItem('auth-storage')
    const state = rawAuthState ? JSON.parse(rawAuthState) as { state?: { user?: { tenant_id?: number } } } : null
    return state?.state?.user?.tenant_id ?? null
  })
  expect(selectedTenantID).toBe(selectedTenant.tenant_id)
  expect(new URL(page.url()).search).toBe('')
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
  startAuthTimeline(landingPage)
  await landingPage.goto(`${landingUrl}/en`, { waitUntil: 'domcontentloaded' })
  const landingAccountButton = landingPage.getByRole('button', { name: new RegExp(ownerName) })
  await expect(landingAccountButton).toBeVisible({ timeout: 30_000 })
  await landingAccountButton.click()
  const accountPanel = landingPage.locator('#sv-account-panel')
  await expect(accountPanel.getByText('WS05 Tenant A Restaurant')).toBeVisible()
  expect(await accountPanel.innerText()).not.toMatch(/phone|user id|tenant id|session|token/i)
  await expect(accountPanel.getByRole('link', { name: /dashboard/i })).toHaveAttribute('href', `${dashboardUrl}/en/dashboard`)
  await expect(accountPanel.getByRole('link', { name: /settings/i })).toHaveAttribute('href', `${dashboardUrl}/en/dashboard/settings/account`)
  await captureScreenshot(landingPage, 'landing-authenticated-account-menu')

  const dashboardTab = await page.context().newPage()
  startAuthTimeline(dashboardTab)
  await dashboardTab.goto(`${dashboardUrl}/en/dashboard`, { waitUntil: 'domcontentloaded' })
  await expectSignedInDashboard(dashboardTab)
  await dashboardTab.reload({ waitUntil: 'domcontentloaded' })
  await expectSignedInDashboard(dashboardTab)
  await landingPage.reload({ waitUntil: 'domcontentloaded' })
  await expect(landingPage.getByRole('button', { name: new RegExp(ownerName) })).toBeVisible({ timeout: 30_000 })
  await page.reload({ waitUntil: 'domcontentloaded' })
  await expectSignedInDashboard(page)
  await dashboardTab.close()

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
  const toggledDarkTheme = !initialDarkTheme
  await expect.poll(() => landingPage.evaluate(() => document.documentElement.classList.contains('dark'))).toBe(toggledDarkTheme)
  await landingPage.reload({ waitUntil: 'domcontentloaded' })
  await expect.poll(() => landingPage.evaluate(() => document.documentElement.classList.contains('dark'))).toBe(toggledDarkTheme)
  await captureScreenshot(landingPage, 'landing-authenticated-theme-toggle')
  await landingPage.locator('.sv-theme-toggle').click()
  await expect.poll(() => landingPage.evaluate(() => document.documentElement.classList.contains('dark'))).toBe(initialDarkTheme)

  await landingPage.locator('a.sv-language').click()
  await expectPathname(landingPage, '/ar')
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
  await expect(mobileMenuButton).toHaveAttribute('aria-label', /.+/)
  await expect(mobileMenuButton).toHaveAttribute('aria-expanded', 'false')
  await mobileMenuButton.click()
  await expect(mobileDrawer).toBeVisible()
  await expect(mobileDrawer).toHaveAttribute('role', 'dialog')
  await expect(mobileDrawer).toHaveAttribute('aria-modal', 'true')
  await expect(mobileDrawer).toHaveAttribute('aria-label', /.+/)
  await expect(mobileMenuButton).toHaveAttribute('aria-expanded', 'true')
  const mobileAccount = mobileDrawer.locator('.sv-mobile-account')
  await expect(mobileAccount.getByText('WS05 Tenant A Restaurant')).toBeVisible()
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

test('Redis outage returns a private 503 without logging out, then recovers on retry', async ({ page }) => {
  test.skip(!composeProject, 'requires the isolated GitHub Actions authentication Compose project')
  test.setTimeout(180_000)

  expect(composeProject).toMatch(/^servio-auth-gha-\d+$/)
  const redisContainerIDs = execFileSync('docker', [
    'ps',
    '--filter', `label=com.docker.compose.project=${composeProject}`,
    '--filter', 'label=com.docker.compose.service=redis',
    '--format', '{{.ID}}',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  expect(redisContainerIDs).toHaveLength(1)

  const redisContainerID = redisContainerIDs[0]
  const redisLabels = JSON.parse(execFileSync('docker', [
    'inspect', '--format', '{{json .Config.Labels}}', redisContainerID,
  ], { encoding: 'utf8' })) as Record<string, string>
  expect(redisLabels['com.docker.compose.project']).toBe(composeProject)
  expect(redisLabels['com.docker.compose.service']).toBe('redis')
  expect(redisLabels['com.pos.disposable']).toBe('true')
  expect(redisLabels['com.pos.workstream']).toBe('WS05')

  const loginResponse = await signIn(page, ownerPassword)
  expect(loginResponse).toBe(200)
  await expectSignedInDashboard(page)

  let redisStopped = false
  try {
    execFileSync('docker', ['exec', redisContainerID, 'redis-cli', 'SAVE'], { stdio: 'ignore' })
    execFileSync('docker', ['stop', '--time', '1', redisContainerID], { stdio: 'ignore' })
    redisStopped = true

    const outageResponsePromise = page.waitForResponse((response) => {
      return new URL(response.url()).pathname === '/api/v1/auth/session'
    }, { timeout: 45_000 })
    await page.evaluate(() => window.dispatchEvent(new Event('focus')))
    const outageResponse = await outageResponsePromise
    expect(outageResponse.status()).toBe(503)
    expect(outageResponse.headers()['cache-control']).toContain('private')
    expect(outageResponse.headers()['set-cookie']).toBeUndefined()
    await expect(page.getByRole('alert')).toBeVisible()
    await expectPathname(page, '/en/dashboard')
    await expect.poll(() => page.evaluate(() => {
      const rawAuthState = localStorage.getItem('auth-storage')
      const authState = rawAuthState ? JSON.parse(rawAuthState) as { state?: { user?: unknown; token?: string } } : null
      return { hasUser: Boolean(authState?.state?.user), hasStoredToken: Boolean(authState?.state?.token) }
    })).toEqual({ hasUser: true, hasStoredToken: false })
    await captureScreenshot(page, 'dashboard-session-unavailable')

    execFileSync('docker', ['start', redisContainerID], { stdio: 'ignore' })
    redisStopped = false
    await expect.poll(() => execFileSync('docker', [
      'inspect', '--format', '{{.State.Health.Status}}', redisContainerID,
    ], { encoding: 'utf8' }).trim(), { timeout: 30_000 }).toBe('healthy')

    const recoveryResponsePromise = page.waitForResponse((response) => {
      return new URL(response.url()).pathname === '/api/v1/auth/session' && response.status() === 200
    }, { timeout: 30_000 })
    await page.getByRole('alert').getByRole('button', { name: /try again/i }).click()
    await recoveryResponsePromise
    await expectSignedInDashboard(page)
    await captureScreenshot(page, 'dashboard-session-recovered')
  } finally {
    if (redisStopped) execFileSync('docker', ['start', redisContainerID], { stdio: 'ignore' })
  }
})
