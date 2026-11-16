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
  const errors = []
  page.on('pageerror', (error) => errors.push(`page: ${error.message}`))
  page.on('console', (message) => {
    if (message.type() === 'error') errors.push(`console: ${message.text()}`)
  })
  page.on('response', (response) => {
    const url = response.url()
    if (url.startsWith(apiBase) && response.status() >= 500) {
      errors.push(`api ${response.status()}: ${new URL(url).pathname}`)
    }
  })
  await page.goto(`/${locale}/dashboard/${route}`, { waitUntil: 'domcontentloaded' })
  await expect(page).toHaveURL(new RegExp(`/${locale}/dashboard/${route.replaceAll('/', '\\/')}`))
  await expect(page.locator('main')).toBeVisible()
  await expect.poll(() => page.locator('html').getAttribute('dir')).toBe(locale === 'ar' ? 'rtl' : 'ltr')
  return errors
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
    const errors = await openFinancialPage(page, 'en', 'accounting/chart-of-accounts')
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
    await dialog.locator('input[type="number"]').nth(0).fill('50')
    await dialog.locator('input[type="number"]').nth(1).fill('50')
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
    expect(errors).toEqual([])
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
    const balanceSheetResponse = page.waitForResponse((response) =>
      response.url().includes('/reports/financial/balance-sheet') && response.status() === 200
    )
    const errors = await openFinancialPage(page, 'ar', 'reports/financial/balance-sheet')
    const balanceSheetPayload = await (await balanceSheetResponse).json()
    const balanceSheetData = balanceSheetPayload.data || balanceSheetPayload
    expect(Math.round(balanceSheetData.total_assets * 100)).toBe(5000)
    expect(balanceSheetData.integrity_status).toBe('UNRECONCILED')
    await expect(page.getByRole('heading', { name: /الميزانية العمومية|Balance Sheet/i })).toBeVisible()
    await expect(page.getByText(/غير متطابق|UNRECONCILED|unreconciled/i).first()).toBeVisible()
    await page.screenshot({ path: testInfo.outputPath('ar-desktop-balance-sheet.png'), fullPage: true })

    await page.goto('/ar/dashboard/reports/financial/profit-loss')
    await expect(page.getByRole('heading', { name: /الأرباح والخسائر|Profit.?Loss/i })).toBeVisible()
    await page.goto('/ar/dashboard/reports/financial/cash-flow')
    await expect(page.locator('main')).toBeVisible()
    await assertNoDocumentOverflow(page)
    expect(errors).toEqual([])
  })

  test('accountant EN mobile navigates journal, P&L, cash-flow, and downloads PDF', async ({ page }, testInfo) => {
    await setViewport(page, 375, 812)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')
    const errors = await openFinancialPage(page, 'en', 'accounting/journal-entries')
    await expect(page.getByRole('heading', { name: /Journal Entries/i })).toBeVisible()
    await page.goto('/en/dashboard/reports/financial/profit-loss')
    await expect(page.getByRole('heading', { name: /Profit.?Loss/i })).toBeVisible()
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
    expect(errors).toEqual([])
    await page.screenshot({ path: testInfo.outputPath('en-mobile-profit-loss.png'), fullPage: true })
  })

  test('accountant AR mobile preserves RTL and exports CSV through the dashboard', async ({ page }, testInfo) => {
    await setViewport(page, 375, 812)
    const session = await login(page, identities.accountant)
    expect(session.role).toBe('accountant')
    await openFinancialPage(page, 'ar', 'accounting/advanced')
    await expect(page.getByRole('heading', { name: /Advanced Accounting Tools/i })).toBeVisible()
    const downloadEvent = page.waitForEvent('download')
    await page.getByRole('button', { name: /Export COA CSV/i }).click()
    const download = await downloadEvent
    expect(download.suggestedFilename()).toMatch(/\.csv$/i)
    const csvPath = testInfo.outputPath('chart-of-accounts.csv')
    await download.saveAs(csvPath)
    const csv = fs.readFileSync(csvPath, 'utf8')
    expect(csv).toContain(fixture.account_name)
    expect(csv).toContain(`'${fixture.formula_name}`)
    await expect(page.locator('html')).toHaveAttribute('dir', 'rtl')
    await assertNoDocumentOverflow(page)
    await page.screenshot({ path: testInfo.outputPath('ar-mobile-accounting-tools.png'), fullPage: true })
  })

  test('cashier AR mobile remains denied by server-side accounting RBAC', async ({ page }, testInfo) => {
    await setViewport(page, 375, 812)
    const session = await login(page, identities.cashier)
    expect(session.role).toBe('cashier')
    await openFinancialPage(page, 'ar', 'accounting/chart-of-accounts')
    await expect(page.getByRole('heading', { name: /غير مصرح|Access denied/i })).toBeVisible()
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
