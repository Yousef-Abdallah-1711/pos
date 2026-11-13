import assert from 'node:assert/strict'
import { createRequire } from 'node:module'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { mkdir, writeFile } from 'node:fs/promises'

const require = createRequire(path.join(process.cwd(), 'package.json'))
const { chromium } = require('@playwright/test')
const outputDir = path.dirname(fileURLToPath(import.meta.url))
const screenshotDir = path.join(outputDir, 'rtl-screenshots')
const baseUrl = process.env.PHASE00_DASHBOARD_URL || 'http://localhost:33654'
const browserPath = process.env.PHASE00_CHROMIUM_PATH || 'C:\\Users\\yosea\\AppData\\Local\\ms-playwright\\chromium-1254\\chrome-win64\\chrome.exe'
const report = {
  source: 'isolated local Next.js production build; synthetic API fixtures for page rendering only',
  scenarios: [],
  consoleErrors: [],
  pageErrors: [],
}

await mkdir(screenshotDir, { recursive: true })
const browser = await chromium.launch({ headless: true, executablePath: browserPath })
const page = await browser.newPage()
page.on('console', (message) => { if (message.type() === 'error') report.consoleErrors.push(message.text()) })
page.on('pageerror', (error) => report.pageErrors.push(error.message))
await page.route('**/api/v1/**', async (route) => {
  const url = new URL(route.request().url())
  if (url.pathname.endsWith('/auth/session')) {
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ id: 710, user_id: 710, email: 'rtl-check@example.invalid', tenant_id: 710, restaurant_id: 810, role: 'owner' }) })
  }
  if (url.pathname.endsWith('/tenant/status')) {
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data: { can_access: true, status: 'active' } }) })
  }
  if (url.pathname.endsWith('/users/710/settings')) {
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data: {} }) })
  }
  if (url.pathname.endsWith('/reports/financial/balance-sheet')) {
    const data = {
      as_of_date: '2026-11-11',
      assets: [{ account_code: 'P00-RTL-CASH', account_name: 'Cash on Hand', amount: 100 }],
      liabilities: [], equity: [], total_assets: 100, total_liabilities: 0, total_equity: 0,
      equation_balanced: false, integrity_status: 'UNRECONCILED', equation_difference: 100,
      integrity_check_basis: 'existing_round2_equation_check',
    }
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data }) })
  }
  return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data: [] }) })
})

const viewports = [
  { locale: 'ar', width: 390, height: 844 },
  { locale: 'ar', width: 768, height: 1024 },
  { locale: 'ar', width: 1440, height: 1000 },
  { locale: 'en', width: 390, height: 844 },
]

try {
  for (const viewport of viewports) {
    await page.setViewportSize({ width: viewport.width, height: viewport.height })
    await page.goto(`${baseUrl}/${viewport.locale}/dashboard/reports/financial/balance-sheet`)
    await page.getByRole('alert').waitFor()
    const measurement = await page.evaluate(() => {
      const localeRoot = document.querySelector('[lang][dir]')
      const main = document.querySelector('main')
      const table = document.querySelector('table')
      const warning = document.querySelector('[role="alert"]')
      const firstHeader = document.querySelector('th')
      const navigation = document.querySelector('nav')
      const computed = (element) => element ? getComputedStyle(element).direction : null
      return {
        viewportWidth: innerWidth,
        documentWidth: document.documentElement.scrollWidth,
        documentLang: document.documentElement.lang,
        documentDirAttribute: document.documentElement.getAttribute('dir'),
        documentComputedDirection: computed(document.documentElement),
        bodyDirAttribute: document.body.getAttribute('dir'),
        bodyComputedDirection: computed(document.body),
        localeRootLang: localeRoot?.getAttribute('lang') ?? null,
        localeRootDir: localeRoot?.getAttribute('dir') ?? null,
        localeRootComputedDirection: computed(localeRoot),
        mainComputedDirection: computed(main),
        navigationComputedDirection: computed(navigation),
        warningComputedDirection: computed(warning),
        warningTextAlign: warning ? getComputedStyle(warning).textAlign : null,
        tableComputedDirection: computed(table),
        tableHeaderTextAlign: firstHeader ? getComputedStyle(firstHeader).textAlign : null,
        pageText: document.body.innerText,
      }
    })
    const scenario = {
      name: `${viewport.locale}-${viewport.width}px`,
      result: measurement.documentWidth <= viewport.width && measurement.localeRootDir === (viewport.locale === 'ar' ? 'rtl' : 'ltr') ? 'PASS' : 'FAIL',
      horizontalOverflow: measurement.documentWidth > viewport.width,
      ...measurement,
    }
    if (viewport.locale === 'ar') {
      assert.equal(measurement.localeRootDir, 'rtl')
      assert.equal(measurement.localeRootComputedDirection, 'rtl')
      assert.equal(measurement.mainComputedDirection, 'rtl')
      assert.equal(measurement.tableComputedDirection, 'rtl')
      assert.equal(measurement.documentWidth > viewport.width, false)
      assert.ok(measurement.pageText.includes('غير متطابقة'))
      assert.ok(measurement.pageText.includes('الأصول'))
    } else {
      assert.equal(measurement.localeRootDir, 'ltr')
      assert.equal(measurement.localeRootComputedDirection, 'ltr')
      assert.equal(measurement.documentWidth > viewport.width, false)
      assert.ok(measurement.pageText.includes('unreconciled'))
    }
    await page.screenshot({ path: path.join(screenshotDir, `p00-${viewport.locale}-${viewport.width}.png`), fullPage: true })
    report.scenarios.push({ ...scenario, pageText: undefined })
  }
} catch (error) {
  report.failure = error.stack || error.message
  throw error
} finally {
  report.consoleErrors = [...new Set(report.consoleErrors)]
  report.pageErrors = [...new Set(report.pageErrors)]
  report.executedAtUtc = new Date().toISOString()
  await writeFile(path.join(outputDir, 'rtl-viewports.json'), JSON.stringify(report, null, 2))
  await browser.close()
}
