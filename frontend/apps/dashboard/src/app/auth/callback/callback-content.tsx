'use client'

import { useEffect, useRef } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { requestDashboardCallbackSession } from '@/lib/authSession'
import { getConfig } from '@/lib/config'

export default function AuthCallbackContent() {
  const setAuth = useAuthStore((state) => state.setAuth)
  const redirectedRef = useRef(false)

  useEffect(() => {
    let active = true
    void requestDashboardCallbackSession(getConfig().apiUrl)
      .then((user) => {
        if (!active || redirectedRef.current) return
        redirectedRef.current = true
        setAuth(user, '')
        window.history.replaceState({}, document.title, '/auth/callback')
        window.location.replace('/en/dashboard')
      })
      .catch(() => {
        if (!active || redirectedRef.current) return
        redirectedRef.current = true
        window.location.replace(`${getConfig().landingUrl}/en/auth/login`)
      })

    return () => { active = false }
  }, [setAuth])

  return (
    <div className="min-h-screen flex items-center justify-center" aria-live="polite">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
    </div>
  )
}
