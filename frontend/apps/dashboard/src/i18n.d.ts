import type englishMessages from './i18n/messages/en.json'
import type subscription from '@pos-saas/i18n/locales/en/subscription.json'

declare module 'next-intl' {
  interface AppConfig {
    Locale: 'en' | 'ar'
    Messages: typeof englishMessages & { subscription: typeof subscription }
  }
}
