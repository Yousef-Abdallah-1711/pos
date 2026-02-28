'use client'

import React from 'react'
import { ChevronRight } from 'lucide-react'
import type { HeroSectionProps } from '../../types'
import { getMockDataForComponent } from '../../mockData'

/**
 * HeroSection: Main hero banner component for websites
 */
export const HeroSection: React.FC<HeroSectionProps> = ({
  title_en = '',
  title_ar = '',
  subtitle_en = '',
  subtitle_ar = '',
  description_en = '',
  description_ar = '',
  background_image_url,
  config = {},
  isArabic = false,
}) => {
  const mockData = getMockDataForComponent('hero')

  const title = isArabic ? title_ar || mockData.title_ar : title_en || mockData.title_en
  const subtitle = isArabic ? subtitle_ar || mockData.subtitle_ar : subtitle_en || mockData.subtitle_en
  const description = isArabic ? description_ar || mockData.description_ar : description_en || mockData.description_en
  const bgImage = background_image_url || mockData.background_image_url || ''

  const ctaText = config.cta_button_text || (isArabic ? 'Learn More' : 'Learn More')
  const ctaUrl = config.cta_button_url || '#'
  const height =
    config.height === 'large'
      ? 'min-h-[34rem]'
      : config.height === 'small'
        ? 'min-h-[24rem]'
        : 'min-h-[28rem]'

  const overlayOpacity = config.overlay_opacity !== undefined ? config.overlay_opacity : 0.3
  const textAlignment = config.text_alignment || 'center'

  const textAlignClass =
    textAlignment === 'left' ? 'text-left' : textAlignment === 'right' ? 'text-right' : 'text-center'

  const justifyClass =
    textAlignment === 'left'
      ? 'items-start'
      : textAlignment === 'right'
        ? 'items-end'
        : 'items-center'

  return (
    <section
      dir={isArabic ? 'rtl' : 'ltr'}
      className={`hero-section relative w-full overflow-hidden flex ${height}`}
      style={{
        backgroundImage: bgImage
          ? `linear-gradient(rgba(0,0,0,${overlayOpacity}), rgba(0,0,0,${overlayOpacity})), url(${bgImage})`
          : `linear-gradient(135deg, rgba(17,24,39,${Math.max(overlayOpacity, 0.7)}), rgba(59,130,246,${Math.max(overlayOpacity, 0.7)}))`,
        backgroundSize: 'cover',
        backgroundPosition: 'center',
      }}
    >
      <div className="relative z-10 w-full max-w-6xl mx-auto px-6 py-14 sm:py-20 flex">
        <div className={`w-full max-w-3xl flex flex-col gap-4 sm:gap-5 ${textAlignClass} ${justifyClass}`}>
          {title && <h1 className="text-3xl sm:text-5xl md:text-6xl font-bold text-white leading-tight break-words">{title}</h1>}

          {subtitle && <p className="text-lg sm:text-2xl text-gray-100 break-words">{subtitle}</p>}

          {description && <p className="text-base sm:text-lg text-gray-200 max-w-2xl break-words">{description}</p>}

          {ctaText && (
            <a
              href={ctaUrl}
              className={`btn-primary inline-flex items-center gap-2 text-base sm:text-lg w-fit ${isArabic ? 'flex-row-reverse' : ''}`}
            >
              {ctaText}
              <ChevronRight className={`h-5 w-5 ${isArabic ? 'rotate-180' : ''}`} />
            </a>
          )}
        </div>
      </div>
    </section>
  )
}

export default HeroSection
