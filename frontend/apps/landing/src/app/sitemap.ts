import type { MetadataRoute } from 'next'
import { getLandingOrigin } from '@/lib/seo'
import { getFeatures } from '@/content/features'
export default function sitemap(): MetadataRoute.Sitemap {
 const origin = getLandingOrigin()
 const paths = ['', '/features', '/how-it-works', '/pricing', '/demo', ...getFeatures('en').map(f => `/features/${f.slug}`)]
 return paths.flatMap(path => (['en', 'ar'] as const).map(locale => ({ url: new URL(`/${locale}${path}`, origin).toString(), alternates: { languages: { en: new URL(`/en${path}`, origin).toString(), ar: new URL(`/ar${path}`, origin).toString() } } })))
}
