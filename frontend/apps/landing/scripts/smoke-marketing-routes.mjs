import { writeFile } from 'node:fs/promises'
import path from 'node:path'
import { chromium } from 'playwright'

const baseURL = (process.env.LANDING_QA_URL ?? 'http://localhost:3001').replace(/\/$/, '')
const output = path.resolve('../../../docs/landing-interactive-experience/artifacts/route-smoke.json')
const featureSlugs = ['pos', 'kitchen-display', 'online-ordering', 'website-builder', 'inventory', 'hr', 'analytics']
const routes = ['features', ...featureSlugs.map(slug => `features/${slug}`), 'how-it-works', 'pricing', 'demo', 'contact', 'privacy', 'auth/login', 'auth/register']
const widths = [320, 375, 430, 820, 1024, 1280, 1920]
const results = []
const browser = await chromium.launch({ channel: 'chrome', headless: true })
try {
  const page = await browser.newPage({ viewport: { width: 390, height: 844 }, deviceScaleFactor: 1 })
  for (const locale of ['en', 'ar']) {
    await page.setViewportSize({ width: 390, height: 844 })
    for (const route of routes) {
      const response = await page.goto(`${baseURL}/${locale}/${route}`, { waitUntil: 'domcontentloaded', timeout: 45_000 })
      await page.locator('main h1').first().waitFor({ timeout: 45_000 })
      const state = await page.evaluate(() => ({ lang: document.documentElement.lang, dir: document.documentElement.dir, width: document.documentElement.scrollWidth, h1: document.querySelectorAll('main h1').length, header: Boolean(document.querySelector('header')), brokenImages: [...document.images].filter(image => image.complete && image.naturalWidth === 0).length }))
      const passed = response?.status() === 200 && state.lang === locale && state.dir === (locale === 'ar' ? 'rtl' : 'ltr') && state.width <= 390 && state.h1 === 1 && state.header && state.brokenImages === 0
      results.push({ kind: 'route', locale, route, status: response?.status(), ...state, passed })
      process.stdout.write(`${passed ? 'PASS' : 'FAIL'} /${locale}/${route}\n`)
    }
    for (const width of widths) {
      await page.setViewportSize({ width, height: 900 })
      const response = await page.goto(`${baseURL}/${locale}`, { waitUntil: 'domcontentloaded', timeout: 45_000 })
      await page.locator('main h1').waitFor({ timeout: 45_000 })
      const state = await page.evaluate(() => ({ documentWidth: document.documentElement.scrollWidth, bodyWidth: document.body.scrollWidth, menuVisible: getComputedStyle(document.querySelector('button[aria-controls="sv-mobile-navigation"]')).display !== 'none' }))
      const passed = response?.status() === 200 && state.documentWidth <= width && state.bodyWidth <= width
      results.push({ kind: 'width', locale, width, status: response?.status(), ...state, passed })
      process.stdout.write(`${passed ? 'PASS' : 'FAIL'} /${locale} ${width}px\n`)
    }
  }
  await writeFile(output, `${JSON.stringify({ baseURL, timestamp: new Date().toISOString(), results }, null, 2)}\n`)
  if (results.some(result => !result.passed)) process.exitCode = 1
} finally {
  await browser.close()
}
