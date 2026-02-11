import { useState, useEffect, useCallback } from 'react'
import { api } from '@/lib/api'
import { apiCache } from '@/lib/apiCache'
import { debugError } from '@/lib/debug'

export type NotificationType = 'low_stock' | 'order' | 'employee' | 'leave' | 'inventory' | 'attendance' | 'system'
export type NotificationModule = 'products' | 'orders' | 'hr' | 'inventory' | 'system'
export type NotificationPriority = 'low' | 'normal' | 'high' | 'critical'

export interface Notification {
  id: number
  type: NotificationType
  module: NotificationModule
  title: string
  message: string
  description?: string
  is_read: boolean
  read_at?: string
  related_entity_type?: string
  related_entity_id?: number
  priority: NotificationPriority
  action_url?: string
  action_label?: string
  icon_name?: string
  color?: string
  created_at: string
  updated_at: string
  expires_at?: string
}

export interface NotificationListResponse {
  total: number
  unread: number
  page: number
  page_size: number
  notifications: Notification[]
}

export interface NotificationStats {
  total: number
  unread: number
  read: number
  by_module: Record<string, number>
  by_priority: Record<string, number>
}

export interface NotificationFilters {
  module?: string
  type?: string
  is_read?: boolean
  priority?: string
  page?: number
  limit?: number
  sort?: string
  order?: string
}

export function useNotifications() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [totalCount, setTotalCount] = useState(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [stats, setStats] = useState<NotificationStats | null>(null)

  const fetchNotifications = useCallback(
    async (filters: NotificationFilters = {}, forceRefresh = false) => {
      try {
        setLoading(true)
        setError(null)

        // Check cache first if not forcing refresh
        const cacheKey = `notifications_${JSON.stringify(filters)}`
        if (!forceRefresh) {
          const cached = apiCache.get(cacheKey)
          if (cached) {
            setNotifications(cached.notifications)
            setUnreadCount(cached.unread)
            setTotalCount(cached.total)
            setLoading(false)
            return
          }
        }

        // Build query params
        const params = new URLSearchParams()
        if (filters.module) params.append('module', filters.module)
        if (filters.type) params.append('type', filters.type)
        if (filters.is_read !== undefined) params.append('is_read', String(filters.is_read))
        if (filters.priority) params.append('priority', filters.priority)
        if (filters.page) params.append('page', String(filters.page))
        if (filters.limit) params.append('limit', String(filters.limit))
        if (filters.sort) params.append('sort', filters.sort)
        if (filters.order) params.append('order', filters.order)

        // Fetch from API
        const response = await api.get<NotificationListResponse>(
          `/notifications?${params.toString()}`
        )
        const data = response.data

        // Cache for 5 minutes (300000ms)
        apiCache.set(cacheKey, data, 300000)

        setNotifications(data.notifications)
        setUnreadCount(data.unread)
        setTotalCount(data.total)
      } catch (err) {
        debugError('Failed to fetch notifications:', err)
        setError('Failed to load notifications')
        setNotifications([])
        setUnreadCount(0)
      } finally {
        setLoading(false)
      }
    },
    []
  )

  const fetchStats = useCallback(async () => {
    try {
      const response = await api.get<NotificationStats>('/notifications/stats')
      setStats(response.data)
    } catch (err) {
      debugError('Failed to fetch notification stats:', err)
    }
  }, [])

  const fetchUnreadCount = useCallback(async () => {
    try {
      const response = await api.get<{ unread: number }>('/notifications/unread-count')
      setUnreadCount(response.data.unread)
    } catch (err) {
      debugError('Failed to fetch unread count:', err)
    }
  }, [])

  const markAsRead = useCallback(
    async (notificationId: number) => {
      try {
        await api.post(`/notifications/${notificationId}/read`, {})

        // Update local state
        setNotifications((prev) =>
          prev.map((n) =>
            n.id === notificationId ? { ...n, is_read: true, read_at: new Date().toISOString() } : n
          )
        )

        // Decrement unread count
        setUnreadCount((prev) => Math.max(0, prev - 1))

        // Invalidate cache
        apiCache.remove('notifications_{}')

        // Fetch fresh unread count
        await fetchUnreadCount()
      } catch (err) {
        debugError('Failed to mark notification as read:', err)
      }
    },
    [fetchUnreadCount]
  )

  const markAsUnread = useCallback(
    async (notificationId: number) => {
      try {
        await api.post(`/notifications/${notificationId}/unread`, {})

        // Update local state
        setNotifications((prev) =>
          prev.map((n) =>
            n.id === notificationId ? { ...n, is_read: false, read_at: undefined } : n
          )
        )

        // Increment unread count
        setUnreadCount((prev) => prev + 1)

        // Invalidate cache
        apiCache.remove('notifications_{}')

        // Fetch fresh unread count
        await fetchUnreadCount()
      } catch (err) {
        debugError('Failed to mark notification as unread:', err)
      }
    },
    [fetchUnreadCount]
  )

  const markAllAsRead = useCallback(async () => {
    try {
      await api.post('/notifications/mark-all-read', {})

      // Update local state
      setNotifications((prev) =>
        prev.map((n) => ({ ...n, is_read: true, read_at: new Date().toISOString() }))
      )

      // Reset unread count
      setUnreadCount(0)

      // Invalidate cache
      apiCache.remove('notifications_{}')

      // Fetch fresh data
      await fetchStats()
    } catch (err) {
      debugError('Failed to mark all notifications as read:', err)
    }
  }, [fetchStats])

  const deleteNotification = useCallback(async (notificationId: number) => {
    try {
      await api.delete(`/notifications/${notificationId}`)

      // Update local state
      const deletedNotification = notifications.find((n) => n.id === notificationId)
      setNotifications((prev) => prev.filter((n) => n.id !== notificationId))

      // Update unread count if deleted notification was unread
      if (deletedNotification && !deletedNotification.is_read) {
        setUnreadCount((prev) => Math.max(0, prev - 1))
      }

      // Invalidate cache
      apiCache.remove('notifications_{}')
    } catch (err) {
      debugError('Failed to delete notification:', err)
    }
  }, [notifications])

  const refreshNotifications = useCallback(
    async (filters?: NotificationFilters) => {
      // Invalidate all notification caches
      apiCache.remove('notifications_{}')
      await fetchNotifications(filters, true)
    },
    [fetchNotifications]
  )

  const clearAllCaches = useCallback(() => {
    apiCache.remove('notifications_{}')
  }, [])

  // Initial fetch on mount
  useEffect(() => {
    fetchNotifications()
    fetchStats()
    fetchUnreadCount()

    // Auto-refresh every 5 minutes
    const refreshInterval = setInterval(() => {
      fetchNotifications({}, true)
      fetchUnreadCount()
    }, 300000)

    // Refresh when page becomes visible
    const handleVisibilityChange = () => {
      if (!document.hidden) {
        fetchNotifications({}, true)
        fetchUnreadCount()
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)

    return () => {
      clearInterval(refreshInterval)
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [fetchNotifications, fetchStats, fetchUnreadCount])

  return {
    notifications,
    unreadCount,
    totalCount,
    loading,
    error,
    stats,
    fetchNotifications,
    fetchStats,
    fetchUnreadCount,
    markAsRead,
    markAsUnread,
    markAllAsRead,
    deleteNotification,
    refreshNotifications,
    clearAllCaches,
  }
}
