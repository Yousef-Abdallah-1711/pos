'use client'

import { Check } from 'lucide-react'
import type { OrderStatus } from '@/components/commerce/OrderStatusBadge'

interface TimelineStep {
  status: OrderStatus
  timestamp?: string
}

interface OrderTimelineProps {
  locale: 'en' | 'ar'
  currentStatus: OrderStatus
  steps?: TimelineStep[]
}

const sequence: OrderStatus[] = ['pending', 'confirmed', 'preparing', 'out_for_delivery', 'delivered']

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

export function OrderTimeline({ locale, currentStatus, steps = [] }: OrderTimelineProps) {
  const currentIndex = sequence.indexOf(currentStatus)
  const isArabic = locale === 'ar'

  return (
    <div className={`space-y-4 ${isArabic ? 'text-right' : ''}`}>
      {sequence.map((status, index) => {
        const completed = currentIndex > index
        const current = currentIndex === index
        const detail = steps.find((step) => step.status === status)
        const label = isArabic ? labelsAr[status] : labelsEn[status]

        return (
          <div key={status} className={`flex gap-3 ${isArabic ? 'flex-row-reverse' : ''}`}>
            <div className="relative flex flex-col items-center">
              <div
                className={`z-10 flex h-6 w-6 items-center justify-center rounded-full border-2 ${
                  completed
                    ? 'border-green-600 bg-green-600 text-white'
                    : current
                      ? 'border-[hsl(var(--theme-primary))] bg-[hsl(var(--theme-primary))] text-white ring-4 ring-[hsl(var(--theme-primary))]/20'
                      : 'border-[hsl(var(--theme-border))] bg-[hsl(var(--theme-background))] text-transparent'
                }`}
              >
                {completed ? <Check className="h-3.5 w-3.5" /> : '.'}
              </div>
              {index < sequence.length - 1 ? (
                <span className={`h-8 w-0.5 ${completed ? 'bg-green-500' : 'bg-[hsl(var(--theme-border))]'}`} />
              ) : null}
            </div>
            <div className="pt-0.5">
              <p className={`text-sm font-semibold ${current ? 'text-[hsl(var(--theme-primary))]' : 'text-[hsl(var(--theme-text))]'}`}>{label}</p>
              {detail?.timestamp ? (
                <p className="text-xs text-[hsl(var(--theme-text))]/65">
                  {new Date(detail.timestamp).toLocaleString(isArabic ? 'ar-EG' : 'en-US')}
                </p>
              ) : null}
            </div>
          </div>
        )
      })}
    </div>
  )
}

