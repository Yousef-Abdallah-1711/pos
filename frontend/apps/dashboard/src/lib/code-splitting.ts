/**
 * Code Splitting Utilities
 *
 * This module provides utilities for route-based code splitting in Next.js.
 * It helps reduce initial bundle size by lazy-loading route-specific code.
 */

import dynamic from 'next/dynamic'
import { ComponentType } from 'react'
import { debugWarn } from '@/lib/debug'

// Loading skeleton component
export const LoadingFallback = () => (
  <div className="flex items-center justify-center h-screen">
    <div className="text-center">
      <div className="inline-block">
        <div className="w-12 h-12 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin mb-4" />
      </div>
      <p className="text-gray-600 text-sm">Loading...</p>
    </div>
  </div>
)

/**
 * Create a code-split route component
 * Automatically handles loading state with skeleton
 */
export function createDynamicRoute<P extends object>(
  importFn: () => Promise<{ default: ComponentType<P> }>,
  options?: {
    loading?: ComponentType
    ssr?: boolean
    delay?: number
  }
) {
  return dynamic(importFn, {
    loading: options?.loading || LoadingFallback,
    ssr: options?.ssr !== false,
    ...(options?.delay && { delay: options.delay })
  })
}

/**
 * Map of dynamic imports for each dashboard route
 * These will be lazy-loaded when the route is accessed
 */
export const routeCodeSplits = {
  // Dashboard
  dashboard: {
    index: () => import('@/app/[locale]/dashboard/page'),

    // Products section
    products: {
      list: () => import('@/app/[locale]/dashboard/products/page'),
      new: () => import('@/app/[locale]/dashboard/products/new/page'),
      edit: () => import('@/app/[locale]/dashboard/products/[id]/edit/page'),
    },

    // Categories section
    categories: {
      list: () => import('@/app/[locale]/dashboard/categories/page'),
      new: () => import('@/app/[locale]/dashboard/categories/new/page'),
      edit: () => import('@/app/[locale]/dashboard/categories/[id]/edit/page'),
    },

    // Theme Builder section
    themeBuilder: {
      list: () => import('@/app/[locale]/dashboard/theme-builder/page'),
      editor: () => import('@/app/[locale]/dashboard/theme-builder/editor/[id]/page'),
    },

    // HR section
    hr: {
      employees: () => import('@/app/[locale]/dashboard/hr/employees/page'),
      attendance: () => import('@/app/[locale]/dashboard/hr/attendance/page'),
      leaves: () => import('@/app/[locale]/dashboard/hr/leaves/page'),
      payroll: () => import('@/app/[locale]/dashboard/hr/payroll/page'),
    },

    // Other sections
    inventory: () => import('@/app/[locale]/dashboard/inventory/page'),
    notifications: () => import('@/app/[locale]/dashboard/notifications/page'),
    settings: () => import('@/app/[locale]/dashboard/settings/page'),
  }
}

/**
 * Preload a specific route to improve performance
 * Call this on route hover or link visibility
 */
export function preloadRoute(
  importFn: () => Promise<any>,
  delay: number = 0
): void {
  if (typeof window === 'undefined') return

  const timer = setTimeout(() => {
    importFn().catch((err) => {
      debugWarn('Failed to preload route:', err)
    })
  }, delay)

  // Cleanup timer if component unmounts
  return () => clearTimeout(timer)
}

/**
 * Batch preload multiple routes
 * Useful for preloading common navigation paths
 */
export function preloadRoutes(
  importFns: Array<() => Promise<any>>,
  options?: {
    staggerDelay?: number
    maxConcurrent?: number
  }
): void {
  if (typeof window === 'undefined') return

  const { staggerDelay = 1000, maxConcurrent = 3 } = options || {}

  let loadingCount = 0
  let queueIndex = 0

  const loadNext = () => {
    if (queueIndex >= importFns.length) return

    if (loadingCount < maxConcurrent) {
      const importFn = importFns[queueIndex++]
      loadingCount++

      importFn()
        .catch((err) => {
          debugWarn('Failed to preload route:', err)
        })
        .finally(() => {
          loadingCount--
          loadNext()
        })

      setTimeout(loadNext, staggerDelay)
    }
  }

  loadNext()
}

/**
 * Common navigation routes for preloading on app load
 */
export const commonRoutes = [
  () => import('@/app/[locale]/dashboard/page'),
  () => import('@/app/[locale]/dashboard/products/page'),
  () => import('@/app/[locale]/dashboard/inventory/page'),
]

/**
 * Heavy routes that should be preloaded on hover
 */
export const heavyRoutes = [
  () => import('@/app/[locale]/dashboard/theme-builder/page'),
  () => import('@/app/[locale]/dashboard/hr/employees/page'),
]
