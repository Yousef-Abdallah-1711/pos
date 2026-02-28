'use client'

import { memo } from 'react'
import type { POSProduct } from '@/lib/api/posApi'
import { selectPOSItemByProductId, usePOSStore } from '@/stores/posStore'
import type { POSProductCardSize } from '@/stores/posStore'

interface POSProductCardProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  product: POSProduct
  size?: POSProductCardSize
  keyboardIndex: number | null
  onAddProduct: (product: POSProduct) => void
}

const formatPrice = (amount: number, locale: 'en' | 'ar') =>
  new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount)

function POSProductCardComponent({
  locale,
  t,
  product,
  size = 'small',
  keyboardIndex,
  onAddProduct,
}: POSProductCardProps) {
  const cartItem = usePOSStore(selectPOSItemByProductId(product.id))
  const quantityInCart = cartItem?.quantity || 0
  const isOutOfStock =
    product.is_available === false ||
    (product.track_inventory === true && (product.quantity_in_stock || 0) <= 0)

  const displayName =
    locale === 'ar' && product.name_ar ? product.name_ar : product.name
  const isLarge = size === 'large'

  return (
    <button
      type="button"
      disabled={isOutOfStock}
      onClick={() => onAddProduct(product)}
      className={`relative flex flex-col items-start justify-between rounded-xl border text-left transition active:scale-95 ${
        isLarge ? 'min-h-[190px] p-4' : 'min-h-[140px] p-3'
      } ${
        isOutOfStock
          ? 'cursor-not-allowed border-gray-200 bg-gray-100 text-gray-400 dark:border-gray-800 dark:bg-gray-900 dark:text-gray-500'
          : 'border-gray-200 bg-white text-gray-900 shadow-sm hover:border-primary/50 hover:shadow-md dark:border-gray-800 dark:bg-gray-950 dark:text-gray-100'
      }`}
    >
      {keyboardIndex !== null && (
        <span className="absolute left-2 top-2 rounded bg-black/80 px-1.5 py-0.5 text-[10px] font-semibold text-white">
          F{keyboardIndex + 1}
        </span>
      )}

      {quantityInCart > 0 && (
        <span className="absolute right-2 top-2 rounded-full bg-primary px-2 py-0.5 text-[10px] font-semibold text-white">
          {quantityInCart}
        </span>
      )}

      <div className="w-full">
        {product.image ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={product.image}
            alt={displayName}
            loading="lazy"
            className={`mb-2 w-full rounded-md object-cover ${
              isLarge ? 'h-20' : 'h-14'
            }`}
          />
        ) : (
          <div
            className={`mb-2 flex w-full items-center justify-center rounded-md bg-gray-100 font-semibold text-gray-500 dark:bg-gray-800 dark:text-gray-300 ${
              isLarge ? 'h-20 text-2xl' : 'h-14 text-xl'
            }`}
          >
            {displayName.charAt(0)}
          </div>
        )}
        <p className={`line-clamp-2 font-semibold ${isLarge ? 'text-base' : 'text-sm'}`}>
          {displayName}
        </p>
      </div>

      <div className="mt-2 flex w-full items-center justify-between">
        <span className={`${isLarge ? 'text-base' : 'text-sm'} font-bold text-primary`}>
          {formatPrice(product.price, locale)}
        </span>
        {isOutOfStock && (
          <span className="rounded bg-red-100 px-1.5 py-0.5 text-[10px] font-semibold text-red-700 dark:bg-red-900/40 dark:text-red-300">
            {t('pos.outOfStock', 'Out of stock')}
          </span>
        )}
      </div>
    </button>
  )
}

export const POSProductCard = memo(POSProductCardComponent)
POSProductCard.displayName = 'POSProductCard'

export default POSProductCard
