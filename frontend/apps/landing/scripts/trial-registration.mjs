// Creates synthetic accounts only when explicitly enabled against the intended stack.
import { createRequire } from 'node:module'
import { randomBytes } from 'node:crypto'
import fs from 'node:fs/promises'
import path from 'node:path'
const { chromium } = createRequire(import.meta.url)('playwright')
const base = process.env.SERVIO_MARKETING_URL
const api = process.env.SERVIO_MARKETING_API_URL
if (process.env.SERVIO_SYNTHETIC_TRIAL !== '1' || !base || !api) throw new Error('Explicit synthetic trial flag and landing/API URLs are required')
const output = path.join(process.env.SERVIO_BROWSER_EVIDENCE || '/tmp/servio-trial-acceptance', 'trial-registration')
const checks = []
function check(name, passed, detail = '') {
 checks.push({ name, passed: !!passed, detail })
 if (!passed) throw new Error(`${name}: ${detail}`)
}
async function register(page, locale, billing, otherCategory) {
 const slug = `website-qa-${Date.now()}-${randomBytes(3).toString('hex')}`
 const password = randomBytes(18).toString('hex')
 await page.goto(`${base}/${locale}/auth/register?billing=${billing}`)
 await page.waitForFunction(() => {
  const input = document.querySelector('input[type=email]')
  return input && Object.keys(input).some(key => key.startsWith('__reactProps') && typeof input[key]?.onChange === 'function')
 })
 for (const [name, value] of Object.entries({ name: 'Synthetic Website Owner', email: `${slug}@example.test`, restaurantName: `Synthetic Website Acceptance ${slug}`, restaurantSlug: slug, password, confirmPassword: password })) await page.locator(`input[name=${name}]`).fill(value)
 await page.getByText(locale === 'ar' ? 'النطاق الفرعي متاح' : 'Subdomain is available', { exact: true }).waitFor()
 const [response] = await Promise.all([
  page.waitForResponse(r => r.request().method() === 'POST' && new URL(r.url()).pathname.endsWith('/auth/register')),
  page.locator('button[type=submit]').click(),
 ])
 check(`registration accepted ${locale}/${billing}`, response.status() === 201 || response.status() === 200, `HTTP ${response.status()}`)
 await page.waitForURL(`**/${locale}/dashboard**`, { timeout: 60000 })
 check(`owner reaches dashboard ${locale}/${billing}`, page.url().includes(`/${locale}/dashboard`))
 const subscription = await page.evaluate(async origin => {
  const response = await fetch(`${origin}/subscription/current`, { credentials: 'include' })
  const body = await response.json()
  return { status: response.status, data: body.data || body }
 }, api)
 check(`trial subscription readable ${locale}`, subscription.status === 200)
 check(`billing preference persisted ${locale}`, subscription.data.preferred_billing_period === billing)
 check(`trial status ${locale}`, subscription.data.subscription_status === 'trial' && subscription.data.can_access)
 const duration = new Date(subscription.data.trial_end_date).getTime() - Date.now()
 check(`14-day entitlement ${locale}`, duration > 14 * 86400000 - 180000 && duration <= 14 * 86400000)
 check(`checkout unavailable ${locale}`, subscription.data.checkout_available === false)
 await page.screenshot({ path: path.join(output, `${locale}-${billing}-dashboard.png`) })
 const category = await firstAction(page, locale, otherCategory)
 const loggedOut = await page.evaluate(async origin => (await fetch(`${origin}/auth/logout`, { method: 'POST', credentials: 'include' })).ok, api)
 check(`logout accepted ${locale}`, loggedOut)
 const afterLogout = await page.evaluate(async origin => (await fetch(`${origin}/auth/session`, { credentials: 'include' })).status, api)
 check(`logged-out access denied ${locale}`, afterLogout === 401)
 return category
}
async function firstAction(page, locale, otherCategory) {
 const dashboard = new URL(page.url()).origin
 const category = `Synthetic First Menu ${locale}`
 await page.goto(`${dashboard}/${locale}/dashboard/categories/new`)
 await page.waitForFunction(() => {
  const input = document.querySelector('input[name=name]')
  return input && Object.keys(input).some(key => key.startsWith('__reactProps') && typeof input[key]?.onChange === 'function')
 })
 await page.locator('input[name=name]').fill(category)
 const [response] = await Promise.all([
  page.waitForResponse(response => response.request().method() === 'POST' && new URL(response.url()).pathname.endsWith('/categories')),
  page.locator('button[type=submit]').click(),
 ])
 check(`first menu action accepted ${locale}`, response.ok())
 await page.waitForURL(`**/${locale}/dashboard/categories`)
 await page.getByRole('cell', { name: category, exact: true }).waitFor()
 check(`first menu action visible ${locale}`, true)
 if (otherCategory) check('new workspace menu isolation', await page.getByRole('cell', { name: otherCategory, exact: true }).count() === 0)
 await page.screenshot({ path: path.join(output, `${locale}-first-menu-action.png`) })
 return category
}
async function main() {
 await fs.mkdir(output, { recursive: true })
 const browser = await chromium.launch({ headless: true })
 try {
  let previousCategory
  for (const [locale, billing] of [['en', 'annual'], ['ar', 'monthly']]) {
   const context = await browser.newContext({ ignoreHTTPSErrors: process.env.SERVIO_DISPOSABLE_BROWSER === '1', viewport: { width: 1440, height: 1000 } })
   try { previousCategory = await register(await context.newPage(), locale, billing, previousCategory) } finally { await context.close() }
  }
 } finally {
  await browser.close()
  await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ checks }, null, 2))
 }
}
main().catch(error => { console.error(error.message); process.exitCode = 1 })
