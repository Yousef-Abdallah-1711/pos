/**
 * HeroEditor Component
 * Specialized editor for hero/banner components with background image, overlay, and CTA
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'

interface HeroEditorProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}

interface HeroConfig {
  backgroundImage?: string
  backgroundPosition?: 'top' | 'center' | 'bottom'
  backgroundSize?: 'cover' | 'contain' | 'stretch'
  overlayColor?: string
  overlayOpacity?: number
  textColor?: string
  textAlign?: 'left' | 'center' | 'right'
  minHeight?: number
  ctaButtonText?: string
  ctaButtonColor?: string
  ctaButtonHoverColor?: string
  ctaButtonUrl?: string
  ctaButtonStyle?: 'solid' | 'outline' | 'ghost'
  ctaButtonSize?: 'sm' | 'md' | 'lg'
  contentPadding?: number
}

/**
 * HeroEditor Component
 */
export function HeroEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: HeroEditorProps): JSX.Element {
  const [heroConfig, setHeroConfig] = useState<HeroConfig>(
    (component.config as unknown as HeroConfig) || {}
  )

  const [previewMode, setPreviewMode] = useState(false)

  /**
   * Handle config change
   */
  const handleConfigChange = useCallback(
    (newConfig: Partial<HeroConfig>) => {
      const updatedConfig = { ...heroConfig, ...newConfig }
      setHeroConfig(updatedConfig)

      const updatedComponent: ThemeComponent = {
        ...component,
        config: updatedConfig,
      }

      onChange(updatedComponent)

      if (onPreview) {
        onPreview(updatedComponent)
      }
    },
    [heroConfig, component, onChange, onPreview]
  )

  const handleBackgroundImageChange = (url: string) => {
    handleConfigChange({ backgroundImage: url })
  }

  const handleOverlayOpacityChange = (opacity: number) => {
    handleConfigChange({ overlayOpacity: Math.max(0, Math.min(1, opacity / 100)) })
  }

  const handleMinHeightChange = (height: number) => {
    handleConfigChange({ minHeight: Math.max(300, height) })
  }

  const handleCTAChange = (field: keyof HeroConfig, value: string | number) => {
    handleConfigChange({ [field]: value })
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
        <button
          onClick={() => setPreviewMode(false)}
          className={`px-4 py-2 font-medium transition-colors ${
            !previewMode
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          }`}
        >
          Settings
        </button>
        <button
          onClick={() => setPreviewMode(true)}
          className={`px-4 py-2 font-medium transition-colors ${
            previewMode
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          }`}
        >
          Preview
        </button>
      </div>

      {!previewMode && (
        <div className="space-y-6">
          {/* Background Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Background</h3>

            {/* Background Image URL */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Background Image URL
              </label>
              <input
                type="url"
                value={heroConfig.backgroundImage || ''}
                onChange={(e) => handleBackgroundImageChange(e.target.value)}
                placeholder="https://example.com/hero-image.jpg"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              {heroConfig.backgroundImage && (
                <div className="mt-2 h-32 rounded-lg overflow-hidden">
                  <img
                    src={heroConfig.backgroundImage}
                    alt="Background preview"
                    className="w-full h-full object-cover"
                    onError={(e) => {
                      (e.target as HTMLImageElement).style.display = 'none'
                    }}
                  />
                </div>
              )}
            </div>

            {/* Background Position */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Background Position
              </label>
              <select
                value={heroConfig.backgroundPosition || 'center'}
                onChange={(e) =>
                  handleConfigChange({
                    backgroundPosition: e.target.value as 'top' | 'center' | 'bottom',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="top">Top</option>
                <option value="center">Center</option>
                <option value="bottom">Bottom</option>
              </select>
            </div>

            {/* Background Size */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Background Size
              </label>
              <select
                value={heroConfig.backgroundSize || 'cover'}
                onChange={(e) =>
                  handleConfigChange({
                    backgroundSize: e.target.value as 'cover' | 'contain' | 'stretch',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="cover">Cover</option>
                <option value="contain">Contain</option>
                <option value="stretch">Stretch</option>
              </select>
            </div>

            {/* Minimum Height */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Minimum Height (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="300"
                  max="800"
                  step="50"
                  value={heroConfig.minHeight || 400}
                  onChange={(e) => handleMinHeightChange(parseInt(e.target.value))}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {heroConfig.minHeight || 400}px
                </span>
              </div>
            </div>
          </div>

          {/* Overlay Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Overlay</h3>

            {/* Overlay Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Overlay Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={heroConfig.overlayColor || '#000000'}
                  onChange={(e) => handleConfigChange({ overlayColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={heroConfig.overlayColor || '#000000'}
                  onChange={(e) => handleConfigChange({ overlayColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Overlay Opacity */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Overlay Opacity
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="5"
                  value={(heroConfig.overlayOpacity || 0.5) * 100}
                  onChange={(e) => handleOverlayOpacityChange(parseInt(e.target.value))}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {Math.round((heroConfig.overlayOpacity || 0.5) * 100)}%
                </span>
              </div>
            </div>
          </div>

          {/* Text Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Text</h3>

            {/* Text Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Text Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={heroConfig.textColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ textColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={heroConfig.textColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ textColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Text Alignment */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Text Alignment
              </label>
              <div className="flex gap-2">
                {(['left', 'center', 'right'] as const).map((align) => (
                  <button
                    key={align}
                    onClick={() => handleConfigChange({ textAlign: align })}
                    className={`flex-1 px-3 py-2 rounded-lg font-medium transition-colors ${
                      heroConfig.textAlign === align
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {align === 'left' ? '◀' : align === 'center' ? '●' : '▶'}
                  </button>
                ))}
              </div>
            </div>

            {/* Content Padding */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Content Padding (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="10"
                  value={heroConfig.contentPadding || 40}
                  onChange={(e) => handleConfigChange({ contentPadding: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {heroConfig.contentPadding || 40}px
                </span>
              </div>
            </div>
          </div>

          {/* CTA Button Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">CTA Button</h3>

            {/* CTA Button Text */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Text
              </label>
              <input
                type="text"
                value={heroConfig.ctaButtonText || ''}
                onChange={(e) => handleCTAChange('ctaButtonText', e.target.value)}
                placeholder="Learn More"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* CTA Button URL */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button URL
              </label>
              <input
                type="url"
                value={heroConfig.ctaButtonUrl || ''}
                onChange={(e) => handleCTAChange('ctaButtonUrl', e.target.value)}
                placeholder="https://example.com"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* CTA Button Style */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Style
              </label>
              <select
                value={heroConfig.ctaButtonStyle || 'solid'}
                onChange={(e) =>
                  handleCTAChange('ctaButtonStyle', e.target.value as 'solid' | 'outline' | 'ghost')
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="solid">Solid</option>
                <option value="outline">Outline</option>
                <option value="ghost">Ghost</option>
              </select>
            </div>

            {/* CTA Button Size */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Size
              </label>
              <div className="flex gap-2">
                {(['sm', 'md', 'lg'] as const).map((size) => (
                  <button
                    key={size}
                    onClick={() => handleCTAChange('ctaButtonSize', size)}
                    className={`flex-1 px-3 py-2 rounded-lg font-medium transition-colors ${
                      heroConfig.ctaButtonSize === size
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {size.toUpperCase()}
                  </button>
                ))}
              </div>
            </div>

            {/* CTA Button Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={heroConfig.ctaButtonColor || '#3b82f6'}
                  onChange={(e) => handleCTAChange('ctaButtonColor', e.target.value)}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={heroConfig.ctaButtonColor || '#3b82f6'}
                  onChange={(e) => handleCTAChange('ctaButtonColor', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* CTA Button Hover Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Hover Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={heroConfig.ctaButtonHoverColor || '#2563eb'}
                  onChange={(e) => handleCTAChange('ctaButtonHoverColor', e.target.value)}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={heroConfig.ctaButtonHoverColor || '#2563eb'}
                  onChange={(e) => handleCTAChange('ctaButtonHoverColor', e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>
        </div>
      )}

      {previewMode && (
        <div className="rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div className="hero-preview-shell">
            <div className="hero-preview-content">
              <h1 className="text-4xl font-bold mb-4">{component.title?.en || 'Hero Title'}</h1>
              <p className="text-lg mb-6 max-w-2xl">{component.subtitle?.en || 'Hero subtitle goes here'}</p>
              {heroConfig.ctaButtonText && (
                <button
                  className="hero-preview-cta font-semibold rounded-lg transition-colors"
                >
                  {heroConfig.ctaButtonText}
                </button>
              )}
            </div>
          </div>
          <style jsx>{`
            .hero-preview-shell {
              background-image: linear-gradient(
                  rgba(${hexToRgb(heroConfig.overlayColor || '#000000')}, ${heroConfig.overlayOpacity || 0.5}),
                  rgba(${hexToRgb(heroConfig.overlayColor || '#000000')}, ${heroConfig.overlayOpacity || 0.5})
                ),
                url('${heroConfig.backgroundImage}');
              background-size: ${heroConfig.backgroundSize || 'cover'};
              background-position: ${heroConfig.backgroundPosition || 'center'};
              min-height: ${heroConfig.minHeight || 400}px;
              display: flex;
              align-items: center;
              justify-content: ${
                heroConfig.textAlign === 'left'
                  ? 'flex-start'
                  : heroConfig.textAlign === 'right'
                    ? 'flex-end'
                    : 'center'
              };
              padding: ${heroConfig.contentPadding || 40}px;
            }

            .hero-preview-content {
              color: ${heroConfig.textColor || '#ffffff'};
              text-align: ${heroConfig.textAlign || 'center'};
            }

            .hero-preview-cta {
              background-color: ${
                heroConfig.ctaButtonStyle === 'solid'
                  ? heroConfig.ctaButtonColor || '#3b82f6'
                  : 'transparent'
              };
              color: ${
                heroConfig.ctaButtonStyle === 'solid'
                  ? '#ffffff'
                  : heroConfig.ctaButtonColor || '#3b82f6'
              };
              border-color: ${heroConfig.ctaButtonColor || '#3b82f6'};
              border-style: solid;
              border-width: ${heroConfig.ctaButtonStyle === 'outline' ? '2px' : '0'};
              padding: ${
                heroConfig.ctaButtonSize === 'sm'
                  ? '0.5rem 1rem'
                  : heroConfig.ctaButtonSize === 'lg'
                    ? '1rem 2rem'
                    : '0.75rem 1.5rem'
              };
              font-size: ${
                heroConfig.ctaButtonSize === 'sm'
                  ? '0.875rem'
                  : heroConfig.ctaButtonSize === 'lg'
                    ? '1.125rem'
                    : '1rem'
              };
            }

            .hero-preview-cta:hover {
              background-color: ${
                heroConfig.ctaButtonStyle === 'solid'
                  ? heroConfig.ctaButtonHoverColor || '#2563eb'
                  : 'transparent'
              };
            }
          `}</style>
        </div>
      )}
    </div>
  )
}

/**
 * Convert hex color to RGB values
 */
function hexToRgb(hex: string): string {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  if (result) {
    const r = parseInt(result[1], 16)
    const g = parseInt(result[2], 16)
    const b = parseInt(result[3], 16)
    return `${r}, ${g}, ${b}`
  }
  return '0, 0, 0'
}

export default HeroEditor
