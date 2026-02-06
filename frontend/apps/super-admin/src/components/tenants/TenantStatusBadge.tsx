'use client'

import type { SubscriptionStatus } from '@/types'

interface TenantStatusBadgeProps {
  status: SubscriptionStatus
  isDisabled?: boolean
  className?: string
}

export function TenantStatusBadge({
  status,
  isDisabled = false,
  className = '',
}: TenantStatusBadgeProps) {
  // Override status if disabled
  const displayStatus = isDisabled ? 'disabled' : status

  const getStatusStyles = (
    status: SubscriptionStatus | 'disabled'
  ): { bg: string; text: string; label: string } => {
    switch (status) {
      case 'trial':
        return {
          bg: 'bg-blue-100',
          text: 'text-blue-800',
          label: 'Trial',
        }
      case 'active':
        return {
          bg: 'bg-green-100',
          text: 'text-green-800',
          label: 'Active',
        }
      case 'expired':
        return {
          bg: 'bg-red-100',
          text: 'text-red-800',
          label: 'Expired',
        }
      case 'suspended':
        return {
          bg: 'bg-orange-100',
          text: 'text-orange-800',
          label: 'Suspended',
        }
      case 'disabled':
        return {
          bg: 'bg-gray-100',
          text: 'text-gray-800',
          label: 'Disabled',
        }
      default:
        return {
          bg: 'bg-gray-100',
          text: 'text-gray-800',
          label: 'Unknown',
        }
    }
  }

  const styles = getStatusStyles(displayStatus)

  return (
    <span
      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${styles.bg} ${styles.text} ${className}`}
    >
      {styles.label}
    </span>
  )
}
