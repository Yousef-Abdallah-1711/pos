import { getConfig } from '@/lib/config'

export type WebsiteQrTarget = 'home' | 'menu'

function normalizeBaseUrl(value: string) {
  return value.replace(/\/+$/, '')
}

function normalizeLocaleSegment(locale: string) {
  return locale.trim().replace(/^\/+|\/+$/g, '')
}

export function getWebsiteBaseUrl(websiteHost = getConfig().restaurantDomain) {
  const normalizedHost = normalizeBaseUrl(websiteHost.trim())

  if (
    normalizedHost.startsWith('http://') ||
    normalizedHost.startsWith('https://')
  ) {
    return normalizedHost
  }

  if (normalizedHost.includes('localhost')) {
    return `http://${normalizedHost}`
  }

  return `https://${normalizedHost}`
}

export function buildWebsiteQrTargetUrl(
  target: WebsiteQrTarget,
  locale: string,
  websiteHost?: string
) {
  const base = getWebsiteBaseUrl(websiteHost)
  const normalizedLocale = normalizeLocaleSegment(locale)
  const localePrefix = normalizedLocale ? `/${normalizedLocale}` : ''

  if (target === 'menu') {
    return `${base}${localePrefix}/menu`
  }

  return `${base}${localePrefix}`
}
