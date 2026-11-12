const ENV_KEYS = [
  'NODE_ENV',
  'NEXT_PUBLIC_API_BASE_URL',
  'NEXT_PUBLIC_API_URL',
  'NEXT_PUBLIC_API_TIMEOUT',
  'NEXT_PUBLIC_DEBUG_LOGS',
  'DEBUG_LOGS',
] as const

const ORIGINAL_ENV = { ...process.env }

function resetRelevantEnv() {
  for (const key of ENV_KEYS) {
    delete process.env[key]
  }
}

function loadApiModule(): typeof import('./api') {
  jest.resetModules()
  return require('./api')
}

describe('restaurant API config resolver', () => {
  beforeEach(() => {
    resetRelevantEnv()
  })

  afterEach(() => {
    process.env = { ...ORIGINAL_ENV }
    jest.resetModules()
  })

  it('prefers NEXT_PUBLIC_API_BASE_URL and normalizes trailing slashes', () => {
    Object.assign(process.env, { NODE_ENV: 'production' })
    process.env.NEXT_PUBLIC_API_BASE_URL = 'https://api.example.com/api/v1///'
    process.env.NEXT_PUBLIC_API_URL = 'https://legacy.example.com/api/v1/'

    const { getApiBaseUrl } = loadApiModule()

    expect(getApiBaseUrl()).toBe('https://api.example.com/api/v1')
  })

  it('uses development fallback when env is missing outside production', () => {
    Object.assign(process.env, { NODE_ENV: 'development' })

    const { getApiBaseUrl } = loadApiModule()

    expect(getApiBaseUrl()).toBe('http://localhost:8080/api/v1')
  })

  it('throws in production when API URL env is missing', () => {
    Object.assign(process.env, { NODE_ENV: 'production' })

    const { getApiBaseUrl, validateApiConfig } = loadApiModule()

    expect(() => getApiBaseUrl()).toThrow(/NEXT_PUBLIC_API_BASE_URL/)
    expect(() => validateApiConfig()).toThrow(/NEXT_PUBLIC_API_BASE_URL/)
  })

  it('builds endpoint URLs using resolved base URL', () => {
    Object.assign(process.env, { NODE_ENV: 'production' })
    process.env.NEXT_PUBLIC_API_BASE_URL = 'https://api.example.com/api/v1/'

    const { buildApiUrl } = loadApiModule()

    expect(buildApiUrl('/health')).toBe('https://api.example.com/api/v1/health')
  })
})
