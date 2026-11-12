import { describe, expect, it, vi } from 'vitest'
import {
  createNotificationCacheKey,
  createAuthenticatedScopeKey,
  dedupeInFlightRequest,
} from './notificationRequestCache'

describe('notification request cache keys', () => {
  it('keeps list cache entries separate for each tenant, restaurant, and user', () => {
    const key = createNotificationCacheKey(
      { tenantId: 2, restaurantId: 4, userId: 8 },
      {}
    )

    expect(key).not.toBe(createNotificationCacheKey({ tenantId: 3, restaurantId: 4, userId: 8 }, {}))
    expect(key).not.toBe(createNotificationCacheKey({ tenantId: 2, restaurantId: 5, userId: 8 }, {}))
    expect(key).not.toBe(createNotificationCacheKey({ tenantId: 2, restaurantId: 4, userId: 9 }, {}))
  })

  it('uses a stable scope key when optional identity fields are missing', () => {
    expect(createAuthenticatedScopeKey({ tenantId: null, restaurantId: undefined, userId: null }))
      .toBe('tenant:none:restaurant:none:user:none')
  })
})

describe('dedupeInFlightRequest', () => {
  it('shares concurrent requests with the same key and releases the key after settlement', async () => {
    let resolveRequest: ((value: string) => void) | undefined
    const request = vi.fn()
      .mockImplementationOnce(() => new Promise<string>((resolve) => { resolveRequest = resolve }))
      .mockResolvedValue('fresh')
    const first = dedupeInFlightRequest('user-scope:stats', request)
    const second = dedupeInFlightRequest('user-scope:stats', request)

    await Promise.resolve()
    expect(request).toHaveBeenCalledTimes(1)
    resolveRequest?.('fresh')
    await expect(Promise.all([first, second])).resolves.toEqual(['fresh', 'fresh'])

    await expect(dedupeInFlightRequest('user-scope:stats', request)).resolves.toBe('fresh')
    expect(request).toHaveBeenCalledTimes(2)
  })

  it('does not share requests across identity scopes', async () => {
    const request = vi.fn(async (value: string) => value)

    await expect(Promise.all([
      dedupeInFlightRequest('tenant:2:user:8:stats', () => request('tenant-2')),
      dedupeInFlightRequest('tenant:3:user:8:stats', () => request('tenant-3')),
    ])).resolves.toEqual(['tenant-2', 'tenant-3'])
    expect(request).toHaveBeenCalledTimes(2)
  })
})
