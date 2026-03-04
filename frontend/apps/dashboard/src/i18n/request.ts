import { getRequestConfig } from 'next-intl/server'
import { defaultLocale, locales, type Locale } from '@/i18n/config'

function resolveLocale(input: string | undefined): Locale {
  if (!input) {
    return defaultLocale
  }
  return locales.includes(input as Locale) ? (input as Locale) : defaultLocale
}

export default getRequestConfig(async ({ requestLocale }) => {
  const locale = resolveLocale(await requestLocale)
  const messages =
    locale === 'ar'
      ? (await import('@/i18n/messages/ar.json')).default
      : (await import('@/i18n/messages/en.json')).default

  return {
    locale,
    messages,
  }
})
