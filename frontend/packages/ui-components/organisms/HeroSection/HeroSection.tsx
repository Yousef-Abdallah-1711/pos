'use client'

import React from 'react'
import { ChevronRight } from 'lucide-react'
import type { HeroSectionProps } from './HeroSection.types'

/**
 * HeroSection: Main hero banner component for websites
 *
 * Displays a full-width hero section with:
 * - Background image support
 * - Title, subtitle, and description (bilingual)
 * - Call-to-action button
 * - Configurable overlay opacity and text alignment
 * - RTL support for Arabic
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
  const title = isArabic ? title_ar : title_en
  const subtitle = isArabic ? subtitle_ar : subtitle_en
  const description = isArabic ? description_ar : description_en
  const ctaText = config.cta_button_text || (isArabic ? 'اطلب الآن' : 'Learn More')
  const ctaUrl = config.cta_button_url || '#'
  const height = config.height || '500px'
  const overlayOpacity = config.overlay_opacity || 0.3
  const textAlignment = config.text_alignment || 'center'
  const direction = isArabic ? 'rtl' : 'ltr'
  const backgroundImage = background_image_url ? `url(${background_image_url})` : 'none'

  return (
    <section className="hero-section-component relative w-full overflow-hidden flex items-center justify-center">
      {/* Overlay */}
      <div className="hero-overlay absolute inset-0 bg-black" />

      {/* Content */}
      <div className="hero-content relative z-10 max-w-4xl mx-auto px-6 text-center w-full">
        {title && (
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
            {title}
          </h1>
        )}

        {subtitle && (
          <p className="text-xl md:text-2xl text-gray-200 mb-6">
            {subtitle}
          </p>
        )}

        {description && (
          <p className="hero-description text-lg text-gray-300 mb-8 max-w-2xl mx-auto">
            {description}
          </p>
        )}

        {ctaText && (
          <a
            href={ctaUrl}
            className="btn-primary inline-flex items-center gap-2 text-lg"
          >
            {ctaText}
            <ChevronRight className="h-5 w-5" />
          </a>
        )}
      </div>

      <style>{`
        .hero-section-component {
          height: ${height};
          background-image: ${backgroundImage};
          background-size: cover;
          background-position: center;
          direction: ${direction};
        }

        .hero-overlay {
          opacity: ${overlayOpacity};
        }

        .hero-content {
          text-align: ${textAlignment};
        }

        .hero-description {
          margin-left: ${textAlignment === 'center' ? 'auto' : '0'};
          margin-right: ${textAlignment === 'center' ? 'auto' : '0'};
        }
      `}</style>
    </section>
  )
}

export default HeroSection