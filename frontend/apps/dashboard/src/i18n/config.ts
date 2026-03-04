export const defaultLocale = 'en' as const
export const locales = ['en', 'ar'] as const

export type Locale = (typeof locales)[number]

export const localeNames: Record<Locale, string> = {
  en: 'English',
  ar: 'عربي',
}

export const rtlLocales = ['ar'] as const
export const isRTL = (locale: Locale) => rtlLocales.includes(locale as any)
