'use client'

import { ClipboardList, Eraser, PackageOpen } from 'lucide-react'
import type { RestaurantTable } from '@/lib/api/tablesApi'
import { selectPOSHeldOrders, selectPOSItems, selectPOSOrderNotes, selectPOSTotalItems, usePOSStore } from '@/stores/posStore'
import { POSCartItem } from './POSCartItem'
import { POSCustomerInput } from './POSCustomerInput'
import { POSOrderSummary } from './POSOrderSummary'
import { POSPaymentSelector } from './POSPaymentSelector'
import { POSSubmitButton } from './POSSubmitButton'
import { POSTableSelector } from './POSTableSelector'

interface POSCartPanelProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  tables: RestaurantTable[]
  isSubmitting: boolean
  onSubmitOrder: () => void
  onHoldOrder: () => void
  onRecallOrder: () => void
  onClearCart: () => void
  onChangeCartItemQuantity?: (productId: number, quantity: number) => void
  onRemoveCartItem?: (productId: number) => void
}

export function POSCartPanel({
  locale,
  t,
  tables,
  isSubmitting,
  onSubmitOrder,
  onHoldOrder,
  onRecallOrder,
  onClearCart,
  onChangeCartItemQuantity,
  onRemoveCartItem,
}: POSCartPanelProps) {
  const items = usePOSStore(selectPOSItems)
  const heldOrders = usePOSStore(selectPOSHeldOrders)
  const totalItems = usePOSStore(selectPOSTotalItems)
  const orderNotes = usePOSStore(selectPOSOrderNotes)
  const updateQuantity = usePOSStore((state) => state.updateQuantity)
  const removeItem = usePOSStore((state) => state.removeItem)
  const setOrderNotes = usePOSStore((state) => state.setOrderNotes)
  const isRTL = locale === 'ar'

  const handleChangeQuantity = (productId: number, quantity: number) => {
    if (onChangeCartItemQuantity) {
      onChangeCartItemQuantity(productId, quantity)
      return
    }
    updateQuantity(productId, quantity)
  }

  const handleRemoveItem = (productId: number) => {
    if (onRemoveCartItem) {
      onRemoveCartItem(productId)
      return
    }
    removeItem(productId)
  }

  return (
    <aside className="flex h-full flex-col rounded-xl border border-gray-200 bg-white p-3 shadow-sm dark:border-gray-800 dark:bg-gray-950">
      <div className={`flex items-center justify-between gap-2 ${isRTL ? 'flex-row-reverse text-right' : ''}`}>
        <div className={`flex items-center gap-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
          <ClipboardList className="h-4 w-4 text-primary" />
          <h2 className="text-sm font-bold">{t('pos.cart', 'Cart')}</h2>
          <span className="rounded-full bg-primary px-2 py-0.5 text-[10px] font-semibold text-white">
            {totalItems}
          </span>
        </div>
        <div className={`flex items-center gap-1 ${isRTL ? 'flex-row-reverse' : ''}`}>
          <button
            type="button"
            onClick={onHoldOrder}
            disabled={items.length === 0}
            className="rounded border border-amber-300 bg-amber-50 px-2 py-1 text-[10px] font-semibold text-amber-700 transition hover:bg-amber-100 disabled:cursor-not-allowed disabled:opacity-50 dark:border-amber-600/40 dark:bg-amber-500/10 dark:text-amber-300"
          >
            {t('pos.holdOrder', 'Hold')} (Ctrl+H)
          </button>
          {heldOrders.length > 0 && (
            <button
              type="button"
              onClick={onRecallOrder}
              className="rounded border border-sky-300 bg-sky-50 px-2 py-1 text-[10px] font-semibold text-sky-700 transition hover:bg-sky-100 dark:border-sky-600/40 dark:bg-sky-500/10 dark:text-sky-300"
            >
              {t('pos.recallOrder', 'Recall')} (Ctrl+R)
            </button>
          )}
          <button
            type="button"
            onClick={onClearCart}
            disabled={items.length === 0}
            className="rounded border border-red-300 bg-red-50 px-2 py-1 text-[10px] font-semibold text-red-700 transition hover:bg-red-100 disabled:cursor-not-allowed disabled:opacity-50 dark:border-red-600/40 dark:bg-red-500/10 dark:text-red-300"
          >
            <span className={`inline-flex items-center gap-1 ${isRTL ? 'flex-row-reverse' : ''}`}>
              <Eraser className="h-3 w-3" />
              {t('pos.clearCart', 'Clear')} (Esc)
            </span>
          </button>
        </div>
      </div>

      <div className="mt-3 min-h-0 flex-1 overflow-auto space-y-2">
        {items.length === 0 ? (
          <div className="flex h-full min-h-[180px] flex-col items-center justify-center rounded-lg border border-dashed border-gray-300 bg-gray-50 p-4 text-center text-gray-500 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300">
            <PackageOpen className="mb-2 h-6 w-6" />
            <p className="text-sm font-medium">{t('pos.emptyCart', 'Cart is empty')}</p>
            <p className="text-xs">{t('pos.tapToAdd', 'Tap a product to add it')}</p>
          </div>
        ) : (
          items.map((item) => (
            <POSCartItem
              key={item.productId}
              locale={locale}
              item={item}
              onDecrease={() =>
                handleChangeQuantity(item.productId, item.quantity - 1)
              }
              onIncrease={() =>
                handleChangeQuantity(item.productId, item.quantity + 1)
              }
              onRemove={() => handleRemoveItem(item.productId)}
            />
          ))
        )}
      </div>

      {items.length > 0 && (
        <div className="mt-3 space-y-3 border-t border-gray-200 pt-3 dark:border-gray-800">
          <POSTableSelector locale={locale} t={t} tables={tables} />
          <POSCustomerInput locale={locale} t={t} />
          <POSPaymentSelector locale={locale} t={t} />

          <div className="space-y-1">
            <label className={`block text-xs font-semibold text-gray-600 dark:text-gray-300 ${isRTL ? 'text-right' : ''}`}>
              {t('pos.orderNotes', 'Order notes')}
            </label>
            <textarea
              value={orderNotes}
              onChange={(event) => setOrderNotes(event.target.value)}
              rows={2}
              className={`w-full rounded border border-gray-300 bg-white px-2 py-1.5 text-sm text-gray-900 outline-none transition focus:border-primary focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 ${
                isRTL ? 'text-right' : 'text-left'
              }`}
              placeholder={t('pos.notes', 'Notes...')}
            />
          </div>

          <POSOrderSummary locale={locale} t={t} />
          <POSSubmitButton t={t} isSubmitting={isSubmitting} onSubmit={onSubmitOrder} />
        </div>
      )}
    </aside>
  )
}

export default POSCartPanel
