'use client'

import type { POSCategory, POSProduct } from '@/lib/api/posApi'
import type { RestaurantTable } from '@/lib/api/tablesApi'
import { AlertCircle, Loader2 } from 'lucide-react'
import { selectPOSSelectedCategoryId, usePOSStore } from '@/stores/posStore'
import { POSCartPanel } from './POSCartPanel'
import { POSCategoryTabs } from './POSCategoryTabs'
import { POSHeader } from './POSHeader'
import { POSProductGrid } from './POSProductGrid'
import { POSSearchBar } from './POSSearchBar'
import type { POSProductCardSize } from '@/stores/posStore'

interface POSTerminalProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  loading: boolean
  error: string | null
  isOffline: boolean
  queuedOrderCount: number
  isReplayingQueue: boolean
  categories: POSCategory[]
  visibleProducts: POSProduct[]
  productCardSize: POSProductCardSize
  tables: RestaurantTable[]
  onAddProduct: (product: POSProduct) => void
  onRetry: () => void
  onSubmitOrder: () => void
  onHoldOrder: () => void
  onRecallOrder: () => void
  onClearCart: () => void
  onChangeCartItemQuantity?: (productId: number, quantity: number) => void
  onRemoveCartItem?: (productId: number) => void
  autoPrintEnabled: boolean
  onToggleAutoPrint: () => void
  isSubmitting: boolean
}

export function POSTerminal({
  locale,
  t,
  loading,
  error,
  isOffline,
  queuedOrderCount,
  isReplayingQueue,
  categories,
  visibleProducts,
  productCardSize,
  tables,
  onAddProduct,
  onRetry,
  onSubmitOrder,
  onHoldOrder,
  onRecallOrder,
  onClearCart,
  onChangeCartItemQuantity,
  onRemoveCartItem,
  autoPrintEnabled,
  onToggleAutoPrint,
  isSubmitting,
}: POSTerminalProps) {
  const selectedCategoryId = usePOSStore(selectPOSSelectedCategoryId)
  const setSelectedCategoryId = usePOSStore((state) => state.setSelectedCategoryId)
  const isRTL = locale === 'ar'

  return (
    <div
      dir={isRTL ? 'rtl' : 'ltr'}
      className="flex h-screen min-h-screen flex-col bg-gray-100 p-3 text-gray-900 dark:bg-slate-950 dark:text-gray-100"
    >
      <POSHeader
        locale={locale}
        t={t}
        isOffline={isOffline}
        queuedOrderCount={queuedOrderCount}
        isReplayingQueue={isReplayingQueue}
        autoPrintEnabled={autoPrintEnabled}
        onToggleAutoPrint={onToggleAutoPrint}
      />

      {isOffline && (
        <div className={`mt-2 rounded-lg border border-amber-300 bg-amber-50 px-3 py-2 text-xs font-medium text-amber-800 dark:border-amber-500/40 dark:bg-amber-500/10 dark:text-amber-200 ${isRTL ? 'text-right' : ''}`}>
          {t('pos.offlineBanner', 'Offline mode enabled. Orders will be queued and synced when online.')}
        </div>
      )}

      <div className={`mt-3 grid min-h-0 flex-1 gap-3 ${isRTL ? 'lg:[grid-template-columns:minmax(0,1fr)_minmax(0,1.85fr)]' : 'lg:[grid-template-columns:minmax(0,1.85fr)_minmax(0,1fr)]'}`}>
        <section className="flex min-h-0 flex-col gap-2">
          <POSSearchBar locale={locale} t={t} />
          <POSCategoryTabs
            locale={locale}
            t={t}
            categories={categories}
            selectedCategoryId={selectedCategoryId}
            onSelectCategory={setSelectedCategoryId}
          />

          {loading && visibleProducts.length === 0 ? (
            <div className="flex h-full min-h-[280px] items-center justify-center rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900">
              <div className="inline-flex items-center gap-2 text-sm text-gray-500 dark:text-gray-300">
                <Loader2 className="h-4 w-4 animate-spin" />
                <span>{t('common.loading', 'Loading...')}</span>
              </div>
            </div>
          ) : error ? (
            <div className="flex h-full min-h-[280px] flex-col items-center justify-center rounded-xl border border-red-300 bg-red-50 p-4 text-center dark:border-red-500/40 dark:bg-red-500/10">
              <AlertCircle className="mb-2 h-6 w-6 text-red-600 dark:text-red-300" />
              <p className="text-sm font-semibold text-red-700 dark:text-red-200">{error}</p>
              <button
                type="button"
                onClick={onRetry}
                className="mt-3 rounded border border-red-300 bg-white px-3 py-1.5 text-xs font-semibold text-red-700 hover:bg-red-50 dark:border-red-500/40 dark:bg-red-950 dark:text-red-200 dark:hover:bg-red-900"
              >
                {t('common.retry', 'Retry')}
              </button>
            </div>
          ) : (
            <POSProductGrid
              locale={locale}
              t={t}
              products={visibleProducts}
              productCardSize={productCardSize}
              onAddProduct={onAddProduct}
            />
          )}
        </section>

        <POSCartPanel
          locale={locale}
          t={t}
          tables={tables}
          isSubmitting={isSubmitting}
          onSubmitOrder={onSubmitOrder}
          onHoldOrder={onHoldOrder}
          onRecallOrder={onRecallOrder}
          onClearCart={onClearCart}
          onChangeCartItemQuantity={onChangeCartItemQuantity}
          onRemoveCartItem={onRemoveCartItem}
        />
      </div>
    </div>
  )
}

export default POSTerminal
