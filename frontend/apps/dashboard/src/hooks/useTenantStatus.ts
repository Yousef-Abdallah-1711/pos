'use client'

import { useMemo } from 'react'
import { useTenantStatusStore } from '@/stores/tenantStatusStore'
import { SubscriptionStatus } from '@/types/tenant-status'

export interface UseTenantStatusReturn {
  status: SubscriptionStatus | null
  trialEndDate: string | null
  daysRemaining: number | null
  canAccess: boolean
  blockReason: string | null
  isBlocked: boolean
  isLoading: boolean
  error: string | null
  checkStatus: () => Promise<void>
  getBlockMessage: () => string
}

export function useTenantStatus(): UseTenantStatusReturn {
  const store = useTenantStatusStore()

  const isBlocked = useMemo(() => !store.canAccess, [store.canAccess])

  const getBlockMessage = useMemo(() => {
    return () => {
      switch (store.blockReason) {
        case 'trial_expired':
          return 'Your trial period has expired. Please contact support to subscribe and continue using the platform.'
        case 'expired':
          return 'Your subscription has expired. Update your billing plan to restore full dashboard access.'
        case 'tenant_disabled':
          return 'Your account has been temporarily disabled. Please contact support for assistance.'
        case 'tenant_not_found':
          return 'Unable to verify your account. Please contact support.'
        default:
          return store.error || 'Your account access is currently restricted. Please contact support.'
      }
    }
  }, [store.blockReason, store.error])

  return {
    status: store.status,
    trialEndDate: store.trialEndDate,
    daysRemaining: store.daysRemaining,
    canAccess: store.canAccess,
    blockReason: store.blockReason,
    isBlocked,
    isLoading: store.isLoading,
    error: store.error,
    checkStatus: store.checkStatus,
    getBlockMessage,
  }
}
