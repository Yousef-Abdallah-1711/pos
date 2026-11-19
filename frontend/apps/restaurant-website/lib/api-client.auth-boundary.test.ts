import type { AxiosAdapter, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { apiClient } from './api-client'
import { rememberOrderToken } from './order-access'

describe('public storefront API credential boundary', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('removes legacy staff bearer storage and sends only the order capability', async () => {
    rememberOrderToken({ id: 7, tracking_token: 'synthetic-order-capability' })
    localStorage.setItem('auth-token', 'legacy-staff-bearer')
    localStorage.setItem('tenant-id', '12')
    localStorage.setItem('restaurant-id', '34')

    let captured: InternalAxiosRequestConfig | undefined
    const adapter: AxiosAdapter = async (config) => {
      captured = config
      return {
        data: { success: true },
        status: 200,
        statusText: 'OK',
        headers: {},
        config,
      } satisfies AxiosResponse
    }

    await apiClient.get('/public/orders/7/track', {
      adapter,
      headers: { Authorization: 'Bearer explicitly-supplied-staff-token' },
    })

    expect(captured).toBeDefined()
    expect(captured?.headers.get('Authorization')).toBeFalsy()
    expect(captured?.headers.get('X-Order-Token')).toBe('synthetic-order-capability')
    expect(captured?.headers.get('X-Tenant-ID')).toBe('12')
    expect(captured?.headers.get('X-Restaurant-ID')).toBe('34')
    expect(localStorage.getItem('auth-token')).toBeNull()
  })
})
