const { test, expect } = require('@playwright/test')
const fs = require('node:fs')
const path = require('node:path')

const apiBase = process.env.SERVIO_FINANCIAL_API_URL
if (!apiBase) throw new Error('Disposable financial API URL is required')
const dashboardOrigin = new URL(process.env.SERVIO_FINANCIAL_DASHBOARD_URL).origin
const fixtureFile = process.env.SERVIO_FINANCIAL_BROWSER_FIXTURE_FILE
if (!fixtureFile) throw new Error('Private financial browser fixture file is required')
const fixture = JSON.parse(fs.readFileSync(fixtureFile, 'utf8'))
const identities = {
  accountant: fixture.accountant,
  cashier: fixture.cashier,
  secondTenant: fixture.second_tenant,
}

let postedJournalId
const browserDiagnostics = new WeakMap()

const expectedNotificationDenials = new Set([
  '/api/v1/notifications',
  '/api/v1/notifications/stats',
  '/api/v1/notifications/unread-count',
])

test.beforeEach(async ({ page }) => {
  const diagnostics = { unexpected: [], notification403Responses: 0, generic403ConsoleErrors: 0 }
  browserDiagnostics.set(page, diagnostics)

  page.on('pageerror', (error) => diagnostics.unexpected.push(`page: ${error.message}`))
  page.on('console', (message) => {
    if (message.type() !== 'error') return
    if (message.text() === 'Failed to load resource: the server responded with a status of 403 (Forbidden)') {
      diagnostics.generic403ConsoleErrors += 1
      return
    }
    diagnostics.unexpected.push(`console: ${message.text()}`)
  })
  page.on('response', (response) => {
    if (!response.url().startsWith(apiBase)) return
    const pathname = new URL(response.url()).pathname
    if (response.status() === 403 && response.request().method() === 'GET' && expectedNotificationDenials.has(pathname)) {
      diagnostics.notification403Responses += 1
      return
    }
    if (response.status() >= 400) {
      diagnostics.unexpected.push(`api ${response.status()}: ${response.request().method()} ${pathname}`)
    }
  })
})

test.afterEach(async ({ page }) => {
  const diagnostics = browserDiagnostics.get(page)
  expect(diagnostics.generic403ConsoleErrors, 'console 403s must correspond to the documented notification permission probes')
    .toBe(diagnostics.notification403Responses)
  expect(diagnostics.unexpected, 'no unexpected page, console, or API errors').toEqual([])
})

async function login(page, identity) {
  const response = await page.request.post(`${apiBase}/auth/login`, {
    headers: { Origin: dashboardOrigin },
    data: { email: identity.email, password: identity.password },
  })
  expect(response.status(), 'synthetic browser identity login').toBe(200)
  const sessionResponse = await page.request.get(`${apiBase}/auth/session`)
  expect(sessionResponse.status(), 'authenticated browser session').toBe(200)
  const session = await sessionResponse.json()
  return session.data || session
}

async function setViewport(page, width, height) {
  await page.setViewportSize({ width, height })
}

async function openFinancialPage(page, locale, route) {
  await page.goto(`/${locale}/dashboard/${route}`, { waitUntil: 'domcontentloaded' })
  await expect(page).toHaveURL(new RegExp(`/${locale}/dashboard/${route.replaceAll('/', '\\/')}`))
  await expect(page.locator('main')).toBeVisible()
  await expect(page.locator('html')).toHaveAttribute('lang', locale)
  await expect.poll(() => page.locator('html').getAttribute('dir')).toBe(locale === 'ar' ? 'rtl' : 'ltr')
}

async function assertNoDocumentOverflow(page) {
  const dimensions = await page.evaluate(() => ({
    viewportWidth: window.innerWidth,
    documentWidth: document.documentElement.scrollWidth,
  }))
  expect(dimensions.documentWidth, JSON.stringify(dimensions)).toBeLessThanOrEqual(dimensions.viewportWidth)
}

test.describe('authenticated financial browser acceptance', () => {
  test('accountant EN desktop creates/posts a balanced journal and reconciles the visible ledger', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')
    await openFinancialPage(page, 'en', 'accounting/chart-of-accounts')
    await expect(page.getByRole('heading', { name: /Chart of Accounts/i })).toBeVisible()
    await expect(page.getByText(fixture.account_name, { exact: false }).first()).toBeVisible()
    await page.getByPlaceholder(/Search accounts/i).fill(fixture.account_name)
    await expect(page.getByText(fixture.account_name, { exact: false }).first()).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('en-desktop-chart-of-accounts.png'), fullPage: true })

    await page.goto('/en/dashboard/accounting/journal-entries')
    await page.getByRole('button', { name: /New journal entry/i }).click()
    const dialog = page.getByRole('dialog')
    await dialog.locator('input').nth(2).fill('Browser acceptance balanced journal')
    await dialog.locator('tbody select').nth(0).selectOption(String(fixture.primary_account_id))
    await dialog.locator('tbody select').nth(1).selectOption(String(fixture.revenue_account_id))
    const journalLines = dialog.locator('tbody tr')
    const debitAmount = journalLines.nth(0).locator('input[type="number"]').nth(0)
    const creditAmount = journalLines.nth(1).locator('input[type="number"]').nth(1)
    await debitAmount.fill('50')
    await creditAmount.fill('50')
    await expect(debitAmount).toHaveValue('50')
    await expect(creditAmount).toHaveValue('50')
    await expect(dialog.getByText('Balanced', { exact: true })).toBeVisible()
    const createResponse = page.waitForResponse((response) =>
      response.url().endsWith('/accounting/journal-entries') && response.request().method() === 'POST'
    )
    await dialog.getByRole('button', { name: /Save journal entry/i }).click()
    expect((await createResponse).status()).toBe(201)
    const row = page.locator('tr').filter({ hasText: 'Browser acceptance balanced journal' })
    await expect(row).toBeVisible()
    await row.getByRole('link', { name: /Open/i }).click()
    await page.getByRole('button', { name: /^Post$/i }).click()
    await expect(page.getByText('posted', { exact: true })).toBeVisible()
    const journalHref = page.url()
    postedJournalId = Number(journalHref.split('/').pop())
    expect(postedJournalId).toBeGreaterThan(0)

    await page.goto('/en/dashboard/accounting/general-ledger')
    await expect(page.getByRole('heading', { name: /General Ledger/i })).toBeVisible()
    await expect(page.getByRole('tab', { name: /Trial Balance/i })).toBeVisible()
    await page.getByRole('tab', { name: /Trial Balance/i }).click()
    await expect(page.getByText(/Debit/i).first()).toBeVisible()
    await expect(page.locator('table').last()).toContainText(/50\.00/)
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-ledger.png'), fullPage: true })
  })

  test('cashier EN desktop is denied financial screens and posting', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const session = await login(page, identities.cashier)
    expect(session.role).toBe('cashier')
    await openFinancialPage(page, 'en', 'accounting/journal-entries')
    await expect(page.getByRole('heading', { name: /Access denied/i })).toBeVisible()
    await expect(page.getByRole('button', { name: /New journal entry/i })).toHaveCount(0)
    if (!postedJournalId) throw new Error('Accountant journal prerequisite did not run')
    const response = await page.request.post(`${apiBase}/accounting/journal-entries/${postedJournalId}/post`, {
      headers: { Origin: dashboardOrigin },
      data: {},
    })
    expect(response.status()).toBe(403)
    const exportResponse = await page.request.get(`${apiBase}/accounting/csv/chart-of-accounts/export`)
    expect(exportResponse.status()).toBe(403)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-cashier-denied.png'), fullPage: true })
  })

  test('accountant AR desktop renders financial reports and preserves unreconciled status', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')
    const balanceSheetResponse = page.waitForResponse((response) => {
      const url = new URL(response.url())
      return response.url().startsWith(apiBase) &&
        url.pathname.endsWith('/reports/financial/balance-sheet') &&
        response.request().method() === 'GET' && response.status() === 200
    })
    await openFinancialPage(page, 'ar', 'reports/financial/balance-sheet')
    const balanceSheetPayload = await (await balanceSheetResponse).json()
    const balanceSheetData = balanceSheetPayload.data || balanceSheetPayload
    expect(Number.isFinite(balanceSheetData.total_assets)).toBe(true)
    expect(Number.isFinite(balanceSheetData.total_liabilities)).toBe(true)
    expect(Number.isFinite(balanceSheetData.total_equity)).toBe(true)
    expect(balanceSheetData.integrity_status).toBe('UNRECONCILED')
    expect(balanceSheetData.equation_balanced).toBe(false)
    expect(balanceSheetData.integrity_check_basis).toBe('existing_round2_equation_check')
    expect(Math.round((balanceSheetData.total_assets - balanceSheetData.total_liabilities - balanceSheetData.total_equity) * 100))
      .toBe(Math.round(balanceSheetData.equation_difference * 100))
    await expect(page.getByRole('heading', { name: /الميزانية العمومية|Balance Sheet/i })).toBeVisible()
    await expect(page.getByText(/غير متطابق|UNRECONCILED|unreconciled/i).first()).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('ar-desktop-balance-sheet.png'), fullPage: true })

    await page.goto('/ar/dashboard/reports/financial/profit-loss')
    await expect(page.getByRole('heading', { name: /الأرباح والخسائر|Profit\s*&\s*Loss/i })).toBeVisible()
    await page.goto('/ar/dashboard/reports/financial/cash-flow')
    await expect(page.locator('main')).toBeVisible()
    await assertNoDocumentOverflow(page)
  })

  test('accountant EN mobile navigates journal, P&L, cash-flow, and downloads PDF', async ({ page }, testInfo) => {
    await setViewport(page, 375, 812)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')
    await openFinancialPage(page, 'en', 'accounting/journal-entries')
    await expect(page.locator('main').getByRole('heading', { name: 'Journal Entries', exact: true })).toBeVisible()
    await page.goto('/en/dashboard/reports/financial/profit-loss')
    await expect(page.getByRole('heading', { name: 'Profit & Loss', exact: true })).toBeVisible()
    const pdfButton = page.getByRole('button', { name: /Export PDF/i })
    await expect(pdfButton).toBeVisible()
    const downloadEvent = page.waitForEvent('download')
    await pdfButton.click()
    const download = await downloadEvent
    expect(download.suggestedFilename()).toMatch(/\.pdf$/i)
    const pdfPath = testInfo.outputPath('profit-loss.pdf')
    await download.saveAs(pdfPath)
    expect(fs.readFileSync(pdfPath).subarray(0, 4).toString()).toBe('%PDF')
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-mobile-profit-loss.png'), fullPage: true })
  })

  test('accountant AR mobile preserves RTL and exports CSV through the dashboard', async ({ page }, testInfo) => {
    await setViewport(page, 375, 812)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')
    await openFinancialPage(page, 'ar', 'accounting/advanced')
    await expect(page.getByRole('heading', { name: 'أدوات المحاسبة المتقدمة', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'تسوية البنك', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'القيود المتكررة', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'الموازنات وأسعار الصرف', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'جدولة مدفوعات الموردين والتذكيرات', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'الفاتورة الإلكترونية والاستيراد المتوافق مع Excel', exact: true })).toBeVisible()
    await expect(page.getByLabel('معرّف الحساب البنكي')).toBeVisible()
    const downloadEvent = page.waitForEvent('download')
    await page.getByRole('button', { name: 'تصدير دليل الحسابات CSV', exact: true }).click()
    const download = await downloadEvent
    expect(download.suggestedFilename()).toMatch(/\.csv$/i)
    const csvPath = testInfo.outputPath('chart-of-accounts.csv')
    await download.saveAs(csvPath)
    const csv = fs.readFileSync(csvPath, 'utf8')
    expect(csv).toContain(fixture.account_name)
    const escapedFormulaCell = `QA-FORMULA-1,"'${fixture.formula_name.replaceAll('"', '""')}"`
    expect(csv).toContain(escapedFormulaCell)
    await expect(page.locator('html')).toHaveAttribute('dir', 'rtl')
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('ar-mobile-accounting-tools.png'), fullPage: true })
  })

  test('accountant EN desktop reviews AR/AP registers, aging reports, and payroll', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')

    await openFinancialPage(page, 'en', 'accounting/customer-invoices')
    await expect(page.getByRole('heading', { name: 'Customer Invoices', exact: true })).toBeVisible()
    await expect(page.getByRole('button', { name: 'New customer invoice', exact: true })).toBeEnabled()
    await expect(page.getByRole('button', { name: 'Record payment', exact: true })).toBeEnabled()
    await expect(page.getByRole('heading', { name: 'No invoices found', exact: true })).toBeVisible()
    await page.getByRole('tab', { name: 'Received Payments', exact: true }).click()
    await expect(page.getByRole('heading', { name: 'No received payments yet', exact: true })).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('en-desktop-customer-invoices-payments.png'), fullPage: true })

    await page.goto('/en/dashboard/accounting/supplier-invoices')
    await expect(page.getByRole('heading', { name: 'Supplier Invoices', exact: true })).toBeVisible()
    await expect(page.getByRole('button', { name: 'New supplier invoice', exact: true })).toBeEnabled()
    await page.getByRole('tab', { name: 'Payments Made', exact: true }).click()
    await expect(page.getByRole('heading', { name: 'No supplier payments yet', exact: true })).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('en-desktop-supplier-invoices-payments.png'), fullPage: true })

    await page.goto('/en/dashboard/reports/financial/receivables-aging')
    await expect(page.getByRole('heading', { name: 'Receivables Aging', exact: true })).toBeVisible()
    await expect(page.getByRole('table')).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('en-desktop-receivables-aging.png'), fullPage: true })

    await page.goto('/en/dashboard/reports/financial/payables-aging')
    await expect(page.getByRole('heading', { name: 'Payables Aging', exact: true })).toBeVisible()
    await expect(page.getByRole('table')).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('en-desktop-payables-aging.png'), fullPage: true })

    await page.goto('/en/dashboard/hr/payroll')
    await expect(page.getByRole('heading', { name: 'Payroll', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'No Salary Records', exact: true })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Create Salary', exact: true })).toBeVisible()
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-payroll.png'), fullPage: true })
  })

  test('cashier AR mobile remains denied by server-side accounting RBAC', async ({ page }, testInfo) => {
    await setViewport(page, 375, 812)
    const session = await login(page, identities.cashier)
    expect(session.role).toBe('cashier')
    await openFinancialPage(page, 'ar', 'accounting/chart-of-accounts')
    await expect(page.getByRole('heading', { name: /تم رفض الوصول|Access denied/i })).toBeVisible()
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('ar-mobile-cashier-denied.png'), fullPage: true })
  })

  test('second tenant cannot read or export the first tenant account', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const session = await login(page, identities.secondTenant)
    expect(session.tenant_id).not.toBe(fixture.primary_tenant_id)
    await openFinancialPage(page, 'en', 'accounting/chart-of-accounts')
    const accountList = await page.request.get(`${apiBase}/accounting/accounts`)
    expect(accountList.status()).toBe(200)
    expect(await accountList.text()).not.toContain(fixture.account_name)
    const foreignAccount = await page.request.get(`${apiBase}/accounting/accounts/${fixture.primary_account_id}`)
    expect([403, 404]).toContain(foreignAccount.status())
    const foreignExport = await page.request.get(`${apiBase}/accounting/csv/chart-of-accounts/export`)
    expect(foreignExport.status()).toBe(200)
    expect(await foreignExport.text()).not.toContain(fixture.account_name)
    await expect(page.getByRole('heading', { name: /Chart of Accounts/i })).toBeVisible()
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-second-tenant.png'), fullPage: true })
  })
})
