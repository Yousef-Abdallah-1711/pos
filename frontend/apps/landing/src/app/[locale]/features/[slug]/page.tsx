import Link from 'next/link'
import { ViewEvent } from '@/components/marketing/ViewEvent'
import { notFound } from 'next/navigation'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { ProductCapture } from '@/components/marketing/ProductCapture'
import { TrialLink } from '@/components/marketing/TrialLink'
import { FinalCTA } from '@/components/marketing/ClosingSections'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'
type Props = { params: Promise<{ locale: Locale; slug: string }> }
export function generateStaticParams() { return getFeatures('en').map(({ slug }) => ({ slug })) }
export async function generateMetadata({ params }: Props) {
 const { locale, slug } = await params; const feature = getFeatures(locale).find(f => f.slug === slug)
 if (!feature) notFound()
 return marketingMetadata(locale, `/features/${slug}`, feature.title, feature.description)
}
export default async function FeaturePage({ params }: Props) {
 const { locale, slug } = await params; const feature = getFeatures(locale).find(f => f.slug === slug)
 if (!feature) notFound()
 const c = getSiteCopy(locale)
 return <MarketingShell><ViewEvent locale={locale} source={slug} /><section className="sv-section sv-feature-hero"><div className="sv-shell">
 <Link className="sv-text-link" href={`/${locale}/features`}>{c.platform}</Link><p className="sv-eyebrow">{feature.eyebrow}</p>
 <h1>{feature.headline}</h1><p className="sv-lead">{feature.description}</p><TrialLink locale={locale} source={`feature-${slug}`} />
 <ProductCapture locale={locale} image={feature.image} alt={feature.title} />
 <div className="sv-feature-grid">{feature.steps.map(([title, body], i) => <article className="sv-feature-card" key={title}><p className="sv-eyebrow">0{i+1}</p><h2>{title}</h2><p>{body}</p></article>)}</div>
 <p className="sv-note">{feature.boundary}</p></div></section><FinalCTA locale={locale} /></MarketingShell>
}
