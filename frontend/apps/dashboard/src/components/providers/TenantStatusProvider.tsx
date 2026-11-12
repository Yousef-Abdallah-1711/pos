'use client'

import { ReactNode, useEffect } from 'react'
import { useTenantStatusStore } from '@/stores/tenantStatusStore'
import { debugLog } from '@/lib/debug'

interface TenantStatusProviderProps {
  children: ReactNode
  checkInterval?: number // in milliseconds, default 5 minutes
}

export function TenantStatusProvider({
  children,
  checkInterval = 5 * 60 * 1000, // 5 minutes
}: TenantStatusProviderProps) {
  const checkStatus = useTenantStatusStore((state) => state.checkStatus)

  // Check status on mount and periodically
  useEffect(() => {
    debugLog('[TenantStatusProvider] Mounting - checking tenant status')
    let intervalId: ReturnType<typeof setInterval> | null = null

    const clearPolling = () => {
      if (intervalId) {
        clearInterval(intervalId)
        intervalId = null
      }
    }

    const runCheck = () => {
      debugLog('[TenantStatusProvider] Checking tenant status')
      void checkStatus()
    }

    const startPolling = () => {
      clearPolling()
      if (document.visibilityState !== 'visible') {
        return
      }

      intervalId = setInterval(() => {
        if (document.visibilityState !== 'visible') {
          return
        }
        debugLog('[TenantStatusProvider] Periodic check - checking tenant status')
        void checkStatus()
      }, checkInterval)
    }

    runCheck()
    startPolling()

    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        runCheck()
        startPolling()
        return
      }

      clearPolling()
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)

    return () => {
      clearPolling()
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [checkStatus, checkInterval])

  // Re-check on window focus
  useEffect(() => {
    const handleFocus = () => {
      debugLog('[TenantStatusProvider] Window focused - rechecking tenant status')
      checkStatus()
    }

    window.addEventListener('focus', handleFocus)
    return () => {
      window.removeEventListener('focus', handleFocus)
    }
  }, [checkStatus])

  return <>{children}</>
}
