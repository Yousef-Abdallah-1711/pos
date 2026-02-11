/**
 * React Hook for Web Vitals Monitoring
 *
 * Usage:
 * const vitals = useWebVitals()
 * // Access: vitals.lcp, vitals.fid, vitals.cls, etc.
 */

import { useEffect, useRef, useState } from 'react'
import { debugLog, debugWarn } from '@/lib/debug'
import {
  WebVitalMetric,
  observeLCP,
  observeFID,
  observeCLS,
  observeTTFB,
  observeFCP,
  observeINP,
  formatMetric,
  sendMetricToAnalytics,
} from '@/lib/web-vitals'

export interface WebVitalsData {
  lcp?: WebVitalMetric
  fid?: WebVitalMetric
  cls?: WebVitalMetric
  ttfb?: WebVitalMetric
  fcp?: WebVitalMetric
  inp?: WebVitalMetric
}

/**
 * Hook to monitor all Web Vitals
 */
export function useWebVitals() {
  const [vitals, setVitals] = useState<WebVitalsData>({})
  const unsubscribesRef = useRef<Array<() => void>>([])

  useEffect(() => {
    // LCP Observer
    unsubscribesRef.current.push(
      observeLCP((metric) => {
        setVitals((prev) => ({
          ...prev,
          lcp: metric,
        }))
        debugLog(formatMetric(metric))
      })
    )

    // FID Observer
    unsubscribesRef.current.push(
      observeFID((metric) => {
        setVitals((prev) => ({
          ...prev,
          fid: metric,
        }))
        debugLog(formatMetric(metric))
      })
    )

    // CLS Observer
    unsubscribesRef.current.push(
      observeCLS((metric) => {
        setVitals((prev) => ({
          ...prev,
          cls: metric,
        }))
        debugLog(formatMetric(metric))
      })
    )

    // TTFB Observer
    unsubscribesRef.current.push(
      observeTTFB((metric) => {
        setVitals((prev) => ({
          ...prev,
          ttfb: metric,
        }))
        debugLog(formatMetric(metric))
      })
    )

    // FCP Observer
    unsubscribesRef.current.push(
      observeFCP((metric) => {
        setVitals((prev) => ({
          ...prev,
          fcp: metric,
        }))
        debugLog(formatMetric(metric))
      })
    )

    // INP Observer (if supported)
    if ('PerformanceObserver' in window) {
      try {
        unsubscribesRef.current.push(
          observeINP((metric) => {
            setVitals((prev) => ({
              ...prev,
              inp: metric,
            }))
            debugLog(formatMetric(metric))
          })
        )
      } catch (error) {
        debugWarn('[Web Vitals] INP not supported')
      }
    }

    // Cleanup
    return () => {
      unsubscribesRef.current.forEach((unsub) => unsub())
    }
  }, [])

  return vitals
}

/**
 * Hook to send Web Vitals to analytics
 */
export function useWebVitalsAnalytics() {
  const vitals = useWebVitals()

  useEffect(() => {
    // Send each metric to analytics as it's collected
    Object.values(vitals).forEach((metric) => {
      if (metric) {
        sendMetricToAnalytics(metric)
      }
    })
  }, [vitals])

  return vitals
}

/**
 * Hook to get performance score (0-100)
 */
export function usePerformanceScore() {
  const vitals = useWebVitals()
  const [score, setScore] = useState(100)

  useEffect(() => {
    let newScore = 100

    // LCP scoring
    if (vitals.lcp) {
      if (vitals.lcp.rating === 'poor') newScore -= 30
      else if (vitals.lcp.rating === 'needs-improvement') newScore -= 15
    }

    // FID/INP scoring
    if (vitals.fid) {
      if (vitals.fid.rating === 'poor') newScore -= 20
      else if (vitals.fid.rating === 'needs-improvement') newScore -= 10
    }

    if (vitals.inp) {
      if (vitals.inp.rating === 'poor') newScore -= 20
      else if (vitals.inp.rating === 'needs-improvement') newScore -= 10
    }

    // CLS scoring
    if (vitals.cls) {
      if (vitals.cls.rating === 'poor') newScore -= 25
      else if (vitals.cls.rating === 'needs-improvement') newScore -= 12
    }

    // TTFB scoring
    if (vitals.ttfb) {
      if (vitals.ttfb.rating === 'poor') newScore -= 10
      else if (vitals.ttfb.rating === 'needs-improvement') newScore -= 5
    }

    setScore(Math.max(0, newScore))
  }, [vitals])

  return {
    score,
    rating: score >= 90 ? 'excellent' : score >= 75 ? 'good' : score >= 50 ? 'fair' : 'poor',
  }
}

/**
 * Hook to track user interaction metrics
 */
export function useInteractionMetrics() {
  const [interactions, setInteractions] = useState({
    clicks: 0,
    keystrokes: 0,
    scrolls: 0,
  })

  useEffect(() => {
    const handleClick = () => {
      setInteractions((prev) => ({
        ...prev,
        clicks: prev.clicks + 1,
      }))
    }

    const handleKeydown = () => {
      setInteractions((prev) => ({
        ...prev,
        keystrokes: prev.keystrokes + 1,
      }))
    }

    const handleScroll = () => {
      setInteractions((prev) => ({
        ...prev,
        scrolls: prev.scrolls + 1,
      }))
    }

    window.addEventListener('click', handleClick)
    window.addEventListener('keydown', handleKeydown)
    window.addEventListener('scroll', handleScroll, { passive: true })

    return () => {
      window.removeEventListener('click', handleClick)
      window.removeEventListener('keydown', handleKeydown)
      window.removeEventListener('scroll', handleScroll)
    }
  }, [])

  return interactions
}

/**
 * Hook to monitor long tasks
 */
export function useLongTasks() {
  const [longTasks, setLongTasks] = useState<Array<{ duration: number; start: number }>>([])

  useEffect(() => {
    if (!('PerformanceObserver' in window)) {
      return
    }

    try {
      const observer = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry: any) => {
          setLongTasks((prev) => [
            ...prev,
            {
              duration: entry.duration,
              start: entry.startTime,
            },
          ])
        })
      })

      observer.observe({ entryTypes: ['longtask'] })

      return () => observer.disconnect()
    } catch (error) {
      debugWarn('[Long Tasks] Not supported:', error)
    }
  }, [])

  return longTasks
}

/**
 * Hook to get resource loading metrics
 */
export function useResourceMetrics() {
  const [resources, setResources] = useState({
    total: 0,
    cached: 0,
    fromNetwork: 0,
    totalSize: 0, // in KB
    avgLoadTime: 0, // in ms
  })

  useEffect(() => {
    const entries = performance.getEntriesByType('resource')

    const cached = entries.filter((entry: any) => entry.transferSize === 0).length
    const fromNetwork = entries.length - cached

    const totalSize = Math.round(
      entries.reduce((sum, entry: any) => sum + (entry.transferSize || 0), 0) / 1024
    )

    const avgLoadTime = Math.round(
      entries.reduce((sum, entry: any) => sum + (entry.duration || 0), 0) / entries.length
    )

    setResources({
      total: entries.length,
      cached,
      fromNetwork,
      totalSize,
      avgLoadTime,
    })
  }, [])

  return resources
}
