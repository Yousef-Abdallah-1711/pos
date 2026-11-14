import type { Metadata } from 'next'
import type { Locale } from '@/content/site'
export function marketingMetadata(locale: Locale, path: string, title: string, description: string): Metadata {
  return { title: `${title} | SERVIO`, description,
    alternates: { canonical: `/${locale}${path}`, languages: { en: `/en${path}`, ar: `/ar${path}` } },
    openGraph: { title: `${title} | SERVIO`, description, siteName: 'SERVIO', type: 'website', locale: locale === 'ar' ? 'ar_EG' : 'en_US' } }
}
