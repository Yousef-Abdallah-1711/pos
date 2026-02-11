'use client'

import { useEffect } from 'react'
import { usePathname } from 'next/navigation'
import { debugLog, debugWarn } from '@/lib/debug'

/**
 * PerformanceMonitor Component
 * Tracks page load time for each route/page
 * Shows warnings in console if page takes longer than 3 seconds
 */
export function PerformanceMonitor() {
  const pathname = usePathname()

  useEffect(() => {
    const startTime = performance.now()

    return () => {
      const endTime = performance.now()
      const loadTime = endTime - startTime

      const emoji = loadTime > 5000 ? '❌' : loadTime > 3000 ? '⚠️' : '✅'

      debugLog(
        `${emoji} [Page Load] ${pathname} - ${Math.round(loadTime)}ms`
      )

      // Send to analytics
      if (typeof window !== 'undefined' && window.gtag) {
        window.gtag('event', 'page_load_time', {
          page_path: pathname,
          load_time_ms: Math.round(loadTime),
        })
      }

      // Warn if page is slow
      if (loadTime > 3000) {
        debugWarn(
          `⚠️ [Performance] Page ${pathname} is slow (${Math.round(loadTime)}ms)`
        )
      }
    }
  }, [pathname])

  return null
}
