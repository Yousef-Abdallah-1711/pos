'use client'

import { useCallback, useEffect, useState } from 'react'
import type { POSCategory, POSProduct } from '@/lib/api/posApi'

type CachedRecord<T> = T & { __cachedAt: number }

interface POSProductCacheState {
  products: POSProduct[]
  categories: POSCategory[]
  isStale: boolean
  refresh: () => Promise<void>
  writeProducts: (products: POSProduct[]) => Promise<boolean>
  writeCategories: (categories: POSCategory[]) => Promise<boolean>
}

const DB_NAME = 'pos-product-db'
const DB_VERSION = 1
const PRODUCTS_STORE = 'products'
const CATEGORIES_STORE = 'categories'
const CACHE_TTL_MS = 5 * 60 * 1000

const canUseIndexedDB = () =>
  typeof window !== 'undefined' && typeof window.indexedDB !== 'undefined'

const openProductDB = (): Promise<IDBDatabase | null> => {
  if (!canUseIndexedDB()) {
    return Promise.resolve(null)
  }

  return new Promise((resolve) => {
    try {
      const request = window.indexedDB.open(DB_NAME, DB_VERSION)

      request.onupgradeneeded = () => {
        const db = request.result
        if (!db.objectStoreNames.contains(PRODUCTS_STORE)) {
          db.createObjectStore(PRODUCTS_STORE, { keyPath: 'id' })
        }
        if (!db.objectStoreNames.contains(CATEGORIES_STORE)) {
          db.createObjectStore(CATEGORIES_STORE, { keyPath: 'id' })
        }
      }

      request.onsuccess = () => resolve(request.result)
      request.onerror = () => resolve(null)
    } catch {
      resolve(null)
    }
  })
}

const readStore = async <T extends { id: number }>(
  storeName: string
): Promise<{ data: T[]; newestCachedAt: number | null }> => {
  const db = await openProductDB()
  if (!db) {
    return { data: [], newestCachedAt: null }
  }

  return new Promise((resolve) => {
    const tx = db.transaction(storeName, 'readonly')
    const store = tx.objectStore(storeName)
    const request = store.getAll()

    request.onsuccess = () => {
      const records = (request.result || []) as CachedRecord<T>[]
      const newestCachedAt =
        records.length > 0
          ? Math.max(...records.map((record) => Number(record.__cachedAt) || 0))
          : null

      resolve({
        data: records.map((record) => {
          const { __cachedAt: _ignore, ...item } = record
          return item as unknown as T
        }),
        newestCachedAt,
      })
    }

    request.onerror = () => resolve({ data: [], newestCachedAt: null })
    tx.oncomplete = () => db.close()
    tx.onerror = () => db.close()
  })
}

const writeStore = async <T extends { id: number }>(
  storeName: string,
  items: T[]
): Promise<boolean> => {
  const db = await openProductDB()
  if (!db) {
    return false
  }

  const cachedAt = Date.now()

  return new Promise<boolean>((resolve) => {
    const tx = db.transaction(storeName, 'readwrite')
    const store = tx.objectStore(storeName)
    const clearRequest = store.clear()

    clearRequest.onerror = () => resolve(false)
    clearRequest.onsuccess = () => {
      for (const item of items) {
        store.put({
          ...item,
          __cachedAt: cachedAt,
        } satisfies CachedRecord<T>)
      }
    }

    tx.oncomplete = () => resolve(true)
    tx.onerror = () => resolve(false)
  }).finally(() => {
    db.close()
  })
}

export const usePOSProductCache = (): POSProductCacheState => {
  const [products, setProducts] = useState<POSProduct[]>([])
  const [categories, setCategories] = useState<POSCategory[]>([])
  const [isStale, setIsStale] = useState(true)

  const refresh = useCallback(async () => {
    const [productsResult, categoriesResult] = await Promise.all([
      readStore<POSProduct>(PRODUCTS_STORE),
      readStore<POSCategory>(CATEGORIES_STORE),
    ])

    setProducts(productsResult.data)
    setCategories(categoriesResult.data)

    const newestCachedAt = Math.max(
      productsResult.newestCachedAt || 0,
      categoriesResult.newestCachedAt || 0
    )
    if (newestCachedAt <= 0) {
      setIsStale(true)
      return
    }

    setIsStale(Date.now() - newestCachedAt > CACHE_TTL_MS)
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const writeProducts = useCallback(async (nextProducts: POSProduct[]) => {
    const stored = await writeStore(PRODUCTS_STORE, nextProducts)
    if (stored) {
      setProducts(nextProducts)
      setIsStale(false)
    }
    return stored
  }, [])

  const writeCategories = useCallback(async (nextCategories: POSCategory[]) => {
    const stored = await writeStore(CATEGORIES_STORE, nextCategories)
    if (stored) {
      setCategories(nextCategories)
      setIsStale(false)
    }
    return stored
  }, [])

  return {
    products,
    categories,
    isStale,
    refresh,
    writeProducts,
    writeCategories,
  }
}

export default usePOSProductCache
