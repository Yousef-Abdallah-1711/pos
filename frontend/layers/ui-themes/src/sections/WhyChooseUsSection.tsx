'use client'

import React from 'react'
import type { WhyChooseUsSectionProps } from '../types'
import { getMockDataForComponent } from '../mockData'

/**
 * WhyChooseUsSection: Highlight reasons to choose your business
 */
export const WhyChooseUsSection: React.FC<WhyChooseUsSectionProps> = ({
  title_en = 'Why Choose Us',
  title_ar = 'Why Choose Us',
  description_en = '',
  description_ar = '',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const description = isArabic ? description_ar : description_en
  const layout = config.layout || 'grid'
  const columns = config.columns || 3

  const mockData = getMockDataForComponent('why_us')
  const items = config.items && config.items.length > 0 ? config.items : mockData.items || []

  const gridClass =
    layout === 'grid'
      ? columns === 3
        ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
        : columns === 2
          ? 'grid-cols-1 sm:grid-cols-2'
          : 'grid-cols-1'
      : 'grid-cols-1'

  return (
    <section
      dir={isArabic ? 'rtl' : 'ltr'}
      className="why-us-section w-full py-12 sm:py-16 px-4 sm:px-6 bg-gray-50 dark:bg-gray-900"
    >
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-10 sm:mb-12">
          {title && <h2 className="text-3xl sm:text-4xl font-bold mb-4 text-gray-900 dark:text-gray-100">{title}</h2>}
          {description && <p className="text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">{description}</p>}
        </div>

        {items.length > 0 ? (
          <div className={`grid gap-6 sm:gap-8 ${gridClass}`}>
            {items.map((item: any, idx: number) => (
              <article
                key={idx}
                className={`rounded-xl border border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 shadow-sm hover:shadow-md transition ${
                  isArabic ? 'text-right' : 'text-left'
                }`}
              >
                {item.icon && <div className="text-3xl sm:text-4xl mb-4">{item.icon}</div>}

                <h3 className="text-xl font-bold mb-2 text-primary break-words">
                  {isArabic ? item.title_ar : item.title_en}
                </h3>

                <p className="text-gray-600 dark:text-gray-300 break-words">
                  {isArabic ? item.description_ar : item.description_en}
                </p>
              </article>
            ))}
          </div>
        ) : (
          <div className="rounded-xl border border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 text-center text-gray-500 dark:text-gray-400">
            <p>No features available</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default WhyChooseUsSection
