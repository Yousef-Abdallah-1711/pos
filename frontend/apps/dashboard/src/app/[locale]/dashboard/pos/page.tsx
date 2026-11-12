'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import { Loader2, Lock } from 'lucide-react'
import { toast } from 'sonner'
import { POSReceipt } from './components/POSReceipt'
import { POSTerminal } from './components/POSTerminal'
import { usePOSAudio } from '@/hooks/usePOSAudio'
import { usePOSBarcodeScanner } from '@/hooks/usePOSBarcodeScanner'
import { usePOSBackgroundSync } from '@/hooks/usePOSBackgroundSync'
import { usePOSData } from '@/hooks/usePOSData'
import { usePOSKeyboard } from '@/hooks/usePOSKeyboard'
import { usePOSOffline } from '@/hooks/usePOSOffline'
import { usePOSPerformance } from '@/hooks/usePOSPerformance'
import { usePOSProductIndex } from '@/hooks/usePOSProductIndex'
import { type POSReceiptData, usePOSPrint } from '@/hooks/usePOSPrint'
import { useUserPermissions } from '@/hooks/useRbac'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'
import { MODULE_IDS, PERMISSION_LEVELS } from '@/lib/rbac'
import { debugError } from '@/lib/debug'
import { type CreatePOSOrderPayload, posApi, type POSProduct } from '@/lib/api/posApi'
import { useAuthStore } from '@/stores/authStore'
import {
  selectPOSAutoPrintEnabled,
  selectPOSCustomerName,
  selectPOSCustomerPhone,
  selectPOSFocusModeEnabled,
  selectPOSHideOutOfStockProducts,
  selectPOSHeldOrders,
  selectPOSItems,
  selectPOSOrderNotes,
  selectPOSPaymentMethod,
  selectPOSProductCardSize,
  selectPOSSearchQuery,
  selectPOSSubtotal,
  selectPOSTableId,
  selectPOSTableName,
  selectPOSTax,
  selectPOSTotal,
  usePOSStore,
} from '@/stores/posStore'

const getErrorMessage = (error: unknown, fallback: string) => {
  if (error && typeof error === 'object' && 'response' in error) {
    const responseError = (error as { response?: { data?: { error?: string; message?: string } } })
      .response?.data
    if (responseError?.error) {
      return responseError.error
    }
    if (responseError?.message) {
      return responseError.message
    }
  }

  if (error instanceof Error && error.message) {
    return error.message
  }

  return fallback
}

const getPerfNow = () => {
  if (typeof window !== 'undefined' && window.performance) {
    return window.performance.now()
  }
  return Date.now()
}

export default function POSPage() {
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)

  const { checkPermission, loading: permissionLoading } = useUserPermissions()
  const canUsePOS = checkPermission(MODULE_IDS.ORDERS, PERMISSION_LEVELS.WRITE)

  const {
    products,
    categories,
    tables,
    loading: dataLoading,
    error: dataError,
    refresh,
  } = usePOSData()

  const {
    isOffline,
    queuedOrderCount,
    isReplaying,
    enqueueOrder,
    replayQueuedOrders,
  } = usePOSOffline()

  const { playSound } = usePOSAudio()
  const { printReceipt } = usePOSPrint()
  const { isDevMode, recentMetrics, markPageLoad, measureAsync, measureSync, recordMetric } =
    usePOSPerformance()

  const { registerOrderSync } = usePOSBackgroundSync({
    onSyncMessage: (message) => {
      if (message.type === 'pos-sync-complete') {
        if ((message.processed ?? 0) > 0) {
          toast.success(
            t('pos.offlineSyncSuccess', `${message.processed} queued order(s) synced`)
          )
        }
        void replayQueuedOrders()
        return
      }

      toast.error(t('pos.offlineSyncFailed', 'Some queued orders could not be synced'))
    },
  })

  const user = useAuthStore((state) => state.user)
  const items = usePOSStore(selectPOSItems)
  const tableId = usePOSStore(selectPOSTableId)
  const tableName = usePOSStore(selectPOSTableName)
  const customerName = usePOSStore(selectPOSCustomerName)
  const customerPhone = usePOSStore(selectPOSCustomerPhone)
  const paymentMethod = usePOSStore(selectPOSPaymentMethod)
  const orderNotes = usePOSStore(selectPOSOrderNotes)
  const subtotal = usePOSStore(selectPOSSubtotal)
  const tax = usePOSStore(selectPOSTax)
  const total = usePOSStore(selectPOSTotal)
  const selectedCategoryId = usePOSStore((state) => state.selectedCategoryId)
  const searchQuery = usePOSStore(selectPOSSearchQuery)
  const heldOrders = usePOSStore(selectPOSHeldOrders)
  const autoPrintEnabled = usePOSStore(selectPOSAutoPrintEnabled)
  const focusModeEnabled = usePOSStore(selectPOSFocusModeEnabled)
  const hideOutOfStockProducts = usePOSStore(selectPOSHideOutOfStockProducts)
  const productCardSize = usePOSStore(selectPOSProductCardSize)

  const addItem = usePOSStore((state) => state.addItem)
  const removeItem = usePOSStore((state) => state.removeItem)
  const updateQuantity = usePOSStore((state) => state.updateQuantity)
  const clearCart = usePOSStore((state) => state.clearCart)
  const holdCurrentOrder = usePOSStore((state) => state.holdCurrentOrder)
  const recallOrder = usePOSStore((state) => state.recallOrder)
  const toggleAutoPrint = usePOSStore((state) => state.toggleAutoPrint)
  const setTableContext = usePOSStore((state) => state.setTableContext)

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [latestReceipt, setLatestReceipt] = useState<POSReceiptData | null>(null)
  const wasOfflineRef = useRef(isOffline)
  const wasDataLoadingRef = useRef(dataLoading)
  const productFetchStartedAtRef = useRef(getPerfNow())
  const initialTableContextAppliedRef = useRef(false)

  const { search, getByBarcode } = usePOSProductIndex(products)

  const visibleProducts = useMemo(() => {
    const bySearch = search(searchQuery)
    if (!selectedCategoryId) {
      return bySearch
    }
    return bySearch.filter((product) => product.category_id === selectedCategoryId)
  }, [search, searchQuery, selectedCategoryId])

  const displayProducts = useMemo(() => {
    if (!hideOutOfStockProducts) {
      return visibleProducts
    }

    return visibleProducts.filter((product) => {
      const isOutOfStock =
        product.is_available === false ||
        (product.track_inventory === true && (product.quantity_in_stock || 0) <= 0)
      return !isOutOfStock
    })
  }, [hideOutOfStockProducts, visibleProducts])

  useEffect(() => {
    const startedAt = getPerfNow()
    const bySearch = search(searchQuery)
    if (selectedCategoryId) {
      bySearch.filter((product) => product.category_id === selectedCategoryId)
    }
    recordMetric('pos.search', getPerfNow() - startedAt)
  }, [recordMetric, search, searchQuery, selectedCategoryId])

  useEffect(() => {
    if (initialTableContextAppliedRef.current || tables.length === 0) {
      return
    }

    const tableParam = searchParams.get('table')
    const parsedTableID = tableParam ? Number(tableParam) : NaN
    if (!Number.isFinite(parsedTableID) || parsedTableID <= 0) {
      initialTableContextAppliedRef.current = true
      return
    }

    const selectedTable = tables.find((table) => table.id === parsedTableID)
    if (selectedTable) {
      setTableContext(selectedTable.id, selectedTable.name || null)
    }
    initialTableContextAppliedRef.current = true
  }, [searchParams, setTableContext, tables])

  const handleAddProduct = (product: POSProduct) => {
    const name = locale === 'ar' && product.name_ar ? product.name_ar : product.name
    measureSync('pos.add-to-cart', () => {
      addItem({
        id: product.id,
        name,
        price: product.price,
        image: product.image,
        track_inventory: product.track_inventory,
        quantity_in_stock: product.quantity_in_stock,
      })
    })
    void playSound('add')
  }

  const handleChangeCartItemQuantity = (productId: number, quantity: number) => {
    const currentItem = items.find((item) => item.productId === productId)
    if (!currentItem) {
      return
    }

    if (quantity <= 0) {
      removeItem(productId)
      void playSound('remove')
      return
    }

    updateQuantity(productId, quantity)
    if (quantity < currentItem.quantity) {
      void playSound('remove')
    } else if (quantity > currentItem.quantity) {
      void playSound('add')
    }
  }

  const handleRemoveCartItem = (productId: number) => {
    removeItem(productId)
    void playSound('remove')
  }

  const buildOrderPayload = (): CreatePOSOrderPayload => ({
    order_source: 'pos',
    customer_name: customerName || 'Walk-in',
    customer_phone: customerPhone.trim() || '0000000000',
    payment_method: paymentMethod,
    table_id: tableId || undefined,
    notes: orderNotes || undefined,
    items: items.map((item) => ({
      product_id: item.productId,
      quantity: item.quantity,
      special_instructions: item.notes || undefined,
    })),
  })

  const buildReceiptData = (orderNumber: string, createdAtISO?: string | null): POSReceiptData => ({
    locale,
    restaurantName: t('pos.receiptRestaurant', 'Restaurant'),
    orderNumber,
    createdAtISO: createdAtISO || new Date().toISOString(),
    paymentMethod,
    subtotal,
    tax,
    total,
    tableName: tableName || undefined,
    customerName: customerName || undefined,
    cashierName: user?.name || undefined,
    items: items.map((item) => ({
      name: item.name,
      quantity: item.quantity,
      unitPrice: item.price,
      totalPrice: item.price * item.quantity,
      notes: item.notes,
    })),
  })

  const handleSubmitOrder = async () => {
    if (items.length === 0 || isSubmitting) {
      return
    }

    setIsSubmitting(true)
    try {
      await measureAsync('pos.order-submit', async () => {
        const payload = buildOrderPayload()

        if (isOffline) {
          const queued = await enqueueOrder(payload)
          if (!queued) {
            throw new Error(t('pos.queueFailed', 'Failed to queue order for offline sync'))
          }

          void registerOrderSync()
          toast.success(t('pos.orderQueued', 'Order queued - will sync when online'))
          void playSound('success')
          clearCart()
          return
        }

        const createdOrder = await posApi.createOrder(payload)
        if (paymentMethod === 'cash' && createdOrder.id > 0) {
          try {
            await posApi.updateOrderPaymentStatus(createdOrder.id, 'paid')
          } catch (error) {
            debugError('[POS] Cash payment auto-mark failed:', error)
          }
        }

        toast.success(
          `${t('pos.orderSuccess', 'Order placed!')} #${createdOrder.order_number}`
        )

        const receiptData = buildReceiptData(
          createdOrder.order_number,
          createdOrder.created_at || null
        )
        setLatestReceipt(receiptData)
        void playSound('success')

        if (autoPrintEnabled) {
          const printed = await printReceipt(receiptData)
          if (!printed) {
            toast.error(t('pos.printFailed', 'Receipt printing failed'))
          }
        }

        clearCart()
      })
    } catch (error) {
      toast.error(getErrorMessage(error, t('pos.orderFailed', 'Failed to place order')))
      void playSound('error')
      debugError('[POS] Submit order failed:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleClearCart = () => {
    clearCart()
  }

  const handleHoldOrder = () => {
    const holdID = holdCurrentOrder()
    if (holdID) {
      toast.success(t('pos.orderHeld', 'Order held'))
    }
  }

  const handleRecallOrder = () => {
    const latestHeldOrder = heldOrders[0]
    if (!latestHeldOrder) {
      return
    }
    recallOrder(latestHeldOrder.id)
    toast.success(t('pos.orderRecalled', 'Held order recalled'))
  }

  usePOSBarcodeScanner({
    enabled: canUsePOS && !permissionLoading,
    getByBarcode,
    onProductFound: (product) => {
      handleAddProduct(product)
    },
    onBarcodeNotFound: (barcode) => {
      const template = t(
        'pos.productNotFoundBarcode',
        'Product not found: {barcode}'
      )
      toast.error(template.replace('{barcode}', barcode))
      void playSound('error')
    },
  })

  usePOSKeyboard({
    visibleProducts: displayProducts,
    cartItemsCount: items.length,
    isSubmitting,
    heldOrdersCount: heldOrders.length,
    onQuickAddProduct: handleAddProduct,
    onSubmitOrder: handleSubmitOrder,
    onClearCart: handleClearCart,
    onHoldOrder: handleHoldOrder,
    onRecallLastHeldOrder: handleRecallOrder,
    onFocusSearch: () => {
      window.dispatchEvent(new Event('pos:focus-search'))
    },
  })

  useEffect(() => {
    if (process.env.NODE_ENV !== 'production') {
      return
    }

    if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
      return
    }

    const scope = `/${locale}/dashboard/pos`
    void navigator.serviceWorker.register('/pos-sw.js', { scope }).catch((error) => {
      debugError('[POS] Service Worker registration failed:', error)
    })
  }, [locale])

  useEffect(() => {
    if (permissionLoading || dataLoading) {
      return
    }
    markPageLoad()
  }, [dataLoading, markPageLoad, permissionLoading])

  useEffect(() => {
    if (dataLoading && !wasDataLoadingRef.current) {
      productFetchStartedAtRef.current = getPerfNow()
    }

    if (!dataLoading && wasDataLoadingRef.current) {
      recordMetric('pos.product-fetch', getPerfNow() - productFetchStartedAtRef.current)
    }

    wasDataLoadingRef.current = dataLoading
  }, [dataLoading, recordMetric])

  useEffect(() => {
    const becameOnline = wasOfflineRef.current && !isOffline
    wasOfflineRef.current = isOffline
    if (!becameOnline || queuedOrderCount <= 0) {
      return
    }

    void (async () => {
      const result = await replayQueuedOrders()
      if (result.processed > 0) {
        toast.success(
          t('pos.offlineSyncSuccess', `${result.processed} queued order(s) synced`)
        )
      }
      if (result.failed > 0) {
        toast.error(t('pos.offlineSyncFailed', 'Some queued orders could not be synced'))
      }
    })()
  }, [isOffline, queuedOrderCount, replayQueuedOrders, t])

  if (permissionLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-100 dark:bg-slate-950">
        <Loader2 className="h-7 w-7 animate-spin text-primary" />
      </div>
    )
  }

  if (!canUsePOS) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-100 px-4 dark:bg-slate-950">
        <div className="w-full max-w-md rounded-lg border border-red-300 bg-red-50 p-5 dark:border-red-500/40 dark:bg-red-500/10">
          <div className="flex items-start gap-3">
            <Lock className="mt-0.5 h-5 w-5 text-red-600 dark:text-red-300" />
            <div>
              <h1 className="text-base font-semibold text-red-700 dark:text-red-100">
                {t('pos.accessDenied', 'Access denied')}
              </h1>
              <p className="mt-1 text-sm text-red-700/90 dark:text-red-200">
                {t(
                  'pos.accessDeniedText',
                  'You do not have permission to use Cashier POS.'
                )}
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <>
      <POSTerminal
        locale={locale}
        t={t}
        loading={dataLoading}
        error={dataError}
        isOffline={isOffline}
        queuedOrderCount={queuedOrderCount}
        isReplayingQueue={isReplaying}
        categories={categories}
        visibleProducts={displayProducts}
        productCardSize={productCardSize}
        tables={tables}
        onAddProduct={handleAddProduct}
        onRetry={refresh}
        onSubmitOrder={handleSubmitOrder}
        onHoldOrder={handleHoldOrder}
        onRecallOrder={handleRecallOrder}
        onClearCart={handleClearCart}
        onChangeCartItemQuantity={handleChangeCartItemQuantity}
        onRemoveCartItem={handleRemoveCartItem}
        autoPrintEnabled={autoPrintEnabled}
        onToggleAutoPrint={toggleAutoPrint}
        isSubmitting={isSubmitting}
      />
      <POSReceipt data={latestReceipt} />
      {isDevMode && !focusModeEnabled && recentMetrics.length > 0 ? (
        <div className="pointer-events-none fixed bottom-3 right-3 z-50 w-56 rounded-lg border border-sky-200 bg-white/95 p-2 text-[11px] shadow-lg backdrop-blur dark:border-sky-700/60 dark:bg-slate-900/95">
          <p className="mb-1 font-semibold text-sky-700 dark:text-sky-300">
            {t('pos.perfBadge', 'POS Perf')}
          </p>
          <div className="space-y-1 text-gray-700 dark:text-gray-100">
            {recentMetrics.map((metric, index) => (
              <div
                key={`${metric.metric}-${metric.recordedAt}-${index}`}
                className="flex items-center justify-between gap-2"
              >
                <span className="truncate">{metric.metric.replace('pos.', '')}</span>
                <span className="tabular-nums">{metric.durationMs}ms</span>
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </>
  )
}
