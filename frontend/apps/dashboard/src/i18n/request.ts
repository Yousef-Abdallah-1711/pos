import { getRequestConfig } from 'next-intl/server'
import { defaultLocale, locales, type Locale } from '@/i18n/config'
import subscriptionEn from '@pos-saas/i18n/locales/en/subscription.json'
import subscriptionAr from '@pos-saas/i18n/locales/ar/subscription.json'

function resolveLocale(input: string | undefined): Locale {
  if (!input) {
    return defaultLocale
  }
  return locales.includes(input as Locale) ? (input as Locale) : defaultLocale
}

export default getRequestConfig(async ({ requestLocale }) => {
  const locale = resolveLocale(await requestLocale)
  const appMessages =
    locale === 'ar'
      ? (await import('@/i18n/messages/ar.json')).default
      : (await import('@/i18n/messages/en.json')).default
  const subscription = locale === 'ar' ? subscriptionAr : subscriptionEn

  return {
    locale,
    messages: { ...appMessages, subscription },
  }
})
