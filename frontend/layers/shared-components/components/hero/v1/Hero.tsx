'use client'

import React from 'react'
import type { HeroProps } from './types'
import { mockData } from './mockData'

/**
 * Hero: Main hero banner component for websites
 *
 * Displays a full-width hero section with:
 * - Background image support
 * - Title, subtitle, and description
 * - Call-to-action button
 * - Configurable overlay opacity and text alignment
 * - RTL support for Arabic
 */
export const Hero: React.FC<HeroProps> = ({
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
  // Use provided values or fallback to mock data
  const title = isArabic ? (title_ar || mockData.title_ar) : (title_en || mockData.title_en)
  const subtitle = isArabic ? (subtitle_ar || mockData.subtitle_ar) : (subtitle_en || mockData.subtitle_en)
  const description = isArabic ? (description_ar || mockData.description_ar) : (description_en || mockData.description_en)
  const bgImage = background_image_url || mockData.background_image_url || ''

  const ctaText = config.cta_button_text || (isArabic ? 'اطلب الآن' : 'Learn More')
  const ctaUrl = config.cta_button_url || '#'
  const height = config.height === 'large' ? '600px' : config.height === 'medium' ? '500px' : config.height === 'small' ? '400px' : (config.height || '500px')
  const overlayOpacity = config.overlay_opacity !== undefined ? config.overlay_opacity : 0.3
  const textAlignment = config.text_alignment || 'center'
  const direction = isArabic ? 'rtl' : 'ltr'
  const backgroundValue = bgImage ? `url(${bgImage})` : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'

  return (
    <section className="hero-section relative w-full overflow-hidden flex items-center justify-center">
      {/* Overlay */}
      <div className="hero-overlay absolute inset-0 bg-black" />

      {/* Content */}
      <div className="hero-content relative z-10 max-w-4xl mx-auto px-6 w-full">
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
          <p className="hero-description text-lg text-gray-300 mb-8 max-w-2xl">
            {description}
          </p>
        )}

        {ctaText && (
          <a
            href={ctaUrl}
            className="btn-primary inline-flex items-center gap-2 text-lg"
          >
            {ctaText}
            <span aria-hidden="true">→</span>
          </a>
        )}
      </div>

      <style>{`
        .hero-section {
          height: ${height};
          background-image: ${backgroundValue};
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

export default Hero

// Component Definition for Registry
import { ComponentDefinition } from '../../../registry/types'
import manifest from '../manifest.json'

export const definition: ComponentDefinition = {
  id: 'hero',
  name: 'Hero Section',
  version: '1.0.0',
  component: Hero,
  mockData,
  props: {
    title_en: 'string',
    title_ar: 'string',
    subtitle_en: 'string',
    subtitle_ar: 'string',
    description_en: 'string',
    description_ar: 'string',
    background_image_url: 'string',
    isArabic: 'boolean',
    config: 'object',
  },
  manifest: manifest as any,
}