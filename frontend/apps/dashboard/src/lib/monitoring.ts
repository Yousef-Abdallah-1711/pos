import { debugLog, debugError } from '@/lib/debug'

export interface PerformanceMetric {
  name: string
  value: number
  rating: 'good' | 'needs-improvement' | 'poor'
  description: string
}

const ENABLE_DEBUG_LOGS = process.env.NEXT_PUBLIC_DEBUG_LOGS === 'true'

/**
 * Report all Web Vitals metrics
 * These are Google's Core Web Vitals that affect SEO and UX
 * Only works in browser environment
 */
export async function reportWebVitals() {
  // Only run in browser
  if (typeof window === 'undefined') return

  try {
    // Dynamically import web-vitals
    const webVitals = await import('web-vitals')

    // Largest Contentful Paint (LCP)
    // Measures when the largest visible element is painted
    // Target: < 2.5s
    webVitals.onLCP((metric) => {
      const rating =
        metric.value < 2500
          ? 'good'
          : metric.value < 4000
            ? 'needs-improvement'
            : 'poor'

      logMetric({
        name: 'LCP',
        value: Math.round(metric.value),
        rating,
        description: 'Largest Contentful Paint - Page load speed',
      })
    })

    // Interaction to Next Paint (INP)
    // Measures interaction latency - replaces FID in web-vitals v5
    // Target: < 200ms
    if (webVitals.onINP) {
      webVitals.onINP((metric) => {
        const rating =
          metric.value < 200
            ? 'good'
            : metric.value < 500
              ? 'needs-improvement'
              : 'poor'

        logMetric({
          name: 'INP',
          value: Math.round(metric.value),
          rating,
          description: 'Interaction to Next Paint - Interaction responsiveness',
        })
      })
    }

    // Cumulative Layout Shift (CLS)
    // Measures unexpected layout shifts
    // Target: < 0.1
    webVitals.onCLS((metric) => {
      const rating =
        metric.value < 0.1
          ? 'good'
          : metric.value < 0.25
            ? 'needs-improvement'
            : 'poor'

      logMetric({
        name: 'CLS',
        value: parseFloat(metric.value.toFixed(3)),
        rating,
        description: 'Cumulative Layout Shift - Layout stability',
      })
    })

    // First Contentful Paint (FCP)
    // Measures when first content is painted
    // Target: < 1.8s
    webVitals.onFCP((metric) => {
      const rating =
        metric.value < 1800
          ? 'good'
          : metric.value < 3000
            ? 'needs-improvement'
            : 'poor'

      logMetric({
        name: 'FCP',
        value: Math.round(metric.value),
        rating,
        description: 'First Contentful Paint - First content painted',
      })
    })

    // Time to First Byte (TTFB)
    // Measures server response time
    // Target: < 600ms
    webVitals.onTTFB((metric) => {
      const rating =
        metric.value < 600
          ? 'good'
          : metric.value < 1200
            ? 'needs-improvement'
            : 'poor'

      logMetric({
        name: 'TTFB',
        value: Math.round(metric.value),
        rating,
        description: 'Time to First Byte - Server response time',
      })
    })
  } catch (error) {
    if (ENABLE_DEBUG_LOGS) {
      debugError('Failed to load web-vitals:', error)
    }
  }
}

/**
 * Log a performance metric with emoji indicators
 */
function logMetric(metric: PerformanceMetric) {
  const emoji =
    metric.rating === 'good'
      ? '✅'
      : metric.rating === 'needs-improvement'
        ? '⚠️'
        : '❌'

  const unit = metric.name === 'CLS' ? '' : 'ms'

  if (ENABLE_DEBUG_LOGS) {
    debugLog(
      `${emoji} [${metric.name}] ${metric.value}${unit} - ${metric.description}`
    )
  }

  // Send to analytics if available
  if (typeof window !== 'undefined' && window.gtag) {
    window.gtag('event', 'web_vital', {
      metric_name: metric.name,
      metric_value: metric.value,
      metric_rating: metric.rating,
    })
  }
}

/**
 * Track page load time for specific routes
 */
export function trackPageLoad(pageName: string) {
  const startTime = performance.now()

  return () => {
    const endTime = performance.now()
    const loadTime = endTime - startTime

    const emoji = loadTime > 3000 ? '⚠️' : loadTime > 5000 ? '❌' : '✅'

    if (ENABLE_DEBUG_LOGS) {
      debugLog(
        `${emoji} [${pageName}] Load time: ${Math.round(loadTime)}ms`
      )
    }

    // Send to analytics
    if (typeof window !== 'undefined' && window.gtag) {
      window.gtag('event', 'page_load_time', {
        page_name: pageName,
        load_time: Math.round(loadTime),
      })
    }
  }
}

/**
 * Measure API response time
 */
export function measureApiCall(apiName: string) {
  const startTime = performance.now()

  return (success: boolean = true) => {
    const endTime = performance.now()
    const duration = endTime - startTime

    const emoji = duration > 2000 ? '⚠️' : '✅'

    if (ENABLE_DEBUG_LOGS) {
      debugLog(
        `${emoji} [API] ${apiName}: ${Math.round(duration)}ms ${success ? '✓' : '✗'}`
      )
    }

    // Send to analytics
    if (typeof window !== 'undefined' && window.gtag) {
      window.gtag('event', 'api_call', {
        api_name: apiName,
        duration: Math.round(duration),
        success,
      })
    }
  }
}

declare global {
  interface Window {
    gtag?: (
      command: string,
      event: string,
      params: Record<string, any>
    ) => void
  }
}
