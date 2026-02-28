'use client'

import { selectPOSSubtotal, selectPOSTax, selectPOSTotal, usePOSStore } from '@/stores/posStore'

interface POSOrderSummaryProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
}

const formatPrice = (amount: number, locale: 'en' | 'ar') =>
  new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount)

export function POSOrderSummary({ locale, t }: POSOrderSummaryProps) {
  const subtotal = usePOSStore(selectPOSSubtotal)
  const tax = usePOSStore(selectPOSTax)
  const total = usePOSStore(selectPOSTotal)
  const isRTL = locale === 'ar'

  return (
    <div className="space-y-1 rounded-lg border border-gray-200 bg-gray-50 p-3 text-sm dark:border-gray-700 dark:bg-gray-900">
      <div className={`flex items-center justify-between tabular-nums ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
        <span className="text-gray-600 dark:text-gray-300">{t('pos.subtotal', 'Subtotal')}</span>
        <span className="font-medium">{formatPrice(subtotal, locale)}</span>
      </div>
      <div className={`flex items-center justify-between tabular-nums ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
        <span className="text-gray-600 dark:text-gray-300">{t('pos.tax', 'Tax')} (14%)</span>
        <span className="font-medium">{formatPrice(tax, locale)}</span>
      </div>
      <div className={`flex items-center justify-between border-t border-gray-200 pt-2 font-bold tabular-nums dark:border-gray-700 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
        <span>{t('pos.total', 'Total')}</span>
        <span className="text-primary">{formatPrice(total, locale)}</span>
      </div>
    </div>
  )
}

export default POSOrderSummary
