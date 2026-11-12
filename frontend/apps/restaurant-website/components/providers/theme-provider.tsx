'use client'

import { useEffect, ReactNode, useMemo } from 'react'
import { ThemeProvider as NextThemesProvider, useTheme as useNextTheme } from 'next-themes'
import { useThemeStore } from '@/lib/store/theme-store'
import { memoHexToHSL } from '@/lib/utils/hex-to-hsl'
import { applyTypographyVariables } from '@/lib/utils/theme-typography'

/**
 * ThemeProvider component
 * Injects CSS variables when theme loads
 * Watches store for theme changes and updates DOM
 */
function ThemeCSSInjector() {
  const { currentTheme } = useThemeStore()
  const { resolvedTheme } = useNextTheme()
  const isDark = resolvedTheme === 'dark'

  const colorVariables = useMemo(() => {
    const colors = (currentTheme?.colors ?? {}) as Record<string, string | undefined>
    const primary = isDark ? colors.primary_dark ?? colors.primary : colors.primary
    const secondary = isDark ? colors.secondary_dark ?? colors.secondary : colors.secondary
    const accent = isDark ? colors.accent_dark ?? colors.accent : colors.accent
    const background = isDark ? colors.background_dark ?? colors.background : colors.background
    const surface = isDark ? colors.surface_dark ?? colors.surface_color ?? colors.background : colors.surface_color ?? colors.background
    const text = isDark ? colors.text_dark ?? colors.text : colors.text
    const border = isDark ? colors.border_dark ?? colors.border : colors.border
    const shadow = isDark ? colors.shadow_dark ?? colors.shadow : colors.shadow
    const headerBg = isDark ? colors.header_bg_dark ?? colors.background_dark ?? colors.header_bg : colors.header_bg ?? colors.primary
    const headerText = isDark
      ? colors.header_text_dark ?? colors.text_dark ?? colors.header_text ?? '#ffffff'
      : colors.header_text ?? '#ffffff'
    const footerBg = isDark ? colors.footer_bg_dark ?? colors.background_dark ?? colors.footer_bg : colors.footer_bg ?? colors.text
    const footerText = isDark
      ? colors.footer_text_dark ?? colors.text_dark ?? colors.footer_text ?? '#ffffff'
      : colors.footer_text ?? '#ffffff'

    return {
      '--theme-primary': memoHexToHSL(primary ?? '#3b82f6'),
      '--theme-secondary': memoHexToHSL(secondary ?? '#1e40af'),
      '--theme-accent': memoHexToHSL(accent ?? '#f59e0b'),
      '--theme-bg': memoHexToHSL(background ?? '#ffffff'),
      '--theme-background': memoHexToHSL(background ?? '#ffffff'),
      '--theme-surface': memoHexToHSL(surface ?? '#f8fafc'),
      '--theme-surface-elevated': memoHexToHSL(surface ?? '#ffffff'),
      '--theme-text': memoHexToHSL(text ?? '#111827'),
      '--theme-text-muted': memoHexToHSL(text ?? '#6b7280'),
      '--theme-border': memoHexToHSL(border ?? '#e5e7eb'),
      '--theme-shadow': memoHexToHSL(shadow ?? '#111827'),
      '--theme-header-bg': memoHexToHSL(headerBg ?? '#111827'),
      '--theme-header-text': memoHexToHSL(headerText),
      '--theme-footer-bg': memoHexToHSL(footerBg ?? '#111827'),
      '--theme-footer-text': memoHexToHSL(footerText),
    }
  }, [currentTheme?.slug, (currentTheme as { version?: number } | null)?.version, isDark])

  /**
   * Apply theme colors and typography to DOM
   * Uses requestAnimationFrame to batch updates and avoid layout thrashing.
   */
  useEffect(() => {
    if (!currentTheme || typeof document === 'undefined') return

    requestAnimationFrame(() => {
      const root = document.documentElement
      for (const [key, value] of Object.entries(colorVariables)) {
        root.style.setProperty(key, value)
      }
      applyTypographyVariables(currentTheme.typography)
    })
  }, [currentTheme, colorVariables])

  return null
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <ThemeCSSInjector />
      {children}
    </NextThemesProvider>
  )
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
