'use client'

import { selectPOSCustomerName, selectPOSCustomerPhone, usePOSStore } from '@/stores/posStore'

interface POSCustomerInputProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
}

export function POSCustomerInput({ locale, t }: POSCustomerInputProps) {
  const customerName = usePOSStore(selectPOSCustomerName)
  const customerPhone = usePOSStore(selectPOSCustomerPhone)
  const setCustomer = usePOSStore((state) => state.setCustomer)
  const isRTL = locale === 'ar'

  return (
    <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
      <div className="space-y-1">
        <label className={`block text-xs font-semibold text-gray-600 dark:text-gray-300 ${isRTL ? 'text-right' : ''}`}>
          {t('pos.customerName', 'Customer name')}
        </label>
        <input
          type="text"
          value={customerName}
          onChange={(event) => setCustomer(event.target.value, customerPhone)}
          placeholder={t('pos.customerNamePlaceholder', 'Walk-in')}
          className={`h-9 w-full rounded border border-gray-300 bg-white px-2 text-sm text-gray-900 outline-none transition focus:border-primary focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 ${
            isRTL ? 'text-right' : 'text-left'
          }`}
        />
      </div>
      <div className="space-y-1">
        <label className={`block text-xs font-semibold text-gray-600 dark:text-gray-300 ${isRTL ? 'text-right' : ''}`}>
          {t('pos.customerPhone', 'Phone')}
        </label>
        <input
          type="tel"
          value={customerPhone}
          onChange={(event) => setCustomer(customerName, event.target.value)}
          placeholder={t('pos.customerPhonePlaceholder', 'Optional')}
          className={`h-9 w-full rounded border border-gray-300 bg-white px-2 text-sm text-gray-900 outline-none transition focus:border-primary focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 ${
            isRTL ? 'text-right' : 'text-left'
          }`}
        />
      </div>
    </div>
  )
}

export default POSCustomerInput
