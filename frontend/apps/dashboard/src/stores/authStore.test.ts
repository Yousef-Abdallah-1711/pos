import { afterEach, describe, expect, it } from 'vitest'
import { apiCache } from '@/lib/apiCache'
import { useKDSStore } from './kdsStore'
import { usePOSStore } from './posStore'
import { usePreferencesStore } from './preferencesStore'
import { useThemeBuilderStore } from './themeBuilderStore'
import { useAuthStore } from './authStore'

const tenantAUser = {
  id: 41,
  email: 'owner-a@example.test',
  name: 'Owner A',
  role: 'owner',
  tenant_id: 7,
  restaurant_id: 12,
}

function seedTenantAClientStores() {
  usePOSStore.getState().setCustomer('Synthetic Tenant A Customer', '+10000000000')
  usePOSStore.getState().addItem({ id: 12, name: 'Tenant A item', price: 10 })
  usePOSStore.getState().holdCurrentOrder()
  usePOSStore.getState().setCustomer('Active Synthetic Customer', '+10000000001')
  usePOSStore.getState().addItem({ id: 13, name: 'Active Tenant A item', price: 12 })
  usePOSStore.getState().setAutoPrintEnabled(true)
  usePOSStore.getState().setSoundEnabled(false)
  usePreferencesStore.getState().setSettings({
    id: 1,
    userId: tenantAUser.id,
    language: 'en',
    theme: 'dark',
    primaryColor: '#123456',
    secondaryColor: '#234567',
    accentColor: '#345678',
    emailNotifications: true,
    pushNotifications: true,
    smsNotifications: false,
  })
  useThemeBuilderStore.getState().updateIdentity({ siteTitle: 'Tenant A Private Theme' })
  useKDSStore.getState().setOrders([{
    id: 1,
    tenant_id: 7,
    restaurant_id: 12,
    order_id: 1,
    order_number: 'A-1',
    customer_name: 'Tenant A Customer',
    status: 'new',
    priority: 0,
    station_id: 'main',
    items: [],
    received_at: '2026-11-19T00:00:00.000Z',
    elapsed_seconds: 0,
    created_at: '2026-11-19T00:00:00.000Z',
    updated_at: '2026-11-19T00:00:00.000Z',
  }])
  useKDSStore.getState().setStats({
    active_orders: 1,
    avg_prep_minutes: 0,
    completed_today: 0,
    new_orders: 1,
    in_progress: 0,
    ready: 0,
  })
  localStorage.setItem('kds-bootstrap-snapshot-v1', '{"synthetic":"tenant-a"}')
  localStorage.setItem('rbac-permission-cache-v1', '{"synthetic":"tenant-a"}')
  localStorage.setItem('pos-products-cache', '{"synthetic":"tenant-a"}')
}

function expectTenantAClientStateCleared() {
  const pos = usePOSStore.getState()
  expect(pos.items).toEqual([])
  expect(pos.heldOrders).toEqual([])
  expect(pos.customerName).toBe('')
  expect(pos.customerPhone).toBe('')
  expect(pos.autoPrintEnabled).toBe(true)
  expect(pos.soundEnabled).toBe(false)
  expect(usePreferencesStore.getState().settings).toBeNull()
  expect(useThemeBuilderStore.getState().identity.siteTitle).toBe('My Restaurant')
  expect(useKDSStore.getState().orders).toEqual([])
  expect(useKDSStore.getState().stats).toBeNull()
  expect(localStorage.getItem('kds-bootstrap-snapshot-v1')).toBeNull()
  expect(localStorage.getItem('rbac-permission-cache-v1')).toBeNull()
  expect(localStorage.getItem('pos-products-cache')).toBeNull()
}

afterEach(() => {
  useAuthStore.getState().logout()
  apiCache.clearAll()
  localStorage.removeItem('auth-storage')
})

describe('auth store cache isolation', () => {
  it('clears cached API responses when the authenticated scope changes', () => {
    useAuthStore.getState().setAuth(tenantAUser, '')
    apiCache.set('categories', [{ id: 12, name: 'Tenant A category' }])

    useAuthStore.getState().setAuth({
      ...tenantAUser,
      id: 52,
      email: 'owner-b@example.test',
      name: 'Owner B',
      tenant_id: 8,
      restaurant_id: 19,
    }, '')

    expect(apiCache.get('categories')).toBeNull()
  })

  it('keeps the cache during session revalidation for the same identity', () => {
    useAuthStore.getState().setAuth(tenantAUser, '')
    apiCache.set('categories', [{ id: 12, name: 'Tenant A category' }])

    useAuthStore.getState().setAuth({ ...tenantAUser, name: 'Owner A Updated' }, '')

    expect(apiCache.get('categories')).toEqual([{ id: 12, name: 'Tenant A category' }])
  })

  it('clears cached API responses on logout', () => {
    useAuthStore.getState().setAuth(tenantAUser, '')
    apiCache.set('categories', [{ id: 12, name: 'Tenant A category' }])

    useAuthStore.getState().logout()

    expect(apiCache.get('categories')).toBeNull()
  })

  it('clears tenant-owned client state while preserving workstation settings', () => {
    useAuthStore.getState().setAuth(tenantAUser, '')
    seedTenantAClientStores()

    useAuthStore.getState().setAuth({
      ...tenantAUser,
      id: 52,
      email: 'owner-b@example.test',
      name: 'Owner B',
      tenant_id: 8,
      restaurant_id: 19,
    }, '')

    expectTenantAClientStateCleared()
  })
})
