'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import {
  cancelLandingSessionRequest,
  discoverLandingSession,
  signOutLandingSession,
  subscribeToSessionChanges,
  type LandingSessionState,
} from '@/lib/auth-session'

export function useLandingSession(pathname: string) {
  const [session, setSession] = useState<LandingSessionState>({ status: 'checking' })
  const [isSigningOut, setIsSigningOut] = useState(false)
  const mountedRef = useRef(false)
  const requestGenerationRef = useRef(0)

  const revalidate = useCallback(async (force = false) => {
    const requestGeneration = ++requestGenerationRef.current
    setSession({ status: 'checking' })
    const result = await discoverLandingSession(force)
    if (mountedRef.current && requestGeneration === requestGenerationRef.current) setSession(result)
  }, [])

  const retry = useCallback(() => revalidate(true), [revalidate])

  const signOut = useCallback(async () => {
    requestGenerationRef.current += 1
    setIsSigningOut(true)
    setSession({ status: 'checking' })
    try {
      await signOutLandingSession()
      if (mountedRef.current) setSession({ status: 'guest' })
      return true
    } catch {
      if (mountedRef.current) setSession({ status: 'temporarily-unavailable', reason: 'sign-out' })
      return false
    } finally {
      if (mountedRef.current) setIsSigningOut(false)
    }
  }, [])

  useEffect(() => {
    mountedRef.current = true
    const initialCheck = window.setTimeout(() => { void revalidate(true) }, 0)
    return () => {
      window.clearTimeout(initialCheck)
      mountedRef.current = false
      requestGenerationRef.current += 1
      cancelLandingSessionRequest()
    }
  }, [revalidate])

  useEffect(() => {
    const navigationCheck = window.setTimeout(() => { void revalidate() }, 0)
    return () => window.clearTimeout(navigationCheck)
  }, [pathname, revalidate])

  useEffect(() => {
    let refreshTimeout: number | undefined
    let wasInactive = document.visibilityState !== 'visible'
    const refreshWhenVisible = () => {
      if (document.visibilityState !== 'visible') {
        wasInactive = true
        return
      }
      if (!wasInactive) return
      wasInactive = false
      window.clearTimeout(refreshTimeout)
      refreshTimeout = window.setTimeout(() => { void revalidate(true) }, 100)
    }
    const markInactive = () => { wasInactive = true }
    const unsubscribe = subscribeToSessionChanges(() => { void revalidate(true) })
    window.addEventListener('focus', refreshWhenVisible)
    window.addEventListener('blur', markInactive)
    document.addEventListener('visibilitychange', refreshWhenVisible)

    return () => {
      unsubscribe()
      window.clearTimeout(refreshTimeout)
      window.removeEventListener('focus', refreshWhenVisible)
      window.removeEventListener('blur', markInactive)
      document.removeEventListener('visibilitychange', refreshWhenVisible)
    }
  }, [revalidate])

  return { session, retry, signOut, isSigningOut }
}
