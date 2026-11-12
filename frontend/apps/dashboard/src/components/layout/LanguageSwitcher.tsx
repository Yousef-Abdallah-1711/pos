'use client'

import { useRouter, usePathname } from 'next/navigation'
import { Globe, Check } from 'lucide-react'
import { useState, useRef, useEffect } from 'react'
import { usePreferencesStore } from '@/stores/preferencesStore'

const LANGUAGE_NAMES: Record<string, string> = {
  'en': 'English',
  'ar': 'عربي'
}

export function LanguageSwitcher() {
  const router = useRouter()
  const pathname = usePathname()
  // Extract locale from pathname (e.g., /en/dashboard -> 'en')
  const locale = (pathname.split('/')[1] || 'en') as 'en' | 'ar'
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const languages = [
    { code: 'en' },
    { code: 'ar' },
  ]

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleLanguageChange = (newLocale: string) => {
    if (newLocale === locale) return // Don't switch if same language

    // Update global preferences store FIRST (this applies dir and persists)
    const { settings, setSettings } = usePreferencesStore.getState()
    if (settings) {
      setSettings({ ...settings, language: newLocale as 'en' | 'ar' })
    }

    // Replace the locale in the current pathname
    const pathParts = pathname.split('/')
    pathParts[1] = newLocale
    const newPathname = pathParts.join('/')

    // Push to new locale URL
    router.push(newPathname)
    setIsOpen(false)
  }

  const isRTL = locale === 'ar'

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        aria-label={locale === 'ar' ? 'تغيير اللغة' : 'Change language'}
        className="flex items-center gap-2 p-2 sm:px-3 sm:py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
      >
        <Globe className="w-4 h-4" />
        <span className="hidden text-sm font-medium sm:inline">{LANGUAGE_NAMES[locale]}</span>
      </button>

      {isOpen && (
        <div className={`absolute ${isRTL ? 'left-0' : 'right-0'} mt-2 w-48 bg-white dark:bg-gray-900 rounded-lg shadow-lg border border-gray-200 dark:border-gray-800 z-50`}>
          {languages.map((lang) => (
            <button
              key={lang.code}
              onClick={() => handleLanguageChange(lang.code)}
              className="w-full flex items-center justify-between px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
            >
              <span className={`text-sm ${locale === lang.code ? 'font-semibold text-primary-600' : 'text-gray-700 dark:text-gray-300'}`}>
                {LANGUAGE_NAMES[lang.code]}
              </span>
              {locale === lang.code && <Check className="w-4 h-4 text-primary-600" />}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
