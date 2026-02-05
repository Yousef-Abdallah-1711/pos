'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams, usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'

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

    console.log('=== CALLBACK: Processing auth ===')
    console.log('Token:', !!token, 'User:', !!userParam)
    console.log('Locale:', locale)

    if (token && userParam) {
      try {
        const user = JSON.parse(decodeURIComponent(userParam))
        console.log('Saving auth for user:', user.email)

        // Save auth to store
        setAuth(user, token)

        console.log('Redirecting to dashboard...')
        // Redirect to localized dashboard
        router.replace(`/${locale}/dashboard`)
      } catch (error) {
        console.error('Auth callback error:', error)
        window.location.href = `${getConfig().landingUrl}/auth/login`
      }
    } else {
      console.log('No auth data, redirecting to login')
      window.location.href = `${getConfig().landingUrl}/auth/login`
    }
  }, [searchParams, setAuth, router, locale])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
  )
}
