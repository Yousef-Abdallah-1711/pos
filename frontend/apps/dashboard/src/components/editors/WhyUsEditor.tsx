/**
 * WhyUsEditor Component
 * Specialized editor for why-us/features/benefits section components
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'

interface WhyUsEditorProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}

interface Feature {
  id: string
  title: string
  description: string
  icon?: string
  iconType?: 'emoji' | 'text' | 'image-url'
}

interface WhyUsConfig {
  layout?: 'grid' | 'list' | 'carousel'
  columns?: 2 | 3 | 4
  features?: Feature[]
  backgroundColor?: string
  cardBackgroundColor?: string
  titleColor?: string
  descriptionColor?: string
  iconSize?: number
  iconBackgroundColor?: string
  iconShape?: 'circle' | 'square' | 'rounded-square'
  spacing?: number
  borderRadius?: number
  showDividers?: boolean
  cardStyle?: 'minimal' | 'shadow' | 'border'
  hoverEffect?: 'lift' | 'glow' | 'color-change' | 'none'
}

/**
 * WhyUsEditor Component
 */
export function WhyUsEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: WhyUsEditorProps): JSX.Element {
  const [config, setConfig] = useState<WhyUsConfig>(
    (component.config as unknown as WhyUsConfig) || {
      layout: 'grid',
      columns: 3,
      features: [],
    }
  )

  const [previewMode, setPreviewMode] = useState(false)
  const [editingFeatureId, setEditingFeatureId] = useState<string | null>(null)

  /**
   * Handle config change
   */
  const handleConfigChange = useCallback(
    (newConfig: Partial<WhyUsConfig>) => {
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

  const addFeature = () => {
    const newFeature: Feature = {
      id: `feature-${Date.now()}`,
      title: 'Feature Title',
      description: 'Feature description goes here',
      icon: 'â­',
      iconType: 'emoji',
    }

    const updated = [...(config.features || []), newFeature]
    handleConfigChange({ features: updated })
  }

  const updateFeature = (featureId: string, updates: Partial<Feature>) => {
    const updated = (config.features || []).map((f) =>
      f.id === featureId ? { ...f, ...updates } : f
    )
    handleConfigChange({ features: updated })
  }

  const deleteFeature = (featureId: string) => {
    const updated = (config.features || []).filter((f) => f.id !== featureId)
    handleConfigChange({ features: updated })
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
          {/* Layout Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Layout</h3>

            {/* Layout Type */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Layout Type
              </label>
              <select
                value={config.layout || 'grid'}
                onChange={(e) =>
                  handleConfigChange({ layout: e.target.value as 'grid' | 'list' | 'carousel' })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="grid">Grid</option>
                <option value="list">List</option>
                <option value="carousel">Carousel</option>
              </select>
            </div>

            {/* Columns (for grid layout) */}
            {config.layout === 'grid' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Columns
                </label>
                <div className="flex gap-2">
                  {[2, 3, 4].map((cols) => (
                    <button
                      key={cols}
                      onClick={() => handleConfigChange({ columns: cols as 2 | 3 | 4 })}
                      className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                        (config.columns || 3) === cols
                          ? 'bg-blue-600 text-white'
                          : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                      }`}
                    >
                      {cols}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Spacing */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Spacing (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="40"
                  step="4"
                  value={config.spacing || 24}
                  onChange={(e) => handleConfigChange({ spacing: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.spacing || 24}px
                </span>
              </div>
            </div>
          </div>

          {/* Icon Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Icon Settings</h3>

            {/* Icon Size */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Icon Size (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="24"
                  max="96"
                  step="8"
                  value={config.iconSize || 48}
                  onChange={(e) => handleConfigChange({ iconSize: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.iconSize || 48}px
                </span>
              </div>
            </div>

            {/* Icon Shape */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Icon Shape
              </label>
              <div className="flex gap-2">
                {(['circle', 'square', 'rounded-square'] as const).map((shape) => (
                  <button
                    key={shape}
                    onClick={() => handleConfigChange({ iconShape: shape })}
                    className={`flex-1 px-4 py-2 rounded-lg font-medium text-sm transition-colors ${
                      (config.iconShape || 'circle') === shape
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {shape === 'circle' ? 'â—' : shape === 'square' ? 'â– ' : 'â—†'}
                  </button>
                ))}
              </div>
            </div>

            {/* Icon Background Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Icon Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.iconBackgroundColor || '#e0f2fe'}
                  onChange={(e) =>
                    handleConfigChange({ iconBackgroundColor: e.target.value })
                  }
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.iconBackgroundColor || '#e0f2fe'}
                  onChange={(e) =>
                    handleConfigChange({ iconBackgroundColor: e.target.value })
                  }
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Style Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Style</h3>

            {/* Card Style */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Card Style
              </label>
              <select
                value={config.cardStyle || 'minimal'}
                onChange={(e) =>
                  handleConfigChange({ cardStyle: e.target.value as 'minimal' | 'shadow' | 'border' })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="minimal">Minimal</option>
                <option value="shadow">Shadow</option>
                <option value="border">Border</option>
              </select>
            </div>

            {/* Hover Effect */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Hover Effect
              </label>
              <select
                value={config.hoverEffect || 'lift'}
                onChange={(e) =>
                  handleConfigChange({
                    hoverEffect: e.target.value as 'lift' | 'glow' | 'color-change' | 'none',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="none">None</option>
                <option value="lift">Lift</option>
                <option value="glow">Glow</option>
                <option value="color-change">Color Change</option>
              </select>
            </div>

            {/* Border Radius */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Border Radius (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="20"
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

            {/* Show Dividers */}
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showDividers ?? false}
                onChange={(e) => handleConfigChange({ showDividers: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Dividers</span>
            </label>
          </div>

          {/* Colors Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Colors</h3>

            {/* Background */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.backgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ backgroundColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.backgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ backgroundColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Title Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Title Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.titleColor || '#1f2937'}
                  onChange={(e) => handleConfigChange({ titleColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.titleColor || '#1f2937'}
                  onChange={(e) => handleConfigChange({ titleColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Description Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Description Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.descriptionColor || '#6b7280'}
                  onChange={(e) => handleConfigChange({ descriptionColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.descriptionColor || '#6b7280'}
                  onChange={(e) => handleConfigChange({ descriptionColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Features Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-gray-900 dark:text-white">Features</h3>
              <button
                onClick={addFeature}
                className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
              >
                + Add Feature
              </button>
            </div>

            <div className="space-y-3 max-h-96 overflow-y-auto">
              {(config.features || []).map((feature) => (
                <div
                  key={feature.id}
                  className="border border-gray-300 dark:border-gray-600 rounded-lg p-3 space-y-2"
                >
                  {editingFeatureId === feature.id ? (
                    <>
                      <input
                        type="text"
                        value={feature.icon || ''}
                        onChange={(e) =>
                          updateFeature(feature.id, { icon: e.target.value })
                        }
                        placeholder="Icon (emoji or text)"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <input
                        type="text"
                        value={feature.title}
                        onChange={(e) =>
                          updateFeature(feature.id, { title: e.target.value })
                        }
                        placeholder="Title"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <textarea
                        value={feature.description}
                        onChange={(e) =>
                          updateFeature(feature.id, { description: e.target.value })
                        }
                        placeholder="Description"
                        rows={2}
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <div className="flex gap-2">
                        <button
                          onClick={() => setEditingFeatureId(null)}
                          className="flex-1 px-3 py-1 rounded text-sm bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
                        >
                          Done
                        </button>
                        <button
                          onClick={() => deleteFeature(feature.id)}
                          className="flex-1 px-3 py-1 rounded text-sm bg-red-600 hover:bg-red-700 text-white font-medium transition-colors"
                        >
                          Delete
                        </button>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex items-start justify-between">
                        <div className="flex items-start gap-2 flex-1">
                          <span className="text-2xl">{feature.icon || 'â­'}</span>
                          <div>
                            <p className="font-medium text-gray-900 dark:text-white">
                              {feature.title}
                            </p>
                            <p className="text-sm text-gray-600 dark:text-gray-400">
                              {feature.description}
                            </p>
                          </div>
                        </div>
                        <button
                          onClick={() => setEditingFeatureId(feature.id)}
                          className="px-3 py-1 rounded text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                        >
                          Edit
                        </button>
                      </div>
                    </>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {previewMode && (
        <div className="whyus-preview-shell rounded-lg border border-gray-200 dark:border-gray-700 p-8">
          <h2 className="whyus-preview-title text-2xl font-bold mb-8">
            {component.title?.en || 'Why Choose Us'}
          </h2>

          <div className="whyus-preview-layout">
            {(config.features || []).map((feature) => (
              <div key={feature.id} className="whyus-preview-card">
                <div className="whyus-preview-icon">
                  {feature.icon || '⭐'}
                </div>

                <h3 className="whyus-preview-item-title font-semibold text-lg mb-2">
                  {feature.title}
                </h3>

                <p className="whyus-preview-item-description text-sm">
                  {feature.description}
                </p>

                {config.showDividers && (
                  <div className="mt-4 pt-4 border-t border-gray-300 dark:border-gray-600" />
                )}
              </div>
            ))}
          </div>
          <style jsx>{`
            .whyus-preview-shell {
              background-color: ${config.backgroundColor || '#ffffff'};
            }

            .whyus-preview-title,
            .whyus-preview-item-title {
              color: ${config.titleColor || '#1f2937'};
            }

            .whyus-preview-item-description {
              color: ${config.descriptionColor || '#6b7280'};
            }

            .whyus-preview-layout {
              display: ${config.layout === 'grid' ? 'grid' : 'flex'};
              grid-template-columns: ${
                config.layout === 'grid' ? `repeat(${config.columns || 3}, 1fr)` : 'none'
              };
              flex-direction: ${config.layout === 'list' ? 'column' : 'row'};
              gap: ${config.spacing || 24}px;
            }

            .whyus-preview-card {
              border-radius: ${config.borderRadius || 8}px;
              padding: 1.5rem;
              ${config.cardStyle === 'shadow' ? 'box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);' : ''}
              ${config.cardStyle === 'border' ? 'border: 1px solid #e5e7eb;' : ''}
            }

            .whyus-preview-icon {
              background-color: ${config.iconBackgroundColor || '#e0f2fe'};
              width: ${config.iconSize || 48}px;
              height: ${config.iconSize || 48}px;
              display: flex;
              align-items: center;
              justify-content: center;
              border-radius: ${
                config.iconShape === 'circle'
                  ? '9999px'
                  : config.iconShape === 'rounded-square'
                    ? '0.5rem'
                    : '0'
              };
              font-size: ${(config.iconSize || 48) * 0.6}px;
              margin-bottom: 1rem;
            }
          `}</style>
        </div>
      )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default WhyUsEditor
