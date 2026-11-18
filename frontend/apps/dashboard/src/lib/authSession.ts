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

export interface DashboardSessionUser {
  id: number
  user_id: number
  email: string
  name: string
  tenant_id: number
  restaurant_id: number
  restaurant_name?: string
  role: string
  is_super_admin: false
}

export const isDashboardSessionUser = (value: unknown): value is DashboardSessionUser => {
  if (!value || typeof value !== 'object') return false
  const candidate = value as Partial<DashboardSessionUser>
  return typeof candidate.id === 'number'
    && Number.isSafeInteger(candidate.id)
    && candidate.id > 0
    && typeof candidate.user_id === 'number'
    && Number.isSafeInteger(candidate.user_id)
    && candidate.user_id === candidate.id
    && typeof candidate.email === 'string'
    && candidate.email.trim().length > 0
    && typeof candidate.name === 'string'
    && typeof candidate.tenant_id === 'number'
    && Number.isSafeInteger(candidate.tenant_id)
    && candidate.tenant_id > 0
    && typeof candidate.restaurant_id === 'number'
    && Number.isSafeInteger(candidate.restaurant_id)
    && candidate.restaurant_id >= 0
    && typeof candidate.role === 'string'
    && candidate.role.trim().length > 0
    && candidate.is_super_admin === false
}

let callbackSessionRequest: Promise<DashboardSessionUser> | null = null

export function requestDashboardCallbackSession(
  apiUrl: string,
  fetcher: typeof fetch = fetch,
): Promise<DashboardSessionUser> {
  if (callbackSessionRequest) return callbackSessionRequest

  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 8_000)
  const request = Promise.resolve()
    .then(() => fetcher(`${apiUrl}/auth/session`, {
      credentials: 'include',
      cache: 'no-store',
      signal: controller.signal,
    }))
    .then(async (response) => {
      if (!response.ok) throw new Error('Session unavailable')
      const session: unknown = await response.json()
      if (!isDashboardSessionUser(session)) throw new Error('Session is not a staff session')
      return session
    })
    .finally(() => {
      clearTimeout(timeout)
      if (callbackSessionRequest === request) callbackSessionRequest = null
    })

  callbackSessionRequest = request
  return request
}

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
