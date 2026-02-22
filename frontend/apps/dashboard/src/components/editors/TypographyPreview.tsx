'use client'

import { useMemo } from 'react'
import { TypographySettings } from '@/types/theme'

interface TypographyPreviewProps {
  typography: TypographySettings
  textColor?: string
}

interface FontScaleItem {
  name: string
  size: number
  lineHeight: number
  weight: number
  sampleText: string
}

export function TypographyPreview({ typography, textColor = '#000000' }: TypographyPreviewProps) {
  const fontScale = useMemo<FontScaleItem[]>(() => {
    const base = typography.baseFontSize

    return [
      {
        name: 'H1 - Heading 1',
        size: Math.round(base * 2.5), // 2.5x base
        lineHeight: 1.2,
        weight: 700,
        sampleText: 'The Quick Brown Fox'
      },
      {
        name: 'H2 - Heading 2',
        size: Math.round(base * 2), // 2x base
        lineHeight: 1.3,
        weight: 700,
        sampleText: 'The Quick Brown Fox'
      },
      {
        name: 'H3 - Heading 3',
        size: Math.round(base * 1.5), // 1.5x base
        lineHeight: 1.3,
        weight: 600,
        sampleText: 'The Quick Brown Fox'
      },
      {
        name: 'Body Text',
        size: base,
        lineHeight: typography.lineHeight,
        weight: 400,
        sampleText:
          'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.'
      },
      {
        name: 'Small Text',
        size: Math.round(base * 0.875), // 0.875x base
        lineHeight: 1.4,
        weight: 400,
        sampleText: 'This is smaller text, often used for captions or helper text'
      },
      {
        name: 'Extra Small',
        size: Math.round(base * 0.75), // 0.75x base
        lineHeight: 1.4,
        weight: 400,
        sampleText: 'Extra small text for labels or annotations'
      }
    ]
  }, [typography.baseFontSize, typography.lineHeight])

  const sampleRules = fontScale
    .map(
      (item, index) => `
        .typography-sample-${index} {
          font-size: ${item.size}px;
          line-height: ${item.lineHeight};
          font-weight: ${item.weight};
        }
      `
    )
    .join('\n')

  return (
    <div className="typography-preview space-y-6">
      <div>
        <h3 className="font-semibold text-sm text-gray-900 dark:text-white mb-4">
          Typography Scale Preview
        </h3>

        {/* Font Settings Summary */}
        <div className="grid grid-cols-2 gap-3 mb-6 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
          <div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mb-1">Font Family</p>
            <p className="typography-font-family text-sm font-medium text-gray-900 dark:text-white truncate">
              {typography.fontFamily.split(',')[0]}
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mb-1">Base Size</p>
            <p className="text-sm font-medium text-gray-900 dark:text-white">
              {typography.baseFontSize}px
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mb-1">Line Height</p>
            <p className="text-sm font-medium text-gray-900 dark:text-white">
              {typography.lineHeight}
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mb-1">Border Radius</p>
            <p className="text-sm font-medium text-gray-900 dark:text-white">
              {typography.borderRadius}px
            </p>
          </div>
        </div>

        {/* Typography Scale */}
        <div className="space-y-4">
          {fontScale.map((item, index) => (
            <div
              key={`${item.name}-${index}`}
              className="p-4 border border-gray-200 dark:border-gray-700 rounded-lg"
            >
              {/* Label */}
              <div className="flex items-baseline justify-between mb-2">
                <h4 className="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">
                  {item.name}
                </h4>
                <span className="text-xs text-gray-500 dark:text-gray-500">
                  {item.size}px / {item.weight} weight
                </span>
              </div>

              {/* Sample Text */}
              <p className={`typography-sample typography-sample-${index} dark:text-white`}>
                {item.sampleText}
              </p>

              {/* Metrics */}
              <div className="mt-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                <div className="text-xs text-gray-500 dark:text-gray-400 space-y-1">
                  <div className="flex justify-between">
                    <span>Line Height:</span>
                    <span className="font-mono">{item.lineHeight}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Letter Spacing:</span>
                    <span className="font-mono">0px</span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Font Pairing Suggestions */}
        <div className="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700 rounded-lg">
          <h4 className="text-xs font-medium text-blue-900 dark:text-blue-300 mb-2">
            💡 Font Pairing Tips:
          </h4>
          <ul className="text-xs text-blue-800 dark:text-blue-400 space-y-1 list-disc list-inside">
            <li>Pair serif fonts with sans-serif for contrast</li>
            <li>Use max 2-3 different fonts in your design</li>
            <li>Ensure sufficient contrast with background color</li>
            <li>Test different sizes across breakpoints</li>
            <li>Maintain consistent line heights for readability (1.4 - 1.8)</li>
          </ul>
        </div>
      </div>
      <style jsx>{`
        .typography-font-family,
        .typography-sample {
          font-family: ${typography.fontFamily};
        }

        .typography-sample {
          color: ${textColor};
        }

        ${sampleRules}
      `}</style>
    </div>
  )
}
