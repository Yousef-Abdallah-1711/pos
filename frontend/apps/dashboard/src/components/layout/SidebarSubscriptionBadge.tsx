'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Crown, Clock, AlertTriangle } from 'lucide-react'
import { useTenantStatus } from '@/hooks/useTenantStatus'
import { getLocaleFromPath } from '@/lib/translations'

export function SidebarSubscriptionBadge() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const billingHref = `/${locale}/dashboard/settings/billing`
  const { status, daysRemaining, isLoading } = useTenantStatus()

  if (isLoading || !status) return null

  const getStatusConfig = () => {
    if (status === 'active') {
      return {
        label: 'Pro Plan',
        icon: <Crown className="h-4 w-4" />,
        bgColor: 'bg-green-50 dark:bg-green-900/20',
        textColor: 'text-green-700 dark:text-green-400',
        borderColor: 'border-green-200 dark:border-green-800',
      }
    }

    if (status === 'trial') {
      const days = daysRemaining ?? 0
      if (days > 7) {
        return {
          label: 'Pro Trial',
          sublabel: `${days} days left`,
          icon: <Clock className="h-4 w-4" />,
          bgColor: 'bg-blue-50 dark:bg-blue-900/20',
          textColor: 'text-blue-700 dark:text-blue-400',
          borderColor: 'border-blue-200 dark:border-blue-800',
        }
      }
      if (days > 3) {
        return {
          label: 'Pro Trial',
          sublabel: `${days} days left`,
          icon: <Clock className="h-4 w-4" />,
          bgColor: 'bg-yellow-50 dark:bg-yellow-900/20',
          textColor: 'text-yellow-700 dark:text-yellow-400',
          borderColor: 'border-yellow-200 dark:border-yellow-800',
        }
      }
      if (days > 0) {
        const remainingText = `${days} day${days === 1 ? '' : 's'} left`
        return {
          label: 'Pro Trial',
          sublabel: days < 3 ? `${remainingText} • Subscribe now` : remainingText,
          icon: <AlertTriangle className="h-4 w-4" />,
          bgColor: 'bg-red-50 dark:bg-red-900/20',
          textColor: 'text-red-700 dark:text-red-400',
          borderColor: 'border-red-200 dark:border-red-800',
        }
      }
      return {
        label: 'Trial Expired',
        icon: <AlertTriangle className="h-4 w-4" />,
        sublabel: 'Tap to subscribe',
        bgColor: 'bg-orange-50 dark:bg-orange-900/20',
        textColor: 'text-orange-700 dark:text-orange-400',
        borderColor: 'border-orange-200 dark:border-orange-800',
      }
    }

    if (status === 'expired') {
      return {
        label: 'Subscription Expired',
        icon: <AlertTriangle className="h-4 w-4" />,
        sublabel: 'Tap to subscribe',
        bgColor: 'bg-orange-50 dark:bg-orange-900/20',
        textColor: 'text-orange-700 dark:text-orange-400',
        borderColor: 'border-orange-200 dark:border-orange-800',
      }
    }

    return {
      label: status.charAt(0).toUpperCase() + status.slice(1),
      icon: <Crown className="h-4 w-4" />,
      bgColor: 'bg-gray-50 dark:bg-gray-700',
      textColor: 'text-gray-600 dark:text-gray-400',
      borderColor: 'border-gray-200 dark:border-gray-600',
    }
  }

  const config = getStatusConfig()

  return (
    <Link href={billingHref} className="mx-4 mb-4 block">
      <div
        className={`p-3 rounded-lg border transition-all hover:shadow-sm hover:brightness-95 ${config.bgColor} ${config.borderColor}`}
      >
        <div className={`flex items-center gap-2 ${config.textColor}`}>
          {config.icon}
          <span className="text-sm font-medium">{config.label}</span>
        </div>
        {config.sublabel && (
          <p className={`text-xs mt-1 ${config.textColor} opacity-80`}>
            {config.sublabel}
          </p>
        )}
      </div>
    </Link>
  )
}
