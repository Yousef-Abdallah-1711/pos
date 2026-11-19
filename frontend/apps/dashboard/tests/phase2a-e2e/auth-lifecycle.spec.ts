import { execFileSync } from 'node:child_process'
import { createHash } from 'node:crypto'
import { mkdir, readFile, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { expect, test, type Page, type Route } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const landingUrl = (process.env.WS05_LANDING_URL || 'http://localhost:13111').replace(/\/$/, '')
const landingRedirectUrl = (process.env.WS05_LANDING_REDIRECT_URL || 'http://localhost:13111').replace(/\/$/, '')
const dashboardUrl = (process.env.WS05_DASHBOARD_URL || 'http://localhost:13112').replace(/\/$/, '')
const apiOrigin = new URL(process.env.WS05_API_URL || 'http://localhost:18082/api/v1').origin
const ownerEmail = 'phase2a-owner-a@example.test'
const ownerName = 'WS05 Tenant A Owner'
const otherTenantOwnerEmail = 'phase2a-owner-b@example.test'
const ownerPassword = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'
const superAdminPassword = process.env.E2E_SUPERADMIN_PASSWORD || 'WS05-SuperAdmin-Only-2026!'
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

function resetIsolatedAuthRateLimitCounters() {
  if (!composeProject) return
  if (!/^servio-auth-gha-\d+$/.test(composeProject)) {
    throw new Error('Refusing to reset rate limits outside the disposable authentication test project')
  }

  const redisContainerIDs = execFileSync('docker', [
    'ps', '-q',
    '--filter', `label=com.docker.compose.project=${composeProject}`,
    '--filter', 'label=com.docker.compose.service=redis',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  if (redisContainerIDs.length !== 1) {
    throw new Error(`Expected one disposable Redis container, found ${redisContainerIDs.length}`)
  }

  const redisContainerID = redisContainerIDs[0]
  const rateLimitKeys = execFileSync('docker', [
    'exec', redisContainerID, 'redis-cli', '--scan', '--pattern', 'pos-saas:auth:rate-limit:v1:*',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  if (rateLimitKeys.length > 0) {
    execFileSync('docker', ['exec', redisContainerID, 'redis-cli', 'DEL', ...rateLimitKeys], { stdio: 'ignore' })
  }
}

function isolatedHealthyRedisContainerID() {
  if (!/^servio-auth-gha-\d+$/.test(composeProject)) {
    throw new Error('Refusing to inspect Redis outside the disposable authentication test project')
  }
  const redisContainerIDs = execFileSync('docker', [
    'ps', '-q',
    '--filter', `label=com.docker.compose.project=${composeProject}`,
    '--filter', 'label=com.docker.compose.service=redis',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  if (redisContainerIDs.length !== 1) {
    throw new Error(`Expected one disposable Redis container, found ${redisContainerIDs.length}`)
  }

  const redisContainerID = redisContainerIDs[0]
  const labels = JSON.parse(execFileSync('docker', [
    'inspect', '--format', '{{json .Config.Labels}}', redisContainerID,
  ], { encoding: 'utf8' })) as Record<string, string>
  const health = execFileSync('docker', [
    'inspect', '--format', '{{.State.Health.Status}}', redisContainerID,
  ], { encoding: 'utf8' }).trim()
  if (labels['com.pos.disposable'] !== 'true' || labels['com.pos.workstream'] !== 'WS05' || health !== 'healthy') {
    throw new Error('The shared Redis service is not the healthy WS05 disposable fixture')
  }
  return redisContainerID
}

async function restartIsolatedAPIForFreshLoginWindow() {
  if (!composeProject) return
  if (!/^servio-auth-gha-\d+$/.test(composeProject)) {
    throw new Error('Refusing to restart an API outside the disposable authentication test project')
  }

  const apiContainerIDs = execFileSync('docker', [
    'ps', '-q',
    '--filter', `label=com.docker.compose.project=${composeProject}`,
    '--filter', 'label=com.docker.compose.service=api',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  if (apiContainerIDs.length !== 1) {
    throw new Error(`Expected one disposable API container, found ${apiContainerIDs.length}`)
  }

  const apiContainerID = apiContainerIDs[0]
  const labels = JSON.parse(execFileSync('docker', [
    'inspect', '--format', '{{json .Config.Labels}}', apiContainerID,
  ], { encoding: 'utf8' })) as Record<string, string>
  if (labels['com.pos.disposable'] !== 'true' || labels['com.pos.workstream'] !== 'WS05') {
    throw new Error('Refusing to restart an API without the WS05 disposable fixture labels')
  }

  execFileSync('docker', ['restart', apiContainerID], { stdio: 'ignore' })
  await expect.poll(() => execFileSync('docker', [
    'inspect', '--format', '{{.State.Health.Status}}', apiContainerID,
  ], { encoding: 'utf8' }).trim(), { timeout: 60_000 }).toBe('healthy')
}

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
  // Browser scenarios share one loopback source IP. Reset only the disposable
  // Redis rate-limit windows between tests; each scenario still exercises the
  // real configured limits internally without inheriting earlier test traffic.
  resetIsolatedAuthRateLimitCounters()
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

async function expectSignedInDashboard(page: Page, allowArabicLocale = false) {
  await expect.poll(async () => page.evaluate((allowArabicLocale) => {
    const rawAuthState = localStorage.getItem('auth-storage')
    const authState = rawAuthState ? JSON.parse(rawAuthState) as { state?: { user?: unknown; token?: string } } : null
    const query = new URLSearchParams(location.search)
    const hasCredentialQuery = ['token', 'jwt', 'user'].some((name) => query.has(name))
    const pathname = location.pathname
    const expectedDashboardPath = allowArabicLocale
      ? /^\/(?:en|ar)\/dashboard$/.test(pathname)
      : pathname === '/en/dashboard'
    return { expectedDashboardPath, hasUser: Boolean(authState?.state?.user), hasStoredToken: Boolean(authState?.state?.token), hasCredentialQuery }
  }, allowArabicLocale).catch(() => ({ expectedDashboardPath: false, hasUser: false, hasStoredToken: false, hasCredentialQuery: false })), { timeout: 90_000 })
    .toEqual({ expectedDashboardPath: true, hasUser: true, hasStoredToken: false, hasCredentialQuery: false })
  const sessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie?.httpOnly).toBe(true)
  expect(sessionCookie?.expires).toBe(-1)
}

async function signIn(
  page: Page,
  password: string,
  loginUrl = `${landingUrl}/en/auth/login`,
  email = ownerEmail,
) {
  await routeContainerApiToCandidate(page)
  const landingSessionCheck = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/session'
  }, { timeout: 30_000 })
  const loginResponse = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/login'
  }, { timeout: 30_000 })
  await page.goto(loginUrl, { waitUntil: 'domcontentloaded' })
  await expect(page.locator('form')).toHaveAttribute('method', 'post')
  await landingSessionCheck
  await page.getByRole('textbox', { name: /email/i }).fill(email)
  await page.getByRole('textbox', { name: /password/i }).fill(password)
  await page.getByRole('button', { name: /^sign in$/i }).click()
  return (await loginResponse).status()
}

type ScopedCategory = { id: number; tenant_id: number; restaurant_id: number }

async function captureNextCategoriesResponse<T>(page: Page) {
  const routeMatcher = '**/api/v1/categories*'
  let resolveCapture!: (result: { status: number; body: T }) => void
  let rejectCapture!: (error: unknown) => void
  const response = new Promise<{ status: number; body: T }>((resolve, reject) => {
    resolveCapture = resolve
    rejectCapture = reject
  })

  const captureResponse = async (route: Route) => {
    const request = route.request()
    const url = new URL(request.url())
    if (request.method() !== 'GET' || url.pathname !== '/api/v1/categories') {
      await route.continue()
      return
    }

    try {
      // Read the upstream body through Playwright's route API before returning it
      // to Chromium. This keeps the evidence available if Next navigates while
      // the browser is still processing the response.
      const upstream = await route.fetch()
      const status = upstream.status()
      const body = await upstream.json() as T
      await route.fulfill({ response: upstream })
      await page.unroute(routeMatcher, captureResponse)
      resolveCapture({ status, body })
    } catch (error) {
      await route.abort().catch(() => undefined)
      await page.unroute(routeMatcher, captureResponse).catch(() => undefined)
      rejectCapture(error)
    }
  }

  await page.route(routeMatcher, captureResponse)
  return { response }
}

async function signInAndLoadTenantCategories(page: Page) {
  expect(await signIn(page, ownerPassword)).toBe(200)
  await expectSignedInDashboard(page)

  const categoriesCapture = await captureNextCategoriesResponse<ScopedCategory[]>(page)
  await page.goto(`${dashboardUrl}/en/dashboard/products/new`, { waitUntil: 'domcontentloaded' })
  const response = await categoriesCapture.response
  expect(response.status).toBe(200)
  await expect(page.getByRole('combobox').first()).toBeVisible()
  const categories = response.body
  expect(categories.length).toBeGreaterThan(0)
  return categories
}

async function seedTenantAScopedBrowserState(page: Page) {
  await page.evaluate(() => {
    localStorage.setItem('pos-storage', JSON.stringify({
      state: {
        items: [{ productId: 901, name: 'Synthetic Tenant A Item', price: 10, quantity: 1 }],
        customerName: 'Synthetic Tenant A Customer',
        customerPhone: '+10000000000',
        heldOrders: [{ id: 'synthetic-a', customerName: 'Synthetic Tenant A Customer' }],
      },
      version: 0,
    }))
    localStorage.setItem('preferences-storage', JSON.stringify({ state: { settings: { userId: 410041 } }, version: 0 }))
    localStorage.setItem('theme-builder-storage', JSON.stringify({ state: { identity: { siteTitle: 'Synthetic Tenant A Theme' } }, version: 0 }))
    localStorage.setItem('pos-products-cache', JSON.stringify([{ id: 901, name: 'Synthetic Tenant A Item' }]))
    localStorage.setItem('kds-bootstrap-snapshot-v1', JSON.stringify({ orders: [{ customer_name: 'Synthetic Tenant A Customer' }] }))
  })
}

async function loginOtherTenantInSecondTab(page: Page) {
  const otherTenantPage = await page.context().newPage()
  await otherTenantPage.bringToFront()
  startAuthTimeline(otherTenantPage)
  expect(await signIn(otherTenantPage, ownerPassword, `${landingUrl}/en/auth/login`, otherTenantOwnerEmail)).toBe(200)
  await expectSignedInDashboard(otherTenantPage)
  return otherTenantPage
}

async function readCategoriesAfterIdentityRefresh(page: Page) {
  const sessionResponsePromise = page.waitForResponse((response) => {
    const url = new URL(response.url())
    return url.pathname === '/api/v1/auth/session' && response.request().method() === 'GET'
  })
  const categoriesCapture = await captureNextCategoriesResponse<ScopedCategory[]>(page)
  await page.bringToFront()
  // Headless Chromium can keep multiple tabs logically visible without emitting
  // the OS-driven focus event. Dispatch the same browser event the provider
  // listens for after bringing this tab forward so the revalidation path is
  // deterministic in CI.
  await expect.poll(() => page.evaluate(() => document.visibilityState)).toBe('visible')
  await page.evaluate(() => window.dispatchEvent(new Event('focus')))
  const sessionResponse = await sessionResponsePromise
  const categoriesResponse = await categoriesCapture.response
  expect(sessionResponse.status()).toBe(200)
  expect(categoriesResponse.status).toBe(200)
  return { sessionResponse, categories: categoriesResponse.body }
}

async function readTenantStateCleanup(page: Page) {
  return page.evaluate(() => {
    const rawPOSState = localStorage.getItem('pos-storage')
    const posState = rawPOSState ? JSON.parse(rawPOSState).state ?? {} : {}
    const posStateText = JSON.stringify(posState)
    return {
      priorCustomerAbsent: !posStateText.includes('Synthetic Tenant A Customer'),
      priorItemAbsent: !posStateText.includes('Synthetic Tenant A Item'),
      heldOrdersEmpty: Array.isArray(posState.heldOrders) && posState.heldOrders.length === 0,
      tenantCatalogCacheAbsent: localStorage.getItem('pos-products-cache') === null,
      kdsSnapshotAbsent: localStorage.getItem('kds-bootstrap-snapshot-v1') === null,
      priorPreferencesAbsent: !localStorage.getItem('preferences-storage')?.includes('410041'),
      themeStorageCleared: !localStorage.getItem('theme-builder-storage')?.includes('Synthetic Tenant A Theme'),
    }
  })
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
  await page.waitForTimeout(1_700)
  const initialSessionRequests = authTimelines.get(page)?.events.filter((event) => {
    return event.kind === 'request' && event.method === 'GET' && event.path === '/api/v1/auth/session'
  }).length
  expect(initialSessionRequests).toBe(1)
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

test('rejects malicious redirect destinations after a real login', async ({ page }) => {
  test.setTimeout(120_000)
  await routeContainerApiToCandidate(page)
  const attacks = [
    { key: 'next', value: '//attacker.example/steal' },
    { key: 'next', value: '%2F%2Fattacker.example%2Fsteal' },
    { key: 'returnTo', value: 'https://attacker.example/' },
    { key: 'redirect', value: 'https://localhost:8443/en/dashboard' },
    { key: 'callbackUrl', value: 'https://localhost.attacker.example/en/dashboard' },
    { key: 'redirect', value: 'http://localhost:13111@attacker.example/' },
    { key: 'callbackUrl', value: '//attacker.example/callback' },
    { key: 'continue', value: String.raw`\\attacker.example\\login` },
    { key: 'next', value: 'javascript:alert(1)' },
  ]

  for (const attack of attacks) {
    // Start on the explicit Arabic route so locale negotiation cannot drop
    // the attack query while redirecting /en to the browser's preferred locale.
    const loginURL = new URL(`${landingUrl}/ar/auth/login`)
    loginURL.searchParams.set(attack.key, attack.value)
    const loginResponse = page.waitForResponse((response) => new URL(response.url()).pathname === '/api/v1/auth/login')
    await page.goto(loginURL.toString(), { waitUntil: 'domcontentloaded' })
    expect(new URL(page.url()).origin).toBe(new URL(landingUrl).origin)
    expect(new URL(page.url()).pathname).toBe('/ar/auth/login')
    expect(new URL(page.url()).searchParams.get(attack.key)).toBe(attack.value)
    await expect(page.locator('form')).toHaveAttribute('method', 'post')
    await page.locator('input[type="email"]').fill(ownerEmail)
    await page.locator('input[type="password"]').fill(ownerPassword)
    await page.locator('form button[type="submit"]').click()
    expect((await loginResponse).status()).toBe(200)
    await expectSignedInDashboard(page, true)
    expect(new URL(page.url()).origin).toBe(new URL(dashboardUrl).origin)
    expect(new URL(page.url()).pathname).toMatch(/^\/(?:en|ar)\/dashboard$/)
    expect(new URL(page.url()).search).toBe('')

    const cookie = (await page.context().cookies(apiOrigin)).find((candidate) => candidate.name === 'poos_session')
    expect(cookie).toBeDefined()
    const logout = await page.request.post(`${apiOrigin}/api/v1/auth/logout`, {
      headers: { Cookie: `poos_session=${cookie!.value}`, Origin: landingUrl },
    })
    expect(logout.status()).toBe(204)
    expect((await page.context().cookies(apiOrigin)).some((candidate) => candidate.name === 'poos_session')).toBe(false)
    // Remove the authenticated dashboard app before the next login navigation.
    // Its session monitor can otherwise race the next page.goto after logout.
    await page.goto('about:blank')
  }

  for (const malformedLocalePath of ['/xx/auth/login', '/../../attacker.example/auth/login', '/%2F%2Fattacker.example/auth/login']) {
    const response = await page.request.get(`${landingUrl}${malformedLocalePath}?next=https%3A%2F%2Fattacker.example`, {
      maxRedirects: 0,
    })
    expect(new URL(response.url()).origin).toBe(new URL(landingUrl).origin)
    if (response.status() >= 300 && response.status() < 400) {
      const location = response.headers().location
      expect(location).toBeTruthy()
      const destination = new URL(location!, landingUrl)
      expect(destination.origin).toBe(new URL(landingUrl).origin)
      expect(destination.pathname).not.toMatch(/^\/\/attacker\.example/)
    } else {
      expect([400, 404]).toContain(response.status())
    }
  }
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('rejects cross-site writes across browser session and login routes', async ({ page }) => {
  test.setTimeout(180_000)
  expect(await signIn(page, ownerPassword)).toBe(200)
  await expectSignedInDashboard(page)
  const sessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie).toBeDefined()
  const staffCookieHeader = `poos_session=${sessionCookie!.value}`
  const unsafeStaffRoutes = [
    { method: 'POST', path: '/api/v1/products', data: { name_en: 'CSRF probe' } },
    { method: 'PUT', path: '/api/v1/products/1', data: { name_en: 'CSRF probe' } },
    { method: 'DELETE', path: '/api/v1/products/1' },
  ] as const
  const origins = [undefined, 'https://attacker.example']
  const results: Array<{ method: string; path: string; origin: string; status: number }> = []

  for (const route of unsafeStaffRoutes) {
    for (const origin of origins) {
      const headers: Record<string, string> = { Cookie: staffCookieHeader }
      if (origin) headers.Origin = origin
      const response = route.method === 'POST'
        ? await page.request.post(`${apiOrigin}${route.path}`, { headers, data: route.data })
        : route.method === 'PUT'
          ? await page.request.put(`${apiOrigin}${route.path}`, { headers, data: route.data })
          : await page.request.delete(`${apiOrigin}${route.path}`, { headers })
      results.push({ method: route.method, path: route.path, origin: origin ?? 'missing', status: response.status() })
      expect(response.status()).toBe(403)
    }
  }

  const protectedPublicAuthRoutes = [
    { path: '/api/v1/auth/register', data: {} },
    { path: '/api/v1/auth/login', data: { email: 'invalid@example.test', password: 'synthetic' } },
    { path: '/api/v1/auth/login/confirm', data: { email: 'invalid@example.test', tenant_id: 1, selection_token: 'invalid' } },
    { path: '/api/v1/super-admin/auth/login', data: { email: 'invalid@example.test', password: 'synthetic' } },
  ]
  for (const route of protectedPublicAuthRoutes) {
    for (const origin of origins) {
      const headers: Record<string, string> = { 'Content-Type': 'application/json' }
      if (origin) headers.Origin = origin
      const response = await page.request.post(`${apiOrigin}${route.path}`, { headers, data: route.data })
      results.push({ method: 'POST', path: route.path, origin: origin ?? 'missing', status: response.status() })
      expect(response.status()).toBe(403)
    }
  }

  for (const route of [
    { method: 'POST', path: '/api/v1/auth/refresh', data: {} },
    { method: 'POST', path: '/api/v1/auth/logout', data: {} },
  ]) {
    for (const origin of origins) {
      const headers: Record<string, string> = { Cookie: staffCookieHeader }
      if (origin) headers.Origin = origin
      const response = await page.request.post(`${apiOrigin}${route.path}`, { headers, data: route.data })
      results.push({ method: route.method, path: route.path, origin: origin ?? 'missing', status: response.status() })
      expect(response.status()).toBe(403)
    }
  }

  for (const route of [
    { method: 'POST', path: '/api/v1/super-admin/auth/logout', cookie: 'poos_super_admin_session=synthetic-invalid' },
    { method: 'PUT', path: '/api/v1/super-admin/settings/ip-whitelist', cookie: 'poos_super_admin_session=synthetic-invalid' },
  ]) {
    for (const origin of origins) {
      const headers: Record<string, string> = { Cookie: route.cookie }
      if (origin) headers.Origin = origin
      const response = route.method === 'POST'
        ? await page.request.post(`${apiOrigin}${route.path}`, { headers, data: {} })
        : await page.request.put(`${apiOrigin}${route.path}`, { headers, data: {} })
      results.push({ method: route.method, path: route.path, origin: origin ?? 'missing', status: response.status() })
      expect(response.status()).toBe(403)
    }
  }

  for (const origin of origins) {
    const headers: Record<string, string> = {}
    if (origin) headers.Origin = origin
    const tokenEnd = await page.request.post(`${apiOrigin}/api/v1/impersonation/end`, { headers, data: {} })
    results.push({ method: 'POST', path: '/api/v1/impersonation/end', origin: origin ?? 'missing', status: tokenEnd.status() })
    expect(tokenEnd.status()).toBe(403)
  }

  const orderPreflightResults = []
  for (const origin of origins) {
    const headers: Record<string, string> = {
      'Access-Control-Request-Method': 'POST',
      'Access-Control-Request-Headers': 'content-type,x-tenant-id,x-restaurant-id',
    }
    if (origin) headers.Origin = origin
    const response = await page.request.fetch(`${apiOrigin}/api/v1/public/orders`, { method: 'OPTIONS', headers })
    const expectedStatus = origin ? 403 : 200
    expect(response.status()).toBe(expectedStatus)
    orderPreflightResults.push({ origin: origin ?? 'missing', status: response.status() })
  }
  const storefrontOrigin = 'http://phase2a-restaurant-a.localhost:13113'
  const storefrontOrderPreflight = await page.request.fetch(`${apiOrigin}/api/v1/public/orders`, {
    method: 'OPTIONS',
    headers: {
      Origin: storefrontOrigin,
      'Access-Control-Request-Method': 'POST',
      'Access-Control-Request-Headers': 'content-type,x-tenant-id,x-restaurant-id,x-order-token',
    },
  })
  expect(storefrontOrderPreflight.status()).toBe(200)
  expect(storefrontOrderPreflight.headers()['access-control-allow-origin']).toBe(storefrontOrigin)
  orderPreflightResults.push({ origin: storefrontOrigin, status: storefrontOrderPreflight.status() })

  const trustedRefresh = await page.request.post(`${apiOrigin}/api/v1/auth/refresh`, {
    headers: { Cookie: staffCookieHeader, Origin: landingUrl },
  })
  expect(trustedRefresh.status()).toBe(204)
  expect(trustedRefresh.headers()['set-cookie']).toContain('poos_session=')
  expect(trustedRefresh.headers()['cache-control']).toContain('no-store')

  const currentCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(currentCookie).toBeDefined()
  const trustedLogout = await page.request.post(`${apiOrigin}/api/v1/auth/logout`, {
    headers: { Cookie: `poos_session=${currentCookie!.value}`, Origin: landingUrl },
  })
  expect(trustedLogout.status()).toBe(204)
  await expect(page.context().cookies(apiOrigin)).resolves.not.toContainEqual(expect.objectContaining({ name: 'poos_session' }))

  await test.info().attach('unsafe-route-csrf-matrix.json', {
    body: JSON.stringify({ results, orderPreflightResults, trustedRefreshStatus: trustedRefresh.status(), trustedLogoutStatus: trustedLogout.status() }, null, 2),
    contentType: 'application/json',
  })
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('rejects missing and untrusted origins on every registered staff and Super Admin unsafe route', async ({ page }, testInfo) => {
  test.skip(!composeProject, 'requires the isolated GitHub Actions authentication Compose project')
  test.setTimeout(180_000)
  expect(await signIn(page, ownerPassword)).toBe(200)
  const sessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie).toBeDefined()

  const workspace = process.env.GITHUB_WORKSPACE
  if (!workspace) throw new Error('The route-manifest check requires the CI workspace root')
  const manifest = await readFile(path.join(workspace, 'backend/cmd/api/routes.manifest'), 'utf8')
  const routePatterns = [...new Set(manifest.split(/\r?\n/).map((line) => line.trim()).filter(Boolean))]
  const exclusions = [
    { prefix: '/api/v1/public/orders', reason: 'anonymous customer capability routes are covered by storefront E2E' },
    { prefix: '/api/v1/translations/', reason: 'translation upload uses bearer authorization, not browser session cookies' },
    { prefix: '/api/v1/mock/paymob/', reason: 'mock payment endpoint is unauthenticated and has no ambient session authority' },
    { prefix: '/api/v1/webhook/paymob', reason: 'payment webhook is server-to-server and uses webhook verification' },
    { prefix: '/api/v1/auth/forgot-password', reason: 'password reset is a public endpoint with no authenticated session authority' },
  ]
  const unsafeRoutes = routePatterns.flatMap((pattern) => {
    const match = /^(POST|PUT|PATCH|DELETE) (\S+)$/.exec(pattern)
    if (!match || exclusions.some(({ prefix }) => match[2].startsWith(prefix))) return []
    return [{ method: match[1], pattern: match[2] }]
  })
  expect(unsafeRoutes.length).toBeGreaterThanOrEqual(140)

  const superAdminLogin = await page.request.post(`${apiOrigin}/api/v1/super-admin/auth/login`, {
    headers: { Origin: landingUrl },
    data: { email: 'phase2a-superadmin@example.test', password: superAdminPassword },
  })
  expect(superAdminLogin.status()).toBe(200)
  const superAdminCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_super_admin_session')
  expect(superAdminCookie?.httpOnly).toBe(true)
  expect(superAdminCookie?.expires).toBe(-1)
  const superAdminLoginBody = await superAdminLogin.json()
  const superAdminExpiry = Date.parse(superAdminLoginBody.data.expires_at)
  expect(superAdminExpiry).toBeGreaterThan(Date.now())
  expect(superAdminExpiry).toBeLessThanOrEqual(Date.now() + 8 * 60 * 60 * 1000)
  const staffSessionWithSuperAdminCookie = await page.request.get(`${apiOrigin}/api/v1/auth/session`, {
    headers: { Cookie: `poos_super_admin_session=${superAdminCookie!.value}` },
  })
  expect(staffSessionWithSuperAdminCookie.status()).toBe(401)
  const superAdminSessionWithStaffCookie = await page.request.get(`${apiOrigin}/api/v1/super-admin/auth/session`, {
    headers: { Cookie: `poos_session=${sessionCookie!.value}` },
  })
  expect(superAdminSessionWithStaffCookie.status()).toBe(401)
  const staffRoutes = unsafeRoutes.filter((route) => !route.pattern.startsWith('/api/v1/super-admin/'))
  const superAdminRoutes = unsafeRoutes.filter((route) => route.pattern.startsWith('/api/v1/super-admin/'))
  expect(staffRoutes.length).toBeGreaterThan(0)
  expect(superAdminRoutes.length).toBeGreaterThan(0)

  const results: Array<{ method: string; path: string; origin: string; status: number }> = []
  for (const route of unsafeRoutes) {
    const concretePath = route.pattern.replace(/\{([^/{}]+)\}/g, (_token, name: string) => {
      return name.toLowerCase() === 'email' ? 'csrf-probe%40example.test' : '1'
    })
    for (const origin of [undefined, 'https://attacker.example']) {
      const routeCookie = route.pattern.startsWith('/api/v1/super-admin/')
        ? superAdminCookie
        : sessionCookie
      expect(routeCookie).toBeDefined()
      const headers: Record<string, string> = {
        Cookie: `${route.pattern.startsWith('/api/v1/super-admin/') ? 'poos_super_admin_session' : 'poos_session'}=${routeCookie!.value}`,
        'Content-Type': 'application/json',
      }
      if (origin) headers.Origin = origin
      const response = await page.request.fetch(`${apiOrigin}${concretePath}`, {
        method: route.method,
        headers,
        ...(route.method === 'POST' || route.method === 'PUT' || route.method === 'PATCH' ? { data: {} } : {}),
      })
      results.push({ method: route.method, path: concretePath, origin: origin ?? 'missing', status: response.status() })
      expect(response.status(), `${route.method} ${route.pattern} with ${origin ?? 'missing Origin'}`).toBe(403)
    }
  }

  await testInfo.attach('registered-staff-super-admin-unsafe-route-csrf-matrix.json', {
    body: JSON.stringify({
      sourceManifest: 'backend/cmd/api/routes.manifest',
      routePatternsChecked: unsafeRoutes.length,
      staffRoutesChecked: staffRoutes.length,
      superAdminRoutesChecked: superAdminRoutes.length,
      staffSessionWithSuperAdminCookieStatus: staffSessionWithSuperAdminCookie.status(),
      superAdminSessionWithStaffCookieStatus: superAdminSessionWithStaffCookie.status(),
      excludedRoutes: exclusions,
      results,
    }, null, 2),
    contentType: 'application/json',
  })
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
    tenants: Array<{ tenant_id: number; tenant_name: string; default_restaurant_id?: number }>
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
    const state = rawAuthState ? JSON.parse(rawAuthState) as { state?: { user?: { tenant_id?: number; restaurant_id?: number } } } : null
    return {
      tenantID: state?.state?.user?.tenant_id ?? null,
      restaurantID: state?.state?.user?.restaurant_id ?? null,
    }
  })
  expect(selectedTenantID).toEqual({
    tenantID: selectedTenant.tenant_id,
    restaurantID: selectedTenant.default_restaurant_id,
  })
  expect(new URL(page.url()).search).toBe('')
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('cross-portal account changes clear cached categories in an open Dashboard tab', async ({ page }, testInfo) => {
  test.setTimeout(120_000)
  const categoriesA = await signInAndLoadTenantCategories(page)
  await seedTenantAScopedBrowserState(page)

  // A logout through Landing clears the shared cookie but leaves this Dashboard tab mounted.
  const logout = await page.context().request.post(`${apiOrigin}/api/v1/auth/logout`, {
    headers: { origin: landingUrl },
  })
  expect(logout.status()).toBe(204)

  const otherTenantLoginPage = await loginOtherTenantInSecondTab(page)
  const { sessionResponse, categories: categoriesB } = await readCategoriesAfterIdentityRefresh(page)
  expect(categoriesB.length).toBeGreaterThan(0)
  expect(new Set(categoriesA.map((category) => category.id))).not.toEqual(
    new Set(categoriesB.map((category) => category.id)),
  )
  await expect.poll(() => page.evaluate(() => {
    const raw = localStorage.getItem('auth-storage')
    return raw ? JSON.parse(raw).state?.user?.email ?? null : null
  })).toBe(otherTenantOwnerEmail)
  const clearedTenantState = await readTenantStateCleanup(page)
  expect(clearedTenantState).toEqual({
    priorCustomerAbsent: true,
    priorItemAbsent: true,
    heldOrdersEmpty: true,
    tenantCatalogCacheAbsent: true,
    kdsSnapshotAbsent: true,
    priorPreferencesAbsent: true,
    themeStorageCleared: true,
  })

  await testInfo.attach('dashboard-cross-account-cache-isolation.json', {
    body: JSON.stringify({
      sourceSha: process.env.GITHUB_SHA ?? 'local-uncommitted-candidate',
      initialTenantCategoriesLoaded: categoriesA.length,
      landingLogoutStatus: logout.status(),
      secondTenantCategoriesLoaded: categoriesB.length,
      sessionRefreshStatus: sessionResponse.status(),
      categoryResponseChangedWithIdentity: true,
      priorCategoriesReused: false,
      tenantScopedBrowserStateCleared: clearedTenantState,
    }, null, 2),
    contentType: 'application/json',
  })
  await otherTenantLoginPage.close()
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
  await landingPage.bringToFront()
  const landingAccountButton = landingPage.getByRole('button', { name: new RegExp(ownerName) })
  await expect(landingAccountButton).toBeVisible({ timeout: 30_000 })
  await landingAccountButton.click()
  const accountPanel = landingPage.locator('#sv-account-panel')
  await expect(accountPanel.getByText('WS05 Tenant A Restaurant')).toBeVisible()
  expect(await accountPanel.innerText()).not.toMatch(/phone|user id|tenant id|session|token/i)
  await expect(accountPanel.getByRole('link', { name: /dashboard/i })).toHaveAttribute('href', `${dashboardUrl}/en/dashboard`)
  await expect(accountPanel.getByRole('link', { name: /settings/i })).toHaveAttribute('href', `${dashboardUrl}/en/dashboard/settings/account`)
  await captureScreenshot(landingPage, 'landing-authenticated-account-menu')
  await landingPage.keyboard.press('Escape')
  await expect(accountPanel).toBeHidden()
  await expect(landingAccountButton).toBeFocused()
  await landingAccountButton.press('Enter')
  await expect(accountPanel).toBeVisible()

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
  const accessibleTree = await mobileDrawer.ariaSnapshot()
  expect(accessibleTree).toContain('dialog')
  expect(accessibleTree).toContain('navigation')
  expect(accessibleTree).toContain('link')
  const mobileAccount = mobileDrawer.locator('.sv-mobile-account')
  await expect(mobileAccount.getByText('WS05 Tenant A Restaurant')).toBeVisible()
  await expect(mobileAccount.locator('a').nth(0)).toHaveAttribute('href', `${dashboardUrl}/ar/dashboard`)
  await expect(mobileAccount.locator('a').nth(1)).toHaveAttribute('href', `${dashboardUrl}/ar/dashboard/settings/account`)
  const signOutBounds = await mobileAccount.locator('.sv-mobile-account__signout').boundingBox()
  expect(signOutBounds?.height ?? 0).toBeGreaterThanOrEqual(44)
  const dialogCloseButton = mobileDrawer.locator('button').first()
  await expect(dialogCloseButton).toBeFocused()
  const mobileNavControls = mobileDrawer.locator('nav').locator('button:visible, a[href]:visible')
  await mobileNavControls.first().focus()
  await landingPage.keyboard.press('ArrowDown')
  await expect(mobileNavControls.nth(1)).toBeFocused()
  await landingPage.keyboard.press('ArrowUp')
  await expect(mobileNavControls.first()).toBeFocused()
  await dialogCloseButton.focus()
  await landingPage.keyboard.press('Tab')
  expect(await mobileDrawer.evaluate((element) => element.contains(document.activeElement))).toBe(true)
  await landingPage.keyboard.press('Shift+Tab')
  expect(await mobileDrawer.evaluate((element) => element.contains(document.activeElement))).toBe(true)
  await mobileDrawer.locator('a[href]:visible, button:not([disabled]):visible').last().focus()
  await landingPage.keyboard.press('Tab')
  await expect(dialogCloseButton).toBeFocused()
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
    const recoveryAlert = page.locator('main[role="alert"]')
    await expect(recoveryAlert).toBeVisible()
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
    await recoveryAlert.getByRole('button', { name: /try again/i }).click()
    await recoveryResponsePromise
    await expectSignedInDashboard(page)
    await captureScreenshot(page, 'dashboard-session-recovered')
  } finally {
    if (redisStopped) execFileSync('docker', ['start', redisContainerID], { stdio: 'ignore' })
  }
})

test('expired Redis session TTL recovers the browser to the canonical login page', async ({ page }, testInfo) => {
  test.skip(!composeProject, 'requires the isolated GitHub Actions authentication Compose project')
  test.setTimeout(150_000)
  expect(await signIn(page, ownerPassword)).toBe(200)
  await expectSignedInDashboard(page)

  const sessionCookie = (await page.context().cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  expect(sessionCookie).toBeDefined()
  const redisContainerIDs = execFileSync('docker', [
    'ps',
    '--filter', `label=com.docker.compose.project=${composeProject}`,
    '--filter', 'label=com.docker.compose.service=redis',
    '--filter', 'label=com.pos.disposable=true',
    '--filter', 'label=com.pos.workstream=WS05',
    '--format', '{{.ID}}',
  ], { encoding: 'utf8' }).trim().split(/\r?\n/).filter(Boolean)
  expect(redisContainerIDs).toHaveLength(1)
  const redisContainerID = redisContainerIDs[0]
  const labels = JSON.parse(execFileSync('docker', [
    'inspect', '--format', '{{json .Config.Labels}}', redisContainerID,
  ], { encoding: 'utf8' })) as Record<string, string>
  expect(labels['com.docker.compose.project']).toBe(composeProject)
  expect(labels['com.docker.compose.service']).toBe('redis')
  expect(labels['com.pos.disposable']).toBe('true')
  expect(labels['com.pos.workstream']).toBe('WS05')

  const sessionKey = `auth-session:v1:${createHash('sha256').update(sessionCookie!.value).digest('hex')}`
  const shortened = execFileSync('docker', [
    'exec', redisContainerID, 'redis-cli', 'PEXPIRE', sessionKey, '1',
  ], { encoding: 'utf8' }).trim()
  expect(shortened).toBe('1')
  await expect.poll(() => execFileSync('docker', [
    'exec', redisContainerID, 'redis-cli', 'EXISTS', sessionKey,
  ], { encoding: 'utf8' }).trim(), { timeout: 10_000 }).toBe('0')

  const expiredSessionResponse = page.waitForResponse((response) => {
    return new URL(response.url()).pathname === '/api/v1/auth/session' && response.status() === 401
  }, { timeout: 30_000 })
  await page.evaluate(() => window.dispatchEvent(new Event('focus')))
  const response = await expiredSessionResponse
  await expect.poll(() => new URL(page.url()).pathname, { timeout: 30_000 }).toMatch(/^\/(?:en|ar)\/auth\/login$/)
  expect(new URL(page.url()).origin).toBe(new URL(landingUrl).origin)
  await expect(page.locator('form')).toBeVisible()
  expect(response.headers()['cache-control']).toContain('no-store')
  const replay = await page.request.get(`${apiOrigin}/api/v1/auth/session`, {
    headers: { Cookie: `poos_session=${sessionCookie!.value}`, Origin: landingUrl },
  })
  expect(replay.status()).toBe(401)
  await testInfo.attach('expired-session-browser-recovery.json', {
    body: JSON.stringify({ redisTTLForcedToOneMillisecond: shortened === '1', expiredKeyAbsent: true, expiredSessionStatus: response.status(), browserReturnedTo: new URL(page.url()).origin + new URL(page.url()).pathname, oldCookieReplayStatus: replay.status() }, null, 2),
    contentType: 'application/json',
  })
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})

test('authenticated Landing bootstrap has no guest flash or hydration mismatch', async ({ page, context }, testInfo) => {
  test.setTimeout(150_000)
  // On PR #82 alone, the disposable API uses its in-memory limiter until PR
  // #83's Redis-backed shared limiter is integrated. Let the earlier synthetic
  // login requests age out of the 20/minute window, then use the guarded API
  // restart below to verify that an existing session survives with shared Redis.
  await new Promise((resolve) => setTimeout(resolve, 61_000))
  const redisContainerBefore = isolatedHealthyRedisContainerID()
  expect(await signIn(page, ownerPassword)).toBe(200)
  const originalSessionCookie = (await context.cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  if (!originalSessionCookie) throw new Error('The browser did not retain the synthetic staff session cookie')
  const sessionBeforeRestart = await page.request.get(`${apiOrigin}/api/v1/auth/session`, {
    headers: { Cookie: `poos_session=${originalSessionCookie.value}`, Origin: landingUrl },
  })
  if (sessionBeforeRestart.status() !== 200) throw new Error('The synthetic session was not valid before the API restart')

  await restartIsolatedAPIForFreshLoginWindow()
  const redisContainerAfter = isolatedHealthyRedisContainerID()
  if (redisContainerAfter !== redisContainerBefore) throw new Error('The API restart replaced the shared Redis fixture')
  const retainedSessionCookie = (await context.cookies(apiOrigin)).find((cookie) => cookie.name === 'poos_session')
  if (!retainedSessionCookie || retainedSessionCookie.value !== originalSessionCookie.value) {
    throw new Error('The browser session cookie changed during the API restart')
  }
  const sessionAfterRestart = await page.request.get(`${apiOrigin}/api/v1/auth/session`, {
    headers: { Cookie: `poos_session=${retainedSessionCookie.value}`, Origin: landingUrl },
  })
  if (sessionAfterRestart.status() !== 200) throw new Error('The same session did not resolve after the API restart')
  await testInfo.attach('redis-session-survives-api-restart.json', {
    body: JSON.stringify({
      sourceSha: process.env.GITHUB_SHA ?? 'local-uncommitted-candidate',
      environment: 'disposable WS05 API and Redis services',
      sessionBeforeRestartStatus: sessionBeforeRestart.status(),
      sessionAfterRestartStatus: sessionAfterRestart.status(),
      sameRedisContainerRemainedHealthy: true,
      sameBrowserCookieRetained: true,
      cookieValueCaptured: false,
    }, null, 2),
    contentType: 'application/json',
  })
  const oldSessionLogout = await page.request.post(`${apiOrigin}/api/v1/auth/logout`, {
    headers: { Cookie: `poos_session=${retainedSessionCookie.value}`, Origin: landingUrl },
  })
  if (oldSessionLogout.status() !== 204) throw new Error('The temporary restart-continuity session could not be revoked')
  await context.clearCookies()

  expect(await signIn(page, ownerPassword)).toBe(200)
  await expectSignedInDashboard(page)

  const landingPage = await context.newPage()
  startAuthTimeline(landingPage)
  const hydrationErrors: string[] = []
  landingPage.on('console', (message) => {
    if (message.type() !== 'error') return
    if (/hydration|did not match|server html|text content does not match/i.test(message.text())) hydrationErrors.push(message.text())
  })
  landingPage.on('pageerror', (error) => {
    if (/hydration|did not match|server html|text content does not match/i.test(error.message)) hydrationErrors.push(error.message)
  })
  await routeContainerApiToCandidate(landingPage)
  await landingPage.route('**/api/v1/auth/session', async (route) => {
    await new Promise((resolve) => setTimeout(resolve, 1_000))
    await route.continue()
  })
  const htmlResponse = await landingPage.goto(`${landingUrl}/en/`, { waitUntil: 'domcontentloaded' })
  expect(htmlResponse?.status()).toBe(200)
  const sharedHTML = await htmlResponse?.text() ?? ''
  expect(sharedHTML).not.toContain(ownerEmail)
  expect(sharedHTML).not.toContain(ownerName)

  const accountActions = landingPage.locator('.sv-desktop-account-actions')
  const bootstrappingStatus = accountActions.getByRole('status')
  await expect(bootstrappingStatus).toHaveAttribute('aria-busy', 'true')
  await expect(accountActions.getByRole('link', { name: /log in/i })).toHaveCount(0)
  const signedInButton = landingPage.getByRole('button', { name: new RegExp(ownerName) })
  await expect(signedInButton).toBeVisible({ timeout: 30_000 })
  expect(hydrationErrors).toEqual([])
  await testInfo.attach('landing-bootstrap-sanitized.json', {
    body: JSON.stringify({ serverHTMLContainsProfile: sharedHTML.includes(ownerEmail) || sharedHTML.includes(ownerName), guestLoginVisibleDuringSessionCheck: false, accountButtonAfterSessionRestore: await signedInButton.count(), hydrationMismatchCount: hydrationErrors.length }, null, 2),
    contentType: 'application/json',
  })
  await landingPage.close()
  await page.unrouteAll({ behavior: 'ignoreErrors' })
})
