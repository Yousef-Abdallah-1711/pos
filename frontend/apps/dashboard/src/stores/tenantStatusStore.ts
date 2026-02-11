'use client'

import { create } from 'zustand'
import { TenantStatusState, TenantStatusResponse, SubscriptionStatus } from '@/types/tenant-status'
import { api } from '@/lib/api'
import { debugLog, debugError, debugWarn } from '@/lib/debug'

interface TenantStatusStoreState extends TenantStatusState {
  setStatus: (status: Partial<TenantStatusState>) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  checkStatus: () => Promise<void>
  clear: () => void
}

export const useTenantStatusStore = create<TenantStatusStoreState>((set) => ({
  // Initial state
  status: null,
  trialEndDate: null,
  daysRemaining: null,
  canAccess: true,
  blockReason: null,
  isLoading: false,
  error: null,
  lastChecked: null,

  // Actions
  setStatus: (status: Partial<TenantStatusState>) => {
    set((state) => ({
      ...state,
      ...status,
    }))
  },

  setLoading: (loading: boolean) => {
    set({ isLoading: loading })
  },

  setError: (error: string | null) => {
    set({ error })
  },

  checkStatus: async () => {
    set({ isLoading: true, error: null })
    try {
      debugLog('[TenantStatusStore] Checking tenant status...')
      const response = await api.get<TenantStatusResponse>('/tenant/status')

      if (response.data.success) {
        const data = response.data.data
        debugLog('[TenantStatusStore] Status check successful:', {
          canAccess: data.can_access,
          blockReason: data.block_reason,
        })

        set({
          canAccess: data.can_access,
          blockReason: data.block_reason || null,
          status: (data.status as SubscriptionStatus) || null,
          trialEndDate: data.trial_end_date || null,
          daysRemaining: data.days_remaining || null,
          isLoading: false,
          lastChecked: new Date(),
        })
      } else {
        debugWarn('[TenantStatusStore] Status check returned success: false')
        set({
          isLoading: false,
          error: 'Failed to check tenant status',
        })
      }
    } catch (error: any) {
      debugError('[TenantStatusStore] Error checking tenant status:', error)

      // Handle 403 Forbidden with TENANT_ACCESS_BLOCKED
      if (error.response?.status === 403) {
        const errorCode = error.response?.data?.error?.code
        const reason = error.response?.data?.error?.reason || 'unknown'
        const message = error.response?.data?.error?.message || 'Access denied'

        debugLog('[TenantStatusStore] Tenant is blocked:', { errorCode, reason, message })

        set({
          canAccess: false,
          blockReason: reason,
          isLoading: false,
          error: message,
          lastChecked: new Date(),
        })
      } else {
        set({
          isLoading: false,
          error: error.message || 'Failed to check tenant status',
        })
      }
    }
  },

  clear: () => {
    set({
      status: null,
      trialEndDate: null,
      daysRemaining: null,
      canAccess: true,
      blockReason: null,
      isLoading: false,
      error: null,
      lastChecked: null,
    })
  },
}))
