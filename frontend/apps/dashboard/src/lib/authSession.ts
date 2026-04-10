'use client'

import { getConfig } from '@/lib/config'
import { getLocaleFromPath } from '@/lib/translations'
import { COOKIE_SESSION_TOKEN } from '@/stores/authStore'

export const hasAuthenticatedSession = (user: unknown, token: string | null | undefined) => {
  return !!user || !!token
}

export const hasUsableBearerToken = (token: string | null | undefined) => {
  return !!token && token !== COOKIE_SESSION_TOKEN
}

export const getDashboardLoginUrl = (pathname?: string) => {
  const currentPath = pathname
    ?? (typeof window !== 'undefined' ? window.location.pathname : '/en/dashboard')
  const locale = getLocaleFromPath(currentPath)
  return `${getConfig().landingUrl}/${locale}/auth/login`
}
