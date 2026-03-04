'use client'

import { Badge } from '@/components/ui/badge'

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'preparing'
  | 'ready'
  | 'out_for_delivery'
  | 'delivered'
  | 'cancelled'

interface OrderStatusBadgeProps {
  status: OrderStatus
  locale?: 'en' | 'ar'
}

const statusStyles: Record<OrderStatus, string> = {
  pending: 'bg-yellow-100 text-yellow-900 dark:bg-yellow-900/30 dark:text-yellow-200',
  confirmed: 'bg-blue-100 text-blue-900 dark:bg-blue-900/30 dark:text-blue-200',
  preparing: 'bg-blue-100 text-blue-900 dark:bg-blue-900/30 dark:text-blue-200',
  ready: 'bg-emerald-100 text-emerald-900 dark:bg-emerald-900/30 dark:text-emerald-200',
  out_for_delivery: 'bg-orange-100 text-orange-900 dark:bg-orange-900/30 dark:text-orange-200',
  delivered: 'bg-green-100 text-green-900 dark:bg-green-900/30 dark:text-green-200',
  cancelled: 'bg-red-100 text-red-900 dark:bg-red-900/30 dark:text-red-200',
}

const labelsEn: Record<OrderStatus, string> = {
  pending: 'Pending',
  confirmed: 'Confirmed',
  preparing: 'Preparing',
  ready: 'Ready',
  out_for_delivery: 'Out for Delivery',
  delivered: 'Delivered',
  cancelled: 'Cancelled',
}

const labelsAr: Record<OrderStatus, string> = {
  pending: 'قيد الانتظار',
  confirmed: 'تم التأكيد',
  preparing: 'جاري التحضير',
  ready: 'جاهز',
  out_for_delivery: 'في الطريق',
  delivered: 'تم التوصيل',
  cancelled: 'ملغي',
}

export function OrderStatusBadge({ status, locale = 'en' }: OrderStatusBadgeProps) {
  return (
    <Badge className={statusStyles[status]}>
      {locale === 'ar' ? labelsAr[status] : labelsEn[status]}
    </Badge>
  )
}

