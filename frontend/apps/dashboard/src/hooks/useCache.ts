/**
 * React Hook for Client-Side Caching
 *
 * Provides easy cache integration in React components
 */

import { useCallback, useEffect, useState } from 'react'
import { CacheManager, NamespacedCache, caches } from '@/lib/cache-manager'

/**
 * Hook to use a cached value
 */
export function useCachedValue<T>(
  key: string,
  compute: () => Promise<T> | T,
  options?: {
    ttl?: number
    cache?: CacheManager | NamespacedCache
    skip?: boolean
  }
) {
  const [value, setValue] = useState<T | undefined>(() => {
    if (options?.skip) return undefined
    return options?.cache?.get<T>(key)
  })

  const [isLoading, setIsLoading] = useState(!value && !options?.skip)
  const [error, setError] = useState<Error | null>(null)

  const cache = options?.cache || caches.computed

  useEffect(() => {
    if (options?.skip) return

    let mounted = true

    const loadValue = async () => {
      try {
        setIsLoading(true)
        setError(null)

        const result = await (cache as any).getOrCompute(
          key,
          compute,
          options?.ttl
        )

        if (mounted) {
          setValue(result)
        }
      } catch (err) {
        if (mounted) {
          const error = err instanceof Error ? err : new Error(String(err))
          setError(error)
        }
      } finally {
        if (mounted) {
          setIsLoading(false)
        }
      }
    }

    // Check if already in cache
    const cached = (cache as any).get(key) as T | undefined
    if (cached !== undefined) {
      setValue(cached)
      setIsLoading(false)
      return
    }

    loadValue()

    // Subscribe to cache updates
    const unsubscribe = (cache as any).subscribe(key, (newValue: T) => {
      if (mounted) {
        setValue(newValue)
      }
    })

    return () => {
      mounted = false
      unsubscribe()
    }
  }, [key, compute, options?.cache, options?.skip, options?.ttl, cache])

  const recompute = useCallback(async () => {
    try {
      setIsLoading(true)
      setError(null)

      // Delete from cache to force recompute
      ;(cache as any).delete(key)

      const result = await (cache as any).getOrCompute(
        key,
        compute,
        options?.ttl
      )

      setValue(result)
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
    } finally {
      setIsLoading(false)
    }
  }, [key, compute, options?.cache, options?.ttl, cache])

  return {
    value,
    isLoading,
    error,
    recompute,
  }
}

/**
 * Hook to manually manage cache
 */
export function useCache(
  cacheInstance?: CacheManager | NamespacedCache
) {
  const cache = cacheInstance || caches.computed

  const set = useCallback(
    <T,>(key: string, value: T, ttl?: number) => {
      return (cache as any).set(key, value, ttl)
    },
    [cache]
  )

  const get = useCallback(
    <T = any,>(key: string) => {
      return (cache as any).get(key) as T | undefined
    },
    [cache]
  )

  const has = useCallback(
    (key: string) => {
      return (cache as any).has(key)
    },
    [cache]
  )

  const remove = useCallback(
    (key: string) => {
      return (cache as any).delete(key)
    },
    [cache]
  )

  const clear = useCallback(() => {
    ;(cache as any).clear()
  }, [cache])

  return {
    set,
    get,
    has,
    remove,
    clear,
  }
}

/**
 * Hook to monitor cache statistics
 */
export function useCacheStats(
  cacheInstance?: CacheManager | NamespacedCache,
  updateInterval: number = 5000
) {
  const [stats, setStats] = useState(() => {
    return (cacheInstance || caches.computed).getStats?.() || {
      totalEntries: 0,
      totalSize: 0,
      hitRate: 0,
      missRate: 0,
      hits: 0,
      misses: 0,
    }
  })

  const cache = cacheInstance || caches.computed

  useEffect(() => {
    const timer = setInterval(() => {
      const newStats = (cache as any).getStats?.() || stats
      setStats(newStats)
    }, updateInterval)

    return () => clearInterval(timer)
  }, [cache, updateInterval])

  return stats
}

/**
 * Hook for API caching with React Query integration pattern
 */
export function useCachedAPI<T>(
  url: string,
  options?: {
    ttl?: number
    staleTime?: number
    skipCache?: boolean
  }
) {
  const [data, setData] = useState<T | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const { ttl = 5 * 60 * 1000, skipCache } = options || {}

  useEffect(() => {
    let mounted = true

    const fetchData = async () => {
      try {
        // Check cache first
        if (!skipCache) {
          const cached = caches.api.get<T>(url)
          if (cached !== undefined) {
            setData(cached)
            setIsLoading(false)
            return
          }
        }

        setIsLoading(true)
        const response = await fetch(url)

        if (!response.ok) {
          throw new Error(`Failed to fetch: ${response.statusText}`)
        }

        const result = (await response.json()) as T

        if (mounted) {
          // Cache the result
          caches.api.set(url, result, ttl)
          setData(result)
          setError(null)
        }
      } catch (err) {
        if (mounted) {
          const error = err instanceof Error ? err : new Error(String(err))
          setError(error)

          // Try to get stale data from cache
          const stale = caches.api.get<T>(url)
          if (stale !== undefined) {
            setData(stale)
          }
        }
      } finally {
        if (mounted) {
          setIsLoading(false)
        }
      }
    }

    fetchData()

    return () => {
      mounted = false
    }
  }, [url, ttl, skipCache])

  const refetch = useCallback(async () => {
    try {
      setIsLoading(true)
      caches.api.delete(url)

      const response = await fetch(url)

      if (!response.ok) {
        throw new Error(`Failed to fetch: ${response.statusText}`)
      }

      const result = (await response.json()) as T

      caches.api.set(url, result, ttl)
      setData(result)
      setError(null)
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err))
      setError(error)
    } finally {
      setIsLoading(false)
    }
  }, [url, ttl])

  return { data, isLoading, error, refetch }
}
