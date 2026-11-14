
import type { JSX } from 'react'
/**
 * BilingualContext
 * Manages bilingual state and language switching across the application
 */

'use client'

import React, { createContext, useContext, useCallback, useEffect } from 'react'
import type { Language } from '@/types/bilingual'

/**
 * Bilingual Context Type
 */
interface BilingualContextType {
  /** Current active language */
  currentLanguage: Language
  /** Set the current language */
  setCurrentLanguage: (language: Language) => void
  /** Supported languages */
  supportedLanguages: Language[]
  /** Check if current language is RTL */
  isRTL: boolean
  /** Get text direction attribute value */
  dir: 'rtl' | 'ltr'
}

/**
 * Create the context
 */
const BilingualContext = createContext<BilingualContextType | undefined>(undefined)

/**
 * Props for BilingualProvider
 */
interface BilingualProviderProps {
  children: React.ReactNode
  defaultLanguage?: Language
  supportedLanguages?: Language[]
}

/**
 * BilingualProvider Component
 * Provides bilingual context to the entire application
 */
export function BilingualProvider({
  children,
  defaultLanguage = 'en',
  supportedLanguages = ['en', 'ar'],
}: BilingualProviderProps): JSX.Element {
  const [currentLanguage, setLanguage] = React.useState<Language>(defaultLanguage)

  /**
   * Handle language change
   * Updates the HTML dir attribute and localStorage
   */
  const setCurrentLanguage = useCallback((language: Language) => {
    // Update state
    setLanguage(language)

    // Update HTML element dir attribute
    if (typeof document !== 'undefined') {
      document.documentElement.dir = language === 'ar' ? 'rtl' : 'ltr'
      document.documentElement.lang = language
    }

    // Save to localStorage
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('app-language', language)
    }
  }, [])

  /**
   * Initialize language from localStorage or system preference
   */
  useEffect(() => {
    if (typeof localStorage === 'undefined') return

    // Try to get saved language
    const savedLanguage = localStorage.getItem('app-language') as Language | null

    if (savedLanguage && supportedLanguages.includes(savedLanguage)) {
      setCurrentLanguage(savedLanguage)
    } else if (typeof navigator !== 'undefined') {
      // Check browser language preference
      const browserLang = navigator.language.split('-')[0].toLowerCase()
      if (browserLang === 'ar') {
        setCurrentLanguage('ar')
      }
    }
  }, [supportedLanguages, setCurrentLanguage])

  const isRTL = currentLanguage === 'ar'
  const dir = isRTL ? 'rtl' : 'ltr'

  const value: BilingualContextType = {
    currentLanguage,
    setCurrentLanguage,
    supportedLanguages,
    isRTL,
    dir,
  }

  return <BilingualContext.Provider value={value}>{children}</BilingualContext.Provider>
}

/**
 * Hook to use the BilingualContext
 */
export function useBilingual(): BilingualContextType {
  const context = useContext(BilingualContext)

  if (context === undefined) {
    throw new Error('useBilingual must be used within a BilingualProvider')
  }

  return context
}

/**
 * Hook to get only the current language
 */
export function useLanguage(): Language {
  return useBilingual().currentLanguage
}

/**
 * Hook to get language setter
 */
export function useSetLanguage(): (language: Language) => void {
  return useBilingual().setCurrentLanguage
}

/**
 * Hook to check if current language is RTL
 */
export function useIsRTL(): boolean {
  return useBilingual().isRTL
}

/**
 * Hook to get dir attribute value
 */
export function useDir(): 'rtl' | 'ltr' {
  return useBilingual().dir
}

export default BilingualContext
