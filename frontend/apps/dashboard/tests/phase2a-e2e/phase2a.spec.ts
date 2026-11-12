import { expect, request, test, type APIRequestContext, type Page } from '@playwright/test'

const api = process.env.WS05_API_URL || 'http://localhost:18080/api/v1'
const storefront = process.env.WS05_STOREFRONT_URL || 'http://phase2a-restaurant-a.localhost:13103'
const password = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'

async function login(page: Page, tenant: 'a' | 'b' = 'a') {
  const email = `phase2a-owner-${tenant}@example.test`
  const response = await page.request.post(`${api}/auth/login`, {
    data: { email, password },
    headers: { Origin: 'http://localhost:13102' },
  })
  expect(response.ok(), `login failed for ${email}`).toBeTruthy()
  return await response.json() as { user: { id: number; tenant_id: number; restaurant_id: number } }
}

test('tenant A cannot see or mutate tenant B products', async ({ page }) => {
  const tenantA = await request.newContext()
  const tenantB = await request.newContext()
  try {
    await loginWithContext(tenantA, 'a')
    const tenantBUser = await loginWithContext(tenantB, 'b')
    const tenantBProducts = await tenantB.get(`${api}/products`)
    const tenantBProduct = (await tenantBProducts.json()).products.find((product: { name_en: string }) => product.name_en === 'WS05 E2E Burger B')
    expect(tenantBProduct).toBeTruthy()

    const listing = await tenantA.get(`${api}/products`, {
      headers: { 'X-Tenant-ID': String(tenantBUser.tenant_id), 'X-Restaurant-ID': String(tenantBUser.restaurant_id) },
    })
    expect(listing.ok()).toBeTruthy()
    expect(await listing.text()).not.toContain('WS05 E2E Burger B')
    const deletion = await tenantA.delete(`${api}/products/${tenantBProduct.id}`, {
      headers: { Origin: 'http://localhost:13102' },
    })
    expect([403, 404]).toContain(deletion.status())
    const stillExists = await tenantB.get(`${api}/products`)
    expect(await stillExists.text()).toContain('WS05 E2E Burger B')
  } finally {
    await tenantA.dispose()
    await tenantB.dispose()
  }
})

async function loginWithContext(context: APIRequestContext, tenant: 'a' | 'b') {
  const response = await context.post(`${api}/auth/login`, {
    data: { email: `phase2a-owner-${tenant}@example.test`, password },
    headers: { Origin: 'http://localhost:13102' },
  })
  expect(response.ok()).toBeTruthy()
  const auth = await response.json() as { user: { tenant_id: number; restaurant_id: number } }
  return auth.user
}

test('POS browser flow creates an order for the signed-in tenant', async ({ page }) => {
  test.setTimeout(120_000)
  page.on('response', (response) => {
    const pathname = new URL(response.url()).pathname
    if (pathname === '/api/v1/pos/orders' || pathname === '/api/v1/admin/orders') {
      console.log(`POS_ORDER_RESPONSE ${pathname} ${response.status()}`)
    }
  })
  page.on('console', (message) => {
    if (message.type() === 'error') console.log(`POS_BROWSER_ERROR ${message.text().split('\n')[0]}`)
  })
  await login(page)
  await page.goto('/en/dashboard/pos', { waitUntil: 'domcontentloaded', timeout: 90_000 })
  await page.getByRole('button', { name: /WS05 E2E Burger A/ }).click()
  await page.getByRole('button', { name: /Place Order/ }).click()
  await expect(page.getByText(/Order placed!/)).toBeVisible({ timeout: 15_000 })
})

test('inventory page loads and stock adjustment persists in PostgreSQL', async ({ page }) => {
  await login(page)
  const initial = await page.request.get(`${api}/products`)
  const initialBody = await initial.json()
  const fixtureProduct = initialBody.products.find((product: { name_en: string }) => product.name_en === 'WS05 E2E Burger A')
  expect(fixtureProduct).toBeTruthy()
  const initialStock = fixtureProduct.quantity_in_stock
  await page.goto('/en/dashboard/inventory')
  const productRow = page.getByRole('row', { name: /WS05 E2E Burger A/ })
  await expect(productRow).toContainText('WS05 E2E Burger A')
  await expect(page.getByText('WS05 E2E Burger B', { exact: true })).toHaveCount(0)
  await productRow.getByRole('button', { name: 'Adjust Stock' }).click()
  await page.getByRole('spinbutton').fill('1')
  await page.getByRole('textbox', { name: 'Add any additional notes...' }).fill('Automated WS05 test')
  await page.getByRole('button', { name: 'Adjust Inventory' }).click()
  await expect(page.getByRole('dialog')).toHaveCount(0)

  const products = await page.request.get(`${api}/products`)
  const body = await products.json()
  expect(body.products.find((product: { id: number }) => product.id === fixtureProduct.id).quantity_in_stock).toBe(initialStock + 1)
})

test('concurrent inventory adjustments serialize and rejected adjustments have no side effects', async ({ page }) => {
  await login(page)
  const headers = { Origin: 'http://localhost:13102' }
  const productResponse = await page.request.get(`${api}/products`)
  expect(productResponse.ok()).toBeTruthy()
  const products = await productResponse.json()
  const product = products.products.find((item: { name_en: string }) => item.name_en === 'WS05 E2E Burger A')
  expect(product).toBeTruthy()

  const rejected = await page.request.post(`${api}/inventory/adjust`, {
    headers,
    data: { product_id: product.id, quantity_change: -100_000, reason: 'adjustment', notes: 'WS05 rejected adjustment' },
  })
  expect(rejected.status()).toBe(409)
  const afterRejectedResponse = await page.request.get(`${api}/products`)
  const afterRejected = await afterRejectedResponse.json()
  expect(afterRejected.products.find((item: { id: number }) => item.id === product.id).quantity_in_stock).toBe(product.quantity_in_stock)

  const adjustmentCount = 6
  const adjustments = await Promise.all(Array.from({ length: adjustmentCount }, () => page.request.post(`${api}/inventory/adjust`, {
    headers,
    data: { product_id: product.id, quantity_change: 1, reason: 'adjustment', notes: 'WS05 concurrent adjustment' },
  })))
  expect(adjustments.map((response) => response.status())).toEqual(Array(adjustmentCount).fill(200))

  const finalResponse = await page.request.get(`${api}/products`, { headers })
  const finalProducts = await finalResponse.json()
  expect(finalProducts.products.find((item: { id: number }) => item.id === product.id).quantity_in_stock).toBe(product.quantity_in_stock + adjustmentCount)
})

test('storefront order endpoint accepts a scoped public order', async ({ page }) => {
  test.setTimeout(120_000)
  const menuResponse = await page.request.get(`${api}/public/restaurants/phase2a-restaurant-a/menu`)
  expect(menuResponse.ok()).toBeTruthy()
  const menu = await menuResponse.json()
  const product = menu.products.find((item: { name_en: string }) => item.name_en === 'WS05 E2E Burger A')
  expect(product).toBeTruthy()
  await page.goto(`${storefront}/en`, { waitUntil: 'domcontentloaded', timeout: 90_000 })
  await expect(page.getByRole('heading', { name: 'WS05 Tenant A Restaurant', level: 1 })).toBeVisible()
  const orderResponse = await page.request.post(`${api}/public/orders`, {
    headers: { 'X-Tenant-ID': String(product.tenant_id), 'X-Restaurant-ID': String(product.restaurant_id) },
    data: {
      customer_name: 'WS05 Browser Test',
      customer_phone: '5550100300',
      delivery_address: '42 Test Avenue',
      payment_method: 'cash',
      order_source: 'website',
      items: [{ product_id: product.id, quantity: 1, special_instructions: '' }],
    },
  })
  expect(orderResponse.status()).toBe(201)
  expect((await orderResponse.json()).data.order_number).toBeTruthy()
})

test('super-admin login workflow is available in the isolated stack', async ({ page }) => {
  test.setTimeout(120_000)
  page.on('response', (response) => {
    const pathname = new URL(response.url()).pathname
    if (pathname.includes('/super-admin/auth/')) console.log(`SUPER_ADMIN_AUTH_RESPONSE ${pathname} ${response.status()}`)
  })
  page.on('console', (message) => {
    if (message.type() === 'error') console.log(`SUPER_ADMIN_BROWSER_ERROR ${message.text().split('\n')[0]}`)
  })
  await page.goto(`${process.env.WS05_SUPER_ADMIN_URL || 'http://localhost:13104'}/login`, { waitUntil: 'load', timeout: 110_000 })
  await expect(page.getByRole('textbox', { name: /email/i })).toBeVisible()
  await expect(page.getByLabel(/password/i)).toBeVisible()
  const superAdminEmail = process.env.E2E_SUPERADMIN_EMAIL || 'phase2a-superadmin@example.test'
  const superAdminPassword = process.env.E2E_SUPERADMIN_PASSWORD || 'WS05-SuperAdmin-Only-2026!'
  await page.getByRole('textbox', { name: /email/i }).fill(superAdminEmail)
  await page.getByLabel(/password/i).fill(superAdminPassword)
  await page.getByRole('button', { name: /^login$/i }).click()
  await expect(page).toHaveURL(/\/dashboard(?:$|\?)/, { timeout: 15_000 })
  await expect(page.getByRole('heading', { name: 'Welcome back!' })).toBeVisible()

  await page.getByRole('link', { name: 'Tenants', exact: true }).click()
  await expect(page.getByRole('heading', { name: 'Tenants', level: 1 })).toBeVisible()
  await expect(page.getByText(/phase2a-restaurant-a/)).toBeVisible()
  await expect(page.getByText(/phase2a-restaurant-b/)).toBeVisible()
})

test('super-admin tenant dashboard access fails closed until actor attribution is supported', async ({ page }) => {
  const apiBase = process.env.WS05_API_URL || 'http://localhost:18080/api/v1'
  const superAdminOrigin = process.env.WS05_SUPER_ADMIN_URL || 'http://localhost:3004'
  const loginResponse = await page.request.post(`${apiBase}/super-admin/auth/login`, {
    headers: { Origin: superAdminOrigin },
    data: {
      email: process.env.E2E_SUPERADMIN_EMAIL || 'phase2a-superadmin@example.test',
      password: process.env.E2E_SUPERADMIN_PASSWORD || 'WS05-SuperAdmin-Only-2026!',
    },
  })
  expect(loginResponse.ok()).toBeTruthy()
  const superAdminCookie = (await page.context().cookies('http://localhost:18080')).find((cookie) => cookie.name === 'poos_super_admin_session')
  expect(superAdminCookie?.httpOnly).toBe(true)
  const superAdminSession = await page.request.get(`${apiBase}/super-admin/auth/session`)
  expect(superAdminSession.ok()).toBeTruthy()
  await page.goto(`${superAdminOrigin}/dashboard/tenants`, { waitUntil: 'domcontentloaded' })
  await expect(page.getByRole('heading', { name: 'Tenants', level: 1 })).toBeVisible({ timeout: 30_000 })
  await page.getByRole('link', { name: 'Tenants', exact: true }).click()

  const tenantRow = page.getByRole('row').filter({ hasText: 'phase2a-restaurant-a' })
  await expect(tenantRow).toBeVisible()
  await tenantRow.getByRole('link', { name: 'View' }).click()
  await expect(page.getByRole('heading', { name: 'Tenant A Restaurant' })).toBeVisible()
  const selectedTenantId = new URL(page.url()).pathname.split('/').pop()
  if (!selectedTenantId || !/^\d+$/.test(selectedTenantId)) throw new Error('Selected tenant route must include a numeric ID')
  await expect(page.getByRole('button', { name: 'Access Paused' })).toBeDisabled()
  await expect(page.getByRole('button', { name: 'Access Paused' })).toHaveAttribute('title', /actor-attributed support sessions/)

  const tenantCookieBefore = (await page.context().cookies('http://localhost:18080')).find((cookie) => cookie.name === 'poos_session')
  expect(tenantCookieBefore).toBeUndefined()
  const refusedHandoff = await page.request.post(`${apiBase}/super-admin/tenants/${selectedTenantId}/dashboard-login`, {
    headers: { Origin: superAdminOrigin },
  })
  expect(refusedHandoff.status()).toBe(503)
  expect((await refusedHandoff.json()).error).toContain('actor-attributed support sessions')
  expect(refusedHandoff.headers()['set-cookie'] || '').not.toContain('poos_session=')
  const tenantCookieAfter = (await page.context().cookies('http://localhost:18080')).find((cookie) => cookie.name === 'poos_session')
  expect(tenantCookieAfter).toBeUndefined()
  const superAdminSessionAfter = await page.request.get(`${apiBase}/super-admin/auth/session`)
  expect(superAdminSessionAfter.ok()).toBeTruthy()
  await expect(page).toHaveURL(new RegExp(`/dashboard/tenants/${selectedTenantId}$`))
  await expect(page.getByRole('heading', { name: 'Tenant A Restaurant' })).toBeVisible()
})
