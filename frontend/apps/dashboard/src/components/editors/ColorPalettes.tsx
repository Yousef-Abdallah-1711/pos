'use client'

import { useState } from 'react'
import { usePathname } from 'next/navigation'
import { Palette, Copy, Check } from 'lucide-react'
import { GlobalColors } from '@/types/theme'
import { createTranslator, getLocaleFromPath } from '@/lib/translations'

interface ColorPalette {
  id: string
  name: string
  category: 'professional' | 'creative' | 'minimal' | 'bold'
  colors: GlobalColors
}

interface ColorPalettesProps {
  onApply: (colors: Partial<GlobalColors>) => void
}

const PRESET_PALETTES: ColorPalette[] = [
  {
    id: 'professional-blue',
    name: 'Professional Blue',
    category: 'professional',
    colors: {
      primary: '#0066cc',
      secondary: '#4a90e2',
      accent: '#64b5f6',
      background: '#ffffff',
      text: '#2d3748',
      border: '#e2e8f0',
      shadow: '#cbd5e0'
    }
  },
  {
    id: 'elegant-dark',
    name: 'Elegant Dark',
    category: 'professional',
    colors: {
      primary: '#1a1a2e',
      secondary: '#16213e',
      accent: '#0f3460',
      background: '#e94560',
      text: '#ffffff',
      border: '#1f2937',
      shadow: '#4b5563'
    }
  },
  {
    id: 'fresh-green',
    name: 'Fresh Green',
    category: 'creative',
    colors: {
      primary: '#2ecc71',
      secondary: '#27ae60',
      accent: '#1abc9c',
      background: '#f8f9fa',
      text: '#2c3e50',
      border: '#ecf0f1',
      shadow: '#bdc3c7'
    }
  },
  {
    id: 'vibrant-orange',
    name: 'Vibrant Orange',
    category: 'bold',
    colors: {
      primary: '#ff6b35',
      secondary: '#f7931e',
      accent: '#ffc107',
      background: '#ffffff',
      text: '#212529',
      border: '#e9ecef',
      shadow: '#dee2e6'
    }
  },
  {
    id: 'minimal-white',
    name: 'Minimal White',
    category: 'minimal',
    colors: {
      primary: '#000000',
      secondary: '#333333',
      accent: '#666666',
      background: '#ffffff',
      text: '#000000',
      border: '#e0e0e0',
      shadow: '#cccccc'
    }
  },
  {
    id: 'luxury-purple',
    name: 'Luxury Purple',
    category: 'professional',
    colors: {
      primary: '#6a0dad',
      secondary: '#9b59b6',
      accent: '#bb8fce',
      background: '#f8f9fa',
      text: '#2c3e50',
      border: '#ecf0f1',
      shadow: '#bdc3c7'
    }
  },
  {
    id: 'ocean-blue',
    name: 'Ocean Blue',
    category: 'creative',
    colors: {
      primary: '#006994',
      secondary: '#0891b2',
      accent: '#06b6d4',
      background: '#f0f9ff',
      text: '#0c4a6e',
      border: '#e0f2fe',
      shadow: '#bae6fd'
    }
  },
  {
    id: 'warm-sunset',
    name: 'Warm Sunset',
    category: 'bold',
    colors: {
      primary: '#e63946',
      secondary: '#f77f00',
      accent: '#fcbf49',
      background: '#fffbf0',
      text: '#1d3557',
      border: '#f1faee',
      shadow: '#a8dadc'
    }
  }
]

export function ColorPalettes({ onApply }: ColorPalettesProps) {
  const pathname = usePathname()
  const locale = getLocaleFromPath(pathname)
  const t = createTranslator(locale)
  const isRTL = locale === 'ar'
  const [copiedPaletteId, setCopiedPaletteId] = useState<string | null>(null)
  const [selectedCategory, setSelectedCategory] = useState<'all' | ColorPalette['category']>('all')
  const getSafePaletteId = (id: string) => id.replace(/[^a-zA-Z0-9_-]/g, '-')

  const filteredPalettes =
    selectedCategory === 'all'
      ? PRESET_PALETTES
      : PRESET_PALETTES.filter((p) => p.category === selectedCategory)

  const categories: Array<{ value: 'all' | ColorPalette['category']; label: string }> = [
    { value: 'all', label: t('themeBuilderEditor.palettes.all', 'All Palettes') },
    { value: 'professional', label: t('themeBuilderEditor.palettes.professional', 'Professional') },
    { value: 'creative', label: t('themeBuilderEditor.palettes.creative', 'Creative') },
    { value: 'minimal', label: t('themeBuilderEditor.palettes.minimal', 'Minimal') },
    { value: 'bold', label: t('themeBuilderEditor.palettes.bold', 'Bold') }
  ]

  const handleCopyPalette = (paletteId: string) => {
    setCopiedPaletteId(paletteId)
    setTimeout(() => setCopiedPaletteId(null), 2000)
  }

  const paletteColorRules = PRESET_PALETTES.map((palette) => {
    const safeId = getSafePaletteId(palette.id)
    const stripRules = Object.values(palette.colors)
      .map(
        (color, idx) => `
          .palette-strip-${safeId}-${idx} {
            background-color: ${color};
          }
        `
      )
      .join('\n')

    const detailRules = Object.entries(palette.colors)
      .map(
        ([key, color]) => `
          .palette-detail-${safeId}-${key} {
            background-color: ${color};
          }
        `
      )
      .join('\n')

    return `${stripRules}\n${detailRules}`
  }).join('\n')

  return (
    <div className={`space-y-4 ${isRTL ? 'text-right' : ''}`}>
      <div>
        <h3 className="font-semibold text-sm text-gray-900 dark:text-white mb-3">
          {t('themeBuilderEditor.colorPalettes', 'Color Palettes')}
        </h3>

        {/* Category Filter */}
        <div className="flex flex-wrap gap-2 mb-4">
          {categories.map((category) => (
            <button
              key={category.value}
              onClick={() => setSelectedCategory(category.value)}
              className={`px-3 py-1 rounded-full text-xs font-medium transition ${
                selectedCategory === category.value
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              {category.label}
            </button>
          ))}
        </div>

        {/* Palettes Grid */}
        <div className="grid grid-cols-1 gap-3">
          {filteredPalettes.map((palette) => (
            <div
              key={palette.id}
              className="rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden hover:shadow-md transition"
            >
              {/* Palette Preview */}
              <div className="flex h-20">
                {Object.values(palette.colors).map((color, idx) => (
                  <div
                    key={idx}
                    className={`flex-1 palette-strip-${getSafePaletteId(palette.id)}-${idx}`}
                    title={color}
                  />
                ))}
              </div>

              {/* Palette Info & Actions */}
              <div className="p-3 bg-white dark:bg-gray-800">
                <div className={`flex items-center justify-between mb-2 ${isRTL ? 'flex-row-reverse' : ''}`}>
                  <h4 className="text-sm font-medium text-gray-900 dark:text-white">
                    {palette.name}
                  </h4>
                  <span className="text-xs px-2 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 capitalize">
                    {t(`themeBuilderEditor.palettes.${palette.category}`, palette.category)}
                  </span>
                </div>

                {/* Buttons */}
                <div className="flex gap-2">
                  <button
                    onClick={() => onApply(palette.colors)}
                    className="flex-1 px-3 py-2 bg-primary-600 hover:bg-primary-700 text-white text-xs font-medium rounded transition flex items-center justify-center gap-1"
                  >
                    <Palette className="h-3 w-3" />
                    {t('themeBuilderEditor.apply', 'Apply')}
                  </button>
                  <button
                    onClick={() => handleCopyPalette(palette.id)}
                    className="px-3 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 text-xs font-medium rounded transition flex items-center justify-center gap-1"
                    title={t('themeBuilderEditor.copyPaletteColors', 'Copy palette colors')}
                  >
                    {copiedPaletteId === palette.id ? (
                      <Check className="h-3 w-3" />
                    ) : (
                      <Copy className="h-3 w-3" />
                    )}
                  </button>
                </div>

                {/* Color Details */}
                <div className="mt-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                  <div className="grid grid-cols-7 gap-1">
                    {Object.entries(palette.colors).map(([key, color]) => (
                      <div
                        key={key}
                        className="flex flex-col items-center"
                        title={`${key}: ${color}`}
                      >
                        <div
                          className={`w-6 h-6 rounded border border-gray-300 dark:border-gray-600 mb-1 palette-detail-${getSafePaletteId(palette.id)}-${key}`}
                        />
                        <span className="text-xs text-gray-500 dark:text-gray-400 capitalize">
                          {key === 'background' ? 'bg' : key === 'text' ? 'txt' : key.substring(0, 3)}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Info Box */}
        <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700 rounded-lg">
          <p className="text-xs text-blue-700 dark:text-blue-300">
            {t('themeBuilderEditor.paletteTip', 'Tip: Click "Apply" to instantly apply a palette to all theme colors, then fine-tune individual colors.')}
          </p>
        </div>

        <style jsx>{`
          ${paletteColorRules}
        `}</style>
      </div>
    </div>
  )
}
