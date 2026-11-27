import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { chromium, firefox, webkit } from 'playwright'

const baseURL = process.env.SERVIO_QA_URL
if (!baseURL) throw new Error('Set SERVIO_QA_URL to the landing application under test')
const evidence = path.resolve(process.env.SERVIO_QA_EVIDENCE ?? '../../../docs/landing-interactive-experience/artifacts/operational-revision')
await mkdir(evidence, { recursive: true })

const profiles = [
  { width: 320, height: 568 }, { width: 360, height: 800 },
  { width: 375, height: 667 }, { width: 375, height: 812 },
  { width: 390, height: 844 }, { width: 393, height: 852 },
  { width: 412, height: 915 }, { width: 414, height: 896 },
  { width: 430, height: 932 }, { width: 768, height: 1024 },
  { width: 820, height: 1180 }, { width: 834, height: 1194 },
  { width: 1024, height: 1366 }, { width: 1024, height: 768 },
  { width: 1280, height: 720 }, { width: 1366, height: 768 },
  { width: 1440, height: 900 }, { width: 1536, height: 864 },
  { width: 1600, height: 900 }, { width: 1920, height: 1080 },
  { width: 2560, height: 1440 }, { width: 1280, height: 600 },
  { width: 1440, height: 650 }, { width: 1920, height: 700 },
  { width: 844, height: 390 }, { width: 1180, height: 820 },
]
const results = []
const reducedResults = []
const routeResults = []

for (const [engineName, engine] of Object.entries({ chromium, firefox, webkit }).filter(([name]) => !process.env.SERVIO_QA_SKIP_MATRIX && (!process.env.SERVIO_QA_ENGINES || process.env.SERVIO_QA_ENGINES.split(',').includes(name)))) {
  const executablePath = process.env[`SERVIO_QA_${engineName.toUpperCase()}_EXECUTABLE`]
  const browser = await engine.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
  try {
    const cases = engineName === 'chromium' && !process.env.SERVIO_QA_QUICK ? profiles : [profiles[4], profiles[16]]
    for (const locale of ['en', 'ar']) for (const theme of ['light', 'dark']) for (const profile of cases) {
      const context = await browser.newContext({ viewport: profile, reducedMotion: 'no-preference' })
      await context.addInitScript(value => localStorage.setItem('servio-marketing-theme', value), theme)
      const page = await context.newPage()
      const issues = []
      page.on('pageerror', error => issues.push(error.message))
      const response = await page.goto(`${baseURL}/${locale}`, { waitUntil: 'domcontentloaded', timeout: 45000 })
      await page.locator('.cv-bento-card').first().waitFor()
      await page.evaluate(() => document.fonts.ready)
      await page.waitForFunction(expected => document.querySelector('.cv-bento-card .sv-capture-image')?.dataset.sourceTheme === expected, theme, { timeout: 10000 })
      if (profile.width > 900 && profile.height >= 760) await page.waitForFunction(() => document.querySelectorAll('.pin-spacer').length === 1, undefined, { timeout: 10000 })
      await page.evaluate(async () => {
        const images = [...document.querySelectorAll('.cv-bento-card img')]
        for (const image of images) image.loading = 'eager'
        await Promise.allSettled(images.map(image => image.decode()))
      })
      const initial = await page.evaluate(() => {
        const grid = document.querySelector('.cv-bento-grid')
        const viewport = document.querySelector('.cv-showcase-viewport')
        return {
          documentWidth: document.documentElement.scrollWidth,
          viewportWidth: innerWidth,
          language: document.documentElement.lang,
          direction: document.documentElement.dir,
          theme: document.documentElement.classList.contains('dark') ? 'dark' : 'light',
          sections: document.querySelectorAll('.cv-page > section').length,
          cards: grid?.children.length,
          gridColumns: getComputedStyle(grid).gridTemplateColumns,
          captureLocale: grid?.querySelector('.sv-capture-image')?.dataset.sourceLocale,
          captureTheme: grid?.querySelector('.sv-capture-image')?.dataset.sourceTheme,
          bridge: document.querySelector('.cv-bento-bridge')?.getAttribute('href'),
          trackWidth: document.querySelector('.cv-showcase-track')?.scrollWidth,
          viewportClient: viewport?.clientWidth,
          missingBentoImages: [...grid.querySelectorAll('img')].filter(image => !image.naturalWidth).length,
          cardBounds: [...grid.children].map(card => {
            const rect = card.getBoundingClientRect()
            return { left: rect.left, right: rect.right, width: rect.width, height: rect.height, scrollHeight: card.scrollHeight, clientHeight: card.clientHeight }
          }),
        }
      })
      if (response?.status() !== 200) issues.push(`HTTP ${response?.status()}`)
      if (initial.documentWidth > initial.viewportWidth + 1) issues.push('horizontal page overflow')
      if (initial.sections !== 13 || initial.cards !== 3) issues.push('section or card count')
      if (initial.language !== locale || initial.direction !== (locale === 'ar' ? 'rtl' : 'ltr')) issues.push('locale/direction')
      if (initial.theme !== theme || initial.captureLocale !== locale || initial.captureTheme !== theme) issues.push('theme/capture matrix')
      if (initial.bridge !== '#how-orders-move') issues.push('order story bridge')
      if (initial.missingBentoImages || initial.cardBounds.some(card => card.width <= 0 || card.height <= 0 || card.left < -1 || card.right > initial.viewportWidth + 1 || card.scrollHeight > card.clientHeight + 2)) issues.push('bento image or card bounds')

      const tallDesktop = profile.width > 900 && profile.height >= 760
      let showcase = null
      if (tallDesktop) {
        await page.evaluate(() => {
          const stage = document.querySelector('.cv-showcase-stage')
          const top = stage.getBoundingClientRect().top + scrollY - 89
          scrollTo({ top, behavior: 'instant' })
        })
        await page.waitForTimeout(650)
        showcase = await page.evaluate(() => {
          const rect = selector => {
            const { top, right, bottom, left } = document.querySelector(selector).getBoundingClientRect()
            return { top, right, bottom, left }
          }
          return {
            header: rect('.sv-header'), stage: rect('.cv-showcase-stage'),
            slide: rect('.cv-feature-slide'), copy: rect('.cv-feature-slide-copy'),
            action: rect('.cv-feature-slide-copy a'), capture: rect('.cv-feature-slide-capture'),
            image: rect('.cv-feature-slide-capture img'), imageFrame: rect('.cv-feature-slide-capture .sv-capture-image'),
            pinSpacers: document.querySelectorAll('.pin-spacer').length,
          }
        })
        if (showcase.pinSpacers !== 1) issues.push('pin spacer count')
        if (showcase.stage.top < showcase.header.bottom - 2 || showcase.slide.bottom > profile.height + 2 || showcase.copy.top < showcase.slide.top - 2 || showcase.action.bottom > showcase.slide.bottom + 2 || showcase.capture.top < showcase.slide.top - 2 || showcase.capture.bottom > profile.height + 2) issues.push('pinned slide clipped')
        if (showcase.image.top < showcase.imageFrame.top - 2 || showcase.image.bottom > showcase.imageFrame.bottom + 2 || showcase.image.left < showcase.imageFrame.left - 2 || showcase.image.right > showcase.imageFrame.right + 2) issues.push('product image cropped')
        if (engineName === 'chromium' && profile.width === 1440) {
          await page.screenshot({ path: path.join(evidence, `${locale}-${theme}-showcase-first.png`) })
        }
        await page.evaluate(() => {
          const stage = document.querySelector('.cv-showcase-stage')
          const viewport = document.querySelector('.cv-showcase-viewport')
          const track = document.querySelector('.cv-showcase-track')
          const spacer = stage.parentElement
          const top = Number.parseFloat(getComputedStyle(stage).getPropertyValue('--cv-showcase-top'))
          const start = spacer.getBoundingClientRect().top + scrollY - top
          scrollTo({ top: start + track.scrollWidth - viewport.clientWidth - 2, behavior: 'instant' })
        })
        await page.waitForTimeout(700)
        const final = await page.evaluate(() => {
          const viewport = document.querySelector('.cv-showcase-viewport').getBoundingClientRect()
          const slide = document.querySelector('.cv-feature-slide:last-child').getBoundingClientRect()
          return {
            viewport: { left: viewport.left, right: viewport.right },
            slide: { left: slide.left, right: slide.right },
            progress: document.querySelector('.cv-showcase-controls span')?.textContent,
          }
        })
        showcase.final = final
        if (final.slide.left < final.viewport.left - 3 || final.slide.right > final.viewport.right + 3 || !final.progress?.includes('07 / 07')) issues.push('final slide or progress')
        if (engineName === 'chromium' && profile.width === 1440) await page.screenshot({ path: path.join(evidence, `${locale}-${theme}-showcase-last.png`) })
        if (profile.width === 1440 && profile.height === 900) {
          await page.evaluate(() => {
            const stage = document.querySelector('.cv-showcase-stage')
            const viewport = document.querySelector('.cv-showcase-viewport')
            const track = document.querySelector('.cv-showcase-track')
            const spacer = stage.parentElement
            const top = Number.parseFloat(getComputedStyle(stage).getPropertyValue('--cv-showcase-top'))
            const start = spacer.getBoundingClientRect().top + scrollY - top
            const first = track.children[0].getBoundingClientRect()
            const second = track.children[1].getBoundingClientRect()
            const offset = document.dir === 'rtl' ? first.right - second.right : second.left - first.left
            scrollTo({ top: start + Math.min(offset, track.scrollWidth - viewport.clientWidth), behavior: 'instant' })
          })
          await page.waitForTimeout(550)
          if (!(await page.locator('.cv-showcase-controls span').textContent())?.includes('02 / 07')) issues.push('second-slide progress')
          await page.evaluate(() => {
            const stage = document.querySelector('.cv-showcase-stage')
            const top = Number.parseFloat(getComputedStyle(stage).getPropertyValue('--cv-showcase-top'))
            const start = stage.parentElement.getBoundingClientRect().top + scrollY - top
            scrollTo({ top: start, behavior: 'instant' })
          })
          await page.waitForTimeout(550)
          if (!(await page.locator('.cv-showcase-controls span').textContent())?.includes('01 / 07')) issues.push('reverse-scroll progress')
        }
      } else {
        const native = await page.evaluate(() => ({ pinSpacers: document.querySelectorAll('.pin-spacer').length, overflow: getComputedStyle(document.querySelector('.cv-showcase-viewport')).overflowX }))
        showcase = native
        if (native.pinSpacers !== 0 || native.overflow !== 'auto') issues.push('native showcase fallback')
        if (profile.width === 390 && profile.height === 844) {
          const step = page.locator('.cv-bento-card--1 .cv-bento-story-controls button').nth(1)
          await step.click()
          await page.waitForFunction(() => document.querySelectorAll('.cv-bento-card--1 .cv-bento-story-controls button')[1]?.getAttribute('aria-pressed') === 'true')
          await page.locator('.cv-feature-slide:nth-child(2)').evaluate(slide => slide.scrollIntoView({ behavior: 'instant', block: 'nearest', inline: 'start' }))
          await page.waitForTimeout(1200)
          if (!(await page.locator('.cv-showcase-controls span').textContent())?.includes('02 / 07')) issues.push('native scroll progress')
        }
      }

      if (engineName === 'chromium' && profile.width === 390) {
        const beforeSample = await page.locator('.cv-bento-card--1 .cv-bento-sample strong').textContent()
        await page.locator('.cv-bento-story-controls button').nth(2).click()
        const pressed = await page.locator('.cv-bento-card--1 .cv-bento-story-controls button').nth(2).getAttribute('aria-pressed')
        if (pressed !== 'true') issues.push('bento step selection')
        if ((await page.locator('.cv-bento-card--1 .cv-bento-sample strong').textContent()) === beforeSample) issues.push('bento sample state')
        await page.locator('.cv-bento-card--1 .cv-bento-replay').click()
        const reset = await page.locator('.cv-bento-card--1 .cv-bento-story-controls button').first().getAttribute('aria-pressed')
        if (reset !== 'true') issues.push('bento replay')
        await page.locator('.cv-problem').screenshot({ path: path.join(evidence, `${locale}-${theme}-bento-mobile.png`) })
      }
      const result = { engine: engineName, locale, theme, ...profile, initial, showcase, issues }
      results.push(result)
      process.stdout.write(`${issues.length ? 'FAIL' : 'PASS'} ${engineName} ${locale} ${theme} ${profile.width}x${profile.height}${issues.length ? `: ${issues.join(', ')}` : ''}\n`)
      await context.close()
    }
  } finally {
    await browser.close()
  }
}

if (!process.env.SERVIO_QA_ENGINES || process.env.SERVIO_QA_ENGINES.split(',').includes('chromium')) {
  const executablePath = process.env.SERVIO_QA_CHROMIUM_EXECUTABLE
  const browser = await chromium.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
  try {
    for (const locale of ['en', 'ar']) for (const theme of ['light', 'dark']) {
      const context = await browser.newContext({ viewport: { width: 1440, height: 900 }, reducedMotion: 'reduce' })
      await context.addInitScript(value => localStorage.setItem('servio-marketing-theme', value), theme)
      const page = await context.newPage()
      await page.goto(`${baseURL}/${locale}`, { waitUntil: 'domcontentloaded' })
      await page.locator('.cv-showcase-track').waitFor()
      const state = await page.evaluate(() => ({
        pinSpacers: document.querySelectorAll('.pin-spacer').length,
        overflow: getComputedStyle(document.querySelector('.cv-showcase-viewport')).overflowX,
      }))
      await page.locator('.cv-feature-slide:last-child a').first().focus()
      await page.waitForTimeout(550)
      const focus = await page.evaluate(() => {
        const viewport = document.querySelector('.cv-showcase-viewport').getBoundingClientRect()
        const link = document.querySelector('.cv-feature-slide:last-child a').getBoundingClientRect()
        return { linkLeft: link.left, linkRight: link.right, viewportLeft: viewport.left, viewportRight: viewport.right }
      })
      const issues = []
      if (state.pinSpacers !== 0 || state.overflow !== 'auto') issues.push('reduced-motion fallback')
      if (focus.linkLeft < focus.viewportLeft - 2 || focus.linkRight > focus.viewportRight + 2) issues.push('keyboard focus reveal')
      reducedResults.push({ locale, theme, state, focus, issues })
      process.stdout.write(`${issues.length ? 'FAIL' : 'PASS'} reduced ${locale} ${theme}${issues.length ? `: ${issues.join(', ')}` : ''}\n`)
      await context.close()
    }

    for (const locale of ['en', 'ar']) {
      const context = await browser.newContext({ viewport: { width: 390, height: 844 } })
      const page = await context.newPage()
      await page.goto(`${baseURL}/${locale}`, { waitUntil: 'domcontentloaded' })
      const routes = await page.evaluate(locale => [...new Set([...document.querySelectorAll('a[href]')]
        .map(link => new URL(link.href).pathname)
        .filter(value => value === `/${locale}` || value.startsWith(`/${locale}/`)))], locale)
      for (const route of routes) {
        const response = await page.goto(`${baseURL}${route}`, { waitUntil: 'domcontentloaded' })
        const check = await page.evaluate(() => ({ width: document.documentElement.scrollWidth, viewport: innerWidth, heading: !!document.querySelector('main h1'), footer: !!document.querySelector('.sv-footer') }))
        const issues = []
        if (response?.status() !== 200 || check.width > check.viewport + 1 || !check.heading || !check.footer) issues.push('route response/layout')
        routeResults.push({ locale, route, status: response?.status(), ...check, issues })
        process.stdout.write(`${issues.length ? 'FAIL' : 'PASS'} route ${route}${issues.length ? `: ${issues.join(', ')}` : ''}\n`)
      }
      await context.close()
    }
  } finally {
    await browser.close()
  }
}

const failures = [...results, ...reducedResults, ...routeResults].filter(result => result.issues.length)
const output = { baseURL, capturedAt: new Date().toISOString(), results, reducedResults, routeResults, failures }
await writeFile(path.join(evidence, 'results.json'), `${JSON.stringify(output, null, 2)}\n`)
if (failures.length) process.exitCode = 1
