'use client'

import { useCallback, useMemo } from 'react'
import type { POSProduct } from '@/lib/api/posApi'

interface IndexedPOSProduct extends POSProduct {
  _lowerName: string
}

interface POSProductIndex {
  indexedProducts: IndexedPOSProduct[]
  search: (query: string) => POSProduct[]
  getByBarcode: (barcode: string) => POSProduct | undefined
}

const normalizeText = (value: string) => value.trim().toLowerCase()

export const usePOSProductIndex = (products: POSProduct[]): POSProductIndex => {
  const indexedProducts = useMemo<IndexedPOSProduct[]>(
    () =>
      products.map((product) => ({
        ...product,
        _lowerName: normalizeText(product.name || ''),
      })),
    [products]
  )

  const barcodeMap = useMemo(() => {
    const nextMap = new Map<string, POSProduct>()
    for (const product of products) {
      if (typeof product.barcode === 'string' && product.barcode.trim()) {
        nextMap.set(product.barcode.trim(), product)
      }
    }
    return nextMap
  }, [products])

  const search = useCallback(
    (query: string): POSProduct[] => {
      const normalizedQuery = normalizeText(query)
      if (!normalizedQuery) {
        return products
      }

      return indexedProducts
        .filter((product) => product._lowerName.includes(normalizedQuery))
        .map(({ _lowerName: _ignore, ...product }) => product)
    },
    [indexedProducts, products]
  )

  const getByBarcode = useCallback(
    (barcode: string) => barcodeMap.get(barcode.trim()),
    [barcodeMap]
  )

  return {
    indexedProducts,
    search,
    getByBarcode,
  }
}

export default usePOSProductIndex
