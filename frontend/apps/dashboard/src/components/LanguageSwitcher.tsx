
import type { JSX } from 'react'
/**
 * LanguageSwitcher Component
 * Allows users to switch between supported languages
 */

'use client'

import React, { useState, useRef, useEffect } from 'react'
import { useBilingual } from '@/context/BilingualContext'
import { LANGUAGES, type Language } from '@/types/bilingual'

interface LanguageSwitcherProps {
  /** Variant style (button, dropdown, tabs) */
  variant?: 'button' | 'dropdown' | 'tabs'
  /** Position for dropdown (top, bottom, left, right) */
  dropdownPosition?: 'top' | 'bottom' | 'left' | 'right'
  /** Additional CSS classes */
  className?: string
  /** Show language names (default: true) */
  showNames?: boolean
  /** Show native language names (default: false) */
  showNativeNames?: boolean
  /** Size of the switcher (sm, md, lg) */
  size?: 'sm' | 'md' | 'lg'
}

/**
 * LanguageSwitcher Component
 */
export function LanguageSwitcher({
  variant = 'dropdown',
  dropdownPosition = 'bottom',
  className = '',
  showNames = true,
  showNativeNames = false,
  size = 'md',
}: LanguageSwitcherProps): JSX.Element {
  const { currentLanguage, setCurrentLanguage, supportedLanguages, isRTL } = useBilingual()
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  /**
   * Handle language change
   */
  const handleLanguageChange = (lang: Language) => {
    setCurrentLanguage(lang)
    setIsOpen(false)
  }

  /**
   * Close dropdown when clicking outside
   */
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const sizeClasses = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-2 text-sm',
    lg: 'px-4 py-2 text-base',
  }

  if (variant === 'tabs') {
    return (
      <div className={`flex gap-2 ${className}`}>
        {supportedLanguages.map((lang) => {
          const langConfig = LANGUAGES[lang]
          const isActive = currentLanguage === lang

          return (
            <button
              key={lang}
              onClick={() => handleLanguageChange(lang)}
              className={`
                ${sizeClasses[size]}
                rounded-lg font-medium transition-all
                ${
                  isActive
                    ? 'bg-blue-600 text-white shadow-lg'
                    : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
                }
              `}
              title={langConfig.name}
            >
              <span>{lang.toUpperCase()}</span>
              {showNativeNames && <span className="text-xs block">{langConfig.nativeName}</span>}
            </button>
          )
        })}
      </div>
    )
  }

  if (variant === 'dropdown') {
    return (
      <div ref={dropdownRef} className={`relative inline-block ${className}`}>
        {/* Dropdown Button */}
        <button
          onClick={() => setIsOpen(!isOpen)}
          className={`
            ${sizeClasses[size]}
            flex items-center gap-2
            bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300
            border border-gray-300 dark:border-gray-600
            rounded-lg font-medium
            hover:bg-gray-50 dark:hover:bg-gray-600
            focus:outline-none focus:ring-2 focus:ring-blue-500
            transition-all
          `}
        >
          <span className="font-bold">{currentLanguage.toUpperCase()}</span>
          <svg
            className={`w-4 h-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
          </svg>
        </button>

        {/* Dropdown Menu */}
        {isOpen && (
          <div
            className={`
              absolute z-50 min-w-max
              bg-white dark:bg-gray-700
              border border-gray-300 dark:border-gray-600
              rounded-lg shadow-lg
              py-2
              ${dropdownPosition === 'bottom' ? 'top-full mt-2' : ''}
              ${dropdownPosition === 'top' ? 'bottom-full mb-2' : ''}
              ${isRTL ? 'right-0' : 'left-0'}
            `}
          >
            {supportedLanguages.map((lang) => {
              const langConfig = LANGUAGES[lang]
              const isActive = currentLanguage === lang

              return (
                <button
                  key={lang}
                  onClick={() => handleLanguageChange(lang)}
                  className={`
                    w-full text-${isRTL ? 'right' : 'left'} px-4 py-2
                    flex items-center gap-3
                    font-medium transition-colors
                    ${
                      isActive
                        ? 'bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200'
                        : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600'
                    }
                  `}
                >
                  {/* Language flag or icon */}
                  <span className="text-lg">
                    {lang === 'en' ? '🇺🇸' : lang === 'ar' ? '🇸🇦' : '🌍'}
                  </span>

                  <div className="flex-1">
                    <div className="font-medium">{showNames ? langConfig.name : lang.toUpperCase()}</div>
                    {showNativeNames && <div className="text-xs opacity-75">{langConfig.nativeName}</div>}
                  </div>

                  {/* Active indicator */}
                  {isActive && (
                    <svg className="w-5 h-5 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                      <path
                        fillRule="evenodd"
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                        clipRule="evenodd"
                      />
                    </svg>
                  )}
                </button>
              )
            })}
          </div>
        )}
      </div>
    )
  }

  // Default: button variant
  return (
    <button
      onClick={() => handleLanguageChange(supportedLanguages.find((l) => l !== currentLanguage) || supportedLanguages[0]!)}
      className={`
        ${sizeClasses[size]}
        inline-flex items-center gap-2
        bg-blue-600 hover:bg-blue-700 text-white
        rounded-lg font-medium
        focus:outline-none focus:ring-2 focus:ring-blue-500
        transition-all
        ${className}
      `}
    >
      <span>
        {currentLanguage === 'en' ? '🇺🇸' : currentLanguage === 'ar' ? '🇸🇦' : '🌍'}
      </span>
      {showNames && <span>{LANGUAGES[currentLanguage].name}</span>}
    </button>
  )
}

export default LanguageSwitcher
