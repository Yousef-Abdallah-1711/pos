import { expect, request, test, type APIRequestContext, type Page } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const api = process.env.WS05_API_URL || 'http://localhost:18082/api/v1'
const storefront = process.env.WS05_STOREFRONT_URL || 'http://phase2a-restaurant-a.localhost:13113'
const password = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'

test.beforeEach(async ({ page }) => {
  await routeContainerApiToCandidate(page)
})

async function login(page: Page, tenant: 'a' | 'b' = 'a') {
  const email = `phase2a-owner-${tenant}@example.test`
  const response = await page.request.post(`${api}/auth/login`, {
    data: { email, password },
    headers: { Origin: 'http://localhost:13112' },
  })
  expect(response.ok(), `login failed for ${email}`).toBeTruthy()
  return await response.json() as { user: { id: number; tenant_id: number; restaurant_id: number } }
}

test('tenant A cannot see or mutate tenant B products', async ({ page }) => {
  const tenantA = await request.newContext()
  const tenantB = await request.newContext()
  const employee = await request.newContext()
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

    const ownedBeforeResponse = await tenantA.get(`${api}/products`)
    expect(ownedBeforeResponse.ok()).toBeTruthy()
    const ownedBefore = await ownedBeforeResponse.json()
    const ownedProduct = ownedBefore.products.find((product: { name_en: string }) => product.name_en === 'WS05 E2E Burger A')
    expect(ownedProduct).toBeTruthy()
    if (!ownedProduct) throw new Error('Tenant A product fixture was not found')
    const categoriesBeforeResponse = await tenantA.get(`${api}/categories`)
    expect(categoriesBeforeResponse.ok()).toBeTruthy()
    const categoriesBefore = await categoriesBeforeResponse.json() as Array<{ id: number; name: string; is_active: boolean }>
    expect(categoriesBefore.length).toBeGreaterThan(0)
    const ownedCategory = categoriesBefore.find((category) => category.id === ownedProduct.category_id)
    expect(ownedCategory).toBeTruthy()
    if (!ownedCategory) throw new Error('Tenant A category fixture was not found')
    const employeeEmail = `ws05-product-employee-${Date.now()}@example.test`
    const employeeCreate = await tenantA.post(`${api}/users`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { name: 'WS05 Product Employee', email: employeeEmail, phone: '5550100403', password: 'WS05-Product-Employee-2026!' },
    })
    expect(employeeCreate.status()).toBe(200)
    const employeeLogin = await employee.post(`${api}/auth/login`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { email: employeeEmail, password: 'WS05-Product-Employee-2026!' },
    })
    expect(employeeLogin.ok()).toBeTruthy()

    const employeeListing = await employee.get(`${api}/products`)
    expect(employeeListing.ok()).toBeTruthy()
    expect(await employeeListing.text()).toContain('WS05 E2E Burger A')
    expect(await employeeListing.text()).not.toContain('WS05 E2E Burger B')
    expect((await employee.get(`${api}/products/${ownedProduct.id}`)).status()).toBe(200)
    const employeeCategoriesResponse = await employee.get(`${api}/categories`)
    expect(employeeCategoriesResponse.ok()).toBeTruthy()
    expect(await employeeCategoriesResponse.text()).toContain(ownedCategory.name)
    expect((await employee.post(`${api}/categories`, {
      headers: { Origin: 'http://localhost:13112' }, data: { name: 'WS05 Unauthorized Category', is_active: true },
    })).status()).toBe(403)
    expect((await employee.put(`${api}/categories/${ownedCategory.id}`, {
      headers: { Origin: 'http://localhost:13112' }, data: { name: 'WS05 Unauthorized Category Update' },
    })).status()).toBe(403)
    expect((await employee.delete(`${api}/categories/${ownedCategory.id}`, {
      headers: { Origin: 'http://localhost:13112' },
    })).status()).toBe(403)
    expect((await employee.post(`${api}/products`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { category_id: ownedProduct.category_id, name_en: 'WS05 Unauthorized Product', price: 1 },
    })).status()).toBe(403)
    expect((await employee.put(`${api}/products/${ownedProduct.id}`, {
      headers: { Origin: 'http://localhost:13112' }, data: { name_en: 'WS05 Unauthorized Product Update' },
    })).status()).toBe(403)
    expect((await employee.delete(`${api}/products/${ownedProduct.id}`, {
      headers: { Origin: 'http://localhost:13112' },
    })).status()).toBe(403)
    const ownedAfterResponse = await tenantA.get(`${api}/products`)
    const ownedAfter = await ownedAfterResponse.json()
    expect(ownedAfter.total).toBe(ownedBefore.total)
    expect(ownedAfter.products.find((product: { id: number }) => product.id === ownedProduct.id)).toMatchObject({
      name_en: 'WS05 E2E Burger A',
      quantity_in_stock: ownedProduct.quantity_in_stock,
    })
    const categoriesAfterResponse = await tenantA.get(`${api}/categories`)
    expect(categoriesAfterResponse.ok()).toBeTruthy()
    const categoriesAfter = await categoriesAfterResponse.json() as Array<{ id: number; name: string; is_active: boolean }>
    expect(categoriesAfter).toEqual(categoriesBefore)
    const deletion = await tenantA.delete(`${api}/products/${tenantBProduct.id}`, {
      headers: { Origin: 'http://localhost:13112' },
    })
    expect([403, 404]).toContain(deletion.status())
    const stillExists = await tenantB.get(`${api}/products`)
    expect(await stillExists.text()).toContain('WS05 E2E Burger B')
  } finally {
    await employee.dispose()
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
      headers: { Origin: 'http://localhost:13112' },
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
      headers: { Origin: 'http://localhost:13112' }, data: { name: 'Unauthorized Change', phone: '5550100499' },
    })
    expect(foreignUpdate.status()).toBe(404)
    const foreignDelete = await tenantA.delete(`${api}/users/${foreignUser.id}`, { headers: { Origin: 'http://localhost:13112' } })
    expect(foreignDelete.status()).toBe(404)
    const foreignAfterResponse = await tenantB.get(`${api}/users`)
    const foreignAfter = await foreignAfterResponse.json() as Array<{ id: number; name: string; email: string }>
    expect(foreignAfter.find((user) => user.id === foreignUser.id)).toMatchObject({ name: foreignUser.name, email: foreignUser.email })

    const update = await tenantA.put(`${api}/users/${createBody.data.id}`, {
      headers: { Origin: 'http://localhost:13112' }, data: { name: 'WS05 CRUD Updated', phone: '5550100401' },
    })
    expect(update.status()).toBe(200)
    const updatedList = await tenantA.get(`${api}/users`)
    const updatedUsers = await updatedList.json() as Array<{ id: number; name: string; email: string }>
    expect(updatedUsers.find((user) => user.id === createBody.data.id)?.name).toBe('WS05 CRUD Updated')

    const employeeLogin = await employee.post(`${api}/auth/login`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { email, password: 'WS05-User-CRUD-2026!' },
    })
    expect(employeeLogin.ok()).toBeTruthy()
    const employeeAuth = await employeeLogin.json() as { user: { id: number } }
    const employeeList = await employee.get(`${api}/users`)
    expect(employeeList.status()).toBe(403)
    const employeeCreate = await employee.post(`${api}/users`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { name: 'Unauthorized Employee', email: `ws05-employee-forbidden-${Date.now()}@example.test`, phone: '5550100402', password: 'WS05-Forbidden-2026!' },
    })
    expect(employeeCreate.status()).toBe(403)
    const employeeUpdate = await employee.put(`${api}/users/${employeeAuth.user.id}`, {
      headers: { Origin: 'http://localhost:13112' }, data: { name: 'Unauthorized Employee Update', phone: '5550100498' },
    })
    expect(employeeUpdate.status()).toBe(403)
    const employeeDelete = await employee.delete(`${api}/users/${employeeAuth.user.id}`, { headers: { Origin: 'http://localhost:13112' } })
    expect(employeeDelete.status()).toBe(403)
    const afterEmployeeAttempts = await tenantA.get(`${api}/users`)
    const afterEmployeeUsers = await afterEmployeeAttempts.json() as Array<{ id: number; name: string; email: string }>
    expect(afterEmployeeUsers.find((user) => user.id === employeeAuth.user.id)).toMatchObject({ name: 'WS05 CRUD Updated', email })

    const deletion = await tenantA.delete(`${api}/users/${createBody.data.id}`, { headers: { Origin: 'http://localhost:13112' } })
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

test('sibling restaurant owner can manage tenant-wide users without seeing foreign tenants', async () => {
  const tenantA = await request.newContext()
  const siblingA = await request.newContext()
  const tenantB = await request.newContext()
  try {
    const primaryAuth = await loginWithContext(tenantA, 'a')
    await loginWithContext(tenantB, 'b')
    const siblingLogin = await siblingA.post(`${api}/auth/login`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { email: 'phase2a-sibling-owner-a@example.test', password },
    })
    expect(siblingLogin.ok()).toBeTruthy()
    const siblingAuth = await siblingLogin.json() as { user: { tenant_id: number; restaurant_id: number } }
    expect(siblingAuth.user.tenant_id).toBe(primaryAuth.tenant_id)
    expect(siblingAuth.user.restaurant_id).not.toBe(primaryAuth.restaurant_id)

    const foreignUsersResponse = await tenantB.get(`${api}/users`)
    expect(foreignUsersResponse.ok()).toBeTruthy()
    const foreignUsers = await foreignUsersResponse.json() as Array<{ id: number; email: string }>
    const foreignOwner = foreignUsers.find((user) => user.email === 'phase2a-owner-b@example.test')
    expect(foreignOwner).toBeTruthy()
    if (!foreignOwner) throw new Error('Tenant B user fixture is missing')

    const email = `ws05-sibling-user-crud-${Date.now()}@example.test`
    const create = await tenantA.post(`${api}/users`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { name: 'Primary Restaurant User', email, phone: '5550100410', password: 'WS05-Sibling-CRUD-2026!' },
    })
    expect(create.status()).toBe(200)
    const created = await create.json() as { data: { id: number } }

    const siblingListResponse = await siblingA.get(`${api}/users`)
    expect(siblingListResponse.ok()).toBeTruthy()
    const siblingList = await siblingListResponse.json() as Array<{ id: number; email: string; name: string }>
    expect(siblingList.find((user) => user.id === created.data.id)?.email).toBe(email)
    expect(siblingList.some((user) => user.id === foreignOwner.id)).toBe(false)

    const update = await siblingA.put(`${api}/users/${created.data.id}`, {
      headers: { Origin: 'http://localhost:13112' },
      data: { name: 'Updated by Sibling Restaurant', phone: '5550100411' },
    })
    expect(update.status()).toBe(200)
    const afterUpdate = await siblingA.get(`${api}/users`)
    const updatedList = await afterUpdate.json() as Array<{ id: number; name: string }>
    expect(updatedList.find((user) => user.id === created.data.id)?.name).toBe('Updated by Sibling Restaurant')

    const deletion = await siblingA.delete(`${api}/users/${created.data.id}`, { headers: { Origin: 'http://localhost:13112' } })
    expect(deletion.status()).toBe(200)
    const primaryList = await tenantA.get(`${api}/users`)
    const remainingUsers = await primaryList.json() as Array<{ id: number }>
    expect(remainingUsers.some((user) => user.id === created.data.id)).toBe(false)
    const foreignAfter = await tenantB.get(`${api}/users`)
    const foreignAfterUsers = await foreignAfter.json() as Array<{ id: number; email: string }>
    expect(foreignAfterUsers.find((user) => user.id === foreignOwner.id)?.email).toBe(foreignOwner.email)
  } finally {
    await tenantA.dispose()
    await siblingA.dispose()
    await tenantB.dispose()
  }
})

async function loginWithContext(context: APIRequestContext, tenant: 'a' | 'b') {
  const response = await context.post(`${api}/auth/login`, {
    data: { email: `phase2a-owner-${tenant}@example.test`, password },
    headers: { Origin: 'http://localhost:13112' },
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
  const headers = { Origin: 'http://localhost:13112' }
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
  if (!orderNumberText) throw new Error('Checkout did not return an order number')
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
  const orderId = order.data?.id as number
  expect(orderId).toBeTruthy()
  const rememberedToken = await page.evaluate(({ number, id }) => {
    const store = JSON.parse(localStorage.getItem('order-access-tokens') || '{}')
    return { byId: store.byId?.[String(id)], byNumber: store.byNumber?.[number] }
  }, { number: orderNumberText, id: orderId })
  expect(rememberedToken.byNumber).toBeTruthy()
  expect(rememberedToken.byId).toBe(rememberedToken.byNumber)
  expect(page.url()).not.toContain(rememberedToken.byNumber)

  await page.goto(`${storefrontUrl}/en/orders/track/${encodeURIComponent(orderNumberText)}`, {
    waitUntil: 'domcontentloaded',
    timeout: 90_000,
  })
  await expect(page.getByRole('heading', { name: 'Order Status' })).toBeVisible({ timeout: 30_000 })
  await expect(page.getByRole('heading', { name: 'Order Details' })).toBeVisible()
  await expect(page.getByText('WS05 E2E Burger A')).toBeVisible()

  const orderPath = `${api}/public/orders/${orderId}`
  const scopedHeaders = {
    'X-Tenant-ID': String(fixtureProduct.tenant_id),
    'X-Restaurant-ID': String(fixtureProduct.restaurant_id),
  }
  const tenantBMenuResponse = await page.request.get(`${api}/public/restaurants/phase2a-restaurant-b/menu`)
  expect(tenantBMenuResponse.ok()).toBeTruthy()
  const tenantBMenu = await tenantBMenuResponse.json()
  const tenantBProduct = tenantBMenu.products.find((item: { name_en: string }) => item.name_en === 'WS05 E2E Burger B')
  expect(tenantBProduct).toBeTruthy()
  const foreignScopeHeaders = {
    'X-Tenant-ID': String(tenantBProduct.tenant_id),
    'X-Restaurant-ID': String(tenantBProduct.restaurant_id),
    'X-Order-Token': rememberedToken.byNumber,
  }
  const siblingRestaurantMenu = await page.request.get(`${api}/public/restaurants/phase2a-restaurant-a2/menu`)
  expect(siblingRestaurantMenu.ok()).toBeTruthy()
  const siblingRestaurant = (await siblingRestaurantMenu.json()).restaurant
  const siblingScopeHeaders = {
    'X-Tenant-ID': String(fixtureProduct.tenant_id),
    'X-Restaurant-ID': String(siblingRestaurant.id),
    'X-Order-Token': rememberedToken.byNumber,
  }
  const siblingScopeTrack = await page.request.get(`${orderPath}/track`, { headers: siblingScopeHeaders })
  expect(siblingScopeTrack.status()).toBe(404)
  const siblingScopeCancel = await page.request.delete(orderPath, {
    headers: siblingScopeHeaders,
    data: { reason: 'sibling restaurant public token scope test' },
  })
  expect(siblingScopeCancel.status()).toBe(404)
  const foreignScopeTrack = await page.request.get(`${orderPath}/track`, { headers: foreignScopeHeaders })
  expect(foreignScopeTrack.status()).toBe(404)
  const foreignScopeCancel = await page.request.delete(orderPath, {
    headers: foreignScopeHeaders,
    data: { reason: 'cross-tenant public token scope test' },
  })
  expect(foreignScopeCancel.status()).toBe(404)
  const noTokenTrack = await page.request.get(`${orderPath}/track`, { headers: scopedHeaders })
  expect(noTokenTrack.status()).toBe(404)
  const wrongTokenTrack = await page.request.get(`${orderPath}/track`, {
    headers: { ...scopedHeaders, 'X-Order-Token': 'incorrect-public-order-token' },
  })
  expect(wrongTokenTrack.status()).toBe(404)
  const validTrack = await page.request.get(`${orderPath}/track`, {
    headers: { ...scopedHeaders, 'X-Order-Token': rememberedToken.byNumber },
  })
  expect(validTrack.ok()).toBeTruthy()
  const validTrackBody = await validTrack.json()
  expect(validTrackBody.data.order_number).toBe(orderNumberText)

  const concurrentTokenAttempts = (await Promise.all(Array.from({ length: 8 }, (_, index) => {
    const foreignHeaders = index % 2 === 0 ? siblingScopeHeaders : foreignScopeHeaders
    return Promise.all([
      page.request.get(`${orderPath}/track`, { headers: foreignHeaders }),
      page.request.delete(orderPath, {
        headers: foreignHeaders,
        data: { reason: 'concurrent foreign-token isolation test' },
      }),
    ])
  }))).flat()
  for (const response of concurrentTokenAttempts) expect(response.status()).toBe(404)
  const afterConcurrentDenials = await page.request.get(`${orderPath}/track`, {
    headers: { ...scopedHeaders, 'X-Order-Token': rememberedToken.byNumber },
  })
  expect(afterConcurrentDenials.ok()).toBeTruthy()
  expect((await afterConcurrentDenials.json()).data.status).toBe(validTrackBody.data.status)

  const unauthorizedCancel = await page.request.delete(orderPath, {
    headers: { ...scopedHeaders, 'X-Order-Token': 'incorrect-public-order-token' },
    data: { reason: 'unauthorized token test' },
  })
  expect(unauthorizedCancel.status()).toBe(404)
  const authorizedCancel = await page.request.delete(orderPath, {
    headers: { ...scopedHeaders, 'X-Order-Token': rememberedToken.byNumber },
    data: { reason: 'public order token lifecycle test' },
  })
  expect(authorizedCancel.ok()).toBeTruthy()
  let trackingThrottleResponse: Awaited<ReturnType<typeof page.request.get>> | undefined
  for (let attempt = 0; attempt < 70; attempt += 1) {
    trackingThrottleResponse = await page.request.get(`${orderPath}/track`, {
      headers: { ...scopedHeaders, 'X-Order-Token': rememberedToken.byNumber },
    })
    if (trackingThrottleResponse.status() === 429) break
  }
  expect(trackingThrottleResponse?.status()).toBe(429)
  expect(trackingThrottleResponse?.headers()['retry-after']).toBeTruthy()

  let cancellationThrottleResponse: Awaited<ReturnType<typeof page.request.delete>> | undefined
  for (let attempt = 0; attempt < 25; attempt += 1) {
    cancellationThrottleResponse = await page.request.delete(orderPath, {
      headers: { ...scopedHeaders, 'X-Order-Token': rememberedToken.byNumber },
      data: { reason: 'rate limit verification' },
    })
    if (cancellationThrottleResponse.status() === 429) break
  }
  expect(cancellationThrottleResponse?.status()).toBe(429)
  expect(cancellationThrottleResponse?.headers()['retry-after']).toBeTruthy()

  const cancelledStatus = await page.request.get(`${orderPath}/status`, {
    headers: { ...scopedHeaders, 'X-Order-Token': rememberedToken.byNumber },
  })
  expect(cancelledStatus.ok()).toBeTruthy()
  expect((await cancelledStatus.json()).data.status).toBe('cancelled')
  expect(legacyThemeLookupCount).toBe(0)
})

test('super-admin login workflow is available in the isolated stack', async ({ page }) => {
  test.setTimeout(180_000)
  page.on('framenavigated', (frame) => {
    if (frame === page.mainFrame()) console.log(`SUPER_ADMIN_NAVIGATION ${new URL(frame.url()).pathname}`)
  })
  page.on('request', (request) => {
    const pathname = new URL(request.url()).pathname
    if (pathname === '/dashboard' || pathname.startsWith('/dashboard?')) console.log(`SUPER_ADMIN_DASHBOARD_REQUEST ${request.method()} ${pathname}`)
  })
  page.on('response', (response) => {
    const pathname = new URL(response.url()).pathname
    if (pathname === '/dashboard' || pathname.startsWith('/dashboard?')) console.log(`SUPER_ADMIN_DASHBOARD_RESPONSE ${response.status()} ${pathname}`)
  })
  page.on('pageerror', (error) => console.log(`SUPER_ADMIN_PAGE_ERROR ${error.name}: ${error.message}`))
  page.on('requestfailed', (request) => console.log(`SUPER_ADMIN_REQUEST_FAILED ${request.method()} ${new URL(request.url()).pathname} ${request.failure()?.errorText ?? 'unknown'}`))
  page.on('response', (response) => {
    const pathname = new URL(response.url()).pathname
    if (pathname.includes('/super-admin/auth/')) console.log(`SUPER_ADMIN_AUTH_RESPONSE ${pathname} ${response.status()}`)
  })
  page.on('console', (message) => {
    if (message.type() === 'error') console.log(`SUPER_ADMIN_BROWSER_ERROR ${message.text().split('\n')[0]}`)
  })
  await page.goto(`${process.env.WS05_SUPER_ADMIN_URL || 'http://localhost:13114'}/login`, { waitUntil: 'load', timeout: 110_000 })
  await expect(page.getByRole('textbox', { name: /email/i })).toBeVisible()
  await expect(page.getByLabel(/password/i)).toBeVisible()
  const superAdminEmail = process.env.E2E_SUPERADMIN_EMAIL || 'phase2a-superadmin@example.test'
  const superAdminPassword = process.env.E2E_SUPERADMIN_PASSWORD || 'WS05-SuperAdmin-Only-2026!'
  await page.getByRole('textbox', { name: /email/i }).fill(superAdminEmail)
  await page.getByLabel(/password/i).fill(superAdminPassword)
  await page.getByRole('button', { name: /^login$/i }).click()
  await expect(page).toHaveURL(/\/dashboard(?:$|\?)/, { timeout: 60_000 })
  await expect(page.getByRole('heading', { name: 'Welcome back!' })).toBeVisible()

  await page.getByRole('link', { name: 'Tenants', exact: true }).click()
  await expect(page.getByRole('heading', { name: 'Tenants', level: 1 })).toBeVisible()
  await expect(page.getByText(/phase2a-restaurant-a/)).toBeVisible()
  await expect(page.getByText(/phase2a-restaurant-b/)).toBeVisible()
})

test('super-admin tenant dashboard access fails closed until actor attribution is supported', async ({ page }) => {
  const apiBase = process.env.WS05_API_URL || 'http://localhost:18082/api/v1'
  const superAdminOrigin = process.env.WS05_SUPER_ADMIN_URL || 'http://localhost:13114'
  const loginResponse = await page.request.post(`${apiBase}/super-admin/auth/login`, {
    headers: { Origin: superAdminOrigin },
    data: {
      email: process.env.E2E_SUPERADMIN_EMAIL || 'phase2a-superadmin@example.test',
      password: process.env.E2E_SUPERADMIN_PASSWORD || 'WS05-SuperAdmin-Only-2026!',
    },
  })
  expect(loginResponse.ok()).toBeTruthy()
  const superAdminCookie = (await page.context().cookies('http://localhost:18082')).find((cookie) => cookie.name === 'poos_super_admin_session')
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

  const tenantCookieBefore = (await page.context().cookies('http://localhost:18082')).find((cookie) => cookie.name === 'poos_session')
  expect(tenantCookieBefore).toBeUndefined()
  const refusedHandoff = await page.request.post(`${apiBase}/super-admin/tenants/${selectedTenantId}/dashboard-login`, {
    headers: { Origin: superAdminOrigin },
  })
  expect(refusedHandoff.status()).toBe(503)
  expect((await refusedHandoff.json()).error).toContain('actor-attributed support sessions')
  expect(refusedHandoff.headers()['set-cookie'] || '').not.toContain('poos_session=')
  const tenantCookieAfter = (await page.context().cookies('http://localhost:18082')).find((cookie) => cookie.name === 'poos_session')
  expect(tenantCookieAfter).toBeUndefined()
  const superAdminSessionAfter = await page.request.get(`${apiBase}/super-admin/auth/session`)
  expect(superAdminSessionAfter.ok()).toBeTruthy()
  await expect(page).toHaveURL(new RegExp(`/dashboard/tenants/${selectedTenantId}$`))
  await expect(page.getByRole('heading', { name: 'Tenant A Restaurant' })).toBeVisible()
})
