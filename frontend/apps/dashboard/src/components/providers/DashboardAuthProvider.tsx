'use client'

import { ReactNode, useEffect, useState } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { getConfig } from '@/lib/config';

// Helper to write debug logs to localStorage
const writeDebugLog = (msg: string) => {
  try {
    const logs = JSON.parse(localStorage.getItem('debug_logs') || '[]')
    logs.push({ time: new Date().toISOString(), msg })
    localStorage.setItem('debug_logs', JSON.stringify(logs.slice(-20))) // Keep last 20 logs
  } catch (e) {
    // Ignore
  }
}

const DashboardAuthProvider = ({ children }: { children: ReactNode }) => {
  // VERY FIRST THING: Check debug mode immediately
  const debugModeActive = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('debug') === 'true'

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

    writeDebugLog('MOUNT EFFECT STARTED')

    // First, try to initialize from URL parameters (for super admin login redirect)
    writeDebugLog('Checking for auth params in URL...')
    const params = new URLSearchParams(window.location.search)
    const authToken = params.get('auth_token')
    const authUserStr = params.get('auth_user')

    writeDebugLog(`auth_token present: ${!!authToken}, auth_user present: ${!!authUserStr}`)

    let foundUrlAuth = false

    if (authToken && authUserStr) {
      writeDebugLog('✅ Found auth params in URL!')
      try {
        const authUser = JSON.parse(decodeURIComponent(authUserStr))
        writeDebugLog(`Parsed user: ${authUser.email}`)
        writeDebugLog(`Token length: ${authToken.length}`)

        // Store in localStorage in the exact format Zustand expects
        const zustandData = {
          state: {
            user: authUser,
            token: authToken
          },
          version: 0
        }
        localStorage.setItem('auth-storage', JSON.stringify(zustandData))
        writeDebugLog('✅ Auth stored in localStorage')

        // Also call setAuth to update the Zustand store in memory
        setAuth(authUser, authToken)
        writeDebugLog('✅ Called setAuth()')

        // Remove auth params from URL for cleanliness, but preserve debug parameter
        const hasDebug = new URLSearchParams(window.location.search).get('debug') === 'true'
        const cleanUrl = window.location.pathname + (hasDebug ? '?debug=true' : '')
        window.history.replaceState({}, document.title, cleanUrl)
        writeDebugLog('✅ URL cleaned, debug param preserved')
        foundUrlAuth = true
      } catch (error) {
        writeDebugLog(`❌ Error parsing auth params: ${error}`)
      }
    } else {
      writeDebugLog('No auth params in URL, checking localStorage...')
      const savedAuth = localStorage.getItem('auth-storage')
      writeDebugLog(`Auth from localStorage: ${savedAuth ? 'exists' : 'not found'}`)
    }

    setUrlAuthProcessed(true)
    writeDebugLog('urlAuthProcessed = true')

    // Set a small delay to ensure Zustand has fully hydrated
    const timer = setTimeout(() => {
      setIsHydrated(true)
    }, 50)

    return () => {
      clearTimeout(timer)
    }
  }, [setAuth])

  useEffect(() => {
    // ALWAYS run this immediately to process URL params and detect debug mode
    // Don't wait for hydration if we're in debug mode
    if (typeof window !== 'undefined') {
      const urlParams = new URLSearchParams(window.location.search)
      const debugParam = urlParams.get('debug')
      const isDebugMode = debugParam === 'true'

      writeDebugLog(`Debug mode: ${isDebugMode}`)

      // If in debug mode, skip auth check entirely
      if (isDebugMode) {
        writeDebugLog('⏭️  Skipping auth check - debug mode active')
        return
      }
    }

    // Only run regular auth check after hydration and URL params processed
    if (!isHydrated || !urlAuthProcessed) {
      writeDebugLog(`Waiting... isHydrated=${isHydrated}, urlAuthProcessed=${urlAuthProcessed}`)
      return
    }

    writeDebugLog('========== AUTH CHECK (FINAL) ==========')
    writeDebugLog(`isHydrated=${isHydrated}, user=${user?.email || 'null'}, token=${token ? 'present' : 'null'}, isAuthenticated=${isAuthenticated}`)

    // Check localStorage directly as it's the source of truth
    const savedAuth = localStorage.getItem('auth-storage')
    writeDebugLog(`localStorage auth-storage exists: ${!!savedAuth}`)

    let localStorageAuth = null
    if (savedAuth) {
      try {
        const parsed = JSON.parse(savedAuth)
        const hasUser = !!parsed.state?.user
        const hasToken = !!parsed.state?.token
        writeDebugLog(`localStorage content: hasUser=${hasUser}, hasToken=${hasToken}, email=${parsed.state?.user?.email}`)

        // If localStorage has auth, use that as source of truth
        if (hasUser && hasToken) {
          localStorageAuth = parsed.state
          writeDebugLog('✅ Found valid auth in localStorage')
        }
      } catch (e) {
        writeDebugLog(`Error parsing localStorage: ${e}`)
      }
    }

    // Determine final authentication state
    // Use Zustand state if available, otherwise use localStorage
    const finalUser = user || localStorageAuth?.user
    const finalToken = token || localStorageAuth?.token
    const finalIsAuthenticated = !!(finalUser && finalToken)

    writeDebugLog(`Final auth: user=${finalUser?.email || 'null'}, token=${finalToken ? 'present' : 'null'}, authenticated=${finalIsAuthenticated}`)

    if (!finalIsAuthenticated) {
      // Double-check debug mode before redirecting
      const urlParams = new URLSearchParams(window.location.search)
      const debugParam = urlParams.get('debug')
      if (debugParam === 'true') {
        writeDebugLog('⚠️  Debug mode detected - NOT redirecting to login')
        return
      }

      writeDebugLog('❌ NOT AUTHENTICATED - Redirecting to login')
      const landingUrl = getConfig().landingUrl
      const redirectUrl = `${landingUrl}/auth/login`
      writeDebugLog(`Redirecting to: ${redirectUrl}`)
      // Use setTimeout to ensure state is fully updated before redirect
      setTimeout(() => {
        window.location.href = redirectUrl
      }, 100)
    } else {
      writeDebugLog('✅ AUTHENTICATED - Allowing access to dashboard')
    }
  }, [isHydrated, urlAuthProcessed, isAuthenticated, user, token])

  // Debug mode: show debug logs if ?debug=true in URL
  if (typeof window !== 'undefined') {
    const urlParams = new URLSearchParams(window.location.search)
    const debugParam = urlParams.get('debug')
    const isDebugMode = debugParam === 'true'

    if (isDebugMode) {
      // Add a visual marker to show debug mode is active
      if (typeof document !== 'undefined') {
        document.body.style.backgroundColor = '#111827' // dark gray
      }

      let logs = []
      try {
        logs = JSON.parse(localStorage.getItem('debug_logs') || '[]')
      } catch (e) {
        logs = []
      }
      return (
        <div className="min-h-screen bg-gray-900 text-white p-8 font-mono text-sm">
          <div className="max-w-6xl mx-auto">
            <h1 className="text-2xl font-bold mb-6">Debug Logs</h1>

            {/* URL and Params Info */}
            <div className="bg-gray-800 p-4 rounded mb-6 border border-gray-600">
              <h2 className="text-lg font-bold mb-3 text-green-400">Current URL Info</h2>
              <div className="space-y-2 text-sm">
                <div>
                  <span className="text-gray-400">Full URL:</span>
                  <p className="text-blue-300 break-all font-mono text-xs mt-1">{window.location.href}</p>
                </div>
                <div>
                  <span className="text-gray-400">Search Params:</span>
                  <p className="text-yellow-300 font-mono text-xs mt-1">{window.location.search || '(none)'}</p>
                </div>
                <div>
                  <span className="text-gray-400">Has auth_token:</span>
                  <span className={new URLSearchParams(window.location.search).get('auth_token') ? 'text-green-400' : 'text-red-400'}>
                    {new URLSearchParams(window.location.search).get('auth_token') ? '✅ Yes' : '❌ No'}
                  </span>
                </div>
                <div>
                  <span className="text-gray-400">Has auth_user:</span>
                  <span className={new URLSearchParams(window.location.search).get('auth_user') ? 'text-green-400' : 'text-red-400'}>
                    {new URLSearchParams(window.location.search).get('auth_user') ? '✅ Yes' : '❌ No'}
                  </span>
                </div>
                <div>
                  <span className="text-gray-400">Has debug param:</span>
                  <span className={new URLSearchParams(window.location.search).get('debug') === 'true' ? 'text-green-400' : 'text-red-400'}>
                    {new URLSearchParams(window.location.search).get('debug') === 'true' ? '✅ Yes' : '❌ No'}
                  </span>
                </div>
              </div>
            </div>

            <div className="bg-gray-800 p-6 rounded mb-6 max-h-96 overflow-y-auto">
              <table className="w-full text-left">
                <thead className="border-b border-gray-700 mb-4">
                  <tr>
                    <th className="pb-2">Time</th>
                    <th className="pb-2">Message</th>
                  </tr>
                </thead>
                <tbody className="space-y-1">
                  {logs.length === 0 ? (
                    <tr>
                      <td colSpan={2} className="text-gray-400 py-4">
                        No debug logs available. Loading...
                      </td>
                    </tr>
                  ) : (
                    logs.map((log: any, i: number) => (
                      <tr key={i} className="hover:bg-gray-700">
                        <td className="text-gray-400 min-w-fit pr-4">{new Date(log.time).toLocaleTimeString()}</td>
                        <td className={log.msg.includes('❌') ? 'text-red-400' : log.msg.includes('✅') ? 'text-green-400' : 'text-gray-200'}>
                          {log.msg}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            <div className="space-y-3">
              <button
                onClick={() => {
                  window.location.href = window.location.href.replace('?debug=true', '')
                }}
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded font-medium"
              >
                Try Dashboard
              </button>
              <button
                onClick={() => {
                  localStorage.removeItem('debug_logs')
                  window.location.reload()
                }}
                className="px-6 py-2 bg-gray-600 hover:bg-gray-700 rounded font-medium ml-2"
              >
                Clear & Reload
              </button>
              <p className="text-gray-400 text-xs mt-4">
                Status: isHydrated={isHydrated}, urlAuthProcessed={urlAuthProcessed}, isAuthenticated={isAuthenticated}
              </p>
            </div>
          </div>
        </div>
      )
    }
  }

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
