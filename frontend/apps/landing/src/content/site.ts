import type { Locale } from '@pos-saas/i18n'
import english from '@/i18n/messages/en/marketing.json'
import arabic from '@/i18n/messages/ar/marketing.json'

export type { Locale } from '@pos-saas/i18n'
export type SiteCopy = typeof english

const copy: Record<Locale, SiteCopy> = { en: english, ar: arabic }

export function getSiteCopy(locale: Locale): SiteCopy {
  return copy[locale]
}
