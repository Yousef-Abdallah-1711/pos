import http from 'k6/http'
import { check, group, sleep } from 'k6'

function isPrivateIPv4(host) {
  const parts = host.split('.')
  if (parts.length !== 4 || parts.some((part) => !/^\d{1,3}$/.test(part))) return false
  const octets = parts.map(Number)
  if (octets.some((part) => part < 0 || part > 255)) return false
  const [first, second] = octets
  return first === 10 || first === 127 ||
    (first === 172 && second >= 16 && second <= 31) ||
    (first === 192 && second === 168)
}

function isValidIPv6(host) {
  if (!/^[0-9a-f:]+$/i.test(host) || host.length === 0 || (host.match(/::/g) || []).length > 1) return false
  const groups = host.includes('::') ? host.split('::') : [host]
  const groupCount = groups.reduce((count, part) => count + (part ? part.split(':').length : 0), 0)
  if (host.includes('::') ? groupCount >= 8 : groupCount !== 8) return false
  return groups.every((part) => !part || part.split(':').every((group) => /^[0-9a-f]{1,4}$/i.test(group)))
}

function isLocalOrPrivateHost(rawHost) {
  const host = rawHost.toLowerCase()
  if (host === 'localhost' || host.endsWith('.localhost') || host.endsWith('.test')) return true

  if (isPrivateIPv4(host)) return true

  if (!isValidIPv6(host)) return false
  if (host === '::1') return true
  const firstGroup = Number.parseInt(host.split(':')[0] || host.split('::')[0], 16)
  return (firstGroup >= 0xfc00 && firstGroup <= 0xfdff) ||
    (firstGroup >= 0xfe80 && firstGroup <= 0xfebf)
}

function safeBaseURL(raw, label) {
  const match = typeof raw === 'string' && raw.match(/^(https?):\/\/([^/?#]+)\/?$/i)
  if (!match || /\s/.test(raw)) {
    throw new Error(`${label} must be an http(s) origin with no credentials, query, fragment, or path`)
  }

  const scheme = match[1].toLowerCase()
  const authority = match[2]
  if (authority.includes('@')) {
    throw new Error(`${label} must not contain URL credentials`)
  }

  let host
  let port = ''
  if (authority.startsWith('[')) {
    const close = authority.indexOf(']')
    if (close < 0) throw new Error(`${label} contains an invalid address`)
    host = authority.slice(1, close).toLowerCase()
    const suffix = authority.slice(close + 1)
    if (suffix && !/^:\d+$/.test(suffix)) throw new Error(`${label} contains an invalid port`)
    port = suffix
  } else {
    const hostPort = authority.match(/^([^:]+)(?::(\d+))?$/)
    if (!hostPort) throw new Error(`${label} contains an invalid host or port`)
    host = hostPort[1].toLowerCase()
    port = hostPort[2] ? `:${hostPort[2]}` : ''
  }

  if (port) {
    const portNumber = Number(port.slice(1))
    if (!Number.isInteger(portNumber) || portNumber < 1 || portNumber > 65535) {
      throw new Error(`${label} contains an invalid port`)
    }
  }
  if (host.startsWith('[') || host.includes('%')) {
    throw new Error(`${label} contains an invalid host`)
  }

  const reservedTestName = host === 'localhost' || host.endsWith('.localhost') || host.endsWith('.test')
  const labels = host.split('.')
  const validLabels = labels.every((part) => /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/.test(part))
  if (!isLocalOrPrivateHost(host) ||
      (reservedTestName && !validLabels) ||
      host === 'rakmyat.online' || host.endsWith('.rakmyat.online')) {
    throw new Error(`${label} must target localhost, an internal test hostname, or a private IP; public targets are blocked`)
  }

  const renderedHost = host.includes(':') ? `[${host}]` : host
  return `${scheme}://${renderedHost}${port}`
}

export function integerEnvironment(name, defaultValue, minimum, maximum) {
  const raw = __ENV[name]
  if (raw === undefined || raw === '') return defaultValue
  if (!/^\d+$/.test(raw)) throw new Error(`${name} must be an integer between ${minimum} and ${maximum}`)
  const value = Number(raw)
  if (!Number.isSafeInteger(value) || value < minimum || value > maximum) {
    throw new Error(`${name} must be an integer between ${minimum} and ${maximum}`)
  }
  return value
}

export const API_BASE_URL = safeBaseURL(__ENV.API_BASE_URL || 'http://127.0.0.1:8080', 'API_BASE_URL')
export const LANDING_BASE_URL = __ENV.LANDING_BASE_URL
  ? safeBaseURL(__ENV.LANDING_BASE_URL, 'LANDING_BASE_URL')
  : ''
export const RESTAURANT_SLUG = (__ENV.RESTAURANT_SLUG || '').trim()
export const THINK_TIME_SECONDS = Number(__ENV.THINK_TIME_SECONDS || '1')
if (!Number.isFinite(THINK_TIME_SECONDS) || THINK_TIME_SECONDS <= 0) {
  throw new Error('THINK_TIME_SECONDS must be a positive number')
}

const cookieFile = __ENV.SESSION_COOKIE_VALUES_FILE || ''
export const SESSION_COOKIE_VALUES = cookieFile
  ? open(cookieFile).split(/\r?\n/).map((value) => value.trim()).filter(Boolean)
  : []

if (SESSION_COOKIE_VALUES.some((value) => !/^[A-Za-z0-9_-]{43}$/.test(value))) {
  throw new Error('SESSION_COOKIE_VALUES_FILE must contain one raw opaque session cookie value per line')
}
if (new Set(SESSION_COOKIE_VALUES).size !== SESSION_COOKIE_VALUES.length) {
  throw new Error('SESSION_COOKIE_VALUES_FILE must contain unique synthetic session values')
}

export function requireDedicatedEnvironment(maxVUs, requestRate = 0) {
  if ((maxVUs > 100 || requestRate > 100) && __ENV.POS_LOAD_TEST_ENVIRONMENT !== 'dedicated') {
    throw new Error('More than 100 VUs or 100 RPS requires POS_LOAD_TEST_ENVIRONMENT=dedicated')
  }
}

export function requireSessionPool(minimumValues) {
  if (SESSION_COOKIE_VALUES.length < minimumValues) {
    throw new Error(`SESSION_COOKIE_VALUES_FILE must contain at least ${minimumValues} synthetic session values`)
  }
}

export function publicRequestRateCeiling(maxVUs) {
  const requestsPerIteration = 1 +
    (RESTAURANT_SLUG ? 1 : 0) +
    (SESSION_COOKIE_VALUES.length > 0 ? 1 : 0) +
    (LANDING_BASE_URL ? 1 : 0)
  return Math.ceil(maxVUs * requestsPerIteration / THINK_TIME_SECONDS)
}

function sessionCookieForCurrentVU() {
  return SESSION_COOKIE_VALUES[(__VU - 1) % SESSION_COOKIE_VALUES.length]
}

function requestParams(name, includeSession = false) {
  const headers = {}
  if (includeSession) headers.Cookie = `poos_session=${sessionCookieForCurrentVU()}`
  return { headers, tags: { name }, timeout: __ENV.REQUEST_TIMEOUT || '5s', redirects: 0 }
}

export function protectedSessionCheck() {
  group('Authenticated session validation', () => {
    const response = http.get(`${API_BASE_URL}/api/v1/auth/session`, requestParams('GET /api/v1/auth/session', true))
    check(response, { 'session is valid': (result) => result.status === 200 })
  })
}

export function publicIteration() {
  group('API health', () => {
    const response = http.get(`${API_BASE_URL}/api/v1/health`, requestParams('GET /api/v1/health'))
    check(response, { 'health endpoint responds': (result) => result.status === 200 })
  })

  if (RESTAURANT_SLUG) {
    group('Anonymous restaurant menu', () => {
      const slug = encodeURIComponent(RESTAURANT_SLUG)
      const response = http.get(
        `${API_BASE_URL}/api/v1/public/restaurants/${slug}/menu`,
        requestParams('GET /api/v1/public/restaurants/:slug/menu')
      )
      check(response, { 'public menu responds': (result) => result.status === 200 })
    })
  }

  if (SESSION_COOKIE_VALUES.length > 0) protectedSessionCheck()

  if (LANDING_BASE_URL) {
    group('Anonymous Landing page', () => {
      const response = http.get(`${LANDING_BASE_URL}/`, requestParams('GET /'))
      check(response, { 'Landing page responds': (result) => result.status === 200 })
    })
  }

  sleep(THINK_TIME_SECONDS)
}
