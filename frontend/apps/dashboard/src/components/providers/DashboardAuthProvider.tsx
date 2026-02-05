'use client'

import { ReactNode, useEffect, useState } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { getConfig } from '@/lib/config';

const DashboardAuthProvider = ({ children }: { children: ReactNode }) => {
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const setAuth = useAuthStore((state) => state.setAuth)
  const [isHydrated, setIsHydrated] = useState(false)
  const [urlAuthProcessed, setUrlAuthProcessed] = useState(false)

  // Check if authenticated by checking if user and token exist
  const isAuthenticated = !!(user && token)

  // Wait for Zustand to hydrate from localStorage on mount
  useEffect(() => {
    // Check if we're in the browser
    if (typeof window === 'undefined') {
      return
    }

    // First, try to initialize from URL parameters (for super admin login redirect)
    const params = new URLSearchParams(window.location.search)
    const authToken = params.get('auth_token')
    const authUserStr = params.get('auth_user')

    if (authToken && authUserStr) {
      try {
        const authUser = JSON.parse(decodeURIComponent(authUserStr))

        // Call setAuth to update the Zustand store
        // The persist middleware will automatically save to localStorage
        setAuth(authUser, authToken)

        // Wait a moment for Zustand to save to localStorage, then clean URL
        setTimeout(() => {
          // Remove auth params from URL for cleanliness
          const cleanUrl = window.location.pathname
          window.history.replaceState({}, document.title, cleanUrl)
        }, 100)
      } catch (error) {
        // Error parsing auth params - continue with normal flow
      }
    }

    setUrlAuthProcessed(true)

    // Wait for Zustand persist middleware to hydrate from localStorage
    // This is crucial for the store to load saved auth data
    const timer = setTimeout(() => {
      setIsHydrated(true)
    }, 100)

    return () => {
      clearTimeout(timer)
    }
  }, [setAuth])

  useEffect(() => {
    // Only run regular auth check after hydration and URL params processed
    if (!isHydrated || !urlAuthProcessed) {
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
