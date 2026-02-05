'use client'

import { useState, useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { useHasPermission } from '@/hooks/useRbac'
import { useNotifications, Notification, NotificationModule, NotificationPriority } from '@/hooks/useNotifications'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { Button } from '@/components/ui/button'
import { NotificationActions } from '@/components/notifications/NotificationActions'
import { Bell, Trash2, CheckCheck, Clock, AlertTriangle, Info, CheckCircle, AlertCircle, Settings, Lock, Loader2 } from 'lucide-react'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'

const notificationModuleIcons: Record<NotificationModule, any> = {
  products: AlertTriangle,
  orders: Clock,
  hr: Info,
  inventory: AlertCircle,
  system: Bell,
}

const priorityColors: Record<NotificationPriority, string> = {
  low: 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-700',
  normal: 'bg-gray-50 dark:bg-gray-800 border-gray-200 dark:border-gray-700',
  high: 'bg-orange-50 dark:bg-orange-900/20 border-orange-200 dark:border-orange-700',
  critical: 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-700',
}

const priorityTextColors: Record<NotificationPriority, string> = {
  low: 'text-blue-700 dark:text-blue-300',
  normal: 'text-gray-700 dark:text-gray-300',
  high: 'text-orange-700 dark:text-orange-300',
  critical: 'text-red-700 dark:text-red-300',
}

export default function NotificationsPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'

  // Check RBAC permissions for Notifications module
  const { hasPermission: canViewNotifications, loading: canViewLoading } = useHasPermission(
    MODULE_IDS.NOTIFICATIONS,
    PERMISSION_LEVELS.READ,
  )
  const { hasPermission: canDeleteNotifications, loading: canDeleteLoading } = useHasPermission(
    MODULE_IDS.NOTIFICATIONS,
    PERMISSION_LEVELS.DELETE,
  )
  const permissionLoading = canViewLoading || canDeleteLoading

  const {
    notifications,
    unreadCount,
    totalCount,
    loading,
    stats,
    fetchNotifications,
    markAsRead,
    markAsUnread,
    markAllAsRead,
    deleteNotification,
  } = useNotifications()

  // Show permission denied if user can't view notifications
  if (permissionLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  if (!canViewNotifications) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <Bell className="h-8 w-8" />
            {t('notifications.title')}
          </h1>
        </div>

        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                Access Denied
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                You do not have permission to view notifications. Please contact your administrator.
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const [filterModule, setFilterModule] = useState<string>('')
  const [filterPriority, setFilterPriority] = useState<string>('')
  const [showUnreadOnly, setShowUnreadOnly] = useState(false)
  const [sortBy, setSortBy] = useState<'created_at' | 'priority'>('created_at')
  const [page, setPage] = useState(1)

  // Apply filters and sorting
  useEffect(() => {
    fetchNotifications({
      module: filterModule || undefined,
      priority: filterPriority || undefined,
      is_read: showUnreadOnly ? false : undefined,
      sort: sortBy,
      page,
      limit: 20,
    })
  }, [filterModule, filterPriority, showUnreadOnly, sortBy, page, fetchNotifications])

  const handleNotificationClick = (notification: Notification) => {
    if (!notification.is_read) {
      markAsRead(notification.id)
    }
    if (notification.action_url) {
      window.location.href = notification.action_url
    }
  }

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
      return locale === 'ar' ? 'أمس' : 'Yesterday'
    } else {
      return date.toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US', {
        month: 'short',
        day: 'numeric',
        year: date.getFullYear() !== today.getFullYear() ? 'numeric' : undefined,
      })
    }
  }

  const sortedNotifications = [...notifications].sort((a, b) => {
    if (sortBy === 'priority') {
      const priorityOrder = { critical: 0, high: 1, normal: 2, low: 3 }
      return (priorityOrder[a.priority as NotificationPriority] || 99) -
        (priorityOrder[b.priority as NotificationPriority] || 99)
    }
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })

  return (
    <div className={`space-y-6 ${isRTL ? 'rtl' : 'ltr'}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
            <Bell className="h-8 w-8" />
            {t('notifications.title')}
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-2">
            {t('notifications.unreadCount').replace('{unread}', String(unreadCount)).replace('{total}', String(totalCount))}
          </p>
        </div>

        <div className="flex gap-2">
          {unreadCount > 0 && (
            <Button
              onClick={() => markAllAsRead()}
              variant="outline"
              className="gap-2"
            >
              <CheckCheck className="h-4 w-4" />
              {t('notifications.markAllAsRead')}
            </Button>
          )}
        </div>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <p className="text-sm text-gray-600 dark:text-gray-400">{t('notifications.total')}</p>
            <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">{stats.total}</p>
          </div>
          <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 border border-blue-200 dark:border-blue-700">
            <p className="text-sm text-blue-700 dark:text-blue-300">{t('notifications.unread')}</p>
            <p className="text-2xl font-bold text-blue-900 dark:text-blue-100 mt-1">{stats.unread}</p>
          </div>
          <div className="bg-green-50 dark:bg-green-900/20 rounded-lg p-4 border border-green-200 dark:border-green-700">
            <p className="text-sm text-green-700 dark:text-green-300">{t('notifications.read')}</p>
            <p className="text-2xl font-bold text-green-900 dark:text-green-100 mt-1">{stats.read}</p>
          </div>
          <div className="bg-red-50 dark:bg-red-900/20 rounded-lg p-4 border border-red-200 dark:border-red-700">
            <p className="text-sm text-red-700 dark:text-red-300">{t('notifications.critical')}</p>
            <p className="text-2xl font-bold text-red-900 dark:text-red-100 mt-1">
              {stats.by_priority?.critical || 0}
            </p>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-lg p-6 border border-gray-200 dark:border-gray-700 space-y-4">
        <h2 className="font-semibold text-gray-900 dark:text-white">{t('notifications.filters')}</h2>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {/* Module Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('notifications.module')}
            </label>
            <select
              value={filterModule}
              onChange={(e) => {
                setFilterModule(e.target.value)
                setPage(1)
              }}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="">{t('notifications.all')}</option>
              <option value="products">{t('notifications.products')}</option>
              <option value="orders">{t('notifications.orders')}</option>
              <option value="hr">{t('notifications.hr')}</option>
              <option value="inventory">{t('notifications.inventory')}</option>
              <option value="system">{t('notifications.system')}</option>
            </select>
          </div>

          {/* Priority Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('notifications.priority')}
            </label>
            <select
              value={filterPriority}
              onChange={(e) => {
                setFilterPriority(e.target.value)
                setPage(1)
              }}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="">{t('notifications.all')}</option>
              <option value="critical">{t('notifications.critical')}</option>
              <option value="high">{t('notifications.high')}</option>
              <option value="normal">{t('notifications.normal')}</option>
              <option value="low">{t('notifications.low')}</option>
            </select>
          </div>

          {/* Read Status Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('notifications.status')}
            </label>
            <select
              value={showUnreadOnly ? 'unread' : ''}
              onChange={(e) => {
                setShowUnreadOnly(e.target.value === 'unread')
                setPage(1)
              }}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="">{t('notifications.all')}</option>
              <option value="unread">{t('notifications.unreadOnly')}</option>
            </select>
          </div>

          {/* Sort By */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('notifications.sortBy')}
            </label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as 'created_at' | 'priority')}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="created_at">{t('notifications.newest')}</option>
              <option value="priority">{t('notifications.priority')}</option>
            </select>
          </div>
        </div>
      </div>

      {/* Notifications List */}
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500 mx-auto mb-4"></div>
            <p className="text-gray-500 dark:text-gray-400">{t('notifications.loading')}</p>
          </div>
        </div>
      ) : sortedNotifications.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-lg p-12 border border-gray-200 dark:border-gray-700 text-center">
          <Bell className="h-12 w-12 mx-auto mb-4 text-gray-400" />
          <p className="text-gray-500 dark:text-gray-400">
            {t('notifications.noNotifications')}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {sortedNotifications.map((notification) => {
            const ModuleIcon = notificationModuleIcons[notification.module] || Bell
            return (
              <div
                key={notification.id}
                className={`${priorityColors[notification.priority]} border rounded-lg p-4 transition-all hover:shadow-md cursor-pointer ${
                  !notification.is_read ? 'ring-2 ring-blue-500' : ''
                }`}
                onClick={() => handleNotificationClick(notification)}
              >
                <div className="flex items-start gap-4">
                  <ModuleIcon className={`h-6 w-6 mt-1 flex-shrink-0 ${priorityTextColors[notification.priority]}`} />

                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <h3 className={`font-semibold ${priorityTextColors[notification.priority]}`}>
                          {notification.title}
                        </h3>
                        <p className="text-sm text-gray-700 dark:text-gray-300 mt-1">{notification.message}</p>
                        {notification.description && (
                          <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">{notification.description}</p>
                        )}
                      </div>
                      {!notification.is_read && (
                        <div className="h-3 w-3 rounded-full bg-blue-500 flex-shrink-0 mt-1"></div>
                      )}
                    </div>

                    <div className="flex items-center gap-4 mt-2 flex-wrap">
                      <span className="text-xs text-gray-600 dark:text-gray-400 flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {formatDate(notification.created_at)}
                      </span>
                      <span className={`text-xs px-2 py-1 rounded-full ${priorityColors[notification.priority]}`}>
                        {notification.priority.charAt(0).toUpperCase() + notification.priority.slice(1)}
                      </span>
                      <span className="text-xs px-2 py-1 bg-gray-200 dark:bg-gray-700 rounded-full text-gray-700 dark:text-gray-300 capitalize">
                        {notification.module}
                      </span>
                    </div>
                  </div>

                  <div className="flex-shrink-0">
                    <NotificationActions
                      notificationId={notification.id}
                      isRead={notification.is_read}
                      onMarkAsRead={markAsRead}
                      onMarkAsUnread={markAsUnread}
                      onDelete={deleteNotification}
                      canDelete={canDeleteNotifications}
                      translations={{
                        markAsRead: t('notifications.markAsRead'),
                        markAsUnread: t('notifications.markAsUnread'),
                        delete: t('notifications.delete'),
                      }}
                      isRTL={isRTL}
                    />
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      )}

      {/* Pagination */}
      {Math.ceil(totalCount / 20) > 1 && (
        <div className="flex items-center justify-center gap-2 mt-6">
          <Button
            variant="outline"
            disabled={page === 1}
            onClick={() => setPage(Math.max(1, page - 1))}
          >
            {t('notifications.previous')}
          </Button>
          <span className="text-sm text-gray-600 dark:text-gray-400">
            {t('notifications.page')} {page} {t('notifications.of')} {Math.ceil(totalCount / 20)}
          </span>
          <Button
            variant="outline"
            disabled={page >= Math.ceil(totalCount / 20)}
            onClick={() => setPage(Math.min(Math.ceil(totalCount / 20), page + 1))}
          >
            {t('notifications.next')}
          </Button>
        </div>
      )}
    </div>
  )
}
