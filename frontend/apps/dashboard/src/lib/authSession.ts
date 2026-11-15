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

export const getDashboardLoginUrl = (pathname?: string) => {
  const currentPath = pathname
    ?? (typeof window !== 'undefined' ? window.location.pathname : `/${defaultLocale}/dashboard`)
  const locale = getLocaleFromPath(currentPath)
  return `${getConfig().landingUrl}/${locale}/auth/login`
}
