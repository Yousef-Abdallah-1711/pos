'use client'

import { useEffect, ReactNode, useMemo } from 'react'
import { useThemeStore } from '@/lib/store/theme-store'
import { hexToHsl } from '@/lib/utils/theme-colors'
import { applyTypographyVariables } from '@/lib/utils/theme-typography'

/**
 * ThemeProvider component
 * Injects CSS variables when theme loads
 * Watches store for theme changes and updates DOM
 */
export function ThemeProvider({ children }: { children: ReactNode }) {
  const { currentTheme, isLoading, error } = useThemeStore()

  // Memoize color conversions to avoid recalculating expensive HSL conversions
  const colorVariables = useMemo(
    () => ({
      '--theme-primary': hexToHsl(currentTheme?.colors.primary || '#000000'),
      '--theme-secondary': hexToHsl(currentTheme?.colors.secondary || '#000000'),
      '--theme-accent': hexToHsl(currentTheme?.colors.accent || '#000000'),
      '--theme-background': hexToHsl(currentTheme?.colors.background || '#ffffff'),
      '--theme-text': hexToHsl(currentTheme?.colors.text || '#000000'),
      '--theme-border': hexToHsl(currentTheme?.colors.border || '#e5e7eb'),
      '--theme-shadow': hexToHsl(currentTheme?.colors.shadow || '#000000'),
    }),
    [currentTheme]
  )

  /**
   * Apply theme colors and typography to DOM
   * Uses requestAnimationFrame to batch updates and avoid layout thrashing
   * Color conversions are memoized, so they only recalculate when theme changes
   */
  useEffect(() => {
    if (!currentTheme || typeof document === 'undefined') return

    // Inject color variables
    requestAnimationFrame(() => {
      const root = document.documentElement

      // Build cssText for all variables at once (single update instead of 7 setProperty calls)
      const cssVariables = Object.entries(colorVariables)
        .map(([key, value]) => `${key}: ${value}`)
        .join('; ')

      // Apply all at once using cssText (batch update = better performance)
      root.style.cssText = root.style.cssText
        .split(';')
        .filter(prop => !colorVariables.hasOwnProperty(prop.split(':')[0]?.trim()))
        .concat([cssVariables])
        .join('; ')

      // Apply typography variables (handles font, sizes, etc.)
      applyTypographyVariables(currentTheme.typography)
    })
  }, [currentTheme, colorVariables])

  return <>{children}</>
}

/**
 * Theme Loading Skeleton
 * Shows while theme is loading
 */
export function ThemeLoadingSkeleton() {
  return (
    <div
      className="fixed inset-0 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-900 dark:to-gray-800 animate-pulse"
      aria-hidden="true"
    />
  )
}

/**
 * Theme Error Fallback
 * Shows when theme fails to load (used with error boundary)
 */
export function ThemeErrorFallback({ error, retry }: { error?: string; retry?: () => void }) {
  return (
    <div className="flex items-center justify-center min-h-screen bg-white dark:bg-gray-900">
      <div className="max-w-md w-full mx-auto px-6 py-8 bg-gray-50 dark:bg-gray-800 rounded-lg shadow-lg">
        <div className="text-center">
          <svg
            className="mx-auto h-12 w-12 text-red-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 8v4m0 4v.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>

          <h2 className="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
            Theme Loading Failed
          </h2>

          {error && (
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              {error}
            </p>
          )}

          <p className="mt-4 text-sm text-gray-600 dark:text-gray-400">
            Using default theme. The app will continue to work normally.
          </p>

          {retry && (
            <button
              onClick={retry}
              className="mt-6 w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-colors"
            >
              Try Again
            </button>
          )}
        </div>
      </div>
    </div>
  )
}

/**
 * Provider wrapper with error boundary handling
 * Use this instead of ThemeProvider if you want loading/error states
 */
export function ThemeProviderWithFallback({ children }: { children: ReactNode }) {
  const { currentTheme, isLoading, error } = useThemeStore()

  // Show fallback while loading
  if (isLoading && !currentTheme) {
    return (
      <>
        <ThemeLoadingSkeleton />
        {children}
      </>
    )
  }

  // Show error state if failed (but keep children rendered)
  if (error && !currentTheme) {
    return (
      <>
        <ThemeErrorFallback error={error} />
        {children}
      </>
    )
  }

  return <ThemeProvider>{children}</ThemeProvider>
}
