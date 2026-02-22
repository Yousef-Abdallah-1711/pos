'use client'

import { createContext, useContext, useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { translations, type Language, type Translations } from '@/lib/translations'

interface LanguageContextType {
  language: Language
  setLanguage: (lang: Language) => void
  t: Translations
  dir: 'ltr' | 'rtl'
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined)

interface LanguageProviderProps {
  children: React.ReactNode
  locale: Language
}

export function LanguageProvider({ children, locale }: LanguageProviderProps) {
  const router = useRouter()
  const pathname = usePathname()
  const dir = locale === 'ar' ? 'rtl' : 'ltr'

  useEffect(() => {
    // Update HTML attributes for RTL support
    document.documentElement.dir = dir
    document.documentElement.lang = locale
  }, [locale, dir])

  const setLanguage = (lang: Language) => {
    // Swap locale segment in the current URL
    const segments = pathname.split('/')
    // segments[0] is '', segments[1] is the current locale
    segments[1] = lang
    const newPath = segments.join('/')
    router.push(newPath)
  }

  return (
    <LanguageContext.Provider
      value={{
        language: locale,
        setLanguage,
        t: translations[locale],
        dir,
      }}
    >
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const context = useContext(LanguageContext)
  if (context === undefined) {
    throw new Error('useLanguage must be used within a LanguageProvider')
  }
  return context
}
