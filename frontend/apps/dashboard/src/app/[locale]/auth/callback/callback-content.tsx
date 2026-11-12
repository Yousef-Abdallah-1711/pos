'use client'

import { useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'

export default function AuthCallbackContent() {
  const pathname = usePathname()
  const setAuth = useAuthStore((state) => state.setAuth)

  useEffect(() => {
    let active = true
    const locale = pathname.split('/')[1] === 'ar' ? 'ar' : 'en'

    fetch(`${getConfig().apiUrl}/auth/session`, { credentials: 'include', cache: 'no-store' })
      .then(async (response) => {
        if (!response.ok) throw new Error('Session unavailable')
        return response.json()
      })
      .then((user: unknown) => {
        if (!active) return
        setAuth(user as Parameters<typeof setAuth>[0], '')
        window.history.replaceState({}, document.title, pathname)
        // This callback crosses the landing/dashboard app boundary. A document
        // navigation avoids carrying a stale App Router tree from the callback.
        window.location.replace(`/${locale}/dashboard`)
      })
      .catch(() => {
        if (active) window.location.replace(`${getConfig().landingUrl}/${locale}/auth/login`)
      })

    return () => { active = false }
  }, [pathname, setAuth])

  return (
    <div className="min-h-screen flex items-center justify-center" aria-live="polite">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
    </div>
  )
}
