import { getRequestConfig } from 'next-intl/server'
import enMessages from '@/i18n/messages/en.json'
import arMessages from '@/i18n/messages/ar.json'
import { DEFAULT_LOCALE, LOCALES, type Locale } from '@/i18n/config'

const messages: Record<Locale, typeof enMessages> = {
  en: enMessages,
  ar: arMessages,
}

export default getRequestConfig(async ({ requestLocale }) => {
  const requested = await requestLocale
  const locale = LOCALES.includes(requested as Locale)
    ? (requested as Locale)
    : DEFAULT_LOCALE

  return {
    locale,
    messages: messages[locale],
  }
})

