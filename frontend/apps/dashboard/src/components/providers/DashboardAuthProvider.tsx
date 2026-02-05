'use client'

import { ReactNode, useEffect, useState, useRef } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { getConfig } from '@/lib/config';

const DashboardAuthProvider = ({ children }: { children: ReactNode }) => {
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const setAuth = useAuthStore((state) => state.setAuth)
  const [isHydrated, setIsHydrated] = useState(false)
  const [urlAuthProcessed, setUrlAuthProcessed] = useState(false)
  const urlAuthApplied = useRef(false) // Track if URL auth was applied this session

  // Check if authenticated by checking if user and token exist
  const isAuthenticated = !!(user && token)

  // Wait for Zustand to hydrate from localStorage on mount
  useEffect(() => {
    // Check if we're in the browser
    if (typeof window === 'undefined') {
      return
    }

    const processAuth = async () => {
      // First, try to initialize from URL parameters (for super admin login redirect)
      const params = new URLSearchParams(window.location.search)
      const authToken = params.get('auth_token')
      const authUserStr = params.get('auth_user')
      const impersonationToken = params.get('impersonation_token')

      // Handle regular auth token + user
      if (authToken && authUserStr) {
        try {
          const authUser = JSON.parse(decodeURIComponent(authUserStr))

          // Mark that we applied URL auth this session
          urlAuthApplied.current = true

          // Call setAuth to update the Zustand store
          // The persist middleware will automatically save to localStorage
          setAuth(authUser, authToken)

          // Wait a moment for Zustand to save to localStorage, then clean URL
          await new Promise(resolve => setTimeout(resolve, 100))
          const cleanUrl = window.location.pathname
          window.history.replaceState({}, document.title, cleanUrl)
        } catch (error) {
          // Error parsing auth params - continue with normal flow
          console.error('Failed to process auth token:', error)
        }
      }
      // Handle impersonation token
      else if (impersonationToken) {
        try {
          const apiUrl = getConfig().apiUrl
          const response = await fetch(`${apiUrl}/impersonation/validate?token=${encodeURIComponent(impersonationToken)}`)

          if (response.ok) {
            const result = await response.json()
            if (result.success && result.data) {
              const { user: impersonatedUser, session_id } = result.data

              // Mark that we applied URL auth this session
              urlAuthApplied.current = true

              // Use the session_id as the token for impersonation
              setAuth(impersonatedUser, session_id)

              // Clean URL
              await new Promise(resolve => setTimeout(resolve, 100))
              const cleanUrl = window.location.pathname
              window.history.replaceState({}, document.title, cleanUrl)
            }
          } else {
            console.error('Token validation failed with status:', response.status)
          }
        } catch (error) {
          console.error('Failed to validate impersonation token:', error)
          // Continue with normal flow
        }
      }
    }

    const initAuth = async () => {
      // Process URL auth params if any
      await processAuth()

      // Then mark as processed
      setUrlAuthProcessed(true)

      // Wait for Zustand persist middleware to hydrate from localStorage
      // This is crucial for the store to load saved auth data
      await new Promise(resolve => setTimeout(resolve, 100))
      setIsHydrated(true)
    }

    initAuth()

    return () => {
      // Cleanup if needed
    }
  }, [setAuth])

  useEffect(() => {
    // Only run regular auth check after hydration and URL params processed
    if (!isHydrated || !urlAuthProcessed) {
      return
    }

    // If URL auth was just applied this session, trust Zustand state (not localStorage)
    // This avoids race conditions where localStorage might have stale data
    if (urlAuthApplied.current && user && token) {
      // Auth from URL params was successfully applied
      return
    }

    // Check localStorage directly as it's the source of truth
    const savedAuth = localStorage.getItem('auth-storage')

    let localStorageAuth = null
    if (savedAuth) {
      try {
        const parsed = JSON.parse(savedAuth)
        const hasUser = !!parsed.state?.user
        const hasToken = !!parsed.state?.token

        // If localStorage has auth, use that as source of truth
        if (hasUser && hasToken) {
          localStorageAuth = parsed.state
        }
      } catch (e) {
        // Error parsing localStorage
      }
    }

    // Determine final authentication state
    // Use Zustand state if available, otherwise use localStorage
    const finalUser = user || localStorageAuth?.user
    const finalToken = token || localStorageAuth?.token
    const finalIsAuthenticated = !!(finalUser && finalToken)

    if (!finalIsAuthenticated) {
      const landingUrl = getConfig().landingUrl
      const redirectUrl = `${landingUrl}/auth/login`
      // Use setTimeout to ensure state is fully updated before redirect
      setTimeout(() => {
        window.location.href = redirectUrl
      }, 100)
    }
  }, [isHydrated, urlAuthProcessed, isAuthenticated, user, token])

  // If not hydrated or not authenticated, show loading while redirecting
  if (!isHydrated || !urlAuthProcessed || !isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  // User is authenticated and hydrated, show dashboard
  return <>{children}</>
};

export default DashboardAuthProvider;
