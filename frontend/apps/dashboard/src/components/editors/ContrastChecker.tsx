'use client'

import { useMemo } from 'react'
import { AlertCircle, CheckCircle, AlertTriangle } from 'lucide-react'
import { GlobalColors } from '@/types/theme'

interface ContrastCheckerProps {
  colors: GlobalColors
  showDetails?: boolean
}

interface ContrastResult {
  colorKey: string
  textColor: string
  backgroundColor: string
  contrastRatio: number
  wcagAA: boolean
  wcagAAA: boolean
  largeTextAA: boolean
  largeTextAAA: boolean
}

/**
 * Calculate relative luminance of a color using WCAG formula
 * https://www.w3.org/TR/WCAG20/#relativeluminancedef
 */
function getLuminance(hexColor: string): number {
  // Parse hex color
  const hex = hexColor.replace('#', '')
  const r = parseInt(hex.substring(0, 2), 16) / 255
  const g = parseInt(hex.substring(2, 4), 16) / 255
  const b = parseInt(hex.substring(4, 6), 16) / 255

  // Apply gamma correction
  const [rs, gs, bs] = [r, g, b].map((c) => {
    return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4)
  })

  // Calculate relative luminance
  return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs
}

/**
 * Calculate contrast ratio between two colors
 * https://www.w3.org/TR/WCAG20/#contrast-ratiodef
 */
function getContrastRatio(color1: string, color2: string): number {
  const lum1 = getLuminance(color1)
  const lum2 = getLuminance(color2)

  const lighter = Math.max(lum1, lum2)
  const darker = Math.min(lum1, lum2)

  return (lighter + 0.05) / (darker + 0.05)
}

/**
 * Check if contrast ratio meets WCAG standards
 */
function checkWCAG(ratio: number): {
  wcagAA: boolean
  wcagAAA: boolean
  largeTextAA: boolean
  largeTextAAA: boolean
} {
  return {
    wcagAA: ratio >= 4.5, // Normal text
    wcagAAA: ratio >= 7, // Enhanced contrast
    largeTextAA: ratio >= 3, // Large text (18pt+)
    largeTextAAA: ratio >= 4.5 // Large enhanced text
  }
}

export function ContrastChecker({ colors, showDetails = true }: ContrastCheckerProps) {
  const contrastResults = useMemo<ContrastResult[]>(() => {
    const textColor = colors.text
    const bgColor = colors.background

    // Check contrast between text and background
    const primaryContrast = getContrastRatio(colors.primary, bgColor)
    const secondaryContrast = getContrastRatio(colors.secondary, bgColor)
    const accentContrast = getContrastRatio(colors.accent, bgColor)
    const textContrast = getContrastRatio(textColor, bgColor)
    const borderContrast = getContrastRatio(colors.border, bgColor)

    return [
      {
        colorKey: 'primary',
        textColor: colors.primary,
        backgroundColor: bgColor,
        contrastRatio: primaryContrast,
        ...checkWCAG(primaryContrast)
      },
      {
        colorKey: 'secondary',
        textColor: colors.secondary,
        backgroundColor: bgColor,
        contrastRatio: secondaryContrast,
        ...checkWCAG(secondaryContrast)
      },
      {
        colorKey: 'accent',
        textColor: colors.accent,
        backgroundColor: bgColor,
        contrastRatio: accentContrast,
        ...checkWCAG(accentContrast)
      },
      {
        colorKey: 'text',
        textColor: textColor,
        backgroundColor: bgColor,
        contrastRatio: textContrast,
        ...checkWCAG(textContrast)
      },
      {
        colorKey: 'border',
        textColor: colors.border,
        backgroundColor: bgColor,
        contrastRatio: borderContrast,
        ...checkWCAG(borderContrast)
      }
    ]
  }, [colors])

  const allPass = contrastResults.every((r) => r.wcagAA)
  const allPassAAA = contrastResults.every((r) => r.wcagAAA)
  const failures = contrastResults.filter((r) => !r.wcagAA)
  const swatchRules = contrastResults
    .map(
      (result) => `
        .contrast-swatch-${result.colorKey} {
          background-color: ${result.textColor};
        }
      `
    )
    .join('\n')

  return (
    <div className="space-y-4">
      {/* Summary */}
      <div
        className={`p-4 rounded-lg border-2 flex items-start gap-3 ${
          allPassAAA
            ? 'bg-green-50 border-green-200 dark:bg-green-900/20 dark:border-green-700'
            : allPass
              ? 'bg-blue-50 border-blue-200 dark:bg-blue-900/20 dark:border-blue-700'
              : 'bg-red-50 border-red-200 dark:bg-red-900/20 dark:border-red-700'
        }`}
      >
        {allPassAAA ? (
          <CheckCircle className="h-5 w-5 text-green-600 dark:text-green-400 flex-shrink-0 mt-0.5" />
        ) : allPass ? (
          <CheckCircle className="h-5 w-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
        ) : (
          <AlertTriangle className="h-5 w-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
        )}
        <div className="flex-1">
          <h4 className="font-semibold text-sm mb-1">
            {allPassAAA
              ? 'WCAG AAA - Enhanced Contrast'
              : allPass
                ? 'WCAG AA - Minimum Contrast'
                : `${failures.length} Contrast Issues`}
          </h4>
          <p className="text-xs text-gray-600 dark:text-gray-400">
            {allPassAAA
              ? 'All colors meet enhanced accessibility standards'
              : allPass
                ? 'All colors meet minimum accessibility standards'
                : `${failures.length} color combination(s) fail WCAG AA standard (4.5:1)`}
          </p>
        </div>
      </div>

      {/* Details */}
      {showDetails && (
        <div className="space-y-2">
          {contrastResults.map((result) => (
            <div
              key={result.colorKey}
              className="p-3 rounded-lg border border-gray-200 dark:border-gray-700"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex items-center gap-2 flex-1">
                  <div
                    className={`w-8 h-8 rounded border-2 border-gray-300 dark:border-gray-600 contrast-swatch-${result.colorKey}`}
                  />
                  <div className="flex-1">
                    <p className="text-sm font-medium capitalize text-gray-900 dark:text-white">
                      {result.colorKey} on Background
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      Ratio: {result.contrastRatio.toFixed(2)}:1
                    </p>
                  </div>
                </div>

                <div className="flex gap-1">
                  {result.wcagAA && (
                    <div
                      title="Meets WCAG AA (4.5:1)"
                      className="px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 text-xs font-medium rounded"
                    >
                      AA
                    </div>
                  )}
                  {result.wcagAAA && (
                    <div
                      title="Meets WCAG AAA (7:1)"
                      className="px-2 py-1 bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300 text-xs font-medium rounded"
                    >
                      AAA
                    </div>
                  )}
                  {!result.wcagAA && (
                    <div
                      title="Does not meet WCAG AA"
                      className="px-2 py-1 bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-300 text-xs font-medium rounded flex items-center gap-1"
                    >
                      <AlertCircle className="h-3 w-3" />
                      Fail
                    </div>
                  )}
                </div>
              </div>

              {/* Detailed breakdown */}
              <div className="mt-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                <div className="grid grid-cols-2 gap-2 text-xs">
                  <div className="flex items-center gap-1">
                    <div
                      className={`w-3 h-3 rounded-full ${
                        result.largeTextAA ? 'bg-green-500' : 'bg-gray-300'
                      }`}
                    />
                    <span className="text-gray-600 dark:text-gray-400">
                      Large AA {result.largeTextAA ? '✓' : '✗'}
                    </span>
                  </div>
                  <div className="flex items-center gap-1">
                    <div
                      className={`w-3 h-3 rounded-full ${
                        result.largeTextAAA ? 'bg-green-500' : 'bg-gray-300'
                      }`}
                    />
                    <span className="text-gray-600 dark:text-gray-400">
                      Large AAA {result.largeTextAAA ? '✓' : '✗'}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Recommendations */}
      {failures.length > 0 && (
        <div className="p-3 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-700 rounded-lg">
          <p className="text-xs font-medium text-yellow-800 dark:text-yellow-300 mb-2">
            💡 Recommendations:
          </p>
          <ul className="text-xs text-yellow-700 dark:text-yellow-400 space-y-1 list-disc list-inside">
            {failures.map((f) => (
              <li key={f.colorKey}>
                Increase contrast for {f.colorKey} (current: {f.contrastRatio.toFixed(2)}:1, need: 4.5:1)
              </li>
            ))}
          </ul>
        </div>
      )}
      <style jsx>{`
        ${swatchRules}
      `}</style>
    </div>
  )
}
