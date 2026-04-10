'use client'

import { useEffect } from 'react'
import { useSearchParams, usePathname } from 'next/navigation'
import { getDashboardLoginUrl } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'
import { warmupCurrentUserPermissions } from '@/hooks/useRbac'
import { getConfig } from '@/lib/config'
import { debugLog, debugError } from '@/lib/debug'

export default function AuthCallbackContent() {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const setAuth = useAuthStore((state) => state.setAuth)

  useEffect(() => {
    const token = searchParams.get('token')
    const userParam = searchParams.get('user')

    debugLog('=== CALLBACK: Processing auth ===')
    debugLog('Token:', !!token, 'User:', !!userParam)

    if (token && userParam) {
      try {
        const user = JSON.parse(decodeURIComponent(userParam))
        debugLog('Saving auth for user:', user.email)

        // Save auth to store
        setAuth(user, token)
        void warmupCurrentUserPermissions()

        // Extract locale from pathname or use default
        const pathSegments = pathname.split('/')
        const locale = ['en', 'ar'].includes(pathSegments[1]) ? pathSegments[1] : 'en'

        debugLog('Redirecting to dashboard for locale:', locale)
        // Redirect to dashboard with proper locale prefix
        window.location.href = `${getConfig().dashboardUrl}/${locale}/dashboard`
      } catch (error) {
        debugError('Auth callback error:', error)
        window.location.href = getDashboardLoginUrl(pathname)
      }
    } else {
      debugLog('No auth data, redirecting to login')
      window.location.href = getDashboardLoginUrl(pathname)
    }
  }, [searchParams, setAuth, pathname])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
  )
}
