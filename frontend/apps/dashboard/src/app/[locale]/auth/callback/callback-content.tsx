'use client'

import { useEffect, useRef } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { requestDashboardCallbackSession } from '@/lib/authSession'
import { getConfig } from '@/lib/config'

export default function AuthCallbackContent() {
  const pathname = usePathname()
  const setAuth = useAuthStore((state) => state.setAuth)
  const redirectedRef = useRef(false)

  useEffect(() => {
    let active = true
    const locale = pathname.split('/')[1] === 'ar' ? 'ar' : 'en'

    void requestDashboardCallbackSession(getConfig().apiUrl)
      .then((user) => {
        if (!active || redirectedRef.current) return
        redirectedRef.current = true
        setAuth(user, '')
        window.history.replaceState({}, document.title, pathname)
        // This callback crosses the landing/dashboard app boundary. A document
        // navigation avoids carrying a stale App Router tree from the callback.
        window.location.replace(`/${locale}/dashboard`)
      })
      .catch(() => {
        if (!active || redirectedRef.current) return
        redirectedRef.current = true
        window.location.replace(`${getConfig().landingUrl}/${locale}/auth/login`)
      })

    return () => { active = false }
  }, [pathname, setAuth])

  return (
    <div className="min-h-screen flex items-center justify-center" aria-live="polite">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
    </div>
  )
}
