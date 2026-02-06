'use client'

import { ChefHat, Lock, LogOut } from 'lucide-react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'
import { usePathname } from 'next/navigation'
import { getLocaleFromPath } from '@/lib/translations'

export function BlockedSidebar() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const logout = useAuthStore((state) => state.logout)
  const isRTL = locale === 'ar'

  const handleLogout = () => {
    logout()
    const landingUrl = getConfig().landingUrl
    window.location.href = `${landingUrl}/auth/login`
  }

  return (
    <aside
      className={`fixed inset-y-0 ${isRTL ? 'right-0' : 'left-0'} z-40 w-64 bg-white dark:bg-gray-800 ${
        isRTL ? 'border-l' : 'border-r'
      } border-gray-200 dark:border-gray-700`}
    >
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center gap-2 opacity-50">
          <ChefHat className="h-8 w-8 text-gray-400 dark:text-gray-500" />
          <span className="text-xl font-bold text-gray-400 dark:text-gray-500">POS Dashboard</span>
        </div>
      </div>

      {/* Blocked Notice */}
      <div className="p-4">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex items-start gap-3">
          <Lock className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-sm font-semibold text-red-900 dark:text-red-300">Access Restricted</p>
            <p className="text-xs text-red-800 dark:text-red-400 mt-1">
              Your account access is temporarily restricted. Please contact support.
            </p>
          </div>
        </div>
      </div>

      {/* Disabled Menu Items */}
      <nav className="p-4">
        <ul className="space-y-2">
          <li>
            <div className="flex items-center gap-3 px-4 py-3 rounded-lg opacity-50 cursor-not-allowed text-gray-400 dark:text-gray-600">
              <div className="h-5 w-5 rounded bg-gray-200 dark:bg-gray-700"></div>
              <span className="font-medium text-sm">Dashboard</span>
            </div>
          </li>
          <li>
            <div className="flex items-center gap-3 px-4 py-3 rounded-lg opacity-50 cursor-not-allowed text-gray-400 dark:text-gray-600">
              <div className="h-5 w-5 rounded bg-gray-200 dark:bg-gray-700"></div>
              <span className="font-medium text-sm">Orders</span>
            </div>
          </li>
          <li>
            <div className="flex items-center gap-3 px-4 py-3 rounded-lg opacity-50 cursor-not-allowed text-gray-400 dark:text-gray-600">
              <div className="h-5 w-5 rounded bg-gray-200 dark:bg-gray-700"></div>
              <span className="font-medium text-sm">Products</span>
            </div>
          </li>
          <li>
            <div className="flex items-center gap-3 px-4 py-3 rounded-lg opacity-50 cursor-not-allowed text-gray-400 dark:text-gray-600">
              <div className="h-5 w-5 rounded bg-gray-200 dark:bg-gray-700"></div>
              <span className="font-medium text-sm">Staff</span>
            </div>
          </li>
          <li>
            <div className="flex items-center gap-3 px-4 py-3 rounded-lg opacity-50 cursor-not-allowed text-gray-400 dark:text-gray-600">
              <div className="h-5 w-5 rounded bg-gray-200 dark:bg-gray-700"></div>
              <span className="font-medium text-sm">Settings</span>
            </div>
          </li>
        </ul>
      </nav>

      {/* Logout Button - Always Active */}
      <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
        <Button onClick={handleLogout} variant="outline" size="sm" className="w-full text-xs">
          <LogOut className="h-4 w-4 mr-2" />
          Sign Out
        </Button>
      </div>
    </aside>
  )
}
