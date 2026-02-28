'use client'

import { DependencyList, useCallback, useEffect, useRef, useState } from 'react'

interface UseReportsOptions<T> {
  enabled: boolean
  deps: DependencyList
  fetcher: () => Promise<T>
  pollIntervalMs?: number
  getErrorMessage: (error: unknown) => string
  onError?: (message: string, error: unknown, silent: boolean) => void
}

export function useReports<T>({
  enabled,
  deps,
  fetcher,
  pollIntervalMs = 30000,
  getErrorMessage,
  onError,
}: UseReportsOptions<T>) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [lastUpdatedAt, setLastUpdatedAt] = useState<Date | null>(null)
  const fetcherRef = useRef(fetcher)
  const getErrorMessageRef = useRef(getErrorMessage)
  const onErrorRef = useRef(onError)

  useEffect(() => {
    fetcherRef.current = fetcher
  }, [fetcher])

  useEffect(() => {
    getErrorMessageRef.current = getErrorMessage
  }, [getErrorMessage])

  useEffect(() => {
    onErrorRef.current = onError
  }, [onError])

  const refetch = useCallback(
    async (silent = false) => {
      if (!enabled) {
        setLoading(false)
        return
      }

      if (!silent) setLoading(true)
      if (silent) setRefreshing(true)
      setError(null)

      try {
        const nextData = await fetcherRef.current()
        setData(nextData)
        setLastUpdatedAt(new Date())
      } catch (err) {
        const message = getErrorMessageRef.current(err)
        setError(message)
        onErrorRef.current?.(message, err, silent)
      } finally {
        setLoading(false)
        setRefreshing(false)
      }
    },
    [enabled]
  )

  useEffect(() => {
    if (!enabled) {
      setLoading(false)
      return
    }
    void refetch(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, refetch, ...deps])

  useEffect(() => {
    if (!enabled || pollIntervalMs <= 0) return
    const id = window.setInterval(() => {
      void refetch(true)
    }, pollIntervalMs)
    return () => window.clearInterval(id)
  }, [enabled, pollIntervalMs, refetch])

  return {
    data,
    loading,
    refreshing,
    error,
    lastUpdatedAt,
    refetch,
    setData,
  }
}
