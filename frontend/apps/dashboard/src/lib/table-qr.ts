import { getConfig } from '@/lib/config'

function normalizeBaseUrl(value: string) {
  return value.replace(/\/+$/, '')
}

export function getRestaurantBaseUrl(restaurantDomain = getConfig().restaurantDomain) {
  const normalizedDomain = normalizeBaseUrl(restaurantDomain)

  if (
    normalizedDomain.startsWith('http://') ||
    normalizedDomain.startsWith('https://')
  ) {
    return normalizedDomain
  }

  if (normalizedDomain.includes('localhost')) {
    return `http://${normalizedDomain}`
  }

  return `https://${normalizedDomain}`
}

export function buildTableMenuUrl(
  tableID: number,
  locale: string,
  restaurantDomain?: string
) {
  return `${getRestaurantBaseUrl(restaurantDomain)}/${locale}/table/${tableID}/menu`
}

export function sanitizeQrFileName(value: string) {
  return value
    .replace(/[^a-zA-Z0-9-_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}
