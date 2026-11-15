import { defaultLocale, locales } from '@pos-saas/i18n'

export const LOCALES = locales
export const DEFAULT_LOCALE = defaultLocale

export const LOCALE_NAMES = {
  en: 'English',
  ar: 'العربية',
};

export type Locale = (typeof LOCALES)[number];
