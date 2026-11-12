import type { MetadataRoute } from 'next'
import { getLandingOrigin } from '@/lib/seo'

export default function robots(): MetadataRoute.Robots {
  const origin = getLandingOrigin()

  return {
    rules: [{ userAgent: '*', allow: '/', disallow: ['/en/auth/', '/ar/auth/'] }],
    sitemap: new URL('/sitemap.xml', origin).toString(),
  }
}
