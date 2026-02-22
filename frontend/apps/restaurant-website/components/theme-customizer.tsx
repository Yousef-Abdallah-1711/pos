'use client'

import { useState } from 'react'
import { ColorPresets, hexToHsl } from '@pos-saas/theme'
import { Button } from '@pos-saas/ui'

interface ThemeCustomizerProps {
  onThemeChange?: (colors: { primary: string; secondary: string; accent: string }) => void
}

export function ThemeCustomizer({ onThemeChange }: ThemeCustomizerProps) {
  const [selectedPreset, setSelectedPreset] = useState<keyof typeof ColorPresets>('slate')
  const [customColors, setCustomColors] = useState({
    primary: ColorPresets[selectedPreset].primary,
    secondary: ColorPresets[selectedPreset].secondary,
    accent: ColorPresets[selectedPreset].accent,
  })

  const handlePresetSelect = (preset: keyof typeof ColorPresets) => {
    setSelectedPreset(preset)
    const newColors = {
      primary: ColorPresets[preset].primary,
      secondary: ColorPresets[preset].secondary,
      accent: ColorPresets[preset].accent,
    }
    setCustomColors(newColors)
    onThemeChange?.(newColors)
  }

  const handleColorChange = (color: 'primary' | 'secondary' | 'accent', value: string) => {
    const hslColor = value.startsWith('#') ? hexToHsl(value) : value
    const newColors = { ...customColors, [color]: hslColor }
    setCustomColors(newColors)
    onThemeChange?.(newColors)
  }

  const presetColorStyles = Object.entries(ColorPresets)
    .map(
      ([key, preset]) => `
        .preset-${key}-primary { background-color: ${preset.primary}; }
        .preset-${key}-secondary { background-color: ${preset.secondary}; }
      `
    )
    .join('\n')

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold mb-4">Color Presets</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {Object.entries(ColorPresets).map(([key, preset]) => (
            <button
              key={key}
              onClick={() => handlePresetSelect(key as keyof typeof ColorPresets)}
              className={`p-4 rounded-lg border-2 transition-all ${
                selectedPreset === key ? 'border-primary ring-2 ring-primary' : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex gap-2 mb-2">
                <div className={`w-6 h-6 rounded preset-${key}-primary`} />
                <div className={`w-6 h-6 rounded preset-${key}-secondary`} />
              </div>
              <p className="text-sm font-medium">{preset.name}</p>
            </button>
          ))}
        </div>
      </div>

      <div>
        <h3 className="text-lg font-semibold mb-4">Custom Colors</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {(['primary', 'secondary', 'accent'] as const).map((color) => (
            <div key={color} className="space-y-2">
              <label className="block text-sm font-medium capitalize">{color} Color</label>
              <div className="flex gap-2">
                <input
                  type="color"
                  value={customColors[color]}
                  onChange={(e) => handleColorChange(color, e.target.value)}
                  className="h-10 w-20 rounded cursor-pointer border border-gray-200"
                />
                <input
                  type="text"
                  value={customColors[color]}
                  onChange={(e) => handleColorChange(color, e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-200 rounded text-sm"
                  placeholder="hsl(0 0% 50%)"
                />
              </div>
            </div>
          ))}
        </div>
      </div>

      <div>
        <h3 className="text-lg font-semibold mb-4">Preview</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="preview-primary p-8 rounded-lg text-white font-semibold">Primary</div>
          <div className="preview-secondary p-8 rounded-lg text-white font-semibold">Secondary</div>
          <div className="preview-accent p-8 rounded-lg text-white font-semibold">Accent</div>
        </div>
      </div>

      <div className="flex gap-2 pt-4">
        <Button variant="default">Save Theme</Button>
        <Button variant="outline">Reset to Default</Button>
      </div>

      <style jsx>{`
        .preview-primary {
          background-color: ${customColors.primary};
        }

        .preview-secondary {
          background-color: ${customColors.secondary};
        }

        .preview-accent {
          background-color: ${customColors.accent};
        }

        ${presetColorStyles}
      `}</style>
    </div>
  )
}
