// src/lib/config.ts
//
// Production values should be set via environment variables:
//   NEXT_PUBLIC_LANDING_URL       - e.g. https://landing.example.com
//   NEXT_PUBLIC_DASHBOARD_URL     - e.g. https://dashboard.example.com
//   NEXT_PUBLIC_API_URL           - e.g. https://api.example.com/api/v1
//   NEXT_PUBLIC_RESTAURANT_DOMAIN - e.g. example.com (used for tenant subdomains)
//
// Localhost fallbacks are provided for local development only.

const DEV_DEFAULTS = {
  landingUrl: 'http://localhost:3001',
  dashboardUrl: 'http://localhost:3002',
  apiUrl: 'http://localhost:8080/api/v1',
  restaurantDomain: 'localhost:3003',
} as const

function isProductionRuntime() {
  return process.env.NODE_ENV === 'production'
}

function getBrowserOrigin() {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin
  }
  return ''
}

function normalizeUrl(value: string) {
  return value.replace(/\/+$/, '')
}

function resolveUrlValue(
  label: string,
  configuredValue: string | undefined,
  devFallback: string,
  options?: { allowCurrentOriginInProduction?: boolean }
) {
  const trimmed = configuredValue?.trim()
  if (trimmed) {
    return normalizeUrl(trimmed)
  }

  if (!isProductionRuntime()) {
    return normalizeUrl(devFallback)
  }

  if (options?.allowCurrentOriginInProduction) {
    const origin = getBrowserOrigin()
    if (origin) {
      return normalizeUrl(origin)
    }
  }

  const keyMessage = label === 'API URL'
    ? 'NEXT_PUBLIC_API_URL'
    : label === 'Landing URL'
      ? 'NEXT_PUBLIC_LANDING_URL'
      : 'NEXT_PUBLIC_DASHBOARD_URL'

  throw new Error(`[landing/config] Missing ${label} in production. Set ${keyMessage}.`)
}

function resolveDomainValue(configuredValue: string | undefined) {
  const trimmed = configuredValue?.trim()
  if (trimmed) return trimmed
  if (!isProductionRuntime()) return DEV_DEFAULTS.restaurantDomain
  throw new Error('[landing/config] Missing restaurant domain in production. Set NEXT_PUBLIC_RESTAURANT_DOMAIN.')
}

export const getConfig = () => {
  const isDevelopment = process.env.NODE_ENV === 'development'

  return {
    isDevelopment,
    landingUrl: resolveUrlValue('Landing URL', process.env.NEXT_PUBLIC_LANDING_URL, DEV_DEFAULTS.landingUrl, {
      allowCurrentOriginInProduction: true,
    }),
    dashboardUrl: resolveUrlValue('Dashboard URL', process.env.NEXT_PUBLIC_DASHBOARD_URL, DEV_DEFAULTS.dashboardUrl),
    apiUrl: resolveUrlValue('API URL', process.env.NEXT_PUBLIC_API_URL, DEV_DEFAULTS.apiUrl),
    restaurantDomain: resolveDomainValue(process.env.NEXT_PUBLIC_RESTAURANT_DOMAIN),
  }
}

export const getApiBaseUrl = () => getConfig().apiUrl
export const getLandingBaseUrl = () => getConfig().landingUrl
export const getDashboardBaseUrl = () => getConfig().dashboardUrl
export const getRestaurantDomain = () => getConfig().restaurantDomain
