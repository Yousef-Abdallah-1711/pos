import 'next-intl'

declare module 'next-intl' {
  interface AppConfig {
    Locale: 'en' | 'ar'
  }
}
