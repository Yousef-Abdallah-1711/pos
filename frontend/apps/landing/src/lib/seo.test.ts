import { afterEach, describe, expect, it, vi } from 'vitest'
import { getLandingOrigin } from './seo'

describe('getLandingOrigin', () => {
  afterEach(() => {
    vi.unstubAllEnvs()
  })

  it('uses the configured landing origin', () => {
    vi.stubEnv('NEXT_PUBLIC_LANDING_URL', 'https://www.poos.example/')
    vi.stubEnv('NODE_ENV', 'production')

    expect(getLandingOrigin().origin).toBe('https://www.poos.example')
  })

  it('rejects a non-HTTPS production origin', () => {
    vi.stubEnv('NEXT_PUBLIC_LANDING_URL', 'http://localhost:3001')
    vi.stubEnv('NODE_ENV', 'production')

    expect(() => getLandingOrigin()).toThrow(/must use HTTPS/)
  })

  it('uses localhost only in development when unset', () => {
    vi.stubEnv('NEXT_PUBLIC_LANDING_URL', '')
    vi.stubEnv('NODE_ENV', 'development')

    expect(getLandingOrigin().origin).toBe('http://localhost:3001')
  })
})
