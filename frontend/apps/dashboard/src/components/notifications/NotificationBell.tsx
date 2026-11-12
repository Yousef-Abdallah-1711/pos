'use client'

import { useState, useRef, useEffect } from 'react'
import Link from 'next/link'
import { Bell, AlertTriangle } from 'lucide-react'
import { useNotifications } from '@/hooks/useNotifications'
import { usePathname } from 'next/navigation'
import { getLocaleFromPath, createTranslator } from '@/lib/translations'

// Legacy shape: this bell predates the generic notifications model and renders low-stock products.
// It is not imported anywhere; kept compiling until the owner decides whether to delete it.
interface LowStockProduct {
  id: number
  name_en: string
  name_ar?: string
  quantity_in_stock: number
  low_stock_threshold: number
  created_at: string
}

export function NotificationBell() {
  const { notifications: rawNotifications, unreadCount: count, fetchNotifications } = useNotifications()
  const notifications = rawNotifications as unknown as LowStockProduct[]
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'

  // Fetch notifications on mount and set up auto-refresh
  useEffect(() => {
    fetchNotifications()

    // Auto-refresh every 5 minutes (300000ms)
    const refreshInterval = setInterval(() => {
      fetchNotifications({}, true)
    }, 300000)

    // Refresh when page becomes visible (e.g., user returns from another tab)
    const handleVisibilityChange = () => {
      if (!document.hidden) {
        fetchNotifications({}, true)
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)

    return () => {
      clearInterval(refreshInterval)
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [fetchNotifications])

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const today = new Date()
    const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000)

    if (date.toDateString() === today.toDateString()) {
      return date.toLocaleTimeString(locale === 'ar' ? 'ar-EG' : 'en-US', {
        hour: '2-digit',
        minute: '2-digit',
      })
    } else if (date.toDateString() === yesterday.toDateString()) {
      return t('notifications.yesterday')
    } else {
      return date.toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US', {
        month: 'short',
        day: 'numeric',
      })
    }
  }

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="relative p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition"
        title="Notifications"
      >
        <Bell className="h-5 w-5 text-gray-700 dark:text-gray-300" />
        {count > 0 && (
          <span className="absolute top-0 right-0 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-white transform translate-x-1/2 -translate-y-1/2 bg-red-500 rounded-full">
            {count > 99 ? '99+' : count}
          </span>
        )}
      </button>

      {/* Notification Dropdown Panel */}
      {isOpen && (
        <div
          className={`absolute ${
            isRTL ? 'left-0' : 'right-0'
          } mt-2 w-80 max-h-96 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col z-50`}
        >
          {/* Header */}
          <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900">
            <h3 className="text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <AlertTriangle className="h-4 w-4 text-yellow-500" />
              {t('notifications.lowStockAlerts')}
            </h3>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              {count === 0
                ? t('notifications.noAlerts')
                : t('notifications.productsNeedRestocking').replace('{count}', String(count))}
            </p>
          </div>

          {/* Notifications List */}
          <div className="flex-1 overflow-y-auto">
            {notifications.length === 0 ? (
              <div className="flex items-center justify-center h-32 text-gray-500 dark:text-gray-400">
                <div className="text-center">
                  <Bell className="h-8 w-8 mx-auto mb-2 opacity-50" />
                  <p className="text-sm">
                    {t('notifications.noNotifications')}
                  </p>
                </div>
              </div>
            ) : (
              <div className="divide-y divide-gray-200 dark:divide-gray-700">
                {notifications.map((product: LowStockProduct) => (
                  <Link
                    key={product.id}
                    href={`/${locale}/dashboard/products/${product.id}/edit`}
                    onClick={() => setIsOpen(false)}
                    className="block px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {product.name_en}
                        </p>
                        {product.name_ar && (
                          <p className="text-xs text-gray-500 dark:text-gray-400" dir="rtl">
                            {product.name_ar}
                          </p>
                        )}
                        <div className="mt-1 flex items-center gap-2">
                          <span className="text-xs font-semibold text-red-600 dark:text-red-400">
                            {product.quantity_in_stock} {t('inventory.units')}
                          </span>
                          <span className="text-xs text-gray-500 dark:text-gray-400">
                            {t('notifications.threshold')}: {product.low_stock_threshold}
                          </span>
                        </div>
                      </div>
                      <span className="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap ms-2">
                        {formatDate(product.created_at)}
                      </span>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>

          {/* Footer */}
          {notifications.length > 0 && (
            <div className="border-t border-gray-200 dark:border-gray-700 px-4 py-3 bg-gray-50 dark:bg-gray-900">
              <Link
                href={`/${locale}/dashboard/products`}
                onClick={() => setIsOpen(false)}
                className="text-xs font-medium text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300"
              >
                {t('notifications.viewAllProducts')}
              </Link>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
