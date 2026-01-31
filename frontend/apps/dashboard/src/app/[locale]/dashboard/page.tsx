'use client'

import { useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'
import { Store, Users, ShoppingCart, TrendingUp } from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

export default function DashboardPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)

  // Check if authenticated by checking if user and token exist
  const isAuthenticated = !!(user && token)

  useEffect(() => {
    // Skip redirect if in debug mode
    if (typeof window !== 'undefined') {
      const debugMode = new URLSearchParams(window.location.search).get('debug') === 'true'
      if (debugMode) {
        return
      }
    }

    if (!isAuthenticated) {
      window.location.href = `${getConfig().landingUrl}/auth/login`
    }
  }, [isAuthenticated, user])

  // Allow rendering if in debug mode, even if not authenticated
  const debugMode = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('debug') === 'true'

  if (!isAuthenticated && !debugMode) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  const stats = [
    {
      name: t('dashboard.statistics'),
      value: '$0.00',
      icon: TrendingUp,
      color: 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400',
      link: '/dashboard',
    },
    {
      name: t('dashboard.orders'),
      value: 'View All',
      icon: ShoppingCart,
      color: 'bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400',
      link: '/dashboard/orders',
    },
    {
      name: t('dashboard.restaurants'),
      value: t('dashboard.restaurantsCount'),
      icon: Store,
      color: 'bg-purple-100 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400',
      link: '/dashboard',
    },
    {
      name: t('dashboard.staff'),
      value: t('dashboard.staffCount'),
      icon: Users,
      color: 'bg-orange-100 text-orange-600 dark:bg-orange-900/30 dark:text-orange-400',
      link: '/dashboard',
    },
  ]

  return (
    <div>
      {/* Welcome Section */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          {t('dashboard.welcome')}, {user?.name}! 👋
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          {t('dashboard.overview')}
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => (
          <div
            key={stat.name}
            className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition-shadow"
          >
            <div className="flex items-center justify-between mb-4">
              <div className={`p-3 rounded-lg ${stat.color}`}>
                <stat.icon className="h-6 w-6" />
              </div>
            </div>
            <p className="text-2xl font-bold text-gray-900 dark:text-white mb-1">
              {stat.value}
            </p>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {stat.name}
            </p>
          </div>
        ))}
      </div>

      {/* Getting Started */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-8 border border-gray-200 dark:border-gray-700">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
          {t('dashboard.gettingStarted')}
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-6">
          {t('dashboard.gettingStartedDesc')}
        </p>

        <div className="space-y-4">
          <div className="flex items-start gap-3 p-4 rounded-lg bg-gray-50 dark:bg-gray-900/50">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 flex items-center justify-center font-semibold">
              1
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                {t('dashboard.step1Title')}
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {t('dashboard.step1Desc')}
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg bg-gray-50 dark:bg-gray-900/50">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 flex items-center justify-center font-semibold">
              2
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                {t('dashboard.step2Title')}
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {t('dashboard.step2Desc')}
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg bg-gray-50 dark:bg-gray-900/50">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 flex items-center justify-center font-semibold">
              3
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                {t('dashboard.step3Title')}
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {t('dashboard.step3Desc')}
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3 p-4 rounded-lg bg-gray-50 dark:bg-gray-900/50">
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 flex items-center justify-center font-semibold">
              4
            </div>
            <div>
              <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                {t('dashboard.step4Title')}
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {t('dashboard.step4Desc')}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
