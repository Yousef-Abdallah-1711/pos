'use client'

import type { POSProduct } from '@/lib/api/posApi'
import { Box } from 'lucide-react'
import { POSProductCard } from './POSProductCard'
import { POSVirtualGrid } from './POSVirtualGrid'
import type { POSProductCardSize } from '@/stores/posStore'

interface POSProductGridProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  products: POSProduct[]
  productCardSize?: POSProductCardSize
  onAddProduct: (product: POSProduct) => void
}

export function POSProductGrid({
  locale,
  t,
  products,
  productCardSize = 'small',
  onAddProduct,
}: POSProductGridProps) {
  if (products.length === 0) {
    return (
      <div className="flex h-full min-h-[260px] flex-col items-center justify-center rounded-xl border border-dashed border-gray-300 bg-white text-gray-500 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300">
        <Box className="mb-2 h-6 w-6" />
        <p className="text-sm font-medium">{t('pos.noProducts', 'No products found')}</p>
      </div>
    )
  }

  if (products.length > 50) {
    return (
      <POSVirtualGrid
        locale={locale}
        t={t}
        products={products}
        productCardSize={productCardSize}
        onAddProduct={onAddProduct}
      />
    )
  }

  return (
    <div className="h-full overflow-auto rounded-xl border border-gray-200 bg-gray-50 p-2 dark:border-gray-800 dark:bg-gray-950">
      <div
        className="grid gap-2 [contain:content]"
        style={{
          gridTemplateColumns:
            productCardSize === 'large'
              ? 'repeat(auto-fill,minmax(190px,1fr))'
              : 'repeat(auto-fill,minmax(130px,1fr))',
        }}
      >
        {products.map((product, index) => (
          <POSProductCard
            key={product.id}
            locale={locale}
            t={t}
            product={product}
            size={productCardSize}
            keyboardIndex={index < 12 ? index : null}
            onAddProduct={onAddProduct}
          />
        ))}
      </div>
    </div>
  )
}

export default POSProductGrid
