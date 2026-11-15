import english from '@/i18n/messages/en/legacy.json'
import arabic from '@/i18n/messages/ar/legacy.json'
import type { Locale } from '@pos-saas/i18n'

export const translations = { en: english, ar: arabic }
export type Language = Locale
export type Translations = typeof english
