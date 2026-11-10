'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Clock, CreditCard, Lock, Mail } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { getDashboardLoginUrl } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'
import { getLocaleFromPath } from '@/lib/translations'

interface BlockedOverlayProps {
  blockReason: 'trial_expired' | 'tenant_disabled' | string
  supportEmail?: string
}

export function BlockedOverlay({ blockReason, supportEmail = 'support@pos-saas.com' }: BlockedOverlayProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const logout = useAuthStore((state) => state.logout)
  const billingHref = `/${locale}/dashboard/settings/billing`

  const getContent = () => {
    switch (blockReason) {
      case 'trial_expired':
      case 'expired':
        return {
          icon: <Clock className="h-16 w-16 text-orange-500" />,
          title: blockReason === 'expired' ? 'Your Subscription Has Expired' : 'Your Trial Has Expired',
          description:
            blockReason === 'expired'
              ? 'Your subscription has expired. Renew now to restore full access to your dashboard.'
              : 'Your 15-day free trial has ended. Subscribe now to keep full access to your dashboard.',
          ctaNote: 'Continue with all features for just 500 EGP/month.',
        }
      case 'tenant_disabled':
        return {
          icon: <Lock className="h-16 w-16 text-red-500" />,
          title: 'Account Access Disabled',
          description:
            'Your account has been temporarily disabled. Please contact support for assistance.',
          ctaNote: '',
        }
      default:
        return {
          icon: <Lock className="h-16 w-16 text-red-500" />,
          title: 'Access Restricted',
          description: 'Your account access is currently restricted. Please contact support.',
          ctaNote: '',
        }
    }
  }

  const content = getContent()

  const handleLogout = () => {
    logout()
    window.location.href = getDashboardLoginUrl(pathname)
  }

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-4rem)] bg-gray-50 dark:bg-gray-900">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-2xl max-w-md w-full mx-4">
        <div className="p-8">
          {/* Icon */}
          <div className="flex justify-center mb-6">{content.icon}</div>

          {/* Title */}
          <h2 className="text-2xl font-bold text-center text-gray-900 dark:text-white mb-4">
            {content.title}
          </h2>

          {/* Description */}
          <p className="text-center text-gray-600 dark:text-gray-300 mb-8">{content.description}</p>

          {blockReason === 'trial_expired' || blockReason === 'expired' ? (
            <div className="space-y-3 mb-4">
              <Button asChild className="w-full bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600">
                <Link href={billingHref}>
                  <CreditCard className="h-4 w-4 mr-2" />
                  Subscribe Now - 500 EGP/month
                </Link>
              </Button>
              <p className="text-xs text-center text-gray-500 dark:text-gray-400">{content.ctaNote}</p>
              <button
                onClick={() => {
                  window.location.href = `mailto:${supportEmail}`
                }}
                className="w-full text-center text-sm text-blue-600 dark:text-blue-400 hover:underline"
              >
                Contact Support
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2 mb-4">
              <Button
                onClick={() => {
                  window.location.href = `mailto:${supportEmail}`
                }}
                className="w-full"
              >
                <Mail className="h-4 w-4 mr-2" />
                Contact Support
              </Button>
            </div>
          )}

          {/* Support Email Display */}
          <p className="text-center text-sm text-gray-500 dark:text-gray-400 mb-6">
            {supportEmail}
          </p>

          {/* Divider */}
          <div className="border-t border-gray-200 dark:border-gray-700 mb-4"></div>

          {/* Logout Link */}
          <button
            onClick={handleLogout}
            className="w-full text-center text-sm text-blue-600 dark:text-blue-400 hover:underline font-medium"
          >
            Sign Out
          </button>
        </div>
      </div>
    </div>
  )
}
