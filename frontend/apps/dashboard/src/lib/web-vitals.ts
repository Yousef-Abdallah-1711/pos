/**
 * Web Vitals Performance Monitoring
 *
 * Tracks Core Web Vitals metrics:
 * - LCP (Largest Contentful Paint) - target: <2.5s
 * - FID (First Input Delay) - target: <100ms
 * - CLS (Cumulative Layout Shift) - target: <0.1
 * - TTFB (Time to First Byte) - target: <600ms
 * - FCP (First Contentful Paint) - target: <1.8s
 */

import { debugError, debugWarn } from '@/lib/debug'

export interface WebVitalMetric {
  name: string
  value: number
  rating: 'good' | 'needs-improvement' | 'poor'
  delta?: number
  id?: string
  navigationType?: string
}

export interface PerformanceThresholds {
  lcp: { good: number; poor: number }
  fid: { good: number; poor: number }
  cls: { good: number; poor: number }
  ttfb: { good: number; poor: number }
  fcp: { good: number; poor: number }
}

// Web Vitals thresholds (Google recommendations)
export const PERFORMANCE_THRESHOLDS: PerformanceThresholds = {
  lcp: { good: 2500, poor: 4000 },      // milliseconds
  fid: { good: 100, poor: 300 },        // milliseconds
  cls: { good: 0.1, poor: 0.25 },       // unitless
  ttfb: { good: 600, poor: 1800 },      // milliseconds
  fcp: { good: 1800, poor: 3000 },      // milliseconds
}

// Type definitions for Web Vitals API
interface PerformanceEntryHandler {
  (entry: PerformanceEntry): void
}

interface VitalsCallback {
  (metric: WebVitalMetric): void
}

/**
 * Utility function to convert bytes to readable format
 */
function bytesToKB(bytes: number): string {
  return (bytes / 1024).toFixed(2)
}

/**
 * Determine if metric is good, needs improvement, or poor
 */
function getRating(
  name: string,
  value: number
): 'good' | 'needs-improvement' | 'poor' {
  const thresholds = PERFORMANCE_THRESHOLDS[name.toLowerCase() as keyof PerformanceThresholds]

  if (!thresholds) {
    return 'needs-improvement'
  }

  if (value <= thresholds.good) {
    return 'good'
  }

  if (value <= thresholds.poor) {
    return 'needs-improvement'
  }

  return 'poor'
}

/**
 * Get Largest Contentful Paint (LCP)
 *
 * Measures when the largest content element becomes visible
 * Target: < 2.5 seconds
 */
export function observeLCP(callback: VitalsCallback): () => void {
  if (!('PerformanceObserver' in window)) {
    debugWarn('[Web Vitals] PerformanceObserver not supported')
    return () => {}
  }

  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()
      const lastEntry = entries[entries.length - 1]

      const lcpEntry = lastEntry as PerformanceEntry & { renderTime?: number; loadTime?: number }
      const lcpTime = lcpEntry.renderTime || lcpEntry.loadTime || 0
      const metric: WebVitalMetric = {
        name: 'LCP',
        value: lcpTime,
        rating: getRating('lcp', lcpTime),
        id: `lcp-${Date.now()}`,
        navigationType: 'navigation',
      }

      callback(metric)
    })

    observer.observe({
      entryTypes: ['largest-contentful-paint'],
      buffered: true,
    })

    return () => observer.disconnect()
  } catch (error) {
    debugError('[Web Vitals] LCP observation failed:', error)
    return () => {}
  }
}

/**
 * Get First Input Delay (FID)
 *
 * Measures delay between user input and response
 * Target: < 100 milliseconds
 *
 * Note: Being replaced by INP in newer browsers
 */
export function observeFID(callback: VitalsCallback): () => void {
  if (!('PerformanceObserver' in window)) {
    debugWarn('[Web Vitals] PerformanceObserver not supported')
    return () => {}
  }

  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()

      entries.forEach((entry: any) => {
        const metric: WebVitalMetric = {
          name: 'FID',
          value: entry.processingDuration,
          rating: getRating('fid', entry.processingDuration),
          id: `fid-${Date.now()}`,
          navigationType: 'navigation',
        }

        callback(metric)
      })
    })

    observer.observe({
      entryTypes: ['first-input'],
      buffered: true,
    })

    return () => observer.disconnect()
  } catch (error) {
    debugError('[Web Vitals] FID observation failed:', error)
    return () => {}
  }
}

/**
 * Get Cumulative Layout Shift (CLS)
 *
 * Measures unexpected layout shifts
 * Target: < 0.1 (unitless)
 */
export function observeCLS(callback: VitalsCallback): () => void {
  if (!('PerformanceObserver' in window)) {
    debugWarn('[Web Vitals] PerformanceObserver not supported')
    return () => {}
  }

  try {
    let clsValue = 0

    const observer = new PerformanceObserver((list) => {
      list.getEntries().forEach((entry: any) => {
        if (!entry.hadRecentInput) {
          clsValue += entry.value

          const metric: WebVitalMetric = {
            name: 'CLS',
            value: clsValue,
            rating: getRating('cls', clsValue),
            id: `cls-${Date.now()}`,
            navigationType: 'navigation',
          }

          callback(metric)
        }
      })
    })

    observer.observe({
      entryTypes: ['layout-shift'],
      buffered: true,
    })

    return () => observer.disconnect()
  } catch (error) {
    debugError('[Web Vitals] CLS observation failed:', error)
    return () => {}
  }
}

/**
 * Get Time to First Byte (TTFB)
 *
 * Measures time to first byte of response
 * Target: < 600 milliseconds
 */
export function observeTTFB(callback: VitalsCallback): () => void {
  try {
    const navigationEntry = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming

    if (navigationEntry) {
      const ttfb = navigationEntry.responseStart - navigationEntry.fetchStart

      const metric: WebVitalMetric = {
        name: 'TTFB',
        value: ttfb,
        rating: getRating('ttfb', ttfb),
        id: `ttfb-${Date.now()}`,
        navigationType: 'navigation',
      }

      callback(metric)
    }

    return () => {}
  } catch (error) {
    debugError('[Web Vitals] TTFB observation failed:', error)
    return () => {}
  }
}

/**
 * Get First Contentful Paint (FCP)
 *
 * Measures when first content appears
 * Target: < 1.8 seconds
 */
export function observeFCP(callback: VitalsCallback): () => void {
  if (!('PerformanceObserver' in window)) {
    debugWarn('[Web Vitals] PerformanceObserver not supported')
    return () => {}
  }

  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()
      const lastEntry = entries[entries.length - 1]

      const metric: WebVitalMetric = {
        name: 'FCP',
        value: lastEntry.startTime,
        rating: getRating('fcp', lastEntry.startTime),
        id: `fcp-${Date.now()}`,
        navigationType: 'navigation',
      }

      callback(metric)
    })

    observer.observe({
      entryTypes: ['paint'],
      buffered: true,
    })

    return () => observer.disconnect()
  } catch (error) {
    debugError('[Web Vitals] FCP observation failed:', error)
    return () => {}
  }
}

/**
 * Get Interaction to Next Paint (INP)
 *
 * Measures responsiveness to interactions
 * Target: < 200 milliseconds
 * Replaces FID in newer browsers
 */
export function observeINP(callback: VitalsCallback): () => void {
  if (!('PerformanceObserver' in window)) {
    debugWarn('[Web Vitals] PerformanceObserver not supported')
    return () => {}
  }

  try {
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries()

      entries.forEach((entry: any) => {
        const metric: WebVitalMetric = {
          name: 'INP',
          value: entry.duration,
          rating: entry.duration <= 200 ? 'good' : entry.duration <= 500 ? 'needs-improvement' : 'poor',
          id: `inp-${Date.now()}`,
          navigationType: 'navigation',
        }

        callback(metric)
      })
    })

    observer.observe({
      entryTypes: ['interaction'],
      buffered: true,
    })

    return () => observer.disconnect()
  } catch (error) {
    debugError('[Web Vitals] INP observation failed:', error)
    return () => {}
  }
}

/**
 * Get all Core Web Vitals at once
 */
export function observeAllWebVitals(callback: VitalsCallback): () => void {
  const unsubscribes = [
    observeLCP(callback),
    observeFID(callback),
    observeCLS(callback),
    observeTTFB(callback),
    observeFCP(callback),
    observeINP(callback),
  ]

  // Return cleanup function
  return () => {
    unsubscribes.forEach((unsub) => unsub())
  }
}

/**
 * Get performance summary
 */
export function getPerformanceSummary(): Record<string, any> {
  const navigationEntry = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming

  return {
    // Navigation timing
    dnslookup: navigationEntry.domainLookupEnd - navigationEntry.domainLookupStart,
    tcp: navigationEntry.connectEnd - navigationEntry.connectStart,
    ttfb: navigationEntry.responseStart - navigationEntry.fetchStart,
    download: navigationEntry.responseEnd - navigationEntry.responseStart,
    domParsing: navigationEntry.domInteractive - (navigationEntry as PerformanceNavigationTiming & { domLoading?: number }).domLoading!,
    domInteractive: navigationEntry.domInteractive - navigationEntry.fetchStart,
    domComplete: navigationEntry.domComplete - navigationEntry.fetchStart,
    loadComplete: navigationEntry.loadEventEnd - navigationEntry.fetchStart,

    // Resource timing
    resourcesCount: performance.getEntriesByType('resource').length,
    resourcesSize: Math.round(
      performance.getEntriesByType('resource')
        .reduce((sum, entry: any) => sum + (entry.transferSize || 0), 0) / 1024
    ),

    // Memory (if available)
    memory: (performance as any).memory ? {
      usedJSHeapSize: bytesToKB((performance as any).memory.usedJSHeapSize),
      totalJSHeapSize: bytesToKB((performance as any).memory.totalJSHeapSize),
      jsHeapSizeLimit: bytesToKB((performance as any).memory.jsHeapSizeLimit),
    } : null,
  }
}

/**
 * Format metric for logging
 */
export function formatMetric(metric: WebVitalMetric): string {
  const value =
    metric.name === 'CLS'
      ? metric.value.toFixed(3)
      : `${metric.value.toFixed(0)}ms`

  const ratingEmoji = {
    good: '✅',
    'needs-improvement': '⚠️',
    poor: '❌',
  }

  return `${ratingEmoji[metric.rating]} ${metric.name}: ${value} (${metric.rating})`
}

/**
 * Send metrics to analytics service
 */
export function sendMetricToAnalytics(metric: WebVitalMetric): void {
  // Send to your analytics service
  if (typeof window !== 'undefined' && (window as any).gtag) {
    (window as any).gtag('event', metric.name.toLowerCase(), {
      value: Math.round(metric.value),
      event_category: 'web_vitals',
      event_label: metric.id,
      non_interaction: true,
    })
  }

  // Send to custom analytics
  if (typeof navigator !== 'undefined' && navigator.sendBeacon) {
    const data = new FormData()
    data.append('name', metric.name)
    data.append('value', String(metric.value))
    data.append('rating', metric.rating)

    navigator.sendBeacon('/api/v1/analytics/metrics', data)
  }
}
