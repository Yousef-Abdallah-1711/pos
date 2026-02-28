'use client'

import { useCallback, useRef, useState } from 'react'
import { recordPerformanceMetric } from '@/lib/performance'

export type POSPerformanceMetric =
  | 'pos.page-load'
  | 'pos.add-to-cart'
  | 'pos.search'
  | 'pos.order-submit'
  | 'pos.product-fetch'

export interface POSPerformanceSample {
  metric: POSPerformanceMetric
  durationMs: number
  recordedAt: number
}

declare global {
  interface Window {
    __posMetrics?: POSPerformanceSample[]
  }
}

const MAX_RECENT_METRICS = 5
const MAX_GLOBAL_METRICS = 100

const now = () => {
  if (typeof window !== 'undefined' && window.performance) {
    return window.performance.now()
  }
  return Date.now()
}

const appendGlobalMetric = (sample: POSPerformanceSample) => {
  if (typeof window === 'undefined') {
    return
  }

  const existing = Array.isArray(window.__posMetrics) ? window.__posMetrics : []
  window.__posMetrics = [sample, ...existing].slice(0, MAX_GLOBAL_METRICS)
}

export const usePOSPerformance = () => {
  const pageStartRef = useRef(now())
  const pageLoadRecordedRef = useRef(false)
  const [recentMetrics, setRecentMetrics] = useState<POSPerformanceSample[]>([])
  const isDevMode = process.env.NODE_ENV !== 'production'

  const recordMetric = useCallback(
    (metric: POSPerformanceMetric, durationMs: number) => {
      if (!Number.isFinite(durationMs) || durationMs < 0) {
        return
      }

      const rounded = Math.round(durationMs)
      recordPerformanceMetric(metric, rounded)

      const sample: POSPerformanceSample = {
        metric,
        durationMs: rounded,
        recordedAt: Date.now(),
      }

      appendGlobalMetric(sample)

      if (isDevMode) {
        setRecentMetrics((current) => [sample, ...current].slice(0, MAX_RECENT_METRICS))
      }
    },
    [isDevMode]
  )

  const measureSync = useCallback(
    <T,>(metric: POSPerformanceMetric, callback: () => T): T => {
      const startedAt = now()
      const result = callback()
      recordMetric(metric, now() - startedAt)
      return result
    },
    [recordMetric]
  )

  const measureAsync = useCallback(
    async <T,>(metric: POSPerformanceMetric, callback: () => Promise<T>): Promise<T> => {
      const startedAt = now()
      try {
        return await callback()
      } finally {
        recordMetric(metric, now() - startedAt)
      }
    },
    [recordMetric]
  )

  const markPageLoad = useCallback(() => {
    if (pageLoadRecordedRef.current) {
      return
    }
    pageLoadRecordedRef.current = true
    recordMetric('pos.page-load', now() - pageStartRef.current)
  }, [recordMetric])

  return {
    isDevMode,
    recentMetrics,
    recordMetric,
    measureSync,
    measureAsync,
    markPageLoad,
  }
}

export default usePOSPerformance
