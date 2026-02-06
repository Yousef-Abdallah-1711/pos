'use client'

import { ReactNode, useEffect, useState, useRef } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { useTenantStatusStore } from '@/stores/tenantStatusStore';
import { getConfig } from '@/lib/config';
import { BlockedOverlay } from '@/components/blocked/BlockedOverlay';

const DashboardAuthProvider = ({ children }: { children: ReactNode }) => {
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const setAuth = useAuthStore((state) => state.setAuth)
  const canAccess = useTenantStatusStore((state) => state.canAccess)
  const blockReason = useTenantStatusStore((state) => state.blockReason)
  const checkStatus = useTenantStatusStore((state) => state.checkStatus)
  const [isHydrated, setIsHydrated] = useState(false)
  const [urlAuthProcessed, setUrlAuthProcessed] = useState(false)
  const urlAuthApplied = useRef(false) // Track if URL auth was applied this session
  const impersonationSuccess = useRef(false) // Track if impersonation was successful

  // Check if authenticated by checking if user and token exist
  const isAuthenticated = !!(user && token)

  console.log('[DashboardAuthProvider] Render - user:', !!user, 'token:', !!token, 'isAuthenticated:', isAuthenticated, 'isHydrated:', isHydrated, 'urlAuthProcessed:', urlAuthProcessed)

  // Wait for Zustand to hydrate from localStorage on mount
  useEffect(() => {
    // Check if we're in the browser
    if (typeof window === 'undefined') {
      return
    }

    console.log('[DashboardAuthProvider] First useEffect running')

    const processAuth = async () => {
      // First, try to initialize from URL parameters (for super admin login redirect)
      const params = new URLSearchParams(window.location.search)
      const authToken = params.get('auth_token')
      const authUserStr = params.get('auth_user')
      const impersonationToken = params.get('impersonation_token')

      console.log('[DashboardAuthProvider] URL params - authToken:', !!authToken, 'impersonationToken:', !!impersonationToken)

      // Handle regular auth token + user
      if (authToken && authUserStr) {
        try {
          const authUser = JSON.parse(decodeURIComponent(authUserStr))
          console.log('[DashboardAuthProvider] Processing regular auth token')

          // Mark that we applied URL auth this session
          urlAuthApplied.current = true

          // Call setAuth to update the Zustand store
          // The persist middleware will automatically save to localStorage
          setAuth(authUser, authToken)

          // Wait a moment for Zustand to save to localStorage, then clean URL
          await new Promise(resolve => setTimeout(resolve, 100))
          const cleanUrl = window.location.pathname
          window.history.replaceState({}, document.title, cleanUrl)
          console.log('[DashboardAuthProvider] Regular auth token processed successfully')
        } catch (error) {
          // Error parsing auth params - continue with normal flow
          console.error('[DashboardAuthProvider] Failed to process auth token:', error)
        }
      }
      // Handle impersonation token
      else if (impersonationToken) {
        console.log('[DashboardAuthProvider] Processing impersonation token...')
        try {
          const apiUrl = getConfig().apiUrl
          console.log('[DashboardAuthProvider] Validating token at:', `${apiUrl}/impersonation/validate`)
          const response = await fetch(`${apiUrl}/impersonation/validate?token=${encodeURIComponent(impersonationToken)}`)

          console.log('[DashboardAuthProvider] Validation response status:', response.status)

          if (response.ok) {
            const result = await response.json()
            console.log('[DashboardAuthProvider] Validation result:', result)

            if (result.success && result.data) {
              const { user: impersonatedUser } = result.data
              console.log('[DashboardAuthProvider] Impersonated user:', impersonatedUser.email)

              // Mark that we applied URL auth this session
              urlAuthApplied.current = true
              impersonationSuccess.current = true

              // Store impersonation token in sessionStorage for API requests
              sessionStorage.setItem('impersonation_token', impersonationToken)
              console.log('[DashboardAuthProvider] Stored impersonation token in sessionStorage')

              // Create a dummy token for auth store (needed for UI to show logged in)
              // The actual auth will happen via the impersonation token header
              setAuth(impersonatedUser, 'impersonated')
              console.log('[DashboardAuthProvider] Called setAuth with impersonated user')

              // Clean URL
              await new Promise(resolve => setTimeout(resolve, 100))
              const cleanUrl = window.location.pathname
              window.history.replaceState({}, document.title, cleanUrl)
              console.log('[DashboardAuthProvider] URL cleaned, impersonation complete')
            } else {
              console.error('[DashboardAuthProvider] Validation result missing success or data:', result)
            }
          } else {
            console.error('[DashboardAuthProvider] Token validation failed with status:', response.status)
            const errorText = await response.text()
            console.error('[DashboardAuthProvider] Error response:', errorText)
          }
        } catch (error) {
          console.error('[DashboardAuthProvider] Failed to validate impersonation token:', error)
          // Continue with normal flow
        }
      } else {
        console.log('[DashboardAuthProvider] No URL auth params found')
      }
    }

    const initAuth = async () => {
      console.log('[DashboardAuthProvider] initAuth starting')
      // Process URL auth params if any
      await processAuth()

      // Wait longer for Zustand persist middleware to write to localStorage
      // This is crucial for the API interceptor to pick up the token
      console.log('[DashboardAuthProvider] Waiting 200ms for Zustand to persist...')
      await new Promise(resolve => setTimeout(resolve, 200))

      // Then mark as processed - this will trigger the second useEffect
      console.log('[DashboardAuthProvider] Setting urlAuthProcessed = true')
      setUrlAuthProcessed(true)

      // Wait for hydration state to update
      await new Promise(resolve => setTimeout(resolve, 100))
      console.log('[DashboardAuthProvider] Setting isHydrated = true')
      setIsHydrated(true)

      console.log('[DashboardAuthProvider] initAuth complete - urlAuthApplied:', urlAuthApplied.current, 'impersonationSuccess:', impersonationSuccess.current)
    }

    initAuth()

    return () => {
      // Cleanup if needed
    }
  }, [setAuth])

  useEffect(() => {
    console.log('[DashboardAuthProvider] Second useEffect - isHydrated:', isHydrated, 'urlAuthProcessed:', urlAuthProcessed)

    // Only run regular auth check after hydration and URL params processed
    if (!isHydrated || !urlAuthProcessed) {
      console.log('[DashboardAuthProvider] Not ready yet, skipping auth check')
      return
    }

    console.log('[DashboardAuthProvider] Auth check - urlAuthApplied:', urlAuthApplied.current, 'impersonationSuccess:', impersonationSuccess.current, 'user:', !!user, 'token:', token)

    // Check tenant status after successful auth
    if (user && token) {
      console.log('[DashboardAuthProvider] Authenticated, checking tenant status')
      checkStatus()
    }

    // If impersonation was successful, don't redirect
    if (impersonationSuccess.current) {
      console.log('[DashboardAuthProvider] Impersonation was successful, staying authenticated')
      return
    }

    // If URL auth was just applied this session AND we have auth in Zustand state, we're good
    // This avoids checking localStorage which might not have been written yet
    if (urlAuthApplied.current) {
      if (user && token) {
        // Auth from URL params was successfully applied - stay authenticated
        console.log('[DashboardAuthProvider] URL auth applied and have user/token, staying authenticated')
        return
      } else {
        // URL auth was applied but failed - redirect to login
        console.log('[DashboardAuthProvider] URL auth applied but NO user/token, redirecting to login')
        const landingUrl = getConfig().landingUrl
        const redirectUrl = `${landingUrl}/auth/login`
        window.location.href = redirectUrl
        return
      }
    }

    // If no URL auth was applied, check localStorage for existing session
    const savedAuth = localStorage.getItem('auth-storage')
    console.log('[DashboardAuthProvider] Checking localStorage, savedAuth exists:', !!savedAuth)

    let localStorageAuth = null
    if (savedAuth) {
      try {
        const parsed = JSON.parse(savedAuth)
        const hasUser = !!parsed.state?.user
        const hasToken = !!parsed.state?.token
        console.log('[DashboardAuthProvider] localStorage parsed - hasUser:', hasUser, 'hasToken:', hasToken)

        // If localStorage has auth, use that as source of truth
        if (hasUser && hasToken) {
          localStorageAuth = parsed.state
        }
      } catch (e) {
        console.error('[DashboardAuthProvider] Error parsing localStorage:', e)
      }
    }

    // Determine final authentication state
    // Use Zustand state if available, otherwise use localStorage
    const finalUser = user || localStorageAuth?.user
    const finalToken = token || localStorageAuth?.token
    const finalIsAuthenticated = !!(finalUser && finalToken)

    console.log('[DashboardAuthProvider] Final auth check - finalUser:', !!finalUser, 'finalToken:', !!finalToken, 'finalIsAuthenticated:', finalIsAuthenticated)

    if (!finalIsAuthenticated) {
      console.log('[DashboardAuthProvider] NOT authenticated, redirecting to login')
      const landingUrl = getConfig().landingUrl
      const redirectUrl = `${landingUrl}/auth/login`
      window.location.href = redirectUrl
    } else {
      console.log('[DashboardAuthProvider] Authenticated, staying on dashboard')
    }
  }, [isHydrated, urlAuthProcessed, user, token, checkStatus])

  // If not hydrated or not authenticated, show loading while redirecting
  // But if impersonation was successful, we should be authenticated
  const shouldShowLoading = !isHydrated || !urlAuthProcessed || (!isAuthenticated && !impersonationSuccess.current)

  console.log('[DashboardAuthProvider] Render check - shouldShowLoading:', shouldShowLoading, 'isHydrated:', isHydrated, 'urlAuthProcessed:', urlAuthProcessed, 'isAuthenticated:', isAuthenticated, 'impersonationSuccess:', impersonationSuccess.current)

  if (shouldShowLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  // If tenant is blocked, show blocked overlay
  if (!canAccess && blockReason) {
    console.log('[DashboardAuthProvider] Tenant is blocked, showing overlay')
    return <BlockedOverlay blockReason={blockReason} />
  }

  // User is authenticated and hydrated, show dashboard
  console.log('[DashboardAuthProvider] Showing dashboard content')
  return <>{children}</>
};

export default DashboardAuthProvider;
