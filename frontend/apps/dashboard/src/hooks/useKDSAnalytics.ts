'use client'

import { useEffect, useRef, useState } from 'react'
import { AxiosError } from 'axios'
import { api } from '@/lib/api'

export interface KDSPrepTrendPoint {
  date: string
  avg_prep_minutes: number
  completed_orders: number
}

export interface KDSHourLoadPoint {
  hour: number
  order_count: number
}

export interface KDSAnalytics {
  average_prep_minutes: number
  total_completed_orders: number
  prep_time_trend: KDSPrepTrendPoint[]
  busiest_hours: KDSHourLoadPoint[]
}

interface ApiResponse<T> {
  success?: boolean
  data?: T
  error?: string
  message?: string
}

const DEFAULT_ANALYTICS: KDSAnalytics = {
  average_prep_minutes: 0,
  total_completed_orders: 0,
  prep_time_trend: [],
  busiest_hours: [],
}

const RETRY_ATTEMPTS = 3
const RETRY_BASE_DELAY_MS = 300

const extractData = <T,>(payload: unknown, fallback: T): T => {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return ((payload as ApiResponse<T>).data ?? fallback) as T
  }
  return (payload as T) ?? fallback
}

const getErrorMessage = (error: unknown, fallback: string) => {
  if (!error || typeof error !== 'object') {
    return fallback
  }

  const axiosError = error as AxiosError<ApiResponse<unknown>>
  const responseError =
    axiosError.response?.data?.error ||
    axiosError.response?.data?.message

  return responseError || axiosError.message || fallback
}

const isTransientError = (error: unknown) => {
  if (!error || typeof error !== 'object') {
    return false
  }

  const axiosError = error as AxiosError<ApiResponse<unknown>>
  const statusCode = axiosError.response?.status
  if (typeof statusCode === 'number') {
    return statusCode >= 500 || statusCode === 429 || statusCode === 408
  }

  return true
}

const sleep = (ms: number) =>
  new Promise<void>((resolve) => {
    setTimeout(resolve, ms)
  })

const withRetry = async <T,>(operation: () => Promise<T>, attempts = RETRY_ATTEMPTS): Promise<T> => {
  let lastError: unknown

  for (let attempt = 0; attempt < attempts; attempt += 1) {
    try {
      return await operation()
    } catch (error) {
      lastError = error
      if (attempt === attempts - 1 || !isTransientError(error)) {
        throw error
      }

      await sleep(RETRY_BASE_DELAY_MS * Math.pow(2, attempt))
    }
  }

  throw lastError
}

export const useKDSAnalytics = (enabled: boolean) => {
  const [analytics, setAnalytics] = useState<KDSAnalytics>(DEFAULT_ANALYTICS)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [available, setAvailable] = useState(true)
  const mountedRef = useRef(true)

  useEffect(() => {
    return () => {
      mountedRef.current = false
    }
  }, [])

  const fetchAnalytics = async () => {
    if (!enabled || !available) {
      return null
    }

    try {
      setLoading(true)
      const response = await withRetry(() => api.get('/kds/analytics'))
      const payload = extractData<KDSAnalytics>(response.data, DEFAULT_ANALYTICS)
      const normalized: KDSAnalytics = {
        ...DEFAULT_ANALYTICS,
        ...payload,
        prep_time_trend: Array.isArray(payload.prep_time_trend) ? payload.prep_time_trend : [],
        busiest_hours: Array.isArray(payload.busiest_hours) ? payload.busiest_hours : [],
      }

      if (!mountedRef.current) {
        return null
      }

      setAnalytics(normalized)
      setAvailable(true)
      setError(null)
      return normalized
    } catch (err) {
      if (!mountedRef.current) {
        return null
      }
      const axiosError = err as AxiosError<ApiResponse<unknown>>
      if (axiosError.response?.status === 404) {
        // Support older backend instances that do not expose analytics yet.
        setAnalytics(DEFAULT_ANALYTICS)
        setAvailable(false)
        setError(null)
        return null
      }
      setError(getErrorMessage(err, 'Failed to load KDS analytics'))
      return null
    } finally {
      if (mountedRef.current) {
        setLoading(false)
      }
    }
  }

  useEffect(() => {
    if (!enabled) {
      setAnalytics(DEFAULT_ANALYTICS)
      setError(null)
      setAvailable(true)
      setLoading(false)
      return
    }

    void fetchAnalytics()
    // Only re-run when enablement changes; manual refresh uses returned function.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled])

  return {
    analytics,
    loading,
    error,
    analyticsAvailable: available,
    refreshAnalytics: fetchAnalytics,
  }
}
