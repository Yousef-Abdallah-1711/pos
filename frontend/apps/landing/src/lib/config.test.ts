import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const ENV_KEYS = [
  'NODE_ENV',
  'NEXT_PUBLIC_LANDING_URL',
  'NEXT_PUBLIC_DASHBOARD_URL',
  'NEXT_PUBLIC_API_URL',
  'NEXT_PUBLIC_RESTAURANT_DOMAIN',
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

describe('landing URL config resolver', () => {
  beforeEach(() => {
    resetRelevantEnv()
  })

  afterEach(() => {
    process.env = { ...ORIGINAL_ENV }
    vi.unstubAllEnvs()
    vi.resetModules()
  })

  it('uses localhost dev defaults outside production', async () => {
    vi.stubEnv('NODE_ENV', 'development')
    const { getConfig } = await loadConfigModule()

    expect(getConfig()).toMatchObject({
      isDevelopment: true,
      landingUrl: 'http://localhost:3001',
      dashboardUrl: 'http://localhost:3002',
      apiUrl: 'http://localhost:8080/api/v1',
      restaurantDomain: 'localhost:3003',
    })
  })

  it('normalizes configured production URLs', async () => {
    vi.stubEnv('NODE_ENV', 'production')
    process.env.NEXT_PUBLIC_LANDING_URL = 'https://www.example.com///'
    process.env.NEXT_PUBLIC_DASHBOARD_URL = 'https://dashboard.example.com/'
    process.env.NEXT_PUBLIC_API_URL = 'https://api.example.com/api/v1/'
    process.env.NEXT_PUBLIC_RESTAURANT_DOMAIN = 'example.com'

    const { getConfig } = await loadConfigModule()

    expect(getConfig()).toMatchObject({
      isDevelopment: false,
      landingUrl: 'https://www.example.com',
      dashboardUrl: 'https://dashboard.example.com',
      apiUrl: 'https://api.example.com/api/v1',
      restaurantDomain: 'example.com',
    })
  })

  it('throws in production when API URL is missing', async () => {
    vi.stubEnv('NODE_ENV', 'production')
    process.env.NEXT_PUBLIC_LANDING_URL = 'https://www.example.com'
    process.env.NEXT_PUBLIC_DASHBOARD_URL = 'https://dashboard.example.com'
    process.env.NEXT_PUBLIC_RESTAURANT_DOMAIN = 'example.com'

    const { getApiBaseUrl } = await loadConfigModule()
    expect(() => getApiBaseUrl()).toThrow(/NEXT_PUBLIC_API_URL/)
  })

  it('throws in production when restaurant domain is missing', async () => {
    vi.stubEnv('NODE_ENV', 'production')
    process.env.NEXT_PUBLIC_LANDING_URL = 'https://www.example.com'
    process.env.NEXT_PUBLIC_DASHBOARD_URL = 'https://dashboard.example.com'
    process.env.NEXT_PUBLIC_API_URL = 'https://api.example.com/api/v1'

    const { getRestaurantDomain } = await loadConfigModule()
    expect(() => getRestaurantDomain()).toThrow(/NEXT_PUBLIC_RESTAURANT_DOMAIN/)
  })
})
