'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { OrderStatusBadge, type OrderStatus } from '@/components/commerce/OrderStatusBadge'
import { formatDate, formatPrice } from '@/lib/utils/format'

export interface OrdersTableItem {
  id: number
  orderNumber: string
  status: OrderStatus
  totalAmount: number
  createdAt: string
  itemsCount?: number
}

interface OrdersTableProps {
  locale: 'en' | 'ar'
  orders: OrdersTableItem[]
}

export function OrdersTable({ locale, orders }: OrdersTableProps) {
  const isArabic = locale === 'ar'

  if (orders.length === 0) {
    return (
      <div className="rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-surface))] p-8 text-center text-[hsl(var(--theme-text))]/70">
        {isArabic ? 'لا توجد طلبات' : 'No orders found'}
      </div>
    )
  }

  return (
    <>
      <div className="hidden overflow-hidden rounded-xl border border-[hsl(var(--theme-border))] md:block">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>#</TableHead>
              <TableHead>{isArabic ? 'التاريخ' : 'Date'}</TableHead>
              <TableHead>{isArabic ? 'العناصر' : 'Items'}</TableHead>
              <TableHead>{isArabic ? 'الحالة' : 'Status'}</TableHead>
              <TableHead>{isArabic ? 'الإجمالي' : 'Total'}</TableHead>
              <TableHead>{isArabic ? 'الإجراءات' : 'Actions'}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {orders.map((order) => (
              <TableRow key={order.id}>
                <TableCell className="font-semibold">{order.orderNumber}</TableCell>
                <TableCell>{formatDate(order.createdAt, locale)}</TableCell>
                <TableCell>{order.itemsCount ?? '-'}</TableCell>
                <TableCell><OrderStatusBadge status={order.status} locale={locale} /></TableCell>
                <TableCell>{formatPrice(order.totalAmount, locale)}</TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <Button asChild size="sm" variant="outline">
                      <Link href={`/${locale}/orders/${order.id}`}>{isArabic ? 'عرض' : 'View'}</Link>
                    </Button>
                    {order.status === 'delivered' ? (
                      <Button asChild size="sm">
                        <Link href={`/${locale}/menu`}>{isArabic ? 'إعادة الطلب' : 'Reorder'}</Link>
                      </Button>
                    ) : (
                      <Button asChild size="sm">
                        <Link href={`/${locale}/orders/track/${order.orderNumber}`}>{isArabic ? 'تتبع' : 'Track'}</Link>
                      </Button>
                    )}
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <div className="space-y-3 md:hidden">
        {orders.map((order) => (
          <div key={order.id} className={`rounded-xl border border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-surface))] p-4 ${isArabic ? 'text-right' : ''}`}>
            <div className={`mb-2 flex items-start justify-between ${isArabic ? 'flex-row-reverse' : ''}`}>
              <p className="font-semibold">{order.orderNumber}</p>
              <OrderStatusBadge status={order.status} locale={locale} />
            </div>
            <p className="text-sm text-[hsl(var(--theme-text))]/70">{formatDate(order.createdAt, locale)}</p>
            <p className="mt-1 text-sm text-[hsl(var(--theme-text))]/70">
              {isArabic ? 'العناصر' : 'Items'}: {order.itemsCount ?? '-'}
            </p>
            <p className="mt-1 font-semibold text-[hsl(var(--theme-primary))]">{formatPrice(order.totalAmount, locale)}</p>
            <div className={`mt-3 flex gap-2 ${isArabic ? 'flex-row-reverse' : ''}`}>
              <Button asChild size="sm" variant="outline" className="flex-1">
                <Link href={`/${locale}/orders/${order.id}`}>{isArabic ? 'عرض' : 'View'}</Link>
              </Button>
              <Button asChild size="sm" className="flex-1">
                <Link href={`/${locale}/orders/track/${order.orderNumber}`}>{isArabic ? 'تتبع' : 'Track'}</Link>
              </Button>
            </div>
          </div>
        ))}
      </div>
    </>
  )
}

