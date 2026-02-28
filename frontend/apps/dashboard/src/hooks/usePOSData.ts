'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { type POSCategory, posApi, type POSProduct } from '@/lib/api/posApi'
import { type RestaurantTable, tablesApi } from '@/lib/api/tablesApi'
import { debugError, debugWarn } from '@/lib/debug'
import { usePOSProductCache } from '@/hooks/usePOSProductCache'

interface POSDataState {
  products: POSProduct[]
  categories: POSCategory[]
  tables: RestaurantTable[]
  loading: boolean
  error: string | null
}

interface CachedPOSCatalog<T> {
  data: T[]
  savedAt: number
}

const PRODUCTS_CACHE_KEY = 'pos-products-cache'
const CATEGORIES_CACHE_KEY = 'pos-categories-cache'

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const readCache = <T>(key: string): T[] => {
  if (typeof window === 'undefined') {
    return []
  }

  try {
    const raw = window.localStorage.getItem(key)
    if (!raw) {
      return []
    }

    const parsed = JSON.parse(raw) as CachedPOSCatalog<T> | T[]
    if (Array.isArray(parsed)) {
      return parsed
    }

    if (isRecord(parsed) && Array.isArray(parsed.data)) {
      return parsed.data
    }

    return []
  } catch {
    return []
  }
}

const writeCache = <T>(key: string, data: T[]) => {
  if (typeof window === 'undefined') {
    return
  }

  try {
    const payload: CachedPOSCatalog<T> = {
      data,
      savedAt: Date.now(),
    }
    window.localStorage.setItem(key, JSON.stringify(payload))
  } catch {
    // Ignore storage failures.
  }
}

const hasData = <T>(list: T[] | null | undefined): list is T[] =>
  Array.isArray(list) && list.length > 0

const isHttpStatus = (error: unknown, status: number): boolean =>
  Boolean(
    error &&
      typeof error === 'object' &&
      'response' in error &&
      (error as { response?: { status?: number } }).response?.status === status
  )

export const usePOSData = () => {
  const hydratedRef = useRef(false)
  const fetchingRef = useRef(false)
  const memoryProductsRef = useRef<POSProduct[] | null>(null)
  const memoryCategoriesRef = useRef<POSCategory[] | null>(null)
  const memoryTablesRef = useRef<RestaurantTable[] | null>(null)

  const {
    products: indexedProducts,
    categories: indexedCategories,
    isStale: indexedCacheStale,
    refresh: refreshIndexedCache,
    writeProducts: writeIndexedProducts,
    writeCategories: writeIndexedCategories,
  } = usePOSProductCache()

  const [state, setState] = useState<POSDataState>({
    products: [],
    categories: [],
    tables: [],
    loading: true,
    error: null,
  })

  const fetchFreshData = useCallback(async () => {
    if (fetchingRef.current) {
      return
    }

    fetchingRef.current = true
    setState((current) => ({
      ...current,
      loading:
        current.products.length === 0 ||
        current.categories.length === 0 ||
        current.tables.length === 0,
      error: null,
    }))

    try {
      const [productsResult, categoriesResult, tablesResult] = await Promise.allSettled([
        posApi.getProducts(),
        posApi.getCategories(),
        tablesApi.list().catch((error: unknown) => {
          // Some environments do not expose the table-management endpoints yet.
          // POS should remain usable with takeaway mode when tables are unavailable.
          if (isHttpStatus(error, 404)) {
            debugWarn('[usePOSData] Tables API returned 404. Continuing with empty table list.')
            return [] as RestaurantTable[]
          }
          throw error
        }),
      ])

      setState((current) => {
        const nextProducts =
          productsResult.status === 'fulfilled' ? productsResult.value : current.products
        const nextCategories =
          categoriesResult.status === 'fulfilled' ? categoriesResult.value : current.categories
        const nextTables =
          tablesResult.status === 'fulfilled' ? tablesResult.value : current.tables

        memoryProductsRef.current = nextProducts
        memoryCategoriesRef.current = nextCategories
        memoryTablesRef.current = nextTables

        if (productsResult.status === 'fulfilled') {
          writeCache(PRODUCTS_CACHE_KEY, productsResult.value)
          void writeIndexedProducts(productsResult.value)
        }

        if (categoriesResult.status === 'fulfilled') {
          writeCache(CATEGORIES_CACHE_KEY, categoriesResult.value)
          void writeIndexedCategories(categoriesResult.value)
        }

        const errors: string[] = []
        if (productsResult.status === 'rejected') {
          debugError('[usePOSData] Failed to fetch POS products:', productsResult.reason)
          errors.push('products')
        }
        if (categoriesResult.status === 'rejected') {
          debugError('[usePOSData] Failed to fetch POS categories:', categoriesResult.reason)
          errors.push('categories')
        }
        if (tablesResult.status === 'rejected') {
          debugError('[usePOSData] Failed to fetch POS tables:', tablesResult.reason)
          errors.push('tables')
        }

        const catalogRequestFailed =
          productsResult.status === 'rejected' || categoriesResult.status === 'rejected'

        return {
          products: nextProducts,
          categories: nextCategories,
          tables: nextTables,
          loading: false,
          error:
            catalogRequestFailed && nextProducts.length === 0 && nextCategories.length === 0
              ? 'Unable to load POS data right now.'
              : null,
        }
      })
    } finally {
      fetchingRef.current = false
    }
  }, [writeIndexedCategories, writeIndexedProducts])

  useEffect(() => {
    if (hydratedRef.current) {
      return
    }
    hydratedRef.current = true

    const inMemoryProducts = memoryProductsRef.current
    const inMemoryCategories = memoryCategoriesRef.current
    const inMemoryTables = memoryTablesRef.current
    const localProducts = readCache<POSProduct>(PRODUCTS_CACHE_KEY)
    const localCategories = readCache<POSCategory>(CATEGORIES_CACHE_KEY)

    if (hasData(inMemoryProducts) || hasData(inMemoryCategories)) {
      setState((current) => ({
        ...current,
        products: hasData(inMemoryProducts) ? inMemoryProducts : current.products,
        categories: hasData(inMemoryCategories) ? inMemoryCategories : current.categories,
        tables: hasData(inMemoryTables) ? inMemoryTables : current.tables,
        loading: false,
      }))
    } else if (hasData(localProducts) || hasData(localCategories)) {
      setState((current) => ({
        ...current,
        products: hasData(localProducts) ? localProducts : current.products,
        categories: hasData(localCategories) ? localCategories : current.categories,
        loading: false,
      }))
    }

    void refreshIndexedCache()
    void fetchFreshData()
  }, [fetchFreshData, refreshIndexedCache])

  useEffect(() => {
    // Phase 2 priority: IndexedDB cache should hydrate before API finishes.
    if (!hasData(indexedProducts) && !hasData(indexedCategories)) {
      return
    }

    setState((current) => ({
      ...current,
      products:
        current.products.length === 0 && hasData(indexedProducts)
          ? indexedProducts
          : current.products,
      categories:
        current.categories.length === 0 && hasData(indexedCategories)
          ? indexedCategories
          : current.categories,
      loading: false,
      error:
        indexedCacheStale && current.error
          ? 'Showing cached POS catalog. Live update is temporarily unavailable.'
          : current.error,
    }))
  }, [indexedCacheStale, indexedCategories, indexedProducts])

  return {
    products: state.products,
    categories: state.categories,
    tables: state.tables,
    loading: state.loading,
    error: state.error,
    refresh: fetchFreshData,
  }
}

export default usePOSData
