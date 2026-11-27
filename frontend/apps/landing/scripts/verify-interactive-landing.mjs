import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { chromium, firefox, webkit } from 'playwright'

const baseURL = (process.env.LANDING_QA_URL ?? 'http://localhost:3001').replace(/\/$/, '')
const output = path.resolve(process.env.LANDING_QA_OUTPUT ?? '../../../docs/landing-interactive-experience/artifacts/acceptance')
const sizes = [{ name: 'desktop', width: 1440, height: 900 }, { name: 'tablet', width: 768, height: 900 }, { name: 'mobile', width: 390, height: 844 }]
const allCases = sizes.flatMap(size => ['en', 'ar'].flatMap(locale => ['light', 'dark'].map(theme => ({ ...size, locale, theme }))))
const cases = process.env.LANDING_QA_CASE ? allCases.filter(item => `${item.locale}-${item.theme}-${item.name}` === process.env.LANDING_QA_CASE) : allCases
const results = []
const smoke = []
await mkdir(output, { recursive: true })

async function loadAllSections(page) {
  await page.evaluate(async () => {
    for (const section of document.querySelectorAll('main section')) {
      window.scrollTo({ top: section.getBoundingClientRect().top + window.scrollY, left: 0, behavior: 'instant' })
      await new Promise(resolve => setTimeout(resolve, 110))
    }
    window.scrollTo({ top: 0, left: 0, behavior: 'instant' })
  })
  await page.waitForTimeout(200)
}

async function inspectPage(page) {
  return page.evaluate(() => { window.scrollTo({ left: 1000, top: window.scrollY, behavior: 'instant' }); const horizontalReach = Math.abs(window.scrollX); window.scrollTo({ left: 0, top: window.scrollY, behavior: 'instant' }); return ({
    lang: document.documentElement.lang,
    dir: document.documentElement.dir,
    theme: document.documentElement.classList.contains('dark') ? 'dark' : 'light',
    width: innerWidth,
    documentWidth: document.documentElement.scrollWidth,
    sectionCount: document.querySelectorAll('main section').length,
    headingCount: document.querySelectorAll('main h1').length,
    imagesBroken: [...document.images].filter(image => image.complete && image.naturalWidth === 0).map(image => image.currentSrc),
    pins: document.querySelectorAll('.pin-spacer').length,
    footerGroups: getComputedStyle(document.querySelector('.sv-footer-groups')).display,
    showcaseOverflow: getComputedStyle(document.querySelector('.cv-showcase-viewport')).overflowX,
    showcaseSectionOverflow: getComputedStyle(document.querySelector('.cv-showcase')).overflowX,
    horizontalScroll: window.scrollX,
    horizontalReach,
    bodyWidth: document.body.scrollWidth,
  }) })
}

async function checkInteraction(page, locale, width) {
  const checks = {}
  const journey = page.locator('.cv-story-step button')
  await journey.nth(4).click()
  await page.waitForFunction(() => document.querySelector('.cv-story-stage-head')?.textContent?.includes('05 / 07'), undefined, { timeout: 5_000 }).catch(async () => { await journey.nth(4).click() })
  checks.journey = await page.locator('.cv-story-stage-head').innerText()
  checks.journeySelected = checks.journey.includes('05 / 07')
  const demo = page.locator('.cv-workflow-demo--online .cv-demo-tabs button')
  await demo.nth(1).click()
  await page.waitForFunction(() => document.querySelector('.cv-workflow-demo--online .cv-demo-tabs button:nth-child(2)')?.getAttribute('aria-pressed') === 'true', undefined, { timeout: 5_000 }).catch(async () => { await demo.nth(1).click() })
  checks.online = await demo.nth(1).getAttribute('aria-pressed') === 'true'
  const ecosystemNode = page.locator('.cv-ecosystem-node').nth(2)
  await ecosystemNode.click()
  await page.waitForFunction(() => document.querySelectorAll('.cv-ecosystem-node')[2]?.getAttribute('aria-pressed') === 'true', undefined, { timeout: 5_000 }).catch(async () => { await ecosystemNode.click() })
  checks.ecosystem = await ecosystemNode.getAttribute('aria-pressed') === 'true'
  await page.locator('.cv-faq summary').first().click()
  checks.faq = await page.locator('.cv-faq details').first().getAttribute('open') !== null
  if (width < 901) {
    const menu = page.locator('button[aria-controls="sv-mobile-navigation"]')
    await menu.click()
    checks.menuOpened = await menu.getAttribute('aria-expanded') === 'true'
    await page.keyboard.press('Escape')
    checks.menuClosed = await menu.getAttribute('aria-expanded') === 'false'
  }
  const account = page.locator('main a[href*="/auth/register"]').first()
  checks.registrationHref = await account.getAttribute('href')
  checks.locale = locale
  return checks
}

const browser = await chromium.launch({ channel: 'chrome', headless: true })
try {
  for (const item of cases) {
    const context = await browser.newContext({ viewport: { width: item.width, height: item.height }, deviceScaleFactor: 1, reducedMotion: 'no-preference' })
    await context.addInitScript(theme => localStorage.setItem('servio-marketing-theme', theme), item.theme)
    const page = await context.newPage()
    const errors = []
    const unauthorized = []
    page.on('pageerror', error => errors.push(error.message))
    page.on('response', response => {
      if (response.status() === 401) unauthorized.push(response.url())
    })
    page.on('console', message => {
      if (message.type() === 'error' && !message.text().includes('ERR_CONNECTION_REFUSED') && !message.text().includes('status of 401 (Unauthorized)')) errors.push(message.text())
    })
    const response = await page.goto(`${baseURL}/${item.locale}`, { waitUntil: 'domcontentloaded', timeout: 60_000 })
    await page.locator('main h1').waitFor({ timeout: 60_000 })
    await page.evaluate(() => document.fonts.ready)
    await loadAllSections(page)
    await page.waitForTimeout(350)
    const checks = await inspectPage(page)
    const screenshot = `${item.locale}-${item.theme}-${item.name}.png`
    await page.screenshot({ path: path.join(output, screenshot), fullPage: true, animations: 'disabled', timeout: 90_000 })
    const interactions = item.name === 'mobile' || item.name === 'desktop' ? await checkInteraction(page, item.locale, item.width) : null
    const passed = response?.status() === 200 && checks.lang === item.locale && checks.theme === item.theme && checks.documentWidth <= item.width && checks.bodyWidth <= item.width && checks.horizontalReach === 0 && checks.sectionCount === 13 && checks.headingCount === 1 && !checks.imagesBroken.length && !errors.length && unauthorized.every(url => new URL(url).pathname === '/api/v1/auth/session') && (!interactions || Object.values(interactions).filter(value => typeof value === 'boolean').every(Boolean))
    const result = { ...item, status: response?.status(), ...checks, errors, unauthorized, interactions, screenshot, passed }
    results.push(result)
    process.stdout.write(`${passed ? 'PASS' : 'FAIL'} ${item.locale} ${item.theme} ${item.name}: width ${checks.documentWidth}/${item.width}, images ${checks.imagesBroken.length}, errors ${errors.length}\n`)
    await context.close()
  }
  for (const engine of [{ name: 'firefox', launcher: firefox }, { name: 'webkit', launcher: webkit }]) {
    try {
      const instance = await engine.launcher.launch({ headless: true })
      for (const locale of ['en', 'ar']) for (const theme of ['light', 'dark']) for (const width of [390, 1440]) {
        const context = await instance.newContext({ viewport: { width, height: 900 } })
        await context.addInitScript(value => localStorage.setItem('servio-marketing-theme', value), theme)
        const page = await context.newPage()
        const response = await page.goto(`${baseURL}/${locale}`, { waitUntil: 'domcontentloaded', timeout: 30_000 })
        await page.locator('main h1').waitFor({ timeout: 30_000 })
        await page.waitForTimeout(350)
        const checks = await inspectPage(page)
        const passed = response?.status() === 200 && checks.lang === locale && checks.theme === theme && checks.documentWidth <= width && checks.sectionCount === 13 && checks.headingCount === 1
        smoke.push({ engine: engine.name, available: true, locale, theme, width, status: response?.status(), documentWidth: checks.documentWidth, passed })
        process.stdout.write(`${passed ? 'PASS' : 'FAIL'} ${engine.name} ${locale} ${theme} ${width}px\n`)
        await context.close()
      }
      await instance.close()
    } catch (error) {
      smoke.push({ engine: engine.name, available: false, reason: String(error).split('\n')[0] })
    }
  }
  await writeFile(path.join(output, 'results.json'), `${JSON.stringify({ baseURL, timestamp: new Date().toISOString(), results, smoke }, null, 2)}\n`)
  if (results.some(result => !result.passed) || smoke.some(item => item.available && !item.passed)) process.exitCode = 1
} finally {
  await browser.close()
}
