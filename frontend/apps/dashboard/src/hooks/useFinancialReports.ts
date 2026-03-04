'use client'

import { DependencyList, useCallback, useEffect, useRef, useState } from 'react'

type CacheEntry<T> = {
  data: T
  timestamp: number
}

const reportCache = new Map<string, CacheEntry<unknown>>()

interface UseFinancialReportOptions<T> {
  cacheKey: string
  deps: DependencyList
  enabled?: boolean
  ttlMs?: number
  fetcher: () => Promise<T>
}

export function useFinancialReport<T>({
  cacheKey,
  deps,
  enabled = true,
  ttlMs = 5 * 60 * 1000,
  fetcher,
}: UseFinancialReportOptions<T>) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState<boolean>(enabled)
  const [error, setError] = useState<string | null>(null)
  const fetcherRef = useRef(fetcher)

  useEffect(() => {
    fetcherRef.current = fetcher
  }, [fetcher])

  const load = useCallback(
    async (force = false) => {
      if (!enabled) {
        setLoading(false)
        return
      }

      const cached = reportCache.get(cacheKey) as CacheEntry<T> | undefined
      if (!force && cached && Date.now()-cached.timestamp <= ttlMs) {
        setData(cached.data)
        setLoading(false)
        setError(null)
        return
      }

      setLoading(true)
      setError(null)
      try {
        const next = await fetcherRef.current()
        reportCache.set(cacheKey, { data: next, timestamp: Date.now() })
        setData(next)
      } catch (err: any) {
        setError(err?.response?.data?.error || err?.message || 'Failed to load report')
      } finally {
        setLoading(false)
      }
    },
    [cacheKey, enabled, ttlMs]
  )

  useEffect(() => {
    void load(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [load, ...deps])

  return {
    data,
    loading,
    error,
    refetch: () => load(true),
    setData,
  }
}

