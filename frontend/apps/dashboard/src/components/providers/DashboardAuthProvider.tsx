'use client'

import { ReactNode, useEffect, useState } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { useTenantStatusStore } from '@/stores/tenantStatusStore'
import { getConfig } from '@/lib/config'

interface SessionUser {
  id: number
  user_id: number
  email: string
  name: string
  tenant_id: number
  restaurant_id: number
  role: string
  is_super_admin: boolean
}

function isSessionUser(value: unknown): value is SessionUser {
  if (!value || typeof value !== 'object') return false
  const candidate = value as Partial<SessionUser>
  return typeof candidate.id === 'number'
    && typeof candidate.email === 'string'
    && typeof candidate.tenant_id === 'number'
    && typeof candidate.role === 'string'
}

const DashboardAuthProvider = ({ children }: { children: ReactNode }) => {
  const setAuth = useAuthStore((state) => state.setAuth)
  const user = useAuthStore((state) => state.user)
  const checkStatus = useTenantStatusStore((state) => state.checkStatus)
  const [isReady, setIsReady] = useState(false)

  useEffect(() => {
    let active = true
    const locale = window.location.pathname.split('/')[1] === 'ar' ? 'ar' : 'en'
    const originalFetch = window.fetch.bind(window)
    const apiOrigin = new URL(getConfig().apiUrl).origin
    window.localStorage.removeItem('token')
    window.localStorage.removeItem('auth-storage')
    window.sessionStorage.removeItem('impersonation_token')
    window.fetch = (input, init) => {
      const requestURL = typeof input === 'string' ? input : input instanceof URL ? input.href : input.url
      const parsedURL = new URL(requestURL, window.location.href)
      const isApiRequest = parsedURL.origin === apiOrigin || parsedURL.pathname.startsWith('/api/v1/')
      if (isApiRequest) {
        const headers = new Headers(init?.headers ?? (input instanceof Request ? input.headers : undefined))
        // Browser auth is cookie-backed. Never forward legacy JWTs from persisted UI state.
        headers.delete('Authorization')
        return originalFetch(input, { ...init, headers, credentials: 'include' })
      }
      return originalFetch(input, init)
    }

    originalFetch(`${getConfig().apiUrl}/auth/session`, { credentials: 'include', cache: 'no-store' })
      .then(async (response) => {
        if (!response.ok) throw new Error('Session unavailable')
        const result: unknown = await response.json()
        if (!isSessionUser(result)) throw new Error('Invalid session response')
        return result
      })
      .then((sessionUser) => {
        if (!active) return
        setAuth(sessionUser, '')
        setIsReady(true)
      })
      .catch(() => {
        if (active) window.location.replace(`${getConfig().landingUrl}/${locale}/auth/login`)
      })

    return () => {
      active = false
      window.fetch = originalFetch
    }
  }, [setAuth])

  useEffect(() => {
    if (isReady && user) checkStatus()
  }, [checkStatus, isReady, user])

  if (!isReady) {
    return (
      <div className="min-h-screen flex items-center justify-center" aria-live="polite">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
      </div>
    )
  }

  return <>{children}</>
}

export default DashboardAuthProvider
