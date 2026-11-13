'use client'

import { useState, useEffect, useCallback } from 'react'
import { usePathname, useParams } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/stores/authStore'
import { useUserPermissions } from '@/hooks/useRbac'
import { api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { DriverAssignmentModal } from '@/components/driver-assignment-modal'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { toast } from 'sonner'
import {
  ArrowLeft,
  Package,
  Truck,
  Clock,
  DollarSign,
  MapPin,
  Phone,
  Mail,
  Loader,
  CheckCircle,
  Navigation,
  Lock,
  Download,
  FileText
} from 'lucide-react'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { exportOrderToPDF } from '@/lib/pdf-export'
import { debugError } from '@/lib/debug'
import { hasAuthenticatedSession } from '@/lib/authSession'

interface OrderItem {
  id: number
  product_id: number
  product_name: string
  quantity: number
  unit_price: number
  subtotal: number
  special_instructions?: string
}

interface StatusHistory {
  id: number
  status: string
  notes?: string
  changed_at: string
  changed_by?: string
}

interface Order {
  id: number
  order_number: string
  order_source?: string
  customer_name: string
  customer_email: string
  customer_phone: string
  delivery_address: string
  delivery_notes?: string
  status: string
  payment_status: string
  payment_method: string
  subtotal: number
  tax: number
  delivery_fee: number
  discount: number
  total: number
  created_at: string
  table?: {
    id: number
    number: number
    name?: string
    location?: string
    capacity?: number
    is_assigned: boolean
  } | null
  items: OrderItem[]
  status_history: StatusHistory[]
}

const STATUS_OPTIONS = [
  { value: 'pending', label: 'Pending' },
  { value: 'confirmed', label: 'Confirmed' },
  { value: 'preparing', label: 'Preparing' },
  { value: 'ready', label: 'Ready' },
  { value: 'out_for_delivery', label: 'Out for Delivery' },
  { value: 'delivered', label: 'Delivered' },
  { value: 'cancelled', label: 'Cancelled' },
]

const STATUS_TRANSITIONS: Record<string, string[]> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['preparing', 'cancelled'],
  preparing: ['ready'],
  ready: ['out_for_delivery', 'cancelled'],
  out_for_delivery: ['delivered'],
  delivered: [],
  cancelled: [],
}

const getStatusOptions = (currentStatus: string) => {
  const allowed = STATUS_TRANSITIONS[currentStatus] || []
  const allowedSet = new Set([currentStatus, ...allowed])
  const filtered = STATUS_OPTIONS.filter((option) => allowedSet.has(option.value))
  return filtered.length > 0 ? filtered : STATUS_OPTIONS
}

const STATUS_COLORS: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
  confirmed: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  preparing: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
  ready: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  out_for_delivery: 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-400',
  delivered: 'bg-green-200 text-green-900 dark:bg-green-900/40 dark:text-green-300',
  cancelled: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
}

const PAYMENT_COLORS: Record<string, string> = {
  pending: 'text-yellow-600 dark:text-yellow-400',
  paid: 'text-green-600 dark:text-green-400',
  failed: 'text-red-600 dark:text-red-400',
  refunded: 'text-gray-600 dark:text-gray-400',
}

export default function OrderDetailsPage() {
  const params = useParams()
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  // Check RBAC permissions for Orders module
  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canViewOrders = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.READ)
  const canUpdateOrders = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)

  const orderId = params.id as string
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [updatingStatus, setUpdatingStatus] = useState(false)
  const [selectedStatus, setSelectedStatus] = useState<string>('')
  const [isAssignmentModalOpen, setIsAssignmentModalOpen] = useState(false)
  const user = useAuthStore((state) => state.user)
  const token = useAuthStore((state) => state.token)
  const hasSession = hasAuthenticatedSession(user, token)

  const fetchOrder = useCallback(async () => {
    setLoading(true)
    setLoadError(null)
    try {
      const response = await api.get(`/admin/orders/${orderId}`)
      const orderData = response.data
      setOrder(orderData)
      setSelectedStatus(orderData.status)
    } catch (error: any) {
      setOrder(null)
      const status = error?.response?.status
      const message =
        status === 404
          ? t('orders.notFound', 'Order not found')
          : t('orders.loadFailed', 'Failed to load order details')
      setLoadError(message)
      toast.error(message)
      debugError('[ORDER_DETAILS] Error:', error)
    } finally {
      setLoading(false)
    }
  }, [orderId, t])

  useEffect(() => {
    if (permissionLoading) {
      return
    }
    if (!hasSession || !canViewOrders) {
      setLoading(false)
      return
    }

    void fetchOrder()
  }, [hasSession, canViewOrders, permissionLoading, fetchOrder])

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
        <Link href={`/${locale}/dashboard/orders`}>
          <Button variant="outline" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            {t('common.back')}
          </Button>
        </Link>

        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
          <div className="flex gap-3 items-start">
            <Lock className="h-5 w-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
            <div>
              <h2 className="font-semibold text-red-900 dark:text-red-300">
                {t('orders.accessDeniedTitle', 'Access Denied')}
              </h2>
              <p className="text-sm text-red-800 dark:text-red-400 mt-1">
                {t('orders.accessDeniedText', 'You do not have permission to view order details. Please contact your administrator.')}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const handleUpdateStatus = async () => {
    if (!selectedStatus || selectedStatus === order?.status) {
      toast.error('Please select a different status')
      return
    }

    setUpdatingStatus(true)
    try {
      await api.put(`/admin/orders/${orderId}/status`, {
        status: selectedStatus,
        notes: `Status updated to ${selectedStatus}`,
      })
      toast.success(t('orders.statusUpdateSuccess'))
      fetchOrder()
    } catch (error) {
      toast.error(t('orders.statusUpdateFailed'))
      debugError('[UPDATE_STATUS] Error:', error)
    } finally {
      setUpdatingStatus(false)
    }
  }

  const handleExportPDF = async () => {
    if (!order) return
    try {
      await exportOrderToPDF(order as any, `${order.order_number}.pdf`)
      toast.success('PDF exported successfully')
    } catch (error: any) {
      debugError('[EXPORT_PDF] Error:', error)
      toast.error(error.message || 'Failed to export PDF')
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-96">
        <Loader className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!order) {
    return (
      <div className="space-y-6">
        <Link href={`/${locale}/dashboard/orders`}>
          <Button variant="outline" size="sm">
            <ArrowLeft className="h-4 w-4 mr-2" />
            {t('common.back')}
          </Button>
        </Link>

        <div className="rounded-lg border border-gray-200 bg-white p-6 text-center dark:border-gray-700 dark:bg-gray-800">
          <FileText className="mx-auto h-10 w-10 text-gray-400" />
          <h2 className="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
            {loadError || t('orders.notFound', 'Order not found')}
          </h2>
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
            {t(
              'orders.notFoundDescription',
              'This order may have been removed or the link is invalid.'
            )}
          </p>
        </div>
      </div>
    )
  }

  const canAssignDriver = canUpdateOrders && !['cancelled', 'delivered'].includes(order.status)
  const isDineInOrder =
    order.order_source === 'dine_in' || (order.table?.is_assigned ?? false)

  const orderSourceLabel = (() => {
    switch (order.order_source) {
      case 'dine_in':
        return t('orders.sources.dine_in', 'Dine-in')
      case 'mobile_app':
        return t('orders.sources.mobile_app', 'Mobile')
      case 'phone':
        return t('orders.sources.phone', 'Phone')
      case 'website':
      default:
        return t('orders.sources.website', 'Online')
    }
  })()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link href={`/${locale}/dashboard/orders`}>
            <Button variant="outline" size="sm">
              <ArrowLeft className="h-4 w-4 mr-2" />
              {t('common.back')}
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {t('orders.orderNumber')} {order.order_number}
            </h1>
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <span className="inline-flex rounded-full border border-indigo-200 bg-indigo-50 px-2 py-1 text-xs font-semibold text-indigo-700 dark:border-indigo-500/40 dark:bg-indigo-500/15 dark:text-indigo-200">
                {orderSourceLabel}
              </span>
              {order.table?.is_assigned && (
                <span className="inline-flex rounded-full border border-amber-300 bg-amber-50 px-2 py-1 text-xs font-semibold text-amber-800 dark:border-amber-500/40 dark:bg-amber-500/15 dark:text-amber-200">
                  {t('navigation.tables', 'Tables')} #{order.table.number}
                  {order.table.name ? ` - ${order.table.name}` : ''}
                </span>
              )}
            </div>
            <p className="text-gray-500 dark:text-gray-400">
              {new Date(order.created_at).toLocaleDateString(locale === 'ar' ? 'ar-EG' : 'en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
              })}
            </p>
          </div>
        </div>

        {/* Export Actions */}
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleExportPDF}
            title="Export order as PDF invoice"
          >
            <Download className="h-4 w-4 mr-2" />
            <span className="hidden sm:inline">PDF</span>
          </Button>
        </div>
      </div>

      <Tabs defaultValue="details" className="space-y-6">
        <TabsList>
          <TabsTrigger value="details">{t('orders.details.orderInfo')}</TabsTrigger>
          <TabsTrigger value="drivers">{t('navigation.drivers')}</TabsTrigger>
        </TabsList>
        <TabsContent value="details">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Order Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* Order Items */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2 text-gray-900 dark:text-white">
              <Package className="h-5 w-5" />
              {t('orders.details.orderItems')}
            </h2>
            <div className="space-y-4">
              {order.items.map((item) => (
                <div key={item.id} className="flex justify-between items-start border-b border-gray-200 dark:border-gray-700 pb-4">
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-900 dark:text-white">{item.product_name}</h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {t('orders.details.quantity')}: {item.quantity} × ${item.unit_price.toFixed(2)}
                    </p>
                    {item.special_instructions && (
                      <p className="text-sm text-gray-600 dark:text-gray-300 italic mt-1">
                        {item.special_instructions}
                      </p>
                    )}
                  </div>
                  <div className="font-semibold text-gray-900 dark:text-white">
                    ${item.subtotal.toFixed(2)}
                  </div>
                </div>
              ))}
            </div>

            {/* Order Totals */}
            <div className="mt-6 space-y-2 border-t border-gray-200 dark:border-gray-700 pt-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600 dark:text-gray-400">{t('orders.details.subtotal')}</span>
                <span className="text-gray-900 dark:text-white">${order.subtotal.toFixed(2)}</span>
              </div>
              {order.tax > 0 && (
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600 dark:text-gray-400">{t('orders.details.tax')}</span>
                  <span className="text-gray-900 dark:text-white">${order.tax.toFixed(2)}</span>
                </div>
              )}
              <div className="flex justify-between text-sm">
                <span className="text-gray-600 dark:text-gray-400">{t('orders.details.deliveryFee')}</span>
                <span className="text-gray-900 dark:text-white">${order.delivery_fee.toFixed(2)}</span>
              </div>
              {order.discount > 0 && (
                <div className="flex justify-between text-sm text-green-600 dark:text-green-400">
                  <span>{t('orders.details.discount')}</span>
                  <span>-${order.discount.toFixed(2)}</span>
                </div>
              )}
              <div className="flex justify-between text-lg font-bold pt-2 border-t border-gray-200 dark:border-gray-700 text-gray-900 dark:text-white">
                <span>{t('orders.details.totalAmount')}</span>
                <span>${order.total.toFixed(2)}</span>
              </div>
            </div>
          </div>

          {/* Status History */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2 text-gray-900 dark:text-white">
              <Clock className="h-5 w-5" />
              {t('orders.details.statusHistory')}
            </h2>
            <div className="space-y-3">
              {order.status_history && order.status_history.length > 0 ? (
                order.status_history.map((history) => (
                  <div key={history.id} className="flex items-start gap-3">
                    <div className="w-2 h-2 bg-primary rounded-full mt-2 flex-shrink-0"></div>
                    <div className="flex-1">
                      <p className="font-medium text-gray-900 dark:text-white capitalize">
                        {history.status.replace(/_/g, ' ')}
                      </p>
                      {history.notes && (
                        <p className="text-sm text-gray-600 dark:text-gray-400">{history.notes}</p>
                      )}
                      <p className="text-xs text-gray-400 dark:text-gray-500">
                        {new Date(history.changed_at).toLocaleString(locale === 'ar' ? 'ar-EG' : 'en-US')}
                      </p>
                    </div>
                  </div>
                ))
              ) : (
                <p className="text-gray-500 dark:text-gray-400">{t('messages.noData')}</p>
              )}
            </div>
          </div>
        </div>

        {/* Right Column - Customer & Actions */}
        <div className="space-y-6">
          {/* Order Status */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
              {t('orders.details.orderInfo')}
            </h2>
            <div className="space-y-4">
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">
                  {t('orders.status')}
                </p>
                <span className={`inline-block px-3 py-1 text-sm font-semibold rounded-full ${STATUS_COLORS[order.status] || 'bg-gray-100 text-gray-800'}`}>
                  {order.status.replace(/_/g, ' ').charAt(0).toUpperCase() + order.status.replace(/_/g, ' ').slice(1)}
                </span>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">
                  {t('orders.paymentStatus')}
                </p>
                <span className={`text-sm font-semibold ${PAYMENT_COLORS[order.payment_status]}`}>
                  {order.payment_status.charAt(0).toUpperCase() + order.payment_status.slice(1)}
                </span>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">
                  {t('orders.details.paymentMethod')}
                </p>
                <p className="text-sm text-gray-900 dark:text-white">
                  {order.payment_method.replace(/_/g, ' ').charAt(0).toUpperCase() + order.payment_method.replace(/_/g, ' ').slice(1)}
                </p>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">
                  {t('orders.orderSource', 'Order Source')}
                </p>
                <p className="text-sm text-gray-900 dark:text-white">
                  {orderSourceLabel}
                </p>
              </div>
            </div>
          </div>

          {/* Customer Info */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
              {t('orders.details.customerInfo')}
            </h2>
            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <Mail className="h-5 w-5 text-gray-400 mt-0.5 flex-shrink-0" />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">{order.customer_name}</p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">{order.customer_email}</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Phone className="h-5 w-5 text-gray-400 mt-0.5 flex-shrink-0" />
                <p className="text-sm text-gray-900 dark:text-white">{order.customer_phone}</p>
              </div>
              {!isDineInOrder ? (
                <div className="flex items-start gap-3">
                  <MapPin className="h-5 w-5 text-gray-400 mt-0.5 flex-shrink-0" />
                  <div>
                    <p className="text-sm text-gray-900 dark:text-white">{order.delivery_address}</p>
                    {order.delivery_notes && (
                      <p className="text-sm text-gray-600 dark:text-gray-400 italic mt-1">
                        {order.delivery_notes}
                      </p>
                    )}
                  </div>
                </div>
              ) : (
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  {t('orders.dineInNoDelivery', 'Dine-in order, delivery details are not required.')}
                </p>
              )}
            </div>
          </div>

          {order.table?.is_assigned && (
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
                {t('navigation.tables', 'Tables')}
              </h2>
              <div className="space-y-2 text-sm text-gray-800 dark:text-gray-200">
                <p>
                  <span className="text-gray-500 dark:text-gray-400">{t('tables.tableNumber')}: </span>
                  {order.table.number}
                </p>
                <p>
                  <span className="text-gray-500 dark:text-gray-400">{t('tables.tableName')}: </span>
                  {order.table.name || t('common.notAvailable', 'N/A')}
                </p>
                <p>
                  <span className="text-gray-500 dark:text-gray-400">{t('tables.location')}: </span>
                  {order.table.location
                    ? t(`tables.locations.${order.table.location}`, order.table.location)
                    : t('common.notAvailable', 'N/A')}
                </p>
                <p>
                  <span className="text-gray-500 dark:text-gray-400">{t('tables.capacity')}: </span>
                  {order.table.capacity || 0}
                </p>
              </div>
            </div>
          )}

          {/* Delivery Tracking */}
          {(order.status === 'out_for_delivery' || order.status === 'delivered') && (
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white flex items-center gap-2">
                <Navigation className="h-5 w-5" />
                Delivery Tracking
              </h2>
              <div className="space-y-3">
                <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
                  <p className="text-sm text-blue-800 dark:text-blue-200">
                    {order.status === 'out_for_delivery'
                      ? '🚗 Your order is on the way! Our driver will arrive soon.'
                      : '✅ Your order has been delivered. Thank you for your purchase!'}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Update Status */}
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
              {t('orders.updateStatus')}
            </h2>
            {!canUpdateOrders ? (
              <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
                <p className="text-sm text-yellow-800 dark:text-yellow-200">
                  You don&apos;t have permission to update order status. Contact your administrator to enable this feature.
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                <select
                  value={selectedStatus}
                  onChange={(e) => setSelectedStatus(e.target.value)}
                  disabled={updatingStatus}
                  className="w-full px-3 py-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary disabled:opacity-50"
                >
                  {getStatusOptions(order.status).map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
                <Button
                  onClick={handleUpdateStatus}
                  disabled={updatingStatus || selectedStatus === order.status}
                  className="w-full"
                >
                  {updatingStatus ? (
                    <>
                      <Loader className="h-4 w-4 mr-2 animate-spin" />
                      {t('common.saving')}
                    </>
                  ) : (
                    <>
                      <CheckCircle className="h-4 w-4 mr-2" />
                      {t('orders.updateStatus')}
                    </>
                  )}
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>
        </TabsContent>
        <TabsContent value="drivers">
          <div className="space-y-6">
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                    {t('drivers.assignDriver')}
                  </h2>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {t('drivers.manageDrivers')}
                  </p>
                </div>
                <div className="flex flex-wrap gap-2">
                  <Link href={`/${locale}/dashboard/drivers`}>
                    <Button variant="outline" size="sm">
                      {t('navigation.drivers')}
                    </Button>
                  </Link>
                  <Button
                    size="sm"
                    onClick={() => setIsAssignmentModalOpen(true)}
                    disabled={!canAssignDriver}
                  >
                    {t('drivers.assignDriver')}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </TabsContent>
      </Tabs>
      <DriverAssignmentModal
        isOpen={isAssignmentModalOpen}
        orderId={order.id}
        orderNumber={order.order_number}
        onClose={() => setIsAssignmentModalOpen(false)}
        onAssignSuccess={fetchOrder}
      />
    </div>
  )
}
