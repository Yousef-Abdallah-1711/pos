/**
 * CTAEditor Component
 * Specialized editor for call-to-action button sections
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'
import { previewText } from '@/lib/previewText'

interface CTAConfig {
  primaryButtonText?: string
  primaryButtonUrl?: string
  primaryButtonColor?: string
  primaryButtonHoverColor?: string
  primaryButtonStyle?: 'solid' | 'outline' | 'ghost'
  primaryButtonSize?: 'sm' | 'md' | 'lg' | 'xl'

  secondaryButtonText?: string
  secondaryButtonUrl?: string
  secondaryButtonColor?: string
  secondaryButtonHoverColor?: string
  secondaryButtonStyle?: 'solid' | 'outline' | 'ghost'
  secondaryButtonSize?: 'sm' | 'md' | 'lg' | 'xl'

  showSecondaryButton?: boolean
  buttonLayout?: 'stacked' | 'horizontal'

  backgroundColor?: string
  backgroundImage?: string
  overlayColor?: string
  overlayOpacity?: number

  titleColor?: string
  descriptionColor?: string

  borderRadius?: number
  padding?: number

  alignment?: 'left' | 'center' | 'right'

  animation?: 'none' | 'pulse' | 'bounce' | 'glow'
  animationSpeed?: 'slow' | 'normal' | 'fast'

  linkTarget?: '_self' | '_blank'

  minHeight?: number
}

/**
 * CTAEditor Component
 */
export function CTAEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}): JSX.Element {
  const [config, setConfig] = useState<CTAConfig>(
    (component.config as unknown as CTAConfig) || {}
  )

  const [previewMode, setPreviewMode] = useState(false)

  const handleConfigChange = useCallback(
    (newConfig: Partial<CTAConfig>) => {
      const updatedConfig = { ...config, ...newConfig }
      setConfig(updatedConfig)

      const updatedComponent: ThemeComponent = {
        ...component,
        config: updatedConfig,
      }

      onChange(updatedComponent)

      if (onPreview) {
        onPreview(updatedComponent)
      }
    },
    [config, component, onChange, onPreview]
  )

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
          {/* Primary Button Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Primary Button</h3>

            {/* Button Text */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Text
              </label>
              <input
                type="text"
                value={config.primaryButtonText || ''}
                onChange={(e) => handleConfigChange({ primaryButtonText: e.target.value })}
                placeholder="Order Now"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Button URL */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button URL
              </label>
              <input
                type="url"
                value={config.primaryButtonUrl || ''}
                onChange={(e) => handleConfigChange({ primaryButtonUrl: e.target.value })}
                placeholder="https://example.com/order"
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            {/* Button Style */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Style
              </label>
              <select
                value={config.primaryButtonStyle || 'solid'}
                onChange={(e) =>
                  handleConfigChange({
                    primaryButtonStyle: e.target.value as 'solid' | 'outline' | 'ghost',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="solid">Solid Fill</option>
                <option value="outline">Outline Only</option>
                <option value="ghost">Ghost (Text Only)</option>
              </select>
            </div>

            {/* Button Size */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Size
              </label>
              <div className="flex gap-2">
                {(['sm', 'md', 'lg', 'xl'] as const).map((size) => (
                  <button
                    key={size}
                    onClick={() => handleConfigChange({ primaryButtonSize: size })}
                    className={`flex-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                      (config.primaryButtonSize || 'lg') === size
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {size.toUpperCase()}
                  </button>
                ))}
              </div>
            </div>

            {/* Button Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.primaryButtonColor || '#3b82f6'}
                  onChange={(e) => handleConfigChange({ primaryButtonColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.primaryButtonColor || '#3b82f6'}
                  onChange={(e) => handleConfigChange({ primaryButtonColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Button Hover Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Hover Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.primaryButtonHoverColor || '#2563eb'}
                  onChange={(e) => handleConfigChange({ primaryButtonHoverColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.primaryButtonHoverColor || '#2563eb'}
                  onChange={(e) => handleConfigChange({ primaryButtonHoverColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Secondary Button Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showSecondaryButton ?? false}
                onChange={(e) => handleConfigChange({ showSecondaryButton: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Show Secondary Button
              </span>
            </label>

            {config.showSecondaryButton && (
              <>
                {/* Button Text */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Button Text
                  </label>
                  <input
                    type="text"
                    value={config.secondaryButtonText || ''}
                    onChange={(e) => handleConfigChange({ secondaryButtonText: e.target.value })}
                    placeholder="Learn More"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                {/* Button URL */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Button URL
                  </label>
                  <input
                    type="url"
                    value={config.secondaryButtonUrl || ''}
                    onChange={(e) => handleConfigChange({ secondaryButtonUrl: e.target.value })}
                    placeholder="https://example.com"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                {/* Button Style */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Button Style
                  </label>
                  <select
                    value={config.secondaryButtonStyle || 'outline'}
                    onChange={(e) =>
                      handleConfigChange({
                        secondaryButtonStyle: e.target.value as 'solid' | 'outline' | 'ghost',
                      })
                    }
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="solid">Solid Fill</option>
                    <option value="outline">Outline Only</option>
                    <option value="ghost">Ghost (Text Only)</option>
                  </select>
                </div>

                {/* Button Color */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Button Color
                  </label>
                  <div className="flex items-center gap-3">
                    <input
                      type="color"
                      value={config.secondaryButtonColor || '#9ca3af'}
                      onChange={(e) => handleConfigChange({ secondaryButtonColor: e.target.value })}
                      className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                    />
                    <input
                      type="text"
                      value={config.secondaryButtonColor || '#9ca3af'}
                      onChange={(e) => handleConfigChange({ secondaryButtonColor: e.target.value })}
                      className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                        bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                        focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>

                {/* Button Layout */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Button Layout
                  </label>
                  <div className="flex gap-2">
                    {(['stacked', 'horizontal'] as const).map((layout) => (
                      <button
                        key={layout}
                        onClick={() => handleConfigChange({ buttonLayout: layout })}
                        className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                          (config.buttonLayout || 'horizontal') === layout
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                        }`}
                      >
                        {layout === 'stacked' ? 'Stacked' : 'Side by Side'}
                      </button>
                    ))}
                  </div>
                </div>
              </>
            )}
          </div>

          {/* Background & Styling */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Styling</h3>

            {/* Background Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.backgroundColor || '#f3f4f6'}
                  onChange={(e) => handleConfigChange({ backgroundColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.backgroundColor || '#f3f4f6'}
                  onChange={(e) => handleConfigChange({ backgroundColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Alignment */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Content Alignment
              </label>
              <div className="flex gap-2">
                {(['left', 'center', 'right'] as const).map((align) => (
                  <button
                    key={align}
                    onClick={() => handleConfigChange({ alignment: align })}
                    className={`flex-1 px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
                      (config.alignment || 'center') === align
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {align === 'left' ? '◀' : align === 'center' ? '●' : '▶'}
                  </button>
                ))}
              </div>
            </div>

            {/* Padding */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Section Padding (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="20"
                  max="100"
                  step="10"
                  value={config.padding || 60}
                  onChange={(e) => handleConfigChange({ padding: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.padding || 60}px
                </span>
              </div>
            </div>

            {/* Border Radius */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Border Radius (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="30"
                  step="2"
                  value={config.borderRadius || 8}
                  onChange={(e) => handleConfigChange({ borderRadius: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.borderRadius || 8}px
                </span>
              </div>
            </div>

            {/* Minimum Height */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Section Height (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="200"
                  max="600"
                  step="50"
                  value={config.minHeight || 300}
                  onChange={(e) => handleConfigChange({ minHeight: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.minHeight || 300}px
                </span>
              </div>
            </div>
          </div>

          {/* Animation Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Animation</h3>

            {/* Animation Type */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Button Animation
              </label>
              <select
                value={config.animation || 'none'}
                onChange={(e) =>
                  handleConfigChange({ animation: e.target.value as 'none' | 'pulse' | 'bounce' | 'glow' })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="none">None</option>
                <option value="pulse">Pulse (Soft breathing)</option>
                <option value="bounce">Bounce (Up and down)</option>
                <option value="glow">Glow (Glowing effect)</option>
              </select>
            </div>

            {config.animation !== 'none' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Animation Speed
                </label>
                <div className="flex gap-2">
                  {(['slow', 'normal', 'fast'] as const).map((speed) => (
                    <button
                      key={speed}
                      onClick={() => handleConfigChange({ animationSpeed: speed })}
                      className={`flex-1 px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
                        (config.animationSpeed || 'normal') === speed
                          ? 'bg-blue-600 text-white'
                          : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                      }`}
                    >
                      {speed.charAt(0).toUpperCase() + speed.slice(1)}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Link Behavior */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Link Behavior</h3>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Open Link In
              </label>
              <div className="flex gap-2">
                {(['_self', '_blank'] as const).map((target) => (
                  <button
                    key={target}
                    onClick={() => handleConfigChange({ linkTarget: target })}
                    className={`flex-1 px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
                      (config.linkTarget || '_blank') === target
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {target === '_self' ? 'Same Tab' : 'New Tab'}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {previewMode && (
        <div className="cta-preview-shell rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="max-w-2xl">
            <h2 className="text-4xl font-bold text-gray-900 dark:text-white text-center mb-4">
              {previewText(component.title, 'Ready to Get Started?')}
            </h2>

            {component.subtitle && (
              <p className="cta-preview-description text-lg text-center mb-8">
                {previewText(component.subtitle)}
              </p>
            )}

            <div className="cta-preview-buttons">
              {/* Primary Button */}
              <button
                className="cta-preview-primary font-semibold transition-all hover:shadow-lg"
              >
                {config.primaryButtonText || 'Get Started'}
              </button>

              {/* Secondary Button */}
              {config.showSecondaryButton && (
                <button
                  className="cta-preview-secondary font-semibold transition-all hover:shadow-lg"
                >
                  {config.secondaryButtonText || 'Learn More'}
                </button>
              )}
            </div>
          </div>
          <style jsx>{`
            .cta-preview-shell {
              background-color: ${config.backgroundColor || '#f3f4f6'};
              min-height: ${config.minHeight || 300}px;
              padding: ${config.padding || 60}px;
              display: flex;
              flex-direction: column;
              align-items: ${
                config.alignment === 'left'
                  ? 'flex-start'
                  : config.alignment === 'right'
                    ? 'flex-end'
                    : 'center'
              };
              justify-content: center;
            }

            .cta-preview-description {
              color: ${config.descriptionColor || '#6b7280'};
            }

            .cta-preview-buttons {
              display: flex;
              flex-direction: ${config.buttonLayout === 'stacked' ? 'column' : 'row'};
              gap: 1rem;
              justify-content: center;
            }

            .cta-preview-primary {
              background-color: ${
                config.primaryButtonStyle === 'solid'
                  ? config.primaryButtonColor || '#3b82f6'
                  : 'transparent'
              };
              color: ${
                config.primaryButtonStyle === 'solid'
                  ? '#ffffff'
                  : config.primaryButtonColor || '#3b82f6'
              };
              border-color: ${config.primaryButtonColor || '#3b82f6'};
              border-style: solid;
              border-width: ${config.primaryButtonStyle === 'outline' ? '2px' : '0'};
              border-radius: ${config.borderRadius || 8}px;
              padding: ${
                config.primaryButtonSize === 'sm'
                  ? '0.5rem 1.5rem'
                  : config.primaryButtonSize === 'md'
                    ? '0.75rem 2rem'
                    : config.primaryButtonSize === 'lg'
                      ? '1rem 2.5rem'
                      : '1.25rem 3rem'
              };
              font-size: ${
                config.primaryButtonSize === 'sm'
                  ? '0.875rem'
                  : config.primaryButtonSize === 'md'
                    ? '1rem'
                    : config.primaryButtonSize === 'lg'
                      ? '1.125rem'
                      : '1.25rem'
              };
            }

            .cta-preview-primary:hover {
              background-color: ${
                config.primaryButtonStyle === 'solid'
                  ? config.primaryButtonHoverColor || '#2563eb'
                  : 'transparent'
              };
            }

            .cta-preview-secondary {
              background-color: ${
                config.secondaryButtonStyle === 'solid'
                  ? config.secondaryButtonColor || '#9ca3af'
                  : 'transparent'
              };
              color: ${
                config.secondaryButtonStyle === 'solid'
                  ? '#ffffff'
                  : config.secondaryButtonColor || '#9ca3af'
              };
              border-color: ${config.secondaryButtonColor || '#9ca3af'};
              border-style: solid;
              border-width: ${config.secondaryButtonStyle === 'outline' ? '2px' : '0'};
              border-radius: ${config.borderRadius || 8}px;
              padding: ${
                config.secondaryButtonSize === 'sm'
                  ? '0.5rem 1.5rem'
                  : config.secondaryButtonSize === 'md'
                    ? '0.75rem 2rem'
                    : config.secondaryButtonSize === 'lg'
                      ? '1rem 2.5rem'
                      : '1.25rem 3rem'
              };
              font-size: ${
                config.secondaryButtonSize === 'sm'
                  ? '0.875rem'
                  : config.secondaryButtonSize === 'md'
                    ? '1rem'
                    : config.secondaryButtonSize === 'lg'
                      ? '1.125rem'
                      : '1.25rem'
              };
            }

            .cta-preview-secondary:hover {
              background-color: ${
                config.secondaryButtonStyle === 'solid'
                  ? config.secondaryButtonHoverColor || '#6b7280'
                  : 'transparent'
              };
            }
          `}</style>
        </div>
      )}
    </div>
  )
}

export default CTAEditor
