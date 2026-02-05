import { LOCALES, DEFAULT_LOCALE, type Locale } from '@/i18n/config';

// Import translations
import en from '@/i18n/messages/en.json';
import ar from '@/i18n/messages/ar.json';

const messages: Record<Locale, any> = {
  en,
  ar,
};

/**
 * Get locale from pathname
 * @param pathname - Current URL pathname (e.g., "/en/menu" or "/ar/orders")
 * @returns Locale ('en' or 'ar')
 */
export function getLocaleFromPath(pathname: string): Locale {
  const segments = pathname.split('/').filter(Boolean);
  const locale = (segments[0] || '').toLowerCase();
  return (LOCALES.includes(locale as Locale) ? locale : DEFAULT_LOCALE) as Locale;
}

/**
 * Get direction (ltr/rtl) from locale
 */
export function getDirectionFromLocale(locale: Locale): 'ltr' | 'rtl' {
  return locale === 'ar' ? 'rtl' : 'ltr';
}

/**
 * Create translator function for specific locale
 * @param locale - Locale to translate for
 * @returns Function that translates keys like 'common.menu'
 */
export function createTranslator(locale: Locale) {
  const localeMessages = messages[locale] || messages[DEFAULT_LOCALE];

  return function t(key: string, defaultValue?: string): string {
    const keys = key.split('.');
    let value: any = localeMessages;

    for (const k of keys) {
      value = value?.[k];
    }

    if (typeof value === 'string') {
      return value;
    }

    // Return default or the key itself if not found
    return defaultValue || key;
  };
}

/**
 * Get all translations for a locale
 */
export function getMessages(locale: Locale) {
  return messages[locale] || messages[DEFAULT_LOCALE];
}

/**
 * Export supported locales
 */
export { LOCALES, DEFAULT_LOCALE, type Locale };
