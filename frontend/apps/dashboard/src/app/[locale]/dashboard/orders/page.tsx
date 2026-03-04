'use client'

import { useState, useEffect } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { useUserPermissions } from '@/hooks/useRbac'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { DriverAssignmentModal } from '@/components/driver-assignment-modal'
import { toast } from 'sonner'
import {
  Eye,
  XCircle,
  Filter,
  ChevronLeft,
  ChevronRight,
  Loader,
  AlertCircle,
  Lock,
  Download,
  Truck
} from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { exportOrdersToCSV } from '@/lib/pdf-export'
import { debugError } from '@/lib/debug'

interface Order {
  id: number
  order_number: string
  order_source?: string
  customer_name: string
  customer_email: string
  customer_phone: string
  delivery_address: string
  table_id?: number
  table_number?: number
  table_name?: string
  table_location?: string
  status: string
  payment_status: string
  payment_method: string
  subtotal: number
  tax: number
  delivery_fee: number
  discount: number
  total: number
  items_count: number
  created_at: string
}

export default function OrdersPage() {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  // Check RBAC permissions for Orders module
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewOrders = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ)
  const canDeleteOrders = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.DELETE)

  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [paymentFilter, setPaymentFilter] = useState<string>('all')
  const [tableFilter, setTableFilter] = useState<string>(() => searchParams.get('table_id') || 'all')
  const [searchQuery, setSearchQuery] = useState<string>('')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalOrders, setTotalOrders] = useState(0)
  const [isAssignmentModalOpen, setIsAssignmentModalOpen] = useState(false)
  const [selectedOrderForAssignment, setSelectedOrderForAssignment] = useState<Order | null>(null)
  const token = useAuthStore((state) => state.token)

  const fetchOrders = async () => {
    setLoading(true)
    setError(null)
    try {
      const params: any = {
        page,
        page_size: 20,
      }

      if (statusFilter !== 'all') {
        params.status = statusFilter
      }

      if (paymentFilter !== 'all') {
        params.payment_status = paymentFilter
      }

      if (tableFilter !== 'all') {
        params.table_id = tableFilter
      }

      if (searchQuery.trim()) {
        params.search = searchQuery.trim()
      }

      const response = await api.get('/admin/orders', { params })
      setOrders(response.data?.orders || [])
      setTotalOrders(response.data?.total || 0)
      setTotalPages(Math.ceil((response.data?.total || 0) / 20))
    } catch (error: any) {
      const errorMsg = error.response?.data?.error || error.message || t('orders.loadFailed');
      setError(errorMsg);
      toast.error(errorMsg)
      debugError('[ORDERS_PAGE] Error:', error)
    } finally {
      setLoading(false)
    }
  }

  // Move useEffect BEFORE early returns to avoid Hook order violation
  useEffect(() => {
    if (permissionLoading) {
      return
    }
    if (!token || !canViewOrders) {
      setLoading(false)
      return
    }
    fetchOrders()
  }, [token, statusFilter, paymentFilter, tableFilter, page, searchQuery, canViewOrders, permissionLoading])

  useEffect(() => {
    const urlTableID = searchParams.get('table_id')
    if (!urlTableID) {
      return
    }

    if (urlTableID !== tableFilter) {
      setTableFilter(urlTableID)
      setPage(1)
    }
  }, [searchParams, tableFilter])

  if (permissionLoading) {
    return (
      <div className="flex justify-center items-center h-96">
        <Loader className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  // Show permission denied if user can't view orders
  if (!canViewOrders) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            {t('orders.title')}
          </h1>
        </div>

        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                {t('orders.accessDeniedTitle', 'Access Denied')}
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                {t('orders.accessDeniedText', 'You do not have permission to view orders. Please contact your administrator.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const handleCancelOrder = async (orderId: number) => {
    if (!canDeleteOrders) {
      toast.error(t('orders.permissionDenied', 'You do not have permission to cancel orders.'))
      return
    }
    if (!confirm(t('orders.confirmCancel'))) return

    try {
      await api.delete(`/admin/orders/${orderId}`)
      setOrders(orders.filter((o) => o.id !== orderId))
      setTotalOrders(totalOrders - 1)
      toast.success(t('orders.cancelSuccess'))
    } catch (error) {
      toast.error(t('orders.cancelFailed'))
      debugError('[CANCEL_ORDER] Error:', error)
    }
  }

  const handleExportCSV = () => {
    try {
      const filename = `orders-export-${new Date().toISOString().split('T')[0]}.csv`
      exportOrdersToCSV(orders as any, filename, locale)
      toast.success(t('orders.exportCsvSuccess', 'Orders exported to CSV'))
    } catch (error: any) {
      debugError('[EXPORT_CSV] Error:', error)
      toast.error(t('orders.exportCsvFailed', 'Failed to export CSV'))
    }
  }

  const getStatusBadge = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
      confirmed: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
      preparing: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
      ready: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
      out_for_delivery: 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-400',
      delivered: 'bg-green-200 text-green-900 dark:bg-green-900/40 dark:text-green-300',
      cancelled: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
    }

    const statusLabel =
      locale === 'ar'
        ? t(`orders.filters.${status}`, status.replace(/_/g, ' '))
        : status.replace(/_/g, ' ').charAt(0).toUpperCase() + status.replace(/_/g, ' ').slice(1)

    return (
      <span className={`px-3 py-1 text-xs font-semibold rounded-full ${colors[status] || 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400'}`}>
        {statusLabel}
      </span>
    )
  }

  const getPaymentBadge = (paymentStatus: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
      paid: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
      failed: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
      refunded: 'bg-gray-100 text-gray-700 dark:bg-gray-900/30 dark:text-gray-400',
    }

    const paymentLabel =
      locale === 'ar'
        ? t(`orders.paymentFilters.${paymentStatus}`, paymentStatus)
        : paymentStatus.charAt(0).toUpperCase() + paymentStatus.slice(1)

    return (
      <span className={`px-3 py-1 text-xs font-semibold rounded-full ${colors[paymentStatus] || 'bg-gray-100 dark:bg-gray-900/30'}`}>
        {paymentLabel}
      </span>
    )
  }

  const getOrderSourceBadge = (orderSource?: string) => {
    const source = orderSource || 'website'
    const config: Record<string, { label: string; className: string }> = {
      website: {
        label: t('orders.sources.website', 'Online'),
        className: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
      },
      dine_in: {
        label: t('orders.sources.dine_in', 'Dine-in'),
        className: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300',
      },
      mobile_app: {
        label: t('orders.sources.mobile_app', 'Mobile'),
        className: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300',
      },
      phone: {
        label: t('orders.sources.phone', 'Phone'),
        className: 'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-300',
      },
      pos: {
        label: t('orders.sources.pos', 'POS'),
        className: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300',
      },
    }

    const resolved = config[source] || config.website
    return (
      <span className={`inline-flex rounded-full px-2 py-0.5 text-[11px] font-semibold ${resolved.className}`}>
        {resolved.label}
      </span>
    )
  }

  const isRTL = locale === 'ar'
  const iconGapClass = isRTL ? 'ml-1' : 'mr-1'
  const rangeStart = totalOrders > 0 ? (page - 1) * 20 + 1 : 0
  const rangeEnd = Math.min(page * 20, totalOrders)
  const tableOptions = Array.from(
    new Map(
      orders
        .filter((order) => typeof order.table_id === 'number' && typeof order.table_number === 'number')
        .map((order) => [
          order.table_id as number,
          {
            tableID: order.table_id as number,
            label: `${t('navigation.tables', 'Tables')} #${order.table_number}${order.table_name ? ` - ${order.table_name}` : ''}`,
          },
        ])
    ).values()
  )
  const tableFilterOptions =
    tableFilter !== 'all' &&
    !tableOptions.some((option) => String(option.tableID) === tableFilter)
      ? [
          ...tableOptions,
          {
            tableID: Number.parseInt(tableFilter, 10),
            label: `${t('navigation.tables', 'Tables')} #${tableFilter}`,
          },
        ]
      : tableOptions

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">
            {t('orders.title')}
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1">
            {t('orders.manageOrders')}
          </p>
        </div>
        {orders.length > 0 && (
          <Button
            variant="outline"
            size="sm"
            onClick={handleExportCSV}
            title={t('orders.exportCsvTitle', 'Export orders list as CSV')}
          >
            <Download className={`h-4 w-4 ${isRTL ? 'ml-2' : 'mr-2'}`} />
            <span className="hidden sm:inline">CSV</span>
          </Button>
        )}
      </div>

      {/* Filters and Search */}
      <div className="flex flex-col gap-4">
        {/* Search Bar */}
        <div className="flex items-center gap-2">
          <input
            type="text"
            placeholder={t('common.search')}
            value={searchQuery}
            onChange={(e) => {
              setSearchQuery(e.target.value)
              setPage(1)
            }}
            className="flex-1 px-4 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        {/* Status and Payment/Table Filters */}
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
          <div className="flex items-center gap-2">
            <Filter className="h-5 w-5 text-gray-400" />
            <select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value)
                setPage(1)
              }}
              className="px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="all">{t('orders.filters.all')}</option>
              <option value="pending">{t('orders.filters.pending')}</option>
              <option value="confirmed">{t('orders.filters.confirmed')}</option>
              <option value="preparing">{t('orders.filters.preparing')}</option>
              <option value="ready">{t('orders.filters.ready')}</option>
              <option value="out_for_delivery">{t('orders.filters.out_for_delivery')}</option>
              <option value="delivered">{t('orders.filters.delivered')}</option>
              <option value="cancelled">{t('orders.filters.cancelled')}</option>
            </select>
          </div>

          <div>
            <select
              value={paymentFilter}
              onChange={(e) => {
                setPaymentFilter(e.target.value)
                setPage(1)
              }}
              className="px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="all">{t('orders.paymentFilters.all')}</option>
              <option value="pending">{t('orders.paymentFilters.pending')}</option>
              <option value="paid">{t('orders.paymentFilters.paid')}</option>
              <option value="failed">{t('orders.paymentFilters.failed')}</option>
              <option value="refunded">{t('orders.paymentFilters.refunded')}</option>
            </select>
          </div>

          <div>
            <select
              value={tableFilter}
              onChange={(e) => {
                setTableFilter(e.target.value)
                setPage(1)
              }}
              className="px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="all">{t('orders.tableFilterAll', 'All Tables')}</option>
              {tableFilterOptions.map((option) => (
                <option key={option.tableID} value={option.tableID}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 flex gap-3">
          <AlertCircle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <h3 className="font-semibold text-red-800 dark:text-red-300">{t('orders.errorLoadingTitle', 'Error Loading Orders')}</h3>
            <p className="text-sm text-red-700 dark:text-red-400 mt-1">{error}</p>
          </div>
        </div>
      )}

      {/* Orders Table */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        {loading ? (
          <div className="flex justify-center items-center h-96">
            <Loader className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 dark:bg-gray-700">
                  <tr>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.orderNumber')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.customer')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.items')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.total')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.status')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.paymentStatus')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.createdAt')}
                    </th>
                    <th className={`px-6 py-3 text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider ${isRTL ? 'text-right' : 'text-left'}`}>
                      {t('orders.actions')}
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                  {orders.length === 0 ? (
                    <tr>
                      <td colSpan={8} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                        {t('orders.noOrders')}
                      </td>
                    </tr>
                  ) : (
                    orders.map((order) => (
                      <tr key={order.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <Link
                            href={`/${locale}/dashboard/orders/${order.id}`}
                            className="text-primary hover:underline font-medium"
                          >
                            {order.order_number}
                          </Link>
                          <div className="mt-1 flex flex-wrap items-center gap-1">
                            {getOrderSourceBadge(order.order_source)}
                            {typeof order.table_number === 'number' && (
                              <span className="inline-flex rounded-full border border-amber-300 bg-amber-50 px-2 py-0.5 text-[11px] font-semibold text-amber-800 dark:border-amber-500/40 dark:bg-amber-500/20 dark:text-amber-200">
                                {t('navigation.tables', 'Tables')} #{order.table_number}
                              </span>
                            )}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="text-sm">
                            <div className="font-medium text-gray-900 dark:text-white">
                              {order.customer_name}
                            </div>
                            <div className="text-gray-500 dark:text-gray-400">{order.customer_email}</div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                          {order.items_count || 0}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                          ${order.total.toFixed(2)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {getStatusBadge(order.status)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {getPaymentBadge(order.payment_status)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                          {new Date(order.created_at).toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm space-y-1">
                          <div className="flex flex-col gap-1">
                            <Link href={`/${locale}/dashboard/orders/${order.id}`}>
                              <Button variant="outline" size="sm" className="w-full">
                                <Eye className={`h-4 w-4 ${iconGapClass}`} />
                                {t('orders.viewDetails')}
                              </Button>
                            </Link>
                            {order.status !== 'cancelled' && order.status !== 'delivered' && order.status !== 'out_for_delivery' && (
                              <Button
                                variant="outline"
                                size="sm"
                                className="w-full text-primary hover:bg-primary/10"
                                onClick={() => {
                                  setSelectedOrderForAssignment(order)
                                  setIsAssignmentModalOpen(true)
                                }}
                              >
                                <Truck className={`h-4 w-4 ${iconGapClass}`} />
                                {t('drivers.assignDriver', 'Assign Driver')}
                              </Button>
                            )}
                            {order.status !== 'cancelled' && order.status !== 'delivered' && (
                              canDeleteOrders ? (
                                <Button
                                  variant="outline"
                                  size="sm"
                                  className="w-full text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20"
                                  onClick={() => handleCancelOrder(order.id)}
                                >
                                  <XCircle className={`h-4 w-4 ${iconGapClass}`} />
                                  {t('orders.cancelOrder')}
                                </Button>
                              ) : (
                                <Button
                                  size="sm"
                                  variant="destructive"
                                  disabled
                                  title={t('orders.permissionDeniedCancel', "You don't have permission to cancel orders")}
                                  className="w-full"
                                >
                                  <XCircle className={`h-4 w-4 ${iconGapClass}`} />
                                  {t('orders.cancelOrder')}
                                </Button>
                              )
                            )}
                          </div>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className={`flex justify-center items-center gap-4 p-6 border-t border-gray-200 dark:border-gray-700 ${isRTL ? 'flex-row-reverse' : ''}`}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage(Math.max(1, page - 1))}
                  disabled={page === 1}
                >
                  {isRTL ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
                </Button>
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  {`${t('orders.paginationPage', 'Page')} ${page} ${t('orders.paginationOf', 'of')} ${totalPages}`}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setPage(Math.min(totalPages, page + 1))}
                  disabled={page === totalPages}
                >
                  {isRTL ? <ChevronLeft className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                </Button>
              </div>
            )}

            {/* Summary */}
            {totalOrders > 0 && (
              <div className="px-6 py-3 bg-gray-50 dark:bg-gray-700/50 border-t border-gray-200 dark:border-gray-700 text-sm text-gray-600 dark:text-gray-400">
                {`${t('orders.showing', 'Showing')} ${rangeStart} ${t('orders.to', 'to')} ${rangeEnd} ${t('orders.of', 'of')} ${totalOrders} ${t('orders.ordersLabel', 'orders')}`}
              </div>
            )}
          </>
        )}
      </div>

      {/* Driver Assignment Modal */}
      {selectedOrderForAssignment && (
        <DriverAssignmentModal
          isOpen={isAssignmentModalOpen}
          orderId={selectedOrderForAssignment.id}
          orderNumber={selectedOrderForAssignment.order_number}
          onClose={() => {
            setIsAssignmentModalOpen(false)
            setSelectedOrderForAssignment(null)
          }}
          onAssignSuccess={() => {
            fetchOrders()
          }}
        />
      )}
    </div>
  )
}
