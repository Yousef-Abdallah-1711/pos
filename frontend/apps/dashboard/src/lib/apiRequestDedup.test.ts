import { describe, expect, it } from 'vitest'
import { buildDedupIdentityScope, buildScopedGetRequestDedupKey } from './apiRequestDedup'

describe('buildDedupIdentityScope', () => {
  it('requires a concrete user ID before a request can be coalesced', () => {
    expect(buildDedupIdentityScope(null)).toBeNull()
    expect(buildDedupIdentityScope({})).toBeNull()
    expect(buildDedupIdentityScope({ tenant_id: 2, role: 'owner' })).toBeNull()
    expect(buildDedupIdentityScope({ id: 8, tenant_id: 2, restaurant_id: 4, role: 'owner' })).toBe(
      '["8","2","4","owner",false]',
    )
  })
})

describe('buildScopedGetRequestDedupKey', () => {
  it('separates otherwise identical reads by authenticated identity and tenant headers', () => {
    const request = { method: 'get', url: '/orders', params: { page: 1 } }
    const tenantTwo = buildScopedGetRequestDedupKey(request, 'user:8:tenant:2:restaurant:4')
    const tenantThree = buildScopedGetRequestDedupKey(request, 'user:8:tenant:3:restaurant:4')
    const restaurantFive = buildScopedGetRequestDedupKey(
      { ...request, headers: { 'X-Tenant-ID': 2, 'X-Restaurant-ID': 5 } },
      'user:8:tenant:2:restaurant:4',
    )

    expect(tenantTwo).not.toBe(tenantThree)
    expect(tenantTwo).not.toBe(restaurantFive)
  })

  it('does not deduplicate writes, unknown identity, or caller-supplied bearer requests', () => {
    expect(buildScopedGetRequestDedupKey({ method: 'post', url: '/orders' }, 'user:8')).toBeNull()
    expect(buildScopedGetRequestDedupKey({ method: 'get', url: '/orders' }, null)).toBeNull()
    expect(buildScopedGetRequestDedupKey({
      method: 'get',
      url: '/orders',
      headers: { authorization: 'Bearer synthetic-token' },
    }, 'user:8')).toBeNull()
  })
})
