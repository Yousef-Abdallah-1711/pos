'use client'

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import type { POSProduct } from '@/lib/api/posApi'
import { POSProductCard } from './POSProductCard'
import type { POSProductCardSize } from '@/stores/posStore'

interface POSVirtualGridProps {
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  products: POSProduct[]
  productCardSize?: POSProductCardSize
  onAddProduct: (product: POSProduct) => void
}

const SMALL_CARD_MIN_WIDTH = 130
const SMALL_CARD_MIN_HEIGHT = 140
const LARGE_CARD_MIN_WIDTH = 190
const LARGE_CARD_MIN_HEIGHT = 190
const GRID_GAP = 8
const ROW_OVERSCAN = 2

export function POSVirtualGrid({
  locale,
  t,
  products,
  productCardSize = 'small',
  onAddProduct,
}: POSVirtualGridProps) {
  const scrollRef = useRef<HTMLDivElement | null>(null)
  const [columns, setColumns] = useState(1)
  const cardMinWidth =
    productCardSize === 'large' ? LARGE_CARD_MIN_WIDTH : SMALL_CARD_MIN_WIDTH
  const rowEstimatedHeight =
    (productCardSize === 'large' ? LARGE_CARD_MIN_HEIGHT : SMALL_CARD_MIN_HEIGHT) + GRID_GAP

  const recalculateColumns = useCallback(() => {
    const width = scrollRef.current?.clientWidth ?? 0
    if (width <= 0) {
      return
    }

    const nextColumns = Math.max(
      1,
      Math.floor((width + GRID_GAP) / (cardMinWidth + GRID_GAP))
    )
    setColumns((current) => (current === nextColumns ? current : nextColumns))
  }, [cardMinWidth])

  useEffect(() => {
    recalculateColumns()

    const container = scrollRef.current
    if (!container) {
      return
    }

    if (typeof ResizeObserver !== 'undefined') {
      const observer = new ResizeObserver(() => {
        recalculateColumns()
      })
      observer.observe(container)
      return () => {
        observer.disconnect()
      }
    }

    const onResize = () => {
      recalculateColumns()
    }
    window.addEventListener('resize', onResize)
    return () => {
      window.removeEventListener('resize', onResize)
    }
  }, [recalculateColumns])

  const rowCount = useMemo(
    () => Math.ceil(products.length / columns),
    [products.length, columns]
  )

  const rowVirtualizer = useVirtualizer({
    count: rowCount,
    getScrollElement: () => scrollRef.current,
    estimateSize: () => rowEstimatedHeight,
    overscan: ROW_OVERSCAN,
  })

  const virtualRows = rowVirtualizer.getVirtualItems()

  return (
    <div
      ref={scrollRef}
      className="h-full overflow-auto rounded-xl border border-gray-200 bg-gray-50 p-2 dark:border-gray-800 dark:bg-gray-950"
    >
      <div
        className="relative w-full"
        style={{ height: `${rowVirtualizer.getTotalSize()}px` }}
      >
        {virtualRows.map((virtualRow) => {
          const rowStart = virtualRow.index * columns
          const rowProducts = products.slice(rowStart, rowStart + columns)
          const emptySlots = Math.max(0, columns - rowProducts.length)

          return (
            <div
              key={virtualRow.key}
              ref={rowVirtualizer.measureElement}
              className="absolute left-0 top-0 grid gap-2 [contain:content]"
              style={{
                width: '100%',
                transform: `translateY(${virtualRow.start}px)`,
                gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))`,
              }}
            >
              {rowProducts.map((product, offset) => {
                const globalIndex = rowStart + offset
                return (
                  <POSProductCard
                    key={product.id}
                    locale={locale}
                    t={t}
                    product={product}
                    size={productCardSize}
                    keyboardIndex={globalIndex < 12 ? globalIndex : null}
                    onAddProduct={onAddProduct}
                  />
                )
              })}

              {Array.from({ length: emptySlots }).map((_, slotIndex) => (
                <div
                  key={`empty-${virtualRow.index}-${slotIndex}`}
                  aria-hidden="true"
                  className={`rounded-xl border border-transparent ${
                    productCardSize === 'large' ? 'min-h-[190px]' : 'min-h-[140px]'
                  }`}
                />
              ))}
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default POSVirtualGrid
