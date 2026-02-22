/**
 * Performance Monitor Component
 *
 * Displays Web Vitals and performance metrics in real-time
 * Useful for development and debugging
 */

'use client'

import { useWebVitals, usePerformanceScore, useResourceMetrics } from '@/hooks/useWebVitals'
import { PERFORMANCE_THRESHOLDS } from '@/lib/web-vitals'
import { useEffect, useState } from 'react'

interface PerformanceMonitorProps {
  position?: 'top-right' | 'bottom-right' | 'bottom-left' | 'top-left'
  compact?: boolean
  showOnProduction?: boolean
}

/**
 * Main Performance Monitor Component
 */
export function PerformanceMonitor({
  position = 'bottom-right',
  compact = false,
  showOnProduction = false,
}: PerformanceMonitorProps) {
  // Only show in development unless explicitly enabled
  if (
    typeof window !== 'undefined' &&
    process.env.NODE_ENV === 'production' &&
    !showOnProduction
  ) {
    return null
  }

  const vitals = useWebVitals()
  const { score, rating } = usePerformanceScore()
  const resources = useResourceMetrics()
  const [isOpen, setIsOpen] = useState(!compact)

  const positionClass = {
    'top-right': 'top-4 right-4',
    'bottom-right': 'bottom-4 right-4',
    'bottom-left': 'bottom-4 left-4',
    'top-left': 'top-4 left-4',
  }[position]

  return (
    <div className={`performance-monitor fixed ${positionClass} z-50 font-mono text-xs`}>
      {compact && (
        <button
          onClick={() => setIsOpen(!isOpen)}
          className="bg-gray-900 text-white px-3 py-2 rounded-lg hover:bg-gray-800 transition"
        >
          {isOpen ? 'Hide' : 'Show'} Metrics
        </button>
      )}

      {isOpen && (
        <div className="mt-2 bg-gray-900 text-white p-4 rounded-lg shadow-xl max-w-md max-h-96 overflow-y-auto">
          {/* Score Card */}
          <div className="mb-4 pb-4 border-b border-gray-700">
            <div className="flex items-center justify-between mb-2">
              <span className="font-bold">Performance Score</span>
              <span className={`font-bold text-lg ${getScoreColor(score)}`}>
                {score}
              </span>
            </div>
            <div className="bg-gray-800 rounded-full h-2 w-full overflow-hidden">
              <div className={`performance-score-fill h-full transition-all ${getScoreColor(score)}`} />
            </div>
            <span className="text-xs text-gray-400 mt-1 block capitalize">
              {rating}
            </span>
          </div>

          {/* Core Web Vitals */}
          <div className="mb-4 pb-4 border-b border-gray-700">
            <h3 className="font-bold mb-2 text-sm">Core Web Vitals</h3>

            {/* LCP */}
            {vitals.lcp && (
              <MetricRow
                name="LCP"
                value={vitals.lcp.value}
                unit="ms"
                rating={vitals.lcp.rating}
                threshold={PERFORMANCE_THRESHOLDS.lcp}
              />
            )}

            {/* FID */}
            {vitals.fid && (
              <MetricRow
                name="FID"
                value={vitals.fid.value}
                unit="ms"
                rating={vitals.fid.rating}
                threshold={PERFORMANCE_THRESHOLDS.fid}
              />
            )}

            {/* INP */}
            {vitals.inp && (
              <MetricRow
                name="INP"
                value={vitals.inp.value}
                unit="ms"
                rating={vitals.inp.rating}
                threshold={{ good: 200, poor: 500 }}
              />
            )}

            {/* CLS */}
            {vitals.cls && (
              <MetricRow
                name="CLS"
                value={vitals.cls.value}
                unit=""
                rating={vitals.cls.rating}
                threshold={PERFORMANCE_THRESHOLDS.cls}
                decimals={3}
              />
            )}

            {/* TTFB */}
            {vitals.ttfb && (
              <MetricRow
                name="TTFB"
                value={vitals.ttfb.value}
                unit="ms"
                rating={vitals.ttfb.rating}
                threshold={PERFORMANCE_THRESHOLDS.ttfb}
              />
            )}

            {/* FCP */}
            {vitals.fcp && (
              <MetricRow
                name="FCP"
                value={vitals.fcp.value}
                unit="ms"
                rating={vitals.fcp.rating}
                threshold={PERFORMANCE_THRESHOLDS.fcp}
              />
            )}
          </div>

          {/* Resource Metrics */}
          <div className="mb-4 pb-4 border-b border-gray-700">
            <h3 className="font-bold mb-2 text-sm">Resources</h3>

            <div className="space-y-1 text-gray-300">
              <div className="flex justify-between">
                <span>Total:</span>
                <span>{resources.total} files</span>
              </div>
              <div className="flex justify-between">
                <span>Cached:</span>
                <span className="text-green-400">{resources.cached}</span>
              </div>
              <div className="flex justify-between">
                <span>Network:</span>
                <span className="text-yellow-400">{resources.fromNetwork}</span>
              </div>
              <div className="flex justify-between">
                <span>Total Size:</span>
                <span>{resources.totalSize} KB</span>
              </div>
              <div className="flex justify-between">
                <span>Avg Load:</span>
                <span>{resources.avgLoadTime} ms</span>
              </div>
            </div>
          </div>

          {/* Legend */}
          <div className="pt-2 border-t border-gray-700">
            <div className="flex gap-3 text-xs mb-2">
              <span className="flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-green-500" /> Good
              </span>
              <span className="flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-yellow-500" /> Needs Work
              </span>
              <span className="flex items-center gap-1">
                <span className="w-2 h-2 rounded-full bg-red-500" /> Poor
              </span>
            </div>
          </div>
        </div>
      )}
      <style jsx>{`
        .performance-score-fill {
          width: ${score}%;
        }
      `}</style>
    </div>
  )
}

/**
 * Metric Row Component
 */
interface MetricRowProps {
  name: string
  value: number
  unit: string
  rating: 'good' | 'needs-improvement' | 'poor'
  threshold?: { good: number; poor: number }
  decimals?: number
}

function MetricRow({
  name,
  value,
  unit,
  rating,
  threshold,
  decimals = 0,
}: MetricRowProps) {
  const formattedValue =
    decimals > 0 ? value.toFixed(decimals) : Math.round(value)

  return (
    <div className="flex justify-between items-center py-1 text-gray-300 text-xs">
      <span>{name}</span>
      <div className="flex items-center gap-2">
        <span className={`font-mono ${getRatingColor(rating)}`}>
          {formattedValue}{unit}
        </span>
        <span className={`w-2 h-2 rounded-full ${getRatingBgColor(rating)}`} />
      </div>
    </div>
  )
}

/**
 * Helper: Get color class for score
 */
function getScoreColor(score: number): string {
  if (score >= 90) return 'text-green-400'
  if (score >= 75) return 'text-yellow-400'
  return 'text-red-400'
}

/**
 * Helper: Get color class for rating
 */
function getRatingColor(rating: string): string {
  switch (rating) {
    case 'good':
      return 'text-green-400'
    case 'needs-improvement':
      return 'text-yellow-400'
    case 'poor':
      return 'text-red-400'
    default:
      return 'text-gray-400'
  }
}

/**
 * Helper: Get background color class for rating
 */
function getRatingBgColor(rating: string): string {
  switch (rating) {
    case 'good':
      return 'bg-green-500'
    case 'needs-improvement':
      return 'bg-yellow-500'
    case 'poor':
      return 'bg-red-500'
    default:
      return 'bg-gray-500'
  }
}

/**
 * Lightweight metrics badge (minimal UI)
 */
export function PerformanceBadge() {
  const { score, rating } = usePerformanceScore()

  if (process.env.NODE_ENV === 'production') {
    return null
  }

  return (
    <div
      className={`fixed bottom-4 left-4 w-12 h-12 rounded-full flex items-center justify-center font-bold text-white text-xs ${
        score >= 90
          ? 'bg-green-500'
          : score >= 75
            ? 'bg-yellow-500'
            : 'bg-red-500'
      }`}
      title={`Performance: ${score} (${rating})`}
    >
      {score}
    </div>
  )
}

/**
 * Metrics Summary Card
 */
export function PerformanceMetricsCard() {
  const vitals = useWebVitals()
  const { score } = usePerformanceScore()
  const resources = useResourceMetrics()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) return null

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 shadow-sm">
      <h2 className="text-lg font-bold mb-4">Performance Metrics</h2>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        {/* Score */}
        <MetricCard
          title="Score"
          value={score}
          suffix="/100"
          color={score >= 90 ? 'green' : score >= 75 ? 'yellow' : 'red'}
        />

        {/* LCP */}
        {vitals.lcp && (
          <MetricCard
            title="LCP"
            value={Math.round(vitals.lcp.value)}
            suffix="ms"
            color={vitals.lcp.rating === 'good' ? 'green' : vitals.lcp.rating === 'needs-improvement' ? 'yellow' : 'red'}
            target="<2.5s"
          />
        )}

        {/* FID/INP */}
        {(vitals.fid || vitals.inp) && (
          <MetricCard
            title={vitals.inp ? 'INP' : 'FID'}
            value={Math.round((vitals.inp?.value || vitals.fid?.value) ?? 0)}
            suffix="ms"
            color={
              (vitals.inp?.rating || vitals.fid?.rating) === 'good'
                ? 'green'
                : (vitals.inp?.rating || vitals.fid?.rating) === 'needs-improvement'
                  ? 'yellow'
                  : 'red'
            }
            target={vitals.inp ? '<200ms' : '<100ms'}
          />
        )}

        {/* CLS */}
        {vitals.cls && (
          <MetricCard
            title="CLS"
            value={vitals.cls.value.toFixed(2)}
            suffix=""
            color={vitals.cls.rating === 'good' ? 'green' : vitals.cls.rating === 'needs-improvement' ? 'yellow' : 'red'}
            target="<0.1"
          />
        )}

        {/* Resources */}
        <MetricCard
          title="Resources"
          value={resources.total}
          suffix={` (${resources.cached} cached)`}
        />

        {/* Size */}
        <MetricCard
          title="Total Size"
          value={resources.totalSize}
          suffix=" KB"
        />
      </div>
    </div>
  )
}

/**
 * Individual Metric Card
 */
interface MetricCardProps {
  title: string
  value: string | number
  suffix?: string
  color?: 'green' | 'yellow' | 'red'
  target?: string
}

function MetricCard({ title, value, suffix = '', color, target }: MetricCardProps) {
  const colorClass = {
    green: 'text-green-600',
    yellow: 'text-yellow-600',
    red: 'text-red-600',
  }[color || 'gray']

  return (
    <div>
      <p className="text-sm text-gray-600 mb-1">{title}</p>
      <p className={`text-2xl font-bold ${colorClass}`}>
        {value}
        {suffix}
      </p>
      {target && <p className="text-xs text-gray-500 mt-1">Target: {target}</p>}
    </div>
  )
}
