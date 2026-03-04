'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams, usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { warmupCurrentUserPermissions } from '@/hooks/useRbac'
import { getConfig } from '@/lib/config'
import { debugLog, debugError } from '@/lib/debug'

export default function AuthCallbackContent() {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const setAuth = useAuthStore((state) => state.setAuth)

  // Extract locale from pathname (e.g., /en/auth/callback -> 'en')
  const locale = pathname.split('/')[1] || 'en'

  useEffect(() => {
    const token = searchParams.get('token')
    const userParam = searchParams.get('user')

    debugLog('=== CALLBACK: Processing auth ===')
    debugLog('Token:', !!token, 'User:', !!userParam)
    debugLog('Locale:', locale)

    if (token && userParam) {
      try {
        const user = JSON.parse(decodeURIComponent(userParam))
        debugLog('Saving auth for user:', user.email)

        // Save auth to store
        setAuth(user, token)
        void warmupCurrentUserPermissions()

        debugLog('Redirecting to dashboard...')
        // Redirect to localized dashboard
        router.replace(`/${locale}/dashboard`)
      } catch (error) {
        debugError('Auth callback error:', error)
        window.location.href = `${getConfig().landingUrl}/auth/login`
      }
    } else {
      debugLog('No auth data, redirecting to login')
      window.location.href = `${getConfig().landingUrl}/auth/login`
    }
  }, [searchParams, setAuth, router, locale])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
  )
}
