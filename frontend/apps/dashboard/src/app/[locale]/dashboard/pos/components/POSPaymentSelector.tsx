'use client'

import { Banknote, CreditCard, Wallet } from 'lucide-react'
import { selectPOSPaymentMethod, usePOSStore } from '@/stores/posStore'

interface POSPaymentSelectorProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
}

export function POSPaymentSelector({ locale, t }: POSPaymentSelectorProps) {
  const paymentMethod = usePOSStore(selectPOSPaymentMethod)
  const setPaymentMethod = usePOSStore((state) => state.setPaymentMethod)
  const isRTL = locale === 'ar'

  const methods = [
    { value: 'cash' as const, label: t('pos.cash', 'Cash'), icon: Banknote },
    { value: 'card' as const, label: t('pos.card', 'Card'), icon: CreditCard },
    { value: 'wallet' as const, label: t('pos.wallet', 'Wallet'), icon: Wallet },
  ]

  return (
    <div className="space-y-1">
      <label className={`block text-xs font-semibold text-gray-600 dark:text-gray-300 ${isRTL ? 'text-right' : ''}`}>
        {t('pos.paymentMethod', 'Payment')}
      </label>
      <div className={`grid grid-cols-3 gap-2 ${isRTL ? 'text-right' : ''}`}>
        {methods.map((method) => {
          const Icon = method.icon
          const active = paymentMethod === method.value

          return (
            <button
              key={method.value}
              type="button"
              onClick={() => setPaymentMethod(method.value)}
              className={`inline-flex h-9 items-center justify-center gap-1 rounded border text-xs font-semibold transition ${
                active
                  ? 'border-primary bg-primary text-white'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              <Icon className="h-3.5 w-3.5" />
              <span>{method.label}</span>
            </button>
          )
        })}
      </div>
    </div>
  )
}

export default POSPaymentSelector
