/* Run against a preview or the deployed landing URL. Captures contain public synthetic data. */
import { createRequire } from 'node:module'
const { chromium } = createRequire(import.meta.url)('playwright')
import fs from 'node:fs/promises'
import path from 'node:path'
const base = process.env.SERVIO_MARKETING_URL
if (!base) throw new Error('Set SERVIO_MARKETING_URL to the intended preview or deployment URL')
const output = process.env.SERVIO_MARKETING_EVIDENCE || (process.env.SERVIO_BROWSER_EVIDENCE ? path.join(process.env.SERVIO_BROWSER_EVIDENCE, 'marketing') : '/tmp/servio-marketing-acceptance')
const widths = [360, 390, 768, 1024, 1440, 1920]
const routes = ['', '/features', '/features/pos', '/features/kitchen-display', '/features/online-ordering', '/features/website-builder', '/features/inventory', '/features/hr', '/features/analytics', '/how-it-works', '/pricing', '/demo', '/contact', '/privacy']
const checks = []
function check(name, condition, detail = '') {
 checks.push({ name, passed: !!condition, detail })
 if (!condition) throw new Error(`${name}: ${detail}`)
}
async function load(page, locale, route) {
 const response = await page.goto(`${base}/${locale}${route}`, { waitUntil: 'domcontentloaded', timeout: 60000 })
 check(`HTTP ${locale}${route}`, response.status() === 200, response.status())
 await page.locator('main h1').waitFor()
 check(`heading ${locale}${route}`, await page.locator('main h1').count() === 1)
 const overflow = await page.evaluate(() => document.documentElement.scrollWidth > innerWidth)
 check(`overflow ${locale}${route}`, !overflow)
}
async function capture(page, name) {
 await page.evaluate(async () => {
  for (const image of document.querySelectorAll('main img')) {
   if (image.closest('[hidden]')) continue
   image.scrollIntoView()
   await image.decode().catch(() => undefined)
  }
  window.scrollTo(0, 0)
 })
 await page.screenshot({ path: path.join(output, `${name}.png`), fullPage: page.viewportSize().width <= 1440 })
}
async function interactions(page, locale) {
 await load(page, locale, '/pricing')
 await page.locator('input[value=annual]').check()
 const href = await page.locator('main a.sv-button').getAttribute('href')
 check(`annual preference ${locale}`, href.endsWith('?billing=annual'), href)
 await page.locator('main a.sv-button').click()
 await page.waitForURL('**/auth/register?billing=annual')
 await page.waitForFunction(() => /annual|سنوي/.test(document.querySelector('main')?.textContent || ''))
 check(`registration annual ${locale}`, (await page.locator('main').textContent()).includes(locale === 'ar' ? 'سنوي' : 'annual'))
 await load(page, locale, '/demo')
 const tabs = page.getByRole('tab')
 await tabs.first().focus()
 await page.keyboard.press(locale === 'ar' ? 'ArrowLeft' : 'ArrowRight')
 check(`keyboard tabs ${locale}`, await tabs.nth(1).getAttribute('aria-selected') === 'true')
 await page.setViewportSize({ width: 390, height: 844 })
 const menu = page.locator('.sv-menu-button')
 await menu.click()
 check(`mobile menu ${locale}`, await menu.getAttribute('aria-expanded') === 'true')
 await page.keyboard.press('Escape')
 check(`menu escape ${locale}`, await menu.getAttribute('aria-expanded') === 'false')
}
async function main() {
 await fs.mkdir(output, { recursive: true })
 const browser = await chromium.launch({ headless: true, ...(process.env.SERVIO_BROWSER_EXECUTABLE ? { executablePath: process.env.SERVIO_BROWSER_EXECUTABLE } : {}), args: ['--no-sandbox'] })
 try {
  let page = await browser.newPage({ ignoreHTTPSErrors: process.env.SERVIO_DISPOSABLE_BROWSER === '1', viewport: { width: 1440, height: 1000 }, reducedMotion: 'reduce' })
  const errors = []
  page.on('pageerror', error => errors.push(error.message))
  for (const locale of ['en', 'ar']) {
   for (const width of widths) {
    await page.setViewportSize({ width, height: 1000 })
    await load(page, locale, '')
    check(`sections ${locale}/${width}`, await page.locator('main section').count() === 14)
    await capture(page, `${locale}-${width}`)
    await page.close()
    page = await browser.newPage({ ignoreHTTPSErrors: process.env.SERVIO_DISPOSABLE_BROWSER === '1', viewport: { width, height: 1000 }, reducedMotion: 'reduce' })
    page.on('pageerror', error => errors.push(error.message))
   }
   for (const route of routes.slice(1)) await load(page, locale, route)
   await interactions(page, locale)
  }
  check('browser exceptions', errors.length === 0, errors.join('; '))
 } finally {
  await browser.close()
  await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ base, checks }, null, 2))
 }
}
main().catch(error => { console.error(error.message); process.exitCode = 1 })
