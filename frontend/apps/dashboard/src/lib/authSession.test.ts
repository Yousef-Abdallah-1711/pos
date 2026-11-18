import { afterEach, describe, expect, it, vi } from 'vitest'

import { COOKIE_SESSION_MARKER } from '@/stores/authStore'

describe('auth session helpers', () => {
  afterEach(() => {
    vi.resetModules()
    vi.unstubAllEnvs()
  })

  it('treats cookie-backed user sessions as authenticated', async () => {
    const { hasAuthenticatedSession } = await import('@/lib/authSession')

    expect(hasAuthenticatedSession({ id: 1 }, null)).toBe(true)
    expect(hasAuthenticatedSession(null, COOKIE_SESSION_MARKER)).toBe(true)
    expect(hasAuthenticatedSession(null, null)).toBe(false)
  })

  it('does not treat the cookie-session placeholder as a bearer token', async () => {
    const { hasUsableBearerToken } = await import('@/lib/authSession')

    expect(hasUsableBearerToken('real-token')).toBe(true)
    expect(hasUsableBearerToken(COOKIE_SESSION_MARKER)).toBe(false)
    expect(hasUsableBearerToken(null)).toBe(false)
  })

  it('distinguishes expiry, authorization denial, and temporary session failures', async () => {
    const { classifySessionBootstrapStatus } = await import('@/lib/authSession')

    expect(classifySessionBootstrapStatus(200)).toBe('ready')
    expect(classifySessionBootstrapStatus(401)).toBe('expired')
    expect(classifySessionBootstrapStatus(403)).toBe('forbidden')
    expect(classifySessionBootstrapStatus(429)).toBe('unavailable')
    expect(classifySessionBootstrapStatus(503)).toBe('unavailable')
  })

  it('shares one validated callback session request across Strict Mode effect replay', async () => {
    const { requestDashboardCallbackSession } = await import('@/lib/authSession')
    const fetcher = vi.fn<typeof fetch>().mockResolvedValue(new Response(JSON.stringify({
      id: 21,
      user_id: 21,
      email: 'owner@example.test',
      name: 'Owner',
      tenant_id: 12,
      restaurant_id: 34,
      role: 'owner',
      is_super_admin: false,
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    const [first, replay] = await Promise.all([
      requestDashboardCallbackSession('https://api.example.test/api/v1', fetcher),
      requestDashboardCallbackSession('https://api.example.test/api/v1', fetcher),
    ])

    expect(fetcher).toHaveBeenCalledTimes(1)
    expect(fetcher).toHaveBeenCalledWith('https://api.example.test/api/v1/auth/session', expect.objectContaining({
      credentials: 'include',
      cache: 'no-store',
      signal: expect.any(AbortSignal),
    }))
    expect(first).toEqual(replay)
  })

  it('rejects a Super Admin identity from the staff Dashboard callback', async () => {
    const { requestDashboardCallbackSession } = await import('@/lib/authSession')
    const fetcher = vi.fn<typeof fetch>().mockResolvedValue(new Response(JSON.stringify({
      id: 21,
      user_id: 21,
      email: 'admin@example.test',
      name: 'Admin',
      tenant_id: 12,
      restaurant_id: 0,
      role: 'super_admin',
      is_super_admin: true,
    }), { status: 200, headers: { 'Content-Type': 'application/json' } }))

    await expect(requestDashboardCallbackSession('https://api.example.test/api/v1', fetcher))
      .rejects.toThrow('Session is not a staff session')
  })

  it('builds locale-aware login urls', async () => {
    vi.stubEnv('NEXT_PUBLIC_LANDING_URL', 'http://localhost:3001')
    const { getDashboardLoginUrl } = await import('@/lib/authSession')

    expect(getDashboardLoginUrl('/ar/dashboard/orders')).toBe('http://localhost:3001/ar/auth/login')
    expect(getDashboardLoginUrl('/dashboard')).toBe('http://localhost:3001/ar/auth/login')
  })
})
