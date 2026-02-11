'use client';

import { useEffect } from 'react';
import { debugWarn } from '../debug';
import type { Metric } from 'web-vitals';

/**
 * Hook to capture Core Web Vitals and send to monitoring endpoint
 * Tracks: LCP, FID, CLS, FCP, TTFB
 */
export function useWebVitals() {
  useEffect(() => {
    const handleMetric = async (metric: Metric) => {
      // Get monitoring endpoint from env or use default
      const endpoint = process.env.NEXT_PUBLIC_METRICS_API || '/api/metrics';

      try {
        // Send metric to backend for storage/analysis
        await fetch(endpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            // Web Vital specific fields
            name: metric.name, // 'LCP', 'FID', 'CLS', 'FCP', 'TTFB'
            value: metric.value, // Numeric value in milliseconds
            rating: metric.rating, // 'good', 'needs-improvement', 'poor'
            delta: metric.delta, // Change since last measurement

            // Context fields
            url: typeof window !== 'undefined' ? window.location.href : '',
            pathname: typeof window !== 'undefined' ? window.location.pathname : '',
            timestamp: new Date().toISOString(),
            userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',

            // Performance context
            pageTitle: typeof document !== 'undefined' ? document.title : '',
            referrer: typeof document !== 'undefined' ? document.referrer : '',
          }),
        });
      } catch (error) {
        // Silently fail - don't break app if monitoring fails
        if (process.env.NODE_ENV === 'development') {
          debugWarn('Failed to send Web Vital metric:', metric.name, error);
        }
      }
    };

    // Dynamically import web-vitals at runtime
    // web-vitals v5 exports onCLS, onINP (replaces onFID), onFCP, onLCP, onTTFB
    import('web-vitals').then((vitals: any) => {
      // Subscribe to all Core Web Vitals
      vitals.onCLS?.(handleMetric);
      vitals.onINP?.(handleMetric); // INP replaces FID in v5
      vitals.onFID?.(handleMetric); // Fallback for older versions
      vitals.onFCP?.(handleMetric);
      vitals.onLCP?.(handleMetric);
      vitals.onTTFB?.(handleMetric);
    }).catch((error) => {
      if (process.env.NODE_ENV === 'development') {
        debugWarn('Failed to load web-vitals:', error);
      }
    });
  }, []);
}

/**
 * Hook for advanced analytics tracking
 * Sends comprehensive performance data
 */
export function usePerformanceMetrics() {
  useEffect(() => {
    // Wait for page to fully load
    window.addEventListener('load', () => {
      // Get performance data from Navigation Timing API
      const perfData = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;

      if (perfData) {
        const metrics = {
          // Navigation timings
          dns: perfData.domainLookupEnd - perfData.domainLookupStart,
          tcp: perfData.connectEnd - perfData.connectStart,
          ttfb: perfData.responseStart - perfData.requestStart,
          download: perfData.responseEnd - perfData.responseStart,
          domInteractive: perfData.domInteractive - perfData.fetchStart,
          domComplete: perfData.domComplete - perfData.fetchStart,
          loadComplete: perfData.loadEventEnd - perfData.fetchStart,

          // Resource counts
          resourceCount: performance.getEntriesByType('resource').length,

          // Memory (if available)
          memory: (performance as any).memory ? {
            usedJSHeapSize: (performance as any).memory.usedJSHeapSize,
            totalJSHeapSize: (performance as any).memory.totalJSHeapSize,
            jsHeapSizeLimit: (performance as any).memory.jsHeapSizeLimit,
          } : null,
        };

        // Send to monitoring endpoint
        const endpoint = process.env.NEXT_PUBLIC_METRICS_API || '/api/metrics';
        fetch(endpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            type: 'performance',
            metrics,
            url: window.location.href,
            timestamp: new Date().toISOString(),
          }),
        }).catch(() => {
          // Silently fail
        });
      }
    });
  }, []);
}
