'use client'

import React, { useState, useCallback, useEffect, useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { Search, Star, TrendingUp } from 'lucide-react'
import type { ThemeComponent, TypographySettings } from '@/types/theme'
import { googleFontsService, type GoogleFont } from '@/services/googleFontsService'
import { useUpdateTypography, useThemeTypography } from '@/stores'
import { getLocaleFromPath } from '@/lib/translations'

interface FontPickerProps {
  component?: ThemeComponent
  onChange?: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
  useStore?: boolean
}

interface FontConfig {
  fontFamily?: string
  fontWeight?: number
  fontStyle?: 'normal' | 'italic'
  fontWeights?: number[]
  fontSize?: number
  lineHeight?: number
  letterSpacing?: number
  preloadFonts?: boolean
}

type FontCategory = 'all' | 'serif' | 'sans-serif' | 'monospace' | 'display' | 'handwriting'

/**
 * FontPicker Component
 * Allows users to select and preview Google Fonts for the theme
 * Supports both local state (component prop) and themeBuilderStore (useStore prop)
 */
export function FontPicker({
  component,
  onChange,
  onPreview,
  className = '',
  useStore = false,
}: FontPickerProps): JSX.Element {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const isRTL = locale === 'ar'
  const tx = (en: string, ar: string) => {
    if (!isRTL) return en
        return /[ÃØÙðâ]/.test(ar) ? en : ar
  }

  const storeTypography = useThemeTypography()
  const updateTypography = useUpdateTypography()

  const initialFontConfig: FontConfig = useMemo(() => {
    if (useStore && storeTypography) {
      return {
        fontFamily: storeTypography.fontFamily,
        fontWeight: storeTypography.fontWeights?.[0] || 400,
        fontStyle: storeTypography.fontStyle as 'normal' | 'italic' || 'normal',
        fontWeights: storeTypography.fontWeights || [400],
        fontSize: storeTypography.fontSize || 16,
        lineHeight: storeTypography.lineHeight || 1.5,
        letterSpacing: 0,
        preloadFonts: true,
      }
    }
    return (component?.config as unknown as FontConfig) || {
      fontWeight: 400,
      fontStyle: 'normal',
      fontSize: 16,
      lineHeight: 1.6,
      letterSpacing: 0,
      fontWeights: [400],
    }
  }, [useStore, storeTypography, component])

  const [fontConfig, setFontConfig] = useState<FontConfig>(initialFontConfig)

  const [category, setCategory] = useState<FontCategory>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [fonts, setFonts] = useState<GoogleFont[]>([])
  const [selectedFont, setSelectedFont] = useState<GoogleFont | null>(null)
  const [previewMode, setPreviewMode] = useState(false)
  const [recentFonts, setRecentFonts] = useState<GoogleFont[]>([])
  const [fontLoaded, setFontLoaded] = useState(false)

  // Load fonts on mount and when category/search changes
  useEffect(() => {
    let filtered: GoogleFont[] = []

    if (searchQuery) {
      filtered = googleFontsService.searchFonts(searchQuery)
    } else if (category === 'all') {
      filtered = googleFontsService.getPopularFonts()
    } else {
      filtered = googleFontsService.getFontsByCategory(category as any)
    }

    setFonts(filtered)
  }, [category, searchQuery])

  // Load selected font and recent fonts
  useEffect(() => {
    const recent = googleFontsService.getRecentFonts()
    setRecentFonts(recent)

    if (fontConfig.fontFamily) {
      const font = googleFontsService.getFont(fontConfig.fontFamily)
      if (font) {
        setSelectedFont(font)
        // Load the font
        googleFontsService.loadFont(font.family, fontConfig.fontWeights || [400]).then(() => {
          setFontLoaded(true)
        })
      }
    }
  }, [fontConfig.fontFamily])

  const handleSelectFont = useCallback(
    (font: GoogleFont) => {
      const newConfig: FontConfig = {
        ...fontConfig,
        fontFamily: font.id,
        fontWeights: [400, 700],
      }

      setFontConfig(newConfig)
      setSelectedFont(font)
      googleFontsService.addRecentFont(font)

      if (useStore) {
        updateTypography({
          fontFamily: font.id,
          fontWeights: [400, 700],
        })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
        onPreview?.(updated)
      }
    },
    [fontConfig, component, onChange, onPreview, useStore, updateTypography]
  )

  const handleWeightChange = useCallback(
    (weight: number) => {
      const weights = fontConfig.fontWeights || [400]
      const newWeights = weights.includes(weight)
        ? weights.filter((w) => w !== weight)
        : [...weights, weight]

      const newConfig: FontConfig = {
        ...fontConfig,
        fontWeights: newWeights.sort((a, b) => a - b),
        fontWeight: weight,
      }

      setFontConfig(newConfig)

      if (useStore) {
        updateTypography({
          fontWeights: newWeights.sort((a, b) => a - b),
        })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [fontConfig, component, onChange, useStore, updateTypography]
  )

  const handleStyleChange = useCallback(
    (style: 'normal' | 'italic') => {
      const newConfig: FontConfig = {
        ...fontConfig,
        fontStyle: style,
      }

      setFontConfig(newConfig)

      if (useStore) {
        updateTypography({ fontStyle: style })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [fontConfig, component, onChange, useStore, updateTypography]
  )

  const handleFontSizeChange = useCallback(
    (size: number) => {
      const newConfig: FontConfig = {
        ...fontConfig,
        fontSize: size,
      }

      setFontConfig(newConfig)

      if (useStore) {
        updateTypography({ fontSize: size })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [fontConfig, component, onChange, useStore, updateTypography]
  )

  const handleLineHeightChange = useCallback(
    (height: number) => {
      const newConfig: FontConfig = {
        ...fontConfig,
        lineHeight: height,
      }

      setFontConfig(newConfig)

      if (useStore) {
        updateTypography({ lineHeight: height })
      } else if (component && onChange) {
        const updated: ThemeComponent = {
          ...component,
          config: newConfig,
        }
        onChange(updated)
      }
    },
    [fontConfig, component, onChange, useStore, updateTypography]
  )

  const previewText =
    "Pack my box with five dozen liquor jugs. The quick brown fox jumps over the lazy dog. Sphinx of black quartz, judge my vow."
  const selectedFontFamily = selectedFont ? selectedFont.family.replace(/'/g, "\\'") : ''
  const previewFontFamily = selectedFont ? `'${selectedFontFamily}', sans-serif` : 'system-ui'
  const previewFontSize = fontConfig.fontSize || 16
  const previewFontStyle = fontConfig.fontStyle || 'normal'
  const previewFontWeight = fontConfig.fontWeight || 400
  const previewLineHeight = fontConfig.lineHeight || 1.6
  const previewLetterSpacing = fontConfig.letterSpacing || 0
  const getSafeFontId = (id: string) => id.replace(/[^a-zA-Z0-9_-]/g, '-')
  const fontFamilyRules = fonts
    .map(
      (font) => `
        .font-family-${getSafeFontId(font.id)} {
          font-family: '${font.family.replace(/'/g, "\\'")}', sans-serif;
        }
      `
    )
    .join('\n')

  // Preview Mode
  if (previewMode) {
    return (
      <div dir={isRTL ? 'rtl' : 'ltr'} className={`w-full space-y-4 ${className} ${isRTL ? 'text-right' : ''}`}>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{tx('Font Preview', 'Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â§Ãƒâ„¢Ã…Â Ãƒâ„¢Ã¢â‚¬Â ÃƒËœÃ‚Â© ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}</h3>
          <button
            onClick={() => setPreviewMode(false)}
            className="px-3 py-2 text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition"
          >
            {tx('Back to Settings', 'ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â¹Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â¯ÃƒËœÃ‚Â© ÃƒËœÃ‚Â¥Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã¢â‚¬Â° ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â¥ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â¯ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¯ÃƒËœÃ‚Â§ÃƒËœÃ‚Âª')}
          </button>
        </div>

        <div className="space-y-6">
          {/* Font Info */}
          {selectedFont && (
            <div className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
              <h4 className="font-semibold text-gray-900 dark:text-white mb-2">{selectedFont.family}</h4>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                {selectedFont.category} - {selectedFont.weights.length} weights
              </p>
              <div className="space-y-2 text-xs text-gray-600 dark:text-gray-400">
                <p>
                  <strong>{tx('Selected Weight:', 'ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â²Ãƒâ„¢Ã¢â‚¬Â  ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â­ÃƒËœÃ‚Â¯ÃƒËœÃ‚Â¯:')}</strong> {fontConfig.fontWeight || 400}
                </p>
                <p>
                  <strong>{tx('Font Size:', 'ÃƒËœÃ‚Â­ÃƒËœÃ‚Â¬Ãƒâ„¢Ã¢â‚¬Â¦ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·:')}</strong> {fontConfig.fontSize || 16}px
                </p>
                <p>
                  <strong>{tx('Line Height:', 'ÃƒËœÃ‚Â§ÃƒËœÃ‚Â±ÃƒËœÃ‚ÂªÃƒâ„¢Ã‚ÂÃƒËœÃ‚Â§ÃƒËœÃ‚Â¹ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â³ÃƒËœÃ‚Â·ÃƒËœÃ‚Â±:')}</strong> {fontConfig.lineHeight || 1.6}
                </p>
              </div>
            </div>
          )}

          {/* Preview Text */}
            <div className="space-y-3">
            <h4 className="font-semibold text-gray-900 dark:text-white">{tx('Preview', 'Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â§Ãƒâ„¢Ã…Â Ãƒâ„¢Ã¢â‚¬Â ÃƒËœÃ‚Â©')}</h4>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
              <p className="font-preview-text text-gray-900 dark:text-white">
                {previewText}
              </p>
            </div>
          </div>

          {/* Heading Preview */}
          <div className="space-y-3">
            <h4 className="font-semibold text-gray-900 dark:text-white">{tx('Heading (24px, Bold)', 'ÃƒËœÃ‚Â¹Ãƒâ„¢Ã¢â‚¬Â Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Â  (24pxÃƒËœÃ…â€™ ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â±Ãƒâ„¢Ã…Â ÃƒËœÃ‚Â¶)')}</h4>
            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
              <h1
                className="font-preview-text font-preview-heading text-gray-900 dark:text-white"
              >
                {tx('This is a heading', 'Ãƒâ„¢Ã¢â‚¬Â¡ÃƒËœÃ‚Â°ÃƒËœÃ‚Â§ ÃƒËœÃ‚Â¹Ãƒâ„¢Ã¢â‚¬Â Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Â  ÃƒËœÃ‚ÂªÃƒËœÃ‚Â¬ÃƒËœÃ‚Â±Ãƒâ„¢Ã…Â ÃƒËœÃ‚Â¨Ãƒâ„¢Ã…Â ')}
              </h1>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // Editor Mode
  return (
    <div dir={isRTL ? 'rtl' : 'ltr'} className={`w-full space-y-6 ${className} ${isRTL ? 'text-right' : ''}`}>
      {/* Header */}
      <div>
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-2">{tx('Font Picker', 'ÃƒËœÃ‚Â§ÃƒËœÃ‚Â®ÃƒËœÃ‚ÂªÃƒâ„¢Ã…Â ÃƒËœÃ‚Â§ÃƒËœÃ‚Â± ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}</h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          {tx('Choose from popular Google Fonts and customize weight, size, and styling.', 'ÃƒËœÃ‚Â§ÃƒËœÃ‚Â®ÃƒËœÃ‚ÂªÃƒËœÃ‚Â± Ãƒâ„¢Ã¢â‚¬Â¦Ãƒâ„¢Ã¢â‚¬Â  ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â· Google ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â´ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¦ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â© Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â®ÃƒËœÃ‚ÂµÃƒËœÃ‚Âµ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â²Ãƒâ„¢Ã¢â‚¬Â  Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â­ÃƒËœÃ‚Â¬Ãƒâ„¢Ã¢â‚¬Â¦ Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚ÂªÃƒâ„¢Ã¢â‚¬Â ÃƒËœÃ‚Â³Ãƒâ„¢Ã…Â Ãƒâ„¢Ã¢â‚¬Å¡.')}
        </p>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className={`absolute top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400 ${isRTL ? 'right-3' : 'left-3'}`} />
        <input
          type="text"
          placeholder={tx('Search fonts...', 'ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¨ÃƒËœÃ‚Â­ÃƒËœÃ‚Â« ÃƒËœÃ‚Â¹Ãƒâ„¢Ã¢â‚¬Â  ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â·...')}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className={`w-full py-2 border border-gray-200 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white text-sm ${isRTL ? 'pr-10 pl-4' : 'pl-10 pr-4'}`}
        />
      </div>

      {/* Category Tabs */}
      <div className="flex gap-2 overflow-x-auto pb-2">
            {(['all', 'sans-serif', 'serif', 'monospace'] as const).map((cat) => (
          <button
            key={cat}
            onClick={() => setCategory(cat)}
            className={`px-3 py-1.5 rounded-lg whitespace-nowrap text-sm font-medium transition ${
              category === cat
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
            }`}
          >
            {cat === 'all'
              ? tx('All', 'ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã†â€™Ãƒâ„¢Ã¢â‚¬Å¾')
              : cat === 'sans-serif'
                ? tx('Sans Serif', 'ÃƒËœÃ‚Â¨ÃƒËœÃ‚Â¯Ãƒâ„¢Ã‹â€ Ãƒâ„¢Ã¢â‚¬Â  ÃƒËœÃ‚Â²Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¦ÃƒËœÃ‚Â¯')
                : cat === 'serif'
                  ? tx('Serif', 'Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â²ÃƒËœÃ‚Â®ÃƒËœÃ‚Â±Ãƒâ„¢Ã‚Â')
                  : tx('Monospace', 'ÃƒËœÃ‚Â£ÃƒËœÃ‚Â­ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¯Ãƒâ„¢Ã…Â  ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â³ÃƒËœÃ‚Â§Ãƒâ„¢Ã‚ÂÃƒËœÃ‚Â©')}
          </button>
        ))}
      </div>

      {/* Recent Fonts */}
      {!searchQuery && recentFonts.length > 0 && (
        <div className="space-y-2">
          <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300">{tx('Recent', 'ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â£ÃƒËœÃ‚Â®Ãƒâ„¢Ã…Â ÃƒËœÃ‚Â±ÃƒËœÃ‚Â©')}</h3>
          <div className="grid grid-cols-2 gap-2">
            {recentFonts.map((font) => (
              <button
                key={font.id}
                onClick={() => handleSelectFont(font)}
                className={`p-3 rounded-lg text-left transition ${
                  selectedFont?.id === font.id
                    ? 'bg-blue-100 dark:bg-blue-900 border-2 border-blue-600 dark:border-blue-400'
                    : 'bg-gray-100 dark:bg-gray-800 border-2 border-transparent hover:border-gray-300 dark:hover:border-gray-700'
                }`}
              >
                <div className="text-sm font-semibold text-gray-900 dark:text-white">{font.family}</div>
                <div className="text-xs text-gray-600 dark:text-gray-400">{font.category}</div>
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Font List */}
      <div className="space-y-2">
        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300">
          {searchQuery ? tx('Search Results', 'Ãƒâ„¢Ã¢â‚¬Â ÃƒËœÃ‚ÂªÃƒËœÃ‚Â§ÃƒËœÃ‚Â¦ÃƒËœÃ‚Â¬ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â¨ÃƒËœÃ‚Â­ÃƒËœÃ‚Â«') : category === 'all' ? tx('Popular Fonts', 'ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â· ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â´ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¦ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â©') : `${category} ${tx('Fonts', 'ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â·')}`}
          <span className={`text-gray-500 dark:text-gray-500 ${isRTL ? 'mr-2' : 'ml-2'}`}>({fonts.length})</span>
        </h3>

        {fonts.length === 0 ? (
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">
            <p className="text-sm">{tx('No fonts found. Try a different search.', 'Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã¢â‚¬Â¦ Ãƒâ„¢Ã…Â ÃƒËœÃ‚ÂªÃƒâ„¢Ã¢â‚¬Â¦ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â«Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â± ÃƒËœÃ‚Â¹Ãƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã¢â‚¬Â° ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â·. ÃƒËœÃ‚Â¬ÃƒËœÃ‚Â±Ãƒâ„¢Ã¢â‚¬ËœÃƒËœÃ‚Â¨ ÃƒËœÃ‚Â¨ÃƒËœÃ‚Â­ÃƒËœÃ‚Â«Ãƒâ„¢Ã¢â‚¬Â¹ÃƒËœÃ‚Â§ Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â®ÃƒËœÃ‚ÂªÃƒâ„¢Ã¢â‚¬Å¾Ãƒâ„¢Ã‚ÂÃƒâ„¢Ã¢â‚¬Â¹ÃƒËœÃ‚Â§.')}</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-2 max-h-64 overflow-y-auto">
            {fonts.map((font) => (
              <button
                key={font.id}
                onClick={() => handleSelectFont(font)}
                className={`p-3 rounded-lg text-left transition border-2 ${
                  selectedFont?.id === font.id
                    ? 'bg-blue-100 dark:bg-blue-900 border-blue-600 dark:border-blue-400'
                    : 'bg-gray-50 dark:bg-gray-800 border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                }`}
              >
                <div className={`font-family-${getSafeFontId(font.id)}`}>
                  <div className="font-semibold text-gray-900 dark:text-white">{font.family}</div>
                </div>
                <div className="flex items-center justify-between mt-1">
                  <span className="text-xs text-gray-600 dark:text-gray-400">{font.category}</span>
                  {font.popularity && font.popularity <= 3 && (
                    <div className="flex items-center gap-1 text-yellow-500">
                      <TrendingUp className="h-3 w-3" />
                      <span className="text-xs">{tx('Popular', 'ÃƒËœÃ‚Â´ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¦ÃƒËœÃ‚Â¹')}</span>
                    </div>
                  )}
                </div>
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Font Configuration */}
      {selectedFont && (
        <div className="bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4">
          <h3 className="font-semibold text-gray-900 dark:text-white">{tx('Font Configuration', 'ÃƒËœÃ‚Â¥ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â¯ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¯ÃƒËœÃ‚Â§ÃƒËœÃ‚Âª ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}</h3>

          {/* Font Weight Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Font Weights', 'ÃƒËœÃ‚Â£Ãƒâ„¢Ã‹â€ ÃƒËœÃ‚Â²ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Â  ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}
            </label>
            <div className="flex flex-wrap gap-2">
              {selectedFont.weights.map((weight) => (
                <button
                  key={weight}
                  onClick={() => handleWeightChange(weight)}
                  className={`px-3 py-1.5 rounded-lg text-sm font-medium transition ${
                    fontConfig.fontWeights?.includes(weight)
                      ? 'bg-blue-600 text-white'
                      : 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
                  }`}
                >
                  {weight}
                </button>
              ))}
            </div>
          </div>

          {/* Font Style */}
          {selectedFont.italics && (
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {tx('Font Style', 'Ãƒâ„¢Ã¢â‚¬Â Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â· ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}
              </label>
              <div className="flex gap-2">
                {(['normal', 'italic'] as const).map((style) => (
                  <button
                    key={style}
                    onClick={() => handleStyleChange(style)}
                    className={`px-4 py-2 rounded-lg text-sm font-medium transition ${
                      fontConfig.fontStyle === style
                        ? 'bg-blue-600 text-white'
                        : 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:border-gray-300'
                    }`}
                  >
                    {style === 'normal' ? tx('Normal', 'ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¯Ãƒâ„¢Ã…Â ') : tx('Italic', 'Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â§ÃƒËœÃ‚Â¦Ãƒâ„¢Ã¢â‚¬Å¾')}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Font Size */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Font Size', 'ÃƒËœÃ‚Â­ÃƒËœÃ‚Â¬Ãƒâ„¢Ã¢â‚¬Â¦ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}: {fontConfig.fontSize || 16}px
            </label>
            <input
              type="range"
              aria-label="Font Size"
              min="10"
              max="32"
              value={fontConfig.fontSize || 16}
              onChange={(e) => handleFontSizeChange(parseInt(e.target.value))}
              className="w-full"
            />
          </div>

          {/* Line Height */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {tx('Line Height', 'ÃƒËœÃ‚Â§ÃƒËœÃ‚Â±ÃƒËœÃ‚ÂªÃƒâ„¢Ã‚ÂÃƒËœÃ‚Â§ÃƒËœÃ‚Â¹ ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â³ÃƒËœÃ‚Â·ÃƒËœÃ‚Â±')}: {(fontConfig.lineHeight || 1.6).toFixed(1)}
            </label>
            <input
              type="range"
              aria-label="Line Height"
              min="1"
              max="2"
              step="0.1"
              value={fontConfig.lineHeight || 1.6}
              onChange={(e) => handleLineHeightChange(parseFloat(e.target.value))}
              className="w-full"
            />
          </div>
        </div>
      )}

      {/* Preview Button */}
      <div className="flex gap-3">
        <button
          onClick={() => setPreviewMode(true)}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition"
        >
          {tx('Preview Font', 'Ãƒâ„¢Ã¢â‚¬Â¦ÃƒËœÃ‚Â¹ÃƒËœÃ‚Â§Ãƒâ„¢Ã…Â Ãƒâ„¢Ã¢â‚¬Â ÃƒËœÃ‚Â© ÃƒËœÃ‚Â§Ãƒâ„¢Ã¢â‚¬Å¾ÃƒËœÃ‚Â®ÃƒËœÃ‚Â·')}
        </button>
      </div>

      {/* Live Preview */}
      {selectedFont && (
        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <p className="font-preview-text text-gray-900 dark:text-white">
            The quick brown fox jumps over the lazy dog
          </p>
        </div>
      )}

      <style jsx>{`
        .font-preview-text {
          font-family: ${previewFontFamily};
          font-size: ${previewFontSize}px;
          font-style: ${previewFontStyle};
          font-weight: ${previewFontWeight};
          line-height: ${previewLineHeight};
          letter-spacing: ${previewLetterSpacing}px;
        }

        .font-preview-heading {
          font-size: 24px;
          font-weight: 700;
        }

        ${fontFamilyRules}
      `}</style>
    </div>
  )
}
