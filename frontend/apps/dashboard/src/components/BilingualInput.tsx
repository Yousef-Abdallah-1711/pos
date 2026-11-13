
import type { JSX } from 'react'
/**
 * BilingualInput Component
 * Provides dual-language input with EN/AR tabs and RTL support
 */

import React, { useState, useCallback } from 'react'
import { BilingualText, Language, LANGUAGES, getTextByLanguage } from '@/types/bilingual'
import { bilingualFieldValidators } from '@/lib/validation/theme-schemas'

interface BilingualInputProps {
  /** Current bilingual text value */
  value: BilingualText
  /** Change handler */
  onChange: (value: BilingualText) => void
  /** Label for the field */
  label?: string
  /** Placeholder text (uses the language-specific placeholder) */
  placeholder?: string
  /** Whether the field is required */
  required?: boolean
  /** Whether the input is disabled */
  disabled?: boolean
  /** Maximum character length for each language */
  maxLength?: number
  /** Error message to display */
  error?: string
  /** Help text below the input */
  helpText?: string
  /** Supported languages (defaults to both EN and AR) */
  supportedLanguages?: Language[]
  /** Input type (text, textarea, etc.) */
  inputType?: 'text' | 'textarea' | 'email' | 'url'
  /** Additional CSS classes */
  className?: string
  /** On blur handler */
  onBlur?: (field: Language) => void
}

/**
 * BilingualInput Component
 * Allows editing of bilingual (EN/AR) text with proper RTL support
 */
export function BilingualInput({
  value,
  onChange,
  label,
  placeholder,
  required = false,
  disabled = false,
  maxLength = 500,
  error,
  helpText,
  supportedLanguages = ['en', 'ar'],
  inputType = 'text',
  className = '',
  onBlur,
}: BilingualInputProps): JSX.Element {
  // State for current active language and field-specific errors
  const [currentLanguage, setCurrentLanguage] = useState<Language>('en')
  const [fieldErrors, setFieldErrors] = useState<Record<Language, string | null>>({
    en: null,
    ar: null,
  })

  /**
   * Handle text change for the current language
   */
  const handleChange = useCallback(
    (newText: string) => {
      // Validate the input
      const validation =
        currentLanguage === 'en'
          ? bilingualFieldValidators.englishText(newText)
          : bilingualFieldValidators.arabicText(newText)

      setFieldErrors((prev) => ({
        ...prev,
        [currentLanguage]: validation.isValid ? null : validation.error,
      }))

      // Update the value
      onChange({
        ...value,
        [currentLanguage]: newText,
      })
    },
    [value, currentLanguage, onChange]
  )

  /**
   * Handle language tab change
   */
  const handleTabChange = useCallback((lang: Language) => {
    setCurrentLanguage(lang)
  }, [])

  /**
   * Handle blur event
   */
  const handleBlur = useCallback(() => {
    onBlur?.(currentLanguage)
  }, [currentLanguage, onBlur])

  // Get the current text based on active language
  const currentText = getTextByLanguage(value, currentLanguage)
  const currentLanguageConfig = LANGUAGES[currentLanguage]
  const isRTL = currentLanguageConfig.direction === 'rtl'
  const currentFieldError = fieldErrors[currentLanguage]

  return (
    <div className={`bilingual-input-wrapper ${className}`}>
      {/* Label */}
      {label && (
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}

      {/* Language Tabs */}
      <div className="flex border-b border-gray-200 dark:border-gray-700 mb-4 bg-gray-50 dark:bg-gray-800 rounded-t-lg">
        {supportedLanguages.map((lang) => {
          const langConfig = LANGUAGES[lang]
          const isActive = currentLanguage === lang
          const hasError = fieldErrors[lang] !== null

          return (
            <button
              key={lang}
              onClick={() => handleTabChange(lang)}
              disabled={disabled}
              className={`
                flex-1 px-4 py-3 text-center font-medium transition-colors
                ${isActive
                  ? 'border-b-2 border-blue-500 text-blue-600 dark:text-blue-400 bg-white dark:bg-gray-700'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'}
                ${hasError ? 'border-red-500' : ''}
                ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
              `}
            >
              <div className="flex items-center justify-center gap-2">
                <span>{langConfig.nativeName}</span>
                {hasError && (
                  <span className="w-2 h-2 bg-red-500 rounded-full" title={fieldErrors[lang] || ''} />
                )}
              </div>
            </button>
          )
        })}
      </div>

      {/* Input Container */}
      <div className={`relative ${isRTL ? 'dir-rtl' : 'dir-ltr'}`}>
        {inputType === 'textarea' ? (
          <textarea
            value={currentText}
            onChange={(e) => handleChange(e.target.value)}
            onBlur={handleBlur}
            placeholder={placeholder}
            disabled={disabled}
            maxLength={maxLength}
            dir={isRTL ? 'rtl' : 'ltr'}
            className={`
              w-full px-4 py-3 border rounded-lg font-medium
              placeholder-gray-400 dark:placeholder-gray-500
              dark:bg-gray-700 dark:text-white
              focus:outline-none focus:ring-2 transition-all
              resize-none h-32
              ${
                currentFieldError || error
                  ? 'border-red-500 focus:ring-red-200 dark:focus:ring-red-900'
                  : 'border-gray-300 dark:border-gray-600 focus:ring-blue-200 dark:focus:ring-blue-900'
              }
              ${disabled ? 'opacity-50 cursor-not-allowed bg-gray-100 dark:bg-gray-800' : ''}
            `}
          />
        ) : (
          <input
            type={inputType}
            value={currentText}
            onChange={(e) => handleChange(e.target.value)}
            onBlur={handleBlur}
            placeholder={placeholder}
            disabled={disabled}
            maxLength={maxLength}
            dir={isRTL ? 'rtl' : 'ltr'}
            className={`
              w-full px-4 py-3 border rounded-lg font-medium
              placeholder-gray-400 dark:placeholder-gray-500
              dark:bg-gray-700 dark:text-white
              focus:outline-none focus:ring-2 transition-all
              ${
                currentFieldError || error
                  ? 'border-red-500 focus:ring-red-200 dark:focus:ring-red-900'
                  : 'border-gray-300 dark:border-gray-600 focus:ring-blue-200 dark:focus:ring-blue-900'
              }
              ${disabled ? 'opacity-50 cursor-not-allowed bg-gray-100 dark:bg-gray-800' : ''}
            `}
          />
        )}

        {/* Character Counter */}
        <div className={`absolute bottom-3 text-xs text-gray-500 dark:text-gray-400 ${isRTL ? 'left-4' : 'right-4'}`}>
          {currentText.length} / {maxLength}
        </div>
      </div>

      {/* Error Message */}
      {(currentFieldError || error) && (
        <div className="mt-2 text-sm text-red-600 dark:text-red-400 flex items-start gap-2">
          <span className="text-red-500 mt-0.5">!</span>
          <span>{currentFieldError || error}</span>
        </div>
      )}

      {/* Help Text */}
      {helpText && !error && !currentFieldError && (
        <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">{helpText}</p>
      )}

      {/* Language Info */}
      <div className="mt-3 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
        <p className="text-xs text-gray-600 dark:text-gray-300 flex items-center gap-2">
          <span className="text-blue-600 dark:text-blue-400">ℹ</span>
          Editing <strong>{LANGUAGES[currentLanguage].name}</strong>
          {isRTL && ' (Right-to-Left)'}
          {!isRTL && ' (Left-to-Right)'}
        </p>
      </div>

      {/* Language Summary */}
      <div className="mt-4 grid grid-cols-2 gap-2">
        {supportedLanguages.map((lang) => {
          const langConfig = LANGUAGES[lang]
          const langText = getTextByLanguage(value, lang)
          const langError = fieldErrors[lang]

          return (
            <div
              key={lang}
              className={`
                p-2 rounded border text-xs
                ${langError ? 'border-red-200 bg-red-50 dark:bg-red-900/20' : 'border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800'}
              `}
            >
              <div className="font-medium text-gray-700 dark:text-gray-300 mb-1">
                {langConfig.nativeName} ({lang.toUpperCase()})
              </div>
              <div className="text-gray-600 dark:text-gray-400 truncate" dir={langConfig.direction}>
                {langText || <em className="text-gray-400">Not set</em>}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default BilingualInput
