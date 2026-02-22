'use client'

import { X } from 'lucide-react'
import { useThemeStore } from '@/hooks/useThemeStore'
import { ThemePreset } from '@/types/theme'

interface PresetGalleryProps {
  isOpen: boolean
  onClose: () => void
  onSelectPreset: (preset: ThemePreset) => void
}

const PRESET_TEMPLATES: ThemePreset[] = [
  {
    id: 'modern-blue',
    name: 'Modern Blue',
    description: 'Clean and professional blue palette',
    colors: {
      primary: '#3b82f6',
      secondary: '#1e40af',
      accent: '#0ea5e9',
      background: '#ffffff',
      text: '#1f2937',
      border: '#e5e7eb',
      shadow: '#1f2937',
    },
    typography: {
      fontFamily: 'Inter',
      baseFontSize: 14,
      borderRadius: 8,
      lineHeight: 1.5,
    },
  },
  {
    id: 'warm-orange',
    name: 'Warm Orange',
    description: 'Warm and inviting orange tones',
    colors: {
      primary: '#d97706',
      secondary: '#b45309',
      accent: '#f59e0b',
      background: '#fefce8',
      text: '#78350f',
      border: '#fcd34d',
      shadow: '#78350f',
    },
    typography: {
      fontFamily: 'Georgia',
      baseFontSize: 16,
      borderRadius: 12,
      lineHeight: 1.6,
    },
  },
  {
    id: 'fresh-green',
    name: 'Fresh Green',
    description: 'Natural and fresh green palette',
    colors: {
      primary: '#059669',
      secondary: '#047857',
      accent: '#10b981',
      background: '#f0fdf4',
      text: '#1f2937',
      border: '#d1fae5',
      shadow: '#1f2937',
    },
    typography: {
      fontFamily: 'Poppins',
      baseFontSize: 15,
      borderRadius: 10,
      lineHeight: 1.6,
    },
  },
  {
    id: 'dark-purple',
    name: 'Dark Purple',
    description: 'Elegant dark purple theme',
    colors: {
      primary: '#7c3aed',
      secondary: '#6d28d9',
      accent: '#a78bfa',
      background: '#1f2937',
      text: '#f3f4f6',
      border: '#4b5563',
      shadow: '#000000',
    },
    typography: {
      fontFamily: 'Playfair Display',
      baseFontSize: 16,
      borderRadius: 6,
      lineHeight: 1.7,
    },
  },
  {
    id: 'minimal-gray',
    name: 'Minimal Gray',
    description: 'Minimalist grayscale palette',
    colors: {
      primary: '#374151',
      secondary: '#6b7280',
      accent: '#9ca3af',
      background: '#f9fafb',
      text: '#1f2937',
      border: '#e5e7eb',
      shadow: '#6b7280',
    },
    typography: {
      fontFamily: 'Helvetica',
      baseFontSize: 14,
      borderRadius: 4,
      lineHeight: 1.5,
    },
  },
  {
    id: 'vibrant-pink',
    name: 'Vibrant Pink',
    description: 'Bold and vibrant pink accent',
    colors: {
      primary: '#ec4899',
      secondary: '#be185d',
      accent: '#f472b6',
      background: '#fdf2f8',
      text: '#831843',
      border: '#fbcfe8',
      shadow: '#831843',
    },
    typography: {
      fontFamily: 'Poppins',
      baseFontSize: 15,
      borderRadius: 12,
      lineHeight: 1.6,
    },
  },
]

export function PresetGallery({ isOpen, onClose, onSelectPreset }: PresetGalleryProps) {
  if (!isOpen) return null

  const getSafePresetId = (id: string) => id.replace(/[^a-zA-Z0-9_-]/g, '-')
  const presetColorRules = PRESET_TEMPLATES.map((preset) => {
    const safeId = getSafePresetId(preset.id)
    return Object.values(preset.colors)
      .slice(0, 7)
      .map(
        (color, idx) => `
          .preset-${safeId}-color-${idx} {
            background-color: ${color};
          }
        `
      )
      .join('\n')
  }).join('\n')

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-2xl max-w-5xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="sticky top-0 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 p-6 flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Theme Presets</h2>
            <p className="text-gray-600 dark:text-gray-400 mt-1">
              Choose from our professionally designed theme templates
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        {/* Preset Grid */}
        <div className="p-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {PRESET_TEMPLATES.map((preset) => (
            <button
              key={preset.id}
              onClick={() => {
                onSelectPreset(preset)
                onClose()
              }}
              className="group text-left rounded-lg border-2 border-gray-200 dark:border-gray-600 hover:border-primary-400 dark:hover:border-primary-500 transition overflow-hidden"
            >
              {/* Color Preview */}
              <div className="h-32 flex gap-1 p-3 bg-gray-100 dark:bg-gray-700">
                {Object.values(preset.colors).slice(0, 7).map((color, idx) => (
                  <div
                    key={idx}
                    className={`flex-1 rounded preset-${getSafePresetId(preset.id)}-color-${idx}`}
                  />
                ))}
              </div>

              {/* Info */}
              <div className="p-4">
                <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400 transition">
                  {preset.name}
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                  {preset.description}
                </p>
                <div className="mt-3">
                  <p className="text-xs text-gray-500 dark:text-gray-500">
                    {preset.typography.fontFamily}
                  </p>
                </div>
              </div>
            </button>
          ))}
        </div>

        {/* Footer */}
        <div className="sticky bottom-0 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 p-6">
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Select a preset to apply it to your current theme. You can customize all settings after.
          </p>
        </div>
        <style jsx>{`
          ${presetColorRules}
        `}</style>
      </div>
    </div>
  )
}
