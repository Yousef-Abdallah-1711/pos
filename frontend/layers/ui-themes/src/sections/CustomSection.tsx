'use client'

import React from 'react'
import type { CustomSectionProps } from '../types'

/**
 * CustomSection: Render custom HTML content
 */
export const CustomSection: React.FC<CustomSectionProps> = ({
  title,
  content = '',
  html_content = '',
  config = {},
  isArabic = false,
}) => {
  const contentToRender = html_content || content

  if (!contentToRender) {
    return (
      <section className="w-full py-16 px-4 bg-gray-100">
        <div className="max-w-4xl mx-auto text-center">
          <p className="text-gray-500">{isArabic ? 'Ù„Ø§ ÙŠÙˆØ¬Ø¯ Ù…Ø­ØªÙˆÙ‰' : 'No content'}</p>
        </div>
      </section>
    )
  }

  return (
    <section className="custom-section w-full py-16 px-4">
      <div className="max-w-4xl mx-auto">
        {title && <h2 className="text-4xl font-bold mb-6">{title}</h2>}

        <div dangerouslySetInnerHTML={{ __html: contentToRender }} className="prose max-w-none" />
      </div>

      <style>{`
        .custom-section {
          direction: ${isArabic ? 'rtl' : 'ltr'};
        }
      `}</style>
    </section>
  )
}

export default CustomSection
