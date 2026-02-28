'use client'

import { memo } from 'react'
import { Minus, Plus, X } from 'lucide-react'
import type { POSCartItem as POSCartLine } from '@/stores/posStore'

interface POSCartItemProps {
  locale: 'en' | 'ar'
  item: POSCartLine
  onIncrease: () => void
  onDecrease: () => void
  onRemove: () => void
}

const formatPrice = (amount: number, locale: 'en' | 'ar') =>
  new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount)

function POSCartItemComponent({
  locale,
  item,
  onIncrease,
  onDecrease,
  onRemove,
}: POSCartItemProps) {
  const isRTL = locale === 'ar'

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-2 dark:border-gray-700 dark:bg-gray-900">
      <div className={`flex items-start justify-between gap-2 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-semibold">{item.name}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {formatPrice(item.price, locale)}
          </p>
        </div>
        <button
          type="button"
          onClick={onRemove}
          className="rounded p-1 text-gray-500 transition hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      <div className={`mt-2 flex items-center justify-between ${isRTL ? 'flex-row-reverse' : ''}`}>
        <div className={`flex items-center gap-1 ${isRTL ? 'flex-row-reverse' : ''}`}>
          <button
            type="button"
            onClick={onDecrease}
            className="inline-flex h-7 w-7 items-center justify-center rounded border border-gray-300 text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:text-gray-200 dark:hover:bg-gray-800"
          >
            <Minus className="h-3.5 w-3.5" />
          </button>
          <span className="min-w-6 text-center text-sm font-semibold tabular-nums">{item.quantity}</span>
          <button
            type="button"
            onClick={onIncrease}
            className="inline-flex h-7 w-7 items-center justify-center rounded border border-gray-300 text-gray-700 hover:bg-gray-100 dark:border-gray-700 dark:text-gray-200 dark:hover:bg-gray-800"
          >
            <Plus className="h-3.5 w-3.5" />
          </button>
        </div>

        <span className="text-sm font-semibold text-primary">
          {formatPrice(item.price * item.quantity, locale)}
        </span>
      </div>
    </div>
  )
}

export const POSCartItem = memo(POSCartItemComponent)
POSCartItem.displayName = 'POSCartItem'

export default POSCartItem
