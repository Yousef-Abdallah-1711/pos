'use client'

import Link from 'next/link'
import { useEffect, useMemo, useRef, useState } from 'react'
import { usePathname } from 'next/navigation'
import {
  AlertCircle,
  ArrowRight,
  Bell,
  Boxes,
  ClipboardList,
  CreditCard,
  Loader2,
  Monitor,
  Package,
  Palette,
  Settings,
  ShoppingCart,
  Store,
  TrendingUp,
  Truck,
  Users,
} from 'lucide-react'

import { Button } from '@/components/ui/button'
import { api } from '@/lib/api'
import { getConfig } from '@/lib/config'
import { debugError } from '@/lib/debug'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { useAuthStore } from '@/stores/authStore'
import { useHasPermission } from '@/hooks/useRbac'
import { fetchNotificationStats } from '@/hooks/useNotifications'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'

interface DashboardOrder {
  id: number
  order_number?: string
  customer_name?: string
  status?: string
  payment_status?: string
  total?: number
  created_at?: string
}

interface TenantStatusSnapshot {
  status: string | null
  canAccess: boolean
  blockReason: string | null
  daysRemaining: number | null
  planName: string | null
}

interface DashboardOverviewData {
  totalOrders: number
  pendingOrders: number
  productsCount: number
  categoriesCount: number
  employeesCount: number
  driversCount: number
  unreadNotifications: number
  totalNotifications: number
  restaurantsCount: number
  recentPaidRevenue: number
  recentOrders: DashboardOrder[]
  tableStats: {
    total: number
    available: number
    occupied: number
    reserved: number
  }
  tablesFeatureEnabled: boolean
  tenantStatus: TenantStatusSnapshot
  unavailableSources: string[]
  lastUpdatedAt: string | null
}

const EMPTY_OVERVIEW: DashboardOverviewData = {
  totalOrders: 0,
  pendingOrders: 0,
  productsCount: 0,
  categoriesCount: 0,
  employeesCount: 0,
  driversCount: 0,
  unreadNotifications: 0,
  totalNotifications: 0,
  restaurantsCount: 0,
  recentPaidRevenue: 0,
  recentOrders: [],
  tableStats: {
    total: 0,
    available: 0,
    occupied: 0,
    reserved: 0,
  },
  tablesFeatureEnabled: true,
  tenantStatus: {
    status: null,
    canAccess: true,
    blockReason: null,
    daysRemaining: null,
    planName: null,
  },
  unavailableSources: [],
  lastUpdatedAt: null,
}

function toNumber(value: unknown, fallback = 0): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : fallback
}

function getHttpStatus(value: unknown): number | null {
  if (!value || typeof value !== 'object') {
    return null
  }

  const status = (value as { response?: { status?: number } }).response?.status
  return typeof status === 'number' ? status : null
}

const TABLES_ENDPOINT_BACKOFF_KEY = 'dashboard.tables-stats-backoff-until'
const TABLES_ENDPOINT_BACKOFF_MS = 5 * 60 * 1000

function isTablesStatsBackedOff() {
  if (typeof window === 'undefined') {
    return false
  }

  try {
    const raw = window.localStorage.getItem(TABLES_ENDPOINT_BACKOFF_KEY)
    if (!raw) {
      return false
    }
    const until = Number(raw)
    return Number.isFinite(until) && until > Date.now()
  } catch {
    return false
  }
}

function setTablesStatsBackoff() {
  if (typeof window === 'undefined') {
    return
  }

  try {
    window.localStorage.setItem(
      TABLES_ENDPOINT_BACKOFF_KEY,
      String(Date.now() + TABLES_ENDPOINT_BACKOFF_MS)
    )
  } catch {
    // Ignore storage write failures.
  }
}

function extractArrayFromPayload(payload: any, candidateKeys: string[] = []): any[] {
  if (Array.isArray(payload)) return payload
  for (const key of candidateKeys) {
    if (Array.isArray(payload?.[key])) return payload[key]
    if (Array.isArray(payload?.data?.[key])) return payload.data[key]
  }
  if (Array.isArray(payload?.data)) return payload.data
  return []
}

function extractOrdersPayload(payload: any) {
  const direct = payload?.orders ? payload : payload?.data?.orders ? payload.data : payload
  const orders = extractArrayFromPayload(direct, ['orders']) as DashboardOrder[]
  const total = toNumber(direct?.total, orders.length)
  return { orders, total }
}

function extractTotalFromCollection(payload: any, collectionKey: string) {
  const direct = payload?.[collectionKey] || payload?.total ? payload : payload?.data?.[collectionKey] || payload?.data?.total ? payload.data : payload
  const items = extractArrayFromPayload(direct, [collectionKey])
  const total = toNumber(direct?.total, items.length)
  return { items, total }
}

function extractTenantStatus(payload: any): TenantStatusSnapshot {
  const direct = payload?.success === true && payload?.data ? payload.data : payload?.data?.success === true && payload?.data?.data ? payload.data.data : (payload?.data ?? payload)

  return {
    status: typeof direct?.status === 'string' ? direct.status : null,
    canAccess: direct?.can_access !== false,
    blockReason: typeof direct?.block_reason === 'string' ? direct.block_reason : null,
    daysRemaining: typeof direct?.days_remaining === 'number' ? direct.days_remaining : null,
    planName: typeof direct?.plan?.name === 'string'
      ? direct.plan.name
      : typeof direct?.current_plan?.name === 'string'
        ? direct.current_plan.name
        : null,
  }
}

function formatCurrency(value: number, currency = 'EGP', locale: 'en' | 'ar' = 'en') {
  try {
    return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en', {
      style: 'currency',
      currency,
      maximumFractionDigits: 2,
    }).format(value)
  } catch {
    return `${value.toFixed(2)} ${currency}`
  }
}

function formatDateTime(value: string | undefined, locale: 'en' | 'ar' = 'en') {
  if (!value) return locale === 'ar' ? 'وقت غير معروف' : 'Unknown time'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return locale === 'ar' ? 'وقت غير معروف' : 'Unknown time'
  return date.toLocaleString(locale === 'ar' ? 'ar-EG' : 'en-US')
}

function orderStatusClass(status?: string) {
  const normalized = (status || '').toLowerCase()
  const map: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
    confirmed: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
    preparing: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300',
    ready: 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-300',
    delivered: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
    completed: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
    cancelled: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
    canceled: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
  }
  return map[normalized] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
}

function paymentStatusClass(status?: string) {
  const normalized = (status || '').toLowerCase()
  const map: Record<string, string> = {
    paid: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300',
    success: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300',
    pending: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300',
    failed: 'bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-300',
    unpaid: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
  }
  return map[normalized] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
}

export default function DashboardPage() {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)

  // Permission checks for each module
  const { hasPermission: canViewReports } = useHasPermission(MODULE_IDS.REPORTS, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewOrders } = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewProducts } = useHasPermission(MODULE_IDS.PRODUCTS, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewCategories } = useHasPermission(MODULE_IDS.CATEGORIES, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewHR } = useHasPermission(MODULE_IDS.HR, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewTables } = useHasPermission(MODULE_IDS.TABLES, PERMISSION_LEVELS.READ)
  const { hasPermission: canViewKDS } = useHasPermission(MODULE_IDS.KDS, PERMISSION_LEVELS.READ)
  const { hasPermission: canUsePOS } = useHasPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)

  const [overview, setOverview] = useState<DashboardOverviewData>(EMPTY_OVERVIEW)
  const [isLoading, setIsLoading] = useState(true)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [overviewError, setOverviewError] = useState<string | null>(null)
  const hasLoadedOverviewRef = useRef(false)

  const isAuthenticated = !!(user && token)
  const debugMode = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('debug') === 'true'

  const dashboardRoute = (path: string) => `/${locale}${path}`

  useEffect(() => {
    if (debugMode) return
    if (!isAuthenticated) {
      window.location.href = `${getConfig().landingUrl}/${locale}/auth/login`
    }
  }, [debugMode, isAuthenticated, locale])

  useEffect(() => {
    if (!isAuthenticated && !debugMode) return

    let cancelled = false

    const loadOverview = async () => {
      const firstLoad = !hasLoadedOverviewRef.current
      if (firstLoad) {
        setIsLoading(true)
      } else {
        setIsRefreshing(true)
      }
      setOverviewError(null)

      const shouldFetchTableStats = canViewTables && !isTablesStatsBackedOff()

      const requests = await Promise.allSettled([
        api.get('/admin/orders', { params: { page: 1, page_size: 5 } }),
        api.get('/admin/orders', { params: { page: 1, page_size: 1, status: 'pending' } }),
        api.get('/products', { params: { page: 1, page_size: 1 } }),
        api.get('/categories'),
        api.get('/hr/employees'),
        api.get('/admin/drivers'),
        fetchNotificationStats(),
        api.get('/tenant/status'),
        shouldFetchTableStats ? api.get('/tables/stats') : Promise.resolve({ data: { data: null } }),
      ])

      if (cancelled) return

      const unavailableSources: string[] = []

      const recentOrdersResult = requests[0]
      const pendingOrdersResult = requests[1]
      const productsResult = requests[2]
      const categoriesResult = requests[3]
      const employeesResult = requests[4]
      const driversResult = requests[5]
      const notificationsResult = requests[6]
      const tenantStatusResult = requests[7]
      const tableStatsResult = requests[8]

      const recentOrdersPayload = recentOrdersResult.status === 'fulfilled'
        ? extractOrdersPayload(recentOrdersResult.value.data)
        : (unavailableSources.push('orders'), { orders: [] as DashboardOrder[], total: 0 })

      const pendingOrdersPayload = pendingOrdersResult.status === 'fulfilled'
        ? extractOrdersPayload(pendingOrdersResult.value.data)
        : (unavailableSources.push('pending-orders'), { total: 0 })

      const productsPayload = productsResult.status === 'fulfilled'
        ? extractTotalFromCollection(productsResult.value.data, 'products')
        : (unavailableSources.push('products'), { total: 0, items: [] })

      const categoriesPayload = categoriesResult.status === 'fulfilled'
        ? { items: extractArrayFromPayload(categoriesResult.value.data, ['categories']), total: 0 }
        : (unavailableSources.push('categories'), { items: [] as any[], total: 0 })

      const employeesPayload = employeesResult.status === 'fulfilled'
        ? { items: extractArrayFromPayload(employeesResult.value.data, ['employees']) }
        : (unavailableSources.push('employees'), { items: [] as any[] })

      const driversPayload = driversResult.status === 'fulfilled'
        ? { items: extractArrayFromPayload(driversResult.value.data, ['drivers']) }
        : (unavailableSources.push('drivers'), { items: [] as any[] })

      const notificationsStats = notificationsResult.status === 'fulfilled'
        ? notificationsResult.value
        : (unavailableSources.push('notifications'), {
          total: 0,
          unread: 0,
          read: 0,
          by_module: {},
          by_priority: {},
        })

      const tenantStatus = tenantStatusResult.status === 'fulfilled'
        ? extractTenantStatus(tenantStatusResult.value.data)
        : (unavailableSources.push('tenant-status'), EMPTY_OVERVIEW.tenantStatus)

      const tableStatsRouteMissing =
        tableStatsResult.status === 'rejected' &&
        getHttpStatus(tableStatsResult.reason) === 404

      if (tableStatsRouteMissing) {
        setTablesStatsBackoff()
      }

      const tableStatsPayload = tableStatsResult.status === 'fulfilled'
        ? (tableStatsResult.value.data?.data ?? tableStatsResult.value.data ?? {})
        : {}

      if (canViewTables && tableStatsResult.status === 'rejected' && !tableStatsRouteMissing) {
        unavailableSources.push('tables')
      }

      const recentPaidRevenue = recentOrdersPayload.orders
        .filter((order) => ['paid', 'success'].includes((order.payment_status || '').toLowerCase()))
        .reduce((sum, order) => sum + toNumber(order.total, 0), 0)

      setOverview({
        totalOrders: recentOrdersPayload.total,
        pendingOrders: pendingOrdersPayload.total ?? 0,
        productsCount: productsPayload.total,
        categoriesCount: categoriesPayload.items.length,
        employeesCount: employeesPayload.items.length,
        driversCount: driversPayload.items.length,
        unreadNotifications: toNumber(notificationsStats?.unread, 0),
        totalNotifications: toNumber(notificationsStats?.total, 0),
        restaurantsCount: user?.restaurant_id ? 1 : 0,
        recentPaidRevenue,
        recentOrders: recentOrdersPayload.orders,
        tableStats: {
          total: toNumber(tableStatsPayload?.total_tables, 0),
          available: toNumber(tableStatsPayload?.available, 0),
          occupied: toNumber(tableStatsPayload?.occupied, 0),
          reserved: toNumber(tableStatsPayload?.reserved, 0),
        },
        tablesFeatureEnabled:
          canViewTables &&
          shouldFetchTableStats &&
          !tableStatsRouteMissing &&
          tableStatsResult.status === 'fulfilled',
        tenantStatus,
        unavailableSources,
        lastUpdatedAt: new Date().toISOString(),
      })
      hasLoadedOverviewRef.current = true

      if (unavailableSources.length > 0) {
        setOverviewError(`Some dashboard widgets could not load: ${unavailableSources.join(', ')}`)
      }

      setIsLoading(false)
      setIsRefreshing(false)
    }

    loadOverview().catch((error) => {
      if (cancelled) return
      debugError('[DashboardPage] Failed to load overview:', error)
      setOverviewError('Failed to load dashboard overview data.')
      setIsLoading(false)
      setIsRefreshing(false)
    })

    return () => {
      cancelled = true
    }
  }, [debugMode, isAuthenticated, token, user?.restaurant_id, canViewTables])

  const setupTasks = useMemo(() => {
    const items = [
      {
        key: 'restaurant',
        done: overview.restaurantsCount > 0,
        title: t('dashboard.setup.restaurantConnectedTitle', 'Restaurant connected'),
        description: t('dashboard.setup.restaurantConnectedDesc', 'Your tenant and restaurant are linked and ready to manage.'),
        href: dashboardRoute('/dashboard/settings/restaurant'),
      },
      {
        key: 'categories',
        done: overview.categoriesCount > 0,
        title: t('dashboard.setup.createCategoriesTitle', 'Create categories'),
        description: t('dashboard.setup.createCategoriesDesc', 'Organize your menu into sections like main courses, drinks, and desserts.'),
        href: dashboardRoute('/dashboard/categories/new'),
      },
      {
        key: 'products',
        done: overview.productsCount > 0,
        title: t('dashboard.setup.addProductsTitle', 'Add menu products'),
        description: t('dashboard.setup.addProductsDesc', 'Create products with pricing and inventory to start taking orders.'),
        href: dashboardRoute('/dashboard/products/new'),
      },
      {
        key: 'staff',
        done: overview.employeesCount > 0,
        title: t('dashboard.setup.addTeamTitle', 'Add your team'),
        description: t('dashboard.setup.addTeamDesc', 'Invite employees and assign roles for operations and cashier workflows.'),
        href: dashboardRoute('/dashboard/hr/employees'),
      },
      {
        key: 'theme',
        done: false,
        title: t('dashboard.setup.customizeThemeTitle', 'Customize storefront theme'),
        description: t('dashboard.setup.customizeThemeDesc', 'Adjust homepage sections, branding, header, and footer in Theme Builder.'),
        href: dashboardRoute('/dashboard/theme-builder'),
      },
    ]

    return items
  }, [overview.categoriesCount, overview.employeesCount, overview.productsCount, overview.restaurantsCount, locale])

  const completedSetupCount = setupTasks.filter((task) => task.done).length

  const quickActions = useMemo(() => {
    const actions = [
      canUsePOS
        ? { label: t('navigation.pos', 'Cashier POS'), href: dashboardRoute('/dashboard/pos'), icon: CreditCard }
        : null,
      canViewKDS
        ? { label: t('navigation.kds', 'Kitchen Display'), href: dashboardRoute('/dashboard/kds'), icon: Monitor }
        : null,
      canViewOrders
        ? { label: t('dashboard.quickActions.viewOrders', 'View Orders'), href: dashboardRoute('/dashboard/orders'), icon: ShoppingCart }
        : null,
      canViewProducts
        ? { label: t('dashboard.quickActions.addProduct', 'Add Product'), href: dashboardRoute('/dashboard/products/new'), icon: Package }
        : null,
      canViewCategories
        ? { label: t('dashboard.quickActions.addCategory', 'Add Category'), href: dashboardRoute('/dashboard/categories/new'), icon: Boxes }
        : null,
      { label: t('navigation.themeBuilder', 'Theme Builder'), href: dashboardRoute('/dashboard/theme-builder'), icon: Palette },
      canViewHR
        ? { label: t('dashboard.quickActions.team', 'Team'), href: dashboardRoute('/dashboard/hr/employees'), icon: Users }
        : null,
      { label: t('settings.title', 'Settings'), href: dashboardRoute('/dashboard/settings'), icon: Settings },
    ]

    return actions.filter((action) => action !== null)
  }, [canUsePOS, canViewCategories, canViewHR, canViewKDS, canViewOrders, canViewProducts, locale, t])

  const metricCards = useMemo(() => {
    const allCards = [
      {
        key: 'recentRevenue',
        label: t('dashboard.metrics.recentPaidRevenue', 'Recent Paid Revenue'),
        helper: t('dashboard.metrics.recentPaidRevenueHelp', 'From the latest loaded orders'),
        value: formatCurrency(overview.recentPaidRevenue, 'EGP', locale),
        icon: TrendingUp,
        color: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
        href: dashboardRoute('/dashboard/orders'),
      },
      {
        key: 'orders',
        label: t('dashboard.metrics.orders', 'Orders'),
        helper: t('dashboard.metrics.ordersHelp', 'Total orders in system'),
        value: overview.totalOrders.toLocaleString(),
        icon: ShoppingCart,
        color: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
        href: dashboardRoute('/dashboard/orders'),
      },
      {
        key: 'pendingOrders',
        label: t('dashboard.metrics.pendingOrders', 'Pending Orders'),
        helper: t('dashboard.metrics.pendingOrdersHelp', 'Awaiting processing'),
        value: overview.pendingOrders.toLocaleString(),
        icon: ClipboardList,
        color: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
        href: dashboardRoute('/dashboard/orders'),
      },
      {
        key: 'restaurants',
        label: t('dashboard.metrics.restaurants', 'Restaurants'),
        helper: t('dashboard.metrics.restaurantsHelp', 'Linked restaurants'),
        value: overview.restaurantsCount.toLocaleString(),
        icon: Store,
        color: 'bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300',
        href: dashboardRoute('/dashboard/settings/restaurant'),
      },
      {
        key: 'products',
        label: t('dashboard.metrics.products', 'Products'),
        helper: t('dashboard.metrics.productsHelp', 'Menu items'),
        value: overview.productsCount.toLocaleString(),
        icon: Package,
        color: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-900/30 dark:text-cyan-300',
        href: dashboardRoute('/dashboard/products'),
      },
      {
        key: 'categories',
        label: t('dashboard.metrics.categories', 'Categories'),
        helper: t('dashboard.metrics.categoriesHelp', 'Menu sections'),
        value: overview.categoriesCount.toLocaleString(),
        icon: Boxes,
        color: 'bg-fuchsia-100 text-fuchsia-700 dark:bg-fuchsia-900/30 dark:text-fuchsia-300',
        href: dashboardRoute('/dashboard/categories'),
      },
      {
        key: 'staff',
        label: t('dashboard.metrics.staff', 'Staff'),
        helper: t('dashboard.metrics.staffHelp', 'Employees'),
        value: overview.employeesCount.toLocaleString(),
        icon: Users,
        color: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300',
        href: dashboardRoute('/dashboard/hr/employees'),
      },
      {
        key: 'drivers',
        label: t('dashboard.metrics.drivers', 'Drivers'),
        helper: t('dashboard.metrics.driversHelp', 'Delivery team'),
        value: overview.driversCount.toLocaleString(),
        icon: Truck,
        color: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300',
        href: dashboardRoute('/dashboard/drivers'),
      },
      {
        key: 'alerts',
        label: t('dashboard.metrics.unreadNotifications', 'Unread Notifications'),
        helper: `${overview.totalNotifications.toLocaleString(locale === 'ar' ? 'ar-EG' : 'en-US')} ${t('dashboard.metrics.totalNotifications', 'total notifications')}`,
        value: overview.unreadNotifications.toLocaleString(),
        icon: Bell,
        color: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300',
        href: dashboardRoute('/dashboard/notifications'),
      },
    ]

    // Filter based on user permissions
    return allCards.filter((card) => {
      if (card.key === 'recentRevenue' && !canViewReports) return false
      if ((card.key === 'orders' || card.key === 'pendingOrders' || card.key === 'drivers') && !canViewOrders) return false
      if (card.key === 'products' && !canViewProducts) return false
      if (card.key === 'categories' && !canViewCategories) return false
      if (card.key === 'staff' && !canViewHR) return false
      // 'restaurants' and 'alerts' are always visible
      return true
    })
  }, [canViewReports, canViewOrders, canViewProducts, canViewCategories, canViewHR, overview, locale, t])

  if (!isAuthenticated && !debugMode) {
    return (
      <div className="min-h-[40vh] flex items-center justify-center">
        <div className="flex items-center gap-3 text-gray-600 dark:text-gray-300">
          <Loader2 className="h-5 w-5 animate-spin" />
          <span>{t('auth.redirectingToLogin', 'Redirecting to login...')}</span>
        </div>
      </div>
    )
  }

  const userDisplayName = user?.name || user?.email || 'User'

  return (
    <div className="space-y-6">
      {overviewError && (
        <div className="rounded-xl border border-amber-200 bg-amber-50 p-4 text-amber-900 dark:border-amber-900/50 dark:bg-amber-950/40 dark:text-amber-200">
          <div className="flex items-start gap-3">
            <AlertCircle className="mt-0.5 h-5 w-5 flex-shrink-0" />
            <div>
              <p className="font-semibold">{t('dashboard.partialLoadTitle', 'Dashboard data partially loaded')}</p>
              <p className="text-sm opacity-90">{overviewError}</p>
            </div>
          </div>
        </div>
      )}

      <section className="rounded-2xl border border-gray-200 bg-gradient-to-br from-white to-gray-50 p-6 shadow-sm dark:border-gray-800 dark:from-gray-900 dark:to-gray-950 sm:p-8">
        <div className="flex flex-col gap-6 xl:flex-row xl:items-start xl:justify-between">
          <div className="space-y-4">
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
                {t('dashboard.welcome', 'Welcome')}, {userDisplayName}
              </h1>
              <p className="mt-2 text-sm text-gray-600 dark:text-gray-400 sm:text-base">
                {t('dashboard.overviewWithLiveData', 'Here is an overview of your restaurant performance with live data from orders, menu, staff, and notifications.')}
              </p>
            </div>

            <div className="flex flex-wrap gap-2">
              <span className="inline-flex items-center rounded-full border border-gray-200 bg-white px-3 py-1 text-xs font-medium text-gray-700 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300">
                {t('dashboard.roleLabel', 'Role')}: <span className={`${locale === 'ar' ? 'mr-1' : 'ml-1'} capitalize`}>{user?.role || 'owner'}</span>
              </span>
              <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ${
                overview.tenantStatus.canAccess
                  ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300'
                  : 'bg-rose-100 text-rose-800 dark:bg-rose-900/30 dark:text-rose-300'
              }`}>
                {t('dashboard.accessLabel', 'Access')}: {overview.tenantStatus.canAccess ? t('dashboard.status.active', 'Active') : t('dashboard.status.blocked', 'Blocked')}
              </span>
              {overview.tenantStatus.status && (
                <span className="inline-flex items-center rounded-full border border-indigo-200 bg-indigo-50 px-3 py-1 text-xs font-medium text-indigo-800 dark:border-indigo-900/40 dark:bg-indigo-900/20 dark:text-indigo-300">
                  {t('dashboard.subscriptionLabel', 'Subscription')}: {overview.tenantStatus.status}
                </span>
              )}
              {typeof overview.tenantStatus.daysRemaining === 'number' && (
                <span className="inline-flex items-center rounded-full border border-amber-200 bg-amber-50 px-3 py-1 text-xs font-medium text-amber-800 dark:border-amber-900/40 dark:bg-amber-900/20 dark:text-amber-300">
                  {overview.tenantStatus.daysRemaining} {t('dashboard.daysRemaining', 'day(s) remaining')}
                </span>
              )}
            </div>
          </div>

          <div className="flex w-full flex-col gap-3 xl:w-auto xl:min-w-[320px]">
            <div className="grid grid-cols-2 gap-2">
              {quickActions.slice(0, 4).map((action) => (
                <Link key={action.label} href={action.href}>
                  <Button variant="outline" className="w-full justify-start gap-2 bg-white/70 dark:bg-gray-800/70">
                    <action.icon className="h-4 w-4" />
                    <span className="truncate">{action.label}</span>
                  </Button>
                </Link>
              ))}
            </div>
            <div className="grid grid-cols-2 gap-2">
              {quickActions.slice(4).map((action) => (
                <Link key={action.label} href={action.href}>
                  <Button variant="outline" className="w-full justify-start gap-2 bg-white/70 dark:bg-gray-800/70">
                    <action.icon className="h-4 w-4" />
                    <span className="truncate">{action.label}</span>
                  </Button>
                </Link>
              ))}
            </div>
            <div className="text-xs text-gray-500 dark:text-gray-400">
              {overview.lastUpdatedAt ? `${t('reports.common.lastUpdated', 'Last updated')}: ${formatDateTime(overview.lastUpdatedAt, locale)}` : t('dashboard.loadingData', 'Loading dashboard data...')}
              {isRefreshing && ` ${t('dashboard.refreshing', '(refreshing)')}`}
            </div>
          </div>
        </div>
      </section>

      {isLoading && !overview.lastUpdatedAt ? (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, index) => (
            <div key={index} className="h-32 animate-pulse rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900" />
          ))}
        </div>
      ) : (
        <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 2xl:grid-cols-3">
          {metricCards.map((card) => (
            <Link
              key={card.key}
              href={card.href}
              className="group rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md dark:border-gray-800 dark:bg-gray-900"
            >
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-medium text-gray-600 dark:text-gray-400">{card.label}</p>
                  <p className="mt-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{card.value}</p>
                  <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">{card.helper}</p>
                </div>
                <div className={`rounded-lg p-3 ${card.color}`}>
                  <card.icon className="h-5 w-5" />
                </div>
              </div>
              <div className="mt-4 flex items-center gap-1 text-sm font-medium text-primary">
                {t('common.open', 'Open')}
                <ArrowRight className={`h-4 w-4 transition-transform ${locale === 'ar' ? 'rotate-180 group-hover:-translate-x-0.5' : 'group-hover:translate-x-0.5'}`} />
              </div>
            </Link>
          ))}
        </section>
      )}

      <section className="grid grid-cols-1 gap-6 xl:grid-cols-3">
        {canViewOrders && (
          <div className="xl:col-span-2 rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-gray-800 dark:bg-gray-900">
            <div className="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-gray-800">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{t('dashboard.recentOrdersTitle', 'Recent Orders')}</h2>
                <p className="text-sm text-gray-500 dark:text-gray-400">{t('dashboard.recentOrdersDesc', 'Latest orders from your storefront and POS channels')}</p>
              </div>
              <Link href={dashboardRoute('/dashboard/orders')}>
                <Button variant="outline" size="sm">{t('dashboard.viewAllOrders', 'View all orders')}</Button>
              </Link>
            </div>

            {overview.recentOrders.length === 0 ? (
              <div className="px-5 py-12 text-center">
                <ShoppingCart className="mx-auto h-10 w-10 text-gray-300 dark:text-gray-600" />
                <p className="mt-3 font-medium text-gray-900 dark:text-white">{t('dashboard.noOrdersYet', 'No orders yet')}</p>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  {t('dashboard.noOrdersDesc', 'Orders will appear here once customers start placing them.')}
                </p>
                <div className="mt-4">
                  <Link href={dashboardRoute('/dashboard/theme-builder')}>
                    <Button size="sm" className="gap-2">
                      <Palette className="h-4 w-4" />
                      {t('dashboard.quickActions.openThemeBuilder', 'Open Theme Builder')}
                    </Button>
                  </Link>
                </div>
              </div>
            ) : (
              <div className="divide-y divide-gray-100 dark:divide-gray-800">
                {overview.recentOrders.map((order) => (
                  <div key={order.id} className="flex flex-col gap-3 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <p className="font-semibold text-gray-900 dark:text-white">
                          #{order.order_number || order.id}
                        </p>
                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${orderStatusClass(order.status)}`}>
                          {order.status ? t(`orders.filters.${order.status.toLowerCase()}`, order.status) : t('dashboard.unknown', 'Unknown')}
                        </span>
                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${paymentStatusClass(order.payment_status)}`}>
                          {order.payment_status ? t(`orders.paymentFilters.${order.payment_status.toLowerCase()}`, order.payment_status) : t('orders.paymentFilters.pending', 'Pending')}
                        </span>
                      </div>
                      <p className="mt-1 truncate text-sm text-gray-600 dark:text-gray-400">
                        {order.customer_name || t('dashboard.unknownCustomer', 'Unknown customer')} - {formatDateTime(order.created_at, locale)}
                      </p>
                    </div>

                    <div className="flex items-center justify-between gap-4 sm:justify-end">
                      <p className="text-sm font-semibold text-gray-900 dark:text-white">
                        {formatCurrency(toNumber(order.total, 0), 'EGP', locale)}
                      </p>
                      <Link href={dashboardRoute(`/dashboard/orders/${order.id}`)}>
                        <Button variant="outline" size="sm">{t('common.open', 'Open')}</Button>
                      </Link>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        <div className="space-y-6">
          {canViewTables && overview.tablesFeatureEnabled && (
            <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-800 dark:bg-gray-900">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {t('dashboard.tablesOverviewTitle', 'Tables Overview')}
                </h2>
                <Link href={dashboardRoute('/dashboard/tables')}>
                  <Button variant="outline" size="sm">
                    {t('dashboard.manageTables', 'Manage Tables')}
                  </Button>
                </Link>
              </div>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {t('dashboard.tablesOverviewDesc', 'Track dine-in table availability and occupancy at a glance.')}
              </p>

              <div className="mt-4 grid grid-cols-3 gap-3">
                <div className="rounded-lg bg-emerald-50 px-3 py-2 text-center dark:bg-emerald-900/20">
                  <p className="text-xs text-emerald-700 dark:text-emerald-300">{t('tables.stats.available')}</p>
                  <p className="mt-1 text-xl font-bold text-emerald-800 dark:text-emerald-200">
                    {overview.tableStats.available}
                  </p>
                </div>
                <div className="rounded-lg bg-red-50 px-3 py-2 text-center dark:bg-red-900/20">
                  <p className="text-xs text-red-700 dark:text-red-300">{t('tables.stats.occupied')}</p>
                  <p className="mt-1 text-xl font-bold text-red-800 dark:text-red-200">
                    {overview.tableStats.occupied}
                  </p>
                </div>
                <div className="rounded-lg bg-amber-50 px-3 py-2 text-center dark:bg-amber-900/20">
                  <p className="text-xs text-amber-700 dark:text-amber-300">{t('tables.stats.reserved')}</p>
                  <p className="mt-1 text-xl font-bold text-amber-800 dark:text-amber-200">
                    {overview.tableStats.reserved}
                  </p>
                </div>
              </div>

              <p className="mt-3 text-xs text-gray-500 dark:text-gray-400">
                {t('tables.stats.total')}: {overview.tableStats.total}
              </p>
            </div>
          )}

          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-800 dark:bg-gray-900">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{t('dashboard.workspaceStatus', 'Workspace Status')}</h2>
            <div className="mt-4 space-y-3">
              <div className="flex items-center justify-between gap-3 rounded-lg bg-gray-50 px-3 py-2 dark:bg-gray-800/70">
                <span className="text-sm text-gray-600 dark:text-gray-400">{t('dashboard.tenantAccess', 'Tenant Access')}</span>
                <span className={`text-sm font-semibold ${overview.tenantStatus.canAccess ? 'text-emerald-600 dark:text-emerald-300' : 'text-rose-600 dark:text-rose-300'}`}>
                  {overview.tenantStatus.canAccess ? t('dashboard.status.active', 'Active') : t('dashboard.status.blocked', 'Blocked')}
                </span>
              </div>
              <div className="flex items-center justify-between gap-3 rounded-lg bg-gray-50 px-3 py-2 dark:bg-gray-800/70">
                <span className="text-sm text-gray-600 dark:text-gray-400">{t('dashboard.plan', 'Plan')}</span>
                <span className="text-sm font-semibold text-gray-900 dark:text-white">
                  {overview.tenantStatus.planName || overview.tenantStatus.status || t('dashboard.unknown', 'Unknown')}
                </span>
              </div>
              <div className="flex items-center justify-between gap-3 rounded-lg bg-gray-50 px-3 py-2 dark:bg-gray-800/70">
                <span className="text-sm text-gray-600 dark:text-gray-400">{t('dashboard.notifications', 'Notifications')}</span>
                <span className="text-sm font-semibold text-gray-900 dark:text-white">
                  {overview.unreadNotifications} {t('dashboard.unread', 'unread')} / {overview.totalNotifications} {t('dashboard.total', 'total')}
                </span>
              </div>
              <div className="flex items-center justify-between gap-3 rounded-lg bg-gray-50 px-3 py-2 dark:bg-gray-800/70">
                <span className="text-sm text-gray-600 dark:text-gray-400">{t('dashboard.dataSources', 'Data sources')}</span>
                <span className={`text-sm font-semibold ${overview.unavailableSources.length ? 'text-amber-600 dark:text-amber-300' : 'text-emerald-600 dark:text-emerald-300'}`}>
                  {overview.unavailableSources.length ? `${overview.unavailableSources.length} ${t('dashboard.unavailable', 'unavailable')}` : t('dashboard.allHealthy', 'All healthy')}
                </span>
              </div>
            </div>
          </div>

          <div className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-800 dark:bg-gray-900">
            <div className="flex items-center justify-between gap-3">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{t('dashboard.gettingStarted', 'Getting Started')}</h2>
              <span className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {completedSetupCount}/{setupTasks.length}
              </span>
            </div>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {t('dashboard.setupChecklistDesc', 'Complete the setup checklist to make the dashboard fully operational.')}
            </p>

            <div className="mt-4 space-y-3">
              {setupTasks.map((task, index) => (
                <Link
                  key={task.key}
                  href={task.href}
                  className="group flex items-start gap-3 rounded-xl border border-gray-200 p-3 transition hover:border-primary/40 hover:bg-primary/5 dark:border-gray-800 dark:hover:border-primary/30 dark:hover:bg-primary/10"
                >
                  <div className={`mt-0.5 flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full text-xs font-bold ${
                    task.done
                      ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
                      : 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300'
                  }`}>
                    {task.done ? '✓' : index + 1}
                  </div>
                  <div className="min-w-0">
                    <p className="font-medium text-gray-900 dark:text-white">{task.title}</p>
                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{task.description}</p>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
