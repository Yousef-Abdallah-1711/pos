import 'next-intl'
import type enMessages from './i18n/messages/en.json'
import type enSubscription from '@pos-saas/i18n/locales/en/subscription.json'

declare module 'next-intl' {
  interface AppConfig {
    Locale: 'en' | 'ar'
    Messages: typeof enMessages & { subscription: typeof enSubscription }
  }
}
