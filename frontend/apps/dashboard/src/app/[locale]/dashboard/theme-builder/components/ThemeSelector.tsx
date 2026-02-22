'use client'

import { useState } from 'react'
import { Copy, Trash2, Download, Upload, Plus } from 'lucide-react'
import { useThemeStore } from '@/hooks/useThemeStore'

const PRESET_THEMES = [
  {
    id: 'modern',
    name: 'Modern',
    colors: {
      primary: '#3b82f6',
      secondary: '#1e40af',
      accent: '#0ea5e9',
      background: '#ffffff',
      text: '#1f2937',
      border: '#e5e7eb',
      shadow: '#1f2937',
    },
  },
  {
    id: 'warm',
    name: 'Warm',
    colors: {
      primary: '#d97706',
      secondary: '#b45309',
      accent: '#f59e0b',
      background: '#fefce8',
      text: '#78350f',
      border: '#fcd34d',
      shadow: '#78350f',
    },
  },
  {
    id: 'fresh',
    name: 'Fresh',
    colors: {
      primary: '#059669',
      secondary: '#047857',
      accent: '#10b981',
      background: '#f0fdf4',
      text: '#1f2937',
      border: '#d1fae5',
      shadow: '#1f2937',
    },
  },
]

export function ThemeSelector() {
  const currentTheme = useThemeStore((state) => state.currentTheme)
  const themes = useThemeStore((state) => state.themes)
  const setCurrentTheme = useThemeStore((state) => state.setCurrentTheme)

  const [showActions, setShowActions] = useState<string | null>(null)
  const getSafeThemeId = (id: string) => id.replace(/[^a-zA-Z0-9_-]/g, '-')

  const handleApplyPreset = (preset: any) => {
    if (!currentTheme) return

    setCurrentTheme({
      ...currentTheme,
      name: preset.name,
      colors: preset.colors,
    })
  }

  const presetColorRules = PRESET_THEMES.map((preset) => {
    const safeId = getSafeThemeId(preset.id)
    return Object.values(preset.colors)
      .slice(0, 3)
      .map(
        (color, idx) => `
          .preset-${safeId}-color-${idx} {
            background-color: ${String(color)};
          }
        `
      )
      .join('\n')
  }).join('\n')

  const customThemeColorRules = themes
    .map((theme) => {
      const safeId = getSafeThemeId(theme.id)
      const defaultColors = {
        primary: '#3b82f6',
        secondary: '#1e40af',
        accent: '#0ea5e9',
      }
      const themeColors = theme.colors || defaultColors
      const palette = [
        themeColors.primary || defaultColors.primary,
        themeColors.secondary || defaultColors.secondary,
        themeColors.accent || defaultColors.accent,
      ]

      return palette
        .map(
          (color, idx) => `
            .theme-${safeId}-color-${idx} {
              background-color: ${color};
            }
          `
        )
        .join('\n')
    })
    .join('\n')

  return (
    <div className="h-full flex flex-col bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
      {/* Header */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Themes</h2>

        <button className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg text-sm font-medium transition">
          <Plus className="h-4 w-4" />
          New Theme
        </button>
      </div>

      {/* Presets */}
      <div className="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">Presets</h3>
        <div className="space-y-2">
          {PRESET_THEMES.map((preset) => (
            <button
              key={preset.id}
              onClick={() => handleApplyPreset(preset)}
              className="w-full p-3 rounded-lg border-2 border-gray-200 dark:border-gray-600 hover:border-primary-400 dark:hover:border-primary-500 transition group"
            >
              <div className="flex gap-2 mb-2">
                {Object.values(preset.colors)
                  .slice(0, 3)
                  .map((color, idx) => (
                    <div
                      key={idx}
                      className={`w-6 h-6 rounded preset-${getSafeThemeId(preset.id)}-color-${idx}`}
                    />
                  ))}
              </div>
              <p className="text-xs font-medium text-gray-700 dark:text-gray-300 group-hover:text-primary-600 dark:group-hover:text-primary-400">
                {preset.name}
              </p>
            </button>
          ))}
        </div>
      </div>

      {/* Custom Themes */}
      <div className="p-6 flex-1 overflow-y-auto border-b border-gray-200 dark:border-gray-700">
        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
          Custom Themes ({themes.length})
        </h3>

        <div className="space-y-2">
          {themes.length === 0 ? (
            <p className="text-xs text-gray-500 dark:text-gray-400">No custom themes yet</p>
          ) : (
            themes.map((theme) => {
              // Provide default colors if theme.colors is missing
              const defaultColors = {
                primary: '#3b82f6',
                secondary: '#1e40af',
                accent: '#0ea5e9',
              }
              const themeColors = theme.colors || defaultColors
              const colorPalette = [
                themeColors.primary || defaultColors.primary,
                themeColors.secondary || defaultColors.secondary,
                themeColors.accent || defaultColors.accent,
              ]

              return (
                <div
                  key={theme.id}
                  className={`relative p-3 rounded-lg border-2 transition cursor-pointer ${
                    currentTheme?.id === theme.id
                      ? 'border-primary-600 bg-primary-50 dark:bg-primary-900/20'
                      : 'border-gray-200 dark:border-gray-600 hover:border-primary-300 dark:hover:border-primary-700'
                  }`}
                  onClick={() => setCurrentTheme(theme)}
                >
                  <div className="flex items-start gap-2 mb-2">
                    {colorPalette.map((color, idx) => (
                      <div
                        key={idx}
                        className={`w-5 h-5 rounded theme-${getSafeThemeId(theme.id)}-color-${idx}`}
                      />
                    ))}
                  </div>
                  <p className="text-xs font-medium text-gray-700 dark:text-gray-300 truncate">
                    {theme.name}
                </p>

                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    setShowActions(showActions === theme.id ? null : theme.id)
                  }}
                  className="mt-2 text-xs text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                >
                  Actions ▼
                </button>

                {showActions === theme.id && (
                  <div className="mt-2 space-y-1 pt-2 border-t border-gray-200 dark:border-gray-600">
                    <button className="w-full text-left px-2 py-1 text-xs hover:bg-gray-100 dark:hover:bg-gray-700 rounded flex items-center gap-2 transition">
                      <Copy className="h-3 w-3" />
                      Duplicate
                    </button>
                    <button className="w-full text-left px-2 py-1 text-xs hover:bg-gray-100 dark:hover:bg-gray-700 rounded flex items-center gap-2 transition">
                      <Download className="h-3 w-3" />
                      Export
                    </button>
                    <button className="w-full text-left px-2 py-1 text-xs hover:bg-red-100 dark:hover:bg-red-900/30 rounded flex items-center gap-2 text-red-600 dark:text-red-400 transition">
                      <Trash2 className="h-3 w-3" />
                      Delete
                    </button>
                  </div>
                )}
              </div>
              )
            })
          )}
        </div>
      </div>
      <style jsx>{`
        ${presetColorRules}
        ${customThemeColorRules}
      `}</style>

      {/* Import/Export */}
      <div className="p-6 border-t border-gray-200 dark:border-gray-700 space-y-2">
        <button className="w-full flex items-center justify-center gap-2 px-4 py-2 border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg text-sm transition">
          <Upload className="h-4 w-4" />
          Import
        </button>
      </div>
    </div>
  )
}
