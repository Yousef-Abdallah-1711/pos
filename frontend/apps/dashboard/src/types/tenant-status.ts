// Subscription status types
export type SubscriptionStatus = 'trial' | 'active' | 'expired' | 'suspended' | 'disabled'

// Tenant status response from API
export interface TenantStatusResponse {
  success: boolean
  data: {
    can_access: boolean
    block_reason: string | null
    is_maintenance?: boolean
    status?: SubscriptionStatus
    trial_end_date?: string | null
    days_remaining?: number | null
    is_disabled?: boolean
    subscription_expires_at?: string | null
    last_payment_at?: string | null
    plan?: {
      id: string
      name?: string
      price: number
      currency: string
    }
  }
}

// Tenant status store state
export interface TenantStatusState {
  status: SubscriptionStatus | null
  trialEndDate: string | null
  daysRemaining: number | null
  canAccess: boolean
  blockReason: string | null
  isLoading: boolean
  error: string | null
  lastChecked: Date | null
}
