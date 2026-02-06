'use client'

import { Clock, Lock, Mail } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'

interface BlockedOverlayProps {
  blockReason: 'trial_expired' | 'tenant_disabled' | string
  supportEmail?: string
}

export function BlockedOverlay({ blockReason, supportEmail = 'support@pos-saas.com' }: BlockedOverlayProps) {
  const logout = useAuthStore((state) => state.logout)

  const getContent = () => {
    switch (blockReason) {
      case 'trial_expired':
        return {
          icon: <Clock className="h-16 w-16 text-orange-500" />,
          title: 'Your Trial Has Expired',
          description:
            'Your 15-day free trial has ended. Please contact support to subscribe and continue using the platform.',
          buttonText: 'Contact Support',
        }
      case 'tenant_disabled':
        return {
          icon: <Lock className="h-16 w-16 text-red-500" />,
          title: 'Account Access Disabled',
          description:
            'Your account has been temporarily disabled. Please contact support for assistance.',
          buttonText: 'Contact Support',
        }
      default:
        return {
          icon: <Lock className="h-16 w-16 text-red-500" />,
          title: 'Access Restricted',
          description: 'Your account access is currently restricted. Please contact support.',
          buttonText: 'Contact Support',
        }
    }
  }

  const content = getContent()

  const handleLogout = () => {
    logout()
    const landingUrl = getConfig().landingUrl
    window.location.href = `${landingUrl}/auth/login`
  }

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center">
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

          {/* Contact Support Button */}
          <div className="flex items-center gap-2 mb-4">
            <Button
              onClick={() => {
                window.location.href = `mailto:${supportEmail}`
              }}
              className="w-full"
            >
              <Mail className="h-4 w-4 mr-2" />
              {content.buttonText}
            </Button>
          </div>

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
