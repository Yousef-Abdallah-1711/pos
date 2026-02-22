'use client'

import { useState, useMemo } from 'react'
import { Palette, ChevronRight, Layers } from 'lucide-react'
import { ThemeData } from '@/types/theme'
import { getComponentCount, getThemeCategory, getThemeTags } from '@/lib/themeLoader'
import { debugError } from '@/lib/debug'

interface PresetsGalleryProps {
  presets: ThemeData[]
  onSelectPreset: (preset: ThemeData) => void
  isLoading?: boolean
}

function getSafePresetId(id: string): string {
  return id.replace(/[^a-zA-Z0-9_-]/g, '-')
}

// Helper function to extract colors from preset data
function getPresetColors(preset: ThemeData) {
  try {
    // Try to get colors directly
    if (preset.colors) {
      return preset.colors
    }

    // Fallback default colors
    return {
      primary: '#3b82f6',
      secondary: '#1e40af',
      accent: '#0ea5e9',
      background: '#ffffff',
      text: '#1f2937'
    }
  } catch (err) {
    debugError('Error parsing preset colors:', err)
    return {
      primary: '#3b82f6',
      secondary: '#1e40af',
      accent: '#0ea5e9',
      background: '#ffffff',
      text: '#1f2937'
    }
  }
}

// Helper function to get component count
function getThemeComponents(preset: ThemeData): number {
  try {
    if (preset.components && Array.isArray(preset.components)) {
      return preset.components.filter(c => c.enabled !== false).length
    }
    return 0
  } catch (err) {
    debugError('Error getting component count:', err)
    return 0
  }
}

export function PresetsGallery({ presets, onSelectPreset, isLoading = false }: PresetsGalleryProps) {
  const [hoveredId, setHoveredId] = useState<string | null>(null)

  const categorizedPresets = useMemo(() => {
    const categories: Record<string, ThemeData[]> = {
      'Professional': [],
      'Luxury': [],
      'Modern': [],
      'Casual': [],
      'Playful': [],
      'Other': []
    }

    presets.forEach((preset) => {
      // Try to get category from theme metadata
      try {
        const category = getThemeCategory(preset)
        if (categories[category]) {
          categories[category].push(preset)
        } else {
          categories['Other'].push(preset)
        }
      } catch {
        categories['Other'].push(preset)
      }
    })

    return Object.entries(categories).filter(([_, items]) => items.length > 0)
  }, [presets])

  const presetColorRules = useMemo(
    () =>
      presets
        .map((preset) => {
          const safeId = getSafePresetId(preset.id)
          const colors = getPresetColors(preset)
          const palette = [
            colors.primary,
            colors.secondary,
            colors.accent,
            colors.background,
            colors.text,
          ]

          return palette
            .map(
              (color, idx) => `
                .preset-color-${safeId}-${idx} {
                  background-color: ${color || '#ccc'};
                }
              `
            )
            .join('\n')
        })
        .join('\n'),
    [presets]
  )

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="animate-spin">
          <Palette className="h-8 w-8 text-gray-400" />
        </div>
      </div>
    )
  }

  if (presets.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <Palette className="h-8 w-8 mx-auto mb-2 opacity-50" />
        <p>No theme presets available</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
          Choose a Theme Template
        </h2>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Start with a pre-designed theme and customize it to match your brand
        </p>
      </div>

      {/* Presets by Category */}
      {categorizedPresets.map(([category, items]) => (
        <div key={category}>
          <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
            {category}
          </h3>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            {items.map((preset) => (
              <div
                key={preset.id}
                onMouseEnter={() => setHoveredId(preset.id)}
                onMouseLeave={() => setHoveredId(null)}
                className="group cursor-pointer"
                onClick={() => onSelectPreset(preset)}
              >
                {/* Preset Card */}
                <div className="relative bg-white dark:bg-gray-800 rounded-lg border-2 border-gray-200 dark:border-gray-700 overflow-hidden hover:border-blue-400 dark:hover:border-blue-500 transition-all hover:shadow-lg">
                  {/* Color Preview */}
                  <div className="flex h-24 gap-1 p-2">
                    {(() => {
                      const colors = getPresetColors(preset)
                      return [
                        colors.primary,
                        colors.secondary,
                        colors.accent,
                        colors.background,
                        colors.text
                      ].map((color, idx) => (
                        <div
                          key={idx}
                          className={`flex-1 rounded-md shadow-sm hover:shadow-md transition preset-color-${getSafePresetId(preset.id)}-${idx}`}
                          title={color || 'No color'}
                        />
                      ))
                    })()}
                  </div>

                  {/* Info */}
                  <div className="p-3 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-start justify-between gap-2 mb-1">
                      <h4 className="font-semibold text-sm text-gray-900 dark:text-white flex-1">
                        {preset.name}
                      </h4>
                      <div className="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                        <Layers className="h-3 w-3" />
                        {getThemeComponents(preset)}
                      </div>
                    </div>
                    {preset.description && (
                      <p className="text-xs text-gray-600 dark:text-gray-400 line-clamp-2 mb-2">
                        {preset.description}
                      </p>
                    )}
                    <div className="flex flex-wrap gap-1">
                      {(() => {
                        try {
                          const tags = getThemeTags(preset)
                          return tags.slice(0, 2).map((tag, idx) => (
                            <span key={idx} className="text-xs px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded">
                              {tag}
                            </span>
                          ))
                        } catch {
                          return null
                        }
                      })()}
                    </div>
                  </div>

                  {/* Hover Overlay */}
                  {hoveredId === preset.id && (
                    <div className="absolute inset-0 bg-blue-600/10 dark:bg-blue-500/10 flex items-center justify-center backdrop-blur-sm">
                      <div className="flex items-center gap-1 text-blue-600 dark:text-blue-400 font-medium">
                        Use Template
                        <ChevronRight className="h-4 w-4" />
                      </div>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
      <style jsx>{`
        ${presetColorRules}
      `}</style>

      {/* Info Box */}
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700 rounded-lg">
        <p className="text-sm text-blue-700 dark:text-blue-300">
          💡 <strong>Tip:</strong> Click on any template to start creating a new theme based on it.
          You can then customize all colors, fonts, and content to match your restaurant's brand.
        </p>
      </div>
    </div>
  )
}
