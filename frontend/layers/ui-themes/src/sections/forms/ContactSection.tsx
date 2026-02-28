'use client'

import React from 'react'
import type { ContactSectionProps } from '../../types'
import { getMockDataForComponent } from '../../mockData'

/**
 * ContactSection: Contact information and form
 */
export const ContactSection: React.FC<ContactSectionProps> = ({
  title_en = 'Contact Us',
  title_ar = 'Contact Us',
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

  const mockData = getMockDataForComponent('contact')
  const displayPhone = phone || mockData.phone || ''
  const displayEmail = email || mockData.email || ''
  const displayAddress = address || (isArabic ? mockData.address_ar : mockData.address_en) || ''

  return (
    <section
      dir={isArabic ? 'rtl' : 'ltr'}
      className="contact-section w-full py-12 sm:py-16 px-4 sm:px-6 bg-white dark:bg-gray-950"
    >
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-10 sm:mb-12">
          {title && <h2 className="text-3xl sm:text-4xl font-bold mb-4 text-gray-900 dark:text-gray-100">{title}</h2>}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-10 md:gap-12">
          <div className={isArabic ? 'text-right' : 'text-left'}>
            {displayPhone && (
              <div className="mb-6">
                <h3 className="font-bold mb-2 text-gray-900 dark:text-gray-100">Phone</h3>
                <a href={`tel:${displayPhone}`} className="text-primary hover:underline break-all">
                  {displayPhone}
                </a>
              </div>
            )}

            {displayEmail && (
              <div className="mb-6">
                <h3 className="font-bold mb-2 text-gray-900 dark:text-gray-100">Email</h3>
                <a href={`mailto:${displayEmail}`} className="text-primary hover:underline break-all">
                  {displayEmail}
                </a>
              </div>
            )}

            {displayAddress && (
              <div className="mb-6">
                <h3 className="font-bold mb-2 text-gray-900 dark:text-gray-100">Address</h3>
                <p className="text-gray-600 dark:text-gray-300 break-words">{displayAddress}</p>
              </div>
            )}
          </div>

          {showForm && (
            <div>
              <form className="space-y-4" onSubmit={(e) => e.preventDefault()}>
                <input
                  type="text"
                  placeholder={isArabic ? 'Your name' : 'Your name'}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:border-primary"
                />
                <input
                  type="email"
                  placeholder={isArabic ? 'Your email' : 'Your email'}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:border-primary"
                />
                <textarea
                  placeholder={isArabic ? 'Your message' : 'Your message'}
                  rows={5}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:border-primary"
                />
                <button type="submit" className="btn-primary w-full">
                  {isArabic ? 'Send' : 'Send'}
                </button>
              </form>
            </div>
          )}
        </div>

        {showMap && (
          <div className="mt-10 sm:mt-12 bg-gray-200 dark:bg-gray-800 rounded-lg h-72 sm:h-96 flex items-center justify-center text-center px-4">
            <p className="text-gray-500 dark:text-gray-400">Map will be displayed here</p>
          </div>
        )}
      </div>
    </section>
  )
}

export default ContactSection
