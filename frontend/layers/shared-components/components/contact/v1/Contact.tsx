'use client'

import React from 'react'
import type { ContactProps } from './types'
import { mockData } from './mockData'

/**
 * Contact: Contact information and form
 *
 * Displays:
 * - Contact form
 * - Phone and email
 * - Address
 * - Optional map integration
 * - RTL support
 */
export const Contact: React.FC<ContactProps> = ({
  title_en = 'Contact Us',
  title_ar = 'اتصل بنا',
  phone = '',
  email = '',
  address_en = '',
  address_ar = '',
  config = {},
  isArabic = false,
}) => {
  const title = isArabic ? title_ar : title_en
  const address = isArabic ? address_ar : address_en
  const showForm = config.show_form !== false
  const showMap = config.show_map !== false

  // Use provided values or fallback to mock data
  const displayPhone = phone || mockData.phone || ''
  const displayEmail = email || mockData.email || ''
  const displayAddress = address || (isArabic ? mockData.address_ar : mockData.address_en) || ''

  return (
    <section
      className="w-full py-16 px-4 bg-white"
      dir={isArabic ? 'rtl' : 'ltr'}
    >
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="text-center mb-12">
          {title && <h2 className="text-4xl font-bold mb-4">{title}</h2>}
        </div>

        <div className="grid md:grid-cols-2 gap-12">
          {/* Contact Info */}
          <div>
            {displayPhone && (
              <div className="mb-6">
                <h3 className="font-bold mb-2">{isArabic ? 'الهاتف' : 'Phone'}</h3>
                <a href={`tel:${displayPhone}`} className="text-primary hover:underline">
                  {displayPhone}
                </a>
              </div>
            )}

            {displayEmail && (
              <div className="mb-6">
                <h3 className="font-bold mb-2">{isArabic ? 'البريد الإلكتروني' : 'Email'}</h3>
                <a href={`mailto:${displayEmail}`} className="text-primary hover:underline">
                  {displayEmail}
                </a>
              </div>
            )}

            {displayAddress && (
              <div className="mb-6">
                <h3 className="font-bold mb-2">{isArabic ? 'العنوان' : 'Address'}</h3>
                <p className="text-gray-600">{displayAddress}</p>
              </div>
            )}
          </div>

          {/* Contact Form */}
          {showForm && (
            <div>
              <form className="space-y-4">
                <input
                  type="text"
                  placeholder={isArabic ? 'اسمك' : 'Your name'}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-primary"
                />
                <input
                  type="email"
                  placeholder={isArabic ? 'بريدك الإلكتروني' : 'Your email'}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-primary"
                />
                <textarea
                  placeholder={isArabic ? 'رسالتك' : 'Your message'}
                  rows={5}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-primary"
                />
                <button type="submit" className="btn-primary w-full">
                  {isArabic ? 'إرسال' : 'Send'}
                </button>
              </form>
            </div>
          )}
        </div>

        {/* Map */}
        {showMap && (
          <div className="mt-12 bg-gray-200 rounded-lg h-96 flex items-center justify-center">
            <p className="text-gray-500">{isArabic ? 'الخريطة ستظهر هنا' : 'Map will be displayed here'}</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default Contact

// Component Definition for Registry
import { ComponentDefinition } from '../../../registry/types'
import manifest from '../manifest.json'

export const definition: ComponentDefinition = {
  id: 'contact',
  name: 'Contact Section',
  version: '1.0.0',
  component: Contact,
  mockData,
  props: {
    title_en: 'string',
    title_ar: 'string',
    phone: 'string',
    email: 'string',
    address_en: 'string',
    address_ar: 'string',
    isArabic: 'boolean',
    config: 'object',
  },
  manifest: manifest as any,
}
