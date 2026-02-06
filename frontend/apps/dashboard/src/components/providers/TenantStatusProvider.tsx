'use client'

import { ReactNode, useEffect } from 'react'
import { useTenantStatusStore } from '@/stores/tenantStatusStore'

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
    console.log('[TenantStatusProvider] Mounting - checking tenant status')

    // Initial check
    checkStatus()

    // Set up periodic checks
    const intervalId = setInterval(() => {
      console.log('[TenantStatusProvider] Periodic check - checking tenant status')
      checkStatus()
    }, checkInterval)

    return () => {
      clearInterval(intervalId)
    }
  }, [checkStatus, checkInterval])

  // Re-check on window focus
  useEffect(() => {
    const handleFocus = () => {
      console.log('[TenantStatusProvider] Window focused - rechecking tenant status')
      checkStatus()
    }

    window.addEventListener('focus', handleFocus)
    return () => {
      window.removeEventListener('focus', handleFocus)
    }
  }, [checkStatus])

  return <>{children}</>
}
