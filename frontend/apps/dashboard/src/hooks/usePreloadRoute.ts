/**
 * Hook for preloading routes on link hover
 *
 * Usage:
 * const { onMouseEnter, isLoading } = usePreloadRoute(() => import('@/pages/heavy-route'))
 * <Link onMouseEnter={onMouseEnter}>Navigate</Link>
 */

import { useCallback, useState } from 'react'
import { debugWarn } from '@/lib/debug'

interface UsePreloadRouteOptions {
  delay?: number
  onSuccess?: () => void
  onError?: (error: Error) => void
}

export function usePreloadRoute(
  importFn: () => Promise<any>,
  options?: UsePreloadRouteOptions
) {
  const [isLoading, setIsLoading] = useState(false)
  const [isPreloaded, setIsPreloaded] = useState(false)

  const preload = useCallback(async () => {
    if (isPreloaded || isLoading) return

    setIsLoading(true)

    try {
      // Add delay if specified
      if (options?.delay) {
        await new Promise((resolve) => setTimeout(resolve, options.delay))
      }

      await importFn()
      setIsPreloaded(true)
      options?.onSuccess?.()
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Failed to preload route')
      debugWarn('Route preload failed:', err)
      options?.onError?.(err)
    } finally {
      setIsLoading(false)
    }
  }, [importFn, isLoading, isPreloaded, options])

  return {
    preload,
    isLoading,
    isPreloaded,
    // Convenient handlers for common events
    onMouseEnter: preload,
    onFocus: preload,
    onTouchStart: preload,
  }
}

/**
 * Hook for preloading multiple routes in sequence
 */
export function usePreloadRoutes(
  importFns: Array<() => Promise<any>>,
  options?: UsePreloadRouteOptions & { concurrency?: number }
) {
  const [isLoading, setIsLoading] = useState(false)
  const [preloadedCount, setPreloadedCount] = useState(0)

  const preload = useCallback(async () => {
    if (isLoading) return

    setIsLoading(true)
    const { concurrency = 3 } = options || {}

    try {
      // Process routes with concurrency limit
      for (let i = 0; i < importFns.length; i += concurrency) {
        const batch = importFns.slice(i, i + concurrency)
        await Promise.all(batch.map((fn) => fn()))
        setPreloadedCount(i + batch.length)
      }

      options?.onSuccess?.()
    } catch (error) {
      const err = error instanceof Error ? error : new Error('Failed to preload routes')
      debugWarn('Routes preload failed:', err)
      options?.onError?.(err)
    } finally {
      setIsLoading(false)
    }
  }, [importFns, isLoading, options])

  return {
    preload,
    isLoading,
    preloadedCount,
    totalCount: importFns.length,
    progress: (preloadedCount / importFns.length) * 100,
  }
}
