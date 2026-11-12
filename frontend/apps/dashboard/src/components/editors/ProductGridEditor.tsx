/**
 * ProductGridEditor Component
 * Specialized editor for product grid/showcase components
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'
import { previewText } from '@/lib/previewText'

interface ProductGridEditorProps {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  availableProducts?: Array<{ id: number; name: string; price: number; image: string }>
  className?: string
}

interface ProductGridConfig {
  columns?: 2 | 3 | 4
  itemsPerPage?: number
  productIds?: number[]
  showPrice?: boolean
  showDescription?: boolean
  showRating?: boolean
  showCategory?: boolean
  cardStyle?: 'minimal' | 'detailed' | 'image-focus'
  borderStyle?: 'none' | 'shadow' | 'border'
  borderRadius?: number
  spacing?: number
  sortBy?: 'featured' | 'newest' | 'popular' | 'price-low' | 'price-high'
  priceColor?: string
  backgroundColor?: string
  cardBackgroundColor?: string
  hoverEffect?: 'lift' | 'overlay' | 'glow' | 'none'
  showCTA?: boolean
  ctaText?: string
  ctaVariant?: 'button' | 'link'
  ctaColor?: string
}

/**
 * ProductGridEditor Component
 */
export function ProductGridEditor({
  component,
  onChange,
  onPreview,
  availableProducts = [],
  className = '',
}: ProductGridEditorProps): JSX.Element {
  const [gridConfig, setGridConfig] = useState<ProductGridConfig>(
    (component.config as unknown as ProductGridConfig) || {}
  )

  const [previewMode, setPreviewMode] = useState(false)
  const [selectedProducts, setSelectedProducts] = useState<number[]>(gridConfig.productIds || [])

  /**
   * Handle config change
   */
  const handleConfigChange = useCallback(
    (newConfig: Partial<ProductGridConfig>) => {
      const updatedConfig = { ...gridConfig, ...newConfig }
      setGridConfig(updatedConfig)

      const updatedComponent: ThemeComponent = {
        ...component,
        config: updatedConfig,
      }

      onChange(updatedComponent)

      if (onPreview) {
        onPreview(updatedComponent)
      }
    },
    [gridConfig, component, onChange, onPreview]
  )

  const handleProductToggle = (productId: number) => {
    const updated = selectedProducts.includes(productId)
      ? selectedProducts.filter((id) => id !== productId)
      : [...selectedProducts, productId]

    setSelectedProducts(updated)
    handleConfigChange({ productIds: updated })
  }

  const handleSelectAllProducts = () => {
    const allIds = availableProducts.map((p) => p.id)
    setSelectedProducts(allIds)
    handleConfigChange({ productIds: allIds })
  }

  const handleClearProducts = () => {
    setSelectedProducts([])
    handleConfigChange({ productIds: [] })
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
          {/* Grid Layout Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Grid Layout</h3>

            {/* Columns */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Number of Columns
              </label>
              <div className="flex gap-2">
                {[2, 3, 4].map((cols) => (
                  <button
                    key={cols}
                    onClick={() => handleConfigChange({ columns: cols as 2 | 3 | 4 })}
                    className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                      (gridConfig.columns || 3) === cols
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                    }`}
                  >
                    {cols} Col
                  </button>
                ))}
              </div>
            </div>

            {/* Items Per Page */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Items Per Page
              </label>
              <select
                value={gridConfig.itemsPerPage || 12}
                onChange={(e) =>
                  handleConfigChange({ itemsPerPage: parseInt(e.target.value) })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="6">6</option>
                <option value="9">9</option>
                <option value="12">12</option>
                <option value="15">15</option>
                <option value="20">20</option>
              </select>
            </div>

            {/* Spacing */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Spacing Between Items (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="40"
                  step="4"
                  value={gridConfig.spacing || 16}
                  onChange={(e) => handleConfigChange({ spacing: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {gridConfig.spacing || 16}px
                </span>
              </div>
            </div>

            {/* Border Radius */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Card Border Radius (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="20"
                  step="2"
                  value={gridConfig.borderRadius || 8}
                  onChange={(e) => handleConfigChange({ borderRadius: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {gridConfig.borderRadius || 8}px
                </span>
              </div>
            </div>
          </div>

          {/* Display Options Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-3">
            <h3 className="font-semibold text-gray-900 dark:text-white mb-3">Display Options</h3>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={gridConfig.showPrice ?? true}
                onChange={(e) => handleConfigChange({ showPrice: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Price</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={gridConfig.showDescription ?? true}
                onChange={(e) => handleConfigChange({ showDescription: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Description</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={gridConfig.showRating ?? false}
                onChange={(e) => handleConfigChange({ showRating: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Rating</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={gridConfig.showCategory ?? false}
                onChange={(e) => handleConfigChange({ showCategory: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Category</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={gridConfig.showCTA ?? true}
                onChange={(e) => handleConfigChange({ showCTA: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show CTA Button</span>
            </label>
          </div>

          {/* Card Style Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Card Style</h3>

            {/* Card Style Type */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Style Type
              </label>
              <select
                value={gridConfig.cardStyle || 'detailed'}
                onChange={(e) =>
                  handleConfigChange({
                    cardStyle: e.target.value as 'minimal' | 'detailed' | 'image-focus',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="minimal">Minimal - Image + Title</option>
                <option value="detailed">Detailed - Full Info</option>
                <option value="image-focus">Image Focus - Large Preview</option>
              </select>
            </div>

            {/* Border Style */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Border Style
              </label>
              <select
                value={gridConfig.borderStyle || 'shadow'}
                onChange={(e) =>
                  handleConfigChange({
                    borderStyle: e.target.value as 'none' | 'shadow' | 'border',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="none">None</option>
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
                value={gridConfig.hoverEffect || 'lift'}
                onChange={(e) =>
                  handleConfigChange({
                    hoverEffect: e.target.value as 'lift' | 'overlay' | 'glow' | 'none',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="none">None</option>
                <option value="lift">Lift - Card elevates</option>
                <option value="overlay">Overlay - Dark overlay</option>
                <option value="glow">Glow - Border glow</option>
              </select>
            </div>
          </div>

          {/* Colors Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Colors</h3>

            {/* Background Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Grid Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={gridConfig.backgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ backgroundColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={gridConfig.backgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ backgroundColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Card Background Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Card Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={gridConfig.cardBackgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ cardBackgroundColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={gridConfig.cardBackgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ cardBackgroundColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Price Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Price Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={gridConfig.priceColor || '#ef4444'}
                  onChange={(e) => handleConfigChange({ priceColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={gridConfig.priceColor || '#ef4444'}
                  onChange={(e) => handleConfigChange({ priceColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* CTA Color */}
            {gridConfig.showCTA && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  CTA Button Color
                </label>
                <div className="flex items-center gap-3">
                  <input
                    type="color"
                    value={gridConfig.ctaColor || '#3b82f6'}
                    onChange={(e) => handleConfigChange({ ctaColor: e.target.value })}
                    className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                  />
                  <input
                    type="text"
                    value={gridConfig.ctaColor || '#3b82f6'}
                    onChange={(e) => handleConfigChange({ ctaColor: e.target.value })}
                    className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            )}
          </div>

          {/* Sort Options Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Sorting</h3>

            {/* Sort By */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Default Sort Order
              </label>
              <select
                value={gridConfig.sortBy || 'featured'}
                onChange={(e) =>
                  handleConfigChange({
                    sortBy: e.target.value as 'featured' | 'newest' | 'popular' | 'price-low' | 'price-high',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="featured">Featured</option>
                <option value="newest">Newest First</option>
                <option value="popular">Most Popular</option>
                <option value="price-low">Price: Low to High</option>
                <option value="price-high">Price: High to Low</option>
              </select>
            </div>
          </div>

          {/* Product Selection Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-gray-900 dark:text-white">Select Products</h3>
              <div className="flex gap-2">
                <button
                  onClick={handleSelectAllProducts}
                  className="px-3 py-1 text-sm rounded-lg bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200 hover:bg-blue-200 dark:hover:bg-blue-800 transition-colors"
                >
                  All
                </button>
                <button
                  onClick={handleClearProducts}
                  className="px-3 py-1 text-sm rounded-lg bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                >
                  Clear
                </button>
              </div>
            </div>

            {availableProducts.length === 0 ? (
              <p className="text-sm text-gray-500 dark:text-gray-400">No products available</p>
            ) : (
              <div className="grid grid-cols-2 gap-3 max-h-64 overflow-y-auto">
                {availableProducts.map((product) => (
                  <label
                    key={product.id}
                    className="flex items-center gap-2 cursor-pointer p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
                  >
                    <input
                      type="checkbox"
                      checked={selectedProducts.includes(product.id)}
                      onChange={() => handleProductToggle(product.id)}
                      className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
                    />
                    <div>
                      <p className="text-sm font-medium text-gray-900 dark:text-white">
                        {product.name}
                      </p>
                      <p className="text-xs text-gray-600 dark:text-gray-400">
                        ${product.price.toFixed(2)}
                      </p>
                    </div>
                  </label>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {previewMode && (
        <div className="product-grid-preview-shell rounded-lg border border-gray-200 dark:border-gray-700 p-8">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
            {previewText(component.title, 'Our Products')}
          </h2>

          <div className="product-grid-preview-layout">
            {availableProducts.slice(0, gridConfig.itemsPerPage || 12).map((product) => (
              <div
                key={product.id}
                className="product-grid-preview-card overflow-hidden hover:shadow-lg transition-shadow"
              >
                <img
                  src={product.image}
                  alt={product.name}
                  className="product-grid-preview-image w-full h-48 object-cover"
                  onError={(e) => {
                    e.currentTarget.classList.add('image-error')
                  }}
                />
                <div className="p-4">
                  <h3 className="font-semibold text-gray-900 dark:text-white mb-2">
                    {product.name}
                  </h3>
                  {gridConfig.showPrice && (
                    <p className="product-grid-preview-price text-lg font-bold mb-3">
                      ${product.price.toFixed(2)}
                    </p>
                  )}
                  {gridConfig.showCTA && (
                    <button
                      className="product-grid-preview-cta hover:opacity-90 transition-opacity"
                    >
                      {gridConfig.ctaText || 'View Product'}
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
          <style jsx>{`
            .product-grid-preview-shell {
              background-color: ${gridConfig.backgroundColor || '#ffffff'};
            }

            .product-grid-preview-layout {
              display: grid;
              grid-template-columns: repeat(${gridConfig.columns || 3}, 1fr);
              gap: ${gridConfig.spacing || 16}px;
            }

            .product-grid-preview-card {
              background-color: ${gridConfig.cardBackgroundColor || '#ffffff'};
              border-radius: ${gridConfig.borderRadius || 8}px;
              ${gridConfig.borderStyle === 'shadow' ? 'box-shadow: 0 2px 8px rgba(0,0,0,0.1);' : ''}
              ${gridConfig.borderStyle === 'border' ? 'border: 1px solid #e5e7eb;' : ''}
            }

            .product-grid-preview-image.image-error {
              background-color: #f3f4f6;
            }

            .product-grid-preview-price {
              color: ${gridConfig.priceColor || '#ef4444'};
            }

            .product-grid-preview-cta {
              background-color: ${gridConfig.ctaColor || '#3b82f6'};
              width: 100%;
              padding: 0.5rem 1rem;
              border-radius: 0.375rem;
              color: #ffffff;
              font-weight: 500;
              cursor: pointer;
            }
          `}</style>
        </div>
      )}
    </div>
  )
}

export default ProductGridEditor
