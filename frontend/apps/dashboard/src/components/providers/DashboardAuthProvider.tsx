'use client'

import { ReactNode, useCallback, useEffect, useRef, useState } from 'react'
import { useTranslations } from 'next-intl'
import { useAuthStore } from '@/stores/authStore'
import { useTenantStatusStore } from '@/stores/tenantStatusStore'
import { classifySessionBootstrapStatus, getDashboardLoginUrl } from '@/lib/authSession'
import { getConfig } from '@/lib/config'

interface SessionUser {
  id: number
  user_id: number
  email: string
  name: string
  tenant_id: number
  restaurant_id: number
  restaurant_name?: string
  role: string
  is_super_admin: boolean
}

type BootstrapState = 'checking' | 'ready' | 'unavailable' | 'forbidden'

const SESSION_REQUEST_TIMEOUT_MS = 8_000
let sessionRequest: Promise<Response> | null = null

function isSessionUser(value: unknown): value is SessionUser {
  if (!value || typeof value !== 'object') return false
  const candidate = value as Partial<SessionUser>
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

function requestSession(apiUrl: string, fetcher: typeof fetch): Promise<Response> {
  if (sessionRequest) return sessionRequest
  const controller = new AbortController()
  const timeout = window.setTimeout(() => controller.abort(), SESSION_REQUEST_TIMEOUT_MS)
  const request = Promise.resolve()
    .then(() => fetcher(`${apiUrl}/auth/session`, { credentials: 'include', cache: 'no-store', signal: controller.signal }))
    .finally(() => {
      window.clearTimeout(timeout)
      if (sessionRequest === request) sessionRequest = null
    })
  sessionRequest = request
  return request
}

const DashboardAuthProvider = ({ children }: { children: ReactNode }) => {
  const setAuth = useAuthStore((state) => state.setAuth)
  const clearAuth = useAuthStore((state) => state.logout)
  const user = useAuthStore((state) => state.user)
  const checkStatus = useTenantStatusStore((state) => state.checkStatus)
  const clearTenantStatus = useTenantStatusStore((state) => state.clear)
  const translate = useTranslations('auth')
  const [bootstrapState, setBootstrapStateState] = useState<BootstrapState>('checking')
  const [retryCount, setRetryCount] = useState(0)
  const bootstrapStateRef = useRef<BootstrapState>('checking')
  const setBootstrapState = useCallback((nextState: BootstrapState) => {
    bootstrapStateRef.current = nextState
    setBootstrapStateState(nextState)
  }, [])

  useEffect(() => {
    let active = true
    const locale = window.location.pathname.split('/')[1] === 'ar' ? 'ar' : 'en'
    const originalFetch = window.fetch.bind(window)
    const apiOrigin = new URL(getConfig().apiUrl).origin
    window.localStorage.removeItem('token')
    window.localStorage.removeItem('auth-storage')
    window.sessionStorage.removeItem('impersonation_token')
    window.fetch = (input, init) => fetchWithCookieSession(input, init, apiOrigin, originalFetch)

    void restoreDashboardSession({
      active: () => active,
      apiUrl: getConfig().apiUrl,
      fetcher: originalFetch,
      locale,
      setAuth,
      clearAuth,
      clearTenantStatus,
      setBootstrapState,
    })

    return () => {
      active = false
      window.fetch = originalFetch
    }
  }, [clearAuth, clearTenantStatus, retryCount, setAuth, setBootstrapState])

  useEffect(() => {
    let active = true
    let refreshTimeout: number | undefined
    const revalidateWhenVisible = () => {
      if (document.visibilityState !== 'visible' || bootstrapStateRef.current === 'checking') return
      window.clearTimeout(refreshTimeout)
      refreshTimeout = window.setTimeout(() => {
        if (!active) return
        const locale = window.location.pathname.split('/')[1] === 'ar' ? 'ar' : 'en'
        setBootstrapState('checking')
        void restoreDashboardSession({
          active: () => active,
          apiUrl: getConfig().apiUrl,
          fetcher: window.fetch.bind(window),
          locale,
          setAuth,
          clearAuth,
          clearTenantStatus,
          setBootstrapState,
        })
      }, 100)
    }
    window.addEventListener('focus', revalidateWhenVisible)
    document.addEventListener('visibilitychange', revalidateWhenVisible)

    return () => {
      active = false
      window.clearTimeout(refreshTimeout)
      window.removeEventListener('focus', revalidateWhenVisible)
      document.removeEventListener('visibilitychange', revalidateWhenVisible)
    }
  }, [clearAuth, clearTenantStatus, setAuth, setBootstrapState])

  useEffect(() => {
    if (bootstrapState === 'ready' && user) checkStatus()
  }, [bootstrapState, checkStatus, user])

  if (bootstrapState === 'checking') {
    return <div className="min-h-screen flex items-center justify-center" role="status" aria-live="polite">
      <span className="sr-only">{translate('checkingSession')}</span>
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" aria-hidden="true" />
    </div>
  }

  if (bootstrapState === 'unavailable' || bootstrapState === 'forbidden') {
    return <SessionRecoveryScreen
      message={translate(bootstrapState === 'forbidden' ? 'sessionForbidden' : 'sessionUnavailable')}
      retryLabel={translate('retrySession')}
      onRetry={() => {
        setBootstrapState('checking')
        setRetryCount(count => count + 1)
      }}
    />
  }

  return <>{children}</>
}

interface SessionRestoreOptions {
  active: () => boolean
  apiUrl: string
  fetcher: typeof fetch
  locale: string
  setAuth: (user: SessionUser, token: string) => void
  clearAuth: () => void
  clearTenantStatus: () => void
  setBootstrapState: (state: BootstrapState) => void
}

async function restoreDashboardSession(options: SessionRestoreOptions) {
  try {
    const response = await requestSession(options.apiUrl, options.fetcher)
    if (!options.active()) return

    const disposition = classifySessionBootstrapStatus(response.status)
    if (disposition === 'expired') {
      options.clearAuth()
      options.clearTenantStatus()
      window.location.replace(getDashboardLoginUrl(`/${options.locale}/dashboard`))
      return
    }
    if (disposition === 'forbidden') {
      options.clearTenantStatus()
      options.setBootstrapState('forbidden')
      return
    }
    if (disposition !== 'ready') {
      options.setBootstrapState('unavailable')
      return
    }

    const sessionUser: unknown = await response.json()
    if (!options.active()) return
    if (!isSessionUser(sessionUser)) {
      options.setBootstrapState('unavailable')
      return
    }
    options.clearTenantStatus()
    options.setAuth(sessionUser, '')
    options.setBootstrapState('ready')
  } catch {
    if (options.active()) options.setBootstrapState('unavailable')
  }
}

function fetchWithCookieSession(input: RequestInfo | URL, init: RequestInit | undefined, apiOrigin: string, originalFetch: typeof fetch) {
  const requestURL = typeof input === 'string' ? input : input instanceof URL ? input.href : input.url
  const parsedURL = new URL(requestURL, window.location.href)
  const isApiRequest = parsedURL.origin === apiOrigin || (parsedURL.origin === window.location.origin && parsedURL.pathname.startsWith('/api/v1/'))
  if (!isApiRequest) return originalFetch(input, init)

  const headers = new Headers(init?.headers ?? (input instanceof Request ? input.headers : undefined))
  headers.delete('Authorization')
  return originalFetch(input, { ...init, headers, credentials: 'include' })
}

function SessionRecoveryScreen({ message, retryLabel, onRetry }: { message: string; retryLabel: string; onRetry: () => void }) {
  return <main className="min-h-screen flex flex-col items-center justify-center gap-4 p-6 text-center" role="alert">
    <p className="max-w-xl text-muted-foreground">{message}</p>
    <button className="min-h-12 rounded-md bg-primary px-5 py-3 text-primary-foreground focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2" type="button" onClick={onRetry}>
      {retryLabel}
    </button>
  </main>
}

export default DashboardAuthProvider
