'use client'

import React from 'react'
import type { TestimonialsSectionProps } from '../types'
import { getMockDataForComponent } from '../mockData'

/**
 * TestimonialsSection: Customer testimonials and reviews
 */
export const TestimonialsSection: React.FC<TestimonialsSectionProps> = ({
  title_en = 'What Our Customers Say',
  title_ar = 'What Our Customers Say',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const layout = config.layout || 'grid'

  const mockData = getMockDataForComponent('testimonials')
  const items = config.items && config.items.length > 0 ? config.items : mockData.items || []

  return (
    <section
      dir={isArabic ? 'rtl' : 'ltr'}
      className="testimonials-section w-full py-12 sm:py-16 px-4 sm:px-6 bg-gray-50 dark:bg-gray-900"
    >
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-10 sm:mb-12">
          {title && <h2 className="text-3xl sm:text-4xl font-bold mb-4 text-gray-900 dark:text-gray-100">{title}</h2>}
        </div>

        {items.length > 0 ? (
          <div className={`grid gap-6 sm:gap-8 ${layout === 'grid' ? 'grid-cols-1 md:grid-cols-3' : 'grid-cols-1'}`}>
            {items.map((item: any, idx: number) => (
              <article
                key={idx}
                className={`rounded-xl border border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 shadow-sm hover:shadow-md transition ${
                  isArabic ? 'text-right' : 'text-left'
                }`}
              >
                {item.rating && (
                  <div className="mb-3 text-amber-500" aria-label={`Rating ${item.rating} of 5`}>
                    {'*'.repeat(item.rating)}
                  </div>
                )}

                <p className="text-gray-600 dark:text-gray-300 mb-4 italic break-words">
                  "{isArabic ? item.content_ar : item.content_en}"
                </p>

                {item.author && <p className="font-bold text-gray-900 dark:text-gray-100 break-words">{item.author}</p>}
              </article>
            ))}
          </div>
        ) : (
          <div className="rounded-xl border border-gray-100 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 text-center text-gray-500 dark:text-gray-400">
            <p>No testimonials available</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default TestimonialsSection
