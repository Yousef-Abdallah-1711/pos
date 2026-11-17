'use client'

import { getConfig } from '@/lib/config'
import { defaultLocale } from '@pos-saas/i18n'
import { getLocaleFromPath } from '@/lib/translations'
import { COOKIE_SESSION_MARKER } from '@/stores/authStore'

export const hasAuthenticatedSession = (user: unknown, token: string | null | undefined) => {
  return !!user || !!token
}

export const hasUsableBearerToken = (token: string | null | undefined) => {
  return !!token && token !== COOKIE_SESSION_MARKER
}

export type SessionBootstrapDisposition = 'ready' | 'expired' | 'forbidden' | 'unavailable'

export const classifySessionBootstrapStatus = (status: number): SessionBootstrapDisposition => {
  if (status === 200) return 'ready'
  if (status === 401) return 'expired'
  if (status === 403) return 'forbidden'
  return 'unavailable'
}

export const getDashboardLoginUrl = (pathname?: string) => {
  const currentPath = pathname
    ?? (typeof window !== 'undefined' ? window.location.pathname : `/${defaultLocale}/dashboard`)
  const locale = getLocaleFromPath(currentPath)
  return `${getConfig().landingUrl}/${locale}/auth/login`
}
