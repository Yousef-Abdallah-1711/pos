'use client'

import React from 'react'
import type { ProductsSectionProps } from '../types'
import { getMockDataForComponent } from '../mockData'

/**
 * ProductsSection: Display products in grid or list layout
 */
export const ProductsSection: React.FC<ProductsSectionProps> = ({
  title_en = 'Our Products',
  title_ar = 'Our Products',
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

  const mockData = getMockDataForComponent('products')
  const items = config.items && config.items.length > 0 ? config.items : mockData.items || []

  const gridClass =
    layout === 'grid'
      ? columns === 3
        ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
        : columns === 2
          ? 'grid-cols-1 sm:grid-cols-2'
          : 'grid-cols-1'
      : 'grid-cols-1'

  const formatPrice = (value: unknown) => {
    if (typeof value !== 'number') {
      return ''
    }
    return `${value.toFixed(2)} EGP`
  }

  return (
    <section
      dir={isArabic ? 'rtl' : 'ltr'}
      className="products-section w-full py-12 sm:py-16 px-4 sm:px-6 bg-white dark:bg-gray-950"
    >
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-10 sm:mb-12">
          {title && <h2 className="text-3xl sm:text-4xl font-bold mb-4 text-gray-900 dark:text-gray-100">{title}</h2>}
          {description && <p className="text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">{description}</p>}
        </div>

        {items.length > 0 ? (
          <div className={`grid gap-6 sm:gap-8 ${gridClass}`}>
            {items.map((product: any, idx: number) => (
              <article
                key={product.id || idx}
                className="rounded-xl border border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm hover:shadow-md transition overflow-hidden"
              >
                {showImages && product.image && (
                  <div className="relative w-full h-44 sm:h-48 bg-gray-200 dark:bg-gray-700 overflow-hidden">
                    <img src={product.image} alt={product.name || 'Product image'} className="w-full h-full object-cover" />
                  </div>
                )}

                <div className={`p-4 ${isArabic ? 'text-right' : 'text-left'}`}>
                  {product.name && <h3 className="text-lg font-bold mb-2 text-gray-900 dark:text-gray-100 break-words">{product.name}</h3>}

                  {product.description && <p className="text-sm text-gray-600 dark:text-gray-300 mb-4 break-words">{product.description}</p>}

                  {product.category && (
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-3 uppercase tracking-wide">{product.category}</p>
                  )}

                  <div className={`flex items-center gap-3 ${isArabic ? 'justify-between flex-row-reverse' : 'justify-between'}`}>
                    {showPrices && typeof product.price === 'number' && (
                      <span className="text-lg sm:text-xl font-bold text-primary">{formatPrice(product.price)}</span>
                    )}
                    <button className="bg-primary text-white px-4 py-2 rounded-lg hover:opacity-90 transition text-sm sm:text-base">
                      {isArabic ? 'Add' : 'Add'}
                    </button>
                  </div>
                </div>
              </article>
            ))}
          </div>
        ) : (
          <div className="rounded-xl border border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 p-6 text-center text-gray-500 dark:text-gray-400">
            <p>No products available</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default ProductsSection
