/* Run against a preview or the deployed landing URL. Captures contain public synthetic data. */
import { createRequire } from 'node:module'
const require = createRequire(import.meta.url)
const { chromium, firefox } = require('playwright')
const browserName = process.env.SERVIO_MARKETING_BROWSER || 'chromium'
if (!['chromium', 'firefox'].includes(browserName)) throw new Error('Unsupported marketing browser')
import fs from 'node:fs/promises'
import path from 'node:path'
const base = process.env.SERVIO_MARKETING_URL
if (!base) throw new Error('Set SERVIO_MARKETING_URL to the intended preview or deployment URL')
const output = process.env.SERVIO_MARKETING_EVIDENCE || (process.env.SERVIO_BROWSER_EVIDENCE ? path.join(process.env.SERVIO_BROWSER_EVIDENCE, 'marketing') : '/tmp/servio-marketing-acceptance')
const widths = [360, 390, 768, 1024, 1440, 1920]
const routes = ['', '/features', '/features/pos', '/features/kitchen-display', '/features/online-ordering', '/features/website-builder', '/features/inventory', '/features/hr', '/features/analytics', '/how-it-works', '/pricing', '/demo', '/contact', '/privacy']
const checks = []
const descriptions = new Map()
function check(name, condition, detail = '') {
 checks.push({ name, passed: !!condition, detail })
 if (!condition) throw new Error(`${name}: ${detail}`)
}
async function load(page, locale, route) {
 const response = await page.goto(`${base}/${locale}${route}`, { waitUntil: 'domcontentloaded', timeout: 60000 })
 check(`HTTP ${locale}${route}`, response.status() === 200, response.status())
 const html = await response.text()
 const root = html.match(/<html[^>]*>/)?.[0] || ''
 check(`server locale ${locale}${route}`, root.includes(`lang="${locale}"`) && root.includes(`dir="${locale === 'ar' ? 'rtl' : 'ltr'}"`), root)
 await page.locator('main h1').waitFor()
 const canonical = await page.locator('link[rel=canonical]').getAttribute('href')
 check(`canonical ${locale}${route}`, canonical === `${base}/${locale}${route}`, canonical)
 const description = await page.locator('meta[name="description"]').getAttribute('content')
 const duplicate = [...descriptions].find(([url, text]) => url !== canonical && text === description)
 check(`unique description ${locale}${route}`, !!description && !duplicate, duplicate?.[0] || '')
 descriptions.set(canonical, description)
 const alternates = await page.locator('link[rel=alternate][hreflang]').evaluateAll(links => links.map(link => ({ language: link.hreflang, href: link.href })))
 check(`hreflang ${locale}${route}`, ['en', 'ar'].every(language => alternates.some(link => link.language === language && link.href === `${base}/${language}${route}`)))
 const schema = JSON.parse(await page.locator('script[type="application/ld+json"]').textContent())
 check(`structured data ${locale}${route}`, schema['@type'] === 'WebSite' && schema.name === 'SERVIO')
 check(`heading ${locale}${route}`, await page.locator('main h1').count() === 1)
 const overflow = await page.evaluate(() => document.documentElement.scrollWidth > innerWidth)
 check(`overflow ${locale}${route}`, !overflow)
 await page.addScriptTag({ path: require.resolve('axe-core/axe.min.js') })
 const violations = await page.evaluate(async () => {
  const result = await window.axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'] } })
  return result.violations.map(item => ({ id: item.id, impact: item.impact, nodes: item.nodes.map(node => node.target) }))
 })
 check(`accessibility ${locale}${route}`, violations.length === 0, JSON.stringify(violations))
}
async function capture(page, name) {
 const missing = []
 const images = page.locator('main img')
 for (let index = 0; index < await images.count(); index++) {
  const image = images.nth(index)
  if (await image.evaluate(element => !!element.closest('[hidden]'))) continue
  await image.scrollIntoViewIfNeeded()
  const handle = await image.elementHandle()
  try {
   await page.waitForFunction(element => element.complete && element.naturalWidth > 0, handle, { timeout: 15000 })
  } catch (error) {
   missing.push({ ...await image.evaluate(element => ({ src: element.currentSrc || element.src, complete: element.complete, width: element.naturalWidth })), error: error.message.slice(0, 160) })
   break
  } finally { await handle.dispose() }
 }
 await page.evaluate(() => window.scrollTo(0, 0))
 await page.screenshot({ path: path.join(output, `${name}.png`), fullPage: page.viewportSize().width <= 1440 })
 check(`captures loaded ${name}`, missing.length === 0, JSON.stringify(missing))
}
async function focusChecks(page, locale) {
 await load(page, locale, '')
 await page.keyboard.press('Tab')
 for (const [name, selector] of [['light', '.sv-hero .sv-button'], ['dark', '#kitchen .sv-section-actions a'], ['final', '.sv-final .sv-button'], ['footer', '.sv-footer nav a']]) {
  const target = page.locator(selector).first()
  await target.focus()
  await target.scrollIntoViewIfNeeded()
  const result = await target.evaluate(element => {
   const luminance = color => {
    const channels = color.match(/[\d.]+/g).slice(0, 3).map(Number).map(value => value / 255).map(value => value <= .04045 ? value / 12.92 : ((value + .055) / 1.055) ** 2.4)
    return channels.reduce((total, value, index) => total + value * [.2126, .7152, .0722][index], 0)
   }
   let parent = element.parentElement
   while (parent && (getComputedStyle(parent).backgroundColor === 'rgba(0, 0, 0, 0)' || getComputedStyle(parent).backgroundColor === 'transparent')) parent = parent.parentElement
   const style = getComputedStyle(element)
   const levels = [luminance(style.outlineColor), luminance(parent ? getComputedStyle(parent).backgroundColor : 'rgb(255, 255, 255)')].sort((a, b) => a - b)
   const rect = element.getBoundingClientRect()
   const visible = document.elementFromPoint(rect.x + rect.width / 2, rect.y + rect.height / 2)
   return { contrast: (levels[1] + .05) / (levels[0] + .05), width: Number.parseFloat(style.outlineWidth), focusVisible: element.matches(':focus-visible'), unobscured: !!visible && element.contains(visible) }
  })
  check(`keyboard focus ${locale}/${name}`, result.focusVisible && result.width >= 2 && result.contrast >= 3 && result.unobscured, JSON.stringify(result))
 }
}
async function interactions(page, locale) {
 await page.setViewportSize({ width: 390, height: 844 })
 await load(page, locale, '/pricing')
 await page.locator('input[value=annual]').check()
 const href = await page.locator('#pricing a.sv-button').getAttribute('href')
 check(`annual preference ${locale}`, href.endsWith('?billing=annual'), href)
 check(`annual charge ${locale}`, (await page.locator('#pricing .sv-price-value strong').first().textContent()) === (locale === 'ar' ? '٥٬٠٠٠' : '5,000'))
 check(`annual mobile overflow ${locale}`, await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth))
 await capture(page, `${locale}-pricing-annual-390`)
 await page.locator('#pricing a.sv-button').click()
 await page.waitForURL('**/auth/register?billing=annual')
 await page.waitForFunction(() => /annual|سنوي/.test(document.querySelector('main')?.textContent || ''))
 check(`registration annual ${locale}`, (await page.locator('main').textContent()).includes(locale === 'ar' ? 'سنوي' : 'annual'))
 await load(page, locale, '')
 const faq = page.locator('#faq details').first()
 await faq.locator('summary').click()
 check(`FAQ opens ${locale}`, await faq.getAttribute('open') !== null)
 await page.locator('.sv-language').click()
 const other = locale === 'en' ? 'ar' : 'en'
 await page.waitForURL(`**/${other}`)
 await page.waitForFunction(direction => document.documentElement.dir === direction, other === 'ar' ? 'rtl' : 'ltr')
 check(`language switch ${locale}`, page.url().endsWith(`/${other}`))
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
 await focusChecks(page, locale)
}
async function createPage(browser, viewport) {
 const page = await browser.newPage({ ignoreHTTPSErrors: process.env.SERVIO_DISPOSABLE_BROWSER === '1', viewport, reducedMotion: 'reduce' })
 // Routing disables HTTP cache, so each raw-SSR assertion reads a fresh response.
 // https://playwright.dev/docs/api/class-browsercontext#browser-context-route
 await page.route('**/*', route => route.continue())
 return page
}
async function main() {
 await fs.mkdir(output, { recursive: true })
 const browserType = browserName === 'firefox' ? firefox : chromium
 const browser = await browserType.launch({ headless: true, ...(process.env.SERVIO_BROWSER_EXECUTABLE ? { executablePath: process.env.SERVIO_BROWSER_EXECUTABLE } : {}), ...(browserName === 'chromium' ? { args: ['--no-sandbox'] } : {}) })
 try {
  let page = await createPage(browser, { width: 1440, height: 1000 })
  const errors = []
  const networkFailures = []
  const observe = page => {
   page.on('pageerror', error => errors.push(error.message))
   page.on('console', message => { if (message.type() === 'error') errors.push(message.text()) })
   page.on('response', response => { if (response.status() >= 400) networkFailures.push(`${response.status()} ${new URL(response.url()).pathname}`) })
  }
  observe(page)
  for (const locale of ['en', 'ar']) {
   for (const width of widths) {
    await page.setViewportSize({ width, height: 1000 })
    await load(page, locale, '')
    check(`sections ${locale}/${width}`, await page.locator('main section').count() === 14)
    await capture(page, `${locale}-${width}`)
    for (const route of routes.slice(1)) await load(page, locale, route)
    await page.close()
    page = await createPage(browser, { width, height: 1000 })
    observe(page)
   }
   await interactions(page, locale)
  }
  check('browser exceptions', errors.length === 0, errors.join('; '))
  check('failed network responses', networkFailures.length === 0, networkFailures.join('; '))
 } finally {
  await browser.close()
  await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ base, browser: browserName, checks }, null, 2))
 }
}
main().catch(error => { console.error(error.message); process.exitCode = 1 })
