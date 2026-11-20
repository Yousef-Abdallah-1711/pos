const { test, expect } = require('@playwright/test')
const fs = require('node:fs')
const path = require('node:path')

const apiBase = process.env.SERVIO_FINANCIAL_API_URL
if (!apiBase) throw new Error('Disposable financial API URL is required')
const dashboardOrigin = new URL(process.env.SERVIO_FINANCIAL_DASHBOARD_URL).origin
const fixtureFile = process.env.SERVIO_FINANCIAL_BROWSER_FIXTURE_FILE
if (!fixtureFile) throw new Error('Private financial browser fixture file is required')
const fixture = JSON.parse(fs.readFileSync(fixtureFile, 'utf8'))
const browserJournalAmountCents = 5000
const identities = {
  accountant: fixture.accountant,
  cashier: fixture.cashier,
  payroll: fixture.payroll,
  secondTenant: fixture.second_tenant,
}

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

async function createAccountantDraftJournal(page, identity, description) {
  const session = await login(page, identity)
  expect(session.role).toBe('accountant')

  const createResponse = await page.request.post(`${apiBase}/accounting/journal-entries`, {
    headers: { Origin: dashboardOrigin },
    data: {
      entry_date: new Date().toISOString().slice(0, 10),
      description,
      entry_type: 'manual',
      lines: [
        { account_id: fixture.primary_account_id, debit_amount: 50, credit_amount: 0, line_order: 1 },
        { account_id: fixture.revenue_account_id, debit_amount: 0, credit_amount: 50, line_order: 2 },
      ],
    },
  })
  expect(createResponse.status(), 'accountant creates a balanced authorization target').toBe(201)
  const payload = await createResponse.json()
  const journal = payload.data || payload
  const journalId = Number(journal.id)
  expect(journalId).toBeGreaterThan(0)

  return journalId
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

async function expectFinancialMetric(page, label, amountCents, locale) {
  const formatter = new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: 'EGP',
    maximumFractionDigits: 2,
  })
  const value = page.getByText(label, { exact: true }).locator('xpath=following-sibling::div')
  await expect(value).toHaveText(formatter.format(amountCents / 100))
}

async function expectBalanceSheetMetrics(page, locale, expected) {
  const labels = locale === 'ar'
    ? ['إجمالي الأصول', 'إجمالي الالتزامات', 'إجمالي حقوق الملكية']
    : ['Total Assets', 'Total Liabilities', 'Total Equity']
  await expectFinancialMetric(page, labels[0], expected.assets_cents, locale)
  await expectFinancialMetric(page, labels[1], expected.liabilities_cents, locale)
  await expectFinancialMetric(page, labels[2], expected.equity_cents, locale)
}

async function assertNoDocumentOverflow(page) {
  const dimensions = await page.evaluate(() => ({
    viewportWidth: window.innerWidth,
    documentWidth: document.documentElement.scrollWidth,
  }))
  expect(dimensions.documentWidth, JSON.stringify(dimensions)).toBeLessThanOrEqual(dimensions.viewportWidth)
}

async function assertMobileFinancialContentWithinViewport(page) {
  const geometry = await page.evaluate(() => {
    const main = document.querySelector('main')?.getBoundingClientRect()
    const heading = document.querySelector('main h1')?.getBoundingClientRect()
    const tableViewports = [...document.querySelectorAll('main table')].map((table) => {
      const viewport = table.closest('.overflow-x-auto') || table
      const bounds = viewport.getBoundingClientRect()
      return { left: bounds.left, right: bounds.right }
    })

    return {
      viewportWidth: window.innerWidth,
      documentWidth: document.documentElement.scrollWidth,
      main: main ? { left: main.left, right: main.right } : null,
      heading: heading ? { left: heading.left, right: heading.right } : null,
      tableViewports,
    }
  })

  expect(geometry.documentWidth, JSON.stringify(geometry)).toBeLessThanOrEqual(geometry.viewportWidth)
  expect(geometry.main, JSON.stringify(geometry)).not.toBeNull()
  expect(geometry.heading, JSON.stringify(geometry)).not.toBeNull()
  expect(geometry.tableViewports, 'financial report must expose its table viewport').not.toHaveLength(0)

  for (const element of [geometry.main, geometry.heading, ...geometry.tableViewports]) {
    expect(element.left, JSON.stringify(geometry)).toBeGreaterThanOrEqual(0)
    expect(element.right, JSON.stringify(geometry)).toBeLessThanOrEqual(geometry.viewportWidth)
  }
}

async function assertFinancialContentClearsSidebar(page) {
  await expect(page.locator('main table').first(), 'financial content table mounts before layout checks').toBeVisible()
  const layout = await page.evaluate(() => {
    const heading = document.querySelector('main h1')?.getBoundingClientRect()
    const sidebar = [...document.querySelectorAll('aside.fixed')]
      .map((element) => ({
        bounds: element.getBoundingClientRect(),
        style: window.getComputedStyle(element),
      }))
      .find(({ bounds, style }) => (
        bounds.width > 0
        && bounds.height > 0
        && bounds.right > 0
        && bounds.left < window.innerWidth
        && style.visibility !== 'hidden'
        && style.display !== 'none'
      ))?.bounds
    const main = document.querySelector('main')?.getBoundingClientRect()
    const content = [...document.querySelectorAll('main table')].map((table) => {
      const viewport = table.closest('.overflow-x-auto')?.getBoundingClientRect() || table.getBoundingClientRect()
      return {
        left: viewport.left,
        right: viewport.right,
        overlapsSidebar: sidebar
          ? viewport.left < sidebar.right && viewport.right > sidebar.left
          : true,
      }
    })
    if (!heading || !sidebar || !main) {
      return { foundHeading: Boolean(heading), foundSidebar: Boolean(sidebar), foundMain: Boolean(main), content }
    }
    const intersectsSidebar = (bounds) => bounds.left < sidebar.right && bounds.right > sidebar.left
    return {
      foundHeading: true,
      foundSidebar: true,
      foundMain: true,
      headingOverlaps: intersectsSidebar(heading),
      mainOverlaps: intersectsSidebar(main),
      content,
    }
  })
  expect(layout, JSON.stringify(layout)).toMatchObject({
    foundHeading: true,
    foundSidebar: true,
    foundMain: true,
    headingOverlaps: false,
    mainOverlaps: false,
  })
  expect(layout.content, 'financial tables must render in the unobstructed main pane').not.toHaveLength(0)
  for (const content of layout.content) {
    expect(content.overlapsSidebar, JSON.stringify(layout)).toBe(false)
  }
}

test.describe('authenticated financial browser acceptance', () => {
  test('accountant EN desktop posts a balanced asset transfer and reconciles the visible ledger', async ({ page }, testInfo) => {
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
    // The counter-entry is another debit-normal asset so aggregate Balance Sheet values stay stable across test selection/order.
    await dialog.locator('tbody select').nth(1).selectOption(String(fixture.contra_asset_account_id))
    const journalLines = dialog.locator('tbody tr')
    const debitAmount = journalLines.nth(0).locator('input[type="number"]').nth(0)
    const creditAmount = journalLines.nth(1).locator('input[type="number"]').nth(1)
    const browserJournalAmount = String(browserJournalAmountCents / 100)
    await debitAmount.fill(browserJournalAmount)
    await creditAmount.fill(browserJournalAmount)
    await expect(debitAmount).toHaveValue(browserJournalAmount)
    await expect(creditAmount).toHaveValue(browserJournalAmount)
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
    const postedJournalId = Number(journalHref.split('/').pop())
    expect(postedJournalId).toBeGreaterThan(0)

    const duplicatePostResponse = await page.request.post(
      `${apiBase}/accounting/journal-entries/${postedJournalId}/post`,
      { headers: { Origin: dashboardOrigin }, data: {} },
    )
    expect(duplicatePostResponse.status()).toBe(400)
    expect((await duplicatePostResponse.json()).error).toBe('only draft journal entries can be posted')
    const ledgerResponse = await page.request.get(
      `${apiBase}/accounting/ledger?account_id=${fixture.primary_account_id}`,
    )
    expect(ledgerResponse.status()).toBe(200)
    const ledgerPayload = await ledgerResponse.json()
    const postedJournalRows = (ledgerPayload.data?.lines || []).filter(
      (line) => line.journal_entry_id === postedJournalId,
    )
    expect(postedJournalRows).toHaveLength(1)
    expect(postedJournalRows[0].account_id).toBe(fixture.primary_account_id)
    expect(Math.round(postedJournalRows[0].debit_amount * 100)).toBe(browserJournalAmountCents)
    expect(Math.round(postedJournalRows[0].credit_amount * 100)).toBe(0)
    const contraAssetLedgerResponse = await page.request.get(
      `${apiBase}/accounting/ledger?account_id=${fixture.contra_asset_account_id}`,
    )
    expect(contraAssetLedgerResponse.status()).toBe(200)
    const contraAssetLedgerPayload = await contraAssetLedgerResponse.json()
    const postedCounterRows = (contraAssetLedgerPayload.data?.lines || []).filter(
      (line) => line.journal_entry_id === postedJournalId,
    )
    expect(postedCounterRows).toHaveLength(1)
    expect(postedCounterRows[0].account_id).toBe(fixture.contra_asset_account_id)
    expect(Math.round(postedCounterRows[0].debit_amount * 100)).toBe(0)
    expect(Math.round(postedCounterRows[0].credit_amount * 100)).toBe(browserJournalAmountCents)

    await page.goto('/en/dashboard/accounting/general-ledger')
    await expect(page.getByRole('heading', { name: /General Ledger/i })).toBeVisible()
    await assertFinancialContentClearsSidebar(page)
    await expect(page.getByRole('tab', { name: /Trial Balance/i })).toBeVisible()
    await page.getByRole('tab', { name: /Trial Balance/i }).click()
    await expect(page.getByText(/Debit/i).first()).toBeVisible()
    await expect(page.locator('table').last()).toContainText(/50\.00/)
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-ledger.png'), fullPage: true })
  })

  test('cashier EN desktop is denied financial screens and posting', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const draftJournalId = await createAccountantDraftJournal(page, identities.accountant, 'Cashier authorization target')
    const session = await login(page, identities.cashier)
    expect(session.role).toBe('cashier')
    await openFinancialPage(page, 'en', 'accounting/journal-entries')
    await expect(page.getByRole('heading', { name: /Access denied/i })).toBeVisible()
    await expect(page.getByRole('button', { name: /New journal entry/i })).toHaveCount(0)
    const response = await page.request.post(`${apiBase}/accounting/journal-entries/${draftJournalId}/post`, {
      headers: { Origin: dashboardOrigin },
      data: {},
    })
    expect(response.status()).toBe(403)
    const exportResponse = await page.request.get(`${apiBase}/accounting/csv/chart-of-accounts/export`)
    expect(exportResponse.status()).toBe(403)
    for (const format of ['csv', 'pdf']) {
      const reportExportResponse = await page.request.post(`${apiBase}/reports/export`, {
        headers: { Origin: dashboardOrigin },
        data: { report_type: 'sales', format },
      })
      expect(reportExportResponse.status(), `cashier ${format.toUpperCase()} export is denied`).toBe(403)
    }
    await page.screenshot({ path: testInfo.outputPath('en-desktop-cashier-denied.png'), fullPage: true })
  })

  test('accountant EN/AR desktop and mobile Balance Sheet values match the SQL oracle', async ({ page }, testInfo) => {
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
    expect(Math.round(balanceSheetData.total_assets * 100)).toBe(fixture.expected_balance_sheet.assets_cents)
    expect(Math.round(balanceSheetData.total_liabilities * 100)).toBe(fixture.expected_balance_sheet.liabilities_cents)
    expect(Math.round(balanceSheetData.total_equity * 100)).toBe(fixture.expected_balance_sheet.equity_cents)
    expect(Math.round((balanceSheetData.total_assets - balanceSheetData.total_liabilities - balanceSheetData.total_equity) * 100))
      .toBe(Math.round(balanceSheetData.equation_difference * 100))
    await expectBalanceSheetMetrics(page, 'ar', fixture.expected_balance_sheet)
    await expect(page.getByRole('heading', { name: /الميزانية العمومية|Balance Sheet/i })).toBeVisible()
    await assertFinancialContentClearsSidebar(page)
    await expect(page.getByText(/غير متطابق|UNRECONCILED|unreconciled/i).first()).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('ar-desktop-balance-sheet.png'), fullPage: true })

    await setViewport(page, 375, 812)
    await expectBalanceSheetMetrics(page, 'ar', fixture.expected_balance_sheet)
    await assertMobileFinancialContentWithinViewport(page)
    await page.screenshot({ path: testInfo.outputPath('ar-mobile-balance-sheet.png'), fullPage: true })

    await openFinancialPage(page, 'en', 'reports/financial/balance-sheet')
    await expect(page.getByRole('heading', { name: 'Balance Sheet', exact: true })).toBeVisible()
    await expectBalanceSheetMetrics(page, 'en', fixture.expected_balance_sheet)
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-balance-sheet.png'), fullPage: true })

    await setViewport(page, 375, 812)
    await expectBalanceSheetMetrics(page, 'en', fixture.expected_balance_sheet)
    await assertMobileFinancialContentWithinViewport(page)
    await page.screenshot({ path: testInfo.outputPath('en-mobile-balance-sheet.png'), fullPage: true })

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

    await setViewport(page, 320, 812)
    await expect(page.getByRole('heading', { name: 'Profit & Loss', exact: true })).toBeVisible()
    await assertMobileFinancialContentWithinViewport(page)
    await page.screenshot({ path: testInfo.outputPath('en-mobile-profit-loss-320.png'), fullPage: true })
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
    await expect(page.getByRole('button', { name: 'Create Salary', exact: true })).toBeDisabled()
    const unauthorizedSalaryCreation = await page.request.post(`${apiBase}/hr/salaries`, {
      headers: { Origin: dashboardOrigin },
      data: {},
    })
    expect(unauthorizedSalaryCreation.status()).toBe(403)
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('en-desktop-payroll.png'), fullPage: true })
  })

  test('accountant EN desktop creates an AR invoice, receives a partial payment, and sees the open balance in aging', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    await login(page, identities.accountant)
    await openFinancialPage(page, 'en', 'accounting/customer-invoices')

    await page.getByRole('button', { name: 'New customer invoice', exact: true }).click()
    const invoiceDialog = page.getByRole('dialog')
    await invoiceDialog.getByText('Customer', { exact: true }).locator('..').locator('select').selectOption({ label: fixture.customer_name })
    await invoiceDialog.getByText('Status', { exact: true }).locator('..').locator('select').selectOption('sent')
    await invoiceDialog.getByText('Description', { exact: true }).locator('..').locator('input').fill('Browser synthetic service')
    await invoiceDialog.getByText('Quantity', { exact: true }).locator('..').locator('input').fill('1')
    await invoiceDialog.getByText('Unit price', { exact: true }).locator('..').locator('input').fill('100')
    const createInvoice = page.waitForResponse((response) =>
      new URL(response.url()).pathname === '/api/v1/accounting/customer-invoices' && response.request().method() === 'POST'
    )
    await invoiceDialog.getByRole('button', { name: 'Create customer invoice', exact: true }).click()
    const invoiceResponse = await createInvoice
    expect(invoiceResponse.status()).toBe(201)
    const invoicePayload = await invoiceResponse.json()
    expect(invoicePayload.data.total_amount).toBe(100)
    expect(invoicePayload.data.balance_due).toBe(100)
    const invoiceNumber = invoicePayload.data.invoice_number
    const invoiceRow = page.locator('tbody tr').filter({ hasText: fixture.customer_name }).filter({ hasText: invoiceNumber })
    await expect(invoiceRow).toContainText('100.00')

    await page.getByRole('button', { name: 'Record payment', exact: true }).click()
    const paymentDialog = page.getByRole('dialog')
    await paymentDialog.getByText('Customer invoice', { exact: true }).locator('..').locator('select').selectOption({ label: `${invoiceNumber} • ${fixture.customer_name}` })
    await paymentDialog.getByText('Amount', { exact: true }).locator('..').locator('input').fill('40')
    await paymentDialog.getByText('Payment method', { exact: true }).locator('..').locator('select').selectOption('cash')
    const createPayment = page.waitForResponse((response) =>
      new URL(response.url()).pathname === '/api/v1/accounting/payments/received' && response.request().method() === 'POST'
    )
    await paymentDialog.getByRole('button', { name: 'Record payment received', exact: true }).click()
    const paymentResponse = await createPayment
    expect(paymentResponse.status()).toBe(201)
    const paymentPayload = await paymentResponse.json()
    expect(paymentPayload.data.id).toBeGreaterThan(0)
    await expect(invoiceRow).toContainText('60.00')

    await page.goto('/en/dashboard/reports/financial/receivables-aging')
    await expect(page.getByRole('heading', { name: 'Receivables Aging', exact: true })).toBeVisible()
    const receivablesAgingRow = page.locator('tbody tr').filter({ hasText: fixture.customer_name })
    await expect(receivablesAgingRow).toBeVisible()
    await expect(receivablesAgingRow.getByRole('cell').last()).toContainText('60.00')
    await page.screenshot({ path: testInfo.outputPath('en-desktop-ar-partial-payment-aging.png'), fullPage: true })
  })

  test('accountant EN desktop approves an AP invoice, pays partially, and sees remaining balance in payables aging', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    await login(page, identities.accountant)
    await openFinancialPage(page, 'en', 'accounting/supplier-invoices')

    await page.getByRole('button', { name: 'New supplier invoice', exact: true }).click()
    const invoiceDialog = page.getByRole('dialog')
    await invoiceDialog.getByText('Supplier', { exact: true }).locator('..').locator('select').selectOption({ label: fixture.supplier_name })
    await invoiceDialog.getByText('Supplier reference', { exact: true }).locator('..').locator('input').fill('BROWSER-AP-PARTIAL')
    await invoiceDialog.getByText('Description', { exact: true }).locator('..').locator('input').fill('Browser synthetic supplies')
    await invoiceDialog.getByText('Quantity', { exact: true }).locator('..').locator('input').fill('1')
    await invoiceDialog.getByText('Unit price', { exact: true }).locator('..').locator('input').fill('100')
    const createInvoice = page.waitForResponse((response) =>
      new URL(response.url()).pathname === '/api/v1/accounting/supplier-invoices' && response.request().method() === 'POST'
    )
    await invoiceDialog.getByRole('button', { name: 'Create supplier invoice', exact: true }).click()
    const invoiceResponse = await createInvoice
    expect(invoiceResponse.status()).toBe(201)
    const invoicePayload = await invoiceResponse.json()
    expect(invoicePayload.data.total_amount).toBe(100)
    const invoiceNumber = invoicePayload.data.invoice_number
    let invoiceRow = page.locator('tbody tr').filter({ hasText: fixture.supplier_name }).filter({ hasText: invoiceNumber })
    const approvalResponseWait = page.waitForResponse((response) =>
      new URL(response.url()).pathname === `/api/v1/accounting/supplier-invoices/${invoicePayload.data.id}/approve` && response.request().method() === 'POST'
    )
    await invoiceRow.getByRole('button', { name: 'Approve', exact: true }).click()
    const approvalResponse = await approvalResponseWait
    expect(approvalResponse.status()).toBe(200)
    await expect(invoiceRow).toContainText('approved')

    await page.getByRole('button', { name: 'Record payment', exact: true }).click()
    const paymentDialog = page.getByRole('dialog')
    await paymentDialog.getByText('Supplier invoice', { exact: true }).locator('..').locator('select').selectOption({ label: `${invoiceNumber} • ${fixture.supplier_name}` })
    await paymentDialog.getByText('Amount', { exact: true }).locator('..').locator('input').fill('40')
    await paymentDialog.getByText('Payment method', { exact: true }).locator('..').locator('select').selectOption('cash')
    const createPayment = page.waitForResponse((response) =>
      new URL(response.url()).pathname === '/api/v1/accounting/payments/made' && response.request().method() === 'POST'
    )
    await paymentDialog.getByRole('button', { name: 'Record payment made', exact: true }).click()
    const paymentResponse = await createPayment
    expect(paymentResponse.status()).toBe(201)
    const paymentPayload = await paymentResponse.json()
    expect(paymentPayload.data.id).toBeGreaterThan(0)
    await expect(invoiceRow).toContainText('60.00')

    const overpaymentResponse = await page.request.post(`${apiBase}/accounting/payments/made`, {
      headers: { Origin: dashboardOrigin },
      data: {
        payment_date: new Date().toISOString().slice(0, 10),
        supplier_id: invoicePayload.data.supplier_id,
        invoice_id: invoicePayload.data.id,
        amount: 60.01,
        payment_method: 'cash',
        reference: 'BROWSER-AP-OVERPAYMENT',
      },
    })
    expect(overpaymentResponse.status()).toBe(400)
    expect((await overpaymentResponse.json()).error).toBe('payment amount exceeds invoice balance')
    await page.reload()
    await expect(invoiceRow).toContainText('60.00')

    await page.goto('/en/dashboard/reports/financial/payables-aging')
    const agingRow = page.locator('tbody tr').filter({ hasText: fixture.supplier_name })
    await expect(agingRow).toBeVisible()
    await expect(agingRow).toContainText('60.00')
    await page.screenshot({ path: testInfo.outputPath('en-desktop-ap-partial-payment-aging.png'), fullPage: true })
  })

  test('payroll manager EN desktop creates and marks a synthetic salary paid', async ({ page }, testInfo) => {
    await setViewport(page, 1440, 900)
    const session = await login(page, identities.payroll)
    expect(session.role).toBe('payroll-manager')
    await openFinancialPage(page, 'en', 'hr/payroll')
    await page.getByRole('button', { name: 'Create Salary', exact: true }).click()
    const salaryDialog = page.getByRole('dialog')
    await salaryDialog.locator('select[name="employee_id"]').selectOption({ label: fixture.employee_name })
    const today = new Date()
    const payPeriodStart = new Date(Date.UTC(today.getUTCFullYear(), today.getUTCMonth(), 1)).toISOString().slice(0, 10)
    const payPeriodEnd = today.toISOString().slice(0, 10)
    await salaryDialog.locator('input[name="pay_period_start"]').fill(payPeriodStart)
    await salaryDialog.locator('input[name="pay_period_end"]').fill(payPeriodEnd)
    await salaryDialog.locator('input[name="base_salary"]').fill('1500')
    const createSalary = page.waitForResponse((response) =>
      new URL(response.url()).pathname === '/api/v1/hr/salaries' && response.request().method() === 'POST'
    )
    await salaryDialog.getByRole('button', { name: /Create Salary/i }).click()
    const salaryResponse = await createSalary
    expect(salaryResponse.status()).toBe(201)
    const salaryPayload = await salaryResponse.json()
    expect(salaryPayload.data.net_salary).toBe(1500)
    const salaryId = salaryPayload.data.id
    const salaryRow = page.locator('tbody tr').filter({ hasText: fixture.employee_name })
    await expect(salaryRow).toContainText('Pending')
    const markPaid = page.waitForResponse((response) =>
      new URL(response.url()).pathname === `/api/v1/hr/salaries/${salaryId}/mark-paid` && response.request().method() === 'POST'
    )
    await salaryRow.getByRole('button', { name: 'Mark Paid', exact: true }).click()
    const paidResponse = await markPaid
    expect(paidResponse.status()).toBe(200)
    await expect(salaryRow).toContainText('Paid')
    await page.screenshot({ path: testInfo.outputPath('en-desktop-payroll-salary-paid.png'), fullPage: true })
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
