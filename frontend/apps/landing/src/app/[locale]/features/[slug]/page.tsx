import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ViewEvent } from '@/components/marketing/ViewEvent'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { ProductCapture } from '@/components/marketing/ProductCapture'
import { TrialLink } from '@/components/marketing/TrialLink'
import { FinalCTA } from '@/components/marketing/ClosingSections'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'

type Props = { params: Promise<{ locale: Locale; slug: string }> }

export function generateStaticParams() {
  return getFeatures('en').map(({ slug }) => ({ slug }))
}

function FeatureQuestions({ locale, faqs }: { locale: Locale; faqs: Array<{ question: string; answer: string }> }) {
  return <section className="sv-section sv-feature-faq" aria-labelledby="feature-faq-heading">
    <p className="sv-eyebrow">{getSiteCopy(locale).featureFaqTitle}</p>
    <div className="sv-faq">{faqs.map(item => <details key={item.question}>
      <summary>{item.question}<span aria-hidden="true">+</span></summary><p>{item.answer}</p>
    </details>)}</div>
  </section>
}

export async function generateMetadata({ params }: Props) {
  const { locale, slug } = await params
  const feature = getFeatures(locale).find(item => item.slug === slug)
  if (!feature) notFound()
  return marketingMetadata(locale, `/features/${slug}`, feature.title, feature.description)
}

export default async function FeaturePage({ params }: Props) {
  const { locale, slug } = await params
  const feature = getFeatures(locale).find(item => item.slug === slug)
  if (!feature) notFound()
  const copy = getSiteCopy(locale)
  const related = getFeatures(locale).filter(item => feature.related.includes(item.slug))
  return <MarketingShell locale={locale}><ViewEvent locale={locale} source={slug} />
    <section className={`sv-section sv-feature-hero sv-feature-page sv-feature-page--${feature.slug}`}>
      <div className="sv-shell">
        <Link className="sv-text-link" href={`/${locale}/features`}>{copy.platform}</Link>
        <p className="sv-eyebrow">{feature.eyebrow}</p><h1>{feature.headline}</h1>
        <p className="sv-lead">{feature.description}</p><TrialLink locale={locale} source={`feature-${slug}`} />
        <ProductCapture locale={locale} image={feature.image} alt={feature.title} priority />
        <section className="sv-feature-detail" aria-labelledby="feature-steps-heading">
          <p className="sv-eyebrow">{copy.featureStepsTitle}</p><h2 id="feature-steps-heading">{feature.title}</h2>
          <ol className="sv-feature-steps">{feature.steps.map((step, index) => <li key={step.title}>
            <span>0{index + 1}</span><h3>{step.title}</h3><p>{step.description}</p>
          </li>)}</ol>
          <p className="sv-note">{feature.boundary}</p>
        </section>
        <section className="sv-feature-related" aria-labelledby="feature-related-heading">
          <h2 id="feature-related-heading">{copy.featureRelatedTitle}</h2>
          <nav className="sv-related-links" aria-label={copy.featureRelatedTitle}>
            {related.map(item => <Link key={item.slug} href={`/${locale}/features/${item.slug}`}>{item.title}</Link>)}
          </nav>
        </section>
        <FeatureQuestions locale={locale} faqs={feature.faqs} />
      </div>
    </section>
    <FinalCTA locale={locale} />
  </MarketingShell>
}
