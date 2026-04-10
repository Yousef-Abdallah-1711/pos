'use client'

import { Lock, Mail, LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { getDashboardLoginUrl } from '@/lib/authSession'
import { useAuthStore } from '@/stores/authStore'

interface AccountDisabledScreenProps {
  disabledReason?: string | null
  supportEmail?: string
}

export function AccountDisabledScreen({
  disabledReason,
  supportEmail = 'support@pos-saas.com',
}: AccountDisabledScreenProps) {
  const logout = useAuthStore((state) => state.logout)

  const handleLogout = () => {
    logout()
    window.location.href = getDashboardLoginUrl()
  }

  return (
    <div className="fixed inset-0 bg-gradient-to-b from-red-50 to-white dark:from-gray-900 dark:to-gray-800 z-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full text-center">
        {/* Icon */}
        <div className="flex justify-center mb-6">
          <div className="bg-red-100 dark:bg-red-900/30 rounded-full p-4">
            <Lock className="h-12 w-12 text-red-600" />
          </div>
        </div>

        {/* Title */}
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">Account Disabled</h1>

        {/* Main Description */}
        <p className="text-gray-600 dark:text-gray-300 mb-4">
          Your account has been temporarily disabled and you cannot access the dashboard at this time.
        </p>

        {/* Reason (if provided) */}
        {disabledReason && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-6 text-left">
            <p className="text-sm font-semibold text-red-900 dark:text-red-300 mb-2">Reason:</p>
            <p className="text-sm text-red-800 dark:text-red-400">{disabledReason}</p>
          </div>
        )}

        {/* Explanation */}
        <p className="text-gray-600 dark:text-gray-300 mb-8">
          Please contact our support team for assistance in resolving this issue and regaining access to your account.
        </p>

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
