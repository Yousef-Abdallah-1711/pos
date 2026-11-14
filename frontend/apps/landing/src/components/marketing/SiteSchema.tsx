import { getLandingOrigin } from '@/lib/seo'

export function SiteSchema() {
  const schema = {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: 'SERVIO',
    url: getLandingOrigin().origin,
    inLanguage: ['en', 'ar'],
  }
  // Controlled metadata only; escape markup delimiters before script serialization.
  return <script type="application/ld+json">{JSON.stringify(schema).replace(/</g, '\\u003c')}</script>
}
