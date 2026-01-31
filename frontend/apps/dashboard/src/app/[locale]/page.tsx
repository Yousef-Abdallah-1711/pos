'use client'

import { useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'

export default function LocaleHomePage() {
  const router = useRouter()
  const pathname = usePathname()
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)

  // Extract locale from pathname (e.g., /en -> 'en')
  const pathParts = pathname.split('/').filter(Boolean)
  const locale = pathParts.length > 0 ? pathParts[0] : 'en'

  // Check if this is a valid locale (simple check - just lowercase letters)
  const isValidLocale = /^[a-z]{2}$/.test(locale)

  // Check if authenticated by checking if user and token exist
  const isAuthenticated = !!(user && token)

  useEffect(() => {
    // Don't redirect if locale is invalid (e.g., /dashboard was interpreted as locale)
    if (!isValidLocale) {
      window.location.href = `${getConfig().dashboardUrl}/en/dashboard`
      return
    }

    // Check for debug mode
    const debugMode = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('debug') === 'true'
    if (debugMode) {
      router.replace(`/${locale}/dashboard${window.location.search}`)
      return
    }

    if (isAuthenticated) {
      router.replace(`/${locale}/dashboard`)
    } else {
      window.location.href = `${getConfig().landingUrl}/auth/login`
    }
  }, [isAuthenticated, user, router, locale, isValidLocale])

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
  )
}
