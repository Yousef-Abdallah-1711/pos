'use client'

import React from 'react'
import type { WhyUsProps } from './types'
import { mockData } from './mockData'

/**
 * WhyUs: Highlight reasons to choose your business
 *
 * Features:
 * - Multiple items with icons/titles/descriptions
 * - Grid or flex layouts
 * - Customizable items
 * - RTL support
 */
export const WhyUs: React.FC<WhyUsProps> = ({
  title_en = 'Why Choose Us',
  title_ar = 'لماذا اختياراتنا',
  description_en = '',
  description_ar = '',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const description = isArabic ? description_ar : description_en
  const layout = config.layout || 'grid'
  const columns = config.columns || 3

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
          {description && <p className="text-gray-600 max-w-2xl mx-auto">{description}</p>}
        </div>

        {/* Features Grid */}
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
            {items.map((item: any, idx: number) => (
              <div key={idx} className="bg-white rounded-lg p-6 shadow-sm hover:shadow-md transition text-center">
                {/* Icon */}
                {item.icon && (
                  <div className="text-4xl mb-4">
                    {item.icon}
                  </div>
                )}

                {/* Title */}
                <h3 className="text-xl font-bold mb-2 text-primary">
                  {isArabic ? item.title_ar : item.title_en}
                </h3>

                {/* Description */}
                <p className="text-gray-600">
                  {isArabic ? item.description_ar : item.description_en}
                </p>
              </div>
            ))}
          </div>
        ) : (
          <div className="bg-white rounded-lg p-6 text-center text-gray-500">
            <p>{isArabic ? 'لا توجد ميزات متاحة' : 'No features available'}</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default WhyUs

// Component Definition for Registry
import { ComponentDefinition } from '../../../registry/types'
import manifest from '../manifest.json'

export const definition: ComponentDefinition = {
  id: 'why_us',
  name: 'Why Us Section',
  version: '1.0.0',
  component: WhyUs,
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
