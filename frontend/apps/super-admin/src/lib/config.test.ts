import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const ENV_KEYS = [
  'NODE_ENV',
  'NEXT_PUBLIC_SUPER_ADMIN_API_URL',
  'NEXT_PUBLIC_DASHBOARD_URL',
] as const

const ORIGINAL_ENV = { ...process.env }

function resetRelevantEnv() {
  for (const key of ENV_KEYS) {
    delete process.env[key]
  }
}

async function loadConfigModule() {
  vi.resetModules()
  return import('./config')
}

describe('super-admin URL config resolver', () => {
  beforeEach(() => {
    resetRelevantEnv()
  })

  afterEach(() => {
    process.env = { ...ORIGINAL_ENV }
    vi.resetModules()
  })

  it('uses localhost dev defaults outside production', async () => {
    process.env.NODE_ENV = 'development'
    const { getSuperAdminApiBaseUrl, getDashboardBaseUrl } = await loadConfigModule()

    expect(getSuperAdminApiBaseUrl()).toBe('http://localhost:8080/api/v1/super-admin')
    expect(getDashboardBaseUrl()).toBe('http://localhost:3002')
  })

  it('normalizes configured production URLs', async () => {
    process.env.NODE_ENV = 'production'
    process.env.NEXT_PUBLIC_SUPER_ADMIN_API_URL = 'https://api.example.com/api/v1/super-admin/'
    process.env.NEXT_PUBLIC_DASHBOARD_URL = 'https://dashboard.example.com//'

    const { getSuperAdminApiBaseUrl, getDashboardBaseUrl } = await loadConfigModule()

    expect(getSuperAdminApiBaseUrl()).toBe('https://api.example.com/api/v1/super-admin')
    expect(getDashboardBaseUrl()).toBe('https://dashboard.example.com')
  })

  it('throws in production when required values are missing', async () => {
    process.env.NODE_ENV = 'production'
    const { getSuperAdminApiBaseUrl, getDashboardBaseUrl } = await loadConfigModule()

    expect(() => getSuperAdminApiBaseUrl()).toThrow(/NEXT_PUBLIC_SUPER_ADMIN_API_URL/)
    expect(() => getDashboardBaseUrl()).toThrow(/NEXT_PUBLIC_DASHBOARD_URL/)
  })
})
