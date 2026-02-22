'use client'

import React from 'react'
import type { ProductsProps } from './types'
import { mockData } from './mockData'

/**
 * Products: Display products in grid or list layout
 *
 * Shows products with:
 * - Flexible grid/list layouts
 * - Product images and prices
 * - Category filtering
 * - Responsive design
 * - RTL support
 */
export const Products: React.FC<ProductsProps> = ({
  title_en = 'Our Products',
  title_ar = 'منتجاتنا',
  description_en = '',
  description_ar = '',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const description = isArabic ? description_ar : description_en
  const layout = config.layout || 'grid'
  const columns = config.columns || 3
  const showPrices = config.show_prices !== false
  const showImages = config.show_images !== false

  // Use provided items or fallback to mock data
  const items = config.items && config.items.length > 0 ? config.items : (mockData.items || [])

  return (
    <section
      className="w-full py-16 px-4 bg-white"
      dir={isArabic ? 'rtl' : 'ltr'}
    >
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="text-center mb-12">
          {title && <h2 className="text-4xl font-bold mb-4">{title}</h2>}
          {description && <p className="text-gray-600 max-w-2xl mx-auto">{description}</p>}
        </div>

        {/* Products Grid */}
        {items.length > 0 ? (
          <div
            className={`grid gap-8 ${
              layout === 'grid'
                ? columns === 3
                  ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
                  : columns === 2
                  ? 'grid-cols-1 sm:grid-cols-2'
                  : 'grid-cols-1'
                : 'grid-cols-1'
            }`}
          >
            {items.map((product: any) => (
              <div key={product.id} className="bg-white rounded-lg shadow-sm hover:shadow-md transition overflow-hidden">
                {/* Product Image */}
                {showImages && product.image && (
                  <div className="relative w-full h-48 bg-gray-200 overflow-hidden">
                    <img
                      src={product.image}
                      alt={product.name}
                      className="w-full h-full object-cover"
                    />
                  </div>
                )}

                {/* Product Info */}
                <div className="p-4">
                  {product.name && (
                    <h3 className="text-lg font-bold mb-2 text-gray-900">
                      {product.name}
                    </h3>
                  )}

                  {product.description && (
                    <p className="text-sm text-gray-600 mb-4">
                      {product.description}
                    </p>
                  )}

                  {product.category && (
                    <p className="text-xs text-gray-500 mb-3 uppercase tracking-wide">
                      {product.category}
                    </p>
                  )}

                  <div className="flex items-center justify-between">
                    {showPrices && product.price && (
                      <span className="text-xl font-bold text-primary">
                        ${product.price.toFixed(2)}
                      </span>
                    )}
                    <button className="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark transition">
                      {isArabic ? 'إضافة' : 'Add'}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="bg-gray-100 rounded-lg p-6 text-center text-gray-500">
            <p>{isArabic ? 'لا توجد منتجات متاحة' : 'No products available'}</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default Products

// Component Definition for Registry
import { ComponentDefinition } from '../../../registry/types'
import manifest from '../manifest.json'

export const definition: ComponentDefinition = {
  id: 'products',
  name: 'Products Section',
  version: '1.0.0',
  component: Products,
  mockData,
  props: {
    title_en: 'string',
    title_ar: 'string',
    description_en: 'string',
    description_ar: 'string',
    isArabic: 'boolean',
    config: 'object',
  },
  manifest: manifest as any,
}
