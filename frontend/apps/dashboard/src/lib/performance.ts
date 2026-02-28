'use client'

import { debugLog } from '@/lib/debug'

const STORAGE_KEY = 'dashboard-performance-metrics-v1'
const MAX_SAMPLES_PER_METRIC = 120

type MetricSamples = Record<string, number[]>

const isBrowser = () => typeof window !== 'undefined'

const readMetrics = (): MetricSamples => {
  if (!isBrowser()) {
    return {}
  }

  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return {}
    }

    const parsed = JSON.parse(raw) as MetricSamples
    if (!parsed || typeof parsed !== 'object') {
      return {}
    }
    return parsed
  } catch {
    return {}
  }
}

const writeMetrics = (metrics: MetricSamples) => {
  if (!isBrowser()) {
    return
  }

  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(metrics))
  } catch {
    // Ignore localStorage errors.
  }
}

const percentile = (values: number[], p: number): number => {
  if (values.length === 0) {
    return 0
  }

  const sorted = [...values].sort((a, b) => a - b)
  const index = Math.min(sorted.length - 1, Math.max(0, Math.ceil((p / 100) * sorted.length) - 1))
  return sorted[index]
}

export const recordPerformanceMetric = (metric: string, valueMs: number) => {
  if (!isBrowser() || !Number.isFinite(valueMs) || valueMs < 0) {
    return
  }

  const metrics = readMetrics()
  const current = metrics[metric] ?? []
  const next = [...current, Math.round(valueMs)]
  if (next.length > MAX_SAMPLES_PER_METRIC) {
    next.splice(0, next.length - MAX_SAMPLES_PER_METRIC)
  }

  metrics[metric] = next
  writeMetrics(metrics)

  const p95 = percentile(next, 95)
  debugLog(`[Perf] ${metric}: last=${Math.round(valueMs)}ms p95=${Math.round(p95)}ms samples=${next.length}`)
}

export const startPerformanceTimer = (metric: string) => {
  if (!isBrowser()) {
    return () => undefined
  }

  const startedAt = window.performance.now()
  return () => {
    const duration = window.performance.now() - startedAt
    recordPerformanceMetric(metric, duration)
  }
}

export const getPerformanceSummary = (metric: string) => {
  const values = readMetrics()[metric] ?? []
  if (values.length === 0) {
    return {
      count: 0,
      avg: 0,
      p95: 0,
      max: 0,
    }
  }

  const total = values.reduce((sum, value) => sum + value, 0)
  return {
    count: values.length,
    avg: total / values.length,
    p95: percentile(values, 95),
    max: Math.max(...values),
  }
}
