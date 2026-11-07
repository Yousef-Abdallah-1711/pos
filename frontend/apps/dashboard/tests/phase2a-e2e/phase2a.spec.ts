import { expect, request, test, type APIRequestContext, type Page } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const api = process.env.WS05_API_URL || 'http://localhost:18081/api/v1'
const storefront = process.env.WS05_STOREFRONT_URL || 'http://phase2a-restaurant-a.localhost:13103'
const password = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'

test.beforeEach(async ({ page }) => {
  await routeContainerApiToCandidate(page)
})

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

test('authenticated user CRUD is tenant-scoped and foreign IDs leave rows unchanged', async ({ page }) => {
  const tenantA = await request.newContext()
  const tenantB = await request.newContext()
  const employee = await request.newContext()
  try {
    await loginWithContext(tenantA, 'a')
    await loginWithContext(tenantB, 'b')

    const foreignBeforeResponse = await tenantB.get(`${api}/users`)
    expect(foreignBeforeResponse.ok()).toBeTruthy()
    const foreignBefore = await foreignBeforeResponse.json() as Array<{ id: number; name: string; email: string }>
    const foreignUser = foreignBefore.find((user) => user.email === 'phase2a-owner-b@example.test')
    expect(foreignUser).toBeTruthy()
    if (!foreignUser) throw new Error('Tenant B user fixture is missing')

    const email = `ws05-user-crud-${Date.now()}@example.test`
    const create = await tenantA.post(`${api}/users`, {
      headers: { Origin: 'http://localhost:13102' },
      data: { name: 'WS05 CRUD Original', email, phone: '5550100400', password: 'WS05-User-CRUD-2026!' },
    })
    expect(create.status()).toBe(200)
    const createBody = await create.json() as { data: { id: number; email: string } }
    expect(createBody.data.email).toBe(email)

    const tenantAList = await tenantA.get(`${api}/users`)
    expect(tenantAList.ok()).toBeTruthy()
    const tenantAUsers = await tenantAList.json() as Array<{ id: number; name: string; email: string }>
    const createdUser = tenantAUsers.find((user) => user.id === createBody.data.id)
    expect(createdUser?.name).toBe('WS05 CRUD Original')
    expect(tenantAUsers.some((user) => user.email === 'phase2a-owner-b@example.test')).toBe(false)

    const foreignUpdate = await tenantA.put(`${api}/users/${foreignUser.id}`, {
      headers: { Origin: 'http://localhost:13102' }, data: { name: 'Unauthorized Change', phone: '5550100499' },
    })
    expect(foreignUpdate.status()).toBe(404)
    const foreignDelete = await tenantA.delete(`${api}/users/${foreignUser.id}`, { headers: { Origin: 'http://localhost:13102' } })
    expect(foreignDelete.status()).toBe(404)
    const foreignAfterResponse = await tenantB.get(`${api}/users`)
    const foreignAfter = await foreignAfterResponse.json() as Array<{ id: number; name: string; email: string }>
    expect(foreignAfter.find((user) => user.id === foreignUser.id)).toMatchObject({ name: foreignUser.name, email: foreignUser.email })

    const update = await tenantA.put(`${api}/users/${createBody.data.id}`, {
      headers: { Origin: 'http://localhost:13102' }, data: { name: 'WS05 CRUD Updated', phone: '5550100401' },
    })
    expect(update.status()).toBe(200)
    const updatedList = await tenantA.get(`${api}/users`)
    const updatedUsers = await updatedList.json() as Array<{ id: number; name: string; email: string }>
    expect(updatedUsers.find((user) => user.id === createBody.data.id)?.name).toBe('WS05 CRUD Updated')

    const employeeLogin = await employee.post(`${api}/auth/login`, {
      headers: { Origin: 'http://localhost:13102' },
      data: { email, password: 'WS05-User-CRUD-2026!' },
    })
    expect(employeeLogin.ok()).toBeTruthy()
    const employeeAuth = await employeeLogin.json() as { user: { id: number } }
    const employeeList = await employee.get(`${api}/users`)
    expect(employeeList.status()).toBe(403)
    const employeeCreate = await employee.post(`${api}/users`, {
      headers: { Origin: 'http://localhost:13102' },
      data: { name: 'Unauthorized Employee', email: `ws05-employee-forbidden-${Date.now()}@example.test`, phone: '5550100402', password: 'WS05-Forbidden-2026!' },
    })
    expect(employeeCreate.status()).toBe(403)
    const employeeUpdate = await employee.put(`${api}/users/${employeeAuth.user.id}`, {
      headers: { Origin: 'http://localhost:13102' }, data: { name: 'Unauthorized Employee Update', phone: '5550100498' },
    })
    expect(employeeUpdate.status()).toBe(403)
    const employeeDelete = await employee.delete(`${api}/users/${employeeAuth.user.id}`, { headers: { Origin: 'http://localhost:13102' } })
    expect(employeeDelete.status()).toBe(403)
    const afterEmployeeAttempts = await tenantA.get(`${api}/users`)
    const afterEmployeeUsers = await afterEmployeeAttempts.json() as Array<{ id: number; name: string; email: string }>
    expect(afterEmployeeUsers.find((user) => user.id === employeeAuth.user.id)).toMatchObject({ name: 'WS05 CRUD Updated', email })

    const deletion = await tenantA.delete(`${api}/users/${createBody.data.id}`, { headers: { Origin: 'http://localhost:13102' } })
    expect(deletion.status()).toBe(200)
    const finalList = await tenantA.get(`${api}/users`)
    const finalUsers = await finalList.json() as Array<{ id: number }>
    expect(finalUsers.some((user) => user.id === createBody.data.id)).toBe(false)
  } finally {
    await employee.dispose()
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

test('storefront menu cart and checkout complete an order in the browser', async ({ page }) => {
  test.setTimeout(120_000)
  let legacyThemeLookupCount = 0
  page.on('request', (request) => {
    if (new URL(request.url()).pathname.startsWith('/api/v1/public/themes/')) legacyThemeLookupCount += 1
  })
  const storefrontUrl = process.env.WS05_STOREFRONT_URL || 'http://phase2a-restaurant-a.localhost:3003'
  const menuResponse = await page.request.get(`${api}/public/restaurants/phase2a-restaurant-a/menu`)
  expect(menuResponse.ok()).toBeTruthy()
  const menu = await menuResponse.json()
  const fixtureProduct = menu.products.find((item: { name_en: string }) => item.name_en === 'WS05 E2E Burger A')
  expect(fixtureProduct).toBeTruthy()

  await page.goto(`${storefrontUrl}/en/menu`, { waitUntil: 'domcontentloaded', timeout: 90_000 })
  const productTitle = page.getByRole('heading', { name: 'WS05 E2E Burger A' })
  await expect(productTitle).toBeVisible({ timeout: 30_000 })
  const productCard = productTitle.locator('xpath=../..')
  await productCard.getByRole('button', { name: 'Add', exact: true }).click()

  await page.goto(`${storefrontUrl}/en/cart`, { waitUntil: 'domcontentloaded' })
  await expect(page.getByRole('heading', { name: 'WS05 E2E Burger A' })).toBeVisible()
  await page.getByRole('link', { name: /proceed to checkout/i }).first().click()
  await expect(page).toHaveURL(/\/en\/checkout$/, { timeout: 30_000 })
  await page.locator('input[type="text"]').first().fill('WS05 Browser Checkout')
  await page.locator('input[type="tel"]').fill('5550100399')
  await page.locator('textarea').first().fill('42 Test Avenue')
  await page.getByRole('button', { name: /place order/i }).click()

  const orderNumber = page.locator('.checkout-success-order')
  await expect(orderNumber).toBeVisible({ timeout: 30_000 })
  const orderNumberText = (await orderNumber.textContent())?.trim()
  expect(orderNumberText).toBeTruthy()
  const lookup = await page.request.get(`${api}/public/orders/number/${orderNumberText}/lookup`, {
    headers: {
      'X-Tenant-ID': String(fixtureProduct.tenant_id),
      'X-Restaurant-ID': String(fixtureProduct.restaurant_id),
      'X-Customer-Phone': '5550100399',
    },
  })
  expect(lookup.ok()).toBeTruthy()
  const order = await lookup.json()
  expect(order.data?.customer_name || order.data?.customerName).toBe('WS05 Browser Checkout')
  expect(legacyThemeLookupCount).toBe(0)
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
  const apiBase = process.env.WS05_API_URL || 'http://localhost:18081/api/v1'
  const superAdminOrigin = process.env.WS05_SUPER_ADMIN_URL || 'http://localhost:13104'
  const loginResponse = await page.request.post(`${apiBase}/super-admin/auth/login`, {
    headers: { Origin: superAdminOrigin },
    data: {
      email: process.env.E2E_SUPERADMIN_EMAIL || 'phase2a-superadmin@example.test',
      password: process.env.E2E_SUPERADMIN_PASSWORD || 'WS05-SuperAdmin-Only-2026!',
    },
  })
  expect(loginResponse.ok()).toBeTruthy()
  const superAdminCookie = (await page.context().cookies('http://localhost:18081')).find((cookie) => cookie.name === 'poos_super_admin_session')
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

  const tenantCookieBefore = (await page.context().cookies('http://localhost:18081')).find((cookie) => cookie.name === 'poos_session')
  expect(tenantCookieBefore).toBeUndefined()
  const refusedHandoff = await page.request.post(`${apiBase}/super-admin/tenants/${selectedTenantId}/dashboard-login`, {
    headers: { Origin: superAdminOrigin },
  })
  expect(refusedHandoff.status()).toBe(503)
  expect((await refusedHandoff.json()).error).toContain('actor-attributed support sessions')
  expect(refusedHandoff.headers()['set-cookie'] || '').not.toContain('poos_session=')
  const tenantCookieAfter = (await page.context().cookies('http://localhost:18081')).find((cookie) => cookie.name === 'poos_session')
  expect(tenantCookieAfter).toBeUndefined()
  const superAdminSessionAfter = await page.request.get(`${apiBase}/super-admin/auth/session`)
  expect(superAdminSessionAfter.ok()).toBeTruthy()
  await expect(page).toHaveURL(new RegExp(`/dashboard/tenants/${selectedTenantId}$`))
  await expect(page.getByRole('heading', { name: 'Tenant A Restaurant' })).toBeVisible()
})
