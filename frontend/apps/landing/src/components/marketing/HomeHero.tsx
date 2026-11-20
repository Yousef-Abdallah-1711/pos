import Link from 'next/link'
import { ArrowDown, ArrowRight } from 'lucide-react'
import { getSiteCopy, type Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'
import { TrialLink } from './TrialLink'

export function HomeHero({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <section className="sv-hero">
    <div className="sv-shell sv-hero-grid">
      <div className="sv-hero-copy">
        <p className="sv-eyebrow"><span className="sv-status-dot" />{c.eyebrow}</p>
        <h1>{c.heroStart}<br /><em>{c.heroEnd}</em></h1>
        <p className="sv-lead">{c.heroBody}</p>
        <div className="sv-actions"><TrialLink locale={locale} source="hero" />
          <a className="sv-text-link" href="#order-journey">{c.heroSecondary}<ArrowDown size={17} aria-hidden="true" /></a>
        </div>
      </div>
      <div className="sv-hero-stage" aria-label={c.workflowAccessibleLabel}>
        <div className="sv-stage-rail"><span>01 / {c.counter}</span><span>02 / {c.kitchen}</span></div>
        <p className="sv-stage-source">{c.captureSource}</p>
        <ProductCapture locale={locale} image="pos" alt={c.posCaptureAlt} priority className="sv-hero-pos" />
        <div className="sv-handoff"><span />{c.handoff}</div>
        <ProductCapture locale={locale} image="kitchen" alt={c.kitchenCaptureAlt} priority className="sv-hero-kitchen" />
      </div>
    </div>
    <div className="sv-shell sv-hero-foot"><p>{c.heroNote}</p>
      <Link href={`/${locale}/how-it-works`}>{c.how}<ArrowRight size={17} className="sv-directional" aria-hidden="true" /></Link>
    </div>
  </section>
}
