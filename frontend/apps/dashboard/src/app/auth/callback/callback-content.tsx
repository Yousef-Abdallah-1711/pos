'use client'

import { useEffect } from 'react'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'

export default function AuthCallbackContent() {
  const setAuth = useAuthStore((state) => state.setAuth)

  useEffect(() => {
    let active = true
    fetch(`${getConfig().apiUrl}/auth/session`, { credentials: 'include', cache: 'no-store' })
      .then(async (response) => {
        if (!response.ok) throw new Error('Session unavailable')
        return response.json()
      })
      .then((user: unknown) => {
        if (!active) return
        setAuth(user as Parameters<typeof setAuth>[0], '')
        window.history.replaceState({}, document.title, '/auth/callback')
        window.location.replace('/en/dashboard')
      })
      .catch(() => {
        if (active) window.location.replace(`${getConfig().landingUrl}/en/auth/login`)
      })

    return () => { active = false }
  }, [setAuth])

  return (
    <div className="min-h-screen flex items-center justify-center" aria-live="polite">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
    </div>
  )
}
