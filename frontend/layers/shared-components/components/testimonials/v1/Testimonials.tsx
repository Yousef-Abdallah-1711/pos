'use client'

import React from 'react'
import type { TestimonialsProps } from './types'
import { mockData } from './mockData'

/**
 * Testimonials: Customer testimonials and reviews
 *
 * Features:
 * - Carousel or grid layout
 * - Star ratings
 * - Customer quotes
 * - RTL support
 */
export const Testimonials: React.FC<TestimonialsProps> = ({
  title_en = 'What Our Customers Say',
  title_ar = 'ما يقوله عملاؤنا',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const layout = config.layout || 'grid'

  // Use provided items or fallback to mock data
  const items = config.items && config.items.length > 0 ? config.items : (mockData.items || [])

  return (
    <section
      className="w-full py-16 px-4 bg-gray-50"
      dir={isArabic ? 'rtl' : 'ltr'}
    >
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="text-center mb-12">
          {title && <h2 className="text-4xl font-bold mb-4">{title}</h2>}
        </div>

        {/* Testimonials Grid */}
        {items.length > 0 ? (
          <div className={`grid gap-8 ${layout === 'grid' ? 'grid-cols-1 md:grid-cols-3' : 'grid-cols-1'}`}>
            {items.map((item: any, idx: number) => (
              <div key={idx} className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition">
                {/* Rating */}
                {item.rating && (
                  <div className="mb-3">
                    {'⭐'.repeat(item.rating)}
                  </div>
                )}

                {/* Quote */}
                <p className="text-gray-600 mb-4 italic">"{isArabic ? item.content_ar : item.content_en}"</p>

                {/* Author */}
                {item.author && <p className="font-bold text-gray-900">{item.author}</p>}
              </div>
            ))}
          </div>
        ) : (
          <div className="bg-white rounded-lg p-6 text-center text-gray-500">
            <p>{isArabic ? 'لا توجد تقييمات متاحة' : 'No testimonials available'}</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default Testimonials

// Component Definition for Registry
import { ComponentDefinition } from '../../../registry/types'
import manifest from '../manifest.json'

export const definition: ComponentDefinition = {
  id: 'testimonials',
  name: 'Testimonials Section',
  version: '1.0.0',
  component: Testimonials,
  mockData,
  props: {
    title_en: 'string',
    title_ar: 'string',
    isArabic: 'boolean',
    config: 'object',
  },
  manifest: manifest as any,
}
