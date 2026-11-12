import type { MetadataRoute } from 'next'
import { getLandingOrigin } from '@/lib/seo'

export default function sitemap(): MetadataRoute.Sitemap {
  const origin = getLandingOrigin()

  return (['en', 'ar'] as const).map((locale) => ({
    url: new URL(`/${locale}`, origin).toString(),
    alternates: {
      languages: {
        en: new URL('/en', origin).toString(),
        ar: new URL('/ar', origin).toString(),
      },
    },
  }))
}
