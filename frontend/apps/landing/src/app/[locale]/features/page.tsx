import Link from 'next/link'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { FinalCTA } from '@/components/marketing/ClosingSections'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale }> }
export async function generateMetadata({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return marketingMetadata(locale, '/features', c.explorerTitle, c.explorerBody)
}
export default async function FeaturesPage({ params }: Props) {
 const { locale } = await params; const c = getSiteCopy(locale)
 return <MarketingShell><section className="sv-section sv-feature-hero"><div className="sv-shell">
 <p className="sv-eyebrow">{c.platform}</p><h1>{c.explorerTitle}</h1><p className="sv-lead">{c.explorerBody}</p>
 <div className="sv-feature-grid">{getFeatures(locale).map(feature => <article className="sv-feature-card" key={feature.slug}>
 <h2>{feature.title}</h2><p>{feature.description}</p><Link className="sv-text-link" href={`/${locale}/features/${feature.slug}`}>{c.details}<span aria-hidden="true">↗</span></Link>
 </article>)}</div></div></section><FinalCTA locale={locale} /></MarketingShell>
}
