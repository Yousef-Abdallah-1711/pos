/* Run against a preview or the deployed landing URL. Captures contain public synthetic data. */
import { createRequire } from 'node:module'
const require = createRequire(import.meta.url)
const { chromium, firefox } = require('playwright')
const browserName = process.env.SERVIO_MARKETING_BROWSER || 'chromium'
if (!['chromium', 'firefox'].includes(browserName)) throw new Error('Unsupported marketing browser')
import fs from 'node:fs/promises'
import path from 'node:path'
const colorScheme = process.env.SERVIO_COLOR_SCHEME || 'light'
if (!['light', 'dark'].includes(colorScheme)) throw new Error('Unsupported browser color scheme')
const expectedTheme = process.env.SERVIO_THEME || 'dark'
if (!['light', 'dark'].includes(expectedTheme)) throw new Error('Unsupported marketing theme')
const base = process.env.SERVIO_MARKETING_URL
if (!base) throw new Error('Set SERVIO_MARKETING_URL to the intended preview or deployment URL')
const output = process.env.SERVIO_MARKETING_EVIDENCE || (process.env.SERVIO_BROWSER_EVIDENCE ? path.join(process.env.SERVIO_BROWSER_EVIDENCE, 'marketing') : '/tmp/servio-marketing-acceptance')
const widths = [360, 390, 768, 1024, 1440, 1920]
const routes = ['', '/features', '/features/pos', '/features/kitchen-display', '/features/online-ordering', '/features/website-builder', '/features/inventory', '/features/hr', '/features/analytics', '/how-it-works', '/pricing', '/demo', '/contact', '/privacy']
const checks = []
const descriptions = new Map()
const navigationScreenshots = []
const sourceCommit = process.env.SERVIO_SOURCE_SHA || 'unknown'
function check(name, condition, detail = '') {
 checks.push({ name, passed: !!condition, detail })
 if (!condition) throw new Error(`${name}: ${detail}`)
}
async function load(page, locale, route) {
 let response
 let navigationError
 for (let attempt = 0; attempt < 2; attempt++) {
  try {
   response = await page.goto(`${base}/${locale}${route}`, { waitUntil: 'domcontentloaded', timeout: 60000 })
   navigationError = undefined
   break
  } catch (error) {
   navigationError = error
   if (attempt === 0) await page.waitForTimeout(250)
  }
 }
 if (navigationError) throw navigationError
 if (!response) throw new Error(`No response loading ${locale}${route}`)
 check(`HTTP ${locale}${route}`, response.status() === 200, response.status())
 const html = await response.text()
 const root = html.match(/<html[^>]*>/)?.[0] || ''
 check(`server locale ${locale}${route}`, root.includes(`lang="${locale}"`) && root.includes(`dir="${locale === 'ar' ? 'rtl' : 'ltr'}"`), root)
 // Next.js may briefly retain a hidden streamed tree during navigation. Wait
 // for one rendered heading; persistent duplicates must still time out/fail.
 await page.waitForFunction(() => {
  const headings = document.querySelectorAll('main h1')
  return headings.length === 1 && headings[0].getClientRects().length > 0
 })
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), expectedTheme)
 check(`resolved theme ${locale}${route}`, await page.evaluate(theme => document.documentElement.classList.contains(theme), expectedTheme))
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
 check(`shared header and footer ${locale}${route}`, await page.locator('.sv-desktop-links[aria-label]').count() === 1 && await page.locator('.sv-mobile-navigation nav[aria-label]').count() === 1 && await page.locator('.sv-footer nav[aria-label]').count() === 1)
 const overflow = await page.evaluate(() => document.documentElement.scrollWidth > innerWidth)
 check(`overflow ${locale}${route}`, !overflow)
 await page.addScriptTag({ path: require.resolve('axe-core/axe.min.js') })
 const violations = await page.evaluate(async () => {
  const result = await window.axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'] } })
  return result.violations.map(item => ({ id: item.id, impact: item.impact, nodes: item.nodes.map(node => ({ target: node.target, summary: node.failureSummary, contrast: node.any.map(check => check.data).filter(Boolean) })) }))
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
async function captureNavigationEvidence(browser, page, { locale, variant, viewport, frame }) {
 const directory = path.join(output, 'navigation')
 await fs.mkdir(directory, { recursive: true })
 const filename = `${locale}-${expectedTheme}-${variant}-${viewport.width}x${viewport.height}.png`
 const rawPath = path.join(directory, filename)
 await page.screenshot({ path: rawPath, fullPage: false, animations: 'disabled' })
 const evidence = {
  variant,
  application: 'public marketing website',
  route: new URL(page.url()).pathname,
  sourceEnvironment: new URL(base).origin,
  sourceCommit,
  userRole: 'anonymous public visitor',
  authenticationVerified: false,
  dataSanitized: true,
  locale,
  direction: locale === 'ar' ? 'rtl' : 'ltr',
  theme: expectedTheme,
  viewport,
  timestamp: new Date().toISOString(),
  screenshot: path.relative(output, rawPath),
  frame: null,
 }
 if (frame) {
  const image = (await fs.readFile(rawPath)).toString('base64')
  const label = `${locale.toUpperCase()} · ${evidence.direction.toUpperCase()} · ${expectedTheme.toUpperCase()} · ${variant.replaceAll('-', ' ')}`
  const outer = { width: viewport.width + 120, height: viewport.height + 200 }
  const framePage = await browser.newPage({ viewport: outer, deviceScaleFactor: 1 })
  const shape = frame === 'phone' ? 'phone' : frame === 'tablet' ? 'tablet' : 'desktop'
  await framePage.setContent(`<!doctype html><html lang="${locale}" dir="${evidence.direction}"><meta charset="utf-8"><style>
    *{box-sizing:border-box}html,body{margin:0;min-height:100%;background:#eeeaf3;color:#241d31;font-family:Arial,sans-serif}body{display:grid;place-items:center;padding:24px}.board{display:grid;justify-items:center;gap:12px}.label{font-size:12px;font-weight:700;letter-spacing:.04em;text-align:center}.device{padding:12px;background:#17131f;border:1px solid #625a6d;box-shadow:0 24px 60px #241d3130}.device.phone{border-width:11px;border-radius:42px}.device.tablet{border-width:12px;border-radius:25px}.device.desktop{border-width:10px;border-radius:14px}.chrome{height:22px;margin:-2px 0 8px;border-radius:5px;background:#302a39;color:#ded5e8;font-size:10px;padding:4px 8px}.screen{display:block;width:${viewport.width}px;height:${viewport.height}px;object-fit:fill;background:#fff}.source{font-size:10px;color:#625b70;text-align:center}
    </style><body><main class="board"><div class="label">${label}</div><div class="device ${shape}">${frame === 'desktop' ? '<div class="chrome">SERVIO · genuine browser capture</div>' : ''}<img class="screen" width="${viewport.width}" height="${viewport.height}" src="data:image/png;base64,${image}" alt="Unmodified SERVIO browser capture"></div><div class="source">Presentation frame only · raw: ${filename}</div></main></body></html>`)
  const frameName = filename.replace(/\.png$/, '-framed.png')
  const framePath = path.join(directory, frameName)
  await framePage.screenshot({ path: framePath, fullPage: true, animations: 'disabled' })
  await framePage.close()
  evidence.frame = path.relative(output, framePath)
 }
 navigationScreenshots.push(evidence)
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
async function interactions(browser, page, locale) {
 await page.setViewportSize({ width: 390, height: 844 })
 await page.goto(`${base}/${locale}`, { waitUntil: 'domcontentloaded' })
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), expectedTheme)
 const otherTheme = expectedTheme === 'dark' ? 'light' : 'dark'
 await page.locator('.sv-theme-toggle').click()
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), otherTheme)
 await page.reload({ waitUntil: 'domcontentloaded' })
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), otherTheme)
 await page.locator('.sv-theme-toggle').click()
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), expectedTheme)
 check(`theme toggle and persistence ${locale}`, await page.evaluate(theme => document.documentElement.classList.contains(theme), expectedTheme))
 await load(page, locale, '/pricing')
 await page.locator('input[value=annual]').check()
 const href = await page.locator('#pricing a.sv-button').getAttribute('href')
 check(`annual preference ${locale}`, href.endsWith('?billing=annual'), href)
 check(`annual charge ${locale}`, (await page.locator('#pricing .sv-price-value strong').first().textContent()) === (locale === 'ar' ? '٥٬٠٠٠' : '5,000'))
 check(`annual mobile overflow ${locale}`, await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth))
 await capture(page, `${locale}-pricing-annual-390`)
 await page.locator('#pricing a.sv-button').click()
 await page.waitForURL('**/auth/register?billing=annual')
 await page.waitForFunction(() => /annual|سنوي/i.test(document.querySelector('main')?.textContent || ''))
 const registrationText = await page.locator('main').textContent()
 check(`registration annual ${locale}`, /annual|سنوي/i.test(registrationText || ''))
 check(`registration theme persistence ${locale}`, await page.evaluate(theme => document.documentElement.classList.contains(theme), expectedTheme))
 if (!await page.evaluate(() => Boolean(window.axe))) await page.addScriptTag({ path: require.resolve('axe-core/axe.min.js') })
 const registrationContrast = await page.evaluate(async () => {
  const result = await window.axe.run(document, { runOnly: { type: 'rule', values: ['color-contrast'] } })
  return result.violations.map(item => ({ id: item.id, nodes: item.nodes.map(node => ({ target: node.target, summary: node.failureSummary })) }))
 })
 check(`registration contrast ${locale}`, registrationContrast.length === 0, JSON.stringify(registrationContrast))
 await capture(page, `register-annual-${locale}-390`)
 await load(page, locale, '')
 const faq = page.locator('#faq details').first()
 await faq.locator('summary').click()
 check(`FAQ opens ${locale}`, await faq.getAttribute('open') !== null)
 await page.locator('.sv-language').click()
 const other = locale === 'en' ? 'ar' : 'en'
 await page.waitForURL(`**/${other}`)
 await page.waitForFunction(direction => document.documentElement.dir === direction, other === 'ar' ? 'rtl' : 'ltr')
 check(`language switch ${locale}`, page.url().endsWith(`/${other}`))
 await load(page, locale, '/features/analytics')
 await page.locator('.sv-language').click()
 await page.waitForURL(`${base}/${other}/features/analytics`)
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), expectedTheme)
 check(`locale switch preserves equivalent feature route ${locale}`, page.url() === `${base}/${other}/features/analytics`)
 check(`locale switch preserves theme ${locale}`, await page.evaluate(theme => document.documentElement.classList.contains(theme), expectedTheme))
 await page.locator('.sv-language').click()
 await page.waitForURL(`${base}/${locale}/features/analytics`)

 await page.setViewportSize({ width: 1440, height: 1000 })
 await load(page, locale, '')
 await page.waitForLoadState('networkidle')
 const menuApiRequests = []
 const trackMenuApiRequest = request => { if (/\/api(?:\/|$)/.test(new URL(request.url()).pathname)) menuApiRequests.push(new URL(request.url()).pathname) }
 page.on('request', trackMenuApiRequest)
 const platform = page.locator('.sv-platform-trigger')
 await platform.click()
 await page.waitForTimeout(100)
 page.off('request', trackMenuApiRequest)
 check(`mega menu has no API dependency ${locale}`, menuApiRequests.length === 0, JSON.stringify(menuApiRequests))
 check(`desktop platform disclosure ${locale}`, await platform.getAttribute('aria-expanded') === 'true')
 const expectedFeatures = ['pos', 'kitchen-display', 'online-ordering', 'website-builder', 'inventory', 'hr', 'analytics'].map(slug => `${base}/${locale}/features/${slug}`)
 const renderedFeatureHrefs = await page.locator('.sv-mega-group a').evaluateAll(links => links.map(link => link.href))
 check(`mega menu feature route coverage ${locale}`, expectedFeatures.every(href => renderedFeatureHrefs.includes(href)) && renderedFeatureHrefs.length === expectedFeatures.length, JSON.stringify(renderedFeatureHrefs))
 check(`mega menu groups ${locale}`, await page.locator('.sv-mega-group').count() === 3)
 const desktopMenuViolations = await page.evaluate(async () => (await window.axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'] } })).violations.map(item => ({ id: item.id, nodes: item.nodes.map(node => node.target) })))
 check(`desktop mega menu accessibility ${locale}`, desktopMenuViolations.length === 0, JSON.stringify(desktopMenuViolations))
 await captureNavigationEvidence(browser, page, { locale, variant: 'desktop-mega-menu-open', viewport: { width: 1440, height: 1000 }, frame: 'desktop' })
 await page.keyboard.press('Escape')
 check(`desktop menu Escape and focus return ${locale}`, await platform.getAttribute('aria-expanded') === 'false' && await platform.evaluate(element => document.activeElement === element))
 await platform.focus()
 await page.keyboard.press('Enter')
 await page.keyboard.press('Tab')
 const focusedMenuItem = await page.locator('.sv-mega-group a').first().evaluate(element => element === document.activeElement && element.matches(':focus-visible'))
 check(`desktop menu keyboard focus ${locale}`, focusedMenuItem)
 await captureNavigationEvidence(browser, page, { locale, variant: 'desktop-mega-menu-keyboard-focused', viewport: { width: 1440, height: 1000 }, frame: 'desktop' })
 await page.keyboard.press('Escape')
 await page.mouse.move(8, page.viewportSize().height - 8)
 await platform.hover()
 check(`desktop menu pointer hover ${locale}`, await platform.getAttribute('aria-expanded') === 'true')
 await page.locator('.sv-main').hover({ position: { x: 8, y: 8 } })
 await page.waitForTimeout(120)
 check(`desktop menu closes after pointer leave ${locale}`, await platform.getAttribute('aria-expanded') === 'false')
 await platform.click()
 check(`desktop menu click open ${locale}`, await platform.getAttribute('aria-expanded') === 'true')
 await platform.click()
 check(`desktop menu click toggle ${locale}`, await platform.getAttribute('aria-expanded') === 'false')
 await platform.click()
 await page.locator('.sv-main').click({ position: { x: 8, y: 8 } })
 check(`desktop menu outside click ${locale}`, await platform.getAttribute('aria-expanded') === 'false')
 for (const slug of ['pos', 'kitchen-display', 'online-ordering', 'website-builder', 'inventory', 'hr', 'analytics']) {
  await load(page, locale, '')
  await page.locator('.sv-platform-trigger').click()
  const destination = `${base}/${locale}/features/${slug}`
  await Promise.all([page.waitForURL(destination), page.locator(`.sv-mega-group a[href="/${locale}/features/${slug}"]`).click()])
  await page.locator('main h1').waitFor()
  check(`mega menu destination ${locale}/${slug}`, page.url() === destination)
 }
 await load(page, locale, '')
 const footerTargets = await page.locator('.sv-footer nav a').evaluateAll(links => [...new Set(links.map(link => link.href))])
 check(`footer has unique verified destinations ${locale}`, footerTargets.length >= 15 && footerTargets.every(href => new URL(href).origin === new URL(base).origin))
 for (const destination of footerTargets) {
  const destinationPath = new URL(destination).pathname
  await load(page, locale, '')
  const route = destinationPath.replace(new RegExp(`^/${locale}`), '')
  const link = page.locator(`.sv-footer nav a[href="${destinationPath}"]`).first()
  await Promise.all([page.waitForURL(destination), link.click()])
  check(`footer destination ${locale}${route || '/'}`, page.url() === destination)
 }

 await page.setViewportSize({ width: 390, height: 844 })
 await load(page, locale, '')
 const menu = page.locator('.sv-menu-button')
 await menu.click()
 check(`mobile menu ${locale}`, await menu.getAttribute('aria-expanded') === 'true')
 const mobileControls = page.locator('.sv-mobile-navigation__controls button')
 await mobileControls.nth(1).click()
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), otherTheme)
 check(`mobile drawer theme control ${locale}`, await page.evaluate(theme => document.documentElement.classList.contains(theme), otherTheme))
 await mobileControls.nth(1).click()
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), expectedTheme)
 await page.locator('.sv-mobile-navigation__accordion').click()
 const mobileFeatureHrefs = await page.locator('.sv-mobile-feature-links a[href*="/features/"]').evaluateAll(links => links.map(link => link.href))
 check(`mobile feature accordion routes ${locale}`, expectedFeatures.every(href => mobileFeatureHrefs.includes(href)) && mobileFeatureHrefs.length === expectedFeatures.length)
 for (const slug of ['pos', 'kitchen-display', 'online-ordering', 'website-builder', 'inventory', 'hr', 'analytics']) {
  await load(page, locale, '')
  await page.locator('.sv-menu-button').click()
  await page.locator('.sv-mobile-navigation__accordion').click()
  const destination = `${base}/${locale}/features/${slug}`
  await Promise.all([page.waitForURL(destination), page.locator(`.sv-mobile-feature-links a[href="/${locale}/features/${slug}"]`).click()])
  await page.locator('main h1').waitFor()
  check(`mobile feature destination ${locale}/${slug}`, page.url() === destination)
 }
 await load(page, locale, '')
 await menu.click()
 await page.locator('.sv-mobile-navigation__accordion').click()
 check(`mobile navigation overflow ${locale}`, await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth))
 const mobileMenuViolations = await page.evaluate(async () => (await window.axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'] } })).violations.map(item => ({ id: item.id, nodes: item.nodes.map(node => node.target) })))
 check(`mobile navigation accessibility ${locale}`, mobileMenuViolations.length === 0, JSON.stringify(mobileMenuViolations))
 await captureNavigationEvidence(browser, page, { locale, variant: 'mobile-features-expanded', viewport: { width: 390, height: 844 }, frame: 'phone' })
 await page.keyboard.press('Escape')
 await page.waitForFunction(() => document.querySelector('.sv-menu-button')?.getAttribute('aria-expanded') === 'false' && document.activeElement === document.querySelector('.sv-menu-button'))
 check(`mobile menu Escape and focus return ${locale}`, await menu.getAttribute('aria-expanded') === 'false' && await menu.evaluate(element => document.activeElement === element))
 await menu.click()
 check(`mobile drawer locks body scroll ${locale}`, await page.evaluate(() => document.body.style.overflow === 'hidden'))
 await page.locator('.sv-mobile-backdrop').click({ position: { x: 8, y: 120 } })
 await page.waitForFunction(() => document.querySelector('.sv-menu-button')?.getAttribute('aria-expanded') === 'false' && document.activeElement === document.querySelector('.sv-menu-button'))
 check(`mobile backdrop closes and restores focus ${locale}`, await menu.getAttribute('aria-expanded') === 'false' && await menu.evaluate(element => document.activeElement === element) && await page.evaluate(() => document.body.style.overflow !== 'hidden'))
 await menu.click()
 await mobileControls.nth(0).click()
 const otherLocale = locale === 'en' ? 'ar' : 'en'
 await page.waitForURL(`${base}/${otherLocale}`)
 await page.waitForFunction(theme => document.documentElement.classList.contains(theme), expectedTheme)
 check(`mobile drawer locale control ${locale}`, page.url() === `${base}/${otherLocale}` && await page.evaluate(theme => document.documentElement.classList.contains(theme), expectedTheme))
 await page.locator('.sv-language').click()
 await page.waitForURL(`${base}/${locale}`)

 for (const [name, viewport] of [['tablet-portrait', { width: 768, height: 1024 }], ['tablet-landscape', { width: 1024, height: 768 }]]) {
  await page.setViewportSize(viewport)
  await load(page, locale, '')
  const tabletMenu = page.locator('.sv-menu-button')
  await tabletMenu.click()
  check(`${name} navigation opens ${locale}`, await tabletMenu.getAttribute('aria-expanded') === 'true' && await page.locator('.sv-mobile-navigation').isVisible())
  check(`${name} navigation fits ${locale}`, await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth))
  if (name === 'tablet-portrait') await captureNavigationEvidence(browser, page, { locale, variant: name, viewport, frame: 'tablet' })
  await page.keyboard.press('Escape')
 }

 for (const [name, viewport] of [['mobile', { width: 390, height: 844 }], ['desktop', { width: 1440, height: 1000 }]]) {
  await page.setViewportSize(viewport)
  await load(page, locale, '')
  const footerPath = path.join(output, 'navigation', `${locale}-${expectedTheme}-${name}-footer.png`)
  await fs.mkdir(path.dirname(footerPath), { recursive: true })
  await page.locator('.sv-footer').screenshot({ path: footerPath, animations: 'disabled' })
  navigationScreenshots.push({
   application: 'public marketing website', route: `/${locale}`, sourceEnvironment: new URL(base).origin,
   variant: `${name}-footer`,
   sourceCommit, userRole: 'anonymous public visitor', authenticationVerified: false, dataSanitized: true,
   locale, direction: locale === 'ar' ? 'rtl' : 'ltr', theme: expectedTheme, viewport,
   timestamp: new Date().toISOString(), screenshot: path.relative(output, footerPath), frame: null,
  })
 }

 await load(page, locale, '/demo')
 const tabs = page.getByRole('tab')
 await tabs.first().focus()
 await page.keyboard.press(locale === 'ar' ? 'ArrowLeft' : 'ArrowRight')
 check(`keyboard tabs ${locale}`, await tabs.nth(1).getAttribute('aria-selected') === 'true')
 for (const width of [320, 360, 375, 390, 430, 600, 768, 820, 1024, 1180, 1200, 1280, 1440, 1920]) {
  const height = width < 600 ? 844 : width === 1024 || width === 1180 ? 768 : width < 1200 ? 1024 : width === 1920 ? 1080 : 900
  await page.setViewportSize({ width, height })
  await load(page, locale, '')
  const desktop = width >= 1200
  const desktopLinks = page.locator('.sv-desktop-links')
  const mobileTrigger = page.locator('.sv-menu-button')
  check(`navigation breakpoint ${locale}/${width}`, desktop ? await desktopLinks.isVisible() && !await mobileTrigger.isVisible() : !await desktopLinks.isVisible() && await mobileTrigger.isVisible())
  if (desktop) {
   const trigger = page.locator('.sv-platform-trigger')
   await trigger.click()
   const bounds = await page.locator('.sv-mega-menu').boundingBox()
   check(`desktop mega menu viewport fit ${locale}/${width}`, !!bounds && bounds.x >= 0 && bounds.x + bounds.width <= width)
   await page.keyboard.press('Escape')
  } else {
   await mobileTrigger.click()
   check(`responsive drawer opens without overflow ${locale}/${width}`, await page.locator('.sv-mobile-navigation').isVisible() && await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth))
   await page.keyboard.press('Escape')
  }
 }
 await focusChecks(page, locale)
}
async function createPage(browser, viewport) {
 const page = await browser.newPage({ colorScheme, ignoreHTTPSErrors: process.env.SERVIO_DISPOSABLE_BROWSER === '1', viewport, reducedMotion: 'reduce' })
 if (process.env.SERVIO_THEME) await page.addInitScript(({ theme, origin }) => {
  if (location.origin === origin && !localStorage.getItem('servio-marketing-theme')) localStorage.setItem('servio-marketing-theme', theme)
 }, { theme: expectedTheme, origin: new URL(base).origin })
 // Routing disables HTTP cache, so each raw-SSR assertion reads a fresh response.
 // https://playwright.dev/docs/api/class-browsercontext#browser-context-route
 await page.route('**/*', route => route.continue())
 return page
}
async function captureComparisonBoards(browser) {
 for (const locale of ['en', 'ar']) {
  const variants = ['mobile-features-expanded', 'tablet-portrait', 'desktop-mega-menu-open']
  const sources = variants.map(variant => navigationScreenshots.find(item => item.locale === locale && item.variant === variant))
  if (sources.some(item => !item)) throw new Error(`Missing navigation screenshots for ${locale} comparison board`)
  const images = await Promise.all(sources.map(async item => (await fs.readFile(path.join(output, item.screenshot))).toString('base64')))
  const labels = ['Mobile 390 × 844', 'Tablet portrait 768 × 1024', 'Desktop mega menu 1440 × 1000']
  const cards = images.map((image, index) => `<figure><figcaption>${labels[index]}</figcaption><img src="data:image/png;base64,${image}" alt="${locale.toUpperCase()} ${labels[index]} genuine browser capture"><small>Raw source: ${sources[index].screenshot}</small></figure>`).join('')
  const boardPage = await browser.newPage({ viewport: { width: 2000, height: 1200 }, deviceScaleFactor: 1 })
  await boardPage.setContent(`<!doctype html><html lang="${locale}" dir="${locale === 'ar' ? 'rtl' : 'ltr'}"><meta charset="utf-8"><style>*{box-sizing:border-box}html,body{margin:0;background:#eeeaf3;color:#241d31;font-family:Arial,sans-serif}body{padding:28px}.board{display:grid;gap:22px}.heading{font-size:20px;font-weight:700}.cards{display:grid;grid-template-columns:.62fr .9fr 1.5fr;gap:18px;align-items:start}figure{display:grid;gap:10px;justify-items:center;margin:0;padding:14px;border-radius:14px;background:#fff;box-shadow:0 14px 36px #241d311a}figcaption{font-size:13px;font-weight:700;text-align:center}img{display:block;width:100%;height:800px;object-fit:contain;object-position:top center}small{font-size:10px;color:#625b70;text-align:center;overflow-wrap:anywhere}</style><body><main class="board"><div class="heading">SERVIO navigation · ${locale.toUpperCase()} · ${expectedTheme.toUpperCase()} · ${sourceCommit}</div><div class="cards">${cards}</div></main></body></html>`)
  const filename = `${locale}-${expectedTheme}-mobile-tablet-desktop-comparison-board.png`
  const boardPath = path.join(output, 'navigation', filename)
  await boardPage.screenshot({ path: boardPath, fullPage: true, animations: 'disabled' })
  await boardPage.close()
  navigationScreenshots.push({
   application: 'public marketing website', variant: 'mobile-tablet-desktop-comparison-board',
   route: `/${locale}`, sourceEnvironment: new URL(base).origin, sourceCommit,
   userRole: 'anonymous public visitor', authenticationVerified: false, dataSanitized: true,
   locale, direction: locale === 'ar' ? 'rtl' : 'ltr', theme: expectedTheme,
   viewport: { width: 2000, height: 1200 }, timestamp: new Date().toISOString(),
   screenshot: path.relative(output, boardPath), frame: null,
   sourceScreenshots: sources.map(item => item.screenshot),
  })
 }
}
async function main() {
 await fs.mkdir(output, { recursive: true })
 const browserType = browserName === 'firefox' ? firefox : chromium
 const browser = await browserType.launch({ headless: true, ...(process.env.SERVIO_BROWSER_EXECUTABLE ? { executablePath: process.env.SERVIO_BROWSER_EXECUTABLE } : {}), ...(browserName === 'chromium' ? { args: ['--no-sandbox'] } : {}) })
 try {
  let page = await createPage(browser, { width: 1440, height: 1000 })
  const errors = []
  const recoveredRscNavigations = []
  const networkFailures = []
  const observe = page => {
   page.on('pageerror', error => {
    const details = Object.fromEntries(Object.getOwnPropertyNames(error).map(key => [key, error[key]]))
    errors.push(`${page.url()}: ${JSON.stringify({ name: error.name, message: error.message, cause: error.cause, details })}`)
   })
   page.on('console', message => {
    if (message.type() !== 'error') return
    const text = message.text()
    const rscFallback = text.match(/^Failed to fetch RSC payload for (https?:\/\/\S+)\. Falling back to browser navigation\./)
    if (rscFallback) {
     recoveredRscNavigations.push({ pageUrl: page.url(), destination: rscFallback[1] })
     return
    }
    errors.push(`${page.url()}: ${text} @ ${JSON.stringify(message.location())}`)
   })
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
   await interactions(browser, page, locale)
  }
  await captureComparisonBoards(browser)
  const expectedPaths = new Set(['en', 'ar'].flatMap(locale => routes.map(route => `/${locale}${route}`)))
  const invalidFallbacks = recoveredRscNavigations.filter(({ pageUrl, destination }) => {
   try {
    const destinationUrl = new URL(destination)
    return !expectedPaths.has(destinationUrl.pathname) || new URL(pageUrl).pathname !== destinationUrl.pathname
   } catch {
    return true
   }
  })
  check('recovered Next.js navigation fallbacks use verified routes', invalidFallbacks.length === 0, JSON.stringify(invalidFallbacks))
  check('browser exceptions', errors.length === 0, errors.join('; '))
  check('failed network responses', networkFailures.length === 0, networkFailures.join('; '))
 } finally {
  await browser.close()
  await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ base, browser: browserName, color_scheme: colorScheme, theme: expectedTheme, checks }, null, 2))
  await fs.writeFile(path.join(output, 'navigation-evidence.json'), JSON.stringify({ sourceCommit, screenshots: navigationScreenshots }, null, 2))
 }
}
main().catch(error => { console.error(error.message); process.exitCode = 1 })
