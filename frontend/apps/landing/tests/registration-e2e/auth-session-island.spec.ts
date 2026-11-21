import { expect, test, type BrowserContext } from '@playwright/test'

const landingURL = (process.env.REG_E2E_LANDING_URL ?? 'http://localhost:14211').replace(/\/$/, '')
const apiOrigin = process.env.REG_E2E_API_URL ?? 'http://localhost:19482/api/v1'
const cookieValue = 'A'.repeat(43)
const profile = {
  id: 713,
  is_super_admin: false,
  name: 'Synthetic Owner',
  email: 'synthetic-owner@example.test',
  restaurant_name: 'Synthetic Kitchen',
}

const corsHeaders = {
  'access-control-allow-origin': landingURL,
  'access-control-allow-credentials': 'true',
  'cache-control': 'private, no-store',
  'content-type': 'application/json',
}

async function addSyntheticSessionCookie(context: BrowserContext) {
  await context.addCookies([{
    name: 'poos_session',
    value: cookieValue,
    url: new URL(apiOrigin).origin,
    httpOnly: true,
    secure: false,
    sameSite: 'Lax',
  }])
}

test('restores the server session in Landing and a new tab, then propagates confirmed logout', async ({ page, context }) => {
  test.setTimeout(120_000)
  await page.setViewportSize({ width: 1440, height: 900 })
  await addSyntheticSessionCookie(context)

  let authenticated = true
  let sessionRequests = 0
  let logoutRequests = 0
  const observedRequests: Array<{ url: string; cookiePresent: boolean }> = []
  await context.route('**/api/v1/auth/session', async (route) => {
    sessionRequests += 1
    const headers = await route.request().allHeaders()
    observedRequests.push({ url: route.request().url(), cookiePresent: Boolean(headers.cookie?.includes(`poos_session=${cookieValue}`)) })
    await route.fulfill({
      status: authenticated ? 200 : 401,
      headers: corsHeaders,
      body: authenticated ? JSON.stringify(profile) : '',
    })
  })
  await context.route('**/api/v1/auth/logout', async (route) => {
    logoutRequests += 1
    authenticated = false
    await route.fulfill({
      status: 204,
      headers: {
        ...corsHeaders,
        'set-cookie': 'poos_session=; Max-Age=0; Path=/; HttpOnly; SameSite=Lax',
      },
    })
  })

  const response = await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  const sharedHTML = await response?.text() ?? ''
  expect(sharedHTML).not.toContain(profile.email)
  expect(sharedHTML).not.toContain(profile.name)

  const accountButton = page.getByRole('button', { name: 'Account menu for Synthetic Owner' })
  await expect(accountButton).toBeVisible()
  expect(observedRequests[0]).toEqual({ url: `${apiOrigin}/auth/session`, cookiePresent: true })
  await accountButton.click()
  const desktopPanel = page.locator('#sv-account-panel')
  await expect(desktopPanel.getByText(profile.email, { exact: true })).toBeVisible()
  await expect(desktopPanel.getByText(profile.restaurant_name, { exact: true })).toBeVisible()
  await expect(page.getByRole('link', { name: 'Dashboard' })).toHaveAttribute('href', /\/en\/dashboard$/)
  await expect(page.getByRole('link', { name: 'Account settings' })).toHaveAttribute('href', /\/en\/dashboard\/settings\/account$/)

  const desktopCapture = process.env.AUTH_E2E_DESKTOP_CAPTURE
  if (desktopCapture) {
    await page.addStyleTag({ content: 'nextjs-portal { display: none !important; }' })
    await page.screenshot({ path: desktopCapture, animations: 'disabled' })
  }

  const secondTab = await context.newPage()
  await secondTab.setViewportSize({ width: 390, height: 844 })
  await secondTab.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  await expect(secondTab.getByRole('button', { name: 'Open menu' })).toBeVisible()
  await secondTab.getByRole('button', { name: 'Open menu' }).click()
  const mobileDialog = secondTab.getByRole('dialog')
  await expect(mobileDialog.getByText(profile.email, { exact: true })).toBeVisible()
  await expect(mobileDialog.getByRole('link', { name: 'Dashboard' })).toBeVisible()
  await expect(mobileDialog.getByRole('link', { name: /^Log in$/i })).toHaveCount(0)

  await page.getByRole('button', { name: 'Sign out' }).click()
  const desktopAccountActions = page.locator('.sv-desktop-account-actions')
  await expect(desktopAccountActions.getByRole('link', { name: /^Log in$/i })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Account menu for Synthetic Owner' })).toHaveCount(0)
  await expect(mobileDialog.getByRole('link', { name: /^Log in$/i })).toBeVisible()
  await expect.poll(async () => (await context.cookies(apiOrigin)).some(cookie => cookie.name === 'poos_session')).toBe(false)
  expect(logoutRequests).toBe(1)
  expect(sessionRequests).toBeGreaterThanOrEqual(2)
})

test('keeps a temporary API outage distinct from an expired session and retries safely', async ({ page }) => {
  let status = 503
  await page.route('**/api/v1/auth/session', async route => route.fulfill({
    status,
    headers: corsHeaders,
    body: status === 200 ? JSON.stringify(profile) : '',
  }))

  await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  const unavailable = page.getByRole('status')
  await expect(unavailable).toContainText('Sign-in status is temporarily unavailable.')
  const desktopAccountActions = page.locator('.sv-desktop-account-actions')
  await expect(desktopAccountActions.getByRole('link', { name: /^Log in$/i })).toBeVisible()
  await expect(desktopAccountActions.getByRole('link', { name: /^Sign up$/i })).toBeVisible()
  await expect(desktopAccountActions.getByRole('button', { name: 'Retry' })).toBeVisible()

  status = 401
  await page.getByRole('button', { name: 'Retry' }).click()
  await expect(desktopAccountActions.getByRole('link', { name: /^Log in$/i })).toBeVisible()
  await expect(desktopAccountActions.getByRole('link', { name: /^Sign up$/i })).toBeVisible()
  await expect(desktopAccountActions.getByRole('link', { name: /Start your 14-day free trial/i })).toBeVisible()
})

test('keeps safe account links during a DNS failure and does not retry until requested', async ({ page }) => {
  const hydrationWarnings: string[] = []
  page.on('console', message => {
    if (/hydration|did not match|server rendered html/i.test(message.text())) hydrationWarnings.push(message.text())
  })
  page.on('pageerror', error => {
    if (/hydration|did not match|server rendered html/i.test(error.message)) hydrationWarnings.push(error.message)
  })

  let sessionRequests = 0
  let shouldFail = true
  await page.route('**/api/v1/auth/session', async route => {
    sessionRequests += 1
    if (shouldFail) {
      await route.abort('namenotresolved')
      return
    }
    await route.fulfill({ status: 200, headers: corsHeaders, body: JSON.stringify(profile) })
  })

  await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  await expect(page.getByRole('status')).toContainText('Sign-in status is temporarily unavailable.')
  const accountActions = page.locator('.sv-desktop-account-actions')
  await expect(accountActions.getByRole('link', { name: /^Log in$/i })).toBeVisible()
  await expect(accountActions.getByRole('link', { name: /^Sign up$/i })).toBeVisible()
  await page.waitForTimeout(1000)
  expect(sessionRequests).toBe(1)

  shouldFail = false
  await accountActions.getByRole('button', { name: 'Retry' }).click()
  await expect(page.getByRole('button', { name: 'Account menu for Synthetic Owner' })).toBeVisible()
  expect(hydrationWarnings).toEqual([])

  await page.setViewportSize({ width: 390, height: 844 })
  shouldFail = true
  await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  await page.getByRole('button', { name: 'Open menu' }).click()
  const mobileDialog = page.getByRole('dialog')
  await expect(mobileDialog.getByRole('status')).toContainText('Sign-in status is temporarily unavailable.')
  await expect(mobileDialog.getByRole('link', { name: /^Log in$/i })).toBeVisible()
  await expect(mobileDialog.getByRole('link', { name: /^Sign up$/i })).toBeVisible()
  shouldFail = false
  await mobileDialog.getByRole('button', { name: 'Retry' }).click()
  await expect(mobileDialog.getByText(profile.email, { exact: true })).toBeVisible()
})

test('shows localized login and registration routes for anonymous users on mobile and desktop', async ({ page }) => {
  await page.route('**/api/v1/auth/session', route => route.fulfill({ status: 401, headers: corsHeaders, body: '' }))

  for (const width of [320, 375, 390, 768, 1024, 1440]) {
    await page.setViewportSize({ width, height: 900 })
    for (const locale of ['en', 'ar'] as const) {
      await page.goto(`${landingURL}/${locale}/`, { waitUntil: 'domcontentloaded' })
      if (width < 1200) {
        await page.getByRole('button', { name: locale === 'en' ? 'Open menu' : 'فتح القائمة' }).click()
      }

      const entry = width < 1200 ? page.getByRole('dialog') : page.locator('.sv-desktop-account-actions')
      const loginLabel = locale === 'en' ? 'Log in' : 'تسجيل الدخول'
      const signUpLabel = locale === 'en' ? 'Sign up' : 'إنشاء حساب'
      await expect(entry.getByRole('link', { name: loginLabel, exact: true })).toHaveAttribute('href', `/${locale}/auth/login`)
      await expect(entry.getByRole('link', { name: signUpLabel, exact: true })).toHaveAttribute('href', `/${locale}/auth/register`)
      await expect(entry.getByRole('button', { name: /account menu|قائمة الحساب/i })).toHaveCount(0)
      await expect(page.getByRole('status')).toHaveCount(0)
    }
  }
})

test('preserves the account menu across mobile sizes, Arabic RTL, keyboard use, and theme changes', async ({ page, context }) => {
  test.setTimeout(120_000)
  await addSyntheticSessionCookie(context)
  await context.route('**/api/v1/auth/session', async route => {
    await route.fulfill({ status: 200, headers: corsHeaders, body: JSON.stringify(profile) })
  })

  for (const width of [320, 390, 768, 1024, 1440]) {
    await page.setViewportSize({ width, height: 900 })
    await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
    await expect.poll(() => page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)

    if (width < 1200) {
      const openMenu = page.getByRole('button', { name: 'Open menu' })
      await expect(openMenu).toBeVisible()
      await openMenu.click()
      const dialog = page.getByRole('dialog')
      await expect(dialog.getByText(profile.email, { exact: true })).toBeVisible()
      const signOut = dialog.getByRole('button', { name: 'Sign out' })
      const bounds = await signOut.boundingBox()
      expect(bounds?.height ?? 0).toBeGreaterThanOrEqual(44)
      if (width === 320) {
        const mobileCapture = process.env.AUTH_E2E_MOBILE_CAPTURE
        if (mobileCapture) {
          await page.addStyleTag({ content: 'nextjs-portal { display: none !important; }' })
          await page.screenshot({ path: mobileCapture, animations: 'disabled' })
        }
        await page.keyboard.press('Escape')
        await expect(dialog).toBeHidden()
        await expect(openMenu).toBeFocused()
      } else {
        await page.keyboard.press('Escape')
      }
    } else {
      const account = page.getByRole('button', { name: 'Account menu for Synthetic Owner' })
      await expect(account).toBeVisible()
      await account.click()
      await page.keyboard.press('Escape')
      await expect(account).toBeFocused()
    }
  }

  await page.setViewportSize({ width: 320, height: 844 })
  await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  await page.getByRole('button', { name: 'Open menu' }).click()
  await page.getByRole('button', { name: 'Switch to light mode' }).click()
  await expect(page.locator('html')).toHaveClass(/light/)
  await page.getByRole('button', { name: 'Switch to Arabic' }).click()
  await expect(page).toHaveURL(/\/ar\/?$/)
  await expect(page.locator('html')).toHaveAttribute('dir', 'rtl')
  await page.getByRole('button', { name: 'فتح القائمة' }).click()
  const arabicDialog = page.getByRole('dialog')
  await expect(arabicDialog.getByText(profile.email, { exact: true })).toBeVisible()
  await expect.poll(() => page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
})

test('does not claim logout when the server cannot confirm it', async ({ page, context }) => {
  await addSyntheticSessionCookie(context)
  await context.route('**/api/v1/auth/session', async route => route.fulfill({
    status: 200,
    headers: corsHeaders,
    body: JSON.stringify(profile),
  }))
  await context.route('**/api/v1/auth/logout', async route => route.fulfill({
    status: 503,
    headers: corsHeaders,
    body: '',
  }))

  await page.goto(`${landingURL}/en/`, { waitUntil: 'domcontentloaded' })
  await page.getByRole('button', { name: 'Account menu for Synthetic Owner' }).click()
  await page.getByRole('button', { name: 'Sign out' }).click()
  await expect(page.getByRole('status')).toContainText('We could not confirm sign out. Try again.')
  await expect(page.locator('.sv-desktop-account-actions').getByRole('link', { name: /^Log in$/i })).toHaveCount(0)
  await expect(page.getByRole('button', { name: /Account menu for Synthetic Owner/ })).toHaveCount(0)
  await expect.poll(async () => (await context.cookies(apiOrigin)).some(cookie => cookie.name === 'poos_session')).toBe(true)
})
