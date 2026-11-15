import { defaultLocale, locales } from '@pos-saas/i18n'

export { defaultLocale, locales }
export type Locale = (typeof locales)[number]

export const localeNames: Record<Locale, string> = {
  en: 'English',
  ar: 'عربي',
}

export const rtlLocales = ['ar'] as const
export const isRTL = (locale: Locale) => rtlLocales.includes(locale as (typeof rtlLocales)[number])
