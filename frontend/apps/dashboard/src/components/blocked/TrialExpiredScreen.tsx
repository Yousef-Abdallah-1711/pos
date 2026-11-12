'use client'

import { Clock, Mail, LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { getDashboardLoginUrl } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'

interface TrialExpiredScreenProps {
  trialEndDate?: string | null
  supportEmail?: string
}

export function TrialExpiredScreen({
  trialEndDate,
  supportEmail = 'support@pos-saas.com',
}: TrialExpiredScreenProps) {
  const logout = useAuthStore((state) => state.logout)

  const daysSinceExpiration = trialEndDate
    ? Math.floor((new Date().getTime() - new Date(trialEndDate).getTime()) / (1000 * 60 * 60 * 24))
    : 0

  const handleLogout = () => {
    logout()
    window.location.href = getDashboardLoginUrl()
  }

  return (
    <div className="fixed inset-0 bg-gradient-to-b from-orange-50 to-white dark:from-gray-900 dark:to-gray-800 z-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full text-center">
        {/* Icon */}
        <div className="flex justify-center mb-6">
          <div className="bg-orange-100 dark:bg-orange-900/30 rounded-full p-4">
            <Clock className="h-12 w-12 text-orange-600" />
          </div>
        </div>

        {/* Title */}
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Trial Expired</h1>

        {/* Expiration Info */}
        <p className="text-gray-600 dark:text-gray-300 mb-2">
          Your 15-day free trial ended {daysSinceExpiration} {daysSinceExpiration === 1 ? 'day' : 'days'} ago.
        </p>

        {/* Description */}
        <p className="text-gray-600 dark:text-gray-300 mb-8">
          To continue using all features and manage your restaurant, please contact our support team to set up a paid subscription.
        </p>

        {/* Benefits List */}
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 mb-8 text-left">
          <h3 className="font-semibold text-gray-900 dark:text-white mb-4">What you're missing:</h3>
          <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-300">
            <li>✓ Full access to all dashboard features</li>
            <li>✓ Order management and tracking</li>
            <li>✓ Inventory and menu management</li>
            <li>✓ Employee and HR tools</li>
            <li>✓ Advanced analytics and reporting</li>
          </ul>
        </div>

        {/* Contact Support Button */}
        <Button
          onClick={() => {
            window.location.href = `mailto:${supportEmail}`
          }}
          className="w-full mb-3"
          size="lg"
        >
          <Mail className="h-5 w-5 mr-2" />
          Contact Support
        </Button>

        {/* Support Email */}
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">{supportEmail}</p>

        {/* Logout Button */}
        <Button onClick={handleLogout} variant="outline" className="w-full">
          <LogOut className="h-4 w-4 mr-2" />
          Sign Out
        </Button>
      </div>
    </div>
  )
}
