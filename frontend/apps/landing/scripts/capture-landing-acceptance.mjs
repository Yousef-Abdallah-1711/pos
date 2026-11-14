import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { createRequire } from 'node:module'
const { chromium } = createRequire(import.meta.url)('playwright')

const baseURL = (process.env.POOS_ACCEPTANCE_URL ?? 'http://localhost:38691').replace(/\/$/, '')
const evidenceRoot = path.resolve(process.env.POOS_ACCEPTANCE_EVIDENCE ?? '../../../audit-reports/poos-final-acceptance-2026-11-11/evidence')
const viewports = [320, 375, 390, 430, 768, 820, 1024, 1280, 1440, 1920]
const routes = [
  { name: 'landing', path: (locale) => `/${locale}` },
  { name: 'login', path: (locale) => `/${locale}/auth/login` },
  { name: 'register', path: (locale) => `/${locale}/auth/register` },
]
const results = []
const browserErrors = []
await mkdir(evidenceRoot, { recursive: true })

const browser = await chromium.launch({ headless: true })
try {
  const page = await browser.newPage({ deviceScaleFactor: 1, ignoreHTTPSErrors: process.env.SERVIO_DISPOSABLE_BROWSER === '1' })
  page.on('console', (message) => { if (message.type() === 'error') browserErrors.push({ type: 'console', text: message.text() }) })
  page.on('pageerror', (error) => browserErrors.push({ type: 'pageerror', text: error.message }))
  await page.addInitScript(() => {
    const state = { cls: 0, lcp: 0 }
    Object.defineProperty(window, '__poosPerformance', { value: state, configurable: true })
    try {
      new PerformanceObserver((entries) => {
        for (const entry of entries.getEntries()) if (!entry.hadRecentInput) state.cls += entry.value
      }).observe({ type: 'layout-shift', buffered: true })
      new PerformanceObserver((entries) => {
        for (const entry of entries.getEntries()) state.lcp = entry.startTime
      }).observe({ type: 'largest-contentful-paint', buffered: true })
    } catch {}
  })

  for (const locale of ['en', 'ar']) {
    for (const route of routes) {
      for (const width of viewports) {
        await page.setViewportSize({ width, height: 900 })
        const response = await page.goto(`${baseURL}${route.path(locale)}`, { waitUntil: 'networkidle', timeout: 45_000 })
        await page.evaluate(() => document.fonts.ready)
        await page.evaluate(async () => {
          for (const image of document.images) {
            if (image.closest('[hidden]')) continue
            image.scrollIntoView()
            try { await image.decode() } catch { /* Failed decoding is recorded below. */ }
          }
          window.scrollTo(0, 0)
          await new Promise(resolve => requestAnimationFrame(resolve))
        })
        await page.waitForTimeout(350)
        const checks = await page.evaluate(() => {
          const header = document.querySelector('header')
          const headerRect = header?.getBoundingClientRect()
          const authContent = document.querySelector('.auth-layout')
          const authRect = authContent?.getBoundingClientRect()
          const images = [...document.images].filter(image => !image.closest('[hidden]'))
          const navigation = performance.getEntriesByType('navigation')[0]
          const resources = performance.getEntriesByType('resource')
          const metrics = window.__poosPerformance ?? { cls: 0, lcp: 0 }
          return {
            viewportWidth: innerWidth,
            documentWidth: document.documentElement.scrollWidth,
            bodyWidth: document.body.scrollWidth,
            headerVisible: Boolean(headerRect && headerRect.width > 0 && headerRect.height > 0),
            authFormBelowHeader: !authRect || !headerRect || authRect.top >= headerRect.bottom - 1,
            direction: document.documentElement.dir,
            language: document.documentElement.lang,
            images: images.map((image) => ({ src: image.currentSrc, loaded: image.complete && image.naturalWidth > 0, width: image.naturalWidth, height: image.naturalHeight })),
            performance: {
              ttfbMs: Math.round(navigation?.responseStart - navigation?.requestStart || 0),
              fcpMs: Math.round(performance.getEntriesByName('first-contentful-paint')[0]?.startTime || 0),
              lcpMs: Math.round(metrics.lcp),
              cls: Number(metrics.cls.toFixed(4)),
              resourceTransferBytes: resources.reduce((total, entry) => total + entry.transferSize, 0),
            },
            title: document.title,
            canonical: document.querySelector('link[rel="canonical"]')?.href ?? null,
            alternates: [...document.querySelectorAll('link[rel="alternate"][hreflang]')].map((link) => ({ language: link.hreflang, href: link.href })),
          }
        })
        const filename = `${route.name}-${locale}-${width}.png`
        await page.screenshot({ path: path.join(evidenceRoot, filename), fullPage: true, animations: 'disabled' })
        const failedImages = checks.images.filter((image) => !image.loaded)
        const result = {
          route: route.path(locale),
          locale,
          width,
          status: response?.status() ?? null,
          overflow: checks.documentWidth > width || checks.bodyWidth > width,
          headerVisible: checks.headerVisible,
          authFormBelowHeader: checks.authFormBelowHeader,
          failedImages,
          screenshot: filename,
          ...checks,
        }
        results.push(result)
        if (result.status !== 200 || result.overflow || !result.headerVisible || !result.authFormBelowHeader || failedImages.length) {
          process.stderr.write(`VISUAL_CHECK_FAIL ${JSON.stringify(result)}\n`)
        }
      }
    }
  }

  for (const width of [320, 390, 768, 1024]) {
    await page.setViewportSize({ width, height: 900 })
    await page.goto(`${baseURL}/en`, { waitUntil: 'networkidle' })
    const menu = page.locator('button[aria-controls="sv-mobile-navigation"]')
    if (await menu.isVisible()) {
      await menu.click()
      const opened = await menu.getAttribute('aria-expanded')
      await page.keyboard.press('Escape')
      const closed = await menu.getAttribute('aria-expanded')
      results.push({ route: '/en', width, interaction: 'mobile menu open and Escape close', passed: opened === 'true' && closed === 'false' })
    } else {
      results.push({ route: '/en', width, interaction: 'mobile menu', passed: width >= 1024 })
    }
  }

  const endpointChecks = []
  for (const endpoint of ['/opengraph-image', '/robots.txt', '/sitemap.xml', '/en/route-that-does-not-exist']) {
    const response = await page.request.get(`${baseURL}${endpoint}`)
    endpointChecks.push({ path: endpoint, status: response.status(), contentType: response.headers()['content-type'] ?? null })
  }
  const output = { baseURL, capturedAt: new Date().toISOString(), viewports, routes: results, browserErrors, endpointChecks }
  await writeFile(path.join(evidenceRoot, 'acceptance-results.json'), `${JSON.stringify(output, null, 2)}\n`, 'utf8')
  const failures = results.filter((result) => result.passed === false || result.status && (result.status !== 200 || result.overflow || !result.headerVisible || !result.authFormBelowHeader || result.failedImages?.length || result.passed === false))
  if (failures.length || browserErrors.length || endpointChecks.some((item) => item.path === '/opengraph-image' ? item.status !== 200 || !item.contentType?.startsWith('image/png') : item.path === '/robots.txt' ? item.status !== 200 || !item.contentType?.startsWith('text/plain') : item.path === '/sitemap.xml' ? item.status !== 200 || !item.contentType?.includes('xml') : item.status !== 404)) {
    process.exitCode = 1
  }
  process.stdout.write(`${JSON.stringify({ pages: results.filter((item) => item.route && item.status).length, failures: failures.length, browserErrors: browserErrors.length, endpointChecks }, null, 2)}\n`)
} finally {
  await browser.close()
}
