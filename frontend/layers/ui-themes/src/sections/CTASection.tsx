'use client'

import React from 'react'
import type { CTASectionProps } from '../types'
import { getMockDataForComponent } from '../mockData'

/**
 * CTASection: Call-to-action section
 */
export const CTASection: React.FC<CTASectionProps> = ({
  title_en = 'Ready to Get Started?',
  title_ar = 'Ready to Get Started?',
  description_en = '',
  description_ar = '',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const description = isArabic ? description_ar : description_en

  const mockData = getMockDataForComponent('cta')
  const displayTitle = title || (isArabic ? mockData.title_ar : mockData.title_en) || ''
  const displayDescription = description || (isArabic ? mockData.description_ar : mockData.description_en) || ''

  const buttonText =
    config.button_text_en || config.button_text_ar
      ? isArabic
        ? config.button_text_ar
        : config.button_text_en
      : isArabic
        ? mockData.button_text_ar || 'Order Now'
        : mockData.button_text_en || 'Order Now'

  const buttonUrl = config.button_url || mockData.button_url || '#'
  const backgroundColor = config.background_color || '#3b82f6'

  return (
    <section
      dir={isArabic ? 'rtl' : 'ltr'}
      className="cta-section w-full py-14 sm:py-20 px-4 sm:px-6 text-white"
      style={{ backgroundColor }}
    >
      <div className="max-w-4xl mx-auto text-center">
        {displayTitle && <h2 className="text-3xl sm:text-4xl font-bold mb-4 break-words">{displayTitle}</h2>}

        {displayDescription && <p className="text-base sm:text-lg mb-8 opacity-90 break-words">{displayDescription}</p>}

        {buttonText && (
          <a href={buttonUrl} className="btn-accent inline-block text-base sm:text-lg font-semibold w-full sm:w-auto text-center">
            {buttonText}
          </a>
        )}
      </div>
    </section>
  )
}

export default CTASection
