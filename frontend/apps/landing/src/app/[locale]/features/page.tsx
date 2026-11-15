import Link from 'next/link'
import { MarketingShell } from '@/components/marketing/MarketingShell'
import { ProductCapture } from '@/components/marketing/ProductCapture'
import { TrialLink } from '@/components/marketing/TrialLink'
import { FinalCTA } from '@/components/marketing/ClosingSections'
import { getFeatures } from '@/content/features'
import { getSiteCopy, type Locale } from '@/content/site'
import { marketingMetadata } from '@/lib/marketingMetadata'

type Props = { params: Promise<{ locale: Locale }> }

export async function generateMetadata({ params }: Props) {
  const { locale } = await params
  const copy = getSiteCopy(locale)
  return marketingMetadata(locale, '/features', copy.directoryTitle, copy.directoryBody)
}

function Workflow({ label, steps }: { label: string; steps: string[][] }) {
  return <div className="sv-feature-directory__workflow">
    {steps.map(([number, title, description]) => <article key={number}>
      <p className="sv-eyebrow">{number}</p><h3>{title}</h3><p>{description}</p>
    </article>)}
    <span className="sr-only">{label}</span>
  </div>
}

function Showcase({ locale, feature, index }: {
  locale: Locale
  feature: ReturnType<typeof getFeatures>[number]
  index: number
}) {
  return <article className={`sv-feature-showcase sv-feature-showcase--${feature.slug}`} id={feature.slug}>
    <div className="sv-feature-showcase__copy">
      <p className="sv-eyebrow">0{index + 1} / {feature.eyebrow}</p>
      <h3>{feature.headline}</h3><p>{feature.description}</p>
      <ol className="sv-feature-steps">{feature.steps.map((step, stepIndex) => <li key={step.title}>
        <span>0{stepIndex + 1}</span><strong>{step.title}</strong><p>{step.description}</p>
      </li>)}</ol>
      <p className="sv-note">{feature.boundary}</p>
      <Link className="sv-text-link" href={`/${locale}/features/${feature.slug}`}>
        {getSiteCopy(locale).details}<span aria-hidden="true">↗</span>
      </Link>
    </div>
    <div className="sv-feature-showcase__capture">
      <ProductCapture locale={locale} image={feature.image} alt={feature.title} />
    </div>
  </article>
}

export default async function FeaturesPage({ params }: Props) {
  const { locale } = await params
  const copy = getSiteCopy(locale)
  const features = getFeatures(locale)
  const hero = features[0]
  return <MarketingShell>
    <section className="sv-section sv-feature-directory">
      <div className="sv-shell">
        <div className="sv-feature-directory__hero">
          <div><p className="sv-eyebrow">{copy.directoryEyebrow}</p><h1>{copy.directoryTitle}</h1>
            <p className="sv-lead">{copy.directoryBody}</p>
            <TrialLink locale={locale} source="features-directory-hero" />
          </div>
          <ProductCapture locale={locale} image={hero.image} alt={hero.title} priority />
        </div>
        <section aria-labelledby="workflow-heading">
          <h2 id="workflow-heading">{copy.workflowTitle}</h2><p className="sv-lead">{copy.workflowBody}</p>
          <Workflow label={copy.workflowAccessibleLabel} steps={copy.journey} />
        </section>
        <section aria-labelledby="tasks-heading" className="sv-feature-directory__tasks">
          <p className="sv-eyebrow">{copy.tasksTitle}</p><h2 id="tasks-heading">{copy.tasksBody}</h2>
          <nav className="sv-feature-directory__task-nav" aria-label={copy.tasksTitle}>
            {features.map((feature, index) => <a key={feature.slug} href={`#${feature.slug}`}>
              <span>0{index + 1}</span>{feature.title}
            </a>)}
          </nav>
          <div className="sv-feature-gallery">
            {features.map((feature, index) => <Showcase key={feature.slug} locale={locale} feature={feature} index={index} />)}
          </div>
        </section>
        <section className="sv-section sv-feature-faq" aria-labelledby="directory-faq-heading">
          <p className="sv-eyebrow">{copy.faqTitle}</p><h2 id="directory-faq-heading">{copy.faqBody}</h2>
          {features.slice(0, 3).map(feature => <div key={feature.slug} className="sv-feature-faq__group">
            <h3>{feature.title}</h3>{feature.faqs.slice(0, 1).map(faq => <details key={faq.question}>
              <summary>{faq.question}<span aria-hidden="true">+</span></summary><p>{faq.answer}</p>
            </details>)}
          </div>)}
        </section>
      </div>
    </section>
    <FinalCTA locale={locale} />
  </MarketingShell>
}
