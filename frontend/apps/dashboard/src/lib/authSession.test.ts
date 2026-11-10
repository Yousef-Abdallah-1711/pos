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

  it('builds locale-aware login urls', async () => {
    vi.stubEnv('NEXT_PUBLIC_LANDING_URL', 'http://localhost:3001')
    const { getDashboardLoginUrl } = await import('@/lib/authSession')

    expect(getDashboardLoginUrl('/ar/dashboard/orders')).toBe('http://localhost:3001/ar/auth/login')
    expect(getDashboardLoginUrl('/dashboard')).toBe('http://localhost:3001/en/auth/login')
  })
})
