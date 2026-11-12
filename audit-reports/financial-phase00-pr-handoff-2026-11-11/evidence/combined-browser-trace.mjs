import assert from 'node:assert/strict'
import { createRequire } from 'node:module'
import path from 'node:path'
import { mkdir, writeFile } from 'node:fs/promises'

const require = createRequire(path.join(process.cwd(), 'package.json'))
const { chromium } = require('@playwright/test')
const reportDir = path.resolve(process.cwd(), '..', '..', '..', 'audit-reports', 'financial-remediation-execution-2026-11-11')
const screenshotDir = path.join(reportDir, 'screenshots')
const resultDir = path.join(reportDir, 'test-results')
const baseUrl = process.env.PHASE00_DASHBOARD_URL || 'http://localhost:33653'
const browserPath = process.env.PHASE00_CHROMIUM_PATH || 'C:\\Users\\yosea\\AppData\\Local\\ms-playwright\\chromium-1254\\chrome-win64\\chrome.exe'

await mkdir(screenshotDir, { recursive: true })
await mkdir(resultDir, { recursive: true })

const browser = await chromium.launch({ headless: true, executablePath: browserPath })
const page = await browser.newPage({ viewport: { width: 1440, height: 1000 } })
const evidence = { source: 'synthetic API fixtures; isolated local Next.js production build', requests: [], responses: [], consoleErrors: [], pageErrors: [], scenarios: [] }
let reportMode = 'unreconciled'
let delayedReport = false

const fixtures = {
  unreconciled: {
    as_of_date: '2026-11-11',
    assets: [{ account_code: '1110', account_name: 'Cash on Hand', amount: 100 }],
    liabilities: [], equity: [], total_assets: 100, total_liabilities: 0, total_equity: 0,
    equation_balanced: false, integrity_status: 'UNRECONCILED', equation_difference: 100,
    integrity_check_basis: 'existing_round2_equation_check',
  },
  balanced: {
    as_of_date: '2026-11-11',
    assets: [{ account_code: '1110', account_name: 'Cash on Hand', amount: 100 }],
    liabilities: [{ account_code: '2101', account_name: 'Payables', amount: 40 }],
    equity: [{ account_code: '3200', account_name: 'Retained Earnings', amount: 60 }],
    total_assets: 100, total_liabilities: 40, total_equity: 60,
    equation_balanced: true, integrity_status: 'RECONCILED', equation_difference: 0,
    integrity_check_basis: 'existing_round2_equation_check',
  },
}

page.on('console', (message) => {
  if (message.type() === 'error') evidence.consoleErrors.push(message.text())
})
page.on('pageerror', (error) => evidence.pageErrors.push({ url: page.url(), message: error.message }))
page.on('response', (response) => {
  if (response.url().includes('/api/v1/')) evidence.responses.push({ status: response.status(), url: response.url() })
})

// The dashboard redirects unauthorized users to the separate landing app. Fulfill
// that route in this isolated browser harness; no landing service is part of this test.
await page.route('**/auth/login**', async (route) => route.fulfill({
  status: 200,
  contentType: 'text/html',
  body: '<!doctype html><html lang="en"><body>synthetic login page</body></html>',
}))

await page.route('**/api/v1/**', async (route) => {
  const url = new URL(route.request().url())
  evidence.requests.push({ method: route.request().method(), path: `${url.pathname}${url.search}` })

  if (url.pathname.endsWith('/auth/session')) {
    return route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ id: 1, user_id: 1, email: 'synthetic-owner@example.invalid', name: 'Synthetic Owner', tenant_id: 71001, restaurant_id: 81001, role: 'owner', is_super_admin: false }),
    })
  }
  if (url.pathname.endsWith('/tenant/status')) {
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data: { can_access: true, status: 'active' } }) })
  }
  if (url.pathname.endsWith('/users/1/settings')) {
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data: {} }) })
  }
  if (url.pathname.endsWith('/reports/financial/balance-sheet')) {
    if (reportMode === 'unauthorized') return route.fulfill({ status: 401, contentType: 'application/json', body: JSON.stringify({ success: false, error: 'Unauthorized' }) })
    if (reportMode === 'error') return route.fulfill({ status: 500, contentType: 'application/json', body: JSON.stringify({ success: false, error: 'Synthetic report error' }) })
    if (delayedReport) await new Promise((resolve) => setTimeout(resolve, 1800))
    const data = reportMode === 'unknown'
      ? { ...fixtures.balanced, integrity_status: 'FUTURE_STATUS' }
      : fixtures[reportMode]
    return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data }) })
  }
  return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ success: true, data: [] }) })
})

try {
  await page.goto(`${baseUrl}/en/dashboard/reports/financial/balance-sheet`)
  await page.getByText('This balance sheet is unreconciled').waitFor()
  let body = await page.locator('body').innerText()
  const unreconciled = {
    alert: body.includes('This balance sheet is unreconciled'),
    sourceFigurePreserved: body.includes('Cash on Hand') && body.includes('100.00'),
    pdfExportButtons: await page.getByRole('button', { name: /export pdf/i }).count(),
  }
  assert.deepEqual(unreconciled, { alert: true, sourceFigurePreserved: true, pdfExportButtons: 0 })
  evidence.scenarios.push({ name: 'unreconciled', result: 'PASS', ...unreconciled })
  await page.screenshot({ path: path.join(screenshotDir, 'p00-unreconciled-en-desktop.png'), fullPage: true })

  reportMode = 'balanced'
  await page.reload()
  await page.getByRole('status').waitFor()
  body = await page.locator('body').innerText()
  const balanced = {
    accountantCaveat: body.includes('not an accountant'),
    figuresPreserved: body.includes('100.00') && body.includes('40.00') && body.includes('60.00'),
    pdfExportButtons: await page.getByRole('button', { name: /export pdf/i }).count(),
  }
  assert.deepEqual(balanced, { accountantCaveat: true, figuresPreserved: true, pdfExportButtons: 1 })
  let downloadResult = 'NOT VERIFIED'
  try {
    const downloadPromise = page.waitForEvent('download', { timeout: 20000 })
    await page.getByRole('button', { name: /export pdf/i }).click()
    const download = await downloadPromise
    const pdfPath = path.join(resultDir, 'synthetic-balanced-balance-sheet.pdf')
    await download.saveAs(pdfPath)
    downloadResult = `PASS: ${path.basename(pdfPath)}`
  } catch (error) {
    downloadResult = `FAIL: ${error.message}`
  }
  assert.ok(downloadResult.startsWith('PASS'), downloadResult)
  evidence.scenarios.push({ name: 'balanced', result: 'PASS', ...balanced, pdfDownload: downloadResult })
  await page.screenshot({ path: path.join(screenshotDir, 'p00-balanced-en-desktop.png'), fullPage: true })

  reportMode = 'unknown'
  await page.reload()
  await page.getByText('This balance sheet is unreconciled').waitFor()
  body = await page.locator('body').innerText()
  const unknown = {
    warningVisible: body.includes('This balance sheet is unreconciled') && body.includes('does not balance'),
    pdfExportButtons: await page.getByRole('button', { name: /export pdf/i }).count(),
  }
  assert.deepEqual(unknown, { warningVisible: true, pdfExportButtons: 0 })
  evidence.scenarios.push({ name: 'unknown-integrity-status', result: 'PASS', ...unknown })

  reportMode = 'unreconciled'
  await page.setViewportSize({ width: 390, height: 844 })
  await page.goto(`${baseUrl}/ar/dashboard/reports/financial/balance-sheet`)
  await page.getByRole('alert').waitFor()
  body = await page.locator('body').innerText()
  const rtl = await page.evaluate(() => ({
    viewportWidth: innerWidth,
    scrollWidth: document.documentElement.scrollWidth,
    direction: document.documentElement.dir || getComputedStyle(document.body).direction,
  }))
  const mobileRtl = { arabicWarning: body.includes('غير متطابقة'), ...rtl, horizontalOverflow: rtl.scrollWidth > rtl.viewportWidth }
  assert.equal(mobileRtl.arabicWarning, true)
  assert.equal(mobileRtl.horizontalOverflow, false)
  evidence.scenarios.push({ name: 'mobile-rtl-390px', result: 'PASS', ...mobileRtl })
  await page.screenshot({ path: path.join(screenshotDir, 'p00-unreconciled-ar-mobile.png'), fullPage: true })

  reportMode = 'unreconciled'
  delayedReport = true
  await page.setViewportSize({ width: 1440, height: 1000 })
  await page.goto(`${baseUrl}/en/dashboard/reports/financial/balance-sheet`)
  const loadingVisible = await page.getByText('Loading report').isVisible().catch(() => false)
  await page.getByText('This balance sheet is unreconciled').waitFor()
  delayedReport = false
  evidence.scenarios.push({ name: 'loading-state', result: loadingVisible ? 'PASS' : 'FAIL', loadingVisible })
  assert.equal(loadingVisible, true)

  reportMode = 'error'
  await page.reload()
  await page.getByText('Synthetic report error').waitFor()
  evidence.scenarios.push({ name: 'api-error-state', result: 'PASS', userVisibleError: true })

  reportMode = 'unauthorized'
  await page.reload()
  await page.waitForURL(/\/auth\/login(?:\?.*)?$/, { timeout: 15000 })
  const unauthorizedBody = await page.locator('body').innerText()
  assert.equal(unauthorizedBody.includes('Cash on Hand'), false)
  evidence.scenarios.push({ name: 'unauthorized-report-response', result: 'PASS', redirectedToLogin: true, financialFiguresWithheld: true })
} catch (error) {
  evidence.failure = error.stack || error.message
  throw error
} finally {
  evidence.requestCount = evidence.requests.length
  evidence.consoleErrors = [...new Set(evidence.consoleErrors)]
  evidence.pageErrors = [...new Set(evidence.pageErrors)]
  await writeFile(path.join(resultDir, 'browser-acceptance.json'), JSON.stringify(evidence, null, 2))
  await browser.close()
}

