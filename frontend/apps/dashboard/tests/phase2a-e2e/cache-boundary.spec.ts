import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { expect, test, type BrowserContext, type BrowserContextOptions, type TestInfo } from '@playwright/test'

const apiDirectUrl = 'https://localhost:18443/api/v1'
const apiCacheUrl = 'https://localhost:18453/api/v1'
const landingCacheUrl = 'https://localhost:13151'
const dashboardCacheUrl = 'https://localhost:13152'
const ownerAEmail = 'phase2a-owner-a@example.test'
const ownerAName = 'WS05 Tenant A Owner'
const ownerBEmail = 'phase2a-owner-b@example.test'
const ownerBName = 'WS05 Tenant B Owner'
const password = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'

type ResponseLike = {
  status(): number
  headers(): Record<string, string>
  text(): Promise<string>
}

type CacheEvidence = {
  principal: string
  route: string
  variant: string
  status: number
  cacheStatus: string
  cacheKey: string
  cacheControl: string
  vary: string
  setCookie: boolean
  containsOwnerAIdentity: boolean
  containsOwnerBIdentity: boolean
}

type HtmlRequest = {
  context: BrowserContext
  principal: string
  url: string
  route: string
  variant: string
  expectedCacheStatus?: string
}

const evidence: CacheEvidence[] = []

async function observe(response: ResponseLike, principal: string, route: string, variant: string) {
  const headers = response.headers()
  const body = await response.text()
  const entry: CacheEvidence = {
    principal,
    route,
    variant,
    status: response.status(),
    cacheStatus: headers['x-servio-cache-status'] ?? 'MISSING',
    cacheKey: headers['x-servio-cache-key'] ?? 'MISSING',
    cacheControl: headers['cache-control'] ?? '',
    vary: headers.vary ?? '',
    setCookie: Boolean(headers['set-cookie']),
    containsOwnerAIdentity: body.includes(ownerAEmail) || body.includes(ownerAName),
    containsOwnerBIdentity: body.includes(ownerBEmail) || body.includes(ownerBName),
  }
  evidence.push(entry)
  return { entry, body }
}

async function signIn(context: BrowserContext, email: string) {
  const login = () => context.request.post(`${apiDirectUrl}/auth/login`, {
    data: { email, password },
    headers: { origin: dashboardCacheUrl },
  })
  let response = await login()
  if (response.status() === 429) {
    const retryAfterSeconds = Number(response.headers()['retry-after'])
    expect(Number.isInteger(retryAfterSeconds) && retryAfterSeconds > 0 && retryAfterSeconds <= 60).toBe(true)
    await new Promise<void>((resolve) => setTimeout(resolve, (retryAfterSeconds + 1) * 1000))
    response = await login()
  }
  expect(response.status()).toBe(200)
  expect((await context.cookies('https://localhost')).some((cookie) => cookie.name === 'poos_session')).toBe(true)
}

async function requestMenu(context: BrowserContext, principal: string, restaurant: string, locale: 'en' | 'ar') {
  const response = await context.request.get(
    `${apiCacheUrl}/public/restaurants/phase2a-restaurant-${restaurant}/menu?lang=${locale}`,
    { headers: { 'accept-language': locale } },
  )
  const result = await observe(response, principal, 'public-menu', `${restaurant}-${locale}`)
  expect(response.status()).toBe(200)
  // The cache exposes only a SHA-256 identifier, never the raw URL or query data.
  expect(result.entry.cacheKey).toMatch(/^[a-f0-9]{64}$/)
  expect(result.entry.containsOwnerAIdentity || result.entry.containsOwnerBIdentity).toBe(false)
  return result.entry
}

async function requestSession(context: BrowserContext, principal: string) {
  const response = await context.request.get(`${apiCacheUrl}/auth/session`)
  const result = await observe(response, principal, 'auth-session', 'profile')
  expect(response.status()).toBe(200)
  expect(result.entry.cacheStatus).toBe('BYPASS')
  expect(result.entry.cacheControl.toLowerCase()).toContain('private')
  expect(result.entry.cacheControl.toLowerCase()).toContain('no-store')
  return result.body
}

async function requestHtml({ context, principal, url, route, variant, expectedCacheStatus }: HtmlRequest) {
  const page = await context.newPage()
  const response = await page.goto(url, { waitUntil: 'domcontentloaded' })
  if (!response) throw new Error(`Expected an HTTP response for ${route}`)
  const result = await observe(response, principal, route, variant)
  // Validate that evidence contains a digest rather than a URL or raw cache input.
  expect(result.entry.cacheKey).toMatch(/^[a-f0-9]{64}$/)
  if (expectedCacheStatus) expect(result.entry.cacheStatus).toBe(expectedCacheStatus)
  if (principal === 'anonymous') {
    expect(result.entry.containsOwnerAIdentity || result.entry.containsOwnerBIdentity).toBe(false)
  } else if (principal === 'user-a') {
    expect(result.entry.containsOwnerBIdentity).toBe(false)
  } else if (principal === 'user-b') {
    expect(result.entry.containsOwnerAIdentity).toBe(false)
  }
  await page.close()
  return result
}

async function requestProtectedDashboard(context: BrowserContext, principal: string) {
  const response = await context.request.get(`${dashboardCacheUrl}/en/dashboard`, { maxRedirects: 0 })
  const result = await observe(response, principal, 'dashboard-html', 'protected-route')
  expect(response.status()).toBeGreaterThanOrEqual(200)
  expect(response.status()).toBeLessThan(400)
  expect(result.entry.cacheStatus).toBe('BYPASS')
  // A key digest proves a cache decision exists without exposing the path input.
  expect(result.entry.cacheKey).toMatch(/^[a-f0-9]{64}$/)
  if (principal === 'user-a') expect(result.entry.containsOwnerBIdentity).toBe(false)
  if (principal === 'user-b') expect(result.entry.containsOwnerAIdentity).toBe(false)
}

async function warmLandingVariants(anonymous: BrowserContext) {
  const english = await requestHtml({ context: anonymous, principal: 'anonymous', url: `${landingCacheUrl}/en`, route: 'landing-html', variant: 'en' })
  const englishAgain = await requestHtml({ context: anonymous, principal: 'anonymous', url: `${landingCacheUrl}/en`, route: 'landing-html', variant: 'en-repeat' })
  expect(englishAgain.entry.cacheStatus).toBe(english.entry.cacheStatus === 'MISS' ? 'HIT' : 'BYPASS')

  const arabic = await requestHtml({ context: anonymous, principal: 'anonymous', url: `${landingCacheUrl}/ar`, route: 'landing-html', variant: 'ar' })
  const arabicAgain = await requestHtml({ context: anonymous, principal: 'anonymous', url: `${landingCacheUrl}/ar`, route: 'landing-html', variant: 'ar-repeat' })
  expect(arabicAgain.entry.cacheStatus).toBe(arabic.entry.cacheStatus === 'MISS' ? 'HIT' : 'BYPASS')
  expect(english.entry.cacheKey).not.toBe(arabic.entry.cacheKey)
  return english.entry.cacheStatus
}

async function verifyAuthenticatedCacheBypass(context: BrowserContext, principal: string) {
  expect((await requestMenu(context, principal, 'a', 'en')).cacheStatus).toBe('BYPASS')
  await requestHtml({ context, principal, url: `${landingCacheUrl}/en`, route: 'landing-html', variant: 'en', expectedCacheStatus: 'BYPASS' })
  await requestProtectedDashboard(context, principal)
  return requestSession(context, principal)
}

async function verifyPublicMenuVariants(anonymous: BrowserContext) {
  const english = await requestMenu(anonymous, 'anonymous', 'a', 'en')
  expect(english.cacheStatus).toBe('MISS')
  expect((await requestMenu(anonymous, 'anonymous', 'a', 'en')).cacheStatus).toBe('HIT')
  const arabic = await requestMenu(anonymous, 'anonymous', 'a', 'ar')
  expect(arabic.cacheStatus).toBe('MISS')
  expect((await requestMenu(anonymous, 'anonymous', 'a', 'ar')).cacheStatus).toBe('HIT')
  const restaurantB = await requestMenu(anonymous, 'anonymous', 'b', 'en')
  expect(restaurantB.cacheStatus).toBe('MISS')
  expect((await requestMenu(anonymous, 'anonymous', 'b', 'en')).cacheStatus).toBe('HIT')
  expect(new Set([english.cacheKey, arabic.cacheKey, restaurantB.cacheKey]).size).toBe(3)
}

async function verifySessionProfiles(ownerA: BrowserContext, ownerB: BrowserContext) {
  await signIn(ownerA, ownerAEmail)
  await signIn(ownerB, ownerBEmail)
  const sessionA = await verifyAuthenticatedCacheBypass(ownerA, 'user-a')
  const sessionB = await verifyAuthenticatedCacheBypass(ownerB, 'user-b')
  expect(sessionA).toContain(ownerAEmail)
  expect(sessionA).toContain(ownerAName)
  expect(sessionA).not.toContain(ownerBEmail)
  expect(sessionB).toContain(ownerBEmail)
  expect(sessionB).toContain(ownerBName)
  expect(sessionB).not.toContain(ownerAEmail)
}

async function verifyAnonymousAfterUsers(anonymous: BrowserContext, warmStatus: string) {
  const menu = await requestMenu(anonymous, 'anonymous', 'a', 'en')
  expect(menu.cacheStatus).toBe('HIT')
  const response = await requestHtml({ context: anonymous, principal: 'anonymous', url: `${landingCacheUrl}/en`, route: 'landing-html', variant: 'after-users' })
  expect(response.entry.cacheStatus).toBe(warmStatus === 'MISS' ? 'HIT' : 'BYPASS')
  expect(response.entry.containsOwnerAIdentity || response.entry.containsOwnerBIdentity).toBe(false)
}

async function saveEvidence(testInfo: TestInfo) {
  const report = {
    sourceSha: process.env.GITHUB_SHA ?? 'local-uncommitted-candidate',
    environment: 'isolated-disposable-services-with-local-cache-equivalent-proxies',
    limitation: 'This verifies application headers and cache-boundary behavior; it does not verify a live CDN configuration.',
    observations: evidence,
  }
  const serialized = JSON.stringify(report, null, 2)
  await testInfo.attach('cache-boundary-evidence.json', { body: serialized, contentType: 'application/json' })
  const directory = process.env.WS05_AUTH_EVIDENCE_DIR
  if (directory) {
    await mkdir(directory, { recursive: true })
    await writeFile(path.join(directory, 'cache-boundary-evidence.json'), serialized, { mode: 0o600 })
  }
}

test('shared cache isolates public variants from authenticated sessions and HTML', async ({ browser }, testInfo) => {
  test.setTimeout(120_000)
  evidence.length = 0
  const contextOptions: BrowserContextOptions = { ignoreHTTPSErrors: true, permissions: ['local-network-access'] }
  const anonymous = await browser.newContext(contextOptions)
  const ownerA = await browser.newContext(contextOptions)
  const ownerB = await browser.newContext(contextOptions)

  try {
    await verifyPublicMenuVariants(anonymous)
    const landingWarmStatus = await warmLandingVariants(anonymous)
    await verifySessionProfiles(ownerA, ownerB)
    await verifyAnonymousAfterUsers(anonymous, landingWarmStatus)
  } finally {
    await Promise.all([anonymous.close(), ownerA.close(), ownerB.close()])
    await saveEvidence(testInfo)
  }
})
